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

// @API HSS GET /v5/{project_id}/setting/alarm-config
func DataSourceSettingAlarmConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSettingAlarmConfigurationRead,

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
			"alarm_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"topic_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"daily_alarm": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"realtime_alarm": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"alarm_level": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ignored_event_class_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceSettingAlarmConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/setting/alarm-config"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	if epsId != "" {
		getPath = fmt.Sprintf("%s?enterprise_project_id=%s", getPath, epsId)
	}

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving the alarm configuration: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("alarm_type", utils.PathSearch("alarm_type", respBody, nil)),
		d.Set("display_name", utils.PathSearch("display_name", respBody, nil)),
		d.Set("topic_urn", utils.PathSearch("topic_urn", respBody, nil)),
		d.Set("daily_alarm", utils.PathSearch("daily_alarm", respBody, nil)),
		d.Set("realtime_alarm", utils.PathSearch("realtime_alarm", respBody, nil)),
		d.Set("alarm_level", utils.PathSearch("alarm_level", respBody, nil)),
		d.Set("ignored_event_class_list", utils.PathSearch("ignored_event_class_list", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
