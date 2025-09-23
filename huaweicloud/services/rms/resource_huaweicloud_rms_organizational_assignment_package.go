// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product Config
// ---------------------------------------------------------------

package rms

import (
	"context"
	"encoding/json"
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

// @API Config POST /v1/resource-manager/organizations/{organization_id}/conformance-packs
// @API Config GET /v1/resource-manager/organizations/{organization_id}/conformance-packs/statuses
// @API Config GET /v1/resource-manager/organizations/{organization_id}/conformance-packs/{conformance_pack_id}
// @API Config DELETE /v1/resource-manager/organizations/{organization_id}/conformance-packs/{conformance_pack_id}
// @API Config PUT /v1/resource-manager/organizations/{organization_id}/conformance-packs/{conformance_pack_id}
func ResourceOrgAssignmentPackage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrgAssignmentPackageCreate,
		ReadContext:   resourceOrgAssignmentPackageRead,
		UpdateContext: resourceOrgAssignmentPackageUpdate,
		DeleteContext: resourceOrgAssignmentPackageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceOrgAssignmentPackageImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the organization.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the assignment package name.`,
			},
			"excluded_accounts": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Specifies the excluded accounts for conformance package deployment.`,
			},
			"template_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  `Specifies the name of a predefined conformance package.`,
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
				Description: `Specifies the OBS address of a conformance package.`,
			},
			"vars_structure": {
				Type:        schema.TypeSet,
				Elem:        orgAssignmentPackageVarStructureSchema(),
				Optional:    true,
				Computed:    true,
				Description: `Specifies the parameters of a conformance package.`,
			},
			"owner_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creator of an organization conformance package.`,
			},
			"org_conformance_pack_urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the unique identifier of resources in an organization conformance package.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time of an organization conformance package.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the latest update time of an organization conformance package.`,
			},
		},
	}
}

func orgAssignmentPackageVarStructureSchema() *schema.Resource {
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

func resourceOrgAssignmentPackageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createOrgAssignmentPackage: Create RMS organizational assignment package.
	var (
		createOrgAssignmentPackageHttpUrl = "v1/resource-manager/organizations/{organization_id}/conformance-packs"
		createOrgAssignmentPackageProduct = "rms"
	)
	createOrgAssignmentPackageClient, err := cfg.NewServiceClient(createOrgAssignmentPackageProduct, region)
	if err != nil {
		return diag.Errorf("error creating Config client: %s", err)
	}

	createOrgAssignmentPackagePath := createOrgAssignmentPackageClient.Endpoint + createOrgAssignmentPackageHttpUrl
	createOrgAssignmentPackagePath = strings.ReplaceAll(createOrgAssignmentPackagePath, "{organization_id}",
		d.Get("organization_id").(string))

	createOrgAssignmentPackageOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createOpts, err := buildCreateOrgAssignmentPackageBodyParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	createOrgAssignmentPackageOpt.JSONBody = utils.RemoveNil(createOpts)
	createOrgAssignmentPackageResp, err := createOrgAssignmentPackageClient.Request("POST",
		createOrgAssignmentPackagePath, &createOrgAssignmentPackageOpt)
	if err != nil {
		return diag.Errorf("error creating RMS organizational assignment package: %s", err)
	}

	createOrgAssignmentPackageRespBody, err := utils.FlattenResponse(createOrgAssignmentPackageResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("org_conformance_pack_id", createOrgAssignmentPackageRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating RMS organizational assignment package: ID is not found in API response")
	}
	d.SetId(id)

	stateConf := &resource.StateChangeConf{
		Target:       []string{"CREATE_SUCCESSFUL", "ROLLBACK_SUCCESSFUL"},
		Pending:      []string{"CREATE_IN_PROGRESS"},
		Refresh:      refreshDeployStatus(d, createOrgAssignmentPackageClient),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for RMS organizational assignment package (%s) to be created: %s",
			id, err)
	}

	return resourceOrgAssignmentPackageRead(ctx, d, meta)
}

func buildCreateOrgAssignmentPackageBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	varsStructure, err := buildOrgAssignmentPackageRequestBodyVarStructure(d.Get("vars_structure"))
	if err != nil {
		return nil, err
	}

	bodyParams := map[string]interface{}{
		"name":              d.Get("name"),
		"excluded_accounts": utils.ValueIgnoreEmpty(d.Get("excluded_accounts")),
		"template_key":      utils.ValueIgnoreEmpty(d.Get("template_key")),
		"template_body":     utils.ValueIgnoreEmpty(d.Get("template_body")),
		"template_uri":      utils.ValueIgnoreEmpty(d.Get("template_uri")),
		"vars_structure":    varsStructure,
	}
	return bodyParams, nil
}

func buildOrgAssignmentPackageRequestBodyVarStructure(rawParams interface{}) ([]map[string]interface{}, error) {
	rawArray := rawParams.(*schema.Set).List()
	if len(rawArray) == 0 {
		return nil, nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		raw := v.(map[string]interface{})
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
	return rst, nil
}

func resourceOrgAssignmentPackageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getOrgAssignmentPackage: Query the RMS organizational assignment package
	var (
		getOrgAssignmentPackageHttpUrl = "v1/resource-manager/organizations/{organization_id}/conformance-packs/{conformance_pack_id}"
		getOrgAssignmentPackageProduct = "rms"
	)
	getOrgAssignmentPackageClient, err := cfg.NewServiceClient(getOrgAssignmentPackageProduct, region)
	if err != nil {
		return diag.Errorf("error creating Config client: %s", err)
	}

	getOrgAssignmentPackagePath := getOrgAssignmentPackageClient.Endpoint + getOrgAssignmentPackageHttpUrl
	getOrgAssignmentPackagePath = strings.ReplaceAll(getOrgAssignmentPackagePath, "{organization_id}",
		d.Get("organization_id").(string))
	getOrgAssignmentPackagePath = strings.ReplaceAll(getOrgAssignmentPackagePath, "{conformance_pack_id}", d.Id())

	getOrgAssignmentPackageOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getOrgAssignmentPackageResp, err := getOrgAssignmentPackageClient.Request("GET", getOrgAssignmentPackagePath,
		&getOrgAssignmentPackageOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RMS organizational assignment package")
	}

	getOrgAssignmentPackageRespBody, err := utils.FlattenResponse(getOrgAssignmentPackageResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("organization_id", utils.PathSearch("organization_id", getOrgAssignmentPackageRespBody, nil)),
		d.Set("name", utils.PathSearch("org_conformance_pack_name", getOrgAssignmentPackageRespBody, nil)),
		d.Set("owner_id", utils.PathSearch("owner_id", getOrgAssignmentPackageRespBody, nil)),
		d.Set("organization_id", utils.PathSearch("organization_id", getOrgAssignmentPackageRespBody, nil)),
		d.Set("org_conformance_pack_urn", utils.PathSearch("org_conformance_pack_urn",
			getOrgAssignmentPackageRespBody, nil)),
		d.Set("excluded_accounts", utils.PathSearch("excluded_accounts", getOrgAssignmentPackageRespBody, nil)),
		d.Set("vars_structure", flattenGetOrgAssignmentPackageResponseBodyVarStructure(getOrgAssignmentPackageRespBody)),
		d.Set("created_at", utils.PathSearch("created_at", getOrgAssignmentPackageRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", getOrgAssignmentPackageRespBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetOrgAssignmentPackageResponseBodyVarStructure(resp interface{}) []interface{} {
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

func resourceOrgAssignmentPackageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updateOrgAssignmentPackage: Update an existing RMS organizational assignment package
	var (
		updateOrgAssignmentPackageHttpUrl = "v1/resource-manager/organizations/{organization_id}/conformance-packs/{conformance_pack_id}"
		updateOrgAssignmentPackageProduct = "rms"
		conformancePackId                 = d.Id()
	)
	updateOrgAssignmentPackageClient, err := cfg.NewServiceClient(updateOrgAssignmentPackageProduct, region)
	if err != nil {
		return diag.Errorf("error creating Config client: %s", err)
	}

	updateOrgAssignmentPackagePath := updateOrgAssignmentPackageClient.Endpoint + updateOrgAssignmentPackageHttpUrl
	updateOrgAssignmentPackagePath = strings.ReplaceAll(updateOrgAssignmentPackagePath, "{organization_id}",
		d.Get("organization_id").(string))
	updateOrgAssignmentPackagePath = strings.ReplaceAll(updateOrgAssignmentPackagePath, "{conformance_pack_id}", conformancePackId)

	updateOrgAssignmentPackageOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	updateOpts, err := buildUpdateOrgAssignmentPackageBodyParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	updateOrgAssignmentPackageOpt.JSONBody = utils.RemoveNil(updateOpts)
	_, err = updateOrgAssignmentPackageClient.Request("PUT",
		updateOrgAssignmentPackagePath, &updateOrgAssignmentPackageOpt)
	if err != nil {
		return diag.Errorf("error updating RMS organizational assignment package: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Target:       []string{"UPDATE_SUCCESSFUL", "ROLLBACK_SUCCESSFUL"},
		Pending:      []string{"UPDATE_IN_PROGRESS"},
		Refresh:      refreshDeployStatus(d, updateOrgAssignmentPackageClient),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for RMS organizational assignment package (%s) to be updated: %s",
			conformancePackId, err)
	}

	return resourceOrgAssignmentPackageRead(ctx, d, meta)
}

func buildUpdateOrgAssignmentPackageBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	varsStructure, err := buildOrgAssignmentPackageRequestBodyVarStructure(d.Get("vars_structure"))
	if err != nil {
		return nil, err
	}

	bodyParams := map[string]interface{}{
		"name":              d.Get("name"),
		"excluded_accounts": utils.ValueIgnoreEmpty(d.Get("excluded_accounts")),
		"vars_structure":    varsStructure,
	}
	return bodyParams, nil
}

func resourceOrgAssignmentPackageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteOrgAssignmentPackage: Delete an existing RMS organizational assignment package
	var (
		deleteOrgAssignmentPackageHttpUrl = "v1/resource-manager/organizations/{organization_id}/conformance-packs/{conformance_pack_id}"
		deleteOrgAssignmentPackageProduct = "rms"
	)
	deleteOrgAssignmentPackageClient, err := cfg.NewServiceClient(deleteOrgAssignmentPackageProduct, region)
	if err != nil {
		return diag.Errorf("error creating Config client: %s", err)
	}

	deleteOrgAssignmentPackagePath := deleteOrgAssignmentPackageClient.Endpoint + deleteOrgAssignmentPackageHttpUrl
	deleteOrgAssignmentPackagePath = strings.ReplaceAll(deleteOrgAssignmentPackagePath, "{organization_id}",
		d.Get("organization_id").(string))
	deleteOrgAssignmentPackagePath = strings.ReplaceAll(deleteOrgAssignmentPackagePath, "{conformance_pack_id}", d.Id())

	deleteOrgAssignmentPackageOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = deleteOrgAssignmentPackageClient.Request("DELETE", deleteOrgAssignmentPackagePath,
		&deleteOrgAssignmentPackageOpt)
	if err != nil {
		return diag.Errorf("error deleting RMS organizational assignment package: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Target:       []string{""},
		Pending:      []string{"DELETE_IN_PROGRESS"},
		Refresh:      refreshDeployStatus(d, deleteOrgAssignmentPackageClient),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for RMS organizational assignment package (%s) to be deleted: %s", d.Id(), err)
	}

	return nil
}

func refreshDeployStatus(d *schema.ResourceData, client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (result interface{}, state string, err error) {
		// getDeployStatus: Query the RMS organizational assignment package
		var (
			getDeployStatusHttpUrl = "v1/resource-manager/organizations/{organization_id}/conformance-packs/statuses"
		)

		getDeployStatusPath := client.Endpoint + getDeployStatusHttpUrl
		getDeployStatusPath = strings.ReplaceAll(getDeployStatusPath, "{organization_id}", d.Get("organization_id").(string))

		conformancePackName := d.Get("name").(string)
		getDeployStatusQueryParams := buildGetDeployStatusQueryParams(conformancePackName)
		getDeployStatusPath += getDeployStatusQueryParams

		getDeployStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		getDeployStatusResp, err := client.Request("GET", getDeployStatusPath, &getDeployStatusOpt)
		if err != nil {
			return nil, "", err
		}
		getDeployStatusRespBody, err := utils.FlattenResponse(getDeployStatusResp)
		if err != nil {
			return nil, "", err
		}
		return getDeployStatusRespBody, utils.PathSearch("statuses|[0].state", getDeployStatusRespBody, "").(string), nil
	}
}

func buildGetDeployStatusQueryParams(conformancePackName string) string {
	return fmt.Sprintf("?conformance_pack_name=%s", conformancePackName)
}

// resourceOrgAssignmentPackageImportState use to import an id with format <organization_id>/<id>
func resourceOrgAssignmentPackageImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	if !strings.Contains(d.Id(), "/") {
		return []*schema.ResourceData{d}, nil
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <organization_id>/<id>")
	}

	d.SetId(parts[1])
	mErr := multierror.Append(
		nil,
		d.Set("organization_id", parts[0]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set values in import state, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
