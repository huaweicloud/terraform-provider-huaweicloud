package nat

import (
	"context"
	"fmt"
	"log"
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

// @API NAT POST /v2/{project_id}/snat_rules
// @API NAT GET /v2/{project_id}/snat_rules/{snat_rule_id}
// @API NAT PUT /v2/{project_id}/snat_rules/{snat_rule_id}
// @API NAT DELETE /v2/{project_id}/nat_gateways/{nat_gateway_id}/snat_rules/{snat_rule_id}
// @API EIP GET /v1/{project_id}/publicips/{publicip_id}
func ResourcePublicSnatRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePublicSnatRuleCreate,
		ReadContext:   resourcePublicSnatRuleRead,
		UpdateContext: resourcePublicSnatRuleUpdate,
		DeleteContext: resourcePublicSnatRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the SNAT rule is located.",
			},
			"floating_ip_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ExactlyOneOf:     []string{"floating_ip_id", "global_eip_id"},
				DiffSuppressFunc: utils.SuppressSnatFiplistDiffs,
				Description:      "The IDs of floating IPs connected by SNAT rule.",
			},
			"global_eip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The IDs (separated by commas) of global EIPs connected by SNAT rule.",
			},
			"nat_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "schema: Required; The ID of the gateway to which the SNAT rule belongs.",
			},
			"source_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The resource type of the SNAT rule.",
			},
			"subnet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"cidr", "network_id"},
				Description:  "The network IDs of subnet connected by SNAT rule (VPC side).",
			},
			"cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"subnet_id", "network_id"},
				Description:  "The CIDR block connected by SNAT rule (DC side).",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the SNAT rule.",
			},
			"floating_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The floating IP addresses (separated by commas) connected by SNAT rule.",
			},
			"global_eip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The global EIP addresses (separated by commas) connected by SNAT rule.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the SNAT rule.",
			},
			"freezed_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The frozen EIP associated with the SNAT rule.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the SNAT rule.",
			},

			// deprecated
			"network_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "schema: Deprecated; Use 'subnet_id' instead.",
			},
		},
	}
}

func buildCreateSnatRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	snatRuleBodyParams := map[string]interface{}{
		"nat_gateway_id": d.Get("nat_gateway_id"),
		"floating_ip_id": utils.ValueIgnoreEmpty(d.Get("floating_ip_id")),
		"global_eip_id":  utils.ValueIgnoreEmpty(d.Get("global_eip_id")),
		"network_id":     utils.ValueIgnoreEmpty(buildSubnetId(d.Get("subnet_id").(string), d.Get("network_id").(string))),
		"cidr":           utils.ValueIgnoreEmpty(d.Get("cidr")),
		"description":    utils.ValueIgnoreEmpty(d.Get("description")),
	}

	subnetId := d.Get("subnet_id").(string)
	sourceType := d.Get("source_type").(int)
	if sourceType == 1 && subnetId != "" {
		log.Printf("[WARN] in the DC (Direct Connect) scenario (source_type is 1), only the parameter 'cidr' " +
			"is valid, and the parameter 'subnet_id' must be empty")
	}

	snatRuleBodyParams["source_type"] = utils.ValueIgnoreEmpty(sourceType)

	return map[string]interface{}{
		"snat_rule": snatRuleBodyParams,
	}
}

func buildSubnetId(subnetId, networkId string) string {
	if subnetId != "" {
		return subnetId
	}

	return networkId
}

func waitingForSnatRuleStateRefresh(client *golangsdk.ServiceClient, ruleId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := GetSnatRule(client, ruleId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "Resource Not Found", "COMPLETED", nil
			}

			return resp, "", err
		}

		status := utils.PathSearch("snat_rule.status", resp, "").(string)

		if utils.StrSliceContains([]string{"INACTIVE", "EIP_FREEZED"}, status) {
			return resp, "", fmt.Errorf("unexpect status (%s)", status)
		}

		if utils.StrSliceContains(targets, status) {
			return resp, "COMPLETED", nil
		}

		return resp, "PENDING", nil
	}
}

func resourcePublicSnatRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/snat_rules"
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateSnatRuleBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SNAT rule: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleId := utils.PathSearch("snat_rule.id", respBody, "").(string)
	if ruleId == "" {
		return diag.Errorf("error creating SNAT rule: ID is not found in API response")
	}

	d.SetId(ruleId)

	err = waitingForSnatRuleStateCompleted(ctx, client, d.Timeout(schema.TimeoutCreate), d.Id(), []string{"ACTIVE"})
	if err != nil {
		return diag.Errorf("error waiting for the SNAT rule (%s) creation to complete: %s", d.Id(), err)
	}

	return resourcePublicSnatRuleRead(ctx, d, meta)
}

func GetSnatRule(client *golangsdk.ServiceClient, ruleId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/snat_rules/{snat_rule_id}"
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

func resourcePublicSnatRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	respBody, err := GetSnatRule(client, d.Id())
	if err != nil {
		// If the SNAT rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving SNAT rule")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("nat_gateway_id", utils.PathSearch("snat_rule.nat_gateway_id", respBody, nil)),
		d.Set("floating_ip_id", utils.PathSearch("snat_rule.floating_ip_id", respBody, nil)),
		d.Set("floating_ip_address", utils.PathSearch("snat_rule.floating_ip_address", respBody, nil)),
		d.Set("global_eip_id", utils.PathSearch("snat_rule.global_eip_id", respBody, nil)),
		d.Set("global_eip_address", utils.PathSearch("snat_rule.global_eip_address", respBody, nil)),
		d.Set("source_type", utils.PathSearch("snat_rule.source_type", respBody, nil)),
		d.Set("subnet_id", utils.PathSearch("snat_rule.network_id", respBody, nil)),
		d.Set("cidr", utils.PathSearch("snat_rule.cidr", respBody, nil)),
		d.Set("status", utils.PathSearch("snat_rule.status", respBody, nil)),
		d.Set("description", utils.PathSearch("snat_rule.description", respBody, nil)),
		d.Set("freezed_ip_address", utils.PathSearch("snat_rule.freezed_ip_address", respBody, nil)),
		d.Set("created_at", utils.PathSearch("snat_rule.created_at", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateSnatRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	dnatRuleBodyParams := map[string]interface{}{
		"nat_gateway_id": d.Get("nat_gateway_id"),
		"description":    d.Get("description"),
	}

	return dnatRuleBodyParams
}

func updateSnatRule(client *golangsdk.ServiceClient, ruleId string, snatRuleBodyParams map[string]interface{}) error {
	httpUrl := "v2/{project_id}/snat_rules/{snat_rule_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{snat_rule_id}", ruleId)
	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"snat_rule": snatRuleBodyParams,
		},
	}

	_, err := client.Request("PUT", updatePath, &opts)
	return err
}

func getEipAddress(client *golangsdk.ServiceClient, eipId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/publicips/{publicip_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{publicip_id}", eipId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func resourcePublicSnatRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		ruleId = d.Id()
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	updateOpts := buildUpdateSnatRuleBodyParams(d)

	if d.HasChange("floating_ip_id") {
		eipClient, err := cfg.NewServiceClient("vpc", region)
		if err != nil {
			return diag.Errorf("error creating VPC v1 client: %s", err)
		}

		eipIds := d.Get("floating_ip_id").(string)
		eipList := strings.Split(eipIds, ",")
		eipAddrs := make([]string, len(eipList))

		// get EIP address from ID
		for i, eipId := range eipList {
			eipResp, err := getEipAddress(eipClient, eipId)
			if err != nil {
				return diag.Errorf("error fetching EIP (%s): %s", eipId, err)
			}

			eipAddrs[i] = utils.PathSearch("publicip.public_ip_address", eipResp, "").(string)
		}

		updateOpts["public_ip_address"] = strings.Join(eipAddrs, ",")
	}

	if d.HasChange("global_eip_id") {
		updateOpts["global_eip_id"] = utils.ValueIgnoreEmpty(d.Get("global_eip_id"))
	}

	err = updateSnatRule(client, d.Id(), updateOpts)
	if err != nil {
		return diag.Errorf("error updating SNAT rule: %s", err)
	}

	err = waitingForSnatRuleStateCompleted(ctx, client, d.Timeout(schema.TimeoutUpdate), ruleId, []string{"ACTIVE"})
	if err != nil {
		return diag.Errorf("error waiting for the SNAT rule (%s) update to complete: %s", ruleId, err)
	}

	return resourcePublicSnatRuleRead(ctx, d, meta)
}

func resourcePublicSnatRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		gatewayId = d.Get("nat_gateway_id").(string)
		ruleId    = d.Id()
		httpUrl   = "v2/{project_id}/nat_gateways/{nat_gateway_id}/snat_rules/{snat_rule_id}"
	)

	client, err := cfg.NewServiceClient("nat", region)
	if err != nil {
		return diag.Errorf("error creating NAT v2 client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{nat_gateway_id}", gatewayId)
	deletePath = strings.ReplaceAll(deletePath, "{snat_rule_id}", ruleId)
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		// If the SNAT rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting SNAT rule")
	}

	err = waitingForSnatRuleStateCompleted(ctx, client, d.Timeout(schema.TimeoutDelete), ruleId, nil)
	if err != nil {
		return diag.Errorf("error waiting for the SNAT rule (%s) deletion to complete: %s", ruleId, err)
	}

	return nil
}

func waitingForSnatRuleStateCompleted(ctx context.Context, client *golangsdk.ServiceClient, t time.Duration,
	ruleId string, targets []string) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      waitingForSnatRuleStateRefresh(client, ruleId, targets),
		Timeout:      t,
		Delay:        3 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
