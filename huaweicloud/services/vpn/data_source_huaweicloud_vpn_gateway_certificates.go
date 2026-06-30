package vpn

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPN GET /v5/{project_id}/vpn-gateway-certificates
func DataSourceVpnGatewayCertificates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpnGatewayCertificatesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     vpnGatewayCertificatesSchema(),
			},
		},
	}
}

func vpnGatewayCertificatesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vgw_id": {
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

func dataSourceVpnGatewayCertificatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("vpn", region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	httpUrl := "v5/{project_id}/vpn-gateway-certificates"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving VPN gateway certificates: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("certificates", flattenGetVpnGatewayCertificatesBody(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetVpnGatewayCertificatesBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("vpn_gateway_certificates", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":                              utils.PathSearch("id", v, nil),
			"name":                            utils.PathSearch("name", v, nil),
			"project_id":                      utils.PathSearch("project_id", v, nil),
			"vgw_id":                          utils.PathSearch("vgw_id", v, nil),
			"status":                          utils.PathSearch("status", v, nil),
			"issuer":                          utils.PathSearch("issuer", v, nil),
			"signature_algorithm":             utils.PathSearch("signature_algorithm", v, nil),
			"certificate_serial_number":       utils.PathSearch("certificate_serial_number", v, nil),
			"certificate_subject":             utils.PathSearch("certificate_subject", v, nil),
			"certificate_expire_time":         utils.PathSearch("certificate_expire_time", v, nil),
			"certificate_chain_serial_number": utils.PathSearch("certificate_chain_serial_number", v, nil),
			"certificate_chain_subject":       utils.PathSearch("certificate_chain_subject", v, nil),
			"certificate_chain_expire_time":   utils.PathSearch("certificate_chain_expire_time", v, nil),
			"enc_certificate_serial_number":   utils.PathSearch("enc_certificate_serial_number", v, nil),
			"enc_certificate_subject":         utils.PathSearch("enc_certificate_subject", v, nil),
			"enc_certificate_expire_time":     utils.PathSearch("enc_certificate_expire_time", v, nil),
			"created_at":                      utils.PathSearch("created_at", v, nil),
			"updated_at":                      utils.PathSearch("updated_at", v, nil),
		})
	}
	return res
}
