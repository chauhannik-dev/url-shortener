package main

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"

	"math/rand"

	"github.com/catinello/base62"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShortenedRequestBody struct {
	URL string `json:"url"`
}

type ShortenedResponseURL struct {
	Key      string `json:"key"`
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

func connectToMongo(connectionString string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func hashUrl(url string, salt string) string {
	hash := sha256.Sum256([]byte(url + salt))
	hashInt := int(binary.BigEndian.Uint64(hash[:8]))

	hashInt = int(math.Abs(float64(hashInt))) // - ve hash results in an empty string

	return base62.Encode(hashInt) // shortened url
}

func validateRequest(r *http.Request) (string, string) {
	var rawMessage map[string]json.RawMessage
	err := json.NewDecoder(r.Body).Decode(&rawMessage)

	if err != nil {
		return "", "failed_to_decode"
	}

	if _, ok := rawMessage["url"]; !ok {
		return "", "url_missing"
	}

	if len(rawMessage) > 1 {
		return "", "only_url_allowed"
	}

	var reqBody ShortenedRequestBody
	err = json.Unmarshal(rawMessage["url"], &reqBody.URL)

	if err != nil || reqBody.URL == "" {
		return "", "invalid_url_format"
	}

	return reqBody.URL, ""
}

func shortenUrl(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	originalUrl, err := validateRequest(r)

	if err != "" {
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	var response ShortenedResponseURL

	filter := bson.M{"long_url": originalUrl}
	var existingRecord bson.M
	findErr := collection.FindOne(context.TODO(), filter).Decode(&existingRecord)

	if findErr == nil {
		response.Key = existingRecord["key"].(string)
		response.ShortURL = existingRecord["short_url"].(string)
		response.LongURL = originalUrl

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return

	} else if findErr != mongo.ErrNoDocuments { // Any other error other than no docs found
		http.Error(w, "Something went wrong: "+findErr.Error(), http.StatusInternalServerError)
		return
	}

	var salt string = ""
	for {
		hashedKey := hashUrl(originalUrl, salt)
		shortURL := fmt.Sprintf("http://short.ly/%s", hashedKey)

		// check if the hash key is already present
		filter := bson.M{"key": hashedKey}
		var result bson.M
		err := collection.FindOne(context.TODO(), filter).Decode(&result)

		if err == mongo.ErrNoDocuments {
			_, err := collection.InsertOne(context.TODO(), bson.M{
				"key":       hashedKey,
				"short_url": shortURL,
				"long_url":  originalUrl,
			})

			if err != nil {
				http.Error(w, "Failed to create short URL: "+err.Error(), http.StatusInternalServerError)
				return
			}

			response.Key = hashedKey
			response.ShortURL = shortURL
			response.LongURL = originalUrl

			break
		} else if err != nil { // If any db related issues occur, we handle it here.
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		salt = fmt.Sprintf("%d", rand.Intn(10000)) // generate a salt for a new key
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func redirectToURL(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	shortURL := r.URL.Path[1:]

	if shortURL == "" {
		http.Error(w, "encoded_url_missing", http.StatusBadRequest)
		return
	}

	var findDocument bson.M
	filter := bson.M{"key": shortURL}
	findErr := collection.FindOne(context.TODO(), filter).Decode(&findDocument)

	if findErr == mongo.ErrNoDocuments {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	} else if findErr != nil {
		http.Error(w, "Database error: "+findErr.Error(), http.StatusInternalServerError)
		return
	}

	redirectionURL := findDocument["long_url"].(string)
	w.Header().Set("Location", redirectionURL)
	w.WriteHeader(http.StatusFound)
}

func deleteURL(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	key := r.URL.Path[1:]

	if key == "" {
		http.Error(w, "short_url_missing", http.StatusBadRequest)
		return
	}

	filter := bson.M{"key": key}
	result, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		http.Error(w, "db_error"+err.Error(), http.StatusNotFound)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "encoded_url_not_found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(result)
}

func main() {
	client, err := connectToMongo("mongodb+srv://cn2blue:cn2blue123@cluster0.80bge.mongodb.net/url_shortener?retryWrites=true&w=majority")

	if err != nil {
		log.Fatal("Couldn't connect to the database")
	}

	defer client.Disconnect(context.TODO())

	collection := client.Database("url_shortener").Collection("urls")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			shortenUrl(w, r, collection)
		case http.MethodGet:
			redirectToURL(w, r, collection)
		case http.MethodDelete:
			deleteURL(w, r, collection)
		default:
			http.Error(w, "Invalid route entered", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server is listening on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Failed to start the server: ", err)
	}
}
