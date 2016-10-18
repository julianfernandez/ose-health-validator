echo "Building health-validator docker image"

TAG=$1
PROXY=""

if [ "$#" -gt 1 ]; then
    PROXY="--build-arg https_proxy=$2"
fi

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o health-validator ./go/*.go 

if [ $? = 0 ]
then
	docker build ${PROXY} -t juliusmore/ose-health-validator:${TAG} -f ./image/Dockerfile .
fi