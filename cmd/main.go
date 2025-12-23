package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"linkedin-automation/accept"
	"linkedin-automation/auth"
	"linkedin-automation/browser"
	"linkedin-automation/config"
	"linkedin-automation/connect"
	"linkedin-automation/messaging"
	"linkedin-automation/search"
	"linkedin-automation/stealth"
	"linkedin-automation/storage"
)

func main() {
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("🚀 LinkedIn Automation POC - Config-Driven")
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Load configuration
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("✓ Configuration loaded successfully")

	// Override credentials from environment variables if provided
	if email := os.Getenv("LINKEDIN_EMAIL"); email != "" {
		cfg.Credentials.Email = email
	}
	if password := os.Getenv("LINKEDIN_PASSWORD"); password != "" {
		cfg.Credentials.Password = password
	}

	if cfg.Credentials.Email == "" || cfg.Credentials.Password == "" {
		log.Fatal("❌ LinkedIn credentials not found. Set LINKEDIN_EMAIL and LINKEDIN_PASSWORD environment variables.")
	}

	// Set credentials as environment variables for auth module
	os.Setenv("LINKEDIN_EMAIL", cfg.Credentials.Email)
	os.Setenv("LINKEDIN_PASSWORD", cfg.Credentials.Password)

	// Initialize rate limiter
	rateLimiter := stealth.NewRateLimiter(time.Duration(cfg.RateLimiting.MinIntervalBetweenRunsMinutes) * time.Minute)
	rateLimiter.WaitIfNeeded()

	// Initialize database
	db, err := storage.InitDB(cfg.Storage.DatabaseFile)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer db.Close()
	log.Println("✓ Database initialized")

	// Initialize browser
	browserInstance, page := browser.NewBrowser()
	defer browserInstance.Close()

	// Login
	log.Println("\n📍 Phase 1: Authentication")
	err = auth.Login(page)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	time.Sleep(3 * time.Second)

	// Phase 2: Auto-accept pending connections
	if cfg.AutoAccept.Enabled {
		log.Println("\n📍 Phase 2: Auto-Accept Pending Connections")
		acceptedCount, err := accept.AcceptPendingConnections(page, &cfg.AutoAccept)
		if err != nil {
			log.Printf("⚠️  Auto-accept error: %v", err)
		} else {
			rateLimiter.RecordAction("connection_accepted")
			log.Printf("✓ Accepted %d connections", acceptedCount)
		}
		time.Sleep(3 * time.Second)
	}

	// Phase 3: Auto-reply to messages
	if cfg.Messaging.AutoReplyEnabled {
		log.Println("\n📍 Phase 3: Auto-Reply to Messages")
		repliesCount, err := messaging.AutoReplyToMessages(page, &cfg.Messaging)
		if err != nil {
			log.Printf("⚠️  Auto-reply error: %v", err)
		} else {
			rateLimiter.RecordAction("message_sent")
			log.Printf("✓ Sent %d auto-replies", repliesCount)
		}
		time.Sleep(3 * time.Second)
	}

	// Phase 4: Search for profiles
	log.Println("\n📍 Phase 4: Search & Profile Discovery")
	log.Printf("Search criteria: Keywords=%v, JobTitles=%v, Companies=%v, Locations=%v",
		cfg.Search.Keywords, cfg.Search.JobTitles, cfg.Search.Companies, cfg.Search.Locations)

	profiles, err := search.FindProfiles(page, &cfg.Search)
	if err != nil {
		log.Fatalf("Profile search failed: %v", err)
	}

	if len(profiles) == 0 {
		log.Println("⚠️  No profiles found matching search criteria")
		return
	}

	// Phase 5: Send connection requests
	log.Println("\n📍 Phase 5: Sending Connection Requests")
	connectionsSent := 0
	maxConnections := cfg.Connection.MaxRequestsPerSession

	if !rateLimiter.CheckLimit("connection_request", cfg.Connection.MaxRequestsPerDay) {
		log.Println("⚠️  Daily connection limit reached - skipping connection requests")
	} else {
		for _, profileURL := range profiles {
			if connectionsSent >= maxConnections {
				log.Printf("✓ Reached session limit: %d connections", maxConnections)
				break
			}

			// Check if already processed
			if storage.IsProfileProcessed(db, profileURL) {
				log.Printf("  ⏭️  Skipping already processed: %s", profileURL)
				continue
			}

			// Check daily limit
			if !rateLimiter.CheckLimit("connection_request", cfg.Connection.MaxRequestsPerDay) {
				log.Println("⚠️  Daily limit reached - stopping")
				break
			}

			log.Printf("  [%d/%d] Sending request to: %s", connectionsSent+1, maxConnections, profileURL)

			err := connect.SendConnectionRequest(page, profileURL)
			if err != nil {
				log.Printf("  ⚠️  Failed: %v", err)
				continue
			}

			storage.MarkProfileProcessed(db, profileURL)
			rateLimiter.RecordAction("connection_request")
			connectionsSent++

			// Variable delay between requests
			baseDelay := cfg.Connection.DelayBetweenRequestsMin
			variance := cfg.Connection.DelayBetweenRequestsMax - cfg.Connection.DelayBetweenRequestsMin
			delay := time.Duration(baseDelay+connectionsSent*2+rand.Intn(variance)) * time.Second

			log.Printf("  ⏳ Waiting %v before next request...", delay.Round(time.Second))
			time.Sleep(delay)
		}
	}

	// Phase 6: Send follow-up messages
	if cfg.Messaging.SendFollowupToNewConnections {
		log.Println("\n📍 Phase 6: Sending Follow-up Messages")
		err = messaging.SendFollowUps(page, db, &cfg.Messaging)
		if err != nil {
			log.Printf("⚠️  Messaging error: %v", err)
		}
	}

	// Summary
	log.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("📊 Session Summary:")
	log.Printf("   • Profiles found: %d", len(profiles))
	log.Printf("   • Connection requests sent: %d", connectionsSent)
	log.Printf("   • Daily connection requests: %d/%d",
		rateLimiter.GetActionCount("connection_request"),
		cfg.Connection.MaxRequestsPerDay)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Println("✅ Automation run completed successfully")
}
