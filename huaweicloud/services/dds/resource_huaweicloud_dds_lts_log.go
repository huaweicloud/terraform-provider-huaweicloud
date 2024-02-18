package dds

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDS POST /v3/{project_id}/instances/logs/lts-configs
// @API DDS GET /v3/{project_id}/instances/logs/lts-configs
// @API DDS DELETE /v3/{project_id}/instances/logs/lts-configs
func ResourceDdsLtsLog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDdsLtsLogCreate,
		ReadContext:   resourceDdsLtsLogRead,
		DeleteContext: resourceDdsLtsLogDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
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
				Description: `Specifies the ID of the DDS instance.`,
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
				ForceNew:    true,
				Description: `Specifies the ID of the LTS log group.`,
			},
			"lts_stream_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the LTS log stream.`,
			},
		},
	}
}

func buildCreateDdsLtsLogBodyParams(d *schema.ResourceData) map[string]interface{} {
	ltsConfigs := map[string]interface{}{
		"lts_configs": []map[string]interface{}{
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

func resourceDdsLtsLogCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createDdsLtsLogHttpUrl = "v3/{project_id}/instances/logs/lts-configs"
		createDdsLtsLogProduct = "dds"
	)
	createDdsLtsLogClient, err := cfg.NewServiceClient(createDdsLtsLogProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	createDdsLtsLogPath := createDdsLtsLogClient.Endpoint + createDdsLtsLogHttpUrl
	createDdsLtsLogPath = strings.ReplaceAll(createDdsLtsLogPath, "{project_id}", createDdsLtsLogClient.ProjectID)

	createDdsLtsLogOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	createDdsLtsLogOpt.JSONBody = utils.RemoveNil(buildCreateDdsLtsLogBodyParams(d))

	instanceID := d.Get("instance_id").(string)

	retryFunc := func() (interface{}, bool, error) {
		createDdsLtsLogResp, err := createDdsLtsLogClient.Request("POST", createDdsLtsLogPath, &createDdsLtsLogOpt)
		retry, err := handleMultiOperationsError(err)
		return createDdsLtsLogResp, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(createDdsLtsLogClient, instanceID),
		WaitTarget:   []string{"normal"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error associating DDS with LTS log: %s", err)
	}

	d.SetId(instanceID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"false"},
		Target:     []string{"true"},
		Refresh:    ddsLtsConfigRefreshFunc(createDdsLtsLogClient, instanceID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for associating DDS (%s) with LTS log to complete: %s ", instanceID, err)
	}

	return resourceDdsLtsLogRead(ctx, d, meta)
}

func resourceDdsLtsLogRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error
	var (
		getDdsLtsLogHttpUrl = "v3/{project_id}/instances/logs/lts-configs"
		getDdsLtsLogProduct = "dds"
	)

	getDdsLtsLogClient, err := cfg.NewServiceClient(getDdsLtsLogProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	getDdsLtsLogPath := getDdsLtsLogClient.Endpoint + getDdsLtsLogHttpUrl
	getDdsLtsLogPath = strings.ReplaceAll(getDdsLtsLogPath, "{project_id}", getDdsLtsLogClient.ProjectID)

	getDdsLtsLogResp, err := pagination.ListAllItems(
		getDdsLtsLogClient,
		"offset",
		getDdsLtsLogPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving DDS LTS configs: %s", err)
	}

	getDdsLtsLogRespJson, err := json.Marshal(getDdsLtsLogResp)
	if err != nil {
		return diag.Errorf("error marshaling DDS LTS configs: %s", err)
	}

	var getDdsLtsLogRespBody interface{}
	err = json.Unmarshal(getDdsLtsLogRespJson, &getDdsLtsLogRespBody)
	if err != nil {
		return diag.Errorf("error unmarshaling DDS LTS configs: %s", err)
	}

	jsonPath := fmt.Sprintf("instance_lts_configs[?instance.id=='%s']|[0].lts_configs|[0]", d.Id())
	ltsConfig := utils.PathSearch(jsonPath, getDdsLtsLogRespBody, nil)
	if !utils.PathSearch("enabled", ltsConfig, false).(bool) {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", d.Id()),
		d.Set("log_type", utils.PathSearch("log_type", ltsConfig, nil)),
		d.Set("lts_group_id", utils.PathSearch("lts_group_id", ltsConfig, nil)),
		d.Set("lts_stream_id", utils.PathSearch("lts_stream_id", ltsConfig, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDeleteDdsLtsLogBodyParams(d *schema.ResourceData) map[string]interface{} {
	ltsConfigs := map[string]interface{}{
		"lts_configs": []map[string]interface{}{
			{
				"instance_id": d.Get("instance_id"),
				"log_type":    d.Get("log_type"),
			},
		},
	}
	return ltsConfigs
}

func resourceDdsLtsLogDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteDdsLtsLogHttpUrl = "v3/{project_id}/instances/logs/lts-configs"
		deleteDdsLtsLogProduct = "dds"
	)
	deleteDdsLtsLogClient, err := cfg.NewServiceClient(deleteDdsLtsLogProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	deleteDdsLtsLogPath := deleteDdsLtsLogClient.Endpoint + deleteDdsLtsLogHttpUrl
	deleteDdsLtsLogPath = strings.ReplaceAll(deleteDdsLtsLogPath, "{project_id}", deleteDdsLtsLogClient.ProjectID)

	deleteDdsLtsLogOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	deleteDdsLtsLogOpt.JSONBody = utils.RemoveNil(buildDeleteDdsLtsLogBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		deleteDdsLtsLogResp, err := deleteDdsLtsLogClient.Request("DELETE", deleteDdsLtsLogPath, &deleteDdsLtsLogOpt)
		retry, err := handleMultiOperationsError(err)
		return deleteDdsLtsLogResp, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(deleteDdsLtsLogClient, d.Id()),
		WaitTarget:   []string{"normal"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error unassociating DDS with LTS log: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"true"},
		Target:     []string{"false"},
		Refresh:    ddsLtsConfigRefreshFunc(deleteDdsLtsLogClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for unassociating DDS (%s) with LTS log to complete: %s ", d.Id(), err)
	}

	return resourceDdsLtsLogRead(ctx, d, meta)
}

func ddsLtsConfigRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getDdsLtsLogHttpUrl := "v3/{project_id}/instances/logs/lts-configs"
		getDdsLtsLogPath := client.Endpoint + getDdsLtsLogHttpUrl
		getDdsLtsLogPath = strings.ReplaceAll(getDdsLtsLogPath, "{project_id}", client.ProjectID)

		getDdsLtsLogResp, err := pagination.ListAllItems(
			client,
			"offset",
			getDdsLtsLogPath,
			&pagination.QueryOpts{MarkerField: ""})
		if err != nil {
			return nil, "QUERY ERROR", err
		}

		getDdsLtsLogRespJson, err := json.Marshal(getDdsLtsLogResp)
		if err != nil {
			return nil, "PARSE ERROR", err
		}

		var getDdsLtsLogRespBody interface{}
		err = json.Unmarshal(getDdsLtsLogRespJson, &getDdsLtsLogRespBody)
		if err != nil {
			return nil, "PARSE ERROR", err
		}

		jsonPath := fmt.Sprintf("instance_lts_configs[?instance.id=='%s']|[0].lts_configs|[0]", instanceID)
		ltsConfig := utils.PathSearch(jsonPath, getDdsLtsLogRespBody, nil)
		enabled := utils.PathSearch("enabled", ltsConfig, false).(bool)
		return ltsConfig, strconv.FormatBool(enabled), nil
	}
}
