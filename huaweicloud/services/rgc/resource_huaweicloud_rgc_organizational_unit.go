package rgc

import (
	"context"
	"errors"
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

var organizationalUnitNonUpdatableParams = []string{"organizational_unit_name", "parent_organizational_unit_id"}

// @API RGC GET /v1/managed-organization/managed-organizational-units/{managed_organizational_unit_id}
// @API RGC POST /v1/managed-organization/managed-organizational-units
// @API RGC DELETE /v1/managed-organization/managed-organizational-units/{managed_organizational_unit_id}
func ResourceOrganizationalUnit() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationalUnitCreate,
		UpdateContext: resourceOrganizationalUnitUpdate,
		ReadContext:   resourceOrganizationalUnitRead,
		DeleteContext: resourceOrganizationalUnitDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(organizationalUnitNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"organizational_unit_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_organizational_unit_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_organizational_unit_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"manage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organizational_unit_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organizational_unit_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organizational_unit_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOrganizationalUnitCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	unitId, err := createOrganizationalUnit(d, meta)
	if err != nil {
		return diag.Errorf("error creating RGC organizational unit: %s", err)
	}
	d.SetId(unitId)

	return resourceOrganizationalUnitRead(ctx, d, meta)
}

func resourceOrganizationalUnitRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getOrganizationUnitUrl     = "v1/managed-organization/managed-organizational-units/{managed_organizational_unit_id}"
		getOrganizationUnitProduct = "rgc"
	)

	getOrganizationUnitClient, err := cfg.NewServiceClient(getOrganizationUnitProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	getOrganizationUnitPath := getOrganizationUnitClient.Endpoint + getOrganizationUnitUrl
	getOrganizationUnitPath = strings.ReplaceAll(getOrganizationUnitPath, "{managed_organizational_unit_id}", d.Id())

	getOrganizationUnitOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getOrganizationUnitResp, err := getOrganizationUnitClient.Request("GET", getOrganizationUnitPath, &getOrganizationUnitOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving organizational unit")
	}

	getOrganizationUnitRespBody, err := utils.FlattenResponse(getOrganizationUnitResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("organizational_unit_id", d.Id()),
		d.Set("organizational_unit_name", utils.PathSearch("organizational_unit_name", getOrganizationUnitRespBody, nil)),
		d.Set("parent_organizational_unit_id", utils.PathSearch("parent_organizational_unit_id", getOrganizationUnitRespBody, nil)),
		d.Set("parent_organizational_unit_name", utils.PathSearch("parent_organizational_unit_name", getOrganizationUnitRespBody, nil)),
		d.Set("manage_account_id", utils.PathSearch("manage_account_id", getOrganizationUnitRespBody, nil)),
		d.Set("organizational_unit_type", utils.PathSearch("organizational_unit_type", getOrganizationUnitRespBody, nil)),
		d.Set("organizational_unit_status", utils.PathSearch("organizational_unit_status", getOrganizationUnitRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceOrganizationalUnitUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceOrganizationalUnitDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := deleteOrganizationalUnit(ctx, d, meta)
	if err != nil {
		return diag.Errorf("error deleting organizational unit: %s", err)
	}
	return nil
}

func createOrganizationalUnit(d *schema.ResourceData, meta interface{}) (string, error) {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createOrganizationalUnitUrl     = "v1/managed-organization/managed-organizational-units"
		createOrganizationalUnitProduct = "rgc"
	)
	createOrganizationalUnitClient, err := cfg.NewServiceClient(createOrganizationalUnitProduct, region)
	if err != nil {
		return "", fmt.Errorf("error creating RGC client: %s", err)
	}

	createOrganizationalUnitPath := createOrganizationalUnitClient.Endpoint + createOrganizationalUnitUrl

	createOrganizationalUnitOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"name":      d.Get("organizational_unit_name").(string),
			"parent_id": d.Get("parent_organizational_unit_id").(string),
		},
		OkCodes: []int{
			201,
		},
	}
	createOrganizationalUnitResp, err := createOrganizationalUnitClient.Request("POST",
		createOrganizationalUnitPath, &createOrganizationalUnitOpt)
	if err != nil {
		return "", err
	}

	createOrganizationalUnitRespBody, err := utils.FlattenResponse(createOrganizationalUnitResp)
	if err != nil {
		return "", err
	}

	unitId := utils.PathSearch("organizational_unit_id", createOrganizationalUnitRespBody, "").(string)
	if unitId == "" {
		return "", errors.New("unable to find the organizational unit ID from the API response")
	}

	return unitId, nil
}

func deleteOrganizationalUnit(_ context.Context, d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteOrganizationalUnitUrl     = "v1/managed-organization/managed-organizational-units/{managed_organizational_unit_id}"
		deleteOrganizationalUnitProduct = "rgc"
	)
	deleteOrganizationalUnitClient, err := cfg.NewServiceClient(deleteOrganizationalUnitProduct, region)
	if err != nil {
		return fmt.Errorf("error creating RGC client: %n", err)
	}

	deleteOrganizationalUnitPath := deleteOrganizationalUnitClient.Endpoint + deleteOrganizationalUnitUrl
	deleteOrganizationalUnitPath = strings.ReplaceAll(deleteOrganizationalUnitPath, "{managed_organizational_unit_id}", d.Id())
	deleteOrganizationalUnitOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = deleteOrganizationalUnitClient.Request("DELETE",
		deleteOrganizationalUnitPath, &deleteOrganizationalUnitOpt)
	if err != nil {
		return err
	}

	return nil
}
