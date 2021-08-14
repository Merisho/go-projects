# Just a chat with file transfer feature
This is merely a chat implemented over TCP with capability of transferring files

Files are transferred asynchronously, so while the one is being sent, you can still send and receive messages.
A TCP connection is multiplexed with the help of SPPP (see below).

Clients automatically receive files that are sent.
All the files are downloaded in the directory from which the client executable is started.

## Server
```
go run cmd/server/server.go
```

Server starts TCP listener on `1337` port.

## Client
```
go run cmd/client/client.go [host] [port]
```
By default, the client connects to `localhost:1337`.

## SPPP - Simple Pet Project Protocol
This is the thing that does connection multiplexing.
It is stream based, so whenever you write anything to SPPP connection, internally happens the following:

1. a write stream is created with its own unique ID
2. the data is split into chunks which are sent one by one
3. once the sending is done, the write stream is closed and `end` chunk is sent to receiver

On the recipient side:

1. once there is an incoming stream, the read stream is created
2. the read stream is read till the `end` chunk or till the read timeout is reached.
The timeout is reached when read stream has not received any data in specified amount of time.

SPPP is developed to support 2 types of data: bare streams and atomic messages (which is just an abstraction over streams).
The latter has its own timeout for entire message read (timeout for receiving `end` chunk), so even if read stream receives data chunks in time, the entire message read may still be not completed in time.
