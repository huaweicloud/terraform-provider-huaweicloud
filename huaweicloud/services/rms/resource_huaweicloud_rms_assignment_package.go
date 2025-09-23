// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RMS
// ---------------------------------------------------------------

package rms

import (
	"context"
	"encoding/json"
	"log"
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

// @API Config POST /v1/resource-manager/domains/{domain_id}/conformance-packs
// @API Config DELETE /v1/resource-manager/domains/{domain_id}/conformance-packs/{id}
// @API Config GET /v1/resource-manager/domains/{domain_id}/conformance-packs/{id}
// @API Config PUT /v1/resource-manager/domains/{domain_id}/conformance-packs/{id}
func ResourceAssignmentPackage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssignmentPackageCreate,
		ReadContext:   resourceAssignmentPackageRead,
		UpdateContext: resourceAssignmentPackageUpdate,
		DeleteContext: resourceAssignmentPackageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the assignment package name.`,
			},
			"agency_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the agency name.`,
			},
			"template_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  `Specifies the name of a built-in assignment package template.`,
				ExactlyOneOf: []string{"template_key", "template_body", "template_uri"},
			},
			"template_body": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the content of a custom assignment package.`,
			},
			"template_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the URL address of the OBS bucket where an assignment package template was stored.`,
			},
			"vars_structure": {
				Type:        schema.TypeSet,
				Elem:        assignmentPackageParameterSchema(),
				Optional:    true,
				Computed:    true,
				Description: `Specifies the parameters of an assignment package.`,
			},
			"stack_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the unique ID of a resource stack.`,
			},
			"stack_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of a resource stack.`,
			},
			"deployment_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the deployment ID.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the deployment status of an assignment package.`,
			},
		},
	}
}

func assignmentPackageParameterSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"var_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the name of a parameter.`,
			},
			"var_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the value of a parameter.`,
			},
		},
	}
	return &sc
}

func resourceAssignmentPackageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createAssignmentPackage: Create a RMS assignment package.
	var (
		createAssignmentPackageHttpUrl = "v1/resource-manager/domains/{domain_id}/conformance-packs"
		createAssignmentPackageProduct = "rms"
	)
	createAssignmentPackageClient, err := cfg.NewServiceClient(createAssignmentPackageProduct, region)
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	createAssignmentPackagePath := createAssignmentPackageClient.Endpoint + createAssignmentPackageHttpUrl
	createAssignmentPackagePath = strings.ReplaceAll(createAssignmentPackagePath, "{domain_id}", cfg.DomainID)

	createAssignmentPackageOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createOpts, err := buildCreateAssignmentPackageBodyParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	createAssignmentPackageOpt.JSONBody = utils.RemoveNil(createOpts)
	log.Printf("[DEBUG] Create RMS assignment package options: %#v", createAssignmentPackageOpt)
	createAssignmentPackageResp, err := createAssignmentPackageClient.Request("POST",
		createAssignmentPackagePath, &createAssignmentPackageOpt)
	if err != nil {
		return diag.Errorf("error creating RMS assignment package: %s", err)
	}

	createAssignmentPackageRespBody, err := utils.FlattenResponse(createAssignmentPackageResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", createAssignmentPackageRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating RMS assignment package: ID is not found in API response")
	}
	d.SetId(id)

	stateConf := &resource.StateChangeConf{
		Target:       []string{"CREATE_SUCCESSFUL", "ROLLBACK_SUCCESSFUL"},
		Pending:      []string{"CREATE_IN_PROGRESS"},
		Refresh:      rmsAssignmentPackageStateRefreshFunc(createAssignmentPackageClient, cfg, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for RMS assignment Package (%s) to be created: %s", id, err)
	}

	return resourceAssignmentPackageRead(ctx, d, meta)
}

func buildCreateAssignmentPackageBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	varsStructure, err := buildAssignmentPackageRequestBodyParameter(d.Get("vars_structure"))
	if err != nil {
		return nil, err
	}

	bodyParams := map[string]interface{}{
		"name":           d.Get("name"),
		"agency_name":    utils.ValueIgnoreEmpty(d.Get("agency_name")),
		"template_key":   utils.ValueIgnoreEmpty(d.Get("template_key")),
		"template_body":  utils.ValueIgnoreEmpty(d.Get("template_body")),
		"template_uri":   utils.ValueIgnoreEmpty(d.Get("template_uri")),
		"vars_structure": varsStructure,
	}
	return bodyParams, err
}

func buildAssignmentPackageRequestBodyParameter(rawParams interface{}) ([]map[string]interface{}, error) {
	rawArray := rawParams.(*schema.Set).List()
	if len(rawArray) == 0 {
		return nil, nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		if raw, ok := v.(map[string]interface{}); ok {
			var value interface{}
			err := json.Unmarshal([]byte(raw["var_value"].(string)), &value)
			if err != nil {
				return nil, err
			}
			rst[i] = map[string]interface{}{
				"var_key":   utils.ValueIgnoreEmpty(raw["var_key"]),
				"var_value": value,
			}
		}
	}
	return rst, nil
}

func resourceAssignmentPackageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getAssignmentPackage: Query the RMS assignment package
	var (
		getAssignmentPackageProduct = "rms"
	)
	getAssignmentPackageClient, err := cfg.NewServiceClient(getAssignmentPackageProduct, region)
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	getAssignmentPackageRespBody, err := getAssignmentPackage(getAssignmentPackageClient, cfg, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RMS assignment package")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("name", utils.PathSearch("name", getAssignmentPackageRespBody, nil)),
		d.Set("stack_id", utils.PathSearch("stack_id", getAssignmentPackageRespBody, nil)),
		d.Set("stack_name", utils.PathSearch("stack_name", getAssignmentPackageRespBody, nil)),
		d.Set("deployment_id", utils.PathSearch("deployment_id", getAssignmentPackageRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getAssignmentPackageRespBody, nil)),
		d.Set("vars_structure", flattenGetAssignmentPackageResponseBodyParameter(getAssignmentPackageRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetAssignmentPackageResponseBodyParameter(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("vars_structure", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"var_key":   utils.PathSearch("var_key", v, nil),
			"var_value": utils.JsonToString(utils.PathSearch("var_value", v, nil)),
		})
	}
	return rst
}

func resourceAssignmentPackageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	assignmentPackageId := d.Id()

	// updateAssignmentPackage: Update an existing RMS assignment package
	var (
		updateAssignmentPackageHttpUrl = "v1/resource-manager/domains/{domain_id}/conformance-packs/{id}"
		updateAssignmentPackageProduct = "rms"
	)
	updateAssignmentPackageClient, err := cfg.NewServiceClient(updateAssignmentPackageProduct, region)
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	updateAssignmentPackagePath := updateAssignmentPackageClient.Endpoint + updateAssignmentPackageHttpUrl
	updateAssignmentPackagePath = strings.ReplaceAll(updateAssignmentPackagePath, "{domain_id}", cfg.DomainID)
	updateAssignmentPackagePath = strings.ReplaceAll(updateAssignmentPackagePath, "{id}", assignmentPackageId)

	updateAssignmentPackageOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	updateOpts, err := buildUpdateAssignmentPackageBodyParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateAssignmentPackageOpt.JSONBody = utils.RemoveNil(updateOpts)
	log.Printf("[DEBUG] Update RMS assignment package options: %#v", updateAssignmentPackageOpt)
	_, err = updateAssignmentPackageClient.Request("PUT", updateAssignmentPackagePath, &updateAssignmentPackageOpt)
	if err != nil {
		return diag.Errorf("error updating RMS assignment package: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Target:       []string{"UPDATE_SUCCESSFUL", "ROLLBACK_SUCCESSFUL"},
		Pending:      []string{"UPDATE_IN_PROGRESS"},
		Refresh:      rmsAssignmentPackageStateRefreshFunc(updateAssignmentPackageClient, cfg, assignmentPackageId),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for RMS assignment Package (%s) to be updated: %s", assignmentPackageId, err)
	}

	return resourceAssignmentPackageRead(ctx, d, meta)
}

func buildUpdateAssignmentPackageBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	varsStructure, err := buildAssignmentPackageRequestBodyParameter(d.Get("vars_structure"))
	if err != nil {
		return nil, err
	}

	bodyParams := map[string]interface{}{
		"name":           d.Get("name"),
		"vars_structure": varsStructure,
	}
	return bodyParams, err
}

func resourceAssignmentPackageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteAssignmentPackage: Delete an existing RMS assignment package
	var (
		deleteAssignmentPackageHttpUrl = "v1/resource-manager/domains/{domain_id}/conformance-packs/{id}"
		deleteAssignmentPackageProduct = "rms"
	)
	deleteAssignmentPackageClient, err := cfg.NewServiceClient(deleteAssignmentPackageProduct, region)
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	deleteAssignmentPackagePath := deleteAssignmentPackageClient.Endpoint + deleteAssignmentPackageHttpUrl
	deleteAssignmentPackagePath = strings.ReplaceAll(deleteAssignmentPackagePath, "{domain_id}", cfg.DomainID)
	deleteAssignmentPackagePath = strings.ReplaceAll(deleteAssignmentPackagePath, "{id}", d.Id())

	deleteAssignmentPackageOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = deleteAssignmentPackageClient.Request("DELETE", deleteAssignmentPackagePath,
		&deleteAssignmentPackageOpt)
	if err != nil {
		return diag.Errorf("error deleting RMS assignment package: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Target:       []string{"DELETED"},
		Pending:      []string{"DELETE_IN_PROGRESS"},
		Refresh:      rmsAssignmentPackageStateRefreshFunc(deleteAssignmentPackageClient, cfg, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for RMS assignment Package (%s) to be deleted: %s", d.Id(), err)
	}

	return nil
}

func rmsAssignmentPackageStateRefreshFunc(client *golangsdk.ServiceClient, cfg *config.Config, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getAssignmentPackageRespBody, err := getAssignmentPackage(client, cfg, id)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
				return "Resource Not Found", "DELETED", nil
			}
			return nil, "", err
		}
		status := utils.PathSearch("status", getAssignmentPackageRespBody, "").(string)
		return getAssignmentPackageRespBody, status, nil
	}
}

func getAssignmentPackage(client *golangsdk.ServiceClient, cfg *config.Config, id string) (interface{}, error) {
	// getAssignmentPackage: Query the RMS assignment package
	var (
		getAssignmentPackageHttpUrl = "v1/resource-manager/domains/{domain_id}/conformance-packs/{id}"
	)

	getAssignmentPackagePath := client.Endpoint + getAssignmentPackageHttpUrl
	getAssignmentPackagePath = strings.ReplaceAll(getAssignmentPackagePath, "{domain_id}", cfg.DomainID)
	getAssignmentPackagePath = strings.ReplaceAll(getAssignmentPackagePath, "{id}", id)

	getAssignmentPackageOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getAssignmentPackageResp, err := client.Request("GET", getAssignmentPackagePath,
		&getAssignmentPackageOpt)

	if err != nil {
		return getAssignmentPackageResp, err
	}

	return utils.FlattenResponse(getAssignmentPackageResp)
}
