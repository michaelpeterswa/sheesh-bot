# -=-=-=-=-=-=- Compile Image -=-=-=-=-=-=-

FROM golang:1.16 AS stage-compile

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./cmd/sheesh-bot
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/sheesh-bot

# -=-=-=-=-=-=- Final Image -=-=-=-=-=-=-

FROM alpine:latest 

WORKDIR /root/
COPY --from=stage-compile /go/src/app/sheesh-bot ./

RUN apk --no-cache add ca-certificates

EXPOSE 8080

ENTRYPOINT [ "./sheesh-bot" ]  