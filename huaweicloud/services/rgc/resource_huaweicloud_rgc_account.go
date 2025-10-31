package rgc

import (
	"context"
	"fmt"
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

var accountNonUpdatableParams = []string{"name", "email", "phone", "parent_organizational_unit_name",
	"parent_organizational_unit_id", "identity_store_user_name", "identity_store_email", "blueprint",
	"blueprint.*.blueprint_product_id", "blueprint.*.blueprint_product_version", "blueprint.*.variables",
	"blueprint.*.is_blueprint_has_multi_account_resource",
}

// @API RGC POST /v1/managed-organization/managed-accounts
// @API RGC GET /v1/managed-organization/{operation_id}
// @API Organizations GET /v1/organizations/accounts
// @API Organizations GET /v1/organizations/accounts/{account_id}
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
			Create: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(accountNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"identity_store_user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"identity_store_email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_organizational_unit_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_organizational_unit_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"phone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"blueprint": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blueprint_product_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"blueprint_product_version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"variables": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"is_blueprint_has_multi_account_resource": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
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

	// createAccount: create RGC account
	var (
		createAccountHttpUrl = "v1/managed-organization/managed-accounts"
		createAccountProduct = "rgc"
	)
	createAccountClient, err := cfg.NewServiceClient(createAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
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

	operationID := utils.PathSearch("operation_id", createAccountRespBody, "").(string)
	if operationID == "" {
		return diag.Errorf("unable to find the account operation ID from the API response")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"IN_PROGRESS"},
		Target:       []string{"SUCCEEDED"},
		Refresh:      accountStateRefreshFunc(createAccountClient, operationID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	accountName := d.Get("name").(string)
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for RGC account (%s) to create: %s", accountName, err)
	}

	id, err := getAccountIDbyName(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	return resourceAccountRead(ctx, d, meta)
}

func getAccountIDbyName(d *schema.ResourceData, meta interface{}) (string, error) {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	listAccountsHttpUrl := "v1/organizations/accounts"
	listAccountsProduct := "organizations"
	listAccountsClient, err := cfg.NewServiceClient(listAccountsProduct, region)
	if err != nil {
		return "", fmt.Errorf("error creating Organizations client: %s", err)
	}

	listAccountsPath := listAccountsClient.Endpoint + listAccountsHttpUrl
	listAccountsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	filterName := d.Get("name").(string)
	var marker string
	var queryPath string

	for {
		queryPath = listAccountsPath + buildListAccountsQueryParams(marker)
		listAccountsResp, err := listAccountsClient.Request("GET", queryPath, &listAccountsOpt)
		if err != nil {
			return "", fmt.Errorf("retrieving Organizations accounts: %s", err)
		}

		listAccountsRespBody, err := utils.FlattenResponse(listAccountsResp)
		if err != nil {
			return "", err
		}

		jsonPath := fmt.Sprintf("accounts[?name=='%s']|[0].id", filterName)
		id := utils.PathSearch(jsonPath, listAccountsRespBody, nil)
		if id != nil {
			return id.(string), nil
		}

		marker = utils.PathSearch("page_info.next_marker", listAccountsRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	return "", fmt.Errorf("unbale to find the id of account: %s", filterName)
}

func buildListAccountsQueryParams(marker string) string {
	// the default value of limit is 200
	res := "?limit=200"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	return res
}

func accountStateRefreshFunc(client *golangsdk.ServiceClient, operationId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getAccountStatusHttpUrl := "v1/managed-organization/{operation_id}"
		getAccountStatusPath := client.Endpoint + getAccountStatusHttpUrl
		getAccountStatusPath = strings.ReplaceAll(getAccountStatusPath, "{operation_id}", operationId)

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

		status := utils.PathSearch("status", getAccountStatusRespBody, nil)
		if status == nil {
			message := utils.PathSearch("message", getAccountStatusRespBody, nil)
			return nil, "", fmt.Errorf("status: %s; message: %s", status, message)
		}

		return getAccountStatusRespBody, status.(string), nil
	}
}

func buildCreateAccountBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"account_name":                    d.Get("name"),
		"account_email":                   d.Get("email"),
		"phone":                           utils.ValueIgnoreEmpty(d.Get("phone")),
		"identity_store_user_name":        d.Get("identity_store_user_name"),
		"identity_store_email":            d.Get("identity_store_email"),
		"parent_organizational_unit_name": d.Get("parent_organizational_unit_name"),
		"parent_organizational_unit_id":   d.Get("parent_organizational_unit_id"),
		"blueprint":                       buildBlueprintBodyParams(d),
	}

	return bodyParams
}

func buildBlueprintBodyParams(d *schema.ResourceData) map[string]interface{} {
	if v, ok := d.GetOk("blueprint"); ok {
		blueprint := v.([]interface{})
		if b, ok := blueprint[0].(map[string]interface{}); ok {
			return map[string]interface{}{
				"blueprint_product_id":                    utils.ValueIgnoreEmpty(b["blueprint_product_id"]),
				"blueprint_product_version":               utils.ValueIgnoreEmpty(b["blueprint_product_version"]),
				"variables":                               utils.ValueIgnoreEmpty(b["variables"]),
				"is_blueprint_has_multi_account_resource": utils.ValueIgnoreEmpty(b["is_blueprint_has_multi_account_resource"]),
			}
		}
	}

	return nil
}

func resourceAccountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getAccount: Query RGC account via organizations API
	var (
		getAccountProduct = "organizations"
	)
	getAccountClient, err := cfg.NewServiceClient(getAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	getAccountHttpUrl := "v1/organizations/accounts/{account_id}"
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

	mErr = multierror.Append(
		mErr,
		d.Set("name", utils.PathSearch("account.name", getAccountRespBody, nil)),
		d.Set("urn", utils.PathSearch("account.urn", getAccountRespBody, nil)),
		d.Set("joined_at", utils.PathSearch("account.joined_at", getAccountRespBody, nil)),
		d.Set("joined_method", utils.PathSearch("account.join_method", getAccountRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAccountUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// un-enroll Account: un-enroll RGC account
	var (
		unEnrollAccountHttpUrl = "v1/managed-organization/managed-accounts/{managed_account_id}/un-enroll"
		unEnrollAccountProduct = "rgc"
	)
	unEnrollAccountClient, err := cfg.NewServiceClient(unEnrollAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	unEnrollAccountPath := unEnrollAccountClient.Endpoint + unEnrollAccountHttpUrl
	unEnrollAccountPath = strings.ReplaceAll(unEnrollAccountPath, "{managed_account_id}", d.Id())

	unEnrollAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	unEnrollAccountResp, err := unEnrollAccountClient.Request("POST", unEnrollAccountPath, &unEnrollAccountOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error un-enrolling Account")
	}

	unEnrollAccountRespBody, err := utils.FlattenResponse(unEnrollAccountResp)
	if err != nil {
		return diag.FromErr(err)
	}

	operationID := utils.PathSearch("operation_id", unEnrollAccountRespBody, "").(string)
	if operationID == "" {
		return diag.Errorf("unable to find the account operation ID from the API response")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"IN_PROGRESS"},
		Target:       []string{"SUCCEEDED"},
		Refresh:      accountStateRefreshFunc(unEnrollAccountClient, operationID),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	accountName := d.Get("name").(string)
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for RGC account (%s) to un-enroll: %s", accountName, err)
	}

	// deleteAccount: close RGC account via organizations API
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
		return common.CheckDeletedDiag(d, err, "error deleting account")
	}

	return nil
}
