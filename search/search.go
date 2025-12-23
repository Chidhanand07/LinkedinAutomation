package search

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"linkedin-automation/config"
	"linkedin-automation/stealth"

	"github.com/go-rod/rod"
)

// FindProfiles searches for profiles using config-based filters
func FindProfiles(page *rod.Page, cfg *config.SearchConfig) ([]string, error) {
	log.Println("🔍 Starting profile search with configured filters...")

	var allProfiles []string
	seen := make(map[string]bool)

	// Use custom search URLs if provided, otherwise build from config
	searchURLs := cfg.CustomSearchURLs
	if len(searchURLs) == 0 {
		searchURLs = buildSearchURLs(cfg)
	}

	for _, searchURL := range searchURLs {
		log.Printf("Searching: %s", searchURL)

		profiles, err := searchWithURL(page, searchURL, cfg, seen)
		if err != nil {
			log.Printf("⚠️  Error searching URL %s: %v", searchURL, err)
			continue
		}

		allProfiles = append(allProfiles, profiles...)

		// Check if we've reached max profiles
		if len(allProfiles) >= cfg.MaxProfiles {
			log.Printf("✓ Reached maximum profile limit: %d", cfg.MaxProfiles)
			break
		}

		// Delay between search URLs
		time.Sleep(stealth.RandomizedTimingDelay("navigate"))
	}

	// Trim to max profiles
	if len(allProfiles) > cfg.MaxProfiles {
		allProfiles = allProfiles[:cfg.MaxProfiles]
	}

	log.Printf("✓ Found %d total profiles", len(allProfiles))
	return allProfiles, nil
}

// buildSearchURLs creates search URLs from config criteria
func buildSearchURLs(cfg *config.SearchConfig) []string {
	var urls []string
	baseURL := "https://www.linkedin.com/search/results/people/?"

	// Build URLs for each keyword
	for _, keyword := range cfg.Keywords {
		params := url.Values{}
		params.Add("keywords", keyword)

		// Add connection level filter
		if cfg.ConnectionLevel != "" {
			params.Add("network", fmt.Sprintf("[\"%s\"]", cfg.ConnectionLevel))
		}

		searchURL := baseURL + params.Encode()
		urls = append(urls, searchURL)
	}

	// If no keywords, create a basic people search
	if len(urls) == 0 {
		urls = append(urls, "https://www.linkedin.com/search/results/people/")
	}

	return urls
}

// searchWithURL performs search for a specific URL with pagination
func searchWithURL(page *rod.Page, searchURL string, cfg *config.SearchConfig, seen map[string]bool) ([]string, error) {
	var profiles []string

	for pageNum := 1; pageNum <= cfg.MaxPages; pageNum++ {
		// Navigate to search page
		paginatedURL := searchURL
		if pageNum > 1 {
			separator := "&"
			if !strings.Contains(searchURL, "?") {
				separator = "?"
			}
			paginatedURL = fmt.Sprintf("%s%spage=%d", searchURL, separator, pageNum)
		}

		log.Printf("  Page %d: %s", pageNum, paginatedURL)
		page.MustNavigate(paginatedURL)
		time.Sleep(stealth.RandomizedTimingDelay("navigate"))

		// Random scroll to load dynamic content
		stealth.RandomPageScroll(page)
		time.Sleep(time.Duration(rand.Intn(2000)+1000) * time.Millisecond)

		// Extract profile links
		links, err := page.Elements("a.app-aware-link")
		if err != nil {
			return profiles, fmt.Errorf("failed to find profile links: %w", err)
		}

		if len(links) == 0 {
			log.Printf("  No more results on page %d", pageNum)
			break
		}

		foundOnPage := 0
		for _, link := range links {
			href, err := link.Attribute("href")
			if err != nil || href == nil {
				continue
			}

			profileURL := *href
			if !isValidProfileURL(profileURL) {
				continue
			}

			// Normalize URL (remove query params)
			profileURL = normalizeProfileURL(profileURL)

			// Skip duplicates
			if cfg.SkipDuplicates && seen[profileURL] {
				continue
			}

			// Apply filters
			if !matchesFilters(page, link, cfg) {
				continue
			}

			seen[profileURL] = true
			profiles = append(profiles, profileURL)
			foundOnPage++

			// Check if reached max profiles
			if len(profiles) >= cfg.MaxProfiles {
				return profiles, nil
			}
		}

		log.Printf("  ✓ Found %d profiles on page %d", foundOnPage, pageNum)

		// If no profiles found on this page, stop pagination
		if foundOnPage == 0 {
			break
		}

		// Delay before next page
		if pageNum < cfg.MaxPages {
			delay := time.Duration(rand.Intn(3000)+2000) * time.Millisecond
			time.Sleep(delay)
		}
	}

	return profiles, nil
}

// isValidProfileURL checks if URL is a valid LinkedIn profile
func isValidProfileURL(urlStr string) bool {
	return strings.Contains(urlStr, "/in/") &&
		!strings.Contains(urlStr, "overlay") &&
		!strings.Contains(urlStr, "miniProfile") &&
		!strings.Contains(urlStr, "/company/")
}

// normalizeProfileURL removes query parameters from profile URL
func normalizeProfileURL(urlStr string) string {
	if idx := strings.Index(urlStr, "?"); idx != -1 {
		urlStr = urlStr[:idx]
	}
	return strings.TrimRight(urlStr, "/")
}

// matchesFilters applies configured filters to profile
func matchesFilters(page *rod.Page, link *rod.Element, cfg *config.SearchConfig) bool {
	// Get parent element text for filtering
	parent, err := link.Parent()
	if err != nil || parent == nil {
		return true
	}

	text := ""
	if t, err := parent.Text(); err == nil {
		text = strings.ToLower(t)
	}

	// Filter by job titles
	if len(cfg.JobTitles) > 0 {
		matched := false
		for _, title := range cfg.JobTitles {
			if strings.Contains(text, strings.ToLower(title)) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// Filter by companies
	if len(cfg.Companies) > 0 {
		matched := false
		for _, company := range cfg.Companies {
			if strings.Contains(text, strings.ToLower(company)) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// Filter by locations
	if len(cfg.Locations) > 0 {
		matched := false
		for _, location := range cfg.Locations {
			if strings.Contains(text, strings.ToLower(location)) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	return true
}
