package apigateway

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/apigw/shared/v1/groups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API APIG PUT /v1.0/apigw/api-groups/{id}
// @API APIG DELETE /v1.0/apigw/api-groups/{id}
// @API APIG GET /v1.0/apigw/api-groups/{id}
// @API APIG POST /v1.0/apigw/api-groups
func ResourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAPIGatewayGroupCreate,
		ReadContext:   resourceAPIGatewayGroupRead,
		UpdateContext: resourceAPIGatewayGroupUpdate,
		DeleteContext: resourceAPIGatewayGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAPIGatewayGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	apigwClient, err := cfg.ApiGatewayV1Client(region)
	if err != nil {
		return diag.Errorf("error creating API Gateway client: %s", err)
	}

	createOpts := &groups.CreateOpts{
		Name:   d.Get("name").(string),
		Remark: d.Get("description").(string),
	}

	log.Printf("[DEBUG] create options: %#v", createOpts)
	v, err := groups.Create(apigwClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating API Gateway Group: %s", err)
	}

	// Store the ID now
	d.SetId(v.ID)

	return resourceAPIGatewayGroupRead(ctx, d, meta)
}

func resourceAPIGatewayGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	apigwClient, err := cfg.ApiGatewayV1Client(region)
	if err != nil {
		return diag.Errorf("error creating API Gateway client: %s", err)
	}

	v, err := groups.Get(apigwClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "API Gateway API")
	}

	log.Printf("[DEBUG] retrieved API Group %s: %+v", d.Id(), v)
	mErr := multierror.Append(
		d.Set("name", v.Name),
		d.Set("description", v.Remark),
		d.Set("status", v.Status),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAPIGatewayGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	apigwClient, err := cfg.ApiGatewayV1Client(region)
	if err != nil {
		return diag.Errorf("error creating API Gateway client: %s", err)
	}

	updateOpts := groups.UpdateOpts{
		Name:   d.Get("name").(string),
		Remark: d.Get("description").(string),
	}

	_, err = groups.Update(apigwClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating API Gateway Group: %s", err)
	}

	return resourceAPIGatewayGroupRead(ctx, d, meta)
}

func resourceAPIGatewayGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	apigwClient, err := cfg.ApiGatewayV1Client(region)
	if err != nil {
		return diag.Errorf("error creating API Gateway client: %s", err)
	}

	if err := groups.Delete(apigwClient, d.Id()).ExtractErr(); err != nil {
		return common.CheckDeletedDiag(d, err, "API apis")
	}

	return nil
}
