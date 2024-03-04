package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS POST /v3/{project_id}/{engine}/instances/logs/lts-configs
// @API RDS GET /v3/{project_id}/{engine}/instances/logs/lts-configs
// @API RDS DELETE /v3/{project_id}/{engine}/instances/logs/lts-configs
func ResourceRdsLtsLog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsLtsLogCreateOrUpdate,
		ReadContext:   resourceRdsLtsLogRead,
		UpdateContext: resourceRdsLtsLogCreateOrUpdate,
		DeleteContext: resourceRdsLtsLogDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRdsLtsLogImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the RDS instance.`,
			},
			"engine": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the engine of the RDS instance.`,
			},
			"log_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the type of the LTS log.`,
			},
			"lts_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the LTS log group.`,
			},
			"lts_stream_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the LTS log stream.`,
			},
		},
	}
}

func buildRdsLtsLogBodyParams(d *schema.ResourceData) map[string]interface{} {
	ltsConfigs := map[string]interface{}{
		"log_configs": []map[string]interface{}{
			{
				"instance_id":   d.Get("instance_id"),
				"log_type":      d.Get("log_type"),
				"lts_group_id":  d.Get("lts_group_id"),
				"lts_stream_id": d.Get("lts_stream_id"),
			},
		},
	}
	return ltsConfigs
}

func resourceRdsLtsLogCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	engine := d.Get("engine").(string)
	instanceID := d.Get("instance_id").(string)
	logType := d.Get("log_type").(string)

	var (
		rdsLtsLogHttpUrl = "v3/{project_id}/{engine}/instances/logs/lts-configs"
		rdsLtsLogProduct = "rds"
	)
	rdsLtsLogClient, err := cfg.NewServiceClient(rdsLtsLogProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	rdsLtsLogPath := rdsLtsLogClient.Endpoint + rdsLtsLogHttpUrl
	rdsLtsLogPath = strings.ReplaceAll(rdsLtsLogPath, "{project_id}", rdsLtsLogClient.ProjectID)
	rdsLtsLogPath = strings.ReplaceAll(rdsLtsLogPath, "{engine}", engine)

	rdsLtsLogOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	rdsLtsLogOpt.JSONBody = utils.RemoveNil(buildRdsLtsLogBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		_, err = rdsLtsLogClient.Request("POST", rdsLtsLogPath, &rdsLtsLogOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(rdsLtsLogClient, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error associating RDS with LTS log: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", instanceID, logType))

	return resourceRdsLtsLogRead(ctx, d, meta)
}

func resourceRdsLtsLogRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error
	var (
		getRdsLtsLogHttpUrl = "v3/{project_id}/{engine}/instances/logs/lts-configs"
		getRdsLtsLogProduct = "rds"
	)

	getRdsLtsLogClient, err := cfg.NewServiceClient(getRdsLtsLogProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	getRdsLtsLogPath := getRdsLtsLogClient.Endpoint + getRdsLtsLogHttpUrl
	getRdsLtsLogPath = strings.ReplaceAll(getRdsLtsLogPath, "{project_id}", getRdsLtsLogClient.ProjectID)
	getRdsLtsLogPath = strings.ReplaceAll(getRdsLtsLogPath, "{engine}", d.Get("engine").(string))
	getRdsLtsLogPath += fmt.Sprintf("?instance_id=%s", d.Get("instance_id").(string))

	getRdsLtsLogResp, err := pagination.ListAllItems(
		getRdsLtsLogClient,
		"offset",
		getRdsLtsLogPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving RDS LTS configs: %s", err)
	}
	getRdsLtsLogRespJson, err := json.Marshal(getRdsLtsLogResp)
	if err != nil {
		return diag.Errorf("error marshaling RDS LTS configs: %s", err)
	}
	var getRdsLtsLogRespBody interface{}
	err = json.Unmarshal(getRdsLtsLogRespJson, &getRdsLtsLogRespBody)
	if err != nil {
		return diag.Errorf("error unmarshaling RDS LTS configs: %s", err)
	}
	jsonPath := fmt.Sprintf("instance_lts_configs[0].lts_configs[?log_type=='%s']|[0]", d.Get("log_type").(string))
	ltsConfig := utils.PathSearch(jsonPath, getRdsLtsLogRespBody, nil)
	if !utils.PathSearch("enabled", ltsConfig, false).(bool) {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("engine", utils.PathSearch("instance_lts_configs[0].instance.engine_name", getRdsLtsLogRespBody, nil)),
		d.Set("instance_id", utils.PathSearch("instance_lts_configs[0].instance.id", getRdsLtsLogRespBody, nil)),
		d.Set("log_type", utils.PathSearch("log_type", ltsConfig, nil)),
		d.Set("lts_group_id", utils.PathSearch("lts_group_id", ltsConfig, nil)),
		d.Set("lts_stream_id", utils.PathSearch("lts_stream_id", ltsConfig, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDeleteRdsLtsLogBodyParams(d *schema.ResourceData) map[string]interface{} {
	ltsConfigs := map[string]interface{}{
		"log_configs": []map[string]interface{}{
			{
				"instance_id": d.Get("instance_id"),
				"log_type":    d.Get("log_type"),
			},
		},
	}
	return ltsConfigs
}

func resourceRdsLtsLogDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	engine := d.Get("engine").(string)
	instanceID := d.Get("instance_id").(string)

	var (
		deleteRdsLtsLogHttpUrl = "v3/{project_id}/{engine}/instances/logs/lts-configs"
		deleteRdsLtsLogProduct = "rds"
	)
	deleteRdsLtsLogClient, err := cfg.NewServiceClient(deleteRdsLtsLogProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	deleteRdsLtsLogPath := deleteRdsLtsLogClient.Endpoint + deleteRdsLtsLogHttpUrl
	deleteRdsLtsLogPath = strings.ReplaceAll(deleteRdsLtsLogPath, "{project_id}", deleteRdsLtsLogClient.ProjectID)
	deleteRdsLtsLogPath = strings.ReplaceAll(deleteRdsLtsLogPath, "{engine}", engine)

	deleteRdsLtsLogOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	deleteRdsLtsLogOpt.JSONBody = utils.RemoveNil(buildDeleteRdsLtsLogBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		_, err = deleteRdsLtsLogClient.Request("DELETE", deleteRdsLtsLogPath, &deleteRdsLtsLogOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(deleteRdsLtsLogClient, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error unassociating RDS with LTS log: %s", err)
	}

	return resourceRdsLtsLogRead(ctx, d, meta)
}

func resourceRdsLtsLogImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<log_type>")
	}

	cfg := meta.(*config.Config)
	client, err := cfg.RdsV3Client(cfg.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	instance, err := GetRdsInstanceByID(client, parts[0])
	if err != nil {
		return nil, fmt.Errorf("error getting RDS instance: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("log_type", parts[1]),
		d.Set("engine", strings.ToLower(instance.DataStore.Type)),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
