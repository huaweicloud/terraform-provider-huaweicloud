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

// @API GeminiDB GET /v3/{project_id}/instances/{instance_id}/sessions
func DataSourceInstanceSessions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceSessionsRead,

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
			"node_sessions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}

func buildInstanceSessionsQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("node_id"); ok {
		queryParams = fmt.Sprintf("%s&node_id=%v", queryParams, v)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func dataSourceInstanceSessionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v3/{project_id}/instances/{instance_id}/sessions"
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath += buildInstanceSessionsQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving the instance sessions: %s", err)
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
		d.Set("node_sessions", flattenInstanceSessions(
			utils.PathSearch("node_sessions", getRespBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenInstanceSessions(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"node_id": utils.PathSearch("node_id", v, nil),
			"sessions": flattenINodeSessions(
				utils.PathSearch("sessions", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenINodeSessions(resp []interface{}) []interface{} {
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
