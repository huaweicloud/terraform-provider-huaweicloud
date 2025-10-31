package iam

import (
	"context"
	"io"
	"strconv"

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
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"metadata_file": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIdentityKeystoneMetadataFileRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	keystoneMetadataFilePath := iamClient.Endpoint + "v3-ext/auth/OS-FEDERATION/SSO/metadata?unsigned=" +
		strconv.FormatBool(d.Get("unsigned").(bool))
	options := golangsdk.RequestOpts{KeepResponseBody: true}
	response, err := iamClient.Request("GET", keystoneMetadataFilePath, &options)
	if err != nil {
		return diag.Errorf("error keystoneMetadataFile: %s", err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generate UUID: %s", err)
	}
	d.SetId(id)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("metadata_file", string(body)); err != nil {
		return diag.Errorf("error set metadata_file field: %s", err)
	}
	return nil
}
