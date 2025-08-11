package secmaster

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

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/workflows/instances
func DataSourceWorkflowInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWorkflowInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"workflow_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dataclass_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"playbook_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"defence_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"trigger_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"from_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"to_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workflow": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name_en": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"dataclass": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"en_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"playbook": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"en_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"trigger_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"retry_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"defense_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dataobject_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildWorkflowInstancesQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=1000"

	if v, ok := d.GetOk("workflow_id"); ok {
		queryParams = fmt.Sprintf("%s&workflow_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("id"); ok {
		queryParams = fmt.Sprintf("%s&id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("dataclass_id"); ok {
		queryParams = fmt.Sprintf("%s&dataclass_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("playbook_id"); ok {
		queryParams = fmt.Sprintf("%s&playbook_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("defence_id"); ok {
		queryParams = fmt.Sprintf("%s&defence_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("trigger_type"); ok {
		queryParams = fmt.Sprintf("%s&trigger_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("from_date"); ok {
		queryParams = fmt.Sprintf("%s&from_date=%v", queryParams, v)
	}
	if v, ok := d.GetOk("to_date"); ok {
		queryParams = fmt.Sprintf("%s&to_date=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_key"); ok {
		queryParams = fmt.Sprintf("%s&sort_key=%v", queryParams, v)
	}
	if v, ok := d.GetOk("sort_dir"); ok {
		queryParams = fmt.Sprintf("%s&sort_dir=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceWorkflowInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/workflows/instances"
		result      = make([]interface{}, 0)
		offset      = 0
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workspace_id}", workspaceId)
	getPath += buildWorkflowInstancesQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving workflow instances: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		instances := utils.PathSearch("instances", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(instances) == 0 {
			break
		}

		result = append(result, instances...)
		offset += len(instances)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance", flattenWorkflowInstances(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenWorkflowInstances(instancesResp []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(instancesResp))
	for _, v := range instancesResp {
		rst = append(rst, map[string]interface{}{
			"id":            utils.PathSearch("id", v, nil),
			"name":          utils.PathSearch("name", v, nil),
			"workflow":      flattenInstancesWorkflow(utils.PathSearch("workflow", v, nil)),
			"dataclass":     flattenInstancesDataclass(utils.PathSearch("dataclass", v, nil)),
			"playbook":      flattenInstancesPlaybook(utils.PathSearch("playbook", v, nil)),
			"trigger_type":  utils.PathSearch("trigger_type", v, nil),
			"status":        utils.PathSearch("status", v, nil),
			"start_time":    utils.PathSearch("start_time", v, nil),
			"end_time":      utils.PathSearch("end_time", v, nil),
			"retry_count":   utils.PathSearch("retry_count", v, nil),
			"defense_id":    utils.PathSearch("defense_id", v, nil),
			"dataobject_id": utils.PathSearch("dataobject_id", v, nil),
		})
	}

	return rst
}

func flattenInstancesWorkflow(workflow interface{}) []interface{} {
	if workflow == nil {
		return nil
	}

	result := map[string]interface{}{
		"id":      utils.PathSearch("id", workflow, nil),
		"name":    utils.PathSearch("name", workflow, nil),
		"name_en": utils.PathSearch("name_en", workflow, nil),
		"version": utils.PathSearch("version", workflow, nil),
	}

	return []interface{}{result}
}

func flattenInstancesDataclass(dataclass interface{}) []interface{} {
	if dataclass == nil {
		return nil
	}

	result := map[string]interface{}{
		"id":      utils.PathSearch("id", dataclass, nil),
		"name":    utils.PathSearch("name", dataclass, nil),
		"en_name": utils.PathSearch("en_name", dataclass, nil),
	}

	return []interface{}{result}
}

func flattenInstancesPlaybook(playbook interface{}) []interface{} {
	if playbook == nil {
		return nil
	}

	result := map[string]interface{}{
		"id":      utils.PathSearch("id", playbook, nil),
		"name":    utils.PathSearch("name", playbook, nil),
		"en_name": utils.PathSearch("en_name", playbook, nil),
	}

	return []interface{}{result}
}
