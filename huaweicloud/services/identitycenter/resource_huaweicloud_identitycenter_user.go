package identitycenter

import (
	"context"
	"encoding/json"
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

// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/users
// @API IdentityCenter GET /v1/identity-stores/{identity_store_id}/users/{user_id}
// @API IdentityCenter PUT /v1/identity-stores/{identity_store_id}/users/{user_id}
// @API IdentityCenter DELETE /v1/identity-stores/{identity_store_id}/users/{user_id}
// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/users/{user_id}/enable
// @API IdentityCenter POST /v1/identity-stores/{identity_store_id}/users/{user_id}/disable
func ResourceIdentityCenterUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterUserCreate,
		UpdateContext: resourceIdentityCenterUserUpdate,
		ReadContext:   resourceIdentityCenterUserRead,
		DeleteContext: resourceIdentityCenterUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityCenterUserImportState,
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
			"phone_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the phone number of the user.`,
			},
			"user_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of the user.`,
			},
			"title": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the title of the user.`,
			},
			"addresses": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the addresses information of the user.`,
				MaxItems:    1,
				Elem:        identityCenterUserAddressesSchema(),
			},
			"enterprise": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the enterprise information of the user.`,
				MaxItems:    1,
				Elem:        identityCenterUserEnterpriseSchema(),
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the user.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the user.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time of the user.`,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The updater of the user.`,
			},
			"email_verified": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the email is verified.`,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether the user is enabled.`,
			},
		},
	}
}

func identityCenterUserAddressesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"country": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the country of the user.`,
			},
			"formatted": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies a string containing a formatted version of the address to be displayed.`,
			},
			"locality": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the locality of the user.`,
			},
			"postal_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the postal code of the user.`,
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the region of the user.`,
			},
			"street_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the street address of the user.`,
			},
		},
	}
	return &sc
}

func identityCenterUserEnterpriseSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"cost_center": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the cost center of the enterprise.`,
			},
			"department": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the department of the enterprise.`,
			},
			"division": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the division of the enterprise.`,
			},
			"employee_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the employee number of the enterprise.`,
			},
			"organization": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the organization of the enterprise.`,
			},
			"manager": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the manager of the enterprise.`,
			},
		},
	}
	return &sc
}

func resourceIdentityCenterUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createIdentityCenterUser: create IdentityCenter user
	var (
		createIdentityCenterUserHttpUrl = "v1/identity-stores/{identity_store_id}/users"
		disableUserHttpUrl              = "v1/identity-stores/{identity_store_id}/users/{user_id}/disable"
		createIdentityCenterUserProduct = "identitystore"
	)
	client, err := cfg.NewServiceClient(createIdentityCenterUserProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	createIdentityCenterUserPath := client.Endpoint + createIdentityCenterUserHttpUrl
	createIdentityCenterUserPath = strings.ReplaceAll(createIdentityCenterUserPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))

	createIdentityCenterUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createIdentityCenterUserOpt.JSONBody = utils.RemoveNil(buildCreateIdentityCenterUserBodyParams(d))
	createIdentityCenterUserResp, err := client.Request("POST",
		createIdentityCenterUserPath, &createIdentityCenterUserOpt)
	if err != nil {
		return diag.Errorf("error creating Identity Center User: %s", err)
	}

	createIdentityCenterUserRespBody, err := utils.FlattenResponse(createIdentityCenterUserResp)
	if err != nil {
		return diag.FromErr(err)
	}

	userId := utils.PathSearch("user_id", createIdentityCenterUserRespBody, "").(string)
	if userId == "" {
		return diag.Errorf("unable to find the Identity Center user ID from the API response")
	}
	d.SetId(userId)

	if !d.Get("enabled").(bool) {
		disableUserPath := client.Endpoint + disableUserHttpUrl
		disableUserPath = strings.ReplaceAll(disableUserPath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))
		disableUserPath = strings.ReplaceAll(disableUserPath, "{user_id}", d.Id())

		disableUserOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		_, err = client.Request("POST", disableUserPath, &disableUserOpt)
		if err != nil {
			return diag.Errorf("error disabling Identity Center User: %s", err)
		}
	}

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
		"title":     utils.ValueIgnoreEmpty(d.Get("title")),
		"user_type": utils.ValueIgnoreEmpty(d.Get("user_type")),
		"phone_numbers": []map[string]interface{}{{
			"primary": true,
			"type":    "Work",
			"value":   utils.ValueIgnoreEmpty(d.Get("phone_number")),
		}},
		"addresses": []map[string]interface{}{
			{
				"country":        utils.ValueIgnoreEmpty(d.Get("addresses.0.country")),
				"formatted":      utils.ValueIgnoreEmpty(d.Get("addresses.0.formatted")),
				"locality":       utils.ValueIgnoreEmpty(d.Get("addresses.0.locality")),
				"postal_code":    utils.ValueIgnoreEmpty(d.Get("addresses.0.postal_code")),
				"region":         utils.ValueIgnoreEmpty(d.Get("addresses.0.region")),
				"street_address": utils.ValueIgnoreEmpty(d.Get("addresses.0.street_address")),
			},
		},
		"enterprise": map[string]interface{}{
			"cost_center":     utils.ValueIgnoreEmpty(d.Get("enterprise.0.cost_center")),
			"department":      utils.ValueIgnoreEmpty(d.Get("enterprise.0.department")),
			"division":        utils.ValueIgnoreEmpty(d.Get("enterprise.0.division")),
			"employee_number": utils.ValueIgnoreEmpty(d.Get("enterprise.0.employee_number")),
			"manager": map[string]interface{}{
				"value": utils.ValueIgnoreEmpty(d.Get("enterprise.0.manager")),
			},
			"organization": utils.ValueIgnoreEmpty(d.Get("enterprise.0.organization")),
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
		d.Set("phone_number", utils.PathSearch("phone_numbers|[0].value", getIdentityCenterUserRespBody, nil)),
		d.Set("title", utils.PathSearch("title", getIdentityCenterUserRespBody, nil)),
		d.Set("user_type", utils.PathSearch("user_type", getIdentityCenterUserRespBody, nil)),
		d.Set("created_by", utils.PathSearch("created_by", getIdentityCenterUserRespBody, nil)),
		d.Set("updated_by", utils.PathSearch("updated_by", getIdentityCenterUserRespBody, nil)),
		d.Set("email_verified", utils.PathSearch("email_verified", getIdentityCenterUserRespBody, false)),
		d.Set("enabled", utils.PathSearch("enabled", getIdentityCenterUserRespBody, false)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("created_at", getIdentityCenterUserRespBody, float64(0)).(float64))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("updated_at", getIdentityCenterUserRespBody, float64(0)).(float64))/1000, false)),
		d.Set("addresses", flattenIdentityCenterUserAddresses(utils.PathSearch("addresses|[0]", getIdentityCenterUserRespBody, nil))),
		d.Set("enterprise", flattenIdentityCenterUserEnterprise(utils.PathSearch("enterprise", getIdentityCenterUserRespBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenIdentityCenterUserAddresses(address interface{}) []map[string]interface{} {
	if address == nil || len(address.(map[string]interface{})) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"country":        utils.PathSearch("country", address, nil),
			"formatted":      utils.PathSearch("formatted", address, nil),
			"locality":       utils.PathSearch("locality", address, nil),
			"postal_code":    utils.PathSearch("postal_code", address, nil),
			"region":         utils.PathSearch("region", address, nil),
			"street_address": utils.PathSearch("street_address", address, nil),
		},
	}
}

func flattenIdentityCenterUserEnterprise(enterprise interface{}) []map[string]interface{} {
	if enterprise == nil {
		return nil
	}

	// enterprise format: {"enterprise":{"manager":{}}}
	enterpriseMap := enterprise.(map[string]interface{})
	manager := utils.PathSearch("manager.value", enterprise, nil)
	if len(enterpriseMap) == 1 && manager == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"cost_center":     utils.PathSearch("cost_center", enterprise, nil),
			"department":      utils.PathSearch("department", enterprise, nil),
			"division":        utils.PathSearch("division", enterprise, nil),
			"employee_number": utils.PathSearch("employee_number", enterprise, nil),
			"organization":    utils.PathSearch("organization", enterprise, nil),
			"manager":         manager,
		},
	}
}

func resourceIdentityCenterUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updateIdentityCenterUser: update Identity Center user
	var (
		updateIdentityCenterUserHttpUrl = "v1/identity-stores/{identity_store_id}/users/{user_id}"
		enableUserHttpUrl               = "v1/identity-stores/{identity_store_id}/users/{user_id}/enable"
		disableUserHttpUrl              = "v1/identity-stores/{identity_store_id}/users/{user_id}/disable"
		updateIdentityCenterUserProduct = "identitystore"
	)
	client, err := cfg.NewServiceClient(updateIdentityCenterUserProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	updateIdentityCenterUserPath := client.Endpoint + updateIdentityCenterUserHttpUrl
	updateIdentityCenterUserPath = strings.ReplaceAll(updateIdentityCenterUserPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))
	updateIdentityCenterUserPath = strings.ReplaceAll(updateIdentityCenterUserPath, "{user_id}", d.Id())

	updateIdentityCenterUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	params := utils.RemoveNil(buildUpdateIdentityCenterUserBodyParams(d))
	if params != nil {
		ops := params["operations"]
		if ops != nil && len(ops.([]map[string]interface{})) > 0 {
			updateIdentityCenterUserOpt.JSONBody = params
			_, err = client.Request("PUT", updateIdentityCenterUserPath, &updateIdentityCenterUserOpt)
			if err != nil {
				return diag.Errorf("error updating Identity Center User: %s", err)
			}
		}
	}

	if d.HasChange("enabled") {
		var actionPath string
		var actionHttpUrl string
		var action string

		if d.Get("enabled").(bool) {
			actionHttpUrl = enableUserHttpUrl
			action = "enabling"
		} else {
			actionHttpUrl = disableUserHttpUrl
			action = "disabling"
		}

		actionPath = client.Endpoint + actionHttpUrl
		actionPath = strings.ReplaceAll(actionPath, "{identity_store_id}", fmt.Sprintf("%v", d.Get("identity_store_id")))
		actionPath = strings.ReplaceAll(actionPath, "{user_id}", d.Id())

		actionOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		_, err = client.Request("POST", actionPath, &actionOpt)
		if err != nil {
			return diag.Errorf("error %s Identity Center User: %s", action, err)
		}
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

	if d.HasChange("phone_number") {
		updateValue := []map[string]interface{}{{
			"primary": true,
			"type":    "Work",
			"value":   utils.ValueIgnoreEmpty(d.Get("phone_number")),
		}}
		updateValueJson, _ := json.Marshal(updateValue)
		operations = append(operations, map[string]interface{}{
			"attribute_path":  "phone_numbers",
			"attribute_value": string(updateValueJson),
		})
	}

	if d.HasChange("addresses") {
		updateValue := []map[string]interface{}{{
			"country":        utils.ValueIgnoreEmpty(d.Get("addresses.0.country")),
			"region":         utils.ValueIgnoreEmpty(d.Get("addresses.0.region")),
			"locality":       utils.ValueIgnoreEmpty(d.Get("addresses.0.locality")),
			"postal_code":    utils.ValueIgnoreEmpty(d.Get("addresses.0.postal_code")),
			"street_address": utils.ValueIgnoreEmpty(d.Get("addresses.0.street_address")),
			"formatted":      utils.ValueIgnoreEmpty(d.Get("addresses.0.formatted")),
		}}
		updateValueJson, _ := json.Marshal(updateValue)
		operations = append(operations, map[string]interface{}{
			"attribute_path":  "addresses",
			"attribute_value": string(updateValueJson),
		})
	}

	if d.HasChange("user_type") {
		operations = append(operations, map[string]interface{}{
			"attribute_path":  "user_type",
			"attribute_value": utils.ValueIgnoreEmpty(d.Get("user_type")),
		})
	}

	if d.HasChange("title") {
		operations = append(operations, map[string]interface{}{
			"attribute_path":  "title",
			"attribute_value": utils.ValueIgnoreEmpty(d.Get("title")),
		})
	}

	if d.HasChange("enterprise") {
		updateValue := map[string]interface{}{
			"cost_center":     utils.ValueIgnoreEmpty(d.Get("enterprise.0.cost_center")),
			"department":      utils.ValueIgnoreEmpty(d.Get("enterprise.0.department")),
			"division":        utils.ValueIgnoreEmpty(d.Get("enterprise.0.division")),
			"employee_number": utils.ValueIgnoreEmpty(d.Get("enterprise.0.employee_number")),
			"organization":    utils.ValueIgnoreEmpty(d.Get("enterprise.0.organization")),
			"manager": map[string]interface{}{
				"value": utils.ValueIgnoreEmpty(d.Get("enterprise.0.manager")),
			},
		}
		updateValueJson, _ := json.Marshal(updateValue)
		operations = append(operations, map[string]interface{}{
			"attribute_path":  "enterprise",
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
