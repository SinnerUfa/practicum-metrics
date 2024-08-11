package server

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	codes "github.com/SinnerUfa/practicum-metric/internal/codes"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Hundlers(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	type pair struct {
		req  string
		want want
	}
	testsGetVoid := []pair{
		{
			req: "/value//test0",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			req: "/value/gauge/test0",
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			req: "/value/counter/test1",
			want: want{
				code: http.StatusNotFound,
			},
		},
		{
			req: "/value/abc/test1",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	testsPostVoid := []pair{
		{
			req: "/update/gauge/test0/100",
			want: want{
				code:        http.StatusOK,
				response:    "",
				contentType: "text/plain",
			},
		},
		{
			req: "/update/counter/test1/200",
			want: want{
				code:        http.StatusOK,
				response:    "",
				contentType: "text/plain",
			},
		},
		{
			req: "/update/counter/test1/abc",
			want: want{
				code:        http.StatusBadRequest,
				response:    codes.ErrRepParseInt.Error() + "\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			req: "/update/gauge/test1/abc",
			want: want{
				code:        http.StatusBadRequest,
				response:    codes.ErrRepParseFloat.Error() + "\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			req: "/update/abc/test1/100",
			want: want{
				code:        http.StatusBadRequest,
				response:    codes.ErrRepMetricNotSupported.Error() + "\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	testsGet := []pair{
		{
			req: "/value/gauge/test0",
			want: want{
				code:        http.StatusOK,
				response:    "100",
				contentType: "text/plain",
			},
		},
		{
			req: "/value/counter/test1",
			want: want{
				code:        http.StatusOK,
				response:    "200",
				contentType: "text/plain",
			},
		},
	}
	rep, _ := repository.New(context.Background(), repository.Config{Restore: false})
	store := rep.Storage()
	for _, test := range testsGetVoid {
		t.Run("testsGetVoid", func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.req, nil)
			w := httptest.NewRecorder()
			Routes(store).ServeHTTP(w, request)

			res := w.Result()

			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
		})
	}
	slog.Debug("rep void ", "list", store.GetList())
	for _, test := range testsPostVoid {
		t.Run("testsPostVoid", func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, test.req, nil)
			w := httptest.NewRecorder()
			Routes(store).ServeHTTP(w, request)

			res := w.Result()

			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
	slog.Debug("rep posted", "list", store.GetList())
	for _, test := range testsGet {
		t.Run("testsGet", func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.req, nil)
			w := httptest.NewRecorder()
			Routes(rep.Storage()).ServeHTTP(w, request)

			res := w.Result()

			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
	slog.Debug("rep endless", "list", store.GetList())
}
