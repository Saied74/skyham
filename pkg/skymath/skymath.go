package skymath

import (
	"fmt"
	"math"
)

//Euler is a 3x3 matrix for 3D Euler rotation
type Euler [3][3]float64

//Vec is the 3D position vector for positioning planets
type Vec [3]float64

//E1 is Euler rotation around the x axis
func E1(a float64) Euler {
	e := Euler{}
	a = ToRadians(a)
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
	a = ToRadians(a)
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
	a = ToRadians(a)
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

//Vadd is vector addition
func Vadd(a, b Vec) Vec {
	return Vec{
		a[0] + b[0],
		a[1] + b[1],
		a[2] + b[2],
	}
}

//Vsub is vector subtraction
func Vsub(a, b Vec) Vec {
	return Vec{
		a[0] - b[0],
		a[1] - b[1],
		a[2] - b[2],
	}
}

//CalcBetaEpsilon calculates bearing and elevation of the location
func CalcBetaEpsilon(v Vec) (beta, epsilon float64) {
	x := v[0]
	y := v[1]
	z := v[2]

	d := math.Sqrt(x*x + y*y)
	sBeta := x / d
	cBeta := y / d
	fmt.Println(sBeta)
	fmt.Println(cBeta)
	if sBeta > 0.0 && cBeta > 0.0 {
		beta = math.Asin(sBeta)
	}
	if sBeta > 0.0 && cBeta < 0.0 {

		beta = math.Pi/2 + math.Asin(-cBeta)
	}
	if sBeta < 0.0 && cBeta < 0.0 {
		beta = math.Pi + math.Asin(-sBeta)
	}
	if sBeta < 0.0 && cBeta > 0.0 {
		beta = math.Pi + math.Pi/2 + math.Asin(cBeta)
	}
	tEpsilon := z / d
	// beta = math.Asin(sBeta)
	beta = ToDegrees(beta)
	// fmt.Printf("Raw Beta: %f\n", beta)
	// if beta < 0.0 {
	// 	beta += 360.0
	// }
	epsilon = math.Atan(tEpsilon)
	epsilon = ToDegrees(epsilon)
	return beta, epsilon
}

//ToRadians converts degrees to radians modulo 2Pi
func ToRadians(d float64) float64 {
	return math.Mod((d/180)*math.Pi, 2.0*math.Pi)
}

//ToDegrees converts radians to degrees modulo 360
func ToDegrees(r float64) float64 {
	return math.Mod((r/math.Pi)*180, 360.0)
}
