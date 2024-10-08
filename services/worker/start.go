package worker

import (
	"context"
	"dots-api/bootstrap"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/valve"
	"github.com/urfave/cli/v2"
)

// Boot ...
type Boot struct {
	*bootstrap.App
}

var (
	// Flags ...
	Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "host",
			Value: "127.0.0.1:4000",
			Usage: "Run API service with custom host",
		},
	}
)

// Start main function to run the http host
func (app Boot) Start(c *cli.Context) error {
	var err error

	host := c.String("host")
	if len(host) == 0 {
		host = app.Config.GetString("app.host")
	}
	if app.Debug {
		log.Printf("Check Appointment Service -> Running on Debug Mode: On at host [%v]", host)
	}

	// gracefull shutdown handler
	valv := valve.New()
	baseCtx := valv.Context()

	// start new app
	r := chi.NewRouter()
	cr := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-Timezone",
			"X-CHANNEL",
			"X-PLAYER",
			"Access-Control-Allow-Headers",
			"X-Requested-With",
			"application/json",
			"Cache-Control",
			"multipart/form-data; boundary=<calculated when request is sent>",
			"multipart/form-data",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cr.Handler)
	if app.Debug {
		r.Use(middleware.Logger)
	}
	r.Use(app.Recoverer)
	r.Use(app.NotfoundMiddleware)

	// handle grace full shutdown
	srv := http.Server{Addr: host, Handler: r}
	srv.BaseContext = func(_ net.Listener) context.Context {
		return baseCtx
	}
	sng := make(chan os.Signal, 1)
	signal.Notify(sng, os.Interrupt)

	go func() {
		for range sng {
			fmt.Println("shutting down..")
			err = valv.Shutdown(20 * time.Second)
			if err != nil {
				log.Println("Can't shutdown this server until all process are done!")
			}

			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			err = srv.Shutdown(ctx)
			if err != nil {
				log.Println("Can't shutdown this server until all process are done!")
			}
			select {
			case <-time.After(21 * time.Second):
				fmt.Println("not all connections done")
			case <-ctx.Done():

			}
		}
	}()

	return srv.ListenAndServe()
}
