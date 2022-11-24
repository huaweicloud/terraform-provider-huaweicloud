// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product VPN
// ---------------------------------------------------------------

package vpn

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/jmespath/go-jmespath"
)

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
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9_-]+$`),
						"the input is invalid"),
					validation.StringLenBetween(0, 64),
				),
			},
			"gateway_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  `The VPN gateway ID.`,
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"gateway_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  `The VPN gateway IP ID.`,
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"vpn_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The connection type. The value can be **policy**, **static** or **bgp**.`,
				ValidateFunc: validation.StringInSlice([]string{
					"policy", "static", "bgp",
				}, false),
				DiffSuppressFunc: utils.SuppressCaseDiffs,
			},
			"customer_gateway_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  `The customer gateway ID.`,
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"peer_subnets": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `The customer subnets.`,
			},
			"psk": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  `The pre-shared key.`,
				ValidateFunc: validation.StringLenBetween(8, 128),
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
			"enterprise_project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  `The enterprise project ID.`,
				ValidateFunc: validation.StringLenBetween(1, 64),
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
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the VPN connection.`,
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
				ValidateFunc: validation.StringInSlice([]string{
					"sha1", "md5", "sha2-256", "sha2-384", "sha2-512",
				}, false),
			},
			"encryption_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The encryption algorithm, 3DES is less secure, please use them with caution.`,
				ValidateFunc: validation.StringInSlice([]string{
					"3des", "aes-128", "aes-192", "aes-256", "aes-128-gcm-16", "aes-256-gcm-16", "aes-128-gcm-128", "aes-256-gcm-128",
				}, false),
			},
			"pfs": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The DH key group used by PFS.`,
				ValidateFunc: validation.StringInSlice([]string{
					"group1", "group2", "group5", "group14", "group16", "group19", "group20", "group21",
				}, false),
			},
			"ike_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The IKE negotiation version.`,
				ValidateFunc: validation.StringInSlice([]string{
					"v1", "v2",
				}, false),
			},
			"lifetime_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  `The life cycle of SA in seconds, when the life cycle expires, IKE SA will be automatically updated.`,
				ValidateFunc: validation.IntBetween(60, 604800),
			},
			"local_id_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The local ID type.`,
				ValidateFunc: validation.StringInSlice([]string{
					"ip", "fqdn",
				}, false),
			},
			"local_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  `The local ID.`,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"peer_id_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The peer ID type.`,
				ValidateFunc: validation.StringInSlice([]string{
					"ip", "fqdn", "any",
				}, false),
			},
			"peer_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  `The peer ID.`,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"phase1_negotiation_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The negotiation mode, only works when the ike_version is v1.`,
				ValidateFunc: validation.StringInSlice([]string{
					"main", "aggressive",
				}, false),
			},
			"authentication_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The authentication method during IKE negotiation. Only **pre-share** supported for now.`,
				ValidateFunc: validation.StringInSlice([]string{
					"pre-share",
				}, false),
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
				ValidateFunc: validation.StringInSlice([]string{
					"sha1", "md5", "sha2-256", "sha2-384", "sha2-512",
				}, false),
			},
			"encryption_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The encryption algorithm, 3DES is less secure, please use them with caution.`,
				ValidateFunc: validation.StringInSlice([]string{
					"3des", "aes-128", "aes-192", "aes-256", "aes-128-gcm-16", "aes-256-gcm-16", "aes-128-gcm-128", "aes-256-gcm-128",
				}, false),
			},
			"pfs": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The DH key group used by PFS.`,
				ValidateFunc: validation.StringInSlice([]string{
					"group1", "group2", "group5", "group14", "group15", "group16", "group19", "group20", "group21",
				}, false),
			},
			"lifetime_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  `The lifecycle time of Ipsec tunnel in seconds.`,
				ValidateFunc: validation.IntBetween(30, 604800),
			},
			"transform_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The transform protocol. Only **esp** supported for now.`,
				ValidateFunc: validation.StringInSlice([]string{
					"esp",
				}, false),
			},
			"encapsulation_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The encapsulation mode, only **tunnel** supported for now.`,
				ValidateFunc: validation.StringInSlice([]string{
					"tunnel",
				}, false),
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
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	// createConnection: Create a VPN Connection.
	var (
		createConnectionHttpUrl = "v5/{project_id}/vpn-connection"
		createConnectionProduct = "vpn"
	)
	createConnectionClient, err := config.NewServiceClient(createConnectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating Connection Client: %s", err)
	}

	createConnectionPath := createConnectionClient.Endpoint + createConnectionHttpUrl
	createConnectionPath = strings.ReplaceAll(createConnectionPath, "{project_id}", createConnectionClient.ProjectID)

	createConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	createConnectionOpt.JSONBody = utils.RemoveNil(buildCreateConnectionBodyParams(d, config))
	createConnectionResp, err := createConnectionClient.Request("POST", createConnectionPath, &createConnectionOpt)
	if err != nil {
		return diag.Errorf("error creating Connection: %s", err)
	}

	createConnectionRespBody, err := utils.FlattenResponse(createConnectionResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("vpn_connection.id", createConnectionRespBody)
	if err != nil {
		return diag.Errorf("error creating Connection: ID is not found in API response")
	}
	d.SetId(id.(string))

	err = createConnectionWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Create of Connection (%s) to complete: %s", d.Id(), err)
	}
	return resourceConnectionRead(ctx, d, meta)
}

func buildCreateConnectionBodyParams(d *schema.ResourceData, config *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"vpn_connection": buildCreateConnectionVpnConnectionChildBody(d, config),
	}
	return bodyParams
}

func buildCreateConnectionVpnConnectionChildBody(d *schema.ResourceData, config *config.Config) map[string]interface{} {
	params := map[string]interface{}{
		"name":                  utils.ValueIngoreEmpty(d.Get("name")),
		"vgw_id":                utils.ValueIngoreEmpty(d.Get("gateway_id")),
		"vgw_ip":                utils.ValueIngoreEmpty(d.Get("gateway_ip")),
		"style":                 utils.ValueIngoreEmpty(d.Get("vpn_type")),
		"cgw_id":                utils.ValueIngoreEmpty(d.Get("customer_gateway_id")),
		"peer_subnets":          utils.ValueIngoreEmpty(d.Get("peer_subnets")),
		"psk":                   utils.ValueIngoreEmpty(d.Get("psk")),
		"tunnel_local_address":  utils.ValueIngoreEmpty(d.Get("tunnel_local_address")),
		"tunnel_peer_address":   utils.ValueIngoreEmpty(d.Get("tunnel_peer_address")),
		"enable_nqa":            utils.ValueIngoreEmpty(d.Get("enable_nqa")),
		"enterprise_project_id": utils.ValueIngoreEmpty(common.GetEnterpriseProjectID(d, config)),
		"ikepolicy":             buildCreateConnectionIkepolicyChildBody(d),
		"ipsecpolicy":           buildCreateConnectionIpsecpolicyChildBody(d),
		"policy_rules":          buildCreateConnectionPolicyRulesChildBody(d),
	}
	return params
}

func buildCreateConnectionIkepolicyChildBody(d *schema.ResourceData) map[string]interface{} {
	rawParams := d.Get("ikepolicy").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	raw := rawParams[0].(map[string]interface{})
	params := map[string]interface{}{
		"authentication_algorithm": utils.ValueIngoreEmpty(raw["authentication_algorithm"]),
		"encryption_algorithm":     utils.ValueIngoreEmpty(raw["encryption_algorithm"]),
		"pfs":                      utils.ValueIngoreEmpty(raw["pfs"]),
		"ike_version":              utils.ValueIngoreEmpty(raw["ike_version"]),
		"lifetime_seconds":         utils.ValueIngoreEmpty(raw["lifetime_seconds"]),
		"local_id_type":            utils.ValueIngoreEmpty(raw["local_id_type"]),
		"local_id":                 utils.ValueIngoreEmpty(raw["local_id"]),
		"peer_id_type":             utils.ValueIngoreEmpty(raw["peer_id_type"]),
		"peer_id":                  utils.ValueIngoreEmpty(raw["peer_id"]),
		"phase1_negotiation_mode":  utils.ValueIngoreEmpty(raw["phase1_negotiation_mode"]),
		"authentication_method":    utils.ValueIngoreEmpty(raw["authentication_method"]),
	}

	return params
}

func buildCreateConnectionIpsecpolicyChildBody(d *schema.ResourceData) map[string]interface{} {
	rawParams := d.Get("ipsecpolicy").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	raw := rawParams[0].(map[string]interface{})
	params := map[string]interface{}{
		"authentication_algorithm": utils.ValueIngoreEmpty(raw["authentication_algorithm"]),
		"encryption_algorithm":     utils.ValueIngoreEmpty(raw["encryption_algorithm"]),
		"pfs":                      utils.ValueIngoreEmpty(raw["pfs"]),
		"lifetime_seconds":         utils.ValueIngoreEmpty(raw["lifetime_seconds"]),
		"transform_protocol":       utils.ValueIngoreEmpty(raw["transform_protocol"]),
		"encapsulation_mode":       utils.ValueIngoreEmpty(raw["encapsulation_mode"]),
	}

	return params
}

func buildCreateConnectionPolicyRulesChildBody(d *schema.ResourceData) []map[string]interface{} {
	rawParams := d.Get("policy_rules").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	params := make([]map[string]interface{}, len(rawParams))
	for i, raw := range rawParams {
		rawMap := raw.(map[string]interface{})
		params[i] = map[string]interface{}{
			"rule_index":  utils.ValueIngoreEmpty(rawMap["rule_index"]),
			"source":      utils.ValueIngoreEmpty(rawMap["source"]),
			"destination": utils.ValueIngoreEmpty(rawMap["destination"]),
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
				return nil, "ERROR", fmt.Errorf("error creating Connection Client: %s", err)
			}

			createConnectionWaitingPath := createConnectionWaitingClient.Endpoint + createConnectionWaitingHttpUrl
			createConnectionWaitingPath = strings.ReplaceAll(createConnectionWaitingPath, "{project_id}", createConnectionWaitingClient.ProjectID)
			createConnectionWaitingPath = strings.ReplaceAll(createConnectionWaitingPath, "{id}", d.Id())

			createConnectionWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			createConnectionWaitingResp, err := createConnectionWaitingClient.Request("GET", createConnectionWaitingPath, &createConnectionWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createConnectionWaitingRespBody, err := utils.FlattenResponse(createConnectionWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw, err := jmespath.Search(`vpn_connection.status`, createConnectionWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `vpn_connection.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"ACTIVE",
				"DOWN",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return createConnectionWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"ERROR",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return createConnectionWaitingRespBody, status, nil
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

func resourceConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	var mErr *multierror.Error

	// getConnection: Query the VPN Connection detail
	var (
		getConnectionHttpUrl = "v5/{project_id}/vpn-connection/{id}"
		getConnectionProduct = "vpn"
	)
	getConnectionClient, err := config.NewServiceClient(getConnectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating Connection Client: %s", err)
	}

	getConnectionPath := getConnectionClient.Endpoint + getConnectionHttpUrl
	getConnectionPath = strings.ReplaceAll(getConnectionPath, "{project_id}", getConnectionClient.ProjectID)
	getConnectionPath = strings.ReplaceAll(getConnectionPath, "{id}", d.Id())

	getConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getConnectionResp, err := getConnectionClient.Request("GET", getConnectionPath, &getConnectionOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Connection")
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
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetConnectionResponseBodyCreateRequestIkePolicy(resp interface{}) []interface{} {
	var rst []interface{}
	curJson, err := jmespath.Search("vpn_connection.ikepolicy", resp)
	if err != nil {
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
		},
	}
	return rst
}

func flattenGetConnectionResponseBodyCreateRequestIpsecPolicy(resp interface{}) []interface{} {
	var rst []interface{}
	curJson, err := jmespath.Search("vpn_connection.ipsecpolicy", resp)
	if err != nil {
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
	config := meta.(*config.Config)
	region := config.GetRegion(d)

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
		var (
			updateConnectionHttpUrl = "v5/{project_id}/vpn-connection/{id}"
			updateConnectionProduct = "vpn"
		)
		updateConnectionClient, err := config.NewServiceClient(updateConnectionProduct, region)
		if err != nil {
			return diag.Errorf("error creating Connection Client: %s", err)
		}

		updateConnectionPath := updateConnectionClient.Endpoint + updateConnectionHttpUrl
		updateConnectionPath = strings.ReplaceAll(updateConnectionPath, "{project_id}", updateConnectionClient.ProjectID)
		updateConnectionPath = strings.ReplaceAll(updateConnectionPath, "{id}", d.Id())

		updateConnectionOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateConnectionOpt.JSONBody = utils.RemoveNil(buildUpdateConnectionBodyParams(d, config))
		_, err = updateConnectionClient.Request("PUT", updateConnectionPath, &updateConnectionOpt)
		if err != nil {
			return diag.Errorf("error updating Connection: %s", err)
		}
		err = updateConnectionWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the Update of Connection (%s) to complete: %s", d.Id(), err)
		}
	}
	return resourceConnectionRead(ctx, d, meta)
}

func buildUpdateConnectionBodyParams(d *schema.ResourceData, config *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"vpn_connection": buildUpdateConnectionVpnConnectionChildBody(d),
	}
	return bodyParams
}

func buildUpdateConnectionVpnConnectionChildBody(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"cgw_id":               utils.ValueIngoreEmpty(d.Get("customer_gateway_id")),
		"enable_nqa":           utils.ValueIngoreEmpty(d.Get("enable_nqa")),
		"name":                 utils.ValueIngoreEmpty(d.Get("name")),
		"peer_subnets":         utils.ValueIngoreEmpty(d.Get("peer_subnets")),
		"psk":                  utils.ValueIngoreEmpty(d.Get("psk")),
		"tunnel_local_address": utils.ValueIngoreEmpty(d.Get("tunnel_local_address")),
		"tunnel_peer_address":  utils.ValueIngoreEmpty(d.Get("tunnel_peer_address")),
		"ikepolicy":            buildUpdateConnectionIkepolicyChildBody(d),
		"ipsecpolicy":          buildUpdateConnectionIpsecpolicyChildBody(d),
		"policy_rules":         buildCreateConnectionPolicyRulesChildBody(d),
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
		"authentication_algorithm": utils.ValueIngoreEmpty(raw["authentication_algorithm"]),
		"authentication_method":    utils.ValueIngoreEmpty(raw["authentication_method"]),
		"encryption_algorithm":     utils.ValueIngoreEmpty(raw["encryption_algorithm"]),
		"ike_version":              utils.ValueIngoreEmpty(raw["ike_version"]),
		"lifetime_seconds":         utils.ValueIngoreEmpty(raw["lifetime_seconds"]),
		"local_id":                 utils.ValueIngoreEmpty(raw["local_id"]),
		"local_id_type":            utils.ValueIngoreEmpty(raw["local_id_type"]),
		"peer_id":                  utils.ValueIngoreEmpty(raw["peer_id"]),
		"peer_id_type":             utils.ValueIngoreEmpty(raw["peer_id_type"]),
		"pfs":                      utils.ValueIngoreEmpty(raw["pfs"]),
		"phase1_negotiation_mode":  utils.ValueIngoreEmpty(raw["phase1_negotiation_mode"]),
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
		"authentication_algorithm": utils.ValueIngoreEmpty(raw["authentication_algorithm"]),
		"encapsulation_mode":       utils.ValueIngoreEmpty(raw["encapsulation_mode"]),
		"encryption_algorithm":     utils.ValueIngoreEmpty(raw["encryption_algorithm"]),
		"lifetime_seconds":         utils.ValueIngoreEmpty(raw["lifetime_seconds"]),
		"pfs":                      utils.ValueIngoreEmpty(raw["pfs"]),
		"transform_protocol":       utils.ValueIngoreEmpty(raw["transform_protocol"]),
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
				return nil, "ERROR", fmt.Errorf("error creating Connection Client: %s", err)
			}

			updateConnectionWaitingPath := updateConnectionWaitingClient.Endpoint + updateConnectionWaitingHttpUrl
			updateConnectionWaitingPath = strings.ReplaceAll(updateConnectionWaitingPath, "{project_id}", updateConnectionWaitingClient.ProjectID)
			updateConnectionWaitingPath = strings.ReplaceAll(updateConnectionWaitingPath, "{id}", d.Id())

			updateConnectionWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			updateConnectionWaitingResp, err := updateConnectionWaitingClient.Request("GET", updateConnectionWaitingPath, &updateConnectionWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			updateConnectionWaitingRespBody, err := utils.FlattenResponse(updateConnectionWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw, err := jmespath.Search(`vpn_connection.status`, updateConnectionWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `vpn_connection.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"ACTIVE",
				"DOWN",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return updateConnectionWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"ERROR",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return updateConnectionWaitingRespBody, status, nil
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
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	// deleteConnection: Delete an existing VPN Connection
	var (
		deleteConnectionHttpUrl = "v5/{project_id}/vpn-connection/{id}"
		deleteConnectionProduct = "vpn"
	)
	deleteConnectionClient, err := config.NewServiceClient(deleteConnectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating Connection Client: %s", err)
	}

	deleteConnectionPath := deleteConnectionClient.Endpoint + deleteConnectionHttpUrl
	deleteConnectionPath = strings.ReplaceAll(deleteConnectionPath, "{project_id}", deleteConnectionClient.ProjectID)
	deleteConnectionPath = strings.ReplaceAll(deleteConnectionPath, "{id}", d.Id())

	deleteConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = deleteConnectionClient.Request("DELETE", deleteConnectionPath, &deleteConnectionOpt)
	if err != nil {
		return diag.Errorf("error deleting Connection: %s", err)
	}

	err = deleteConnectionWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the Delete of Connection (%s) to complete: %s", d.Id(), err)
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
				return nil, "ERROR", fmt.Errorf("error creating Connection Client: %s", err)
			}

			deleteConnectionWaitingPath := deleteConnectionWaitingClient.Endpoint + deleteConnectionWaitingHttpUrl
			deleteConnectionWaitingPath = strings.ReplaceAll(deleteConnectionWaitingPath, "{project_id}", deleteConnectionWaitingClient.ProjectID)
			deleteConnectionWaitingPath = strings.ReplaceAll(deleteConnectionWaitingPath, "{id}", d.Id())

			deleteConnectionWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			deleteConnectionWaitingResp, err := deleteConnectionWaitingClient.Request("GET", deleteConnectionWaitingPath, &deleteConnectionWaitingOpt)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return deleteConnectionWaitingResp, "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			deleteConnectionWaitingRespBody, err := utils.FlattenResponse(deleteConnectionWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw, err := jmespath.Search(`vpn_connection.status`, deleteConnectionWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `vpn_connection.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			unexpectedStatus := []string{
				"ERROR",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
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
