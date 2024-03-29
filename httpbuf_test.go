package httpbuf_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/matryer/is"
	"github.com/matthewmueller/httpbuf"
)

func TestHeadersNormal(t *testing.T) {
	is := is.New(t)
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	is.NoErr(err)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-A", "A")
		w.Write([]byte("Hello, world!"))
		w.Header().Add("X-B", "B")
	})
	h.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(res.StatusCode, 200)
	is.Equal(res.Header.Get("X-A"), "A")
	is.Equal(res.Header.Get("X-B"), "")
}

func TestHeadersWrapped(t *testing.T) {
	is := is.New(t)
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	is.NoErr(err)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-A", "A")
		w.Write([]byte("Hello, world!"))
		w.Header().Add("X-B", "B")
	})
	rw := httpbuf.Wrap(rec)
	h.ServeHTTP(rw, req)
	rw.Flush()
	res := rec.Result()
	is.Equal(res.StatusCode, 200)
	is.Equal(res.Header.Get("X-A"), "A")
	is.Equal(res.Header.Get("X-B"), "B")
}

func TestWriteStatusNormal(t *testing.T) {
	is := is.New(t)
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	is.NoErr(err)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-A", "A")
		w.WriteHeader(201)
		w.Write([]byte("Hello, world!"))
		w.Header().Add("X-B", "B")
	})
	h.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(res.StatusCode, 201)
	is.Equal(res.Header.Get("X-A"), "A")
	is.Equal(res.Header.Get("X-B"), "")
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "Hello, world!")
}

func TestWriteStatusWrapped(t *testing.T) {
	is := is.New(t)
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	is.NoErr(err)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-A", "A")
		w.WriteHeader(201)
		w.Write([]byte("Hello, world!"))
		w.Header().Add("X-B", "B")
	})
	rw := httpbuf.Wrap(rec)
	h.ServeHTTP(rw, req)
	rw.Flush()
	res := rec.Result()
	is.Equal(res.StatusCode, 201)
	is.Equal(res.Header.Get("X-A"), "A")
	is.Equal(res.Header.Get("X-B"), "B")
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "Hello, world!")
}

func TestFlushNormal(t *testing.T) {
	is := is.New(t)
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	is.NoErr(err)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-A", "A")
		w.WriteHeader(201)
		w.Write([]byte("Hello, world!"))
		flush, ok := w.(http.Flusher)
		if ok {
			flush.Flush()
			flush.Flush()
		}
		w.Header().Add("X-B", "B")
	})
	rw := httpbuf.Wrap(rec)
	h.ServeHTTP(rw, req)
	rw.Flush()
	res := rec.Result()
	is.Equal(res.StatusCode, 201)
	is.Equal(res.Header.Get("X-A"), "A")
	is.Equal(res.Header.Get("X-B"), "")
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "Hello, world!")
}

func TestFlushWrapped(t *testing.T) {
	is := is.New(t)
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	is.NoErr(err)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
		w.Header().Add("X-A", "A")
		flush, ok := w.(http.Flusher)
		if ok {
			flush.Flush()
			w.Write([]byte("yoyo"))
			flush.Flush()
			w.Write([]byte("zzz"))
		}
		w.Header().Add("X-B", "B")
	})
	rw := httpbuf.Wrap(rec)
	h.ServeHTTP(rw, req)
	is.Equal(rw.Status, 200)
	is.Equal(rw.Headers.Get("X-A"), "A")
	is.Equal(rw.Headers.Get("X-B"), "")
	is.Equal(string(rw.Body), "Hello, world!yoyozzz")
	rw.Flush()
	res := rec.Result()
	is.Equal(res.StatusCode, 200)
	is.Equal(res.Header.Get("X-A"), "A")
	is.Equal(res.Header.Get("X-B"), "")
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "Hello, world!yoyozzz")
}

func TestFlushStatusWrapped(t *testing.T) {
	is := is.New(t)
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	is.NoErr(err)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
		w.WriteHeader(201)
		w.Header().Add("X-A", "A")
		flush, ok := w.(http.Flusher)
		if ok {
			flush.Flush()
			w.Write([]byte("yoyo"))
			flush.Flush()
		}
		w.Header().Add("X-B", "B")
	})
	rw := httpbuf.Wrap(rec)
	h.ServeHTTP(rw, req)
	rw.Flush()
	res := rec.Result()
	is.Equal(res.StatusCode, 201)
	is.Equal(res.Header.Get("X-A"), "A")
	is.Equal(res.Header.Get("X-B"), "")
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "Hello, world!yoyo")
}

func TestMultipleWriteHeaders(t *testing.T) {
	is := is.New(t)
	req, err := http.NewRequest("POST", "/", nil)
	rec := httptest.NewRecorder()
	is.NoErr(err)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusSeeOther)
	})
	rw := httpbuf.Wrap(rec)
	h.ServeHTTP(rw, req)
	rw.Flush()
	res := rec.Result()
	is.Equal(res.StatusCode, http.StatusSeeOther)
	req, err = http.NewRequest("POST", "/", nil)
	is.NoErr(err)
	rec = httptest.NewRecorder()
	rw = httpbuf.Wrap(rec)
	h.ServeHTTP(rw, req)
	rw.Flush()
	res = rec.Result()
	is.Equal(res.StatusCode, http.StatusSeeOther)
}

func TestFileServer(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"index.html": &fstest.MapFile{
			Data: []byte("Hello, world!"),
		},
	}
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := httpbuf.Wrap(w)
			next.ServeHTTP(rw, r)
			is.Equal(rw.Status, 200)
			is.Equal(string(rw.Body), "Hello, world!")
			rw.Flush()
		})
	}
	handler := middleware(http.FileServer(http.FS(fsys)))
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(res.StatusCode, 200)
	body, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(body), "Hello, world!")
}
