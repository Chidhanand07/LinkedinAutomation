package stealth

import (
	"math"
	"math/rand"
	"time"

	"github.com/go-rod/rod"
)

// HumanMove simulates human-like mouse movement using Bézier curves
// with variable speed, natural overshoot, and micro-corrections
func HumanMove(page *rod.Page, targetX, targetY float64) {
	// Get current position (or start from random corner)
	startX := rand.Float64() * 100
	startY := rand.Float64() * 100

	// Generate Bézier control points for natural curve
	cp1X := startX + rand.Float64()*(targetX-startX)*0.3 + (rand.Float64()-0.5)*100
	cp1Y := startY + rand.Float64()*(targetY-startY)*0.3 + (rand.Float64()-0.5)*100
	cp2X := startX + rand.Float64()*(targetX-startX)*0.7 + (rand.Float64()-0.5)*100
	cp2Y := startY + rand.Float64()*(targetY-startY)*0.7 + (rand.Float64()-0.5)*100

	// Natural overshoot (occasionally go past target slightly)
	overshootX := targetX
	overshootY := targetY
	if rand.Float64() > 0.7 {
		overshootX += (rand.Float64() - 0.5) * 20
		overshootY += (rand.Float64() - 0.5) * 20
	}

	// Variable speed with realistic acceleration/deceleration
	steps := rand.Intn(30) + 40
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)

		// Ease in-out for realistic acceleration
		eased := easeInOutCubic(t)

		// Cubic Bézier curve calculation
		nx := math.Pow(1-eased, 3)*startX +
			3*math.Pow(1-eased, 2)*eased*cp1X +
			3*(1-eased)*math.Pow(eased, 2)*cp2X +
			math.Pow(eased, 3)*overshootX

		ny := math.Pow(1-eased, 3)*startY +
			3*math.Pow(1-eased, 2)*eased*cp1Y +
			3*(1-eased)*math.Pow(eased, 2)*cp2Y +
			math.Pow(eased, 3)*overshootY

		// Add slight jitter for natural imperfection
		nx += (rand.Float64() - 0.5) * 2
		ny += (rand.Float64() - 0.5) * 2

		page.Mouse.MustMoveTo(nx, ny)

		// Variable timing between movements
		delay := rand.Intn(8) + 3
		if i < 5 || i > steps-5 {
			delay += rand.Intn(10) // Slower at start/end
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	// Micro-corrections: small adjustments after reaching target
	if rand.Float64() > 0.6 {
		for j := 0; j < rand.Intn(3)+1; j++ {
			corrX := targetX + (rand.Float64()-0.5)*6
			corrY := targetY + (rand.Float64()-0.5)*6
			page.Mouse.MustMoveTo(corrX, corrY)
			time.Sleep(time.Duration(rand.Intn(30)+20) * time.Millisecond)
		}
	}

	// Final precise position
	page.Mouse.MustMoveTo(targetX, targetY)
}

// HumanHover simulates natural hovering behavior
func HumanHover(page *rod.Page, x, y float64, duration time.Duration) {
	startTime := time.Now()
	for time.Since(startTime) < duration {
		// Slight drift while hovering
		driftX := x + (rand.Float64()-0.5)*10
		driftY := y + (rand.Float64()-0.5)*10
		page.Mouse.MustMoveTo(driftX, driftY)
		time.Sleep(time.Duration(rand.Intn(150)+100) * time.Millisecond)
	}
}

// RandomMouseWander simulates idle cursor movement
func RandomMouseWander(page *rod.Page) {
	if rand.Float64() > 0.4 { // 60% chance to wander
		offsetX := (rand.Float64() - 0.5) * 200
		offsetY := (rand.Float64() - 0.5) * 200
		page.Mouse.MustMoveTo(offsetX, offsetY)
		time.Sleep(time.Duration(rand.Intn(500)+300) * time.Millisecond)
	}
}

// easeInOutCubic provides smooth acceleration/deceleration
func easeInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	return 1 - math.Pow(-2*t+2, 3)/2
}
