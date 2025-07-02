package events

type MoneyAdded struct {
	Amount float32
}

type MoneySpent struct {
	Amount float32
}

type InsufficientFunds struct {
}

type PurchaseSuccessful struct {
}

type BuyAttempt struct {
	Cost float32
	Item string
	Name string
}

type NewPurchase struct {
	Purchase string
	Type     string
}
