package gaussdb

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

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/slow-sql-execute-node
func DataSourceGaussDbSlowSqlNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDbSlowSqlNodesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     slowSqlNodesSchema(),
			},
		},
	}
}

func slowSqlNodesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"component_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGaussDbSlowSqlNodesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/slow-sql-execute-node"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getOpt.JSONBody = utils.RemoveNil(buildGaussDbSlowSqlNodesBodyParams(d))
	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB slow SQL nodes: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.Errorf("error flattening GaussDB slow SQL nodes response: %s", err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("nodes", flattenGaussDbSlowSqlNodesBody(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGaussDbSlowSqlNodesBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"action": d.Get("action"),
	}
	return bodyParams
}

func flattenGaussDbSlowSqlNodesBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("nodes", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"node_id":        utils.PathSearch("node_id", v, nil),
			"node_name":      utils.PathSearch("node_name", v, nil),
			"role":           utils.PathSearch("role", v, nil),
			"instance_id":    utils.PathSearch("instance_id", v, nil),
			"component_type": utils.PathSearch("component_type", v, nil),
		})
	}
	return res
}
