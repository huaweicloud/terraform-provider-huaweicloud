package dds

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

// @API DDS GET /v3/{project_id}/quotas
func DataSourceDdsQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDdsQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"quotas": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the quotas information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"used": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of the used instances.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the quota resource type.`,
						},
						"mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the instance type.`,
						},
						"quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the existing quota.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceDdsQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	httpUrl := "v3/{project_id}/quotas"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DDS quotas: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.Errorf("error flattening response: %s", err)
	}

	resources := utils.PathSearch("quotas.resources", getRespBody, make([]interface{}, 0)).([]interface{})
	rst := make([]map[string]interface{}, 0, len(resources))
	for _, resource := range resources {
		rst = append(rst, map[string]interface{}{
			"type":  utils.PathSearch("type", resource, nil),
			"mode":  utils.PathSearch("mode", resource, nil),
			"quota": utils.PathSearch("quota", resource, nil),
			"used":  utils.PathSearch("used", resource, nil),
		})
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("quotas", rst),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
