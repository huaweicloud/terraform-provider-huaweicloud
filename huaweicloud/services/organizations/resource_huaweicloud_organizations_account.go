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

var accountNotFoundErrCodes = []string{
	"Organizations.1325", // The account is being closed, cannot be closed repeatedly.
}

// @API Organizations POST /v2/organizations/accounts
// @API Organizations POST /v1/organizations/accounts
// @API Organizations GET /v1/organizations/create-account-status/{create_account_status_id}
// @API Organizations GET /v1/organizations/entities
// @API Organizations GET /v1/organizations/accounts/{account_id}
// @API Organizations GET /v1/organizations/{resource_type}/{resource_id}/tags
// @API Organizations POST /v1/organizations/accounts/{account_id}/move
// @API Organizations PATCH /v1/organizations/accounts/{account_id}
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
				Description: `The name of the account.`,
			},
			"agency_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The agency name of the account.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the account.`,
			},
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the organizational unit or root to which the account belongs.`,
			},
			"tags": common.TagsSchema(`The key/value pairs to be associated with the account`),
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The email address of the account.`,
			},
			"phone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The mobile number of the account.`,
			},
			"intl_number_prefix": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The prefix of a mobile number.`,
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
				Description: `How the account joined the organization.`,
			},
		},
	}
}

func createAccount(client *golangsdk.ServiceClient, version string, d *schema.ResourceData) (interface{}, error) {
	httpUrl := fmt.Sprintf("%s/organizations/accounts", version)
	createPath := client.Endpoint + httpUrl
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateAccountBodyParams(d)),
	}
	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		respBody interface{}
	)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	if d.Get("email") != "" || d.Get("phone") != "" {
		respBody, err = createAccount(client, "v1", d)
	} else {
		respBody, err = createAccount(client, "v2", d)
	}

	if err != nil {
		return diag.Errorf("error creating account (%s): %s", d.Get("name"), err)
	}

	// we cannot get the account ID in API response, retrieve it from ShowCreateAccountStatus API
	statusID := utils.PathSearch("create_account_status.id", respBody, "").(string)
	if statusID == "" {
		return diag.Errorf("unable to find the account creation status ID from the API response")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      accountStateRefreshFunc(client, statusID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	accountStatusRespBody, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for account (%v) to be created: %s", d.Get("name"), err)
	}

	accountId := utils.PathSearch("create_account_status.account_id", accountStatusRespBody, "").(string)
	if accountId == "" {
		return diag.Errorf("unable to find the account ID from the API response")
	}
	d.SetId(accountId)

	if v, ok := d.GetOk("parent_id"); ok {
		parentID, err := getParentIdByAccountId(client, accountId)
		if err != nil {
			return diag.FromErr(err)
		}

		if v.(string) != parentID {
			err = moveAccount(client, accountId, parentID, v.(string))
			if err != nil {
				return diag.Errorf("error moving account (%s) to organization unit (%v): %s", accountId, v, err)
			}
		}
	}

	return resourceAccountRead(ctx, d, meta)
}

func accountStateRefreshFunc(client *golangsdk.ServiceClient, accountStatusId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		httpUrl := "v1/organizations/create-account-status/{create_account_status_id}"
		getPath := client.Endpoint + httpUrl
		getPath = strings.ReplaceAll(getPath, "{create_account_status_id}", accountStatusId)

		getOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		}
		resp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return nil, "failed", err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, "failed", err
		}

		state := utils.PathSearch("create_account_status.state", respBody, "").(string)
		if state == "succeeded" {
			return respBody, "COMPLETED", nil
		}

		if state == "failed" {
			return respBody, state, fmt.Errorf("state: %s; failure_reason: %v", state,
				utils.PathSearch("create_account_status.failure_reason", respBody, ""))
		}

		return respBody, state, nil
	}
}

func buildCreateAccountBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"email":       utils.ValueIgnoreEmpty(d.Get("email")),
		"phone":       utils.ValueIgnoreEmpty(d.Get("phone")),
		"agency_name": utils.ValueIgnoreEmpty(d.Get("agency_name")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"tags":        utils.ValueIgnoreEmpty(utils.ExpandResourceTags(d.Get("tags").(map[string]interface{}))),
	}
	return bodyParams
}

func GetAccountInfoById(client *golangsdk.ServiceClient, accountId string) (interface{}, error) {
	getHttpUrl := "v1/organizations/accounts/{account_id}?with_register_contact_info=true"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{account_id}", accountId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	status := utils.PathSearch("account.status", respBody, "").(string)
	if status == "" || status == "pending_closure" || status == "suspended" {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/organizations/accounts/{account_id}",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the account (%s) does not exist", accountId)),
			},
		}
	}

	return utils.PathSearch("account", respBody, nil), nil
}

func resourceAccountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	accountId := d.Id()
	account, err := GetAccountInfoById(client, accountId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
			"error retrieving account")
	}

	parentID, err := getParentIdByAccountId(client, accountId)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("parent_id", parentID),
		d.Set("name", utils.PathSearch("name", account, nil)),
		d.Set("email", utils.PathSearch("email", account, nil)),
		d.Set("phone", utils.PathSearch("mobile_phone", account, nil)),
		d.Set("description", utils.PathSearch("description", account, nil)),
		d.Set("intl_number_prefix", utils.PathSearch("intl_number_prefix", account, nil)),
		d.Set("urn", utils.PathSearch("urn", account, nil)),
		d.Set("joined_at", utils.PathSearch("joined_at", account, nil)),
		d.Set("joined_method", utils.PathSearch("join_method", account, nil)),
	)

	tagMap, err := getTags(client, accountsType, accountId)
	if err != nil {
		log.Printf("[WARN] error fetching tags of Organizations account (%s): %s", accountId, err)
	} else {
		mErr = multierror.Append(mErr, d.Set("tags", tagMap))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	accountId := d.Id()
	if d.HasChange("parent_id") {
		oldVal, newVal := d.GetChange("parent_id")
		err = moveAccount(client, accountId, oldVal.(string), newVal.(string))
		if err != nil {
			return diag.Errorf("error moving account (%s) to organizational unit (%v): %s", accountId, newVal, err)
		}
	}

	if d.HasChange("description") {
		err = updateAccount(client, d)
		if err != nil {
			return diag.Errorf("error updating account (%s): %s", accountId, err)
		}
	}

	if d.HasChange("tags") {
		err = updateTags(d, client, accountsType, accountId, "tags")
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
	updateHttpUrl := "v1/organizations/accounts/{account_id}"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{account_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateAccountBodyParams(d),
	}
	_, err := client.Request("PATCH", updatePath, &updateOpt)
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
	moveHttpUrl := "v1/organizations/accounts/{account_id}/move"
	movePath := client.Endpoint + moveHttpUrl
	movePath = strings.ReplaceAll(movePath, "{account_id}", accountId)

	moveOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildMoveAccountBodyParams(sourceParentID, destinationParentID),
	}
	_, err := client.Request("POST", movePath, &moveOpt)
	return err
}

func resourceAccountDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1/organizations/accounts/{account_id}/close"
	)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{account_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected401ErrInto404Err(
				common.ConvertExpected400ErrInto404Err(err, "error_code", accountNotFoundErrCodes...),
				"error_code", organizationNotFoundErrCodes...,
			),
			"error deleting account",
		)
	}

	return nil
}

func getParentIdByAccountId(client *golangsdk.ServiceClient, accountID string) (string, error) {
	getParentHttpUrl := "v1/organizations/entities?child_id={account_id}"
	getParentPath := client.Endpoint + getParentHttpUrl
	getParentPath = strings.ReplaceAll(getParentPath, "{account_id}", accountID)

	getParentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getAccountResp, err := client.Request("GET", getParentPath, &getParentOpt)
	if err != nil {
		return "", fmt.Errorf("error retrieving parent by account ID: %s", accountID)
	}
	respBody, err := utils.FlattenResponse(getAccountResp)
	if err != nil {
		return "", err
	}

	return utils.PathSearch("entities|[0].id", respBody, "").(string), nil
}
