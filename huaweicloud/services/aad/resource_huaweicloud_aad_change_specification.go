package aad

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableChangeSpecificationParams = []string{
	"instance_id",
	"upgrade_data",
	"upgrade_data.*.basic_bandwidth",
	"upgrade_data.*.elastic_bandwidth",
	"upgrade_data.*.service_bandwidth",
	"upgrade_data.*.port_num",
	"upgrade_data.*.bind_domain_num",
}

// @API AAD PUT /v2/aad/instance
func ResourceChangeSpecification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChangeSpecificationCreate,
		ReadContext:   resourceChangeSpecificationRead,
		UpdateContext: resourceChangeSpecificationUpdate,
		DeleteContext: resourceChangeSpecificationDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableChangeSpecificationParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the AAD instance ID.`,
			},
			"upgrade_data": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `Specifies the upgrade data.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"basic_bandwidth": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the basic bandwidth (Gbps).`,
						},
						"elastic_bandwidth": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the elastic bandwidth (Gbps).`,
						},
						"service_bandwidth": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `Specifies the service bandwidth (Mbps).`,
						},
						"port_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `Specifies the port number.`,
						},
						"bind_domain_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `Specifies the bind domain number.`,
						},
					},
				},
			},
		},
	}
}

func buildUpgradeDataBodyParams(d *schema.ResourceData) interface{} {
	rawArray := d.Get("upgrade_data").([]interface{})
	if len(rawArray) == 0 {
		// When the length of the parameter array is 0, an empty map is returned.
		return make(map[string]interface{})
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		// When the length of the parameter array is 0, an empty map is returned.
		return make(map[string]interface{})
	}

	rst := map[string]interface{}{
		"basic_bandwidth":   utils.ValueIgnoreEmpty(rawMap["basic_bandwidth"]),
		"elastic_bandwidth": utils.ValueIgnoreEmpty(rawMap["elastic_bandwidth"]),
		"service_bandwidth": utils.ValueIgnoreEmpty(rawMap["service_bandwidth"]),
		"port_num":          utils.ValueIgnoreEmpty(rawMap["port_num"]),
		"bind_domain_num":   utils.ValueIgnoreEmpty(rawMap["bind_domain_num"]),
	}

	return utils.RemoveNil(rst)
}

func buildChangeSpecificationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"instance_id":  d.Get("instance_id"),
		"upgrade_data": buildUpgradeDataBodyParams(d),
	}
}

func resourceChangeSpecificationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v2/aad/instance"
		product    = "aad"
		instanceID = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildChangeSpecificationBodyParams(d),
	}

	resp, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating AAD instance specification: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	orderID := utils.PathSearch("order_id", respBody, "").(string)
	if orderID != "" {
		// When the upgrade or subtraction of specifications results in a change in fees, the order ID will be returned.
		bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		if err = common.WaitOrderComplete(ctx, bssClient, orderID, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("the order is not completed while updating AAD instance (%s) specification: %s",
				instanceID, err)
		}
		if _, err = common.WaitOrderAllResourceComplete(ctx, bssClient, orderID, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(instanceID)
	return resourceChangeSpecificationRead(ctx, d, meta)
}

func resourceChangeSpecificationRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceChangeSpecificationUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceChangeSpecificationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to modify AAD specification. Deleting this resource
will not change the current AAD specification, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
