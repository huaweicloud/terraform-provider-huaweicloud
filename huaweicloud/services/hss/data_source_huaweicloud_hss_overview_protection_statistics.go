package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS GET /v5/{project_id}/overview/protection/statistics
func DataSourceOverviewProtectionStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOverviewProtectionStatisticsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vul_library_update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"protect_days": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"threat_library_update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vul_detected_total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"baseline_detected_total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"finger_scan_total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"alarm_detected_total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ransomware_alarm_detected_total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"file_alarm_detected_total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rasp_alarm_detected_total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"wtp_alarm_detected_total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"image_risk_total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"container_alarm_total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"container_firewall_policy_total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"auto_kill_virus_status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func buildOverviewProtectionStatisticsQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return ""
}

func dataSourceOverviewProtectionStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/overview/protection/statistics"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildOverviewProtectionStatisticsQueryParams(epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS overview protection statistics: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("vul_library_update_time", utils.PathSearch("vul_library_update_time", respBody, nil)),
		d.Set("protect_days", utils.PathSearch("protect_days", respBody, nil)),
		d.Set("threat_library_update_time", utils.PathSearch("threat_library_update_time", respBody, nil)),
		d.Set("vul_detected_total_num", utils.PathSearch("vul_detected_total_num", respBody, nil)),
		d.Set("baseline_detected_total_num", utils.PathSearch("baseline_detected_total_num", respBody, nil)),
		d.Set("finger_scan_total_num", utils.PathSearch("finger_scan_total_num", respBody, nil)),
		d.Set("alarm_detected_total_num", utils.PathSearch("alarm_detected_total_num", respBody, nil)),
		d.Set("ransomware_alarm_detected_total_num",
			utils.PathSearch("ransomware_alarm_detected_total_num", respBody, nil)),
		d.Set("file_alarm_detected_total_num", utils.PathSearch("file_alarm_detected_total_num", respBody, nil)),
		d.Set("rasp_alarm_detected_total_num", utils.PathSearch("rasp_alarm_detected_total_num", respBody, nil)),
		d.Set("wtp_alarm_detected_total_num", utils.PathSearch("wtp_alarm_detected_total_num", respBody, nil)),
		d.Set("image_risk_total_num", utils.PathSearch("image_risk_total_num", respBody, nil)),
		d.Set("container_alarm_total_num", utils.PathSearch("container_alarm_total_num", respBody, nil)),
		d.Set("container_firewall_policy_total_num",
			utils.PathSearch("container_firewall_policy_total_num", respBody, nil)),
		d.Set("auto_kill_virus_status", utils.PathSearch("auto_kill_virus_status", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
