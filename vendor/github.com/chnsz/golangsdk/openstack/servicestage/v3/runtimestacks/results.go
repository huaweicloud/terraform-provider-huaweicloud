package runtimestacks

import "github.com/chnsz/golangsdk/pagination"

// RuntimeStack is the structure that represents the detail of the ServiceStage component runtime stack.
type RuntimeStack struct {
	// The runtime stack name.
	Name string `json:"name"`
	// The runtime stack deploy mode.
	DeployMode string `json:"deploy_mode"`
	// The runtime stack version.
	Version string `json:"version"`
	// The runtime stack type, just like Nodejs.
	Type string `json:"type"`
	// The runtime stack status, just like
	Status string `json:"status"`
}

// ApplicationPage is a single page maximum result representing a query by offset page.
type ApplicationPage struct {
	pagination.OffsetPageBase
}
