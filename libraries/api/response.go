package api

import (
	"encoding/json"
	"net/http"
)

// ResponseFormat is used to pass an response in standard format
type ResponseFormat struct {
	StatusCode string      `json:"status_code"`
	Message    string      `json:"status_message"`
	Data       interface{} `json:"data"`
}

// Response converts a Go value to JSON and sends it to the client.
func Response(w http.ResponseWriter, data interface{}, statusCode string, message string, httpCode int) error {

	// Convert the response value to JSON.
	res, err := json.Marshal(ResponseFormat{StatusCode: statusCode, Message: message, Data: data})
	if err != nil {
		return err
	}

	// Respond with the provided JSON.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpCode)
	if _, err := w.Write(res); err != nil {
		return err
	}

	return nil
}

// ResponseOK converts a Go value to JSON and sends it to the client.
func ResponseOK(w http.ResponseWriter, data interface{}, HTTPStatus int) error {
	return Response(w, data, StatusCodeOK, StatusMessageOK, HTTPStatus)
}

// ResponseError sends an error reponse back to the client.
func ResponseError(w http.ResponseWriter, err error) error {

	// If the error was of the type *Error, the handler has
	// a specific status code and error to return.
	if webErr, ok := err.(*Error); ok {
		if err := Response(w, nil, webErr.Status, webErr.MessageStatus, webErr.HTTPStatus); err != nil {
			return err
		}
		return nil
	}

	// If not, the handler sent any arbitrary error value so use 500.
	if err := Response(w, nil, StatusCodeInternalServerError, StatusMessageInternalServerError, http.StatusInternalServerError); err != nil {
		return err
	}
	return nil
}
