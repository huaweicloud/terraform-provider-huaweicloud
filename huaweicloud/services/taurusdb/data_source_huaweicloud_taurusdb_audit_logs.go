package taurusdb

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

// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/audit-logs
func DataSourceTaurusDBAuditLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBAuditLogsRead,

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
			"audit_logs": {
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
						"size": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"begin_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTaurusDBAuditLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/instances/{instance_id}/audit-logs"
		offset  = 0
		limit   = 50
		res     = make([]map[string]interface{}, 0)
	)
	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	getBasePath := client.Endpoint + httpUrl
	getBasePath = strings.ReplaceAll(getBasePath, "{project_id}", client.ProjectID)
	getBasePath = strings.ReplaceAll(getBasePath, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		getPath := getBasePath + buildGetAuditLogsQueryParams(d, offset, limit)
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving TaurusDB audit logs: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		auditLogs := flattenTaurusDBAuditLogs(getRespBody)
		res = append(res, auditLogs...)

		if len(auditLogs) < limit {
			break
		}
		offset += limit
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("audit_logs", res),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetAuditLogsQueryParams(d *schema.ResourceData, offset, limit int) string {
	startTime := d.Get("start_time").(string)
	endTime := d.Get("end_time").(string)
	return fmt.Sprintf("?start_time=%s&end_time=%s&offset=%d&limit=%d", startTime, endTime, offset, limit)
}

func flattenTaurusDBAuditLogs(resp interface{}) []map[string]interface{} {
	auditLogsJson := utils.PathSearch("audit_logs", resp, make([]interface{}, 0))
	auditLogsArray := auditLogsJson.([]interface{})
	result := make([]map[string]interface{}, 0, len(auditLogsArray))
	for _, auditLog := range auditLogsArray {
		result = append(result, map[string]interface{}{
			"id":         utils.PathSearch("id", auditLog, nil),
			"name":       utils.PathSearch("name", auditLog, nil),
			"size":       utils.PathSearch("size", auditLog, nil),
			"begin_time": utils.PathSearch("begin_time", auditLog, nil),
			"end_time":   utils.PathSearch("end_time", auditLog, nil),
		})
	}
	return result
}
