package socket

import (
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	ln, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		t.Fatalf("net.Listen failed: %s", err)
	}
	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}

			t.Cleanup(func() { conn.Close() })
		}
	}()

	// create a new client with the address of the temporary server
	c, err := NewClient(ln.Addr().String())
	if err != nil {
		t.Fatalf("NewClient failed: %s", err)
	}

	// check that conn, readWriter, sizeCh, matchCh and closeCh have been set up correctly
	if c.conn == nil {
		t.Errorf("c.conn was nil")
	}
	if c.readWriter == nil {
		t.Errorf("c.readWriter was nil")
	}
	if c.sizeCh == nil {
		t.Errorf("c.sizeCh was nil")
	}
	if c.matchCh == nil {
		t.Errorf("c.matchCh was nil")
	}
	if c.closeCh == nil {
		t.Errorf("c.closeCh was nil")
	}
}

func TestNewClientFail(t *testing.T) {
	t.Parallel()

	if _, err := NewClient("random bad ip"); err == nil {
		t.Fatalf("expected error, got nothing")
	}
}

func TestDisconnectServer(t *testing.T) {
	t.Parallel()
	// start the server in server-client/server/main.go
	client, err := NewClient(":8080")
	if err != nil {
		t.Fatalf("NewClient failed: %s", err)
	}
	// stop the server in server-client/server/main.go
	time.Sleep(7 * time.Second)
	t.Log(client.readWriter)
	_, _, err = client.readWriter.Read()
	if err != nil {
		if isConnectionResetError(err) || strings.ContainsAny(err.Error(), "closed") ||
			err == io.EOF || strings.ContainsAny(err.Error(), "reset") || strings.ContainsAny(err.Error(), "EOF") {
			t.Log("server disconnected.")
		}
	}
}
