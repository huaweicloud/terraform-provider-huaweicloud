package apigateway

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/shared/v1/environments"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG POST /v1.0/apigw/envs
// @API APIG GET /v1.0/apigw/envs
// @API APIG PUT /v1.0/apigw/envs/{id}
// @API APIG DELETE /v1.0/apigw/envs/{id}
func ResourceEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentCreate,
		ReadContext:   resourceEnvironmentRead,
		UpdateContext: resourceEnvironmentUpdate,
		DeleteContext: resourceEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	apigClient, err := cfg.ApiGatewayV1Client(region)
	if err != nil {
		return diag.Errorf("error creating shared API gateway client (v1) : %s", err)
	}

	createOpts := environments.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
	log.Printf("[DEBUG] Create API gateway shared options: %v", createOpts)

	rst, err := environments.Create(apigClient, createOpts)
	if err != nil {
		return diag.Errorf("error creating shared API gateway environment: %s", err)
	}
	d.SetId(rst.Id)

	return resourceEnvironmentRead(ctx, d, meta)
}

func getEnvironment(client *golangsdk.ServiceClient, envName, envID string) (*environments.Environment, error) {
	envs, err := environments.List(client, environments.ListOpts{
		EnvName: envName,
	})
	if err != nil {
		return nil, err
	}

	for i, v := range envs {
		if v.Id == envID {
			return &envs[i], nil
		}
	}
	return nil, golangsdk.ErrDefault404{}
}

func resourceEnvironmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	apigClient, err := cfg.ApiGatewayV1Client(region)
	if err != nil {
		return diag.Errorf("error creating shared API gateway client (v1) : %s", err)
	}

	n, err := getEnvironment(apigClient, d.Get("name").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving shared APIG environment")
	}
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", n.Name),
		d.Set("description", n.Description),
		d.Set("created_at", n.CreateTime),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	apigClient, err := cfg.ApiGatewayV1Client(region)
	if err != nil {
		return diag.Errorf("error creating shared API gateway client (v1) : %s", err)
	}

	updateOpts := environments.UpdateOpts{
		Name:        d.Get("name").(string),
		Description: utils.String(d.Get("description").(string)),
	}

	_, err = environments.Update(apigClient, d.Id(), updateOpts)
	if err != nil {
		return diag.Errorf("error updating shared API gateway environment: %s", err)
	}

	return resourceEnvironmentRead(ctx, d, meta)
}

func resourceEnvironmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	apigClient, err := cfg.ApiGatewayV1Client(region)
	if err != nil {
		return diag.Errorf("error creating shared API gateway client (v1) : %s", err)
	}

	err = environments.Delete(apigClient, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting shared API gateway environment")
	}
	return nil
}
