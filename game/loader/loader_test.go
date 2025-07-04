package loader

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"strconv"
	"testing"
)

func TestLoadAnimation(t *testing.T) {
	savedFish := entities.SavedFish{FishType: "mollyFish", Size: 1}

	path := "data/animationData/" + savedFish.FishType + "%dAnimation.json"
	path = fmt.Sprintf(path, savedFish.Size)

	_, _, err := LoadAnimation(path)
	if err != nil {
		t.Logf("animation loader error: %q, tested path = %s", err, path)
		t.Fatal()
	}

}

func TestLoadCreature(t *testing.T) {

	res, err := LoadFishSprite(entities.Fish, 1)

	if err != nil {

		t.Log("creature sprite loader error", err)
		t.Fatal()
	}

	t.Log(strconv.Itoa(res.Animation.FirstF))
}
