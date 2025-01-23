package handlers

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/models"
	"url-shortener/internal/service"

	"github.com/julienschmidt/httprouter"
)

const constraint = "CONSTRAINT";

func SaveUrl(app *service.Application) httprouter.Handle {
	return func (w http.ResponseWriter, r *http.Request,  ps httprouter.Params){
		const op = "SaveUrl";
		//TODO: распарсить данные
		var input struct {
			Url 	string `json:"url"`
			Alias 	string `json:"alias"`
		};

		err := app.ReadJSON(w, r, &input);
		if err != nil {
			app.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload");
			app.Log.Error("Error reading JSON in request");
			app.Log.Debug("Error reading JSON in request",
						slog.String("error", err.Error()),
						slog.String("place", op));
		}

		//TODO: сверить данные
		if input.Url == "" || input.Alias == "" {
			app.ErrorResponse(w, http.StatusBadRequest, "Url and alias are required");
			app.Log.Info("Error no content request");
			return;
		}

		url := &models.Url{
			UrlText: input.Url,
			Alias: input.Alias,
		};

		//TODO: запрос в БД
		err = app.StorageUrl.SaveUrl(url);

		if err != nil && err.Error() == constraint {
			app.ErrorResponse(w, http.StatusConflict, "Alias already exists");
			app.Log.Info("Alias already exists",
					slog.String("alias", input.Alias));
			return;
		}

		//TODO: ответ
		err = app.WriteJSON(w, http.StatusCreated, service.Envelope{
				"success": true,
				"message": "Alias created successfully",
		}, nil);

		if err != nil {
			app.ErrorResponse(w, http.StatusInternalServerError,
				"Server internal error");
			app.Log.Error("Error write JSON in response");
			app.Log.Debug("Error write JSON in response",
						slog.String("error", err.Error()),
						slog.String("place", op));
			return;
		}

		app.Log.Info("Alias created",
			slog.String("url", input.Url),
			slog.String("alias", input.Alias));
	}
}