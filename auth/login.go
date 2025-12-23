package auth

import (
	"errors"
	"log"
	"math/rand"
	"os"
	"time"

	"linkedin-automation/stealth"

	"github.com/go-rod/rod"
)

func Login(page *rod.Page) error {
	log.Println("Initiating login with anti-bot stealth techniques...")

	sessionMgr := stealth.NewSessionManager("linkedin_session.json")

	// [STRATEGY 1] Try to restore previous session (avoid login completely)
	if sessionMgr.SessionExists() {
		log.Println("🔄 Found existing session - attempting to restore...")

		if err := sessionMgr.LoadSession(page); err == nil {
			// Navigate to feed to check if session is valid
			page.MustNavigate("https://www.linkedin.com/feed")
			time.Sleep(stealth.RandomizedTimingDelay("navigate"))

			// Check if we're actually logged in (not redirected to login page)
			if !page.MustHas("input[type=password]") && !page.MustHas("#username") {
				log.Println("✓ Session restored successfully - skipping login!")
				return nil
			}

			log.Println("⚠️  Session expired - proceeding with fresh login")
			sessionMgr.DeleteSession() // Clean up invalid session
		}
	}

	// [TECHNIQUE 1] Activity Scheduling - Check business hours
	if !stealth.IsBusinessHours() {
		log.Println("Warning: Operating outside business hours - adding extra delays")
	}

	// Navigate to LinkedIn login
	page.MustNavigate("https://www.linkedin.com/login")

	// [TECHNIQUE 2] Randomized Timing - Initial page load delay
	time.Sleep(stealth.RandomizedTimingDelay("navigate"))

	// [TECHNIQUE 3] Random Scrolling Behavior - Explore page briefly
	stealth.RandomPageScroll(page)

	email := os.Getenv("LINKEDIN_EMAIL")
	password := os.Getenv("LINKEDIN_PASSWORD")

	if email == "" || password == "" {
		return errors.New("credentials not found in environment variables")
	}

	// Locate email input
	emailInput, err := page.Element("#username")
	if err != nil {
		return errors.New("email input not found")
	}

	passwordInput, err := page.Element("#password")
	if err != nil {
		return errors.New("password input not found")
	}

	// [TECHNIQUE 4] Random Mouse Wandering - Simulate reading the page
	stealth.RandomMouseWander(page)
	time.Sleep(stealth.RandomizedTimingDelay("read"))

	// [TECHNIQUE 5] Human-like Mouse Movement with Bézier curves
	log.Println("Moving to email field...")
	emailBox := emailInput.MustShape()
	targetX := emailBox.Box().X + emailBox.Box().Width/2 + (rand.Float64()-0.5)*20
	targetY := emailBox.Box().Y + emailBox.Box().Height/2 + (rand.Float64()-0.5)*10
	stealth.HumanMove(page, targetX, targetY)

	// [TECHNIQUE 6] Mouse Hovering - Brief hover before clicking
	stealth.HumanHover(page, targetX, targetY, time.Duration(rand.Intn(300)+200)*time.Millisecond)
	time.Sleep(stealth.RandomizedTimingDelay("think"))

	// Click email field
	emailInput.MustClick()
	time.Sleep(stealth.RandomizedTimingDelay("click"))

	// [TECHNIQUE 7] Realistic Typing Simulation with typos and corrections
	log.Println("Typing email...")
	stealth.HumanType(emailInput, email)

	// Pause after typing (reading what was typed)
	time.Sleep(stealth.SmartDelay(stealth.RandomizedTimingDelay("think")))

	// Random mouse wander while thinking
	stealth.RandomMouseWander(page)

	// Move to password field with Bézier curve
	log.Println("Moving to password field...")
	passBox := passwordInput.MustShape()
	targetX = passBox.Box().X + passBox.Box().Width/2 + (rand.Float64()-0.5)*20
	targetY = passBox.Box().Y + passBox.Box().Height/2 + (rand.Float64()-0.5)*10
	stealth.HumanMove(page, targetX, targetY)

	// Brief hover before clicking
	stealth.HumanHover(page, targetX, targetY, time.Duration(rand.Intn(200)+150)*time.Millisecond)
	time.Sleep(stealth.RandomizedTimingDelay("think"))

	passwordInput.MustClick()
	time.Sleep(stealth.RandomizedTimingDelay("click"))

	// Type password with realistic patterns
	log.Println("Typing password...")
	stealth.HumanType(passwordInput, password)

	// Longer pause after password (thinking/checking)
	time.Sleep(stealth.SmartDelay(stealth.RandomizedTimingDelay("think")))

	// Occasional mouse movement before clicking submit
	if rand.Float64() > 0.5 {
		stealth.RandomMouseWander(page)
		time.Sleep(time.Duration(rand.Intn(500)+300) * time.Millisecond)
	}

	// Locate and move to submit button
	submitBtn, err := page.Element("button[type=submit]")
	if err != nil {
		return errors.New("login button not found")
	}

	log.Println("Moving to submit button...")
	btnBox := submitBtn.MustShape()
	targetX = btnBox.Box().X + btnBox.Box().Width/2 + (rand.Float64()-0.5)*30
	targetY = btnBox.Box().Y + btnBox.Box().Height/2 + (rand.Float64()-0.5)*15
	stealth.HumanMove(page, targetX, targetY)

	// Hover over button briefly (realistic hesitation)
	stealth.HumanHover(page, targetX, targetY, time.Duration(rand.Intn(400)+300)*time.Millisecond)
	time.Sleep(stealth.RandomizedTimingDelay("think"))

	// Click submit
	log.Println("Clicking submit...")
	submitBtn.MustClick()

	// Wait for initial response
	time.Sleep(3 * time.Second)

	// [STRATEGY 2] Check for security challenges (CAPTCHA, PIN, etc.)
	if challengeType, detected := stealth.CheckForSecurityChallenge(page); detected {
		log.Printf("🔒 Security challenge detected: %s", challengeType)

		switch challengeType {
		case "CAPTCHA":
			// Manual CAPTCHA solving (free method)
			if err := stealth.WaitForManualCaptchaSolve(page, 180*time.Second); err != nil {
				return err
			}
			// Wait a bit more after CAPTCHA is solved
			time.Sleep(stealth.SmartDelay(5 * time.Second))

		case "PIN_VERIFICATION", "EMAIL_VERIFICATION":
			log.Println("⚠️  Verification required - please complete it manually")
			log.Println("Waiting 120 seconds for manual verification...")
			time.Sleep(120 * time.Second)

		case "SUSPICIOUS_ACTIVITY":
			return errors.New("account flagged for suspicious activity - manual review required")
		}
	}

	// Wait for login to complete
	time.Sleep(stealth.SmartDelay(5 * time.Second))

	// Verify successful login
	if page.MustHas("input[type=password]") || page.MustHas("#username") {
		return errors.New("login failed - still on login page")
	}

	// [STRATEGY 1] Save session for future use
	if err := sessionMgr.SaveSession(page); err != nil {
		log.Printf("⚠️  Could not save session: %v", err)
	} else {
		log.Println("✓ Session saved - next run will skip login!")
	}

	log.Println("✓ Login completed successfully with full stealth techniques applied")
	return nil
}
