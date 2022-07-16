package sms

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/sms/v3/sources"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// DataSourceServers is the impl of data/huaweicloud_sms_source_servers
func DataSourceServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServersRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connected": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"agent_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"registered_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vcpus": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"device_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceServersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	smsClient, err := config.SmsV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	listOpts := sources.ListOpts{
		Id:    d.Get("id").(string),
		Ip:    d.Get("ip").(string),
		Name:  d.Get("name").(string),
		State: d.Get("state").(string),
	}

	log.Printf("[DEBUG] filtering SMS source servers by %#v", listOpts)
	allServers, err := sources.List(smsClient, listOpts)
	if err != nil {
		return diag.Errorf("unable to list source servers: %s ", err)
	}

	ids := make([]string, len(allServers))
	stateServers := make([]map[string]interface{}, len(allServers))

	for i, item := range allServers {
		ids[i] = item.Id
		stateServers[i] = flattenSourceServer(item)
	}

	if len(ids) == 1 {
		d.SetId(ids[0])
	} else {
		d.SetId(hashcode.Strings(ids))
	}

	if err := d.Set("servers", stateServers); err != nil {
		diag.Errorf("error setting SMS source servers: %s", err)
	}

	return nil
}

func flattenSourceServer(server sources.SourceServer) map[string]interface{} {
	disks := make([]map[string]interface{}, len(server.InitTargetServer.Disks))
	for i, d := range server.InitTargetServer.Disks {
		disks[i] = map[string]interface{}{
			"device_type": d.DeviceUse,
			"name":        d.Name,
			"size":        convertBytestoMB(d.Size),
		}
	}

	return map[string]interface{}{
		"id":              server.Id,
		"ip":              server.Ip,
		"name":            server.Name,
		"state":           server.State,
		"connected":       server.Connected,
		"agent_version":   server.AgentVersion,
		"os_type":         server.OsType,
		"os_version":      server.OsVersion,
		"vcpus":           server.CPU,
		"memory":          convertBytestoMB(server.Memory),
		"disks":           disks,
		"registered_time": utils.FormatTimeStampUTC(server.AddDate / 1000),
	}
}
