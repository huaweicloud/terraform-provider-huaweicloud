package geminidb

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GeminiDB GET /v3/{project_id}/redis/nodes/{node_id}/sessions
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
			"addr_prefix": {
				Type:     schema.TypeString,
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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cmd": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"age": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"idle": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"addr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fd": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sub": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"psub": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"multi": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildNodeSessionsQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("addr_prefix"); ok {
		queryParams = fmt.Sprintf("%s&addr_prefix=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceNodeSessionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/redis/nodes/{node_id}/sessions?limit=100"
		offset  = 0
		result  = make([]interface{}, 0)
		nodeId  = d.Get("node_id").(string)
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{node_id}", nodeId)
	getPath += buildNodeSessionsQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving the node sessions: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		nodeSessions := utils.PathSearch("sessions", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(nodeSessions) == 0 {
			break
		}

		result = append(result, nodeSessions...)
		offset += len(nodeSessions)
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("sessions", flattenNodeSessions(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenNodeSessions(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"id":    utils.PathSearch("id", v, nil),
			"name":  utils.PathSearch("name", v, nil),
			"cmd":   utils.PathSearch("cmd", v, nil),
			"age":   utils.PathSearch("age", v, nil),
			"idle":  utils.PathSearch("idle", v, nil),
			"db":    utils.PathSearch("db", v, nil),
			"addr":  utils.PathSearch("addr", v, nil),
			"fd":    utils.PathSearch("fd", v, nil),
			"sub":   utils.PathSearch("sub", v, nil),
			"psub":  utils.PathSearch("psub", v, nil),
			"multi": utils.PathSearch("multi", v, nil),
		})
	}

	return result
}
