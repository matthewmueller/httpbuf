package httpbuf

import (
	"bytes"
	"net/http"

	"github.com/felixge/httpsnoop"
)

func Wrap(w http.ResponseWriter) *ResponseWriter {
	rw := &ResponseWriter{
		Body:    new(bytes.Buffer),
		Headers: w.Header().Clone(),
		inner:   w,
	}
	rw.ResponseWriter = httpsnoop.Wrap(w, httpsnoop.Hooks{
		Header: func(_ httpsnoop.HeaderFunc) httpsnoop.HeaderFunc {
			return func() http.Header {
				if rw.wrote {
					return rw.inner.Header()
				}
				return rw.Headers
			}
		},
		WriteHeader: func(_ httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
			return func(code int) {
				if rw.wrote {
					return
				}
				rw.Status = code
			}
		},
		Write: func(_ httpsnoop.WriteFunc) httpsnoop.WriteFunc {
			return func(b []byte) (int, error) {
				return rw.Body.Write(b)
			}
		},
		Flush: func(flush httpsnoop.FlushFunc) httpsnoop.FlushFunc {
			return func() {
				rw.Flush()
			}
		},
	})
	return rw
}

type ResponseWriter struct {
	http.ResponseWriter
	Status  int
	Body    *bytes.Buffer
	Headers http.Header
	inner   http.ResponseWriter
	wrote   bool
}

func (rw *ResponseWriter) Flush() {
	// Only write status code once to avoid: "http: superfluous
	// response.WriteHeader". Not concurrency safe.
	if !rw.wrote {
		headers := rw.inner.Header()
		for k := range rw.Headers {
			headers.Set(k, rw.Headers.Get(k))
		}
		if rw.Status == 0 {
			rw.Status = http.StatusOK
		}
		rw.inner.WriteHeader(rw.Status)
		rw.wrote = true
	}
	rw.Body.WriteTo(rw.inner)
}
