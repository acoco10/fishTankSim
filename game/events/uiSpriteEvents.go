package events

type MoneyAvailable struct {
	CurrentAmount float64
	Amount        float64
}

func (m MoneyAvailable) Type() string {
	return "MoneyAvailable"
}
