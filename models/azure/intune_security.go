// models/azure/intune_security.go
package azure

import (
	"time"
)

// DeviceSecurityAnalysis represents the complete security analysis for a device
type DeviceSecurityAnalysis struct {
	Device            IntuneDevice         `json:"device"`
	AnalysisTimestamp time.Time            `json:"analysisTimestamp"`
	SecurityFindings  []SecurityFinding    `json:"securityFindings"`
	EscalationVectors []EscalationVector   `json:"escalationVectors"`
	BloodHoundData    BloodHoundDeviceData `json:"bloodhoundData"`
	RiskScore         int                  `json:"riskScore"`
	ComplianceStatus  string               `json:"complianceStatus"`
	LastUpdated       time.Time            `json:"lastUpdated"`
}

// IntuneSecurityConfiguration represents security configuration collected from Intune
type IntuneSecurityConfiguration struct {
	DeviceID                  string                    `json:"deviceId"`
	DeviceName                string                    `json:"deviceName"`
	CompliancePolicies        []DeviceCompliancePolicy  `json:"compliancePolicies"`
	ConfigurationProfiles     []DeviceConfiguration     `json:"configurationProfiles"`
	SecurityBaselines         []SecurityBaseline        `json:"securityBaselines"`
	BitLockerStatus           BitLockerStatus           `json:"bitLockerStatus"`
	WindowsDefenderStatus     WindowsDefenderStatus     `json:"windowsDefenderStatus"`
	FirewallStatus            FirewallStatus            `json:"firewallStatus"`
	AppProtectionPolicies     []AppProtectionPolicy     `json:"appProtectionPolicies"`
	ConditionalAccessPolicies []ConditionalAccessPolicy `json:"conditionalAccessPolicies"`
	CollectedAt               time.Time                 `json:"collectedAt"`
}

// SecurityBaseline represents a security baseline configuration
type SecurityBaseline struct {
	ID                   string            `json:"id"`
	DisplayName          string            `json:"displayName"`
	Description          string            `json:"description"`
	SecurityBaselineType string            `json:"securityBaselineType"`
	CreatedDateTime      time.Time         `json:"createdDateTime"`
	LastModifiedDateTime time.Time         `json:"lastModifiedDateTime"`
	Settings             []SecuritySetting `json:"settings"`
}

// SecuritySetting represents an individual security setting
type SecuritySetting struct {
	ID               string      `json:"id"`
	DisplayName      string      `json:"displayName"`
	SettingType      string      `json:"settingType"`
	CurrentValue     interface{} `json:"currentValue"`
	RecommendedValue interface{} `json:"recommendedValue"`
	ComplianceState  string      `json:"complianceState"`
	Severity         string      `json:"severity"`
}

// BitLockerStatus represents BitLocker encryption status
type BitLockerStatus struct {
	EncryptionMethod        string    `json:"encryptionMethod"`
	EncryptionStatus        string    `json:"encryptionStatus"`
	ProtectionStatus        string    `json:"protectionStatus"`
	KeyProtectors           []string  `json:"keyProtectors"`
	RecoveryKeyBackupStatus string    `json:"recoveryKeyBackupStatus"`
	LastStatusUpdate        time.Time `json:"lastStatusUpdate"`
}

// WindowsDefenderStatus represents Windows Defender status
type WindowsDefenderStatus struct {
	AntivirusEnabled             bool      `json:"antivirusEnabled"`
	AntivirusSignatureVersion    string    `json:"antivirusSignatureVersion"`
	AntivirusSignatureLastUpdate time.Time `json:"antivirusSignatureLastUpdate"`
	RealTimeProtectionEnabled    bool      `json:"realTimeProtectionEnabled"`
	BehaviorMonitorEnabled       bool      `json:"behaviorMonitorEnabled"`
	FirewallEnabled              bool      `json:"firewallEnabled"`
	SmartScreenEnabled           bool      `json:"smartScreenEnabled"`
	CloudProtectionEnabled       bool      `json:"cloudProtectionEnabled"`
	TamperProtectionEnabled      bool      `json:"tamperProtectionEnabled"`
}

// FirewallStatus represents Windows Firewall status
type FirewallStatus struct {
	DomainProfile    FirewallProfile `json:"domainProfile"`
	PrivateProfile   FirewallProfile `json:"privateProfile"`
	PublicProfile    FirewallProfile `json:"publicProfile"`
	LastStatusUpdate time.Time       `json:"lastStatusUpdate"`
}

// FirewallProfile represents a specific firewall profile configuration
type FirewallProfile struct {
	Enabled               bool     `json:"enabled"`
	DefaultInboundAction  string   `json:"defaultInboundAction"`
	DefaultOutboundAction string   `json:"defaultOutboundAction"`
	NotificationsEnabled  bool     `json:"notificationsEnabled"`
	StealthModeEnabled    bool     `json:"stealthModeEnabled"`
	ExceptionRules        []string `json:"exceptionRules"`
}

// AppProtectionPolicy represents an app protection policy
type AppProtectionPolicy struct {
	ID                   string                 `json:"id"`
	DisplayName          string                 `json:"displayName"`
	Description          string                 `json:"description"`
	PlatformType         string                 `json:"platformType"`
	CreatedDateTime      time.Time              `json:"createdDateTime"`
	LastModifiedDateTime time.Time              `json:"lastModifiedDateTime"`
	Settings             map[string]interface{} `json:"settings"`
	AssignedGroups       []string               `json:"assignedGroups"`
}

// ConditionalAccessPolicy represents a conditional access policy
type ConditionalAccessPolicy struct {
	ID               string                 `json:"id"`
	DisplayName      string                 `json:"displayName"`
	State            string                 `json:"state"`
	CreatedDateTime  time.Time              `json:"createdDateTime"`
	ModifiedDateTime time.Time              `json:"modifiedDateTime"`
	Conditions       map[string]interface{} `json:"conditions"`
	GrantControls    map[string]interface{} `json:"grantControls"`
	SessionControls  map[string]interface{} `json:"sessionControls"`
}

// IntuneComplianceReport represents a comprehensive compliance report
type IntuneComplianceReport struct {
	TenantID            string                   `json:"tenantId"`
	ReportTimestamp     time.Time                `json:"reportTimestamp"`
	TotalDevices        int                      `json:"totalDevices"`
	CompliantDevices    int                      `json:"compliantDevices"`
	NonCompliantDevices int                      `json:"nonCompliantDevices"`
	DeviceBreakdown     DeviceBreakdown          `json:"deviceBreakdown"`
	SecurityFindings    []SecurityFinding        `json:"securityFindings"`
	TopRisks            []RiskSummary            `json:"topRisks"`
	ComplianceTrends    []ComplianceTrend        `json:"complianceTrends"`
	Recommendations     []SecurityRecommendation `json:"recommendations"`
}

// DeviceBreakdown provides statistics about device types and platforms
type DeviceBreakdown struct {
	Windows      int `json:"windows"`
	MacOS        int `json:"macOS"`
	iOS          int `json:"iOS"`
	Android      int `json:"android"`
	WindowsPhone int `json:"windowsPhone"`
	Other        int `json:"other"`
}

// RiskSummary represents a high-level risk category summary
type RiskSummary struct {
	RiskCategory    string  `json:"riskCategory"`
	AffectedDevices int     `json:"affectedDevices"`
	Severity        string  `json:"severity"`
	Description     string  `json:"description"`
	ImpactScore     float64 `json:"impactScore"`
}

// ComplianceTrend represents compliance status over time
type ComplianceTrend struct {
	Date              time.Time `json:"date"`
	CompliantCount    int       `json:"compliantCount"`
	NonCompliantCount int       `json:"nonCompliantCount"`
	TotalCount        int       `json:"totalCount"`
	ComplianceRate    float64   `json:"complianceRate"`
}

// SecurityRecommendation represents an actionable security recommendation
type SecurityRecommendation struct {
	ID               string   `json:"id"`
	Title            string   `json:"title"`
	Priority         string   `json:"priority"`
	Category         string   `json:"category"`
	Description      string   `json:"description"`
	Impact           string   `json:"impact"`
	Implementation   string   `json:"implementation"`
	AffectedDevices  int      `json:"affectedDevices"`
	EstimatedEffort  string   `json:"estimatedEffort"`
	MITREMitigations []string `json:"mitreMitigations"`
}

// BloodHoundIntuneData represents data formatted specifically for BloodHound ingestion
type BloodHoundIntuneData struct {
	Meta               BloodHoundMeta        `json:"meta"`
	Data               BloodHoundDataWrapper `json:"data"`
	ComputerDomains    []ComputerDomain      `json:"computerDomains"`
	Computers          []Computer            `json:"computers"`
	Users              []BloodHoundUser      `json:"users"`
	Groups             []BloodHoundGroup     `json:"groups"`
	LocalAdmins        []LocalAdmin          `json:"localAdmins"`
	RemoteDesktopUsers []RemoteDesktopUser   `json:"remoteDesktopUsers"`
	DcomUsers          []DcomUser            `json:"dcomUsers"`
	PSRemoteUsers      []PSRemoteUser        `json:"psRemoteUsers"`
	Sessions           []Session             `json:"sessions"`
	RegistryKeys       []RegistryKey         `json:"registryKeys"`
}

// BloodHoundMeta contains metadata about the collection
type BloodHoundMeta struct {
	Type        string    `json:"type"`
	Count       int       `json:"count"`
	Version     string    `json:"version"`
	Methods     int       `json:"methods"`
	CollectedBy string    `json:"collectedBy"`
	CollectedAt time.Time `json:"collectedAt"`
}

// BloodHoundDataWrapper wraps the data arrays for BloodHound
type BloodHoundDataWrapper struct {
	Computers          []Computer          `json:"computers"`
	Users              []BloodHoundUser    `json:"users"`
	Groups             []BloodHoundGroup   `json:"groups"`
	LocalAdmins        []LocalAdmin        `json:"localAdmins"`
	RemoteDesktopUsers []RemoteDesktopUser `json:"remoteDesktopUsers"`
	Sessions           []Session           `json:"sessions"`
}

// Computer represents a computer object for BloodHound
type Computer struct {
	ObjectIdentifier   string               `json:"ObjectIdentifier"`
	PrimaryGroupSID    string               `json:"PrimaryGroupSID"`
	LocalAdmins        []LocalAdminRelation `json:"LocalAdmins"`
	RemoteDesktopUsers []RDPUsersRelation   `json:"RemoteDesktopUsers"`
	DcomUsers          []DcomUsersRelation  `json:"DcomUsers"`
	PSRemoteUsers      []PSRemoteRelation   `json:"PSRemoteUsers"`
	Properties         ComputerProperties   `json:"Properties"`
	Aces               []ACE                `json:"Aces"`
	Sessions           []SessionRelation    `json:"Sessions"`
	RegistryFindings   []RegistryFinding    `json:"RegistryFindings"`
	SecurityFindings   []SecurityFinding    `json:"SecurityFindings"`
}

// ComputerProperties represents properties of a computer
type ComputerProperties struct {
	Name                    string    `json:"name"`
	Domain                  string    `json:"domain"`
	ObjectID                string    `json:"objectid"`
	PrimaryGroupSID         string    `json:"primarygroupsid"`
	HasLAPS                 bool      `json:"haslaps"`
	LastLogon               int64     `json:"lastlogon"`
	LastLogonTimestamp      int64     `json:"lastlogontimestamp"`
	PwdLastSet              int64     `json:"pwdlastset"`
	ServicePrincipalNames   []string  `json:"serviceprincipalnames"`
	Description             string    `json:"description"`
	OperatingSystem         string    `json:"operatingsystem"`
	Enabled                 bool      `json:"enabled"`
	UnconstrainedDelegation bool      `json:"unconstraineddelegation"`
	TrustedToAuth           bool      `json:"trustedtoauth"`
	SamAccountName          string    `json:"samaccountname"`
	DistinguishedName       string    `json:"distinguishedname"`
	IntuneDeviceID          string    `json:"intunedeviceid"`
	ComplianceState         string    `json:"compliancestate"`
	LastSyncDateTime        time.Time `json:"lastsyncdatetime"`
	RiskScore               int       `json:"riskscore"`
}

// BloodHoundUser represents a user object for BloodHound (renamed to avoid conflict)
type BloodHoundUser struct {
	ObjectIdentifier string                   `json:"ObjectIdentifier"`
	PrimaryGroupSID  string                   `json:"PrimaryGroupSID"`
	Properties       BloodHoundUserProperties `json:"Properties"`
	Aces             []ACE                    `json:"Aces"`
	Sessions         []SessionRelation        `json:"Sessions"`
}

// BloodHoundUserProperties represents properties of a user (renamed to avoid conflict)
type BloodHoundUserProperties struct {
	Name                     string   `json:"name"`
	Domain                   string   `json:"domain"`
	ObjectID                 string   `json:"objectid"`
	PrimaryGroupSID          string   `json:"primarygroupsid"`
	HasSPN                   bool     `json:"hasspn"`
	ServicePrincipalNames    []string `json:"serviceprincipalnames"`
	DisplayName              string   `json:"displayname"`
	Email                    string   `json:"email"`
	Title                    string   `json:"title"`
	Department               string   `json:"department"`
	LastLogon                int64    `json:"lastlogon"`
	LastLogonTimestamp       int64    `json:"lastlogontimestamp"`
	PwdLastSet               int64    `json:"pwdlastset"`
	Enabled                  bool     `json:"enabled"`
	PasswordNeverExpires     bool     `json:"passwordneverexpires"`
	PasswordNotRequired      bool     `json:"passwordnotrequired"`
	UserCannotChangePassword bool     `json:"usercannotchangepassword"`
	DontRequirePreAuth       bool     `json:"dontreqpreauth"`
	SamAccountName           string   `json:"samaccountname"`
	DistinguishedName        string   `json:"distinguishedname"`
	UnconstrainedDelegation  bool     `json:"unconstraineddelegation"`
	Sensitive                bool     `json:"sensitive"`
	AllowedToDelegate        []string `json:"allowedtodelegate"`
	AdminCount               bool     `json:"admincount"`
	SIDHistory               []string `json:"sidhistory"`
}

// BloodHoundGroup represents a group object for BloodHound (renamed to avoid conflict)
type BloodHoundGroup struct {
	ObjectIdentifier string                    `json:"ObjectIdentifier"`
	Properties       BloodHoundGroupProperties `json:"Properties"`
	Aces             []ACE                     `json:"Aces"`
	Members          []Member                  `json:"Members"`
}

// BloodHoundGroupProperties represents properties of a group (renamed to avoid conflict)
type BloodHoundGroupProperties struct {
	Name              string `json:"name"`
	Domain            string `json:"domain"`
	ObjectID          string `json:"objectid"`
	Description       string `json:"description"`
	AdminCount        bool   `json:"admincount"`
	SamAccountName    string `json:"samaccountname"`
	DistinguishedName string `json:"distinguishedname"`
}

// LocalAdmin represents a local administrator relationship
type LocalAdmin struct {
	ObjectIdentifier string `json:"ObjectIdentifier"`
	ObjectType       string `json:"ObjectType"`
	ComputerSID      string `json:"ComputerSID"`
}

// LocalAdminRelation represents a local admin relationship for BloodHound
type LocalAdminRelation struct {
	ObjectIdentifier string `json:"ObjectIdentifier"`
	ObjectType       string `json:"ObjectType"`
}

// RemoteDesktopUser represents a remote desktop user relationship
type RemoteDesktopUser struct {
	ObjectIdentifier string `json:"ObjectIdentifier"`
	ObjectType       string `json:"ObjectType"`
	ComputerSID      string `json:"ComputerSID"`
}

// RDPUsersRelation represents an RDP user relationship for BloodHound
type RDPUsersRelation struct {
	ObjectIdentifier string `json:"ObjectIdentifier"`
	ObjectType       string `json:"ObjectType"`
}

// DcomUser represents a DCOM user relationship
type DcomUser struct {
	ObjectIdentifier string `json:"ObjectIdentifier"`
	ObjectType       string `json:"ObjectType"`
	ComputerSID      string `json:"ComputerSID"`
}

// DcomUsersRelation represents a DCOM user relationship for BloodHound
type DcomUsersRelation struct {
	ObjectIdentifier string `json:"ObjectIdentifier"`
	ObjectType       string `json:"ObjectType"`
}

// PSRemoteUser represents a PowerShell remoting user relationship
type PSRemoteUser struct {
	ObjectIdentifier string `json:"ObjectIdentifier"`
	ObjectType       string `json:"ObjectType"`
	ComputerSID      string `json:"ComputerSID"`
}

// PSRemoteRelation represents a PS Remote relationship for BloodHound
type PSRemoteRelation struct {
	ObjectIdentifier string `json:"ObjectIdentifier"`
	ObjectType       string `json:"ObjectType"`
}

// Session represents a user session
type Session struct {
	ComputerSID string `json:"ComputerSID"`
	UserSID     string `json:"UserSID"`
	LogonType   string `json:"LogonType"`
}

// SessionRelation represents a session relationship for BloodHound
type SessionRelation struct {
	UserSID   string `json:"UserSID"`
	LogonType string `json:"LogonType"`
}

// ComputerDomain represents a computer's domain relationship
type ComputerDomain struct {
	ComputerSID string `json:"ComputerSID"`
	DomainSID   string `json:"DomainSID"`
}

// Member represents a group membership
type Member struct {
	ObjectIdentifier string `json:"ObjectIdentifier"`
	ObjectType       string `json:"ObjectType"`
}

// ACE represents an Access Control Entry
type ACE struct {
	PrincipalSID  string `json:"PrincipalSID"`
	PrincipalType string `json:"PrincipalType"`
	RightName     string `json:"RightName"`
	AceType       string `json:"AceType"`
	IsInherited   bool   `json:"IsInherited"`
}

// RegistryKey represents a registry key finding for BloodHound
type RegistryKey struct {
	ComputerSID  string                 `json:"ComputerSID"`
	RegistryPath string                 `json:"RegistryPath"`
	ValueName    string                 `json:"ValueName"`
	ValueData    interface{}            `json:"ValueData"`
	ValueType    string                 `json:"ValueType"`
	SecurityRisk string                 `json:"SecurityRisk"`
	AttackVector string                 `json:"AttackVector"`
	Properties   map[string]interface{} `json:"Properties"`
}
