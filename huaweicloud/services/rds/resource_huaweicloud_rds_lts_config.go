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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var ltsConfigNonUpdatableParams = []string{"instance_id", "engine", "log_type"}

// @API RDS POST /v3/{project_id}/{engine}/instances/logs/lts-configs
// @API RDS GET /v3/{project_id}/{engine}/instances/logs/lts-configs
// @API RDS DELETE /v3/{project_id}/{engine}/instances/logs/lts-configs
func ResourceRdsLtsConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsLtsConfigCreateOrUpdate,
		ReadContext:   resourceRdsLtsConfigRead,
		UpdateContext: resourceRdsLtsConfigCreateOrUpdate,
		DeleteContext: resourceRdsLtsConfigDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceRdsLtsConfigImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(ltsConfigNonUpdatableParams),

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
				Description: `Specifies the ID of the RDS instance.`,
			},
			"engine": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the engine of the RDS instance.`,
			},
			"log_type": {
				Type:        schema.TypeString,
				Required:    true,
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
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildRdsLtsConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
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

func resourceRdsLtsConfigCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	engine := d.Get("engine").(string)
	instanceID := d.Get("instance_id").(string)
	logType := d.Get("log_type").(string)

	var (
		rdsLtsConfigHttpUrl = "v3/{project_id}/{engine}/instances/logs/lts-configs"
		rdsLtsConfigProduct = "rds"
	)
	rdsLtsConfigClient, err := cfg.NewServiceClient(rdsLtsConfigProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	rdsLtsConfigPath := rdsLtsConfigClient.Endpoint + rdsLtsConfigHttpUrl
	rdsLtsConfigPath = strings.ReplaceAll(rdsLtsConfigPath, "{project_id}", rdsLtsConfigClient.ProjectID)
	rdsLtsConfigPath = strings.ReplaceAll(rdsLtsConfigPath, "{engine}", engine)

	rdsLtsConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	rdsLtsConfigOpt.JSONBody = utils.RemoveNil(buildRdsLtsConfigBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		_, err = rdsLtsConfigClient.Request("POST", rdsLtsConfigPath, &rdsLtsConfigOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(rdsLtsConfigClient, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error associating RDS with LTS log: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", instanceID, logType))

	return resourceRdsLtsConfigRead(ctx, d, meta)
}

func resourceRdsLtsConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error
	var (
		getRdsLtsConfigHttpUrl = "v3/{project_id}/{engine}/instances/logs/lts-configs"
		getRdsLtsConfigProduct = "rds"
	)

	getRdsLtsConfigClient, err := cfg.NewServiceClient(getRdsLtsConfigProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	getRdsLtsConfigPath := getRdsLtsConfigClient.Endpoint + getRdsLtsConfigHttpUrl
	getRdsLtsConfigPath = strings.ReplaceAll(getRdsLtsConfigPath, "{project_id}", getRdsLtsConfigClient.ProjectID)
	getRdsLtsConfigPath = strings.ReplaceAll(getRdsLtsConfigPath, "{engine}", d.Get("engine").(string))
	getRdsLtsConfigPath += fmt.Sprintf("?instance_id=%s", d.Get("instance_id").(string))

	getRdsLtsConfigResp, err := pagination.ListAllItems(
		getRdsLtsConfigClient,
		"offset",
		getRdsLtsConfigPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving RDS LTS configs: %s", err)
	}
	getRdsLtsConfigRespJson, err := json.Marshal(getRdsLtsConfigResp)
	if err != nil {
		return diag.Errorf("error marshaling RDS LTS configs: %s", err)
	}
	var getRdsLtsConfigRespBody interface{}
	err = json.Unmarshal(getRdsLtsConfigRespJson, &getRdsLtsConfigRespBody)
	if err != nil {
		return diag.Errorf("error unmarshaling RDS LTS configs: %s", err)
	}
	jsonPath := fmt.Sprintf("instance_lts_configs[0].lts_configs[?log_type=='%s']|[0]", d.Get("log_type").(string))
	ltsConfig := utils.PathSearch(jsonPath, getRdsLtsConfigRespBody, nil)
	if !utils.PathSearch("enabled", ltsConfig, false).(bool) {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("engine", utils.PathSearch("instance_lts_configs[0].instance.engine_name", getRdsLtsConfigRespBody, nil)),
		d.Set("instance_id", utils.PathSearch("instance_lts_configs[0].instance.id", getRdsLtsConfigRespBody, nil)),
		d.Set("log_type", utils.PathSearch("log_type", ltsConfig, nil)),
		d.Set("lts_group_id", utils.PathSearch("lts_group_id", ltsConfig, nil)),
		d.Set("lts_stream_id", utils.PathSearch("lts_stream_id", ltsConfig, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDeleteRdsLtsConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
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

func resourceRdsLtsConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	engine := d.Get("engine").(string)
	instanceID := d.Get("instance_id").(string)

	var (
		deleteRdsLtsConfigHttpUrl = "v3/{project_id}/{engine}/instances/logs/lts-configs"
		deleteRdsLtsConfigProduct = "rds"
	)
	deleteRdsLtsConfigClient, err := cfg.NewServiceClient(deleteRdsLtsConfigProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	deleteRdsLtsConfigPath := deleteRdsLtsConfigClient.Endpoint + deleteRdsLtsConfigHttpUrl
	deleteRdsLtsConfigPath = strings.ReplaceAll(deleteRdsLtsConfigPath, "{project_id}", deleteRdsLtsConfigClient.ProjectID)
	deleteRdsLtsConfigPath = strings.ReplaceAll(deleteRdsLtsConfigPath, "{engine}", engine)

	deleteRdsLtsConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	deleteRdsLtsConfigOpt.JSONBody = utils.RemoveNil(buildDeleteRdsLtsConfigBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		_, err = deleteRdsLtsConfigClient.Request("DELETE", deleteRdsLtsConfigPath, &deleteRdsLtsConfigOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(deleteRdsLtsConfigClient, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error unassociating RDS with LTS log: %s", err)
	}

	return resourceRdsLtsConfigRead(ctx, d, meta)
}

func resourceRdsLtsConfigImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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
		d.Set("engine", strings.ToLower(utils.PathSearch("datastore.type", instance, "").(string))),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
