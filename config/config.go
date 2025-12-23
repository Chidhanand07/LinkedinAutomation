package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds all automation configuration
type Config struct {
	Credentials   CredentialsConfig   `yaml:"credentials"`
	Search        SearchConfig        `yaml:"search"`
	AutoAccept    AutoAcceptConfig    `yaml:"auto_accept"`
	Connection    ConnectionConfig    `yaml:"connection"`
	Messaging     MessagingConfig     `yaml:"messaging"`
	RateLimiting  RateLimitingConfig  `yaml:"rate_limiting"`
	Browser       BrowserConfig       `yaml:"browser"`
	Storage       StorageConfig       `yaml:"storage"`
	Logging       LoggingConfig       `yaml:"logging"`
	Advanced      AdvancedConfig      `yaml:"advanced"`
}

type CredentialsConfig struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

type SearchConfig struct {
	Keywords            []string `yaml:"keywords"`
	JobTitles           []string `yaml:"job_titles"`
	Companies           []string `yaml:"companies"`
	Locations           []string `yaml:"locations"`
	ConnectionLevel     string   `yaml:"connection_level"`
	CurrentCompanyOnly  bool     `yaml:"current_company_only"`
	PastCompanyOnly     bool     `yaml:"past_company_only"`
	MaxPages            int      `yaml:"max_pages"`
	ResultsPerPage      int      `yaml:"results_per_page"`
	MaxProfiles         int      `yaml:"max_profiles"`
	SkipDuplicates      bool     `yaml:"skip_duplicates"`
	SkipAlreadyConnected bool    `yaml:"skip_already_connected"`
	SkipPendingRequests bool     `yaml:"skip_pending_requests"`
	CustomSearchURLs    []string `yaml:"custom_search_urls"`
}

type AutoAcceptConfig struct {
	Enabled                bool   `yaml:"enabled"`
	MaxPerSession          int    `yaml:"max_per_session"`
	SendWelcomeMessage     bool   `yaml:"send_welcome_message"`
	WelcomeMessageTemplate string `yaml:"welcome_message_template"`
}

type ConnectionConfig struct {
	MaxRequestsPerDay        int      `yaml:"max_requests_per_day"`
	MaxRequestsPerHour       int      `yaml:"max_requests_per_hour"`
	MaxRequestsPerSession    int      `yaml:"max_requests_per_session"`
	SendNote                 bool     `yaml:"send_note"`
	NoteRequired             bool     `yaml:"note_required"`
	NoteTemplates            []string `yaml:"note_templates"`
	MaxNoteLength            int      `yaml:"max_note_length"`
	DelayBetweenRequestsMin  int      `yaml:"delay_between_requests_min"`
	DelayBetweenRequestsMax  int      `yaml:"delay_between_requests_max"`
	ViewProfileBeforeConnect bool     `yaml:"view_profile_before_connect"`
	ScrollProfileBeforeConnect bool   `yaml:"scroll_profile_before_connect"`
	HoverOnElements          bool     `yaml:"hover_on_elements"`
}

type MessagingConfig struct {
	Enabled                     bool     `yaml:"enabled"`
	CheckForNewConnections      bool     `yaml:"check_for_new_connections"`
	SendFollowupToNewConnections bool    `yaml:"send_followup_to_new_connections"`
	AutoReplyEnabled            bool     `yaml:"auto_reply_enabled"`
	AutoReplyTemplates          []string `yaml:"auto_reply_templates"`
	ReplyToUnreadOnly           bool     `yaml:"reply_to_unread_only"`
	MaxAutoRepliesPerSession    int      `yaml:"max_auto_replies_per_session"`
	WaitHoursAfterAcceptance    int      `yaml:"wait_hours_after_acceptance"`
	CheckIntervalMinutes        int      `yaml:"check_interval_minutes"`
	MaxMessagesPerDay           int      `yaml:"max_messages_per_day"`
	MaxMessagesPerHour          int      `yaml:"max_messages_per_hour"`
	MaxMessagesPerSession       int      `yaml:"max_messages_per_session"`
	FollowupTemplates           []string `yaml:"followup_templates"`
	TrackResponses              bool     `yaml:"track_responses"`
	MarkAsMessaged              bool     `yaml:"mark_as_messaged"`
	SkipAlreadyMessaged         bool     `yaml:"skip_already_messaged"`
	MaxMessageLength            int      `yaml:"max_message_length"`
	DelayBetweenMessagesMin     int      `yaml:"delay_between_messages_min"`
	DelayBetweenMessagesMax     int      `yaml:"delay_between_messages_max"`
}

type RateLimitingConfig struct {
	MinIntervalBetweenRunsMinutes int  `yaml:"min_interval_between_runs_minutes"`
	MaxActionsPerDay              int  `yaml:"max_actions_per_day"`
	ProfileViewsPerDay            int  `yaml:"profile_views_per_day"`
	SearchesPerDay                int  `yaml:"searches_per_day"`
	PageLoadDelayMin              int  `yaml:"page_load_delay_min"`
	PageLoadDelayMax              int  `yaml:"page_load_delay_max"`
	ActionDelayMin                int  `yaml:"action_delay_min"`
	ActionDelayMax                int  `yaml:"action_delay_max"`
	OperateBusinessHoursOnly      bool `yaml:"operate_business_hours_only"`
	BusinessHoursStart            int  `yaml:"business_hours_start"`
	BusinessHoursEnd              int  `yaml:"business_hours_end"`
	OperateWeekdaysOnly           bool `yaml:"operate_weekdays_only"`
}

type BrowserConfig struct {
	Headless     bool   `yaml:"headless"`
	UserDataDir  string `yaml:"user_data_dir"`
	UseProxy     bool   `yaml:"use_proxy"`
	ProxyURL     string `yaml:"proxy_url"`
	WindowWidth  int    `yaml:"window_width"`
	WindowHeight int    `yaml:"window_height"`
	ChromeBin    string `yaml:"chrome_bin"`
}

type StorageConfig struct {
	DatabaseFile       string `yaml:"database_file"`
	SessionFile        string `yaml:"session_file"`
	AutoBackup         bool   `yaml:"auto_backup"`
	BackupIntervalDays int    `yaml:"backup_interval_days"`
	KeepHistoryDays    int    `yaml:"keep_history_days"`
}

type LoggingConfig struct {
	Level        string `yaml:"level"`
	LogFile      string `yaml:"log_file"`
	LogToConsole bool   `yaml:"log_to_console"`
	LogToFile    bool   `yaml:"log_to_file"`
}

type AdvancedConfig struct {
	MinConnections              int  `yaml:"min_connections"`
	MaxConnections              int  `yaml:"max_connections"`
	MustHaveProfilePicture      bool `yaml:"must_have_profile_picture"`
	MustHaveHeadline            bool `yaml:"must_have_headline"`
	PrioritizeRecentActivity    bool `yaml:"prioritize_recent_activity"`
	PrioritizeMutualConnections bool `yaml:"prioritize_mutual_connections"`
	RetryFailedActions          bool `yaml:"retry_failed_actions"`
	MaxRetries                  int  `yaml:"max_retries"`
	ContinueOnError             bool `yaml:"continue_on_error"`
	ScreenshotOnError           bool `yaml:"screenshot_on_error"`
	ScreenshotDir               string `yaml:"screenshot_dir"`
}

var GlobalConfig *Config

// LoadConfig reads and parses the configuration file
func LoadConfig(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Override with environment variables
	if email := os.Getenv("LINKEDIN_EMAIL"); email != "" {
		cfg.Credentials.Email = email
	}
	if password := os.Getenv("LINKEDIN_PASSWORD"); password != "" {
		cfg.Credentials.Password = password
	}
	if chromeBin := os.Getenv("CHROME_BIN"); chromeBin != "" {
		cfg.Browser.ChromeBin = chromeBin
	}
	if proxyURL := os.Getenv("PROXY"); proxyURL != "" {
		cfg.Browser.ProxyURL = proxyURL
		cfg.Browser.UseProxy = true
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	GlobalConfig = &cfg
	return &cfg, nil
}

// Validate checks if configuration is valid
func (c *Config) Validate() error {
	if c.Credentials.Email == "" || c.Credentials.Password == "" {
		return fmt.Errorf("credentials are required (set LINKEDIN_EMAIL and LINKEDIN_PASSWORD)")
	}

	if c.Search.MaxProfiles <= 0 {
		return fmt.Errorf("search.max_profiles must be greater than 0")
	}

	if c.Connection.MaxRequestsPerDay <= 0 {
		return fmt.Errorf("connection.max_requests_per_day must be greater than 0")
	}

	if len(c.Connection.NoteTemplates) == 0 && c.Connection.SendNote {
		return fmt.Errorf("connection.note_templates cannot be empty when send_note is true")
	}

	if c.Messaging.Enabled && len(c.Messaging.FollowupTemplates) == 0 {
		return fmt.Errorf("messaging.followup_templates cannot be empty when enabled is true")
	}

	return nil
}

// GetConfig returns the global configuration
func GetConfig() *Config {
	if GlobalConfig == nil {
		panic("Configuration not loaded. Call LoadConfig() first.")
	}
	return GlobalConfig
}
