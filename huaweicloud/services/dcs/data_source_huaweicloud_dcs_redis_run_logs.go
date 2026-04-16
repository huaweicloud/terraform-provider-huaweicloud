package dcs

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

// @API DCS GET /v2/{project_id}/instances/{instance_id}/redislog
func DataSourceDcsRedisRunLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsRedisRunLogsRead,

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
			"log_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"file_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dcsRedisRunLogSchema(),
			},
		},
	}
}

func dcsRedisRunLogSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replication_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDcsRedisRunLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("dcs", region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	httpUrl := "v2/{project_id}/instances/{instance_id}/redislog"
	getBasePath := client.Endpoint + httpUrl
	getBasePath = strings.ReplaceAll(getBasePath, "{project_id}", client.ProjectID)
	getBasePath = strings.ReplaceAll(getBasePath, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200, 204},
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	offset := 1
	res := make([]interface{}, 0)
	for {
		getPath := getBasePath + buildGetRedisRunLogsQueryParams(d, offset)
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving DCS redis run logs: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		fileList := flattenListRedisRunLogsBody(getRespBody)
		if len(fileList) == 0 {
			break
		}

		res = append(res, fileList...)
		offset += 2
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("file_list", res),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetRedisRunLogsQueryParams(d *schema.ResourceData, offset int) string {
	res := fmt.Sprintf("?limit=2&offset=%v", offset)
	if v, ok := d.GetOk("log_type"); ok {
		res = fmt.Sprintf("%s&log_type=%v", res, v)
	}

	return res
}

func flattenListRedisRunLogsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("file_list", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":             utils.PathSearch("id", v, nil),
			"file_name":      utils.PathSearch("file_name", v, nil),
			"group_name":     utils.PathSearch("group_name", v, nil),
			"replication_ip": utils.PathSearch("replication_ip", v, nil),
			"status":         utils.PathSearch("status", v, nil),
			"time":           utils.PathSearch("time", v, nil),
			"backup_id":      utils.PathSearch("backup_id", v, nil),
		})
	}
	return res
}
