package organizations

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Organizations GET /v1/organizations/entities
// @API Organizations POST /v1/organizations/accounts/{account_id}/move
// @API Organizations GET /v1/organizations/accounts/{account_id}
// @API Organizations GET /v1/organizations/roots
func ResourceAccountAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccountAssociateCreate,
		UpdateContext: resourceAccountAssociateUpdate,
		ReadContext:   resourceAccountAssociateRead,
		DeleteContext: resourceAccountAssociateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the account.`,
			},
			"parent_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of root or organizational unit in which you want to move the account.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the account.`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uniform resource name of the account.`,
			},
			"joined_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the account was created.`,
			},
			"joined_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `How an account joined an organization.`,
			},
		},
	}
}

func resourceAccountAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	accountId := d.Get("account_id").(string)
	oParentId, err := getParentIdByAccountId(client, accountId)
	if err != nil {
		return diag.FromErr(err)
	}

	nParentId := d.Get("parent_id").(string)
	if oParentId != nParentId {
		err = moveAccount(client, accountId, oParentId, nParentId)
		if err != nil {
			return diag.Errorf("error moving account (%s) to organization unit (%v): %s", accountId, nParentId, err)
		}
	}

	d.SetId(accountId)

	return resourceAccountAssociateRead(ctx, d, meta)
}

func resourceAccountAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	accountId := d.Id()
	account, err := GetAccountById(client, accountId)
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
			"error retrieving account associate",
		)
	}

	parentId, err := getParentIdByAccountId(client, accountId)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("account_id", accountId),
		d.Set("parent_id", parentId),
		d.Set("name", utils.PathSearch("account.name", account, nil)),
		d.Set("urn", utils.PathSearch("account.urn", account, nil)),
		d.Set("joined_at", utils.PathSearch("account.joined_at", account, nil)),
		d.Set("joined_method", utils.PathSearch("account.join_method", account, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAccountAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	accountId := d.Id()
	oParentId, err := getParentIdByAccountId(client, accountId)
	if err != nil {
		return diag.FromErr(err)
	}

	nParentId := d.Get("parent_id").(string)
	err = moveAccount(client, accountId, oParentId, nParentId)
	if err != nil {
		return diag.Errorf("error moving account (%s) to organization unit (%v): %s", accountId, nParentId, err)
	}

	return resourceAccountAssociateRead(ctx, d, meta)
}

func resourceAccountAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	getRootRespBody, err := getRoot(client)
	if err != nil {
		return diag.FromErr(err)
	}

	accountId := d.Id()
	parentId := d.Get("parent_id").(string)
	err = moveAccount(client, accountId, parentId, utils.PathSearch("roots|[0].id", getRootRespBody, "").(string))
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
			fmt.Sprintf("error moving account (%s) to root", accountId))
	}

	return nil
}
