package dcs

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DCS GET /v2/{project_id}/instances/{instance_id}/offline/key-analysis/{task_id}/nodes
func DataSourceOfflineKeyAnalysisNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOfflineKeyAnalysisNodesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the DCS instance.`,
			},
			"task_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the task.`,
			},
			"nodes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of the offline key analysis nodes.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the node ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the node name.`,
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the group name.`,
						},
						"node_ipv6": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the node IP address.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceOfflineKeyAnalysisNodesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v2/{project_id}/instances/{instance_id}/offline/key-analysis/{task_id}/nodes"
		instanceId = d.Get("instance_id").(string)
		taskId     = d.Get("task_id").(string)
	)

	client, err := cfg.NewServiceClient("dcs", region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{task_id}", taskId)
	getOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving the DCS offline key analysis nodes: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	nodeInfos := utils.PathSearch("nodes", getRespBody, make([]interface{}, 0)).([]interface{})

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("nodes", flattenOfflineKeyAnalysisNodes(nodeInfos)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenOfflineKeyAnalysisNodes(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"id":         utils.PathSearch("id", v, nil),
			"name":       utils.PathSearch("name", v, nil),
			"group_name": utils.PathSearch("group_name", v, nil),
			"node_ipv6":  utils.PathSearch("node_ipv6", v, nil),
		})
	}

	return result
}
