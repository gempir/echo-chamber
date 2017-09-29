FROM golang:latest
WORKDIR /go/src/github.com/gempir/echo-chamber
RUN go get github.com/gempir/go-twitch-irc && go get gopkg.in/olivere/elastic.v5
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/gempir/echo-chamber/app .
CMD ["./app"]  
EXPOSE 3333