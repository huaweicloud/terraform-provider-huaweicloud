package vpcattachments

import "github.com/chnsz/golangsdk"

func rootURL(client *golangsdk.ServiceClient, instanceId string) string {
	return client.ServiceURL("enterprise-router", instanceId, "vpc-attachments")
}

func resourceURL(client *golangsdk.ServiceClient, instanceId, attachmentId string) string {
	return client.ServiceURL("enterprise-router", instanceId, "vpc-attachments", attachmentId)
}
