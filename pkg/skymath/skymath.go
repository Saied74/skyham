package skymath

import "math"

//Euler is a 3x3 matrix for 3D Euler rotation
type Euler [3][3]float64

//Vec is the 3D position vector for positioning planets
type Vec [3]float64

//E1 is Euler rotation around the x axis
func E1(a float64) Euler {
	e := Euler{}
	a = (a / 180.0) * math.Pi
	s := math.Sin(a)
	c := math.Cos(a)
	e[0] = [3]float64{1.0, 0.0, 0.0}
	e[1] = [3]float64{0.0, c, -s}
	e[2] = [3]float64{0.0, s, c}
	return e
}

//E2 is Euler rotation aroudn the y axis
func E2(a float64) Euler {
	e := Euler{}
	a = (a / 180.0) * math.Pi
	s := math.Sin(a)
	c := math.Cos(a)
	e[0] = [3]float64{c, 0.0, s}
	e[1] = [3]float64{0.0, 1.0, 0.0}
	e[2] = [3]float64{-s, 0.0, c}
	return e
}

//E3 is Euler rotation around the z axis
func E3(a float64) Euler {
	e := Euler{}
	a = (a / 180.0) * math.Pi
	s := math.Sin(a)
	c := math.Cos(a)
	e[0] = [3]float64{c, -s, 0}
	e[1] = [3]float64{s, c, 0}
	e[2] = [3]float64{0.0, 0.0, 1.0}
	return e
}

//Mply multiplies two 3x3 matrices
func Mply(a, b Euler) Euler {
	c := Euler{}
	for m := 0; m < 3; m++ {
		for n := 0; n < 3; n++ {
			for i := 0; i < 3; i++ {
				c[m][n] += a[m][i] * b[i][n]
			}
		}
	}
	return c
}

//Vply multiplies a matrix (first argument) with a vector
func Vply(a Euler, b Vec) Vec {
	c := Vec{}
	for m := 0; m < 3; m++ {
		for i := 0; i < 3; i++ {
			c[m] += a[m][i] * b[i]
		}
	}
	return c
}
