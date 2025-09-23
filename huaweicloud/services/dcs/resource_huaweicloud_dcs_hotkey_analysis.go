package dcs

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DCS POST /v2/{project_id}/instances/{instance_id}/hotkey-task
// @API DCS GET /v2/{project_id}/instances/{instance_id}
// @API DCS GET /v2/{project_id}/instances/{instance_id}/hotkey-task/{hotkey_id}
// @API DCS DELETE /v2/{project_id}/instances/{instance_id}/hotkey-task/{hotkey_id}
func ResourceHotKeyAnalysis() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHotKeyAnalysisCreate,
		ReadContext:   resourceHotKeyAnalysisRead,
		DeleteContext: resourceHotKeyAnalysisDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDmsHotKeyAnalysisImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
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
				Description: "Specifies the ID of the DCS instance.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the hot key analysis.`,
			},
			"scan_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the mode of the hot key analysis.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time of the hot key analysis.`,
			},
			"started_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time when the hot key analysis started.`,
			},
			"finished_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time when the hot key analysis ended.`,
			},
			"num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the number of the hot key.`,
			},
			"keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the name of the hot key.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the type of the hot key.`,
						},
						"shard": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the shard where the hot key is located.`,
						},
						"db": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the database where the hot key is located.`,
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the size of the key value.`,
						},
						"unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the unit of the hot key.`,
						},
						"freq": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the access frequency of a key within a specific period of time.`,
						},
					},
				},
			},
		},
	}
}

func resourceHotKeyAnalysisCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	// createHotKeyAnalysis: create DCS hot key analysis
	var (
		createHotKeyAnalysisHttpUrl = "v2/{project_id}/instances/{instance_id}/hotkey-task"
		createHotKeyAnalysisProduct = "dcs"
	)
	createHotKeyAnalysisClient, err := cfg.NewServiceClient(createHotKeyAnalysisProduct, region)

	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createHotKeyAnalysisPath := createHotKeyAnalysisClient.Endpoint + createHotKeyAnalysisHttpUrl
	createHotKeyAnalysisPath = strings.ReplaceAll(createHotKeyAnalysisPath, "{project_id}",
		createHotKeyAnalysisClient.ProjectID)
	createHotKeyAnalysisPath = strings.ReplaceAll(createHotKeyAnalysisPath, "{instance_id}", instanceId)

	createHotKeyAnalysisOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	retryFunc := func() (interface{}, bool, error) {
		createHotKeyAnalysisResp, createErr := createHotKeyAnalysisClient.Request("POST", createHotKeyAnalysisPath,
			&createHotKeyAnalysisOpt)
		retry, err := handleOperationError(createErr)
		return createHotKeyAnalysisResp, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     refreshDcsInstanceState(createHotKeyAnalysisClient, instanceId),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return diag.Errorf("error creating DCS hot key analysis: %v", err)
	}

	hotkeyAnalysisRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	analysisId := utils.PathSearch("id", hotkeyAnalysisRespBody, "").(string)
	if analysisId == "" {
		return diag.Errorf("unable to find the analysis ID of the DCS hot key form the API response")
	}
	d.SetId(analysisId)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"waiting", "running"},
		Target:       []string{"success"},
		Refresh:      hotKeyAnalysisStatusRefreshFunc(instanceId, analysisId, createHotKeyAnalysisClient),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the hot key analysis(%s) to complete: %s", analysisId, err)
	}

	return resourceHotKeyAnalysisRead(ctx, d, meta)
}

func resourceHotKeyAnalysisRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getHotKeyAnalysis: query DCS hot key analysis
	var (
		getHotKeyAnalysisProduct = "dcs"
	)
	getHotKeyAnalysisClient, err := cfg.NewServiceClient(getHotKeyAnalysisProduct, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	getHotKeyAnalysisRes, err := getHotKeyAnalysis(getHotKeyAnalysisClient, instanceId, d.Id())

	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DCS.4941"),
			"error retrieving DCS hot key analysis")
	}

	getHotKeyAnalysisRespBody, err := utils.FlattenResponse(getHotKeyAnalysisRes)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("status", utils.PathSearch("status", getHotKeyAnalysisRespBody, nil)),
		d.Set("num", utils.PathSearch("num", getHotKeyAnalysisRespBody, nil)),
		d.Set("scan_type", utils.PathSearch("scan_type", getHotKeyAnalysisRespBody, nil)),
		d.Set("created_at", utils.PathSearch("created_at", getHotKeyAnalysisRespBody, nil)),
		d.Set("started_at", utils.PathSearch("started_at", getHotKeyAnalysisRespBody, nil)),
		d.Set("finished_at", utils.PathSearch("finished_at", getHotKeyAnalysisRespBody, nil)),
		d.Set("keys", flattenHotKeyAnalysisKeys(getHotKeyAnalysisRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceHotKeyAnalysisDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteHotKeyAnalysis: delete DCS hot key analysis
	var (
		deleteHotKeyAnalysisHttpUrl = "v2/{project_id}/instances/{instance_id}/hotkey-task/{hotkey_id}"
		deleteHotKeyAnalysisProduct = "dcs"
	)
	deleteHotKeyAnalysisClient, err := cfg.NewServiceClient(deleteHotKeyAnalysisProduct, region)
	if err != nil {
		return diag.Errorf("error creating DCS Client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	deleteHotKeyAnalysisPath := deleteHotKeyAnalysisClient.Endpoint + deleteHotKeyAnalysisHttpUrl
	deleteHotKeyAnalysisPath = strings.ReplaceAll(deleteHotKeyAnalysisPath, "{project_id}",
		deleteHotKeyAnalysisClient.ProjectID)
	deleteHotKeyAnalysisPath = strings.ReplaceAll(deleteHotKeyAnalysisPath, "{instance_id}", instanceId)
	deleteHotKeyAnalysisPath = strings.ReplaceAll(deleteHotKeyAnalysisPath, "{hotkey_id}", d.Id())

	deleteHotKeyAnalysisOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteHotKeyAnalysisClient.Request("DELETE", deleteHotKeyAnalysisPath, &deleteHotKeyAnalysisOpt)
	if err != nil {
		return diag.Errorf("error deleting DCS hot key analysis: %v", err)
	}

	return nil
}

func resourceDmsHotKeyAnalysisImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<id>")
	}

	d.SetId(parts[1])
	mErr := multierror.Append(
		d.Set("instance_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

func hotKeyAnalysisStatusRefreshFunc(instanceId, hotkeyId string, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getHotKeyAnalysisResp, err := getHotKeyAnalysis(client, instanceId, hotkeyId)
		if err != nil {
			return nil, "", err
		}
		task, err := utils.FlattenResponse(getHotKeyAnalysisResp)
		if err != nil {
			return nil, "", err
		}
		status := utils.PathSearch("status", task, "")
		return task, status.(string), nil
	}
}

func getHotKeyAnalysis(client *golangsdk.ServiceClient, instanceId string, hotkeyId string) (*http.Response, error) {
	// getHotKeyAnalysis: query DCS hot key analysis
	var (
		getHotKeyAnalysisHttpUrl = "v2/{project_id}/instances/{instance_id}/hotkey-task/{hotkey_id}"
	)
	getHotKeyAnalysisPath := client.Endpoint + getHotKeyAnalysisHttpUrl
	getHotKeyAnalysisPath = strings.ReplaceAll(getHotKeyAnalysisPath, "{project_id}", client.ProjectID)
	getHotKeyAnalysisPath = strings.ReplaceAll(getHotKeyAnalysisPath, "{instance_id}", instanceId)
	getHotKeyAnalysisPath = strings.ReplaceAll(getHotKeyAnalysisPath, "{hotkey_id}", hotkeyId)

	getHotKeyAnalysisOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getHotKeyAnalysisResp, err := client.Request("GET", getHotKeyAnalysisPath, &getHotKeyAnalysisOpt)
	return getHotKeyAnalysisResp, err
}

func flattenHotKeyAnalysisKeys(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("keys", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":  utils.PathSearch("name", v, nil),
			"type":  utils.PathSearch("type", v, nil),
			"shard": utils.PathSearch("shard", v, nil),
			"db":    utils.PathSearch("db", v, nil),
			"size":  utils.PathSearch("size", v, nil),
			"unit":  utils.PathSearch("unit", v, nil),
			"freq":  utils.PathSearch("freq", v, nil),
		})
	}
	return rst
}
