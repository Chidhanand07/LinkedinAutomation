Here is a complete, clean, submission-ready README.md that matches your project exactly and defends all design decisions.

You can copy-paste this as-is into README.md.


LinkedIn Automation Proof-of-Concept (Go + Rod)

Project Overview

This project is a Go-based LinkedIn automation proof-of-concept built using the Rod browser automation library.
It demonstrates advanced browser automation techniques, human-like behavior simulation, anti-bot awareness, and clean modular Go architecture.

The goal of this project is not large-scale automation, but to showcase technical understanding of browser control, stealth strategies, and responsible automation design.


Objectives
	•	Demonstrate proficiency in Go and Rod
	•	Implement realistic human-like interaction patterns
	•	Showcase awareness of anti-bot detection mechanisms
	•	Design a clean, modular, maintainable architecture
	•	Persist state to avoid duplicate or unsafe actions
	•	Handle platform security checkpoints responsibly


Technologies Used
	•	Language: Go
	•	Browser Automation: Rod (Chrome DevTools Protocol)
	•	Browser: Google Chrome
	•	Storage: SQLite
	•	Configuration: YAML + Environment Variables
	•	IDE: GoLand


Project Architecture

The project follows a modular, responsibility-driven architecture:

cmd/         → Application entry point
browser/     → Browser launch and configuration
auth/        → LinkedIn authentication logic
search/      → Profile discovery
connect/     → Connection request automation
messaging/   → Follow-up messaging
stealth/     → Human-like behavior & fingerprint masking
storage/     → SQLite-based state persistence
config/      → YAML configuration handling
utils/       → Shared helpers (delay, retry)

Each module has a single responsibility, making the system easy to understand, extend, and maintain.


Core Features

Authentication
	•	Uses environment variables for credentials
	•	Automates LinkedIn login flow
	•	Detects security checkpoints (OTP, CAPTCHA)
	•	Halts execution safely if verification is required

Search & Targeting
	•	Navigates LinkedIn People Search
	•	Collects and deduplicates profile URLs
	•	Limits scope to realistic usage patterns

Connection Requests
	•	Visits individual profiles
	•	Scrolls and interacts naturally
	•	Sends connection requests with optional notes
	•	Enforces pacing to avoid aggressive behavior
	•	Tracks sent requests using SQLite

Messaging
	•	Navigates accepted connections
	•	Sends follow-up messages using a template
	•	Includes delays between actions
	•	Demonstrates post-connection automation flow


Stealth & Anti-Bot Techniques

This project implements multiple anti-detection and human-simulation strategies, including:
	•	Randomized delays between actions
	•	Natural scrolling behavior
	•	Human-like typing patterns
	•	Mouse hovering and movement
	•	Browser fingerprint masking (navigator.webdriver)
	•	User-Agent randomization
	•	Rate limiting between profiles
	•	Session reuse within a single run
	•	State persistence to prevent duplicate actions

These techniques are implemented at a proof-of-understanding level, not for bypassing safeguards.


State Persistence

The application uses SQLite to store:
	•	Processed profile URLs
	•	Timestamps of actions

This ensures:
	•	Duplicate requests are avoided
	•	Safe resumption after interruptions
	•	More realistic automation behavior


Error Handling & Safety
	•	Login failures are detected and logged
	•	Security checkpoints trigger a safe exit
	•	Missing environment variables stop execution
	•	No attempts are made to bypass CAPTCHA or OTP


Ethical Disclaimer

This project is intended strictly as a technical proof-of-concept for demonstrating browser automation, system design, and anti-detection awareness.

It does not attempt to bypass LinkedIn security mechanisms and intentionally halts execution when security checkpoints are detected.

The project is built for educational and evaluation purposes only.



Known Limitations
	•	Search parameters are static (can be extended)
	•	Pagination is limited
	•	Messaging templates are simple
	•	Cookie persistence across runs is not implemented

These limitations are intentional to keep the project focused on architecture and behavior simulation.


Conclusion

This project demonstrates:
	•	Strong understanding of Go-based browser automation
	•	Awareness of bot detection techniques
	•	Clean and maintainable software design
	•	Responsible handling of platform safeguards

It fulfills the assignment requirements as a robust and ethical proof-of-concept.


Author

Chidanandh R
B.E. – Artificial Intelligence & Machine Learning
RNS Institute of Technology
