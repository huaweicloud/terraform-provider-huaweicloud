// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DLI
// ---------------------------------------------------------------

package dli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DLI POST /v1.0/{project_id}/variables
// @API DLI GET /v1.0/{project_id}/variables
// @API DLI PUT /v1.0/{project_id}/variables/{var_name}
// @API DLI DELETE /v1.0/{project_id}/variables/{var_name}
func ResourceGlobalVariable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlobalVariableCreate,
		UpdateContext: resourceGlobalVariableUpdate,
		ReadContext:   resourceGlobalVariableRead,
		DeleteContext: resourceGlobalVariableDelete,
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
				Description: `The name of a Global variable.`,
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The value of Global variable.`,
			},
			"is_sensitive": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Whether to set a variable as a sensitive variable. The default value is **false**.`,
			},
		},
	}
}

func resourceGlobalVariableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createGlobalVariable: create a Global variable.
	var (
		createGlobalVariableHttpUrl = "v1.0/{project_id}/variables"
		createGlobalVariableProduct = "dli"
	)
	createGlobalVariableClient, err := cfg.NewServiceClient(createGlobalVariableProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	createGlobalVariablePath := createGlobalVariableClient.Endpoint + createGlobalVariableHttpUrl
	createGlobalVariablePath = strings.ReplaceAll(createGlobalVariablePath, "{project_id}", createGlobalVariableClient.ProjectID)

	createGlobalVariableOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createGlobalVariableOpt.JSONBody = utils.RemoveNil(buildCreateGlobalVariableBodyParams(d))
	requestResp, err := createGlobalVariableClient.Request("POST", createGlobalVariablePath, &createGlobalVariableOpt)
	if err != nil {
		return diag.Errorf("error creating DLI global variable: %s", err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}
	if !utils.PathSearch("is_success", respBody, true).(bool) {
		return diag.Errorf("unable to create the global variable: %s",
			utils.PathSearch("message", respBody, "Message Not Found"))
	}

	d.SetId(d.Get("name").(string))

	return resourceGlobalVariableRead(ctx, d, meta)
}

func buildCreateGlobalVariableBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"var_name":     utils.ValueIgnoreEmpty(d.Get("name")),
		"var_value":    utils.ValueIgnoreEmpty(d.Get("value")),
		"is_sensitive": utils.ValueIgnoreEmpty(d.Get("is_sensitive")),
	}
	return bodyParams
}

func resourceGlobalVariableRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getGlobalVariable: Query the Global variable.
	var (
		getGlobalVariableHttpUrl = "v1.0/{project_id}/variables"
		getGlobalVariableProduct = "dli"
	)
	getGlobalVariableClient, err := cfg.NewServiceClient(getGlobalVariableProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	getGlobalVariablePath := getGlobalVariableClient.Endpoint + getGlobalVariableHttpUrl
	getGlobalVariablePath = strings.ReplaceAll(getGlobalVariablePath, "{project_id}", getGlobalVariableClient.ProjectID)

	getGlobalVariableResp, err := pagination.ListAllItems(
		getGlobalVariableClient,
		"offset",
		getGlobalVariablePath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DLI global variable")
	}

	getGlobalVariableRespJson, err := json.Marshal(getGlobalVariableResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getGlobalVariableRespBody interface{}
	err = json.Unmarshal(getGlobalVariableRespJson, &getGlobalVariableRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	if !utils.PathSearch("is_success", getGlobalVariableRespBody, true).(bool) {
		return diag.Errorf("unable to query the global variables: %s",
			utils.PathSearch("message", getGlobalVariableRespBody, "Message Not Found"))
	}
	jsonPath := fmt.Sprintf("global_vars[?var_name=='%s']|[0]", d.Id())
	globalVariable := utils.PathSearch(jsonPath, getGlobalVariableRespBody, nil)
	if globalVariable == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("var_name", globalVariable, nil)),
		d.Set("value", utils.PathSearch("var_value", globalVariable, nil)),
		d.Set("is_sensitive", utils.PathSearch("is_sensitive", globalVariable, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceGlobalVariableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateGlobalVariableChanges := []string{
		"value",
	}

	if d.HasChanges(updateGlobalVariableChanges...) {
		// updateGlobalVariable: update Global variable
		var (
			updateGlobalVariableHttpUrl = "v1.0/{project_id}/variables/{var_name}"
			updateGlobalVariableProduct = "dli"
		)
		updateGlobalVariableClient, err := cfg.NewServiceClient(updateGlobalVariableProduct, region)
		if err != nil {
			return diag.Errorf("error creating DLI Client: %s", err)
		}

		updateGlobalVariablePath := updateGlobalVariableClient.Endpoint + updateGlobalVariableHttpUrl
		updateGlobalVariablePath = strings.ReplaceAll(updateGlobalVariablePath, "{project_id}", updateGlobalVariableClient.ProjectID)
		updateGlobalVariablePath = strings.ReplaceAll(updateGlobalVariablePath, "{var_name}", d.Id())

		updateGlobalVariableOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateGlobalVariableOpt.JSONBody = utils.RemoveNil(buildUpdateGlobalVariableBodyParams(d))
		requestResp, err := updateGlobalVariableClient.Request("PUT", updateGlobalVariablePath, &updateGlobalVariableOpt)
		if err != nil {
			return diag.Errorf("error updating DLI global variable: %s", err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return diag.FromErr(err)
		}
		if !utils.PathSearch("is_success", respBody, true).(bool) {
			return diag.Errorf("unable to update the global variable: %s",
				utils.PathSearch("message", respBody, "Message Not Found"))
		}
	}
	return resourceGlobalVariableRead(ctx, d, meta)
}

func buildUpdateGlobalVariableBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"var_value": utils.ValueIgnoreEmpty(d.Get("value")),
	}
	return bodyParams
}

func resourceGlobalVariableDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteGlobalVariable: delete Global variable
	var (
		deleteGlobalVariableHttpUrl = "v1.0/{project_id}/variables/{var_name}"
		deleteGlobalVariableProduct = "dli"
	)
	deleteGlobalVariableClient, err := cfg.NewServiceClient(deleteGlobalVariableProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	deleteGlobalVariablePath := deleteGlobalVariableClient.Endpoint + deleteGlobalVariableHttpUrl
	deleteGlobalVariablePath = strings.ReplaceAll(deleteGlobalVariablePath, "{project_id}", deleteGlobalVariableClient.ProjectID)
	deleteGlobalVariablePath = strings.ReplaceAll(deleteGlobalVariablePath, "{var_name}", d.Id())

	deleteGlobalVariableOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	requestResp, err := deleteGlobalVariableClient.Request("DELETE", deleteGlobalVariablePath, &deleteGlobalVariableOpt)
	if err != nil {
		return diag.Errorf("error deleting DLI global variable: %s", err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}
	if !utils.PathSearch("is_success", respBody, true).(bool) {
		return diag.Errorf("unable to delete the global variable: %s",
			utils.PathSearch("message", respBody, "Message Not Found"))
	}

	return nil
}
