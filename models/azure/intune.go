// models/azure/intune.go
package azure

import (
	"time"
)

// IntuneDevice represents a device managed by Microsoft Intune
type IntuneDevice struct {
	ID                                      string        `json:"id"`
	DeviceName                              string        `json:"deviceName"`
	OperatingSystem                         string        `json:"operatingSystem"`
	OSVersion                               string        `json:"osVersion"`
	ComplianceState                         string        `json:"complianceState"`
	LastSyncDateTime                        time.Time     `json:"lastSyncDateTime"`
	EnrollmentType                          string        `json:"enrollmentType"`
	ManagementAgent                         string        `json:"managementAgent"`
	AzureADDeviceID                         string        `json:"azureADDeviceId"`
	UserPrincipalName                       string        `json:"userPrincipalName"`
	SerialNumber                            string        `json:"serialNumber"`
	Manufacturer                            string        `json:"manufacturer"`
	Model                                   string        `json:"model"`
	TotalStorageSpaceInBytes                int64         `json:"totalStorageSpaceInBytes"`
	FreeStorageSpaceInBytes                 int64         `json:"freeStorageSpaceInBytes"`
	ManagedDeviceName                       string        `json:"managedDeviceName"`
	PartnerReportedThreatState              string        `json:"partnerReportedThreatState"`
	RequireUserEnrollmentApproval           bool          `json:"requireUserEnrollmentApproval"`
	ManagementCertificateExpirationDate     time.Time     `json:"managementCertificateExpirationDate"`
	ICCID                                   string        `json:"iccid"`
	UDID                                    string        `json:"udid"`
	Notes                                   string        `json:"notes"`
	EthernetMacAddress                      string        `json:"ethernetMacAddress"`
	WiFiMacAddress                          string        `json:"wiFiMacAddress"`
	PhysicalMemoryInBytes                   int64         `json:"physicalMemoryInBytes"`
	ProcessorArchitecture                   string        `json:"processorArchitecture"`
	SpecificationVersion                    string        `json:"specificationVersion"`
	JoinType                                string        `json:"joinType"`
	SkuFamily                               string        `json:"skuFamily"`
	SkuNumber                               int           `json:"skuNumber"`
	ManagementFeatures                      string        `json:"managementFeatures"`
	ChromeOSDeviceInfo                      []interface{} `json:"chromeOSDeviceInfo"`
	EnrolledDateTime                        time.Time     `json:"enrolledDateTime"`
	EmailAddress                            string        `json:"emailAddress"`
	UserID                                  string        `json:"userId"`
	UserDisplayName                         string        `json:"userDisplayName"`
	DeviceRegistrationState                 string        `json:"deviceRegistrationState"`
	DeviceCategoryDisplayName               string        `json:"deviceCategoryDisplayName"`
	IsSupervised                            bool          `json:"isSupervised"`
	ExchangeLastSuccessfulSyncDateTime      time.Time     `json:"exchangeLastSuccessfulSyncDateTime"`
	ExchangeAccessState                     string        `json:"exchangeAccessState"`
	ExchangeAccessStateReason               string        `json:"exchangeAccessStateReason"`
	RemoteAssistanceSessionURL              string        `json:"remoteAssistanceSessionUrl"`
	RemoteAssistanceSessionErrorDetails     string        `json:"remoteAssistanceSessionErrorDetails"`
	IsEncrypted                             bool          `json:"isEncrypted"`
	ComplianceGracePeriodExpirationDateTime time.Time     `json:"complianceGracePeriodExpirationDateTime"`
	ManagementAgents                        []string      `json:"managementAgents"`
	LostModeState                           string        `json:"lostModeState"`
	ActivationLockBypassCode                string        `json:"activationLockBypassCode"`
}

// ScriptExecution represents the execution of a PowerShell script on an Intune device
type ScriptExecution struct {
	ID            string     `json:"id"`
	DeviceID      string     `json:"deviceId"`
	Status        string     `json:"status"`
	StartDateTime time.Time  `json:"startDateTime"`
	EndDateTime   *time.Time `json:"endDateTime"`
	ScriptName    string     `json:"scriptName"`
	RunAsAccount  string     `json:"runAsAccount"`
}

// ScriptExecutionResult represents the result of a PowerShell script execution
type ScriptExecutionResult struct {
	ID                                   string    `json:"id"`
	DeviceID                             string    `json:"deviceId"`
	DeviceName                           string    `json:"deviceName"`
	RunState                             string    `json:"runState"`
	ResultMessage                        string    `json:"resultMessage"`
	PreRemediationDetectionScriptOutput  string    `json:"preRemediationDetectionScriptOutput"`
	RemediationScriptOutput              string    `json:"remediationScriptOutput"`
	PostRemediationDetectionScriptOutput string    `json:"postRemediationDetectionScriptOutput"`
	ErrorCode                            int       `json:"errorCode"`
	LastStateUpdateDateTime              time.Time `json:"lastStateUpdateDateTime"`
}

// DeviceCompliancePolicy represents a device compliance policy
type DeviceCompliancePolicy struct {
	ID                   string    `json:"id"`
	DisplayName          string    `json:"displayName"`
	Description          string    `json:"description"`
	Platform             string    `json:"platform"`
	CreatedDateTime      time.Time `json:"createdDateTime"`
	LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
}

// DeviceConfiguration represents a device configuration profile
type DeviceConfiguration struct {
	ID                   string                 `json:"id"`
	DisplayName          string                 `json:"displayName"`
	Description          string                 `json:"description"`
	Platform             string                 `json:"platform"`
	CreatedDateTime      time.Time              `json:"createdDateTime"`
	LastModifiedDateTime time.Time              `json:"lastModifiedDateTime"`
	Settings             map[string]interface{} `json:"settings"`
}

// DeviceRegistryData combines device information with collected registry data
type DeviceRegistryData struct {
	Device       IntuneDevice `json:"device"`
	RegistryData RegistryData `json:"registryData"`
	CollectedAt  time.Time    `json:"collectedAt"`

	// BloodHound specific fields for integration
	BloodHoundData BloodHoundDeviceData `json:"bloodhoundData"`
}

// RegistryData represents the complete registry data collected from a device
type RegistryData struct {
	DeviceInfo         DeviceInfo          `json:"deviceInfo"`
	RegistryData       []RegistryEntry     `json:"registryData"`
	SecurityIndicators SecurityIndicators  `json:"securityIndicators"`
	Summary            RegistryDataSummary `json:"summary"`
}

// DeviceInfo contains basic information about the device where data was collected
type DeviceInfo struct {
	ComputerName  string `json:"computerName"`
	Domain        string `json:"domain"`
	User          string `json:"user"`
	Timestamp     string `json:"timestamp"`
	ScriptVersion string `json:"scriptVersion"`
}

// RegistryEntry represents a single registry path and its collected values
type RegistryEntry struct {
	Path       string                 `json:"path"`
	Purpose    string                 `json:"purpose"`
	Values     map[string]interface{} `json:"values"`
	Accessible bool                   `json:"accessible"`
	Error      *string                `json:"error"`
}

// SecurityIndicators contains analysis results of security-relevant registry settings
type SecurityIndicators struct {
	UACDisabled            bool             `json:"uacDisabled"`
	AutoAdminLogon         bool             `json:"autoAdminLogon"`
	WeakServicePermissions bool             `json:"weakServicePermissions"`
	SuspiciousStartupItems []SuspiciousItem `json:"suspiciousStartupItems"`
}

// SuspiciousItem represents a potentially malicious startup item
type SuspiciousItem struct {
	Path  string `json:"path"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

// RegistryDataSummary provides high-level statistics about the collected data
type RegistryDataSummary struct {
	TotalKeysChecked   int      `json:"totalKeysChecked"`
	AccessibleKeys     int      `json:"accessibleKeys"`
	HighRiskIndicators []string `json:"highRiskIndicators"`
}

// BloodHoundDeviceData contains processed data formatted for BloodHound consumption
type BloodHoundDeviceData struct {
	ObjectIdentifier    string              `json:"ObjectIdentifier"`
	AzureDeviceID       string              `json:"AzureDeviceID"`
	DisplayName         string              `json:"DisplayName"`
	LocalGroups         map[string][]string `json:"LocalGroups"`
	UserRights          map[string][]string `json:"UserRights"`
	Sessions            []SessionInfo       `json:"Sessions"`
	RegistryFindings    []RegistryFinding   `json:"RegistryFindings"`
	SecurityFindings    []SecurityFinding   `json:"SecurityFindings"`
	PrivilegeEscalation []EscalationVector  `json:"PrivilegeEscalation"`
}

// SessionInfo represents active user sessions on the device
type SessionInfo struct {
	UserName    string    `json:"UserName"`
	SessionType string    `json:"SessionType"`
	SessionID   int       `json:"SessionID"`
	State       string    `json:"State"`
	IdleTime    string    `json:"IdleTime"`
	LogonTime   time.Time `json:"LogonTime"`
}

// RegistryFinding represents a specific registry-based security finding
type RegistryFinding struct {
	Category      string      `json:"Category"`
	Finding       string      `json:"Finding"`
	Severity      string      `json:"Severity"`
	RegistryPath  string      `json:"RegistryPath"`
	ValueName     string      `json:"ValueName"`
	CurrentValue  interface{} `json:"CurrentValue"`
	ExpectedValue interface{} `json:"ExpectedValue"`
	Description   string      `json:"Description"`
	Remediation   string      `json:"Remediation"`
	AttackVector  string      `json:"AttackVector"`
}

// SecurityFinding represents high-level security issues identified
type SecurityFinding struct {
	ID              string   `json:"ID"`
	Title           string   `json:"Title"`
	Severity        string   `json:"Severity"`
	Category        string   `json:"Category"`
	Description     string   `json:"Description"`
	Evidence        []string `json:"Evidence"`
	Recommendations []string `json:"Recommendations"`
	MITREAttack     []string `json:"MITREAttack"`
}

// EscalationVector represents a potential privilege escalation path
type EscalationVector struct {
	VectorID      string   `json:"VectorID"`
	Type          string   `json:"Type"`
	Source        string   `json:"Source"`
	Target        string   `json:"Target"`
	Method        string   `json:"Method"`
	RequiredPrivs []string `json:"RequiredPrivs"`
	Complexity    string   `json:"Complexity"`
	Impact        string   `json:"Impact"`
	Conditions    []string `json:"Conditions"`
}

// IntuneAppRegistration represents the Azure app registration for Intune access
type IntuneAppRegistration struct {
	ClientID     string   `json:"clientId"`
	TenantID     string   `json:"tenantId"`
	ClientSecret string   `json:"clientSecret"`
	Permissions  []string `json:"permissions"`
}

// IntuneManagementScript represents a PowerShell script configured in Intune
type IntuneManagementScript struct {
	ID                   string    `json:"id"`
	DisplayName          string    `json:"displayName"`
	Description          string    `json:"description"`
	ScriptContent        string    `json:"scriptContent"`
	CreatedDateTime      time.Time `json:"createdDateTime"`
	LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
	RunAsAccount         string    `json:"runAsAccount"`
	FileName             string    `json:"fileName"`
	RoleScopeTagIds      []string  `json:"roleScopeTagIds"`
}
