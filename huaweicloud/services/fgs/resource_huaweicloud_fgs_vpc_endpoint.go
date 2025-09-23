package fgs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var vpcEndpointNonUpdatableParams = []string{"vpc_id", "subnet_id", "flavor", "xrole"}

// @API FunctionGraph POST /v2/{project_id}/fgs/vpc-endpoint
// @API FunctionGraph DELETE /v2/{project_id}/fgs/vpc-endpoint/{vpc_id}/{subnet_id}
func ResourceFgsVpcEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFgsVpcEndpointCreate,
		ReadContext:   resourceFgsVpcEndpointRead,
		UpdateContext: resourceFgsVpcEndpointUpdate,
		DeleteContext: resourceFgsVpcEndpointDelete,

		CustomizeDiff: config.FlexibleForceNew(vpcEndpointNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the VPC endpoint is located.`,
			},

			// Required parameter(s).
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the VPC to which the VPC endpoint belongs.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the subnet to which the VPC endpoint belongs.`,
			},

			// Optional parameter(s).
			"flavor": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The flavor of the VPC endpoint.`,
			},
			"xrole": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The IAM agency name of the VPC endpoint.`,
			},

			// Attributes.
			"endpoints": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of IP addresses of the VPC endpoint.`,
			},
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain address of the VPC endpoint.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildFgsVpcEndpointBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"vpc_id":    d.Get("vpc_id"),
		"subnet_id": d.Get("subnet_id"),
		"flavor":    utils.ValueIgnoreEmpty(d.Get("flavor")),
		"xrole":     utils.ValueIgnoreEmpty(d.Get("xrole")),
	}
}

func resourceFgsVpcEndpointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	httpUrl := "v2/{project_id}/fgs/vpc-endpoint"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildFgsVpcEndpointBodyParams(d),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph VPC endpoint: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("endpoints", utils.PathSearch("endpoints", respBody, nil)),
		d.Set("address", utils.PathSearch("address", respBody, nil)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving VPC endpoint fields: %s", err)
	}
	return resourceFgsVpcEndpointRead(ctx, d, meta)
}

func resourceFgsVpcEndpointRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceFgsVpcEndpointUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceFgsVpcEndpointRead(ctx, d, meta)
}

func resourceFgsVpcEndpointDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	httpUrl := "v2/{project_id}/fgs/vpc-endpoint/{vpc_id}/{subnet_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{vpc_id}", d.Get("vpc_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{subnet_id}", d.Get("subnet_id").(string))

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// Deleting a not exist VPC endpoint will not report an error.
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting FunctionGraph VPC endpoint: %s", err)
	}

	return nil
}
