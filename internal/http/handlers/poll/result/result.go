package result

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"poll-service/internal/entities"
	"poll-service/internal/repository/repoerrs"
	"poll-service/lib/http/response"
	"poll-service/lib/slogkz"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Resulter interface {
	Result(ctx context.Context, req entities.GetResultRequest) (entities.Poll, error)
}

func New(log *slog.Logger, res Resulter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = `http.handlers.result.New`

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
			slog.String("remote_addr", r.RemoteAddr),
		)

		var resReq entities.GetResultRequest

		if err := render.DecodeJSON(r.Body, &resReq); err != nil {
			log.Error("Can't decode request body", slogkz.Err(err))

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("resReq", resReq))

		valid := validator.New()

		if err := valid.Struct(resReq); err != nil {
			validationError := err.(validator.ValidationErrors)

			render.JSON(w, r, response.ValidationError(validationError))

			return
		}

		poll, err := res.Result(ctx, resReq)
		if err != nil {
			if errors.Is(err, repoerrs.ErrPollNotExist) {
				log.Info("Can't get result: poll does not exist")

				render.JSON(w, r, response.Error("can't get result: poll does not exist"))

				return
			}

			log.Error("can't get result", slogkz.Err(err))

			render.JSON(w, r, response.Error(err.Error()))

			return
		}

		log.Info("Result:", slog.Any("poll", poll))

		render.JSON(w, r, poll)
	}
}
