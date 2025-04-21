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

// @API RDS GET /v3/{project_id}/instances/{instance_id}/hba-info/history
func DataSourcePgHbaChangeRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePgHbaChangeRecordsRead,
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
			"pg_hba_change_records": {
				Type:     schema.TypeList,
				Elem:     pgHbaChangeRecordsSchema(),
				Computed: true,
			},
		},
	}
}

func pgHbaChangeRecordsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fail_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"before_confs": {
				Type:     schema.TypeList,
				Elem:     pgHbaConfsSchema(),
				Computed: true,
			},
			"after_confs": {
				Type:     schema.TypeList,
				Elem:     pgHbaConfsSchema(),
				Computed: true,
			},
		},
	}
	return &sc
}

func pgHbaConfsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mask": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"method": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourcePgHbaChangeRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var mErr *multierror.Error
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/hba-info/history"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))
	listPath += buildListPgHbaChangeRecordsQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("pg_hba_change_records", flattenListPgHbaChangeRecordsBody(listRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListPgHbaChangeRecordsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("start_time"); ok {
		res = fmt.Sprintf("%s&start_time=%v", res, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		res = fmt.Sprintf("%s&end_time=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenListPgHbaChangeRecordsBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	pgHbaChangeRecordsArray := resp.([]interface{})

	rst := make([]interface{}, 0, len(pgHbaChangeRecordsArray))

	for _, v := range pgHbaChangeRecordsArray {
		rst = append(rst, map[string]interface{}{
			"status":       utils.PathSearch("status", v, nil),
			"time":         utils.PathSearch("time", v, nil),
			"before_confs": flattenListPgHbaChangeRecordConfsBody(utils.PathSearch("before_confs", v, nil)),
			"after_confs":  flattenListPgHbaChangeRecordConfsBody(utils.PathSearch("after_confs", v, nil)),
		})
	}
	return rst
}

func flattenListPgHbaChangeRecordConfsBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	confArray := resp.([]interface{})
	rst := make([]interface{}, 0, len(confArray))
	for _, v := range confArray {
		rst = append(rst, map[string]interface{}{
			"type":     utils.PathSearch("type", v, nil),
			"database": utils.PathSearch("database", v, nil),
			"user":     utils.PathSearch("user", v, nil),
			"address":  utils.PathSearch("address", v, nil),
			"mask":     utils.PathSearch("mask", v, nil),
			"method":   utils.PathSearch("method", v, nil),
			"priority": utils.PathSearch("priority", v, nil),
		})
	}
	return rst
}
