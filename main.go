package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

// LocalError represents the error that doesn't involve internet.
type LocalError struct {
	err error
}

func (ce LocalError) Error() string {
	return ce.err.Error()
}

func printUsageAndExit() {
	fmt.Println(`Usage:
	sft [command] [arguments]

Avaliable commands and their arguments:

	listen [ip]:<port>	starts the sft server at current directory.
	connect	<ip>:<port>	connects to a sft server.`)
	os.Exit(1)
}

func ifErrorPrintAndExit(err error) {
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Print(`Available Command:
	ls		list all filenames
	get <filename>	download a file to current directory
	exit		disconnect
	help		print this message
`)
}

func sendHeader(conn net.Conn, header int32) error {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(header))
	_, err := conn.Write(buf)
	return err
}

func sendString(conn net.Conn, ok bool, str string) (err error) {
	bytes := []byte(str)
	if ok {
		err = sendHeader(conn, int32(len(bytes)))
	} else {
		err = sendHeader(conn, int32(-len(bytes)))
	}

	if err != nil {
		return
	}
	_, err = conn.Write([]byte(str))
	return
}

func sendFile(conn net.Conn, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return LocalError{err}
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return LocalError{err}
	}

	msgSize := stat.Size()
	if msgSize > 2147483647 {
		return LocalError{fmt.Errorf("file too large")}
	}

	err = sendHeader(conn, int32(msgSize))
	if err != nil {
		return err
	}

	_, err = io.Copy(conn, f)
	return err
}

func receiveHeader(conn net.Conn) (header int32, err error) {
	buf := make([]byte, 4)

	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return
	}

	header = int32(binary.LittleEndian.Uint32(buf))
	return
}

func receiveString(conn net.Conn) (msg string, err error) {
	header, err := receiveHeader(conn)
	if err != nil {
		return
	}

	var buf []byte
	if header > 0 {
		buf = make([]byte, header)
	} else {
		buf = make([]byte, -header)
	}

	_, err = io.ReadFull(conn, buf)
	if err == nil {
		msg = string(buf)
	}

	return
}

func receiveFile(conn net.Conn, file string) error {
	header, err := receiveHeader(conn)
	if err != nil {
		return err
	}

	if header < 0 {
		buf := make([]byte, -header)
		_, err = io.ReadFull(conn, buf)
		if err != nil {
			return err
		}
		return LocalError{fmt.Errorf(string(buf))}
	}

	f, err := os.Create(file)
	if err != nil {
		return LocalError{err}
	}
	defer f.Close()

	_, err = io.CopyN(f, conn, int64(header))
	return err
}

func clientLoop(stdreader *bufio.Reader, conn net.Conn) (err error) {
	fmt.Print(">")

	cmd, err := stdreader.ReadString('\n')
	if err != nil {
		return
	}

	tokens := strings.Fields(cmd)

	if len(tokens) == 0 {
		return
	}

	switch tokens[0] {
	case "ls":
		if len(tokens) != 1 {
			fmt.Println("Usage: ls")
			return
		}

		err = sendHeader(conn, 0)
		if err != nil {
			return
		}

		msg, err := receiveString(conn)
		if err != nil {
			return err
		}

		fmt.Println(msg)
		return err
	case "get":
		if len(tokens) != 2 {
			fmt.Println("Usage: get <file name>")
			return nil
		}

		err = sendString(conn, true, tokens[1])
		if err != nil {
			return
		}

		err = receiveFile(conn, tokens[1])
		if err != nil {
			return err
		}
		return nil
	case "exit":
		return fmt.Errorf("Goodbye")
	case "help":
		showHelp()
		return nil
	default:
		fmt.Println("Unknow command:", cmd)
		return nil
	}
}

func startClient() {
	fmt.Println("Connecting...")
	conn, err := net.Dial("tcp", os.Args[2])
	ifErrorPrintAndExit(err)
	defer conn.Close()

	stdreader := bufio.NewReader(os.Stdin)
	showHelp()

	for {
		if err = clientLoop(stdreader, conn); err != nil {
			fmt.Println(err.Error())
			if _, ok := err.(LocalError); !ok {
				break
			}
		}
	}
}

func serverLoop(conn net.Conn) error {
	msg, err := receiveString(conn)

	if err != nil {
		return err
	}

	if len(msg) == 0 { // ls

		files, err := ioutil.ReadDir("./")
		if err != nil {
			return sendString(conn, false, err.Error())
		}

		var code bytes.Buffer
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			code.WriteString(f.Name())
			code.WriteString("\n")
		}
		return sendString(conn, true, code.String())
	}

	// get
	return sendFile(conn, msg)
}

func startServer() {
	fmt.Println("Listening " + os.Args[2])
	sock, err := net.Listen("tcp", os.Args[2])
	ifErrorPrintAndExit(err)

	defer sock.Close()

	for {
		conn, err := sock.Accept()
		if err != nil {
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()

			for {
				if err := serverLoop(conn); err != nil {
					fmt.Println(err.Error())

					if _, ok := err.(LocalError); !ok {
						break
					}
				}
			}
		}(conn)
	}
}

func main() {
	if len(os.Args) == 3 {
		switch os.Args[1] {
		case "connect":
			startClient()
			break
		case "listen":
			startServer()
			break
		default:
			printUsageAndExit()
		}
	} else {
		printUsageAndExit()
	}
}
