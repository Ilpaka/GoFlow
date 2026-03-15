package response

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	apperr "goflow/backend/internal/pkg/errors"
)

// ErrorBody is the JSON shape for API errors.
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// SuccessBody is a minimal success envelope; use Data for payloads.
type SuccessBody struct {
	OK   bool `json:"ok"`
	Data any  `json:"data,omitempty"`
}

// WriteJSON writes JSON with the given status and payload (no envelope).
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// WriteSuccess writes { "ok": true, "data": data } with HTTP 200 unless status overrides.
//
// Example:
//
//	response.WriteSuccess(w, http.StatusOK, map[string]any{"user": dto})
func WriteSuccess(w http.ResponseWriter, status int, data any) {
	if status == 0 {
		status = http.StatusOK
	}
	WriteJSON(w, status, SuccessBody{OK: true, Data: data})
}

// WriteError maps err to JSON and HTTP status. Unknown errors become internal.
//
// Example:
//
//	if err := svc.Login(...); err != nil {
//	    response.WriteError(w, r, log, err)
//	    return
//	}
func WriteError(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error) {
	if err == nil {
		WriteJSON(w, http.StatusOK, SuccessBody{OK: true})
		return
	}

	status := apperr.HTTPStatus(err)
	body := ErrorBody{Code: "internal", Message: "internal server error"}

	if ae, ok := apperr.As(err); ok {
		body.Code = string(ae.Kind)
		body.Message = ae.Message
	} else {
		if log != nil {
			log.Error("unhandled error", "err", err, "path", r.URL.Path)
		}
	}

	WriteJSON(w, status, map[string]any{
		"ok":    false,
		"error": body,
	})
}

// WriteErrorWithDetails is like WriteError but attaches optional details (validation fields, etc.).
func WriteErrorWithDetails(w http.ResponseWriter, r *http.Request, log *slog.Logger, err error, details any) {
	if err == nil {
		WriteSuccess(w, http.StatusOK, nil)
		return
	}
	status := apperr.HTTPStatus(err)
	body := ErrorBody{Code: "internal", Message: "internal server error", Details: details}

	if ae, ok := apperr.As(err); ok {
		body.Code = string(ae.Kind)
		body.Message = ae.Message
	} else {
		if log != nil {
			log.Error("unhandled error", "err", err, "path", r.URL.Path)
		}
	}

	WriteJSON(w, status, map[string]any{
		"ok":    false,
		"error": body,
	})
}

// Is reports whether err unwraps to an apperr.Error of the given kind.
func Is(err error, kind apperr.Kind) bool {
	var ae *apperr.Error
	if !errors.As(err, &ae) {
		return false
	}
	return ae.Kind == kind
}
