package vpc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPC POST /v3/{project_id}/vpc/traffic-mirror-filter-rules
// @API VPC GET  /v3/{project_id}/vpc/traffic-mirror-filter-rules/{traffic_mirror_filter_rule_id}
// @API VPC PUT  /v3/{project_id}/vpc/traffic-mirror-filter-rules/{traffic_mirror_filter_rule_id}
// @API VPC DELETE  /v3/{project_id}/vpc/traffic-mirror-filter-rules/{traffic_mirror_filter_rule_id}
func ResourceTrafficMirrorFilterRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTrafficMirrorFilterRuleCreate,
		ReadContext:   resourceTrafficMirrorFilterRuleRead,
		UpdateContext: resourceTrafficMirrorFilterRuleUpdate,
		DeleteContext: resourceTrafficMirrorFilterRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"traffic_mirror_filter_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ethertype": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"source_cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_port_range": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_port_range": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
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

func buildTrafficMirrorFilterRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"protocol":               d.Get("protocol"),
		"ethertype":              d.Get("ethertype"),
		"action":                 d.Get("action"),
		"priority":               d.Get("priority"),
		"source_cidr_block":      utils.ValueIgnoreEmpty(d.Get("source_cidr_block")),
		"destination_cidr_block": utils.ValueIgnoreEmpty(d.Get("destination_cidr_block")),
		"source_port_range":      utils.ValueIgnoreEmpty(d.Get("source_port_range")),
		"destination_port_range": utils.ValueIgnoreEmpty(d.Get("destination_port_range")),
		"description":            d.Get("description"),
	}
	if d.Id() == "" {
		params["traffic_mirror_filter_id"] = d.Get("traffic_mirror_filter_id")
		params["direction"] = d.Get("direction")
	}

	bodyParams := map[string]interface{}{
		"traffic_mirror_filter_rule": params,
	}

	return bodyParams
}

func resourceTrafficMirrorFilterRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	ctreateTrafficMirrorFilterRulePath := client.ResourceBaseURL() + "vpc/traffic-mirror-filter-rules"
	createTrafficMirrorFilterRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	createTrafficMirrorFilterRuleOpt.JSONBody = utils.RemoveNil(buildTrafficMirrorFilterRuleBodyParams(d))
	createTrafficMirrorFilterRuleResp, err := client.Request("POST", ctreateTrafficMirrorFilterRulePath, &createTrafficMirrorFilterRuleOpt)
	if err != nil {
		return diag.Errorf("error creating traffic mirror filter rule: %s", err)
	}

	createTrafficMirrorFilterRuleRespBody, err := utils.FlattenResponse(createTrafficMirrorFilterRuleResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("traffic_mirror_filter_rule.id", createTrafficMirrorFilterRuleRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating traffic mirror filter rule: ID is not found in API response")
	}
	d.SetId(id)

	return resourceTrafficMirrorFilterRuleRead(ctx, d, meta)
}

func resourceTrafficMirrorFilterRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	getTrafficMirrorFilterRulePath := client.ResourceBaseURL() + "vpc/traffic-mirror-filter-rules/" + d.Id()
	getTrafficMirrorFilterRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getTrafficMirrorFilterRuleResp, err := client.Request("GET", getTrafficMirrorFilterRulePath, &getTrafficMirrorFilterRuleOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "VPC traffic mirror filter rule")
	}

	respBody, err := utils.FlattenResponse(getTrafficMirrorFilterRuleResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("traffic_mirror_filter_id", utils.PathSearch("traffic_mirror_filter_rule.traffic_mirror_filter_id", respBody, nil)),
		d.Set("direction", utils.PathSearch("traffic_mirror_filter_rule.direction", respBody, nil)),
		d.Set("protocol", utils.PathSearch("traffic_mirror_filter_rule.protocol", respBody, nil)),
		d.Set("ethertype", utils.PathSearch("traffic_mirror_filter_rule.ethertype", respBody, nil)),
		d.Set("action", utils.PathSearch("traffic_mirror_filter_rule.action", respBody, nil)),
		d.Set("priority", utils.PathSearch("traffic_mirror_filter_rule.priority", respBody, nil)),
		d.Set("source_cidr_block", utils.PathSearch("traffic_mirror_filter_rule.source_cidr_block", respBody, nil)),
		d.Set("source_port_range", utils.PathSearch("traffic_mirror_filter_rule.source_port_range", respBody, nil)),
		d.Set("destination_cidr_block", utils.PathSearch("traffic_mirror_filter_rule.destination_cidr_block", respBody, nil)),
		d.Set("destination_port_range", utils.PathSearch("traffic_mirror_filter_rule.destination_port_range", respBody, nil)),
		d.Set("description", utils.PathSearch("traffic_mirror_filter_rule.description", respBody, nil)),
		d.Set("created_at", utils.PathSearch("traffic_mirror_filter_rule.created_at", respBody, nil)),
		d.Set("updated_at", utils.PathSearch("traffic_mirror_filter_rule.updated_at", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceTrafficMirrorFilterRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	updateTrafficMirrorFilterRuleChanges := []string{
		"description",
		"protocol",
		"ethertype",
		"source_cidr_block",
		"destination_cidr_block",
		"source_port_range",
		"destination_port_range",
		"priority",
		"action",
	}
	if d.HasChanges(updateTrafficMirrorFilterRuleChanges...) {
		updateTrafficMirrorFilterRuleOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateTrafficMirrorFilterRulePath := client.ResourceBaseURL() + "vpc/traffic-mirror-filter-rules/" + d.Id()
		updateTrafficMirrorFilterRuleOpt.JSONBody = utils.RemoveNil(buildTrafficMirrorFilterRuleBodyParams(d))
		_, err = client.Request("PUT", updateTrafficMirrorFilterRulePath, &updateTrafficMirrorFilterRuleOpt)
		if err != nil {
			return diag.Errorf("error updating traffic mirror filter rule: %s", err)
		}
	}

	return resourceTrafficMirrorFilterRuleRead(ctx, d, meta)
}

func resourceTrafficMirrorFilterRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vpcv3", region)
	if err != nil {
		return diag.Errorf("error creating VPC v3 client: %s", err)
	}

	deleteTrafficMirrorFilterRulePath := client.ResourceBaseURL() + "vpc/traffic-mirror-filter-rules/" + d.Id()
	deleteTrafficMirrorFilterRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = client.Request("DELETE", deleteTrafficMirrorFilterRulePath, &deleteTrafficMirrorFilterRuleOpt)
	if err != nil {
		return diag.Errorf("error deleting traffic mirror filter rule: %s", err)
	}
	return nil
}
