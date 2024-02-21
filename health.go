package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
)

const socketPath = "/tmp/healthcheck.socket"
const BUF_SIZE = 4096

type HealthcheckServer struct {
	socketPath      string
	socket          net.Listener
	ctx             context.Context
	cancel          context.CancelFunc
	healthcheckFunc func() error
}

// randomString generates a random string of the specified length using
// characters from the alphabet.
// The function uses the `letters` slice which contains lowercase and uppercase
// letters, as well as digits.
// It creates an empty slice of runes with the desired length and randomly
// selects a character from `letters` for each position.
// Finally, it converts the rune slice to a string and returns the result.
//
// Example:
//
//	s.socketPath = fmt.Sprintf("/tmp/%s.socket", randomString(12))
func randomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, length)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// Init initializes the HealthcheckServer instance by setting the socketPath
// and creating a new context with cancel function.
//
// Example:
//
//	srv := HealthcheckServer{}
//	srv.Init()
func (s *HealthcheckServer) Init(f func() error) {
	s.socketPath = socketPath
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.healthcheckFunc = f
}

func (s *HealthcheckServer) Start() (err error) {
	if _, err := os.Stat(s.socketPath); err == nil {
		err = os.Remove(s.socketPath)
		if err != nil {
			return err
		}
	}
	s.socket, err = net.Listen("unix", s.socketPath)
	if err != nil {
		return err
	}
	for {
		select {
		case <-s.ctx.Done():
			_ = s.socket.Close()
			return nil
		default:
			conn, err := s.socket.Accept()
			if err != nil {
				continue
			}
			go func(conn net.Conn) {
				defer conn.Close()
				inputBuffer := make([]byte, BUF_SIZE)
				n, _ := conn.Read(inputBuffer)
				incomingMessage := strings.TrimSpace(string(inputBuffer[:n-1]))
				if incomingMessage != "ping" {
					_, _ = conn.Write([]byte("send 'ping' to trigger healthcheck\n"))
					return
				}
				err = s.healthcheckFunc()
				if err != nil {
					conn.Write([]byte(fmt.Sprintf("%s", err.Error())))
					return
				}
				conn.Write([]byte("success\n"))
			}(conn)
		}
	}
}

func (s *HealthcheckServer) Stop() {
	s.cancel()
}
