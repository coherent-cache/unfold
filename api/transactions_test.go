package api

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestTransactionsUsesResponseCursorAndParsesV3Fields(t *testing.T) {
	t.Parallel()

	oldClient := Client
	t.Cleanup(func() { Client = oldClient })

	requests := []string{}
	Client = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			requests = append(requests, req.URL.Path+"?"+req.URL.RawQuery)
			switch req.URL.Path {
			case "/api/v1/auth/tokens/refresh":
				return jsonResponse(`{"data":{"access_token":"new-access","refresh_token":"new-refresh"}}`), nil
			case "/api/v3/users/user-123/transactions":
				after := req.URL.Query().Get("after")
				if after == "" {
					return jsonResponse(`{"data":{"transactions":[{"uuid":"txn-1","amount":421,"source_amount":421,"current_balance":0,"txn_timestamp":"2026-03-10T09:00:00Z","mode":"CARD","type":"OUTGOING","source":"CREDIT_CARD","narration":"CARD/abcdef0123456789/ACME FOODS LTD/INR/421.00/OUTGOING/Mar 10 2026 at 09:00:00","merchant":{"name":"ACME FOODS LTD","type":"f1","address":null},"category":{"id":"cat-food","subcategory_id":"sub-takeaway"},"account_id":"cc-account-1","financial_information_provider_id":"","tags":[],"kind":"NORMAL","notes":null,"excluded_from_cash_flow":false,"is_bookmarked":false,"transaction_id":"txn-1","reference":"","summary":"You spent ₹421","before_fold_account":false,"via":null,"account_in":null,"refund":{"status":"NONE","received_on":null},"group_ids":null,"contact_id":null,"is_f1_predicted":{"category_or_subcategory":true,"merchant":true},"currency":"INR","source_currency":"INR","user_manual_added":null,"split_type":"","parent_transaction_id":""}],"after":"cursor-1"}}`), nil
				}
				require.Equal(t, "cursor-1", after)
				return jsonResponse(`{"data":{"transactions":[{"uuid":"txn-2","amount":15000,"source_amount":15000,"current_balance":0,"txn_timestamp":"2026-03-05T06:00:00Z","mode":"CASH","type":"OUTGOING","source":"MANUAL","narration":"CASH/Landlord/15000.00/DEBIT/2026-03-05/2026-03-05 06:00:00","merchant":{"name":"Landlord","type":"custom","address":null},"category":{"id":"cat-bill","subcategory_id":"sub-rent"},"account_id":"cash-account-1","financial_information_provider_id":"","tags":[],"kind":"NORMAL","notes":null,"excluded_from_cash_flow":false,"is_bookmarked":false,"transaction_id":"txn-2","reference":"","summary":"You transferred ₹15K","before_fold_account":false,"via":null,"account_in":null,"refund":{"status":"NONE","received_on":null},"group_ids":null,"contact_id":null,"is_f1_predicted":{"category_or_subcategory":false,"merchant":false},"currency":"INR","source_currency":"INR","user_manual_added":true,"split_type":"","parent_transaction_id":""}],"after":""}}`), nil
			default:
				t.Fatalf("unexpected path: %s", req.URL.Path)
				return nil, nil
			}
		}),
	}

	viper.Set("device_hash", "device-123")
	viper.Set("token.access", "expired-access")
	viper.Set("token.refresh", "refresh-123")

	since := time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC)
	result, err := Transactions("user-123", since, since.AddDate(0, 0, 10))
	require.NoError(t, err)
	require.Len(t, result.Transactions, 2)

	first := result.Transactions[0]
	require.Equal(t, "DEBIT", first.Type)
	require.Equal(t, "CREDIT_CARD", first.Source)
	require.Equal(t, "ACME FOODS LTD", first.Merchant)
	require.Equal(t, "ACME FOODS LTD", first.MerchantName)
	require.Equal(t, "f1", first.MerchantType)
	require.Equal(t, "cat-food", first.CategoryID)
	require.Equal(t, "sub-takeaway", first.SubcategoryID)
	require.Equal(t, "NONE", first.RefundStatus)
	require.True(t, first.F1PredictedCategory)

	second := result.Transactions[1]
	require.Equal(t, "MANUAL", second.Source)
	require.True(t, second.UserManualAdded)

	require.Len(t, requests, 3)
	require.Equal(t, "/api/v1/auth/tokens/refresh?", requests[0])
	require.Equal(t, "/api/v3/users/user-123/transactions?count_by=month&limit=100", requests[1])
	require.Equal(t, "/api/v3/users/user-123/transactions?after=cursor-1&count_by=month&limit=100", requests[2])
}

func jsonResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}
