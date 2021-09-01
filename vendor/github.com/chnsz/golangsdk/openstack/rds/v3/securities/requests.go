package securities

import "github.com/chnsz/golangsdk"

// SSLOpts is a struct which will be used to config the SSL (Secure Sockets Layer).
type SSLOpts struct {
	// Specifies whether to enable SSL.
	//   true: SSL is enabled.
	//   false: SSL is disabled.
	SSLEnable *bool `json:"ssl_option" required:"true"`
}

// SSLOptsBuilder is an interface which to support request body build of
// the ssl configuration of the specifies database.
type SSLOptsBuilder interface {
	ToSSLOptsMap() (map[string]interface{}, error)
}

// ToSSLOptsMap is a method which to build a request body by the SSLOpts.
func (opts SSLOpts) ToSSLOptsMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateSSL is a method to enable or disable the SSL.
func UpdateSSL(client *golangsdk.ServiceClient, instanceId string, opts SSLOptsBuilder) (r SSLUpdateResult) {
	b, err := opts.ToSSLOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(rootURL(client, instanceId, "ssl"), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})

	return
}

// PortOpts is a struct which will be used to config the secure sockets layer.
type PortOpts struct {
	// Specifies the port number.
	// The MySQL port number ranges from 1024 to 65535, excluding 12017 and 33071.
	Port int `json:"port" required:"true"`
}

// PortOptsBuilder is an interface which to support request body build of
// the port configuration of the specifies database.
type PortOptsBuilder interface {
	ToPortOptsMap() (map[string]interface{}, error)
}

// ToPortOptsMap is a method which to build a request body by the DBPortOpts.
func (opts PortOpts) ToPortOptsMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

// UpdatePort is a method to update the port of the database.
func UpdatePort(client *golangsdk.ServiceClient, instanceId string, opts PortOptsBuilder) (r commonResult) {
	b, err := opts.ToPortOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(rootURL(client, instanceId, "port"), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})

	return
}

// SecGroupOpts is a struct which will be used to update the specifies security group.
type SecGroupOpts struct {
	// Specifies the security group ID.
	SecurityGroupId string `json:"security_group_id" required:"true"`
}

// SecGroupOptsBuilder is an interface which to support request body build of the security group updation.
type SecGroupOptsBuilder interface {
	ToSecGroupOptsMap() (map[string]interface{}, error)
}

// ToSecGroupOptsMap is a method which to build a request body by the SecGroupOpts.
func (opts SecGroupOpts) ToSecGroupOptsMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateSecGroup is a method to update the security group which the database belongs.
func UpdateSecGroup(client *golangsdk.ServiceClient, instanceId string, opts SecGroupOptsBuilder) (r commonResult) {
	b, err := opts.ToSecGroupOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(rootURL(client, instanceId, "security-group"), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})

	return
}
