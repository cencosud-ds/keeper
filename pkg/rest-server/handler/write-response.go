package handler

import "net/http"

func writeResponse(w http.ResponseWriter, status int, body string) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)

	_, err := w.Write([]byte(body))
	if err != nil {
		return err
	}

	return nil
}
