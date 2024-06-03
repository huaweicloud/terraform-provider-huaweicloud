// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product IdentityCenter
// ---------------------------------------------------------------

package identitycenter

import (
	"context"
	"encoding/json"
	"fmt"
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

// @API IdentityStore POST /v1/identity-stores/{identity_store_id}/users
// @API IdentityStore GET /v1/identity-stores/{identity_store_id}/users/{user_id}
// @API IdentityStore PUT /v1/identity-stores/{identity_store_id}/users/{user_id}
// @API IdentityStore DELETE /v1/identity-stores/{identity_store_id}/users/{user_id}
func ResourceIdentityCenterUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterUserCreate,
		UpdateContext: resourceIdentityCenterUserUpdate,
		ReadContext:   resourceIdentityCenterUserRead,
		DeleteContext: resourceIdentityCenterUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityCenterUserImportState,
		},

		Description: "schema: Internal",
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
				Description: `Specifies the ID of the identity store`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the username of the user.`,
			},
			"password_mode": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the initialize password mode.`,
			},
			"family_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the family name of the user.`,
			},
			"given_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the given name of the user.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the display name of the user.`,
			},
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the email of the user.`,
			},
		},
	}
}

func resourceIdentityCenterUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createIdentityCenterUser: create IdentityCenter user
	var (
		createIdentityCenterUserHttpUrl = "v1/identity-stores/{identity_store_id}/users"
		createIdentityCenterUserProduct = "identitystore"
	)
	createIdentityCenterUserClient, err := cfg.NewServiceClient(createIdentityCenterUserProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	createIdentityCenterUserPath := createIdentityCenterUserClient.Endpoint + createIdentityCenterUserHttpUrl
	createIdentityCenterUserPath = strings.ReplaceAll(createIdentityCenterUserPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))

	createIdentityCenterUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createIdentityCenterUserOpt.JSONBody = utils.RemoveNil(buildCreateIdentityCenterUserBodyParams(d))
	createIdentityCenterUserResp, err := createIdentityCenterUserClient.Request("POST",
		createIdentityCenterUserPath, &createIdentityCenterUserOpt)
	if err != nil {
		return diag.Errorf("error creating Identity Center User: %s", err)
	}

	createIdentityCenterUserRespBody, err := utils.FlattenResponse(createIdentityCenterUserResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("user_id", createIdentityCenterUserRespBody)
	if err != nil {
		return diag.Errorf("error creating Identity Center User: ID is not found in API response")
	}
	d.SetId(id.(string))

	return resourceIdentityCenterUserRead(ctx, d, meta)
}

func buildCreateIdentityCenterUserBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"user_name":     utils.ValueIgnoreEmpty(d.Get("user_name")),
		"password_mode": utils.ValueIgnoreEmpty(d.Get("password_mode")),
		"display_name":  utils.ValueIgnoreEmpty(d.Get("display_name")),
		"emails": []map[string]interface{}{{
			"primary": true,
			"type":    "Work",
			"value":   utils.ValueIgnoreEmpty(d.Get("email")),
		}},
		"name": map[string]interface{}{
			"family_name": utils.ValueIgnoreEmpty(d.Get("family_name")),
			"given_name":  utils.ValueIgnoreEmpty(d.Get("given_name")),
		},
	}
	return bodyParams
}

func resourceIdentityCenterUserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getIdentityCenterUser: Query Identity Center user
	var (
		getIdentityCenterUserHttpUrl = "v1/identity-stores/{identity_store_id}/users/{user_id}"
		getIdentityCenterUserProduct = "identitystore"
	)
	getIdentityCenterUserClient, err := cfg.NewServiceClient(getIdentityCenterUserProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	getIdentityCenterUserPath := getIdentityCenterUserClient.Endpoint + getIdentityCenterUserHttpUrl
	getIdentityCenterUserPath = strings.ReplaceAll(getIdentityCenterUserPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))
	getIdentityCenterUserPath = strings.ReplaceAll(getIdentityCenterUserPath, "{user_id}", d.Id())

	getIdentityCenterUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getIdentityCenterUserResp, err := getIdentityCenterUserClient.Request("GET", getIdentityCenterUserPath,
		&getIdentityCenterUserOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center User")
	}

	getIdentityCenterUserRespBody, err := utils.FlattenResponse(getIdentityCenterUserResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("user_name", utils.PathSearch("user_name", getIdentityCenterUserRespBody, nil)),
		d.Set("family_name", utils.PathSearch("name.family_name", getIdentityCenterUserRespBody, nil)),
		d.Set("given_name", utils.PathSearch("name.given_name", getIdentityCenterUserRespBody, nil)),
		d.Set("display_name", utils.PathSearch("display_name", getIdentityCenterUserRespBody, nil)),
		d.Set("email", utils.PathSearch("emails|[0].value", getIdentityCenterUserRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceIdentityCenterUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updateIdentityCenterUser: update Identity Center user
	var (
		updateIdentityCenterUserHttpUrl = "v1/identity-stores/{identity_store_id}/users/{user_id}"
		updateIdentityCenterUserProduct = "identitystore"
	)
	updateIdentityCenterUserClient, err := cfg.NewServiceClient(updateIdentityCenterUserProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	updateIdentityCenterUserPath := updateIdentityCenterUserClient.Endpoint + updateIdentityCenterUserHttpUrl
	updateIdentityCenterUserPath = strings.ReplaceAll(updateIdentityCenterUserPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))
	updateIdentityCenterUserPath = strings.ReplaceAll(updateIdentityCenterUserPath, "{user_id}", d.Id())

	updateIdentityCenterUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateIdentityCenterUserOpt.JSONBody = utils.RemoveNil(buildUpdateIdentityCenterUserBodyParams(d))
	_, err = updateIdentityCenterUserClient.Request("PUT", updateIdentityCenterUserPath,
		&updateIdentityCenterUserOpt)
	if err != nil {
		return diag.Errorf("error updating Identity Center User: %s", err)
	}
	return resourceIdentityCenterUserRead(ctx, d, meta)
}

func buildUpdateIdentityCenterUserBodyParams(d *schema.ResourceData) map[string]interface{} {
	operations := make([]map[string]interface{}, 0)
	if d.HasChanges("family_name", "given_name") {
		updateValue := map[string]interface{}{
			"family_name": utils.ValueIgnoreEmpty(d.Get("family_name")),
			"given_name":  utils.ValueIgnoreEmpty(d.Get("given_name")),
		}
		updateValueJson, _ := json.Marshal(updateValue)
		operations = append(operations, map[string]interface{}{
			"attribute_path":  "name",
			"attribute_value": string(updateValueJson),
		})
	}
	if d.HasChange("display_name") {
		operations = append(operations, map[string]interface{}{
			"attribute_path":  "display_name",
			"attribute_value": d.Get("display_name"),
		})
	}
	if d.HasChange("email") {
		updateValue := []map[string]interface{}{{
			"primary": true,
			"type":    "Work",
			"value":   utils.ValueIgnoreEmpty(d.Get("email")),
		}}
		updateValueJson, _ := json.Marshal(updateValue)
		operations = append(operations, map[string]interface{}{
			"attribute_path":  "emails",
			"attribute_value": string(updateValueJson),
		})
	}
	return map[string]interface{}{"operations": operations}
}

func resourceIdentityCenterUserDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteIdentityCenterUser: Delete Identity Center user
	var (
		deleteIdentityCenterUserHttpUrl = "v1/identity-stores/{identity_store_id}/users/{user_id}"
		deleteIdentityCenterUserProduct = "identitystore"
	)
	deleteIdentityCenterUserClient, err := cfg.NewServiceClient(deleteIdentityCenterUserProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	deleteIdentityCenterUserPath := deleteIdentityCenterUserClient.Endpoint + deleteIdentityCenterUserHttpUrl
	deleteIdentityCenterUserPath = strings.ReplaceAll(deleteIdentityCenterUserPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))
	deleteIdentityCenterUserPath = strings.ReplaceAll(deleteIdentityCenterUserPath, "{user_id}", d.Id())

	deleteIdentityCenterUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteIdentityCenterUserClient.Request("DELETE", deleteIdentityCenterUserPath,
		&deleteIdentityCenterUserOpt)
	if err != nil {
		return diag.Errorf("error deleting Identity Center User: %s", err)
	}

	return nil
}

func resourceIdentityCenterUserImportState(_ context.Context, d *schema.ResourceData,
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
