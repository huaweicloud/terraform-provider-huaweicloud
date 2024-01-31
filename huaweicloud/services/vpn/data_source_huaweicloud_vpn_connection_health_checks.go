package vpn

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/vpn/v5/connection_monitors"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPN GET /v5/{project_id}/connection-monitors
func DataSourceVpnConnectionHealthChecks() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceConnectionHealthChecksRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connection_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connection_health_checks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"proto_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceConnectionHealthChecksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.VpnV5Client(region)

	if err != nil {
		return diag.Errorf("error creating vpn v5 client: %s", err)
	}

	body, err := connection_monitors.List(client).Extract()
	if err != nil {
		return diag.Errorf("error retrieving VPN connection health check: %s", err)
	}

	filter := map[string]interface{}{
		"Status":          d.Get("status").(string),
		"VpnConnectionId": d.Get("connection_id").(string),
		"SourceIp":        d.Get("source_ip").(string),
		"DestinationIp":   d.Get("destination_ip").(string),
	}

	filterConnectionHealthChecks, err := utils.FilterSliceWithField(body, filter)
	if err != nil {
		return diag.Errorf("filter VPN connection health check failed: %s", err)
	}

	var chk []map[string]interface{}
	for _, item := range filterConnectionHealthChecks {
		check := item.(connection_monitors.ConnectionMonitor)
		chk = append(chk, flattenGetVpnConnectionHealthCheck(check))
	}

	uuidStr, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuidStr)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("connection_health_checks", chk),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetVpnConnectionHealthCheck(ctm connection_monitors.ConnectionMonitor) map[string]interface{} {
	check := map[string]interface{}{
		"id":             ctm.ID,
		"status":         ctm.Status,
		"connection_id":  ctm.VpnConnectionId,
		"type":           ctm.Type,
		"source_ip":      ctm.SourceIp,
		"destination_ip": ctm.DestinationIp,
		"proto_type":     ctm.ProtoType,
	}
	return check
}
