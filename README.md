# sft

<del>Stupid</del>Simple File Transporter

## How to use

### Server

`sft server <listening ip>:<listening port>`

Example: `sft server 127.0.0.1:8888`

### Client

`sft client <server ip>:<port>`

Example: `sft client 127.0.0.1:8888`

If the connection is established, you can use command `ls` to list the files and command `get` to download a file.

## Transportation Protocal

This program uses TCP/IP for the transportation protocal. It sets the maximum size of one package to 1024, which is specified by constance `packageSize`. 

When the client executes `ls` command, it sends a package with a single byte `0` to server, the server sends a package that starts with a byte `0` and follows by the file-list string back. It assumes that the message does not exceed the maximum size.

When the client executes `get` command, it sends a package that starts with a byte `1` and follows by the file-name string to server, the server sends one or more packages that starts with a byte `0` and follows by the file content back. The server guarantees that the size of the last package is strictly less than the maximum size, even though it may contains only a single byte `0`, which indicates that there's nothing more.

If there's any error on server, it sends a package starts with a byte `1` and follows by the error-message string back.