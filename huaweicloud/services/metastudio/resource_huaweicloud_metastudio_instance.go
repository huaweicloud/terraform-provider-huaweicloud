package metastudio

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableParams = []string{
	"period_type",
	"period_num",
	"is_auto_renew",
	"resource_spec_code",
}

// @API MetaStudio POST /v1/{project_id}/tenants/special-agreements/signed
// @API MetaStudio POST /v1/{project_id}/mss/public/orders
// @API MetaStudio GET /v1/{project_id}/tenants/resources
// @API BSS GET /V2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe

func ResourceMetaStudio() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMetaStudioCreate,
		ReadContext:   resourceMetaStudioRead,
		UpdateContext: resourceMetaStudioUpdate,
		DeleteContext: resourceMetaStudioDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"period_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{2, 3, 6}),
				Description:  `Specifies the charging period unit`,
			},
			"period_num": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 2147483647),
				Description:  `Specifies the number of periods to purchase.`,
			},
			"is_auto_renew": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 1),
				Description:  `Specifies whether to auto-renew the vault when it expires.`,
			},
			"resource_spec_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the resource specification code`,
			},
			"order_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The order ID of resource`,
			},
			"resource_expire_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource expire time`,
			},
			"business_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The business type of resource`,
			},
			"sub_resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The sub-resource type of resource`,
			},
			"is_sub_resource": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether it is a sub-resource. A sub-resource describes the quantity and unit of a subsidiary resource.`,
			},
			"charging_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The billing mode`,
			},
			"amount": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `Total amount`,
			},
			"usage": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `Usage amount`,
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Resource status. 0: Normal, 1: Frozen`,
			},
			"unit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unit of amount`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceMetaStudioCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("metastudio", region)
	if err != nil {
		return diag.Errorf("error creating MetaStudio client: %s", err)
	}
	// sign auto-pay agreements
	if err := signAutoPayAgreeMents(client); err != nil {
		return diag.Errorf("error signing auto-apy agreements: %s", err)
	}
	orderId, err := resourceCreate(client, d)
	if err != nil {
		return diag.Errorf("error creating MetaStudio: %s", err)
	}
	bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}
	timeout := d.Timeout(schema.TimeoutCreate)
	// wait for order complete
	if err := common.WaitOrderComplete(ctx, bssClient, orderId, timeout); err != nil {
		return diag.Errorf("the order (%s) is not completed while creating metaStudio : %v", orderId, err)
	}
	resourceId, err := common.WaitOrderAllResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resourceId)
	return resourceMetaStudioRead(ctx, d, meta)
}

func signAutoPayAgreeMents(client *golangsdk.ServiceClient) error {
	httpUrl := "v1/{project_id}/tenants/special-agreements/signed"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildSpecialAgreementSighReq()),
	}
	resp, err := client.Request("POST", requestPath, &createOpt)
	if err != nil {
		return err
	}
	if resp.StatusCode == 200 {
		return nil
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}
	return fmt.Errorf("failed to sign auto-pay-agreements: %s", respBody)
}

func buildSpecialAgreementSighReq() map[string]interface{} {
	return map[string]interface{}{
		"agreement_type": "AUTO_PAY",
		"to_sign":        true,
	}
}

func resourceCreate(client *golangsdk.ServiceClient, d *schema.ResourceData) (string, error) {
	httpUrl := "v1/{project_id}/mss/public/orders"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateMetaStudioParams(d)),
	}
	resp, err := client.Request("POST", requestPath, &createOpt)
	if err != nil {
		return "", err
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return "", err
	}
	orderId := utils.PathSearch("order_id", respBody, "").(string)
	return orderId, nil
}

func buildCreateMetaStudioParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"cloud_services": []interface{}{
			map[string]interface{}{
				"is_auto_pay":        1,
				"period_type":        d.Get("period_type"),
				"period_num":         d.Get("period_num"),
				"is_auto_renew":      d.Get("is_auto_renew"),
				"subscription_num":   1,
				"resource_spec_code": d.Get("resource_spec_code"),
			},
		},
	}
	return bodyParams
}

func resourceMetaStudioRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("metastudio", region)
	if err != nil {
		return diag.Errorf("error creating MetaStudio client: %s", err)
	}
	resourceDetail, err := GetResourceDetail(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Meta Studio Resource")
	}
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("order_id", utils.PathSearch("order_id", resourceDetail, nil)),
		d.Set("resource_expire_time", utils.PathSearch("resource_expire_time", resourceDetail, nil)),
		d.Set("business_type", utils.PathSearch("business_type", resourceDetail, nil)),
		d.Set("sub_resource_type", utils.PathSearch("sub_resource_type", resourceDetail, nil)),
		d.Set("is_sub_resource", utils.PathSearch("is_sub_resource", resourceDetail, false)),
		d.Set("charging_mode", utils.PathSearch("charging_mode", resourceDetail, nil)),
		d.Set("amount", utils.PathSearch("amount", resourceDetail, nil)),
		d.Set("usage", utils.PathSearch("usage", resourceDetail, nil)),
		d.Set("status", utils.PathSearch("status", resourceDetail, nil)),
		d.Set("unit", utils.PathSearch("unit", resourceDetail, nil)),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting MetaStudio resource fields: %s", mErr)
	}
	return nil
}

func GetResourceDetail(client *golangsdk.ServiceClient, resourceId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/tenants/resources"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = fmt.Sprintf("%s?resource_source=PURCHASED&resource_id=%v", requestPath, resourceId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}
	resources := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
	if len(resources) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}
	return resources[0], nil
}

func resourceMetaStudioDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("metastudio", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}
	resourceId := d.Id()
	if d.Get("charging_mode").(string) != "PERIODIC" {
		errorMsg := `This resource is a one-time action resource. Deleting this 
resource will not change the current resource status, but will only remove the resource information from the 
tfstate file.`
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  errorMsg,
			},
		}
	}
	if err := common.UnsubscribePrePaidResource(d, cfg, []string{resourceId}); err != nil {
		return diag.Errorf("error unsubscribing meta studio resource (%s): %s",
			resourceId, err)
	}
	if err := waitingForResourceDeleteCompleted(ctx, client, d); err != nil {
		return diag.Errorf("error waiting for Workspace APP server (%s) deleted: %s", d.Id(), err)
	}
	return nil
}

func waitingForResourceDeleteCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) interface{} {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			resourceDetail, err := GetResourceDetail(client, d.Id())
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return "deleted", "COMPLETED", nil
				}
				return nil, "ERROR", err
			}
			return resourceDetail, "PENDING", nil
		},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceMetaStudioUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource doesn't support update operation.
	return nil
}
