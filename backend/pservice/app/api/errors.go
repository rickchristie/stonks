package api

type ApiErr string

const (
	ApiErrNone          ApiErr = ""
	ApiErrValidation    ApiErr = "Validation"
	ApiErrNotFound      ApiErr = "NotFound"
	ApiErrInternalError ApiErr = "InternalError"
)
