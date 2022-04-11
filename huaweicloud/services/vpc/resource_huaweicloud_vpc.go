package vpc

import (
	"context"
	"log"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/networking/v1/vpcs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3"
	v3vpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func ResourceVirtualPrivateCloudV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVirtualPrivateCloudCreate,
		ReadContext:   resourceVirtualPrivateCloudRead,
		UpdateContext: resourceVirtualPrivateCloudUpdate,
		DeleteContext: resourceVirtualPrivateCloudDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ //request and response parameters
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: utils.ValidateString64WithChinese,
			},
			"cidr": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: utils.ValidateCIDR,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secondary_cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: utils.ValidateCIDR,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"routes": {
				Type:       schema.TypeList,
				Computed:   true,
				Deprecated: "use huaweicloud_vpc_route_table data source to get all routes",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nexthop": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tags": common.TagsSchema(),
		},
	}
}

func resourceVirtualPrivateCloudCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	vpcClient, err := config.NetworkingV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud VPC client: %s", err)
	}

	createOpts := vpcs.CreateOpts{
		Name:        d.Get("name").(string),
		CIDR:        d.Get("cidr").(string),
		Description: d.Get("description").(string),
	}

	epsID := common.GetEnterpriseProjectID(d, config)
	if epsID != "" {
		createOpts.EnterpriseProjectID = epsID
	}

	n, err := vpcs.Create(vpcClient, createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud VPC: %s", err)
	}

	d.SetId(n.ID)
	log.Printf("[DEBUG] Vpc ID: %s", n.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForVpcActive(vpcClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return fmtp.DiagErrorf(
			"Error waiting for Vpc (%s) to become ACTIVE: %s",
			n.ID, stateErr)
	}

	//set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		vpcV2Client, err := config.NetworkingV2Client(region)
		if err != nil {
			return fmtp.DiagErrorf("Error creating Huaweicloud VPC client: %s", err)
		}
		taglist := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(vpcV2Client, "vpcs", n.ID, taglist).ExtractErr(); tagErr != nil {
			return fmtp.DiagErrorf("Error setting tags of VPC %q: %s", n.ID, tagErr)
		}
	}

	if v, ok := d.GetOk("secondary_cidr"); ok {
		v3Client, err := config.HcVpcV3Client(region)
		if err != nil {
			return diag.Errorf("error creating VPC v3 client: %s", err)
		}

		extendCidr := v.(string)
		if err := addSecondaryCIDR(v3Client, d.Id(), extendCidr); err != nil {
			return diag.Errorf("error adding VPC secondary CIDR: %s", err)
		}
	}

	return resourceVirtualPrivateCloudRead(ctx, d, meta)
}

// GetVpcById is a method to obtain vpc informations from special region through vpc ID.
func GetVpcById(config *config.Config, region, vpcId string) (*vpcs.Vpc, error) {
	client, err := config.NetworkingV1Client(region)
	if err != nil {
		return nil, err
	}

	return vpcs.Get(client, vpcId).Extract()
}

func resourceVirtualPrivateCloudRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	n, err := GetVpcById(config, config.GetRegion(d), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error obtain VPC information")
	}

	d.Set("name", n.Name)
	d.Set("cidr", n.CIDR)
	d.Set("description", n.Description)
	d.Set("enterprise_project_id", n.EnterpriseProjectID)
	d.Set("status", n.Status)
	d.Set("region", config.GetRegion(d))

	// save route tables
	routes := make([]map[string]interface{}, len(n.Routes))
	for i, rtb := range n.Routes {
		route := map[string]interface{}{
			"destination": rtb.DestinationCIDR,
			"nexthop":     rtb.NextHop,
		}
		routes[i] = route
	}
	d.Set("routes", routes)

	// save VirtualPrivateCloudV2 tags
	if vpcV2Client, err := config.NetworkingV2Client(config.GetRegion(d)); err == nil {
		if resourceTags, err := tags.Get(vpcV2Client, "vpcs", d.Id()).Extract(); err == nil {
			tagmap := utils.TagsToMap(resourceTags.Tags)
			if err := d.Set("tags", tagmap); err != nil {
				return fmtp.DiagErrorf("Error saving tags to state for VPC (%s): %s", d.Id(), err)
			}
		} else {
			log.Printf("[WARN] Error fetching tags of VPC (%s): %s", d.Id(), err)
		}
	} else {
		return fmtp.DiagErrorf("Error creating Huaweicloud VPC client: %s", err)
	}

	return nil
}

func resourceVirtualPrivateCloudUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	vpcClient, err := config.NetworkingV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud VPC client: %s", err)
	}

	vpcID := d.Id()
	if d.HasChanges("name", "cidr", "description") {
		updateOpts := vpcs.UpdateOpts{
			Name: d.Get("name").(string),
			CIDR: d.Get("cidr").(string),
		}
		if d.HasChange("description") {
			desc := d.Get("description").(string)
			updateOpts.Description = &desc
		}

		_, err = vpcs.Update(vpcClient, vpcID, updateOpts).Extract()
		if err != nil {
			return fmtp.DiagErrorf("Error updating Huaweicloud VPC: %s", err)
		}
	}

	//update tags
	if d.HasChange("tags") {
		vpcV2Client, err := config.NetworkingV2Client(region)
		if err != nil {
			return fmtp.DiagErrorf("Error creating Huaweicloud VPC client: %s", err)
		}

		tagErr := utils.UpdateResourceTags(vpcV2Client, d, "vpcs", vpcID)
		if tagErr != nil {
			return fmtp.DiagErrorf("Error updating tags of VPC %s: %s", vpcID, tagErr)
		}
	}

	if d.HasChange("secondary_cidr") {
		v3Client, err := config.HcVpcV3Client(region)
		if err != nil {
			return diag.Errorf("error creating VPC v3 client: %s", err)
		}

		old, new := d.GetChange("secondary_cidr")
		preExtendCidr := old.(string)
		newExtendCidr := new.(string)

		if preExtendCidr != "" {
			if err := removeSecondaryCIDR(v3Client, vpcID, preExtendCidr); err != nil {
				return diag.Errorf("error deleting VPC secondary CIDR: %s", err)
			}
		}
		if newExtendCidr != "" {
			if err := addSecondaryCIDR(v3Client, vpcID, newExtendCidr); err != nil {
				return diag.Errorf("error adding VPC secondary CIDR: %s", err)
			}
		}
	}

	return resourceVirtualPrivateCloudRead(ctx, d, meta)
}

func resourceVirtualPrivateCloudDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	config := meta.(*config.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud VPC client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForVpcDelete(vpcClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf("Error deleting Huaweicloud VPC %s: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func waitForVpcActive(vpcClient *golangsdk.ServiceClient, vpcId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := vpcs.Get(vpcClient, vpcId).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.Status == "OK" {
			return n, "ACTIVE", nil
		}

		//If vpc status is other than Ok, send error
		if n.Status == "DOWN" {
			return nil, "", fmtp.Errorf("Vpc status: '%s'", n.Status)
		}

		return n, n.Status, nil
	}
}

func waitForVpcDelete(vpcClient *golangsdk.ServiceClient, vpcId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := vpcs.Get(vpcClient, vpcId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully delete VPC %s", vpcId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = vpcs.Delete(vpcClient, vpcId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully delete VPC %s", vpcId)
				return r, "DELETED", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, "ACTIVE", nil
				}
			}
			return r, "ACTIVE", err
		}

		return r, "ACTIVE", nil
	}
}

func addSecondaryCIDR(client *client.VpcClient, vpcID, cidr string) error {
	reqBody := v3vpc.AddVpcExtendCidrRequestBody{
		Vpc: &v3vpc.AddExtendCidrOption{
			ExtendCidrs: []string{cidr},
		},
	}
	reqOpts := v3vpc.AddVpcExtendCidrRequest{
		VpcId: vpcID,
		Body:  &reqBody,
	}

	log.Printf("[DEBUG] add secondary CIDR %s into VPC %s", cidr, vpcID)
	_, err := client.AddVpcExtendCidr(&reqOpts)
	return err
}

func removeSecondaryCIDR(client *client.VpcClient, vpcID, preCidr string) error {
	reqBody := v3vpc.RemoveVpcExtendCidrRequestBody{
		Vpc: &v3vpc.RemoveExtendCidrOption{
			ExtendCidrs: []string{preCidr},
		},
	}
	reqOpts := v3vpc.RemoveVpcExtendCidrRequest{
		VpcId: vpcID,
		Body:  &reqBody,
	}

	log.Printf("[DEBUG] remove secondary CIDR %s from VPC %s", preCidr, vpcID)
	_, err := client.RemoveVpcExtendCidr(&reqOpts)
	return err
}
