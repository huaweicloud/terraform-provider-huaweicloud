package huaweicloud

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/apigw/v2/instances"
	"github.com/huaweicloud/golangsdk/openstack/networking/v1/eips"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

type Refresh struct {
	Pending, Target                          []string
	Delay, Timeout, MinTimeout, PollInterval time.Duration
}

func ResourceApigInstanceV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceApigInstanceV2Create,
		Read:   resourceApigInstanceV2Read,
		Update: resourceApigInstanceV2Update,
		Delete: resourceApigInstanceV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
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
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^([\u4e00-\u9fa5A-Za-z][\u4e00-\u9fa5A-Za-z-_0-9]{2,63})$"),
					"The name contains of 3 to 64 characters, starting with a letter. Only letters, digits, "+
						"hyphens (-) and underscore (_) are allowed."),
			},
			"edition": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"BASIC", "PROFESSIONAL", "ENTERPRISE", "PLATINUM",
				}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"available_zone_ids": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"bandwidth_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 2000),
			},
			"eip_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"maintain_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^(02|06|10|14|18|22):00:00$"),
					"The start-time format of maintenance window is not 'xx:00:00' or "+
						"the hour is not 02, 06, 10, 14, 18 or 22."),
			},
			"ingress_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"egress_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"supported_features": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpc_ingress_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildMaintainEndTime(maintainStart string) (string, error) {
	regex := regexp.MustCompile("^(02|06|10|14|18|22):00:00$")
	isMatched := regex.MatchString(maintainStart)
	if !isMatched {
		return "", fmt.Errorf("The start-time format of maintenance window is not 'xx:00:00' or " +
			"the hour is not 02, 06, 10, 14, 18 or 22.")
	}
	result := regex.FindStringSubmatch(maintainStart)
	if len(result) < 2 {
		return "", fmt.Errorf("The hour is missing")
	}
	num, err := strconv.Atoi(result[1])
	if err != nil {
		return "", fmt.Errorf("The number (%s) cannot be converted to string", result[1])
	}
	return fmt.Sprintf("%02d:00:00", (num+4)%24), nil
}

func buildApigAvailableZoneIds(d *schema.ResourceData) []string {
	ids := d.Get("available_zone_ids").([]interface{})
	result := make([]string, len(ids))
	for i, v := range ids {
		result[i] = v.(string)
	}
	return result
}

func buildApigInstanceParameters(d *schema.ResourceData, config *config.Config) (instances.CreateOpts, error) {
	opt := instances.CreateOpts{
		Name:                d.Get("name").(string),
		Edition:             d.Get("edition").(string),
		VpcId:               d.Get("vpc_id").(string),
		SubnetId:            d.Get("subnet_id").(string),
		SecurityGroupId:     d.Get("security_group_id").(string),
		Description:         d.Get("description").(string),
		EipId:               d.Get("eip_id").(string),
		BandwidthSize:       d.Get("bandwidth_size").(int), // Bandwidth 0 means turn off the egress access.
		EnterpriseProjectId: GetEnterpriseProjectID(d, config),
		AvailableZoneIds:    buildApigAvailableZoneIds(d),
	}
	if v, ok := d.GetOk("maintain_time"); ok {
		startTime := v.(string)
		opt.MaintainBegin = startTime
		endTime, err := buildMaintainEndTime(startTime)
		if err != nil {
			return opt, err
		}
		opt.MaintainEnd = endTime
	}

	return opt, nil
}

func watiForApigInstanceV2TargetState(d *schema.ResourceData, client *golangsdk.ServiceClient, ref Refresh) error {
	stateConf := &resource.StateChangeConf{
		Pending: ref.Pending,
		Target:  ref.Target,
		Refresh: ApigInstanceV2StateRefreshFunc(client, d.Id()),
		Timeout: ref.Timeout,
		Delay:   ref.Delay,
	}
	if ref.MinTimeout != 0 {
		stateConf.MinTimeout = ref.MinTimeout
	} else {
		stateConf.PollInterval = ref.PollInterval
	}
	_, err := stateConf.WaitForState()
	return err
}

func ApigInstanceV2StateRefreshFunc(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opt := instances.ListOpts{
			Id: id,
		}
		// Some status cannot be read by GET method, just like 'Deleting'.
		// GET method will link to other table (vpc) for query. The response time is not as good as the LIST method.
		allPages, err := instances.List(client, opt).AllPages()
		if err != nil {
			return allPages, "", fmt.Errorf("Error getting APIG v2 dedicated instance by ID (%s): %s", id, err)
		}
		instances, err := instances.ExtractInstances(allPages)
		if len(instances) == 0 {
			return instances, "DELETED", nil
		}
		return instances[0], instances[0].Status, nil
	}
}

func waitForApigInstanceCreateCompleted(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	ref := Refresh{
		Pending:      []string{"Creating", "Initing", "Installing", "Registering"},
		Target:       []string{"Running"},
		Delay:        30 * time.Second,
		Timeout:      d.Timeout(schema.TimeoutCreate),
		PollInterval: 10 * time.Second,
	}
	return watiForApigInstanceV2TargetState(d, client, ref)
}

func resourceApigInstanceV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	opts, err := buildApigInstanceParameters(d, config)
	if err != nil {
		return fmt.Errorf("Error craeting APIG v2 dedicated instance options: %s", err)
	}
	log.Printf("[DEBUG] Create APIG v2 dedicated instance options: %#v", opts)

	client, err := config.ApigV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	v, err := instances.Create(client, opts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 dedicated instance: %s", err)
	}
	d.SetId(v.Id)
	err = waitForApigInstanceCreateCompleted(d, client)
	if err != nil {
		return fmt.Errorf("Error waiting for APIG v2 dedicated instance (%s) to become running: %s", d.Id(), err)
	}
	return resourceApigInstanceV2Read(d, meta)
}

func setApigAvailableZoneIds(d *schema.ResourceData, resp instances.Instance) error {
	idsStr := strings.TrimLeft(resp.AvailableZoneIds, "[")
	idsStr = strings.TrimRight(idsStr, "]")
	idsStr = strings.ReplaceAll(idsStr, " ", "")
	ids := strings.Split(idsStr, ",")
	return d.Set("available_zone_ids", ids)
}

func setApigCreateTimestamp(d *schema.ResourceData, resp instances.Instance) error {
	createTime := time.Unix(resp.CreateTimestamp, 0)
	return d.Set("create_time", createTime.Format(time.RFC3339))
}

func setApigIngressAccess(d *schema.ResourceData, config *config.Config, resp instances.Instance) error {
	if resp.Ipv4IngressEipAddress != "" {
		// The response of ingress acess does not contain eip_id, just the ip address.
		publicAddress := resp.Ipv4IngressEipAddress
		client, err := config.NetworkingV1Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating VPC client: %s", err)
		}
		opt := eips.ListOpts{
			PublicIp: publicAddress,
		}
		allPages, err := eips.List(client, opt).AllPages()
		if err != nil {
			return err
		}
		publicIps, err := eips.ExtractPublicIPs(allPages)
		if err != nil {
			return err
		}
		if len(publicIps) == 0 {
			return fmt.Errorf("Error getting eip id from server by ip address (%s): %s", publicAddress, err)
		}
		return d.Set("eip_id", publicIps[0].ID)
	}
	return d.Set("eip_id", nil)
}

func setApigSupportedFeatures(d *schema.ResourceData, resp instances.Instance) error {
	features := resp.SupportedFeatures
	result := make([]interface{}, len(features))
	for i, v := range features {
		result[i] = v
	}
	return d.Set("supported_features", result)
}

func setApigInstanceParamters(d *schema.ResourceData, config *config.Config, resp instances.Instance) error {
	mErr := multierror.Append(nil,
		d.Set("region", GetRegion(d, config)),
		d.Set("name", resp.Name),
		d.Set("edition", resp.Edition),
		d.Set("vpc_id", resp.VpcId),
		d.Set("subnet_id", resp.SubnetId),
		d.Set("security_group_id", resp.SecurityGroupId),
		d.Set("maintain_time", resp.MaintainBegin),
		d.Set("description", resp.Description),
		d.Set("enterprise_project_id", resp.EnterpriseProjectId),
		d.Set("status", resp.Status),
		d.Set("bandwidth_size", resp.BandwidthSize),
		d.Set("vpc_ingress_address", resp.Ipv4VpcIngressAddress),
		d.Set("egress_address", resp.Ipv4EgressAddress),
		d.Set("ingress_address", resp.Ipv4IngressEipAddress),
		setApigAvailableZoneIds(d, resp),
		setApigCreateTimestamp(d, resp),
		setApigIngressAccess(d, config, resp),
		setApigSupportedFeatures(d, resp),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}
	return nil
}

func getApigInstanceFromServer(d *schema.ResourceData, client *golangsdk.ServiceClient) (*instances.Instance, error) {
	resp, err := instances.Get(client, d.Id()).Extract()
	if err != nil {
		return resp, CheckDeleted(d, err, "APIG v2 dedicated instance")
	}
	log.Printf("[DEBUG] Retrieved APIG v2 dedicated instance (%s): %+v", d.Id(), resp)
	return resp, nil
}

func resourceApigInstanceV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG client: %s", err)
	}
	resp, err := getApigInstanceFromServer(d, client)
	if err != nil {
		return fmt.Errorf("Error getting APIG v2 dedicated instance (%s) form server: %s", d.Id(), err)
	}
	return setApigInstanceParamters(d, config, *resp)
}

func buildApigInstanceUpdateOpts(d *schema.ResourceData) (instances.UpdateOpts, error) {
	opts := instances.UpdateOpts{}
	if d.HasChange("name") {
		opts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		opts.Description = d.Get("description").(string)
	}
	if d.HasChange("maintain_time") {
		startTime := d.Get("maintain_time").(string)
		opts.MaintainBegin = startTime
		endTime, err := buildMaintainEndTime(startTime)
		if err != nil {
			return opts, err
		}
		opts.MaintainEnd = endTime
	}
	if d.HasChange("security_group_id") {
		opts.SecurityGroupId = d.Get("security_group_id").(string)
	}
	log.Printf("[DEBUG] Update options of APIG v2 dedicated instance is: %#v", opts)
	return opts, nil
}

func waitForApigInstanceUpdateCompleted(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	ref := Refresh{
		Pending:    []string{"Updating", "Running"},
		Target:     []string{"Running"},
		Delay:      2 * time.Second,
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		MinTimeout: 2 * time.Second,
	}
	return watiForApigInstanceV2TargetState(d, client, ref)
}

func updateApigInstanceEgressAccess(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	oldVal, newVal := d.GetChange("bandwidth_size")
	// Enable the egress access.
	if oldVal.(int) == 0 {
		size := d.Get("bandwidth_size").(int)
		opts := instances.EgressAccessOpts{
			BandwidthSize: strconv.Itoa(size),
		}
		egress, err := instances.EnableEgressAccess(client, d.Id(), opts).Extract()
		if err != nil {
			return fmt.Errorf("Unable to enable egress access of the dedicated instance (%s), size: %d", d.Id(), size)
		}
		if egress.BandwidthSize != size {
			return fmt.Errorf("Wrong bandwidth size is enabled, size: %d", size)
		}
	}
	// Disable the egress access.
	if newVal.(int) == 0 {
		err := instances.DisableEgressAccess(client, d.Id()).ExtractErr()
		if err != nil {
			return fmt.Errorf("Unable to disable egress bandwidth of the dedicated instance (%s)", d.Id())
		}
		return nil
	}
	// Update the egress nat.
	size := d.Get("bandwidth_size").(int)
	opts := instances.EgressAccessOpts{
		BandwidthSize: strconv.Itoa(size),
	}
	egress, err := instances.UpdateEgressBandwidth(client, d.Id(), opts).Extract()
	if err != nil {
		return fmt.Errorf("Unable to update egress bandwidth of the dedicated instance (%s), size: %d", d.Id(), size)
	}
	if egress.BandwidthSize != size {
		return fmt.Errorf("Wrong bandwidth size is set, size: %d", size)
	}
	return nil
}

func updateApigInstanceIngressAccess(d *schema.ResourceData, client *golangsdk.ServiceClient) (err error) {
	oldVal, newVal := d.GetChange("eip_id")
	// Disable the ingress access.
	// The update logic is to disable first and then enable. Update means thar both oldVal and newVal exist.
	if oldVal.(string) != "" {
		err = instances.DisableIngressAccess(client, d.Id()).ExtractErr()
		if err != nil || newVal.(string) == "" {
			return
		}
	}
	// Enable the ingress access.
	updateOpts := instances.IngressAccessOpts{
		EipId: d.Get("eip_id").(string),
	}
	_, err = instances.EnableIngressAccess(client, d.Id(), updateOpts).Extract()
	return
}

func disableApigInstanceIngressAccess(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	return instances.DisableIngressAccess(client, d.Id()).ExtractErr()
}

func resourceApigInstanceV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}

	// Update egress access
	if d.HasChange("bandwidth_size") {
		if err = updateApigInstanceEgressAccess(d, client); err != nil {
			return fmt.Errorf("Update egress access failed: %s", err)
		}
	}
	// Update ingerss access
	if d.HasChange("eip_id") {
		if err = updateApigInstanceIngressAccess(d, client); err != nil {
			return fmt.Errorf("Update ingress access failed: %s", err)
		}
	}
	// Update APIG v2 instance name, maintain window, description and security group id
	updateOpts, err := buildApigInstanceUpdateOpts(d)
	if err != nil {
		return fmt.Errorf("Unable to get the update options of APIG v2 dedicated instance: %s", err)
	}
	if updateOpts != (instances.UpdateOpts{}) {
		_, err = instances.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating APIG v2 dedicated instance: %s", err)
		}
		err = waitForApigInstanceUpdateCompleted(d, client)
		if err != nil {
			return fmt.Errorf("Error waiting for APIG dedicated instance (%s) to become running: %s", d.Id(), err)
		}
	}
	return resourceApigInstanceV2Read(d, meta)
}

func waitForApigInstanceDeleteCompleted(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	ref := Refresh{
		Pending:      []string{"Deleting"},
		Target:       []string{"DELETED"},
		Delay:        30 * time.Second,
		Timeout:      d.Timeout(schema.TimeoutDelete),
		PollInterval: 10 * time.Second,
	}
	return watiForApigInstanceV2TargetState(d, client, ref)
}

func resourceApigInstanceV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.ApigV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud APIG v2 client: %s", err)
	}
	if err = instances.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmt.Errorf("Unable to delete the APIG v2 dedicated instance (%s): %s", d.Id(), err)
	}
	err = waitForApigInstanceDeleteCompleted(d, client)
	if err != nil {
		return fmt.Errorf("Error deleting APIG v2 dedicated instance (%s): %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}
