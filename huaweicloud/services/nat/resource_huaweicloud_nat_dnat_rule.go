package nat

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

// @API NAT POST /v2/{project_id}/dnat_rules
// @API NAT GET /v2/{project_id}/dnat_rules/{dnat_rule_id}
// @API NAT PUT /v2/{project_id}/dnat_rules/{dnat_rule_id}
// @API NAT DELETE /v2/{project_id}/nat_gateways/{nat_gateway_id}/dnat_rules/{dnat_rule_id}
func ResourcePublicDnatRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePublicDnatRuleCreate,
		ReadContext:   resourcePublicDnatRuleRead,
		UpdateContext: resourcePublicDnatRuleUpdate,
		DeleteContext: resourcePublicDnatRuleDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

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
			"floating_ip_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"floating_ip_id", "global_eip_id"},
				Description:  "The ID of the floating IP address.",
			},
			"global_eip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the global EIP connected by the DNAT rule.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The protocol type.",
			},
			"nat_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the NAT gateway to which the DNAT rule belongs.",
			},
			"internal_service_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ExactlyOneOf: []string{"internal_service_port_range"},
				RequiredWith: []string{"external_service_port"},
				Description:  "The port used by Floating IP provide services for external systems.",
			},
			"internal_service_port_range": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"external_service_port_range"},
				Description:  "The port used by ECSs or BMSs to provide services for external systems.",
			},
			"external_service_port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ExactlyOneOf: []string{"external_service_port_range"},
				Description:  "The port range used by Floating IP provide services for external systems.",
			},
			"external_service_port_range": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The port range used by ECSs or BMSs to provide services for external systems.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the DNAT rule.",
			},
			"port_id": {
				Type:         schema.TypeString,
				ExactlyOneOf: []string{"port_id", "private_ip"},
				Optional:     true,
				Computed:     true,
				Description:  "The port ID of network.",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The private IP address of a user.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the DNAT rule.",
			},
			"floating_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The floating IP address of the DNAT rule.",
			},
			"global_eip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The global EIP address connected by the DNAT rule.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the DNAT rule.",
			},
		},
	}
}

func buildCreateDnatRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	dnatRuleBodyParams := map[string]interface{}{
		"nat_gateway_id":              d.Get("nat_gateway_id"),
		"protocol":                    d.Get("protocol"),
		"internal_service_port":       d.Get("internal_service_port"),
		"external_service_port":       d.Get("external_service_port"),
		"floating_ip_id":              utils.ValueIgnoreEmpty(d.Get("floating_ip_id")),
		"global_eip_id":               utils.ValueIgnoreEmpty(d.Get("global_eip_id")),
		"internal_service_port_range": utils.ValueIgnoreEmpty(d.Get("internal_service_port_range")),
		"external_service_port_range": utils.ValueIgnoreEmpty(d.Get("external_service_port_range")),
		"description":                 utils.ValueIgnoreEmpty(d.Get("description")),
		"port_id":                     utils.ValueIgnoreEmpty(d.Get("port_id")),
		"private_ip":                  utils.ValueIgnoreEmpty(d.Get("private_ip")),
	}

	return map[string]interface{}{
		"dnat_rule": dnatRuleBodyParams,
	}
}

func waitingForDnatRuleStateRefresh(client *golangsdk.ServiceClient, ruleId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := GetDnatRule(client, ruleId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "Resource Not Found", "COMPLETED", nil
			}

			return resp, "", err
		}

		status := utils.PathSearch("dnat_rule.status", resp, "").(string)

		if utils.StrSliceContains([]string{"INACTIVE", "EIP_FREEZED"}, status) {
			return resp, "", fmt.Errorf("unexpect status (%s)", status)
		}

		if utils.StrSliceContains(targets, status) {
			return resp, "COMPLETED", nil
		}

		return resp, "PENDING", nil
	}
}

func resourcePublicDnatRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/dnat_rules"
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateDnatRuleBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DNAT rule: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleId := utils.PathSearch("dnat_rule.id", respBody, "").(string)
	if ruleId == "" {
		return diag.Errorf("error creating DNAT rule: ID is not found in API response")
	}

	d.SetId(ruleId)

	err = waitingForDnatRuleStateCompleted(ctx, client, d.Timeout(schema.TimeoutCreate), ruleId, []string{"ACTIVE"})
	if err != nil {
		return diag.Errorf("error waiting for the DNAT rule (%s) creation to complete: %s", ruleId, err)
	}

	return resourcePublicDnatRuleRead(ctx, d, meta)
}

func GetDnatRule(client *golangsdk.ServiceClient, ruleId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/dnat_rules/{dnat_rule_id}"
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

func resourcePublicDnatRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	respBody, err := GetDnatRule(client, d.Id())
	if err != nil {
		// If the DNAT rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving DNAT rule")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("nat_gateway_id", utils.PathSearch("dnat_rule.nat_gateway_id", respBody, nil)),
		d.Set("floating_ip_id", utils.PathSearch("dnat_rule.floating_ip_id", respBody, nil)),
		d.Set("global_eip_id", utils.PathSearch("dnat_rule.global_eip_id", respBody, nil)),
		d.Set("protocol", utils.PathSearch("dnat_rule.protocol", respBody, nil)),
		d.Set("internal_service_port", utils.PathSearch("dnat_rule.internal_service_port", respBody, nil)),
		d.Set("external_service_port", utils.PathSearch("dnat_rule.external_service_port", respBody, nil)),
		d.Set("internal_service_port_range", utils.PathSearch("dnat_rule.internal_service_port_range", respBody, nil)),
		d.Set("external_service_port_range", utils.PathSearch("dnat_rule.external_service_port_range", respBody, nil)),
		d.Set("description", utils.PathSearch("dnat_rule.description", respBody, nil)),
		d.Set("port_id", utils.PathSearch("dnat_rule.port_id", respBody, nil)),
		d.Set("private_ip", utils.PathSearch("dnat_rule.private_ip", respBody, nil)),
		d.Set("created_at", utils.PathSearch("dnat_rule.created_at", respBody, nil)),
		d.Set("floating_ip_address", utils.PathSearch("dnat_rule.floating_ip_address", respBody, nil)),
		d.Set("global_eip_address", utils.PathSearch("dnat_rule.global_eip_address", respBody, nil)),
		d.Set("status", utils.PathSearch("dnat_rule.status", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateDnatRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	dnatRuleBodyParams := map[string]interface{}{
		"nat_gateway_id":              d.Get("nat_gateway_id"),
		"protocol":                    d.Get("protocol"),
		"internal_service_port":       d.Get("internal_service_port"),
		"external_service_port":       d.Get("external_service_port"),
		"description":                 d.Get("description"),
		"floating_ip_id":              utils.ValueIgnoreEmpty(d.Get("floating_ip_id")),
		"global_eip_id":               utils.ValueIgnoreEmpty(d.Get("global_eip_id")),
		"internal_service_port_range": utils.ValueIgnoreEmpty(d.Get("internal_service_port_range")),
		"external_service_port_range": utils.ValueIgnoreEmpty(d.Get("external_service_port_range")),
		"port_id":                     utils.ValueIgnoreEmpty(d.Get("port_id")),
		"private_ip":                  utils.ValueIgnoreEmpty(d.Get("private_ip")),
	}

	return map[string]interface{}{
		"dnat_rule": dnatRuleBodyParams,
	}
}

func resourcePublicDnatRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		ruleId  = d.Id()
		httpUrl = "v2/{project_id}/dnat_rules/{dnat_rule_id}"
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{dnat_rule_id}", ruleId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateDnatRuleBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating DNAT rule: %s", err)
	}

	err = waitingForDnatRuleStateCompleted(ctx, client, d.Timeout(schema.TimeoutUpdate), ruleId, []string{"ACTIVE"})
	if err != nil {
		return diag.Errorf("error waiting for the DNAT rule (%s) update to complete: %s", ruleId, err)
	}

	return resourcePublicDnatRuleRead(ctx, d, meta)
}

func resourcePublicDnatRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		ruleId    = d.Id()
		gatewayId = d.Get("nat_gateway_id").(string)
		httpUrl   = "v2/{project_id}/nat_gateways/{nat_gateway_id}/dnat_rules/{dnat_rule_id}"
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{nat_gateway_id}", gatewayId)
	deletePath = strings.ReplaceAll(deletePath, "{dnat_rule_id}", ruleId)
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		// If the DNAT rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting DNAT rule")
	}

	err = waitingForDnatRuleStateCompleted(ctx, client, d.Timeout(schema.TimeoutCreate), ruleId, nil)
	if err != nil {
		return diag.Errorf("error waiting for the DNAT rule (%s) deletion to complete: %s", ruleId, err)
	}

	return nil
}

func waitingForDnatRuleStateCompleted(ctx context.Context, client *golangsdk.ServiceClient, t time.Duration,
	ruleId string, targets []string) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      waitingForDnatRuleStateRefresh(client, ruleId, targets),
		Timeout:      t,
		Delay:        3 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
