package dws

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DWS GET /v2/{project_id}/clusters/{cluster_id}/nodes
func DataSourceDwsClusterNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDwsClusterNodesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the DWS cluster ID.`,
			},
			"node_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the node.`,
			},
			"filter_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the query filter criteria.`,
			},
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type corresponding to the ` + "`" + `filter_by` + "`" + ` parameter.`,
			},
			"nodes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All nodes that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the node.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the node.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The current status of the node.`,
						},
						"sub_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The sub-status of the node.`,
						},
						"spec": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The specification of the node.`,
						},
						"inst_create_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The occupancy status of nodes by the cluster.`,
						},
						"alias_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The alias of the node.`,
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The availability zone of the node.`,
						},
					},
				},
			},
		},
	}
}

func queryClusterNodes(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl   = "v2/{project_id}/clusters/{cluster_id}/nodes?offset={offset}"
		clusterId = d.Get("cluster_id").(string)
		offset    = 0
		result    = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{cluster_id}", clusterId)
	queryParams := buildListNodesParams(d)
	listPath += queryParams

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		listPathWithOffset := strings.ReplaceAll(listPath, "{offset}", strconv.Itoa(offset))
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving nodes under specified DWS cluster (%s): %s", clusterId, err)
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		nodes := utils.PathSearch("node_list", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, nodes...)
		if len(result) == int(utils.PathSearch("count", respBody, float64(0)).(float64)) {
			break
		}
		offset += len(nodes)
	}
	return result, nil
}

func buildListNodesParams(d *schema.ResourceData) string {
	res := ""
	if nodeId, ok := d.GetOk("node_id"); ok {
		res = fmt.Sprintf("%s&node_ids=%v", res, nodeId)
	}
	if filterBy, ok := d.GetOk("filter_by"); ok {
		res = fmt.Sprintf("%s&filter_by=%v", res, filterBy)
	}
	if filter, ok := d.GetOk("filter"); ok {
		res = fmt.Sprintf("%s&filter=%v", res, filter)
	}

	return res
}

func dataSourceDwsClusterNodesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	nodes, err := queryClusterNodes(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("nodes", flattenNodes(nodes)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenNodes(all []interface{}) []map[string]interface{} {
	if len(all) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(all))
	for i, v := range all {
		result[i] = map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"name":              utils.PathSearch("name", v, nil),
			"status":            utils.PathSearch("status", v, nil),
			"sub_status":        utils.PathSearch("sub_status", v, nil),
			"spec":              utils.PathSearch("spec", v, nil),
			"inst_create_type":  utils.PathSearch("inst_create_type", v, nil),
			"alias_name":        utils.PathSearch("alias_name", v, nil),
			"availability_zone": utils.PathSearch("az_code", v, nil),
		}
	}
	return result
}
