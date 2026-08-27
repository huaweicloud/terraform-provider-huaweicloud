package rds

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

// @API RDS POST /v3.1/{project_id}/instances/{instance_id}/slow-logs/statistics
func DataSourceRdsSlowLogStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsSlowLogStatisticsRead,

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
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"database": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"slow_log_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     rdsSlowLogSchema(),
			},
		},
	}
}

func rdsSlowLogSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"count": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lock_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rows_sent": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rows_examined": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"database": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"users": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query_sample": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildRdsSlowLogStatisticsBodyParams(d *schema.ResourceData, limit, offset int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"start_time": d.Get("start_time").(string),
		"end_time":   d.Get("end_time").(string),
		"limit":      limit,
		"offset":     offset,
	}

	if v, ok := d.GetOk("type"); ok {
		bodyParams["type"] = v
	}

	if v, ok := d.GetOk("database"); ok {
		bodyParams["database"] = v
	}

	if v, ok := d.GetOk("sort"); ok {
		bodyParams["sort"] = v
	}

	if v, ok := d.GetOk("order"); ok {
		bodyParams["order"] = v
	}

	return bodyParams
}

func dataSourceRdsSlowLogStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3.1/{project_id}/instances/{instance_id}/slow-logs/statistics"
		product = "rds"
		offset  = 0
		limit   = 100
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", d.Get("instance_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		requestOpt.JSONBody = buildRdsSlowLogStatisticsBodyParams(d, limit, offset)

		resp, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving RDS slow log statistics: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		slowLogList := utils.PathSearch("slow_log_list", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, slowLogList...)

		if len(slowLogList) < limit {
			break
		}

		offset += len(slowLogList)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("slow_log_list", flattenRdsSlowLogList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRdsSlowLogList(slowLogList []interface{}) []interface{} {
	if len(slowLogList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(slowLogList))
	for _, v := range slowLogList {
		rst = append(rst, map[string]interface{}{
			"count":         utils.PathSearch("count", v, nil),
			"time":          utils.PathSearch("time", v, nil),
			"lock_time":     utils.PathSearch("lock_time", v, nil),
			"rows_sent":     utils.PathSearch("rows_sent", v, nil),
			"rows_examined": utils.PathSearch("rows_examined", v, nil),
			"database":      utils.PathSearch("database", v, nil),
			"users":         utils.PathSearch("users", v, nil),
			"query_sample":  utils.PathSearch("query_sample", v, nil),
			"client_ip":     utils.PathSearch("client_ip", v, nil),
			"type":          utils.PathSearch("type", v, nil),
		})
	}

	return rst
}
