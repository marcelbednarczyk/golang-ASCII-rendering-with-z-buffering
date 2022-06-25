package main

import (
	"fmt"
	"math"
	"os"
	"time"
)

func main() {
	for i := 0.0; i < 5; i += 0.05 {
		renderFrame(i, 1.0)
		time.Sleep(1 * time.Millisecond)
	}
}

const screenWidth = 30
const screenHeight = 30

const thetaSpacing = 0.07
const phiSpacing = 0.02
const R1 = 1
const R2 = 2
const K2 = 5

// Calculate K1 based on screen size: the maximum x-distance occurs
// roughly at the edge of the torus, which is at x=R1+R2, z=0.  we
// want that to be displaced 3/8ths of the width of the screen, which
// is 3/4th of the way from the center to the side of the screen.
// screenWidth*3/8 = K1*(R1+R2)/(K2+0)
// screenWidth*K2*3/(8*(R1+R2)) = K1
const K1 = screenWidth * K2 * 3 / (8 * (R1 + R2))

func getLimenanceChar(index int) string {
	return []string{"..", ",,", "--", "~~", "::", ";;", "==", "!!", "**", "##", "$$", "@@"}[index]
}

func renderFrame(A float64, B float64) {
	// precompute sines and cosines of A and B
	cosA := cos(A)
	sinA := sin(A)
	cosB := cos(B)
	sinB := sin(B)
	pi := math.Pi

	output := [screenWidth][screenHeight]string{}
	zbuffer := [screenWidth][screenHeight]float64{}
	for i := range output {
		for j := range output[i] {
			output[i][j] = "  "
		}
	}

	// theta goes around the cross-sectional circle of a torus
	for theta := 0.0; theta < 2*pi; theta += thetaSpacing {
		// precompute sines and cosines of theta
		costheta := cos(theta)
		sintheta := sin(theta)

		// phi goes around the center of revolution of a torus
		for phi := 0.0; phi < 2*pi; phi += phiSpacing {
			// precompute sines and cosines of phi
			cosphi := cos(phi)
			sinphi := sin(phi)

			// the x,y coordinate of the circle, before revolving (factored
			// out of the above equations)
			circlex := R2 + R1*costheta
			circley := R1 * sintheta

			// final 3D (x,y,z) coordinate after rotations, directly from
			// our math above
			x := circlex*(cosB*cosphi+sinA*sinB*sinphi) - circley*cosA*sinB
			y := circlex*(sinB*cosphi-sinA*cosB*sinphi) + circley*cosA*cosB
			z := K2 + cosA*circlex*sinphi + circley*sinA
			ooz := 1 / z // "one over z"

			// x and y projection.  note that y is negated here, because y
			// goes up in 3D space but down on 2D displays.
			xp := (int)(screenWidth/2 + K1*ooz*x)
			yp := (int)(screenHeight/2 - K1*ooz*y)

			// calculate luminance.  ugly, but correct.
			L := cosphi*costheta*sinB - cosA*costheta*sinphi - sinA*sintheta + cosB*(cosA*sintheta-costheta*sinA*sinphi)
			// L ranges from -sqrt(2) to +sqrt(2).  If it's < 0, the surface
			// is pointing away from us, so we won't bother trying to plot it.
			if L > 0 {
				// test against the z-buffer.  larger 1/z means the pixel is
				// closer to the viewer than what's already plotted.
				if ooz > zbuffer[xp][yp] {
					zbuffer[xp][yp] = ooz
					luminanceIndex := (int)(L * 8)
					// luminanceIndex is now in the range 0..11 (8*sqrt(2) = 11.3)
					// now we lookup the character corresponding to the
					// luminance and plot it in our output:
					output[xp][yp] = getLimenanceChar(luminanceIndex)
				}
			}
		}
	}

	// now, dump output[] to the screen.
	// bring cursor to "home" location, in just about any currently-used
	// terminal emulation mode
	fmt.Fprint(os.Stdout, "\033[H\033[J")
	for j := 0; j < screenHeight; j++ {
		for i := 0; i < screenWidth; i++ {
			fmt.Print(output[i][j])
		}
		fmt.Println()
	}
}

func cos(a float64) float64 {
	return math.Cos(a)
}

func sin(a float64) float64 {
	return math.Sin(a)
}
