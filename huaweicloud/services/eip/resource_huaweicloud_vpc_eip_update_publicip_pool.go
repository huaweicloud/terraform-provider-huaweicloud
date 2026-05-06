package eip

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var eipUpdatePublicipPoolNonUpdatableParams = []string{
	"publicip_pool_id",
}

// @API EIP PUT /v3/{project_id}/eip/publicip-pools/{publicip_pool_id}
// @API EIP GET /v3/{project_id}/eip/publicip-pools/{publicip_pool_id}
func ResourceEipUpdatePublicipPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEipUpdatePublicipPoolCreate,
		ReadContext:   resourceEipUpdatePublicipPoolRead,
		UpdateContext: resourceEipUpdatePublicipPoolUpdate,
		DeleteContext: resourceEipUpdatePublicipPoolDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(eipUpdatePublicipPoolNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"publicip_pool_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"used": {
				Type:     schema.TypeInt,
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
			"billing_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"order_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"public_border_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shared": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allow_share_bandwidth_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func updateEipPublicipPool(client *golangsdk.ServiceClient, opt map[string]interface{}, poolId string) error {
	requestPath := client.Endpoint + "v3/{project_id}/eip/publicip-pools/{publicip_pool_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{publicip_pool_id}", poolId)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(opt),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)

	return err
}

func buildEipUpdatePublicipPoolBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return map[string]interface{}{
		"publicip_pool": bodyParams,
	}
}

func resourceEipUpdatePublicipPoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		product        = "vpc"
		publicipPoolId = d.Get("publicip_pool_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	if err := updateEipPublicipPool(client, buildEipUpdatePublicipPoolBodyParams(d), publicipPoolId); err != nil {
		return diag.Errorf("error updating EIP public IP pool in creation: %s", err)
	}

	d.SetId(publicipPoolId)

	return resourceEipUpdatePublicipPoolRead(ctx, d, meta)
}

func ReadEipUpdatePublicipPool(client *golangsdk.ServiceClient, poolId string) (interface{}, error) {
	requestPath := client.Endpoint + "v3/{project_id}/eip/publicip-pools/{publicip_pool_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{publicip_pool_id}", poolId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceEipUpdatePublicipPoolRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "vpc"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	respBody, err := ReadEipUpdatePublicipPool(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving EIP public IP pool")
	}

	publicipPool := utils.PathSearch("publicip_pool", respBody, nil)

	mErr := multierror.Append(
		d.Set("publicip_pool_id", utils.PathSearch("id", publicipPool, nil)),
		d.Set("name", utils.PathSearch("name", publicipPool, nil)),
		d.Set("description", utils.PathSearch("description", publicipPool, nil)),
		d.Set("status", utils.PathSearch("status", publicipPool, nil)),
		d.Set("type", utils.PathSearch("type", publicipPool, nil)),
		d.Set("project_id", utils.PathSearch("project_id", publicipPool, nil)),
		d.Set("size", utils.PathSearch("size", publicipPool, nil)),
		d.Set("used", utils.PathSearch("used", publicipPool, nil)),
		d.Set("created_at", utils.PathSearch("created_at", publicipPool, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", publicipPool, nil)),
		d.Set("billing_info", flattenPublicipPoolBillingInfo(publicipPool)),
		d.Set("public_border_group", utils.PathSearch("public_border_group", publicipPool, nil)),
		d.Set("shared", utils.PathSearch("shared", publicipPool, nil)),
		d.Set("tags", flattenPublicipPoolTags(publicipPool)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", publicipPool, nil)),
		d.Set("allow_share_bandwidth_types", utils.PathSearch("allow_share_bandwidth_types", publicipPool, nil)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting public IP pool fields: %s", err)
	}
	return nil
}

func flattenPublicipPoolBillingInfo(publicipPool interface{}) []interface{} {
	if publicipPool == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"order_id":   utils.PathSearch("billing_info.order_id", publicipPool, nil),
			"product_id": utils.PathSearch("billing_info.product_id", publicipPool, nil),
		},
	}
}

func flattenPublicipPoolTags(publicipPool interface{}) []interface{} {
	if publicipPool == nil {
		return nil
	}

	tagsArray := utils.PathSearch("tags", publicipPool, make([]interface{}, 0)).([]interface{})
	if len(tagsArray) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(tagsArray))
	for _, v := range tagsArray {
		result = append(result, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}

	return result
}

func buildUpdatePublicipPoolBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}

	// Passing in an existing name will trigger a name duplicate error.
	if d.HasChange("name") {
		bodyParams["name"] = d.Get("name")
	}

	return map[string]interface{}{
		"publicip_pool": bodyParams,
	}
}

func resourceEipUpdatePublicipPoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		product        = "vpc"
		publicipPoolId = d.Get("publicip_pool_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	if d.HasChanges("name", "description") {
		if err := updateEipPublicipPool(client, buildUpdatePublicipPoolBodyParams(d), publicipPoolId); err != nil {
			return diag.Errorf("error updating EIP public IP pool in update: %s", err)
		}
	}

	return resourceEipUpdatePublicipPoolRead(ctx, d, meta)
}

func resourceEipUpdatePublicipPoolDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to update a public IP pool.
Deleting this resource will not delete the actual public IP pool on the cloud, but will only remove the resource 
information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
