FROM golang:latest

RUN go get github.com/zenwerk/go-wave && go get github.com/gorilla/websocket && go get github.com/spf13/viper

COPY . /opt/go_record_demo

WORKDIR /opt/go_record_demo
CMD go build -o /tmp/___go_build_recv_wav_server_go recv_wav_server.go