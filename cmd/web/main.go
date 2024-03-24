package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/makarellav/codecapsule/internal/models"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
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
	defer db.Close()

	app := &application{logger: logger, snippets: &models.SnippetModel{DB: db}}

	if err != nil {
		app.logger.Error(err.Error())

		os.Exit(1)
	}

	app.logger.Info("starting the server", "addr", *addr)
	err = http.ListenAndServe(*addr, app.routes())

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
