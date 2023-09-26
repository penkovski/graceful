package graceful

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Start the given HTTP server and wait for receiving
// a stop signal to gracefully stop it by waiting for
// {timeout} period of time to complete its current
// in-flight requests.
func Start(srv *http.Server, timeout time.Duration) error {
	done := make(chan error, 1)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c

		ctx := context.Background()
		var cancel context.CancelFunc
		if timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		done <- srv.Shutdown(ctx)
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return <-done
}

// StartCtx start the given HTTP server and wait for receiving
// a stop signal or context cancellation signal to gracefully stop
// it by waiting for {timeout} period of time to complete its current
// in-flight requests.
//
// The {timeout} period is always respected.
func StartCtx(ctx context.Context, srv *http.Server, timeout time.Duration) error {
	done := make(chan error, 1)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

		// wait for a signal or context cancellation
		select {
		case <-c:
		case <-ctx.Done():
		}

		ctx := context.Background()
		var cancel context.CancelFunc
		if timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		done <- srv.Shutdown(ctx)
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return <-done
}
