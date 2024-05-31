package main

import "math"

const (
	h  = 0.001 // the change in runge kutta
	L1 = 1.0   // length 1
	L2 = 1.0   // length 2
	m1 = 1.0   // mass of blob 1
	m2 = 1.0   // mass of blob 2
	g  = 9.81  // gravity
)

// dw/dt function for theta 1
func funcdwdt1(theta1, theta2, w1, w2 float64) float64 {
	cos12 := math.Cos(theta1 - theta2)
	sin12 := math.Sin(theta1 - theta2)
	sin1 := math.Sin(theta1)
	sin2 := math.Sin(theta2)
	denom := math.Pow(cos12, 2)*m2 - m1 - m2
	ans := (L1*m2*cos12*sin12*math.Pow(w1, 2) + L2*m2*sin12*math.Pow(w2, 2) - m2*g*cos12*sin2 + (m1+m2)*g*sin1) / (L1 * denom)
	return ans
}

// dw/dt function for theta 2
func funcdwdt2(theta2, theta1, w1, w2 float64) float64 {
	cos12 := math.Cos(theta1 - theta2)
	sin12 := math.Sin(theta1 - theta2)
	sin1 := math.Sin(theta1)
	sin2 := math.Sin(theta2)
	denom := math.Pow(cos12, 2)*m2 - m1 - m2
	ans2 := -(L2*m2*cos12*sin12*math.Pow(w2, 2) + L1*(m1+m2)*sin12*math.Pow(w1, 2) + (m1+m2)*g*sin1*cos12 - (m1+m2)*g*sin2) / (L2 * denom)
	return ans2
}

// d0/dt function for theta 1
func funcd0dt1(w0 float64) float64 {
	return w0
}

// d0/dt function for theta 2
func funcd0dt2(w0 float64) float64 {
	return w0
}

func step(w1, w2, theta1, theta2 float64) (float64, float64, float64, float64) {
	k1a := h * funcd0dt1(w1)                     // gives theta1
	k1b := h * funcdwdt1(theta1, theta2, w1, w2) // gives omega1
	k1c := h * funcd0dt2(w2)                     // gives theta2
	k1d := h * funcdwdt2(theta2, theta1, w1, w2) // gives omega2

	k2a := h * funcd0dt1(w1+(0.5*k1b))
	k2b := h * funcdwdt1(theta1+(0.5*k1a), theta2, w1, w2)
	k2c := h * funcd0dt2(w2+(0.5*k1d))
	k2d := h * funcdwdt2(theta2+(0.5*k1c), theta1, w1, w2)

	k3a := h * funcd0dt1(w1+(0.5*k2b))
	k3b := h * funcdwdt1(theta1+(0.5*k2a), theta2, w1, w2)
	k3c := h * funcd0dt2(w2+(0.5*k2d))
	k3d := h * funcdwdt2(theta2+(0.5*k2c), theta1, w1, w2)

	k4a := h * funcd0dt1(w1+k3b)
	k4b := h * funcdwdt1(theta1+k3a, theta2, w1, w2)
	k4c := h * funcd0dt2(w2+k3d)
	k4d := h * funcdwdt2(theta2+k3c, theta1, w1, w2)

	//summing the values after the iterations
	theta1 += 1.0 / 6.0 * (k1a + 2*k2a + 2*k3a + k4a)
	w1 += 1.0 / 6.0 * (k1b + 2*k2b + 2*k3b + k4b)
	theta2 += 1.0 / 6.0 * (k1c + 2*k2c + 2*k3c + k4c)
	w2 += 1.0 / 6.0 * (k1d + 2*k2d + 2*k3d + k4d)
	return w1, w2, theta1, theta2
}

func GetPos(w1, w2, theta1, theta2, scale float64) (x1, y1, x2, y2 float64) {
	x1 = L1 * math.Sin(theta1)
	y1 = -L1 * math.Cos(theta1)
	x2 = x1 + L2*math.Sin(theta2)
	y2 = y1 - L2*math.Cos(theta2)

	x1 *= scale
	x2 *= scale
	y1 *= scale
	y2 *= scale

	return x1, y1, x2, y2
}

/*
// https://en.wikipedia.org/wiki/Runge%E2%80%93Kutta_methods
// Runge Kutta 4th order
// ended up that I can't use this because the f(t,y) for this has 4 parts to y...
func RK4(t, y float64, stepSize float64, firstDerivative func(t, y float64) float64) float64 {
	h := stepSize
	f := firstDerivative

	//Runge Kutta standard calculations
	k1 := f(t, y)
	k2 := f(t+h/2, y+h/2*k1)
	k3 := f(t+h/2, y+h/2*k2)
	k4 := f(t+h, y+h*k3)

	return y + h/6.0*(k1+2.0*k2+2.0*k3+k4)
}
*/
