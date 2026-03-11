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

var (
	accountInviteNotFoundErrCodes = []string{
		"Organizations.1401", // The invitation has been canceled.
	}
)

// @API Organizations POST /v1/organizations/accounts/invite
// @API Organizations GET /v1/organizations/handshakes/{handshake_id}
// @API Organizations GET /v1/organizations/accounts/{account_id}
// @API Organizations POST /v1/organizations/handshakes/{handshake_id}/cancel
// @API Organizations POST /v1/organizations/accounts/{account_id}/remove
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
				Description: `The ID of the target account.`,
			},
			"remove_account_on_destroy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to remove the invited account when delete the invitation (handshake).`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uniform resource name of the invitation`,
			},
			"master_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unique ID of the organization's management account.`,
			},
			"master_account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the organization's management account.`,
			},
			"organization_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the organization.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The date and time when an invitation (handshake) request was made.`,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The date and time when an invitation (handshake) request was accepted,
canceled, declined, or expired.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current state of the invitation (handshake).`,
			},
		},
	}
}

func resourceAccountInviteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1/organizations/accounts/invite"
	)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateAccountInviteBodyParams(d),
	}
	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error inviting the account (%s) to join the organization: %s", d.Get("account_id").(string), err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	handshakeId := utils.PathSearch("handshake.id", respBody, "").(string)
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
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	handshake, err := GetAccountInvite(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
			"error retrieving account invitation",
		)
	}

	mErr := multierror.Append(
		d.Set("urn", utils.PathSearch("urn", handshake, nil)),
		d.Set("account_id", utils.PathSearch("target.entity", handshake, "").(string)),
		d.Set("master_account_id", utils.PathSearch("management_account_id", handshake, nil)),
		d.Set("master_account_name", utils.PathSearch("management_account_name", handshake, nil)),
		d.Set("organization_id", utils.PathSearch("organization_id", handshake, nil)),
		d.Set("created_at", utils.PathSearch("created_at", handshake, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", handshake, nil)),
		d.Set("status", utils.PathSearch("status", handshake, "").(string)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetAccountInviteById(client *golangsdk.ServiceClient, handshakeId string) (interface{}, error) {
	httpUrl := "v1/organizations/handshakes/{handshake_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{handshake_id}", handshakeId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	resp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	handshake := utils.PathSearch("handshake", respBody, nil)
	if handshake == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/organizations/handshakes/{handshake_id}",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the invitation (%s) does not exist", handshakeId)),
			},
		}
	}

	return handshake, nil
}

func GetAccountInvite(client *golangsdk.ServiceClient, handshakeId string) (interface{}, error) {
	handshake, err := GetAccountInviteById(client, handshakeId)
	if err != nil {
		return nil, err
	}

	status := utils.PathSearch("status", handshake, "").(string)
	if status == "cancelled" || status == "expired" {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/organizations/handshakes/{handshake_id}",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the invitation (%s) has been cancelled or expired", handshakeId)),
			},
		}
	}

	// the handshake will always exist, so it is necessary to check whether the account can be obtained if the
	// status is accepted.
	if status == "accepted" {
		_, err = GetAccountById(client, utils.PathSearch("target.entity", handshake, "").(string))
		if err != nil {
			return nil, err
		}
	}

	return handshake, nil
}

func GetAccountById(client *golangsdk.ServiceClient, accountId string) (interface{}, error) {
	httpUrl := "v1/organizations/accounts/{account_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{account_id}", accountId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	resp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceAccountInviteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAccountInviteDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	status := d.Get("status")
	if status == "pending" {
		invitationId := d.Id()
		if err = cancelAccountInvite(client, invitationId); err != nil {
			return common.CheckDeletedDiag(
				d,
				common.ConvertExpected400ErrInto404Err(
					common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
					"error_code",
					accountInviteNotFoundErrCodes...,
				),
				fmt.Sprintf("error canceling account invitation (%s)", invitationId),
			)
		}
	} else if status == "accepted" && d.Get("remove_account_on_destroy").(bool) {
		accountId := d.Get("account_id").(string)
		if err = removeAccountInvite(client, accountId); err != nil {
			return common.CheckDeletedDiag(d,
				common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
				fmt.Sprintf("error removing account (%s) from the organization", accountId))
		}
	}

	return nil
}

func cancelAccountInvite(client *golangsdk.ServiceClient, handshakeId string) error {
	httpUrl := "v1/organizations/handshakes/{handshake_id}/cancel"
	cancelPath := client.Endpoint + httpUrl
	cancelPath = strings.ReplaceAll(cancelPath, "{handshake_id}", handshakeId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err := client.Request("POST", cancelPath, &opt)
	return err
}

func removeAccountInvite(client *golangsdk.ServiceClient, accountId string) error {
	httpUrl := "v1/organizations/accounts/{account_id}/remove"
	removePath := client.Endpoint + httpUrl
	removePath = strings.ReplaceAll(removePath, "{account_id}", accountId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err := client.Request("POST", removePath, &opt)
	return err
}
