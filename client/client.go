// Package client provides a high-performance, concurrent M-Pesa Daraja API wrapper.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/freelancer254/mpesa-go/types"
	"github.com/freelancer254/mpesa-go/utils"
	"github.com/go-playground/validator/v10"
)

// Mpesa is the main client for interacting with the M-Pesa Daraja API.
type Mpesa struct {
	baseURL  string
	headers  map[string]string
	client   *http.Client
	mu       sync.RWMutex
	validate *validator.Validate
}

// NewMpesa initializes a new Mpesa client.
func NewMpesa() *Mpesa {
	return &Mpesa{
		baseURL:  "https://api.safaricom.co.ke",
		headers:  make(map[string]string),
		client:   &http.Client{},
		validate: validator.New(),
	}
}

// SetBaseURL sets the base URL for testing purposes.
func (m *Mpesa) SetBaseURL(url string) {
	m.baseURL = url
}

// BaseURL returns the base URL.
func (m *Mpesa) BaseURL() string {
	return m.baseURL
}

// Headers returns the headers.
func (m *Mpesa) Headers() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.headers
}

// setHeaders sets the authorization headers with the provided access token.
func (m *Mpesa) setHeaders(accessToken string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.headers = map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", accessToken),
	}
}

// doRequest performs an HTTP request with the given method, URL, and payload.
func (m *Mpesa) doRequest(ctx context.Context, method, url string, payload interface{}) (*http.Response, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	m.mu.RLock()
	for k, v := range m.headers {
		req.Header.Set(k, v)
	}
	m.mu.RUnlock()

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	return resp, nil
}

// GetAccessToken retrieves an OAuth access token using consumer key and secret.
func (m *Mpesa) GetAccessToken(ctx context.Context, consumerKey string, consumerSecret string) (*types.AccessTokenResponse, error) {
	url := m.baseURL + "/oauth/v1/generate?grant_type=client_credentials"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.SetBasicAuth(consumerKey, consumerSecret)

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}
	defer resp.Body.Close()

	var token types.AccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &token, nil
}

// STKPush initiates a transaction using STK Push.
func (m *Mpesa) STKPush(ctx context.Context, payload types.STKPushRequest) (*types.STKPushResponse, error) {
	if err := m.validate.Struct(payload); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	m.setHeaders(payload.AccessToken)
	payloadMap := map[string]interface{}{
		"BusinessShortCode": payload.BusinessShortCode,
		"Password":          payload.Password,
		"Timestamp":         utils.GetTimestamp(),
		"TransactionType":   "CustomerPayBillOnline",
		"Amount":            payload.Amount,
		"PartyA":            payload.PartyA,
		"PartyB":            payload.PartyB,
		"PhoneNumber":       payload.PhoneNumber,
		"CallBackURL":       payload.CallBackURL,
		"AccountReference":  payload.AccountReference,
		"TransactionDesc":   payload.TransactionDesc,
	}

	url := m.baseURL + "/mpesa/stkpush/v1/processrequest"
	resp, err := m.doRequest(ctx, http.MethodPost, url, payloadMap)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp types.STKPushError
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, fmt.Errorf("failed to decode error response: %w", err)
		}
		if err := m.validate.Struct(errorResp); err != nil {
			return nil, fmt.Errorf("invalid error response: %w", err)
		}
		return nil, fmt.Errorf("STK Push failed: %s (code: %s)", errorResp.Body.StkCallback.ResultDesc, errorResp.Body.StkCallback.ResultCode)
	}

	var response types.STKPushResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &response, nil
}

// STKPushQuery initiates a transaction query for tx initiated using STK Push.
func (m *Mpesa) STKPushQuery(ctx context.Context, payload types.STKPushQueryRequest) (*types.STKPushQueryResponse, error) {
	if err := m.validate.Struct(payload); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	m.setHeaders(payload.AccessToken)
	payloadMap := map[string]interface{}{
		"BusinessShortCode": payload.BusinessShortCode,
		"Password":          payload.Password,
		"Timestamp":         utils.GetTimestamp(),
		"CheckoutRequestID": payload.CheckoutRequestID,
	}

	url := m.baseURL + "/mpesa/stkpushquery/v1/query"
	resp, err := m.doRequest(ctx, http.MethodPost, url, payloadMap)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response types.STKPushQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &response, nil
}

// RegisterURL registers validation and confirmation URLs.
func (m *Mpesa) RegisterURL(ctx context.Context, payload types.RegisterURLRequest) (*types.RegisterURLResponse, error) {
	if err := m.validate.Struct(payload); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	m.setHeaders(payload.AccessToken)
	payloadMap := map[string]interface{}{
		"ShortCode":       payload.ShortCode,
		"ResponseType":    payload.ResponseType,
		"ConfirmationURL": payload.ConfirmationURL,
		"ValidationURL":   payload.ValidationURL,
	}

	url := m.baseURL + "/mpesa/c2b/v2/registerurl"
	resp, err := m.doRequest(ctx, http.MethodPost, url, payloadMap)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response types.RegisterURLResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &response, nil
}

// SimulateTransaction simulates a customer transaction for testing.
func (m *Mpesa) SimulateTransaction(ctx context.Context, payload types.SimulateTransactionRequest) (*types.SimulateTransactionResponse, error) {
	if err := m.validate.Struct(payload); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	m.setHeaders(payload.AccessToken)
	payloadMap := map[string]interface{}{
		"ShortCode":     payload.ShortCode,
		"CommandID":     "CustomerPayBillOnline",
		"Amount":        payload.Amount,
		"Msisdn":        payload.Msisdn,
		"BillRefNumber": payload.BillRefNumber,
	}

	url := m.baseURL + "/mpesa/c2b/v1/simulate"
	resp, err := m.doRequest(ctx, http.MethodPost, url, payloadMap)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response types.SimulateTransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &response, nil
}

// ReverseTransaction reverses a transaction.
func (m *Mpesa) ReverseTransaction(ctx context.Context, payload types.ReverseTransactionRequest) (*types.ReverseTransactionResponse, error) {
	if err := m.validate.Struct(payload); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	m.setHeaders(payload.AccessToken)
	payloadMap := map[string]interface{}{
		"Initiator":              payload.Initiator,
		"SecurityCredential":     payload.SecurityCredential,
		"CommandID":              "TransactionReversal",
		"TransactionID":          payload.TransactionID,
		"Amount":                 payload.Amount,
		"ReceiverParty":          payload.ReceiverParty,
		"ReceiverIdentifierType": payload.ReceiverIdentifierType,
		"ResultURL":              payload.ResultURL,
		"QueueTimeOutURL":        payload.QueueTimeOutURL,
		"Remarks":                payload.Remarks,
		"Occasion":               payload.Occasion,
	}

	url := m.baseURL + "/mpesa/reversal/v1/request"
	resp, err := m.doRequest(ctx, http.MethodPost, url, payloadMap)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response types.ReverseTransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &response, nil
}

// QueryTransaction queries the status of a transaction.
func (m *Mpesa) QueryTransaction(ctx context.Context, payload types.QueryTransactionRequest) (*types.QueryTransactionResponse, error) {
	if err := m.validate.Struct(payload); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	m.setHeaders(payload.AccessToken)
	payloadMap := map[string]interface{}{
		"Initiator":                payload.Initiator,
		"SecurityCredential":       payload.SecurityCredential,
		"CommandID":                "TransactionStatusQuery",
		"TransactionID":            payload.TransactionID,
		"OriginatorConversationID": payload.OriginatorConversationID,
		"PartyA":                   payload.PartyA,
		"IdentifierType":           payload.IdentifierType,
		"ResultURL":                payload.ResultURL,
		"QueueTimeOutURL":          payload.QueueTimeOutURL,
		"Remarks":                  payload.Remarks,
		"Occasion":                 payload.Occasion,
	}

	url := m.baseURL + "/mpesa/transactionstatus/v1/query"
	resp, err := m.doRequest(ctx, http.MethodPost, url, payloadMap)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response types.QueryTransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &response, nil
}

// GetBalance retrieves the paybill account balance.
func (m *Mpesa) GetBalance(ctx context.Context, payload types.GetBalanceRequest) (*types.GetBalanceResponse, error) {
	if err := m.validate.Struct(payload); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	m.setHeaders(payload.AccessToken)
	payloadMap := map[string]interface{}{
		"Initiator":          payload.Initiator,
		"SecurityCredential": payload.SecurityCredential,
		"CommandID":          "AccountBalance",
		"PartyA":             payload.PartyA,
		"IdentifierType":     payload.IdentifierType,
		"Remarks":            payload.Remarks,
		"QueueTimeOutURL":    payload.QueueTimeOutURL,
		"ResultURL":          payload.ResultURL,
	}

	url := m.baseURL + "/mpesa/accountbalance/v1/query"
	resp, err := m.doRequest(ctx, http.MethodPost, url, payloadMap)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response types.GetBalanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &response, nil
}

// B2CSend sends funds from paybill to customer.
func (m *Mpesa) B2CSend(ctx context.Context, payload types.B2CSendRequest) (*types.B2CSendResponse, error) {
	if err := m.validate.Struct(payload); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	m.setHeaders(payload.AccessToken)
	payloadMap := map[string]interface{}{
		"InitiatorName":      payload.InitiatorName,
		"SecurityCredential": payload.SecurityCredential,
		"CommandID":          payload.CommandID,
		"Amount":             payload.Amount,
		"PartyA":             payload.PartyA,
		"PartyB":             payload.PartyB,
		"Remarks":            payload.Remarks,
		"QueueTimeOutURL":    payload.QueueTimeOutURL,
		"ResultURL":          payload.ResultURL,
		"Occasion":           payload.Occasion,
	}

	url := m.baseURL + "/mpesa/b2c/v1/paymentrequest"
	resp, err := m.doRequest(ctx, http.MethodPost, url, payloadMap)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response types.B2CSendResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &response, nil
}

// B2BSend sends funds from paybill to paybill.
func (m *Mpesa) B2BSend(ctx context.Context, payload types.B2BSendRequest) (*types.B2BSendResponse, error) {
	if err := m.validate.Struct(payload); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	m.setHeaders(payload.AccessToken)
	payloadMap := map[string]interface{}{
		"Initiator":              payload.Initiator,
		"SecurityCredential":     payload.SecurityCredential,
		"CommandID":              payload.CommandID,
		"SenderIdentifierType":   payload.SenderIdentifierType,
		"RecieverIdentifierType": payload.ReceiverIdentifierType,
		"Amount":                 payload.Amount,
		"PartyA":                 payload.PartyA,
		"PartyB":                 payload.PartyB,
		"Remarks":                payload.Remarks,
		"AccountReference":       payload.AccountReference,
		"Requester":              payload.Requester,
		"QueueTimeOutURL":        payload.QueueTimeOutURL,
		"ResultURL":              payload.ResultURL,
	}

	url := m.baseURL + "/mpesa/b2b/v1/paymentrequest"
	resp, err := m.doRequest(ctx, http.MethodPost, url, payloadMap)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response types.B2BSendResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &response, nil
}

// RegisterPullAPI registers the pull transaction API.
func (m *Mpesa) RegisterPullAPI(ctx context.Context, payload types.RegisterPullAPIRequest) (*types.RegisterPullAPIResponse, error) {
	if err := m.validate.Struct(payload); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	m.setHeaders(payload.AccessToken)
	payloadMap := map[string]interface{}{
		"ShortCode":       payload.ShortCode,
		"NominatedNumber": payload.NominatedNumber,
		"CallBackURL":     payload.CallBackURL,
	}

	url := m.baseURL + "/pulltransactions/v1/register"
	resp, err := m.doRequest(ctx, http.MethodPost, url, payloadMap)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response types.RegisterPullAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &response, nil
}

// PullTransactions pulls transactions for a shortcode.
func (m *Mpesa) PullTransactions(ctx context.Context, payload types.PullTransactionsRequest) (*types.PullTransactionsResponse, error) {
	if err := m.validate.Struct(payload); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	m.setHeaders(payload.AccessToken)
	payloadMap := map[string]interface{}{
		"ShortCode":   payload.ShortCode,
		"StartDate":   payload.StartDate,
		"EndDate":     payload.EndDate,
		"OffSetValue": payload.OffSetValue,
	}

	url := m.baseURL + "/pulltransactions/v1/query"
	resp, err := m.doRequest(ctx, http.MethodPost, url, payloadMap)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response types.PullTransactionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &response, nil
}
