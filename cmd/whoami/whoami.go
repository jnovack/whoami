package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	release "github.com/jnovack/release"

	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	port string
)

func init() {
	flag.StringVar(&port, "port", "80", "port number")
}

func main() {
	fmt.Println(release.Info())

	flag.Parse()
	http.HandleFunc("/", whoami)
	http.HandleFunc("/api", api)
	http.HandleFunc("/healthcheck", healthHandler)
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

	req.Header.Add("Cache-Control", "must-validate")
	req.Header.Add("Hostname", hostname)

	// fmt.Fprintln(w, "Hostname:", hostname)

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
			if ip.String() != "127.0.0.1" && ip.String() != "::1" && ip.String() != "fe80::1" {
				fmt.Fprintln(w, "IP:", ip)
			}
		}
	}

	environ := os.Environ()
	for _, env := range environ {
		fmt.Fprintln(w, "ENV:", env)
	}

	fmt.Fprintln(w, "BUILD_VERSION:", release.Version)
	fmt.Fprintln(w, "BUILD_COMMIT:", release.Revision)
	fmt.Fprintln(w, "BUILD_RFC3339:", release.BuildRFC3339)

	fmt.Fprintln(w, "TIMESTAMP:", time.Now().Format(time.RFC3339Nano))

	fmt.Fprintln(w, "PROTOCOL:", req.Proto)

	req.Write(w)
}

func api(w http.ResponseWriter, req *http.Request) {
	hostname, _ := os.Hostname()
	type Build struct {
		Application  string `json:"application,omitempty"`
		GoVersion    string `json:"go_version,omitempty"`
		Version      string `json:"version,omitempty"`
		Commit       string `json:"commit,omitempty"`
		BuildRFC3339 string `json:"build_rfc3339,omitempty"`
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
			Application:  release.Application,
			GoVersion:    release.GoVersion,
			Version:      release.Version,
			Commit:       release.Revision,
			BuildRFC3339: release.BuildRFC3339,
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
