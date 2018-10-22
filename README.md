[![wercker status](https://app.wercker.com/status/ddd250bc7c2e2c24bfe85689f8ff45ba/s/master "wercker status")](https://app.wercker.com/project/byKey/ddd250bc7c2e2c24bfe85689f8ff45ba)

# sft

<del>Stupid</del>Simple File Transporter

## Warning

This project is originally my Internet course's homework. Currently it is the project for me to practice CI/CD development. It doesn't have any security or anti-DDoS. DO NOT USE IT IN THE PRODUCTION.

## How to use

### Server

`sft listen [ip]:<port>`

Example: `sft listen 127.0.0.1:8888`, `sft listen :8888`

### Client

`sft connect <ip>:<port>`

Example: `sft connect 127.0.0.1:8888`

If the connection is established, you can use command `ls` to list the files and command `get` to download a file.

## Transportation Protocol

This program uses TCP/IP for the transportation protocol. We send message via one or more TCP/IP packages. The first package, header, is a little-eden representation of a 32-bit unsigned integer. This integer represents the length of the message. The rest of the packages are the chunked message. 

When the client executes `ls` command, it will send an empty message (that is, with only the header of value `0`). When the client executes `get` command, it will send the file name.

When the server receives `ls` command, it sends the file-list string or the error message back. When the server receives `get` command, it sends one or two messages. The first message is the error message. If there's no error, the first message is empty, and the second message is sent, which is the content of the file.

## TODO

* Add docker back.