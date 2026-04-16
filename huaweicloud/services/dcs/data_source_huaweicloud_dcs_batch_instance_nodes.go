package dcs

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DCS GET /v2/{project_id}/instances-logical-nodes
func DataSourceDcsBatchInstanceNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBatchInstanceNodesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dcsBatchInstanceNodesInstanceSchema(),
			},
		},
	}
}

func dcsBatchInstanceNodesInstanceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dcsBatchInstanceNodesInstanceNodeSchema(),
			},
		},
	}
}

func dcsBatchInstanceNodesInstanceNodeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"logical_node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"az_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"priority_weight": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_access": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_remove_ip": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"replication_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dimensions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dcsBatchInstanceNodesInstanceNodeDimensionsSchema(),
			},
		},
	}
}

func dcsBatchInstanceNodesInstanceNodeDimensionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBatchInstanceNodesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("dcs", region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	httpUrl := "v2/{project_id}/instances-logical-nodes"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving DCS batch instance nodes: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("instances", flattenListBatchInstanceNodesBody(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListBatchInstanceNodesBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("instances", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"instance_id": utils.PathSearch("instance_id", v, nil),
			"node_count":  utils.PathSearch("node_count", v, nil),
			"nodes":       flattenListBatchInstanceNodesNodesBody(v),
		})
	}
	return res
}

func flattenListBatchInstanceNodesNodesBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("nodes", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"logical_node_id": utils.PathSearch("logical_node_id", v, nil),
			"name":            utils.PathSearch("name", v, nil),
			"status":          utils.PathSearch("status", v, nil),
			"az_code":         utils.PathSearch("az_code", v, nil),
			"node_role":       utils.PathSearch("node_role", v, nil),
			"node_type":       utils.PathSearch("node_type", v, nil),
			"node_ip":         utils.PathSearch("node_ip", v, nil),
			"node_port":       utils.PathSearch("node_port", v, nil),
			"node_id":         utils.PathSearch("node_id", v, nil),
			"priority_weight": utils.PathSearch("priority_weight", v, nil),
			"is_access":       utils.PathSearch("is_access", v, nil),
			"group_id":        utils.PathSearch("group_id", v, nil),
			"group_name":      utils.PathSearch("group_name", v, nil),
			"is_remove_ip":    utils.PathSearch("is_remove_ip", v, nil),
			"replication_id":  utils.PathSearch("replication_id", v, nil),
			"dimensions":      flattenListBatchInstanceNodesNodesDimensionsBody(v),
		})
	}
	return res
}

func flattenListBatchInstanceNodesNodesDimensionsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("dimensions", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"name":  utils.PathSearch("name", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}
	return res
}
