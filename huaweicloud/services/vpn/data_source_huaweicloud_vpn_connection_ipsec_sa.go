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

// @API VPN GET /v5/{project_id}/vpn-connection/{vpn_connection_id}/ipsec-sa
func DataSourceVpnConnectionIpsecSa() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpnConnectionIpsecSaRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpn_connection_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sa_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     ipsecSaInfoSchema(),
			},
		},
	}
}

func ipsecSaInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_ip_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dest_ip_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"packets_sent": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"packets_recv": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"traffic_sent": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"traffic_recv": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"collected_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVpnConnectionIpsecSaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v5/{project_id}/vpn-connection/{vpn_connection_id}/ipsec-sa"
		product = "vpn"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{vpn_connection_id}", d.Get("vpn_connection_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving VPN connection IPSec SA: %s", err)
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

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("sa_infos", flattenGetIpsecSaBody(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetIpsecSaBody(resp interface{}) []map[string]interface{} {
	jobsJson := utils.PathSearch("sa_infos", resp, make([]interface{}, 0))
	jobsArray := jobsJson.([]interface{})
	if len(jobsArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(jobsArray))
	for _, saInfo := range jobsArray {
		result = append(result, map[string]interface{}{
			"id":             utils.PathSearch("id", saInfo, nil),
			"source_ip_cidr": utils.PathSearch("source_ip_cidr", saInfo, nil),
			"dest_ip_cidr":   utils.PathSearch("dest_ip_cidr", saInfo, nil),
			"packets_sent":   utils.PathSearch("packets_sent", saInfo, nil),
			"packets_recv":   utils.PathSearch("packets_recv", saInfo, nil),
			"traffic_sent":   utils.PathSearch("traffic_sent", saInfo, nil),
			"traffic_recv":   utils.PathSearch("traffic_recv", saInfo, nil),
			"collected_at":   utils.PathSearch("collected_at", saInfo, nil),
		})
	}

	return result
}
