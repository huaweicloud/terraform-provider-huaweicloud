package nat

import (
	"context"
	"fmt"
	"log"
	"strconv"
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

		CustomizeDiff: config.MergeDefaultTags(),

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
				Type:        schema.TypeString,
				Required:    true,
				Description: "The specification of the NAT gateway.",
			},
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),
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
				Description: "The IP address used for the NG port of the NAT gateway.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The enterprise project ID of the NAT gateway.",
			},
			"session_conf": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "The session configuration of the NAT gateway.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tcp_session_expire_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The TCP session expiration time.",
						},
						"udp_session_expire_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The UDP session expiration time.",
						},
						"icmp_session_expire_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The ICMP session expiration time.",
						},
						"tcp_time_wait_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The duration of TIME_WAIT state when TCP connection is closed.",
						},
					},
				},
			},
			"tags": common.TagsSchema(),
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the NAT gateway.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the NAT gateway.",
			},
			"billing_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The order information of the NAT gateway.",
			},
			"dnat_rules_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of DNAT rules on the NAT gateway.",
			},
			"snat_rule_public_ip_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of SNAT rules on the NAT gateway.",
			},
			"pps_max": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of packets that the NAT gateway can receive or send per second.",
			},
			"bps_max": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The bandwidth that the NAT gateway can receive or send per second.",
			},
		},
	}
}

func buildSessionConfigBodyParams(sessionConfig []interface{}) map[string]interface{} {
	if len(sessionConfig) == 0 {
		return nil
	}

	rawSession := sessionConfig[0].(map[string]interface{})
	sessionConfigParams := map[string]interface{}{
		"tcp_session_expire_time":  utils.ValueIgnoreEmpty(rawSession["tcp_session_expire_time"]),
		"udp_session_expire_time":  utils.ValueIgnoreEmpty(rawSession["udp_session_expire_time"]),
		"icmp_session_expire_time": utils.ValueIgnoreEmpty(rawSession["icmp_session_expire_time"]),
		"tcp_time_wait_time":       rawSession["tcp_time_wait_time"],
	}

	return sessionConfigParams
}

func buildPrepaidOptionsBodyParams(d *schema.ResourceData) map[string]interface{} {
	autoRenew, err := strconv.ParseBool(d.Get("auto_renew").(string))
	if err != nil {
		log.Printf("[WARN] error parsing auto-renew to boolean value: %s", err)
	}

	return map[string]interface{}{
		"period_type":   d.Get("period_unit"),
		"period_num":    d.Get("period"),
		"is_auto_renew": autoRenew,
		"is_auto_pay":   true,
	}
}

func buildCreatePublicGatewayBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	natGatewayBodyParams := map[string]interface{}{
		"name":                  d.Get("name"),
		"router_id":             d.Get("vpc_id"),
		"internal_network_id":   d.Get("subnet_id"),
		"spec":                  d.Get("spec"),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"ngport_ip_address":     utils.ValueIgnoreEmpty(d.Get("ngport_ip_address")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		"session_conf":          buildSessionConfigBodyParams(d.Get("session_conf").([]interface{})),
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		natGatewayBodyParams["prepaid_options"] = buildPrepaidOptionsBodyParams(d)
	}

	return map[string]interface{}{
		"nat_gateway": natGatewayBodyParams,
	}
}

func ReadPublicGateway(client *golangsdk.ServiceClient, gatewayID string) (interface{}, error) {
	getPath := client.Endpoint + "v2/{project_id}/nat_gateways/{nat_gateway_id}"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{nat_gateway_id}", gatewayID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func waitingForPublicGatewayStatusActive(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := ReadPublicGateway(client, d.Id())
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("nat_gateway.status", respBody, "").(string)
			if status == "" {
				return nil, "ERROR", fmt.Errorf("status is not found in API response")
			}

			if "INACTIVE" == status {
				return nil, "ERROR", fmt.Errorf("unexpect status (%s)", status)
			}

			if "ACTIVE" == status {
				return "success", "COMPLETED", nil
			}
			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourcePublicGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/nat_gateways"
		product = "nat"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreatePublicGatewayBodyParams(d, cfg)),
	}
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating NAT gateway: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		orderID := utils.PathSearch("order_id", createRespBody, "").(string)
		if orderID == "" {
			return diag.Errorf("error creating prepaid NAT gateway: order_id is not found in API response")
		}

		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}

		if err := common.WaitOrderComplete(ctx, bssClient, orderID, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.FromErr(err)
		}
		resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, orderID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("error waiting for NAT gateway order (%s) complete: %s", orderID, err)
		}
		d.SetId(resourceId)
	} else {
		id := utils.PathSearch("nat_gateway.id", createRespBody, "").(string)
		if id == "" {
			return diag.Errorf("error creating postpaid NAT gateway: ID is not found in API response")
		}
		d.SetId(id)
	}

	if err := waitingForPublicGatewayStatusActive(ctx, client, d, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for NAT gateway (%s) to become active in"+
			" creation operation: %s", d.Id(), err)
	}

	networkClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC v2.0 client: %s", err)
	}
	if err := utils.CreateResourceTags(networkClient, d, "nat_gateways", d.Id()); err != nil {
		return diag.Errorf("error setting tags to the NAT gateway: %s", err)
	}
	return resourcePublicGatewayRead(ctx, d, meta)
}

func resourcePublicGatewayRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "nat"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	networkClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC v2.0 client: %s", err)
	}

	respBody, err := ReadPublicGateway(client, d.Id())
	if err != nil {
		// If the NAT gateway does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving NAT gateway")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("nat_gateway.name", respBody, nil)),
		d.Set("spec", utils.PathSearch("nat_gateway.spec", respBody, nil)),
		d.Set("vpc_id", utils.PathSearch("nat_gateway.router_id", respBody, nil)),
		d.Set("subnet_id", utils.PathSearch("nat_gateway.internal_network_id", respBody, nil)),
		d.Set("description", utils.PathSearch("nat_gateway.description", respBody, nil)),
		d.Set("ngport_ip_address", utils.PathSearch("nat_gateway.ngport_ip_address", respBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("nat_gateway.enterprise_project_id", respBody, nil)),
		d.Set("status", utils.PathSearch("nat_gateway.status", respBody, nil)),
		d.Set("session_conf", flattenGatewaySessionConfig(utils.PathSearch("nat_gateway.session_conf", respBody, nil))),
		d.Set("created_at", utils.PathSearch("nat_gateway.created_at", respBody, nil)),
		d.Set("billing_info", utils.PathSearch("nat_gateway.billing_info", respBody, nil)),
		d.Set("dnat_rules_limit", utils.PathSearch("nat_gateway.dnat_rules_limit", respBody, nil)),
		d.Set("snat_rule_public_ip_limit", utils.PathSearch("nat_gateway.snat_rule_public_ip_limit", respBody, nil)),
		d.Set("pps_max", utils.PathSearch("nat_gateway.pps_max", respBody, nil)),
		d.Set("bps_max", utils.PathSearch("nat_gateway.bps_max", respBody, nil)),
		utils.SetResourceTagsToState(d, networkClient, "nat_gateways", d.Id()),
		d.Set("tags", d.Get("tags")),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving NAT gateway fields: %s", err)
	}
	return nil
}

func flattenGatewaySessionConfig(sessionConfig interface{}) []map[string]interface{} {
	if sessionConfig == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"tcp_session_expire_time":  utils.PathSearch("tcp_session_expire_time", sessionConfig, nil),
			"udp_session_expire_time":  utils.PathSearch("udp_session_expire_time", sessionConfig, nil),
			"icmp_session_expire_time": utils.PathSearch("icmp_session_expire_time", sessionConfig, nil),
			"tcp_time_wait_time":       utils.PathSearch("tcp_time_wait_time", sessionConfig, nil),
		},
	}
}

func buildUpdatePublicGatewayBodyParams(d *schema.ResourceData) map[string]interface{} {
	natGatewayBodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"spec":        utils.ValueIgnoreEmpty(d.Get("spec")),
		"description": d.Get("description"),
	}
	return map[string]interface{}{
		"nat_gateway": natGatewayBodyParams,
	}
}

func buildUpdateSessionConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	sessionConfParams := map[string]interface{}{
		"session_conf": buildSessionConfigBodyParams(d.Get("session_conf").([]interface{})),
	}
	return map[string]interface{}{
		"nat_gateway": sessionConfParams,
	}
}

func resourcePublicGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		product      = "nat"
		chargingMode = d.Get("charging_mode").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	updatePath := client.Endpoint + "v2/{project_id}/nat_gateways/{nat_gateway_id}"
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{nat_gateway_id}", d.Id())

	// Due to API limitations, the prepaid NAT gateway does not support editing.
	if d.HasChanges("name", "spec", "description") && chargingMode != "prePaid" {
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdatePublicGatewayBodyParams(d)),
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating NAT gateway (%s): %s", d.Id(), err)
		}

		if err := waitingForPublicGatewayStatusActive(ctx, client, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error waiting for NAT gateway (%s) to become active in"+
				" update operation: %s", d.Id(), err)
		}
	}

	if d.HasChange("session_conf") {
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateSessionConfigBodyParams(d)),
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating NAT gateway (%s): %s", d.Id(), err)
		}

		if err := waitingForPublicGatewayStatusActive(ctx, client, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error waiting for NAT gateway (%s) to become active in"+
				" update operation: %s", d.Id(), err)
		}
	}

	if d.HasChange("tags") {
		networkClient, err := cfg.NetworkingV2Client(region)
		if err != nil {
			return diag.Errorf("error creating VPC v2.0 client: %s", err)
		}
		err = utils.UpdateResourceTags(networkClient, d, "nat_gateways", d.Id())
		if err != nil {
			return diag.Errorf("error updating tags of the NAT gateway: %s", err)
		}
	}

	// Only prepaid NAT gateway supports editing `auto_renew`.
	if d.HasChange("auto_renew") && chargingMode == "prePaid" {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), d.Id()); err != nil {
			return diag.Errorf("error updating the auto-renew of the NAT gateway (%s): %s", d.Id(), err)
		}
	}

	return resourcePublicGatewayRead(ctx, d, meta)
}

func waitingForPublicGatewayDelete(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := ReadPublicGateway(client, d.Id())
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return "deleted", "COMPLETED", nil
				}
				return nil, "ERROR", err
			}
			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func deletePostpaidGateway(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	deletePath := client.Endpoint + "v2/{project_id}/nat_gateways/{nat_gateway_id}"
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{nat_gateway_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200, 202, 204},
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	return err
}

func resourcePublicGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "nat"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		if err := common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()}); err != nil {
			return diag.Errorf("error unsubscribing NAT gateway (%s): %s", d.Id(), err)
		}
	} else {
		if err := deletePostpaidGateway(client, d); err != nil {
			// If the NAT gateway does not exist, the response HTTP status code of the details API is 404.
			return common.CheckDeletedDiag(d, err, "err deleting NAT gateway")
		}
	}

	if err := waitingForPublicGatewayDelete(ctx, client, d, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for NAT gateway (%s) deleted: %s", d.Id(), err)
	}

	return nil
}
