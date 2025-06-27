// Package client_test contains unit tests for the Mpesa client.
package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/freelancer254/mpesa-go/client"
	"github.com/freelancer254/mpesa-go/types"
)

// mockServer creates a test HTTP server with a custom handler.
func mockServer(t *testing.T, statusCode int, response interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if response != nil {
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}
	}))
}

// TestNewMpesa tests the initialization of the Mpesa client.
func TestNewMpesa(t *testing.T) {
	mpesa := client.NewMpesa()
	if mpesa == nil {
		t.Fatal("NewMpesa returned nil")
	}
	if mpesa.BaseURL() != "https://api.safaricom.co.ke" {
		t.Errorf("expected baseURL to be %s, got %s", "https://api.safaricom.co.ke", mpesa.BaseURL())
	}
	if len(mpesa.Headers()) != 0 {
		t.Errorf("expected headers to be empty, got %v", mpesa.Headers())
	}
}

// TestGetAccessToken_Success tests the GetAccessToken method with a successful response.
func TestGetAccessToken_Success(t *testing.T) {
	ctx := context.Background()
	response := types.AccessTokenResponse{
		AccessToken: "test-token",
		ExpiresIn:   "3600",
	}
	server := mockServer(t, http.StatusOK, response)
	defer server.Close()

	mpesa := client.NewMpesa()
	mpesa.SetBaseURL(server.URL)

	token, err := mpesa.GetAccessToken(ctx, "consumer_key", "consumer_secret")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token.AccessToken != response.AccessToken {
		t.Errorf("expected access token %s, got %s", response.AccessToken, token.AccessToken)
	}
	if token.ExpiresIn != response.ExpiresIn {
		t.Errorf("expected expires in %s, got %s", response.ExpiresIn, token.ExpiresIn)
	}
}

// TestGetAccessToken_Error tests the GetAccessToken method with an error response.
func TestGetAccessToken_Error(t *testing.T) {
	ctx := context.Background()
	server := mockServer(t, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	defer server.Close()

	mpesa := client.NewMpesa()
	mpesa.SetBaseURL(server.URL)

	err, _ := mpesa.GetAccessToken(ctx, "invalid_key", "invalid_secret")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestSTKPush_Success tests the STKPush method with a successful response.
func TestSTKPush_Success(t *testing.T) {
	ctx := context.Background()
	response := types.STKPushResponse{
		MerchantRequestID:   "29115-34620561-1",
		CheckoutRequestID:   "ws_CO_191220191020363925",
		ResponseCode:        "0",
		ResponseDescription: "Success",
		CustomerMessage:     "Accepted",
	}
	server := mockServer(t, http.StatusOK, response)
	defer server.Close()

	mpesa := client.NewMpesa()
	mpesa.SetBaseURL(server.URL)

	payload := types.STKPushRequest{
		AccessToken:       "test-token",
		BusinessShortCode: "123456",
		Password:          "encoded_password",
		Amount:            "100",
		PartyA:            "254700000000",
		PartyB:            "123456",
		PhoneNumber:       "254700000000",
		CallBackURL:       "https://callback.example.com",
		AccountReference:  "Test123",
		TransactionDesc:   "Payment",
	}

	result, err := mpesa.STKPush(ctx, payload)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ResponseCode != response.ResponseCode {
		t.Errorf("expected response code %s, got %s", response.ResponseCode, result.ResponseCode)
	}
	if result.CheckoutRequestID != response.CheckoutRequestID {
		t.Errorf("expected checkout request ID %s, got %s", response.CheckoutRequestID, result.CheckoutRequestID)
	}
}

// TestSTKPush_ValidationError tests the STKPush method with invalid payload.
func TestSTKPush_ValidationError(t *testing.T) {
	ctx := context.Background()
	mpesa := client.NewMpesa()

	payload := types.STKPushRequest{
		AccessToken: "test-token",
		// Missing required fields
	}

	_, err := mpesa.STKPush(ctx, payload)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if err.Error() == "" {
		t.Errorf("expected validation error, got %v", err)
	}
}

// TestRegisterURL_Success tests the RegisterURL method with a successful response.
func TestRegisterURL_Success(t *testing.T) {
	ctx := context.Background()
	response := types.RegisterURLResponse{
		ResponseDescription: "Success",
	}
	server := mockServer(t, http.StatusOK, response)
	defer server.Close()

	mpesa := client.NewMpesa()
	mpesa.SetBaseURL(server.URL)

	payload := types.RegisterURLRequest{
		AccessToken:     "test-token",
		ShortCode:       "123456",
		ResponseType:    "Completed",
		ConfirmationURL: "https://confirm.example.com",
		ValidationURL:   "https://validate.example.com",
	}

	result, err := mpesa.RegisterURL(ctx, payload)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ResponseDescription != response.ResponseDescription {
		t.Errorf("expected response description %s, got %s", response.ResponseDescription, result.ResponseDescription)
	}
}

// TestSimulateTransaction_Success tests the SimulateTransaction method with a successful response.
func TestSimulateTransaction_Success(t *testing.T) {
	ctx := context.Background()
	response := types.SimulateTransactionResponse{
		ConversationID:           "AG_20180324_000066530b914eee3f85",
		OriginatorConversationID: "25344-885903-1",
		ResponseDescription:      "Accept the service request successfully.",
	}
	server := mockServer(t, http.StatusOK, response)
	defer server.Close()

	mpesa := client.NewMpesa()
	mpesa.SetBaseURL(server.URL)

	payload := types.SimulateTransactionRequest{
		AccessToken:   "test-token",
		ShortCode:     "123456",
		Amount:        "100",
		Msisdn:        "254700000000",
		BillRefNumber: "TEST123",
	}

	result, err := mpesa.SimulateTransaction(ctx, payload)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.OriginatorConversationID != response.OriginatorConversationID {
		t.Errorf("expected originator conversation ID %s, got %s", response.OriginatorConversationID, result.OriginatorConversationID)
	}
}

// TestQueryTransaction_Success tests the QueryTransaction method with a successful response.
func TestQueryTransaction_Success(t *testing.T) {
	ctx := context.Background()
	response := types.QueryTransactionResponse{
		ConversationID:           "AG_20180324_000066530b914eee3f85",
		OriginatorConversationID: "25344-885903-1",
		ResponseDescription:      "Accept the service request successfully.",
	}
	server := mockServer(t, http.StatusOK, response)
	defer server.Close()

	mpesa := client.NewMpesa()
	mpesa.SetBaseURL(server.URL)

	payload := types.QueryTransactionRequest{
		AccessToken:        "test-token",
		Initiator:          "test-initiator",
		SecurityCredential: "credential",
		TransactionID:      "TX123",
		PartyA:             "123456",
		IdentifierType:     "4",
		ResultURL:          "https://result.example.com",
		QueueTimeOutURL:    "https://timeout.example.com",
		Remarks:            "Test query",
		Occasion:           "Test",
	}

	result, err := mpesa.QueryTransaction(ctx, payload)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.OriginatorConversationID != response.OriginatorConversationID {
		t.Errorf("expected originator conversation ID %s, got %s", response.OriginatorConversationID, result.OriginatorConversationID)
	}
}

// TestGetBalance_Success tests the GetBalance method with a successful response.
func TestGetBalance_Success(t *testing.T) {
	ctx := context.Background()
	response := types.GetBalanceResponse{
		ConversationID:           "AG_20180324_000066530b914eee3f85",
		OriginatorConversationID: "25344-885903-1",
		ResponseDescription:      "Accept the service request successfully.",
	}
	server := mockServer(t, http.StatusOK, response)
	defer server.Close()

	mpesa := client.NewMpesa()
	mpesa.SetBaseURL(server.URL)

	payload := types.GetBalanceRequest{
		AccessToken:        "test-token",
		Initiator:          "test-initiator",
		SecurityCredential: "credential",
		PartyA:             "123456",
		IdentifierType:     "4",
		Remarks:            "Test balance",
		QueueTimeOutURL:    "https://timeout.example.com",
		ResultURL:          "https://result.example.com",
	}

	result, err := mpesa.GetBalance(ctx, payload)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.OriginatorConversationID != response.OriginatorConversationID {
		t.Errorf("expected originator conversation ID %s, got %s", response.OriginatorConversationID, result.OriginatorConversationID)
	}
}

// TestB2CSend_Success tests the B2CSend method with a successful response.
func TestB2CSend_Success(t *testing.T) {
	ctx := context.Background()
	response := types.B2CSendResponse{
		ConversationID:           "AG_20180324_000066530b914eee3f85",
		OriginatorConversationID: "25344-885903-1",
		ResponseDescription:      "Accept the service request successfully.",
	}
	server := mockServer(t, http.StatusOK, response)
	defer server.Close()

	mpesa := client.NewMpesa()
	mpesa.SetBaseURL(server.URL)

	payload := types.B2CSendRequest{
		AccessToken:        "test-token",
		InitiatorName:      "test-initiator",
		SecurityCredential: "credential",
		CommandID:          "PromotionPayment",
		Amount:             "100",
		PartyA:             "123456",
		PartyB:             "254700000000",
		Remarks:            "Test B2C",
		QueueTimeOutURL:    "https://timeout.example.com",
		ResultURL:          "https://result.example.com",
		Occasion:           "Test",
	}

	result, err := mpesa.B2CSend(ctx, payload)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.OriginatorConversationID != response.OriginatorConversationID {
		t.Errorf("expected originator conversation ID %s, got %s", response.OriginatorConversationID, result.OriginatorConversationID)
	}
}

// TestB2BSend_Success tests the B2BSend method with a successful response.
func TestB2BSend_Success(t *testing.T) {
	ctx := context.Background()
	response := types.B2BSendResponse{
		ConversationID:           "AG_20180324_000066530b914eee3f85",
		OriginatorConversationID: "25344-885903-1",
		ResponseCode:             "0",
		ResponseDescription:      "Accept the service request successfully.",
	}
	server := mockServer(t, http.StatusOK, response)
	defer server.Close()

	mpesa := client.NewMpesa()
	mpesa.SetBaseURL(server.URL)

	payload := types.B2BSendRequest{
		AccessToken:            "test-token",
		Initiator:              "test-initiator",
		SecurityCredential:     "credential",
		CommandID:              "BusinessPayment",
		SenderIdentifierType:   "4",
		ReceiverIdentifierType: "4",
		Amount:                 "100",
		PartyA:                 "123456",
		PartyB:                 "654321",
		Remarks:                "Test B2B",
		AccountReference:       "TEST123",
		Requester:              "254700000000",
		QueueTimeOutURL:        "https://timeout.example.com",
		ResultURL:              "https://result.example.com",
	}

	result, err := mpesa.B2BSend(ctx, payload)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ResponseCode != response.ResponseCode {
		t.Errorf("expected response code %s, got %s", response.ResponseCode, result.ResponseCode)
	}
}

// TestRegisterPullAPI_Success tests the RegisterPullAPI method with a successful response.
func TestRegisterPullAPI_Success(t *testing.T) {
	ctx := context.Background()
	response := types.RegisterPullAPIResponse{
		ResponseRefID:       "18633-7271215-1",
		ResponseStatus:      "1001",
		ShortCode:           "600000",
		ResponseDescription: "ShortCode already Registered",
	}
	server := mockServer(t, http.StatusOK, response)
	defer server.Close()

	mpesa := client.NewMpesa()
	mpesa.SetBaseURL(server.URL)

	payload := types.RegisterPullAPIRequest{
		AccessToken:     "test-token",
		ShortCode:       "600000",
		NominatedNumber: "254700000000",
		CallBackURL:     "https://callback.example.com",
	}

	result, err := mpesa.RegisterPullAPI(ctx, payload)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ResponseRefID != response.ResponseRefID {
		t.Errorf("expected response ref ID %s, got %s", response.ResponseRefID, result.ResponseRefID)
	}
}

// TestPullTransactions_Success tests the PullTransactions method with a successful response.
func TestPullTransactions_Success(t *testing.T) {
	ctx := context.Background()
	response := types.PullTransactionsResponse{
		ResponseRefID:   "26178-42530161-2",
		ResponseCode:    "1000",
		ResponseMessage: "Success",
		Transactions: []types.Transaction{
			{
				TransactionID:    "yzlyrEsRG1",
				TrxDate:          "2020-08-05T10:13:00Z",
				Msisdn:           72200000,
				Sender:           "UAT2",
				TransactionType:  "c2b-pay-bill-debit",
				BillReference:    "37207636392",
				Amount:           "168.00",
				OrganizationName: "Daraja Pull API Test",
			},
		},
	}
	server := mockServer(t, http.StatusOK, response)
	defer server.Close()

	mpesa := client.NewMpesa()
	mpesa.SetBaseURL(server.URL)

	payload := types.PullTransactionsRequest{
		AccessToken: "test-token",
		ShortCode:   "600000",
		StartDate:   "2020-08-01",
		EndDate:     "2020-08-10",
		OffSetValue: "0",
	}

	result, err := mpesa.PullTransactions(ctx, payload)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Transactions) != 1 {
		t.Errorf("expected 1 transaction, got %d", len(result.Transactions))
	}
	if result.Transactions[0].TransactionID != response.Transactions[0].TransactionID {
		t.Errorf("expected transaction ID %s, got %s", response.Transactions[0].TransactionID, result.Transactions[0].TransactionID)
	}
}

// TestConcurrentAccess tests concurrent access to the Mpesa client.
func TestConcurrentAccess(t *testing.T) {
	ctx := context.Background()
	mpesa := client.NewMpesa()
	server := mockServer(t, http.StatusOK, types.STKPushResponse{ResponseCode: "0"})
	defer server.Close()
	mpesa.SetBaseURL(server.URL)

	var wg sync.WaitGroup
	payload := types.STKPushRequest{
		AccessToken:       "test-token",
		BusinessShortCode: "123456",
		Password:          "encoded_password",
		Amount:            "100",
		PartyA:            "254700000000",
		PartyB:            "123456",
		PhoneNumber:       "254700000000",
		CallBackURL:       "https://callback.example.com",
		AccountReference:  "Test123",
		TransactionDesc:   "Payment",
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := mpesa.STKPush(ctx, payload)
			if err != nil {
				t.Errorf("concurrent STKPush failed: %v", err)
			}
		}()
	}
	wg.Wait()
}
