package aliases

import "github.com/chnsz/golangsdk"

// CreateOpts is a structure that used to create a alias for specified version.
type CreateOpts struct {
	// The URN of the function to which the alias and version are belong.
	FunctionUrn string `json:"-" required:"true"`
	// Function alias to be created.
	Name string `json:"name" required:"true"`
	// Version corresponding to the alias.
	Version string `json:"version" required:"true"`
	// Description of the alias.
	Description string `json:"description,omitempty"`
	// The weights configuration of the additional version.
	AdditionalVersionWeights map[string]interface{} `json:"additional_version_weights,omitempty"`
}

var requestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

// Create is a method to create a alias for specified version using given parameters.
func Create(client *golangsdk.ServiceClient, opts CreateOpts) (*Alias, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Alias
	_, err = client.Post(rootURL(client, opts.FunctionUrn), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Get is a method to obtain a specified alias using given parameters.
func Get(client *golangsdk.ServiceClient, functionUrn, aliasName string) (*Alias, error) {
	var r Alias
	_, err := client.Get(resourceURL(client, functionUrn, aliasName), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// List is a method to obtain all aliases using given parameters.
func List(client *golangsdk.ServiceClient, functionUrn string) ([]Alias, error) {
	var r []Alias
	_, err := client.Get(rootURL(client, functionUrn), &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return r, err
}

// UpdateOpts is the structure used to update a specified alias.
type UpdateOpts struct {
	// The URN of the function to which the alias and version are belong.
	FunctionUrn string `json:"-" required:"true"`
	// The name of the alias.
	Name string `json:"-" required:"true"`
	// Version corresponding to the alias.
	Version string `json:"version" required:"true"`
	// Description of the alias.
	Description string `json:"description,omitempty"`
	// The weights configuration of the additional version.
	AdditionalVersionWeights map[string]interface{} `json:"additional_version_weights,omitempty"`
}

// Update is a method to update a specified alias using given parameters.
func Update(client *golangsdk.ServiceClient, opts UpdateOpts) (*Alias, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r Alias
	_, err = client.Put(resourceURL(client, opts.FunctionUrn, opts.Name), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return &r, err
}

// Delete is a method to delete a specified alias.
func Delete(client *golangsdk.ServiceClient, functionUrn, aliasName string) error {
	_, err := client.Delete(resourceURL(client, functionUrn, aliasName), &golangsdk.RequestOpts{
		MoreHeaders: requestOpts.MoreHeaders,
	})
	return err
}
