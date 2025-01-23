package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"url-shortener/internal/service"

	"github.com/julienschmidt/httprouter"
)

func Auth(app *service.Application, h httprouter.Handle) httprouter.Handle {
	return func (w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authorizationHeader := r.Header.Get("Authorization");

		authorizationHeader = strings.TrimSpace(authorizationHeader)

		parts := strings.Split(authorizationHeader, " ")
		app.Log.Debug("d",
			slog.String("1", parts[0]),
			slog.String("2", parts[1]))	
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			app.ErrorResponse(w, http.StatusUnauthorized, "You unauthorized");
			app.Log.Info("Incorrect Auth");
			return;
		}
	
		//TODO: запрос на сервис
		type Input struct {
			Token string `json:"token"`
			IdApi int	`json:"id_api"`
		};

		in := Input{
			Token: parts[1],
			IdApi: app.Cnf.Server.IdApi,
		};

		data, err := json.Marshal(in);

		app.Log.Debug(string(data));

		if err != nil {
			app.ErrorResponse(w, http.StatusInternalServerError, "Server internal error");
			app.Log.Error("Failed parse json");
			app.Log.Debug("Failed parse json",
					slog.String("error", err.Error()));
			return;
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:5555/Verify", bytes.NewBuffer(data));
		req.Header.Set("Content-Type", "application/json");

		if err != nil {
			app.ErrorResponse(w, http.StatusInternalServerError, "Server internal error");
			app.Log.Error("Failed create request");
			app.Log.Debug("Failed create request",
					slog.String("error", err.Error()));
			return;
		}

		client := &http.Client{};
		resp, err := client.Do(req)

		if err != nil {
			app.ErrorResponse(w, http.StatusInternalServerError, "Server internal error");
			app.Log.Error("Failed sent service auth");
			app.Log.Debug("Failed sent service auth",
					slog.String("error", err.Error()));
			return;
		}
		
		if resp.StatusCode == http.StatusOK {
			h(w, r, ps);
			return;
		}

		app.ErrorResponse(w, http.StatusUnauthorized, "You unauthorized");
		app.Log.Info("Token Invalid");
		return;
	}
}