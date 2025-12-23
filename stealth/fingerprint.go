package stealth

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// EnhancedFingerprintMasking applies comprehensive browser fingerprint modifications
// to avoid detection by anti-bot systems
func EnhancedFingerprintMasking(page *rod.Page) {
	// Randomize user agent with realistic combinations
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
	selectedUA := userAgents[rand.Intn(len(userAgents))]

	page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent: selectedUA,
	})

	// Determine platform based on host OS or randomize
	platform := "Win32"
	platformVersion := "10.0.0"
	if runtime.GOOS == "darwin" {
		platform = "MacIntel"
		platformVersion = "10.15.7"
	} else if runtime.GOOS == "linux" {
		platform = "Linux x86_64"
		platformVersion = ""
	}

	// Comprehensive JavaScript injection to mask automation signatures
	page.MustEval(fmt.Sprintf(`(platform, platformVersion, userAgent) => {
		// Remove webdriver flag
		Object.defineProperty(navigator, 'webdriver', { get: () => undefined });

		// Set realistic language preferences
		Object.defineProperty(navigator, 'languages', { 
			get: () => ['en-US', 'en'] 
		});
		Object.defineProperty(navigator, 'language', { 
			get: () => 'en-US' 
		});

		// Set platform
		Object.defineProperty(navigator, 'platform', { 
			get: () => platform 
		});

		// Set realistic hardware concurrency
		Object.defineProperty(navigator, 'hardwareConcurrency', { 
			get: () => %d 
		});

		// Set realistic device memory (GB)
		Object.defineProperty(navigator, 'deviceMemory', { 
			get: () => %d 
		});

		// Mask automation in permissions
		const originalQuery = window.navigator.permissions.query;
		window.navigator.permissions.query = (parameters) => (
			parameters.name === 'notifications' ?
				Promise.resolve({ state: Notification.permission }) :
				originalQuery(parameters)
		);

		// Add realistic plugin array
		Object.defineProperty(navigator, 'plugins', {
			get: () => [
				{
					name: 'Chrome PDF Plugin',
					filename: 'internal-pdf-viewer',
					description: 'Portable Document Format'
				},
				{
					name: 'Chrome PDF Viewer',
					filename: 'mhjfbmdgcfjbbpaeojofohoefgiehjai',
					description: ''
				}
			]
		});

		// Override chrome runtime
		if (!window.chrome) {
			window.chrome = {};
		}
		window.chrome.runtime = {
			OnInstalledReason: {
				CHROME_UPDATE: 'chrome_update',
				INSTALL: 'install',
				SHARED_MODULE_UPDATE: 'shared_module_update',
				UPDATE: 'update'
			},
			OnRestartRequiredReason: {
				APP_UPDATE: 'app_update',
				OS_UPDATE: 'os_update',
				PERIODIC: 'periodic'
			},
			PlatformArch: {
				ARM: 'arm',
				ARM64: 'arm64',
				MIPS: 'mips',
				MIPS64: 'mips64',
				X86_32: 'x86-32',
				X86_64: 'x86-64'
			},
			PlatformOs: {
				ANDROID: 'android',
				CROS: 'cros',
				LINUX: 'linux',
				MAC: 'mac',
				OPENBSD: 'openbsd',
				WIN: 'win'
			}
		};

		// Randomize screen properties slightly
		const screenWidth = %d;
		const screenHeight = %d;
		Object.defineProperty(screen, 'width', { get: () => screenWidth });
		Object.defineProperty(screen, 'height', { get: () => screenHeight });
		Object.defineProperty(screen, 'availWidth', { get: () => screenWidth });
		Object.defineProperty(screen, 'availHeight', { get: () => screenHeight - 40 });

		// Mask timezone
		Object.defineProperty(Intl.DateTimeFormat.prototype, 'resolvedOptions', {
			value: function() {
				return {
					locale: 'en-US',
					calendar: 'gregory',
					numberingSystem: 'latn',
					timeZone: 'America/New_York'
				};
			}
		});

		// Add realistic battery API (if applicable)
		if (navigator.getBattery) {
			navigator.getBattery = () => Promise.resolve({
				charging: true,
				chargingTime: 0,
				dischargingTime: Infinity,
				level: %f,
				addEventListener: () => {},
				removeEventListener: () => {}
			});
		}

		// Connection API
		Object.defineProperty(navigator, 'connection', {
			get: () => ({
				effectiveType: '4g',
				rtt: 50,
				downlink: 10,
				saveData: false
			})
		});
	}`, 
		rand.Intn(8)+4,         // hardwareConcurrency: 4-12 cores
		[]int{4, 8, 16}[rand.Intn(3)], // deviceMemory: 4, 8, or 16 GB
		[]int{1920, 1366, 2560}[rand.Intn(3)], // screen width
		[]int{1080, 768, 1440}[rand.Intn(3)],  // screen height
		0.5+rand.Float64()*0.5,  // battery level: 50-100%
	), platform, platformVersion, selectedUA)
}

// RandomizedTimingDelay implements realistic, randomized delays
// that vary based on action type and simulate human cognitive processing
func RandomizedTimingDelay(actionType string) time.Duration {
	var base, variance int

	switch actionType {
	case "click":
		base = 300
		variance = 400
	case "type_char":
		base = 80
		variance = 120
	case "think":
		base = 800
		variance = 1200
	case "read":
		base = 1500
		variance = 2000
	case "navigate":
		base = 2000
		variance = 3000
	case "scroll":
		base = 500
		variance = 800
	default:
		base = 400
		variance = 600
	}

	delay := base + rand.Intn(variance)
	return time.Duration(delay) * time.Millisecond
}

// IsBusinessHours checks if current time is within typical business hours
// for activity scheduling (9 AM - 6 PM weekdays, local time)
func IsBusinessHours() bool {
	now := time.Now()
	hour := now.Hour()
	weekday := now.Weekday()

	// Weekend check
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}

	// Business hours: 9 AM - 6 PM
	return hour >= 9 && hour < 18
}

// SmartDelay implements intelligent delays that consider time of day
func SmartDelay(baseDelay time.Duration) time.Duration {
	// Add extra delay during non-business hours to seem more cautious
	if !IsBusinessHours() {
		return baseDelay + time.Duration(rand.Intn(2000)+1000)*time.Millisecond
	}
	return baseDelay
}
