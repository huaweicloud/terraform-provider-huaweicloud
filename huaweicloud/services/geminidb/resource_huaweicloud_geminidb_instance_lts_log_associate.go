package geminidb

import (
	"context"
	"errors"
	"fmt"
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

var geminiDbInstanceLtsLogAssociateNonUpdatableParams = []string{"instance_id", "log_type"}

// @API GeminiDB POST /v3/{project_id}/instances/logs/lts-configs
// @API GeminiDB DELETE /v3/{project_id}/instances/logs/lts-configs
// @API GeminiDB GET /v3/{project_id}/instances/logs/lts-configs
// @API GeminiDB POST /v3/{project_id}/instances
// @API GeminiDB GET /v3/{project_id}/jobs
func ResourceGeminiDBInstanceLtsLogAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGeminiDBInstanceLtsLogAssociateCreate,
		ReadContext:   resourceGeminiDBInstanceLtsLogAssociateRead,
		UpdateContext: resourceGeminiDBInstanceLtsLogAssociateUpdate,
		DeleteContext: resourceGeminiDBInstanceLtsLogAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceGeminiDBInstanceLtsLogAssociateImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(geminiDbInstanceLtsLogAssociateNonUpdatableParams),

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
			"log_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lts_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lts_stream_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
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

func resourceGeminiDBInstanceLtsLogAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	logType := d.Get("log_type").(string)
	ltsGroupID := d.Get("lts_group_id").(string)
	ltsStreamID := d.Get("lts_stream_id").(string)

	httpUrl := "v3/{project_id}/instances/logs/lts-configs"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"instance_ids":  []string{instanceID},
			"log_type":      logType,
			"lts_group_id":  ltsGroupID,
			"lts_stream_id": ltsStreamID,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}

	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     geminiDbInstanceStatusRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return diag.Errorf("error associating LTS log stream: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", instanceID, logType))

	return resourceGeminiDBInstanceLtsLogAssociateRead(ctx, d, meta)
}

func resourceGeminiDBInstanceLtsLogAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	logType := d.Get("log_type").(string)

	httpUrl := "v3/{project_id}/instances/logs/lts-configs?instance_id={instance_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving GeminiDB instance LTS log associate")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	listPath := fmt.Sprintf("instance_lts_configs[?instance.id=='%s']|[0].lts_configs", instanceId)
	ltsList := utils.PathSearch(listPath, getRespBody, []interface{}{}).([]interface{})

	filterPath := fmt.Sprintf("[?log_type=='%s']|[0]", logType)
	matchedConfig := utils.PathSearch(filterPath, ltsList, nil)
	if matchedConfig == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving GeminiDB instance LTS log associate")
	}

	enabled := utils.PathSearch("enabled", matchedConfig, false).(bool)
	if !enabled {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "GeminiDB instance LTS log associate is disabled")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("log_type", utils.PathSearch("log_type", matchedConfig, nil)),
		d.Set("lts_group_id", utils.PathSearch("lts_group_id", matchedConfig, nil)),
		d.Set("lts_stream_id", utils.PathSearch("lts_stream_id", matchedConfig, nil)),
		d.Set("enabled", enabled),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceGeminiDBInstanceLtsLogAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	if d.HasChangeExcept("enable_force_new") {
		httpUrl := "v3/{project_id}/instances/logs/lts-configs"
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
			JSONBody:         buildUpdateLtsLogAssociateBodyParams(d, instanceId),
		}

		retryFunc := func() (interface{}, bool, error) {
			res, err := client.Request("POST", updatePath, &updateOpt)
			retry, err := handleMultiOperationsError(err)
			return res, retry, err
		}

		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     geminiDbInstanceStatusRefreshFunc(client, instanceId),
			WaitTarget:   []string{"ACTIVE"},
			Timeout:      d.Timeout(schema.TimeoutCreate),
			DelayTimeout: 10 * time.Second,
			PollInterval: 10 * time.Second,
		})

		if err != nil {
			return diag.Errorf("error updating LTS log configuration: %s", err)
		}
	}

	return resourceGeminiDBInstanceLtsLogAssociateRead(ctx, d, meta)
}

func buildUpdateLtsLogAssociateBodyParams(d *schema.ResourceData, instanceId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instance_ids":  []string{instanceId},
		"log_type":      d.Get("log_type"),
		"lts_group_id":  d.Get("lts_group_id"),
		"lts_stream_id": d.Get("lts_stream_id"),
	}

	return bodyParams
}

func resourceGeminiDBInstanceLtsLogAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	logType := d.Get("log_type").(string)

	httpUrl := "v3/{project_id}/instances/logs/lts-configs"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"instance_ids": []string{instanceID},
			"log_type":     logType,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error disassociating LTS log stream: %s", err)
	}

	return nil
}

func resourceGeminiDBInstanceLtsLogAssociateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <instance_id>/<log_type>")
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("log_type", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
