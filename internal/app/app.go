package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/dubovayaRoshcha/posts-service/config"
	"github.com/dubovayaRoshcha/posts-service/internal/delivery/graphql"

	commentMemoryRepo "github.com/dubovayaRoshcha/posts-service/internal/repository/inmemory/comment"
	postMemoryRepo "github.com/dubovayaRoshcha/posts-service/internal/repository/inmemory/post"

	commentSqlRepo "github.com/dubovayaRoshcha/posts-service/internal/repository/sql/comment"
	postSqlRepo "github.com/dubovayaRoshcha/posts-service/internal/repository/sql/post"

	commentUC "github.com/dubovayaRoshcha/posts-service/internal/usecase/comment"
	postUC "github.com/dubovayaRoshcha/posts-service/internal/usecase/post"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Run(cfg *config.ProjectConfig) {
	r := mux.NewRouter()

	switch cfg.Storage {
	case "inmemory":
		postRepo := postMemoryRepo.NewPostRepo()
		commentRepo := commentMemoryRepo.NewCommentRepo()

		postUseCase := postUC.NewPostUseCase(postRepo)
		commentUseCase := commentUC.NewCommentUseCase(commentRepo, postRepo)

		resolver := graphql.NewResolver(postUseCase, commentUseCase)

		graphql.RegisterHandlers(r, resolver)

	case "postgres":
		dsn := BuildDSN(cfg.Postgres)

		pool, err := pgxpool.New(context.Background(), dsn)
		if err != nil {
			log.Fatalf("cannot create pgx pool: %v", err)
		}
		defer pool.Close()

		postRepo := postSqlRepo.NewPostRepo(pool)
		commentRepo := commentSqlRepo.NewCommentRepo(pool)

		postUseCase := postUC.NewPostUseCase(postRepo)
		commentUseCase := commentUC.NewCommentUseCase(commentRepo, postRepo)

		resolver := graphql.NewResolver(postUseCase, commentUseCase)

		graphql.RegisterHandlers(r, resolver)

	default:
		log.Fatalf("unknown storage: %s", cfg.Storage)
	}

	serverAddress := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	srv := &http.Server{
		Handler: r,
		Addr:    serverAddress,
	}

	log.Printf("start listen: %s", serverAddress)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func BuildDSN(cfg config.Postgres) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.SslMode,
	)
}
