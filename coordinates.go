package main

type position struct {
	x, y float64
}
type route struct {
	id        int64
	positions []position
}

func normalize(r []route, maxWidth, maxHeight int) []route {
	width, height := float64(maxWidth), float64(maxHeight)

	if len(r) == 0 || len(r[0].positions) == 0 {
		return r
	}

	topLeft, bottomRight := bounds(flattenRoute(r))
	offsetX := -topLeft.x
	offsetY := -topLeft.y
	scaleX := width / (bottomRight.x - topLeft.x)
	scaleY := height / (bottomRight.y - topLeft.y)

	// To keep the original dimensions, we just keep the lower scale
	if scaleX < scaleY {
		scaleY = scaleX
	} else {
		scaleX = scaleY
	}

	var newRoutes []route
	for _, oldRoute := range r {
		var newPositions []position
		for _, oldPositions := range oldRoute.positions {
			newPositions = append(newPositions, position{
				x: (oldPositions.x + offsetX) * scaleX,
				y: (oldPositions.y + offsetY) * scaleY,
			})
		}

		newRoutes = append(newRoutes, route{positions: newPositions})
	}

	return newRoutes
}

func bounds(positions []position) (topLeft, bottomRight position) {
	topLeft, bottomRight = positions[0], positions[0]

	for _, position := range positions {
		if position.x < topLeft.x {
			topLeft.x = position.x
		}

		if position.y < topLeft.y {
			topLeft.y = position.y
		}

		if position.x > bottomRight.x {
			bottomRight.x = position.x
		}

		if position.y > bottomRight.y {
			bottomRight.y = position.y
		}
	}

	return topLeft, bottomRight
}

func flattenRoute(r []route) []position {
	var completePositions []position
	for _, oldRoute := range r {
		completePositions = append(completePositions, oldRoute.positions...)
	}

	return completePositions
}

func filter(r []route, ids ...int64) []route {
	var newRoutes []route

skip:
	for _, route := range r {
		for _, id := range ids {
			if route.id == id {
				continue skip
			}
		}

		newRoutes = append(newRoutes, route)
	}

	return newRoutes
}
