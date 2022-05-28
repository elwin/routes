package strava

type Position struct {
	X, Y float64
}

type Route struct {
	Id        int64
	Positions []Position
}

func normalize(r Route, width, height float64) Route {
	topLeft, bottomRight := bounds(r.Positions)
	originalWidth := bottomRight.X - topLeft.X
	scaleX := width / originalWidth
	originalHeight := bottomRight.Y - topLeft.Y
	scaleY := height / originalHeight

	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}

	offsetX := -topLeft.X
	xxx := (width - originalWidth*scale) / 2
	offsetY := -topLeft.Y
	yyyy := (height - originalHeight*scale) / 2


	var scaledRoute Route
	for _, position := range r.Positions {
		scaledRoute.Positions = append(scaledRoute.Positions, Position{
			X: (position.X+offsetX)*scale + xxx,
			Y: height - (position.Y+offsetY)*scale - yyyy,
		})
	}

	return scaledRoute
}

func bounds(positions []Position) (topLeft, bottomRight Position) {
	topLeft, bottomRight = positions[0], positions[0]

	for _, position := range positions {
		if position.X < topLeft.X {
			topLeft.X = position.X
		}

		if position.Y < topLeft.Y {
			topLeft.Y = position.Y
		}

		if position.X > bottomRight.X {
			bottomRight.X = position.X
		}

		if position.Y > bottomRight.Y {
			bottomRight.Y = position.Y
		}
	}

	return topLeft, bottomRight
}

func flattenRoute(r []Route) []Position {
	var completePositions []Position
	for _, oldRoute := range r {
		completePositions = append(completePositions, oldRoute.Positions...)
	}

	return completePositions
}

func filter(r []Route, ids ...int64) []Route {
	var newRoutes []Route

skip:
	for _, route := range r {
		for _, id := range ids {
			if route.Id == id {
				continue skip
			}
		}

		newRoutes = append(newRoutes, route)
	}

	return newRoutes
}
