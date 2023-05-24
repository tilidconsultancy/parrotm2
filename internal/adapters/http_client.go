package adapters

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"
)

const (
	AUTH_HEADER  = "Authorization"
	CONTENT_TYPE = "Content-Type"
)

type (
	DialFunc func(ctx context.Context, network, addr string) (net.Conn, error)
)

func NewHttpClient(tmo time.Duration) *http.Client {
	return &http.Client{
		Timeout:   tmo,
		Transport: NewHttpTransport(tmo),
	}
}

func NewHttpTransport(timeout time.Duration) http.RoundTripper {
	return &http.Transport{
		Proxy:               nil,
		MaxIdleConns:        100,
		MaxConnsPerHost:     100,
		MaxIdleConnsPerHost: 100,
		ForceAttemptHTTP2:   false,
		DialContext:         CustomDialContext(timeout),
	}
}

func CustomDialContext(timeout time.Duration) DialFunc {
	return func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
		cerr := make(chan error)
		cconn := make(chan net.Conn)
		mctx, cancel := context.WithTimeout(ctx, timeout)
		go func() {
			<-mctx.Done()
			if d, ok := ctx.Deadline(); d.Before(time.Now()) && ok {
				cerr <- errors.New("custom dial tcp was break with timeout")
			}
		}()
		go func() {
			defer cancel()
			hsAddress, err := net.ResolveTCPAddr("tcp", addr)
			if err != nil {
				cerr <- err
				return
			}
			conn, err := net.DialTCP(network, nil, hsAddress)
			if err != nil {
				cerr <- err
				return
			}
			cconn <- conn
		}()
		select {
		case conn = <-cconn:
			return conn, nil
		case err = <-cerr:
			return nil, err
		}
	}
}
