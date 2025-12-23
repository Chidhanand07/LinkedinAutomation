package accept

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"linkedin-automation/config"
	"linkedin-automation/stealth"

	"github.com/go-rod/rod"
)

// AcceptPendingConnections navigates to My Network and accepts all pending connection requests
func AcceptPendingConnections(page *rod.Page, cfg *config.AutoAcceptConfig) (int, error) {
	if !cfg.Enabled {
		log.Println("⚠️  Auto-accept is disabled in config")
		return 0, nil
	}

	log.Println("🤝 Starting auto-accept pending connections...")

	// Navigate to My Network page
	page.MustNavigate("https://www.linkedin.com/mynetwork/invitation-manager/")
	time.Sleep(stealth.RandomizedTimingDelay("navigate"))

	// Scroll to load all pending requests
	for i := 0; i < 3; i++ {
		stealth.HumanScroll(page)
		time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)
	}

	accepted := 0

	// Find all "Accept" buttons
	acceptButtons, err := page.Elements("button[aria-label*='Accept']")
	if err != nil {
		return 0, fmt.Errorf("failed to find accept buttons: %w", err)
	}

	if len(acceptButtons) == 0 {
		log.Println("✓ No pending connection requests to accept")
		return 0, nil
	}

	log.Printf("Found %d pending connection requests", len(acceptButtons))

	// Accept connections up to the configured limit
	maxToAccept := cfg.MaxPerSession
	if len(acceptButtons) < maxToAccept {
		maxToAccept = len(acceptButtons)
	}

	for i := 0; i < maxToAccept; i++ {
		// Re-query buttons as DOM changes after each accept
		acceptButtons, err = page.Elements("button[aria-label*='Accept']")
		if err != nil || len(acceptButtons) == 0 {
			break
		}

		button := acceptButtons[0]

		// Get the profile name if possible
		profileName := getProfileNameFromButton(button)
		if profileName != "" {
			log.Printf("  Accepting connection from: %s", profileName)
		} else {
			log.Printf("  Accepting connection %d/%d", i+1, maxToAccept)
		}

		// Move mouse to button and click
		box := button.MustShape().Box()
		targetX := box.X + box.Width/2 + (rand.Float64()-0.5)*10
		targetY := box.Y + box.Height/2 + (rand.Float64()-0.5)*10
		stealth.HumanMove(page, targetX, targetY)

		// Brief hover
		stealth.HumanHover(page, targetX, targetY, time.Duration(rand.Intn(200)+100)*time.Millisecond)

		// Click accept
		button.MustClick()
		accepted++

		// Wait for the action to complete
		time.Sleep(stealth.RandomizedTimingDelay("click"))

		// Send welcome message if enabled
		if cfg.SendWelcomeMessage && cfg.WelcomeMessageTemplate != "" {
			// Note: This would require navigating to the profile and sending a message
			// For now, we'll skip this as it requires additional navigation
			// You can implement this if needed
		}

		// Delay between accepts
		delay := time.Duration(rand.Intn(5000)+3000) * time.Millisecond
		log.Printf("  ⏳ Waiting %v before next accept...", delay.Round(time.Second))
		time.Sleep(delay)
	}

	log.Printf("✓ Accepted %d connection requests", accepted)
	return accepted, nil
}

// getProfileNameFromButton tries to extract the profile name from button context
func getProfileNameFromButton(button *rod.Element) string {
	// Try to find the name in nearby elements
	parent, err := button.Parent()
	if err != nil || parent == nil {
		return ""
	}

	// Look for name elements
	nameElement, err := parent.Element("span[dir='ltr'] span[aria-hidden='true']")
	if err != nil || nameElement == nil {
		return ""
	}

	name, err := nameElement.Text()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(name)
}

// GetPendingRequestCount returns the number of pending connection requests
func GetPendingRequestCount(page *rod.Page) (int, error) {
	// Navigate to My Network
	page.MustNavigate("https://www.linkedin.com/mynetwork/invitation-manager/")
	time.Sleep(3 * time.Second)

	// Try to find the count badge
	countElement, err := page.Element(".mn-invitation-manager__count")
	if err != nil {
		// If no count element, check for accept buttons
		acceptButtons, err := page.Elements("button[aria-label*='Accept']")
		if err != nil {
			return 0, nil
		}
		return len(acceptButtons), nil
	}

	countText, err := countElement.Text()
	if err != nil {
		return 0, nil
	}

	// Parse count
	var count int
	fmt.Sscanf(countText, "%d", &count)
	return count, nil
}
