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

import "github.com/certmichelin/azurehound/v2/enums"

type SiteConfig struct {
	AcrUseManagedIdentityCreds             bool                             `json:"acrUseManagedIdentityCreds,omitempty"`
	AcrUserManagedIdentityID               string                           `json:"acrUserManagedIdentityID,omitempty"`
	AlwaysOn                               bool                             `json:"alwaysOn,omitempty"`
	ApiDefinition                          ApiDefinitionInfo                `json:"apiDefinition,omitempty"`
	ApiManagementConfig                    ApiManagementConfig              `json:"apiManagementConfig,omitempty"`
	AppCommandLine                         string                           `json:"appCommandLine,omitempty"`
	AppSettings                            []NameValuePair                  `json:"appSettings,omitempty"`
	AutoHealEnabled                        bool                             `json:"autoHealEnabled,omitempty"`
	AutoHealRules                          string                           `json:"autoHealRules,omitempty"`
	AutoSwapSlotName                       string                           `json:"autoSwapSlotName,omitempty"`
	AzureStorageAccounts                   map[string]AzureStorageInfoValue `json:"azureStorageAccounts,omitempty"`
	ConnectionStrings                      []ConnStringInfo                 `json:"connectionStrings,omitempty"`
	Cors                                   CorsSettings                     `json:"cors,omitempty"`
	DefaultDocuments                       []string                         `json:"defaultDocuments,omitempty"`
	DetailedErrorLoggingEnabled            bool                             `json:"detailedErrorLoggingEnabled,omitempty"`
	DocumentRoot                           string                           `json:"documentRoot,omitempty"`
	Experiments                            Experiments                      `json:"experiments,omitempty"`
	FtpsState                              enums.FtpsState                  `json:"ftpsState,omitempty"`
	FunctionAppScaleLimit                  int                              `json:"functionAppScaleLimit,omitempty"`
	FunctionsRuntimeScaleMonitoringEnabled bool                             `json:"functionsRuntimeScaleMonitoringEnabled,omitempty"`
	HandlerMappings                        []HandlerMapping                 `json:"handlerMappings,omitempty"`
	HealthCheckPath                        string                           `json:"healthCheckPath,omitempty"`
	Http20Enabled                          bool                             `json:"http20Enabled,omitempty"`
	HttpLoggingEnabled                     bool                             `json:"httpLoggingEnabled,omitempty"`
	IpSecurityRestrictions                 []IpSecurityRestriction          `json:"ipSecurityRestrictions,omitempty"`
	JavaContainer                          string                           `json:"javaContainer,omitempty"`
	JavaContainerVersion                   string                           `json:"javaContainerVersion,omitempty"`
	JavaVersion                            string                           `json:"javaVersion,omitempty"`
	KeyVaultReferenceIdentity              string                           `json:"keyVaultReferenceIdentity,omitempty"`
	Limits                                 SiteLimits                       `json:"limits,omitempty"`
	LinuxFxVersion                         string                           `json:"linuxFxVersion,omitempty"`
	LoadBalancing                          enums.SiteLoadBalancing          `json:"loadBalancing,omitempty"`
	LocalMySqlEnabled                      bool                             `json:"localMySqlEnabled,omitempty"`
	LogsDirectorySizeLimit                 int                              `json:"logsDirectorySizeLimit,omitempty"`
	MachineKey                             SiteMachineKey                   `json:"machineKey,omitempty"`
	ManagedPipelineMode                    enums.ManagedPipelineMode        `json:"managedPipelineMode,omitempty"`
	ManagedServiceIdentityId               int                              `json:"managedServiceIdentityId,omitempty"`
	MinTlsVersion                          enums.SupportedTlsVersions       `json:"minTlsVersion,omitempty"`
	MinimumElasticInstanceCount            int                              `json:"minimumElasticInstanceCount,omitempty"`
	NetFrameworkVersion                    string                           `json:"netFrameworkVersion,omitempty"`
	NodeVersion                            string                           `json:"nodeVersion,omitempty"`
	NumberOfWorkers                        int                              `json:"numberOfWorkers,omitempty"`
	PhpVersion                             string                           `json:"phpVersion,omitempty"`
	PowerShellVersion                      string                           `json:"powerShellVersion,omitempty"`
	PreWarmedInstanceCount                 int                              `json:"preWarmedInstanceCount,omitempty"`
	PublicNetworkAccess                    string                           `json:"publicNetworkAccess,omitempty"`
	PublishingUsername                     string                           `json:"publishingUsername,omitempty"`
	Push                                   PushSettings                     `json:"push,omitempty"`
	PythonVersion                          string                           `json:"pythonVersion,omitempty"`
	RemoteDebuggingEnabled                 bool                             `json:"remoteDebuggingEnabled,omitempty"`
	RemoteDebuggingVersion                 string                           `json:"remoteDebuggingVersion,omitempty"`
	RequestTracingEnabled                  bool                             `json:"requestTracingEnabled,omitempty"`
	RequestTracingExpirationTime           string                           `json:"requestTracingExpirationTime,omitempty"`
	ScmIpSecurityRestrictions              []IpSecurityRestriction          `json:"scmIpSecurityRestrictions,omitempty"`
	ScmIpSecurityRestrictionsUseMain       bool                             `json:"scmIpSecurityRestrictionsUseMain,omitempty"`
	ScmMinTlsVersion                       enums.SupportedTlsVersions       `json:"scmMinTlsVersion,omitempty"`
	ScmType                                enums.ScmType                    `json:"scmType,omitempty"`
	TracingOptions                         string                           `json:"tracingOptions,omitempty"`
	Use32BitWorkerProcess                  bool                             `json:"use32BitWorkerProcess,omitempty"`
	VirtualApplications                    []VirtualApplication             `json:"virtualApplications,omitempty"`
	VnetName                               string                           `json:"vnetName,omitempty"`
	VnetPrivatePortsCount                  int                              `json:"vnetPrivatePortsCount,omitempty"`
	VnetRouteAllEnabled                    bool                             `json:"vnetRouteAllEnabled,omitempty"`
	WebSocketsEnabled                      bool                             `json:"webSocketsEnabled,omitempty"`
	WebsiteTimeZone                        string                           `json:"websiteTimeZone,omitempty"`
	WindowsFxVersion                       string                           `json:"windowsFxVersion,omitempty"`
	XManagedServiceIdentityId              int                              `json:"xManagedServiceIdentityId,omitempty"`

	//Following ones have been found in testing, but not present in the documentation
	AntivirusScanEnabled                   bool        `json:"antivirusScanEnabled,omitempty"`
	AzureMonitorLogCategories              interface{} `json:"azureMonitorLogCategories,omitempty"`
	CustomAppPoolIdentityAdminState        interface{} `json:"customAppPoolIdentityAdminState,omitempty"`
	CustomAppPoolIdentityTenantState       interface{} `json:"customAppPoolIdentityTenantState,omitempty"`
	ElasticWebAppScaleLimit                interface{} `json:"elasticWebAppScaleLimit,omitempty"`
	FileChangeAuditEnabled                 bool        `json:"fileChangeAuditEnabled,omitempty"`
	Http20ProxyFlag                        interface{} `json:"http20ProxyFlag,omitempty"`
	IpSecurityRestrictionsDefaultAction    interface{} `json:"ipSecurityRestrictionsDefaultAction,omitempty"`
	Metadata                               interface{} `json:"metadata,omitempty"`
	MinTlsCipherSuite                      interface{} `json:"minTlsCipherSuite,omitempty"`
	PublishingPassword                     interface{} `json:"publishingPassword,omitempty"`
	RoutingRules                           interface{} `json:"routingRules,omitempty"`
	RuntimeADUser                          interface{} `json:"runtimeADUser,omitempty"`
	RuntimeADUserPassword                  interface{} `json:"runtimeADUserPassword,omitempty"`
	ScmIpSecurityRestrictionsDefaultAction interface{} `json:"scmIpSecurityRestrictionsDefaultAction,omitempty"`
	SitePort                               interface{} `json:"sitePort,omitempty"`
	StorageType                            interface{} `json:"storageType,omitempty"`
	SupportedTlsCipherSuites               interface{} `json:"supportedTlsCipherSuites,omitempty"`
	WinAuthAdminState                      interface{} `json:"winAuthAdminState,omitempty"`
	WinAuthTenantState                     interface{} `json:"winAuthTenantState,omitempty"`
}

type ApiDefinitionInfo struct {
	Url string `json:"url,omitempty"`
}

type ApiManagementConfig struct {
	Id string `json:"id,omitempty"`
}

type CorsSettings struct {
	AllowedOrigins     []string `json:"allowedOrigins,omitempty"`
	SupportCredentials bool     `json:"supportCredentials,omitempty"`
}

type Experiments struct {
	RampUpRules []RampUpRule `json:"rampUpRules,omitempty"`
}

type RampUpRule struct {
	ActionHostName            string `json:"actionHostName,omitempty"`
	ChangeDecisionCallbackUrl string `json:"changeDecisionCallbackUrl,omitempty"`
	ChangeIntervalInMinutes   int    `json:"changeIntervalInMinutes,omitempty"`
	ChangeStep                int    `json:"changeStep,omitempty"`
	MaxReroutePercentage      int    `json:"maxReroutePercentage,omitempty"`
	MinReroutePercentage      int    `json:"minReroutePercentage,omitempty"`
	Name                      string `json:"name,omitempty"`
	ReroutePercentage         int    `json:"reroutePercentage,omitempty"`
}

type HandlerMapping struct {
	Arguments       string `json:"arguments,omitempty"`
	Extension       string `json:"extension,omitempty"`
	ScriptProcessor string `json:"scriptProcessor,omitempty"`
}

type SiteLimits struct {
	MaxDiskSizeInMb  int `json:"maxDiskSizeInMb,omitempty"`
	MaxMemoryInMb    int `json:"maxMemoryInMb,omitempty"`
	MaxPercentageCpu int `json:"maxPercentageCpu,omitempty"`
}
