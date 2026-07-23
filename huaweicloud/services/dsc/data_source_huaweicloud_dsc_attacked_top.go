package dsc

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

// @API DSC GET /v2/{project_id}/sec-ops/situation-dashboard/top-attacks-info
func DataSourceDscAttackedTop() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscAttackedTopRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"attacked_api_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"attacked_api_top": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     attackedApiTopSchema(),
			},
			"attacked_asset_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"attacked_asset_top": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     attackedAssetTopSchema(),
			},
		},
	}
}

func attackedApiTopSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"api_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"application_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attacked_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func attackedAssetTopSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"attacked_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"db_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDscAttackedTopRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/sec-ops/situation-dashboard/top-attacks-info"
		product = "dsc"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC attacked top: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("attacked_api_num", utils.PathSearch("attacked_api_num", respBody, nil)),
		d.Set("attacked_api_top", flattenAttackedApiTop(
			utils.PathSearch("attacked_api_top", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("attacked_asset_num", utils.PathSearch("attacked_asset_num", respBody, nil)),
		d.Set("attacked_asset_top", flattenAttackedAssetTop(
			utils.PathSearch("attacked_asset_top", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAttackedApiTop(items []interface{}) []interface{} {
	if len(items) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(items))
	for _, v := range items {
		rst = append(rst, map[string]interface{}{
			"api_name":         utils.PathSearch("api_name", v, nil),
			"application_type": utils.PathSearch("application_type", v, nil),
			"attacked_num":     utils.PathSearch("attacked_num", v, nil),
			"instance_id":      utils.PathSearch("instance_id", v, nil),
		})
	}

	return rst
}

func flattenAttackedAssetTop(items []interface{}) []interface{} {
	if len(items) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(items))
	for _, v := range items {
		rst = append(rst, map[string]interface{}{
			"attacked_num":  utils.PathSearch("attacked_num", v, nil),
			"db_ip":         utils.PathSearch("db_ip", v, nil),
			"db_name":       utils.PathSearch("db_name", v, nil),
			"db_type":       utils.PathSearch("db_type", v, nil),
			"instance_id":   utils.PathSearch("instance_id", v, nil),
			"instance_name": utils.PathSearch("instance_name", v, nil),
		})
	}

	return rst
}
