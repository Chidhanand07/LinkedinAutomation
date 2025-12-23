package browser

import (
	"log"
	"math/rand"
	"os"
	"time"

	"linkedin-automation/stealth"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func NewBrowser() (*rod.Browser, *rod.Page) {
	rand.Seed(time.Now().UnixNano())

	// Prefer an existing Chrome/Chromium binary to avoid downloading
	// leakless.exe which may be flagged by Windows Defender.
	findChrome := func() string {
		// common install locations
		paths := []string{
			"C:/Program Files/Google/Chrome/Application/chrome.exe",
			"C:/Program Files (x86)/Google/Chrome/Application/chrome.exe",
			"C:/Program Files/Microsoft/Edge/Application/msedge.exe",
			"/usr/bin/google-chrome",
			"/usr/bin/chromium-browser",
			"/usr/bin/chromium",
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		}
		for _, p := range paths {
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
		return ""
	}

	// Create persistent profile directory to maintain trust signals
	profileDir := "./chrome-profile"
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		log.Printf("Warning: Could not create profile directory: %v", err)
		profileDir = "" // Fall back to temporary profile
	}

	l := launcher.New().
		Headless(false).
		// Core anti-detection flags
		Set("disable-blink-features", "AutomationControlled").
		Set("disable-infobars", "true").
		Set("exclude-switches", "enable-automation").
		// Additional stealth flags to reduce CAPTCHA triggers
		Set("disable-features", "IsolateOrigins,site-per-process").
		Set("disable-web-security", "true").
		// Persistent profile for trust signals
		UserDataDir(profileDir).
		// Window size
		Set("window-size", "1920,1080")

	// Optional proxy support via environment variable
	// Example: PROXY=http://username:password@proxy.com:8080
	if proxyURL := os.Getenv("PROXY"); proxyURL != "" {
		log.Printf("Using proxy: %s", proxyURL)
		l = l.Proxy(proxyURL)
	}

	// choose binary: env override preferred
	bin := os.Getenv("CHROME_BIN")
	if bin == "" {
		bin = findChrome()
	}

	if bin == "" {
		log.Fatal("No local Chrome/Edge binary found and CHROME_BIN is not set.\nPlease install Chrome/Edge or set the CHROME_BIN environment variable to avoid Rod downloading leakless.exe which may be blocked by antivirus. Example:\n$env:CHROME_BIN='C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe'\n")
	}

	log.Printf("Using browser binary: %s", bin)
	if profileDir != "" {
		log.Printf("Using persistent profile: %s", profileDir)
	}
	l = l.Bin(bin)

	u := l.MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()
	page := browser.MustPage()

	// Apply comprehensive fingerprint masking on browser initialization
	stealth.EnhancedFingerprintMasking(page)

	log.Println("Browser initialized with anti-detection stealth techniques")
	return browser, page
}
