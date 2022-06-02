package services

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cse/dedicated/v4/auth"
	"github.com/chnsz/golangsdk/openstack/cse/dedicated/v4/instances"
)

// CreateOpts is the structure required by the Create method to create a new dedicated microservice.
type CreateOpts struct {
	// Microservice information.
	Services Service `json:"service" required:"true"`
	// Blacklist and whitelist.
	Rules []Rule `json:"rules,omitempty"`
	// Instance information.
	Instances []instances.CreateOpts `json:"instances,omitempty"`
	// Extended attribute. You can customize a key and value. The value must be at least 1 byte long.
	Tags map[string]interface{} `json:"tags,omitempty"`
}

// Service is an object that specifies the microservice configuration details.
type Service struct {
	// Microservice name, which must be unique in an application. The value contains 1 to 128 characters.
	// Regular expression: ^[a-zA-Z0-9]*$|^[a-zA-Z0-9][a-zA-Z0-9_\-.]*[a-zA-Z0-9]$
	Name string `json:"serviceName" required:"true"`
	// Application ID, which must be unique. The value contains 1 to 160 characters.
	// Regular expression: ^[a-zA-Z0-9]*$|^[a-zA-Z0-9][a-zA-Z0-9_\-.]*[a-zA-Z0-9]$
	AppId string `json:"appId" required:"true"`
	// Microservice version. The value contains 1 to 64 characters. Regular expression: ^[0-9]$|^[0-9]+(.[0-9]+)$
	Version string `json:"version" required:"true"`
	// Microservice ID, which must be unique. The value contains 1 to 64 characters. Regular expression: ^.*$
	ID string `json:"serviceId,omitempty"`
	// Service stage. Value: development, testing, acceptance, or production.
	// Only when the value is development, testing, or acceptance, you can use the API for uploading schemas in batches
	// to add or modify an existing schema. Default value: development.
	Environment *string `json:"environment,omitempty"`
	// Microservice description. The value contains a maximum of 256 characters.
	Description string `json:"description,omitempty"`
	// Microservice level. Value: FRONT, MIDDLE, or BACK.
	Level string `json:"level,omitempty"`
	// Microservice registration mode. Value: SDK, PLATFORM, SIDECAR, or UNKNOWN.
	RegisterBy string `json:"registerBy,omitempty"`
	// Foreign key ID of a microservice access schema. The array length supports a maximum of 100 schemas.
	Schemas []string `json:"schemas,omitempty"`
	// Microservice status. Value: UP or DOWN. Default value: UP.
	Status string `json:"status,omitempty"`
	// Microservice registration time.
	Timestamp string `json:"timestamp,omitempty"`
	// Latest modification time (UTC).
	ModTimestamp string `json:"modTimestamp,omitempty"`
	// Development framework.
	Framework *Framework `json:"framework,omitempty"`
	// Service path.
	Paths []ServicePath `json:"paths,omitempty"`
}

// Framework is an object that specifies the configuration of the microservice framework.
type Framework struct {
	// Microservice development framework. Default value: UNKNOWN.
	Name string `json:"name,omitempty"`
	// Version of the microservice development framework.
	Version string `json:"version,omitempty"`
}

// ServicePath is an object that specifies the configuration of the service path.
type ServicePath struct {
	// Route address.
	Path string `json:"Path,omitempty"`
	// Extended attribute. You can customize a key and value. The value must be at least 1 byte long.
	Property map[string]interface{} `json:"Property,omitempty"`
}

// Rule is an object that specifies the configuration of the black list or white list.
type Rule struct {
	// Customized rule ID.
	RuleId string `json:"ruleId,omitempty"`
	// Rule type. Value: WHITE or BLACK.
	RuleType string `json:"ruleType,omitempty"`
	// If the value starts with tag_xxx, the attributes are filtered by Tag.
	// Otherwise, the attributes are filtered by serviceId, AppId, ServiceName, Version, Description, Level, or Status.
	Attribute string `json:"attribute,omitempty"`
	// Matching rule. The value is a regular expression containing 1 to 64 characters.
	Pattern string `json:"pattern,omitempty"`
	// Rule description.
	Description string `json:"description,omitempty"`
	// Time when a rule is created. This parameter is used only when you query rules.
	Timestamp string `json:"timestamp,omitempty"`
	// Update time.
	ModTimestamp string `json:"modTimestamp,omitempty"`
}

// Create is a method to create a microservice using given parameters.
func Create(c *golangsdk.ServiceClient, opts CreateOpts, token string) (*CreateResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r CreateResp
	_, err = c.Post(rootURL(c), b, &r, &golangsdk.RequestOpts{
		JSONResponse: &r,
		MoreHeaders:  auth.BuildMoreHeaderUsingToken(c, token),
	})
	return &r, err
}

// Get is a method to retrieves a particular configuration based on its unique ID (and token).
func Get(c *golangsdk.ServiceClient, serviceId, token string) (*ServiceResp, error) {
	var r struct {
		Service ServiceResp `json:"service"`
	}
	_, err := c.Get(resourceURL(c, serviceId), &r, &golangsdk.RequestOpts{
		MoreHeaders: auth.BuildMoreHeaderUsingToken(c, token),
	})
	return &r.Service, err
}

// DeleteOpts is the structure required by the Delete method to specified whether microservice is force delete.
type DeleteOpts struct {
	Force bool `q:"force"`
}

// Delete is a method to remove an existing microservice using its unique ID (and token).
func Delete(c *golangsdk.ServiceClient, opts DeleteOpts, engineId, token string) error {
	url := resourceURL(c, engineId)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return err
	}
	url += query.String()

	_, err = c.Delete(url, &golangsdk.RequestOpts{
		MoreHeaders: auth.BuildMoreHeaderUsingToken(c, token),
	})
	return err
}
