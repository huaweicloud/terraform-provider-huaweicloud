package dew

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

// @API DEW GET /v1/{project_id}/dew/cpcs/vm-monitor
func DataSourceVmMonitor() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceVmMonitorRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"metric_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vsm_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"from": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"to": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datapoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"min": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"average": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"sum": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"variance": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"metric_name_output": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"average": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
		},
	}
}

func buildVmMonitorQueryParams(d *schema.ResourceData) string {
	rst := fmt.Sprintf("?namespace=%s&metric_name=%s", d.Get("namespace").(string), d.Get("metric_name").(string))

	if v, ok := d.GetOk("instance_id"); ok {
		rst += fmt.Sprintf("&instance_id=%v", v)
	}

	if v, ok := d.GetOk("vsm_id"); ok {
		rst += fmt.Sprintf("&vsm_id=%v", v)
	}

	if v, ok := d.GetOk("from"); ok {
		rst += fmt.Sprintf("&from=%v", v)
	}

	if v, ok := d.GetOk("to"); ok {
		rst += fmt.Sprintf("&to=%v", v)
	}

	if v, ok := d.GetOk("period"); ok {
		rst += fmt.Sprintf("&period=%v", v)
	}

	if v, ok := d.GetOk("filter"); ok {
		rst += fmt.Sprintf("&filter=%v", v)
	}

	return rst
}

func resourceVmMonitorRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/dew/cpcs/vm-monitor"
	)

	client, err := cfg.NewServiceClient("kms", region)
	if err != nil {
		return diag.Errorf("error creating DEW KMS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildVmMonitorQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving vm monitor: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("datapoints", flattenDatapoints(utils.PathSearch("datapoints", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("metric_name_output", utils.PathSearch("metric_name", respBody, nil)),
		d.Set("max", utils.PathSearch("max", respBody, nil)),
		d.Set("average", utils.PathSearch("average", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDatapoints(datapoints []interface{}) []map[string]interface{} {
	if len(datapoints) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(datapoints))
	for _, dp := range datapoints {
		result = append(result, map[string]interface{}{
			"max":       utils.PathSearch("max", dp, nil),
			"min":       utils.PathSearch("min", dp, nil),
			"average":   utils.PathSearch("average", dp, nil),
			"sum":       utils.PathSearch("sum", dp, nil),
			"variance":  utils.PathSearch("variance", dp, nil),
			"timestamp": utils.PathSearch("timestamp", dp, nil),
			"unit":      utils.PathSearch("unit", dp, nil),
		})
	}

	return result
}
