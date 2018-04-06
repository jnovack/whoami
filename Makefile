
build:
	docker build -t jnovack/whoami .

run:
	docker run -p 80 jnovack/whoami