// Package client provides a high-performance, concurrent Mpesa Daraja API wrapper
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
)

// Mpesa is the main client for interacting with the Mpesa Daraja API

type Mpesa struct {
	baseURL string
	headers map[string]string
	client  *http.Client
	mu      sync.RWMutex
}

func (m *Mpesa) SetBaseURL(url string) {
	m.baseURL = url
}
func (m *Mpesa) BaseURL() string {
	return m.baseURL
}
func (m *Mpesa) Headers() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.headers
}

// NewMpesa initializes a new Mpesa client
func NewMpesa() *Mpesa {
	return &Mpesa{
		baseURL: "https://api.safaricom.co.ke",
		headers: make(map[string]string),
		client:  &http.Client{},
	}
}

// setHeaders sets the authorization headers with the provided access token
func (m *Mpesa) setHeaders(accessToken string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.headers = map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", accessToken),
	}
}

// GetAccessToken retrieves an Oauth access token using consumer key and secret
func (m *Mpesa) GetAccessToken(ctx context.Context, consumerKey string, consumerSecret string) (*types.AccessTokenResponse, error) {
	url := m.baseURL + "/oauth/v1/generate?grant_type=client_credentials"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.SetBasicAuth(consumerKey, consumerSecret)
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token %w", err)
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
	requiredKeys := []string{
		"AccessToken", "BusinessShortCode", "Password", "Amount",
		"PartyA", "PartyB", "PhoneNumber", "CallBackURL",
		"AccountReference", "TransactionDesc",
	}
	cleanedPayload, err := utils.CheckKeys(requiredKeys, payload)
	if err != nil {
		return nil, err
	}
	cleanedPayload["TransactionType"] = "CustomerPayBillOnline"
	cleanedPayload["Timestamp"] = utils.GetTimestamp()

	accessToken, _ := cleanedPayload["AccessToken"].(string)
	m.setHeaders(accessToken)
	delete(cleanedPayload, "AccessToken")

	url := m.baseURL + "/mpesa/stkpush/v1/processrequest"
	body, err := json.Marshal(cleanedPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
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
		return nil, fmt.Errorf("failed to send STK push request: %w", err)
	}
	defer resp.Body.Close()

	var response types.STKPushResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &response, nil
}

// RegisterURL registers validation and confirmation URLs.
func (m *Mpesa) RegisterURL(ctx context.Context, payload types.RegisterURLRequest) (*types.RegisterURLResponse, error) {
	requiredKeys := []string{"AccessToken", "ShortCode", "ResponseType", "ConfirmationURL", "ValidationURL"}
	cleanedPayload, err := utils.CheckKeys(requiredKeys, payload)
	if err != nil {
		return nil, err
	}

	accessToken, _ := cleanedPayload["AccessToken"].(string)
	m.setHeaders(accessToken)
	delete(cleanedPayload, "AccessToken")

	url := m.baseURL + "/mpesa/c2b/v2/registerurl" //Using v2 for all apps using C2B V2
	body, err := json.Marshal(cleanedPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
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
		return nil, fmt.Errorf("failed to send register URL request: %w", err)
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
	requiredKeys := []string{"AccessToken", "ShortCode", "Amount", "Msisdn", "BillRefNumber"}
	cleanedPayload, err := utils.CheckKeys(requiredKeys, payload)
	if err != nil {
		return nil, err
	}
	cleanedPayload["CommandID"] = "CustomerPayBillOnline"

	accessToken, _ := cleanedPayload["AccessToken"].(string)
	m.setHeaders(accessToken)
	delete(cleanedPayload, "AccessToken")

	url := m.baseURL + "/mpesa/c2b/v1/simulate"
	body, err := json.Marshal(cleanedPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
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
		return nil, fmt.Errorf("failed to send simulate transaction request: %w", err)
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
	requiredKeys := []string{
		"AccessToken", "Initiator", "SecurityCredential", "TransactionID",
		"Amount", "ReceiverParty", "ResultURL", "QueueTimeOutURL", "Remarks", "Occasion",
	}
	cleanedPayload, err := utils.CheckKeys(requiredKeys, payload)
	if err != nil {
		return nil, err
	}
	cleanedPayload["CommandID"] = "TransactionReversal"

	accessToken, _ := cleanedPayload["AccessToken"].(string)
	m.setHeaders(accessToken)
	delete(cleanedPayload, "AccessToken")

	url := m.baseURL + "/mpesa/reversal/v1/request"
	body, err := json.Marshal(cleanedPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
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
		return nil, fmt.Errorf("failed to send reverse transaction request: %w", err)
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
	requiredKeys := []string{
		"AccessToken", "Initiator", "SecurityCredential", "TransactionID",
		"PartyA", "ResultURL", "QueueTimeOutURL", "Remarks", "Occasion",
	}
	cleanedPayload, err := utils.CheckKeys(requiredKeys, payload)
	if err != nil {
		return nil, err
	}
	cleanedPayload["CommandID"] = "TransactionStatusQuery"
	cleanedPayload["IdentifierType"] = "4"

	accessToken, _ := cleanedPayload["AccessToken"].(string)
	m.setHeaders(accessToken)
	delete(cleanedPayload, "AccessToken")

	url := m.baseURL + "/mpesa/transactionstatus/v1/query"
	body, err := json.Marshal(cleanedPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
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
		return nil, fmt.Errorf("failed to send query transaction request: %w", err)
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
	requiredKeys := []string{
		"AccessToken", "Initiator", "SecurityCredential", "PartyA",
		"Remarks", "QueueTimeOutURL", "ResultURL",
	}
	cleanedPayload, err := utils.CheckKeys(requiredKeys, payload)
	if err != nil {
		return nil, err
	}
	cleanedPayload["CommandID"] = "AccountBalance"
	cleanedPayload["IdentifierType"] = "4"

	accessToken, _ := cleanedPayload["AccessToken"].(string)
	m.setHeaders(accessToken)
	delete(cleanedPayload, "AccessToken")

	url := m.baseURL + "/mpesa/accountbalance/v1/query"
	body, err := json.Marshal(cleanedPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
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
		return nil, fmt.Errorf("failed to send get balance request: %w", err)
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
	requiredKeys := []string{
		"AccessToken", "InitiatorName", "SecurityCredential", "Amount",
		"PartyA", "PartyB", "Remarks", "QueueTimeOutURL", "ResultURL", "Occasion",
	}
	cleanedPayload, err := utils.CheckKeys(requiredKeys, payload)
	if err != nil {
		return nil, err
	}
	cleanedPayload["CommandID"] = "PromotionPayment"

	accessToken, _ := cleanedPayload["AccessToken"].(string)
	m.setHeaders(accessToken)
	delete(cleanedPayload, "AccessToken")

	url := m.baseURL + "/mpesa/b2c/v1/paymentrequest"
	body, err := json.Marshal(cleanedPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
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
		return nil, fmt.Errorf("failed to send B2C request: %w", err)
	}
	defer resp.Body.Close()

	var response types.B2CSendResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &response, nil
}

// B2BSend sends funds from paybill to paybill/till
func (m *Mpesa) B2BSend(ctx context.Context, payload types.B2BSendRequest) (*types.B2BSendResponse, error) {
	requiredKeys := []string{
		"AccessToken", "Initiator", "SecurityCredential", "CommandID",
		"SenderIdentifierType", "RecieverIdentifierType", "Amount",
		"PartyA", "PartyB", "Remarks", "AccountReference", "Requester",
		"QueueTimeOutURL", "ResultURL",
	}
	cleanedPayload, err := utils.CheckKeys(requiredKeys, payload)
	if err != nil {
		return nil, err
	}

	accessToken, _ := cleanedPayload["AccessToken"].(string)
	m.setHeaders(accessToken)
	delete(cleanedPayload, "AccessToken")

	url := m.baseURL + "/mpesa/b2b/v1/paymentrequest"
	body, err := json.Marshal(cleanedPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
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
		return nil, fmt.Errorf("failed to send B2B request: %w", err)
	}
	defer resp.Body.Close()

	var response types.B2BSendResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &response, nil
}

// RegisterPullAPI registers the pull transaction API. Request API Support to add the product first to the App
func (m *Mpesa) RegisterPullAPI(ctx context.Context, payload types.RegisterPullAPIRequest) (*types.RegisterPullAPIResponse, error) {
	requiredKeys := []string{"AccessToken", "ShortCode", "NominatedNumber", "CallBackURL"}
	cleanedPayload, err := utils.CheckKeys(requiredKeys, payload)
	if err != nil {
		return nil, err
	}
	accessToken, _ := cleanedPayload["AccessToken"].(string)
	m.setHeaders(accessToken)
	delete(cleanedPayload, "AccessToken")
	cleanedPayload["RequestType"] = "Pull"

	url := m.baseURL + "/pulltransactions/v1/register"
	body, err := json.Marshal(cleanedPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
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
		return nil, fmt.Errorf("failed to send register pull API request: %w", err)
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
	requiredKeys := []string{"AccessToken", "ShortCode", "StartDate", "EndDate", "OffSetValue"}
	cleanedPayload, err := utils.CheckKeys(requiredKeys, payload)
	if err != nil {
		return nil, err
	}

	accessToken, _ := cleanedPayload["AccessToken"].(string)
	m.setHeaders(accessToken)
	delete(cleanedPayload, "AccessToken")

	url := m.baseURL + "/pulltransactions/v1/query"
	body, err := json.Marshal(cleanedPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
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
		return nil, fmt.Errorf("failed to send pull transactions request: %w", err)
	}
	defer resp.Body.Close()

	var response types.PullTransactionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &response, nil
}
