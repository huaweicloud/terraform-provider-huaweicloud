// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product VPN
// ---------------------------------------------------------------

package vpn

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

// @API VPN POST /v5/{project_id}/vpn-connection
// @API VPN DELETE /v5/{project_id}/vpn-connection/{id}
// @API VPN GET /v5/{project_id}/vpn-connection/{id}
// @API VPN PUT /v5/{project_id}/vpn-connection/{id}
// @API VPN POST /v5/{project_id}/{resource_type}/{resource_id}/tags/create
// @API VPN POST /v5/{project_id}/{resource_type}/{resource_id}/tags/delete
func ResourceConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectionCreate,
		UpdateContext: resourceConnectionUpdate,
		ReadContext:   resourceConnectionRead,
		DeleteContext: resourceConnectionDelete,
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
				Description: `The name of the VPN connection.`,
			},
			"gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The VPN gateway ID.`,
			},
			"gateway_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The VPN gateway IP ID.`,
			},
			"vpn_type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      `The connection type. The value can be **policy**, **static** or **bgp**.`,
				DiffSuppressFunc: utils.SuppressCaseDiffs(),
			},
			"customer_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The customer gateway ID.`,
			},
			"peer_subnets": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `The customer subnets.`,
			},
			"psk": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The pre-shared key.`,
			},
			"tunnel_local_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The local tunnel address.`,
			},
			"tunnel_peer_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The peer tunnel address.`,
			},
			"enable_nqa": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to enable NQA check.`,
			},
			"ikepolicy": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     ConnectionCreateRequestIkePolicySchema(),
				Optional: true,
				Computed: true,
			},
			"ipsecpolicy": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     ConnectionCreateRequestIpsecPolicySchema(),
				Optional: true,
				Computed: true,
			},
			"policy_rules": {
				Type:        schema.TypeList,
				Elem:        ConnectionPolicyRuleSchema(),
				Optional:    true,
				Computed:    true,
				Description: `The policy rules. Only works when vpn_type is set to **policy**`,
			},
			"tags": common.TagsSchema(),
			"ha_role": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the VPN connection.`,
			},
			// some early users may input this, leaving it optional, but the code won't handle it
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`The enterprise project ID.`,
					utils.SchemaDescInput{
						Computed: true,
					}),
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The create time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time.`,
			},
		},
	}
}

func ConnectionCreateRequestIkePolicySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"authentication_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The authentication algorithm, SHA1 and MD5 are less secure, please use them with caution.`,
			},
			"encryption_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The encryption algorithm, 3DES is less secure, please use them with caution.`,
			},
			"ike_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The IKE negotiation version.`,
			},
			"lifetime_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The life cycle of SA in seconds, when the life cycle expires, IKE SA will be automatically updated.`,
			},
			"local_id_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The local ID type.`,
			},
			"local_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The local ID.`,
			},
			"peer_id_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The peer ID type.`,
			},
			"peer_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The peer ID.`,
			},
			"phase1_negotiation_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The negotiation mode, only works when the ike_version is v1.`,
			},
			"authentication_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The authentication method during IKE negotiation.`,
			},
			"dh_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the DH group used for key exchange in phase 1.`,
			},
			"dpd": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        ConnectionPolicyDPDSchema(),
				Description: `Specifies the dead peer detection (DPD) object.`,
			},
			"pfs": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`The DH key group used by PFS.`,
					utils.SchemaDescInput{
						Deprecated: true,
					}),
			},
		},
	}
	return &sc
}

func ConnectionPolicyDPDSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the interval for retransmitting DPD packets.`,
			},
			"interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the DPD idle timeout period.`,
			},
			"msg": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the format of DPD packets.`,
			},
		},
	}
	return &sc
}

func ConnectionCreateRequestIpsecPolicySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"authentication_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The authentication algorithm, SHA1 and MD5 are less secure, please use them with caution.`,
			},
			"encryption_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The encryption algorithm, 3DES is less secure, please use them with caution.`,
			},
			"pfs": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The DH key group used by PFS.`,
			},
			"lifetime_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The lifecycle time of Ipsec tunnel in seconds.`,
			},
			"transform_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The transform protocol. Only **esp** supported for now.`,
			},
			"encapsulation_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The encapsulation mode, only **tunnel** supported for now.`,
			},
		},
	}
	return &sc
}

func ConnectionPolicyRuleSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"rule_index": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The rule index.`,
			},
			"source": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The source CIDR.`,
			},
			"destination": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `The list of destination CIDRs.`,
			},
		},
	}
	return &sc
}

func resourceConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	// createConnection: Create a VPN Connection.
	var (
		createConnectionHttpUrl = "v5/{project_id}/vpn-connection"
		createConnectionProduct = "vpn"
	)
	createConnectionClient, err := conf.NewServiceClient(createConnectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	createConnectionPath := createConnectionClient.Endpoint + createConnectionHttpUrl
	createConnectionPath = strings.ReplaceAll(createConnectionPath, "{project_id}", createConnectionClient.ProjectID)

	createConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createConnectionOpt.JSONBody = utils.RemoveNil(buildCreateConnectionBodyParams(d))
	createConnectionResp, err := createConnectionClient.Request("POST", createConnectionPath, &createConnectionOpt)
	if err != nil {
		return diag.Errorf("error creating VPN connection: %s", err)
	}

	createConnectionRespBody, err := utils.FlattenResponse(createConnectionResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("vpn_connection.id", createConnectionRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating VPN connection: ID is not found in API response")
	}
	d.SetId(id)

	err = createConnectionWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for creating VPN connection (%s) to complete: %s", d.Id(), err)
	}
	return resourceConnectionRead(ctx, d, meta)
}

func buildCreateConnectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"vpn_connection": buildCreateConnectionVpnConnectionChildBody(d),
	}
	return bodyParams
}

func buildCreateConnectionVpnConnectionChildBody(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"name":                 utils.ValueIgnoreEmpty(d.Get("name")),
		"vgw_id":               utils.ValueIgnoreEmpty(d.Get("gateway_id")),
		"vgw_ip":               utils.ValueIgnoreEmpty(d.Get("gateway_ip")),
		"style":                utils.ValueIgnoreEmpty(d.Get("vpn_type")),
		"cgw_id":               utils.ValueIgnoreEmpty(d.Get("customer_gateway_id")),
		"peer_subnets":         utils.ValueIgnoreEmpty(d.Get("peer_subnets")),
		"psk":                  utils.ValueIgnoreEmpty(d.Get("psk")),
		"tunnel_local_address": utils.ValueIgnoreEmpty(d.Get("tunnel_local_address")),
		"tunnel_peer_address":  utils.ValueIgnoreEmpty(d.Get("tunnel_peer_address")),
		"ikepolicy":            buildCreateConnectionIkepolicyChildBody(d),
		"ipsecpolicy":          buildCreateConnectionIpsecpolicyChildBody(d),
		"policy_rules":         buildCreateConnectionPolicyRulesChildBody(d),
		"tags":                 utils.ValueIgnoreEmpty(utils.ExpandResourceTags(d.Get("tags").(map[string]interface{}))),
		"ha_role":              utils.ValueIgnoreEmpty(d.Get("ha_role")),
	}

	if enableNqa, ok := d.GetOk("enable_nqa"); ok {
		params["enable_nqa"] = enableNqa
	}

	return params
}

func buildCreateConnectionIkepolicyChildBody(d *schema.ResourceData) map[string]interface{} {
	rawParams := d.Get("ikepolicy").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	if raw, ok := rawParams[0].(map[string]interface{}); ok {
		params := map[string]interface{}{
			"authentication_algorithm": utils.ValueIgnoreEmpty(raw["authentication_algorithm"]),
			"encryption_algorithm":     utils.ValueIgnoreEmpty(raw["encryption_algorithm"]),
			"pfs":                      utils.ValueIgnoreEmpty(raw["pfs"]),
			"ike_version":              utils.ValueIgnoreEmpty(raw["ike_version"]),
			"lifetime_seconds":         utils.ValueIgnoreEmpty(raw["lifetime_seconds"]),
			"local_id_type":            utils.ValueIgnoreEmpty(raw["local_id_type"]),
			"local_id":                 utils.ValueIgnoreEmpty(raw["local_id"]),
			"peer_id_type":             utils.ValueIgnoreEmpty(raw["peer_id_type"]),
			"peer_id":                  utils.ValueIgnoreEmpty(raw["peer_id"]),
			"phase1_negotiation_mode":  utils.ValueIgnoreEmpty(raw["phase1_negotiation_mode"]),
			"authentication_method":    utils.ValueIgnoreEmpty(raw["authentication_method"]),
			"dh_group":                 utils.ValueIgnoreEmpty(raw["dh_group"]),
			"dpd":                      buildCreateConnectionDPDChildBody(raw["dpd"]),
		}

		return params
	}

	return nil
}

func buildCreateConnectionDPDChildBody(dpd interface{}) map[string]interface{} {
	rawParams := dpd.([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	if raw, ok := rawParams[0].(map[string]interface{}); ok {
		params := map[string]interface{}{
			"timeout":  utils.ValueIgnoreEmpty(raw["timeout"]),
			"interval": utils.ValueIgnoreEmpty(raw["interval"]),
			"msg":      utils.ValueIgnoreEmpty(raw["msg"]),
		}

		return params
	}
	return nil
}

func buildCreateConnectionIpsecpolicyChildBody(d *schema.ResourceData) map[string]interface{} {
	rawParams := d.Get("ipsecpolicy").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	if raw, ok := rawParams[0].(map[string]interface{}); ok {
		params := map[string]interface{}{
			"authentication_algorithm": utils.ValueIgnoreEmpty(raw["authentication_algorithm"]),
			"encryption_algorithm":     utils.ValueIgnoreEmpty(raw["encryption_algorithm"]),
			"pfs":                      utils.ValueIgnoreEmpty(raw["pfs"]),
			"lifetime_seconds":         utils.ValueIgnoreEmpty(raw["lifetime_seconds"]),
			"transform_protocol":       utils.ValueIgnoreEmpty(raw["transform_protocol"]),
			"encapsulation_mode":       utils.ValueIgnoreEmpty(raw["encapsulation_mode"]),
		}

		return params
	}

	return nil
}

func buildCreateConnectionPolicyRulesChildBody(d *schema.ResourceData) []map[string]interface{} {
	rawParams := d.Get("policy_rules").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	params := make([]map[string]interface{}, len(rawParams))
	for i, raw := range rawParams {
		if rawMap, ok := raw.(map[string]interface{}); ok {
			params[i] = map[string]interface{}{
				"rule_index":  utils.ValueIgnoreEmpty(rawMap["rule_index"]),
				"source":      utils.ValueIgnoreEmpty(rawMap["source"]),
				"destination": utils.ValueIgnoreEmpty(rawMap["destination"]),
			}
		}
	}

	return params
}

func createConnectionWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			config := meta.(*config.Config)
			region := config.GetRegion(d)
			// createConnectionWaiting: missing operation notes
			var (
				createConnectionWaitingHttpUrl = "v5/{project_id}/vpn-connection/{id}"
				createConnectionWaitingProduct = "vpn"
			)
			createConnectionWaitingClient, err := config.NewServiceClient(createConnectionWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating VPN client: %s", err)
			}

			createConnectionWaitingPath := createConnectionWaitingClient.Endpoint + createConnectionWaitingHttpUrl
			createConnectionWaitingPath = strings.ReplaceAll(createConnectionWaitingPath, "{project_id}", createConnectionWaitingClient.ProjectID)
			createConnectionWaitingPath = strings.ReplaceAll(createConnectionWaitingPath, "{id}", d.Id())

			createConnectionWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
			}
			createConnectionWaitingResp, err := createConnectionWaitingClient.Request("GET", createConnectionWaitingPath, &createConnectionWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createConnectionWaitingRespBody, err := utils.FlattenResponse(createConnectionWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw := utils.PathSearch(`vpn_connection.status`, createConnectionWaitingRespBody, nil)
			if statusRaw == nil {
				return nil, "ERROR", fmt.Errorf("error parsing %s from response body", `vpn_connection.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			if status == "ERROR" {
				return createConnectionWaitingRespBody, status, nil
			}

			targetStatus := []string{
				"ACTIVE",
				"DOWN",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return createConnectionWaitingRespBody, "COMPLETED", nil
			}

			return createConnectionWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceConnectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var mErr *multierror.Error

	// getConnection: Query the VPN Connection detail
	var (
		getConnectionHttpUrl = "v5/{project_id}/vpn-connection/{id}"
		getConnectionProduct = "vpn"
	)
	getConnectionClient, err := conf.NewServiceClient(getConnectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	getConnectionPath := getConnectionClient.Endpoint + getConnectionHttpUrl
	getConnectionPath = strings.ReplaceAll(getConnectionPath, "{project_id}", getConnectionClient.ProjectID)
	getConnectionPath = strings.ReplaceAll(getConnectionPath, "{id}", d.Id())

	getConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getConnectionResp, err := getConnectionClient.Request("GET", getConnectionPath, &getConnectionOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VPN connection")
	}

	getConnectionRespBody, err := utils.FlattenResponse(getConnectionResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("vpn_connection.name", getConnectionRespBody, nil)),
		d.Set("gateway_id", utils.PathSearch("vpn_connection.vgw_id", getConnectionRespBody, nil)),
		d.Set("gateway_ip", utils.PathSearch("vpn_connection.vgw_ip", getConnectionRespBody, nil)),
		d.Set("vpn_type", utils.PathSearch("vpn_connection.style", getConnectionRespBody, nil)),
		d.Set("customer_gateway_id", utils.PathSearch("vpn_connection.cgw_id", getConnectionRespBody, nil)),
		d.Set("peer_subnets", utils.PathSearch("vpn_connection.peer_subnets", getConnectionRespBody, nil)),
		d.Set("tunnel_local_address", utils.PathSearch("vpn_connection.tunnel_local_address", getConnectionRespBody, nil)),
		d.Set("tunnel_peer_address", utils.PathSearch("vpn_connection.tunnel_peer_address", getConnectionRespBody, nil)),
		d.Set("enable_nqa", utils.PathSearch("vpn_connection.enable_nqa", getConnectionRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("vpn_connection.enterprise_project_id", getConnectionRespBody, nil)),
		d.Set("ikepolicy", flattenGetConnectionResponseBodyCreateRequestIkePolicy(getConnectionRespBody)),
		d.Set("ipsecpolicy", flattenGetConnectionResponseBodyCreateRequestIpsecPolicy(getConnectionRespBody)),
		d.Set("policy_rules", flattenGetConnectionResponseBodyPolicyRule(getConnectionRespBody)),
		d.Set("status", utils.PathSearch("vpn_connection.status", getConnectionRespBody, nil)),
		d.Set("created_at", utils.PathSearch("vpn_connection.created_at", getConnectionRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("vpn_connection.updated_at", getConnectionRespBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("vpn_connection.tags", getConnectionRespBody, nil))),
		d.Set("ha_role", utils.PathSearch("vpn_connection.ha_role", getConnectionRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetConnectionResponseBodyCreateRequestIkePolicy(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("vpn_connection.ikepolicy", resp, nil)
	if curJson == nil {
		log.Printf("[ERROR] error parsing vpn_connection.ikepolicy from response= %#v", resp)
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"authentication_algorithm": utils.PathSearch("authentication_algorithm", curJson, nil),
			"encryption_algorithm":     utils.PathSearch("encryption_algorithm", curJson, nil),
			"pfs":                      utils.PathSearch("pfs", curJson, nil),
			"ike_version":              utils.PathSearch("ike_version", curJson, nil),
			"lifetime_seconds":         utils.PathSearch("lifetime_seconds", curJson, nil),
			"local_id_type":            utils.PathSearch("local_id_type", curJson, nil),
			"local_id":                 utils.PathSearch("local_id", curJson, nil),
			"peer_id_type":             utils.PathSearch("peer_id_type", curJson, nil),
			"peer_id":                  utils.PathSearch("peer_id", curJson, nil),
			"phase1_negotiation_mode":  utils.PathSearch("phase1_negotiation_mode", curJson, nil),
			"authentication_method":    utils.PathSearch("authentication_method", curJson, nil),
			"dh_group":                 utils.PathSearch("dh_group", curJson, nil),
			"dpd":                      flattenGetConnectionResponseBodyDPD(resp),
		},
	}
	return rst
}

func flattenGetConnectionResponseBodyDPD(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("vpn_connection.ikepolicy.dpd", resp, nil)
	if curJson == nil {
		log.Printf("[ERROR] error parsing vpn_connection.ikepolicy.dpd from response= %#v", resp)
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"timeout":  utils.PathSearch("timeout", curJson, nil),
			"interval": utils.PathSearch("interval", curJson, nil),
			"msg":      utils.PathSearch("msg", curJson, nil),
		},
	}
	return rst
}

func flattenGetConnectionResponseBodyCreateRequestIpsecPolicy(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("vpn_connection.ipsecpolicy", resp, nil)
	if curJson == nil {
		log.Printf("[ERROR] error parsing vpn_connection.ipsecpolicy from response= %#v", resp)
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"authentication_algorithm": utils.PathSearch("authentication_algorithm", curJson, nil),
			"encryption_algorithm":     utils.PathSearch("encryption_algorithm", curJson, nil),
			"pfs":                      utils.PathSearch("pfs", curJson, nil),
			"lifetime_seconds":         utils.PathSearch("lifetime_seconds", curJson, nil),
			"transform_protocol":       utils.PathSearch("transform_protocol", curJson, nil),
			"encapsulation_mode":       utils.PathSearch("encapsulation_mode", curJson, nil),
		},
	}
	return rst
}

func flattenGetConnectionResponseBodyPolicyRule(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("vpn_connection.policy_rules", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"rule_index":  utils.PathSearch("rule_index", v, nil),
			"destination": utils.PathSearch("destination", v, nil),
			"source":      utils.PathSearch("source", v, nil),
		})
	}
	return rst
}

func resourceConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	updateConnectionClient, err := conf.NewServiceClient("vpn", region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	updateConnectionhasChanges := []string{
		"customer_gateway_id",
		"enable_nqa",
		"ikepolicy",
		"ipsecpolicy",
		"name",
		"peer_subnets",
		"policy_rules",
		"psk",
		"tunnel_local_address",
		"tunnel_peer_address",
	}

	if d.HasChanges(updateConnectionhasChanges...) {
		// updateConnection: Update the configuration of VPN Connection
		updateConnectionHttpUrl := "v5/{project_id}/vpn-connection/{id}"

		updateConnectionPath := updateConnectionClient.Endpoint + updateConnectionHttpUrl
		updateConnectionPath = strings.ReplaceAll(updateConnectionPath, "{project_id}", updateConnectionClient.ProjectID)
		updateConnectionPath = strings.ReplaceAll(updateConnectionPath, "{id}", d.Id())

		updateConnectionOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		updateConnectionOpt.JSONBody = utils.RemoveNil(buildUpdateConnectionBodyParams(d))
		_, err = updateConnectionClient.Request("PUT", updateConnectionPath, &updateConnectionOpt)
		if err != nil {
			return diag.Errorf("error updating VPN connection: %s", err)
		}
		err = updateConnectionWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for updating VPN connection (%s) to complete: %s", d.Id(), err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		tagErr := updateTags(updateConnectionClient, d, "vpn-connection", d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of VPN connection (%s): %s", d.Id(), tagErr)
		}
	}
	return resourceConnectionRead(ctx, d, meta)
}

func buildUpdateConnectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"vpn_connection": buildUpdateConnectionVpnConnectionChildBody(d),
	}
	return bodyParams
}

func buildUpdateConnectionVpnConnectionChildBody(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"cgw_id":               utils.ValueIgnoreEmpty(d.Get("customer_gateway_id")),
		"name":                 utils.ValueIgnoreEmpty(d.Get("name")),
		"peer_subnets":         utils.ValueIgnoreEmpty(d.Get("peer_subnets")),
		"psk":                  utils.ValueIgnoreEmpty(d.Get("psk")),
		"tunnel_local_address": utils.ValueIgnoreEmpty(d.Get("tunnel_local_address")),
		"tunnel_peer_address":  utils.ValueIgnoreEmpty(d.Get("tunnel_peer_address")),
		"ikepolicy":            buildUpdateConnectionIkepolicyChildBody(d),
		"ipsecpolicy":          buildUpdateConnectionIpsecpolicyChildBody(d),
		"policy_rules":         buildCreateConnectionPolicyRulesChildBody(d),
	}

	if enableNqa, ok := d.GetOk("enable_nqa"); ok {
		params["enable_nqa"] = enableNqa
	}

	return params
}

func buildUpdateConnectionIkepolicyChildBody(d *schema.ResourceData) map[string]interface{} {
	rawParams := d.Get("ikepolicy").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	raw := rawParams[0].(map[string]interface{})
	params := map[string]interface{}{
		"authentication_algorithm": utils.ValueIgnoreEmpty(raw["authentication_algorithm"]),
		"encryption_algorithm":     utils.ValueIgnoreEmpty(raw["encryption_algorithm"]),
		"ike_version":              utils.ValueIgnoreEmpty(raw["ike_version"]),
		"lifetime_seconds":         utils.ValueIgnoreEmpty(raw["lifetime_seconds"]),
		"local_id_type":            utils.ValueIgnoreEmpty(raw["local_id_type"]),
		"peer_id_type":             utils.ValueIgnoreEmpty(raw["peer_id_type"]),
		"pfs":                      utils.ValueIgnoreEmpty(raw["pfs"]),
		"phase1_negotiation_mode":  utils.ValueIgnoreEmpty(raw["phase1_negotiation_mode"]),
		"dh_group":                 utils.ValueIgnoreEmpty(raw["dh_group"]),
		"dpd":                      buildCreateConnectionDPDChildBody(raw["dpd"]),
	}

	// if the id type is ip, the id must be empty
	if raw["local_id_type"] != "ip" {
		params["local_id"] = utils.ValueIgnoreEmpty(raw["local_id"])
	}

	if raw["peer_id_type"] != "ip" {
		params["peer_id"] = utils.ValueIgnoreEmpty(raw["peer_id"])
	}

	return params
}

func buildUpdateConnectionIpsecpolicyChildBody(d *schema.ResourceData) map[string]interface{} {
	rawParams := d.Get("ipsecpolicy").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	raw := rawParams[0].(map[string]interface{})
	params := map[string]interface{}{
		"authentication_algorithm": utils.ValueIgnoreEmpty(raw["authentication_algorithm"]),
		"encapsulation_mode":       utils.ValueIgnoreEmpty(raw["encapsulation_mode"]),
		"encryption_algorithm":     utils.ValueIgnoreEmpty(raw["encryption_algorithm"]),
		"lifetime_seconds":         utils.ValueIgnoreEmpty(raw["lifetime_seconds"]),
		"pfs":                      utils.ValueIgnoreEmpty(raw["pfs"]),
		"transform_protocol":       utils.ValueIgnoreEmpty(raw["transform_protocol"]),
	}

	return params
}

func updateConnectionWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			config := meta.(*config.Config)
			region := config.GetRegion(d)
			// updateConnectionWaiting: missing operation notes
			var (
				updateConnectionWaitingHttpUrl = "v5/{project_id}/vpn-connection/{id}"
				updateConnectionWaitingProduct = "vpn"
			)
			updateConnectionWaitingClient, err := config.NewServiceClient(updateConnectionWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating VPN client: %s", err)
			}

			updateConnectionWaitingPath := updateConnectionWaitingClient.Endpoint + updateConnectionWaitingHttpUrl
			updateConnectionWaitingPath = strings.ReplaceAll(updateConnectionWaitingPath, "{project_id}", updateConnectionWaitingClient.ProjectID)
			updateConnectionWaitingPath = strings.ReplaceAll(updateConnectionWaitingPath, "{id}", d.Id())

			updateConnectionWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
			}
			updateConnectionWaitingResp, err := updateConnectionWaitingClient.Request("GET", updateConnectionWaitingPath, &updateConnectionWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			updateConnectionWaitingRespBody, err := utils.FlattenResponse(updateConnectionWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw := utils.PathSearch(`vpn_connection.status`, updateConnectionWaitingRespBody, nil)
			if statusRaw == nil {
				return nil, "ERROR", fmt.Errorf("error parsing %s from response body", `vpn_connection.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			if status == "ERROR" {
				return updateConnectionWaitingRespBody, status, nil
			}

			targetStatus := []string{
				"ACTIVE",
				"DOWN",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return updateConnectionWaitingRespBody, "COMPLETED", nil
			}

			return updateConnectionWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceConnectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	// deleteConnection: Delete an existing VPN Connection
	var (
		deleteConnectionHttpUrl = "v5/{project_id}/vpn-connection/{id}"
		deleteConnectionProduct = "vpn"
	)
	deleteConnectionClient, err := conf.NewServiceClient(deleteConnectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	deleteConnectionPath := deleteConnectionClient.Endpoint + deleteConnectionHttpUrl
	deleteConnectionPath = strings.ReplaceAll(deleteConnectionPath, "{project_id}", deleteConnectionClient.ProjectID)
	deleteConnectionPath = strings.ReplaceAll(deleteConnectionPath, "{id}", d.Id())

	deleteConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteConnectionClient.Request("DELETE", deleteConnectionPath, &deleteConnectionOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting VPN connection")
	}

	err = deleteConnectionWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for deleting VPN connection (%s) to complete: %s", d.Id(), err)
	}
	return nil
}

func deleteConnectionWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			config := meta.(*config.Config)
			region := config.GetRegion(d)
			// deleteConnectionWaiting: missing operation notes
			var (
				deleteConnectionWaitingHttpUrl = "v5/{project_id}/vpn-connection/{id}"
				deleteConnectionWaitingProduct = "vpn"
			)
			deleteConnectionWaitingClient, err := config.NewServiceClient(deleteConnectionWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating VPN client: %s", err)
			}

			deleteConnectionWaitingPath := deleteConnectionWaitingClient.Endpoint + deleteConnectionWaitingHttpUrl
			deleteConnectionWaitingPath = strings.ReplaceAll(deleteConnectionWaitingPath, "{project_id}", deleteConnectionWaitingClient.ProjectID)
			deleteConnectionWaitingPath = strings.ReplaceAll(deleteConnectionWaitingPath, "{id}", d.Id())

			deleteConnectionWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
			}
			deleteConnectionWaitingResp, err := deleteConnectionWaitingClient.Request("GET", deleteConnectionWaitingPath, &deleteConnectionWaitingOpt)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
					return "Resource Not Found", "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			deleteConnectionWaitingRespBody, err := utils.FlattenResponse(deleteConnectionWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw := utils.PathSearch(`vpn_connection.status`, deleteConnectionWaitingRespBody, nil)
			if statusRaw == nil {
				return nil, "ERROR", fmt.Errorf("error parsing %s from response body", `vpn_connection.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			if status == "ERROR" {
				return deleteConnectionWaitingRespBody, status, nil
			}

			return deleteConnectionWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
