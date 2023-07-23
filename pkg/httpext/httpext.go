package httpext

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	// ErrCTXDone is used to notify that the reason of the select unlock was context.Done().
	ErrCTXDone = errors.New("parent context was canceled")
	// ErrHTTPListenerErrorReason used to notify that the reason of the select unlock was some execution error.
	ErrHTTPListenerErrorReason = errors.New("error happened wile listening http")
)

// ServerConfig holds settings to be used for the HttpServer.
type ServerConfig struct {
	// Host is the host to listen on.
	Host string `split_words:"true" required:"true"`
	// Port is the port to listen on.
	Port int `split_words:"true" required:"true"`
	// MaxHeaderBytes is the maximum number of bytes the server will
	// read parsing the request header's keys and values, including the
	// request line. It does not limit the size of the request body.
	// It can be set to 0 to disable the limit. Default is 1MB
	MaxHeaderBytes int `split_words:"true" default:"1048576"`
	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body.
	ReadTimeout time.Duration `split_words:"true" default:"30s"`
	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are zero,
	// ReadHeaderTimeout is used. Default is 120s
	IdleTimeout time.Duration `split_words:"true" default:"120s"`
	// ReadHeaderTimeout is the amount of time allowed to read request headers.
	// Default is 10s
	ReadHeaderTimeout time.Duration `split_words:"true" default:"10s"`
	// WriteTimeout is the maximum duration before timing out writes of the response.
	// Default is 60s
	WriteTimeout time.Duration `split_words:"true" default:"60s"`
}

// HTTPServer runs the http server.
type HTTPServer struct {
	srv *http.Server
}

// NewHTTPServer creates a new http server.
func NewHTTPServer(cfg ServerConfig, router http.Handler) HTTPServer {
	return HTTPServer{
		srv: &http.Server{
			Addr:              fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler:           router,
			ReadTimeout:       cfg.ReadTimeout,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       cfg.IdleTimeout,
			MaxHeaderBytes:    cfg.MaxHeaderBytes,
		},
	}
}

// ListenAndServe starts the http server.
func (o HTTPServer) ListenAndServe(ctx context.Context) error {
	if o.srv.Handler == nil {
		return errors.New("no router defined")
	}

	errCh := o.listenAndServe()
	err := <-o.mergeStopChannels(ctx.Done(), errCh)

	if errors.Is(err, ErrHTTPListenerErrorReason) {
		return err
	}

	return o.Close(context.Background())
}

func (o HTTPServer) listenAndServe() <-chan error {
	errCh := make(chan error, 1)
	listenServeFn := func() error {
		switch o.srv.TLSConfig {
		case nil:
			return o.srv.ListenAndServe()
		default:
			// TODO add TLS support
			return o.srv.ListenAndServeTLS("", "")
		}
	}

	go func() {
		defer close(errCh)

		if err := listenServeFn(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	return errCh
}

func (HTTPServer) mergeStopChannels(ctxDone <-chan struct{}, listenErrCh <-chan error) <-chan error {
	reasonCh := make(chan error)
	go func() {
		defer close(reasonCh)
		select {
		case <-ctxDone:
			reasonCh <- ErrCTXDone
		case err := <-listenErrCh:
			reasonCh <- fmt.Errorf("%w: %v", ErrHTTPListenerErrorReason, err)
		}
	}()

	return reasonCh
}

// Close tears down all the HttpServer.
func (o HTTPServer) Close(ctx context.Context) error {
	cancelCtx, cancel := context.WithTimeout(ctx, shutDownTimeout)
	defer cancel()

	if err := o.srv.Shutdown(cancelCtx); err == nil {
		return nil
	}

	return o.srv.Close()
}

const (
	shutDownTimeout = 5 * time.Second
)
