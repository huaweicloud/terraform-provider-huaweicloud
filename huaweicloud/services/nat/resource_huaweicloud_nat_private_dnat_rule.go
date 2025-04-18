package nat

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API NAT POST /v3/{project_id}/private-nat/dnat-rules
// @API NAT GET /v3/{project_id}/private-nat/dnat-rules/{dnat_rule_id}
// @API NAT PUT /v3/{project_id}/private-nat/dnat-rules/{dnat_rule_id}
// @API NAT DELETE /v3/{project_id}/private-nat/dnat-rules/{dnat_rule_id}
func ResourcePrivateDnatRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateDnatRuleCreate,
		ReadContext:   resourcePrivateDnatRuleRead,
		UpdateContext: resourcePrivateDnatRuleUpdate,
		DeleteContext: resourcePrivateDnatRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the DNAT rule is located.",
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The private NAT gateway ID to which the DNAT rule belongs.",
			},
			"transit_ip_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the transit IP for private NAT.",
			},
			"transit_service_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The port of the transit IP.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The protocol type.",
			},
			"backend_interface_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The network interface ID of the transit IP for private NAT.",
				ExactlyOneOf: []string{"backend_private_ip"},
			},
			"backend_private_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The private IP address of the backend instance.",
			},
			"internal_service_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The port of the backend instance.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the DNAT rule.",
			},
			"backend_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of backend instance.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the DNAT rule.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the DNAT rule.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the enterprise project to which the private DNAT rule belongs.",
			},
		},
	}
}

func buildCreatePrivateDnatRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	dnatRuleBodyParams := map[string]interface{}{
		"gateway_id":            d.Get("gateway_id"),
		"transit_ip_id":         d.Get("transit_ip_id"),
		"protocol":              utils.ValueIgnoreEmpty(d.Get("protocol")),
		"network_interface_id":  utils.ValueIgnoreEmpty(d.Get("backend_interface_id")),
		"private_ip_address":    utils.ValueIgnoreEmpty(d.Get("backend_private_ip")),
		"internal_service_port": convertIntToStr(d.Get("internal_service_port").(int)),
		"transit_service_port":  convertIntToStr(d.Get("transit_service_port").(int)),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
	}

	return map[string]interface{}{
		"dnat_rule": dnatRuleBodyParams,
	}
}

func convertIntToStr(param int) string {
	return fmt.Sprintf("%v", param)
}

func resourcePrivateDnatRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/private-nat/dnat-rules"
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreatePrivateDnatRuleBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating private DNAT rule: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleId := utils.PathSearch("dnat_rule.id", respBody, "").(string)
	if ruleId == "" {
		return diag.Errorf("error creating private DNAT rule: ID is not found in API response")
	}

	d.SetId(ruleId)

	return resourcePrivateDnatRuleRead(ctx, d, meta)
}

func GetPrivateDnatRule(client *golangsdk.ServiceClient, ruleId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/private-nat/dnat-rules/{dnat_rule_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{dnat_rule_id}", ruleId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func resourcePrivateDnatRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	respBody, err := GetPrivateDnatRule(client, d.Id())
	if err != nil {
		// If the private DNAT rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving private DNAT rule")
	}
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("gateway_id", utils.PathSearch("dnat_rule.gateway_id", respBody, nil)),
		d.Set("transit_ip_id", utils.PathSearch("dnat_rule.transit_ip_id", respBody, nil)),
		d.Set("description", utils.PathSearch("dnat_rule.description", respBody, nil)),
		d.Set("backend_interface_id", utils.PathSearch("dnat_rule.network_interface_id", respBody, nil)),
		d.Set("protocol", utils.PathSearch("dnat_rule.protocol", respBody, nil)),
		d.Set("backend_private_ip", utils.PathSearch("dnat_rule.private_ip_address", respBody, nil)),
		d.Set("backend_type", utils.PathSearch("dnat_rule.type", respBody, nil)),
		d.Set("created_at", utils.PathSearch("dnat_rule.created_at", respBody, nil)),
		d.Set("updated_at", utils.PathSearch("dnat_rule.updated_at", respBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("dnat_rule.enterprise_project_id", respBody, nil)),
	)

	internalPort := utils.PathSearch("dnat_rule.internal_service_port", respBody, "").(string)
	transitPort := utils.PathSearch("dnat_rule.transit_service_port", respBody, "").(string)

	// Parse the internal service port
	if internalServicePort, err := strconv.Atoi(internalPort); err != nil {
		mErr = multierror.Append(mErr, fmt.Errorf("invalid format for internal service port, want 'string', but '%T'",
			internalPort))
	} else if internalServicePort != 0 {
		mErr = multierror.Append(mErr, d.Set("internal_service_port", internalServicePort))
	}

	// Parse the transit service port
	if transitServicePort, err := strconv.Atoi(transitPort); err != nil {
		mErr = multierror.Append(mErr, fmt.Errorf("invalid format for internal service port, want 'string', but '%T'",
			transitPort))
	} else if transitServicePort != 0 {
		mErr = multierror.Append(mErr, d.Set("transit_service_port", transitServicePort))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdatePrivateDnatRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	dnatRuleBodyParams := map[string]interface{}{
		"transit_ip_id":         d.Get("transit_ip_id"),
		"protocol":              utils.ValueIgnoreEmpty(d.Get("protocol")),
		"internal_service_port": convertIntToStr(d.Get("internal_service_port").(int)),
		"transit_service_port":  convertIntToStr(d.Get("transit_service_port").(int)),
		"description":           d.Get("description"),
	}

	return dnatRuleBodyParams
}

func updatePrivateDnatRule(client *golangsdk.ServiceClient, ruleId string, dnatRuleBodyParams map[string]interface{}) error {
	httpUrl := "v3/{project_id}/private-nat/dnat-rules/{dnat_rule_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{dnat_rule_id}", ruleId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"dnat_rule": dnatRuleBodyParams,
		},
	}
	_, err := client.Request("PUT", updatePath, &opt)
	return err
}

func resourcePrivateDnatRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	updateOpt := buildUpdatePrivateDnatRuleBodyParams(d)
	if d.HasChange("backend_private_ip") {
		updateOpt["private_ip_address"] = utils.ValueIgnoreEmpty(d.Get("backend_private_ip"))
	} else if d.HasChange("backend_interface_id") {
		updateOpt["network_interface_id"] = utils.ValueIgnoreEmpty(d.Get("backend_interface_id"))
	}

	err = updatePrivateDnatRule(client, d.Id(), updateOpt)
	if err != nil {
		return diag.Errorf("error updating private DNAT rule: %s", err)
	}

	return resourcePrivateDnatRuleRead(ctx, d, meta)
}

func resourcePrivateDnatRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/private-nat/dnat-rules/{dnat_rule_id}"
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v3 client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{dnat_rule_id}", d.Id())
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		// If the private DNAT rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting private DNAT rule")
	}

	return nil
}
