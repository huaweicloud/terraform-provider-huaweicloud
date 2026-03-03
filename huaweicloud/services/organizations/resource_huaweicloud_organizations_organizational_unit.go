package organizations

import (
	"context"
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
			// Required parameters.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the organizational unit.`,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `The ID of the root or organizational unit in which you
want to create a new organizational unit.`,
			},
			// Optional parameters.
			"tags": common.TagsSchema(`The key/value pairs associated with the organizational unit.`),
			// Attributes.
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uniform resource name of the organizational unit.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the OU was created.`,
			},
		},
	}
}

func resourceOrganizationalUnitCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/organizations/organizational-units"
	)
	client, err := cfg.NewServiceClient("organizations", region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrganizationalUnitBodyParams(d)),
	}
	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating organizational unit: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	unitId := utils.PathSearch("organizational_unit.id", respBody, "").(string)
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

func GetOrganizationalUnit(client *golangsdk.ServiceClient, ouId string) (interface{}, error) {
	httpUrl := "v1/organizations/organizational-units/{organizational_unit_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{organizational_unit_id}", ouId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	organizationalUnit := utils.PathSearch("organizational_unit", respBody, nil)
	if organizationalUnit == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/organizations/organizational-units/{organizational_unit_id}",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the organizational unit (%s) does not exist", ouId)),
			},
		}
	}

	return organizationalUnit, nil
}

func resourceOrganizationalUnitRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		ouId   = d.Id()
		mErr   *multierror.Error
	)

	client, err := cfg.NewServiceClient("organizations", region)
	if err != nil {
		return diag.Errorf("error creating client: %s", err)
	}

	organizationalUnit, err := GetOrganizationalUnit(client, ouId)
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
			fmt.Sprintf("error retrieving organizational unit (%s)", ouId),
		)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("name", utils.PathSearch("name", organizationalUnit, nil)),
		d.Set("urn", utils.PathSearch("urn", organizationalUnit, nil)),
		d.Set("created_at", utils.PathSearch("created_at", organizationalUnit, nil)),
	)

	tagMap, err := getTags(client, unitType, ouId)
	if err != nil {
		log.Printf("[WARN] error fetching tags of organizational unit (%s): %s", ouId, err)
	} else {
		mErr = multierror.Append(mErr, d.Set("tags", tagMap))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceOrganizationalUnitUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/organizations/organizational-units/{organizational_unit_id}"
		ouId    = d.Id()
	)
	client, err := cfg.NewServiceClient("organizations", region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	if d.HasChange("name") {
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{organizational_unit_id}", ouId)
		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateOrganizationalUnitBodyParams(d)),
		}

		_, err = client.Request("PATCH", updatePath, &opt)
		if err != nil {
			return diag.Errorf("error updating organizational unit (%s): %s", ouId, err)
		}
	}

	if d.HasChange("tags") {
		err = updateTags(d, client, unitType, ouId, "tags")
		if err != nil {
			return diag.Errorf("error updating tags of organizational unit (%s): %s", ouId, err)
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
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/organizations/organizational-units/{organizational_unit_id}"
		ouId    = d.Id()
	)
	client, err := cfg.NewServiceClient("organizations", region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{organizational_unit_id}", ouId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &opt)
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
			fmt.Sprintf("error deleting organizational unit (%s)", ouId),
		)
	}

	return nil
}
