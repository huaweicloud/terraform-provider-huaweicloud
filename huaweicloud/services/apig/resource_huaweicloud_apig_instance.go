package apig

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/instances"
	"github.com/chnsz/golangsdk/openstack/eps/v1/enterpriseprojects"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type Edition string      // The edition of the dedicated instance.
type ProviderType string // The type of the loadbalancer provider.

const (
	// IPv4 Editions
	EditionBasic        Edition = "BASIC"        // Basic Edition instance.
	EditionProfessional Edition = "PROFESSIONAL" // Professional Edition instance.
	EditionEnterprise   Edition = "ENTERPRISE"   // Enterprise Edition instance.
	EditionPlatinum     Edition = "PLATINUM"     // Platinum Edition instance.
	// IPv6 Editions
	Ipv6EditionBasic        Edition = "BASIC_IPv6"        // IPv6 instance of the Basic Edition.
	Ipv6EditionProfessional Edition = "PROFESSIONAL_IPv6" // IPv6 instance of the Professional Edition.
	Ipv6EditionEnterprise   Edition = "ENTERPRISE_IPv6"   // IPv6 instance of the Enterprise Edition.
	Ipv6EditionPlatinum     Edition = "PLATINUM_IPv6"     // IPv6 instance of the Platinum Edition.
)

// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/eip
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/eip
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/nat-eip
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/nat-eip
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/nat-eip
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/instance-tags/action
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/instance-tags
// @API APIG POST /v2/{project_id}/apigw/instances
// @API EIP GET /v1/{project_id}/publicips
// @API APIG POST /v2/{project_id}/apigw/instances{instance_id}/ingress-eip
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/ingress-eip
func ResourceApigInstanceV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceCreate,
		ReadContext:   resourceInstanceRead,
		UpdateContext: resourceInstanceUpdate,
		DeleteContext: resourceInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region in which to create the dedicated instance resource.`,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile("^([\u4e00-\u9fa5A-Za-z][\u4e00-\u9fa5A-Za-z-_0-9]*)$"),
						"The name can only contain letters, digits, hyphens (-) and underscore (_), and must start "+
							"with a letter."),
					validation.StringLenBetween(3, 64),
				),
				Description: `The name of the dedicated instance.`,
			},
			"edition": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(EditionBasic),
					string(EditionProfessional),
					string(EditionEnterprise),
					string(EditionPlatinum),
					string(Ipv6EditionBasic),
					string(Ipv6EditionProfessional),
					string(Ipv6EditionEnterprise),
					string(Ipv6EditionPlatinum),
				}, false),
				Description: `The edition of the dedicated instance.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the VPC used to create the dedicated instance.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the VPC subnet used to create the dedicated instance.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the security group to which the dedicated instance belongs to.`,
			},
			"availability_zones": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `schema: Required; The name list of availability zones for the dedicated instance.`,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile("^[^<>]*$"),
						"The description cannot contain the angle brackets (< and >)."),
					validation.StringLenBetween(0, 255),
				),
				Description: `The description of the dedicated instance.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The enterprise project ID to which the dedicated instance belongs.`,
			},
			"bandwidth_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 2000),
				Description:  `The egress bandwidth size of the dedicated instance.`,
			},
			"ipv6_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Whether public access with an IPv6 address is supported.`,
			},
			"maintain_begin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^(02|06|10|14|18|22):00:00$`),
					"The start-time format of maintenance window is not 'xx:00:00' or "+
						"the hour is not 02, 06, 10, 14, 18 or 22."),
				Description: `The start time of the maintenance time window.`,
			},
			"vpcep_service_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Name of the VPC endpoint service.`,
			},
			"tags": common.TagsSchema(),
			"ingress_bandwidth_size": {
				Type:          schema.TypeInt,
				Optional:      true,
				RequiredWith:  []string{"ingress_bandwidth_charging_mode"},
				ConflictsWith: []string{"eip_id"},
			},
			"ingress_bandwidth_charging_mode": {
				Type:          schema.TypeString,
				Optional:      true,
				RequiredWith:  []string{"ingress_bandwidth_size"},
				ConflictsWith: []string{"eip_id"},
			},
			// Attributes
			"maintain_end": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `End time of the maintenance time window, 4-hour difference between the start time and end time.`,
			},
			"ingress_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ingress EIP address.`,
			},
			"vpc_ingress_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ingress private IP address of the VPC.`,
			},
			"egress_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The egress (NAT) public IP address.`,
			},
			"supported_features": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The supported features of the dedicated instance.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Time when the dedicated instance is created, in RFC-3339 format.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Status of the dedicated instance.`,
			},
			"vpcep_service_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The address (full name) of the VPC endpoint service.`,
			},
			// Deprecated arguments
			"available_zones": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `schema: Deprecated; The name list of availability zones for the dedicated instance.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Deprecated:  "Use 'created_at' instead",
				Description: `schema: Deprecated; Time when the dedicated instance is created.`,
			},
			"eip_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ConflictsWith: []string{
					"ingress_bandwidth_size", "ingress_bandwidth_charging_mode",
				},
				Description: utils.SchemaDesc(
					`The EIP ID associated with the dedicated instance.`,
					utils.SchemaDescInput{
						Deprecated: true,
					}),
			},
			"loadbalancer_provider": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Description: utils.SchemaDesc(
					`The type of loadbalancer provider used by the instance.`,
					utils.SchemaDescInput{
						Computed: true,
					}),
			},
		},
	}
}

func buildMaintainEndTime(maintainStart string) (string, error) {
	result := regexp.MustCompile("^(02|06|10|14|18|22):00:00$").FindStringSubmatch(maintainStart)
	if len(result) < 2 {
		return "", fmt.Errorf("the hour is missing")
	}
	num, err := strconv.Atoi(result[1])
	if err != nil {
		return "", fmt.Errorf("the number (%s) cannot be converted to string", result[1])
	}
	return fmt.Sprintf("%02d:00:00", (num+4)%24), nil
}

func buildInstanceAvailabilityZones(d *schema.ResourceData) ([]string, error) {
	if v, ok := d.GetOk("availability_zones"); ok {
		return utils.ExpandToStringList(v.([]interface{})), nil
	}

	// When 'availability_zones' is omitted, the deprecated parameter 'available_zones' is used.
	if v, ok := d.GetOk("available_zones"); ok {
		return utils.ExpandToStringList(v.([]interface{})), nil
	}

	return nil, fmt.Errorf("The parameter 'availability_zones' must be specified")
}

func buildInstanceCreateOpts(d *schema.ResourceData, cfg *config.Config) (instances.CreateOpts, error) {
	result := instances.CreateOpts{
		Name:                        d.Get("name").(string),
		Edition:                     d.Get("edition").(string),
		VpcId:                       d.Get("vpc_id").(string),
		SubnetId:                    d.Get("subnet_id").(string),
		SecurityGroupId:             d.Get("security_group_id").(string),
		Description:                 d.Get("description").(string),
		EipId:                       d.Get("eip_id").(string),
		BandwidthSize:               d.Get("bandwidth_size").(int), // Bandwidth 0 means turn off the egress access.
		EnterpriseProjectId:         common.GetEnterpriseProjectID(d, cfg),
		Ipv6Enable:                  d.Get("ipv6_enable").(bool),
		LoadbalancerProvider:        d.Get("loadbalancer_provider").(string),
		Tags:                        utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		VpcepServiceName:            d.Get("vpcep_service_name").(string),
		IngressBandwithSize:         d.Get("ingress_bandwidth_size").(int), // BandWidth must be greater than or equal to 5.
		IngressBandwithChargingMode: d.Get("ingress_bandwidth_charging_mode").(string),
	}

	azList, err := buildInstanceAvailabilityZones(d)
	if err != nil {
		return result, err
	}
	result.AvailableZoneIds = azList

	if v, ok := d.GetOk("maintain_begin"); ok {
		startTime := v.(string)
		result.MaintainBegin = startTime
		endTime, err := buildMaintainEndTime(startTime)
		if err != nil {
			return result, err
		}
		result.MaintainEnd = endTime
	}

	log.Printf("[DEBUG] Create options of the dedicated instance is: %#v", result)
	return result, nil
}

func buildTagsUpdateOpts(tags map[string]interface{}, instanceId, action string) *instances.TagsUpdateOpts {
	if len(tags) < 1 {
		return nil
	}
	return &instances.TagsUpdateOpts{
		InstanceId: instanceId,
		Action:     action,
		Tags:       utils.ExpandResourceTags(tags),
	}
}

func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	opts, err := buildInstanceCreateOpts(d, cfg)
	if err != nil {
		return diag.Errorf("error creating the dedicated instance options: %s", err)
	}
	log.Printf("[DEBUG] The CreateOpts of the dedicated instance is: %#v", opts)

	resp, err := instances.Create(client, opts).Extract()
	if err != nil {
		return diag.Errorf("error creating the dedicated instance: %s", err)
	}
	d.SetId(resp.Id)

	instanceId := d.Id()
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      InstanceStateRefreshFunc(client, instanceId, []string{"Running"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the dedicated instance (%s) to become running: %s", instanceId, err)
	}

	if tagsRaw, ok := d.GetOk("tags"); ok {
		err = instances.UpdateTags(client, buildTagsUpdateOpts(tagsRaw.(map[string]interface{}), instanceId, "create"))
		if err != nil {
			return diag.Errorf("error creating instance tags: %s", err)
		}
	}

	return resourceInstanceRead(ctx, d, meta)
}

// parseInstanceAvailabilityZones is a method that used to convert the string returned by the API which contains
// brackets ([ and ]) and space into a list of strings (available_zone code) and save to state.
func parseInstanceAvailabilityZones(azStr string) []string {
	codesStr := strings.TrimLeft(azStr, "[")
	codesStr = strings.TrimRight(codesStr, "]")
	codesStr = strings.ReplaceAll(codesStr, " ", "")

	return strings.Split(codesStr, ",")
}

// The response of ingress access does not contain EIP ID, just the IP address.
func parseInstanceIngressAccess(cfg *config.Config, region, publicAddress string) (*string, error) {
	if publicAddress == "" {
		return nil, nil
	}

	client, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC v1 client: %s", err)
	}

	opt := eips.ListOpts{
		PublicIp:            []string{publicAddress},
		EnterpriseProjectId: "all_granted_eps",
	}
	allPages, err := eips.List(client, opt).AllPages()
	if err != nil {
		return nil, err
	}
	publicIps, err := eips.ExtractPublicIPs(allPages)
	if err != nil {
		return nil, err
	}
	if len(publicIps) > 0 {
		return &publicIps[0].ID, nil
	}

	log.Printf("[WARN] The instance does not synchronize EIP information, got (%s), but not found on the server",
		publicAddress)
	return nil, nil
}

func parseInstanceIpv6Enable(ipv6Address string) bool {
	return ipv6Address != ""
}

func parseVpcepServiceName(serviceName string) string {
	// The format of the service endpoint is the '{region}.{vpcep_service_name}.{service_id}'
	regexExp := `^[\w-]+\.(.*)\.[a-f0-9-]+$`
	result := regexp.MustCompile(regexExp).FindStringSubmatch(serviceName)
	log.Printf("[DEBUG] The result of the regex matching is: %v (length: %d)", result, len(result))
	if len(result) <= 1 {
		return ""
	}
	// For the result of the regex matching, the first element (result[0]) is the full
	// address ({region}.{vpcep_service_name}.{service_id}), the others (result[1:]) are match objects.
	return result[1]
}

func resourceInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	instanceId := d.Id()
	resp, err := instances.Get(client, instanceId).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error getting instance (%s) details form server", instanceId))
	}
	log.Printf("[DEBUG] Retrieved the dedicated instance (%s): %#v", instanceId, resp)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", resp.Name),
		d.Set("edition", resp.Edition),
		d.Set("vpc_id", resp.VpcId),
		d.Set("subnet_id", resp.SubnetId),
		d.Set("security_group_id", resp.SecurityGroupId),
		d.Set("description", resp.Description),
		d.Set("enterprise_project_id", resp.EnterpriseProjectId),
		d.Set("bandwidth_size", resp.BandwidthSize),
		d.Set("ipv6_enable", parseInstanceIpv6Enable(resp.Ipv6IngressEipAddress)),
		d.Set("loadbalancer_provider", resp.LoadbalancerProvider),
		d.Set("availability_zones", parseInstanceAvailabilityZones(resp.AvailableZoneIds)),
		d.Set("maintain_begin", resp.MaintainBegin),
		d.Set("ingress_bandwidth_charging_mode", resp.IngressBandwidthChargingMode),
		// Attributes
		d.Set("maintain_end", resp.MaintainEnd),
		d.Set("ingress_address", resp.Ipv4IngressEipAddress),
		d.Set("vpc_ingress_address", resp.Ipv4VpcIngressAddress),
		d.Set("egress_address", resp.Ipv4EgressAddress),
		d.Set("supported_features", resp.SupportedFeatures),
		d.Set("status", resp.Status),
		d.Set("created_at", utils.FormatTimeStampRFC3339(resp.CreateTimestamp, false)),
		// Deprecated
		d.Set("create_time", utils.FormatTimeStampRFC3339(resp.CreateTimestamp, false)),
	)

	if eipId, err := parseInstanceIngressAccess(cfg, region, resp.Ipv4IngressEipAddress); err != nil {
		mErr = multierror.Append(mErr, err)
	} else {
		mErr = multierror.Append(d.Set("eip_id", eipId))
	}

	if len(resp.EndpointServices) > 0 {
		mErr = multierror.Append(mErr,
			d.Set("vpcep_service_name", parseVpcepServiceName(resp.EndpointServices[0].ServiceName)),
			d.Set("vpcep_service_address", resp.EndpointServices[0].ServiceName),
		)
	}

	ingressBandwidthSize := 0
	if len(resp.PublicIps) > 0 {
		ingressBandwidthSize = resp.PublicIps[0].BandwidthSize
	}
	mErr = multierror.Append(mErr,
		d.Set("ingress_bandwidth_size", ingressBandwidthSize),
	)

	if tagList, err := instances.GetTags(client, instanceId); err != nil {
		log.Printf("[WARN] error querying instance tags: %s", err)
	} else {
		mErr = multierror.Append(d.Set("tags", utils.TagsToMap(tagList)))
	}

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving resource fields of the dedicated instance: %s", mErr)
	}

	return nil
}

func buildInstanceUpdateOpts(d *schema.ResourceData) (instances.UpdateOpts, error) {
	result := instances.UpdateOpts{}
	if d.HasChange("name") {
		result.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		result.Description = utils.String(d.Get("description").(string))
	}
	if d.HasChange("security_group_id") {
		result.SecurityGroupId = d.Get("security_group_id").(string)
	}
	if d.HasChange("vpcep_service_name") {
		result.VpcepServiceName = d.Get("vpcep_service_name").(string)
	}
	if d.HasChange("maintain_begin") {
		startTime := d.Get("maintain_begin").(string)
		result.MaintainBegin = startTime
		endTime, err := buildMaintainEndTime(startTime)
		if err != nil {
			return result, err
		}
		result.MaintainEnd = endTime
	}

	log.Printf("[DEBUG] Update options of the dedicated instance is: %#v", result)
	return result, nil
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
			return fmt.Errorf("unable to enable egress bandwidth of the dedicated instance (%s): %s", d.Id(), err)
		}
		if egress.BandwidthSize != size {
			return fmt.Errorf("the egress bandwidth size change failed, want '%d', but '%d'", size, egress.BandwidthSize)
		}
	}
	// Disable the egress access.
	if newVal.(int) == 0 {
		err := instances.DisableEgressAccess(client, d.Id()).ExtractErr()
		if err != nil {
			return fmt.Errorf("unable to disable egress bandwidth of the dedicated instance (%s)", d.Id())
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
		return fmt.Errorf("unable to update egress bandwidth of the dedicated instance (%s): %s", d.Id(), err)
	}
	if egress.BandwidthSize != size {
		return fmt.Errorf("the egress bandwidth size change failed, want '%d', but '%d'", size, egress.BandwidthSize)
	}
	return nil
}

func updateInstanceIngressAccess(d *schema.ResourceData, client *golangsdk.ServiceClient) (err error) {
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

func updateInstanceTags(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		err              error
		instanceId       = d.Id()
		oldRaws, newRaws = d.GetChange("tags")
		rmTags           = oldRaws.(map[string]interface{})
		addTags          = newRaws.(map[string]interface{})
	)
	if len(rmTags) > 0 {
		err := instances.UpdateTags(client, buildTagsUpdateOpts(rmTags, instanceId, "delete"))
		if err != nil {
			return fmt.Errorf("error deleting instance tags: %s", err)
		}
	}
	if len(addTags) > 0 {
		err = instances.UpdateTags(client, buildTagsUpdateOpts(addTags, instanceId, "create"))
		if err != nil {
			return fmt.Errorf("[WARN] error creating instance tags: %s", err)
		}
	}
	return nil
}

func waitForUpdateIngressEIPCompleted(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, action string) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshInstanceFunc(client, d, action),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
		// When changing the bandwidth billing type, there will be a delay between the EIP unbinding and EIP binding.
		ContinuousTargetOccurence: 2,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func refreshInstanceFunc(client *golangsdk.ServiceClient, d *schema.ResourceData, action string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := instances.Get(client, d.Id()).Extract()
		if err != nil {
			return resp, "", err
		}

		disabledSucc := action == "disabled" && len(resp.PublicIps) == 0
		enabledSucc := action == "enabled" && len(resp.PublicIps) > 0
		if enabledSucc || disabledSucc {
			return resp, "COMPLETED", nil
		}

		return resp, "PENDING", nil
	}
}

func updateElbInstanceIngressAccess(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	oldSizeVal, newSizeVal := d.GetChange("ingress_bandwidth_size")
	oldModeVal, newModeVal := d.GetChange("ingress_bandwidth_charging_mode")
	instanceId := d.Id()
	if oldSizeVal.(int) != 0 || oldModeVal.(string) != "" {
		err := instances.DisableElbIngressAccess(client, instanceId)
		if err != nil {
			return fmt.Errorf("error unbinding ingress EIP of the dedicated instance: %s", err)
		}

		err = waitForUpdateIngressEIPCompleted(ctx, d, client, "disabled")
		if err != nil {
			return fmt.Errorf("error waiting for unbinding ingress EIP completed: %s", err)
		}
	}

	if newSizeVal.(int) == 0 && newModeVal.(string) == "" {
		return nil
	}

	opts := instances.ElbIngressAccessOpts{
		InstanceId:                  instanceId,
		IngressBandwithSize:         newSizeVal.(int),
		IngressBandwithChargingMode: newModeVal.(string),
	}
	_, err := instances.EnableElbIngressAccess(client, opts)
	if err != nil {
		return fmt.Errorf("error enabled ingress bandwidth of the dedicated instance: %s", err)
	}

	err = waitForUpdateIngressEIPCompleted(ctx, d, client, "enabled")
	if err != nil {
		return fmt.Errorf("error waiting for enabling ingress EIP completed: %s", err)
	}
	return nil
}

func resourceInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	instanceId := d.Id()

	// Update egress access
	if d.HasChange("bandwidth_size") {
		if err = updateApigInstanceEgressAccess(d, client); err != nil {
			return diag.Errorf("update egress access failed: %s", err)
		}
	}
	// Update ingerss access
	if d.HasChange("eip_id") {
		if err = updateInstanceIngressAccess(d, client); err != nil {
			return diag.Errorf("update ingress access failed: %s", err)
		}
	}
	// Update instance name, maintain window, description, security group ID and vpcep service name.
	updateOpts, err := buildInstanceUpdateOpts(d)
	if err != nil {
		return diag.Errorf("unable to get the update options of the dedicated instance: %s", err)
	}
	if updateOpts != (instances.UpdateOpts{}) {
		_, err = instances.Update(client, instanceId, updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating the dedicated instance: %s", err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:      []string{"PENDING"},
			Target:       []string{"COMPLETED"},
			Refresh:      InstanceStateRefreshFunc(client, instanceId, []string{"Running"}),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        20 * time.Second,
			PollInterval: 20 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("tags") {
		if err = updateInstanceTags(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := enterpriseprojects.MigrateResourceOpts{
			ResourceId:   instanceId,
			ResourceType: "apig",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := common.MigrateEnterpriseProject(ctx, cfg, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("ingress_bandwidth_size", "ingress_bandwidth_charging_mode") {
		if err = updateElbInstanceIngressAccess(ctx, d, client); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceInstanceRead(ctx, d, meta)
}

func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}
	if err = instances.Delete(client, d.Id()).ExtractErr(); err != nil {
		return diag.Errorf("error deleting the dedicated instance (%s): %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      InstanceStateRefreshFunc(client, d.Id(), nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func InstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := instances.Get(client, instanceId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return resp, "COMPLETED", nil
			}
			return resp, "", err
		}

		if utils.StrSliceContains([]string{"CreateFail", "InitingFailed", "RegisterFailed", "InstallFailed",
			"UpdateFailed", "RollbackFailed", "UnRegisterFailed", "DeleteFailed"}, resp.Status) {
			return resp, "", fmt.Errorf("unexpect status (%s)", resp.Status)
		}

		if utils.StrSliceContains(targets, resp.Status) {
			return resp, "COMPLETED", nil
		}
		return resp, "PENDING", nil
	}
}
