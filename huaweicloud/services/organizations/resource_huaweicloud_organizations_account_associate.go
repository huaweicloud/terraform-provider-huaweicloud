// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product Organizations
// ---------------------------------------------------------------

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
				Description: `Specifies the ID of the account.`,
			},
			"parent_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of root or organizational unit in which you want to move the account.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the account.`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the uniform resource name of the account.`,
			},
			"joined_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time when the account was created.`,
			},
			"joined_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates how an account joined an organization.`,
			},
		},
	}
}

func resourceAccountAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createAccountAssociate: create Organizations account associate
	var (
		createAccountAssociateProduct = "organizations"
	)
	createAccountAssociateClient, err := cfg.NewServiceClient(createAccountAssociateProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	accountID := d.Get("account_id").(string)
	oParentID, err := getParentIdByAccountId(createAccountAssociateClient, accountID)
	if err != nil {
		return diag.FromErr(err)
	}
	nParentID := d.Get("parent_id").(string)
	if oParentID != nParentID {
		err = moveAccount(createAccountAssociateClient, accountID, oParentID, nParentID)
		if err != nil {
			return diag.Errorf("error updating Account: %s", err)
		}
	}

	d.SetId(accountID)

	return resourceAccountAssociateRead(ctx, d, meta)
}

func resourceAccountAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getAccountAssociate: Query Organizations account associate
	var (
		getAccountAssociateHttpUrl = "v1/organizations/accounts/{account_id}"
		getAccountAssociateProduct = "organizations"
	)
	getAccountAssociateClient, err := cfg.NewServiceClient(getAccountAssociateProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	getAccountAssociatePath := getAccountAssociateClient.Endpoint + getAccountAssociateHttpUrl
	getAccountAssociatePath = strings.ReplaceAll(getAccountAssociatePath, "{account_id}", d.Id())

	getAccountAssociateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getAccountAssociateResp, err := getAccountAssociateClient.Request("GET", getAccountAssociatePath,
		&getAccountAssociateOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving AccountAssociate")
	}

	getAccountAssociateRespBody, err := utils.FlattenResponse(getAccountAssociateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	parentID, err := getParentIdByAccountId(getAccountAssociateClient, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("account_id", d.Id()),
		d.Set("parent_id", parentID),
		d.Set("name", utils.PathSearch("account.name", getAccountAssociateRespBody, nil)),
		d.Set("urn", utils.PathSearch("account.urn", getAccountAssociateRespBody, nil)),
		d.Set("joined_at", utils.PathSearch("account.joined_at", getAccountAssociateRespBody,
			nil)),
		d.Set("joined_method", utils.PathSearch("account.join_method", getAccountAssociateRespBody,
			nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAccountAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updateAccountAssociate: update Organizations account associate
	var (
		updateAccountAssociateProduct = "organizations"
	)
	updateAccountAssociateClient, err := cfg.NewServiceClient(updateAccountAssociateProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	oParentID, err := getParentIdByAccountId(updateAccountAssociateClient, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	nParentID := d.Get("parent_id").(string)

	err = moveAccount(updateAccountAssociateClient, d.Id(), oParentID, nParentID)
	if err != nil {
		return diag.Errorf("error updating Account: %s", err)
	}

	return resourceAccountAssociateRead(ctx, d, meta)
}

func resourceAccountAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteAccountAssociate: Delete Organizations account associate
	var (
		deleteAccountAssociateProduct = "organizations"
	)
	deleteAccountAssociateClient, err := cfg.NewServiceClient(deleteAccountAssociateProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	getRootRespBody, err := getRoot(deleteAccountAssociateClient)
	if err != nil {
		return diag.FromErr(err)
	}
	rootId := utils.PathSearch("roots|[0].id", getRootRespBody, "").(string)
	parentID := d.Get("parent_id").(string)
	err = moveAccount(deleteAccountAssociateClient, d.Id(), parentID, rootId)
	if err != nil {
		return diag.Errorf("error updating Account: %s", err)
	}

	return nil
}
