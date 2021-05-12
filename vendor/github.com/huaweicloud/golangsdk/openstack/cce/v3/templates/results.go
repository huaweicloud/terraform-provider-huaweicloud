package templates

import "github.com/huaweicloud/golangsdk"

type commonResult struct {
	golangsdk.Result
}

type ListResutlt struct {
	commonResult
}

type Template struct {
	Kind       string   `json:"kind"`
	ApiVersion string   `json:"apiVersion"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
}

type Metadata struct {
	UID               string `json:"uid"`
	Name              string `json:"name"`
	CreationTimestamp string `json:"creationTimestamp"`
	UpdateTimestamp   string `json:"updateTimestamp"`
}

type Spec struct {
	Type        string     `json:"type"`
	Labels      []string   `json:"labels"`
	LogoURL     string     `json:"logoURL"`
	Description string     `json:"description"`
	Versions    []Versions `json:"versions"`
}

type Versions struct {
	Version string      `json:"version"`
	Input   interface{} `json:"input"`
}

func (r ListResutlt) Extract() ([]Template, error) {
	var s struct {
		Templates []Template `json:"items"`
	}
	err := r.ExtractInto(&s)
	return s.Templates, err
}
