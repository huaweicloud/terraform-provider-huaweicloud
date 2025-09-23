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

// @API DCS POST /v2/{project_id}/instances/{instance_id}/bigkey-task
// @API DCS GET /v2/{project_id}/instances/{instance_id}
// @API DCS GET /v2/{project_id}/instances/{instance_id}/bigkey-task/{bigkey_id}
// @API DCS DELETE /v2/{project_id}/instances/{instance_id}/bigkey-task/{bigkey_id}
func ResourceBigKeyAnalysis() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigKeyAnalysisCreate,
		ReadContext:   resourceBigKeyAnalysisRead,
		DeleteContext: resourceBigKeyAnalysisDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDmsBigKeyAnalysisImportState,
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
				Description: `Indicates the status of the big key analysis.`,
			},
			"scan_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the mode of the big key analysis.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time of the big key analysis.`,
			},
			"started_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time when the big key analysis started.`,
			},
			"finished_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time when the big key analysis ended.`,
			},
			"num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the number of the big key.`,
			},
			"keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the name of the big key.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the type of the big key.`,
						},
						"shard": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates shard where the big key is located.`,
						},
						"db": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the database where the big key is located.`,
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the size of the key value.`,
						},
						"unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the unit of the big key.`,
						},
					},
				},
			},
		},
	}
}

func resourceBigKeyAnalysisCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	// createBigKeyAnalysis: create DCS big key analysis
	var (
		createBigKeyAnalysisHttpUrl = "v2/{project_id}/instances/{instance_id}/bigkey-task"
		createBigKeyAnalysisProduct = "dcs"
	)
	createBigKeyAnalysisClient, err := cfg.NewServiceClient(createBigKeyAnalysisProduct, region)

	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createBigKeyAnalysisPath := createBigKeyAnalysisClient.Endpoint + createBigKeyAnalysisHttpUrl
	createBigKeyAnalysisPath = strings.ReplaceAll(createBigKeyAnalysisPath, "{project_id}",
		createBigKeyAnalysisClient.ProjectID)
	createBigKeyAnalysisPath = strings.ReplaceAll(createBigKeyAnalysisPath, "{instance_id}", instanceId)

	createBigKeyAnalysisOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	retryFunc := func() (interface{}, bool, error) {
		createBigKeyAnalysisResp, createErr := createBigKeyAnalysisClient.Request("POST", createBigKeyAnalysisPath,
			&createBigKeyAnalysisOpt)
		retry, err := handleOperationError(createErr)
		return createBigKeyAnalysisResp, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     refreshDcsInstanceState(createBigKeyAnalysisClient, instanceId),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return diag.Errorf("error creating DCS big key analysis: %v", err)
	}

	bigKeyAnalysisRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	analysisId := utils.PathSearch("id", bigKeyAnalysisRespBody, "").(string)
	if analysisId == "" {
		return diag.Errorf("unable to find the analysis ID of the DCS big key from the API response")
	}
	d.SetId(analysisId)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"waiting", "running"},
		Target:       []string{"success"},
		Refresh:      bigKeyAnalysisStatusRefreshFunc(instanceId, analysisId, createBigKeyAnalysisClient),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the big key analysis(%s) to complete: %s", analysisId, err)
	}

	return resourceBigKeyAnalysisRead(ctx, d, meta)
}

func resourceBigKeyAnalysisRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getBigKeyAnalysis: query DCS big key analysis
	var (
		getBigKeyAnalysisProduct = "dcs"
	)
	getBigKeyAnalysisClient, err := cfg.NewServiceClient(getBigKeyAnalysisProduct, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	getBigKeyAnalysisResp, err := getBigKeyAnalysis(getBigKeyAnalysisClient, instanceId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DCS.4942"),
			"error retrieving DCS big key analysis")
	}

	getBigKeyAnalysisRespBody, err := utils.FlattenResponse(getBigKeyAnalysisResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("status", utils.PathSearch("status", getBigKeyAnalysisRespBody, nil)),
		d.Set("num", utils.PathSearch("num", getBigKeyAnalysisRespBody, nil)),
		d.Set("scan_type", utils.PathSearch("scan_type", getBigKeyAnalysisRespBody, nil)),
		d.Set("created_at", utils.PathSearch("created_at", getBigKeyAnalysisRespBody, nil)),
		d.Set("started_at", utils.PathSearch("started_at", getBigKeyAnalysisRespBody, nil)),
		d.Set("finished_at", utils.PathSearch("finished_at", getBigKeyAnalysisRespBody, nil)),
		d.Set("keys", flattenBigKeyAnalysisKeys(getBigKeyAnalysisRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceBigKeyAnalysisDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteBigKeyAnalysis: delete DCS big key analysis
	var (
		deleteBigKeyAnalysisHttpUrl = "v2/{project_id}/instances/{instance_id}/bigkey-task/{bigkey_id}"
		deleteBigKeyAnalysisProduct = "dcs"
	)
	deleteBigKeyAnalysisClient, err := cfg.NewServiceClient(deleteBigKeyAnalysisProduct, region)
	if err != nil {
		return diag.Errorf("error creating DCS Client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	deleteBigKeyAnalysisPath := deleteBigKeyAnalysisClient.Endpoint + deleteBigKeyAnalysisHttpUrl
	deleteBigKeyAnalysisPath = strings.ReplaceAll(deleteBigKeyAnalysisPath, "{project_id}",
		deleteBigKeyAnalysisClient.ProjectID)
	deleteBigKeyAnalysisPath = strings.ReplaceAll(deleteBigKeyAnalysisPath, "{instance_id}", instanceId)
	deleteBigKeyAnalysisPath = strings.ReplaceAll(deleteBigKeyAnalysisPath, "{bigkey_id}", d.Id())

	deleteBigKeyAnalysisOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteBigKeyAnalysisClient.Request("DELETE", deleteBigKeyAnalysisPath, &deleteBigKeyAnalysisOpt)
	if err != nil {
		return diag.Errorf("error deleting DCS big key analysis: %v", err)
	}

	return nil
}

func resourceDmsBigKeyAnalysisImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<id>")
	}

	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

func bigKeyAnalysisStatusRefreshFunc(instanceId, bigKeyId string, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getBigKeyAnalysisResp, err := getBigKeyAnalysis(client, instanceId, bigKeyId)
		if err != nil {
			return nil, "", err
		}
		task, err := utils.FlattenResponse(getBigKeyAnalysisResp)
		if err != nil {
			return nil, "", err
		}
		status := utils.PathSearch("status", task, "")
		return task, status.(string), nil
	}
}

func getBigKeyAnalysis(client *golangsdk.ServiceClient, instanceId string, bigKeyId string) (*http.Response, error) {
	// getBigKeyAnalysis: query DCS big key analysis
	var (
		getBigKeyAnalysisHttpUrl = "v2/{project_id}/instances/{instance_id}/bigkey-task/{bigkey_id}"
	)
	getBigKeyAnalysisPath := client.Endpoint + getBigKeyAnalysisHttpUrl
	getBigKeyAnalysisPath = strings.ReplaceAll(getBigKeyAnalysisPath, "{project_id}", client.ProjectID)
	getBigKeyAnalysisPath = strings.ReplaceAll(getBigKeyAnalysisPath, "{instance_id}", instanceId)
	getBigKeyAnalysisPath = strings.ReplaceAll(getBigKeyAnalysisPath, "{bigkey_id}", bigKeyId)

	getBigKeyAnalysisOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getBigKeyAnalysisResp, err := client.Request("GET", getBigKeyAnalysisPath, &getBigKeyAnalysisOpt)
	return getBigKeyAnalysisResp, err
}

func flattenBigKeyAnalysisKeys(resp interface{}) []interface{} {
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
		})
	}
	return rst
}
