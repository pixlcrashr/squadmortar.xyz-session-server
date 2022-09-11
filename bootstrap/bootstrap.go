package bootstrap

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/auth"
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/crypto"
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/graphql"
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/graphql/generated"
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/log"
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/session/storage"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Bootstrapper interface {
	Listen() error
}

type bootstrapper struct {
	host                string
	port                uint16
	privateKey          *ecdsa.PrivateKey
	wsKeepAlive         time.Duration
	allowedOrigins      []string
	logger              *zap.Logger
	enablePlayground    bool
	enableIntrospection bool
}

func New(
	host string,
	port uint16,
	privateKeyFilepath string,
	wsKeepAlive time.Duration,
	allowedOrigins []string,
	enablePlayground bool,
	enableIntrospection bool) (Bootstrapper, error) {
	logger, _ := zap.NewProduction()

	privateKey, err := crypto.ParseEcdsaPemPrivateKey(privateKeyFilepath)
	if err != nil {
		return nil, err
	}

	logger.Info(
		"created bootstrapper",
		zap.String("host", host),
		zap.Uint16("port", port),
		zap.Duration("wsKeepAlive", wsKeepAlive),
		zap.Strings("allowedOrigins", allowedOrigins),
		zap.Bool("enablePlayground", enablePlayground),
		zap.Bool("enableIntrospection", enableIntrospection))

	return &bootstrapper{
		host:                host,
		port:                port,
		privateKey:          privateKey,
		wsKeepAlive:         wsKeepAlive,
		allowedOrigins:      allowedOrigins,
		logger:              logger,
		enablePlayground:    enablePlayground,
		enableIntrospection: enableIntrospection,
	}, nil
}

func (b *bootstrapper) Listen() error {
	b.logger.Info("setting up listener")
	config := generated.Config{
		Resolvers: &graphql.Resolver{
			EcdsaKey:       b.privateKey,
			SessionStorage: storage.NewStorage(),
		},
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(config))
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

	if b.enableIntrospection {
		b.logger.Info("enabled introspection. to disable, remove the --enabled-playground flag")
		srv.Use(extension.Introspection{})
	}

	r := chi.NewRouter()

	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger: &log.ChiLogger{
			Logger: b.logger,
		},
		NoColor: true,
	})
	r.Use(middleware.Recoverer)

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Authorization"},
		AllowedMethods:   []string{http.MethodPost, http.MethodGet},
		AllowCredentials: true,
		Debug:            false,
	}).Handler)

	r.Use(auth.AuthenticationMiddleware(&b.privateKey.PublicKey))

	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("This is the Squadmortar.xyz™ session server™ <3."))
		writer.WriteHeader(http.StatusOK)
	})
	r.Handle("/graphql", srv)

	if b.enablePlayground {
		b.logger.Info("enabled playground. to disable, remove the --enable-playground flag")
		r.Handle("/graphql/playground", playground.Handler("Squadmortar Session Server", fmt.Sprintf("http://%s/graphql", b.addr())))
	}

	b.logger.Info("now listening", zap.String("host", b.host), zap.Uint16("port", b.port))
	defer b.logger.Info("stopped http server", zap.String("host", b.host), zap.Uint16("port", b.port))
	return http.ListenAndServe(b.addr(), r)
}

func (b *bootstrapper) addr() string {
	return fmt.Sprintf("%s:%d", b.host, b.port)
}
