package proxy

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

	"io"

	"bufio"

	"github.com/dgageot/docker-machine-proxy/machine"
)

type DockerMachineProxy struct {
	Machine *machine.DockerMachine
}

func (dmp *DockerMachineProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := proxyHttp(dmp.Machine, w, r); err != nil {
		panic(err)
	}
}

func proxyHttp(dm *machine.DockerMachine, writer http.ResponseWriter, r *http.Request) error {
	req, err := http.NewRequest(r.Method, fmt.Sprintf("%s%s", dm.Url, r.URL), r.Body)
	if err != nil {
		return err
	}
	req.Header = r.Header

	tlsConfig, err := dm.ReadTLSConfig()
	if err != nil {
		return err
	}

	underlying, err := tls.Dial("tcp", dm.Url[6:], tlsConfig)
	if err != nil {
		return err
	}

	defer underlying.Close()

	requestErrors := make(chan error)
	go func() {
		requestErrors <- req.Write(underlying)
	}()

	resp, err := http.ReadResponse(bufio.NewReader(underlying), r)
	if err != nil {
		return err
	}

	for k, v := range resp.Header {
		writer.Header()[k] = v
	}
	writer.WriteHeader(resp.StatusCode)

	if resp.StatusCode == 101 {
		err := upgradeToRaw(writer, underlying)
		if err != nil {
			return err
		}
	} else {
		defer resp.Body.Close()

		_, err := io.Copy(writer, resp.Body)
		if err != nil {
			return err
		}
	}

	return <-requestErrors
}

func upgradeToRaw(writer http.ResponseWriter, underlying io.ReadWriteCloser) error {
	hj, ok := writer.(http.Hijacker)
	if !ok {
		return errors.New("Server doesn't support hijacking")
	}

	conn, buf, err := hj.Hijack()
	if err != nil {
		return err
	}

	defer conn.Close()
	buf.Flush()

	toConnErrors := make(chan error)
	go func() {
		_, err := io.Copy(conn, underlying)
		if err != nil {
			toConnErrors <- err
		} else {
			toConnErrors <- conn.(CloseWriter).CloseWrite()
		}
	}()

	_, err = io.Copy(underlying, conn)
	if err != nil {
		return err
	}

	return <-toConnErrors
}

type CloseWriter interface {
	CloseWrite() error
}
