package dew

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DEW GET /v1/{project_id}/csms/events
func DataSourceDewCsmsEvents() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDewCsmsEventsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"events": {
				Type:     schema.TypeList,
				Elem:     eventsSchema(),
				Computed: true,
			},
		},
	}
}

func eventsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"event_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"event_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
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
			"notification": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     notificationSchema(),
			},
		},
	}
	return &sc
}

func notificationSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"target_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceDewCsmsEventsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listEventsHttpUrl = "v1/{project_id}/csms/events"
		listEventsProduct = "kms"
	)
	listEventsClient, err := cfg.NewServiceClient(listEventsProduct, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	listEventsPath := listEventsClient.Endpoint + listEventsHttpUrl
	listEventsPath = strings.ReplaceAll(listEventsPath, "{project_id}", listEventsClient.ProjectID)

	listEventsResp, err := pagination.ListAllItems(
		listEventsClient,
		"marker",
		listEventsPath,
		&pagination.QueryOpts{MarkerField: "name"})

	if err != nil {
		return diag.Errorf("error retrieving CSMS events: %s", err)
	}

	listEventsRespJson, err := json.Marshal(listEventsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listEventsRespBody interface{}
	err = json.Unmarshal(listEventsRespJson, &listEventsRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("events", filterListEventsBody(
			flattenListEventsBody(listEventsRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListEventsBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("events", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		createAt := utils.PathSearch("create_time", v, float64(0))
		updateAt := utils.PathSearch("update_time", v, float64(0))
		rst = append(rst, map[string]interface{}{
			"name":         utils.PathSearch("name", v, nil),
			"event_id":     utils.PathSearch("event_id", v, nil),
			"event_types":  utils.PathSearch("event_types", v, nil),
			"status":       utils.PathSearch("state", v, nil),
			"created_at":   utils.FormatTimeStampRFC3339(int64(createAt.(float64))/1000, false),
			"updated_at":   utils.FormatTimeStampRFC3339(int64(updateAt.(float64))/1000, false),
			"notification": flattenNotification(v),
		})
	}
	return rst
}

func filterListEventsBody(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))

	for _, v := range all {
		if param, ok := d.GetOk("name"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("name", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("event_id"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("event_id", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("status"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("status", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}

	return rst
}

func flattenNotification(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("notification", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"target_type": utils.PathSearch("target_type", curJson, nil),
			"target_id":   utils.PathSearch("target_id", curJson, nil),
			"target_name": utils.PathSearch("target_name", curJson, nil),
		},
	}
	return rst
}
