package dew

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DEW GET /v1/{project_id}/csms/notification-records
func DataSourceNotificationRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNotificationRecordsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_event_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notification_target_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notification_target_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notification_content": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notification_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNotificationRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/csms/notification-records"
	)

	client, err := cfg.NewServiceClient("kms", region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving notification records: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationRecords := utils.PathSearch("records", getRespBody, make([]interface{}, 0)).([]interface{})

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("records", flattenNotificationRecords(notificationRecords)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenNotificationRecords(records []interface{}) []interface{} {
	if len(records) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(records))
	for _, v := range records {
		rst = append(rst, map[string]interface{}{
			"event_name":               utils.PathSearch("event_name", v, nil),
			"trigger_event_type":       utils.PathSearch("trigger_event_type", v, nil),
			"secret_name":              utils.PathSearch("secret_name", v, nil),
			"secret_type":              utils.PathSearch("secret_type", v, nil),
			"notification_target_name": utils.PathSearch("notification_target_name", v, nil),
			"notification_target_id":   utils.PathSearch("notification_target_id", v, nil),
			"notification_content":     utils.PathSearch("notification_content", v, nil),
			"notification_status":      utils.PathSearch("notification_status", v, nil),
			"create_time":              utils.PathSearch("create_time", v, nil),
		})
	}

	return rst
}
