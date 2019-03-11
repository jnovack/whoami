package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	// "github.com/pkg/profile"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	port      string
	version   string
	commit    string
	buildDate string
	buildTime string
)

func init() {
	flag.StringVar(&port, "port", "80", "give me a port number")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	fmt.Printf("jnovack/whoami %s\n", version)
	fmt.Printf(" :: commit %s built on %s at %s ::\n", commit, buildDate, buildTime)

	// defer profile.Start().Stop()
	flag.Parse()
	http.HandleFunc("/", whoami)
	http.HandleFunc("/api", api)
	http.HandleFunc("/echo", echoHandler)
	http.HandleFunc("/health", healthHandler)
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println(" - listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func printBinary(s []byte) {
	fmt.Printf("Received b:")
	for n := 0; n < len(s); n++ {
		fmt.Printf("%d,", s[n])
	}
	fmt.Printf("\n")
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		printBinary(p)
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			return
		}
	}
}

func whoami(w http.ResponseWriter, req *http.Request) {
	u, _ := url.Parse(req.URL.String())
	queryParams := u.Query()

	wait := queryParams.Get("wait")
	if len(wait) > 0 {
		duration, err := time.ParseDuration(wait)
		if err == nil {
			time.Sleep(duration)
		}
	}

	hostname, _ := os.Hostname()
	fmt.Fprintln(w, "Hostname:", hostname)

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// No need to EVER print localhost, we know
			if ip.String() != "127.0.0.1" {
				fmt.Fprintln(w, "IP:", ip)
			}
		}
	}

	environ := os.Environ()
	for _, env := range environ {
		fmt.Fprintln(w, "ENV:", env)
	}

	fmt.Fprintln(w, "VERSION:", version)
	fmt.Fprintln(w, "COMMIT:", commit)
	fmt.Fprintln(w, "BUILD_DATE:", buildDate)
	fmt.Fprintln(w, "BUILD_TIME:", buildTime)

	fmt.Fprintln(w, "TIMESTAMP:", time.Now().Format(time.RFC3339Nano))

	req.Write(w)
}

func api(w http.ResponseWriter, req *http.Request) {
	hostname, _ := os.Hostname()
	type Build struct {
		Version   string `json:"version,omitempty"`
		Commit    string `json:"commit,omitempty"`
		BuildDate string `json:"build_date,omitempty"`
		BuildTime string `json:"build_time,omitempty"`
	}
	data := struct {
		Hostname    string      `json:"hostname,omitempty"`
		IP          []string    `json:"ip,omitempty"`
		Environment []string    `json:"environment,omitempty"`
		Headers     http.Header `json:"headers,omitempty"`
		URL         string      `json:"url,omitempty"`
		Method      string      `json:"method,omitempty"`
		Build       Build       `json:"build,omitempty"`
	}{
		hostname,
		[]string{},
		[]string{},
		req.Header,
		req.URL.RequestURI(),
		req.Method,
		Build{
			Version:   version,
			Commit:    commit,
			BuildDate: buildDate,
			BuildTime: buildTime,
		},
	}

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// No need to EVER print localhost, we know
			if ip.String() != "127.0.0.1" {
				data.IP = append(data.IP, ip.String())
			}
		}
	}

	environ := os.Environ()
	for _, env := range environ {
		data.Environment = append(data.Environment, env)
	}

	json.NewEncoder(w).Encode(data)
}

type healthState struct {
	StatusCode int
}

var currentHealthState = healthState{204}
var mutexHealthState = &sync.RWMutex{}

func healthHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		var statusCode int
		err := json.NewDecoder(req.Body).Decode(&statusCode)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		} else {
			fmt.Printf("Update health check status code [%d]\n", statusCode)
			mutexHealthState.Lock()
			defer mutexHealthState.Unlock()
			currentHealthState.StatusCode = statusCode
		}
	} else {
		mutexHealthState.RLock()
		defer mutexHealthState.RUnlock()
		w.WriteHeader(currentHealthState.StatusCode)
	}
}
