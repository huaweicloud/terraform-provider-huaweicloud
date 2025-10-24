// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product IdentityCenter
// ---------------------------------------------------------------

package identitycenter

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

// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/groups
// @API IdentityCenter DELETE /v1/identity-stores/{identity_store_id}/groups/{group_id}
// @API IdentityCenter GET /v1/identity-stores/{identity_store_id}/groups/{group_id}
// @API IdentityCenter PUT /v1/identity-stores/{identity_store_id}/groups/{group_id}
func ResourceIdentityCenterGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterGroupCreate,
		UpdateContext: resourceIdentityCenterGroupUpdate,
		ReadContext:   resourceIdentityCenterGroupRead,
		DeleteContext: resourceIdentityCenterGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityCenterGroupImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"identity_store_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the identity store.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the group.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the group.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the last update time.`,
			},
		},
	}
}

func resourceIdentityCenterGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createIdentityCenterGroup: create IdentityCenter group
	var (
		createIdentityCenterGroupHttpUrl = "v1/identity-stores/{identity_store_id}/groups"
		createIdentityCenterGroupProduct = "identitystore"
	)
	createIdentityCenterGroupClient, err := cfg.NewServiceClient(createIdentityCenterGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	createIdentityCenterGroupPath := createIdentityCenterGroupClient.Endpoint + createIdentityCenterGroupHttpUrl
	createIdentityCenterGroupPath = strings.ReplaceAll(createIdentityCenterGroupPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))

	createIdentityCenterGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createIdentityCenterGroupOpt.JSONBody = utils.RemoveNil(buildCreateIdentityCenterGroupBodyParams(d))
	createIdentityCenterGroupResp, err := createIdentityCenterGroupClient.Request("POST",
		createIdentityCenterGroupPath, &createIdentityCenterGroupOpt)
	if err != nil {
		return diag.Errorf("error creating Identity Center Group: %s", err)
	}

	createIdentityCenterGroupRespBody, err := utils.FlattenResponse(createIdentityCenterGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	groupId := utils.PathSearch("group_id", createIdentityCenterGroupRespBody, "").(string)
	if groupId == "" {
		return diag.Errorf("unable to find the Identity Center group ID from the API response")
	}
	d.SetId(groupId)

	return resourceIdentityCenterGroupRead(ctx, d, meta)
}

func buildCreateIdentityCenterGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"display_name": utils.ValueIgnoreEmpty(d.Get("name")),
		"description":  utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceIdentityCenterGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getIdentityCenterGroup: query Identity Center group
	var (
		getIdentityCenterGroupHttpUrl = "v1/identity-stores/{identity_store_id}/groups/{group_id}"
		getIdentityCenterGroupProduct = "identitystore"
	)
	getIdentityCenterGroupClient, err := cfg.NewServiceClient(getIdentityCenterGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	getIdentityCenterGroupPath := getIdentityCenterGroupClient.Endpoint + getIdentityCenterGroupHttpUrl
	getIdentityCenterGroupPath = strings.ReplaceAll(getIdentityCenterGroupPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))
	getIdentityCenterGroupPath = strings.ReplaceAll(getIdentityCenterGroupPath, "{group_id}", d.Id())

	getIdentityCenterGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getIdentityCenterGroupResp, err := getIdentityCenterGroupClient.Request("GET", getIdentityCenterGroupPath,
		&getIdentityCenterGroupOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center Group")
	}

	getIdentityCenterGroupRespBody, err := utils.FlattenResponse(getIdentityCenterGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	createAt := utils.PathSearch("created_at", getIdentityCenterGroupRespBody, float64(0)).(float64)
	updateAt := utils.PathSearch("updated_at", getIdentityCenterGroupRespBody, float64(0)).(float64)
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("display_name", getIdentityCenterGroupRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getIdentityCenterGroupRespBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(createAt)/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(int64(updateAt)/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceIdentityCenterGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updateIdentityCenterGroup: update Identity Center group
	var (
		updateIdentityCenterGroupHttpUrl = "v1/identity-stores/{identity_store_id}/groups/{group_id}"
		updateIdentityCenterGroupProduct = "identitystore"
	)
	updateIdentityCenterGroupClient, err := cfg.NewServiceClient(updateIdentityCenterGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	updateIdentityCenterGroupPath := updateIdentityCenterGroupClient.Endpoint + updateIdentityCenterGroupHttpUrl
	updateIdentityCenterGroupPath = strings.ReplaceAll(updateIdentityCenterGroupPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))
	updateIdentityCenterGroupPath = strings.ReplaceAll(updateIdentityCenterGroupPath, "{group_id}", d.Id())

	updateIdentityCenterGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateIdentityCenterGroupOpt.JSONBody = utils.RemoveNil(buildUpdateIdentityCenterGroupBodyParams(d))
	_, err = updateIdentityCenterGroupClient.Request("PUT", updateIdentityCenterGroupPath, &updateIdentityCenterGroupOpt)
	if err != nil {
		return diag.Errorf("error updating Identity Center Group: %s", err)
	}

	return resourceIdentityCenterGroupRead(ctx, d, meta)
}

func buildUpdateIdentityCenterGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	operations := make([]map[string]interface{}, 0)
	operations = append(operations, map[string]interface{}{
		"attribute_path":  "description",
		"attribute_value": d.Get("description"),
	})
	return map[string]interface{}{"operations": operations}
}

func resourceIdentityCenterGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteIdentityCenterGroup: delete Identity Center group
	var (
		deleteIdentityCenterGroupHttpUrl = "v1/identity-stores/{identity_store_id}/groups/{group_id}"
		deleteIdentityCenterGroupProduct = "identitystore"
	)
	deleteIdentityCenterGroupClient, err := cfg.NewServiceClient(deleteIdentityCenterGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	deleteIdentityCenterGroupPath := deleteIdentityCenterGroupClient.Endpoint + deleteIdentityCenterGroupHttpUrl
	deleteIdentityCenterGroupPath = strings.ReplaceAll(deleteIdentityCenterGroupPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))
	deleteIdentityCenterGroupPath = strings.ReplaceAll(deleteIdentityCenterGroupPath, "{group_id}", d.Id())

	deleteIdentityCenterGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteIdentityCenterGroupClient.Request("DELETE", deleteIdentityCenterGroupPath,
		&deleteIdentityCenterGroupOpt)
	if err != nil {
		return diag.Errorf("error deleting Identity Center Group: %s", err)
	}

	return nil
}

func resourceIdentityCenterGroupImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <identity_store_id>/<id>")
	}
	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("identity_store_id", parts[0]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
