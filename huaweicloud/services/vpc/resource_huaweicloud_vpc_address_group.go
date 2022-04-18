package vpc

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc_model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceVpcAddressGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcVpcV3Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud VPC client: %s", err)
	}

	rawAddresses := d.Get("addresses").(*schema.Set).List()
	ipSet := make([]string, len(rawAddresses))
	for i, value := range rawAddresses {
		ipSet[i] = value.(string)
	}

	addressGroupBody := &vpc_model.CreateAddressGroupOption{
		Name:      d.Get("name").(string),
		IpSet:     &ipSet,
		IpVersion: 4,
	}
	if v, ok := d.GetOk("description"); ok {
		desc := v.(string)
		addressGroupBody.Description = &desc
	}

	createOpts := &vpc_model.CreateAddressGroupRequest{
		Body: &vpc_model.CreateAddressGroupRequestBody{
			AddressGroup: addressGroupBody,
		},
	}

	logp.Printf("[DEBUG] Create VPC address group options: %#v", addressGroupBody)
	response, err := client.CreateAddressGroup(createOpts)
	if err != nil {
		return fmtp.DiagErrorf("Error creating VPC address group: %s", err)
	}

	d.SetId(response.AddressGroup.Id)
	return resourceVpcAddressGroupRead(ctx, d, meta)
}

func resourceVpcAddressGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcVpcV3Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud VPC client: %s", err)
	}

	request := &vpc_model.ShowAddressGroupRequest{
		AddressGroupId: d.Id(),
	}

	response, err := client.ShowAddressGroup(request)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error fetching VPC address group")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", response.AddressGroup.Name),
		d.Set("description", response.AddressGroup.Description),
		d.Set("addresses", response.AddressGroup.IpSet),
		d.Set("ip_version", response.AddressGroup.IpVersion),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("Error saving VPC address group: %s", err)
	}

	return nil
}

func resourceVpcAddressGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcVpcV3Client(c.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud VPC client: %s", err)
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

	updateOpts := &vpc_model.UpdateAddressGroupRequest{
		AddressGroupId: d.Id(),
		Body: &vpc_model.UpdateAddressGroupRequestBody{
			AddressGroup: addressGroupBody,
		},
	}

	logp.Printf("[DEBUG] Update VPC address group options: %#v", addressGroupBody)
	_, err = client.UpdateAddressGroup(updateOpts)
	if err != nil {
		return fmtp.DiagErrorf("Error updating VPC address group: %s", err)
	}

	return resourceVpcAddressGroupRead(ctx, d, meta)
}

func resourceVpcAddressGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcVpcV3Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud VPC client: %s", err)
	}

	request := &vpc_model.DeleteAddressGroupRequest{
		AddressGroupId: d.Id(),
	}

	_, err = client.DeleteAddressGroup(request)
	if err != nil {
		return fmtp.DiagErrorf("Error deleting VPC address group: %s", err)
	}

	return nil
}
