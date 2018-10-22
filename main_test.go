package main

import (
	"bytes"
	"testing"
)

func TestSendReceiveHeader(t *testing.T) {
	var buf bytes.Buffer

	const TestHeader = 0x00241311
	sendHeader(&buf, TestHeader)
	b := buf.Bytes()
	if len(b) != 4 || b[0] != 0x11 || b[1] != 0x13 || b[2] != 0x24 || b[3] != 0x00 {
		t.Fail()
	}

	recvHeader, _ := receiveHeader(&buf)
	if recvHeader != TestHeader {
		t.Fail()
	}
}

func TestSendReceiveString(t *testing.T) {
	var buf bytes.Buffer

	sendString(&buf, "314159")
	b := buf.Bytes()

	if len(b) != 6+4 {
		t.Fail()
	}
	if b[0] != 6 || b[1] != 0 || b[2] != 0 || b[3] != 0 {
		t.Fail()
	}
	if b[4] != '3' || b[5] != '1' || b[6] != '4' || b[7] != '1' || b[8] != '5' || b[9] != '9' {
		t.Fail()
	}

	recvString, _ := receiveString(&buf)
	if recvString != "314159" {
		t.Fail()
	}
}
