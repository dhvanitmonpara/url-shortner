package shorten

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"url-shortner/internal/storage"
	"url-shortner/internal/types"
	"url-shortner/internal/utils/response"

	"github.com/go-playground/validator/v10"
	"github.com/segmentio/encoding/json"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating a url")

		var url types.URL

		err := json.NewDecoder(r.Body).Decode(&url)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("request body is empty")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(url); err != nil {

			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastId, err := storage.CreateURL(url.RedirectTO)

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		url, err := storage.GetOriginalURLById(id)

		if err != nil {
			slog.Error("error getting url", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, url)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		url, err := storage.GetURLs()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err.Error())
			return
		}

		response.WriteJson(w, http.StatusOK, url)
	}
}

func RedirectHandler(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		if id == "favicon.ico" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		url, err := storage.GetOriginalURLById(id)
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, err.Error())
			return
		}

		targetURL := url.RedirectTO

		// If it doesn't start with http:// or https://, prepend https://
		if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
			targetURL = "https://" + targetURL
		}

		fmt.Println("redirecting to", targetURL)

		http.Redirect(w, r, targetURL, http.StatusPermanentRedirect)
	}
}

func UpdateUrl(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		var url types.URL

		err := json.NewDecoder(r.Body).Decode(&url)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("request body is empty")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(url); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		updatedUrl, err := storage.UpdateUrl(id, url.RedirectTO)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err.Error())
			return
		}

		response.WriteJson(w, http.StatusOK, updatedUrl)
	}
}

func DeleteUrl(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		if err := storage.DeleteURL(id); err != nil {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(err))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
