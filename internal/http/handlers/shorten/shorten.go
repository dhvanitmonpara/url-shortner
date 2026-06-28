package shorten

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
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
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusOK, url)
	}
}

func RedirectHandler(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		url, err := storage.GetOriginalURLById(id)
		if err != nil {
			response.WriteJson(w, http.StatusNotFound, err)
			return
		}

		http.Redirect(w, r, url.RedirectTO, http.StatusFound)
	}
}
