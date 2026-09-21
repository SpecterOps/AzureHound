// Copyright (C) 2022 Specter Ops, Inc.
//
// This file is part of AzureHound.
//
// AzureHound is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// AzureHound is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package azure

import (
	"github.com/bloodhoundad/azurehound/v2/enums"
)

// Represents an Azure Active Directory (Azure AD) group, which can be a Microsoft 365 group, or a security group.
// For more detail see https://docs.microsoft.com/en-us/graph/api/resources/group?view=graph-rest-1.0
type Group365 struct {
	DirectoryObject

	// Indicates if people external to the organization can send messages to the group.
	// Default value is false.
	// Returned only on $select for GET /groups/{ID}
	AllowExternalSenders bool `json:"allowExternalSenders,omitempty"`

	// The list of sensitivity label pairs (label ID, label name) associated with a Microsoft 365 group.
	// Returned only on $select.
	// Read-only.
	AssignedLabels []AssignedLabel `json:"assignedLabels,omitempty"`

	// The licenses that are assigned to the group.
	// Returned only on $select.
	// Supports $filter (eq)
	// Read-only.
	AssignedLicenses []AssignedLicense `json:"assignedLicenses,omitempty"`

	// Indicates if new members added to the group will be auto-subscribed to receive email notifications.
	// You can set this property in a PATCH request for the group; do not set it in the initial POST request that
	// creates the group.
	// Default value is false.
	// Returned only on $select for GET /groups/{ID}
	AutoSubscribeNewMembers bool `json:"autoSubscribeNewMembers,omitempty"`

	// Describes a classification for the group (such as low, medium or high business impact).
	// Valid values for this property are defined by creating a ClassificationList setting value, based on the template
	// definition.
	// Returned by default.
	// Supports $filter (eq, ne, NOT, ge, le, startsWith)
	Classification string `json:"classification,omitempty"`

	// Timestamp of when the group was created.
	// The value cannot be modified and is automatically populated when the group is created. The Timestamp type
	// represents date and time information using ISO 8601 format and is always in UTC time.
	// For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
	// Returned by default.
	// Supports $filter (eq, ne, NOT, ge, le, in).
	// Read-only.
	CreatedDateTime string `json:"createdDateTime,omitempty"`

	// For some Azure Active Directory objects (user, group, application), if the object is deleted, it is first
	// logically deleted, and this property is updated with the date and time when the object was deleted. Otherwise,
	// this property is null. If the object is restored, this property is updated to null.
	DeletedDateTime string `json:"deletedDateTime,omitempty"`

	// An optional description for the group.
	// Returned by default.
	// Supports $filter (eq, ne, NOT, ge, le, startsWith) and $search.
	Description string `json:"description,omitempty"`

	// The display name for the group.
	// This property is required when a group is created and cannot be cleared during updates.
	// Returned by default.
	// Supports $filter (eq, ne, NOT, ge, le, in, startsWith), $search, and $orderBy.
	DisplayName string `json:"displayName,omitempty"`

	// Timestamp of when the group is set to expire. The value cannot be modified and is automatically populated when
	// the group is created. The Timestamp type represents date and time information using ISO 8601 format and is always
	// in UTC time.
	// For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
	// Returned by default.
	// Supports $filter (eq, ne, NOT, ge, le, in).
	// Read-only.
	ExpirationDateTime string `json:"expirationDateTime,omitempty"`

	// Specifies the group type and its membership.
	// If the collection contains Unified, the group is a Microsoft 365 group; otherwise, it's either a security group
	// or distribution group. For details, see groups overview.
	// If the collection includes DynamicMembership, the group has dynamic membership; otherwise, membership is static.
	// Returned by default.
	// Supports $filter (eq, NOT).
	GroupTypes []string `json:"groupTypes,omitempty"`

	// Indicates whether there are members in this group that have license errors from its group-based license
	// assignment.
	// This property is never returned on a GET operation.
	// You can use it as a $filter argument to get groups that have members with license errors (that is, filter for
	// this property being true)
	// Supports $filter (eq).
	HasMembersWithLicenseErrors bool `json:"hasMembersWithLicenseErrors,omitempty"`

	// True if the group is not displayed in certain parts of the Outlook UI: the Address Book, address lists for
	// selecting message recipients, and the Browse Groups dialog for searching groups; otherwise, false.
	// Default value is false.
	// Returned only on $select for GET /groups/{ID}
	HideFromAddressLists bool `json:"hideFromAddressLists,omitempty"`

	// True if the group is not displayed in Outlook clients, such as Outlook for Windows and Outlook on the web;
	// otherwise, false.
	// Default value is false.
	// Returned only on $select for GET /groups/{ID}
	HideFromOutlookClients bool `json:"hideFromOutlookClients,omitempty"`

	// Indicates whether this group can be assigned to an Azure Active Directory role or not.
	// Optional.
	// This property can only be set while creating the group and is immutable. If set to true, the securityEnabled
	// property must also be set to true and the group cannot be a dynamic group (that is, groupTypes cannot contain
	// DynamicMembership). Only callers in Global administrator and Privileged role administrator roles can set this
	// property. The caller must be assigned the RoleManagement.ReadWrite.Directory permission to set this property or
	// update the membership of such groups. For more, see Using a group to manage Azure AD role assignments
	// Returned by default.
	// Supports $filter (eq, ne, NOT).
	IsAssignableToRole bool `json:"isAssignableToRole,omitempty"`

	// Indicates whether the signed-in user is subscribed to receive email conversations.
	// Default value is true.
	// Returned only on $select for GET /groups/{ID}
	IsSubscribedByMail bool `json:"isSubscribedByMail,omitempty"`

	// Indicates status of the group license assignment to all members of the group.
	// Default value is false.
	// Read-only.
	// Returned only on $select.
	LicenseProcessingState enums.LicenseProcessingState `json:"licenseProcessingState,omitempty"`

	// The SMTP address for the group, for example, "serviceadmins@contoso.onmicrosoft.com".
	// Returned by default.
	// Read-only.
	// Supports $filter (eq, ne, NOT, ge, le, in, startsWith).
	Mail string `json:"mail,omitempty"`

	// Specifies whether the group is mail-enabled.
	// Required.
	// Returned by default.
	// Supports $filter (eq, ne, NOT).
	MailEnabled bool `json:"mailEnabled,omitempty"`

	// The mail alias for the group, unique in the organization.
	// Maximum length is 64 characters.
	// This property can contain only characters in the ASCII character set 0 - 127 except: @ () \ [] " ; : . <> , SPACE
	// Required.
	// Returned by default.
	// Supports $filter (eq, ne, NOT, ge, le, in, startsWith).
	MailNickname string `json:"mailNickname,omitempty"`

	// The rule that determines members for this group if the group is a dynamic group (groupTypes contains
	// DynamicMembership). For more information about the syntax of the membership rule, see Membership Rules syntax.
	// Returned by default.
	// Supports $filter (eq, ne, NOT, ge, le, startsWith).
	MembershipRule string `json:"membershipRule,omitempty"`

	// Indicates whether the dynamic membership processing is on or paused.
	// Returned by default.
	// Supports $filter (eq, ne, NOT, in).
	MembershipRuleProcessingState enums.RuleProcessingState `json:"membershipRuleProcessingState,omitempty"`

	// Indicates the last time at which the group was synced with the on-premises directory.
	// The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time.
	// For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
	// Returned by default.
	// Read-only.
	// Supports $filter (eq, ne, NOT, ge, le, in).
	OnPremisesLastSyncDateTime string `json:"onPremisesLastSyncDateTime,omitempty"`

	// Errors when using Microsoft synchronization product during provisioning.
	// Returned by default.
	// Supports $filter (eq, NOT).
	OnPremisesProvisioningErrors []OnPremisesProvisioningError `json:"onPremisesProvisioningErrors,omitempty"`

	// Contains the on-premises SAM account name synchronized from the on-premises directory.
	// The property is only populated for customers who are synchronizing their on-premises directory to Azure Active
	// Directory via Azure AD Connect.
	// Returned by default.
	// Supports $filter (eq, ne, NOT, ge, le, in, startsWith).
	// Read-only.
	OnPremisesSamAccountName string `json:"onPremisesSamAccountName,omitempty"`

	// Contains the on-premises security identifier (SID) for the group that was synchronized from on-premises to the
	// cloud.
	// Returned by default.
	// Supports $filter on null values.
	// Read-only.
	OnPremisesSecurityIdentifier string `json:"onPremisesSecurityIdentifier,omitempty"`

	// true if this group is synced from an on-premises directory; false if this group was originally synced from an
	// on-premises directory but is no longer synced; null if this object has never been synced from an on-premises
	// directory (default).
	// Returned by default.
	// Read-only.
	// Supports $filter (eq, ne, NOT, in).
	OnPremisesSyncEnabled bool `json:"onPremisesSyncEnabled,omitempty"`

	// The preferred data location for the Microsoft 365 group.
	// By default, the group inherits the group creator's preferred data location. To set this property, the calling
	// user must be assigned one of the following Azure AD roles:
	// - Global Administrator
	// - User Account Administrator
	// - Directory Writer
	// - Exchange Administrator
	// - SharePoint Administrator
	//
	// Nullable.
	// Returned by default.
	PreferredDataLocation string `json:"preferredDataLocation,omitempty"`

	// The preferred language for a Microsoft 365 group.
	// Should follow ISO 639-1 Code; for example en-US.
	// Returned by default.
	// Supports $filter (eq, ne, NOT, ge, le, in, startsWith).
	PreferredLanguage string `json:"preferredLanguage,omitempty"`

	// Email addresses for the group that direct to the same group mailbox.
	// For example: ["SMTP: bob@contoso.com", "smtp: bob@sales.contoso.com"].
	// The any operator is required to filter expressions on multi-valued properties.
	// Returned by default.
	// Read-only.
	// Not nullable.
	// Supports $filter (eq, NOT, ge, le, startsWith).
	ProxyAddresses []string `json:"proxyAddresses,omitempty"`

	// Timestamp of when the group was last renewed.
	// This cannot be modified directly and is only updated via the renew service action.
	// The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time.
	// For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
	// Returned by default.
	// Supports $filter (eq, ne, NOT, ge, le, in).
	// Read-only.
	RenewedDateTime string `json:"renewedDateTime,omitempty"`

	// Specifies the group behaviors that can be set for a Microsoft 365 group during creation.
	// This can be set only as part of creation (POST).
	ResourceBehaviorOptions []enums.ResourceBehavior `json:"resourceBehaviorOptions,omitempty"`

	// Specifies the group resources that are provisioned as part of Microsoft 365 group creation, that are not normally
	// part of default group creation.
	ResourceProvisioningOptions []enums.ResourceProvisioning `json:"resourceProvisioningOptions,omitempty"`

	// Specifies whether the group is a security group.
	// Required.
	// Returned by default.
	// Supports $filter (eq, ne, NOT, in).
	SecurityEnabled bool `json:"securityEnabled,omitempty"`

	// Security identifier of the group, used in Windows scenarios.
	// Returned by default.
	SecurityIdentifier string `json:"securityIdentifier,omitempty"`

	// Specifies a Microsoft 365 group's color theme. Possible values are Teal, Purple, Green, Blue, Pink, Orange or Red
	Theme string `json:"theme,omitempty"`

	// Count of conversations that have received new posts since the signed-in user last visited the group.
	// Returned only on $select for GET /groups/{ID}
	UnseenCount int32 `json:"unseenCount,omitempty"`

	// Specifies the group join policy and group content visibility for groups.
	// Possible values are: Private, Public, or Hiddenmembership.
	// Hiddenmembership can be set only for Microsoft 365 groups, when the groups are created.
	// It can't be updated later. Other values of visibility can be updated after group creation.
	// If visibility value is not specified during group creation on Microsoft Graph, a security group is created as
	// Private by default and Microsoft 365 group is Public. Groups assignable to roles are always Private.
	// Returned by default.
	// Nullable.
	Visibility enums.GroupVisibility `json:"visibility,omitempty"`
}
