package nat

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API NAT POST /v3/{project_id}/private-nat/snat-rules
// @API NAT GET /v3/{project_id}/private-nat/snat-rules/{snat_rule_id}
// @API NAT PUT /v3/{project_id}/private-nat/snat-rules/{snat_rule_id}
// @API NAT DELETE /v3/{project_id}/private-nat/snat-rules/{snat_rule_id}
func ResourcePrivateSnatRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateSnatRuleCreate,
		ReadContext:   resourcePrivateSnatRuleRead,
		UpdateContext: resourcePrivateSnatRuleUpdate,
		DeleteContext: resourcePrivateSnatRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the SNAT rule is located.",
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The private NAT gateway ID to which the SNAT rule belongs.",
			},
			"transit_ip_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the transit IP associated with SNAT rule.",
			},
			"cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"subnet_id"},
				Description:  "The CIDR block of the match rule.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The subnet ID of the match rule.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the SNAT rule.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the SNAT rule.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the SNAT rule.",
			},
			"transit_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The address of the transit IP",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the enterprise project to which the private SNAT rule belongs.",
			},
		},
	}
}

func buildCreatePrivateSnatRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	snatRuleBodyParams := map[string]interface{}{
		"gateway_id":     d.Get("gateway_id"),
		"transit_ip_ids": []string{d.Get("transit_ip_id").(string)},
		"cidr":           utils.ValueIgnoreEmpty(d.Get("cidr")),
		"virsubnet_id":   utils.ValueIgnoreEmpty(d.Get("subnet_id")),
		"description":    utils.ValueIgnoreEmpty(d.Get("description")),
	}

	return map[string]interface{}{
		"snat_rule": snatRuleBodyParams,
	}
}

func resourcePrivateSnatRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/private-nat/snat-rules"
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreatePrivateSnatRuleBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating private SNAT rule: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleId := utils.PathSearch("snat_rule.id", respBody, "").(string)
	if ruleId == "" {
		return diag.Errorf("error creating private SNAT rule: ID is not found in API response")
	}

	d.SetId(ruleId)

	return resourcePrivateSnatRuleRead(ctx, d, meta)
}

func GetPrivateSnatRule(client *golangsdk.ServiceClient, ruleId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/private-nat/snat-rules/{snat_rule_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{snat_rule_id}", ruleId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func resourcePrivateSnatRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	respBody, err := GetPrivateSnatRule(client, d.Id())
	if err != nil {
		// If the private SNAT rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving private SNAT rule")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("gateway_id", utils.PathSearch("snat_rule.gateway_id", respBody, nil)),
		d.Set("transit_ip_id", utils.PathSearch("snat_rule.transit_ip_associations[0].transit_ip_id", respBody, nil)),
		d.Set("description", utils.PathSearch("snat_rule.description", respBody, nil)),
		d.Set("subnet_id", utils.PathSearch("snat_rule.virsubnet_id", respBody, nil)),
		d.Set("cidr", utils.PathSearch("snat_rule.cidr", respBody, nil)),
		d.Set("created_at", utils.PathSearch("snat_rule.created_at", respBody, nil)),
		d.Set("updated_at", utils.PathSearch("snat_rule.updated_at", respBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("snat_rule.enterprise_project_id", respBody, nil)),
		d.Set("transit_ip_address", utils.PathSearch("snat_rule.transit_ip_associations[0].transit_ip_address", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdatePrivateSnatRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	snatRuleBodyParams := map[string]interface{}{
		"transit_ip_ids": []string{d.Get("transit_ip_id").(string)},
		"description":    d.Get("description"),
	}

	return map[string]interface{}{
		"snat_rule": snatRuleBodyParams,
	}
}

func resourcePrivateSnatRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/private-nat/snat-rules/{snat_rule_id}"
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{snat_rule_id}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdatePrivateSnatRuleBodyParams(d),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating private SNAT rule (%s): %s", d.Id(), err)
	}

	return resourcePrivateSnatRuleRead(ctx, d, meta)
}

func resourcePrivateSnatRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/private-nat/snat-rules/{snat_rule_id}"
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{snat_rule_id}", d.Id())
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		// If the private SNAT rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting private SNAT rule")
	}

	return nil
}
