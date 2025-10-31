package rgc

import (
	"context"
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

var enrollAccountNonUpdatableParams = []string{"managed_account_id"}

// @API RGC POST /v1/managed-organization/accounts/{managed_account_id}/enroll
// @API RGC GET /v1/managed-organization/managed-accounts/{managed_account_id}
// @API RGC POST /v1/managed-organization/managed-accounts/{managed_account_id}/update
// @API RGC POST /v1/managed-organization/managed-accounts/{managed_account_id}/un-enroll
func ResourceAccountEnroll() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnrollAccountCreate,
		UpdateContext: resourceEnrollAccountUpdate,
		ReadContext:   resourceEnrollAccountRead,
		DeleteContext: resourceEnrollAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Hour),
			Delete: schema.DefaultTimeout(6 * time.Hour),
		},

		CustomizeDiff: config.FlexibleForceNew(enrollAccountNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"managed_account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_organizational_unit_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_organizational_unit_name": {
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
							Default:  false,
						},
					},
				},
			},
			"landing_zone_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"manage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceEnrollAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		enrollAccountHttpUrl = "v1/managed-organization/accounts/{managed_account_id}/enroll"
		enrollAccountProduct = "rgc"
	)
	enrollAccountClient, err := cfg.NewServiceClient(enrollAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	managedAccountId := d.Get("managed_account_id").(string)
	enrollAccountPath := enrollAccountClient.Endpoint + enrollAccountHttpUrl
	enrollAccountPath = strings.ReplaceAll(enrollAccountPath, "{managed_account_id}", managedAccountId)

	enrollAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	enrollAccountOpt.JSONBody = utils.RemoveNil(buildCreateAccountBodyParams(d))
	enrollAccountResp, err := enrollAccountClient.Request("POST", enrollAccountPath, &enrollAccountOpt)
	if err != nil {
		return diag.Errorf("error enrolling Account: %s", err)
	}

	enrollAccountRespBody, err := utils.FlattenResponse(enrollAccountResp)
	if err != nil {
		return diag.FromErr(err)
	}

	operationID := utils.PathSearch("operation_id", enrollAccountRespBody, "").(string)
	if operationID == "" {
		return diag.Errorf("unable to find the account operation ID from the API response")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"IN_PROGRESS"},
		Target:       []string{"SUCCEEDED"},
		Refresh:      accountStateRefreshFunc(enrollAccountClient, operationID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	d.SetId(managedAccountId)
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for RGC account (%s) to enroll: %s", managedAccountId, err)
	}

	return resourceEnrollAccountRead(ctx, d, meta)
}

func resourceEnrollAccountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getAccountProduct = "rgc"
	)
	getAccountClient, err := cfg.NewServiceClient(getAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating rgc client: %s", err)
	}

	getAccountHttpUrl := "v1/managed-organization/managed-accounts/{managed_account_id}"
	getAccountPath := getAccountClient.Endpoint + getAccountHttpUrl
	getAccountPath = strings.ReplaceAll(getAccountPath, "{managed_account_id}", d.Id())

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

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("landing_zone_version", utils.PathSearch("landing_zone_version", getAccountRespBody, nil)),
		d.Set("manage_account_id", utils.PathSearch("manage_account_id", getAccountRespBody, nil)),
		d.Set("account_type", utils.PathSearch("account_type", getAccountRespBody, nil)),
		d.Set("account_name", utils.PathSearch("account_name", getAccountRespBody, nil)),
		d.Set("owner", utils.PathSearch("owner", getAccountRespBody, nil)),
		d.Set("stage", utils.PathSearch("state", getAccountRespBody, nil)),
		d.Set("created_at", utils.PathSearch("created_at", getAccountRespBody, nil)),
		d.Set("created_at", utils.PathSearch("updated_at", getAccountRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEnrollAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updateAccountHttpUrl = "v1/managed-organization/managed-accounts/{managed_account_id}/update"
		updateAccountProduct = "rgc"
	)
	updateAccountClient, err := cfg.NewServiceClient(updateAccountProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	updateAccountPath := updateAccountClient.Endpoint + updateAccountHttpUrl
	updateAccountPath = strings.ReplaceAll(updateAccountPath, "{managed_account_id}", d.Id())

	updateAccountOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildEnrollAccountBodyParams(d)),
	}
	updateAccountResp, err := updateAccountClient.Request("POST", updateAccountPath, &updateAccountOpt)
	if err != nil {
		return diag.Errorf("error update Account: %s", err)
	}

	updateAccountRespBody, err := utils.FlattenResponse(updateAccountResp)
	if err != nil {
		return diag.FromErr(err)
	}

	operationID := utils.PathSearch("operation_id", updateAccountRespBody, "").(string)
	if operationID == "" {
		return diag.Errorf("unable to find the account operation ID from the API response")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"IN_PROGRESS"},
		Target:       []string{"SUCCEEDED"},
		Refresh:      accountStateRefreshFunc(updateAccountClient, operationID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for RGC account[%s] to update: %s", d.Id(), err)
	}

	return resourceEnrollAccountRead(ctx, d, meta)
}

func resourceEnrollAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for RGC account (%s) to un-enroll: %s", d.Id(), err)
	}

	return nil
}

func buildEnrollAccountBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"parent_organizational_unit_id": d.Get("parent_organizational_unit_id"),
		"blueprint":                     buildBlueprintBodyParams(d),
	}

	if v, ok := d.GetOk("parent_organizational_unit_name"); ok {
		bodyParams["parent_organizational_unit_name"] = v
	}

	return bodyParams
}
