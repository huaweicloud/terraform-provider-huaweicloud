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

// @API IdentityStore DELETE /v1/identity-stores/{identity_store_id}/group-memberships/{membership_id}
// @API IdentityStore GET /v1/identity-stores/{identity_store_id}/group-memberships/{membership_id}
// @API IdentityStore POST /v1/identity-stores/{identity_store_id}/group-memberships
func ResourceGroupMembership() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupMembershipCreate,
		ReadContext:   resourceGroupMembershipRead,
		DeleteContext: resourceGroupMembershipDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGroupMembershipImportState,
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
				Description: `Specifies the ID of the identity store.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the user group.`,
			},
			"member_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the user.`,
			},
		},
	}
}

func resourceGroupMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createGroupMembership: create IdentityCenter group membership
	var (
		createGroupMembershipHttpUrl = "v1/identity-stores/{identity_store_id}/group-memberships"
		createGroupMembershipProduct = "identitystore"
	)
	createGroupMembershipClient, err := cfg.NewServiceClient(createGroupMembershipProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	createGroupMembershipPath := createGroupMembershipClient.Endpoint + createGroupMembershipHttpUrl
	createGroupMembershipPath = strings.ReplaceAll(createGroupMembershipPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))

	createGroupMembershipOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createGroupMembershipOpt.JSONBody = utils.RemoveNil(buildCreateGroupMembershipBodyParams(d))
	createGroupMembershipResp, err := createGroupMembershipClient.Request("POST",
		createGroupMembershipPath, &createGroupMembershipOpt)
	if err != nil {
		return diag.Errorf("error creating Identity Center Group Membership: %s", err)
	}

	createGroupMembershipRespBody, err := utils.FlattenResponse(createGroupMembershipResp)
	if err != nil {
		return diag.FromErr(err)
	}

	membershipId := utils.PathSearch("membership_id", createGroupMembershipRespBody, "").(string)
	if membershipId == "" {
		return diag.Errorf("unable to find the membership ID of the Identity Center group from the API response")
	}
	d.SetId(membershipId)

	return resourceGroupMembershipRead(ctx, d, meta)
}

func buildCreateGroupMembershipBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"group_id":  utils.ValueIgnoreEmpty(d.Get("group_id")),
		"member_id": buildCreateGroupMembershipMemberIdChildBody(d),
	}
	return bodyParams
}

func buildCreateGroupMembershipMemberIdChildBody(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"user_id": utils.ValueIgnoreEmpty(d.Get("member_id")),
	}
	return params
}

func resourceGroupMembershipRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getGroupMembership: query Identity Center group membership
	var (
		getGroupMembershipHttpUrl = "v1/identity-stores/{identity_store_id}/group-memberships/{membership_id}"
		getGroupMembershipProduct = "identitystore"
	)
	getGroupMembershipClient, err := cfg.NewServiceClient(getGroupMembershipProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	getGroupMembershipPath := getGroupMembershipClient.Endpoint + getGroupMembershipHttpUrl
	getGroupMembershipPath = strings.ReplaceAll(getGroupMembershipPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))
	getGroupMembershipPath = strings.ReplaceAll(getGroupMembershipPath, "{membership_id}", d.Id())

	getGroupMembershipOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getGroupMembershipResp, err := getGroupMembershipClient.Request("GET", getGroupMembershipPath,
		&getGroupMembershipOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center Group Membership")
	}

	getGroupMembershipRespBody, err := utils.FlattenResponse(getGroupMembershipResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("group_id", utils.PathSearch("group_id", getGroupMembershipRespBody, nil)),
		d.Set("member_id", utils.PathSearch("member_id.user_id", getGroupMembershipRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceGroupMembershipDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteGroupMembership: delete Identity Center group membership
	var (
		deleteGroupMembershipHttpUrl = "v1/identity-stores/{identity_store_id}/group-memberships/{membership_id}"
		deleteGroupMembershipProduct = "identitystore"
	)
	deleteGroupMembershipClient, err := cfg.NewServiceClient(deleteGroupMembershipProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	deleteGroupMembershipPath := deleteGroupMembershipClient.Endpoint + deleteGroupMembershipHttpUrl
	deleteGroupMembershipPath = strings.ReplaceAll(deleteGroupMembershipPath, "{identity_store_id}",
		fmt.Sprintf("%v", d.Get("identity_store_id")))
	deleteGroupMembershipPath = strings.ReplaceAll(deleteGroupMembershipPath, "{membership_id}", d.Id())

	deleteGroupMembershipOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteGroupMembershipClient.Request("DELETE", deleteGroupMembershipPath,
		&deleteGroupMembershipOpt)
	if err != nil {
		return diag.Errorf("error deleting Identity Center Group Membership: %s", err)
	}

	return nil
}

func resourceGroupMembershipImportState(_ context.Context, d *schema.ResourceData,
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
