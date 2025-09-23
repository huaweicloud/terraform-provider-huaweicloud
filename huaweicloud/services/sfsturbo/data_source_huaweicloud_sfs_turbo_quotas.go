package sfsturbo

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

// @API SFSTurbo GET /v1/{project_id}/sfs-turbo/quotas
func DataSourceSfsTurboQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSfsTurboQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildQuotasSchema(),
			},
		},
	}
}

func buildQuotasSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     buildQuotaResourceSchema(),
			},
		},
	}
}

func buildQuotaResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"max": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"min": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"quota": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceSfsTurboQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/sfs-turbo/quotas"
	)

	client, err := cfg.NewServiceClient("sfs-turbo", region)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving SFS Turbo quotas: %s", err)
	}

	respBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("quotas", flattenQuotas(utils.PathSearch("quotas", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenQuotas(respBody interface{}) []map[string]interface{} {
	if respBody == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"resources": flattenQuotaResources(utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenQuotaResources(respArray []interface{}) []map[string]interface{} {
	if len(respArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(respArray))
	for _, respBody := range respArray {
		result = append(result, map[string]interface{}{
			"max":   utils.PathSearch("max", respBody, nil),
			"min":   utils.PathSearch("min", respBody, nil),
			"quota": utils.PathSearch("quota", respBody, nil),
			"type":  utils.PathSearch("type", respBody, nil),
			"unit":  utils.PathSearch("unit", respBody, nil),
			"used":  utils.PathSearch("used", respBody, nil),
		})
	}

	return result
}
