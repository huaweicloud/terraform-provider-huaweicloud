package dds

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

// @API DDS GET /v3/{project_id}/instances/{instance_id}/conn-statistics
func DataSourceConnectionStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectionStatisticsRead,

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
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_connections": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_inner_connections": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_outer_connections": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"inner_connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"outer_connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildConnectionStatisticsQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("node_id"); ok {
		queryParams = fmt.Sprintf("%s?node_id=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceConnectionStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v3/{project_id}/instances/{instance_id}/conn-statistics"
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath += buildConnectionStatisticsQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving the instance node connections: %s", err)
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
		d.Set("total_connections", utils.PathSearch("total_connections", getRespBody, nil)),
		d.Set("total_inner_connections", utils.PathSearch("total_inner_connections", getRespBody, nil)),
		d.Set("total_outer_connections", utils.PathSearch("total_outer_connections", getRespBody, nil)),
		d.Set("inner_connections", flattenInnerOrOUterConnections(
			utils.PathSearch("inner_connections", getRespBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("outer_connections", flattenInnerOrOUterConnections(
			utils.PathSearch("outer_connections", getRespBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenInnerOrOUterConnections(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"client_ip": utils.PathSearch("client_ip", v, nil),
			"count":     utils.PathSearch("count", v, nil),
		})
	}

	return result
}
