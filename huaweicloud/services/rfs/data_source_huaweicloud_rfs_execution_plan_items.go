package rfs

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

// @API RFS GET /v1/{project_id}/stacks/{stack_name}/execution-plans/{execution_plan_name}
func DataSourceRfsExecutionPlanItems() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRfsExecutionPlanItemsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"stack_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"execution_plan_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"stack_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"execution_plan_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"execution_plan_items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"index": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"drifted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"imported": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attributes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"previous_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildExecutionPlanItemsQueryParams(d *schema.ResourceData) string {
	rst := ""

	if v, ok := d.GetOk("stack_id"); ok {
		rst += fmt.Sprintf("&stack_id=%s", v.(string))
	}

	if v, ok := d.GetOk("execution_plan_id"); ok {
		rst += fmt.Sprintf("&execution_plan_id=%s", v.(string))
	}

	if rst != "" {
		rst = "?" + rst[1:]
	}

	return rst
}

func dataSourceRfsExecutionPlanItemsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v1/{project_id}/stacks/{stack_name}/execution-plans/{execution_plan_name}"
		stackName = d.Get("stack_name").(string)
		planName  = d.Get("execution_plan_name").(string)
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	reqUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate RFS request UUID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{stack_name}", stackName)
	requestPath = strings.ReplaceAll(requestPath, "{execution_plan_name}", planName)
	requestPath += buildExecutionPlanItemsQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Client-Request-Id": reqUUID},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving RFS execution plan items: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(reqUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("execution_plan_items", flattenExecutionPlanItems(
			utils.PathSearch("execution_plan_items", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenExecutionPlanItems(items []interface{}) []interface{} {
	if len(items) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(items))
	for _, v := range items {
		rst = append(rst, map[string]interface{}{
			"resource_type": utils.PathSearch("resource_type", v, nil),
			"resource_name": utils.PathSearch("resource_name", v, nil),
			"index":         utils.PathSearch("index", v, nil),
			"action":        utils.PathSearch("action", v, nil),
			"action_reason": utils.PathSearch("action_reason", v, nil),
			"provider_name": utils.PathSearch("provider_name", v, nil),
			"mode":          utils.PathSearch("mode", v, nil),
			"drifted":       utils.PathSearch("drifted", v, nil),
			"imported":      utils.PathSearch("imported", v, nil),
			"resource_id":   utils.PathSearch("resource_id", v, nil),
			"attributes": flattenExecutionPlanDiffAttributes(
				utils.PathSearch("attributes", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenExecutionPlanDiffAttributes(attributes []interface{}) []interface{} {
	if len(attributes) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(attributes))
	for _, v := range attributes {
		rst = append(rst, map[string]interface{}{
			"name":           utils.PathSearch("name", v, nil),
			"previous_value": utils.PathSearch("previous_value", v, nil),
			"target_value":   utils.PathSearch("target_value", v, nil),
		})
	}

	return rst
}
