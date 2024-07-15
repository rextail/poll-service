package main

import (
	"log/slog"
	"net/http"
	"os"
	"poll-service/internal/repository"
	"poll-service/internal/repository/mongo"
	"poll-service/lib/logger/handlers/slogpretty"
	"poll-service/lib/slogkz"

	midlogger "poll-service/internal/http/middleware/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"poll-service/internal/http/handlers/poll/create"
	"poll-service/internal/http/handlers/poll/poll"
	"poll-service/internal/http/handlers/poll/result"
)

// TODO: env, config
const conn = `mongodb://rextail:s3cr7tp4ss@localhost:8081/`

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	log := setupLogger(envLocal)

	client, err := mongo.New(conn)
	if err != nil {
		log.Error("can't initialize database connection", slogkz.Err(err))
		os.Exit(1)
	}

	repo := repository.New(client)

	log.Info(
		"Initialized repositories",
	)

	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)
	router.Use(midlogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Post("/createPoll", create.New(log, repo))
	router.Post("/poll", poll.New(log, repo, repo))
	router.Post("/getResult", result.New(log, repo))

	srv := &http.Server{
		Addr:    "localhost:8000",
		Handler: router,
	}
	srv.ListenAndServe()

	// jn := `{"poll_id":0,"name":"Best Pokemon!?",
	// 			"variants":
	// 				[{"votes":0,"variant_id":0,"text":"Pickachu"},
	// 				{"votes":0,"variant_id":1,"text":"Charmander"},
	// 				{"votes":0,"variant_id":2,"text":"Mewtwo"}]}`
	// var poll entities.Poll
	// err = json.Unmarshal([]byte(jn), &poll)
	// if err != nil {
	// 	log.Error("can't unmarshal incoming json", slogkz.Err(err))
	// 	os.Exit(1)
	// }
	// err = client.CreatePoll(context.Background(), poll)
	// if err != nil {
	// 	log.Error("unable ", slogkz.Err(err))
	// 	os.Exit(1)
	// }

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
