package web

import (
	"context"
	"encoding/json"
	"net/http"
)

func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {
	return writeResponse(ctx, w, data, statusCode)
}

func RespondError(ctx context.Context, w http.ResponseWriter, err error, statusCode int) error {
	res := struct {
		Error string `json:"error"`
	}{
		err.Error(),
	}

	return writeResponse(ctx, w, res, statusCode)
}

func writeResponse(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
