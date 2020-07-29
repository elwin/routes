package main

type position struct {
	x, y float64
}
type route struct {
	positions []position
}

func normalize(r []route, maxWidth, maxHeight float64) []route {
	if len(r) == 0 || len(r[0].positions) == 0 {
		return r
	}

	topLeft, bottomRight := bounds(flattenRoute(r))
	offsetX := - topLeft.x
	offsetY := -topLeft.y
	scaleX := maxWidth / (bottomRight.x - topLeft.x)
	scaleY := maxHeight / (bottomRight.y - topLeft.y)

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

	for _, point := range positions {
		if point.x < topLeft.x {
			topLeft.x = point.x
		}

		if point.y < topLeft.y {
			topLeft.y = point.y
		}

		if point.x > bottomRight.x {
			bottomRight.x = point.x
		}

		if point.y > bottomRight.y {
			bottomRight.y = point.y
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