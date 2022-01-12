package secrets

import "github.com/chnsz/golangsdk"

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

type CreateSecretOpts struct {
	Name         string `json:"name" required:"true"`
	KmsKeyID     string `json:"kms_key_id,omitempty"`
	Description  string `json:"description,omitempty"`
	SecretBinary string `json:"secret_binary,omitempty" xor:"SecretString"`
	SecretString string `json:"secret_string,omitempty" xor:"SecretBinary"`
}

func Create(c *golangsdk.ServiceClient, opts CreateSecretOpts) (*Secret, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Post(rootURL(c), b, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r SecretRst
		rst.ExtractInto(&r)
		return &r.Secret, nil
	}
	return nil, err
}

func Get(c *golangsdk.ServiceClient, secretName string) (*Secret, error) {
	var rst golangsdk.Result
	_, err := c.Get(resourceURL(c, secretName), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r SecretRst
		rst.ExtractInto(&r)
		return &r.Secret, nil
	}
	return nil, err
}

type UpdateSecretOpts struct {
	KmsKeyID    string  `json:"kms_key_id,omitempty"`
	Description *string `json:"description,omitempty"`
}

func Update(c *golangsdk.ServiceClient, secretName string, opts UpdateSecretOpts) (*Secret, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Put(resourceURL(c, secretName), b, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r SecretRst
		rst.ExtractInto(&r)
		return &r.Secret, nil
	}
	return nil, err
}

func Delete(c *golangsdk.ServiceClient, secretName string) error {
	_, err := c.Delete(resourceURL(c, secretName), &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return err
}

type CreateVersionOpts struct {
	SecretBinary  string   `json:"secret_binary,omitempty" xor:"SecretString"`
	SecretString  string   `json:"secret_string,omitempty" xor:"SecretBinary"`
	VersionStages []string `json:"version_stages,omitempty"`
}

func CreateSecretVersion(c *golangsdk.ServiceClient, secretName string, opts CreateVersionOpts) (*VersionMetadata, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst golangsdk.Result
	_, err = c.Post(versionRootURL(c, secretName), b, &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r struct {
			VersionMetadata VersionMetadata `json:"version_metadata"`
		}
		rst.ExtractInto(&r)
		return &r.VersionMetadata, nil
	}
	return nil, err
}

func ShowSecretVersion(c *golangsdk.ServiceClient, secretName string, versionId string) (*Version, error) {
	var rst golangsdk.Result
	_, err := c.Get(resourceVersionURL(c, secretName, versionId), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r struct {
			Version Version `json:"version"`
		}
		rst.ExtractInto(&r)
		return &r.Version, nil
	}
	return nil, err
}

// ListSecretVersions 查询凭据的版本列表
func ListSecretVersions(c *golangsdk.ServiceClient, secretName string) ([]VersionMetadata, error) {
	var rst golangsdk.Result
	_, err := c.Get(versionRootURL(c, secretName), &rst.Body, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	if err == nil {
		var r struct {
			VersionMetadatas []VersionMetadata `json:"version_metadatas"`
		}
		rst.ExtractInto(&r)
		return r.VersionMetadatas, nil
	}
	return nil, err
}
