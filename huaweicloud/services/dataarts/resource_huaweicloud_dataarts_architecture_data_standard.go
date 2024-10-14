// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DataArts
// ---------------------------------------------------------------

package dataarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio POST /v2/{project_id}/design/standards
// @API DataArtsStudio GET /v2/{project_id}/design/standards
// @API DataArtsStudio PUT /v2/{project_id}/design/standards/{id}
// @API DataArtsStudio GET //v2/{project_id}/design/standards/{id}
// @API DataArtsStudio DELETE /v2/{project_id}/design/standards
func ResourceDataStandard() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataStandardCreate,
		UpdateContext: resourceDataStandardUpdate,
		ReadContext:   resourceDataStandardRead,
		DeleteContext: resourceDataStandardDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDataStandardImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the workspace ID of DataArts Architecture.`,
			},
			"directory_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the directory ID that the data standard belongs to.`,
			},
			"values": {
				Type:        schema.TypeSet,
				Elem:        dataStandardValueSchema(),
				Required:    true,
				Description: `Specifies the value of data standard attributes.`,
			},
			"directory_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the path of the directory.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the data standard.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of creator.`,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of updater.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time of the data standard.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the latest update time of the data standard.`,
			},
			"new_biz": {
				Type:        schema.TypeList,
				Elem:        dataStandardNewBizSchema(),
				Computed:    true,
				Description: `Indicates the biz info of manager.`,
			},
		},
	}
}

func dataStandardValueSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"fd_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the data standard attribute.`,
			},
			"fd_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the value of the data standard attribute.`,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the data standard attribute.`,
			},
			"fd_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the data standard attribute.`,
			},
			"directory_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the directory ID that the attribute belongs to.`,
			},
			"row_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of data standard.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the data standard.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of creator.`,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of updater.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time of the data standard.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the latest update time of the data standard.`,
			},
		},
	}
	return &sc
}

func dataStandardNewBizSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the new biz.`,
			},
			"biz_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the type of the new biz.`,
			},
			"biz_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of data standard.`,
			},
			"biz_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the info of the new biz.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the new biz.`,
			},
			"biz_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the version of the new biz.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time of the new biz.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the latest update time of the new biz.`,
			},
		},
	}
	return &sc
}

func resourceDataStandardCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createDataStandard: create DataArts Architecture data standard
	var (
		createDataStandardHttpUrl = "v2/{project_id}/design/standards"
		createDataStandardProduct = "dataarts"
	)
	createDataStandardClient, err := cfg.NewServiceClient(createDataStandardProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	createDataStandardPath := createDataStandardClient.Endpoint + createDataStandardHttpUrl
	createDataStandardPath = strings.ReplaceAll(createDataStandardPath, "{project_id}", createDataStandardClient.ProjectID)

	createDataStandardOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"workspace": d.Get("workspace_id").(string),
		},
	}

	createDataStandardOpt.JSONBody = utils.RemoveNil(buildCreateDataStandardBodyParams(d))
	createDataStandardResp, err := createDataStandardClient.Request("POST", createDataStandardPath, &createDataStandardOpt)
	if err != nil {
		return diag.Errorf("error creating DataArts Architecture data standard: %s", err)
	}

	createDataStandardRespBody, err := utils.FlattenResponse(createDataStandardResp)
	if err != nil {
		return diag.FromErr(err)
	}

	standardId := utils.PathSearch("data.value.id", createDataStandardRespBody, "").(string)
	if standardId == "" {
		return diag.Errorf("unable to find the DataArts Architecture data standard ID from the API response")
	}
	d.SetId(standardId)

	return resourceDataStandardRead(ctx, d, meta)
}

func buildCreateDataStandardBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"directory_id": d.Get("directory_id"),
		"values":       buildCreateDataStandardRequestBodyValue(d.Get("values").(*schema.Set).List()),
	}
	return bodyParams
}

func buildCreateDataStandardRequestBodyValue(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		raw := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"fd_name":  raw["fd_name"],
			"fd_value": raw["fd_value"],
		}
	}
	return rst
}

func resourceDataStandardRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDataStandard: query DataArts Architecture data standard
	var (
		getDataStandardHttpUrl = "v2/{project_id}/design/standards"
		getDataStandardProduct = "dataarts"
	)
	getDataStandardClient, err := cfg.NewServiceClient(getDataStandardProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	getDataStandardBasePath := getDataStandardClient.Endpoint + getDataStandardHttpUrl
	getDataStandardBasePath = strings.ReplaceAll(getDataStandardBasePath, "{project_id}", getDataStandardClient.ProjectID)

	var currentTotal int
	var dataStandard interface{}
	for {
		getDataStandardPath := getDataStandardBasePath + buildGetDataStandardQueryParams(d, currentTotal)
		getDataStandardOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"workspace": d.Get("workspace_id").(string),
			},
		}

		getDataStandardResp, err := getDataStandardClient.Request("GET", getDataStandardPath, &getDataStandardOpt)

		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving DataStandard")
		}

		getDataStandardRespBody, err := utils.FlattenResponse(getDataStandardResp)
		if err != nil {
			return diag.FromErr(err)
		}

		records := utils.PathSearch("data.value.records", getDataStandardRespBody, make([]interface{}, 0))
		dataStandard = utils.PathSearch(fmt.Sprintf("[?id=='%s']|[0]", d.Id()), records, nil)
		if dataStandard != nil {
			break
		}
		total := utils.PathSearch("data.value.total", getDataStandardRespBody, float64(0)).(float64)
		currentTotal += len(records.([]interface{}))
		if currentTotal == int(total) {
			break
		}
	}

	if dataStandard == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("directory_id", utils.PathSearch("directory_id", dataStandard, nil)),
		d.Set("directory_path", utils.PathSearch("directory_path", dataStandard, nil)),
		d.Set("status", utils.PathSearch("status", dataStandard, nil)),
		d.Set("created_by", utils.PathSearch("create_by", dataStandard, nil)),
		d.Set("updated_by", utils.PathSearch("update_by", dataStandard, nil)),
		d.Set("created_at", utils.PathSearch("create_time", dataStandard, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", dataStandard, nil)),
		d.Set("values", flattenGetDataStandardResponseBodyValue(dataStandard)),
		d.Set("new_biz", flattenGetDataStandardResponseBodyNewBiz(dataStandard)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetDataStandardQueryParams(d *schema.ResourceData, offset int) string {
	return fmt.Sprintf("?directory_id=%v&limit=100&offset=%v", d.Get("directory_id"), offset)
}

func flattenGetDataStandardResponseBodyValue(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("values", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":           utils.PathSearch("id", v, nil),
			"fd_id":        utils.PathSearch("fd_id", v, nil),
			"fd_name":      utils.PathSearch("fd_name", v, nil),
			"fd_value":     utils.PathSearch("fd_value", v, nil),
			"directory_id": utils.PathSearch("directory_id", v, nil),
			"row_id":       utils.PathSearch("row_id", v, nil),
			"status":       utils.PathSearch("status", v, nil),
			"created_by":   utils.PathSearch("create_by", v, nil),
			"updated_by":   utils.PathSearch("update_by", v, nil),
			"created_at":   utils.PathSearch("create_time", v, nil),
			"updated_at":   utils.PathSearch("update_time", v, nil),
		})
	}
	return rst
}

func flattenGetDataStandardResponseBodyNewBiz(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("new_biz", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"biz_type":    utils.PathSearch("biz_type", v, nil),
			"biz_id":      utils.PathSearch("biz_id", v, nil),
			"biz_info":    utils.PathSearch("biz_info", v, nil),
			"status":      utils.PathSearch("status", v, nil),
			"biz_version": utils.PathSearch("biz_version", v, nil),
			"created_at":  utils.PathSearch("create_time", v, nil),
			"updated_at":  utils.PathSearch("update_time", v, nil),
		})
	}
	return rst
}

func resourceDataStandardUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateDataStandardChanges := []string{
		"directory_id",
		"values",
	}

	if d.HasChanges(updateDataStandardChanges...) {
		// updateDataStandard: update DataArts Architecture data standard
		var (
			updateDataStandardHttpUrl = "v2/{project_id}/design/standards/{id}"
			updateDataStandardProduct = "dataarts"
		)
		updateDataStandardClient, err := cfg.NewServiceClient(updateDataStandardProduct, region)
		if err != nil {
			return diag.Errorf("error creating DataArts Studio client: %s", err)
		}

		updateDataStandardPath := updateDataStandardClient.Endpoint + updateDataStandardHttpUrl
		updateDataStandardPath = strings.ReplaceAll(updateDataStandardPath, "{project_id}", updateDataStandardClient.ProjectID)
		updateDataStandardPath = strings.ReplaceAll(updateDataStandardPath, "{id}", d.Id())

		updateDataStandardOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"workspace": d.Get("workspace_id").(string),
			},
		}

		dataStandardRespBody, err := getDataStandard(d, updateDataStandardClient)
		if err != nil {
			return diag.FromErr(err)
		}
		dataStandardNameToIdMap := buildDataStandardNameToIdMap(dataStandardRespBody)
		updateDataStandardOpt.JSONBody = utils.RemoveNil(buildUpdateDataStandardBodyParams(d, dataStandardNameToIdMap))
		_, err = updateDataStandardClient.Request("PUT", updateDataStandardPath, &updateDataStandardOpt)
		if err != nil {
			return diag.Errorf("error updating DataArts Architecture data standard: %s", err)
		}
	}
	return resourceDataStandardRead(ctx, d, meta)
}

func getDataStandard(d *schema.ResourceData, client *golangsdk.ServiceClient) (interface{}, error) {
	var (
		getDataStandardHttpUrl = "v2/{project_id}/design/standards/{id}"
	)
	getDataStandardPath := client.Endpoint + getDataStandardHttpUrl
	getDataStandardPath = strings.ReplaceAll(getDataStandardPath, "{project_id}", client.ProjectID)
	getDataStandardPath = strings.ReplaceAll(getDataStandardPath, "{id}", d.Id())

	getDataStandardOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"workspace": d.Get("workspace_id").(string),
		},
	}

	getDataStandardResp, err := client.Request("GET", getDataStandardPath, &getDataStandardOpt)

	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getDataStandardResp)
}

func buildDataStandardNameToIdMap(resp interface{}) map[string]string {
	rst := make(map[string]string)
	if resp == nil {
		return rst
	}
	curJson := utils.PathSearch("data.value.values", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	for _, v := range curArray {
		id := utils.PathSearch("id", v, "").(string)
		fdName := utils.PathSearch("fd_name", v, "").(string)
		rst[fdName] = id
	}
	return rst
}

func buildUpdateDataStandardBodyParams(d *schema.ResourceData, dataStandardNameToIdMap map[string]string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"id":           d.Id(),
		"directory_id": d.Get("directory_id"),
		"values":       buildUpdateDataStandardRequestBodyValue(d, dataStandardNameToIdMap),
	}
	return bodyParams
}

func buildUpdateDataStandardRequestBodyValue(d *schema.ResourceData, dataStandardNameToIdMap map[string]string) []map[string]interface{} {
	rawArray := d.Get("values").(*schema.Set).List()
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		raw := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"id":           dataStandardNameToIdMap[raw["fd_name"].(string)],
			"fd_name":      raw["fd_name"],
			"fd_value":     raw["fd_value"],
			"directory_id": d.Get("directory_id"),
		}
	}
	return rst
}

func resourceDataStandardDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteDataStandard: delete DataArts Architecture data standard
	var (
		deleteDataStandardHttpUrl = "v2/{project_id}/design/standards"
		deleteDataStandardProduct = "dataarts"
	)
	deleteDataStandardClient, err := cfg.NewServiceClient(deleteDataStandardProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	deleteDataStandardPath := deleteDataStandardClient.Endpoint + deleteDataStandardHttpUrl
	deleteDataStandardPath = strings.ReplaceAll(deleteDataStandardPath, "{project_id}", deleteDataStandardClient.ProjectID)
	deleteDataStandardPath = strings.ReplaceAll(deleteDataStandardPath, "{id}", d.Id())

	deleteDataStandardOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"workspace": d.Get("workspace_id").(string),
		},
	}

	deleteDataStandardOpt.JSONBody = utils.RemoveNil(buildDeleteDataStandardBodyParams(d))
	_, err = deleteDataStandardClient.Request("DELETE", deleteDataStandardPath, &deleteDataStandardOpt)
	if err != nil {
		return diag.Errorf("error deleting DataArts Architecture data standard: %s", err)
	}

	return nil
}

func buildDeleteDataStandardBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"ids": []string{d.Id()},
	}
	return bodyParams
}

// resourceDataStandardImportState use to import an id with format <workspace_id>/<id>
func resourceDataStandardImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <workspace_id>/<id>")
	}

	d.SetId(parts[1])
	mErr := multierror.Append(
		nil,
		d.Set("workspace_id", parts[0]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set values in import state, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
