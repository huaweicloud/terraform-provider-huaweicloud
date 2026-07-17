package cci

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

// @API CCI GET /v1/{project_id}/feature-gates
func DataSourceV2FeatureGates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2FeatureGatesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"feature_gates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     featureGatesSchema(),
			},
		},
	}
}

func featureGatesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"deprecated": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"feature": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceV2FeatureGatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	httpUrl := "v1/{project_id}/feature-gates"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving CCI feature gates: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	featureGates := flattenFeatureGates(getRespBody)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("feature_gates", featureGates),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenFeatureGates(resp interface{}) []interface{} {
	curArray, ok := resp.([]interface{})
	if !ok {
		return nil
	}

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"deprecated":  utils.PathSearch("deprecated", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"feature":     utils.PathSearch("feature", v, nil),
			"type":        utils.PathSearch("type", v, nil),
			"value":       utils.PathSearch("value", v, nil),
		})
	}
	return rst
}
