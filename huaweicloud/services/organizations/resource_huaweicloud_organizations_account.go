// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product Organizations
// ---------------------------------------------------------------

package organizations

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Organizations POST /v1/organizations/accounts
// @API Organizations GET /v1/organizations/accounts/{account_id}
// @API Organizations GET /v1/organizations/{resource_type}/{resource_id}/tags
// @API Organizations POST /v1/organizations/accounts/{account_id}/move
// @API Organizations PATCH /v1/organizations/accounts/{account_id}
// @API Organizations GET /v1/organizations/entities
// @API Organizations GET /v1/organizations/create-account-status/{create_account_status_id}
// @API Organizations POST /v1/organizations/{resource_type}/{resource_id}/tags/delete
// @API Organizations POST /v1/organizations/{resource_type}/{resource_id}/tags/create
// @API Organizations POST /v1/organizations/accounts/{account_id}/close
func ResourceAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccountCreate,
		UpdateContext: resourceAccountUpdate,
		ReadContext:   resourceAccountRead,
		DeleteContext: resourceAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the account.`,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"phone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"agency_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": common.TagsSchema(),
			"intl_number_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"joined_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"joined_method": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createAccount: create Organizations account
	var (
		createAccountHttpUrl = "v1/organizations/accounts"
		createAccountProduct = "organizations"
	)
	createAccountClient, err := cfg.NewServiceClient(createAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	createAccountPath := createAccountClient.Endpoint + createAccountHttpUrl

	createAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createAccountOpt.JSONBody = utils.RemoveNil(buildCreateAccountBodyParams(d))
	createAccountResp, err := createAccountClient.Request("POST", createAccountPath, &createAccountOpt)
	if err != nil {
		return diag.Errorf("error creating Account: %s", err)
	}

	createAccountRespBody, err := utils.FlattenResponse(createAccountResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// we cannot get the account ID in API response, retrieve it from ShowCreateAccountStatus API
	statusID := utils.PathSearch("create_account_status.id", createAccountRespBody, "").(string)
	if statusID == "" {
		return diag.Errorf("error creating Account: state is not found in API response")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"in_progress"},
		Target:       []string{"succeeded"},
		Refresh:      accountStateRefreshFunc(createAccountClient, statusID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	accountName := d.Get("name").(string)
	accountStatusRespBody, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for Organizations account (%s) to create: %s", accountName, err)
	}

	// set the account ID
	accountId := utils.PathSearch("create_account_status.account_id", accountStatusRespBody, "").(string)
	if accountId == "" {
		return diag.Errorf("unable to find the account ID from the API response")
	}
	d.SetId(accountId)

	if v, ok := d.GetOk("parent_id"); ok {
		parentID, err := getParentIdByAccountId(createAccountClient, accountId)
		if err != nil {
			return diag.FromErr(err)
		}
		if v.(string) != parentID {
			err = moveAccount(createAccountClient, d.Id(), parentID, v.(string))
			if err != nil {
				return diag.Errorf("error moving Account %s to organization unit %s: %s", d.Id(), v.(string), err)
			}
		}
	}

	return resourceAccountRead(ctx, d, meta)
}

func accountStateRefreshFunc(client *golangsdk.ServiceClient, accountStatusId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getAccountStatusHttpUrl := "v1/organizations/create-account-status/{create_account_status_id}"
		getAccountStatusPath := client.Endpoint + getAccountStatusHttpUrl
		getAccountStatusPath = strings.ReplaceAll(getAccountStatusPath, "{create_account_status_id}", accountStatusId)

		getAccountStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getAccountStatusResp, err := client.Request("GET", getAccountStatusPath, &getAccountStatusOpt)
		if err != nil {
			return nil, "", err
		}

		getAccountStatusRespBody, err := utils.FlattenResponse(getAccountStatusResp)
		if err != nil {
			return nil, "", err
		}

		state := utils.PathSearch("create_account_status.state", getAccountStatusRespBody, "").(string)

		reason := utils.PathSearch("create_account_status.failure_reason", getAccountStatusRespBody, "").(string)
		if reason != "" {
			return getAccountStatusRespBody, state, fmt.Errorf("state: %s; failure_reason: %s", state, reason)
		}

		return getAccountStatusRespBody, state, nil
	}
}

func buildCreateAccountBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"email":       utils.ValueIgnoreEmpty(d.Get("email")),
		"phone":       utils.ValueIgnoreEmpty(d.Get("phone")),
		"agency_name": utils.ValueIgnoreEmpty(d.Get("agency_name")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"tags":        utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}
	return bodyParams
}

func resourceAccountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getAccount: Query Organizations account
	var (
		getAccountProduct = "organizations"
	)
	getAccountClient, err := cfg.NewServiceClient(getAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	getAccountHttpUrl := "v1/organizations/accounts/{account_id}?with_register_contact_info=true"
	getAccountPath := getAccountClient.Endpoint + getAccountHttpUrl
	getAccountPath = strings.ReplaceAll(getAccountPath, "{account_id}", d.Id())

	getAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getAccountResp, err := getAccountClient.Request("GET", getAccountPath, &getAccountOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Account")
	}

	getAccountRespBody, err := utils.FlattenResponse(getAccountResp)
	if err != nil {
		return diag.FromErr(err)
	}

	status := utils.PathSearch("account.status", getAccountRespBody, "").(string)
	if status == "" || status == "pending_closure" || status == "suspended" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	parentID, err := getParentIdByAccountId(getAccountClient, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("parent_id", parentID),
		d.Set("name", utils.PathSearch("account.name", getAccountRespBody, nil)),
		d.Set("email", utils.PathSearch("account.email", getAccountRespBody, nil)),
		d.Set("phone", utils.PathSearch("account.mobile_phone", getAccountRespBody, nil)),
		d.Set("description", utils.PathSearch("account.description", getAccountRespBody, nil)),
		d.Set("intl_number_prefix", utils.PathSearch("account.intl_number_prefix", getAccountRespBody, nil)),
		d.Set("urn", utils.PathSearch("account.urn", getAccountRespBody, nil)),
		d.Set("joined_at", utils.PathSearch("account.joined_at", getAccountRespBody, nil)),
		d.Set("joined_method", utils.PathSearch("account.join_method", getAccountRespBody, nil)),
	)

	tagMap, err := getTags(getAccountClient, accountsType, d.Id())
	if err != nil {
		log.Printf("[WARN] error fetching tags of Organizations account (%s): %s", d.Id(), err)
	} else {
		mErr = multierror.Append(mErr, d.Set("tags", tagMap))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updateAccount: update Organizations account
	var (
		updateAccountProduct = "organizations"
	)
	updateAccountClient, err := cfg.NewServiceClient(updateAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	if d.HasChange("parent_id") {
		oldVal, newVal := d.GetChange("parent_id")
		err = moveAccount(updateAccountClient, d.Id(), oldVal.(string), newVal.(string))
		if err != nil {
			return diag.Errorf("error moving account: %s", err)
		}
	}

	if d.HasChange("description") {
		err = updateAccount(updateAccountClient, d)
		if err != nil {
			return diag.Errorf("error updating account: %s", err)
		}
	}

	if d.HasChange("tags") {
		err = updateTags(d, updateAccountClient, accountsType, d.Id(), "tags")
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceAccountRead(ctx, d, meta)
}

func buildUpdateAccountBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"description": d.Get("description"),
	}
	return bodyParams
}

func updateAccount(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		updateAccountHttpUrl = "v1/organizations/accounts/{account_id}"
	)
	updateAccountPath := client.Endpoint + updateAccountHttpUrl
	updateAccountPath = strings.ReplaceAll(updateAccountPath, "{account_id}", d.Id())

	updateAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateAccountOpt.JSONBody = utils.RemoveNil(buildUpdateAccountBodyParams(d))
	_, err := client.Request("PATCH", updateAccountPath, &updateAccountOpt)
	return err
}

func buildMoveAccountBodyParams(oldOrganizationsUnitId, newOrganizationsUnitId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"source_parent_id":      oldOrganizationsUnitId,
		"destination_parent_id": newOrganizationsUnitId,
	}
	return bodyParams
}

func moveAccount(client *golangsdk.ServiceClient, accountId, sourceParentID, destinationParentID string) error {
	// moveAccount: update Organizations account
	var (
		moveAccountHttpUrl = "v1/organizations/accounts/{account_id}/move"
	)
	moveAccountPath := client.Endpoint + moveAccountHttpUrl
	moveAccountPath = strings.ReplaceAll(moveAccountPath, "{account_id}", accountId)

	moveAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	moveAccountOpt.JSONBody = utils.RemoveNil(buildMoveAccountBodyParams(sourceParentID, destinationParentID))
	_, err := client.Request("POST", moveAccountPath, &moveAccountOpt)
	return err
}

func resourceAccountDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteAccount: close Organizations account
	var (
		deleteAccountHttpUrl = "v1/organizations/accounts/{account_id}/close"
		deleteAccountProduct = "organizations"
	)
	deleteAccountClient, err := cfg.NewServiceClient(deleteAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	deleteAccountPath := deleteAccountClient.Endpoint + deleteAccountHttpUrl
	deleteAccountPath = strings.ReplaceAll(deleteAccountPath, "{account_id}", d.Id())

	deleteAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteAccountClient.Request("POST", deleteAccountPath, &deleteAccountOpt)
	if err != nil {
		return diag.Errorf("error deleting account: %s", err)
	}

	return nil
}

func getParentIdByAccountId(client *golangsdk.ServiceClient, accountID string) (string, error) {
	getParentHttpUrl := "v1/organizations/entities?child_id={account_id}"
	getParentPath := client.Endpoint + getParentHttpUrl
	getParentPath = strings.ReplaceAll(getParentPath, "{account_id}", accountID)

	getParentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getAccountResp, err := client.Request("GET", getParentPath, &getParentOpt)
	if err != nil {
		return "", fmt.Errorf("error retrieving parent by account_id: %s", accountID)
	}
	getAccountRespBody, err := utils.FlattenResponse(getAccountResp)
	if err != nil {
		return "", err
	}

	id := utils.PathSearch("entities|[0].id", getAccountRespBody, "").(string)

	return id, nil
}
