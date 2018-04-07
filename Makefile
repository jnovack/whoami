
build:
	docker build \
		--build-arg VERSION=`cat .version` \
		--build-arg COMMIT=`git rev-parse --short HEAD` \
		--build-arg BUILD_DATE=`date +%F` \
		--build-arg BUILD_TIME=`date +%T%z` \
		-t jnovack/whoami .

run:
	docker run -p 80 jnovack/whoami