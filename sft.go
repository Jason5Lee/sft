package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

const packageSize int = 1024

func isBlank(ch byte) bool {
	return ch == ' ' || ch == '\n' || ch == '\r'
}
func readCmd(from *bufio.Reader, to chan string) {
	text, _ := from.ReadString('\n')

	p := 0

	for {
		for p < len(text) && isBlank(text[p]) {
			p++
		}

		if p == len(text) {
			break
		}

		q := p + 1
		for q < len(text) && !isBlank(text[q]) {
			q++
		}
		to <- text[p:q]

		p = q + 1
	}
	to <- ""
}

func printUsageAndExit() {
	fmt.Println("Usage:\n\tsft server <listening ip>:<listening port>\n\tsft client <server ip>:<port>")
	os.Exit(1)
}

func ifErrorPrintAndExit(err error) {
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Print("Available Command:\n" +
		"\tls\t\tList all file.\n" +
		"\tget <file name>\tGet a file.\n" +
		"\texit\t\tExit.\n" +
		"\thelp\t\tPrint this message.\n")
}
func startClient() {
	fmt.Println("Connecting...")
	conn, err := net.Dial("tcp", os.Args[2])
	ifErrorPrintAndExit(err)

	defer conn.Close()

	stdreader := bufio.NewReader(os.Stdin)

	showHelp()
	for {
		fmt.Print(">")
		var cmds = make(chan string)
		go readCmd(stdreader, cmds)
		cmd := <-cmds
		if cmd == "" {
			continue
		}

		var code bytes.Buffer
		var fileName string

		switch cmd {
		case "ls":
			if <-cmds != "" {
				fmt.Println("Usage: ls")
				continue
			}

			code.WriteByte(0)
			break
		case "get":
			fileName = <-cmds
			if fileName == "" || <-cmds != "" {
				fmt.Println("Usage: get <file name>")
				continue
			}

			code.WriteByte(1)
			code.WriteString(fileName)

		case "exit":
			goto exit
		case "help":
			showHelp()
			continue
		default:
			fmt.Println("Unknow command:", cmd)
			continue
		}

		_, err = code.WriteTo(conn)

		if err != nil {
			fmt.Println(err)
			continue
		}
		var respond = make([]byte, packageSize)
		n, err := conn.Read(respond)

		if err != nil {
			fmt.Println(err)
			continue
		}

		respond = respond[:n]

		if respond[0] != 0 {
			fmt.Println(string(respond[1:]))
			continue
		}
		switch cmd {
		case "ls":
			fmt.Println(string(respond[1:]))
			break
		case "get":
			file, err := os.Create(fileName)

			if err != nil {
				fmt.Println(err)
				continue
			}

			file.Write(respond[1:])
			for n == packageSize {
				n, err = conn.Read(respond)
				if err != nil {
					fmt.Println(err)
				}
				file.Write(respond[1:n])
			}
			file.Close()
		}
	}
exit:
	fmt.Println("Goodbye!")
}

func startServer() {
	sock, err := net.Listen("tcp", os.Args[2])
	ifErrorPrintAndExit(err)

	defer sock.Close()

	for {
		conn, err := sock.Accept()
		if err != nil {
			continue
		}

		go (func(conn net.Conn) {
			defer conn.Close()

			for {
				var recieve = make([]byte, packageSize)
				n, err := conn.Read(recieve)

				if err != nil {
					continue
				}

				recieve = recieve[:n]

				var errMsg string
				switch recieve[0] {
				case 0: // ls
					var code bytes.Buffer
					files, err := ioutil.ReadDir("./")
					if err != nil {
						errMsg = err.Error()
						break
					}
					code.WriteByte(0)
					for _, f := range files {
						if f.IsDir() {
							continue
						}
						code.WriteString(f.Name())
						code.WriteString("\n")
					}
					code.WriteTo(conn)
					break

				case 1: // get
					file, err := os.Open(string(recieve[1:]))
					if err != nil {
						errMsg = err.Error()
						break
					}
					var buffer = make([]byte, packageSize)
					buffer[0] = 0
					n, _ := file.Read(buffer[1:])
					n++
					conn.Write(buffer[:n])
					for n == packageSize {
						n, _ = file.Read(buffer[1:])
						n++
						conn.Write(buffer[:n])
					}
					file.Close()
				}
				if errMsg != "" {
					var code bytes.Buffer
					code.WriteByte(1)
					code.WriteString(errMsg)
					code.WriteTo(conn)
				}
			}
		})(conn)
	}
}
func main() {
	if len(os.Args) != 3 {
		printUsageAndExit()
	}

	switch os.Args[1] {
	case "client":
		startClient()
		break
	case "server":
		startServer()
		break
	default:
		printUsageAndExit()
	}
}
