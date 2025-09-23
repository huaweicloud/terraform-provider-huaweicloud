package vpn

import (
	"context"
	"fmt"
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

var clientCACertificateNonUpdatableParams = []string{"vpn_server_id", "content"}

// @API VPN POST /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/client-ca-certificates
// @API VPN GET /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/client-ca-certificates/{certificate_id}
// @API VPN PUT /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/client-ca-certificates/{certificate_id}
// @API VPN DELETE /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/client-ca-certificates/{certificate_id}
// @API VPN GET /v5/{project_id}/vpn-servers
func ResourceClientCACertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClientCACertificateCreate,
		UpdateContext: resourceClientCACertificateUpdate,
		ReadContext:   resourceClientCACertificateRead,
		DeleteContext: resourceClientCACertificateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceClientCACertificateImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(clientCACertificateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpn_server_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The VPN server ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of client CA certificate.`,
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `The content of client CA certificate.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"issuer": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The client CA certificate issuer.`,
			},
			"subject": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The client CA certificate subject.`,
			},
			"serial_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The client CA certificate serial number.`,
			},
			"expiration_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The client CA certificate expiration time.`,
			},
			"signature_algorithm": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The signature algorithm of the client CA certificate.`,
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
		},
	}
}

func resourceClientCACertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var (
		createClientCACertificateHttpUrl = "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/client-ca-certificates"
		createClientCACertificateProduct = "vpn"
	)
	clientCACertificateClient, err := conf.NewServiceClient(createClientCACertificateProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	createClientCACertificatePath := clientCACertificateClient.Endpoint + createClientCACertificateHttpUrl
	createClientCACertificatePath = strings.ReplaceAll(createClientCACertificatePath, "{project_id}", clientCACertificateClient.ProjectID)
	createClientCACertificatePath = strings.ReplaceAll(createClientCACertificatePath, "{vpn_server_id}", d.Get("vpn_server_id").(string))

	createClientCACertificateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createClientCACertificateOpt.JSONBody = buildCreateClientCACertificateBodyParams(d)
	createClientCACertificateResp, err := clientCACertificateClient.Request("POST", createClientCACertificatePath, &createClientCACertificateOpt)
	if err != nil {
		return diag.Errorf("error creating VPN client CA certificate: %s", err)
	}

	createClientCACertificateRespBody, err := utils.FlattenResponse(createClientCACertificateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("client_ca_certificate.id", createClientCACertificateRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find client CA certificate ID in API response")
	}
	d.SetId(id)

	serverId := d.Get("vpn_server_id").(string)
	err = waitingForClientCACertificateStateCompleted(ctx, clientCACertificateClient, serverId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for creating VPN client CA certificate (%s) to complete: %s", id, err)
	}

	return resourceClientCACertificateRead(ctx, d, meta)
}

func buildCreateClientCACertificateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"client_ca_certificate": map[string]interface{}{
			"name":    d.Get("name"),
			"content": d.Get("content"),
		},
	}
	return bodyParams
}

func waitingForClientCACertificateStateCompleted(ctx context.Context, client *golangsdk.ServiceClient, serverId string, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING_UPDATE"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			httpUrl := "v5/{project_id}/vpn-servers"
			path := client.Endpoint + httpUrl
			path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)

			opt := golangsdk.RequestOpts{
				KeepResponseBody: true,
			}

			waitingResp, err := client.Request("GET", path, &opt)
			if err != nil {
				return nil, "ERROR", err
			}

			waitingRespBody, err := utils.FlattenResponse(waitingResp)
			if err != nil {
				return nil, "ERROR", err
			}

			findStatusStr := fmt.Sprintf("vpn_servers[?id == '%s']|[0].status", serverId)
			statusRaw := utils.PathSearch(findStatusStr, waitingRespBody, nil)
			if statusRaw == nil {
				return nil, "ERROR", fmt.Errorf("error parsing status from response body")
			}

			status := fmt.Sprintf("%v", statusRaw)

			if status == "ACTIVE" {
				return waitingRespBody, "COMPLETED", nil
			}

			errorStatus := []string{
				"FAULT",
				"FROZEN",
			}
			if utils.StrSliceContains(errorStatus, status) {
				return waitingRespBody, "ERROR", fmt.Errorf("VPN server status is abnormal")
			}

			return waitingRespBody, status, nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceClientCACertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var mErr *multierror.Error

	getClientCACertificateProduct := "vpn"
	getClientCACertificateClient, err := conf.NewServiceClient(getClientCACertificateProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	serverId := d.Get("vpn_server_id").(string)
	id := d.Id()
	getClientCACertificateRespBody, err := GetClientCACertificate(getClientCACertificateClient, serverId, id)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VPN client CA certificate")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("client_ca_certificate.name", getClientCACertificateRespBody, nil)),
		d.Set("issuer", utils.PathSearch("client_ca_certificate.issuer", getClientCACertificateRespBody, nil)),
		d.Set("subject", utils.PathSearch("client_ca_certificate.subject", getClientCACertificateRespBody, nil)),
		d.Set("serial_number", utils.PathSearch("client_ca_certificate.serial_number", getClientCACertificateRespBody, nil)),
		d.Set("expiration_time", utils.PathSearch("client_ca_certificate.expiration_time", getClientCACertificateRespBody, nil)),
		d.Set("signature_algorithm", utils.PathSearch("client_ca_certificate.signature_algorithm", getClientCACertificateRespBody, nil)),
		d.Set("created_at", utils.PathSearch("client_ca_certificate.created_at", getClientCACertificateRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("client_ca_certificate.updated_at", getClientCACertificateRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetClientCACertificate(client *golangsdk.ServiceClient, serverId, id string) (interface{}, error) {
	getClientCACertificateHttpUrl := "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/client-ca-certificates/{certificate_id}"
	getClientCACertificatePath := buildClientCACertificateURL(client, getClientCACertificateHttpUrl, serverId, id)

	getClientCACertificateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getClientCACertificateResp, err := client.Request("GET", getClientCACertificatePath, &getClientCACertificateOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getClientCACertificateResp)
}

func resourceClientCACertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	updateClientCACertificateClient, err := conf.NewServiceClient("vpn", region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	serverId := d.Get("vpn_server_id").(string)
	id := d.Id()
	updateClientCACertificateHttpUrl := "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/client-ca-certificates/{certificate_id}"
	updateClientCACertificatePath := buildClientCACertificateURL(updateClientCACertificateClient, updateClientCACertificateHttpUrl, serverId, id)

	updateClientCACertificateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateClientCACertificateOpt.JSONBody = map[string]interface{}{
		"client_ca_certificate": map[string]interface{}{
			"name": d.Get("name"),
		},
	}
	_, err = updateClientCACertificateClient.Request("PUT", updateClientCACertificatePath, &updateClientCACertificateOpt)
	if err != nil {
		return diag.Errorf("error updating VPN client CA certificate: %s", err)
	}

	err = waitingForClientCACertificateStateCompleted(ctx, updateClientCACertificateClient, serverId, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.Errorf("error waiting for updating VPN client CA certificate (%s) to complete: %s", id, err)
	}

	return resourceClientCACertificateRead(ctx, d, meta)
}

func buildClientCACertificateURL(client *golangsdk.ServiceClient, urlTemplate, serverId, id string) string {
	url := client.Endpoint + urlTemplate
	url = strings.ReplaceAll(url, "{project_id}", client.ProjectID)
	url = strings.ReplaceAll(url, "{vpn_server_id}", serverId)
	url = strings.ReplaceAll(url, "{certificate_id}", id)
	return url
}

func resourceClientCACertificateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var (
		deleteClientCACertificateHttpUrl = "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/client-ca-certificates/{certificate_id}"
		deleteClientCACertificateProduct = "vpn"
	)
	deleteClientCACertificateClient, err := conf.NewServiceClient(deleteClientCACertificateProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	serverId := d.Get("vpn_server_id").(string)
	id := d.Id()
	deleteClientCACertificate := buildClientCACertificateURL(deleteClientCACertificateClient, deleteClientCACertificateHttpUrl, serverId, id)

	deleteClientCACertificateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteClientCACertificateClient.Request("DELETE", deleteClientCACertificate, &deleteClientCACertificateOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting VPN Client CA certificate")
	}

	err = waitingForClientCACertificateStateCompleted(ctx, deleteClientCACertificateClient, serverId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for deleting VPN client CA certificate (%s) to complete: %s", id, err)
	}

	return nil
}

func resourceClientCACertificateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ID format, must be <vpn_server_id>/<id>")
	}

	d.Set("vpn_server_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
