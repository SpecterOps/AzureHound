// models/azure/intune_sessions.go
package azure

import (
	"time"
)

// DeviceSessionData combines device information with collected session data
type DeviceSessionData struct {
	Device      IntuneDevice `json:"device"`
	SessionData SessionData  `json:"sessionData"`
	CollectedAt time.Time    `json:"collectedAt"`

	// BloodHound specific fields for integration
	BloodHoundData BloodHoundSessionData `json:"bloodhoundData"`
}

// SessionData represents the complete session data collected from a device
type SessionData struct {
	DeviceInfo         DeviceInfo          `json:"deviceInfo"`
	ActiveSessions     []ActiveSession     `json:"activeSessions"`
	LoggedOnUsers      []LoggedOnUser      `json:"loggedOnUsers"`
	SecurityIndicators SessionSecurityInfo `json:"securityIndicators"`
	Summary            SessionDataSummary  `json:"summary"`
}

// ActiveSession represents an active user session on the device
type ActiveSession struct {
	SessionID     int       `json:"sessionId"`
	UserName      string    `json:"userName"`
	DomainName    string    `json:"domainName"`
	SessionType   string    `json:"sessionType"`  // Console, RDP, etc.
	SessionState  string    `json:"sessionState"` // Active, Disconnected, etc.
	LogonTime     time.Time `json:"logonTime"`
	IdleTime      string    `json:"idleTime"`
	ClientName    string    `json:"clientName"`    // For RDP sessions
	ClientAddress string    `json:"clientAddress"` // For RDP sessions
	ProcessCount  int       `json:"processCount"`
	IsElevated    bool      `json:"isElevated"` // Running as administrator
}

// LoggedOnUser represents a user with credentials cached on the system
type LoggedOnUser struct {
	UserName         string    `json:"userName"`
	DomainName       string    `json:"domainName"`
	SID              string    `json:"sid"`
	LogonType        string    `json:"logonType"`   // Interactive, Network, Service, etc.
	AuthPackage      string    `json:"authPackage"` // NTLM, Kerberos, etc.
	LogonTime        time.Time `json:"logonTime"`
	LogonServer      string    `json:"logonServer"`
	HasCachedCreds   bool      `json:"hasCachedCreds"`
	IsServiceAccount bool      `json:"isServiceAccount"`
	TokenPrivileges  []string  `json:"tokenPrivileges"`
}

// SessionSecurityInfo contains analysis results of session-related security indicators
type SessionSecurityInfo struct {
	AdminSessionsActive     bool                 `json:"adminSessionsActive"`
	RemoteSessionsActive    bool                 `json:"remoteSessionsActive"`
	ServiceAccountSessions  bool                 `json:"serviceAccountSessions"`
	SuspiciousLogonTypes    []string             `json:"suspiciousLogonTypes"`
	CredentialTheftRisk     string               `json:"credentialTheftRisk"`     // Low, Medium, High, Critical
	PrivilegeEscalationRisk string               `json:"privilegeEscalationRisk"` // Low, Medium, High, Critical
	SuspiciousActivities    []SuspiciousActivity `json:"suspiciousActivities"`
}

// SuspiciousActivity represents potentially malicious session activity
type SuspiciousActivity struct {
	ActivityType string    `json:"activityType"` // Multiple_Admin_Sessions, Unusual_Logon_Time, etc.
	Description  string    `json:"description"`
	RiskLevel    string    `json:"riskLevel"` // Low, Medium, High, Critical
	Evidence     []string  `json:"evidence"`
	DetectedAt   time.Time `json:"detectedAt"`
	UserName     string    `json:"userName"`
	SessionID    int       `json:"sessionId"`
}

// SessionDataSummary provides high-level statistics about the collected session data
type SessionDataSummary struct {
	TotalActiveSessions int      `json:"totalActiveSessions"`
	UniqueUsers         int      `json:"uniqueUsers"`
	AdminSessions       int      `json:"adminSessions"`
	RemoteSessions      int      `json:"remoteSessions"`
	ServiceSessions     int      `json:"serviceSessions"`
	HighRiskIndicators  []string `json:"highRiskIndicators"`
	CredentialExposure  int      `json:"credentialExposure"` // Number of exposed credential sets
}

// BloodHoundSessionData contains processed session data formatted for BloodHound consumption
type BloodHoundSessionData struct {
	ObjectIdentifier   string                    `json:"ObjectIdentifier"`
	AzureDeviceID      string                    `json:"AzureDeviceID"`
	DisplayName        string                    `json:"DisplayName"`
	Sessions           []BloodHoundSession       `json:"Sessions"`
	LoggedOnUsers      []BloodHoundLoggedOnUser  `json:"LoggedOnUsers"`
	CredentialExposure []CredentialExposure      `json:"CredentialExposure"`
	SessionFindings    []SessionSecurityFinding  `json:"SessionFindings"`
	EscalationVectors  []SessionEscalationVector `json:"EscalationVectors"`
}

// BloodHoundSession represents a session in BloodHound format
type BloodHoundSession struct {
	UserSID     string                 `json:"UserSID"`
	UserName    string                 `json:"UserName"`
	DomainName  string                 `json:"DomainName"`
	ComputerSID string                 `json:"ComputerSID"`
	SessionType string                 `json:"SessionType"`
	LogonType   string                 `json:"LogonType"`
	IsElevated  bool                   `json:"IsElevated"`
	LogonTime   time.Time              `json:"LogonTime"`
	ClientName  string                 `json:"ClientName"`
	Properties  map[string]interface{} `json:"Properties"`
}

// BloodHoundLoggedOnUser represents a logged-on user in BloodHound format
type BloodHoundLoggedOnUser struct {
	UserSID         string                 `json:"UserSID"`
	UserName        string                 `json:"UserName"`
	DomainName      string                 `json:"DomainName"`
	ComputerSID     string                 `json:"ComputerSID"`
	LogonType       string                 `json:"LogonType"`
	AuthPackage     string                 `json:"AuthPackage"`
	HasCachedCreds  bool                   `json:"HasCachedCreds"`
	TokenPrivileges []string               `json:"TokenPrivileges"`
	Properties      map[string]interface{} `json:"Properties"`
}

// CredentialExposure represents credentials that could be harvested from sessions
type CredentialExposure struct {
	UserName         string   `json:"userName"`
	DomainName       string   `json:"domainName"`
	SID              string   `json:"sid"`
	ExposureType     string   `json:"exposureType"`     // Interactive, Cached, Service, etc.
	ExposureRisk     string   `json:"exposureRisk"`     // Low, Medium, High, Critical
	ExposureLocation string   `json:"exposureLocation"` // LSASS, Registry, Memory, etc.
	HarvestMethods   []string `json:"harvestMethods"`   // Mimikatz, ProcDump, etc.
	TargetPrivileges []string `json:"targetPrivileges"`
}

// SessionSecurityFinding represents session-based security issues
type SessionSecurityFinding struct {
	ID              string   `json:"ID"`
	Title           string   `json:"Title"`
	Severity        string   `json:"Severity"`
	Category        string   `json:"Category"`
	Description     string   `json:"Description"`
	Evidence        []string `json:"Evidence"`
	Recommendations []string `json:"Recommendations"`
	MITREAttack     []string `json:"MITREAttack"`
	AffectedUsers   []string `json:"AffectedUsers"`
	SessionIDs      []int    `json:"SessionIDs"`
}

// SessionEscalationVector represents privilege escalation paths through sessions
type SessionEscalationVector struct {
	VectorID         string   `json:"VectorID"`
	Type             string   `json:"Type"`           // Session_Hijacking, Credential_Theft, etc.
	Source           string   `json:"Source"`         // Current user/privilege level
	Target           string   `json:"Target"`         // Target user/privilege level
	Method           string   `json:"Method"`         // Token_Impersonation, Credential_Dumping, etc.
	RequiredAccess   []string `json:"RequiredAccess"` // Local_Logon, Debug_Privilege, etc.
	Complexity       string   `json:"Complexity"`     // Low, Medium, High
	Impact           string   `json:"Impact"`         // Low, Medium, High, Critical
	Conditions       []string `json:"Conditions"`     // Session_Active, Admin_Privileges, etc.
	TechnicalDetails string   `json:"TechnicalDetails"`
	SessionID        int      `json:"SessionID"`
	TargetUserSID    string   `json:"TargetUserSID"`
}

// DeviceSessionAnalysis represents the complete session security analysis for a device
type DeviceSessionAnalysis struct {
	Device              IntuneDevice              `json:"device"`
	AnalysisTimestamp   time.Time                 `json:"analysisTimestamp"`
	SessionFindings     []SessionSecurityFinding  `json:"sessionFindings"`
	EscalationVectors   []SessionEscalationVector `json:"escalationVectors"`
	CredentialExposures []CredentialExposure      `json:"credentialExposures"`
	BloodHoundData      BloodHoundSessionData     `json:"bloodhoundData"`
	RiskScore           int                       `json:"riskScore"`
	SecurityPosture     string                    `json:"securityPosture"` // Secure, Moderate, High_Risk, Critical
	LastUpdated         time.Time                 `json:"lastUpdated"`
}

// SessionCollectionScript represents the PowerShell script for session data collection
type SessionCollectionScript struct {
	ID                   string    `json:"id"`
	DisplayName          string    `json:"displayName"`
	Description          string    `json:"description"`
	ScriptContent        string    `json:"scriptContent"`
	CreatedDateTime      time.Time `json:"createdDateTime"`
	LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
	RunAsAccount         string    `json:"runAsAccount"`
	FileName             string    `json:"fileName"`
	CollectionMethods    []string  `json:"collectionMethods"` // quser, wmic, net, etc.
	RequiredPrivileges   []string  `json:"requiredPrivileges"`
}

// SessionMonitoringConfiguration represents configuration for session monitoring
type SessionMonitoringConfiguration struct {
	EnableSessionCollection  bool     `json:"enableSessionCollection"`
	EnableCredentialAnalysis bool     `json:"enableCredentialAnalysis"`
	EnablePrivilegeAnalysis  bool     `json:"enablePrivilegeAnalysis"`
	MonitorServiceAccounts   bool     `json:"monitorServiceAccounts"`
	AlertOnAdminSessions     bool     `json:"alertOnAdminSessions"`
	AlertOnRemoteSessions    bool     `json:"alertOnRemoteSessions"`
	ExcludedUsers            []string `json:"excludedUsers"`
	ExcludedServiceAccounts  []string `json:"excludedServiceAccounts"`
	HighRiskSessionThreshold int      `json:"highRiskSessionThreshold"`
	CollectionInterval       string   `json:"collectionInterval"` // Daily, Weekly, etc.
	RetentionPeriod          string   `json:"retentionPeriod"`    // 30d, 90d, etc.
}

// SessionComplianceReport represents a comprehensive session security report
type SessionComplianceReport struct {
	TenantID                 string                    `json:"tenantId"`
	ReportTimestamp          time.Time                 `json:"reportTimestamp"`
	TotalDevices             int                       `json:"totalDevices"`
	DevicesWithSessions      int                       `json:"devicesWithSessions"`
	DevicesWithAdminSessions int                       `json:"devicesWithAdminSessions"`
	SessionBreakdown         SessionBreakdown          `json:"sessionBreakdown"`
	SecurityFindings         []SessionSecurityFinding  `json:"securityFindings"`
	TopRisks                 []SessionRiskSummary      `json:"topRisks"`
	CredentialExposureTrend  []CredentialExposureTrend `json:"credentialExposureTrend"`
	Recommendations          []SessionRecommendation   `json:"recommendations"`
}

// SessionBreakdown provides statistics about session types and security posture
type SessionBreakdown struct {
	TotalActiveSessions int `json:"totalActiveSessions"`
	InteractiveSessions int `json:"interactiveSessions"`
	RemoteSessions      int `json:"remoteSessions"`
	ServiceSessions     int `json:"serviceSessions"`
	AdminSessions       int `json:"adminSessions"`
	ElevatedSessions    int `json:"elevatedSessions"`
	SuspiciousSessions  int `json:"suspiciousSessions"`
}

// SessionRiskSummary represents a high-level session risk category summary
type SessionRiskSummary struct {
	RiskCategory       string  `json:"riskCategory"`
	AffectedDevices    int     `json:"affectedDevices"`
	AffectedSessions   int     `json:"affectedSessions"`
	ExposedCredentials int     `json:"exposedCredentials"`
	Severity           string  `json:"severity"`
	Description        string  `json:"description"`
	ImpactScore        float64 `json:"impactScore"`
}

// CredentialExposureTrend represents credential exposure trends over time
type CredentialExposureTrend struct {
	Date                time.Time `json:"date"`
	TotalExposedCreds   int       `json:"totalExposedCreds"`
	AdminCredsExposed   int       `json:"adminCredsExposed"`
	ServiceCredsExposed int       `json:"serviceCredsExposed"`
	HighRiskExposures   int       `json:"highRiskExposures"`
	ExposureRate        float64   `json:"exposureRate"`
}

// SessionRecommendation represents actionable session security recommendations
type SessionRecommendation struct {
	ID               string   `json:"id"`
	Title            string   `json:"title"`
	Priority         string   `json:"priority"`
	Category         string   `json:"category"`
	Description      string   `json:"description"`
	Impact           string   `json:"impact"`
	Implementation   string   `json:"implementation"`
	AffectedDevices  []string `json:"affectedDevices"`
	AffectedUsers    []string `json:"affectedUsers"`
	EstimatedEffort  string   `json:"estimatedEffort"`
	MITREMitigations []string `json:"mitreMitigations"`
	TechnicalDetails string   `json:"technicalDetails"`
}
