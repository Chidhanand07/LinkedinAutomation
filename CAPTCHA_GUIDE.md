# CAPTCHA Handling & Anti-Detection Guide

## 🎯 Implemented Solutions

### 1. ✅ Session Persistence (Primary Strategy)
**Avoids login and CAPTCHA completely by reusing authenticated sessions**

- Saves cookies after successful login to `linkedin_session.json`
- Automatically restores session on next run
- Only logs in when session expires or is invalid

**Benefits**:
- No CAPTCHA after first successful login
- Faster execution (skips entire login flow)
- Maintains trust signals with LinkedIn

### 2. ✅ Manual CAPTCHA Solving (Free Method)
**Pauses automation for human intervention when CAPTCHA appears**

- Detects multiple CAPTCHA types: reCAPTCHA, hCAPTCHA, custom challenges
- Provides clear console warnings with instructions
- Waits up to 3 minutes (configurable) for manual solve
- Automatically continues after CAPTCHA is solved

**Supported Challenge Types**:
- CAPTCHA (all types)
- PIN verification
- Email verification
- Suspicious activity warnings

### 3. ✅ Enhanced Anti-Detection (Reduce CAPTCHA Frequency)
**Persistent browser profile maintains trust signals**

- **Persistent profile directory**: `./chrome-profile/`
- **Browser history**: Maintains previous session data
- **Trust signals**: Looks like a real user's browser
- **Additional flags**: Reduces automation signatures

**New Stealth Flags**:
```
--disable-features=IsolateOrigins,site-per-process
--user-data-dir=./chrome-profile
--window-size=1920,1080
```

### 4. ✅ Rate Limiting & Behavioral Patterns
**Prevents account flagging that triggers CAPTCHAs**

- **Minimum interval**: 30 minutes between automation runs
- **Daily limits**: 30 connection requests per day (configurable)
- **Variable delays**: 20-60 seconds between actions
- **Action tracking**: Monitors daily activity counts

### 5. ✅ Proxy Support
**Use residential IPs to avoid datacenter IP detection**

**Setup**:
```powershell
# Set proxy via environment variable
$env:PROXY = "http://username:password@proxy-server.com:8080"

# Or for SOCKS5
$env:PROXY = "socks5://username:password@proxy-server.com:1080"

# Then run
go run .\cmd\main.go
```

**Recommended Proxy Types**:
- ✅ Residential proxies (best for avoiding detection)
- ✅ Mobile proxies (excellent trust signals)
- ⚠️ Datacenter proxies (higher CAPTCHA rate)
- ❌ Free proxies (usually flagged)

## 📋 Usage Instructions

### First Run (With Login)
```powershell
# Set environment variables
$env:LINKEDIN_EMAIL = "your-email@example.com"
$env:LINKEDIN_PASSWORD = "your-password"
$env:CHROME_BIN = "C:\Program Files\Google\Chrome\Application\chrome.exe"

# Optional: Set proxy
$env:PROXY = "http://proxy-server.com:8080"

# Run automation
go run .\cmd\main.go
```

**If CAPTCHA Appears**:
1. Automation will pause with clear warning message
2. Solve the CAPTCHA manually in the browser window
3. Automation continues automatically after solving
4. Session is saved for future runs

### Subsequent Runs (Session Restored)
```powershell
# Just run - no login needed!
go run .\cmd\main.go
```

**Output**:
```
🔄 Found existing session - attempting to restore...
✓ Session restored successfully - skipping login!
```

## 🔧 Configuration Options

### Rate Limiter Settings
Edit [cmd/main.go](cmd/main.go):

```go
// Adjust minimum interval between runs (default: 30 minutes)
rateLimiter := stealth.NewRateLimiter(30 * time.Minute)

// Adjust daily connection limit (default: 30)
maxConnectionsPerDay := 30
```

### CAPTCHA Timeout
Edit [auth/login.go](auth/login.go):

```go
// Adjust manual solve timeout (default: 180 seconds)
stealth.WaitForManualCaptchaSolve(page, 180*time.Second)
```

### Session File Location
Edit [auth/login.go](auth/login.go):

```go
// Change session file path
sessionMgr := stealth.NewSessionManager("linkedin_session.json")
```

## 🛡️ Best Practices to Avoid CAPTCHAs

### 1. Account Warm-up
- Use the account manually for 3-5 days before automation
- Complete profile to 100%
- Build 50+ connections organically
- Join relevant groups

### 2. Gradual Activity Increase
**Week 1**: 5-10 actions/day
**Week 2**: 10-20 actions/day
**Week 3**: 20-30 actions/day
**Week 4+**: Full automation (30-50 actions/day)

### 3. Business Hours Operation
The system automatically adds extra delays outside 9 AM - 6 PM M-F.

**Schedule runs during business hours**:
```powershell
# Windows Task Scheduler - run at 10 AM weekdays
```

### 4. Use Residential Proxies
- Rotate IPs periodically
- Use proxies from your target region
- Avoid sharing proxies across multiple accounts

### 5. Monitor Logs
Watch for patterns:
```
⚠️  Daily limit reached
🔒 Security challenge detected: CAPTCHA
⏱️  Rate limiting: waiting 25m before next action
```

## 📊 Monitoring Session Health

### Check Session File
```powershell
# View saved cookies
Get-Content linkedin_session.json | ConvertFrom-Json
```

### Session Indicators
- ✅ **Large file (>5KB)**: Healthy session with many cookies
- ⚠️ **Small file (<2KB)**: Partial session, may expire soon
- ❌ **Missing file**: No session saved, will require login

### Force Fresh Login
```powershell
# Delete session and login again
Remove-Item linkedin_session.json
go run .\cmd\main.go
```

## 🚨 Troubleshooting

### CAPTCHA Keeps Appearing
1. **Use persistent profile**: Check `./chrome-profile/` exists
2. **Reduce activity**: Lower `maxConnectionsPerDay` to 15-20
3. **Add proxy**: Use residential IP via `$env:PROXY`
4. **Increase delays**: Edit rate limiter to 60 minutes
5. **Manual warm-up**: Use account manually for a week

### Session Not Restoring
1. **Check file**: Ensure `linkedin_session.json` exists
2. **Cookies expired**: Natural after 7-30 days - just login again
3. **IP changed**: Session tied to IP - use same proxy/network
4. **Account security**: LinkedIn may have reset session

### Browser Profile Issues
```powershell
# Clear profile and start fresh
Remove-Item -Recurse -Force .\chrome-profile
```

## 📈 Success Metrics

**Good Indicators**:
- ✅ Session restored 80%+ of runs
- ✅ No CAPTCHAs for 7+ days
- ✅ Consistent daily activity accepted
- ✅ No security warnings

**Warning Signs**:
- ⚠️ CAPTCHA every login
- ⚠️ Frequent PIN verifications  
- ⚠️ "Unusual activity" warnings
- ⚠️ Session expires daily

## 🎓 Advanced: Custom CAPTCHA Solving

For fully unattended automation, integrate paid service (not implemented, but here's the approach):

```go
// Example with 2Captcha (requires API key)
import "github.com/2captcha/2captcha-go"

func SolveRecaptcha(page *rod.Page, apiKey string) error {
    client := api2captcha.NewClient(apiKey)
    // Implementation here
}
```

**Cost**: ~$1-3 per 1000 CAPTCHAs

**Free alternatives** (current implementation):
- Manual solving (pause + human intervention) ✅ Implemented
- Session persistence (avoid CAPTCHAs entirely) ✅ Implemented

---

**Summary**: Session persistence + manual solving provides a completely free, effective solution for CAPTCHA handling while maintaining high success rates.
