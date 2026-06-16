package geminidb

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

// @API GeminiDB GET /v3/{project_id}/redis/nodes/{node_id}/session-statistics
func DataSourceNodeSessionStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNodeSessionStatisticsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"total_connection_count": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active_connection_count": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"top_source_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"top_dbs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNodeSessionStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/redis/nodes/{node_id}/session-statistics"
		nodeId  = d.Get("node_id").(string)
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{node_id}", nodeId)
	getOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving the node session statistics: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("total_connection_count", utils.PathSearch("total_connection_count", getRespBody, nil)),
		d.Set("active_connection_count", utils.PathSearch("active_connection_count", getRespBody, nil)),
		d.Set("top_source_ips", flattenTopSourceIps(
			utils.PathSearch("top_source_ips", getRespBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("top_dbs", flattenTopDbs(
			utils.PathSearch("top_dbs", getRespBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTopSourceIps(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"client_ip":        utils.PathSearch("client_ip", v, nil),
			"connection_count": utils.PathSearch("connection_count", v, nil),
		})
	}

	return result
}

func flattenTopDbs(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"db":               utils.PathSearch("db", v, nil),
			"connection_count": utils.PathSearch("connection_count", v, nil),
		})
	}

	return result
}
