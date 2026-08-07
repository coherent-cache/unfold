package api

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type TransactionsResponse struct {
	Meta struct {
		RequestID string    `json:"request_id"`
		Timestamp time.Time `json:"timestamp"`
		URI       string    `json:"uri"`
	} `json:"meta"`
	Data struct {
		Transactions []struct {
			UUID           string    `json:"uuid"`
			Amount         float64   `json:"amount"`
			SourceAmount   float64   `json:"source_amount"`
			CurrentBalance float64   `json:"current_balance"`
			TxnTimestamp   time.Time `json:"txn_timestamp"`
			Mode           string    `json:"mode"`
			Type           string    `json:"type"`
			Source         string    `json:"source"`
			Narration      string    `json:"narration"`
			Category       *struct {
				ID            string `json:"id"`
				SubcategoryID string `json:"subcategory_id"`
			} `json:"category"`
			Merchant *struct {
				Name    string          `json:"name"`
				Type    string          `json:"type"`
				Address json.RawMessage `json:"address"`
			} `json:"merchant"`
			AccountID                      string          `json:"account_id"`
			FinancialInformationProviderID string          `json:"financial_information_provider_id"`
			Tags                           json.RawMessage `json:"tags"`
			Kind                           string          `json:"kind"`
			Notes                          json.RawMessage `json:"notes"`
			ExcludedFromCashFlow           bool            `json:"excluded_from_cash_flow"`
			IsBookmarked                   bool            `json:"is_bookmarked"`
			TransactionID                  string          `json:"transaction_id"`
			Reference                      string          `json:"reference"`
			Summary                        string          `json:"summary"`
			BeforeFoldAccount              bool            `json:"before_fold_account"`
			Via                            json.RawMessage `json:"via"`
			AccountIn                      json.RawMessage `json:"account_in"`
			Refund                         *struct {
				Status     string          `json:"status"`
				ReceivedOn json.RawMessage `json:"received_on"`
			} `json:"refund"`
			GroupIDs                     json.RawMessage `json:"group_ids"`
			ContactID                    json.RawMessage `json:"contact_id"`
			IsF1Predicted                json.RawMessage `json:"is_f1_predicted"`
			Currency                     string          `json:"currency"`
			SourceCurrency               string          `json:"source_currency"`
			UserManualAdded              *bool           `json:"user_manual_added"`
			SplitType                    string          `json:"split_type"`
			ParentTransactionID          string          `json:"parent_transaction_id"`
			FinancialInformationProvider *struct {
				Name string `json:"name"`
			} `json:"financial_information_provider"`
		} `json:"transactions"`
		Counts []struct {
			Date              string `json:"date"`
			Total             int    `json:"total"`
			BeforeFoldAccount int    `json:"before_fold_account"`
			AfterFoldAccount  int    `json:"after_fold_account"`
		} `json:"counts"`
		Total         int         `json:"total"`
		SearchSummary interface{} `json:"search_summary"`
		After         string      `json:"after"`
	} `json:"data"`
	Error interface{} `json:"error"`
}

type FilteredTransactions struct {
	UUID                 string    `json:"uuid"`
	Amount               float64   `json:"amount"`
	SourceAmount         float64   `json:"source_amount"`
	CurrentBalance       float64   `json:"current_balance"`
	TxnTimestamp         time.Time `json:"txn_timestamp"`
	Type                 string    `json:"type"`
	Source               string    `json:"source"`
	Account              string    `json:"account"`
	AccountID            string    `json:"account_id"`
	Merchant             string    `json:"merchant"`
	MerchantName         string    `json:"merchant_name"`
	MerchantType         string    `json:"merchant_type"`
	MerchantAddress      string    `json:"merchant_address"`
	Narration            string    `json:"narration"`
	Category             string    `json:"category"`
	CategoryID           string    `json:"category_id"`
	Subcategory          string    `json:"subcategory"`
	SubcategoryID        string    `json:"subcategory_id"`
	Tags                 string    `json:"tags"`
	Kind                 string    `json:"kind"`
	Mode                 string    `json:"mode"`
	Reference            string    `json:"reference"`
	Notes                string    `json:"notes"`
	ExcludedFromCashFlow bool      `json:"excluded_from_cash_flow"`
	IsBookmarked         bool      `json:"is_bookmarked"`
	Summary              string    `json:"summary"`
	TransactionID        string    `json:"transaction_id"`
	RefundStatus         string    `json:"refund_status"`
	RefundReceivedOn     string    `json:"refund_received_on"`
	BeforeFoldAccount    bool      `json:"before_fold_account"`
	Via                  string    `json:"via"`
	AccountIn            string    `json:"account_in"`
	ContactID            string    `json:"contact_id"`
	GroupIDs             string    `json:"group_ids"`
	Currency             string    `json:"currency"`
	SourceCurrency       string    `json:"source_currency"`
	UserManualAdded      bool      `json:"user_manual_added"`
	SplitType            string    `json:"split_type"`
	ParentTransactionID  string    `json:"parent_transaction_id"`
	F1PredictedCategory  bool      `json:"f1_predicted_category"`
	F1PredictedMerchant  bool      `json:"f1_predicted_merchant"`
}

type TransactionsReturn struct {
	Transactions []FilteredTransactions
}

func filterTransactions(raw TransactionsResponse, since time.Time) []FilteredTransactions {
	transactions := make([]FilteredTransactions, 0, len(raw.Data.Transactions))

	t := raw.Data.Transactions
	for i := 0; i < len(t); i++ {
		if t[i].TxnTimestamp.Before(since) {
			break
		}

		transaction := FilteredTransactions{
			UUID:                 t[i].UUID,
			Amount:               t[i].Amount,
			SourceAmount:         t[i].SourceAmount,
			Type:                 normalizeFoldAPIType(t[i].Type),
			Source:               t[i].Source,
			Account:              t[i].FinancialInformationProviderID,
			AccountID:            t[i].AccountID,
			Merchant:             t[i].Narration,
			Narration:            t[i].Narration,
			TxnTimestamp:         t[i].TxnTimestamp,
			CurrentBalance:       t[i].CurrentBalance,
			Kind:                 t[i].Kind,
			Mode:                 t[i].Mode,
			Reference:            t[i].Reference,
			ExcludedFromCashFlow: t[i].ExcludedFromCashFlow,
			IsBookmarked:         t[i].IsBookmarked,
			Summary:              t[i].Summary,
			TransactionID:        t[i].TransactionID,
			BeforeFoldAccount:    t[i].BeforeFoldAccount,
			Tags:                 jsonValue(t[i].Tags),
			Notes:                jsonString(t[i].Notes),
			Via:                  jsonValue(t[i].Via),
			AccountIn:            jsonValue(t[i].AccountIn),
			ContactID:            jsonString(t[i].ContactID),
			GroupIDs:             jsonValue(t[i].GroupIDs),
			Currency:             t[i].Currency,
			SourceCurrency:       t[i].SourceCurrency,
			UserManualAdded:      t[i].UserManualAdded != nil && *t[i].UserManualAdded,
			SplitType:            t[i].SplitType,
			ParentTransactionID:  t[i].ParentTransactionID,
		}

		if t[i].FinancialInformationProvider != nil && t[i].FinancialInformationProvider.Name != "" {
			transaction.Account = t[i].FinancialInformationProvider.Name
		}
		if t[i].Category != nil {
			transaction.CategoryID = t[i].Category.ID
			transaction.SubcategoryID = t[i].Category.SubcategoryID
		}
		if t[i].Merchant != nil {
			transaction.MerchantName = t[i].Merchant.Name
			transaction.MerchantType = t[i].Merchant.Type
			transaction.MerchantAddress = jsonValue(t[i].Merchant.Address)
			if transaction.MerchantName != "" {
				transaction.Merchant = transaction.MerchantName
			}
		}
		if t[i].Refund != nil {
			transaction.RefundStatus = t[i].Refund.Status
			transaction.RefundReceivedOn = jsonString(t[i].Refund.ReceivedOn)
		}

		var predicted map[string]bool
		if len(t[i].IsF1Predicted) > 0 && string(t[i].IsF1Predicted) != "null" {
			_ = json.Unmarshal(t[i].IsF1Predicted, &predicted)
			transaction.F1PredictedCategory = predicted["category_or_subcategory"]
			transaction.F1PredictedMerchant = predicted["merchant"]
		}

		transactions = append(transactions, transaction)
	}

	return transactions
}

func Transactions(uuid string, since time.Time, till time.Time) (TransactionsReturn, error) {
	RefreshOrFail()

	_ = till
	req, _ := APIRequest("GET", Url("/v3/users/"+uuid+"/transactions"), nil)

	q := req.URL.Query()
	q.Add("limit", "100")
	q.Add("count_by", "month")
	req.URL.RawQuery = q.Encode()

	resp, err := Client.Do(req)

	if err != nil {
		return TransactionsReturn{}, err
	} else {

		log.Debug().Msgf("Transactions response status: %+v", resp.StatusCode)

		if resp.StatusCode/100 != 2 {
			return TransactionsReturn{}, errors.New(resp.Status)
		}

		data := TransactionsResponse{}
		json.NewDecoder(resp.Body).Decode(&data)

		var ret TransactionsReturn
		ret.Transactions = make([]FilteredTransactions, 0)

		if len(data.Data.Transactions) == 0 {
			return ret, nil
		}

		log.Debug().Msgf("Transactions response body: %+v", data.Data.Transactions[0].TxnTimestamp)
		ret.Transactions = append(ret.Transactions, filterTransactions(data, since)...)

		for len(data.Data.Transactions) > 0 && data.Data.After != "" && data.Data.Transactions[len(data.Data.Transactions)-1].TxnTimestamp.After(since) {
			log.Debug().Msg("Fetching older transactions")

			log.Debug().Msg("New cursor base64: " + data.Data.After)
			q.Set("after", data.Data.After)
			req.URL.RawQuery = q.Encode()

			resp, err := Client.Do(req)
			if err != nil {
				log.Warn().Msg("Failed to fetch older transactions")
				break
			}

			log.Debug().Msgf("Transactions response status: %+v", resp.StatusCode)

			if resp.StatusCode/100 != 2 {
				log.Warn().Msgf("Failed to fetch older transactions, status code: %+v", resp.StatusCode)
				break
			}

			json.NewDecoder(resp.Body).Decode(&data)

			if len(data.Data.Transactions) == 0 {
				break
			}

			log.Debug().Msgf("Transactions response body: %+v", data.Data.Transactions[0].TxnTimestamp)
			ret.Transactions = append(ret.Transactions, filterTransactions(data, since)...)
		}

		return ret, nil
	}
}

func normalizeFoldAPIType(value string) string {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "OUTGOING":
		return "DEBIT"
	case "INCOMING":
		return "CREDIT"
	default:
		return strings.ToUpper(strings.TrimSpace(value))
	}
}

func jsonString(value json.RawMessage) string {
	if len(value) == 0 || string(value) == "null" {
		return ""
	}
	var text string
	if err := json.Unmarshal(value, &text); err == nil {
		return text
	}
	return ""
}

func jsonValue(value json.RawMessage) string {
	if len(value) == 0 || string(value) == "null" {
		return ""
	}
	if text := jsonString(value); text != "" {
		return text
	}
	return string(value)
}
