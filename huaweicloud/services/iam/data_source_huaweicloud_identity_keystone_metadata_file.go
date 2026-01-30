package iam

import (
	"context"
	"io"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// DataSourceIdentityKeystoneMetadataFile
// @API IAM GET /v3-ext/auth/OS-FEDERATION/SSO/metadata
func DataSourceIdentityKeystoneMetadataFile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityKeystoneMetadataFileRead,

		Schema: map[string]*schema.Schema{
			"unsigned": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `Whether to sign the metadata according to the SAML2.0 specification.`,
			},

			// Attribute
			"metadata_file": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The keystone metadata file.`,
			},
		},
	}
}

func getKeystoneMetadataFile(client *golangsdk.ServiceClient, unsigned bool) ([]byte, error) {
	var (
		httpUrl = "v3-ext/auth/OS-FEDERATION/SSO/metadata?unsigned={unsigned}"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{unsigned}", strconv.FormatBool(unsigned))
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(requestResp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func dataSourceIdentityKeystoneMetadataFileRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	respBody, err := getKeystoneMetadataFile(iamClient, d.Get("unsigned").(bool))
	if err != nil {
		return diag.Errorf("error querying Keystone metadata file: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("metadata_file", string(respBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
