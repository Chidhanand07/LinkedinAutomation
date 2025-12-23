# LinkedIn Automation - Complete Config-Driven System

## 🎉 New Features Implemented

### ✅ 1. Auto-Accept Pending Connections
Automatically accepts all incoming connection requests from your "My Network" page.

**Configuration** (`config/config.yaml`):
```yaml
auto_accept:
  enabled: true
  max_per_session: 50  # Max connections to accept per run
  send_welcome_message: true
  welcome_message_template: |
    Hi {firstName}, thanks for connecting!
```

### ✅ 2. Auto-Reply to Messages
Automatically replies to all unread messages with configured templates.

**Configuration**:
```yaml
messaging:
  auto_reply_enabled: true
  auto_reply_templates:
    - "Thanks for reaching out! I'll get back to you soon."
    - "Hi! I appreciate your message."
  reply_to_unread_only: true
  max_auto_replies_per_session: 20
```

### ✅ 3. Config-Driven Search
Search uses configured criteria from `config.yaml`:

**Configuration**:
```yaml
search:
  keywords: ["Software Engineer", "Backend Developer"]
  job_titles: ["Senior Software Engineer", "Lead Developer"]
  companies: ["Google", "Microsoft", "Amazon"]
  locations: ["United States", "San Francisco Bay Area"]
  max_pages: 5
  max_profiles: 50
```

### ✅ 4. Complete Config Integration
All modules now use centralized configuration from `config/config.yaml`.

## 🚀 Quick Start

### 1. Configure Settings

Edit [config/config.yaml](config/config.yaml):

```yaml
# Enable/disable features
auto_accept:
  enabled: true  # Auto-accept incoming connections

messaging:
  auto_reply_enabled: true  # Auto-reply to messages
  send_followup_to_new_connections: true

# Configure search
search:
  keywords: ["Your", "Keywords", "Here"]
  companies: ["Target", "Companies"]
  max_profiles: 50
```

### 2. Set Environment Variables

```powershell
# Required credentials
$env:LINKEDIN_EMAIL = "your-email@example.com"
$env:LINKEDIN_PASSWORD = "your-password"

# Required browser path
$env:CHROME_BIN = "C:\Program Files\Google\Chrome\Application\chrome.exe"

# Optional proxy
$env:PROXY = "http://proxy-server.com:8080"
```

### 3. Run Automation

```powershell
# Run the automation
go run .\cmd\main.go

# Or use the built executable
.\linkedin-automation.exe
```

## 📋 Execution Flow

The automation runs in 6 phases:

```
Phase 1: Authentication
  ├─ Session restore (if available)
  ├─ Login with stealth techniques
  └─ Save session for next run

Phase 2: Auto-Accept Connections
  ├─ Navigate to "My Network"
  ├─ Find pending requests
  └─ Accept up to max_per_session

Phase 3: Auto-Reply to Messages
  ├─ Navigate to "Messaging"
  ├─ Find unread messages
  └─ Reply with random template

Phase 4: Search for Profiles
  ├─ Build search URLs from config
  ├─ Paginate through results
  ├─ Apply filters (job titles, companies, locations)
  └─ Collect profile URLs

Phase 5: Send Connection Requests
  ├─ Visit each profile
  ├─ Check if already connected
  ├─ Send connection with optional note
  └─ Track in database

Phase 6: Send Follow-up Messages
  ├─ Find new connections
  ├─ Send personalized messages
  └─ Track messaged profiles
```

## 🎛️ Configuration Reference

### Auto-Accept Settings

```yaml
auto_accept:
  enabled: true                    # Enable/disable feature
  max_per_session: 50             # Max accepts per run
  send_welcome_message: true      # Send message after accepting
  welcome_message_template: |     # Message template
    Hi {firstName}, thanks for connecting!
```

### Search & Targeting

```yaml
search:
  # Search criteria
  keywords: ["Engineer", "Developer"]
  job_titles: ["Senior Engineer", "Lead"]
  companies: ["Google", "Microsoft"]
  locations: ["San Francisco", "New York"]
  
  # Filters
  connection_level: "2nd"          # 1st, 2nd, 3rd
  current_company_only: false
  
  # Limits
  max_pages: 5                     # Pages to scrape
  results_per_page: 10            # Results per page
  max_profiles: 50                # Total profiles to collect
  
  # Duplicate handling
  skip_duplicates: true
  skip_already_connected: true
  
  # Custom URLs (optional)
  custom_search_urls:
    - "https://www.linkedin.com/search/results/people/?keywords=developer"
```

### Connection Requests

```yaml
connection:
  # Daily limits
  max_requests_per_day: 30
  max_requests_per_hour: 10
  max_requests_per_session: 20
  
  # Note settings
  send_note: true
  note_required: false
  note_templates:
    - "Hi {{FirstName}}, let's connect!"
    - "Hey {{FirstName}}, I'd love to network!"
  
  # Timing
  delay_between_requests_min: 20   # seconds
  delay_between_requests_max: 60
  
  # Stealth behavior
  view_profile_before_connect: true
  scroll_profile_before_connect: true
  hover_on_elements: true
```

### Messaging System

```yaml
messaging:
  enabled: true
  
  # Auto-reply feature
  auto_reply_enabled: true
  auto_reply_templates:
    - "Thanks for reaching out!"
    - "I'll get back to you soon."
  reply_to_unread_only: true
  max_auto_replies_per_session: 20
  
  # Follow-up messages
  send_followup_to_new_connections: true
  wait_hours_after_acceptance: 24
  
  # Limits
  max_messages_per_day: 20
  max_messages_per_session: 15
  
  # Templates
  followup_templates:
    - "Thanks for connecting, {{FirstName}}!"
    - "Hi {{FirstName}}, great to be connected!"
  
  # Tracking
  track_responses: true
  skip_already_messaged: true
  
  # Timing
  delay_between_messages_min: 30   # seconds
  delay_between_messages_max: 90
```

### Rate Limiting

```yaml
rate_limiting:
  min_interval_between_runs_minutes: 30
  max_actions_per_day: 100
  operate_business_hours_only: true
  business_hours_start: 9          # 9 AM
  business_hours_end: 18           # 6 PM
  operate_weekdays_only: true
```

## 📊 Console Output Example

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🚀 LinkedIn Automation POC - Config-Driven
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✓ Configuration loaded successfully
✓ Database initialized
Using browser binary: C:/Program Files/Google/Chrome/Application/chrome.exe
✓ Browser initialized with anti-detection stealth techniques

📍 Phase 1: Authentication
🔄 Found existing session - attempting to restore...
✓ Session restored successfully - skipping login!

📍 Phase 2: Auto-Accept Pending Connections
🤝 Starting auto-accept pending connections...
Found 12 pending connection requests
  Accepting connection from: John Smith
  ⏳ Waiting 4s before next accept...
  Accepting connection from: Jane Doe
✓ Accepted 12 connection requests

📍 Phase 3: Auto-Reply to Messages
💬 Starting auto-reply to messages...
Found 5 unread message thread(s)
  Replying to: Mike Johnson
  Typing reply...
  ✓ Reply sent to Mike Johnson
  ⏳ Waiting 45s before next reply...
✓ Sent 5 auto-replies

📍 Phase 4: Search & Profile Discovery
🔍 Starting profile search with configured filters...
Search criteria: Keywords=[Software Engineer], Companies=[Google Microsoft]
  Page 1: https://www.linkedin.com/search/results/people/...
  ✓ Found 10 profiles on page 1
✓ Found 50 total profiles

📍 Phase 5: Sending Connection Requests
  [1/20] Sending request to: https://www.linkedin.com/in/john-doe
  ✓ Connection request sent
  ⏳ Waiting 25s before next request...
✓ Reached session limit: 20 connections

📍 Phase 6: Sending Follow-up Messages
📨 Starting follow-up messages to new connections...
✓ Sent 10 follow-up messages

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 Session Summary:
   • Profiles found: 50
   • Connection requests sent: 20
   • Daily connection requests: 20/30
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✅ Automation run completed successfully
```

## 🎯 Best Practices

### 1. Start Conservative

**First Week**:
```yaml
connection:
  max_requests_per_day: 15
auto_accept:
  max_per_session: 20
messaging:
  max_messages_per_day: 10
```

**Second Week** (if no issues):
```yaml
connection:
  max_requests_per_day: 25
auto_accept:
  max_per_session: 40
messaging:
  max_messages_per_day: 15
```

**Third Week+**:
```yaml
connection:
  max_requests_per_day: 30-50
auto_accept:
  max_per_session: 50
messaging:
  max_messages_per_day: 20
```

### 2. Use Business Hours

```yaml
rate_limiting:
  operate_business_hours_only: true
  business_hours_start: 9
  business_hours_end: 18
  operate_weekdays_only: true
```

### 3. Vary Your Templates

Add multiple templates for natural variation:

```yaml
messaging:
  auto_reply_templates:
    - "Thanks for reaching out! I'll get back to you soon."
    - "Hi! I appreciate your message. I'll respond shortly."
    - "Thanks for the message! I'm reviewing it now."
    - "Hello! I'll be in touch soon. Thanks for connecting!"
```

### 4. Monitor Daily Limits

The automation tracks and enforces daily limits automatically:

```
📊 Daily connection requests: 28/30
⚠️  Daily limit reached - stopping
```

### 5. Use Proxies for Multiple Accounts

```powershell
# Residential proxy
$env:PROXY = "http://username:password@residential-proxy.com:8080"

# Rotate proxies between runs
```

## 🛠️ Troubleshooting

### Auto-Accept Not Working

1. **Check config**: Ensure `auto_accept.enabled: true`
2. **Check page**: LinkedIn may have changed selectors
3. **Check limits**: May have reached `max_per_session`

### Auto-Reply Not Sending

1. **Check config**: Ensure `auto_reply_enabled: true`
2. **Check templates**: Must have at least one template
3. **No unread messages**: Only replies to unread messages

### Search Not Finding Profiles

1. **Broaden criteria**: Reduce filters in config
2. **Check keywords**: Try different keywords
3. **Increase pages**: Set `max_pages` higher

### Rate Limits Hit

Reduce daily limits:
```yaml
connection:
  max_requests_per_day: 15  # Reduce from 30
messaging:
  max_messages_per_day: 10   # Reduce from 20
```

## 📁 Project Structure

```
d:\AI Development\GO Automation\
├── cmd/
│   └── main.go                 # Main orchestration (all 6 phases)
├── config/
│   ├── config.yaml             # Central configuration
│   └── config.go               # Config loader
├── accept/
│   └── accept.go               # ✨ NEW: Auto-accept module
├── auth/
│   └── login.go                # Authentication + session
├── browser/
│   └── browser.go              # Browser initialization
├── connect/
│   └── connect.go              # Connection requests
├── messaging/
│   └── messaging.go            # ✨ UPDATED: Auto-reply + follow-ups
├── search/
│   └── search.go               # ✨ UPDATED: Config-driven search
├── stealth/
│   ├── captcha.go              # Session + CAPTCHA handling
│   ├── fingerprint.go          # Browser fingerprinting
│   ├── mouse.go                # Mouse movements
│   ├── typing.go               # Typing simulation
│   └── scroll.go               # Scrolling behavior
└── storage/
    └── sqlite.go               # Database operations
```

## 🎓 Template Variables

Use these in your message templates:

```yaml
note_templates:
  - "Hi {{FirstName}}, let's connect!"           # First name
  - "Hey {{FirstName}} from {{Company}}!"        # Company
  - "Interested in {{Industry}}? Let's talk!"    # Industry
  - "Your work with {{Skill}} is impressive!"    # Skills
```

**Note**: Variable replacement requires profile data extraction (not fully implemented in current version).

## 🔐 Security & Privacy

1. **Never commit credentials**: Use environment variables only
2. **Secure config.yaml**: Add to `.gitignore` if it contains sensitive data
3. **Use session files**: Stored in `linkedin_session.json` (add to `.gitignore`)
4. **Proxy recommended**: Use residential proxies for privacy

## 📈 Success Metrics

Monitor your automation's health:

### Good Indicators ✅
- Session restored 80%+ of runs
- No CAPTCHAs for 7+ days
- 90%+ connection acceptance rate
- No security warnings

### Warning Signs ⚠️
- Frequent CAPTCHAs
- Low acceptance rate (<50%)
- Account restrictions
- Messages not sending

## 🎉 Summary

This automation now provides:

✅ **Full config-driven operation** - Everything controlled from `config.yaml`
✅ **Auto-accept connections** - Automatically accept incoming requests
✅ **Auto-reply to messages** - Respond to all unread messages
✅ **Smart search** - Filter by keywords, titles, companies, locations
✅ **Duplicate detection** - Prevent re-processing profiles
✅ **Session persistence** - Skip login after first run
✅ **Rate limiting** - Stay within safe daily limits
✅ **Comprehensive logging** - Track all actions

All features work together seamlessly with human-like stealth techniques!
