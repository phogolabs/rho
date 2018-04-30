package rho

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/goware/errorx"
	uuid "github.com/satori/go.uuid"
)

const (
	// ErrCodeParamRequired is an error code returned when the parameter is missing
	ErrCodeParamRequired = 20101
	// ErrCodeParamInvalid is an error code returned when the parameter's value is an invalid
	ErrCodeParamInvalid = 20102
)

// URLParamUUID returns a request query parameter as UUID
func URLParamUUID(r *http.Request, key string) (uuid.UUID, error) {
	param := chi.URLParam(r, key)

	if param == "" {
		return uuid.Nil, paramRequiredErr(key)
	}

	value, err := uuid.FromString(param)
	if err == nil {
		return value, nil
	}

	err = paramParseErr(key, "UUID", err)
	return uuid.Nil, err
}

// URLParamUUIDOrValue returns a request query parameter as UUID or the
// provided default value if cannot parse the parameter.
func URLParamUUIDOrValue(r *http.Request, key string, value uuid.UUID) uuid.UUID {
	param, err := URLParamUUID(r, key)
	if err != nil {
		param = value
	}

	return param
}

// URLParamUUIDOrNil returns a nil value if cannot parse the UUID parameter
func URLParamUUIDOrNil(r *http.Request, key string) uuid.UUID {
	return URLParamUUIDOrValue(r, key, uuid.Nil)
}

// URLParamInt returns a request query parameter as int64
func URLParamInt(r *http.Request, key string, base, bitSize int) (int64, error) {
	param := chi.URLParam(r, key)

	if param == "" {
		return 0, paramRequiredErr(key)
	}

	value, err := strconv.ParseInt(param, base, bitSize)
	if err == nil {
		return value, nil
	}

	err = paramParseErr(key, "integer number", err)
	return 0, err
}

// URLParamIntOrValue returns a request query parameter as int64 or the
// provided default value if cannot parse the parameter.
func URLParamIntOrValue(r *http.Request, key string, base, bitSize int, value int64) int64 {
	param, err := URLParamInt(r, key, base, bitSize)
	if err != nil {
		param = value
	}

	return param
}

// URLParamUint returns a request query parameter as uint64
func URLParamUint(r *http.Request, key string, base, bitSize int) (uint64, error) {
	param := chi.URLParam(r, key)

	if param == "" {
		return 0, paramRequiredErr(key)
	}

	value, err := strconv.ParseUint(param, base, bitSize)
	if err == nil {
		return value, nil
	}

	err = paramParseErr(key, "unsigned integer number", err)
	return 0, err
}

// URLParamUintOrValue returns a request query parameter as uint64 or the
// provided default value if cannot parse the parameter.
func URLParamUintOrValue(r *http.Request, key string, base, bitSize int, value uint64) uint64 {
	param, err := URLParamUint(r, key, base, bitSize)
	if err != nil {
		param = value
	}

	return param
}

// URLParamFloat returns a request query parameter as float64
func URLParamFloat(r *http.Request, key string, bitSize int) (float64, error) {
	param := chi.URLParam(r, key)

	if param == "" {
		return 0, paramRequiredErr(key)
	}

	value, err := strconv.ParseFloat(param, bitSize)
	if err == nil {
		return value, nil
	}

	err = paramParseErr(key, "float number", err)
	return 0, err
}

// URLParamFloatOrValue returns a request query parameter as float64 or the
// provided default value if cannot parse the parameter.
func URLParamFloatOrValue(r *http.Request, key string, bitSize int, value float64) float64 {
	param, err := URLParamFloat(r, key, bitSize)
	if err != nil {
		param = value
	}

	return param
}

// URLParamTime returns a request query parameter as time.Time
func URLParamTime(r *http.Request, key, format string) (time.Time, error) {
	param := chi.URLParam(r, key)

	if param == "" {
		return time.Time{}, paramRequiredErr(key)
	}

	value, err := time.Parse(format, param)
	if err == nil {
		return value, nil
	}

	info := fmt.Sprintf("Expected date time format '%s'", format)
	err = paramParseErr(key, "date time", err, info)
	return time.Time{}, err
}

// URLParamTimeOrValue returns a request query parameter as time.Time or the
// provided default value if cannot parse the parameter.
func URLParamTimeOrValue(r *http.Request, key, format string, value time.Time) time.Time {
	param, err := URLParamTime(r, key, format)
	if err != nil {
		param = value
	}

	return param
}

func paramRequiredErr(key string) error {
	msg := fmt.Sprintf("Parameter '%s' is required", key)
	err := &ErrorResponse{
		StatusCode: http.StatusBadRequest,
		Err:        errorx.New(ErrCodeParamRequired, msg),
	}
	return err
}

func paramParseErr(key, tname string, err error, details ...string) error {
	info := fmt.Sprintf("Parameter '%s' is not valid %s", key, tname)
	message := append([]string{info}, details...)
	errx := &ErrorResponse{
		StatusCode: http.StatusUnprocessableEntity,
		Err:        errorx.New(ErrCodeParamInvalid, message...),
	}
	errx.Err.Wrap(err)
	return errx
}
