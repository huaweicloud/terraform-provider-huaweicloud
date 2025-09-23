// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product Organizations
// ---------------------------------------------------------------

package organizations

import (
	"context"
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

// @API Organizations POST /v1/organizations/organizational-units
// @API Organizations PATCH /v1/organizations/organizational-units/{organizational_unit_id}
// @API Organizations POST /v1/organizations/{resource_type}/{resource_id}/tags/create
// @API Organizations POST /v1/organizations/{resource_type}/{resource_id}/tags/delete
// @API Organizations GET /v1/organizations/organizational-units/{organizational_unit_id}
// @API Organizations GET /v1/organizations/{resource_type}/{resource_id}/tags
// @API Organizations DELETE /v1/organizations/organizational-units/{organizational_unit_id}
func ResourceOrganizationalUnit() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationalUnitCreate,
		UpdateContext: resourceOrganizationalUnitUpdate,
		ReadContext:   resourceOrganizationalUnitRead,
		DeleteContext: resourceOrganizationalUnitDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the organizational unit.`,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `Specifies the ID of the root or organizational unit in which you
want to create a new organizational unit.`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the uniform resource name of the organizational unit.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time when the OU was created.`,
			},
			"tags": common.TagsSchema(),
		},
	}
}

func resourceOrganizationalUnitCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createOrganizationalUnit: create Organizations organizational unit
	var (
		createOrganizationalUnitHttpUrl = "v1/organizations/organizational-units"
		createOrganizationalUnitProduct = "organizations"
	)
	createOrganizationalUnitClient, err := cfg.NewServiceClient(createOrganizationalUnitProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	createOrganizationalUnitPath := createOrganizationalUnitClient.Endpoint + createOrganizationalUnitHttpUrl

	createOrganizationalUnitOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	createOrganizationalUnitOpt.JSONBody = utils.RemoveNil(buildCreateOrganizationalUnitBodyParams(d))
	createOrganizationalUnitResp, err := createOrganizationalUnitClient.Request("POST",
		createOrganizationalUnitPath, &createOrganizationalUnitOpt)
	if err != nil {
		return diag.Errorf("error creating Organizations organizational unit: %s", err)
	}

	createOrganizationalUnitRespBody, err := utils.FlattenResponse(createOrganizationalUnitResp)
	if err != nil {
		return diag.FromErr(err)
	}

	unitId := utils.PathSearch("organizational_unit.id", createOrganizationalUnitRespBody, "").(string)
	if unitId == "" {
		return diag.Errorf("unable to find the organizational unit ID from the API response")
	}

	d.SetId(unitId)

	return resourceOrganizationalUnitRead(ctx, d, meta)
}

func buildCreateOrganizationalUnitBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":      utils.ValueIgnoreEmpty(d.Get("name")),
		"parent_id": utils.ValueIgnoreEmpty(d.Get("parent_id")),
		"tags":      utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}
	return bodyParams
}

func resourceOrganizationalUnitRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getOrganizationalUnit: Query Organizations organizational unit
	var (
		getOrganizationalUnitHttpUrl = "v1/organizations/organizational-units/{organizational_unit_id}"
		getOrganizationalUnitProduct = "organizations"
	)
	getOrganizationalUnitClient, err := cfg.NewServiceClient(getOrganizationalUnitProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	getOrganizationalUnitPath := getOrganizationalUnitClient.Endpoint + getOrganizationalUnitHttpUrl
	getOrganizationalUnitPath = strings.ReplaceAll(getOrganizationalUnitPath, "{organizational_unit_id}", d.Id())

	getOrganizationalUnitOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getOrganizationalUnitResp, err := getOrganizationalUnitClient.Request("GET",
		getOrganizationalUnitPath, &getOrganizationalUnitOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Organizations organizational unit")
	}

	getOrganizationalUnitRespBody, err := utils.FlattenResponse(getOrganizationalUnitResp)
	if err != nil {
		return diag.FromErr(err)
	}

	organizationalUnit := utils.PathSearch("organizational_unit", getOrganizationalUnitRespBody, nil)
	if organizationalUnit == nil {
		log.Printf("[WARN] failed to get organizational unit, organizational_unit %s is not found "+
			"in API response", d.Id())
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("name", utils.PathSearch("name", organizationalUnit, nil)),
		d.Set("urn", utils.PathSearch("urn", organizationalUnit, nil)),
		d.Set("created_at", utils.PathSearch("created_at", organizationalUnit, nil)),
	)

	tagMap, err := getTags(getOrganizationalUnitClient, unitType, d.Id())
	if err != nil {
		log.Printf("[WARN] error fetching tags of Organizations organizational unit (%s): %s", d.Id(), err)
	} else {
		mErr = multierror.Append(mErr, d.Set("tags", tagMap))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceOrganizationalUnitUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updateOrganizationalUnit: update Organizations organizational unit
	var (
		updateOrganizationalUnitHttpUrl = "v1/organizations/organizational-units/{organizational_unit_id}"
		updateOrganizationalUnitProduct = "organizations"
	)
	updateOrganizationalUnitClient, err := cfg.NewServiceClient(updateOrganizationalUnitProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	if d.HasChange("name") {
		updateOrganizationalUnitPath := updateOrganizationalUnitClient.Endpoint + updateOrganizationalUnitHttpUrl
		updateOrganizationalUnitPath = strings.ReplaceAll(updateOrganizationalUnitPath,
			"{organizational_unit_id}", d.Id())

		updateOrganizationalUnitOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateOrganizationalUnitOpt.JSONBody = utils.RemoveNil(buildUpdateOrganizationalUnitBodyParams(d))
		_, err = updateOrganizationalUnitClient.Request("PATCH",
			updateOrganizationalUnitPath, &updateOrganizationalUnitOpt)
		if err != nil {
			return diag.Errorf("error updating Organizations organizational unit: %s", err)
		}
	}

	if d.HasChange("tags") {
		err = updateTags(d, updateOrganizationalUnitClient, unitType, d.Id(), "tags")
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceOrganizationalUnitRead(ctx, d, meta)
}

func buildUpdateOrganizationalUnitBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name": utils.ValueIgnoreEmpty(d.Get("name")),
	}
	return bodyParams
}

func resourceOrganizationalUnitDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteOrganizationalUnit: Delete Organizations organizational unit
	var (
		deleteOrganizationalUnitHttpUrl = "v1/organizations/organizational-units/{organizational_unit_id}"
		deleteOrganizationalUnitProduct = "organizations"
	)
	deleteOrganizationalUnitClient, err := cfg.NewServiceClient(deleteOrganizationalUnitProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	deleteOrganizationalUnitPath := deleteOrganizationalUnitClient.Endpoint + deleteOrganizationalUnitHttpUrl
	deleteOrganizationalUnitPath = strings.ReplaceAll(deleteOrganizationalUnitPath,
		"{organizational_unit_id}", d.Id())

	deleteOrganizationalUnitOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = deleteOrganizationalUnitClient.Request("DELETE", deleteOrganizationalUnitPath,
		&deleteOrganizationalUnitOpt)
	if err != nil {
		return diag.Errorf("error deleting Organizations organizational unit: %s", err)
	}

	return nil
}
