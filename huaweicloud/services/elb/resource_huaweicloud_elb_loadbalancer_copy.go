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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var loadBalancerCopyNonUpdatableParams = []string{
	"loadbalancer_id", "reuse_pool", "force_delete",
}

// @API ELB POST /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/clone
// @API ELB GET /v3/{project_id}/elb/jobs/{job_id}
// @API ELB PUT /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}
// @API ELB POST /v3/{project_id}/elb/loadbalancers/change-charge-mode
// @API ELB POST /v2.0/{project_id}/loadbalancers/{loadbalancer_id}/tags/action
// @API ELB GET /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}
// @API ELB GET /v2.0/{project_id}/loadbalancers/{loadbalancer_id}/tags
// @API ELB POST /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/availability-zone/{batch-add}
// @API ELB POST /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/availability-zone/{batch-remove}
// @API ELB DELETE /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/force-elb
// @API ELB DELETE /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources-migrat
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
func ResourceLoadBalancerCopy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadBalancerCopyCreate,
		ReadContext:   resourceLoadBalancerCopyRead,
		UpdateContext: resourceLoadBalancerCopyUpdate,
		DeleteContext: resourceLoadBalancerCopyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(loadBalancerCopyNonUpdatableParams),
			customdiff.ValidateChange("charging_mode", func(_ context.Context, o, n, _ any) error {
				// can only update from postPaid
				if o.(string) != n.(string) && (o.(string) == "prePaid") {
					return errors.New("charging_mode can only be updated from postPaid to prePaid")
				}
				return nil
			}),
			customdiff.ValidateChange("period_unit", func(_ context.Context, o, n, _ any) error {
				// can only update from empty
				if o.(string) != n.(string) && o.(string) != "" {
					return errors.New("period_unit can only be updated when changing charging_mode to prePaid")
				}
				return nil
			}),
			customdiff.ValidateChange("period", func(_ context.Context, o, n, _ any) error {
				// can only update from empty
				if o.(int) != n.(int) && o.(int) != 0 {
					return errors.New("period can only be updated when changing charging_mode to prePaid")
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
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ipv4_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ipv4_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ipv6_network_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ipv6_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"backend_subnets": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"l4_flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"l7_flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"reuse_pool": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipv6_bandwidth_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cross_vpc_backend": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
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
			"deletion_protection_enable": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
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
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"loadbalancer_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
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
			"charge_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv4_port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gw_flavor_id": {
				Type:     schema.TypeString,
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
		},
	}
}

func resourceLoadBalancerCopyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl    = "v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/clone"
		product    = "elbv2"
		bssProduct = "bssv2"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{loadbalancer_id}", d.Get("loadbalancer_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateLoadBalancerCopyBodyParams(d, cfg))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating ELB LoadBalancer copy: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	loadBalancerId := utils.PathSearch("loadbalancer_list[0].id", createRespBody, nil)
	if loadBalancerId == nil {
		return diag.Errorf("error creating ELB LoadBalancer copy: ID is not found in API response")
	}
	d.SetId(loadBalancerId.(string))

	jobId := utils.PathSearch("job_id", createRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error creating ELB LoadBalancer copy: job_id is not found in API response")
	}
	err = checkLoadBalancerJobFinish(ctx, client, jobId.(string), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	if err = initializeLoadBalancerParams(ctx, client, d); err != nil {
		return diag.FromErr(err)
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		bssClient, err := cfg.NewServiceClient(bssProduct, region)
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}
		if err = updateChargingMode(ctx, client, bssClient, d, schema.TimeoutCreate); err != nil {
			return diag.FromErr(err)
		}
	}

	if tagRaw := d.Get("tags").(map[string]interface{}); len(tagRaw) > 0 {
		tagList := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(client, "loadbalancers", d.Id(), tagList).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags for ELB LoadBalancer copy %s: %s", d.Id(), tagErr)
		}
	}
	return resourceLoadBalancerCopyRead(ctx, d, meta)
}

func buildCreateLoadBalancerCopyBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	return map[string]interface{}{
		"target_loadbalancer_param": buildCreateLoadBalancerCopyParamBodyParams(d, cfg),
	}
}

func buildCreateLoadBalancerCopyParamBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	return map[string]interface{}{
		"name":                   utils.ValueIgnoreEmpty(d.Get("name")),
		"availability_zone_list": utils.ValueIgnoreEmpty(d.Get("availability_zone").(*schema.Set).List()),
		"vip_subnet_cidr_id":     utils.ValueIgnoreEmpty(d.Get("ipv4_subnet_id").(string)),
		"vip_address":            utils.ValueIgnoreEmpty(d.Get("ipv4_address").(string)),
		"ipv6_vip_virsubnet_id":  utils.ValueIgnoreEmpty(d.Get("ipv6_network_id").(string)),
		"ipv6_vip_address":       utils.ValueIgnoreEmpty(d.Get("ipv6_address").(string)),
		"elb_virsubnet_ids":      utils.ValueIgnoreEmpty(d.Get("backend_subnets").(*schema.Set).List()),
		"l4_flavor_id":           utils.ValueIgnoreEmpty(d.Get("l4_flavor_id").(string)),
		"l7_flavor_id":           utils.ValueIgnoreEmpty(d.Get("l7_flavor_id").(string)),
		"enterprise_project_id":  utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		"reuse_pool":             utils.ValueIgnoreEmpty(d.Get("reuse_pool").(string)),
		"guaranteed":             true,
	}
}

func buildCreateLoadBalancerParamsBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"description":        utils.ValueIgnoreEmpty(d.Get("description")),
		"protection_status":  utils.ValueIgnoreEmpty(d.Get("protection_status")),
		"protection_reason":  utils.ValueIgnoreEmpty(d.Get("protection_reason")),
		"waf_failure_action": utils.ValueIgnoreEmpty(d.Get("waf_failure_action")),
	}
	if v, ok := d.GetOk("ipv6_bandwidth_id"); ok {
		params["ipv6_bandwidth"] = map[string]interface{}{
			"id": v,
		}
	}
	if v, ok := d.GetOk("cross_vpc_backend"); ok && v == "true" {
		params["ip_target_enable"] = true
	}
	if v, ok := d.GetOk("deletion_protection_enable"); ok {
		if v == "true" {
			params["deletion_protection_enable"] = true
		} else {
			params["deletion_protection_enable"] = false
		}
	}
	bodyParams := utils.RemoveNil(params)
	if len(bodyParams) > 0 {
		return map[string]interface{}{
			"loadbalancer": bodyParams,
		}
	}
	return nil
}

func initializeLoadBalancerParams(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	params := buildCreateLoadBalancerParamsBodyParams(d)
	if len(params) > 0 {
		_, err := updateElbLoadBalancerParams(client, d, utils.RemoveNil(params))
		if err != nil {
			return err
		}
		err = waitForLoadBalancer(ctx, client, d.Id(), "ACTIVE", d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceLoadBalancerCopyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return common.CheckDeletedDiag(d, err, "error retrieving ELB LoadBalancer copy")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("loadbalancer.name", getRespBody, nil)),
		d.Set("loadbalancer_type", utils.PathSearch("loadbalancer.loadbalancer_type", getRespBody, nil)),
		d.Set("description", utils.PathSearch("loadbalancer.description", getRespBody, nil)),
		d.Set("availability_zone", utils.PathSearch("loadbalancer.availability_zone_list", getRespBody, nil)),
		d.Set("cross_vpc_backend", strconv.FormatBool(utils.PathSearch("loadbalancer.ip_target_enable", getRespBody, false).(bool))),
		d.Set("ipv4_subnet_id", utils.PathSearch("loadbalancer.vip_subnet_cidr_id", getRespBody, nil)),
		d.Set("ipv4_address", utils.PathSearch("loadbalancer.vip_address", getRespBody, nil)),
		d.Set("ipv6_network_id", utils.PathSearch("loadbalancer.ipv6_vip_virsubnet_id", getRespBody, nil)),
		d.Set("ipv6_address", utils.PathSearch("loadbalancer.ipv6_vip_address", getRespBody, nil)),
		d.Set("ipv4_port_id", utils.PathSearch("loadbalancer.vip_port_id", getRespBody, nil)),
		d.Set("l4_flavor_id", utils.PathSearch("loadbalancer.l4_flavor_id", getRespBody, nil)),
		d.Set("l7_flavor_id", utils.PathSearch("loadbalancer.l7_flavor_id", getRespBody, nil)),
		d.Set("gw_flavor_id", utils.PathSearch("loadbalancer.gw_flavor_id", getRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("loadbalancer.enterprise_project_id", getRespBody, nil)),
		d.Set("backend_subnets", utils.PathSearch("loadbalancer.elb_virsubnet_ids", getRespBody, nil)),
		d.Set("protection_status", utils.PathSearch("loadbalancer.protection_status", getRespBody, nil)),
		d.Set("protection_reason", utils.PathSearch("loadbalancer.protection_reason", getRespBody, nil)),
		d.Set("charge_mode", utils.PathSearch("loadbalancer.charge_mode", getRespBody, nil)),
		d.Set("created_at", utils.PathSearch("loadbalancer.created_at", getRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("loadbalancer.updated_at", getRespBody, nil)),
		d.Set("waf_failure_action", utils.PathSearch("loadbalancer.waf_failure_action", getRespBody, nil)),
		d.Set("vpc_id", utils.PathSearch("loadbalancer.vpc_id", getRespBody, nil)),
		d.Set("elb_virsubnet_type", utils.PathSearch("loadbalancer.elb_virsubnet_type", getRespBody, nil)),
		d.Set("frozen_scene", utils.PathSearch("loadbalancer.frozen_scene", getRespBody, nil)),
		d.Set("operating_status", utils.PathSearch("loadbalancer.operating_status", getRespBody, nil)),
		d.Set("public_border_group", utils.PathSearch("loadbalancer.public_border_group", getRespBody, nil)),
		d.Set("deletion_protection_enable", strconv.FormatBool(
			utils.PathSearch("loadbalancer.deletion_protection_enable", getRespBody, false).(bool))),
	)
	if v := utils.PathSearch("loadbalancer.billing_info", getRespBody, "").(string); len(v) > 0 {
		mErr = multierror.Append(mErr, d.Set("charging_mode", "prePaid"))
	} else {
		mErr = multierror.Append(mErr, d.Set("charging_mode", "postPaid"))
	}

	if resourceTags, err := tags.Get(client, "loadbalancers", d.Id()).Extract(); err == nil {
		tagMap := utils.TagsToMap(resourceTags.Tags)
		mErr = multierror.Append(mErr, d.Set("tags", tagMap))
	} else {
		log.Printf("[WARN] error fetching tags of ELB LoadBalancer(%s): %s", d.Id(), err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceLoadBalancerCopyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	updateLoadBalancerChanges := []string{
		"name",
		"description",
		"cross_vpc_backend",
		"ipv4_address",
		"ipv6_address",
		"l4_flavor_id",
		"l7_flavor_id",
		"protection_status",
		"protection_reason",
		"ipv4_subnet_id",
		"ipv6_network_id",
		"ipv4_subnet_id",
		"waf_failure_action",
		"deletion_protection_enable",
	}
	if d.HasChanges(updateLoadBalancerChanges...) {
		if err = updateLoadBalancerParams(ctx, client, bssClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("ipv6_bandwidth_id") {
		err = updateLoadBalancerIpv6BandwidthId(ctx, client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("backend_subnets") {
		if err = updateLoadBalancerBackendSubnets(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		elbV2Client, err := cfg.ElbV2Client(region)
		if err != nil {
			return diag.Errorf("error creating ELB 2.0 client: %s", err)
		}
		tagErr := utils.UpdateResourceTags(elbV2Client, d, "loadbalancers", d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of LoadBalancer:%s, err:%s", d.Id(), tagErr)
		}
	}

	// update availability zone
	if d.HasChange("availability_zone") {
		if err = updateAvailabilityZone(ctx, cfg, client, d); err != nil {
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

	if d.HasChange("charging_mode") {
		if err = updateChargingMode(ctx, client, bssClient, d, schema.TimeoutUpdate); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceLoadBalancerCopyRead(ctx, d, meta)
}

func updateLoadBalancerIpv6BandwidthId(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	o, n := d.GetChange("ipv6_bandwidth_id")
	if o.(string) != "" {
		params := buildLoadBalancerParamsIpv6BandwidthIdBodyParams(nil)
		err := updateIpv6BandwidthId(ctx, client, d, params)
		if err != nil {
			return err
		}
	}
	if n.(string) != "" {
		params := buildLoadBalancerParamsIpv6BandwidthIdBodyParams(n.(string))
		err := updateIpv6BandwidthId(ctx, client, d, params)
		if err != nil {
			return err
		}
	}
	return nil
}

func buildLoadBalancerParamsIpv6BandwidthIdBodyParams(ipv6BandwidthId interface{}) map[string]interface{} {
	params := map[string]interface{}{
		"id": ipv6BandwidthId,
	}
	bodyParams := map[string]interface{}{
		"ipv6_bandwidth": params,
	}
	return map[string]interface{}{
		"loadbalancer": bodyParams,
	}
}

func updateIpv6BandwidthId(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	params map[string]interface{}) error {
	_, err := updateElbLoadBalancerParams(client, d, params)
	if err != nil {
		return err
	}

	return waitForLoadBalancer(ctx, client, d.Id(), "ACTIVE", d.Timeout(schema.TimeoutUpdate))
}

func updateLoadBalancerBackendSubnets(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	params := buildUpdateLoadBalancerBackendSubnetsBodyParams(d)
	_, err := updateElbLoadBalancerParams(client, d, utils.RemoveNil(params))
	if err != nil {
		return err
	}

	return waitForLoadBalancer(ctx, client, d.Id(), "ACTIVE", d.Timeout(schema.TimeoutUpdate))
}

func buildUpdateLoadBalancerBackendSubnetsBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"elb_virsubnet_ids":     d.Get("backend_subnets").(*schema.Set).List(),
		"vip_subnet_cidr_id":    utils.ValueIgnoreEmpty(d.Get("ipv4_subnet_id")),
		"ipv6_vip_virsubnet_id": utils.ValueIgnoreEmpty(d.Get("ipv6_network_id")),
	}
	bodyParams := map[string]interface{}{
		"loadbalancer": params,
	}
	return bodyParams
}

func updateLoadBalancerParams(ctx context.Context, client, bssClient *golangsdk.ServiceClient, d *schema.ResourceData) error {
	params := buildUpdateLoadBalancerParamsBodyParams(d)
	updateResp, err := updateElbLoadBalancerParams(client, d, utils.RemoveNil(params))
	if err != nil {
		return err
	}

	orderId := utils.PathSearch("order_id", updateResp, "")
	if orderId != "" {
		err = common.WaitOrderComplete(ctx, bssClient, orderId.(string), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	return waitForLoadBalancer(ctx, client, d.Id(), "ACTIVE", d.Timeout(schema.TimeoutUpdate))
}

func buildUpdateLoadBalancerParamsBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
		"description":           d.Get("description"),
		"l4_flavor_id":          utils.ValueIgnoreEmpty(d.Get("l4_flavor_id")),
		"l7_flavor_id":          utils.ValueIgnoreEmpty(d.Get("l7_flavor_id")),
		"protection_status":     utils.ValueIgnoreEmpty(d.Get("protection_status")),
		"protection_reason":     d.Get("protection_reason"),
		"waf_failure_action":    utils.ValueIgnoreEmpty(d.Get("waf_failure_action")),
		"ipv6_vip_virsubnet_id": utils.ValueIgnoreEmpty(d.Get("ipv6_network_id")),
	}
	if d.HasChange("ipv4_address") {
		params["vip_address"] = d.Get("ipv4_address")
	}
	if d.HasChange("ipv6_address") {
		params["ipv6_vip_address"] = d.Get("ipv6_address")
	}
	if d.HasChange("cross_vpc_backend") {
		crossVpcBackend, _ := strconv.ParseBool(d.Get("cross_vpc_backend").(string))
		params["ip_target_enable"] = crossVpcBackend
	}
	if d.HasChange("deletion_protection_enable") {
		deletionProtectionEnable, _ := strconv.ParseBool(d.Get("deletion_protection_enable").(string))
		params["deletion_protection_enable"] = deletionProtectionEnable
	}
	if d.Get("loadbalancer_type") != "gateway" || d.HasChange("ipv4_subnet_id") {
		params["vip_subnet_cidr_id"] = utils.ValueIgnoreEmpty(d.Get("ipv4_subnet_id"))
	}
	bodyParams := map[string]interface{}{
		"loadbalancer": params,
	}
	return bodyParams
}

func updateElbLoadBalancerParams(client *golangsdk.ServiceClient, d *schema.ResourceData, params map[string]interface{}) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/elb/loadbalancers/{loadbalancer_id}"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{loadbalancer_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         params,
	}

	updateResp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return nil, fmt.Errorf("error updating ELB LoadBalancer: %s", err)
	}

	return utils.FlattenResponse(updateResp)
}

func updateChargingMode(ctx context.Context, client, bssClient *golangsdk.ServiceClient, d *schema.ResourceData, timeout string) error {
	var (
		httpUrl = "v3/{project_id}/elb/loadbalancers/change-charge-mode"
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateChargingModeBodyParams(d))

	updateResp, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating ELB LoadBalancer charging mode: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}
	orderId := utils.PathSearch("order_id", updateRespBody, nil)
	if orderId == nil {
		return errors.New("error updating ELB LoadBalancer charging mode: order_id is not found in API response")
	}

	return common.WaitOrderComplete(ctx, bssClient, orderId.(string), d.Timeout(timeout))
}

func buildUpdateChargingModeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"loadbalancer_ids": []string{d.Id()},
		"charge_mode":      "prepaid",
		"prepaid_options":  buildUpdateChargingModePrepaidOptionsBodyParams(d),
	}
	return bodyParams
}

func buildUpdateChargingModePrepaidOptionsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"period_type": d.Get("period_unit").(string),
		"period_num":  d.Get("period").(int),
		"auto_renew":  d.Get("auto_renew").(string),
		"auto_pay":    true,
	}
	return bodyParams
}

func resourceLoadBalancerCopyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "elbv2"
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
				return common.CheckDeletedDiag(d, err, "error force deleting ELB LoadBalancer copy")
			}
		} else {
			if err = deleteLoadBalancer(client, d.Id()); err != nil {
				return common.CheckDeletedDiag(d, err, "error deleting ELB LoadBalancer copy")
			}
		}
	}

	err = waitForLoadBalancer(ctx, client, d.Id(), "DELETED", d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteLoadBalancer(client *golangsdk.ServiceClient, id string) error {
	var (
		httpUrl = "v3/{project_id}/elb/loadbalancers/{loadbalancer_id}"
	)

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{loadbalancer_id}", id)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	return err
}

func deleteLoadBalancerForce(client *golangsdk.ServiceClient, id string) error {
	var (
		httpUrl = "v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/force-elb"
	)

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{loadbalancer_id}", id)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	return err
}

func waitForLoadBalancer(ctx context.Context, client *golangsdk.ServiceClient, id, target string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Target:       []string{target},
		Refresh:      resourceLoadBalancerRefreshFunc(client, id),
		Timeout:      timeout,
		PollInterval: 2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for ELB LoadBalancer (%s) to be %s: %s ", id, strings.ToLower(target), err)
	}
	return nil
}

func resourceLoadBalancerRefreshFunc(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getRespBody, err := getLoadBalancer(client, id)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "", "DELETED", nil
			}
			return nil, "Failed", err
		}

		status := utils.PathSearch("loadbalancer.provisioning_status", getRespBody, "")
		return getRespBody, status.(string), nil
	}
}

func getLoadBalancer(client *golangsdk.ServiceClient, id string) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/elb/loadbalancers/{loadbalancer_id}"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{loadbalancer_id}", id)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func checkLoadBalancerJobFinish(ctx context.Context, client *golangsdk.ServiceClient, jobID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"RUNNING"},
		Target:       []string{"COMPLETE"},
		Refresh:      loadBalancerJobStatusRefreshFunc(client, jobID),
		Timeout:      timeout,
		PollInterval: 2 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for ELB job (%s) to be completed: %s ", jobID, err)
	}
	return nil
}

func loadBalancerJobStatusRefreshFunc(client *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			httpUrl = "v3/{project_id}/elb/jobs/{job_id}"
		)

		getPath := client.Endpoint + httpUrl
		getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
		getPath = strings.ReplaceAll(getPath, "{job_id}", jobId)

		getOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return nil, "Failed", err
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, "Failed", err
		}

		status := utils.PathSearch("job.status", getRespBody, "")
		return getRespBody, status.(string), nil
	}
}
