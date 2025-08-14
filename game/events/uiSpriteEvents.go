package events

type MoneyAvailable struct {
	Amount float64
}

func (m MoneyAvailable) Type() string {
	return "MoneyAvailable"
}
