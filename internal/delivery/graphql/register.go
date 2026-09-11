package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
)

func RegisterHandlers(app *mux.Router, resolver *Resolver) {
	server := handler.NewDefaultServer(
		NewExecutableSchema(
			Config{
				Resolvers: resolver,
			},
		),
	)

	app.Handle("/graphql", server)
	app.Handle("/", playground.Handler("", "/graphql")).Methods(http.MethodGet)
}
