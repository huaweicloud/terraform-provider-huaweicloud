package keypairs

import (
	"github.com/chnsz/golangsdk"
)

// AssociateOpts is the request body of binding an SSH keypair
type AssociateOpts struct {
	// the SSH keypair name
	Name string `json:"keypair_name" required:"true"`
	// Information about the VM to which the key pair is to be bound
	Server EcsServerOpts `json:"server" required:"true"`
}

// DisassociateOpts is the request body of unbinding an SSH keypair
type DisassociateOpts struct {
	// the ID of the VM to which the SSH key pair is to be unbound
	ServerID string `json:"id" required:"true"`
	// server authentication object, this parameter is required when the server is poweron
	Auth *AuthOpts `json:"auth,omitempty"`
}

// EcsServerOpts is object about the VM to which the key pair is to be bound
type EcsServerOpts struct {
	// the ID of the VM to which the SSH key pair is to be bound
	ID string `json:"id" required:"true"`
	// server authentication object, this parameter is required when the server is poweron
	Auth *AuthOpts `json:"auth,omitempty"`
	// whether to disable SSH login on the VM
	DisablePassword *bool `json:"disable_password,omitempty"`
}

// AuthOpts is the object about server authentication
type AuthOpts struct {
	// the Authentication type, the value can be password or keypair
	Type string `json:"type,omitempty"`
	// If type is set to password, this parameter indicates the password.
	// If type is set to keypair, this parameter indicates the private key.
	Key string `json:"key,omitempty"`
}

// Associate is used to bind an SSH key pair to a specified VM
func Associate(c *golangsdk.ServiceClient, opts AssociateOpts) (string, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return "", err
	}

	var r golangsdk.Result
	_, err = c.Post(associateURL(c), b, &r.Body, nil)
	if err != nil {
		return "", err
	}

	var resp TaskResp
	r.ExtractInto(&resp)
	return resp.ID, nil
}

// Disassociate is used to unbind an SSH key pair to a specified VM
func Disassociate(c *golangsdk.ServiceClient, opts DisassociateOpts) (string, error) {
	b, err := golangsdk.BuildRequestBody(opts, "server")
	if err != nil {
		return "", err
	}

	var r golangsdk.Result
	_, err = c.Post(disassociateURL(c), b, &r.Body, nil)
	if err != nil {
		return "", err
	}

	var resp TaskResp
	r.ExtractInto(&resp)
	return resp.ID, nil
}

// GetTask retrieves the keypair task with the provided ID. To extract the task object
// from the response, call the Extract method on the GetResult.
func GetTask(client *golangsdk.ServiceClient, taskID string) (r GetResult) {
	_, r.Err = client.Get(getTaskURL(client, taskID), &r.Body, nil)
	return
}
