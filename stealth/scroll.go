package stealth

import (
	"math/rand"
	"time"

	"github.com/go-rod/rod"
)

// HumanScroll simulates natural scrolling with variable speeds,
// acceleration/deceleration, and occasional scroll-back movements
func HumanScroll(page *rod.Page) {
	// Determine scroll direction and amount
	scrollDistance := rand.Intn(500) + 200
	direction := 1.0 // Downward scroll

	// Occasionally scroll back up (reading back)
	if rand.Float64() < 0.15 {
		direction = -1.0
		scrollDistance = rand.Intn(200) + 100
	}

	// Break scroll into smaller chunks for natural acceleration
	chunks := rand.Intn(4) + 3
	for i := 0; i < chunks; i++ {
		chunkSize := float64(scrollDistance) / float64(chunks)

		// Variable speed per chunk
		if i == 0 || i == chunks-1 {
			chunkSize *= 0.7 // Slower at start/end
		}

		page.Mouse.Scroll(0, direction*chunkSize, 1)

		// Variable pause between chunks
		delay := rand.Intn(150) + 80
		if i == chunks-1 {
			delay += rand.Intn(300) + 200 // Longer pause at end
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	// Occasional micro-scroll adjustments
	if rand.Float64() < 0.3 {
		time.Sleep(time.Duration(rand.Intn(400)+300) * time.Millisecond)
		page.Mouse.Scroll(0, direction*float64(rand.Intn(50)+20), 1)
	}
}

// ScrollToElement scrolls an element into view naturally
func ScrollToElement(page *rod.Page, targetY float64) {
	currentScroll := page.MustEval("window.pageYOffset").Int()
	distance := int(targetY) - currentScroll

	if distance == 0 {
		return
	}

	direction := 1.0
	if distance < 0 {
		direction = -1.0
		distance = -distance
	}

	// Scroll in natural increments
	steps := rand.Intn(5) + 3
	for i := 0; i < steps; i++ {
		stepSize := float64(distance) / float64(steps)

		// Ease in-out
		if i == 0 || i == steps-1 {
			stepSize *= 0.8
		}

		page.Mouse.Scroll(0, direction*stepSize, 1)
		time.Sleep(time.Duration(rand.Intn(120)+80) * time.Millisecond)
	}

	// Brief pause after scrolling
	time.Sleep(time.Duration(rand.Intn(400)+300) * time.Millisecond)
}

// RandomPageScroll simulates occasional page exploration
func RandomPageScroll(page *rod.Page) {
	if rand.Float64() < 0.4 { // 40% chance
		scrolls := rand.Intn(2) + 1
		for i := 0; i < scrolls; i++ {
			HumanScroll(page)
			time.Sleep(time.Duration(rand.Intn(800)+600) * time.Millisecond)
		}
	}
}
