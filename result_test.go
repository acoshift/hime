package hime_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/acoshift/hime"
	"github.com/stretchr/testify/assert"
)

func invokeHandler(h http.Handler, method string, target string, body io.Reader) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, body)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestResult(t *testing.T) {
	t.Parallel()

	t.Run("StatusCode", func(t *testing.T) {
		t.Parallel()

		app := hime.New().
			Handler(hime.H(func(ctx hime.Context) hime.Result {
				return ctx.Status(http.StatusNotFound).String("not found")
			}))

		w := invokeHandler(app, "GET", "/", nil)
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		assert.Equal(t, "not found", w.Body.String())
	})

	t.Run("StatusTest", func(t *testing.T) {
		t.Parallel()

		app := hime.New().
			Handler(hime.H(func(ctx hime.Context) hime.Result {
				return ctx.Status(http.StatusTeapot).StatusText()
			}))

		w := invokeHandler(app, "GET", "/", nil)
		assert.Equal(t, http.StatusTeapot, w.Result().StatusCode)
		assert.Equal(t, http.StatusText(http.StatusTeapot), w.Body.String())
	})

	t.Run("NotFound", func(t *testing.T) {
		t.Parallel()

		app := hime.New().
			Handler(hime.H(func(ctx hime.Context) hime.Result {
				return ctx.NotFound()
			}))

		w := invokeHandler(app, "GET", "/", nil)
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		assert.Equal(t, "404 page not found\n", w.Body.String())
	})

	t.Run("NoContent", func(t *testing.T) {
		t.Parallel()

		app := hime.New().
			Handler(hime.H(func(ctx hime.Context) hime.Result {
				return ctx.NoContent()
			}))

		w := invokeHandler(app, "GET", "/", nil)
		assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
	})

	t.Run("Bytes", func(t *testing.T) {
		t.Parallel()

		app := hime.New().
			Handler(hime.H(func(ctx hime.Context) hime.Result {
				return ctx.Bytes([]byte("hello hime"))
			}))

		w := invokeHandler(app, "GET", "/", nil)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, "hello hime", w.Body.String())
	})

	t.Run("File", func(t *testing.T) {
		t.Parallel()

		app := hime.New().
			Handler(hime.H(func(ctx hime.Context) hime.Result {
				return ctx.File("testdata/file.txt")
			}))

		w := invokeHandler(app, "GET", "/", nil)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, "file content", w.Body.String())
	})

	t.Run("JSON", func(t *testing.T) {
		t.Parallel()

		app := hime.New().
			Handler(hime.H(func(ctx hime.Context) hime.Result {
				return ctx.JSON(map[string]interface{}{"abc": "afg", "bbb": 123})
			}))

		w := invokeHandler(app, "GET", "/", nil)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
		assert.JSONEq(t, `{"abc":"afg","bbb":123}`, w.Body.String())
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		app := hime.New().
			Handler(hime.H(func(ctx hime.Context) hime.Result {
				return ctx.String("hello, hime")
			}))

		w := invokeHandler(app, "GET", "/", nil)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, "hello, hime", w.Body.String())
	})

	t.Run("Nothing", func(t *testing.T) {
		t.Parallel()

		app := hime.New().
			Handler(hime.H(func(ctx hime.Context) hime.Result {
				return ctx.Nothing()
			}))

		w := invokeHandler(app, "GET", "/", nil)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Empty(t, w.Body.String())
	})

	t.Run("BeforeRender", func(t *testing.T) {
		t.Parallel()

		app := hime.New().
			BeforeRender(func(h http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Cache-Control", "public, max-age=3600")
					h.ServeHTTP(w, r)
				})
			}).
			Handler(hime.H(func(ctx hime.Context) hime.Result {
				return ctx.String("hello, hime")
			}))

		w := invokeHandler(app, "GET", "/", nil)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, "public, max-age=3600", w.Header().Get("Cache-Control"))
	})
}