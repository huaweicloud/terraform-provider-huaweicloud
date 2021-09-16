package thesaurus

import "github.com/chnsz/golangsdk"

// loadIKThesaurusURL /v1.0/{project_id}/clusters/{cluster_id}/thesaurus
func loadURL(c *golangsdk.ServiceClient, clusterId string) string {
	return c.ServiceURL("clusters", clusterId, "thesaurus")
}

// queryIKThesaurusStatusURL /v1.0/{project_id}/clusters/{cluster_id}/thesaurus
func getURL(c *golangsdk.ServiceClient, clusterId string) string {
	return c.ServiceURL("clusters", clusterId, "thesaurus")
}

// deleteIKThesaurusURL /v1.0/{project_id}/clusters/{cluster_id}/thesaurus
func deleteURL(c *golangsdk.ServiceClient, clusterId string) string {
	return c.ServiceURL("clusters", clusterId, "thesaurus")
}
