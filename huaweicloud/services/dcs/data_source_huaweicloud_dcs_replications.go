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

// @API DCS GET /v2/{project_id}/instance/{instance_id}/groups/{group_id}/group-nodes-state
func DataSourceReplications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReplicationsRead,

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
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the shard.`,
			},
			"replications": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of replications.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the ID of the replication.`,
						},
						"role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the role of the replication.`,
						},
						"replication_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the IP address of the replication.`,
						},
						"is_replication": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the replica is newly added.`,
						},
						"node_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the ID of the node.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the status of the replication.`,
						},
						"az_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the availability zone where the replication located.`,
						},
						"dimensions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the corresponding monitoring indicator dimension information of the replication.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the monitoring dimension name.`,
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the dimension value.`,
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

func dataSourceReplicationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v2/{project_id}/instance/{instance_id}/groups/{group_id}/group-nodes-state"
		instanceId = d.Get("instance_id").(string)
		groupId    = d.Get("group_id").(string)
	)

	client, err := cfg.NewServiceClient("dcs", region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{group_id}", groupId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving the DCS instance replication state information: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("replications", flattenReplications(
			utils.PathSearch("[*]", getRespBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenReplications(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"id":             utils.PathSearch("replication_id", v, nil),
			"role":           utils.PathSearch("replication_role", v, nil),
			"replication_ip": utils.PathSearch("replication_ip", v, nil),
			"is_replication": utils.PathSearch("is_replication", v, nil),
			"node_id":        utils.PathSearch("node_id", v, nil),
			"status":         utils.PathSearch("status", v, nil),
			"az_code":        utils.PathSearch("az_code", v, nil),
			"dimensions":     flattenReplicationDimensions(utils.PathSearch("dimensions", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenReplicationDimensions(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"name":  utils.PathSearch("name", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}

	return result
}
