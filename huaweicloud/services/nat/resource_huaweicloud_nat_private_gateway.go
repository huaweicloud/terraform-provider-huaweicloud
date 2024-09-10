package nat

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/nat/v3/gateways"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	PrivateSpecTypeSmall      string = "Small"
	PrivateSpecTypeMedium     string = "Medium"
	PrivateSpecTypeLarge      string = "Large"
	PrivateSpecTypeExtraLarge string = "Extra-Large"
)

// @API NAT POST /v3/{project_id}/private-nat/gateways
// @API NAT GET /v3/{project_id}/private-nat/gateways/{gateway_id}
// @API NAT PUT /v3/{project_id}/private-nat/gateways/{gateway_id}
// @API NAT DELETE /v3/{project_id}/private-nat/gateways/{gateway_id}
// @API NAT POST /v3/{project_id}/private-nat-gateways/{resource_id}/tags/action
func ResourcePrivateGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateGatewayCreate,
		ReadContext:   resourcePrivateGatewayRead,
		UpdateContext: resourcePrivateGatewayUpdate,
		DeleteContext: resourcePrivateGatewayDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the private NAT gateway is located.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The network ID of the subnet to which the private NAT gateway belongs.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the private NAT gateway.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the private NAT gateway.",
			},
			"spec": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The specification of the private NAT gateway.",
				ValidateFunc: validation.StringInSlice([]string{
					PrivateSpecTypeSmall,
					PrivateSpecTypeMedium,
					PrivateSpecTypeLarge,
					PrivateSpecTypeExtraLarge,
				}, false),
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ID of the enterprise project to which the private NAT gateway belongs.",
			},
			"tags": common.TagsSchema(),
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the private NAT gateway.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the private NAT gateway.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the private NAT gateway.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the VPC to which the private NAT gateway belongs.",
			},
		},
	}
}

func buildDownLinkVpcs(subnetId string) []gateways.DownLinkVpc {
	return []gateways.DownLinkVpc{
		{
			SubnetId: subnetId,
		},
	}
}

func resourcePrivateGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NatV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	opts := gateways.CreateOpts{
		Name:                d.Get("name").(string),
		DownLinkVpcs:        buildDownLinkVpcs(d.Get("subnet_id").(string)),
		Description:         d.Get("description").(string),
		Spec:                d.Get("spec").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
		Tags:                utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}

	log.Printf("[DEBUG] The create options of the private NAT gateway is: %#v", opts)
	resp, err := gateways.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating private NAT gateway: %s", err)
	}
	d.SetId(resp.ID)

	return resourcePrivateGatewayRead(ctx, d, meta)
}

func resourcePrivateGatewayRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	natClient, err := cfg.NatV3Client(region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	resp, err := gateways.Get(natClient, d.Id())
	if err != nil {
		// If the private NAT gateway does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving private NAT gateway")
	}
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("subnet_id", utils.PathSearch("[0].SubnetId", resp.DownLinkVpcs, nil)),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("spec", resp.Spec),
		d.Set("enterprise_project_id", resp.EnterpriseProjectId),
		d.Set("tags", utils.TagsToMap(resp.Tags)),
		d.Set("created_at", resp.CreatedAt),
		d.Set("updated_at", resp.UpdatedAt),
		d.Set("status", resp.Status),
		d.Set("vpc_id", utils.PathSearch("[0].VpcId", resp.DownLinkVpcs, nil)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving private NAT gateway fields: %s", err)
	}
	return nil
}

func resourcePrivateGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	natClient, err := cfg.NatV3Client(region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	gatewayId := d.Id()
	if d.HasChangeExcept("tags") {
		opts := gateways.UpdateOpts{
			Name:        d.Get("name").(string),
			Description: utils.String(d.Get("description").(string)),
		}
		if d.HasChange("spec") {
			opts.Spec = d.Get("spec").(string)
		}

		log.Printf("[DEBUG] The update options of the private NAT gateway is: %#v", opts)
		_, err = gateways.Update(natClient, gatewayId, opts)
		if err != nil {
			return diag.Errorf("error updating private NAT gateway: %s", err)
		}
	}

	if d.HasChange("tags") {
		err = utils.UpdateResourceTags(natClient, d, "private-nat-gateways", gatewayId)
		if err != nil {
			return diag.Errorf("error updating tags of the private NAT gateway (%s): %s", gatewayId, err)
		}
	}

	return resourcePrivateGatewayRead(ctx, d, meta)
}

func resourcePrivateGatewayDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NatV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	gatewayId := d.Id()
	err = gateways.Delete(client, gatewayId)
	if err != nil {
		// If the private NAT gateway does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting private NAT gateway")
	}

	return nil
}
