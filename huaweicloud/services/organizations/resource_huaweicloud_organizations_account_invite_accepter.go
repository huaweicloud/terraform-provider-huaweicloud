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

// @API Organizations GET /v1/organizations/handshakes/{handshake_id}
// @API Organizations POST /v1/organizations/leave
// @API Organizations POST /v1/received-handshakes/{handshake_id}/accept
func ResourceAccountInviteAccepter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccountInviteAccepterCreate,
		UpdateContext: resourceAccountInviteAccepterUpdate,
		ReadContext:   resourceAccountInviteAccepterRead,
		DeleteContext: resourceAccountInviteAccepterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"invitation_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the unique ID of an invitation (handshake).`,
			},
			"leave_organization_on_destroy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to leave the organization when delete the invitation (handshake).`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the uniform resource name of the invitation`,
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the target account.`,
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
cancelled, declined, or expired.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the current state of the invitation (handshake).`,
			},
		},
	}
}

func resourceAccountInviteAccepterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createAccountInviteAccepter: create Organizations account invite accepter
	var (
		createAccountInviteAccepterHttpUrl = "v1/received-handshakes/{handshake_id}/accept"
		createAccountInviteAccepterProduct = "organizations"
	)
	createAccountInviteAccepterClient, err := cfg.NewServiceClient(createAccountInviteAccepterProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	createAccountInviteAccepterPath := createAccountInviteAccepterClient.Endpoint + createAccountInviteAccepterHttpUrl

	createAccountInviteAccepterPath = strings.ReplaceAll(createAccountInviteAccepterPath, "{handshake_id}",
		fmt.Sprintf("%v", d.Get("invitation_id")))

	createAccountInviteAccepterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createAccountInviteAccepterResp, err := createAccountInviteAccepterClient.Request("POST",
		createAccountInviteAccepterPath, &createAccountInviteAccepterOpt)
	if err != nil {
		return diag.Errorf("error creating AccountInviteAccepter: %s", err)
	}

	createAccountInviteAccepterRespBody, err := utils.FlattenResponse(createAccountInviteAccepterResp)
	if err != nil {
		return diag.FromErr(err)
	}

	handshakeId := utils.PathSearch("handshake.id", createAccountInviteAccepterRespBody, "").(string)
	if handshakeId == "" {
		return diag.Errorf("unable to find the handshake ID of the Organizations account invite accepter from the API response")
	}
	d.SetId(handshakeId)

	return resourceAccountInviteAccepterRead(ctx, d, meta)
}

func resourceAccountInviteAccepterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getAccountInviteAccepter: Query Organizations account invite accepter
	var (
		getAccountInviteAccepterProduct = "organizations"
	)
	getAccountInviteAccepterClient, err := cfg.NewServiceClient(getAccountInviteAccepterProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	getAccountInviteAccepterRespBody, err := getAccountInvite(getAccountInviteAccepterClient, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving AccountInvite")
	}

	handshake := utils.PathSearch("handshake", getAccountInviteAccepterRespBody, nil)
	if handshake == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("invitation_id", d.Id()),
		d.Set("urn", utils.PathSearch("urn", handshake, nil)),
		d.Set("account_id", utils.PathSearch("target.entity", handshake, nil)),
		d.Set("master_account_id", utils.PathSearch("management_account_id", handshake, nil)),
		d.Set("master_account_name", utils.PathSearch("management_account_name", handshake, nil)),
		d.Set("organization_id", utils.PathSearch("organization_id", handshake, nil)),
		d.Set("created_at", utils.PathSearch("created_at", handshake, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", handshake, nil)),
		d.Set("status", utils.PathSearch("status", handshake, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAccountInviteAccepterUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAccountInviteAccepterDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	if !d.Get("leave_organization_on_destroy").(bool) {
		return nil
	}

	// deleteAccountInviteAccepter: Delete Organizations account invite accepter
	var (
		deleteAccountInviteAccepterHttpUrl = "v1/organizations/leave"
		deleteAccountInviteAccepterProduct = "organizations"
	)
	deleteAccountInviteAccepterClient, err := cfg.NewServiceClient(deleteAccountInviteAccepterProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	// if the value of accept is false(decline the invite), it can not be removed
	_, err = getAccountInvite(deleteAccountInviteAccepterClient, d.Id())
	if err != nil {
		return nil
	}

	deleteAccountInviteAccepterPath := deleteAccountInviteAccepterClient.Endpoint + deleteAccountInviteAccepterHttpUrl

	deleteAccountInviteAccepterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteAccountInviteAccepterClient.Request("POST", deleteAccountInviteAccepterPath,
		&deleteAccountInviteAccepterOpt)
	if err != nil {
		return diag.Errorf("error deleting AccountInviteAccepter: %s", err)
	}

	return nil
}
