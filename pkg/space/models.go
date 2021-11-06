package space

import "space-game/pkg/maths"

type CelestialBody struct {
	Name                     string
	Mass                     float64
	internalAbsoluteLocation maths.Point3
	internalAbsoluteVelocity maths.Vector3
}

type Ship struct {
	Fuel        float64
	DeltaV      float64
	Orientation maths.Vector3
	Velocity    maths.Vector3
	Distance    maths.Point3
}
