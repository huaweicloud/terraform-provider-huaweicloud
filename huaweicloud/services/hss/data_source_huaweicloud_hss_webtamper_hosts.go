package hss

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	hssv5model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/hss/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS GET /v5/{project_id}/webtamper/hosts
func DataSourceWebTamperHosts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWebTamperHostsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protect_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rasp_protect_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_bit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protect_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rasp_protect_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"anti_tampering_times": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"detect_tampering_times": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceWebTamperHostsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		epsId       = cfg.GetEnterpriseProjectID(d, "all_granted_eps")
		limit       = int32(20)
		offset      int32
		allWtpHosts []hssv5model.WtpProtectHostResponseInfo
	)

	client, err := cfg.HcHssV5Client(region)
	if err != nil {
		return diag.Errorf("error creating HSS v5 client: %s", err)
	}

	for {
		request := hssv5model.ListWtpProtectHostRequest{
			Region:              region,
			EnterpriseProjectId: utils.String(epsId),
			HostName:            utils.StringIgnoreEmpty(d.Get("name").(string)),
			HostId:              utils.StringIgnoreEmpty(d.Get("host_id").(string)),
			PublicIp:            utils.StringIgnoreEmpty(d.Get("public_ip").(string)),
			PrivateIp:           utils.StringIgnoreEmpty(d.Get("private_ip").(string)),
			GroupName:           utils.StringIgnoreEmpty(d.Get("group_name").(string)),
			OsType:              utils.StringIgnoreEmpty(d.Get("os_type").(string)),
			ProtectStatus:       utils.StringIgnoreEmpty(d.Get("protect_status").(string)),
			Limit:               utils.Int32(limit),
			Offset:              utils.Int32(offset),
		}

		listResp, listErr := client.ListWtpProtectHost(&request)
		if listErr != nil {
			return diag.Errorf("error querying HSS web tamper hosts: %s", listErr)
		}

		if listResp == nil || listResp.DataList == nil || listResp.TotalNum == nil {
			break
		}
		if len(*listResp.DataList) == 0 {
			break
		}
		allWtpHosts = append(allWtpHosts, *listResp.DataList...)

		if int(*listResp.TotalNum) == len(allWtpHosts) {
			break
		}

		offset += limit
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuId)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("hosts", flattenWebTamperHosts(filterWebTamperHosts(allWtpHosts, d))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterWebTamperHosts(wtpHosts []hssv5model.WtpProtectHostResponseInfo, d *schema.ResourceData) []hssv5model.WtpProtectHostResponseInfo {
	if len(wtpHosts) == 0 {
		return nil
	}

	rst := make([]hssv5model.WtpProtectHostResponseInfo, 0, len(wtpHosts))
	for _, v := range wtpHosts {
		if raspProtectStatus, ok := d.GetOk("rasp_protect_status"); ok &&
			fmt.Sprint(raspProtectStatus) != utils.StringValue(v.RaspProtectStatus) {
			continue
		}
		rst = append(rst, v)
	}

	return rst
}

func flattenWebTamperHosts(wtpHosts []hssv5model.WtpProtectHostResponseInfo) []interface{} {
	if len(wtpHosts) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(wtpHosts))
	for _, v := range wtpHosts {
		rst = append(rst, map[string]interface{}{
			"id":                     v.HostId,
			"name":                   v.HostName,
			"public_ip":              v.PublicIp,
			"private_ip":             v.PrivateIp,
			"group_name":             v.GroupName,
			"os_bit":                 v.OsBit,
			"os_type":                v.OsType,
			"protect_status":         v.ProtectStatus,
			"rasp_protect_status":    v.RaspProtectStatus,
			"anti_tampering_times":   v.AntiTamperingTimes,
			"detect_tampering_times": v.DetectTamperingTimes,
		})
	}

	return rst
}
