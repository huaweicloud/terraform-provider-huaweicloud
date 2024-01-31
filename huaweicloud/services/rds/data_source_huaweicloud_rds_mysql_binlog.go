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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS GET /v3/{project_id}/instances/{instance_id}/binlog/clear-policy
func DataSourceRdsMysqlBinlog() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsmysqlBinlogRead,
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
			"binlog_retention_hours": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceRdsmysqlBinlogRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var mErr *multierror.Error
	var (
		getMysqlBinlogHttpUrl = "v3/{project_id}/instances/{instance_id}/binlog/clear-policy"
		getMysqlBinlogProduct = "rds"
	)
	getMysqlBinlogClient, err := cfg.NewServiceClient(getMysqlBinlogProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}
	getMysqlBinlogPath := getMysqlBinlogClient.Endpoint + getMysqlBinlogHttpUrl
	getMysqlBinlogPath = strings.ReplaceAll(getMysqlBinlogPath, "{project_id}", getMysqlBinlogClient.ProjectID)
	getMysqlBinlogPath = strings.ReplaceAll(getMysqlBinlogPath, "{instance_id}", fmt.Sprintf("%v",
		d.Get("instance_id")))
	getMysqlBinlogOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getMysqlBinlogResp, err := getMysqlBinlogClient.Request("GET", getMysqlBinlogPath, &getMysqlBinlogOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	getMysqlBinlogRespBody, err := utils.FlattenResponse(getMysqlBinlogResp)
	if err != nil {
		return diag.FromErr(err)
	}
	retentionHours := utils.PathSearch("binlog_retention_hours", getMysqlBinlogRespBody, nil)

	if retentionHours == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("binlog_retention_hours", retentionHours),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
