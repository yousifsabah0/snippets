package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yousifsabah0/snippets/internal/models/snippets"
	"github.com/yousifsabah0/snippets/internal/models/users"
)

type application struct {
	logger       *slog.Logger
	snippets     *snippets.SnippetModel
	users        *users.UserModel
	templateCace map[string]*template.Template
	formDecoder  *form.Decoder
	session      *scs.SessionManager
}

func main() {
	port := flag.String("port", ":8080", "HTTP network port")
	dsn := flag.String("dsn", "odyssey:odyssey@/snippets?parseTime=true", "Database source name")

	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdin, &slog.HandlerOptions{}))

	db, err := openDB(*dsn)

	defer func() {
		err := db.Close()
		if err != nil {
			logger.Error(err.Error(), "error", err)
			os.Exit(1)
		}
	}()

	if err != nil {
		logger.Error(err.Error(), "error", err)
		os.Exit(1)
	}

	tc, err := newTemplateCaceh()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	session := scs.New()
	session.Store = mysqlstore.New(db)
	session.Lifetime = 12 * time.Hour

	app := &application{
		logger:       logger,
		snippets:     &snippets.SnippetModel{DB: db},
		users:        &users.UserModel{DB: db},
		templateCace: tc,
		formDecoder:  formDecoder,
		session:      session,
	}

	srv := &http.Server{
		Handler:      app.routes(),
		Addr:         *port,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Error(err.Error(), "err", err)
		os.Exit(1)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
