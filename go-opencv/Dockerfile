FROM golang:latest
RUN apt-get update -qq && apt-get install -y  build-essential libopencv-dev
WORKDIR /go/src/github.com/lazywei/
RUN mkdir -p /go/src/github.com/lazywei/go-opencv
RUN cd /go/src/github.com/lazywei/
RUN git clone https://github.com/go-opencv/go-opencv.git
RUN cd go-opencv/samples
RUN ls
CMD go run /go/src/github.com/lazywei/go-opencv/samples/hellocv.go