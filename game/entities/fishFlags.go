package entities

type FishFlags uint32

const (
	// Physics/Movement flags
	FlagIncreasingZ = 1 << iota
	FlagDecreasingZ
	FlagSwimmingDirection
	FlagChangedDirection
	HitEdgeBoundary //
	FlagFacingDirection
)

func (c *CreatureData) IsGoingRight() bool { return c.Flags&FlagSwimmingDirection != 0 }

func (c *CreatureData) IsGoingLeft() bool { return c.Flags&FlagSwimmingDirection == 0 }

func (c *CreatureData) IsFacingLeft() bool { return c.Flags&FlagFacingDirection == 0 }

func (c *CreatureData) IsFacingRight() bool { return c.Flags&FlagFacingDirection != 0 }

func (c *CreatureData) CheckRecentlyChangedDirection() bool {
	return c.Flags&FlagChangedDirection != 0
}

func (c *CreatureData) CheckHitBoundary() bool {
	return c.Flags&HitEdgeBoundary != 0
}

func (c *CreatureData) SetRecentlyChangedDirection() {
	c.Flags |= FlagChangedDirection
}

func (c *CreatureData) ClearRecentlyChangedDirection() {
	c.Flags &^= FlagChangedDirection
}

func (c *CreatureData) IsIncreasingZ() bool {
	return c.Flags&FlagIncreasingZ != 0
}

func (c *CreatureData) IsDecreasingZ() bool {
	return c.Flags&FlagDecreasingZ != 0
}

func (c *CreatureData) SetHitBoundary() {
	c.Flags |= HitEdgeBoundary
}

func (c *CreatureData) ClearHitBoundary() {
	c.Flags &^= HitEdgeBoundary
}

func (c *CreatureData) SetDirection(direction Direction) {
	if direction == Left {
		c.Flags &^= FlagSwimmingDirection
	}

	if direction == Right {
		c.Flags |= FlagSwimmingDirection
	}
	
}

func (c *CreatureData) SetFacingDirection(direction Direction) {
	if direction == Left {
		c.Flags &^= FlagFacingDirection
	}
	if direction == Right {
		c.Flags |= FlagFacingDirection
	}
}

func (c *CreatureData) SetIncreasingZ() {
	c.Flags &^= FlagDecreasingZ
	c.Flags |= FlagIncreasingZ
}
func (c *CreatureData) ClearZChange() {
	c.Flags &^= FlagIncreasingZ
	c.Flags &^= FlagDecreasingZ
}

func (c *CreatureData) SetDecreasingZ() {
	c.Flags &^= FlagIncreasingZ
	c.Flags |= FlagDecreasingZ
}
