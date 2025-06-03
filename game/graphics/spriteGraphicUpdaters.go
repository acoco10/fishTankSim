package graphics

func UpdateWhiteBoardCloth(gs *SpriteGraphic) {
	maxPoint := gs.parameters["max"].([2]float32)
	origin := gs.parameters["origin"].([2]float32)
	direction := gs.parameters["direction"].(string)

	if gs.Sprite.X >= maxPoint[0] && direction == "right" {
		gs.parameters["direction"] = "left"
	}

	if gs.Sprite.X <= origin[0] && direction == "left" {
		gs.parameters["direction"] = "right"
	}

	if gs.Sprite.Y >= maxPoint[1] {
		gs.parameters["direction"] = "stop"
		gs.complete = true
	}

	switch direction {
	case "right":
		gs.Sprite.X += 10
		gs.Sprite.Y++
	case "left":
		gs.Sprite.X -= 10
		gs.Sprite.Y++
	}

}
