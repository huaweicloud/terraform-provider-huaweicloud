package kafka

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Kafka GET /v2/instances/maintain-windows
func DataSourceMaintainWindows() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMaintainWindowsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the maintain windows are located.`,
			},
			"maintain_windows": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the maintain windows.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether this is the default time window.`,
						},
						"begin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The start time of the maintain window.`,
						},
						"end": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The end time of the maintain window.`,
						},
						"seq": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The sequence number.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceMaintainWindowsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/instances/maintain-windows"
	)

	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=utf-8"},
	}

	resp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return diag.Errorf("error retrieving maintain windows: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(randomId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("maintain_windows", flattenMaintainWindows(utils.PathSearch("maintain_windows",
			respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMaintainWindows(maintainWindows []interface{}) []interface{} {
	result := make([]interface{}, 0, len(maintainWindows))
	for _, window := range maintainWindows {
		result = append(result, map[string]interface{}{
			"default": utils.PathSearch("default", window, nil),
			"begin":   utils.PathSearch("begin", window, nil),
			"end":     utils.PathSearch("end", window, nil),
			"seq":     utils.PathSearch("seq", window, nil),
		})
	}

	return result
}
