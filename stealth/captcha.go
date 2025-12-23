package stealth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// SessionManager handles browser session persistence to avoid repeated logins
type SessionManager struct {
	sessionFile string
}

// NewSessionManager creates a new session manager
func NewSessionManager(file string) *SessionManager {
	return &SessionManager{sessionFile: file}
}

// SaveSession saves current browser cookies to file
func (sm *SessionManager) SaveSession(page *rod.Page) error {
	cookies := page.MustCookies()
	data, err := json.MarshalIndent(cookies, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cookies: %w", err)
	}

	if err := os.WriteFile(sm.sessionFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	log.Printf("✓ Session saved to %s", sm.sessionFile)
	return nil
}

// LoadSession restores browser cookies from file
func (sm *SessionManager) LoadSession(page *rod.Page) error {
	data, err := os.ReadFile(sm.sessionFile)
	if err != nil {
		return fmt.Errorf("session file not found: %w", err)
	}

	var cookies []*proto.NetworkCookie
	if err := json.Unmarshal(data, &cookies); err != nil {
		return fmt.Errorf("failed to parse session file: %w", err)
	}

	if len(cookies) == 0 {
		return errors.New("no cookies found in session file")
	}

	// Convert NetworkCookie to NetworkCookieParam
	cookieParams := make([]*proto.NetworkCookieParam, len(cookies))
	for i, cookie := range cookies {
		cookieParams[i] = &proto.NetworkCookieParam{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Secure:   cookie.Secure,
			HTTPOnly: cookie.HTTPOnly,
			SameSite: cookie.SameSite,
			Expires:  cookie.Expires,
		}
	}

	// Set cookies
	if err := page.SetCookies(cookieParams); err != nil {
		return fmt.Errorf("failed to set cookies: %w", err)
	}

	log.Printf("✓ Loaded %d cookies from session", len(cookies))
	return nil
}

// SessionExists checks if a session file exists
func (sm *SessionManager) SessionExists() bool {
	_, err := os.Stat(sm.sessionFile)
	return err == nil
}

// DeleteSession removes the session file
func (sm *SessionManager) DeleteSession() error {
	if !sm.SessionExists() {
		return nil
	}
	return os.Remove(sm.sessionFile)
}

// WaitForManualCaptchaSolve pauses automation and waits for user to solve CAPTCHA manually
func WaitForManualCaptchaSolve(page *rod.Page, timeout time.Duration) error {
	log.Println("⚠️  CAPTCHA DETECTED!")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("Please solve the CAPTCHA manually in the browser window.")
	log.Printf("Waiting up to %v for manual intervention...\n", timeout)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	start := time.Now()
	checkInterval := 1 * time.Second
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			elapsed := time.Since(start)
			remaining := timeout - elapsed

			if elapsed >= timeout {
				return fmt.Errorf("CAPTCHA solve timeout after %v", timeout)
			}

			// Check if common CAPTCHA indicators are gone
			hasCaptcha := detectCaptcha(page)

			if !hasCaptcha {
				log.Println("✓ CAPTCHA appears to be solved!")
				time.Sleep(2 * time.Second) // Brief pause to ensure page is stable
				return nil
			}

			// Log progress every 10 seconds
			if int(elapsed.Seconds())%10 == 0 {
				log.Printf("⏳ Still waiting... (%v remaining)", remaining.Round(time.Second))
			}
		}
	}
}

// detectCaptcha checks for various CAPTCHA indicators on the page
func detectCaptcha(page *rod.Page) bool {
	captchaSelectors := []string{
		"iframe[src*='captcha']",
		"iframe[src*='recaptcha']",
		"iframe[title*='recaptcha']",
		"#captcha",
		".g-recaptcha",
		"[data-sitekey]",
		"#recaptcha",
		"iframe[src*='hcaptcha']",
		".h-captcha",
	}

	for _, selector := range captchaSelectors {
		if page.MustHas(selector) {
			return true
		}
	}

	// Check for common CAPTCHA text
	html := page.MustHTML()
	captchaTexts := []string{
		"I'm not a robot",
		"Please verify you are human",
		"Security check",
		"Complete the security check",
	}

	for _, text := range captchaTexts {
		if contains(html, text) {
			return true
		}
	}

	return false
}

// contains checks if a string contains a substring (case-insensitive helper)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// CheckForSecurityChallenge detects various security challenges
func CheckForSecurityChallenge(page *rod.Page) (challengeType string, detected bool) {
	// Check for PIN verification
	if page.MustHas("input[name=pin]") || page.MustHas("#input__phone_verification_pin") {
		return "PIN_VERIFICATION", true
	}

	// Check for CAPTCHA
	if detectCaptcha(page) {
		return "CAPTCHA", true
	}

	// Check for email verification
	if page.MustHas("input[name='verification-code']") || page.MustHas("#email-pin-challenge") {
		return "EMAIL_VERIFICATION", true
	}

	// Check for suspicious activity warning
	html := page.MustHTML()
	if contains(html, "unusual activity") || contains(html, "suspicious activity") {
		return "SUSPICIOUS_ACTIVITY", true
	}

	return "", false
}

// RateLimiter manages timing between actions to avoid rate limiting
type RateLimiter struct {
	lastAction   time.Time
	minInterval  time.Duration
	actionCounts map[string]int
	resetTime    time.Time
}

// NewRateLimiter creates a rate limiter with specified minimum interval
func NewRateLimiter(minInterval time.Duration) *RateLimiter {
	return &RateLimiter{
		lastAction:   time.Time{},
		minInterval:  minInterval,
		actionCounts: make(map[string]int),
		resetTime:    time.Now().Add(24 * time.Hour),
	}
}

// WaitIfNeeded enforces minimum interval between actions
func (r *RateLimiter) WaitIfNeeded() {
	if r.lastAction.IsZero() {
		r.lastAction = time.Now()
		return
	}

	elapsed := time.Since(r.lastAction)
	if elapsed < r.minInterval {
		waitTime := r.minInterval - elapsed
		log.Printf("⏱️  Rate limiting: waiting %v before next action", waitTime.Round(time.Second))
		time.Sleep(waitTime)
	}

	r.lastAction = time.Now()
}

// RecordAction tracks action counts for daily limits
func (r *RateLimiter) RecordAction(actionType string) {
	// Reset counters if 24 hours have passed
	if time.Now().After(r.resetTime) {
		r.actionCounts = make(map[string]int)
		r.resetTime = time.Now().Add(24 * time.Hour)
		log.Println("📊 Daily action counters reset")
	}

	r.actionCounts[actionType]++
	log.Printf("📈 Action recorded: %s (total today: %d)", actionType, r.actionCounts[actionType])
}

// GetActionCount returns the count for a specific action type
func (r *RateLimiter) GetActionCount(actionType string) int {
	if time.Now().After(r.resetTime) {
		return 0
	}
	return r.actionCounts[actionType]
}

// CheckLimit checks if action limit has been reached
func (r *RateLimiter) CheckLimit(actionType string, maxPerDay int) bool {
	count := r.GetActionCount(actionType)
	if count >= maxPerDay {
		log.Printf("⚠️  Daily limit reached for %s (%d/%d)", actionType, count, maxPerDay)
		return false
	}
	return true
}
