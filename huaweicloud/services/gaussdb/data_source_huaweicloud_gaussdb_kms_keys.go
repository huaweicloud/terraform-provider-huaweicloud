package gaussdb

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB GET /v3/{project_id}/kms/list-keys/{kms_project_name}
func DataSourceGaussDBKmsKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDBKmsKeysRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"kms_project_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDBKmsKeysKeyDetailsSchema(),
			},
			"authorized_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authorized_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func gaussDBKmsKeysKeyDetailsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_key_flag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_spec": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGaussDBKmsKeysRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/kms/list-keys/{kms_project_name}"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{kms_project_name}", d.Get("kms_project_name").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB KMS keys: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("key_details", flattenGaussDBKmsKeysKeyDetails(getRespBody)),
		d.Set("authorized_id", utils.PathSearch("authorized_id", getRespBody, nil)),
		d.Set("authorized_name", utils.PathSearch("authorized_name", getRespBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGaussDBKmsKeysKeyDetails(resp interface{}) []interface{} {
	curJson := utils.PathSearch("key_details", resp, make([]interface{}, 0))
	curArray, ok := curJson.([]interface{})
	if !ok {
		return nil
	}

	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"key_id":           utils.PathSearch("key_id", v, nil),
			"default_key_flag": utils.PathSearch("default_key_flag", v, nil),
			"key_alias":        utils.PathSearch("key_alias", v, nil),
			"key_spec":         utils.PathSearch("key_spec", v, nil),
			"domain_id":        utils.PathSearch("domain_id", v, nil),
			"key_state":        utils.PathSearch("key_state", v, nil),
		})
	}
	return res
}
