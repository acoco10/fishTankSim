package system

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/tasks"
)

type Store struct {
	fishPriceMap map[string]int
	tryingToBuy  string
}

func NewStore(eHub *tasks.EventHub) Store {
	s := Store{}

	fishPriceMap := make(map[string]int)

	fishPriceMap["kirbensis"] = 2
	fishPriceMap["guppy"] = 1

	s.fishPriceMap = fishPriceMap

	s.Subscribe(eHub)
	return s
}

func (s *Store) Subscribe(eHub *tasks.EventHub) {

	eHub.Subscribe(events.BuyAttempt{}, func(e tasks.Event) {
		ev := e.(events.BuyAttempt)
		s.tryingToBuy = ev.Name
		pev := events.MoneySpent{Amount: ev.Cost}
		eHub.Publish(pev)
	})

	eHub.Subscribe(events.PurchaseSuccessful{}, func(e tasks.Event) {
		ev := events.NewPurchase{Purchase: s.tryingToBuy, Type: "Fish"}
		eHub.Publish(ev)
	})

}
