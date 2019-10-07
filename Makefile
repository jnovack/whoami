
build:
	docker build \
		--build-arg VERSION=`git describe --tags --always` \
		--build-arg COMMIT=`git rev-parse --short HEAD` \
		--build-arg BUILD_RFC3339=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		-t jnovack/whoami .

run:
	docker run -p 80 jnovack/whoami
