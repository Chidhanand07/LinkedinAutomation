# LinkedIn Automation System - Complete Guide

## 🎯 System Overview

A comprehensive, configuration-driven LinkedIn automation system with advanced features:
- **Search & Targeting**: Multi-criteria search with pagination
- **Connection Requests**: Personalized notes with template system
- **Auto Messaging**: Follow-up to new connections
- **CAPTCHA Handling**: Session persistence + manual solving
- **Anti-Bot Detection**: 8+ stealth techniques
- **Rate Limiting**: Safety controls to avoid flagging

## 📋 Features

### Search & Targeting
✅ Search by job title, company, location, keywords  
✅ Custom search URLs  
✅ Pagination (process multiple pages)  
✅ Duplicate detection  
✅ Skip already connected profiles  
✅ Profile data extraction (name, headline, company)  
✅ Advanced filtering (connections, profile picture, etc.)  

### Connection Requests
✅ Navigate to profiles programmatically  
✅ Click Connect with precision  
✅ Personalized notes with templates  
✅ Character limit enforcement  
✅ Daily/hourly/session limits  
✅ Profile interaction (view, scroll, hover)  
✅ Request tracking in database  

### Messaging System
✅ Detect new connections  
✅ Auto follow-up messages  
✅ Template system with variables  
✅ Message tracking  
✅ Configurable wait period  
✅ Daily/hourly limits  
✅ Skip already messaged  

### CAPTCHA & Session Management
✅ Session persistence (skip login)  
✅ Manual CAPTCHA solving  
✅ Automatic challenge detection  
✅ Persistent browser profile  
✅ Proxy support  

### Anti-Detection
✅ Bézier mouse movement  
✅ Realistic typing with typos  
✅ Natural scrolling  
✅ Mouse hovering  
✅ Fingerprint masking  
✅ Business hours awareness  
✅ Variable timing  
✅ Rate limiting  

## 🚀 Quick Start

### 1. Setup

```powershell
# Clone or download the project
cd "D:\AI Development\GO Automation"

# Set environment variables
$env:LINKEDIN_EMAIL = "your-email@example.com"
$env:LINKEDIN_PASSWORD = "your-password"
$env:CHROME_BIN = "C:\Program Files\Google\Chrome\Application\chrome.exe"

# Optional: Set proxy
$env:PROXY = "http://proxy-server.com:8080"
```

### 2. Configure

Edit `config/config.yaml` to customize:

```yaml
search:
  keywords: ["Software Engineer", "Developer"]
  max_profiles: 30
  max_pages: 3

connection:
  max_requests_per_day: 25
  send_note: true
  note_templates:
    - "Hi {{FirstName}}, let's connect!"

messaging:
  enabled: true
  max_messages_per_day: 15
```

### 3. Run

```powershell
# First run (may prompt for CAPTCHA)
go run .\cmd\main.go

# Subsequent runs (session restored automatically)
go run .\cmd\main.go
```

## 📖 Configuration Guide

### Template Variables

Use these in `note_templates` and `followup_templates`:

- `{{FirstName}}` - Person's first name
- `{{LastName}}` - Person's last name  
- `{{Company}}` - Current company
- `{{JobTitle}}` - Current job title
- `{{Location}}` - Profile location
- `{{Keywords}}` - Search keywords

**Example**:
```yaml
note_templates:
  - "Hi {{FirstName}}, impressed by your work at {{Company}}!"
  - "Hello {{FirstName}}, I'd love to connect about {{JobTitle}} topics."
```

### Search Configuration

```yaml
search:
  # What to search for
  keywords: ["Software Engineer", "Backend Developer"]
  job_titles: ["Senior Engineer", "Lead Developer"]
  companies: ["Google", "Microsoft", "Amazon"]
  locations: ["San Francisco", "New York"]
  
  # Pagination
  max_pages: 5  # Pages to process per search
  max_profiles: 50  # Total profiles to collect
  
  # Filters
  skip_duplicates: true
  skip_already_connected: true
  
  # Custom URLs (optional)
  custom_search_urls:
    - "https://www.linkedin.com/search/results/people/?keywords=engineer"
```

### Connection Configuration

```yaml
connection:
  # Limits
  max_requests_per_day: 30
  max_requests_per_hour: 10
  max_requests_per_session: 20
  
  # Notes
  send_note: true
  note_templates:
    - "Hi {{FirstName}}, let's connect!"
    - "Hello {{FirstName}}, impressed by {{Company}}!"
  max_note_length: 300
  
  # Timing
  delay_between_requests_min: 20  # seconds
  delay_between_requests_max: 60
  
  # Behavior
  view_profile_before_connect: true
  scroll_profile_before_connect: true
```

### Messaging Configuration

```yaml
messaging:
  enabled: true
  
  # Auto follow-up
  send_followup_to_new_connections: true
  wait_hours_after_acceptance: 24
  
  # Limits
  max_messages_per_day: 20
  max_messages_per_hour: 5
  
  # Templates
  followup_templates:
    - "Thanks for connecting, {{FirstName}}!"
    - "Hi {{FirstName}}, great to be connected!"
  
  # Tracking
  skip_already_messaged: true
```

### Rate Limiting

```yaml
rate_limiting:
  min_interval_between_runs_minutes: 30
  max_actions_per_day: 100
  
  # Business hours (optional - reduces CAPTCHA risk)
  operate_business_hours_only: false
  business_hours_start: 9   # 9 AM
  business_hours_end: 18    # 6 PM
```

## 📊 Database Tracking

All activity is tracked in `linkedin_automation.db`:

- **Profiles**: URL, name, status, timestamps
- **Connection Requests**: Sent requests with dates
- **Messages**: Sent messages with content
- **Action Counters**: Daily limits enforcement

**Query examples**:
```powershell
# View sent requests
sqlite3 linkedin_automation.db "SELECT * FROM profiles WHERE status='connected'"

# View messages
sqlite3 linkedin_automation.db "SELECT * FROM messages"
```

## 🔒 CAPTCHA Handling

### First Time (May See CAPTCHA)

1. Automation starts login
2. If CAPTCHA appears, you'll see:
   ```
   ⚠️  CAPTCHA DETECTED!
   Please solve manually in the browser window.
   Waiting up to 3 minutes...
   ```
3. Solve CAPTCHA manually
4. Automation continues automatically
5. Session saved for future

### Subsequent Runs (No CAPTCHA)

```
🔄 Found existing session - attempting to restore...
✓ Session restored successfully - skipping login!
```

Session file: `linkedin_session.json`

### Force Fresh Login

```powershell
Remove-Item linkedin_session.json
go run .\cmd\main.go
```

## 🛡️ Best Practices

### Avoid Getting Flagged

1. **Start Small**: Begin with 10-15 actions/day, increase gradually
2. **Use Realistic Limits**: Don't exceed 30-50 connections/day
3. **Business Hours**: Run during 9 AM - 6 PM M-F when possible
4. **Session Persistence**: Let it restore sessions (reduces suspicion)
5. **Vary Timing**: Use provided randomization
6. **Profile Warm-up**: Use account manually for a week before automating

### Optimal Settings

```yaml
connection:
  max_requests_per_day: 25  # Conservative limit
  delay_between_requests_min: 25
  delay_between_requests_max: 50

rate_limiting:
  min_interval_between_runs_minutes: 60  # Run hourly max
  operate_business_hours_only: true
```

## 🐛 Troubleshooting

### "Configuration not loaded" Error

```powershell
# Ensure config.yaml exists
Get-Content config\config.yaml
```

### "CAPTCHA timeout" Error

- Solve CAPTCHA faster (3 minute limit)
- Use persistent profile (`./chrome-profile/`)
- Consider residential proxy

### "Daily limit reached" Error

- Expected behavior - protects your account
- Wait 24 hours or increase limits in config
- Check `linkedin_automation.db` for current counts

### Session Won't Restore

```powershell
# Clear profile and session
Remove-Item -Recurse -Force .\chrome-profile
Remove-Item linkedin_session.json
```

### Build Errors

```powershell
# Update dependencies
go get -u ./...
go mod tidy
go build -o linkedin-automation.exe .\cmd
```

## 📁 Project Structure

```
├── auth/           # Login with CAPTCHA handling
├── browser/        # Browser setup with stealth
├── config/         # Configuration system
├── connect/        # Connection request system
├── messaging/      # Messaging system
├── search/         # Search & targeting
├── stealth/        # Anti-detection techniques
├── storage/        # Database operations
├── cmd/            # Main application
├── config.yaml     # Configuration file
└── README.md       # This file
```

## 🔧 Advanced Usage

### Using Proxy

```powershell
$env:PROXY = "http://username:password@proxy.com:8080"
go run .\cmd\main.go
```

### Custom Search URLs

```yaml
search:
  custom_search_urls:
    - "https://www.linkedin.com/search/results/people/?geoUrn=%5B103644278%5D&keywords=engineer"
    - "https://www.linkedin.com/search/results/people/?currentCompany=%5B1441%5D"
```

### Scheduled Runs (Windows Task Scheduler)

Create a PowerShell script `run-automation.ps1`:
```powershell
$env:LINKEDIN_EMAIL = "your-email@example.com"
$env:LINKEDIN_PASSWORD = "your-password"
$env:CHROME_BIN = "C:\Program Files\Google\Chrome\Application\chrome.exe"
cd "D:\AI Development\GO Automation"
go run .\cmd\main.go
```

Schedule it to run 2-3 times per day during business hours.

## 📈 Monitoring

### Check Daily Limits

```powershell
# View log output
Get-Content automation.log

# Check database
sqlite3 linkedin_automation.db "SELECT COUNT(*) FROM profiles WHERE date(created_at) = date('now')"
```

### Success Indicators

- ✅ Session restores successfully
- ✅ No CAPTCHAs after first run
- ✅ Steady connection acceptance rate
- ✅ No "unusual activity" warnings

## 📚 Additional Documentation

- [CAPTCHA_GUIDE.md](CAPTCHA_GUIDE.md) - CAPTCHA handling strategies
- [STEALTH_TECHNIQUES.md](STEALTH_TECHNIQUES.md) - Anti-detection details
- [IMPLEMENTATION_GUIDE.md](IMPLEMENTATION_GUIDE.md) - Technical implementation

## ⚠️ Disclaimer

This tool is for educational purposes. Use responsibly and in accordance with LinkedIn's Terms of Service. The authors are not responsible for any account restrictions or bans resulting from misuse.

## 🤝 Support

For issues or questions, check the troubleshooting section or review the configuration examples.

---

**Built with**: Go, Rod (Chrome DevTools Protocol), SQLite  
**Anti-Detection**: 8+ stealth techniques  
**CAPTCHA**: Session persistence + manual solving (free)  
**Configuration**: Complete YAML-based system
