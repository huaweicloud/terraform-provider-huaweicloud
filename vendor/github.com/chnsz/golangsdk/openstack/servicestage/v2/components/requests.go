package components

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

var requestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// CreateOpts is the structure required by the Create method to create a new component.
type CreateOpts struct {
	// Application component name.
	// The value can contain 2 to 64 characters, including letters, digits, hyphens (-), and underscores (_).
	// It must start with a letter and end with a letter or digit.
	Name string `json:"name" required:"true"`
	// Runtime.
	// The value can be obtained from type_name returned by the API in Obtaining All Supported Runtimes of Application
	// Components.
	Runtime string `json:"runtime" required:"true"`
	// Application component type. Example: Webapp, MicroService, or Common.
	Type string `json:"category" required:"true"`
	// Application component sub-type.
	// Webapp sub-types include Web.
	// MicroService sub-types include Java Chassis, Go Chassis, Mesher, Spring Cloud, and Dubbo.
	// Common sub-type can be empty.
	Framwork string `json:"sub_category,omitempty"`
	// Description.
	// The value can contain up to 128 characters.
	Description string `json:"description,omitempty"`
	// Source of the code or software package.
	Source *Source `json:"source,omitempty"`
	// Component builder.
	Builder *Builder `json:"build,omitempty"`
}

// Source is an object to specified the source information of Open-Scoure codes or package storage.
type Source struct {
	// Type. Option: source code or artifact software package.
	Kind string `json:"kind" required:"true"`
	// The details about the Repository source code and the Artifact software package.
	Spec Spec `json:"spec" required:"true"`
}

// Spec is an object to specified the configuration of repository or storage.
type Spec struct {
	// The parameters of code are as follows:
	// Code repository. Value: GitHub, GitLab, Gitee, or Bitbucket.
	RepoType string `json:"repo_type,omitempty"`
	// Code repository URL. Example: https://github.com/example/demo.git.
	RepoUrl string `json:"repo_url,omitempty"`
	// Authorization name, which can be obtained from the authorization list.
	RepoAuth string `json:"repo_auth,omitempty"`
	// The code's organization. Value: GitHub, GitLab, Gitee, or Bitbucket.
	RepoNamespace string `json:"repo_namespace,omitempty"`
	// Code branch or tag. Default value: master.
	RepoRef string `json:"repo_ref,omitempty"`

	// The parameters of artifact are as follows:
	// Storage mode. Value: swr or obs.
	Storage string `json:"storage,omitempty"`
	// Type. Value: package.
	Type string `json:"type,omitempty"`
	// Address of the software package or source code.
	Url string `json:"url,omitempty"`
	// Authentication mode. Value: iam or none. Default value: iam.
	Auth string `json:"auth,omitempty"`
	// Other attributes of the software package. You need to add these attributes only when you set storage to obs.
	Properties Properties `json:"properties,omitempty"`
}

// Properties is an object to specified the other configuration of the software package for OBS bucket.
type Properties struct {
	// Object Storage Service (OBS) endpoint address. Example: https://obs.region_id.external_domain_name.com.
	Endpoint string `json:"endpoint,omitempty"`
	// Name of the OBS bucket where the software package is stored.
	Bucket string `json:"bucket,omitempty"`
	// Object in the OBS bucket, which is usually the name of the software package.
	// If there is a folder, the path of the folder must be added. Example: test.jar or demo/test.jar.
	Key string `json:"key,omitempty"`
}

// Builder is the component builder, the configuration details refer to 'Parameter'.
type Builder struct {
	// This parameter is provided only when no ID is available during build creation.
	Parameter Parameter `json:"parameters" required:"true"`
}

// Parameter is an object to specified the building configuration of codes or package.
type Parameter struct {
	// Compilation command. By default:
	// When build.sh exists in the root directory, the command is ./build.sh.
	// When build.sh does not exist in the root directory, the command varies depending on the operating system (OS). Example:
	// Java and Tomcat: mvn clean package
	// Nodejs: npm build
	BuildCmd string `json:"build_cmd,omitempty"`
	// Address of the Docker file. By default, the Docker file is in the root directory (./).
	DockerfilePath string `json:"dockerfile_path,omitempty"`
	// Build archive organization. Default value: cas_{project_id}.
	ArtifactNamespace string `json:"artifact_namespace,omitempty"`
	// The ID of the cluster to be built.
	ClusterId string `json:"cluster_id,omitempty"`
	// The name of the cluster to be built.
	ClusterName string `json:"clsuter_name,omitempty"`
	// The type of the cluster to be built.
	ClusterType string `json:"cluster_type,omitempty"`
	// key indicates the key of the tag, and value indicates the value of the tag.
	UsePublicCluster bool `json:"use_public_cluster,omitempty"`
	// key indicates the key of the tag, and value indicates the value of the tag.
	NodeLabelSelector map[string]interface{} `json:"node_label_selector,omitempty"`
}

// Create is a method to create a component in the specified appliation using given parameters.
func Create(c *golangsdk.ServiceClient, appId string, opts CreateOpts) (*Component, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Component
	_, err = c.Post(rootURL(c, appId), b, &r, nil)
	return &r, err
}

// Get is a method to retrieves a particular configuration based on its unique ID.
func Get(c *golangsdk.ServiceClient, appId, componentId string) (*Component, error) {
	var r Component
	_, err := c.Get(resourceURL(c, appId, componentId), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// Number of records to be queried.
	// Value range: 0–100, or 1000.
	// Default value: 1000, indicating that a maximum of 1000 records can be queried and all records are displayed on
	// the same page.
	Limit string `q:"limit"`
	// The offset number.
	Offset int `q:"offset"`
	// Sorting field. By default, query results are sorted by creation time.
	// The following enumerated values are supported: create_time, name, and update_time.
	OrderBy string `q:"order_by"`
	// Descending or ascending order. Default value: desc.
	Order string `q:"order"`
}

// List is a method to query the list of the components using given opts.
func List(c *golangsdk.ServiceClient, appId string, opts ListOpts) ([]Component, error) {
	url := rootURL(c, appId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	pages, err := pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		p := ComponentPage{pagination.OffsetPageBase{PageResult: r}}
		return p
	}).AllPages()

	if err != nil {
		return nil, err
	}
	return ExtractComponents(pages)
}

// UpdateOpts is the structure required by the Update method to update the component configuration.
type UpdateOpts struct {
	// Application component name.
	// The value can contain 2 to 64 characters, including letters, digits, hyphens (-), and underscores (_).
	// It must start with a letter and end with a letter or digit.
	Name string `json:"name,omitempty"`
	// Description.
	// The value can contain up to 128 characters.
	Description *string `json:"description,omitempty"`
	// Source of the code or software package.
	Source *Source `json:"source,omitempty"`
	// Component build.
	Builder *Builder `json:"build,omitempty"`
}

// Update is a method to update the component configuration, such as name, description, builder and code source.
func Update(c *golangsdk.ServiceClient, appId, componentId string, opts UpdateOpts) (*Component, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Component
	_, err = c.Put(resourceURL(c, appId, componentId), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Delete is a method to delete an existing component from a specified application.
func Delete(c *golangsdk.ServiceClient, appId, componentId string) error {
	_, err := c.Delete(resourceURL(c, appId, componentId), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
