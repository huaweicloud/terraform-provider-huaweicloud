// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product UCS
// ---------------------------------------------------------------

package ucs

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

// @API UCS POST /v1/clustergroups
// @API UCS DELETE /v1/clustergroups/{id}
// @API UCS GET /v1/clustergroups/{id}
// @API UCS PUT /v1/clustergroups/{id}/associatedrules
// @API UCS PUT /v1/clustergroups/{id}/description
func ResourceFleet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFleetCreate,
		UpdateContext: resourceFleetUpdate,
		ReadContext:   resourceFleetRead,
		DeleteContext: resourceFleetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the UCS fleet.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the UCS fleet.`,
			},
			"permissions": {
				Type:        schema.TypeList,
				Elem:        FleetPermissionSchema(),
				Optional:    true,
				Description: `Specifies the permissions associated to the cluster.`,
			},
			"cluster_ids": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `Indicates the list of cluster IDs to add to the UCS fleet.`,
			},
		},
	}
}

func FleetPermissionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"policy_ids": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `Specifies the policy IDs.`,
			},
			"namespaces": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `Specifies the namespaces.`,
			},
		},
	}
	return &sc
}

func resourceFleetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createFleet: Create a UCS Fleet.
	var (
		createFleetHttpUrl = "v1/clustergroups"
		createFleetProduct = "ucs"
	)
	createFleetClient, err := cfg.NewServiceClient(createFleetProduct, region)
	if err != nil {
		return diag.Errorf("error creating UCS Client: %s", err)
	}

	createFleetPath := createFleetClient.Endpoint + createFleetHttpUrl

	createFleetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}

	createFleetOpt.JSONBody = utils.RemoveNil(buildCreateFleetBodyParams(d))
	createFleetResp, err := createFleetClient.Request("POST", createFleetPath, &createFleetOpt)
	if err != nil {
		return diag.Errorf("error creating Fleet: %s", err)
	}

	createFleetRespBody, err := utils.FlattenResponse(createFleetResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("uid", createFleetRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating Fleet: ID is not found in API response")
	}
	d.SetId(id)

	if _, ok := d.GetOk("permissions"); ok {
		err := updatePermissions(d, createFleetClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceFleetRead(ctx, d, meta)
}

func buildCreateFleetBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name": utils.ValueIgnoreEmpty(d.Get("name")),
		},
		"spec": buildCreateFleetSpecOpts(d),
	}
	return bodyParams
}

func buildCreateFleetSpecOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}

	return bodyParams
}

func resourceFleetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getFleet: Query the UCS Fleet detail
	var (
		getFleetHttpUrl = "v1/clustergroups/{id}"
		getFleetProduct = "ucs"
	)
	getFleetClient, err := cfg.NewServiceClient(getFleetProduct, region)
	if err != nil {
		return diag.Errorf("error creating UCS Client: %s", err)
	}

	getFleetPath := getFleetClient.Endpoint + getFleetHttpUrl
	getFleetPath = strings.ReplaceAll(getFleetPath, "{id}", d.Id())

	getFleetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getFleetResp, err := getFleetClient.Request("GET", getFleetPath, &getFleetOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Fleet")
	}

	getFleetRespBody, err := utils.FlattenResponse(getFleetResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("name", utils.PathSearch("metadata.name", getFleetRespBody, nil)),
		d.Set("description", utils.PathSearch("spec.description", getFleetRespBody, nil)),
		d.Set("cluster_ids", utils.PathSearch("spec.clusterIds", getFleetRespBody, nil)),
		d.Set("permissions", flattenPermissions(utils.PathSearch("spec.ruleNamespaces", getFleetRespBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPermissions(permissionsRaw interface{}) []map[string]interface{} {
	if permissionsRaw == nil {
		return nil
	}

	permissions := permissionsRaw.([]interface{})
	res := make([]map[string]interface{}, len(permissions))
	for i, v := range permissions {
		res[i] = map[string]interface{}{
			"namespaces": utils.PathSearch("namespaces", v, nil),
			"policy_ids": utils.PathSearch("rules[*].ruleID", v, nil),
		}
	}

	return res
}

func resourceFleetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	updateFleetProduct := "ucs"
	updateFleetClient, err := cfg.NewServiceClient(updateFleetProduct, region)
	if err != nil {
		return diag.Errorf("error creating UCS Client: %s", err)
	}

	if d.HasChange("description") {
		// updateFleet: Update the UCS Fleet
		updateFleetHttpUrl := "v1/clustergroups/{id}/description"

		updateFleetClient, err := cfg.NewServiceClient(updateFleetProduct, region)
		if err != nil {
			return diag.Errorf("error creating UCS Client: %s", err)
		}

		updateFleetPath := updateFleetClient.Endpoint + updateFleetHttpUrl
		updateFleetPath = strings.ReplaceAll(updateFleetPath, "{id}", d.Id())

		updateFleetOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}

		updateFleetOpt.JSONBody = utils.RemoveNil(buildUpdateFleetBodyParams(d))
		_, err = updateFleetClient.Request("PUT", updateFleetPath, &updateFleetOpt)
		if err != nil {
			return diag.Errorf("error updating Fleet: %s", err)
		}
	}

	if d.HasChanges("permissions") {
		// updateFleetPermissions: Update the permissions associated to the UCS Fleet
		err := updatePermissions(d, updateFleetClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceFleetRead(ctx, d, meta)
}

func updatePermissions(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	updateFleetPermissionsHttpUrl := "v1/clustergroups/{id}/associatedrules"
	updateFleetPermissionsPath := client.Endpoint + updateFleetPermissionsHttpUrl
	updateFleetPermissionsPath = strings.ReplaceAll(updateFleetPermissionsPath, "{id}", d.Id())

	updateFleetPermissionsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	updateFleetPermissionsOpt.JSONBody = utils.RemoveNil(buildUpdateFleetPermissionsBodyParams(d))
	_, err := client.Request("PUT", updateFleetPermissionsPath, &updateFleetPermissionsOpt)
	if err != nil {
		return fmt.Errorf("error updating the permissions of fleet: %s", err)
	}

	return nil
}

func buildUpdateFleetBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"description": d.Get("description"),
	}
	return bodyParams
}

func buildUpdateFleetPermissionsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"ruleIDNamespaces": buildUpdateFleetPermissionsRequestBodyPolicy(d.Get("permissions")),
	}
	return bodyParams
}

func buildUpdateFleetPermissionsRequestBodyPolicy(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"ruleIDs":    utils.ValueIgnoreEmpty(raw["policy_ids"]),
				"namespaces": utils.ValueIgnoreEmpty(raw["namespaces"]),
			}
		}
		return rst
	}
	return nil
}

func resourceFleetDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteFleet: Delete an existing UCS Fleet
	var (
		deleteFleetHttpUrl = "v1/clustergroups/{id}"
		deleteFleetProduct = "ucs"
	)
	deleteFleetClient, err := cfg.NewServiceClient(deleteFleetProduct, region)
	if err != nil {
		return diag.Errorf("error creating UCS Client: %s", err)
	}

	deleteFleetPath := deleteFleetClient.Endpoint + deleteFleetHttpUrl
	deleteFleetPath = strings.ReplaceAll(deleteFleetPath, "{id}", d.Id())

	deleteFleetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	_, err = deleteFleetClient.Request("DELETE", deleteFleetPath, &deleteFleetOpt)
	if err != nil {
		return diag.Errorf("error deleting Fleet: %s", err)
	}

	return nil
}
