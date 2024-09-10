package sms

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v1/security/securitygroups"
	"github.com/chnsz/golangsdk/openstack/networking/v1/subnets"
	"github.com/chnsz/golangsdk/openstack/networking/v1/vpcs"
	"github.com/chnsz/golangsdk/openstack/sms/v3/templates"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var AutoCreate string = "autoCreate"

// ResourceServerTemplate is the impl of huaweicloud_sms_server_template
// @API SMS POST /v3/vm/templates
// @API SMS GET /v3/vm/templates/{id}
// @API SMS PUT /v3/vm/templates/{id}
// @API SMS DELETE /v3/vm/templates/{id}
// @API VPC GET /v1/{project_id}/security-groups/{security_group_id}
// @API VPC GET /v1/{project_id}/subnets/{subnet_id}
// @API VPC GET /v1/{project_id}/vpcs/{vpc_id}
func ResourceServerTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerTemplateCreate,
		ReadContext:   resourceServerTemplateRead,
		UpdateContext: resourceServerTemplateUpdate,
		DeleteContext: resourceServerTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"project_id"},
				Computed:     true,
			},
			"project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"region"},
				Computed:     true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "SAS",
			},
			"flavor": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bandwidth_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"target_server_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildVpcOpts(client *golangsdk.ServiceClient, vpcID string) (*templates.VpcRequest, error) {
	if vpcID != "" && vpcID != AutoCreate {
		rst, err := vpcs.Get(client, vpcID).Extract()
		if err != nil {
			return nil, fmt.Errorf("failed to get the VPC %s: %s", vpcID, err)
		}

		return &templates.VpcRequest{
			Id:   rst.ID,
			Name: rst.Name,
		}, nil
	}

	defaultRequest := templates.VpcRequest{
		Id:   AutoCreate,
		Name: AutoCreate,
	}
	return &defaultRequest, nil
}

func buildNicsOpts(client *golangsdk.ServiceClient, nics []string) ([]templates.NicRequest, error) {
	if len(nics) == 0 {
		return []templates.NicRequest{
			{
				Id:   AutoCreate,
				Name: AutoCreate,
			},
		}, nil
	}

	request := make([]templates.NicRequest, len(nics))
	for i, subnetID := range nics {
		if subnetID == AutoCreate {
			request[i] = templates.NicRequest{
				Id:   AutoCreate,
				Name: AutoCreate,
			}
		} else {
			rst, err := subnets.Get(client, subnetID).Extract()
			if err != nil {
				return nil, fmt.Errorf("failed to get the subnet %s: %s", subnetID, err)
			}

			request[i] = templates.NicRequest{
				Id:   rst.ID,
				Name: rst.Name,
				Cidr: rst.CIDR,
			}
		}

	}
	return request, nil
}

func buildSecGroupOpts(client *golangsdk.ServiceClient, sgs []string) ([]templates.SgRequest, error) {
	if len(sgs) == 0 {
		return []templates.SgRequest{
			{
				Id:   AutoCreate,
				Name: AutoCreate,
			},
		}, nil
	}

	request := make([]templates.SgRequest, len(sgs))
	for i, sgID := range sgs {
		if sgID == AutoCreate {
			request[i] = templates.SgRequest{
				Id:   AutoCreate,
				Name: AutoCreate,
			}
		} else {
			rst, err := securitygroups.Get(client, sgID).Extract()
			if err != nil {
				return nil, fmt.Errorf("failed to get the security group %s: %s", sgID, err)
			}

			request[i] = templates.SgRequest{
				Id:   rst.ID,
				Name: rst.Name,
			}
		}

	}
	return request, nil
}

func buildServerTemplateParameters(d *schema.ResourceData, cfg *config.Config) (*templates.TemplateOpts, error) {
	region := cfg.GetRegion(d)
	vpcClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC client: %s", err)
	}

	projectID := d.Get("project_id").(string)
	if projectID == "" {
		// get project ID from config
		projectID = cfg.RegionProjectIDMap[region]
	}

	vpcOpts, err := buildVpcOpts(vpcClient, d.Get("vpc_id").(string))
	if err != nil {
		return nil, err
	}

	sbunetIDs := utils.ExpandToStringList(d.Get("subnet_ids").([]interface{}))
	nicsOpts, err := buildNicsOpts(vpcClient, sbunetIDs)
	if err != nil {
		return nil, err
	}

	secGroupIDs := utils.ExpandToStringList(d.Get("security_group_ids").([]interface{}))
	secGroupOpts, err := buildSecGroupOpts(vpcClient, secGroupIDs)
	if err != nil {
		return nil, err
	}

	name := d.Get("name").(string)
	var targetName string
	if v, ok := d.GetOk("target_server_name"); ok {
		targetName = v.(string)
	} else {
		targetName = name
	}

	createOpts := templates.TemplateOpts{
		IsTemplate:       utils.Bool(true),
		Region:           region,
		ProjectID:        projectID,
		Name:             name,
		TargetServerName: targetName,
		AvailabilityZone: d.Get("availability_zone").(string),
		VolumeType:       d.Get("volume_type").(string),
		Flavor:           d.Get("flavor").(string),
		Vpc:              vpcOpts,
		Nics:             nicsOpts,
		SecurityGroups:   secGroupOpts,
	}

	if v, ok := d.GetOk("bandwidth_size"); ok {
		createOpts.PublicIP = &templates.EipRequest{
			Type:          "5_bgp",
			BandwidthSize: v.(int),
		}
	}

	return &createOpts, nil
}

func resourceServerTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	smsClient, err := config.SmsV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	createOpts, err := buildServerTemplateParameters(d, config)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	id, err := templates.Create(smsClient, createOpts)
	if err != nil {
		return diag.Errorf("error creating SMS server template: %s", err)
	}

	d.SetId(id)
	return resourceServerTemplateRead(ctx, d, meta)
}

func flattenSubnetIDs(nics []templates.NicObject) []string {
	subnets := make([]string, len(nics))
	for i, nic := range nics {
		subnets[i] = nic.Id
	}
	return subnets
}

func flattenSecGroupIDs(groups []templates.SgObject) []string {
	results := make([]string, len(groups))
	for i, group := range groups {
		results[i] = group.Id
	}
	return results
}

func resourceServerTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	smsClient, err := config.SmsV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	temp, err := templates.Get(smsClient, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error fetching SMS server template")
	}

	log.Printf("[DEBUG] Retrieved SMS server template %s: %+v", d.Id(), temp)
	mErr := multierror.Append(
		d.Set("name", temp.Name),
		d.Set("region", temp.Region),
		d.Set("project_id", temp.Projectid),
		d.Set("availability_zone", temp.AvailabilityZone),
		d.Set("target_server_name", temp.TargetServerName),
		d.Set("flavor", temp.Flavor),
		d.Set("volume_type", temp.Volumetype),
		d.Set("bandwidth_size", temp.PublicIP.BandwidthSize),
		d.Set("vpc_id", temp.Vpc.Id),
		d.Set("vpc_name", temp.Vpc.Name),
		d.Set("subnet_ids", flattenSubnetIDs(temp.Nics)),
		d.Set("security_group_ids", flattenSecGroupIDs(temp.SecurityGroups)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting SMS server template fields: %s", err)
	}

	return nil
}

func resourceServerTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	smsClient, err := config.SmsV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	updateOpts, err := buildServerTemplateParameters(d, config)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Update Options: %#v", updateOpts)
	err = templates.Update(smsClient, d.Id(), updateOpts).ExtractErr()
	if err != nil {
		return diag.Errorf("error updating SMS server template: %s", err)
	}

	return resourceServerTemplateRead(ctx, d, meta)
}

func resourceServerTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	smsClient, err := config.SmsV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	err = templates.Delete(smsClient, d.Id()).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SMS server template")
	}

	return nil
}
