// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product Organizations
// ---------------------------------------------------------------

package organizations

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceOrganizationsUnit() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationsUnitCreate,
		UpdateContext: resourceOrganizationsUnitUpdate,
		ReadContext:   resourceOrganizationsUnitRead,
		DeleteContext: resourceOrganizationsUnitDelete,
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
				Description: `Specifies the name of the organization unit.`,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `Specifies the ID of the root or organization unit in which you
want to create a new organizations unit.`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the uniform resource name of the organization unit.`,
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

func resourceOrganizationsUnitCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createOrganizationsUnit: create Organizations unit
	var (
		createOrganizationsUnitHttpUrl = "v1/organizations/organizational-units"
		createOrganizationsUnitProduct = "organizations"
	)
	createOrganizationsUnitClient, err := cfg.NewServiceClient(createOrganizationsUnitProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	createOrganizationsUnitPath := createOrganizationsUnitClient.Endpoint + createOrganizationsUnitHttpUrl

	createOrganizationsUnitOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	createOrganizationsUnitOpt.JSONBody = utils.RemoveNil(buildCreateOrganizationsUnitBodyParams(d))
	createOrganizationsUnitResp, err := createOrganizationsUnitClient.Request("POST",
		createOrganizationsUnitPath, &createOrganizationsUnitOpt)
	if err != nil {
		return diag.Errorf("error creating Organizations: %s", err)
	}

	createOrganizationsUnitRespBody, err := utils.FlattenResponse(createOrganizationsUnitResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("organizational_unit.id", createOrganizationsUnitRespBody)
	if err != nil {
		return diag.Errorf("error creating Organizations: ID is not found in API response")
	}
	d.SetId(id.(string))

	return resourceOrganizationsUnitRead(ctx, d, meta)
}

func buildCreateOrganizationsUnitBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":      utils.ValueIngoreEmpty(d.Get("name")),
		"parent_id": utils.ValueIngoreEmpty(d.Get("parent_id")),
		"tags":      utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}
	return bodyParams
}

func resourceOrganizationsUnitRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getOrganizationsUnit: Query Organizations unit
	var (
		getOrganizationsUnitHttpUrl = "v1/organizations/organizational-units/{organizational_unit_id}"
		getOrganizationsUnitProduct = "organizations"
	)
	getOrganizationsUnitClient, err := cfg.NewServiceClient(getOrganizationsUnitProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	getOrganizationsUnitPath := getOrganizationsUnitClient.Endpoint + getOrganizationsUnitHttpUrl
	getOrganizationsUnitPath = strings.ReplaceAll(getOrganizationsUnitPath, "{organizational_unit_id}",
		fmt.Sprintf("%v", d.Id()))

	getOrganizationsUnitOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getOrganizationsUnitResp, err := getOrganizationsUnitClient.Request("GET",
		getOrganizationsUnitPath, &getOrganizationsUnitOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Organizations")
	}

	getOrganizationsUnitRespBody, err := utils.FlattenResponse(getOrganizationsUnitResp)
	if err != nil {
		return diag.FromErr(err)
	}

	organizationsUnit := utils.PathSearch("organizational_unit", getOrganizationsUnitRespBody, nil)
	if organizationsUnit == nil {
		log.Printf("[WARN] failed to get organizationsUnit, organizational_unit %s is not found in API response", d.Id())
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", organizationsUnit, nil)),
		d.Set("urn", utils.PathSearch("urn", organizationsUnit, nil)),
		d.Set("created_at", utils.PathSearch("created_at", organizationsUnit, nil)),
	)

	tagMap, err := getTags(getOrganizationsUnitClient, unitType, d.Id())
	if err != nil {
		log.Printf("[WARN] error fetching tags of Organization unit (%s): %s", d.Id(), err)
	} else {
		mErr = multierror.Append(mErr, d.Set("tags", tagMap))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceOrganizationsUnitUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updateOrganizationsUnit: update Organizations unit
	var (
		updateOrganizationsUnitHttpUrl = "v1/organizations/organizational-units/{organizational_unit_id}"
		updateOrganizationsUnitProduct = "organizations"
	)
	updateOrganizationsUnitClient, err := cfg.NewServiceClient(updateOrganizationsUnitProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	updateOrganizationsUnitHasChanges := []string{
		"name",
	}

	if d.HasChanges(updateOrganizationsUnitHasChanges...) {
		updateOrganizationsUnitPath := updateOrganizationsUnitClient.Endpoint + updateOrganizationsUnitHttpUrl
		updateOrganizationsUnitPath = strings.ReplaceAll(updateOrganizationsUnitPath, "{organizational_unit_id}",
			fmt.Sprintf("%v", d.Id()))

		updateOrganizationsUnitOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateOrganizationsUnitOpt.JSONBody = utils.RemoveNil(buildUpdateOrganizationsUnitBodyParams(d))
		_, err = updateOrganizationsUnitClient.Request("PATCH",
			updateOrganizationsUnitPath, &updateOrganizationsUnitOpt)
		if err != nil {
			return diag.Errorf("error updating Organizations: %s", err)
		}
	}

	if d.HasChange("tags") {
		err = updateTags(d, updateOrganizationsUnitClient, unitType, d.Id(), "tags")
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceOrganizationsUnitRead(ctx, d, meta)
}

func buildUpdateOrganizationsUnitBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name": utils.ValueIngoreEmpty(d.Get("name")),
	}
	return bodyParams
}

func resourceOrganizationsUnitDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteOrganizationsUnit: Delete Organizations unit
	var (
		deleteOrganizationsUnitHttpUrl = "v1/organizations/organizational-units/{organizational_unit_id}"
		deleteOrganizationsUnitProduct = "organizations"
	)
	deleteOrganizationsUnitClient, err := cfg.NewServiceClient(deleteOrganizationsUnitProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	deleteOrganizationsUnitPath := deleteOrganizationsUnitClient.Endpoint + deleteOrganizationsUnitHttpUrl
	deleteOrganizationsUnitPath = strings.ReplaceAll(deleteOrganizationsUnitPath, "{organizational_unit_id}",
		fmt.Sprintf("%v", d.Id()))

	deleteOrganizationsUnitOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = deleteOrganizationsUnitClient.Request("DELETE", deleteOrganizationsUnitPath, &deleteOrganizationsUnitOpt)
	if err != nil {
		return diag.Errorf("error deleting Organizations: %s", err)
	}

	return nil
}
