// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RDS
// ---------------------------------------------------------------

package rds

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var mysqlBinlogNonUpdatableParams = []string{"instance_id"}

// @API RDS PUT /v3/{project_id}/instances/{instance_id}/binlog/clear-policy
// @API RDS GET /v3/{project_id}/instances
func ResourceMysqlBinlog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMysqlBinlogCreateOrUpdate,
		UpdateContext: resourceMysqlBinlogCreateOrUpdate,
		ReadContext:   resourceMysqlBinlogRead,
		DeleteContext: resourceMysqlBinlogDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(mysqlBinlogNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"binlog_retention_hours": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceMysqlBinlogCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		mysqlBinlogHttpUrl = "v3/{project_id}/instances/{instance_id}/binlog/clear-policy"
		mysqlBinlogProduct = "rds"
	)
	mysqlBinlogClient, err := cfg.NewServiceClient(mysqlBinlogProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	mysqlBinlogPath := mysqlBinlogClient.Endpoint + mysqlBinlogHttpUrl
	mysqlBinlogPath = strings.ReplaceAll(mysqlBinlogPath, "{project_id}",
		mysqlBinlogClient.ProjectID)
	mysqlBinlogPath = strings.ReplaceAll(mysqlBinlogPath, "{instance_id}", instanceId)

	mysqlBinlogOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	requestBody := buildCreateMysqlBinlogBodyParams(d)
	log.Printf("[DEBUG] Create or update RDS Mysql binlog request body: %#v", requestBody)
	mysqlBinlogOpt.JSONBody = utils.RemoveNil(requestBody)

	retryFunc := func() (interface{}, bool, error) {
		_, err = mysqlBinlogClient.Request("PUT", mysqlBinlogPath, &mysqlBinlogOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(mysqlBinlogClient, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 2 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating or updating RDS Mysql binlog: %s", err)
	}

	d.SetId(instanceId)

	return resourceMysqlBinlogRead(ctx, d, meta)
}

func buildCreateMysqlBinlogBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"binlog_retention_hours": d.Get("binlog_retention_hours"),
	}
	return bodyParams
}

func resourceMysqlBinlogRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	instanceId := d.Id()
	getMysqlBinlogPath := getMysqlBinlogClient.Endpoint + getMysqlBinlogHttpUrl
	getMysqlBinlogPath = strings.ReplaceAll(getMysqlBinlogPath, "{project_id}", getMysqlBinlogClient.ProjectID)
	getMysqlBinlogPath = strings.ReplaceAll(getMysqlBinlogPath, "{instance_id}", instanceId)
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

	if retentionHours == nil || int(retentionHours.(float64)) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("binlog_retention_hours", retentionHours),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceMysqlBinlogDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		deleteMysqlBinlogHttpUrl = "v3/{project_id}/instances/{instance_id}/binlog/clear-policy"
		deleteMysqlBinlogProduct = "rds"
	)
	deleteMysqlBinlogClient, err := cfg.NewServiceClient(deleteMysqlBinlogProduct, region)
	if err != nil {
		return diag.Errorf("error deleting RDS client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	deleteMysqlBinlogPath := deleteMysqlBinlogClient.Endpoint + deleteMysqlBinlogHttpUrl
	deleteMysqlBinlogPath = strings.ReplaceAll(deleteMysqlBinlogPath, "{project_id}",
		deleteMysqlBinlogClient.ProjectID)
	deleteMysqlBinlogPath = strings.ReplaceAll(deleteMysqlBinlogPath, "{instance_id}", instanceId)
	deleteMysqlBinlogOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	requestBody := map[string]interface{}{
		"binlog_retention_hours": 0,
	}
	deleteMysqlBinlogOpt.JSONBody = requestBody
	retryFunc := func() (interface{}, bool, error) {
		_, err = deleteMysqlBinlogClient.Request("PUT", deleteMysqlBinlogPath, &deleteMysqlBinlogOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(deleteMysqlBinlogClient, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 1 * time.Second,
		PollInterval: 2 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error deleting RDS Mysql binlog: %s", err)
	}
	return nil
}
