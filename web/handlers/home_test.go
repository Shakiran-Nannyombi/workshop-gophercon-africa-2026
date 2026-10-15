package handlers_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"workshop/internal/common/logs"
	"workshop/web/handlers"

	"github.com/stretchr/testify/assert"
)

func TestHome(t *testing.T) {
	t.Parallel()

	logger := logs.NewProduction(io.Discard, slog.LevelError)
	rec := httptest.NewRecorder()
	handlers.Home(logger)(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil))

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "text/html; charset=utf-8", rec.Header().Get("Content-Type"))
}
