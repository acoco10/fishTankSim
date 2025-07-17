package events

type MoneyAdded struct {
	Amount float64
}

type MoneySpent struct {
	Amount float64
}

type InsufficientFunds struct {
}

type PurchaseSuccessful struct {
}

type BuyAttempt struct {
	Cost float64
	Item string
	Name string
}

type NewPurchase struct {
	Purchase string
	Type     string
}
