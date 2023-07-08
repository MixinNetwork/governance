package session

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

type Error struct {
	Status      int    `json:"status"`
	Code        int    `json:"code"`
	Description string `json:"description"`
	Extra       any    `json:"extra,omitempty"`
	trace       string
}

func (sessionError Error) Error() string {
	str, err := json.Marshal(sessionError)
	if err != nil {
		log.Panicln(err)
	}
	return string(str)
}

func ParseError(err string) (Error, bool) {
	var sessionErr Error
	json.Unmarshal([]byte(err), &sessionErr)
	return sessionErr, sessionErr.Code > 0 && sessionErr.Description != ""
}

func BadRequestError(ctx context.Context) *Error {
	description := "The request body can’t be pasred as valid data."
	return createError(ctx, http.StatusAccepted, http.StatusBadRequest, description, nil)
}

func ServerError(ctx context.Context, err error) *Error {
	description := http.StatusText(http.StatusInternalServerError)
	return createError(ctx, http.StatusInternalServerError, http.StatusInternalServerError, description, err)
}

func NotFoundError(ctx context.Context) *Error {
	description := "The endpoint is not found."
	return createError(ctx, http.StatusAccepted, http.StatusNotFound, description, nil)
}

func AuthorizationError(ctx context.Context) *Error {
	description := "Unauthorized, maybe invalid token."
	return createError(ctx, http.StatusAccepted, 401, description, nil)
}

func ForbiddenError(ctx context.Context) *Error {
	description := http.StatusText(http.StatusForbidden)
	return createError(ctx, http.StatusAccepted, http.StatusForbidden, description, nil)
}

func TransactionError(ctx context.Context, err error) *Error {
	description := http.StatusText(http.StatusInternalServerError)
	return createError(ctx, http.StatusInternalServerError, 10001, description, err)
}

func BadDataError(ctx context.Context) *Error {
	description := "The request data has invalid field."
	return createError(ctx, http.StatusAccepted, 10002, description, nil)
}

func BadDataErrorWithFieldAndData(ctx context.Context, field, reason, data string) *Error {
	description := "The request data has invalid field."
	er := fmt.Errorf("[BAD DATA %s]", data)
	err := createError(ctx, http.StatusAccepted, 10002, description, er)
	err.Extra = map[string]string{
		"field":  field,
		"reason": reason,
	}
	return err
}

func TooManyPendingTransactions(ctx context.Context) *Error {
	return createError(ctx, http.StatusAccepted, 10003, "Too many pending transactions.", nil)
}

func InsufficientAccountError(ctx context.Context) *Error {
	description := "Insufficient account quotas."
	return createError(ctx, http.StatusAccepted, 10301, description, nil)
}

func InsufficientTransactionError(ctx context.Context) *Error {
	description := "Insufficient transaction quotas."
	return createError(ctx, http.StatusAccepted, 10302, description, nil)
}

func InsufficientAmountError(ctx context.Context) *Error {
	description := "Insufficient amount."
	return createError(ctx, http.StatusAccepted, 10303, description, nil)
}

func InsufficientFeeError(ctx context.Context) *Error {
	description := "Insufficient fee."
	return createError(ctx, http.StatusAccepted, 10304, description, nil)
}

func createError(ctx context.Context, status, code int, description string, err error) *Error {
	if serr, ok := err.(*Error); ok {
		return serr
	}
	pc, file, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	trace := fmt.Sprintf("Function: %s\n[ERROR %d] %s\n%s:%d", funcName, code, description, file, line)
	if err != nil {
		if sessionError, ok := err.(Error); ok {
			trace = trace + "\n" + sessionError.trace
		} else {
			trace = trace + "\n" + err.Error()
		}
	}
	return &Error{
		Status:      status,
		Code:        code,
		Description: description,
		trace:       trace,
	}
}
