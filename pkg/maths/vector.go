package maths

import "math"

type Point3 [3]float64
type Vector3 [3]float64

func ScalarProduct(vector Vector3, scalar float64) Vector3 {
	return Vector3{vector[0] * scalar, vector[1] * scalar, vector[2] * scalar}
}

func Magnitude(vector Vector3) float64 {
	return math.Sqrt(math.Pow(vector[0], 2) + math.Pow(vector[1], 2) + math.Pow(vector[2], 2))
}

func Normalize(vector Vector3) Vector3 {
	magnitude := Magnitude(vector)
	return ScalarProduct(vector, 1/magnitude)
}

func Subtract(vectorA, vectorB Vector3) Vector3 {
	return Vector3{vectorA[0] - vectorB[0], vectorA[1] - vectorB[1], vectorA[2] - vectorB[2]}
}

func Add(vectorA, vectorB Vector3) Vector3 {
	return Vector3{vectorA[0] + vectorB[0], vectorA[1] + vectorB[1], vectorA[2] + vectorB[2]}
}
