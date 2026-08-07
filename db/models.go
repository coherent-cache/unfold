package db

import (
	"time"
)

type Transactions struct {
	UUID                 string `gorm:"primaryKey;"`
	Amount               float64
	SourceAmount         float64
	CurrentBalance       float64
	Timestamp            time.Time `gorm:"index:sortTimestamp,sort:desc"`
	Type                 string
	Source               string
	Account              string
	Merchant             string
	MerchantName         string
	MerchantType         string
	MerchantAddress      string
	Narration            string
	Category             string
	CategoryID           string
	Subcategory          string
	SubcategoryID        string
	Tags                 string
	Kind                 string
	Mode                 string
	Reference            string
	Notes                string
	ExcludedFromCashFlow bool
	IsBookmarked         bool
	Summary              string
	TransactionID        string
	RefundStatus         string
	RefundReceivedOn     string
	BeforeFoldAccount    bool
	Via                  string
	AccountIn            string
	ContactID            string
	GroupIDs             string
	Currency             string
	SourceCurrency       string
	UserManualAdded      bool
	SplitType            string
	ParentTransactionID  string
	F1PredictedCategory  bool
	F1PredictedMerchant  bool
	AccountID            string
}
