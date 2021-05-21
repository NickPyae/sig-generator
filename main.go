package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sig-generator/models"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/open-horizon/rsapss-tool/sign"
)

const HZN_KEY_FILE_ERROR_MSG = "HZN_KEY_FILE environment is not exported"
const DEPLOYMENT_IMAGE_ERROR_MSG = "Deployment image must be provided"
const SIGN_ERROR_MSG = "Error in signing deployment string with private key file"
const SERVER_ADDR = "127.0.0.1:8080"

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/encrypt", EncryptHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         SERVER_ADDR,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Server listening at ", SERVER_ADDR)
	log.Fatal(srv.ListenAndServe())
}

func EncryptHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /encrypt")

	var deployment models.Deployment

	data, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(data, &deployment)

	if err != nil {
		errorMessage := err.Error()
		log.Println("Error in request body: ", errorMessage)

		w.WriteHeader(http.StatusBadRequest)

		response := models.ErrorResponseModel{Code: http.StatusBadRequest, Error: errorMessage}
		jsonErr := json.NewEncoder(w).Encode(response)
		if jsonErr != nil {
			log.Println(fmt.Sprintf("unable to encode response, %v", jsonErr))
		}

		return
	}

	log.Println("POST body ", string(data))

	if deployment.Services.Location.Image == "" {
		log.Println(DEPLOYMENT_IMAGE_ERROR_MSG)

		w.WriteHeader(http.StatusBadRequest)

		response := models.ErrorResponseModel{Code: http.StatusBadRequest, Error: DEPLOYMENT_IMAGE_ERROR_MSG}
		jsonErr := json.NewEncoder(w).Encode(response)
		if jsonErr != nil {
			log.Println(fmt.Sprintf("unable to encode response, %v", jsonErr))
		}

		return
	}

	hznKeyFile := strings.TrimSpace(os.Getenv("HZN_KEY_FILE"))

	if hznKeyFile == "" {
		log.Println(HZN_KEY_FILE_ERROR_MSG)

		w.WriteHeader(http.StatusInternalServerError)

		errorMessage := models.ErrorResponseModel{Code: http.StatusInternalServerError, Error: HZN_KEY_FILE_ERROR_MSG}
		jsonErr := json.NewEncoder(w).Encode(errorMessage)
		if jsonErr != nil {
			log.Println(fmt.Sprintf("unable to encode response, %v", jsonErr))
		}

		return
	}

	sig, err := sign.Input(hznKeyFile, data)
	if err != nil {
		errorMessage := err.Error()
		log.Printf("%s: %s", SIGN_ERROR_MSG, errorMessage)

		w.WriteHeader(http.StatusInternalServerError)

		response := models.ErrorResponseModel{Code: http.StatusInternalServerError, Error: fmt.Sprintf("%s %s", SIGN_ERROR_MSG, errorMessage)}
		jsonErr := json.NewEncoder(w).Encode(response)
		if jsonErr != nil {
			log.Println(fmt.Sprintf("unable to encode response, %v", jsonErr))
		}

		return
	}

	log.Printf("Deployment signature: %s ", sig)

	successMessage := models.ResponseModel{DeploymentSignature: sig}
	w.WriteHeader(http.StatusOK)

	jsonErr := json.NewEncoder(w).Encode(successMessage)
	if jsonErr != nil {
		log.Println(fmt.Sprintf("unable to encode response, %v", jsonErr))
	}
}
