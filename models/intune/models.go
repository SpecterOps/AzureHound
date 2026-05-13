// File: models/intune/models.go
// Copyright (C) 2022 SpecterOps
// Data models for Intune integration

package intune

import (
	"time"
)

// ManagedDevice represents an Intune managed device
type ManagedDevice struct {
	Id                    string    `json:"id"`
	DeviceName           string    `json:"deviceName"`
	OperatingSystem      string    `json:"operatingSystem"`
	OSVersion            string    `json:"osVersion"`
	ComplianceState      string    `json:"complianceState"`
	LastSyncDateTime     time.Time `json:"lastSyncDateTime"`
	EnrollmentType       string    `json:"enrollmentType"`
	ManagementAgent      string    `json:"managementAgent"`
	AzureADDeviceId      string    `json:"azureADDeviceId"`
	UserPrincipalName    string    `json:"userPrincipalName"`
	DeviceEnrollmentType string    `json:"deviceEnrollmentType"`
	JoinType             string    `json:"joinType"`
}

// ComplianceState represents device compliance information
type ComplianceState struct {
	Id                                     string    `json:"id"`
	DeviceId                              string    `json:"deviceId"`
	DeviceName                            string    `json:"deviceName"`
	ComplianceGracePeriodExpirationDateTime time.Time `json:"complianceGracePeriodExpirationDateTime"`
	State                                 string    `json:"state"`
	Version                               int       `json:"version"`
	SettingStates                         []ComplianceSettingState `json:"settingStates"`
}

// ComplianceSettingState represents individual compliance setting state
type ComplianceSettingState struct {
	Setting      string `json:"setting"`
	State        string `json:"state"`
	CurrentValue string `json:"currentValue"`
}

// ConfigurationState represents device configuration state
type ConfigurationState struct {
	Id            string    `json:"id"`
	DeviceId      string    `json:"deviceId"`
	DeviceName    string    `json:"deviceName"`
	State         string    `json:"state"`
	Version       int       `json:"version"`
	SettingStates []ConfigurationSettingState `json:"settingStates"`
	PlatformType  string    `json:"platformType"`
}

// ConfigurationSettingState represents individual configuration setting state
type ConfigurationSettingState struct {
	Setting      string `json:"setting"`
	State        string `json:"state"`
	CurrentValue string `json:"currentValue"`
}