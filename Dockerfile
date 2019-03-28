FROM golang

RUN go get github.com/gorilla/mux

# install FFmpeg

RUN apt-get update && apt-get install -y yasm pkg-config && apt-get install -y zlib1g-dev && apt-get install -y libpng-dev


RUN git clone https://github.com/FFmpeg/FFmpeg.git

RUN cd FFmpeg && ./configure && make && make install 

# this is the directory where the terminal will start
WORKDIR /go

CMD ["/bin/bash"]

RUN go get github.com/3d0c/gmf

WORKDIR /go/src/github.com/heaptracetechnology/microservice-ffmpeg

ADD . /go/src/github.com/heaptracetechnology/microservice-ffmpeg

RUN go install github.com/heaptracetechnology/microservice-ffmpeg

RUN chmod 755 /tmp

ENTRYPOINT microservice-ffmpeg

EXPOSE 3000

