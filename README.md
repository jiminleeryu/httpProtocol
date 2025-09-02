## Overview
HTTP/1.1 Protocol: custom HTTP server from scratch

Pull data from file push to internet connection

## Running 
`go run . | tee /tmp/tcp.txt` <- using tee
`go run ./tcplistener/main.go` <- run TCP listener

To connect, run:
`go run ./tcplistener/main.go`

`nc -v localhost $<PORT>`

## TCP vs UDP
- TCP (packet ordering) sends entire JSON in order and complete (sliding window)
- UDP you determine how to break up organize and send data, and how receiver organizes - a lot more protection (more performant): not waiting for an ack 99% of time
- TCP over UDP we can Nack packets, and don't need to wait to send packets back, we can ask to resend lost packets

## HTTP/1.1
- Specify what you're sending, and host
- Requests have headers, same headers

## HTTP Message
    GET (type) /cats (destination) HTTP/1.1\r\n
    Host: url \r\n
    User-Agent: \r\n
    Accept: \r\n
    Content-Length: (how big body is) \r\n
    \r\n
    {
        "body": "body"
    }

