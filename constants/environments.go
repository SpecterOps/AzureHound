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

package constants

// Azure deployment regions
const (
	China   string = "china"
	Cloud   string = "cloud"
	USGovL4 string = "usgovl4"
	USGovL5 string = "usgovl5"
)

type Environment struct {
	ActiveDirectoryAuthority string
	MicrosoftGraphUrl        string
	ResourceManagerUrl       string
}

func AzureCloud() Environment {
	return Environment{
		"https://login.microsoftonline.com",
		"https://graph.microsoft.com",
		"https://management.azure.com",
	}
}

func AzureUSGovernment() Environment {
	return Environment{
		"https://login.microsoftonline.us",
		"https://graph.microsoft.us",
		"https://management.usgovcloudapi.net",
	}
}

func AzureUSGovernmentL5() Environment {
	env := AzureUSGovernment()
	env.MicrosoftGraphUrl = "https://dod-graph.microsoft.us"
	return env
}

func AzureChina() Environment {
	return Environment{
		"https://login.chinacloudapi.cn",
		"https://microsoftgraph.chinacloudapi.cn",
		"https://management.chinacloudapi.cn",
	}
}
