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

// @API DDS GET /v3/{project_id}/nodes/{node_id}/sessions
func DataSourceNodeSessions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNodeSessionsRead,

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
			"plan_summary": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cost_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"sessions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"active": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"operation": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cost_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plan_summary": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildNodeSessionsQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=100"

	if v, ok := d.GetOk("plan_summary"); ok {
		queryParams = fmt.Sprintf("%s&plan_summary=%v", queryParams, v)
	}
	if v, ok := d.GetOk("type"); ok {
		queryParams = fmt.Sprintf("%s&type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("namespace"); ok {
		queryParams = fmt.Sprintf("%s&namespace=%v", queryParams, v)
	}
	if v, ok := d.GetOk("cost_time"); ok {
		queryParams = fmt.Sprintf("%s&cost_time=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceNodeSessionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v3/{project_id}/nodes/{node_id}/sessions"
		nodeId     = d.Get("node_id").(string)
		offset     = 0
		result     = make([]interface{}, 0)
		totalCount float64
	)

	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{node_id}", nodeId)
	getPath += buildNodeSessionsQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving the instance node sessions: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalCount = utils.PathSearch("total_count", getRespBody, float64(0)).(float64)
		dataResp := utils.PathSearch("sessions", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(dataResp) == 0 {
			break
		}

		result = append(result, dataResp...)

		offset += len(dataResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("sessions", flattenNodeSessions(result)),
		d.Set("total_count", totalCount),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenNodeSessions(sessions []interface{}) []interface{} {
	if len(sessions) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(sessions))
	for _, v := range sessions {
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("id", v, nil),
			"active":       utils.PathSearch("active", v, nil),
			"operation":    utils.PathSearch("operation", v, nil),
			"type":         utils.PathSearch("type", v, nil),
			"cost_time":    utils.PathSearch("cost_time", v, nil),
			"plan_summary": utils.PathSearch("plan_summary", v, nil),
			"host":         utils.PathSearch("host", v, nil),
			"client":       utils.PathSearch("client", v, nil),
			"description":  utils.PathSearch("description", v, nil),
			"namespace":    utils.PathSearch("namespace", v, nil),
			"db":           utils.PathSearch("db", v, nil),
			"user":         utils.PathSearch("user", v, nil),
		})
	}

	return result
}
