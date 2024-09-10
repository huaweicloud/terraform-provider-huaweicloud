package vpc

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	vpc_model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPC POST /v3/{project_id}/vpc/address-groups
// @API VPC DELETE /v3/{project_id}/vpc/address-groups/{address_group_id}/force
// @API VPC GET /v3/{project_id}/vpc/address-groups/{address_group_id}
// @API VPC PUT /v3/{project_id}/vpc/address-groups/{address_group_id}
// @API VPC DELETE /v3/{project_id}/vpc/address-groups/{address_group_id}
func ResourceVpcAddressGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcAddressGroupCreate,
		UpdateContext: resourceVpcAddressGroupUpdate,
		DeleteContext: resourceVpcAddressGroupDelete,
		ReadContext:   resourceVpcAddressGroupRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"addresses": {
				// the addresses will be sorted by cloud
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 20,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ip_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      4,
				ValidateFunc: validation.IntInSlice([]int{4, 6}),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_capacity": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"force_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceVpcAddressGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.HcVpcV3Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	ipSet := utils.ExpandToStringListBySet(d.Get("addresses").(*schema.Set))

	addressGroupBody := &vpc_model.CreateAddressGroupOption{
		Name:                d.Get("name").(string),
		IpSet:               &ipSet,
		IpVersion:           int32(d.Get("ip_version").(int)),
		Description:         utils.StringIgnoreEmpty(d.Get("description").(string)),
		EnterpriseProjectId: utils.StringIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		MaxCapacity:         utils.Int32IgnoreEmpty(int32(d.Get("max_capacity").(int))),
	}

	createOpts := &vpc_model.CreateAddressGroupRequest{
		Body: &vpc_model.CreateAddressGroupRequestBody{
			AddressGroup: addressGroupBody,
		},
	}

	log.Printf("[DEBUG] Create VPC address group options: %#v", addressGroupBody)
	response, err := client.CreateAddressGroup(createOpts)
	if err != nil {
		return diag.Errorf("error creating VPC address group: %s", err)
	}

	d.SetId(response.AddressGroup.Id)
	return resourceVpcAddressGroupRead(ctx, d, meta)
}

func resourceVpcAddressGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcVpcV3Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	request := &vpc_model.ShowAddressGroupRequest{
		AddressGroupId: d.Id(),
	}

	response, err := client.ShowAddressGroup(request)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error fetching VPC address group")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", response.AddressGroup.Name),
		d.Set("description", response.AddressGroup.Description),
		d.Set("addresses", response.AddressGroup.IpSet),
		d.Set("ip_version", response.AddressGroup.IpVersion),
		d.Set("max_capacity", response.AddressGroup.MaxCapacity),
		d.Set("enterprise_project_id", response.AddressGroup.EnterpriseProjectId),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceVpcAddressGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcVpcV3Client(c.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	addressGroupBody := &vpc_model.UpdateAddressGroupOption{}

	if d.HasChange("name") {
		groupName := d.Get("name").(string)
		addressGroupBody.Name = &groupName
	}

	if d.HasChange("description") {
		groupDescription := d.Get("description").(string)
		addressGroupBody.Description = &groupDescription
	}

	if d.HasChange("addresses") {
		rawAddresses := d.Get("addresses").(*schema.Set).List()
		ipSet := make([]string, len(rawAddresses))
		for i, value := range rawAddresses {
			ipSet[i] = value.(string)
		}
		addressGroupBody.IpSet = &ipSet
	}

	if d.HasChange("max_capacity") {
		addressGroupBody.MaxCapacity = utils.Int32IgnoreEmpty(int32(d.Get("max_capacity").(int)))
	}

	updateOpts := &vpc_model.UpdateAddressGroupRequest{
		AddressGroupId: d.Id(),
		Body: &vpc_model.UpdateAddressGroupRequestBody{
			AddressGroup: addressGroupBody,
		},
	}

	log.Printf("[DEBUG] Update VPC address group options: %#v", addressGroupBody)
	_, err = client.UpdateAddressGroup(updateOpts)
	if err != nil {
		return diag.Errorf("error updating VPC address group: %s", err)
	}

	return resourceVpcAddressGroupRead(ctx, d, meta)
}

func resourceVpcAddressGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcVpcV3Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	if !d.Get("force_destroy").(bool) {
		request := &vpc_model.DeleteAddressGroupRequest{
			AddressGroupId: d.Id(),
		}
		_, err = client.DeleteAddressGroup(request)
	} else {
		req := &vpc_model.DeleteIpAddressGroupForceRequest{
			AddressGroupId: d.Id(),
		}
		_, err = client.DeleteIpAddressGroupForce(req)
	}

	if err != nil {
		return diag.Errorf("error deleting VPC address group: %s", err)
	}

	return nil
}
