package graceful_test

import (
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/penkovski/graceful"
)

type handler struct {
	requestTime time.Duration
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(h.requestTime)
}

func TestShutdown(t *testing.T) {
	tests := []struct {
		// input
		name    string
		addr    string
		timeout time.Duration
		reqTime time.Duration

		// desired outcome
		err                 error
		maxShutdownDuration time.Duration
	}{
		{
			name:    "without timeout",
			addr:    ":58430",
			timeout: 0,
			reqTime: 200 * time.Millisecond,

			err:                 nil,
			maxShutdownDuration: 300 * time.Millisecond,
		},
		{
			name:    "with timeout higher than request processing time",
			addr:    ":58431",
			timeout: 500 * time.Millisecond,
			reqTime: 200 * time.Millisecond,

			err:                 nil,
			maxShutdownDuration: 300 * time.Millisecond,
		},
		{
			name:    "with timeout lower than request processing time",
			addr:    ":58432",
			timeout: 50 * time.Millisecond,
			reqTime: 200 * time.Millisecond,

			err:                 errors.New("context deadline exceeded"),
			maxShutdownDuration: 100 * time.Millisecond,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			srv := &http.Server{
				Addr:    test.addr,
				Handler: &handler{requestTime: test.reqTime},
			}

			reqerr := make(chan error, 1)
			go func() {
				if err := graceful.Shutdown(srv, test.timeout); err != nil {
					reqerr <- err
				}
			}()

			go func() {
				_, err := http.Get("http://localhost" + test.addr)
				reqerr <- err
			}()

			// wait a while so the HTTP request could be sent
			time.Sleep(50 * time.Millisecond)

			start := time.Now()

			proc, err := os.FindProcess(os.Getpid())
			require.NoError(t, err)
			require.NoError(t, proc.Signal(os.Interrupt))

			err = <-reqerr

			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.True(t, time.Now().Sub(start) < test.maxShutdownDuration)
		})
	}
}
