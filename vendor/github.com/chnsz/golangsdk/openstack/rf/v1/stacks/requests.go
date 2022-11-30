package stacks

import (
	"github.com/chnsz/golangsdk"
)

// CreateOpts is the structure used to create a stack.
type CreateOpts struct {
	// The name of the stack.
	// The value can contain 1 to 64 characters, only letters, digits and hyphens (-) are allowed.
	// The name must start with a letter and end with a letter or digit.
	Name string `json:"stack_name" required:"true"`
	// The agencies authorized to IAC.
	Agencies []Agency `json:"agencies,omitempty"`
	// The description of the stack.
	// The value can contain 0 to 1024 characters.
	Description string `json:"description,omitempty"`
	// The flag of the deletion protection.
	// The resource stack deletion protection is not enabled by default (the resource stack is not allowed to be deleted
	// after the deletion protection is enabled).
	EnableDeletionProtection *bool `json:"enable_deletion_protection,omitempty"`
	// The automatic rollback flag.
	// The automatic rollback of the resource stack is not enabled by default (after automatic rollback is enabled, if
	// the deployment fails, it will automatically roll back and return to the previous stable state).
	EnableAutoRollback *bool `json:"enable_auto_rollback,omitempty"`
	// The HCL template content for deployment resources.
	TemplateBody string `json:"template_body,omitempty"`
	// The OBS address where the HCL template ZIP is located, which describes the target status of the deployment
	// resources.
	TemplateUri string `json:"template_uri,omitempty"`
	// The variable content for deployment resources.
	VarsBody string `json:"vars_body,omitempty"`
	// The variable structures for deployment resources.
	VarsStructure []VarsStructure `json:"vars_structure,omitempty"`
	// The OBS address where the variable ZIP corresponding to the HCL template is located, which describes the target
	// status of the deployment resources.
	VarsUri string `json:"vars_uri,omitempty"`
}

// Agency is an object that represents the IAC agency configuration.
type Agency struct {
	// The name of the provider corresponding to the IAM agency.
	// If the provider_name given by the user contains duplicate values, return 400.
	ProviderName string `json:"provider_name" required:"true"`
	// The name of IAM agency authorized to IAC account.
	// RF will use the agency authority to access and create resources corresponding to the provider.
	AgencyName string `json:"agency_name" required:"true"`
}

// VarsStructure is an object that represents the variable structure details.
type VarsStructure struct {
	// The variable key.
	VarKey string `json:"var_key" required:"true"`
	// The variable value.
	VarValue string `json:"var_value" required:"true"`
	// The encryption configuration.
	Encryption Encryption `json:"encryption,omitempty"`
}

// Encryption is an object that represents the KMS usage for variables.
type Encryption struct {
	// Encryption configuration for variables.
	Kms KmsStructure `json:"kms" required:"true"`
}

// KmsStructure is an object that represents the encrypt structure.
type KmsStructure struct {
	// The ID of the KMD secret key.
	ID string `json:"id" required:"true"`
	// The ciphertext corresponding to the data encryption key.
	CipherText string `json:"cipher_text" required:"true"`
}

// Create is a method to create a stack using create option.
func Create(client *golangsdk.ServiceClient, opts CreateOpts) (*CreateResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var r CreateResp
	_, err = client.Post(rootURL(client), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: client.MoreHeaders,
	})
	return &r, err
}

// ListAll is a method to query all stacks and returns a stack list.
func ListAll(client *golangsdk.ServiceClient) ([]Stack, error) {
	var r listResp
	_, err := client.Get(rootURL(client), &r, &golangsdk.RequestOpts{
		MoreHeaders: client.MoreHeaders,
	})
	if err != nil {
		return nil, err
	}

	return r.Stacks, nil
}

// DeployOpts is the structure used to deploy a terraform script.
type DeployOpts struct {
	// The HCL template content for terraform resources.
	TemplateBody string `json:"template_body,omitempty"`
	// The OBS address where the HCL template ZIP is located, which describes the target status of the terraform
	// resources.
	TemplateUri string `json:"template_uri,omitempty"`
	// The variable content for terraform resources.
	VarsBody string `json:"vars_body,omitempty"`
	// The variable structures for terraform resources.
	VarsStructure []VarsStructure `json:"vars_structure,omitempty"`
	// The OBS address where the variable ZIP corresponding to the HCL template is located, which describes the target
	// status of the deployment resources.
	VarsUri string `json:"vars_uri,omitempty"`
	// The unique ID of the resource stack.
	StackId string `json:"stack_id,omitempty"`
}

// Deploy is a method to deploy a terraform script using given parameters.
func Deploy(client *golangsdk.ServiceClient, stackName string, opts DeployOpts) (string, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return "", err
	}

	var r deployResp
	_, err = client.Post(deploymentURL(client, stackName), b, &r, &golangsdk.RequestOpts{
		MoreHeaders: client.MoreHeaders,
	})
	return r.DeploymentId, err
}

// ListEventsOpts allows to filter list data using given parameters.
type ListEventsOpts struct {
	// The unique ID of the resource stack.
	StackId string `q:"stack_id"`
	// The unique ID of the resource deployment.
	DeploymentId string `q:"deployment_id"`
}

// ListAllEvents is a method to query all events for a specified stack using given parameters.
func ListAllEvents(client *golangsdk.ServiceClient, stackName string, opts ListEventsOpts) ([]StackEvent, error) {
	url := eventURL(client, stackName)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var r listEventsResp
	_, err = client.Get(url, &r, &golangsdk.RequestOpts{
		MoreHeaders: client.MoreHeaders,
	})
	if err != nil {
		return nil, err
	}

	return r.StackEvents, nil
}

// DeleteOpts is the structure used to remove a stack.
type DeleteOpts struct {
	// The stack ID.
	StackId string `json:"stack_id,omitempty"`
}

// Delete is a method to remove a stack using stack name and delete option.
func Delete(client *golangsdk.ServiceClient, stackName string, opts DeleteOpts) error {
	url := resourceURL(client, stackName)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return err
	}
	url += query.String()

	_, err = client.Delete(url, &golangsdk.RequestOpts{
		MoreHeaders: client.MoreHeaders,
	})
	return err
}
