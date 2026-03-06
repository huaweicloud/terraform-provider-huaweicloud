package organizations

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Organizations POST /v1/received-handshakes/{handshake_id}/accept
// @API Organizations GET /v1/organizations/handshakes/{handshake_id}
// @API Organizations POST /v1/organizations/leave
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
				Description: `The unique ID of an invitation (handshake).`,
			},
			"leave_organization_on_destroy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to leave the organization when delete the invitation (handshake).`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uniform resource name of the invitation`,
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the target account.`,
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
cancelled, declined, or expired.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current state of the invitation (handshake).`,
			},
		},
	}
}

func resourceAccountInviteAccepterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1/received-handshakes/{handshake_id}/accept"
	)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{handshake_id}", d.Get("invitation_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error accepting the received account invitation: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	handshakeId := utils.PathSearch("handshake.id", respBody, "").(string)
	if handshakeId == "" {
		return diag.Errorf("unable to find the handshake ID of the Organizations account invite accepter from the API response")
	}
	d.SetId(handshakeId)

	return resourceAccountInviteAccepterRead(ctx, d, meta)
}

func resourceAccountInviteAccepterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	invitationId := d.Id()
	handshake, err := GetAccountInviteById(client, invitationId)
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			// After accepting the invitation, leave the organization.
			common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
			"error retrieving received account invitation",
		)
	}

	mErr := multierror.Append(
		d.Set("invitation_id", invitationId),
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
	if !d.Get("leave_organization_on_destroy").(bool) {
		return nil
	}

	cfg := meta.(*config.Config)
	httpUrl := "v1/organizations/leave"
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	// if the value of accept is false(decline the invite), it can not be removed
	_, err = GetAccountInviteById(client, d.Id())
	if err != nil {
		return nil
	}

	deletePath := client.Endpoint + httpUrl
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error leaving the organization")
	}

	return nil
}
