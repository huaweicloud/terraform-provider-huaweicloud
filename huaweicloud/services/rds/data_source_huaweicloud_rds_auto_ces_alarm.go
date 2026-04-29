package rds

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

// @API RDS GET /v3/{project_id}/auto-ces-alarm
func DataSourceRdsAutoCesAlarm() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsAutoCesAlarmRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     autoCesAlarmEntitiesSchema(),
			},
		},
	}
}

func autoCesAlarmEntitiesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"new_instance_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"switch_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alarm_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"topic_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceRdsAutoCesAlarmRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	httpUrl := "v3/{project_id}/auto-ces-alarm"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildGetAutoCesAlarmQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving RDS auto CES alarm: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("entities", flattenGetAutoCesAlarmBody(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetAutoCesAlarmQueryParams(d *schema.ResourceData) string {
	if v, ok := d.GetOk("engine"); ok {
		return fmt.Sprintf("?engine=%v", v)
	}
	return ""
}

func flattenGetAutoCesAlarmBody(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	entities := utils.PathSearch("entities", resp, make([]interface{}, 0)).([]interface{})
	if len(entities) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(entities))
	for i, entity := range entities {
		result[i] = map[string]interface{}{
			"id":                   utils.PathSearch("id", entity, nil),
			"domain_id":            utils.PathSearch("domain_id", entity, nil),
			"domain_name":          utils.PathSearch("domain_name", entity, nil),
			"project_id":           utils.PathSearch("project_id", entity, nil),
			"project_name":         utils.PathSearch("project_name", entity, nil),
			"engine_name":          utils.PathSearch("engine_name", entity, nil),
			"new_instance_default": utils.PathSearch("new_instance_default", entity, nil),
			"switch_status":        utils.PathSearch("switch_status", entity, nil),
			"alarm_id":             utils.PathSearch("alarm_id", entity, nil),
			"topic_urn":            utils.PathSearch("topic_urn", entity, nil),
			"created_at":           utils.PathSearch("created_at", entity, nil),
			"updated_at":           utils.PathSearch("updated_at", entity, nil),
		}
	}

	return result
}
