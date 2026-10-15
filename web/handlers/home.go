package handlers

import (
	"log/slog"
	"net/http"
	"workshop/web/views/pages"
)

func Home(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := pages.Home().Render(ctx, w); err != nil {
			logger.ErrorContext(ctx, "handlers: rendering home: "+err.Error())
			renderError(ctx, w, logger, http.StatusInternalServerError, "Unable to render page.")
		}
	}
}
