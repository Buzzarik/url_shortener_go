package handlers

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/service"

	"github.com/julienschmidt/httprouter"
)

func GetAlias(app *service.Application) httprouter.Handle {
	return func (w http.ResponseWriter, r *http.Request, ps httprouter.Params){
		const op = "GetAlias";

		//TODO: распарсить ответ
		alias := ps.ByName("alias");
		//TODO: сверить данные
		if (alias == ""){
			app.ErrorResponse(w, http.StatusBadRequest, "Alias is required");
			app.Log.Info("Error no content request");
			return;
		}

		//TODO: запрос к бд
		url, err := app.StorageUrl.GetUrl(alias);
		if (err != nil){
			app.ErrorResponse(w, http.StatusInternalServerError, "Server internal error");
			app.Log.Error("Failed get user from StorageUrl");
			app.Log.Debug("Failed get user from StorageUrl",
						slog.String("error", err.Error()),
						slog.String("place", op));
			return;
		}

		if (url == nil){
			app.ErrorResponse(w, http.StatusNotFound, "Invalid alias");
			app.Log.Info("alias not found",
					slog.String("alias", alias));	
			return;
		}

		//TODO: вывод ответа
		http.Redirect(w, r, url.UrlText, http.StatusFound);
		app.Log.Info("Redirect success", 
				slog.String("alias", url.Alias),
				slog.String("url", url.UrlText));	
	}
}