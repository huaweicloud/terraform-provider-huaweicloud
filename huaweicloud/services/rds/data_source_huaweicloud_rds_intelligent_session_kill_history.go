package rds

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

// @API RDS GET /v3/{project_id}/instances/{instance_id}/ops/intelligent-kill-session/history
func DataSourceIntelligentSessionKillHistory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIntelligentSessionKillHistoryRead,

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
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"history": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     intelligentSessionKillHistoryHistorySchema(),
			},
		},
	}
}

func intelligentSessionKillHistoryHistorySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"download_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceIntelligentSessionKillHistoryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/ops/intelligent-kill-session/history"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	getBasePath := client.Endpoint + httpUrl
	getBasePath = strings.ReplaceAll(getBasePath, "{project_id}", client.ProjectID)
	getBasePath = strings.ReplaceAll(getBasePath, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	pageNum := 1
	res := make([]map[string]interface{}, 0)
	for {
		queryParams, err := buildGetIntelligentSessionKillHistoryQueryParams(d, pageNum)
		if err != nil {
			return diag.FromErr(err)
		}
		getPath := getBasePath + queryParams
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving RDS intelligent session kill history: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		history := flattenIntelligentSessionKillHistory(getRespBody)
		res = append(res, history...)
		if len(history) < 100 {
			break
		}
		pageNum++
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("history", res),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetIntelligentSessionKillHistoryQueryParams(d *schema.ResourceData, pageNum int) (string, error) {
	res := fmt.Sprintf("?page_num=%d&page_size=100", pageNum)
	if v, ok := d.GetOk("start_time"); ok {
		startTime, err := utils.FormatUTCTimeStamp(v.(string))
		if err != nil {
			return "", err
		}
		res = fmt.Sprintf("%s&start_time=%d", res, startTime*1000)
	}
	if v, ok := d.GetOk("end_time"); ok {
		endTime, err := utils.FormatUTCTimeStamp(v.(string))
		if err != nil {
			return "", err
		}
		res = fmt.Sprintf("%s&end_time=%d", res, endTime*1000)
	}
	return res, nil
}

func flattenIntelligentSessionKillHistory(resp interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("history", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]map[string]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"task_id": utils.PathSearch("task_id", v, nil),
			"start_time": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("start_time", v, float64(0)).(float64))/1000, false),
			"end_time": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("end_time", v, float64(0)).(float64))/1000, false),
			"download_link": utils.PathSearch("download_link", v, nil),
		})
	}
	return rst
}
