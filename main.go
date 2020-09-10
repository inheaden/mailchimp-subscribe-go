package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/golang/gddo/httputil/header"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	APIKey   string `required:"true" split_words:"true"`
	Audience string `required:"true" split_words:"true"`
}

type server struct{}

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

	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := `"Content-Type header is not application/json"`
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}

	api := NewMcAPI()

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
	if err := loadFromEnvFile(); err != nil {
		log.Fatalln(err)
	}

	var config config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalln(err)
	}

	s := &server{}
	http.Handle("/", s)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func loadFromEnvFile() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return err
}
