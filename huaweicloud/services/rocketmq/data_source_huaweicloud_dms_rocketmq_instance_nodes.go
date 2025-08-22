package rocketmq

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

// @API DMS GET /v2/{project_id}/{engine}/instances/{instance_id}/nodes
func DataSourceDmsRocketMQInstanceNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsRocketMQInstanceNodesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the instance ID.`,
			},

			// Attributes
			"nodes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of nodes.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The node ID.`,
						},
						"broker_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The broker name.`,
						},
						"broker_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The broker ID.`,
						},
						"address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The private address.`,
						},
						"public_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The public address.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceDmsRocketMQInstanceNodesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)

	httpUrl := "v2/{project_id}/rocketmq/instances/{instance_id}/nodes"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{instance_id}", instanceId)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", path, &opts)
	if err != nil {
		return diag.Errorf("error querying DMS RocketMQ instance nodes: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("nodes", flattenInstanceNodes(utils.PathSearch("nodes", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenInstanceNodes(nodes []interface{}) []map[string]interface{} {
	if len(nodes) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(nodes))
	for _, node := range nodes {
		result = append(result, map[string]interface{}{
			"id":             utils.PathSearch("id", node, nil),
			"broker_name":    utils.PathSearch("broker_name", node, nil),
			"broker_id":      utils.PathSearch("broker_id", node, nil),
			"address":        utils.PathSearch("address", node, nil),
			"public_address": utils.PathSearch("public_address", node, nil),
		})
	}

	return result
}
