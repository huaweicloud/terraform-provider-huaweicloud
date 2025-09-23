// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product MRS
// ---------------------------------------------------------------

package mrs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API MRS GET /v2/{project_id}/data-connectors
// @API MRS POST /v2/{project_id}/data-connectors
// @API MRS DELETE /v2/{project_id}/data-connectors/{connector_id}
// @API MRS PUT /v2/{project_id}/data-connectors/{connector_id}
func ResourceDataConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataConnectionCreate,
		UpdateContext: resourceDataConnectionUpdate,
		ReadContext:   resourceDataConnectionRead,
		DeleteContext: resourceDataConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The data connection name.`,
			},
			"source_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The type of data source.`,
			},
			"source_info": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     dataConnectionSourceInfoSchema(),
				Required: true,
			},
			"used_clusters": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Cluster IDs that use this data connection, separated by commas.`,
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The status of the data connection.`,
			},
		},
	}
}

func dataConnectionSourceInfoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The instance ID of database.`,
			},
			"db_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of database.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The user name for logging in to the database.`,
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `The password for logging in to the database.`,
			},
		},
	}
	return &sc
}

func resourceDataConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createDataConnectionHttpUrl = "v2/{project_id}/data-connectors"
		createDataConnectionProduct = "mrs"
	)
	createDataConnectionClient, err := cfg.NewServiceClient(createDataConnectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating MRS Client: %s", err)
	}

	createDataConnectionPath := createDataConnectionClient.Endpoint + createDataConnectionHttpUrl
	createDataConnectionPath = strings.ReplaceAll(createDataConnectionPath, "{project_id}", createDataConnectionClient.ProjectID)

	createDataConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}

	createDataConnectionOpt.JSONBody = utils.RemoveNil(buildDataConnectionBodyParams(d))
	createDataConnectionResp, err := createDataConnectionClient.Request("POST", createDataConnectionPath, &createDataConnectionOpt)
	if err != nil {
		return diag.Errorf("error creating data connection: %s", err)
	}

	createDataConnectionRespBody, err := utils.FlattenResponse(createDataConnectionResp)
	if err != nil {
		return diag.FromErr(err)
	}

	connectorId := utils.PathSearch("connector_id", createDataConnectionRespBody, "").(string)
	if connectorId == "" {
		return diag.Errorf("unable to find the data connection ID from the API response")
	}
	d.SetId(connectorId)

	return resourceDataConnectionRead(ctx, d, meta)
}

func buildDataConnectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"data_connector": map[string]interface{}{
			"connector_name": utils.ValueIgnoreEmpty(d.Get("name")),
			"source_type":    utils.ValueIgnoreEmpty(d.Get("source_type")),
			"source_info":    buildDataConnectionRequestBodySourceInfo(d.Get("source_info")),
		},
	}
	return bodyParams
}

func buildDataConnectionRequestBodySourceInfo(rawParams interface{}) string {
	if rawArray, ok := rawParams.([]interface{}); ok {
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"rds_instance_id": utils.ValueIgnoreEmpty(raw["db_instance_id"]),
			"db_name":         utils.ValueIgnoreEmpty(raw["db_name"]),
			"user_name":       utils.ValueIgnoreEmpty(raw["user_name"]),
			"password":        utils.ValueIgnoreEmpty(raw["password"]),
		}

		data, _ := json.Marshal(params)
		return string(data)
	}
	return ""
}

func resourceDataConnectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getDataConnectionHttpUrl = "v2/{project_id}/data-connectors"
		getDataConnectionProduct = "mrs"
	)
	getDataConnectionClient, err := cfg.NewServiceClient(getDataConnectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating MRS Client: %s", err)
	}

	getDataConnectionPath := getDataConnectionClient.Endpoint + getDataConnectionHttpUrl
	getDataConnectionPath = strings.ReplaceAll(getDataConnectionPath, "{project_id}", getDataConnectionClient.ProjectID)

	getDataConnectionqueryParams := fmt.Sprintf("?id=%v", d.Id())
	getDataConnectionPath += getDataConnectionqueryParams

	getDataConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getDataConnectionResp, err := getDataConnectionClient.Request("GET", getDataConnectionPath, &getDataConnectionOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving data connection")
	}

	getDataConnectionRespBody, err := utils.FlattenResponse(getDataConnectionResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("data_connectors[?connector_id =='%s']|[0]", d.Id())
	getDataConnectionRespBody = utils.PathSearch(jsonPath, getDataConnectionRespBody, nil)
	if getDataConnectionRespBody == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("connector_name", getDataConnectionRespBody, nil)),
		d.Set("source_type", utils.PathSearch("source_type", getDataConnectionRespBody, nil)),
		d.Set("source_info", flattenGetDataConnectionResponseBodySourceInfo(getDataConnectionRespBody,
			d.Get("source_info.0.password").(string))),
		d.Set("used_clusters", utils.PathSearch("used_clusters", getDataConnectionRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getDataConnectionRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetDataConnectionResponseBodySourceInfo(resp interface{}, password string) []interface{} {
	var rst []interface{}
	curString := utils.PathSearch("source_info", resp, "").(string)

	var curJson map[string]interface{}
	err := json.Unmarshal([]byte(curString), &curJson)
	if err != nil {
		log.Printf("[ERROR] error parsing source_info to json= %#v", resp)
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"db_instance_id": curJson["rds_instance_id"],
			"db_name":        curJson["db_name"],
			"user_name":      curJson["user_name"],
			"password":       password,
		},
	}
	return rst
}

func resourceDataConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateDataConnectionChanges := []string{
		"source_info",
	}

	if d.HasChanges(updateDataConnectionChanges...) {
		var (
			updateDataConnectionHttpUrl = "v2/{project_id}/data-connectors/{id}"
			updateDataConnectionProduct = "mrs"
		)
		updateDataConnectionClient, err := cfg.NewServiceClient(updateDataConnectionProduct, region)
		if err != nil {
			return diag.Errorf("error creating MRS Client: %s", err)
		}

		updateDataConnectionPath := updateDataConnectionClient.Endpoint + updateDataConnectionHttpUrl
		updateDataConnectionPath = strings.ReplaceAll(updateDataConnectionPath, "{project_id}", updateDataConnectionClient.ProjectID)
		updateDataConnectionPath = strings.ReplaceAll(updateDataConnectionPath, "{id}", d.Id())

		updateDataConnectionOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}

		updateDataConnectionOpt.JSONBody = utils.RemoveNil(buildDataConnectionBodyParams(d))
		_, err = updateDataConnectionClient.Request("PUT", updateDataConnectionPath, &updateDataConnectionOpt)
		if err != nil {
			return diag.Errorf("error updating data connection: %s", err)
		}
	}
	return resourceDataConnectionRead(ctx, d, meta)
}

func resourceDataConnectionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteDataConnectionHttpUrl = "v2/{project_id}/data-connectors/{id}"
		deleteDataConnectionProduct = "mrs"
	)
	deleteDataConnectionClient, err := cfg.NewServiceClient(deleteDataConnectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating MRS Client: %s", err)
	}

	deleteDataConnectionPath := deleteDataConnectionClient.Endpoint + deleteDataConnectionHttpUrl
	deleteDataConnectionPath = strings.ReplaceAll(deleteDataConnectionPath, "{project_id}", deleteDataConnectionClient.ProjectID)
	deleteDataConnectionPath = strings.ReplaceAll(deleteDataConnectionPath, "{id}", d.Id())

	deleteDataConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	_, err = deleteDataConnectionClient.Request("DELETE", deleteDataConnectionPath, &deleteDataConnectionOpt)
	if err != nil {
		return diag.Errorf("error deleting data connection: %s", err)
	}

	return nil
}
