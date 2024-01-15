package api

import (
	"database/sql"
	"errors"
	"fmt"
	db "github.com/StanislavIvanovQA/simple_bank/db/sqlc"
	"github.com/StanislavIvanovQA/simple_bank/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) handlerCreateTransfer(context *gin.Context) {
	var req transferRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAcc, valid := s.validAccount(context, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	authPayload := context.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAcc.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		context.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if _, valid = s.validAccount(context, req.ToAccountID, req.Currency); !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := s.store.TransferTransaction(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, result)
}

func (s *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency missmatch: %s expected but was %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
