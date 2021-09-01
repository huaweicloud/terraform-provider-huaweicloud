package keypairs

import (
	"net/http"

	"github.com/chnsz/golangsdk"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// create request.
type CreateOptsBuilder interface {
	ToKeyPairCreateMap() (map[string]interface{}, error)
}
type CreateOpts struct {
	// Name is a friendly name to refer to this KeyPair in other services.
	Name string `json:"name" required:"true"`

	// PublicKey [optional] is a pregenerated OpenSSH-formatted public key.
	// If provided, this key will be imported and no new key will be created.
	PublicKey string `json:"public_key,omitempty"`
}

// ToSecurityGroupsCreateMap converts CreateOpts structures to map[string]interface{}
func (opts CreateOpts) ToKeyPairCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(&opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Create create key pair
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToKeyPairCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	createURL := rootURL(client)

	var resp *http.Response
	resp, r.Err = client.Post(createURL, b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK, http.StatusCreated},
	})
	if r.Err != nil {
		return
	}
	defer resp.Body.Close()

	return
}

// Get get key pair detail
func Get(client *golangsdk.ServiceClient, keyPairID string) (r GetResult) {
	getURL := GetURL(client, keyPairID)

	var resp *http.Response
	resp, r.Err = client.Get(getURL, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	if r.Err != nil {
		return
	}
	defer resp.Body.Close()

	return
}

//Delete delete the key pair
func Delete(client *golangsdk.ServiceClient, keyPairID string) (r DeleteResult) {
	deleteURL := DeleteURL(client, keyPairID)

	var resp *http.Response
	resp, r.Err = client.Delete(deleteURL, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusNoContent},
	})
	if r.Err != nil {
		return
	}
	defer resp.Body.Close()

	return
}

type ListOpts struct {
	Limit  int    `q:"limit"`
	Offset int    `q:"offset"`
	Name   string `q:"name"`
}

type ListKeyPairsOptsBuilder interface {
	ToListKeyPairsQuery() (string, error)
}

func (opts ListOpts) ToListKeyPairsQuery() (string, error) {
	b, err := golangsdk.BuildQueryString(&opts)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func List(client *golangsdk.ServiceClient, opts ListKeyPairsOptsBuilder) (r ListResult) {
	listKeyPairsURL := rootURL(client)
	if opts != nil {
		query, err := opts.ToListKeyPairsQuery()
		if err != nil {
			r.Err = err
			return r
		}
		listKeyPairsURL += query
	}

	_, r.Err = client.Get(listKeyPairsURL, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{http.StatusOK},
	})
	return
}
