package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/makarellav/codecapsule/internal/models"
	"log/slog"
	"net/http"
	"os"
	"text/template"
	"time"
)

type snippetModel interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (*models.Snippet, error)
	Latest() ([]*models.Snippet, error)
}

type userModel interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
}

type application struct {
	logger         *slog.Logger
	snippets       snippetModel
	users          userModel
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "root:password@/codecapsule?parseTime=true", "DSN")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	logger.Info("opening db", "dsn", *dsn)
	db, err := openDB(*dsn)

	if err != nil {
		logger.Error(err.Error())

		os.Exit(1)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()

	if err != nil {
		logger.Error(err.Error())

		os.Exit(1)
	}

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  1 * time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.logger.Info("starting the server", "addr", srv.Addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")

	app.logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
