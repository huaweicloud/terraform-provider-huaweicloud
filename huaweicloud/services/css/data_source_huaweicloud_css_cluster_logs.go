package css

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

// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/logs/search
func DataSourceCssClusterLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterLogsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"keyword": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"time_index": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"log_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     clusterLogsLogSchema(),
			},
			"completed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_log": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func clusterLogsLogSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"content": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"level": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceClusterLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v1.0/{project_id}/clusters/{cluster_id}/logs/search"
		clusterId = d.Get("cluster_id").(string)
	)

	client, err := cfg.NewServiceClient("css", region)
	if err != nil {
		return diag.Errorf("error creating CSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{cluster_id}", clusterId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}
	requestOpt.JSONBody = utils.RemoveNil(buildGetClusterLogsBodyParams(d))

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CSS cluster(%s) logs: %s", clusterId, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	logList := flattenGetClusterLogsBody(respBody)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("log_list", logList),
		d.Set("type", utils.PathSearch("type", respBody, nil)),
		d.Set("completed", utils.PathSearch("completed", respBody, nil)),
		d.Set("instance_log", utils.PathSearch("instanceLog", respBody, "")),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetClusterLogsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instance_name": d.Get("instance_name"),
		"log_type":      d.Get("log_type"),
	}
	if v, ok := d.GetOk("level"); ok {
		bodyParams["level"] = v
	}
	if v, ok := d.GetOk("limit"); ok {
		bodyParams["limit"] = v
	}
	if v, ok := d.GetOk("keyword"); ok {
		bodyParams["keyword"] = v
	}
	if v, ok := d.GetOk("time_index"); ok {
		bodyParams["time_index"] = v
	}
	return bodyParams
}

func flattenGetClusterLogsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("logList", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"content": utils.PathSearch("content", v, nil),
			"date":    utils.PathSearch("date", v, nil),
			"level":   utils.PathSearch("level", v, nil),
		})
	}
	return res
}
