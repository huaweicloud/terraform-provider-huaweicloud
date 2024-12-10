package dns

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DNS POST /v2.1/endpoints
// @API DNS GET /v2.1/endpoints/{endpoint_id}
// @API DNS GET /v2.1/endpoints/{endpoint_id}/ipaddresses
// @API VPC GET /v1/{project_id}/subnets
// @API DNS PUT /v2.1/endpoints/{endpoint_id}
// @API DNS POST /v2.1/endpoints/{endpoint_id}/ipaddresses
// @API DNS DELETE /v2.1/endpoints/{endpoint_id}/ipaddresses/{ipaddress_id}
// @API DNS DELETE /v2.1/endpoints/{endpoint_id}
func ResourceEndpointAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEndpointAssignmentCreate,
		ReadContext:   resourceEndpointAssignmentRead,
		UpdateContext: resourceEndpointAssignmentUpdate,
		DeleteContext: resourceEndpointAssignmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
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
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the endpoint.`,
			},
			"direction": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The direction of the endpoint.`,
			},
			"assignments": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 6,
				MinItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The subnet ID to which the IP address belongs.`,
						},
						"ip_address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The IP address associated with the endpoint.`,
						},
						"ip_address_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the IP address associated with the endpoint.`,
						},
					},
				},
				Description: `The list of the IP addresses of the endpoint.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The VPC ID associated with the endpoint.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current status of the endpoint.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the endpoint, in RFC3339 format.`,
			},
		},
	}
}

func buildCreateEndpointOpts(d *schema.ResourceData, region string) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"direction":   d.Get("direction"),
		"region":      region,
		"ipaddresses": buildEndpointAssignments(d.Get("assignments").(*schema.Set).List()),
	}
}

func buildEndpointAssignments(assignments []interface{}) []interface{} {
	rst := make([]interface{}, len(assignments))
	for i, v := range assignments {
		rst[i] = map[string]interface{}{
			"subnet_id": utils.PathSearch("subnet_id", v, nil),
			"ip":        utils.PathSearch("ip_address", v, nil),
		}
	}
	return rst
}

func resourceEndpointAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf    = meta.(*config.Config)
		region  = conf.GetRegion(d)
		httpUrl = "v2.1/endpoints"
	)

	client, err := conf.NewServiceClient("dns", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateEndpointOpts(d, region)),
	}
	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DNS endpoint: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	endpointId := utils.PathSearch("endpoint.id", respBody, "").(string)
	if endpointId == "" {
		return diag.Errorf("unable to find endpoint ID from API response")
	}
	d.SetId(endpointId)

	stateConf := &resource.StateChangeConf{
		Target:       []string{"ACTIVE"},
		Pending:      []string{"PENDING"},
		Refresh:      refreshEndpointStatus(client, endpointId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for DNS endpoint (%s) creation to completed: %s", endpointId, err)
	}

	return resourceEndpointAssignmentRead(ctx, d, meta)
}

func refreshEndpointStatus(client *golangsdk.ServiceClient, endpointId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		endPoint, err := GetEntpointById(client, endpointId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "Resource Not Found", "DELETED", nil
			}
			return nil, "ERROR", err
		}

		return endPoint, parseStatus(utils.PathSearch("status", endPoint, "").(string)), nil
	}
}

// GetEntpointById is a method used to get endpoint detail by specified endpoint ID.
func GetEntpointById(client *golangsdk.ServiceClient, endpointId string) (interface{}, error) {
	httpUrl := "v2.1/endpoints/{endpoint_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{endpoint_id}", endpointId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	return utils.PathSearch("endpoint", respBody, nil), err
}

func resourceEndpointAssignmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf       = meta.(*config.Config)
		region     = conf.GetRegion(d)
		endpointId = d.Id()
	)
	client, err := conf.NewServiceClient("dns", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	endpointInfo, err := GetEntpointById(client, endpointId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DNS endpoint")
	}

	reqResp, err := getEndpointIpAdressesById(client, endpointId)
	if err != nil {
		return diag.Errorf("error retrieving IP addresses under specified endpoint (%s): %s", endpointId, err)
	}

	vpcId := utils.PathSearch("vpc_id", endpointInfo, "").(string)
	subnetClient, err := conf.NewServiceClient("vpc", region)
	if err != nil {
		return diag.Errorf("error creating VPC client: %s", err)
	}

	subnets, err := getSubnetsByVpcId(subnetClient, vpcId)
	if err != nil {
		return diag.Errorf("error retrieving subnets under specified VPC (%s): %s", vpcId, err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", endpointInfo, nil)),
		d.Set("direction", utils.PathSearch("direction", endpointInfo, nil)),
		d.Set("assignments",
			flattenAssignments(utils.PathSearch("ipaddresses", reqResp, make([]interface{}, 0)).([]interface{}), subnets)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", endpointInfo, nil)),
		d.Set("status", utils.PathSearch("status", endpointInfo, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
			utils.PathSearch("create_time", endpointInfo, "").(string), "2006-01-02T15:04:05.000")/1000, false)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr.ErrorOrNil())
	}

	return nil
}

func getEndpointIpAdressesById(client *golangsdk.ServiceClient, endpointId string) (interface{}, error) {
	httpUrl := "v2.1/endpoints/{endpoint_id}/ipaddresses"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{endpoint_id}", endpointId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func getSubnetsByVpcId(client *golangsdk.ServiceClient, vpcId string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/subnets?limit=100&vpcId={vpc_id}"
		listOpt = golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		marker = ""
		result = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{vpc_id}", vpcId)

	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s&marker=%v", listPath, marker)
		}

		requestResp, err := client.Request("GET", listPathWithMarker, &listOpt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		subnets := utils.PathSearch("subnets", respBody, make([]interface{}, 0)).([]interface{})
		if len(subnets) < 1 {
			break
		}
		result = append(result, subnets...)
		marker = utils.PathSearch("[-1].id", subnets, "").(string)
	}

	return result, nil
}

func flattenAssignments(ipAddresses, subnets []interface{}) []map[string]interface{} {
	rst := make([]map[string]interface{}, len(ipAddresses))
	for i, v := range ipAddresses {
		rst[i] = map[string]interface{}{
			"ip_address":    utils.PathSearch("ip", v, nil),
			"ip_address_id": utils.PathSearch("id", v, nil),
		}

		for _, subnet := range subnets {
			if utils.PathSearch("subnet_id", v, "").(string) == utils.PathSearch("neutron_subnet_id", subnet, "").(string) {
				rst[i]["subnet_id"] = utils.PathSearch("id", subnet, nil)
				break
			}
		}
	}

	return rst
}

func resourceEndpointAssignmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	endpointId := d.Id()
	client, err := conf.NewServiceClient("dns", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	if d.HasChange("name") {
		if err = updateEndpointName(client, endpointId, d.Get("name").(string)); err != nil {
			return diag.Errorf("error updating the name of the endpoint (%s): %s", endpointId, err)
		}
	}

	if d.HasChange("assignments") {
		err = updateEndpointAssignments(client, d, endpointId)
		if err != nil {
			return diag.FromErr(err)
		}

		stateConf := &resource.StateChangeConf{
			Target:       []string{"ACTIVE"},
			Pending:      []string{"PENDING"},
			Refresh:      refreshEndpointStatus(client, d.Id()),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        5 * time.Second,
			PollInterval: 5 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for DNS endpoint (%s) update to completed: %s", d.Id(), err)
		}
	}

	return resourceEndpointAssignmentRead(ctx, d, meta)
}

func updateEndpointName(client *golangsdk.ServiceClient, endpointId, name string) error {
	httpUrl := "v2.1/endpoints/{endpoint_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{endpoint_id}", endpointId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"name": name,
		},
	}

	_, err := client.Request("PUT", updatePath, &getOpt)
	return err
}

func updateEndpointAssignments(client *golangsdk.ServiceClient, d *schema.ResourceData, endpointId string) error {
	var (
		oldRaws, newRaws = d.GetChange("assignments")
		addRaws          = newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set)).List()
		removeRaws       = oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set)).List()
		// The `originNum` indicates the actual number of IP addresses obtained in the interface.
		originNum = oldRaws.(*schema.Set).Len()
		err       error
	)

	// Since the assignment length ranges from `2` to `6`, the logic is to handle the critical value.
	for len(addRaws) > 0 || len(removeRaws) > 0 {
		// If the number of IP address mappings does not reach the upper limit, add an IP address mapping first.
		if originNum < 6 && len(addRaws) > 0 {
			addRaws, err = addEndpointIpAddress(client, addRaws, endpointId)
			if err != nil {
				return err
			}
			originNum++
			continue
		}

		if removeRaws, err = removeEndpointIpAddress(client, removeRaws, endpointId); err != nil {
			return err
		}
		originNum--
	}
	return nil
}

func addEndpointIpAddress(client *golangsdk.ServiceClient, addIpAddresses []interface{}, endpointId string) ([]interface{}, error) {
	httpUrl := "v2.1/endpoints/{endpoint_id}/ipaddresses"
	addPath := client.Endpoint + httpUrl
	addPath = strings.ReplaceAll(addPath, "{endpoint_id}", endpointId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"ipaddress": map[string]interface{}{
				"subnet_id": utils.PathSearch("subnet_id", addIpAddresses[0], nil),
				"ip":        utils.PathSearch("ip_address", addIpAddresses[0], nil),
			},
		},
	}

	_, err := client.Request("POST", addPath, &opt)
	if err != nil {
		return nil, fmt.Errorf("error adding IP addresss for endpoint (%s): %s", endpointId, err)
	}
	return addIpAddresses[1:], nil
}

func removeEndpointIpAddress(client *golangsdk.ServiceClient, rmIpAddresses []interface{}, endpointId string) ([]interface{}, error) {
	httpUrl := "v2.1/endpoints/{endpoint_id}/ipaddresses/{ipaddress_id}"
	addPath := client.Endpoint + httpUrl
	addPath = strings.ReplaceAll(addPath, "{endpoint_id}", endpointId)
	addPath = strings.ReplaceAll(addPath, "{ipaddress_id}", utils.PathSearch("ip_address_id", rmIpAddresses[0], "").(string))
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("DELETE", addPath, &opt)
	if err != nil {
		// The status code is 404 when unbinding a non-existent IP address.
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			return rmIpAddresses[1:], nil
		}
		return nil, fmt.Errorf("error removing IP addresss from endpoint (%s): %s", endpointId, err)
	}
	return rmIpAddresses[1:], nil
}

func resourceEndpointAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf       = meta.(*config.Config)
		region     = conf.GetRegion(d)
		httplUr    = "v2.1/endpoints/{endpoint_id}"
		endpointId = d.Id()
	)
	client, err := conf.NewServiceClient("dns", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	deletePath := client.Endpoint + httplUr
	deletePath = strings.ReplaceAll(deletePath, "{endpoint_id}", endpointId)
	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &opts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting DNS endpoint (%s): %s", endpointId, err))
	}

	stateConf := &resource.StateChangeConf{
		Target: []string{"DELETED"},
		// Endpoints with status "ERROR" are allowed to be deleted.
		Pending:      []string{"ACTIVE", "PENDING", "ERROR"},
		Refresh:      refreshEndpointStatus(client, endpointId),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for DNS endpoint (%s) deletion to completed: %s", d.Id(), err)
	}
	return nil
}
