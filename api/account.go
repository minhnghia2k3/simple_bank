package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/minhnghia2k3/simple_bank/db/sqlc"
	"net/http"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

// createAccount service will validate input from JSON body
// then call CreateAccount from repository layer to store a new user.
func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest

	// Validate context JSON body.
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	/// Get arg from user's request
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	// Call CreateAccount() from repository layer to create a new account.
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest

	// BindURI and assign value to ID
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Query to get account from db.
	account, err := server.store.GetAccount(ctx, req.ID)
	fmt.Println("err: ", err)
	if err != nil {
		// 404 not found
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		// 500 server error
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Response json.
	ctx.JSON(http.StatusOK, account)
}

type listAccountQuery struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type Info struct {
	Page int32 `json:"page"`
	Unit int   `json:"unit"`
}

type responseBody struct {
	Info Info `json:"info"`
	Data any  `json:"data"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var query listAccountQuery

	// Validate query parameters
	err := ctx.ShouldBindQuery(&query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Offset: (query.PageID - 1) * query.PageSize,
		Limit:  query.PageSize,
	}

	// Query to get list accounts
	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	responseBody := responseBody{
		Info: Info{Page: query.PageID, Unit: len(accounts)},
		Data: accounts,
	}

	ctx.JSON(http.StatusOK, responseBody)
}
