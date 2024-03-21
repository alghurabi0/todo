package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go/v4"
	"todo.zaidalghurabi.net/internal/models"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	tasks         *models.TaskModel
	groups        *models.GroupModel
	boards        *models.BoardModel
	users         *models.UserModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	projectId := flag.String("project-id", "todogo-18674", "Google Cloud Project ID")
	credFile := flag.String("cred-file", "./internal/todogo-18674-firebase-adminsdk-5nd28-a2aabcc9ed.json", "Path to the credentials file")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	ctx := context.Background()
	db, auth, err := openDB(ctx, *projectId, *credFile)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a new instance of application containing the dependencies
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		tasks:         &models.TaskModel{DB: db},
		groups:        &models.GroupModel{DB: db},
		boards:        &models.BoardModel{DB: db},
		users:         &models.UserModel{DB: db, Auth: auth},
		templateCache: templateCache,
	}
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("starting the srv and listening on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func openDB(ctx context.Context, projectId string, credFile string) (*firestore.Client, *firebase.App, error) {
	db, err := firestore.NewClient(ctx, projectId, option.WithCredentialsFile(credFile))
	if err != nil {
		return nil, nil, err
	}
	//TODO - ping the database to check if it's connected
	docRef := db.Collection("ping").Doc("test")
	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		return nil, nil, err
	}
	var data map[string]interface{}
	if err := docSnapshot.DataTo(&data); err != nil {
		return nil, nil, err
	}
	expectedValue := "pong"
	if value, ok := data["ping"].(string); !ok || value != expectedValue {
		return nil, nil, fmt.Errorf("ping test failed, expected %s, got %s", expectedValue, value)
	}

	// TODO - use environment variables to set the credentials
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(credFile))
	if err != nil {
		return nil, nil, err
	}

	return db, app, nil
}
