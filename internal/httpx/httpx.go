package httpx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, ErrorResponse{Error: err.Error()})
}

func WriteErrorMessage(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, ErrorResponse{Error: message})
}

func ReadJSON[T any](r *http.Request) (T, error) {
	var data T
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	return data, err
}

func HandleList[T any](serviceFunc func(context.Context) ([]T, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		results, err := serviceFunc(r.Context())
		if err != nil {
			WriteError(w, http.StatusInternalServerError, err)
			return
		}
		WriteJSON(w, http.StatusOK, results)
	}
}

func HandleGetByID[T any](serviceFunc func(context.Context, int32) (T, error), resourceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			ID int32 `json:"id"`
		}

		req, err := ReadJSON[request](r)
		if err != nil {
			WriteError(w, http.StatusBadRequest, err)
			return
		}

		result, err := serviceFunc(r.Context(), req.ID)
		if err != nil {
			WriteErrorMessage(w, http.StatusInternalServerError, fmt.Sprintf("Error getting %s: %v", resourceName, err))
			return
		}
		WriteJSON(w, http.StatusOK, result)
	}
}

func HandleGetByName[T any](serviceFunc func(context.Context, string) ([]T, error), resourceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Name string `json:"name"`
		}

		req, err := ReadJSON[request](r)
		if err != nil {
			WriteError(w, http.StatusBadRequest, err)
			return
		}

		results, err := serviceFunc(r.Context(), req.Name)
		if err != nil {
			WriteErrorMessage(w, http.StatusInternalServerError, fmt.Sprintf("Error getting %s by name: %v", resourceName, err))
			return
		}
		WriteJSON(w, http.StatusOK, results)
	}
}

func HandleCreate[Req any, Res any](serviceFunc func(context.Context, Req) (Res, error), resourceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := ReadJSON[Req](r)
		if err != nil {
			WriteError(w, http.StatusBadRequest, err)
			return
		}

		result, err := serviceFunc(r.Context(), req)
		if err != nil {
			WriteErrorMessage(w, http.StatusInternalServerError, fmt.Sprintf("Error creating %s: %v", resourceName, err))
			return
		}
		WriteJSON(w, http.StatusCreated, result)
	}
}

func HandleDelete[T any](serviceFunc func(context.Context, int32) error, resourceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			ID int32 `json:"id"`
		}

		req, err := ReadJSON[request](r)
		if err != nil {
			WriteError(w, http.StatusBadRequest, err)
			return
		}

		err = serviceFunc(r.Context(), req.ID)
		if err != nil {
			WriteErrorMessage(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting %s: %v", resourceName, err))
			return
		}
		WriteJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("%s deleted successfully", resourceName)})
	}
}
