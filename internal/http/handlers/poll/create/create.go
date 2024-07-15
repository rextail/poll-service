package create

import (
	"context"
	"log/slog"
	"net/http"
	"poll-service/internal/entities"
	"poll-service/lib/http/response"
	"poll-service/lib/slogkz"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type PollCreator interface {
	CreatePoll(ctx context.Context, poll entities.CreatePollRequest) error
}

func New(log *slog.Logger, pollCreator PollCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := `http.handlers.url.create.New`

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
			slog.String("remote_addr", r.RemoteAddr),
		)

		var poll entities.CreatePollRequest

		err := render.DecodeJSON(r.Body, &poll)
		if err != nil {
			log.Error("failed to decode request body", slogkz.Err(err))

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		valid := validator.New()

		log.Info("request body decoded", slog.Any("poll", poll))

		if err = valid.Struct(poll); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("validation failed: incorrect poll body", slogkz.Err(err))

			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		if err = pollCreator.CreatePoll(context.TODO(), poll); err != nil {
			log.Error("can't save poll in the database", slogkz.Err(err))

			render.JSON(w, r, response.Error("save procedure failed"))

			return
		}

		log.Info("Poll has been created")

		render.JSON(w, r, response.OK())

	}
}
