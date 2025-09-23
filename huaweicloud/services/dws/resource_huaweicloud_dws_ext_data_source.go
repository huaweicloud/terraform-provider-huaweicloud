// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DWS
// ---------------------------------------------------------------

package dws

import (
	"context"
	"fmt"
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

// @API DWS GET /v1.0/{project_id}/job/{job_id}
// @API DWS GET /v1.0/{project_id}/clusters/{cluster_id}/ext-data-sources
// @API DWS POST /v1.0/{project_id}/clusters/{cluster_id}/ext-data-sources
// @API DWS DELETE /v1.0/{project_id}/clusters/{cluster_id}/ext-data-sources/{ext_data_source_id}
// @API DWS PUT /v1.0/{project_id}/clusters/{cluster_id}/ext-data-sources/{ext_data_source_id}
func ResourceDwsExtDataSource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDwsExtDataSourceCreate,
		UpdateContext: resourceDwsExtDataSourceUpdate,
		ReadContext:   resourceDwsExtDataSourceRead,
		DeleteContext: resourceDwsExtDataSourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDwsExtDataSourceImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the external data source.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The type of the external data source.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The user name of the external data source.`,
			},
			"data_source_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `ID of the external data source. It is mandatory when **type** is **MRS**.`,
			},
			"user_pwd": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: `The password of the external data source. It is mandatory when **type** is **MRS**.`,
			},
			"connect_info": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The connection information of the external data source. It is mandatory when **type** is **OBS**.`,
			},
			"reboot": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to reboot the cluster.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The description of the external data source.`,
			},
			"configure_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The configure status of the external data source.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the external data source.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the external data source.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The updated time of the external data source.`,
			},
		},
	}
}

func resourceDwsExtDataSourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createDwsExtDataSource: create a DWS external data source.
	var (
		createDwsExtDataSourceHttpUrl = "v1.0/{project_id}/clusters/{cluster_id}/ext-data-sources"
		createDwsExtDataSourceProduct = "dws"
	)
	createDwsExtDataSourceClient, err := cfg.NewServiceClient(createDwsExtDataSourceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	createDwsExtDataSourcePath := createDwsExtDataSourceClient.Endpoint + createDwsExtDataSourceHttpUrl
	createDwsExtDataSourcePath = strings.ReplaceAll(createDwsExtDataSourcePath, "{project_id}", createDwsExtDataSourceClient.ProjectID)
	createDwsExtDataSourcePath = strings.ReplaceAll(createDwsExtDataSourcePath, "{cluster_id}", fmt.Sprintf("%v", d.Get("cluster_id")))

	createDwsExtDataSourceOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
	}
	createDwsExtDataSourceOpt.JSONBody = utils.RemoveNil(buildCreateDwsExtDataSourceBodyParams(d))
	createDwsExtDataSourceResp, err := createDwsExtDataSourceClient.Request("POST", createDwsExtDataSourcePath,
		&createDwsExtDataSourceOpt)
	if err != nil {
		return diag.Errorf("error creating DWS external data source: %s", err)
	}

	createDwsExtDataSourceRespBody, err := utils.FlattenResponse(createDwsExtDataSourceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId := utils.PathSearch("id", createDwsExtDataSourceRespBody, "").(string)
	if dataSourceId == "" {
		return diag.Errorf("unable to find the DWS external data source ID from the API response")
	}
	d.SetId(dataSourceId)

	jobId := utils.PathSearch("job_id", createDwsExtDataSourceRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error creating DWS external data source: job_id is not found in API response")
	}

	err = extDataSourceWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate), jobId.(string))
	if err != nil {
		return diag.Errorf("error waiting for the creation of external data source(%s) to complete: %s", d.Id(), err)
	}

	return resourceDwsExtDataSourceRead(ctx, d, meta)
}

func buildCreateDwsExtDataSourceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"data_source_name": utils.ValueIgnoreEmpty(d.Get("name")),
		"type":             utils.ValueIgnoreEmpty(d.Get("type")),
		"data_source_id":   utils.ValueIgnoreEmpty(d.Get("data_source_id")),
		"user_name":        utils.ValueIgnoreEmpty(d.Get("user_name")),
		"user_pwd":         utils.ValueIgnoreEmpty(d.Get("user_pwd")),
		"connect_info":     utils.ValueIgnoreEmpty(d.Get("connect_info")),
		"description":      utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceDwsExtDataSourceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDwsExtDataSource: Query the DWS external data source.
	extDataSource, err := GetExtDataSource(cfg, region, d, d.Get("type").(string))
	if err != nil {
		// The cluster ID does not exist.
		// "DWS.0001": The cluster ID is a non-standard UUID, the status code is 400.
		// "DWS.0047": The cluster ID is a standard UUID, the status code is 404.
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", ClusterIdIllegalErrCode),
			"error retrieving DWS external data source")
	}

	if extDataSource == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving DWS external data source")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", extDataSource, nil)),
		d.Set("type", utils.PathSearch("type", extDataSource, nil)),
		d.Set("data_source_id", utils.PathSearch("data_source_id", extDataSource, nil)),
		d.Set("user_name", utils.PathSearch("user_name", extDataSource, nil)),
		d.Set("connect_info", utils.PathSearch("connect_info", extDataSource, nil)),
		d.Set("description", utils.PathSearch("description", extDataSource, nil)),
		d.Set("configure_status", utils.PathSearch("configure_status", extDataSource, nil)),
		d.Set("status", utils.PathSearch("status", extDataSource, nil)),
		d.Set("created_at", utils.PathSearch("created", extDataSource, nil)),
		d.Set("updated_at", utils.PathSearch("data_source_updated", extDataSource, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetExtDataSource(cfg *config.Config, region string, d *schema.ResourceData, dsType string) (interface{}, error) {
	var (
		getDwsExtDataSourceHttpUrl = "v1.0/{project_id}/clusters/{cluster_id}/ext-data-sources"
		getDwsExtDataSourceProduct = "dws"
	)
	getDwsExtDataSourceClient, err := cfg.NewServiceClient(getDwsExtDataSourceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS client: %s", err)
	}

	getDwsExtDataSourcePath := getDwsExtDataSourceClient.Endpoint + getDwsExtDataSourceHttpUrl
	getDwsExtDataSourcePath = strings.ReplaceAll(getDwsExtDataSourcePath, "{project_id}", getDwsExtDataSourceClient.ProjectID)
	getDwsExtDataSourcePath = strings.ReplaceAll(getDwsExtDataSourcePath, "{cluster_id}", fmt.Sprintf("%v", d.Get("cluster_id")))

	getDwsExtDataSourcePath += fmt.Sprintf("?type=%s", dsType)

	getDwsExtDataSourceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}
	getDwsExtDataSourceResp, err := getDwsExtDataSourceClient.Request("GET", getDwsExtDataSourcePath, &getDwsExtDataSourceOpt)

	if err != nil {
		return nil, err
	}

	getDwsExtDataSourceRespBody, err := utils.FlattenResponse(getDwsExtDataSourceResp)
	if err != nil {
		return nil, err
	}

	jsonPath := fmt.Sprintf("data_sources[?id=='%s']|[0]", d.Id())
	rawData := utils.PathSearch(jsonPath, getDwsExtDataSourceRespBody, nil)

	return rawData, nil
}

func resourceDwsExtDataSourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateDwsExtDataSourceChanges := []string{
		"reboot",
		"user_name",
	}

	if d.HasChanges(updateDwsExtDataSourceChanges...) {
		// updateDwsExtDataSource: update the DWS external data source.
		var (
			updateDwsExtDataSourceHttpUrl = "v1.0/{project_id}/clusters/{cluster_id}/ext-data-sources/{ext_data_source_id}"
			updateDwsExtDataSourceProduct = "dws"
		)
		updateDwsExtDataSourceClient, err := cfg.NewServiceClient(updateDwsExtDataSourceProduct, region)
		if err != nil {
			return diag.Errorf("error creating DWS client: %s", err)
		}

		updateDwsExtDataSourcePath := updateDwsExtDataSourceClient.Endpoint + updateDwsExtDataSourceHttpUrl
		updateDwsExtDataSourcePath = strings.ReplaceAll(updateDwsExtDataSourcePath, "{project_id}", updateDwsExtDataSourceClient.ProjectID)
		updateDwsExtDataSourcePath = strings.ReplaceAll(updateDwsExtDataSourcePath, "{cluster_id}", fmt.Sprintf("%v", d.Get("cluster_id")))
		updateDwsExtDataSourcePath = strings.ReplaceAll(updateDwsExtDataSourcePath, "{ext_data_source_id}", d.Id())

		updateDwsExtDataSourceOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      requestOpts.MoreHeaders,
		}
		updateDwsExtDataSourceOpt.JSONBody = utils.RemoveNil(buildUpdateDwsExtDataSourceBodyParams(d))
		updateDwsExtDataSourceResp, err := updateDwsExtDataSourceClient.Request("PUT", updateDwsExtDataSourcePath, &updateDwsExtDataSourceOpt)
		if err != nil {
			return diag.Errorf("error updating DWS external data source: %s", err)
		}

		updateDwsExtDataSourceRespBody, err := utils.FlattenResponse(updateDwsExtDataSourceResp)
		if err != nil {
			return diag.FromErr(err)
		}

		jobId := utils.PathSearch("job_id", updateDwsExtDataSourceRespBody, nil)
		if jobId == nil {
			return diag.Errorf("error updating DWS external data source: job_id is not found in API response")
		}

		err = extDataSourceWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate), jobId.(string))
		if err != nil {
			return diag.Errorf("error waiting for the update of external data source(%s) to complete: %s", d.Id(), err)
		}
	}

	return resourceDwsExtDataSourceRead(ctx, d, meta)
}

func buildUpdateDwsExtDataSourceBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"reboot": utils.ValueIgnoreEmpty(d.Get("reboot")),
	}
	if d.Get("type").(string) == "OBS" {
		params["agency"] = utils.ValueIgnoreEmpty(d.Get("user_name"))
	}

	bodyParams := map[string]interface{}{
		"reconfigure": params,
	}
	return bodyParams
}

func resourceDwsExtDataSourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteDwsExtDataSource: delete DWS external data source
	var (
		deleteDwsExtDataSourceHttpUrl = "v1.0/{project_id}/clusters/{cluster_id}/ext-data-sources/{ext_data_source_id}"
		deleteDwsExtDataSourceProduct = "dws"
	)
	deleteDwsExtDataSourceClient, err := cfg.NewServiceClient(deleteDwsExtDataSourceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	deleteDwsExtDataSourcePath := deleteDwsExtDataSourceClient.Endpoint + deleteDwsExtDataSourceHttpUrl
	deleteDwsExtDataSourcePath = strings.ReplaceAll(deleteDwsExtDataSourcePath, "{project_id}", deleteDwsExtDataSourceClient.ProjectID)
	deleteDwsExtDataSourcePath = strings.ReplaceAll(deleteDwsExtDataSourcePath, "{ext_data_source_id}", d.Id())
	deleteDwsExtDataSourcePath = strings.ReplaceAll(deleteDwsExtDataSourcePath, "{cluster_id}", fmt.Sprintf("%v", d.Get("cluster_id")))

	deleteDwsExtDataSourceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}
	deleteDwsExtDataSourceResp, err := deleteDwsExtDataSourceClient.Request("DELETE", deleteDwsExtDataSourcePath, &deleteDwsExtDataSourceOpt)
	if err != nil {
		return diag.Errorf("error deleting DWS external data source: %s", err)
	}

	deleteDwsExtDataSourceRespBody, err := utils.FlattenResponse(deleteDwsExtDataSourceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", deleteDwsExtDataSourceRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error deleting DWS external data source: job_id is not found in API response")
	}

	err = extDataSourceWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete), jobId.(string))
	if err != nil {
		return diag.Errorf("error waiting for the delete of external data source(%s) to complete: %s", d.Id(), err)
	}

	return nil
}

func resourceDwsExtDataSourceImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <cluster_id>/<id>")
	}

	d.Set("cluster_id", parts[0])
	d.SetId(parts[1])

	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// type is null when import
	extDataSource, err := GetExtDataSource(cfg, region, d, "MRS")
	if err != nil {
		return nil, fmt.Errorf("error retrieving DWS external data source")
	}

	if extDataSource == nil {
		extDataSource, err = GetExtDataSource(cfg, region, d, "OBS")
		if err != nil {
			return nil, fmt.Errorf("error retrieving DWS external data source")
		}
	}

	rawType := utils.PathSearch("type", extDataSource, nil)
	if rawType == nil {
		return nil, fmt.Errorf("error import DWS external data source: type is not found in API response")
	}

	d.Set("type", rawType)

	return []*schema.ResourceData{d}, nil
}

func extDataSourceWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{},
	t time.Duration, jobId string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			// extDataSourceWaiting: missing operation notes
			var (
				extDataSourceWaitingHttpUrl = "v1.0/{project_id}/job/{job_id}"
				extDataSourceWaitingProduct = "dws"
			)
			extDataSourceWaitingClient, err := cfg.NewServiceClient(extDataSourceWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating DWS client: %s", err)
			}

			extDataSourceWaitingPath := extDataSourceWaitingClient.Endpoint + extDataSourceWaitingHttpUrl
			extDataSourceWaitingPath = strings.ReplaceAll(extDataSourceWaitingPath, "{project_id}", extDataSourceWaitingClient.ProjectID)
			extDataSourceWaitingPath = strings.ReplaceAll(extDataSourceWaitingPath, "{job_id}", jobId)

			extDataSourceWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				MoreHeaders:      requestOpts.MoreHeaders,
			}
			extDataSourceWaitingResp, err := extDataSourceWaitingClient.Request("GET", extDataSourceWaitingPath, &extDataSourceWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			extDataSourceWaitingRespBody, err := utils.FlattenResponse(extDataSourceWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			status := utils.PathSearch(`status`, extDataSourceWaitingRespBody, "").(string)

			targetStatus := []string{
				"SUCCESS",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return extDataSourceWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"FAIL",
				"ERROR",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return extDataSourceWaitingRespBody, status, nil
			}

			return extDataSourceWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
