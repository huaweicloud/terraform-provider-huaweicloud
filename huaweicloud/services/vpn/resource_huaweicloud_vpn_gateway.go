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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPN POST /v5/{project_id}/vpn-gateways
// @API VPN POST /v5/{project_id}/vpn-gateways/{gateway_id}/certificate
// @API VPN PUT /v5/{project_id}/vpn-gateways/{id}
// @API VPN PUT /v5/{project_id}/vpn-gateways/{gateway_id}/certificate/{certificate_id}
// @API VPN GET /v5/{project_id}/vpn-gateways/{id}
// @API VPN GET /v5/{project_id}/vpn-gateways/{gateway_id}/certificate
// @API VPN DELETE /v5/{project_id}/vpn-gateways/{id}
// @API VPN POST /v5/{project_id}/{resource_type}/{resource_id}/tags/create
// @API VPN POST /v5/{project_id}/{resource_type}/{resource_id}/tags/delete
// @API VPN POST /v5/{project_id}/vpn-gateways/{vgw_id}/update-specification
// @API EIP POST /v3/{project_id}/eip/publicips/{publicip_id}/disassociate-instance
func ResourceGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGatewayCreate,
		UpdateContext: resourceGatewayUpdate,
		ReadContext:   resourceGatewayRead,
		DeleteContext: resourceGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

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
				Description: `The name of the VPN gateway. Only letters, digits, underscores(_) and hypens(-) are supported.`,
			},
			"availability_zones": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				ForceNew:    true,
				Description: `The availability zone IDs.`,
			},
			"flavor": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The flavor of the VPN gateway.`,
			},
			"attachment_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "vpc",
				ForceNew:    true,
				Description: `The attachment type.`,
				ValidateFunc: validation.StringInSlice([]string{
					"vpc", "er",
				}, false),
			},
			"network_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The network type of the VPN gateway.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the VPC to which the VPN gateway is connected.`,
			},
			"local_subnets": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `The local subnets.`,
			},
			"connect_subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The Network ID of the VPC subnet used by the VPN gateway.`,
			},
			"er_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The enterprise router ID to attach with to VPN gateway.`,
			},
			"ha_mode": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				Description:   `The HA mode of the VPN gateway.`,
				ValidateFunc:  validation.StringInSlice([]string{"active-active", "active-standby"}, false),
				ConflictsWith: []string{"master_eip", "slave_eip"},
			},
			"master_eip": {
				Type:         schema.TypeList,
				MaxItems:     1,
				Elem:         GatewayEipSchema(),
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
				RequiredWith: []string{"slave_eip"},
			},
			"eip1": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Elem:          GatewayEipSchema(),
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"master_eip", "slave_eip"},
				RequiredWith:  []string{"eip2"},
			},
			"slave_eip": {
				Type:         schema.TypeList,
				MaxItems:     1,
				Elem:         GatewayEipSchema(),
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
				RequiredWith: []string{"master_eip"},
			},
			"eip2": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Elem:          GatewayEipSchema(),
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"master_eip", "slave_eip"},
				RequiredWith:  []string{"eip1"},
			},
			"access_vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The access VPC ID of the VPN gateway.`,
			},
			"access_subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The access subnet ID of the VPN gateway.`,
			},
			"asn": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     64512,
				ForceNew:    true,
				Description: `The ASN number of BGP`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The enterprise project ID`,
			},
			"access_private_ip_1": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				RequiredWith: []string{"access_private_ip_2"},
			},
			"access_private_ip_2": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				RequiredWith: []string{"access_private_ip_1"},
			},
			"certificate": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     gatewayCertificateSchema(),
			},
			"tags": common.TagsSchema(),
			"delete_eip_on_termination": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Whether to delete the EIP when the VPN gateway is deleted.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of VPN gateway.`,
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
			"er_attachment_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ER attachment ID.`,
			},
			"used_connection_group": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of used connection groups.`,
			},
			"used_connection_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of used connections.`,
			},
		},
	}
}

func gatewayCertificateSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"certificate_chain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enc_certificate": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enc_private_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"issuer": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signature_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_serial_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_subject": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_chain_serial_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_chain_subject": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_chain_expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enc_certificate_serial_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enc_certificate_subject": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enc_certificate_expire_time": {
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

func GatewayEipSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The public IP ID.`,
			},

			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The EIP type. The value can be **5_bgp** and **5_sbgp**.`,
			},
			"bandwidth_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The bandwidth name.`,
			},
			"bandwidth_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Description: `Bandwidth size in Mbit/s. When the flavor is **V300**, the value cannot be greater than **300**.
`,
			},
			"charge_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The charge mode of the bandwidth. The value can be **bandwidth** and **traffic**.`,
				ValidateFunc: validation.StringInSlice([]string{
					"bandwidth", "traffic",
				}, false),
			},

			"bandwidth_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The bandwidth ID.`,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The public IP address.`,
			},
			"ip_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The public IP version.`,
			},
		},
	}
	return &sc
}

func resourceGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createGateway: Create a VPN Gateway.
	var (
		createGatewayHttpUrl            = "v5/{project_id}/vpn-gateways"
		createGatewayCertificateHttpUrl = "v5/{project_id}/vpn-gateways/{gateway_id}/certificate"
		createGatewayProduct            = "vpn"
	)
	createGatewayClient, err := cfg.NewServiceClient(createGatewayProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	createGatewayPath := createGatewayClient.Endpoint + createGatewayHttpUrl
	createGatewayPath = strings.ReplaceAll(createGatewayPath, "{project_id}", createGatewayClient.ProjectID)

	createGatewayOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createGatewayOpt.JSONBody = utils.RemoveNil(buildCreateGatewayBodyParams(d, cfg))
	createGatewayResp, err := createGatewayClient.Request("POST", createGatewayPath, &createGatewayOpt)
	if err != nil {
		return diag.Errorf("error creating gateway: %s", err)
	}

	createGatewayRespBody, err := utils.FlattenResponse(createGatewayResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("vpn_gateway.id", createGatewayRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating VPN gateway: ID is not found in API response")
	}
	d.SetId(id)

	err = createGatewayWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for creating VPN gateway (%s) to complete: %s", d.Id(), err)
	}

	certificateContent := d.Get("certificate").([]interface{})
	if len(certificateContent) == 1 {
		createGatewayCertificatePath := createGatewayClient.Endpoint + createGatewayCertificateHttpUrl
		createGatewayCertificatePath = strings.ReplaceAll(createGatewayCertificatePath, "{project_id}", createGatewayClient.ProjectID)
		createGatewayCertificatePath = strings.ReplaceAll(createGatewayCertificatePath, "{gateway_id}", d.Id())
		createGatewayCertificateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		createGatewayCertificateOpt.JSONBody = utils.RemoveNil(buildCreateGatewayCertificateBodyParams(certificateContent[0]))
		createGatewayCertificateResp, err := createGatewayClient.Request("POST", createGatewayCertificatePath, &createGatewayCertificateOpt)
		if err != nil {
			return diag.Errorf("error creating VPN gateway certificate: %s", err)
		}
		createGatewayCertificateRespBody, err := utils.FlattenResponse(createGatewayCertificateResp)
		if err != nil {
			return diag.FromErr(err)
		}

		mErr := multierror.Append(nil, d.Set("certificate",
			flattenGatewayCertificateResponse(d, createGatewayCertificateRespBody)))
		if mErr.ErrorOrNil() != nil {
			return diag.FromErr(mErr.ErrorOrNil())
		}

		err = waitingForGatewayCertificateStateCompleted(ctx, d, createGatewayClient, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("error waiting for creating VPN gateway (%s) certificate to complete: %s", d.Id(), err)
		}
	}
	return resourceGatewayRead(ctx, d, meta)
}

func buildCreateGatewayCertificateBodyParams(certificateContent interface{}) map[string]interface{} {
	certificateMap := certificateContent.(map[string]interface{})
	bodyParams := map[string]interface{}{
		"certificate": map[string]interface{}{
			"name":              utils.ValueIgnoreEmpty(certificateMap["name"].(string)),
			"certificate":       utils.ValueIgnoreEmpty(certificateMap["content"].(string)),
			"private_key":       utils.ValueIgnoreEmpty(certificateMap["private_key"].(string)),
			"certificate_chain": utils.ValueIgnoreEmpty(certificateMap["certificate_chain"].(string)),
			"enc_certificate":   utils.ValueIgnoreEmpty(certificateMap["enc_certificate"].(string)),
			"enc_private_key":   utils.ValueIgnoreEmpty(certificateMap["enc_private_key"].(string)),
		},
	}
	return bodyParams
}

func buildCreateGatewayBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"vpn_gateway": buildCreateGatewayVpnGatewayChildBody(d, cfg),
	}
	return bodyParams
}

func buildCreateGatewayVpnGatewayChildBody(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	haMode := d.Get("ha_mode").(string)
	masterEIP := buildCreateGatewayEIPChildBody(d, "master_eip")
	slaveEIP := buildCreateGatewayEIPChildBody(d, "slave_eip")

	// default use "active-standby" ha_mode type when declare master_eip and slave_eip
	if haMode == "" && masterEIP != nil && slaveEIP != nil {
		haMode = "active-standby"
	}
	params := map[string]interface{}{
		"attachment_type":       utils.ValueIgnoreEmpty(d.Get("attachment_type")),
		"availability_zone_ids": utils.ValueIgnoreEmpty(d.Get("availability_zones").(*schema.Set).List()),
		"bgp_asn":               utils.ValueIgnoreEmpty(d.Get("asn")),
		"connect_subnet":        utils.ValueIgnoreEmpty(d.Get("connect_subnet")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		"flavor":                utils.ValueIgnoreEmpty(d.Get("flavor")),
		"local_subnets":         utils.ValueIgnoreEmpty(d.Get("local_subnets")),
		"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
		"vpc_id":                utils.ValueIgnoreEmpty(d.Get("vpc_id")),
		"ha_mode":               utils.ValueIgnoreEmpty(haMode),
		"eip1":                  buildCreateGatewayEIPChildBody(d, "eip1"),
		"master_eip":            masterEIP,
		"eip2":                  buildCreateGatewayEIPChildBody(d, "eip2"),
		"slave_eip":             slaveEIP,
		"access_vpc_id":         utils.ValueIgnoreEmpty(d.Get("access_vpc_id")),
		"access_subnet_id":      utils.ValueIgnoreEmpty(d.Get("access_subnet_id")),
		"er_id":                 utils.ValueIgnoreEmpty(d.Get("er_id")),
		"network_type":          utils.ValueIgnoreEmpty(d.Get("network_type")),
		"access_private_ip_1":   utils.ValueIgnoreEmpty(d.Get("access_private_ip_1")),
		"access_private_ip_2":   utils.ValueIgnoreEmpty(d.Get("access_private_ip_2")),
		"tags":                  utils.ValueIgnoreEmpty(utils.ExpandResourceTags(d.Get("tags").(map[string]interface{}))),
	}

	return params
}

func buildCreateGatewayEIPChildBody(d *schema.ResourceData, param string) map[string]interface{} {
	if rawArray, ok := d.Get(param).([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"bandwidth_name": utils.ValueIgnoreEmpty(raw["bandwidth_name"]),
			"bandwidth_size": utils.ValueIgnoreEmpty(raw["bandwidth_size"]),
			"charge_mode":    utils.ValueIgnoreEmpty(raw["charge_mode"]),
			"id":             utils.ValueIgnoreEmpty(raw["id"]),
			"type":           utils.ValueIgnoreEmpty(raw["type"]),
		}
		return params
	}
	return nil
}

func waitingForGatewayCertificateStateCompleted(ctx context.Context, d *schema.ResourceData,
	client *golangsdk.ServiceClient, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"BINDING"},
		Target:       []string{"BOUND"},
		Refresh:      waitForGatewayCertificate(client, d),
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForGatewayCertificate(client *golangsdk.ServiceClient, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getGatewayCertificateHttpUrl := "v5/{project_id}/vpn-gateways/{gateway_id}/certificate"
		certificateContent := d.Get("certificate").([]interface{})
		certificateMap := certificateContent[0].(map[string]interface{})
		certificateID := certificateMap["certificate_id"].(string)

		getGatewayCertificatePath := client.Endpoint + getGatewayCertificateHttpUrl
		getGatewayCertificatePath = strings.ReplaceAll(getGatewayCertificatePath, "{project_id}", client.ProjectID)
		getGatewayCertificatePath = strings.ReplaceAll(getGatewayCertificatePath, "{gateway_id}", d.Id())
		getGatewayCertificatePath = strings.ReplaceAll(getGatewayCertificatePath, "{certificate_id}", certificateID)
		getGatewayCertificateOpt := golangsdk.RequestOpts{KeepResponseBody: true}

		gatewayCertificate, err := client.Request("GET", getGatewayCertificatePath, &getGatewayCertificateOpt)
		if err != nil {
			return nil, "ERROR", err
		}

		body, err := utils.FlattenResponse(gatewayCertificate)
		if err != nil {
			return nil, "ERROR", err
		}

		mErr := multierror.Append(nil, d.Set("certificate",
			flattenGatewayCertificateResponse(d, body)))
		if mErr.ErrorOrNil() != nil {
			return nil, "ERROR", mErr.ErrorOrNil()
		}

		certificate := utils.PathSearch("certificate", body, nil)
		status := utils.PathSearch("certificate.status", body, nil).(string)

		return certificate, status, nil
	}
}

func createGatewayWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			// createGatewayWaiting: missing operation notes
			var (
				createGatewayWaitingHttpUrl = "v5/{project_id}/vpn-gateways/{id}"
				createGatewayWaitingProduct = "vpn"
			)
			createGatewayWaitingClient, err := cfg.NewServiceClient(createGatewayWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating VPN client: %s", err)
			}

			createGatewayWaitingPath := createGatewayWaitingClient.Endpoint + createGatewayWaitingHttpUrl
			createGatewayWaitingPath = strings.ReplaceAll(createGatewayWaitingPath, "{project_id}", createGatewayWaitingClient.ProjectID)
			createGatewayWaitingPath = strings.ReplaceAll(createGatewayWaitingPath, "{id}", d.Id())

			createGatewayWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
			}
			createGatewayWaitingResp, err := createGatewayWaitingClient.Request("GET", createGatewayWaitingPath, &createGatewayWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createGatewayWaitingRespBody, err := utils.FlattenResponse(createGatewayWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw := utils.PathSearch(`vpn_gateway.status`, createGatewayWaitingRespBody, nil)
			if statusRaw == nil {
				return nil, "ERROR", fmt.Errorf("error parsing %s from response body", `vpn_gateway.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"ACTIVE",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return createGatewayWaitingRespBody, "COMPLETED", nil
			}

			pendingStatus := []string{
				"PENDING_CREATE",
			}
			if utils.StrSliceContains(pendingStatus, status) {
				return createGatewayWaitingRespBody, "PENDING", nil
			}

			return createGatewayWaitingRespBody, status, nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceGatewayRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getGateway: Query the VPN gateway detail
	var (
		getGatewayHttpUrl            = "v5/{project_id}/vpn-gateways/{id}"
		getGatewayCertificateHttpUrl = "v5/{project_id}/vpn-gateways/{gateway_id}/certificate"
		getGatewayProduct            = "vpn"
	)
	getGatewayClient, err := cfg.NewServiceClient(getGatewayProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	getGatewayPath := getGatewayClient.Endpoint + getGatewayHttpUrl
	getGatewayPath = strings.ReplaceAll(getGatewayPath, "{project_id}", getGatewayClient.ProjectID)
	getGatewayPath = strings.ReplaceAll(getGatewayPath, "{id}", d.Id())

	getGatewayOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getGatewayResp, err := getGatewayClient.Request("GET", getGatewayPath, &getGatewayOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving gateway")
	}

	getGatewayRespBody, err := utils.FlattenResponse(getGatewayResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("attachment_type", utils.PathSearch("vpn_gateway.attachment_type", getGatewayRespBody, nil)),
		d.Set("availability_zones", utils.PathSearch("vpn_gateway.availability_zone_ids", getGatewayRespBody, nil)),
		d.Set("asn", utils.PathSearch("vpn_gateway.bgp_asn", getGatewayRespBody, nil)),
		d.Set("connect_subnet", utils.PathSearch("vpn_gateway.connect_subnet", getGatewayRespBody, nil)),
		d.Set("created_at", utils.PathSearch("vpn_gateway.created_at", getGatewayRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("vpn_gateway.enterprise_project_id", getGatewayRespBody, nil)),
		d.Set("flavor", utils.PathSearch("vpn_gateway.flavor", getGatewayRespBody, nil)),
		d.Set("local_subnets", utils.PathSearch("vpn_gateway.local_subnets", getGatewayRespBody, nil)),
		d.Set("ha_mode", utils.PathSearch("vpn_gateway.ha_mode", getGatewayRespBody, nil)),
		d.Set("eip1", flattenGetGatewayResponseBodyVPNGatewayBody(getGatewayRespBody, "eip1")),
		d.Set("master_eip", flattenGetGatewayResponseBodyVPNGatewayBody(getGatewayRespBody, "master_eip")),
		d.Set("name", utils.PathSearch("vpn_gateway.name", getGatewayRespBody, nil)),
		d.Set("eip2", flattenGetGatewayResponseBodyVPNGatewayBody(getGatewayRespBody, "eip2")),
		d.Set("slave_eip", flattenGetGatewayResponseBodyVPNGatewayBody(getGatewayRespBody, "slave_eip")),
		d.Set("status", utils.PathSearch("vpn_gateway.status", getGatewayRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("vpn_gateway.updated_at", getGatewayRespBody, nil)),
		d.Set("used_connection_group", utils.PathSearch("vpn_gateway.used_connection_group", getGatewayRespBody, nil)),
		d.Set("used_connection_number", utils.PathSearch("vpn_gateway.used_connection_number", getGatewayRespBody, nil)),
		d.Set("vpc_id", utils.PathSearch("vpn_gateway.vpc_id", getGatewayRespBody, nil)),
		d.Set("access_vpc_id", utils.PathSearch("vpn_gateway.access_vpc_id", getGatewayRespBody, nil)),
		d.Set("access_subnet_id", utils.PathSearch("vpn_gateway.access_subnet_id", getGatewayRespBody, nil)),
		d.Set("er_id", utils.PathSearch("vpn_gateway.er_id", getGatewayRespBody, nil)),
		d.Set("er_attachment_id", utils.PathSearch("vpn_gateway.er_attachment_id", getGatewayRespBody, nil)),
		d.Set("network_type", utils.PathSearch("vpn_gateway.network_type", getGatewayRespBody, nil)),
		d.Set("access_private_ip_1", utils.PathSearch("vpn_gateway.access_private_ip_1", getGatewayRespBody, nil)),
		d.Set("access_private_ip_2", utils.PathSearch("vpn_gateway.access_private_ip_2", getGatewayRespBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("vpn_gateway.tags", getGatewayRespBody, nil))),
	)

	certificateContent := d.Get("certificate").([]interface{})
	if len(certificateContent) == 1 {
		certificateMap := certificateContent[0].(map[string]interface{})
		certificateID := certificateMap["certificate_id"].(string)
		getGatewayCertificatePath := getGatewayClient.Endpoint + getGatewayCertificateHttpUrl
		getGatewayCertificatePath = strings.ReplaceAll(getGatewayCertificatePath, "{project_id}", getGatewayClient.ProjectID)
		getGatewayCertificatePath = strings.ReplaceAll(getGatewayCertificatePath, "{gateway_id}", d.Id())
		getGatewayCertificatePath = strings.ReplaceAll(getGatewayCertificatePath, "{certificate_id}", certificateID)
		getGatewayCertificateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getGatewayCertificateResp, err := getGatewayClient.Request("GET", getGatewayCertificatePath, &getGatewayCertificateOpt)
		if err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "failed to retrieving gateway certificate",
					Detail:   fmt.Sprintf("error retrieving gateway certificate: %s.", err),
				},
			}
		}
		getGatewayCertificateRespBody, err := utils.FlattenResponse(getGatewayCertificateResp)
		if err != nil {
			return diag.FromErr(err)
		}
		mErr = multierror.Append(mErr,
			d.Set("certificate", flattenGatewayCertificateResponse(d, getGatewayCertificateRespBody)),
		)
	}
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGatewayCertificateResponse(d *schema.ResourceData, resp interface{}) []interface{} {
	rst := d.Get("certificate").([]interface{})
	curJson := utils.PathSearch("certificate", resp, nil)
	if curJson == nil {
		log.Printf("[ERROR] error parsing certificate from response")
		return rst
	}

	rstCertificate := rst[0].(map[string]interface{})
	rstCertificate["certificate_id"] = utils.PathSearch("id", curJson, nil)
	rstCertificate["status"] = utils.PathSearch("status", curJson, nil)
	rstCertificate["issuer"] = utils.PathSearch("issuer", curJson, nil)
	rstCertificate["signature_algorithm"] = utils.PathSearch("signature_algorithm", curJson, nil)
	rstCertificate["certificate_serial_number"] = utils.PathSearch("certificate_serial_number", curJson, nil)
	rstCertificate["certificate_subject"] = utils.PathSearch("certificate_subject", curJson, nil)
	rstCertificate["certificate_expire_time"] = utils.PathSearch("certificate_expire_time", curJson, nil)
	rstCertificate["certificate_chain_serial_number"] = utils.PathSearch("certificate_chain_serial_number", curJson, nil)
	rstCertificate["certificate_chain_subject"] = utils.PathSearch("certificate_chain_subject", curJson, nil)
	rstCertificate["certificate_chain_expire_time"] = utils.PathSearch("certificate_chain_expire_time", curJson, nil)
	rstCertificate["enc_certificate_serial_number"] = utils.PathSearch("enc_certificate_serial_number", curJson, nil)
	rstCertificate["enc_certificate_subject"] = utils.PathSearch("enc_certificate_subject", curJson, nil)
	rstCertificate["enc_certificate_expire_time"] = utils.PathSearch("enc_certificate_expire_time", curJson, nil)
	rstCertificate["created_at"] = utils.PathSearch("created_at", curJson, nil)
	rstCertificate["updated_at"] = utils.PathSearch("updated_at", curJson, nil)
	return rst
}

func flattenGetGatewayResponseBodyVPNGatewayBody(resp interface{}, paramName string) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch(fmt.Sprintf("vpn_gateway.%s", paramName), resp, nil)
	if curJson == nil {
		log.Printf("[ERROR] error parsing vpn_gateway.%s from response", paramName)
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"bandwidth_id":   utils.PathSearch("bandwidth_id", curJson, nil),
			"bandwidth_name": utils.PathSearch("bandwidth_name", curJson, nil),
			"bandwidth_size": utils.PathSearch("bandwidth_size", curJson, nil),
			"charge_mode":    utils.PathSearch("charge_mode", curJson, nil),
			"id":             utils.PathSearch("id", curJson, nil),
			"ip_address":     utils.PathSearch("ip_address", curJson, nil),
			"ip_version":     utils.PathSearch("ip_version", curJson, nil),
			"type":           utils.PathSearch("type", curJson, nil),
		},
	}
	return rst
}

func resourceGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	gatewayId := d.Id()

	// updateGateway: Update the configuration of VPN gateway
	var (
		updateGatewayHttpUrl            = "v5/{project_id}/vpn-gateways/{id}"
		updateGatewayCertificateHttpUrl = "v5/{project_id}/vpn-gateways/{gateway_id}/certificate/{certificate_id}"
		updateFlavorHttpUrl             = "v5/{project_id}/vpn-gateways/{vgw_id}/update-specification"
		updateGatewayProduct            = "vpn"
	)
	updateGatewayClient, err := cfg.NewServiceClient(updateGatewayProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	updateGatewayHasChanges := []string{
		"local_subnets",
		"name",
	}

	if d.HasChanges(updateGatewayHasChanges...) {
		updateGatewayPath := updateGatewayClient.Endpoint + updateGatewayHttpUrl
		updateGatewayPath = strings.ReplaceAll(updateGatewayPath, "{project_id}", updateGatewayClient.ProjectID)
		updateGatewayPath = strings.ReplaceAll(updateGatewayPath, "{id}", gatewayId)

		updateGatewayOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		updateGatewayOpt.JSONBody = utils.RemoveNil(buildUpdateGatewayBodyParams(d))
		_, err = updateGatewayClient.Request("PUT", updateGatewayPath, &updateGatewayOpt)
		if err != nil {
			return diag.Errorf("error updating VPN gateway: %s", err)
		}
		err = updateGatewayWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for updating VPN gateway (%s) to complete: %s", gatewayId, err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   gatewayId,
			ResourceType: "vpn-gateway",
			RegionId:     region,
			ProjectId:    updateGatewayClient.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("certificate") {
		certificate := d.Get("certificate").([]interface{})
		certificateMap := certificate[0].(map[string]interface{})
		updateGatewayCertificateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		updateGatewayCertificatePath := updateGatewayClient.Endpoint + updateGatewayCertificateHttpUrl
		updateGatewayCertificatePath = strings.ReplaceAll(updateGatewayCertificatePath, "{project_id}", updateGatewayClient.ProjectID)
		updateGatewayCertificatePath = strings.ReplaceAll(updateGatewayCertificatePath, "{gateway_id}", d.Id())
		updateGatewayCertificatePath = strings.ReplaceAll(updateGatewayCertificatePath, "{certificate_id}", certificateMap["certificate_id"].(string))
		updateGatewayCertificateOpt.JSONBody = utils.RemoveNil(buildUpdateGatewayCertificateBodyParams(certificateMap))
		updateGatewayCertificateResp, err := updateGatewayClient.Request("PUT", updateGatewayCertificatePath, &updateGatewayCertificateOpt)
		if err != nil {
			return diag.Errorf("error updating VPN gateway certificate: %s", err)
		}

		updateGatewayCertificateBody, err := utils.FlattenResponse(updateGatewayCertificateResp)
		if err != nil {
			return diag.FromErr(err)
		}

		mErr := multierror.Append(nil,
			d.Set("certificate", flattenGatewayCertificateResponse(d, updateGatewayCertificateBody)),
		)
		if mErr.ErrorOrNil() != nil {
			return diag.FromErr(mErr.ErrorOrNil())
		}

		err = waitingForGatewayCertificateStateCompleted(ctx, d, updateGatewayClient, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for updating VPN gateway (%s) certificate to complete: %s", d.Id(), err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		tagErr := updateTags(updateGatewayClient, d, "vpn-gateway", d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of VPN gateway (%s): %s", d.Id(), tagErr)
		}
	}

	if d.HasChange("flavor") {
		updateFlavorPath := updateGatewayClient.Endpoint + updateFlavorHttpUrl
		updateFlavorPath = strings.ReplaceAll(updateFlavorPath, "{project_id}", updateGatewayClient.ProjectID)
		updateFlavorPath = strings.ReplaceAll(updateFlavorPath, "{vgw_id}", gatewayId)

		updateGatewayOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		updateGatewayOpt.JSONBody = map[string]interface{}{
			"vpn_gateway": map[string]interface{}{
				"flavor": d.Get("flavor"),
			},
		}
		_, err = updateGatewayClient.Request("POST", updateFlavorPath, &updateGatewayOpt)
		if err != nil {
			return diag.Errorf("error updating VPN gateway flavor: %s", err)
		}
		err = updateGatewayWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for updating VPN gateway (%s) flavor to complete: %s", gatewayId, err)
		}
	}

	return resourceGatewayRead(ctx, d, meta)
}

func buildUpdateGatewayCertificateBodyParams(certificateContent interface{}) map[string]interface{} {
	certificateMap := certificateContent.(map[string]interface{})
	return map[string]interface{}{
		"certificate": map[string]interface{}{
			"name":              utils.ValueIgnoreEmpty(certificateMap["name"].(string)),
			"certificate":       utils.ValueIgnoreEmpty(certificateMap["content"].(string)),
			"private_key":       utils.ValueIgnoreEmpty(certificateMap["private_key"].(string)),
			"certificate_chain": utils.ValueIgnoreEmpty(certificateMap["certificate_chain"].(string)),
			"enc_certificate":   utils.ValueIgnoreEmpty(certificateMap["enc_certificate"].(string)),
			"enc_private_key":   utils.ValueIgnoreEmpty(certificateMap["enc_private_key"].(string)),
		},
	}
}

func buildUpdateGatewayBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"vpn_gateway": buildUpdateGatewayVpnGatewayChildBody(d),
	}
	return bodyParams
}

func buildUpdateGatewayVpnGatewayChildBody(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"local_subnets": utils.ValueIgnoreEmpty(d.Get("local_subnets")),
		"name":          utils.ValueIgnoreEmpty(d.Get("name")),
	}
	return params
}

func updateGatewayWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			// updateGatewayWaiting: missing operation notes
			var (
				updateGatewayWaitingHttpUrl = "v5/{project_id}/vpn-gateways/{id}"
				updateGatewayWaitingProduct = "vpn"
			)
			updateGatewayWaitingClient, err := cfg.NewServiceClient(updateGatewayWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating VPN client: %s", err)
			}

			updateGatewayWaitingPath := updateGatewayWaitingClient.Endpoint + updateGatewayWaitingHttpUrl
			updateGatewayWaitingPath = strings.ReplaceAll(updateGatewayWaitingPath, "{project_id}", updateGatewayWaitingClient.ProjectID)
			updateGatewayWaitingPath = strings.ReplaceAll(updateGatewayWaitingPath, "{id}", d.Id())

			updateGatewayWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
			}
			updateGatewayWaitingResp, err := updateGatewayWaitingClient.Request("GET", updateGatewayWaitingPath, &updateGatewayWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			updateGatewayWaitingRespBody, err := utils.FlattenResponse(updateGatewayWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw := utils.PathSearch(`vpn_gateway.status`, updateGatewayWaitingRespBody, nil)
			if statusRaw == nil {
				return nil, "ERROR", fmt.Errorf("error parsing %s from response body", `vpn_gateway.status`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"ACTIVE",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return updateGatewayWaitingRespBody, "COMPLETED", nil
			}

			pendingStatus := []string{
				"PENDING_UPDATE",
			}
			if utils.StrSliceContains(pendingStatus, status) {
				return updateGatewayWaitingRespBody, "PENDING", nil
			}

			return updateGatewayWaitingRespBody, status, nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	gatewayId := d.Id()

	// deleteGateway: Delete an existing VPN Gateway
	var (
		deleteGatewayHttpUrl = "v5/{project_id}/vpn-gateways/{id}"
		deleteGatewayProduct = "vpn"
	)
	deleteGatewayClient, err := cfg.NewServiceClient(deleteGatewayProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	deleteEipOnTermination := d.Get("delete_eip_on_termination").(bool)
	if !deleteEipOnTermination {
		eipIDs := getEipsToPreserve(d)
		err := disassociateEip(cfg, region, eipIDs)
		if err != nil {
			return diag.Errorf("error deleting VPN gateway (%s): %s", gatewayId, err)
		}
	}

	deleteGatewayPath := deleteGatewayClient.Endpoint + deleteGatewayHttpUrl
	deleteGatewayPath = strings.ReplaceAll(deleteGatewayPath, "{project_id}", deleteGatewayClient.ProjectID)
	deleteGatewayPath = strings.ReplaceAll(deleteGatewayPath, "{id}", gatewayId)

	deleteGatewayOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteGatewayClient.Request("DELETE", deleteGatewayPath, &deleteGatewayOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting VPN gateway")
	}

	err = deleteGatewayWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for deleting VPN gateway (%s) to complete: %s", gatewayId, err)
	}
	return nil
}

func getEipsToPreserve(d *schema.ResourceData) map[string]struct{} {
	eipKeys := []string{"eip1.0.id", "eip2.0.id", "master_eip.0.id", "slave_eip.0.id"}
	eipIDs := make(map[string]struct{})

	for _, key := range eipKeys {
		if eipID := d.Get(key).(string); eipID != "" {
			eipIDs[eipID] = struct{}{}
		}
	}
	return eipIDs
}

func disassociateEip(cfg *config.Config, region string, eipIDs map[string]struct{}) error {
	if len(eipIDs) == 0 {
		return nil
	}

	var (
		disassociateEipHttpUrl = "v3/{project_id}/eip/publicips/{publicip_id}/disassociate-instance"
		disassociateEipProduct = "vpc"
	)

	disassociateEipClient, err := cfg.NewServiceClient(disassociateEipProduct, region)
	if err != nil {
		return fmt.Errorf("error creating VPC client: %s", err)
	}

	disassociateEipPath := disassociateEipClient.Endpoint + disassociateEipHttpUrl
	disassociateEipPath = strings.ReplaceAll(disassociateEipPath, "{project_id}", disassociateEipClient.ProjectID)

	disassociateEipOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for eipID := range eipIDs {
		path := strings.ReplaceAll(disassociateEipPath, "{publicip_id}", eipID)
		_, err = disassociateEipClient.Request("POST", path, &disassociateEipOpt)
		if err != nil {
			return fmt.Errorf("error disassociating EIP (%s) from VPN gateway: %s", eipID, err)
		}
	}

	return nil
}

func deleteGatewayWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			// deleteGatewayWaiting: missing operation notes
			var (
				deleteGatewayWaitingHttpUrl = "v5/{project_id}/vpn-gateways/{id}"
				deleteGatewayWaitingProduct = "vpn"
			)
			deleteGatewayWaitingClient, err := cfg.NewServiceClient(deleteGatewayWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating VPN client: %s", err)
			}

			deleteGatewayWaitingPath := deleteGatewayWaitingClient.Endpoint + deleteGatewayWaitingHttpUrl
			deleteGatewayWaitingPath = strings.ReplaceAll(deleteGatewayWaitingPath, "{project_id}", deleteGatewayWaitingClient.ProjectID)
			deleteGatewayWaitingPath = strings.ReplaceAll(deleteGatewayWaitingPath, "{id}", d.Id())

			deleteGatewayWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
			}
			deleteGatewayWaitingResp, err := deleteGatewayWaitingClient.Request("GET", deleteGatewayWaitingPath, &deleteGatewayWaitingOpt)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
					return "Resource Not Found", "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			deleteGatewayWaitingRespBody, err := utils.FlattenResponse(deleteGatewayWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}

			return deleteGatewayWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func updateTags(client *golangsdk.ServiceClient, d *schema.ResourceData, tagsType string, id string) error {
	oRaw, nRaw := d.GetChange("tags")
	oMap := oRaw.(map[string]interface{})
	nMap := nRaw.(map[string]interface{})

	manageTagsHttpUrl := "v5/{project_id}/{resource_type}/{resource_id}/tags/{action}"
	manageTagsPath := client.Endpoint + manageTagsHttpUrl
	manageTagsPath = strings.ReplaceAll(manageTagsPath, "{project_id}", client.ProjectID)
	manageTagsPath = strings.ReplaceAll(manageTagsPath, "{resource_type}", tagsType)
	manageTagsPath = strings.ReplaceAll(manageTagsPath, "{resource_id}", id)
	manageTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	// remove old tags
	if len(oMap) > 0 {
		manageTagsOpt.JSONBody = map[string]interface{}{
			"tags": utils.ExpandResourceTags(oMap),
		}
		deleteTagsPath := strings.ReplaceAll(manageTagsPath, "{action}", "delete")
		_, err := client.Request("POST", deleteTagsPath, &manageTagsOpt)
		if err != nil {
			return err
		}
	}

	// set new tags
	if len(nMap) > 0 {
		manageTagsOpt.JSONBody = map[string]interface{}{
			"tags": utils.ExpandResourceTags(nMap),
		}
		createTagsPath := strings.ReplaceAll(manageTagsPath, "{action}", "create")
		_, err := client.Request("POST", createTagsPath, &manageTagsOpt)
		if err != nil {
			return err
		}
	}

	return nil
}
