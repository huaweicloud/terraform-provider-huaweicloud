package deh

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DEH GET /v1.0/{project_id}/quota-sets/{tenant_id}
func DataSourceDehQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDehQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"resource": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Indicates the quota resource type.`,
			},
			"quota_set": {
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
						"resource": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the quota resource type.`,
						},
						"hard_limit": {
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

func dataSourceDehQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("deh", region)
	if err != nil {
		return diag.Errorf("error creating DEH client: %s", err)
	}

	httpUrl := "v1.0/{project_id}/quota-sets/{tenant_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{tenant_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	if v, ok := d.GetOk("resource"); ok {
		getPath += fmt.Sprintf("?resource=%s", v)
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DEH quotas: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.Errorf("error flattening response: %s", err)
	}

	resources := utils.PathSearch("quota_set", getRespBody, make([]interface{}, 0)).([]interface{})
	rst := make([]map[string]interface{}, 0, len(resources))
	for _, resource := range resources {
		rst = append(rst, map[string]interface{}{
			"resource":   utils.PathSearch("resource", resource, nil),
			"hard_limit": utils.PathSearch("hard_limit", resource, nil),
			"used":       utils.PathSearch("used", resource, nil),
		})
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("quota_set", rst),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
