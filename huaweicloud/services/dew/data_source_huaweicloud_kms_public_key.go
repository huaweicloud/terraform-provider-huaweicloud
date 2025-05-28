package dew

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DEW POST /v1.0/{project_id}/kms/get-publickey
func DataSourceKmsPublicKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKmsPublicKeyRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"key_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the key.`,
			},
			"sequence": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the request sequence number.`,
			},
			"public_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The information of the public key.`,
			},
		},
	}
}

func buildQueryPublicKeyBodyParams(d *schema.ResourceData) map[string]interface{} {
	requestBody := map[string]interface{}{
		"key_id":   d.Get("key_id"),
		"sequence": utils.ValueIgnoreEmpty(d.Get("sequence")),
	}

	return requestBody
}

func dataSourceKmsPublicKeyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		mErr       *multierror.Error
		getHttpUrl = "v1.0/{project_id}/kms/get-publickey"
		product    = "kms"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildQueryPublicKeyBodyParams(d)),
	}

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving KMS public key: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("key_id").(string))

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("key_id", utils.PathSearch("key_id", getRespBody, nil)),
		d.Set("public_key", utils.PathSearch("public_key", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
