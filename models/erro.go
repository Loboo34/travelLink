// model/errors.go
package model

import "fmt"


type ValidationError struct {
    Message string
}

func (e *ValidationError) Error() string {
    return e.Message
}


type NotFoundError struct {
    Resource string 
    ID       string
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s %q not found", e.Resource, e.ID)
}


type PaymentError struct {
    Message string
}

func (e *PaymentError) Error() string {
    return e.Message
}


type AuthError struct {
    Message string
}

func (e *AuthError) Error() string {
    return e.Message
}

type ConflictError struct {
    Message string
}

func (e *ConflictError) Error() string {
    return e.Message
}