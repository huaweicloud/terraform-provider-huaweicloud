package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/networking/v1/vpcs"

	client "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3"
	v3vpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPC POST /v1/{project_id}/vpcs
// @API VPC GET /v1/{project_id}/vpcs/{id}
// @API VPC PUT /v1/{project_id}/vpcs/{id}
// @API VPC DELETE /v1/{project_id}/vpcs/{id}
// @API VPC PUT /v3/{project_id}/vpc/vpcs/{id}/add-extend-cidr
// @API VPC PUT /v3/{project_id}/vpc/vpcs/{id}/remove-extend-cidr
// @API VPC GET /v3/{project_id}/vpc/vpcs/{id}
// @API VPC POST /v2.0/{project_id}/vpcs/{id}/tags/action
// @API VPC DELETE /v2.0/{project_id}/vpcs/{id}/tags/action
// @API VPC GET /v2.0/{project_id}/vpcs/{id}/tags
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources-migrate
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources/filter
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

		Schema: map[string]*schema.Schema{ // request and response parameters
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
				Description:  "schema: Deprecated; use secondary_cidrs instead",
			},
			"secondary_cidrs": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"secondary_cidr"},
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: utils.ValidateCIDR,
				},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
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
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	v1Client, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC v1 client: %s", err)
	}

	createOpts := vpcs.CreateOpts{
		Name:        d.Get("name").(string),
		CIDR:        d.Get("cidr").(string),
		Description: d.Get("description").(string),
	}

	epsID := cfg.GetEnterpriseProjectID(d)
	if epsID != "" {
		createOpts.EnterpriseProjectID = epsID
	}

	n, err := vpcs.Create(v1Client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating VPC: %s", err)
	}

	d.SetId(n.ID)
	log.Printf("[DEBUG] VPC ID: %s", n.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForVpcActive(v1Client, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return diag.Errorf(
			"error waiting for Vpc (%s) to become ACTIVE: %s",
			n.ID, stateErr)
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		v2Client, err := cfg.NetworkingV2Client(region)
		if err != nil {
			return diag.Errorf("error creating VPC v2 client: %s", err)
		}
		taglist := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(v2Client, "vpcs", n.ID, taglist).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of VPC %q: %s", n.ID, tagErr)
		}
	}

	var extendCidrs []string
	if v, ok := d.GetOk("secondary_cidr"); ok {
		extendCidrs = []string{v.(string)}
	}
	if v, ok := d.GetOk("secondary_cidrs"); ok {
		extendCidrs = utils.ExpandToStringList(v.(*schema.Set).List())
	}

	if len(extendCidrs) > 0 {
		v3Client, err := cfg.HcVpcV3Client(region)
		if err != nil {
			return diag.Errorf("error creating VPC v3 client: %s", err)
		}

		if err := addSecondaryCIDR(v3Client, d.Id(), extendCidrs); err != nil {
			return diag.Errorf("error adding VPC secondary CIDRs: %s", err)
		}
	}

	return resourceVirtualPrivateCloudRead(ctx, d, meta)
}

// GetVpcById is a method to obtain vpc informations from special region through vpc ID.
func GetVpcById(conf *config.Config, region, vpcId string) (*vpcs.Vpc, error) {
	v1Client, err := conf.NetworkingV1Client(region)
	if err != nil {
		return nil, err
	}

	return vpcs.Get(v1Client, vpcId).Extract()
}

func resourceVirtualPrivateCloudRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	n, err := GetVpcById(conf, region, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error obtain VPC information")
	}

	d.Set("name", n.Name)
	d.Set("cidr", n.CIDR)
	d.Set("description", n.Description)
	d.Set("enterprise_project_id", n.EnterpriseProjectID)
	d.Set("status", n.Status)
	d.Set("region", region)

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
	v2Client, err := conf.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}
	if resourceTags, err := tags.Get(v2Client, "vpcs", d.Id()).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagmap); err != nil {
			return diag.Errorf("error saving tags to state for VPC (%s): %s", d.Id(), err)
		}
	} else {
		log.Printf("[WARN] error fetching tags of VPC (%s): %s", d.Id(), err)
	}

	// save VirtualPrivateCloudV3 extend_cidrs
	v3Client, err := conf.HcVpcV3Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	res, err := obtainV3VpcResp(v3Client, d.Id())
	if err != nil {
		return diag.Errorf("error retrieving VPC (%s) v3 detail: %s", d.Id(), err)
	}

	if val, ok := d.GetOk("secondary_cidr"); ok {
		for _, extendCidr := range res.Vpc.ExtendCidrs {
			if extendCidr == val {
				d.Set("secondary_cidr", extendCidr)
				break
			}
		}
	}
	d.Set("secondary_cidrs", res.Vpc.ExtendCidrs)

	return nil
}

func resourceVirtualPrivateCloudUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	v1Client, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}
	v3Client, err := cfg.HcVpcV3Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
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

		_, err = vpcs.Update(v1Client, vpcID, updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating VPC: %s", err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		v2Client, err := cfg.NetworkingV2Client(region)
		if err != nil {
			return diag.Errorf("error creating VPC v2 client: %s", err)
		}

		tagErr := utils.UpdateResourceTags(v2Client, d, "vpcs", vpcID)
		if tagErr != nil {
			return diag.Errorf("error updating tags of VPC %s: %s", vpcID, tagErr)
		}
	}

	if d.HasChange("secondary_cidr") {
		oldValue, newValue := d.GetChange("secondary_cidr")
		preExtendCidr := oldValue.(string)
		newExtendCidr := newValue.(string)
		if preExtendCidr != "" {
			preExtendCidrs := []string{preExtendCidr}
			if err := removeSecondaryCIDR(v3Client, vpcID, preExtendCidrs); err != nil {
				return diag.Errorf("error deleting VPC secondary CIDR: %s", err)
			}
		}
		if newExtendCidr != "" {
			newExtendCidrs := []string{newExtendCidr}
			if err := addSecondaryCIDR(v3Client, vpcID, newExtendCidrs); err != nil {
				return diag.Errorf("error adding VPC secondary CIDR: %s", err)
			}
		}
	}

	if d.HasChanges("secondary_cidrs") {
		oldRaws, newRaws := d.GetChange("secondary_cidrs")
		preExtendCidrs := utils.ExpandToStringListBySet(oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set)))
		newExtendCidrs := utils.ExpandToStringListBySet(newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set)))
		if len(preExtendCidrs) > 0 {
			if err := removeSecondaryCIDR(v3Client, vpcID, preExtendCidrs); err != nil {
				return diag.Errorf("error deleting VPC secondary CIDRs: %s", err)
			}
		}
		if len(newExtendCidrs) > 0 {
			if err := addSecondaryCIDR(v3Client, vpcID, newExtendCidrs); err != nil {
				return diag.Errorf("error adding VPC secondary CIDRs: %s", err)
			}
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   vpcID,
			ResourceType: "vpcs",
			RegionId:     region,
			ProjectId:    v1Client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceVirtualPrivateCloudRead(ctx, d, meta)
}

func resourceVirtualPrivateCloudDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	v1Client, err := conf.NetworkingV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForVpcDelete(v1Client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting VPC %s: %s", d.Id(), err)
	}

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

		// If vpc status is other than Ok, send error
		if n.Status == "DOWN" {
			return nil, "", fmt.Errorf("VPC status: '%s'", n.Status)
		}

		return n, n.Status, nil
	}
}

func waitForVpcDelete(vpcClient *golangsdk.ServiceClient, vpcId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := vpcs.Get(vpcClient, vpcId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] successfully delete VPC %s", vpcId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = vpcs.Delete(vpcClient, vpcId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] successfully delete VPC %s", vpcId)
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

func addSecondaryCIDR(v3Client *client.VpcClient, vpcID string, cidrs []string) error {
	reqBody := v3vpc.AddVpcExtendCidrRequestBody{
		Vpc: &v3vpc.AddExtendCidrOption{
			ExtendCidrs: cidrs,
		},
	}
	reqOpts := v3vpc.AddVpcExtendCidrRequest{
		VpcId: vpcID,
		Body:  &reqBody,
	}

	log.Printf("[DEBUG] add secondary CIDRs %s into VPC %s", cidrs, vpcID)
	_, err := v3Client.AddVpcExtendCidr(&reqOpts)
	return err
}

func removeSecondaryCIDR(v3Client *client.VpcClient, vpcID string, preCidrs []string) error {
	reqBody := v3vpc.RemoveVpcExtendCidrRequestBody{
		Vpc: &v3vpc.RemoveExtendCidrOption{
			ExtendCidrs: preCidrs,
		},
	}
	reqOpts := v3vpc.RemoveVpcExtendCidrRequest{
		VpcId: vpcID,
		Body:  &reqBody,
	}

	log.Printf("[DEBUG] remove secondary CIDRs %s from VPC %s", preCidrs, vpcID)
	_, err := v3Client.RemoveVpcExtendCidr(&reqOpts)
	return err
}

func obtainV3VpcResp(v3Client *client.VpcClient, vpcID string) (*v3vpc.ShowVpcResponse, error) {
	reqOpts := v3vpc.ShowVpcRequest{
		VpcId: vpcID,
	}
	res, err := v3Client.ShowVpc(&reqOpts)
	if err != nil {
		return nil, err
	}

	return res, nil
}
