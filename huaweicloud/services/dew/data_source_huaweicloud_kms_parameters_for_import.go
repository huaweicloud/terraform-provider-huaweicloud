package dew

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DEW POST /v1.0/{project_id}/kms/get-parameters-for-import
func DataSourceKmsParametersForImport() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKmsParametersForImportRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to obtain the KMS parameters.",
			},
			"key_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the key ID.",
			},
			"wrapping_algorithm": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the encryption algorithm of key materials.",
			},
			"sequence": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the `36` bytes sequence number of a request message.",
			},
			"import_token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The key import token.",
			},
			"expiration_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The import parameter expiration time.",
			},
			"public_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public key of the DEK material, in Base64 format.",
			},
		},
	}
}

func buildKmsParametersForImportBodyParams(d *schema.ResourceData) map[string]interface{} {
	requestBody := map[string]interface{}{
		"key_id":             d.Get("key_id"),
		"wrapping_algorithm": d.Get("wrapping_algorithm"),
		"sequence":           utils.ValueIgnoreEmpty(d.Get("sequence")),
	}

	return requestBody
}

func dataSourceKmsParametersForImportRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		getHttpUrl = "v1.0/{project_id}/kms/get-parameters-for-import"
		product    = "kms"
		mErr       *multierror.Error
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	url := client.Endpoint + getHttpUrl
	url = strings.ReplaceAll(url, "{project_id}", client.ProjectID)

	opts := golangsdk.RequestOpts{
		JSONBody:         utils.RemoveNil(buildKmsParametersForImportBodyParams(d)),
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
	}

	resp, err := client.Request("POST", url, &opts)
	if err != nil {
		return diag.Errorf("error retrieving KMS parameters for import: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("key_id", utils.PathSearch("key_id", respBody, nil)),
		d.Set("import_token", utils.PathSearch("import_token", respBody, nil)),
		d.Set("expiration_time", utils.PathSearch("expiration_time", respBody, nil)),
		d.Set("public_key", utils.PathSearch("public_key", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
