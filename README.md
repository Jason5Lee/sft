# sft

<del>Stupid</del>Simple File Transporter

## How to use

### Server

`sft listen [ip]:<port>`

Example: `sft listen 127.0.0.1:8888`, `sft listen :8888`

### Client

`sft connect <ip>:<port>`

Example: `sft connect 127.0.0.1:8888`

If the connection is established, you can use command `ls` to list the files and command `get` to download a file.

## Transportation Protocol

This program uses TCP/IP for the transportation protocol. We use the first 4 bytes as header. The header is a little-eden representation of a 32-bit signed integer. If this integer is positive, it means that the operation success and the integer equals to the length of respond message, if it is negative, it means that the operation failed and the integer equals to the opposite number of the length of the error message.  

When the client executes `ls` command, it will send an empty package (that is, with only the header of value `0`). When the client executes `get` command, it will send the file name to server.

When the server receives `ls` command, it sends the file-list string back. When the server receives `get` command, it sends the content of the file back. Whenever error occurs (i.e. File not found), it sends error message back.