package vpn

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var serverNonUpdatableParams = []string{"p2c_vgw_id", "tunnel_protocol", "ssl_options.0.is_compressed"}

// @API VPN POST /v5/{project_id}/p2c-vpn-gateways/{p2c_vgw_id}/vpn-servers
// @API VPN GET /v5/{project_id}/p2c-vpn-gateways/{p2c_vgw_id}/vpn-servers
// @API VPN PUT /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}
// @API VPN POST /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/client-config/export
// @API VPN POST /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/client-ca-certificates
// @API VPN DELETE /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/client-ca-certificates/{certificate_id}
func ResourceServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerCreate,
		ReadContext:   resourceServerRead,
		UpdateContext: resourceServerUpdate,
		DeleteContext: resourceServerDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceServerImportState,
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(serverNonUpdatableParams),
			func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
				// when `client_auth_type` is LOCAL_PASSWORD, client_ca_certificates should be empty
				clientAuthType := d.Get("client_auth_type").(string)
				clientCaCertificates := d.Get("client_ca_certificates").(*schema.Set).List()
				if (clientAuthType == "" || clientAuthType == "LOCAL_PASSWORD") && len(clientCaCertificates) > 0 {
					return errors.New("client CA certificates should be empty when client auth type is LOCAL_PASSWORD")
				}
				return nil
			},
		),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"p2c_vgw_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of a P2C VPN gateway instance.`,
			},
			"local_subnets": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `Specifies the list of local CIDR blocks.`,
			},
			"client_cidr": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the client CIDR block.`,
			},
			"tunnel_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the tunnel protocol.`,
			},
			"client_auth_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the client authentication mode.`,
			},
			"server_certificate": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the server certificate info.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the certificate ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The server certificate name.`,
						},
						"serial_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The serial number of the server certificate.`,
						},
						"expiration_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The expiration time of the server certificate.`,
						},
						"signature_algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The signature algorithm of the server certificate.`,
						},
						"issuer": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The issuer of the server certificate.`,
						},
						"subject": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The subject of the server certificate.`,
						},
					},
				},
			},
			"ssl_options": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the SSL options.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the protocol.`,
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the port number.`,
						},
						"encryption_algorithm": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the encryption algorithm.`,
						},
						"is_compressed": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: `Specifies whether to compress data.`,
						},
						"authentication_algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The authentication algorithm.`,
						},
					},
				},
			},
			// content is not in return, and name is optional, can not match content with id
			// add use client_ca_certificates_uploaded as computed attributes
			// when update, remove all client_ca_certificates_uploaded and upload all new client_ca_certificates
			"client_ca_certificates": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `Specifies the list of client CA certificates.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the certificate content.`,
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the certificate name.`,
						},
					},
				},
			},
			"client_ca_certificates_uploaded": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of client CA certificates.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the certificate name`,
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The client CA certificate ID.`,
						},
						"issuer": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The issuer of the client CA certificate.`,
						},
						"signature_algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The signature algorithm of the client CA certificate.`,
						},
						"subject": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The subject of the client CA certificate.`,
						},
						"serial_number": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The serial number of the client CA certificate.`,
						},
						"expiration_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The expiration time of the client CA certificate.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the client CA certificate.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The update time of the client CA certificate.`,
						},
					},
				},
			},
			"os_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the OS type.`,
			},
			"client_config": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The client config.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The server status.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.NewServiceClient("vpn", region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	p2cVgwId := d.Get("p2c_vgw_id").(string)
	createServerHttpUrl := "v5/{project_id}/p2c-vpn-gateways/{p2c_vgw_id}/vpn-servers"
	createServerPath := client.Endpoint + createServerHttpUrl
	createServerPath = strings.ReplaceAll(createServerPath, "{project_id}", client.ProjectID)
	createServerPath = strings.ReplaceAll(createServerPath, "{p2c_vgw_id}", p2cVgwId)
	createServerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateServerBodyParams(d)),
	}

	createServerResp, err := client.Request("POST", createServerPath, &createServerOpt)
	if err != nil {
		return diag.Errorf("error creating VPN server: %s", err)
	}
	createServerRespBody, err := utils.FlattenResponse(createServerResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("vpn_server.id", createServerRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find server ID in API response")
	}
	d.SetId(id)

	err = waitForVpnServerActive(ctx, client, p2cVgwId, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for VPN server (%s) to be active: %s", d.Id(), err)
	}

	return resourceServerRead(ctx, d, meta)
}

func buildCreateServerBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"vpn_server": map[string]interface{}{
			"local_subnets":          d.Get("local_subnets").(*schema.Set).List(),
			"client_cidr":            d.Get("client_cidr"),
			"tunnel_protocol":        utils.ValueIgnoreEmpty(d.Get("tunnel_protocol")),
			"client_auth_type":       utils.ValueIgnoreEmpty(d.Get("client_auth_type")),
			"server_certificate":     buildServerBodyParamsServerCertificate(d.Get("server_certificate").([]interface{})),
			"client_ca_certificates": buildServerBodyParamsClientCACertificates(d.Get("client_ca_certificates").(*schema.Set).List()),
			"ssl_options":            buildCreateServerBodyParamsSslOptions(d),
		},
	}
	return bodyParams
}

func buildServerBodyParamsServerCertificate(paramsList []interface{}) map[string]interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	params := paramsList[0].(map[string]interface{})
	rst := map[string]interface{}{
		"id": params["id"],
	}
	return rst
}

func buildServerBodyParamsClientCACertificates(paramsList []interface{}) []map[string]interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0)
	for _, rawParams := range paramsList {
		params := rawParams.(map[string]interface{})
		m := map[string]interface{}{
			"content": params["content"],
			"name":    utils.ValueIgnoreEmpty(params["name"]),
		}
		rst = append(rst, m)
	}

	return rst
}

func buildCreateServerBodyParamsSslOptions(d *schema.ResourceData) interface{} {
	var rst map[string]interface{}
	paramsList := d.Get("ssl_options").([]interface{})
	if len(paramsList) != 0 {
		if params, ok := paramsList[0].(map[string]interface{}); ok {
			rst = utils.RemoveNil(map[string]interface{}{
				"protocol":             utils.ValueIgnoreEmpty(params["protocol"]),
				"port":                 utils.ValueIgnoreEmpty(params["port"]),
				"encryption_algorithm": utils.ValueIgnoreEmpty(params["encryption_algorithm"]),
				"is_compressed":        utils.ValueIgnoreEmpty(params["is_compressed"]),
			})
		}
	}

	// when tunnel protocol (default to SSL) is SSL, ssl options block is required
	tunnelProtocol := d.Get("tunnel_protocol").(string)
	if (tunnelProtocol == "" || tunnelProtocol == "SSL") && (rst == nil || reflect.DeepEqual(rst, map[string]interface{}{})) {
		return &map[string]interface{}{}
	}

	return rst
}

func waitForVpnServerActive(ctx context.Context, client *golangsdk.ServiceClient, p2cVgwId, id string, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"SUCCESS"},
		Refresh: func() (interface{}, string, error) {
			server, err := GetServer(client, p2cVgwId, id)
			if err != nil {
				return nil, "ERROR", err
			}
			status := utils.PathSearch("status", server, nil).(string)
			switch status {
			case "ACTIVE":
				return server, "SUCCESS", nil
			case "FAULT", "FROZEN":
				return server, "FAILED", fmt.Errorf("got abnormal status: %s", status)
			default:
				return server, "PENDING", nil
			}
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func resourceServerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.NewServiceClient("vpn", region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	server, err := GetServer(client, d.Get("p2c_vgw_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VPN server")
	}

	clientConfig, err := getServcerClientConifg(client, d)
	if err != nil {
		log.Printf("[WARN] unable to find clinet config: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("p2c_vgw_id", utils.PathSearch("p2c_vgw_id", server, nil)),
		d.Set("local_subnets", utils.PathSearch("local_subnets", server, nil)),
		d.Set("client_cidr", utils.PathSearch("client_cidr", server, nil)),
		d.Set("tunnel_protocol", utils.PathSearch("tunnel_protocol", server, nil)),
		d.Set("client_auth_type", utils.PathSearch("client_auth_type", server, nil)),
		d.Set("server_certificate", flattenServerServerCertificate(utils.PathSearch("server_certificate", server, nil))),
		d.Set("client_ca_certificates_uploaded", flattenServerClientCACertificates(
			utils.PathSearch("client_ca_certificates", server, make([]interface{}, 0)).([]interface{}))),
		d.Set("ssl_options", flattenServerSslOptions(utils.PathSearch("ssl_options", server, nil))),
		d.Set("status", utils.PathSearch("status", server, nil)),
		d.Set("created_at", utils.PathSearch("created_at", server, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", server, nil)),
		d.Set("client_config", utils.PathSearch("client_config", clientConfig, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetServer(client *golangsdk.ServiceClient, p2cVgwId, serverId string) (interface{}, error) {
	getServerHttpUrl := "v5/{project_id}/p2c-vpn-gateways/{p2c_vgw_id}/vpn-servers"
	getServerPath := client.Endpoint + getServerHttpUrl
	getServerPath = strings.ReplaceAll(getServerPath, "{project_id}", client.ProjectID)
	getServerPath = strings.ReplaceAll(getServerPath, "{p2c_vgw_id}", p2cVgwId)
	getServerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getServerResp, err := client.Request("GET", getServerPath, &getServerOpt)
	if err != nil {
		return nil, err
	}
	getServerRespBody, err := utils.FlattenResponse(getServerResp)
	if err != nil {
		return nil, err
	}

	searchPath := fmt.Sprintf("vpn_servers[?id == '%s']|[0]", serverId)
	server := utils.PathSearch(searchPath, getServerRespBody, nil)
	if server == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return server, nil
}

func getServcerClientConifg(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	getServerHttpUrl := "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/client-config/export"
	getServerPath := client.Endpoint + getServerHttpUrl
	getServerPath = strings.ReplaceAll(getServerPath, "{project_id}", client.ProjectID)
	getServerPath = strings.ReplaceAll(getServerPath, "{vpn_server_id}", d.Id())
	getServerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: utils.RemoveNil(map[string]interface{}{
			"os_type": utils.ValueIgnoreEmpty(d.Get("os_type")),
		}),
	}

	getServerResp, err := client.Request("POST", getServerPath, &getServerOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getServerResp)
}

func flattenServerServerCertificate(params interface{}) []map[string]interface{} {
	if params == nil {
		return nil
	}
	rst := map[string]interface{}{
		"id":                  utils.PathSearch("id", params, nil),
		"name":                utils.PathSearch("name", params, nil),
		"serial_number":       utils.PathSearch("serial_number", params, nil),
		"expiration_time":     utils.PathSearch("expiration_time", params, nil),
		"signature_algorithm": utils.PathSearch("signature_algorithm", params, nil),
		"issuer":              utils.PathSearch("issuer", params, nil),
		"subject":             utils.PathSearch("subject", params, nil),
	}

	return []map[string]interface{}{rst}
}

func flattenServerClientCACertificates(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"name":                utils.PathSearch("name", params, "").(string),
			"id":                  utils.PathSearch("id", params, nil),
			"issuer":              utils.PathSearch("issuer", params, nil),
			"serial_number":       utils.PathSearch("serial_number", params, nil),
			"subject":             utils.PathSearch("subject", params, nil),
			"signature_algorithm": utils.PathSearch("signature_algorithm", params, nil),
			"expiration_time":     utils.PathSearch("expiration_time", params, nil),
			"created_at":          utils.PathSearch("created_at", params, nil),
			"updated_at":          utils.PathSearch("updated_at", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenServerSslOptions(params interface{}) []map[string]interface{} {
	if params == nil {
		return nil
	}
	rst := map[string]interface{}{
		"protocol":                 utils.PathSearch("protocol", params, nil),
		"port":                     utils.PathSearch("port", params, nil),
		"encryption_algorithm":     utils.PathSearch("encryption_algorithm", params, nil),
		"is_compressed":            utils.PathSearch("is_compressed", params, nil),
		"authentication_algorithm": utils.PathSearch("authentication_algorithm", params, nil),
	}

	return []map[string]interface{}{rst}
}

func resourceServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.NewServiceClient("vpn", region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	// remove first for auth type change from CERT to LOCAL_PASSWORD
	if d.HasChange("client_ca_certificates") {
		remove := d.Get("client_ca_certificates_uploaded").([]interface{})
		if len(remove) > 0 {
			err := removeClientCaCerts(ctx, client, d, remove)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	changes := []string{
		"local_subnets",
		"client_cidr",
		"client_auth_type",
		"server_certificate.0.id",
		"ssl_options.0.protocol",
		"ssl_options.0.port",
		"ssl_options.0.encryption_algorithm",
	}
	if d.HasChanges(changes...) {
		updateServerHttpUrl := "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}"
		updateServerPath := client.Endpoint + updateServerHttpUrl
		updateServerPath = strings.ReplaceAll(updateServerPath, "{project_id}", client.ProjectID)
		updateServerPath = strings.ReplaceAll(updateServerPath, "{vpn_server_id}", d.Id())
		updateServerOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateServerBodyParams(d)),
		}

		_, err = client.Request("PUT", updateServerPath, &updateServerOpt)
		if err != nil {
			return diag.Errorf("error updating VPN server: %s", err)
		}
		err = waitForVpnServerActive(ctx, client, d.Get("p2c_vgw_id").(string), d.Id(), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for VPN server (%s) to be active: %s", d.Id(), err)
		}
	}

	// add after updateServer for auth type change from LOCAL_PASSWORD to CERT
	if d.HasChange("client_ca_certificates") {
		add := d.Get("client_ca_certificates").(*schema.Set).List()
		if len(add) > 0 {
			err := addClientCaCerts(ctx, client, d, add)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceServerRead(ctx, d, meta)
}

func removeClientCaCerts(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, remove []interface{}) error {
	for _, rawCert := range remove {
		if cert, ok := rawCert.(map[string]interface{}); ok && cert["id"].(string) != "" {
			deleteCertHttpUrl := "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/client-ca-certificates/{certificate_id}"
			deleteCertPath := client.Endpoint + deleteCertHttpUrl
			deleteCertPath = strings.ReplaceAll(deleteCertPath, "{project_id}", client.ProjectID)
			deleteCertPath = strings.ReplaceAll(deleteCertPath, "{vpn_server_id}", d.Id())
			deleteCertPath = strings.ReplaceAll(deleteCertPath, "{certificate_id}", cert["id"].(string))

			deleteCertOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
			}
			_, err := client.Request("DELETE", deleteCertPath, &deleteCertOpt)
			if err != nil {
				return fmt.Errorf("error removing client CA certificate(%s): %s", cert["id"].(string), err)
			}

			err = waitForVpnServerActive(ctx, client, d.Get("p2c_vgw_id").(string), d.Id(), d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return fmt.Errorf("error waiting for VPN server (%s) to be active: %s", d.Id(), err)
			}
		}
	}

	return nil
}

func addClientCaCerts(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, add []interface{}) error {
	for _, rawCert := range add {
		if cert, ok := rawCert.(map[string]interface{}); ok {
			addCertHttpUrl := "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/client-ca-certificates"
			addCertPath := client.Endpoint + addCertHttpUrl
			addCertPath = strings.ReplaceAll(addCertPath, "{project_id}", client.ProjectID)
			addCertPath = strings.ReplaceAll(addCertPath, "{vpn_server_id}", d.Id())

			addCertOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				JSONBody:         utils.RemoveNil(buildUploadClientCaCertBodyParams(cert)),
			}
			_, err := client.Request("POST", addCertPath, &addCertOpt)
			if err != nil {
				return fmt.Errorf("error adding client CA certificate(%v): %s", cert["content"], err)
			}
			err = waitForVpnServerActive(ctx, client, d.Get("p2c_vgw_id").(string), d.Id(), d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return fmt.Errorf("error waiting for VPN server (%s) to be active: %s", d.Id(), err)
			}
		}
	}

	return nil
}

func buildUploadClientCaCertBodyParams(params map[string]interface{}) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"client_ca_certificate": map[string]interface{}{
			"content": params["content"],
			"name":    utils.ValueIgnoreEmpty(params["name"]),
		},
	}
	return bodyParams
}

func buildUpdateServerBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"vpn_server": map[string]interface{}{
			"local_subnets":      d.Get("local_subnets").(*schema.Set).List(),
			"client_cidr":        d.Get("client_cidr"),
			"client_auth_type":   utils.ValueIgnoreEmpty(d.Get("client_auth_type")),
			"server_certificate": buildServerBodyParamsServerCertificate(d.Get("server_certificate").([]interface{})),
			"ssl_options":        buildUpdateServerBodyParamsSslOptions(d.Get("ssl_options").([]interface{})),
		},
	}
	return bodyParams
}

func buildUpdateServerBodyParamsSslOptions(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	if params, ok := paramsList[0].(map[string]interface{}); ok {
		rst := map[string]interface{}{
			"protocol":             utils.ValueIgnoreEmpty(params["protocol"]),
			"port":                 utils.ValueIgnoreEmpty(params["port"]),
			"encryption_algorithm": utils.ValueIgnoreEmpty(params["encryption_algorithm"]),
		}
		return rst
	}
	return nil
}

func resourceServerDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting vpn server is not supported. The vpn server is only removed from the state, but it remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func resourceServerImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid ID format, must be <p2c_vgw_id>/<server_id>")
	}

	d.Set("p2c_vgw_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
