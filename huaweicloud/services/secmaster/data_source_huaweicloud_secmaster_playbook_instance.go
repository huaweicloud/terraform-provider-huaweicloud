package secmaster

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/instances/{instance_id}
func DataSourcePlaybookInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePlaybookInstanceRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the workspace ID.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the instance ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The playbook instance name.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The project ID.",
			},
			"playbook": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The playbook information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The playbook ID.",
						},
						"version_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The playbook version ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version.",
						},
					},
				},
			},
			"dataclass": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The data class information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The data class ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The data class name.",
						},
					},
				},
			},
			"dataobject": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The data object key field information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name.",
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The playbook instance status.",
			},
			"trigger_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The trigger type.",
			},
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The update time.",
			},
		},
	}
}

func dataSourcePlaybookInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "secmaster"
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/instances/{instance_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", d.Get("instance_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving SecMaster playbook instance: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceId := utils.PathSearch("id", respBody, "").(string)

	d.SetId(instanceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("project_id", utils.PathSearch("project_id", respBody, nil)),
		d.Set("playbook", flattenPlaybookInstancePlaybook(respBody)),
		d.Set("dataclass", flattenPlaybookInstanceDataclass(respBody)),
		d.Set("dataobject", flattenPlaybookInstanceDataObject(respBody)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("trigger_type", utils.PathSearch("trigger_type", respBody, nil)),
		d.Set("start_time", utils.PathSearch("start_time", respBody, nil)),
		d.Set("end_time", utils.PathSearch("end_time", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPlaybookInstancePlaybook(respBody interface{}) []interface{} {
	playbookResp := utils.PathSearch("playbook", respBody, nil)
	if playbookResp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":         utils.PathSearch("id", playbookResp, nil),
			"version_id": utils.PathSearch("version_id", playbookResp, nil),
			"name":       utils.PathSearch("name", playbookResp, nil),
			"version":    utils.PathSearch("version", playbookResp, nil),
		},
	}
}

func flattenPlaybookInstanceDataclass(respBody interface{}) []interface{} {
	dataclassResp := utils.PathSearch("dataclass", respBody, nil)
	if dataclassResp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":   utils.PathSearch("id", dataclassResp, nil),
			"name": utils.PathSearch("name", dataclassResp, nil),
		},
	}
}

func flattenPlaybookInstanceDataObject(respBody interface{}) []interface{} {
	dataObjectResp := utils.PathSearch("dataobject", respBody, nil)
	if dataObjectResp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":   utils.PathSearch("id", dataObjectResp, nil),
			"name": utils.PathSearch("name", dataObjectResp, nil),
		},
	}
}
