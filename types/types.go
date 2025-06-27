// Package types defines the request and response structs for the M-Pesa Daraja API.
package types

// AccessTokenResponse represents the response for an access token request.
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// STKPushRequest represents the payload for an STK Push request.
type STKPushRequest map[string]interface{}

// STKPushResponse represents the response for an STK Push request.
type STKPushResponse struct {
	ResponseCode        int    `json:"ResponseCode"`
	CheckoutRequestID   string `json:"CheckoutRequestID"`
	ResponseDescription string `json:"ResponseDescription"`
	ErrorCode           string `json:"errorCode,omitempty"`
}

// RegisterURLRequest represents the payload for registering URLs.
type RegisterURLRequest map[string]interface{}

// RegisterURLResponse represents the response for registering URLs.
type RegisterURLResponse struct {
	ResponseDescription string `json:"ResponseDescription"`
	ErrorCode           string `json:"errorCode,omitempty"`
}

// SimulateTransactionRequest represents the payload for simulating a transaction.
type SimulateTransactionRequest map[string]interface{}

// SimulateTransactionResponse represents the response for simulating a transaction.
type SimulateTransactionResponse struct {
	ConversationID           string `json:"ConversationID"`
	OriginatorConversationID string `json:"OriginatorCoversationID"`
	ResponseDescription      string `json:"ResponseDescription"`
	ErrorCode                string `json:"errorCode,omitempty"`
}

// ReverseTransactionRequest represents the payload for reversing a transaction.
type ReverseTransactionRequest map[string]interface{}

// ReverseTransactionResponse represents the response for reversing a transaction.
type ReverseTransactionResponse struct {
	Result struct {
		ResultType int `json:"ResultType"`
		ResultCode int `json:"ResultCode"`
	} `json:"Result"`
	ErrorCode string `json:"errorCode,omitempty"`
}

// QueryTransactionRequest represents the payload for querying a transaction.
type QueryTransactionRequest map[string]interface{}

// QueryTransactionResponse represents the response for querying a transaction.
type QueryTransactionResponse struct {
	ConversationID           string `json:"ConversationID"`
	OriginatorConversationID string `json:"OriginatorCoversationID"`
	ResponseDescription      string `json:"ResponseDescription"`
	ErrorCode                string `json:"errorCode,omitempty"`
}

// GetBalanceRequest represents the payload for querying the account balance.
type GetBalanceRequest map[string]interface{}

// GetBalanceResponse represents the response for querying the account balance.
type GetBalanceResponse struct {
	ConversationID           string `json:"ConversationID"`
	OriginatorConversationID string `json:"OriginatorCoversationID"`
	ResponseDescription      string `json:"ResponseDescription"`
	ErrorCode                string `json:"errorCode,omitempty"`
}

// B2CSendRequest represents the payload for a B2C send request.
type B2CSendRequest map[string]interface{}

// B2CSendResponse represents the response for a B2C send request.
type B2CSendResponse struct {
	ConversationID           string `json:"ConversationID"`
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ResponseDescription      string `json:"ResponseDescription"`
	ErrorCode                string `json:"errorCode,omitempty"`
}

// B2BSendRequest represents the payload for a B2B send request.
type B2BSendRequest map[string]interface{}

// B2BSendResponse represents the response for a B2B send request.
type B2BSendResponse struct {
	ConversationID           string `json:"ConversationID"`
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
	ErrorCode                string `json:"errorCode,omitempty"`
}

// RegisterPullAPIRequest represents the payload for registering the pull API.
type RegisterPullAPIRequest map[string]interface{}

// RegisterPullAPIResponse represents the response for registering the pull API.
type RegisterPullAPIResponse struct {
	ResponseRefID       string `json:"ResponseRefID"`
	ResponseStatus      string `json:"ResponseStatus"`
	ShortCode           string `json:"ShortCode"`
	ResponseDescription string `json:"ResponseDescription"`
}

// PullTransactionsRequest represents the payload for pulling transactions.
type PullTransactionsRequest map[string]interface{}

// PullTransactionsResponse represents the response for pulling transactions.
type PullTransactionsResponse struct {
	ResponseRefID   string        `json:"ResponseRefID"`
	ResponseCode    string        `json:"ResponseCode"`
	ResponseMessage string        `json:"ResponseMessage"`
	Transactions    []Transaction `json:"Response"`
}

// Transaction represents a single transaction in the pull transactions response.
type Transaction struct {
	TransactionID    string `json:"transactionId"`
	TrxDate          string `json:"trxDate"`
	Msisdn           int64  `json:"msisdn"`
	Sender           string `json:"sender"`
	TransactionType  string `json:"transactiontype"`
	BillReference    string `json:"billreference"`
	Amount           string `json:"amount"`
	OrganizationName string `json:"organizationname"`
}
