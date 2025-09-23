// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product Organizations
// ---------------------------------------------------------------

package organizations

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

// @API Organizations POST /v1/organizations/accounts/{account_id}/remove
// @API Organizations POST /v1/organizations/accounts/invite
// @API Organizations GET /v1/organizations/accounts/{account_id}
// @API Organizations GET /v1/organizations/handshakes/{handshake_id}
// @API Organizations POST /v1/organizations/handshakes/{handshake_id}/cancel
func ResourceAccountInvite() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccountInviteCreate,
		UpdateContext: resourceAccountInviteUpdate,
		ReadContext:   resourceAccountInviteRead,
		DeleteContext: resourceAccountInviteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the target account.`,
			},
			"remove_account_on_destroy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to remove the invited account when delete the invitation (handshake).`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the uniform resource name of the invitation`,
			},
			"master_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the unique ID of the organization's management account.`,
			},
			"master_account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the organization's management account.`,
			},
			"organization_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the organization.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the date and time when an invitation (handshake) request was made.`,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Indicates the date and time when an invitation (handshake) request was accepted,
canceled, declined, or expired.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the current state of the invitation (handshake).`,
			},
		},
	}
}

func resourceAccountInviteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createAccountInvite: create Organizations account invite
	var (
		createAccountInviteHttpUrl = "v1/organizations/accounts/invite"
		createAccountInviteProduct = "organizations"
	)
	createAccountInviteClient, err := cfg.NewServiceClient(createAccountInviteProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	createAccountInvitePath := createAccountInviteClient.Endpoint + createAccountInviteHttpUrl

	createAccountInviteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createAccountInviteOpt.JSONBody = buildCreateAccountInviteBodyParams(d)
	createAccountInviteResp, err := createAccountInviteClient.Request("POST", createAccountInvitePath,
		&createAccountInviteOpt)
	if err != nil {
		return diag.Errorf("error creating AccountInvite: %s", err)
	}

	createAccountInviteRespBody, err := utils.FlattenResponse(createAccountInviteResp)
	if err != nil {
		return diag.FromErr(err)
	}

	handshakeId := utils.PathSearch("handshake.id", createAccountInviteRespBody, "").(string)
	if handshakeId == "" {
		return diag.Errorf("unable to find the Organizations account invite ID from the API response")
	}
	d.SetId(handshakeId)

	return resourceAccountInviteRead(ctx, d, meta)
}

func buildCreateAccountInviteBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"target": buildCreateAccountInviteTargetChildBody(d),
		"notes":  "",
	}
	return bodyParams
}

func buildCreateAccountInviteTargetChildBody(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"type":   "account",
		"entity": utils.ValueIgnoreEmpty(d.Get("account_id")),
	}
	return params
}

func resourceAccountInviteRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getAccountInvite: Query Organizations account invite
	getAccountInviteProduct := "organizations"
	getAccountInviteClient, err := cfg.NewServiceClient(getAccountInviteProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	getAccountInviteRespBody, err := getAccountInvite(getAccountInviteClient, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving AccountInvite")
	}

	handshake := utils.PathSearch("handshake", getAccountInviteRespBody, nil)
	if handshake == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	status := utils.PathSearch("status", handshake, "")
	if status == "cancelled" || status == "expired" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	accountID := utils.PathSearch("target.entity", handshake, "").(string)

	// the handshake will always exist, so it is necessary to check whether the account can be obtained if the
	// status is accepted
	if status == "accepted" {
		getAccountHttpUrl := "v1/organizations/accounts/{account_id}"
		getAccountPath := getAccountInviteClient.Endpoint + getAccountHttpUrl
		getAccountPath = strings.ReplaceAll(getAccountPath, "{account_id}", accountID)

		getAccountOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		_, err = getAccountInviteClient.Request("GET", getAccountPath, &getAccountOpt)

		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving AccountInvite")
		}
	}

	mErr = multierror.Append(
		mErr,
		d.Set("urn", utils.PathSearch("urn", handshake, nil)),
		d.Set("account_id", accountID),
		d.Set("master_account_id", utils.PathSearch("management_account_id", handshake, nil)),
		d.Set("master_account_name", utils.PathSearch("management_account_name", handshake, nil)),
		d.Set("organization_id", utils.PathSearch("organization_id", handshake, nil)),
		d.Set("created_at", utils.PathSearch("created_at", handshake, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", handshake, nil)),
		d.Set("status", status),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getAccountInvite(client *golangsdk.ServiceClient, handshakeId string) (interface{}, error) {
	// getAccountInvite: Query Organizations account invite
	var (
		getAccountInviteHttpUrl = "v1/organizations/handshakes/{handshake_id}"
	)

	getAccountInvitePath := client.Endpoint + getAccountInviteHttpUrl
	getAccountInvitePath = strings.ReplaceAll(getAccountInvitePath, "{handshake_id}", handshakeId)

	getAccountInviteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getAccountInviteResp, err := client.Request("GET", getAccountInvitePath, &getAccountInviteOpt)

	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getAccountInviteResp)
}

func resourceAccountInviteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAccountInviteDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteAccountInvite: Delete Organizations account invite
	var (
		deleteAccountInviteProduct = "organizations"
	)
	deleteAccountInviteClient, err := cfg.NewServiceClient(deleteAccountInviteProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	status := d.Get("status")
	if status == "pending" {
		err = cancelAccountInvite(deleteAccountInviteClient, d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
	} else if status == "accepted" && d.Get("remove_account_on_destroy").(bool) {
		err = removeAccountInvite(deleteAccountInviteClient, d.Get("account_id").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func cancelAccountInvite(client *golangsdk.ServiceClient, handshakeId string) error {
	// cancelAccountInvite: cancel Organizations account invite
	var (
		cancelAccountInviteHttpUrl = "v1/organizations/handshakes/{handshake_id}/cancel"
	)
	cancelAccountInvitePath := client.Endpoint + cancelAccountInviteHttpUrl
	cancelAccountInvitePath = strings.ReplaceAll(cancelAccountInvitePath, "{handshake_id}", handshakeId)

	cancelAccountInviteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err := client.Request("POST", cancelAccountInvitePath, &cancelAccountInviteOpt)
	if err != nil {
		return fmt.Errorf("error canceling AccountInvite: %s", err)
	}
	return nil
}

func removeAccountInvite(client *golangsdk.ServiceClient, accountId string) error {
	// removeAccount: Remove Organizations account
	var (
		removeAccountHttpUrl = "v1/organizations/accounts/{account_id}/remove"
	)

	removeAccountPath := client.Endpoint + removeAccountHttpUrl
	removeAccountPath = strings.ReplaceAll(removeAccountPath, "{account_id}", accountId)

	removeAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err := client.Request("POST", removeAccountPath, &removeAccountOpt)
	if err != nil {
		return fmt.Errorf("error removing Account: %s", err)
	}
	return nil
}
