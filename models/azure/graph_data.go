// models/azure/graph_data.go
package azure

import (
	"encoding/json"
	"time"
)

// GroupMembershipData represents Azure AD group with its members and owners
type GroupMembershipData struct {
	Group   Group             `json:"group"`
	Members []json.RawMessage `json:"members"`
	Owners  []json.RawMessage `json:"owners"`
}

// UserRoleData represents a user with their role assignments
type UserRoleData struct {
	User            User                `json:"user"`
	RoleAssignments []AppRoleAssignment `json:"roleAssignments"`
}

// DeviceAccessData represents device access permissions and ownership
type DeviceAccessData struct {
	IntuneDevice     IntuneDevice      `json:"intuneDevice"`
	AzureDevice      *Device           `json:"azureDevice,omitempty"`
	RegisteredUsers  []json.RawMessage `json:"registeredUsers"`
	RegisteredOwners []json.RawMessage `json:"registeredOwners"`
}

// SignIn represents sign-in activity data
type SignIn struct {
	ID                      string                 `json:"id"`
	CreatedDateTime         time.Time              `json:"createdDateTime"`
	UserDisplayName         string                 `json:"userDisplayName"`
	UserPrincipalName       string                 `json:"userPrincipalName"`
	UserId                  string                 `json:"userId"`
	AppId                   string                 `json:"appId"`
	AppDisplayName          string                 `json:"appDisplayName"`
	IpAddress               string                 `json:"ipAddress"`
	ClientAppUsed           string                 `json:"clientAppUsed"`
	DeviceDetail            DeviceDetail           `json:"deviceDetail"`
	Location                SignInLocation         `json:"location"`
	RiskDetail              string                 `json:"riskDetail"`
	RiskLevelAggregated     string                 `json:"riskLevelAggregated"`
	RiskLevelDuringSignIn   string                 `json:"riskLevelDuringSignIn"`
	RiskState               string                 `json:"riskState"`
	Status                  SignInStatus           `json:"status"`
	ConditionalAccessStatus string                 `json:"conditionalAccessStatus"`
	AdditionalData          map[string]interface{} `json:"additionalData,omitempty"`
}

// DeviceDetail represents device information from sign-in
type DeviceDetail struct {
	DeviceId        string `json:"deviceId"`
	DisplayName     string `json:"displayName"`
	OperatingSystem string `json:"operatingSystem"`
	Browser         string `json:"browser"`
	IsCompliant     bool   `json:"isCompliant"`
	IsManaged       bool   `json:"isManaged"`
	TrustType       string `json:"trustType"`
}

// SignInLocation represents sign-in location
type SignInLocation struct {
	City            string         `json:"city"`
	State           string         `json:"state"`
	CountryOrRegion string         `json:"countryOrRegion"`
	GeoCoordinates  GeoCoordinates `json:"geoCoordinates"`
}

// GeoCoordinates represents geographic coordinates
type GeoCoordinates struct {
	Altitude  float64 `json:"altitude"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// SignInStatus represents sign-in status
type SignInStatus struct {
	ErrorCode         int    `json:"errorCode"`
	FailureReason     string `json:"failureReason"`
	AdditionalDetails string `json:"additionalDetails"`
}

// BloodHoundGraphData represents data collected via Graph API formatted for BloodHound
type BloodHoundGraphData struct {
	Meta                BloodHoundMeta                 `json:"meta"`
	Data                BloodHoundGraphDataWrapper     `json:"data"`
	GroupMemberships    []BloodHoundGroupMembership    `json:"groupMemberships"`
	UserRoleAssignments []BloodHoundUserRoleAssignment `json:"userRoleAssignments"`
	DeviceOwnerships    []BloodHoundDeviceOwnership    `json:"deviceOwnerships"`
	SignInActivity      []BloodHoundSignInActivity     `json:"signInActivity"`
}

// BloodHoundGraphDataWrapper wraps the core data
type BloodHoundGraphDataWrapper struct {
	Users   []BloodHoundUser   `json:"users"`
	Groups  []BloodHoundGroup  `json:"groups"`
	Devices []BloodHoundDevice `json:"devices"`
}

// BloodHoundGroupMembership represents group membership for BloodHound
type BloodHoundGroupMembership struct {
	GroupId          string `json:"groupId"`
	GroupName        string `json:"groupName"`
	MemberId         string `json:"memberId"`
	MemberName       string `json:"memberName"`
	MemberType       string `json:"memberType"`
	RelationshipType string `json:"relationshipType"` // "MemberOf" or "OwnerOf"
}

// BloodHoundUserRoleAssignment represents user role assignments for BloodHound
type BloodHoundUserRoleAssignment struct {
	UserId          string    `json:"userId"`
	UserName        string    `json:"userName"`
	RoleId          string    `json:"roleId"`
	RoleName        string    `json:"roleName"`
	ResourceId      string    `json:"resourceId"`
	ResourceName    string    `json:"resourceName"`
	AssignmentType  string    `json:"assignmentType"`
	CreatedDateTime time.Time `json:"createdDateTime"`
}

// BloodHoundDeviceOwnership represents device ownership for BloodHound
type BloodHoundDeviceOwnership struct {
	DeviceId        string `json:"deviceId"`
	DeviceName      string `json:"deviceName"`
	UserId          string `json:"userId"`
	UserName        string `json:"userName"`
	OwnershipType   string `json:"ownershipType"` // "RegisteredOwner" or "RegisteredUser"
	ComplianceState string `json:"complianceState"`
}

// BloodHoundSignInActivity represents sign-in activity for BloodHound
type BloodHoundSignInActivity struct {
	UserId            string    `json:"userId"`
	UserName          string    `json:"userName"`
	DeviceId          string    `json:"deviceId"`
	DeviceName        string    `json:"deviceName"`
	AppId             string    `json:"appId"`
	AppName           string    `json:"appName"`
	SignInDateTime    time.Time `json:"signInDateTime"`
	IpAddress         string    `json:"ipAddress"`
	Location          string    `json:"location"`
	RiskLevel         string    `json:"riskLevel"`
	ConditionalAccess string    `json:"conditionalAccess"`
}

// BloodHoundDevice represents a device for BloodHound
type BloodHoundDevice struct {
	ObjectIdentifier string                     `json:"ObjectIdentifier"`
	Properties       BloodHoundDeviceProperties `json:"Properties"`
	RegisteredUsers  []BloodHoundDeviceUser     `json:"RegisteredUsers"`
	RegisteredOwners []BloodHoundDeviceUser     `json:"RegisteredOwners"`
}

// BloodHoundDeviceProperties represents device properties for BloodHound
type BloodHoundDeviceProperties struct {
	Name             string    `json:"name"`
	DisplayName      string    `json:"displayName"`
	ObjectID         string    `json:"objectid"`
	OperatingSystem  string    `json:"operatingsystem"`
	OSVersion        string    `json:"osversion"`
	DeviceId         string    `json:"deviceid"`
	IsCompliant      bool      `json:"iscompliant"`
	IsManaged        bool      `json:"ismanaged"`
	EnrollmentType   string    `json:"enrollmenttype"`
	JoinType         string    `json:"jointype"`
	TrustType        string    `json:"trusttype"`
	LastSyncDateTime time.Time `json:"lastsyncdatetime"`
	CreatedDateTime  time.Time `json:"createddatetime"`
	Enabled          bool      `json:"enabled"`
}

// BloodHoundDeviceUser represents a user associated with a device
type BloodHoundDeviceUser struct {
	ObjectIdentifier string `json:"ObjectIdentifier"`
	ObjectType       string `json:"ObjectType"`
}

// Collection results for Graph API data
type GraphDataCollectionResult struct {
	GroupMemberships    []GroupMembershipData `json:"groupMemberships"`
	UserRoleAssignments []UserRoleData        `json:"userRoleAssignments"`
	DeviceAccess        []DeviceAccessData    `json:"deviceAccess"`
	SignInActivity      []SignIn              `json:"signInActivity"`
	CollectionTime      time.Duration         `json:"collectionTime"`
	TotalGroups         int                   `json:"totalGroups"`
	TotalUsers          int                   `json:"totalUsers"`
	TotalDevices        int                   `json:"totalDevices"`
	TotalSignIns        int                   `json:"totalSignIns"`
	Errors              []string              `json:"errors"`
}
