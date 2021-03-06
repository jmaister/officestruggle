package gamestate

// https://jeremyceri.se/post/sidenote-1-cameras/

type Camera struct {
	// Keep track of where the camera currently starts its viewport from
	X int
	Y int
	// Its width and height
	Width  int
	Height int
}

// ‘targetX’ and ‘targetY’ here are the current position of the player
func (c *Camera) MoveCamera(targetX int, targetY int, mapWidth int, mapHeight int) {
	// Update the camera coordinates to the target coordinates
	x := targetX - c.Width/2
	y := targetY - c.Height/2

	if x < 0 {
		x = 0
	}

	if y < 0 {
		y = 0
	}

	if x > mapWidth-c.Width {
		x = mapWidth - c.Width
	}

	if y > mapHeight-c.Height {
		y = mapHeight - c.Height
	}

	c.X, c.Y = x, y
}

// This function takes a set of map coordinates, and translates them to camera coordinates.
func (c *Camera) ToCameraCoordinates(mapX int, mapY int) (cameraX int, cameraY int) {
	// Convert coordinates on the map, to coordinates on the viewport
	x, y := mapX-c.X, mapY-c.Y

	if x < 0 || y < 0 || x >= c.Width || y >= c.Height {
		return -1, -1
	}

	return x, y
}
