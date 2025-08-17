package system

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"log"
)

type Store struct {
	fishPriceMap map[string]int
	tryingToBuy  string
}

func NewStore(eHub *tasks.EventHub) *Store {
	s := &Store{}

	s.fishPriceMap = make(map[string]int)

	s.fishPriceMap["kirbensis"] = 2
	s.fishPriceMap["guppy"] = 1

	if eHub == nil {
		log.Fatal("is actually not init")
	}
	s.Subscribe(eHub)
	return s
}

func (s *Store) Subscribe(eHub *tasks.EventHub) {
	eHub.Subscribe(events.BuyAttempt{}, func(e tasks.Event) {
		ev := e.(events.BuyAttempt)
		s.tryingToBuy = ev.Item
		pev := events.MoneySpent{Amount: ev.Cost}
		eHub.Publish(pev)
	})

}
