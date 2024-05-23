package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ashishjuyal/banking-lib/errs"
	"github.com/ashishjuyal/banking/dto"
	"github.com/ashishjuyal/banking/mocks/service"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

const dbTSLayout = "2006-01-02 15:04:05"

var ah *AccountHandler
var mockActSrv *service.MockAccountService

func setUp(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockActSrv = service.NewMockAccountService(ctrl)
	ah = &AccountHandler{
		service: mockActSrv,
	}
	return func() {
		defer ctrl.Finish()
		mockActSrv = nil
		ah = nil
	}

}

func Test_new_account_handler_with_DB_internal_error(t *testing.T) {
	setUp(t)

	actReq := dto.NewAccountRequest{
		AccountType: "Saving",
		Amount:      100,
	}
	rsp := dto.NewAccountResponse{
		AccountId: "5555",
	}
	appErr := errs.AppError{
		Code:    http.StatusInternalServerError,
		Message: "Failed to connect to DB",
	}
	mockActSrv.EXPECT().NewAccount(actReq).Return(&rsp, &appErr)

	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	if data, err := json.Marshal(actReq); err != nil {
		t.Error("Failed to Marshal Data, data test error")
	} else {
		router.HandleFunc("/account", ah.NewAccount)
		req, _ := http.NewRequest(http.MethodGet, "/account", bytes.NewBuffer(data))
		router.ServeHTTP(recorder, req)
	}

	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed to get StatusInternalServerError when creating account")
	}

}

func Test_new_account_handler_with_Success(t *testing.T) {
	setUp(t)

	actReq := dto.NewAccountRequest{
		AccountType: "Checking",
		Amount:      200,
	}
	rsp := dto.NewAccountResponse{
		AccountId: "5555",
	}
	mockActSrv.EXPECT().NewAccount(actReq).Return(&rsp, nil)

	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	if data, err := json.Marshal(actReq); err != nil {
		t.Error("Failed to Marshal Data, data test error")
	} else {
		router.HandleFunc("/account", ah.NewAccount)
		req, _ := http.NewRequest(http.MethodGet, "/account", bytes.NewBuffer(data))
		router.ServeHTTP(recorder, req)
	}

	if recorder.Code != http.StatusCreated {
		t.Error("Failed to get Success when creating account")
	}

}

func Test_make_transaction_handler_with_Success(t *testing.T) {
	setUp(t)

	txnReq := dto.TransactionRequest{
		Amount:          100,
		TransactionType: "withdrawal",
		TransactionDate: time.Now().Format(dbTSLayout),
	}
	rsp := dto.TransactionResponse{
		AccountId:       "90720",
		Amount:          100,
		TransactionType: "withdrawal",
		TransactionDate: time.Now().Format(dbTSLayout),
		TransactionId:   "234",
	}

	mockActSrv.EXPECT().MakeTransaction(txnReq).Return(&rsp, nil)

	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	if data, err := json.Marshal(txnReq); err != nil {
		t.Error("Failed to Marshal Data, data test error")
	} else {
		router.HandleFunc("/customers", ah.MakeTransaction)
		req, _ := http.NewRequest(http.MethodGet, "/customers", bytes.NewBuffer(data))
		router.ServeHTTP(recorder, req)
	}

	if recorder.Code != http.StatusOK {
		t.Error("Failed to get Success when making transaction")
	}

}

func Test_make_transaction_handler_with_Failure(t *testing.T) {
	setUp(t)

	txnReq := dto.TransactionRequest{
		Amount:          100,
		TransactionType: "withdrawal",
		TransactionDate: time.Now().Format(dbTSLayout),
	}
	rsp := dto.TransactionResponse{
		AccountId:       "90720",
		Amount:          100,
		TransactionType: "withdrawal",
		TransactionDate: time.Now().Format(dbTSLayout),
		TransactionId:   "234",
	}
	appErr := errs.AppError{
		Code:    http.StatusInternalServerError,
		Message: "Failed to connect to DB",
	}
	mockActSrv.EXPECT().MakeTransaction(txnReq).Return(&rsp, &appErr)

	recorder := httptest.NewRecorder()
	router := mux.NewRouter()

	if data, err := json.Marshal(txnReq); err != nil {
		t.Error("Failed to Marshal Data, data test error")
	} else {
		router.HandleFunc("/customers", ah.MakeTransaction)
		req, _ := http.NewRequest(http.MethodGet, "/customers", bytes.NewBuffer(data))
		router.ServeHTTP(recorder, req)
	}

	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed to get StatusInternalServerError when making transaction")
	}

}
