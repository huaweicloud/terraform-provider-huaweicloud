package cbr

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableChangeOrderParams = []string{
	"resource_id",
	"product_info",
	"product_info.*.product_id",
	"product_info.*.resource_size",
	"product_info.*.resource_size_measure_id",
	"product_info.*.resource_spec_code",
	"promotion_info",
	"cloud_service_console_url",
}

// @API CBR POST /v3/{project_id}/orders/change
func ResourceChangeOrder() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChangeOrderCreate,
		ReadContext:   resourceChangeOrderRead,
		UpdateContext: resourceChangeOrderUpdate,
		DeleteContext: resourceChangeOrderDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableChangeOrderParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the resource to change.`,
			},
			"product_info": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `Specifies product information.`,
				Elem:        productInfoSchema(),
			},
			"promotion_info": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the promotion information for the order.`,
			},
			// This field seems to have no usage scenario and is retained simply to align the API.
			"cloud_service_console_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: utils.SchemaDesc("Specifies the Console URL of the cloud service.", utils.SchemaDescInput{Internal: true}),
			},
			// Internal field to trigger ForceNew manually when needed.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func productInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"product_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the product ID.`,
			},
			"resource_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the size of the resource.`,
			},
			"resource_size_measure_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the measurement unit of the resource size.`,
			},
			"resource_spec_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the spec code of the resource.`,
			},
		},
	}
}

func buildProductInfo(d *schema.ResourceData) map[string]interface{} {
	raw := d.Get("product_info.0").(map[string]interface{})
	return map[string]interface{}{
		"product_id":               raw["product_id"],
		"resource_size":            raw["resource_size"],
		"resource_size_measure_id": raw["resource_size_measure_id"],
		"resource_spec_code":       raw["resource_spec_code"],
	}
}

func buildChangeOrderBodyParams(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{
		"resource_id":               d.Get("resource_id"),
		"product_info":              buildProductInfo(d),
		"is_auto_pay":               true,
		"promotion_info":            utils.ValueIgnoreEmpty(d.Get("promotion_info")),
		"cloud_service_console_url": utils.ValueIgnoreEmpty(d.Get("cloud_service_console_url")),
	}
	return utils.RemoveNil(body)
}

func resourceChangeOrderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/orders/change"
		product = "cbr"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildChangeOrderBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error changing CBR order: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	orderId := utils.PathSearch("orderId", respBody, "").(string)
	if orderId == "" {
		return diag.Errorf("error changing CBR order: order ID not found in API response")
	}

	bssClient, err := cfg.BssV2Client(region)
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}

	if err := common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for CBR change order (%s) to complete: %s", orderId, err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	return resourceChangeOrderRead(ctx, d, meta)
}

func resourceChangeOrderRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceChangeOrderUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceChangeOrderDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to change the order of a CBR resource. Deleting this 
resource will not change the actual order result, but will only remove the resource information from the 
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
