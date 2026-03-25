// handler/errors.go
package handlers

import (
    "errors"
    "net/http"

    model "github.com/Loboo34/travel/models"
    "github.com/Loboo34/travel/utils"
    "go.uber.org/zap"
)

func handleServiceError(w http.ResponseWriter, err error, context string) {
    var validationErr *model.ValidationError
    if errors.As(err, &validationErr) {
        utils.RespondWithError(w, http.StatusBadRequest, validationErr.Error())
        return
    }

    var notFoundErr *model.NotFoundError
    if errors.As(err, &notFoundErr) {
        utils.RespondWithError(w, http.StatusConflict, notFoundErr.Error())
        return
    }

    var paymentErr *model.PaymentError
    if errors.As(err, &paymentErr) {
        utils.RespondWithError(w, http.StatusPaymentRequired, paymentErr.Error())
        return
    }

    // unexpected error — log with context
    utils.Logger.Error(context,
        zap.Error(err),
    )
    utils.RespondWithError(w, http.StatusInternalServerError, "an unexpected error occurred")
}