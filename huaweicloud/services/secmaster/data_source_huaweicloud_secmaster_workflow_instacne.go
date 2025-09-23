package secmaster

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

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/workflows/instances/{instance_id}
func DataSourceWorkflowInstanceDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWorkflowInstanceDetailRead,

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
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"workflow_instance_id": {
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
	}
}

func dataSourceWorkflowInstanceDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		instanceId  = d.Get("instance_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/workflows/instances/{instance_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workspace_id}", workspaceId)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving workflow instance detail: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("workflow_instance_id", utils.PathSearch("id", getRespBody, nil)),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("workflow", flattenWorkflow(utils.PathSearch("workflow", getRespBody, nil))),
		d.Set("dataclass", flattenDataclass(utils.PathSearch("dataclass", getRespBody, nil))),
		d.Set("playbook", flattenPlaybook(utils.PathSearch("playbook", getRespBody, nil))),
		d.Set("trigger_type", utils.PathSearch("trigger_type", getRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
		d.Set("start_time", utils.PathSearch("start_time", getRespBody, nil)),
		d.Set("end_time", utils.PathSearch("end_time", getRespBody, nil)),
		d.Set("retry_count", utils.PathSearch("retry_count", getRespBody, nil)),
		d.Set("defense_id", utils.PathSearch("defense_id", getRespBody, nil)),
		d.Set("dataobject_id", utils.PathSearch("dataobject_id", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenWorkflow(workflow interface{}) []interface{} {
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

func flattenDataclass(dataclass interface{}) []interface{} {
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

func flattenPlaybook(playbook interface{}) []interface{} {
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
