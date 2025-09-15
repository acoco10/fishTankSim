package events

type MoneyAdded struct {
	Amount float64
}

func (m MoneyAdded) Type() string {
	return "MoneyAdded"
}

type MoneySpent struct {
	Amount float64
}

func (m MoneySpent) Type() string {
	return "MoneySpent"
}

type InsufficientFunds struct {
}

func (i InsufficientFunds) Type() string {
	return "InsufficientFunds"
}

type PurchaseSuccessful struct {
	Purchase     string
	PurchaseType uint8
}

func (p PurchaseSuccessful) Type() string {
	return "PurchaseSuccessful"
}

type BuyAttempt struct {
	Cost     float64
	Item     string
	ItemType uint8
}

func (b BuyAttempt) Type() string {
	return "BuyAttempt"
}

type NewPurchase struct {
	Purchase     string
	PurchaseType string
}

func (n NewPurchase) Type() string {
	return "NewPurchase"
}
