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

import "strings"

type DomainServiceLDAPSSettings struct {
	LDAPS          string `json:"ldaps,omitempty"`
	ExternalAccess string `json:"externalAccess,omitempty"`
}

type DomainServiceSecuritySettings struct {
	NTLMV1                   string `json:"ntlmV1,omitempty"`
	TLSV1                    string `json:"tlsV1,omitempty"`
	SyncNTLMPasswords        string `json:"syncNtlmPasswords,omitempty"`
	SyncKerberosPasswords    string `json:"syncKerberosPasswords,omitempty"`
	SyncOnPremPasswords      string `json:"syncOnPremPasswords,omitempty"`
	KerberosRC4Encryption    string `json:"kerberosRc4Encryption,omitempty"`
	KerberosArmoring         string `json:"kerberosArmoring,omitempty"`
	LDAPSigning              string `json:"ldapSigning,omitempty"`
	ChannelBinding           string `json:"channelBinding,omitempty"`
	SyncOnPremSAMAccountName string `json:"syncOnPremSamAccountName,omitempty"`
}

type DomainServiceProperties struct {
	TenantID                string                        `json:"tenantId,omitempty"`
	DomainName              string                        `json:"domainName,omitempty"`
	DomainConfigurationType string                        `json:"domainConfigurationType,omitempty"`
	FilteredSync            string                        `json:"filteredSync,omitempty"`
	SyncScope               string                        `json:"syncScope,omitempty"`
	SyncApplicationID       string                        `json:"syncApplicationId,omitempty"`
	DomainSecuritySettings  DomainServiceSecuritySettings `json:"domainSecuritySettings,omitempty"`
	LDAPSSettings           DomainServiceLDAPSSettings    `json:"ldapsSettings,omitempty"`
}

type DomainService struct {
	Entity

	Location   string                  `json:"location,omitempty"`
	Name       string                  `json:"name,omitempty"`
	Properties DomainServiceProperties `json:"properties,omitempty"`
	Type       string                  `json:"type,omitempty"`
}

func (s DomainService) ResourceGroupName() string {
	parts := strings.Split(s.Id, "/")
	if len(parts) > 4 {
		return parts[4]
	}
	return ""
}

func (s DomainService) ResourceGroupID() string {
	parts := strings.Split(s.Id, "/")
	if len(parts) > 5 {
		return strings.Join(parts[:5], "/")
	}
	return ""
}
