package handlers

import (
	"bytes"
	"encoding/json"
	// "errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	CURRENCY_NGN = "NGN"
	SUBUNIT_NGN  = 100
)

var (
	initTransactionURL   = "https://api.paystack.co/transaction/initialize"
	verifyTransactionURL = "https://api.paystack.co/transaction/verify"
	verifyTransferURL    = "https://api.paystack.co/transfer/verify"
	resolveBankURL       = "https://api.paystack.co/bank/resolve"
	createRecipientURL   = "https://api.paystack.co/transferrecipient"
	initTransferURL      = "https://api.paystack.co/transfer"
)

type initTransactionPayload struct {
	Amount      string `json:"amount"`
	Email       string `json:"email"`
	Reference   string `json:"reference"`
	Currency    string `json:"currency"`
	CallbackURL string `json:"callback_url"`
}

type initTransferPayload struct {
	Source        string `json:"source"`
	Amount        string `json:"amount"`
	Reason        string `json:"reason"`
	RecipientCode string `json:"recipient"`
	Reference     string `json:"reference"`
	Currency      string `json:"currency"`
}

type createRecipientPayload struct {
	Type        string `json:"type"`
	AccountName string `json:"name"`
	AccountNum  string `json:"account_number"`
	BankCode    string `json:"bank_code"`
	Currency    string `json:"currency"`
}

type paystackResponse struct {
	Status  bool         `json:"status"`
	Message string       `json:"message"`
	Data    responseData `json:"data"`
}

type responseData struct {
	Status           string  `json:"status,omitempty"`
	Reference        string  `json:"reference,omitempty"`
	Amount           float64 `json:"amount,omitempty"`
	Currency         string  `json:"currency,omitempty"`
	AccessCode       string  `json:"access_code,omitempty"`
	AuthorizationURL string  `json:"authorization_url,omitempty"`
	AccountName      string  `json:"account_name,omitempty"`
	AccountNumber    string  `json:"account_number,omitempty"`
	RecipientCode    string  `json:"recipient_code,omitempty"`
}

func initTransaction(
	email, paymentRef, callbackURL string,
	amount float64,
	apiKey string,
) (string, error) {
	reqBody := &bytes.Buffer{}
	amountStr := strconv.Itoa(int(amount * SUBUNIT_NGN))
	err := json.NewEncoder(reqBody).Encode(initTransactionPayload{
		Amount:      amountStr,
		Email:       email,
		Reference:   paymentRef,
		Currency:    CURRENCY_NGN,
		CallbackURL: callbackURL,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", initTransactionURL, reqBody)
	if err != nil {
		return "", err
	}
	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Bearer %s", apiKey),
	)
	req.Header.Add(
		"Content-Type",
		"application/json",
	)

	resp, err := execPaystackRequest(req)
	if err != nil {
		return "", err
	}

	return resp.Data.AccessCode, nil
}

func verifyTransaction(paymentRef, apiKey string) (string, error) {
	URL := fmt.Sprintf("%s/%s", verifyTransactionURL, paymentRef)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Bearer %s", apiKey),
	)

	resp, err := execPaystackRequest(req)
	if err != nil {
		return "", err
	}

	return resp.Data.Status, nil
}

func verifyTransfer(paymentRef, apiKey string) (string, error) {
	URL := fmt.Sprintf("%s/%s", verifyTransferURL, paymentRef)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Bearer %s", apiKey),
	)

	resp, err := execPaystackRequest(req)
	if err != nil {
		return "", err
	}

	return resp.Data.Status, nil
}

func resolveAccountNumber(accountNum, bankCode, apiKey string) (string, error) {
	URL := fmt.Sprintf(
		"%s?account_number=%s&bank_code=%s",
		resolveBankURL,
		accountNum,
		bankCode,
	)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Bearer %s", apiKey),
	)

	resp, err := execPaystackRequest(req)
	if err != nil {
		return "", err
	}

	return resp.Data.AccountName, nil
}

func createTransferRecipient(accountName, accountNum, bankCode, apiKey string) (string, error) {
	reqBody := &bytes.Buffer{}
	err := json.NewEncoder(reqBody).Encode(createRecipientPayload{
		Type:        "nuban",
		AccountName: accountName,
		AccountNum:  accountNum,
		BankCode:    bankCode,
		Currency:    CURRENCY_NGN,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", createRecipientURL, reqBody)
	if err != nil {
		return "", err
	}
	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Bearer %s", apiKey),
	)
	req.Header.Add(
		"Content-Type",
		"application/json",
	)

	resp, err := execPaystackRequest(req)
	if err != nil {
		return "", err
	}

	return resp.Data.RecipientCode, nil
}

func deleteRecipient(recipientCode, apiKey string) error {
	URL := fmt.Sprintf("%s/%s", createRecipientURL, recipientCode)
	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest("DELETE", URL, nil)
	if err != nil {
		return err
	}
	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Bearer %s", apiKey),
	)

	var resp *http.Response
	for range 3 {
		resp, err = client.Do(req)
		if err == nil {
			var body []byte
			body, err = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			switch resp.StatusCode {
			case http.StatusOK, http.StatusNotFound:
				return nil
			default:
				err = fmt.Errorf(
					"%s: %d %s",
					URL,
					resp.StatusCode,
					string(body),
				)
			}
		}
	}
	return err
}

func initTransfer(
	paymentRef, recipientCode, transferRemark string,
	amount float64,
	apiKey string,
) (string, error) {
	reqBody := &bytes.Buffer{}
	amountStr := strconv.Itoa(int(amount * SUBUNIT_NGN))
	err := json.NewEncoder(reqBody).Encode(initTransferPayload{
		Amount:        amountStr,
		Source:        "balance",
		Reference:     paymentRef,
		Currency:      CURRENCY_NGN,
		RecipientCode: recipientCode,
		Reason:        transferRemark,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", initTransferURL, reqBody)
	if err != nil {
		return "", err
	}
	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Bearer %s", apiKey),
	)
	req.Header.Add(
		"Content-Type",
		"application/json",
	)

	resp, err := execPaystackRequest(req)
	if err != nil {
		return "", err
	}

	return resp.Data.Status, nil
}

func execPaystackRequest(req *http.Request) (paystackResponse, error) {
	client := http.Client{Timeout: timeout}

	var err error
	var resp *http.Response
	var paystackResp paystackResponse
	for range 3 {
		resp, err = client.Do(req)
		if err == nil {
			defer func() { resp.Body.Close() }()
			var body []byte
			body, err = ioutil.ReadAll(resp.Body)
			if err == nil {
				switch resp.StatusCode {
				case http.StatusOK:
					err = json.Unmarshal(body, &paystackResp)
				default:
					err = fmt.Errorf(
						"%s: %d %s",
						req.URL,
						resp.StatusCode,
						string(body),
					)
					fmt.Println(err)
				}
			}
		}
	}
	return paystackResp, err
}
