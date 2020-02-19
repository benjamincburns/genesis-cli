package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/whiteblock/genesis-cli/pkg/auth"

	"github.com/docker/cli/cli/streams"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/pkg/term"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/whiteblock/httputils/websockets"
	"github.com/whiteblock/utility/common"
)

var onCloseHandlers []func()

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
func handleIn(conn io.ReadWriter, tty bool) <-chan error {
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

func handleOut(conn io.ReadWriter, tty bool) <-chan error {
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
	req := &http.Request{Header: http.Header{}}
	token.SetAuthHeader(req)

	dialer := &websocket.Dialer{}
	raw, _, err := dialer.Dial(conf.AttachExecURL(cmd.Test),
		http.Header{"Authorization": req.Header["Authorization"]})
	if err != nil {
		return err
	}
	err = raw.WriteJSON(cmd)
	if err != nil {
		return err
	}
	defer func() {
		for _, fn := range onCloseHandlers {
			fn()
		}
	}()
	conn := websockets.Wrapper{Conn: raw}
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
