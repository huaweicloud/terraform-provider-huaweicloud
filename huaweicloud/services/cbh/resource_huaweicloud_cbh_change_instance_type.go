package cbh

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CBH PUT /v2/{project_id}/cbs/instance/type
// @API BSS POST /v3/orders/customer-orders/pay
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
func ResourceChangeInstanceType() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChangeInstanceTypeCreate,
		UpdateContext: resourceChangeInstanceTypeUpdate,
		ReadContext:   resourceChangeInstanceTypeRead,
		DeleteContext: resourceChangeInstanceTypeDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"server_id",
			"availability_zone",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
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

func buildChangeInstanceTypeQueryParams(d *schema.ResourceData) string {
	rst := fmt.Sprintf("?server_id=%s", d.Get("server_id").(string))

	if v, ok := d.GetOk("availability_zone"); ok {
		rst += fmt.Sprintf("&availability_zone=%s", v)
	}
	return rst
}

// The API supports configuring the `is_auto_pay` field for automatic payments, but this field doesn't actually work.
// Therefore, it's better to manually pay for orders using the CBC API instead.
func resourceChangeInstanceTypeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v2/{project_id}/cbs/instance/type"
		product  = "cbh"
		serverId = d.Get("server_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildChangeInstanceTypeQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error changing CBH instance type: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	orderId := utils.PathSearch("order_id", respBody, "").(string)
	if orderId == "" {
		return diag.Errorf("error changing CBH instance type: order ID is empty in API response")
	}

	d.SetId(serverId)

	bssClient, err := cfg.NewServiceClient("bssv2", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating BSS client: %s", err)
	}

	if err := cbc.PaySubscriptionOrder(bssClient, orderId); err != nil {
		return diag.Errorf("error paying for the order (%s) of changing the CBH instance type: %s", orderId, err)
	}

	if err := common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	_, err = common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the order (%s) of changing the CBH instance type to complete: %s", orderId, err)
	}

	return resourceChangeInstanceTypeRead(ctx, d, meta)
}

func resourceChangeInstanceTypeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceChangeInstanceTypeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceChangeInstanceTypeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to change CBH instance type.
Deleting this resource will not recover the change CBH instance type, but will only remove the resource information from
the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
