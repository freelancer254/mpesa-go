// Package types defines the request and response structs for the M-Pesa Daraja API.
package types

// AccessTokenResponse represents the response for an access token request.
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

// STKPushRequest represents the payload for an STK Push request.
type STKPushRequest struct {
	AccessToken       string `json:"AccessToken" validate:"required"`
	BusinessShortCode string `json:"BusinessShortCode" validate:"required,numeric"`
	Password          string `json:"Password" validate:"required"`
	Amount            string `json:"Amount" validate:"required,numeric"`
	PartyA            string `json:"PartyA" validate:"required,numeric"`
	PartyB            string `json:"PartyB" validate:"required,numeric"`
	PhoneNumber       string `json:"PhoneNumber" validate:"required,numeric"`
	CallBackURL       string `json:"CallBackURL" validate:"required,url"`
	AccountReference  string `json:"AccountReference" validate:"required"`
	TransactionDesc   string `json:"TransactionDesc" validate:"required"`
}

// STKPushQueryRequest represents the payload for an STK Push Query request.
type STKPushQueryRequest struct {
	AccessToken       string `json:"AccessToken" validate:"required"`
	BusinessShortCode string `json:"BusinessShortCode" validate:"required,numeric"`
	Password          string `json:"Password" validate:"required"`
	Timestamp         string `json:"Timestamp" validate:"required,numeric"`
	CheckoutRequestID string `json:"CheckoutRequestID"`
}

// STKPushResponse represents the response for an STK Push request.
type STKPushResponse struct {
	MerchantRequestID   string `json:"MerchantRequestID" validate:"required"`
	ResponseCode        string `json:"ResponseCode"`
	CheckoutRequestID   string `json:"CheckoutRequestID"`
	ResponseDescription string `json:"ResponseDescription"`
	CustomerMessage     string `json:"CustomerMessage"`
}

// STKPushQueryResponse represents the response for an STK Push Query request.
type STKPushQueryResponse struct {
	MerchantRequestID   string `json:"MerchantRequestID" validate:"required"`
	ResponseCode        string `json:"ResponseCode"`
	CheckoutRequestID   string `json:"CheckoutRequestID"`
	ResponseDescription string `json:"ResponseDescription"`
	ResultCode          string `json:"ResultCode" validate:"required, numeric"` //0, 1032
	ResultDesc          string `json:"ResultDesc" validate:"required"`
}

// STKPushError represents the error response for an STK Push request.
type STKPushError struct {
	Body struct {
		StkCallback struct {
			MerchantRequestID string `json:"MerchantRequestID" validate:"required"`
			CheckoutRequestID string `json:"CheckoutRequestID" validate:"required"`
			ResultCode        string `json:"ResultCode" validate:"required, numeric"`
			ResultDesc        string `json:"ResultDesc" validate:"required"`
		} `json:"stkCallback" validate:"required"`
	} `json:"Body" validate:"required"`
}

// RegisterURLRequest represents the payload for registering URLs.
type RegisterURLRequest struct {
	AccessToken     string `json:"AccessToken" validate:"required"`
	ShortCode       string `json:"ShortCode" validate:"required,numeric"`
	ResponseType    string `json:"ResponseType" validate:"required,eq=Completed|eq=Cancelled"`
	ConfirmationURL string `json:"ConfirmationURL" validate:"required,url"`
	ValidationURL   string `json:"ValidationURL" validate:"required,url"`
}

// RegisterURLResponse represents the response for registering URLs.
type RegisterURLResponse struct {
	OriginatorCoversationID string `json:"OriginatorCoversationID"`
	ResultCode              string `json:"ResultCode" validate:"required, numeric"`
	ResponseDescription     string `json:"ResponseDescription"`
}

// SimulateTransactionRequest represents the payload for simulating a transaction.
type SimulateTransactionRequest struct {
	AccessToken   string `json:"AccessToken" validate:"required"`
	ShortCode     string `json:"ShortCode" validate:"required,numeric"`
	Amount        string `json:"Amount" validate:"required,numeric"`
	Msisdn        string `json:"Msisdn" validate:"required,numeric"`
	BillRefNumber string `json:"BillRefNumber" validate:"required"`
}

// SimulateTransactionResponse represents the response for simulating a transaction.
type SimulateTransactionResponse struct {
	ConversationID           string `json:"ConversationID"`
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ResponseDescription      string `json:"ResponseDescription"`
}

// ReverseTransactionRequest represents the payload for reversing a transaction.
type ReverseTransactionRequest struct {
	AccessToken            string `json:"AccessToken" validate:"required"`
	Initiator              string `json:"Initiator" validate:"required"`
	SecurityCredential     string `json:"SecurityCredential" validate:"required"`
	TransactionID          string `json:"TransactionID" validate:"required"`
	Amount                 string `json:"Amount" validate:"required,numeric"`
	ReceiverParty          string `json:"ReceiverParty" validate:"required,numeric"`
	ReceiverIdentifierType string `json:"ReceiverIdentifierType" validate:"required,numeric"`
	ResultURL              string `json:"ResultURL" validate:"required,url"`
	QueueTimeOutURL        string `json:"QueueTimeOutURL" validate:"required,url"`
	Remarks                string `json:"Remarks" validate:"required"`
	Occasion               string `json:"Occasion" validate:"required"`
}

// ReverseTransactionResponse represents the response for reversing a transaction.
type ReverseTransactionResponse struct {
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ConversationID           string `json:"ConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

// QueryTransactionRequest represents the payload for querying a transaction.
type QueryTransactionRequest struct {
	AccessToken              string `json:"AccessToken" validate:"required"`
	Initiator                string `json:"Initiator" validate:"required"`
	SecurityCredential       string `json:"SecurityCredential" validate:"required"`
	TransactionID            string `json:"TransactionID,omitempty"`
	OriginatorConversationID string `json:"OriginatorConversationID,omitempty"`
	PartyA                   string `json:"PartyA" validate:"required,numeric"`
	IdentifierType           string `json:"IdentifierType" validate:"required,numeric"`
	ResultURL                string `json:"ResultURL" validate:"required,url"`
	QueueTimeOutURL          string `json:"QueueTimeOutURL" validate:"required,url"`
	Remarks                  string `json:"Remarks" validate:"required"`
	Occasion                 string `json:"Occasion" validate:"required"`
}

// QueryTransactionResponse represents the response for querying a transaction.
type QueryTransactionResponse struct {
	ConversationID           string `json:"ConversationID"`
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ResponseDescription      string `json:"ResponseDescription"`
	ResponseCode             string `json:"ResponseCode"`
}

// GetBalanceRequest represents the payload for querying the account balance.
type GetBalanceRequest struct {
	AccessToken        string `json:"AccessToken" validate:"required"`
	Initiator          string `json:"Initiator" validate:"required"`
	SecurityCredential string `json:"SecurityCredential" validate:"required"`
	PartyA             string `json:"PartyA" validate:"required,numeric"`
	IdentifierType     string `json:"IdentifierType" validate:"required,numeric"`
	Remarks            string `json:"Remarks" validate:"required"`
	QueueTimeOutURL    string `json:"QueueTimeOutURL" validate:"required,url"`
	ResultURL          string `json:"ResultURL" validate:"required,url"`
}

// GetBalanceResponse represents the response for querying the account balance.
type GetBalanceResponse struct {
	ConversationID           string `json:"ConversationID"`
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ResponseDescription      string `json:"ResponseDescription"`
	ResponseCode             string `json:"ResponseCode"`
}

// B2CSendRequest represents the payload for a B2C send request.
type B2CSendRequest struct {
	AccessToken        string `json:"AccessToken" validate:"required"`
	InitiatorName      string `json:"InitiatorName" validate:"required"`
	SecurityCredential string `json:"SecurityCredential" validate:"required"`
	CommandID          string `json:"CommandID" validate:"required"`
	Amount             string `json:"Amount" validate:"required,numeric"`
	PartyA             string `json:"PartyA" validate:"required,numeric"`
	PartyB             string `json:"PartyB" validate:"required,numeric"`
	Remarks            string `json:"Remarks" validate:"required"`
	QueueTimeOutURL    string `json:"QueueTimeOutURL" validate:"required,url"`
	ResultURL          string `json:"ResultURL" validate:"required,url"`
	Occasion           string `json:"Occasion" validate:"required"`
}

// B2CSendResponse represents the response for a B2C send request.
type B2CSendResponse struct {
	ConversationID           string `json:"ConversationID"`
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ResponseDescription      string `json:"ResponseDescription"`
	ResponseCode             string `json:"ResponseCode"`
}

// B2BSendRequest represents the payload for a B2B send request.
type B2BSendRequest struct {
	AccessToken            string `json:"AccessToken" validate:"required"`
	Initiator              string `json:"Initiator" validate:"required"`
	SecurityCredential     string `json:"SecurityCredential" validate:"required"`
	CommandID              string `json:"CommandID" validate:"required"`
	SenderIdentifierType   string `json:"SenderIdentifierType" validate:"required,numeric"`
	ReceiverIdentifierType string `json:"RecieverIdentifierType" validate:"required,numeric"`
	Amount                 string `json:"Amount" validate:"required,numeric"`
	PartyA                 string `json:"PartyA" validate:"required,numeric"`
	PartyB                 string `json:"PartyB" validate:"required,numeric"`
	Remarks                string `json:"Remarks" validate:"required"`
	AccountReference       string `json:"AccountReference" validate:"required"`
	Requester              string `json:"Requester" validate:"required,numeric"`
	QueueTimeOutURL        string `json:"QueueTimeOutURL" validate:"required,url"`
	ResultURL              string `json:"ResultURL" validate:"required,url"`
}

// B2BSendResponse represents the response for a B2B send request.
type B2BSendResponse struct {
	ConversationID           string `json:"ConversationID"`
	OriginatorConversationID string `json:"OriginatorConversationID"`
	ResponseCode             string `json:"ResponseCode"`
	ResponseDescription      string `json:"ResponseDescription"`
}

// RegisterPullAPIRequest represents the payload for registering the pull API.
type RegisterPullAPIRequest struct {
	AccessToken     string `json:"AccessToken" validate:"required"`
	ShortCode       string `json:"ShortCode" validate:"required,numeric"`
	NominatedNumber string `json:"NominatedNumber" validate:"required,numeric"`
	CallBackURL     string `json:"CallBackURL" validate:"required,url"`
}

// RegisterPullAPIResponse represents the response for registering the pull API.
type RegisterPullAPIResponse struct {
	ResponseRefID       string `json:"ResponseRefID"`
	ResponseStatus      string `json:"ResponseStatus"`
	ShortCode           string `json:"ShortCode"`
	ResponseDescription string `json:"ResponseDescription"`
}

// PullTransactionsRequest represents the payload for pulling transactions.
type PullTransactionsRequest struct {
	AccessToken string `json:"AccessToken" validate:"required"`
	ShortCode   string `json:"ShortCode" validate:"required,numeric"`
	StartDate   string `json:"StartDate" validate:"required,datetime=2006-01-02"`
	EndDate     string `json:"EndDate" validate:"required,datetime=2006-01-02"`
	OffSetValue string `json:"OffSetValue" validate:"required,numeric"`
}

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
