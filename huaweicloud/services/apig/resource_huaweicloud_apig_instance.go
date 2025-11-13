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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG POST /v2/{project_id}/apigw/instances
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/postpaid-resize
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/eip
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/eip
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/nat-eip
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/nat-eip
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/nat-eip
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/ingress-eip
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/ingress-eip
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/instance-tags
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/custom-ingress-ports
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/custom-ingress-ports
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/custom-ingress-ports/{ingress_port_id}
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

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region in which to create the dedicated instance resource.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the dedicated instance.`,
			},
			"edition": {
				Type:        schema.TypeString,
				Required:    true,
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the dedicated instance.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The enterprise project ID to which the dedicated instance belongs.`,
			},
			"bandwidth_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The egress bandwidth size of the dedicated instance.`,
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
			"custom_ingress_ports": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Specified the list of the instance custom ingress ports.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specified protocol of the custom ingress port.",
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Specified port of the custom ingress port.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the custom ingress port.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The current status of the custom ingress port.",
						},
					},
				},
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

// buildMaintainEndTime is a method that used to calculate the end time based on the start time, with a 4-hour interval
// between two time strings.
func buildMaintainEndTime(maintainStart string) string {
	result := regexp.MustCompile("^(02|06|10|14|18|22):00:00$").FindStringSubmatch(maintainStart)
	if len(result) < 2 {
		log.Printf("the time format of the maintain window (%s) is incorrect", maintainStart)
		return ""
	}
	num, err := strconv.Atoi(result[1])
	if err != nil {
		log.Printf("the maintain window time (%s) cannot convet to number from a string", result[1])
		return ""
	}
	return fmt.Sprintf("%02d:00:00", (num+4)%24)
}

func buildInstanceAvailabilityZones(d *schema.ResourceData) interface{} {
	if v, ok := d.GetOk("availability_zones"); ok {
		return v.([]interface{})
	}

	// When 'availability_zones' is omitted, the deprecated parameter 'available_zones' is used.
	if v, ok := d.GetOk("available_zones"); ok {
		return v.([]interface{})
	}

	return nil
}

func buildCreateInstanceBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	result := map[string]interface{}{
		"instance_name":                   d.Get("name"),
		"spec_id":                         d.Get("edition"),
		"vpc_id":                          d.Get("vpc_id"),
		"subnet_id":                       d.Get("subnet_id"),
		"security_group_id":               d.Get("security_group_id"),
		"available_zone_ids":              utils.ValueIgnoreEmpty(buildInstanceAvailabilityZones(d)),
		"description":                     utils.ValueIgnoreEmpty(d.Get("description")),
		"bandwidth_size":                  d.Get("bandwidth_size"), // Bandwidth 0 means turn off the egress access.
		"enterprise_project_id":           cfg.GetEnterpriseProjectID(d),
		"eip_id":                          utils.ValueIgnoreEmpty(d.Get("eip_id")),
		"ipv6_enable":                     utils.ValueIgnoreEmpty(d.Get("ipv6_enable")),
		"loadbalancer_provider":           utils.ValueIgnoreEmpty(d.Get("loadbalancer_provider")),
		"tags":                            utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
		"vpcep_service_name":              utils.ValueIgnoreEmpty(d.Get("vpcep_service_name")),
		"ingress_bandwidth_size":          utils.ValueIgnoreEmpty(d.Get("ingress_bandwidth_size")), // BandWidth must be greater than or equal to 5.
		"ingress_bandwidth_charging_mode": utils.ValueIgnoreEmpty(d.Get("ingress_bandwidth_charging_mode")),
		"maintain_begin":                  utils.ValueIgnoreEmpty(d.Get("maintain_begin")),
		"maintain_end":                    utils.ValueIgnoreEmpty(buildMaintainEndTime(d.Get("maintain_begin").(string))),
	}
	log.Printf("[DEBUG] The request body of the create method for the dedicated instance is: %#v", result)
	return result
}

func QueryInstanceDetail(client *golangsdk.ServiceClient, instanceId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", createPath, &opt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func instanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := QueryInstanceDetail(client, instanceId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return "not_found", "COMPLETED", nil
			}
			return respBody, "ERROR", err
		}

		statusResp := utils.PathSearch("status", respBody, "").(string)
		if utils.StrSliceContains([]string{"CreateFail", "InitingFailed", "RegisterFailed", "InstallFailed",
			"UpdateFailed", "RollbackFailed", "UnRegisterFailed", "DeleteFailed", "RestartFail", "ResizeFailed"},
			statusResp) {
			return respBody, "ERROR", fmt.Errorf("unexpect status (%s)", statusResp)
		}

		if utils.StrSliceContains(targets, statusResp) {
			return respBody, "COMPLETED", nil
		}
		return "continue", "PENDING", nil
	}
}

func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v2/{project_id}/apigw/instances"
	)
	client, err := cfg.NewServiceClient("apig", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateInstanceBodyParams(d, cfg)),
	}

	requestResp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating dedicated instance: %s", err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}
	instanceId := utils.PathSearch("instance_id", respBody, "").(string)
	d.SetId(instanceId)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      instanceStateRefreshFunc(client, instanceId, []string{"Running"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        1 * time.Minute,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the status of dedicated instance (%s) to become running: %s", instanceId, err)
	}

	if v, ok := d.GetOk("custom_ingress_ports"); ok {
		if err := addCustomIngressPorts(client, instanceId, v.(*schema.Set).List()); err != nil {
			return diag.FromErr(err)
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

func queryInstanceTags(client *golangsdk.ServiceClient, instanceId string) interface{} {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/instance-tags"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", createPath, &opt)
	if err != nil {
		log.Printf("[WARN] error qeurying tag list of the dedicated instance (%s): %s", instanceId, err)
		return nil
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		log.Printf("[ERROR] error retrieving tag list: %s", err)
		return nil
	}
	return utils.PathSearch("tags", respBody, make([]interface{}, 0))
}

func resourceInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	instanceId := d.Id()
	respBody, err := QueryInstanceDetail(client, instanceId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error qeurying dedicated instance (%s) detail", instanceId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("instance_name", respBody, nil)),
		d.Set("edition", utils.PathSearch("spec", respBody, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", respBody, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", respBody, nil)),
		d.Set("security_group_id", utils.PathSearch("security_group_id", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", respBody, nil)),
		d.Set("bandwidth_size", utils.PathSearch("bandwidth_size", respBody, nil)),
		d.Set("ipv6_enable", utils.PathSearch("!!eip_ipv6_address", respBody, nil)),
		d.Set("loadbalancer_provider", utils.PathSearch("loadbalancer_provider", respBody, nil)),
		d.Set("availability_zones", parseInstanceAvailabilityZones(utils.PathSearch("available_zone_ids", respBody, "").(string))),
		d.Set("maintain_begin", utils.PathSearch("maintain_begin", respBody, nil)),
		d.Set("vpcep_service_name", parseVpcepServiceName(utils.PathSearch("endpoint_services[0].service_name", respBody, "").(string))),
		d.Set("ingress_bandwidth_charging_mode", utils.PathSearch("ingress_bandwidth_charging_mode", respBody, nil)),
		d.Set("ingress_bandwidth_size", utils.PathSearch("publicips[0].bandwidth_size", respBody, nil)),
		// Attributes
		d.Set("maintain_end", utils.PathSearch("maintain_end", respBody, nil)),
		d.Set("vpc_ingress_address", utils.PathSearch("ingress_ip", respBody, nil)),
		d.Set("egress_address", utils.PathSearch("nat_eip_address", respBody, nil)),
		d.Set("supported_features", utils.PathSearch("supported_features", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", respBody, float64(0)).(float64))/1000, false)),
		d.Set("vpcep_service_address", utils.PathSearch("endpoint_services[0].service_name", respBody, nil)),
		d.Set("ingress_address", utils.PathSearch("eip_address||publicips[0].ip_address", respBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(queryInstanceTags(client, instanceId))),
		// Deprecated
		d.Set("create_time", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", respBody, float64(0)).(float64))/1000, false)),
	)
	if resp, err := getCustomIngressPorts(client, instanceId); err != nil {
		// This feature is not available in some region, so use log.Printf to record the error.
		log.Printf("[ERROR] unable to find the custom ingerss ports: %s", err)
	} else {
		mErr = multierror.Append(mErr, d.Set("custom_ingress_ports", flattenCustomIngressPorts(resp)))
	}

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving resource fields of the dedicated instance: %s", mErr)
	}
	return nil
}

func flattenCustomIngressPorts(resp interface{}) []map[string]interface{} {
	customIngressPorts := utils.PathSearch("ingress_port_infos", resp, make([]interface{}, 0)).([]interface{})
	if len(customIngressPorts) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(customIngressPorts))
	for i, v := range customIngressPorts {
		result[i] = map[string]interface{}{
			"protocol": utils.PathSearch("protocol", v, ""),
			"port":     utils.PathSearch("ingress_port", v, 0),
			"id":       utils.PathSearch("ingress_port_id", v, ""),
			"status":   utils.PathSearch("status", v, ""),
		}
	}

	return result
}

func buildUpdateInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	result := map[string]interface{}{
		"instance_name":      d.Get("name"),
		"description":        d.Get("description"),
		"security_group_id":  d.Get("security_group_id"),
		"vpcep_service_name": d.Get("vpcep_service_name"),
		"maintain_begin":     utils.ValueIgnoreEmpty(d.Get("maintain_begin")),
		"maintain_end":       utils.ValueIgnoreEmpty(buildMaintainEndTime(d.Get("maintain_begin").(string))),
	}
	log.Printf("[DEBUG] The request body of the update method for the dedicated instance is: %#v", result)
	return result
}

func updateInstanceBasicConfiguration(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}"
		instanceId = d.Id()
	)

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateInstanceBodyParams(d),
	}

	_, err := client.Request("PUT", createPath, &opt)
	if err != nil {
		return fmt.Errorf("error updating dedicated instance: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      instanceStateRefreshFunc(client, instanceId, []string{"Running"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the status of dedicated instance (%s) to become running: %s", instanceId, err)
	}
	return nil
}

func updatePostPaidInstanceEdition(ctx context.Context, client *golangsdk.ServiceClient, instanceId, newSpec string,
	timeout time.Duration) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/postpaid-resize"
	resizePath := client.Endpoint + httpUrl
	resizePath = strings.ReplaceAll(resizePath, "{project_id}", client.ProjectID)
	resizePath = strings.ReplaceAll(resizePath, "{instance_id}", instanceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"spec_id": newSpec,
		},
	}
	_, err := client.Request("POST", resizePath, &opt)
	if err != nil {
		return fmt.Errorf("error updating the specification of the dedicated instance (%s): %s", instanceId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      instanceStateRefreshFunc(client, instanceId, []string{"Running"}),
		Timeout:      timeout,
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the status of dedicated instance (%s) to become running: %s", instanceId, err)
	}

	return nil
}

func updateInstanceEgressAccess(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		err            error
		oldVal, newVal = d.GetChange("bandwidth_size")
		instanceId     = d.Id()
		httpUrl        = "v2/{project_id}/apigw/instances/{instance_id}/nat-eip"
		generalPath    = client.Endpoint + httpUrl
	)
	generalPath = strings.ReplaceAll(generalPath, "{project_id}", client.ProjectID)
	generalPath = strings.ReplaceAll(generalPath, "{instance_id}", instanceId)

	// Enable the egress access.
	if oldVal.(int) == 0 {
		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"bandwidth_size": newVal,
			},
		}

		_, err = client.Request("POST", generalPath, &opt)
		if err != nil {
			return fmt.Errorf("unable to enable egress bandwidth of the dedicated instance (%s): %s", d.Id(), err)
		}
		// After bandwidth enabled, the update process is completed.
		return nil
	}

	// Disable the egress access.
	if newVal.(int) == 0 {
		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		_, err = client.Request("DELETE", generalPath, &opt)
		if err != nil {
			return fmt.Errorf("unable to disable egress bandwidth of the dedicated instance (%s): %s", d.Id(), err)
		}
		// After bandwidth disabled, the update process is completed.
		return nil
	}

	// Update the egress NAT.
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"bandwidth_size": newVal,
		},
	}
	_, err = client.Request("PUT", generalPath, &opt)
	if err != nil {
		return fmt.Errorf("unable to update egress bandwidth of the dedicated instance (%s): %s", d.Id(), err)
	}
	return nil
}

func updateInstanceIngressAccess(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		err            error
		oldVal, newVal = d.GetChange("eip_id")
		instanceId     = d.Id()
		httpUrl        = "v2/{project_id}/apigw/instances/{instance_id}/eip"
		generalPath    = client.Endpoint + httpUrl
	)
	generalPath = strings.ReplaceAll(generalPath, "{project_id}", client.ProjectID)
	generalPath = strings.ReplaceAll(generalPath, "{instance_id}", instanceId)

	// Disable the ingress access.
	// The update logic is to disable first and then enable. Update means that both oldVal and newVal exist.
	if oldVal.(string) != "" {
		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		_, err = client.Request("DELETE", generalPath, &opt)
		if err != nil {
			return fmt.Errorf("unable to disassociate the ingress EIP: %s", err)
		}
		if newVal.(string) == "" {
			return nil
		}
	}
	// Enable the ingress access.
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"eip_id": newVal,
		},
	}
	_, err = client.Request("PUT", generalPath, &opt)
	if err != nil {
		return fmt.Errorf("unable to associate the ingress EIP: %s", err)
	}
	return nil
}

func updateInstanceTags(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		err            error
		oldVal, newVal = d.GetChange("tags")
		rmTags         = oldVal.(map[string]interface{})
		addTags        = newVal.(map[string]interface{})
		instanceId     = d.Id()
		httpUrl        = "v2/{project_id}/apigw/instances/{instance_id}/instance-tags/action"
		generalPath    = client.Endpoint + httpUrl
	)
	generalPath = strings.ReplaceAll(generalPath, "{project_id}", client.ProjectID)
	generalPath = strings.ReplaceAll(generalPath, "{instance_id}", instanceId)

	if len(rmTags) > 0 {
		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"action": "delete",
				"tags":   utils.ExpandResourceTagsMap(rmTags),
			},
			OkCodes: []int{204},
		}
		_, err = client.Request("POST", generalPath, &opt)
		if err != nil {
			return fmt.Errorf("unable to remove the instance tags: %s", err)
		}
	}
	if len(addTags) > 0 {
		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"action": "create",
				"tags":   utils.ExpandResourceTagsMap(addTags),
			},
			OkCodes: []int{204},
		}
		_, err = client.Request("POST", generalPath, &opt)
		if err != nil {
			return fmt.Errorf("unable to add the instance tags: %s", err)
		}
	}
	return nil
}

func waitForElbIngressAccessCompleted(ctx context.Context, client *golangsdk.ServiceClient, instanceId, action string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshInstanceFunc(client, instanceId, action),
		Timeout:      timeout,
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

func refreshInstanceFunc(client *golangsdk.ServiceClient, instanceId, action string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := QueryInstanceDetail(client, instanceId)
		if err != nil {
			return respBody, "ERROR", err
		}

		if action == "disabled" && utils.PathSearch("length(publicips)", respBody, float64(0)).(float64) < 1 ||
			action == "enabled" && utils.PathSearch("length(publicips)", respBody, float64(0)).(float64) > 0 {
			return "matched", "COMPLETED", nil
		}
		return "continue", "PENDING", nil
	}
}

func updateElbInstanceIngressAccess(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		err                    error
		oldSizeVal, newSizeVal = d.GetChange("ingress_bandwidth_size")
		oldModeVal, newModeVal = d.GetChange("ingress_bandwidth_charging_mode")
		instanceId             = d.Id()
		httpUrl                = "v2/{project_id}/apigw/instances/{instance_id}/ingress-eip"
		generalPath            = client.Endpoint + httpUrl
	)
	generalPath = strings.ReplaceAll(generalPath, "{project_id}", client.ProjectID)
	generalPath = strings.ReplaceAll(generalPath, "{instance_id}", instanceId)

	if oldSizeVal.(int) != 0 || oldModeVal.(string) != "" {
		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		_, err = client.Request("DELETE", generalPath, &opt)
		if err != nil {
			return fmt.Errorf("unable to disable the ingress EIP: %s", err)
		}

		err = waitForElbIngressAccessCompleted(ctx, client, instanceId, "disabled", d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("error waiting for ingress EIP unbinding completed: %s", err)
		}
	}

	if newSizeVal.(int) == 0 && newModeVal.(string) == "" {
		return nil
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"bandwidth_size":          newSizeVal,
			"bandwidth_charging_mode": newModeVal,
		},
	}
	_, err = client.Request("POST", generalPath, &opt)
	if err != nil {
		return fmt.Errorf("unable to enable the ingress EIP: %s", err)
	}

	err = waitForElbIngressAccessCompleted(ctx, client, instanceId, "enabled", d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("error waiting for ingress EIP enabling completed: %s", err)
	}
	return nil
}

func resourceInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Id()
	)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	if d.HasChanges("instance_name", "description", "security_group_id", "vpcep_service_name", "maintain_begin", "maintain_end") {
		if err = updateInstanceBasicConfiguration(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	// Update specification
	if d.HasChange("edition") {
		if err = updatePostPaidInstanceEdition(ctx, client, instanceId, d.Get("edition").(string),
			d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.FromErr(err)
		}
	}

	// Update egress access
	if d.HasChange("bandwidth_size") {
		if err = updateInstanceEgressAccess(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	// Update ingerss access
	if d.HasChanges("ingress_bandwidth_size", "ingress_bandwidth_charging_mode") {
		if err = updateElbInstanceIngressAccess(ctx, d, client); err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("eip_id") {
		if err = updateInstanceIngressAccess(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		if err = updateInstanceTags(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   instanceId,
			ResourceType: "apig",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("custom_ingress_ports") {
		oldRaws, newRaws := d.GetChange("custom_ingress_ports")
		err = updateCustomIngressPorts(client, oldRaws, newRaws, instanceId)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceInstanceRead(ctx, d, meta)
}

func addCustomIngressPorts(client *golangsdk.ServiceClient, instanceId string, ingressPorts []interface{}) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/custom-ingress-ports"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	for _, ingressPort := range ingressPorts {
		opts := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"protocol":     utils.PathSearch("protocol", ingressPort, nil),
				"ingress_port": utils.PathSearch("port", ingressPort, 0),
			},
		}
		_, err := client.Request("POST", createPath, &opts)
		if err != nil {
			return fmt.Errorf("error adding custom ingerss port to the dedicated instance: %s", err)
		}
	}
	return nil
}

func getCustomIngressPorts(client *golangsdk.ServiceClient, instanceId string) (interface{}, error) {
	// Currently, a maximum of 40 custom ingress ports can be created. Limit default value is 20.
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/custom-ingress-ports?limit=500"
	httpPath := client.Endpoint + httpUrl
	httpPath = strings.ReplaceAll(httpPath, "{project_id}", client.ProjectID)
	httpPath = strings.ReplaceAll(httpPath, "{instance_id}", instanceId)
	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", httpPath, &opts)
	if err != nil {
		return nil, fmt.Errorf("error getting custom ingerss port form the instance (%s) : %s", instanceId, err)
	}

	return utils.FlattenResponse(resp)
}

func removeCustomIngressPorts(client *golangsdk.ServiceClient, instanceId string, customIngressPorts []interface{}) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/custom-ingress-ports/{ingress_port_id}"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{instance_id}", instanceId)

	for _, v := range customIngressPorts {
		deletePath := strings.ReplaceAll(path, "{ingress_port_id}", utils.PathSearch("id", v, "").(string))
		opts := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		_, err := client.Request("DELETE", deletePath, &opts)
		if err != nil {
			return fmt.Errorf("error removing custom ingerss port form the instance (%s) : %s", instanceId, err)
		}
	}

	return nil
}

func updateCustomIngressPorts(client *golangsdk.ServiceClient, oldRaws, newRaws interface{}, instanceId string) error {
	var (
		addRaws    = newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
		removeRaws = oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))
	)
	if removeRaws.Len() > 0 {
		if err := removeCustomIngressPorts(client, instanceId, removeRaws.List()); err != nil {
			return err
		}
	}

	if addRaws.Len() > 0 {
		if err := addCustomIngressPorts(client, instanceId, addRaws.List()); err != nil {
			return err
		}
	}
	return nil
}

func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}"
		instanceId = d.Id()
	)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &opt)
	if err != nil {
		return diag.Errorf("error deleting dedicated instance (%s): %s", instanceId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      instanceStateRefreshFunc(client, d.Id(), nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        20 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
