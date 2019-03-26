FROM golang

RUN go get github.com/gorilla/mux

WORKDIR /go/src/github.com/heaptracetechnology/microservice-ffmpeg

ADD . /go/src/github.com/heaptracetechnology/microservice-ffmpeg

RUN go install github.com/heaptracetechnology/microservice-ffmpeg

ENTRYPOINT microservice-ffmpeg

EXPOSE 3000