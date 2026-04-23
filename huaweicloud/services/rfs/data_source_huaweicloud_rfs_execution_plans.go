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

// @API RFS GET /v1/{project_id}/stacks/{stack_name}/execution-plans
func DataSourceRfsExecutionPlans() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRfsExecutionPlansRead,

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
			"stack_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"execution_plans": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"stack_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stack_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"execution_plan_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"execution_plan_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"apply_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildExecutionPlansQueryParams(d *schema.ResourceData, marker string) string {
	rst := ""

	if v, ok := d.GetOk("stack_id"); ok {
		rst += fmt.Sprintf("&stack_id=%s", v.(string))
	}

	if marker != "" {
		rst += fmt.Sprintf("&marker=%s", marker)
	}

	if rst != "" {
		rst = "?" + rst[1:]
	}

	return rst
}

func dataSourceRfsExecutionPlansRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v1/{project_id}/stacks/{stack_name}/execution-plans"
		stackName  = d.Get("stack_name").(string)
		allPlans   = make([]interface{}, 0)
		nextMarker string
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{stack_name}", stackName)

	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Client-Request-Id": uuid,
			"Content-Type":      "application/json",
		},
		KeepResponseBody: true,
	}

	for {
		requestPathWithQueryParams := requestPath + buildExecutionPlansQueryParams(d, nextMarker)
		resp, err := client.Request("GET", requestPathWithQueryParams, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving RFS execution plans: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		plans := utils.PathSearch("execution_plans", respBody, make([]interface{}, 0)).([]interface{})
		if len(plans) == 0 {
			break
		}

		allPlans = append(allPlans, plans...)

		nextMarker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if nextMarker == "" {
			break
		}
	}

	d.SetId(uuid)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("execution_plans", flattenRfsExecutionPlans(allPlans)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRfsExecutionPlans(plans []interface{}) []interface{} {
	if len(plans) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(plans))
	for _, plan := range plans {
		rst = append(rst, map[string]interface{}{
			"stack_name":          utils.PathSearch("stack_name", plan, nil),
			"stack_id":            utils.PathSearch("stack_id", plan, nil),
			"execution_plan_id":   utils.PathSearch("execution_plan_id", plan, nil),
			"execution_plan_name": utils.PathSearch("execution_plan_name", plan, nil),
			"description":         utils.PathSearch("description", plan, nil),
			"status":              utils.PathSearch("status", plan, nil),
			"status_message":      utils.PathSearch("status_message", plan, nil),
			"create_time":         utils.PathSearch("create_time", plan, nil),
			"apply_time":          utils.PathSearch("apply_time", plan, nil),
		})
	}
	return rst
}
