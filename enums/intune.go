package enums

// Intune-specific Kind enumerations for data types
const (
    KindAZIntuneDevice         Kind = "AZIntuneDevice"
    KindAZIntuneDeviceCompliance Kind = "AZIntuneDeviceCompliance"
    KindAZIntuneDeviceConfiguration Kind = "AZIntuneDeviceConfiguration"
    KindAZIntuneCompliance     Kind = "AZIntuneCompliance"
)

// Device compliance states
type ComplianceState string

const (
    ComplianceStateCompliant    ComplianceState = "compliant"
    ComplianceStateNoncompliant ComplianceState = "noncompliant"
    ComplianceStateConflict     ComplianceState = "conflict"
    ComplianceStateError        ComplianceState = "error"
    ComplianceStateUnknown      ComplianceState = "unknown"
    ComplianceStateInGracePeriod ComplianceState = "inGracePeriod"
)

// Device enrollment types
type EnrollmentType string

const (
    EnrollmentTypeUserEnrollment            EnrollmentType = "userEnrollment"
    EnrollmentTypeDeviceEnrollmentManager   EnrollmentType = "deviceEnrollmentManager"
    EnrollmentTypeWindowsAzureADJoin        EnrollmentType = "windowsAzureADJoin"
    EnrollmentTypeWindowsAutoEnrollment     EnrollmentType = "windowsAutoEnrollment"
    EnrollmentTypeWindowsCoManagement        EnrollmentType = "windowsCoManagement"
)

// Management agent types
type ManagementAgent string

const (
    ManagementAgentMDM               ManagementAgent = "mdm"
    ManagementAgentIntuneClient      ManagementAgent = "intuneClient"
    ManagementAgentConfigurationManagerClient ManagementAgent = "configurationManagerClient"
    ManagementAgentUnknown           ManagementAgent = "unknown"
)

// Operating system types
type OperatingSystem string

const (
    OperatingSystemWindows     OperatingSystem = "windows"
    OperatingSystemAndroid     OperatingSystem = "android"
    OperatingSystemIOS         OperatingSystem = "iOS"
    OperatingSystemMacOS       OperatingSystem = "macOS"
)