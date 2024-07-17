package poll

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"poll-service/internal/entities"
	"poll-service/lib/http/response"
	"poll-service/lib/slogkz"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Poller interface {
	Poll(ctx context.Context, pollReq entities.PollRequest) error
}

type VoteSaver interface {
	Save(ctx context.Context, voteRequest entities.Request) error
}

func New(log *slog.Logger, poller Poller, saver VoteSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := `http.handlers.url.poll.New`

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		defer cancel()

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
			slog.String("remote_addr", r.RemoteAddr),
		)

		var voteRequest entities.Request

		if err := render.DecodeJSON(r.Body, &voteRequest.Request); err != nil {
			log.Error("Can't decode request body", slogkz.Err(err))

			render.JSON(w, r, response.Error("Can't decode request body"))

			return
		}

		valid := validator.New()

		log.Info("Request body decoded", slog.Any("poll", voteRequest))

		if err := valid.Struct(voteRequest.Request); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("Incorrect request data:", slogkz.Err(err))

			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		voteRequest.RemoteAddr = r.RemoteAddr

		if err := saver.Save(ctx, voteRequest); err != nil {
			log.Error("User's vote can't be taken into account", slogkz.Err(err))

			render.JSON(w, r, response.Error(err.Error()))

			return
		}

		if err := poller.Poll(ctx, voteRequest.Request); err != nil {
			log.Error("Can't save user's vote", slogkz.Err(err))

			render.JSON(w, r, response.Error(err.Error()))

			return
		}

		log.Info("Vote has been taken into account")

		render.JSON(w, r, response.OK())
	}
}
