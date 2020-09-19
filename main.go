package main

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/apex/log"
	"github.com/golang/gddo/httputil/header"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config
type Config struct {
	APIKey      string `required:"true" split_words:"true"`
	ListID      string `required:"true" split_words:"true"`
	APIServer   string `required:"true" split_words:"true"`
	StatusIfNew string `split_words:"true" default:"pending"`
	Port        string `default:"3000"`
}

type server struct {
	config Config
}

type addSubscriberRequest struct {
	Email string
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "UP"}`))

		return
	}

	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := `"Content-Type header is not application/json"`
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	api := NewMcAPI(s.config)

	switch r.Method {
	case "POST":
		var request addSubscriberRequest

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		err := dec.Decode(&request)
		if err != nil {
			msg := err.Error()
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		if !isEmailValid(request.Email) {
			http.Error(w, "Email is malformed", http.StatusBadRequest)
			return
		}

		if err := api.AddSubscriber(request.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "Added subscriber"}`))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	}
}

func main() {
	loadFromEnvFile()

	var config Config
	if err := envconfig.Process("", &config); err != nil {
		log.WithError(err).Fatal("Error")
	}

	s := &server{config: config}
	http.Handle("/", s)
	log.Infof("Listening on port %s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.WithError(err).Fatal("Error")
	}
}

func loadFromEnvFile() {
	err := godotenv.Load()
	if err != nil {
		log.Warn("Error loading .env file")
	}
}
