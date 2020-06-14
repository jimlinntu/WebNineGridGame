FROM ubuntu:18.04
RUN apt-get update
RUN apt-get install -y software-properties-common
RUN apt-get install -y git-core
RUN apt-get install -y curl
RUN add-apt-repository -y ppa:longsleep/golang-backports
RUN apt-get update
RUN apt install -y golang-1.11-go
ENV PATH="/usr/lib/go-1.11/bin:${PATH}"
WORKDIR /root/WebNineGridGame
COPY ./main.go ./
RUN go get -v -d ./...
