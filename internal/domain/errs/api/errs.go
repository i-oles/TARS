package api

import "errors"

var ErrTaskAlreadyExist = errors.New("task already exists in database")

var ErrTaskNotFound = errors.New("task not found in database")
