package security

import "github.com/chnsz/golangsdk"

const rootPath = "OS-SECURITYPOLICY"

func passwordPolicyURL(client *golangsdk.ServiceClient, domainID string) string {
	return client.ServiceURL(rootPath, "domains", domainID, "password-policy")
}

func protectPolicyURL(client *golangsdk.ServiceClient, domainID string) string {
	return client.ServiceURL(rootPath, "domains", domainID, "protect-policy")
}
