package elb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/elb/v3/loadbalancers"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var loadBalancerNonUpdatableParams = []string{
	"loadbalancer_type", "vpc_id", "ipv4_eip_id", "iptype", "bandwidth_charge_mode", "sharetype", "bandwidth_size",
	"bandwidth_id",
}

// @API ELB POST /v3/{project_id}/elb/loadbalancers
// @API ELB POST /v2.0/{project_id}/loadbalancers/{loadbalancer_id}/tags/action
// @API ELB GET /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}
// @API ELB PUT /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}
// @API ELB POST /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/availability-zone/batch-add
// @API ELB POST /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/availability-zone/batch-remove
// @API ELB DELETE /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/force-elb
// @API ELB DELETE /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}
// @API EIP DELETE /v1/{project_id}/publicips/{publicip_id}
// @API ELB POST /v3/{project_id}/elb/loadbalancers/change-charge-mode
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources-migrat
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
func ResourceLoadBalancerV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadBalancerV3Create,
		ReadContext:   resourceLoadBalancerV3Read,
		UpdateContext: resourceLoadBalancerV3Update,
		DeleteContext: resourceLoadBalancerV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(loadBalancerNonUpdatableParams),
			customdiff.ValidateChange("charging_mode", func(_ context.Context, old, new, _ any) error {
				// can only update from postPaid
				if old.(string) != new.(string) && (old.(string) == "prePaid") {
					return fmt.Errorf("charging_mode can only be updated from postPaid to prePaid")
				}
				return nil
			}),
			customdiff.ValidateChange("period_unit", func(_ context.Context, old, new, _ any) error {
				// can only update from empty
				if old.(string) != new.(string) && old.(string) != "" {
					return fmt.Errorf("period_unit can only be updated when changing charging_mode to prePaid")
				}
				return nil
			}),
			customdiff.ValidateChange("period", func(_ context.Context, old, new, _ any) error {
				// can only update from empty
				if old.(int) != new.(int) && old.(int) != 0 {
					return fmt.Errorf("period can only be updated when changing charging_mode to prePaid")
				}
				return nil
			}),
			config.MergeDefaultTags(),
		),

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
			},
			"availability_zone": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"loadbalancer_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"gw_flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cross_vpc_backend": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ipv4_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipv6_network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipv6_bandwidth_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipv4_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ipv6_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ipv4_eip_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ConflictsWith: []string{
					"iptype", "bandwidth_charge_mode", "bandwidth_size", "sharetype", "bandwidth_id",
				},
			},
			"iptype": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"ipv4_eip_id"},
			},
			"bandwidth_charge_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				RequiredWith: []string{
					"iptype", "bandwidth_size", "sharetype",
				},
				ConflictsWith: []string{"ipv4_eip_id", "bandwidth_id"},
			},
			"sharetype": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				RequiredWith: []string{
					"iptype",
				},
				ConflictsWith: []string{"ipv4_eip_id"},
			},
			"bandwidth_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				RequiredWith: []string{
					"iptype", "bandwidth_charge_mode", "sharetype",
				},
				ConflictsWith: []string{"ipv4_eip_id", "bandwidth_id"},
			},
			"bandwidth_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				RequiredWith: []string{
					"iptype",
				},
				ConflictsWith: []string{"ipv4_eip_id", "bandwidth_size", "bandwidth_charge_mode"},
			},
			"l4_flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"l7_flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backend_subnets": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"protection_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"protection_reason": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"deletion_protection_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"waf_failure_action": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			// charge info: charging_mode, period_unit, period, auto_renew, auto_pay
			"charging_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"period"},
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"period_unit"},
			},
			"auto_renew": common.SchemaAutoRenewUpdatable(nil),
			"auto_pay":   common.SchemaAutoPay(nil),
			"charge_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guaranteed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv4_port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv4_eip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_eip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_eip_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"elb_virsubnet_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"frozen_scene": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_border_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"autoscaling_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(``,
					utils.SchemaDescInput{
						Deprecated: true,
					}),
			},
			"min_l7_flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				RequiredWith: []string{
					"l7_flavor_id",
				},
				Description: utils.SchemaDesc(``,
					utils.SchemaDescInput{
						Deprecated: true,
					}),
			},
		},
	}
}

func resourceLoadBalancerV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl    = "v3/{project_id}/elb/loadbalancers"
		product    = "elb"
		bssProduct = "bssv2"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateLoadBalancerBodyParams(d, cfg))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating ELB LoadBalancer: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		id := utils.PathSearch("loadbalancer_id", createRespBody, "").(string)
		if id == "" {
			return diag.Errorf("error creating ELB LoadBalancer: loadbalancer_id is not found in the response")
		}
		d.SetId(id)
		orderId := utils.PathSearch("order_id", createRespBody, "").(string)
		if orderId == "" {
			return diag.Errorf("error creating ELB LoadBalancer: order_id is not found in the response")
		}
		bssClient, err := cfg.NewServiceClient(bssProduct, region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		id := utils.PathSearch("loadbalancer.id", createRespBody, "").(string)
		if id == "" {
			return diag.Errorf("error creating ELB LoadBalancer: loadbalancer id is not found in the response")
		}
		d.SetId(id)
	}
	err = waitForLoadBalancer(ctx, client, d.Id(), "ACTIVE", d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceLoadBalancerV3Read(ctx, d, meta)
}

func buildCreateLoadBalancerBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                       d.Get("name"),
		"availability_zone_list":     d.Get("availability_zone").(*schema.Set).List(),
		"ip_target_enable":           utils.ValueIgnoreEmpty(d.Get("cross_vpc_backend").(bool)),
		"vpc_id":                     utils.ValueIgnoreEmpty(d.Get("vpc_id").(string)),
		"vip_subnet_cidr_id":         utils.ValueIgnoreEmpty(d.Get("ipv4_subnet_id").(string)),
		"ipv6_vip_virsubnet_id":      utils.ValueIgnoreEmpty(d.Get("ipv6_network_id").(string)),
		"vip_address":                utils.ValueIgnoreEmpty(d.Get("ipv4_address").(string)),
		"ipv6_vip_address":           utils.ValueIgnoreEmpty(d.Get("ipv6_address").(string)),
		"l4_flavor_id":               utils.ValueIgnoreEmpty(d.Get("l4_flavor_id").(string)),
		"l7_flavor_id":               utils.ValueIgnoreEmpty(d.Get("l7_flavor_id").(string)),
		"protection_status":          utils.ValueIgnoreEmpty(d.Get("protection_status").(string)),
		"protection_reason":          utils.ValueIgnoreEmpty(d.Get("protection_reason").(string)),
		"loadbalancer_type":          utils.ValueIgnoreEmpty(d.Get("loadbalancer_type").(string)),
		"gw_flavor_id":               utils.ValueIgnoreEmpty(d.Get("gw_flavor_id").(string)),
		"description":                utils.ValueIgnoreEmpty(d.Get("description").(string)),
		"enterprise_project_id":      utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		"deletion_protection_enable": utils.ValueIgnoreEmpty(d.Get("deletion_protection_enable").(bool)),
		"waf_failure_action":         utils.ValueIgnoreEmpty(d.Get("waf_failure_action").(string)),
		"elb_virsubnet_ids":          utils.ValueIgnoreEmpty(d.Get("backend_subnets").(*schema.Set).List()),
		"ipv6_bandwidth":             buildCreateLoadBalancerIpv6BandwidthBodyParams(d),
		"publicip":                   buildCreateLoadBalancerPublicIpBodyParams(d),
		"autoscaling":                buildCreateLoadBalancerAutoscalingBodyParams(d),
		"prepaid_options":            buildCreateLoadBalancerPrepaidOptionsBodyParams(d),
		"tags":                       utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}
	if v, ok := d.GetOk("ipv4_eip_id"); ok {
		bodyParams["publicip_ids"] = []string{v.(string)}
	}

	return map[string]interface{}{
		"loadbalancer": bodyParams,
	}
}

func buildCreateLoadBalancerIpv6BandwidthBodyParams(d *schema.ResourceData) map[string]interface{} {
	v, ok := d.GetOk("ipv6_bandwidth_id")
	if !ok {
		return nil
	}

	bodyParams := map[string]interface{}{
		"id": v.(string),
	}
	return bodyParams
}

func buildCreateLoadBalancerPublicIpBodyParams(d *schema.ResourceData) map[string]interface{} {
	v, ok := d.GetOk("iptype")
	if !ok {
		return nil
	}

	bodyParams := map[string]interface{}{
		"ip_version":   4,
		"network_type": v.(string),
		"bandwidth": map[string]interface{}{
			"id":          utils.ValueIgnoreEmpty(d.Get("bandwidth_id").(string)),
			"name":        utils.ValueIgnoreEmpty(d.Get("name").(string)),
			"size":        utils.ValueIgnoreEmpty(d.Get("bandwidth_size").(int)),
			"charge_mode": utils.ValueIgnoreEmpty(d.Get("bandwidth_charge_mode").(string)),
			"share_type":  utils.ValueIgnoreEmpty(d.Get("sharetype").(string)),
		},
	}
	return bodyParams
}

func buildCreateLoadBalancerAutoscalingBodyParams(d *schema.ResourceData) map[string]interface{} {
	v, ok := d.GetOk("autoscaling_enabled")
	if !ok {
		return nil
	}

	bodyParams := map[string]interface{}{
		"enable":           v.(bool),
		"min_l7_flavor_id": utils.ValueIgnoreEmpty(d.Get("min_l7_flavor_id").(string)),
	}
	return bodyParams
}

func buildCreateLoadBalancerPrepaidOptionsBodyParams(d *schema.ResourceData) map[string]interface{} {
	v, ok := d.GetOk("charging_mode")
	if !ok || v.(string) != "prePaid" {
		return nil
	}

	autoRenew, _ := strconv.ParseBool(d.Get("auto_renew").(string))
	bodyParams := map[string]interface{}{
		"period_type": utils.ValueIgnoreEmpty(d.Get("period_unit").(string)),
		"period_num":  utils.ValueIgnoreEmpty(d.Get("period").(int)),
		"auto_renew":  utils.ValueIgnoreEmpty(autoRenew),
	}
	if d.Get("auto_pay").(string) != "false" {
		bodyParams["auto_pay"] = true
	}
	return bodyParams
}

func resourceLoadBalancerV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "elbv2"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	getRespBody, err := getLoadBalancer(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ELB LoadBalancer")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("loadbalancer.name", getRespBody, nil)),
		d.Set("loadbalancer_type", utils.PathSearch("loadbalancer.loadbalancer_type", getRespBody, nil)),
		d.Set("description", utils.PathSearch("loadbalancer.description", getRespBody, nil)),
		d.Set("availability_zone", utils.PathSearch("loadbalancer.availability_zone_list", getRespBody, nil)),
		d.Set("cross_vpc_backend", utils.PathSearch("loadbalancer.ip_target_enable", getRespBody, false).(bool)),
		d.Set("vpc_id", utils.PathSearch("loadbalancer.vpc_id", getRespBody, nil)),
		d.Set("ipv4_subnet_id", utils.PathSearch("loadbalancer.vip_subnet_cidr_id", getRespBody, nil)),
		d.Set("ipv6_network_id", utils.PathSearch("loadbalancer.ipv6_vip_virsubnet_id", getRespBody, nil)),
		d.Set("ipv4_address", utils.PathSearch("loadbalancer.vip_address", getRespBody, nil)),
		d.Set("ipv4_port_id", utils.PathSearch("loadbalancer.vip_port_id", getRespBody, nil)),
		d.Set("ipv6_address", utils.PathSearch("loadbalancer.ipv6_vip_address", getRespBody, nil)),
		d.Set("l4_flavor_id", utils.PathSearch("loadbalancer.l4_flavor_id", getRespBody, nil)),
		d.Set("l7_flavor_id", utils.PathSearch("loadbalancer.l7_flavor_id", getRespBody, nil)),
		d.Set("gw_flavor_id", utils.PathSearch("loadbalancer.gw_flavor_id", getRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("loadbalancer.enterprise_project_id", getRespBody, nil)),
		d.Set("autoscaling_enabled", utils.PathSearch("loadbalancer.autoscaling.enable", getRespBody, nil)),
		d.Set("min_l7_flavor_id", utils.PathSearch("loadbalancer.autoscaling.min_l7_flavor_id", getRespBody, nil)),
		d.Set("backend_subnets", utils.PathSearch("loadbalancer.elb_virsubnet_ids", getRespBody, nil)),
		d.Set("protection_status", utils.PathSearch("loadbalancer.protection_status", getRespBody, nil)),
		d.Set("protection_reason", utils.PathSearch("loadbalancer.protection_reason", getRespBody, nil)),
		d.Set("charge_mode", utils.PathSearch("loadbalancer.charge_mode", getRespBody, nil)),
		d.Set("elb_virsubnet_type", utils.PathSearch("loadbalancer.elb_virsubnet_type", getRespBody, nil)),
		d.Set("frozen_scene", utils.PathSearch("loadbalancer.frozen_scene", getRespBody, nil)),
		d.Set("operating_status", utils.PathSearch("loadbalancer.operating_status", getRespBody, nil)),
		d.Set("public_border_group", utils.PathSearch("loadbalancer.public_border_group", getRespBody, nil)),
		d.Set("guaranteed", utils.PathSearch("loadbalancer.guaranteed", getRespBody, nil)),
		d.Set("waf_failure_action", utils.PathSearch("loadbalancer.waf_failure_action", getRespBody, nil)),
		d.Set("deletion_protection_enable",
			utils.PathSearch("loadbalancer.deletion_protection_enable", getRespBody, false).(bool)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("loadbalancer.tags", getRespBody, nil))),
		d.Set("created_at", utils.PathSearch("loadbalancer.created_at", getRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("loadbalancer.updated_at", getRespBody, nil)),
	)

	mErr = multierror.Append(mErr, setLoadBalancerEips(d, getRespBody)...)

	if v := utils.PathSearch("loadbalancer.billing_info", getRespBody, "").(string); len(v) > 0 {
		mErr = multierror.Append(mErr, d.Set("charging_mode", "prePaid"))
	} else {
		mErr = multierror.Append(mErr, d.Set("charging_mode", "postPaid"))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func setLoadBalancerEips(d *schema.ResourceData, resp interface{}) []error {
	curArray := utils.PathSearch("loadbalancer.eips", resp, make([]interface{}, 0)).([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	errs := make([]error, 0)
	for _, v := range curArray {
		ipVersion := int(utils.PathSearch("ip_version", v, float64(0)).(float64))
		if ipVersion == 4 {
			errs = append(errs, d.Set("ipv4_eip_id", utils.PathSearch("eip_id", v, nil)))
			errs = append(errs, d.Set("ipv4_eip", utils.PathSearch("eip_address", v, nil)))
		}
		if ipVersion == 6 {
			errs = append(errs, d.Set("ipv6_eip_id", utils.PathSearch("eip_id", v, nil)))
			errs = append(errs, d.Set("ipv6_eip", utils.PathSearch("eip_address", v, nil)))
		}
	}
	return errs
}

func resourceLoadBalancerV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product    = "elbv3"
		bssProduct = "bssv2"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}
	bssClient, err := cfg.NewServiceClient(bssProduct, region)
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}

	if d.HasChange("charging_mode") {
		if err = updateChargingMode(ctx, client, bssClient, d, schema.TimeoutUpdate); err != nil {
			return diag.FromErr(err)
		}
	} else if d.HasChange("auto_renew") {
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), d.Id()); err != nil {
			return diag.Errorf("error updating the auto-renew of the load balancer (%s): %s", d.Id(), err)
		}
	}

	updateLoadBalancerChanges := []string{"name", "description", "cross_vpc_backend", "ipv4_subnet_id", "ipv6_network_id",
		"ipv4_address", "ipv6_address", "l4_flavor_id", "l7_flavor_id", "gw_flavor_id", "autoscaling_enabled",
		"min_l7_flavor_id", "protection_status", "protection_reason", "deletion_protection_enable", "waf_failure_action",
	}
	if d.HasChanges(updateLoadBalancerChanges...) {
		params := buildUpdateLoadBalancerParamsBodyParams(d)
		if err = updateLoadBalancerParams(ctx, client, bssClient, d, params); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("ipv6_bandwidth_id") {
		err = updateLoadBalancerIpv6BandwidthId(ctx, client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// backend_subnets and cross_vpc_backend can not be updated at the same time, otherwise, an error will occur
	if d.HasChange("backend_subnets") {
		if err = updateLoadBalancerBackendSubnets(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		elbV2Client, err := cfg.ElbV2Client(region)
		if err != nil {
			return diag.Errorf("error creating ELB client: %s", err)
		}
		tagErr := utils.UpdateResourceTags(elbV2Client, d, "loadbalancers", d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of LoadBalancer:%s, err:%s", d.Id(), tagErr)
		}
	}

	if d.HasChange("availability_zone") {
		if err = updateAvailabilityZone(ctx, client, bssClient, d); err != nil {
			return diag.Errorf("error updating availability zone of LoadBalancer:%s, err:%s", d.Id(), err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   d.Id(),
			ResourceType: "loadbalancers",
			RegionId:     region,
			ProjectId:    cfg.GetProjectID(region),
		}
		if err = cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceLoadBalancerV3Read(ctx, d, meta)
}

func buildUpdateLoadBalancerParamsBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
		"description":           d.Get("description"),
		"protection_status":     utils.ValueIgnoreEmpty(d.Get("protection_status")),
		"protection_reason":     d.Get("protection_reason"),
		"waf_failure_action":    utils.ValueIgnoreEmpty(d.Get("waf_failure_action")),
		"ipv6_vip_virsubnet_id": utils.ValueIgnoreEmpty(d.Get("ipv6_network_id")),
	}
	if d.HasChange("l4_flavor_id") {
		params["l4_flavor_id"] = d.Get("l4_flavor_id")
	}
	if d.HasChange("l7_flavor_id") {
		params["l7_flavor_id"] = d.Get("l7_flavor_id")
	}
	if d.HasChange("gw_flavor_id") {
		params["gw_flavor_id"] = d.Get("gw_flavor_id")
	}
	if d.HasChange("ipv4_address") {
		params["vip_address"] = d.Get("ipv4_address")
	}
	if d.HasChange("ipv6_address") {
		params["ipv6_vip_address"] = d.Get("ipv6_address")
	}
	if d.HasChange("cross_vpc_backend") {
		params["ip_target_enable"] = d.Get("cross_vpc_backend")
	}
	if d.HasChange("deletion_protection_enable") {
		params["deletion_protection_enable"] = d.Get("deletion_protection_enable")
	}
	if d.Get("loadbalancer_type") != "gateway" || d.HasChange("ipv4_subnet_id") {
		params["vip_subnet_cidr_id"] = utils.ValueIgnoreEmpty(d.Get("ipv4_subnet_id"))
	}
	if d.HasChange("autoscaling_enabled") {
		autoscalingEnabled := d.Get("autoscaling_enabled").(bool)
		autoscaling := map[string]interface{}{
			"enable": autoscalingEnabled,
		}
		if autoscalingEnabled {
			autoscaling["min_l7_flavor_id"] = utils.ValueIgnoreEmpty(d.Get("min_l7_flavor_id").(string))
		} else {
			autoscaling["min_l7_flavor_id"] = ""
		}
		params["autoscaling"] = autoscaling
	} else if d.HasChange("min_l7_flavor_id") && d.Get("autoscaling_enabled").(bool) {
		autoscaling := map[string]interface{}{
			"min_l7_flavor_id": d.Get("min_l7_flavor_id").(string),
		}
		params["autoscaling"] = autoscaling
	}
	if d.Get("charging_mode") == "prePaid" && d.HasChanges("l4_flavor_id", "l7_flavor_id") {
		prepaidPaidParams := map[string]interface{}{
			"change_mode": "immediate",
			"period_type": d.Get("period_unit").(string),
			"period_num":  d.Get("period").(int),
		}
		if d.Get("auto_pay").(string) != "false" {
			prepaidPaidParams["auto_pay"] = true
		}
		params["prepaid_options"] = prepaidPaidParams
	}
	bodyParams := map[string]interface{}{
		"loadbalancer": params,
	}
	return bodyParams
}

func updateAvailabilityZone(ctx context.Context, elbClient, bssClient *golangsdk.ServiceClient,
	d *schema.ResourceData) error {
	oldAvailabilityZone, newAvailabilityZone := d.GetChange("availability_zone")
	addList := newAvailabilityZone.(*schema.Set).Difference(oldAvailabilityZone.(*schema.Set)).List()
	rmList := oldAvailabilityZone.(*schema.Set).Difference(newAvailabilityZone.(*schema.Set)).List()
	err := addAvailabilityZone(ctx, elbClient, bssClient, d, addList)
	if err != nil {
		return err
	}
	err = removeAvailabilityZone(ctx, elbClient, bssClient, d, rmList)
	if err != nil {
		return err
	}
	return nil
}

func addAvailabilityZone(ctx context.Context, elbClient, bssClient *golangsdk.ServiceClient, d *schema.ResourceData,
	addList []interface{}) error {
	if len(addList) == 0 {
		return nil
	}

	var (
		httpUrl = "v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/availability-zone/batch-add"
	)

	updatePath := elbClient.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", elbClient.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{loadbalancer_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"availability_zone_list": utils.ExpandToStringList(addList),
		},
	}

	updateResp, err := elbClient.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error adding AZ to ELB LoadBalancer: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	if d.Get("charging_mode") == "prePaid" {
		orderId := utils.PathSearch("order_id", updateRespBody, "").(string)
		if orderId == "" {
			return errors.New("error adding AZ to ELB LoadBalancer: order_id is not found in the response")
		}
		err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		err = waitForElbV3LoadBalancer(ctx, elbClient, d.Id(), "ACTIVE", nil, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}
	return nil
}

func removeAvailabilityZone(ctx context.Context, elbClient, bssClient *golangsdk.ServiceClient, d *schema.ResourceData,
	rmList []interface{}) error {
	if len(rmList) == 0 {
		return nil
	}

	var (
		httpUrl = "v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/availability-zone/batch-remove"
	)

	updatePath := elbClient.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", elbClient.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{loadbalancer_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"availability_zone_list": utils.ExpandToStringList(rmList),
		},
	}

	updateResp, err := elbClient.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error removing AZ from ELB LoadBalancer: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	if d.Get("charging_mode") == "prePaid" {
		orderId := utils.PathSearch("order_id", updateRespBody, "").(string)
		if orderId == "" {
			return errors.New("error removing AZ from ELB LoadBalancer: order_id is not found in the response")
		}
		err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
		err = waitForElbV3LoadBalancer(ctx, elbClient, d.Id(), "ACTIVE", nil, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceLoadBalancerV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product    = "elb"
		eipProduct = "vpc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		if err = common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()}); err != nil {
			return diag.Errorf("error unsubscribing ELB LoadBalancer : %s", err)
		}
	} else {
		if d.Get("force_delete").(bool) {
			if err = deleteLoadBalancerForce(client, d.Id()); err != nil {
				return common.CheckDeletedDiag(d, err, "error force deleting ELB LoadBalancer")
			}
		} else {
			if err = deleteLoadBalancer(client, d.Id()); err != nil {
				return common.CheckDeletedDiag(d, err, "error deleting ELB LoadBalancer")
			}
		}
	}

	err = waitForLoadBalancer(ctx, client, d.Id(), "DELETED", d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics
	// delete the EIP if necessary
	eipID := d.Get("ipv4_eip_id").(string)
	if _, ok := d.GetOk("iptype"); ok && eipID != "" {
		eipClient, err := cfg.NewServiceClient(eipProduct, region)
		if err != nil {
			clientDiag := diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "failed to create VPC client",
				Detail:   fmt.Sprintf("failed to create VPC client: %s", err),
			}
			diags = append(diags, clientDiag)
		} else {
			err = deleteEip(eipClient, eipID)
			if err != nil {
				var errDefault404 golangsdk.ErrDefault404
				if !errors.As(err, &errDefault404) {
					eipDiag := diag.Diagnostic{
						Severity: diag.Warning,
						Summary:  "failed to delete EIP",
						Detail:   fmt.Sprintf("failed to delete EIP %s: %s", eipID, err),
					}
					diags = append(diags, eipDiag)
				}
			}
		}
	}

	return diags
}

func deleteEip(client *golangsdk.ServiceClient, id string) error {
	var (
		httpUrl = "v1/{project_id}/publicips/{publicip_id}"
	)

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{publicip_id}", id)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	return err
}

func waitForElbV3LoadBalancer(ctx context.Context, elbClient *golangsdk.ServiceClient,
	id string, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for LoadBalancer %s to become %s", id, target)

	stateConf := &resource.StateChangeConf{
		Target:       []string{target},
		Pending:      pending,
		Refresh:      resourceElbV3LoadBalancerRefreshFunc(elbClient, id),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 1 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("error: LoadBalancer %s not found: %s", id, err)
			}
		}
		return fmt.Errorf("error waiting for LoadBalancer %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourceElbV3LoadBalancerRefreshFunc(elbClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		lb, err := loadbalancers.Get(elbClient, id).Extract()
		if err != nil {
			return nil, "", err
		}
		return lb, lb.ProvisioningStatus, nil
	}
}
