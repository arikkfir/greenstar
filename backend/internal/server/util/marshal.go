package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"net/http"
)

const (
	QueryNilValue = "<nil>"
)

func Marshal(w http.ResponseWriter, r *http.Request, statusCode int, v any) error {
	acceptedMediaTypes := parseAcceptHeader(r.Header.Get("Accept"))

	var contentType string
	for _, mediaType := range acceptedMediaTypes {
		if mediaType.Accepts("application/json") {
			contentType = "application/json"
			break
		} else if mediaType.Accepts("application/yaml") {
			contentType = "application/yaml"
			break
		}
	}
	if contentType == "" {
		return ErrNotAcceptable
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)

	switch contentType {
	case "application/json":
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ") // TODO: should depend on dev-mode
		if err := encoder.Encode(v); err != nil {
			return err
		}
	case "application/yaml":
		encoder := yaml.NewEncoder(w)
		defer encoder.Close()
		encoder.SetIndent(2) // TODO: should depend on dev-mode
		if err := encoder.Encode(v); err != nil {
			return err
		}
	}
	return nil
}

func UnmarshalBody(r *http.Request, v any) error {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return fmt.Errorf("%w: empty Content-Type header", ErrUnsupportedMediaType)
	}

	mt, err := parseMediaType(contentType)
	if err != nil {
		return errors.Join(err, ErrUnsupportedMediaType)
	}

	switch mt.Type + "/" + mt.Subtype {
	case "application/json":
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(v); err != nil {
			return errors.Join(fmt.Errorf("failed to unmarshal JSON request: %w", err), ErrInternalError)
		}
		return nil
	case "application/yaml":
		decoder := yaml.NewDecoder(r.Body)
		decoder.KnownFields(false)
		if err := decoder.Decode(v); err != nil {
			return errors.Join(fmt.Errorf("failed to unmarshal YAML request: %w", err), ErrInternalError)
		}
		return nil
	default:
		return ErrUnsupportedMediaType
	}
}
