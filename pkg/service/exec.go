package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/whiteblock/genesis-cli/pkg/auth"

	"github.com/docker/cli/cli/streams"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/pkg/term"
	log "github.com/sirupsen/logrus"
	"github.com/whiteblock/utility/common"
)

type connWrapper struct {
	net.Conn
}

func (cw connWrapper) Close() error {
	log.Info("blocking connection close")
	return nil
}

type clientHijacker struct {
	conn net.Conn
}

func (ch *clientHijacker) Conn() net.Conn {
	return ch.conn
}

var onCloseHandlers []func()

func (ch *clientHijacker) ConnectTLS(network, addr string) (net.Conn, error) {
	log.WithFields(log.Fields{"network": network, "addr": addr}).Debug("opening a tls connection")
	conn, err := tls.Dial(network, addr, &tls.Config{
		ServerName: strings.Split(conf.APIHost(), ":")[0],
	})
	if err == nil {
		ch.conn = connWrapper{Conn: conn}
	}
	return ch.conn, err
}

func (ch *clientHijacker) Connect(network, addr string) (net.Conn, error) {
	log.WithFields(log.Fields{"network": network, "addr": addr}).Debug("opening a tls connection")
	conn, err := net.Dial(network, addr)
	if err == nil {
		ch.conn = connWrapper{Conn: conn}
	}
	return ch.conn, err
}

func ListContainers(testID string) (out []string, err error) {
	return out, auth.Get(conf.ListContainersURL(testID), &out)
}

func PrepareExec(cmd common.Exec) (common.ExecInfo, error) {
	data, err := json.Marshal(cmd)
	if err != nil {
		return common.ExecInfo{}, err
	}
	data, err = auth.Post(conf.PrepareExecURL(cmd.Test), data)
	if err != nil {
		return common.ExecInfo{}, err
	}
	var out common.ExecInfo
	return out, json.Unmarshal(data, &out)
}

func RunDetach(cmd common.ExecAttach) error {
	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	_, err = auth.Post(conf.RunDetachURL(cmd.Test), data)
	return err
}
func handleIn(conn net.Conn, tty bool) <-chan error {
	out := make(chan error)
	go func() {
		if tty {
			input := streams.NewIn(os.Stdin)
			input.SetRawTerminal()
			onCloseHandlers = append(onCloseHandlers, func() {
				input.RestoreTerminal()
			})
			inputStream := ioutils.NewReadCloserWrapper(
				term.NewEscapeProxy(input /*ctrl-p ctrl-q*/, []byte{16, 17}), input.Close)
			_, err := io.Copy(conn, inputStream)
			out <- err
		} else {
			_, err := io.Copy(conn, os.Stdin)
			out <- err
		}

		log.Info("connection to stdin closed")

	}()
	return out
}

func handleOut(conn net.Conn, tty bool) <-chan error {
	out := make(chan error)
	go func() {
		if tty {
			outstream := streams.NewOut(os.Stdout)
			outstream.SetRawTerminal()
			onCloseHandlers = append(onCloseHandlers, func() {
				outstream.RestoreTerminal()
			})
			_, err := io.Copy(outstream, conn)
			out <- err

		} else {
			_, err := stdcopy.StdCopy(os.Stdout, os.Stderr, conn)
			out <- err
		}
		log.Info("connection to stdout/stderr closed")
	}()

	return out
}

func Attach(cmd common.ExecAttach) error {
	token := auth.GetToken()
	if token == nil {
		return fmt.Errorf("not logged in")
	}
	hijacker := &clientHijacker{}
	client := &http.Client{
		Transport: &http.Transport{
			DialTLS: hijacker.ConnectTLS,
			Dial:    hijacker.Connect,
		},
	}
	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("UPGRADE", conf.AttachExecURL(cmd.Test), bytes.NewReader(data))
	if err != nil {
		return err
	}
	token.SetAuthHeader(req)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		data, _ = ioutil.ReadAll(resp.Body)
		return fmt.Errorf("got back status %d: %s", resp.StatusCode, string(data))
	} else {
		log.Info("got back a 200 status code, attaching to the connection")
	}
	conn := hijacker.Conn()
	defer func() {
		for _, fn := range onCloseHandlers {
			fn()
		}
	}()
	if !cmd.Interactive {
		return <-handleOut(conn, cmd.TTY)
	}
	select {
	case err = <-handleOut(conn, cmd.TTY):
	case err = <-handleIn(conn, cmd.TTY):

	}
	if err != nil {
		log.Error(err)
	}

	return nil
}
