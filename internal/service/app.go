package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"url-shortener/internal/config"
	"url-shortener/internal/models"
)

type Envelope map[string]interface{};

type StorageUrl interface {
	SaveUrl(url *models.Url) error;
	GetUrl(alias string) (*models.Url, error);
};

type Application struct {
	Cnf config.Config;
	StorageUrl StorageUrl
	Log *slog.Logger
	//wg     sync.WaitGroup //updated
};

func New(
	cnf config.Config,
	StorageUrl StorageUrl,
	log *slog.Logger,
) *Application {
	return &Application{
		Cnf: cnf,
		StorageUrl: StorageUrl,
		Log: log,
	}
}


func (app *Application) ErrorResponse(w http.ResponseWriter, status int, message interface{}) error{
    const op = "Application.ErrorResponse";
	env := Envelope{"error": message}

    err := app.WriteJSON(w, status, env, nil)
    if err != nil {
        w.WriteHeader(500);
		return fmt.Errorf("%s: %w", op, err);
    }
	return nil;
}


//формирование ответа
func (app* Application) WriteJSON(w http.ResponseWriter, status int, body Envelope, headers http.Header) error {
	const op = "Application.WriteJSON";
	//NOTE: переводим все в json для ответа
    js, err := json.Marshal(body);
    if err != nil {
        return fmt.Errorf("%s: %w", op, err);
    }

    js = append(js, '\n');

	//NOTE: добавляем заголовки, если они будут
    for key, value := range headers {
        w.Header()[key] = value;
    }

    w.Header().Set("Content-Type", "application/json");
    w.WriteHeader(status);
    w.Write(js);

    return nil;
}

//корректное преобразование запроса
func (app *Application)ReadJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
    const op = "Application.ReadJSON";
	//NOTE: ограничиваем размер считывания запроса
	maxBytes := 104856;
    r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes));

    dec := json.NewDecoder(r.Body);
    dec.DisallowUnknownFields();

    err := dec.Decode(dst);
    if err != nil {
        var syntaxError *json.SyntaxError;
        var unmarshalTypeError *json.UnmarshalTypeError;
        var invalidUnmarshalError *json.InvalidUnmarshalError;
        switch {
        case errors.As(err, &syntaxError):
            return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset);

        case errors.Is(err, io.ErrUnexpectedEOF):
            return errors.New("body contains badly-formed JSON");

        case errors.As(err, &unmarshalTypeError):
            if unmarshalTypeError.Field != "" {
                return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field);
            }
            return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset);

        case errors.Is(err, io.EOF):
            return errors.New("body must not be empty");

        case errors.As(err, &invalidUnmarshalError):
            panic(err);
        default:
            return fmt.Errorf("%s: %w", op, err);
        }
    }
	return nil;
}