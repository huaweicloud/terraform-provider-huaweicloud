package acl

import "github.com/chnsz/golangsdk"

const rootPath = "OS-SECURITYPOLICY"

func consoleACLPolicyURL(client *golangsdk.ServiceClient, domainID string) string {
	return client.ServiceURL(rootPath, "domains", domainID, "console-acl-policy")
}

func apiACLPolicyURL(client *golangsdk.ServiceClient, domainID string) string {
	return client.ServiceURL(rootPath, "domains", domainID, "api-acl-policy")
}
