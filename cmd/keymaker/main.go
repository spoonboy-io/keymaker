package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spoonboy-io/keymaker/internal"
	"github.com/spoonboy-io/keymaker/internal/certificate"
	"github.com/spoonboy-io/keymaker/internal/config"
	"github.com/spoonboy-io/keymaker/internal/handlers"
	"github.com/spoonboy-io/keymaker/internal/store"
	"github.com/spoonboy-io/koan"
	"github.com/spoonboy-io/reprise"
	"golang.org/x/sync/errgroup"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var (
	version   = "Development build"
	goversion = "Unknown"
)

var logger *koan.Logger

func init() {
	logger = &koan.Logger{}

	// check/create certificates folder
	tlsPath := filepath.Join(".", internal.TLS_FOLDER)
	if err := os.MkdirAll(tlsPath, os.ModePerm); err != nil {
		logger.FatalError("Problem checking/creating 'certificates' folder", err)
	}

	// add self-signed certificate only if folder empty, if the cert expires it
	// it can be deleted so the code here creates a new cert.pem and key.pem file
	checkExist := fmt.Sprintf("%s/cert.pem", internal.TLS_FOLDER)
	if _, err := os.Stat(checkExist); errors.Is(err, os.ErrNotExist) {
		logger.Info("Creating self-signed TLS certificate for the server")
		if err := certificate.Make(logger); err != nil {
			logger.FatalError("Problem creating the certificate/key", err)
		}
	}
}

func main() {
	reprise.WriteSimple(&reprise.Banner{
		Name:         "Keymaker",
		Description:  "ID/Name Management Server",
		Version:      version,
		GoVersion:    goversion,
		WebsiteURL:   "https://spoonboy.io",
		VcsURL:       "https://github.com/spoonboy-io/keymaker",
		VcsName:      "Github",
		EmailAddress: "hello@spoonboy.io",
	})

	// create context & notifier
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()

	if err := config.CheckConfigExist(internal.CONFIGFILE); err != nil {
		logger.FatalError("config file not found", err)
	}

	fileData, err := config.ReadConfigFile(internal.CONFIGFILE)
	if err != nil {
		logger.FatalError("could not read config file", err)
	}

	cfg, err := config.Parse(fileData)
	if err != nil {
		logger.FatalError("config file could not be parsed", err)
	}

	// add some additional app specific config
	cfg["apiVersion"] = internal.APIVERSION
	// init store
	st := store.New(cfg)

	app := &handlers.App{
		Store:  st,
		Config: cfg,
		Logger: logger,
	}

	// init router & routes
	mux := mux.NewRouter()
	mux.HandleFunc(internal.APIVERSION+"/init/{sequence}", app.InitSequence).Methods("POST")
	mux.HandleFunc(internal.APIVERSION+"/next/{sequence}", app.GetNextInSequence).Methods("GET")
	// TODO need delete end point, possible confirm endpoint

	// create a server
	hostPort := net.JoinHostPort(internal.SRV_HOST, internal.SRV_PORT)
	srvTLS := &http.Server{
		Addr:         hostPort,
		Handler:      mux,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		// start HTTPS server
		logger.Info(fmt.Sprintf("starting HTTPS server on %s", hostPort))
		return srvTLS.ListenAndServeTLS(fmt.Sprintf("%s/cert.pem", internal.TLS_FOLDER), fmt.Sprintf("%s/key.pem", internal.TLS_FOLDER))
	})

	g.Go(func() error {
		<-gCtx.Done()
		return srvTLS.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		shutdownMsg := fmt.Sprintf("server shutdown (exit reason: %s)\n", err)
		fmt.Println("") // break after ^C
		logger.Info(shutdownMsg)
	}
}
