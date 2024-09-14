package nat

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/nat/v2/gateways"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type (
	PublicSpecType string
)

const (
	PublicSpecTypeSmall      PublicSpecType = "1"
	PublicSpecTypeMedium     PublicSpecType = "2"
	PublicSpecTypeLarge      PublicSpecType = "3"
	PublicSpecTypeExtraLarge PublicSpecType = "4"
)

// @API NAT POST /v2/{project_id}/nat_gateways
// @API NAT GET /v2/{project_id}/nat_gateways/{nat_gateway_id}
// @API NAT PUT /v2/{project_id}/nat_gateways/{nat_gateway_id}
// @API NAT DELETE /v2/{project_id}/nat_gateways/{nat_gateway_id}
// @API NAT POST /v2.0/{project_id}/nat_gateways/{nat_gateway_id}/tags/action
// @API NAT GET /v2.0/{project_id}/nat_gateways/{nat_gateway_id}/tags
func ResourcePublicGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePublicGatewayCreate,
		ReadContext:   resourcePublicGatewayRead,
		UpdateContext: resourcePublicGatewayUpdate,
		DeleteContext: resourcePublicGatewayDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the NAT gateway is located.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the VPC to which the NAT gateway belongs.",
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "The network ID of the downstream interface (the next hop of the DVR) " +
					"of the NAT gateway.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The NAT gateway name.",
			},
			"spec": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(PublicSpecTypeSmall),
					string(PublicSpecTypeMedium),
					string(PublicSpecTypeLarge),
					string(PublicSpecTypeExtraLarge),
				}, false),
				Description: "The specification of the NAT gateway.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the NAT gateway.",
			},
			"ngport_ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The private IP address of the NAT gateway.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The enterprise project ID of the NAT gateway.",
			},
			"tags": common.TagsSchema(),
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the NAT gateway.",
			},
		},
	}
}

func publicGatewayStateRefreshFunc(client *golangsdk.ServiceClient, gatewayId string,
	targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := gateways.Get(client, gatewayId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return resp, "COMPLETED", nil
			}
			return resp, "", err
		}

		if utils.StrSliceContains([]string{"INACTIVE"}, resp.Status) {
			return resp, "", fmt.Errorf("unexpect status (%s)", resp.Status)
		}
		if utils.StrSliceContains(targets, resp.Status) {
			return resp, "COMPLETED", nil
		}
		return resp, "PENDING", nil
	}
}

func resourcePublicGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NatGatewayClient(region)
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	opts := gateways.CreateOpts{
		Name:                d.Get("name").(string),
		VpcId:               d.Get("vpc_id").(string),
		InternalNetworkId:   d.Get("subnet_id").(string),
		Spec:                d.Get("spec").(string),
		Description:         d.Get("description").(string),
		NgportIpAddress:     d.Get("ngport_ip_address").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}
	resp, err := gateways.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating NAT gateway: %s", err)
	}
	d.SetId(resp.ID)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      publicGatewayStateRefreshFunc(client, d.Id(), []string{"ACTIVE"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if gatewayTags, ok := d.GetOk("tags"); ok {
		networkClient, err := cfg.NetworkingV2Client(region)
		if err != nil {
			return diag.Errorf("error creating VPC v2.0 client: %s", err)
		}
		taglist := utils.ExpandResourceTags(gatewayTags.(map[string]interface{}))
		err = tags.Create(networkClient, "nat_gateways", d.Id(), taglist).ExtractErr()
		if err != nil {
			return diag.Errorf("error setting tags to the NAT gateway: %s", err)
		}
	}
	return resourcePublicGatewayRead(ctx, d, meta)
}

func resourcePublicGatewayRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NatGatewayClient(region)
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}
	networkClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC v2.0 client: %s", err)
	}

	gatewayId := d.Id()
	resp, err := gateways.Get(client, gatewayId)
	if err != nil {
		// If the NAT gateway does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving NAT Gateway")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", resp.Name),
		d.Set("spec", resp.Spec),
		d.Set("vpc_id", resp.RouterId),
		d.Set("subnet_id", resp.InternalNetworkId),
		d.Set("description", resp.Description),
		d.Set("ngport_ip_address", resp.NgportIpAddress),
		d.Set("enterprise_project_id", resp.EnterpriseProjectId),
		d.Set("status", resp.Status),
	)
	gatewayTags, err := tags.Get(networkClient, "nat_gateways", gatewayId).Extract()
	if err != nil {
		log.Printf("[WARN] Error getting gateway tags: %s", err)
	} else {
		mErr = multierror.Append(mErr, d.Set("tags", utils.TagsToMap(gatewayTags.Tags)))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving NAT gateway fields: %s", err)
	}
	return nil
}

func resourcePublicGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		gatewayId = d.Id()
	)

	if d.HasChangeExcept("tags") {
		client, err := cfg.NatGatewayClient(region)
		if err != nil {
			return diag.Errorf("error creating NAT v2 client: %s", err)
		}
		opts := gateways.UpdateOpts{
			Name:        d.Get("name").(string),
			Spec:        d.Get("spec").(string),
			Description: utils.String(d.Get("description").(string)),
		}
		_, err = gateways.Update(client, gatewayId, opts)
		if err != nil {
			return diag.Errorf("error updating NAT gateway (%s): %s", gatewayId, err)
		}
		stateConf := &resource.StateChangeConf{
			Pending:      []string{"PENDING"},
			Target:       []string{"COMPLETED"},
			Refresh:      publicGatewayStateRefreshFunc(client, gatewayId, []string{"ACTIVE"}),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        5 * time.Second,
			PollInterval: 10 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		networkClient, err := cfg.NetworkingV2Client(region)
		if err != nil {
			return diag.Errorf("error creating VPC v2.0 client: %s", err)
		}
		err = utils.UpdateResourceTags(networkClient, d, "nat_gateways", gatewayId)
		if err != nil {
			return diag.Errorf("error updating tags of the NAT gateway: %s", err)
		}
	}

	return resourcePublicGatewayRead(ctx, d, meta)
}

func resourcePublicGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NatGatewayClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	gatewayId := d.Id()
	err = gateways.Delete(client, gatewayId)
	if err != nil {
		// If the NAT gateway does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "err deleting NAT gateway")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      publicGatewayStateRefreshFunc(client, gatewayId, nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
