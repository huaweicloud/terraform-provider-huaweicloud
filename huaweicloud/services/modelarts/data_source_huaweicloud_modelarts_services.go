// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ModelArts
// ---------------------------------------------------------------

package modelarts

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v1/{project_id}/services
func DataSourceServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceModelartServicesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"service_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Service ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Service name.`,
			},
			"model_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The model ID which the service used.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The workspace ID to which a service belongs.`,
			},
			"infer_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Inference mode of the service.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Service status.`,
			},
			"services": {
				Type:        schema.TypeList,
				Elem:        modelartServicesServicesSchema(),
				Computed:    true,
				Description: `The list of services.`,
			},
		},
	}
}

func modelartServicesServicesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Services ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Services name.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The workspace ID to which a service belongs.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the service.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Service status.`,
			},
			"infer_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Inference mode of the service.`,
			},
			"is_free": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether a free-of-charge flavor is used.`,
			},
			"schedule": {
				Type:        schema.TypeList,
				Elem:        modelartServicesServicesScheduleSchema(),
				Computed:    true,
				Description: `Service scheduling configuration, which can be configured only for real-time services.`,
			},
			"additional_properties": {
				Type:     schema.TypeList,
				Elem:     modelartServicesServicesAdditionalPropertySchema(),
				Computed: true,
			},
			"invocation_times": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Number of service calls.`,
			},
			"failed_times": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Number of failed service calls.`,
			},
			"is_shared": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether a service is subscribed.`,
			},
			"shared_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Number of subscribed services.`,
			},
			"is_opened_sample_collection": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable data collection, which defaults to false.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `User to which a service belongs.`,
			},
		},
	}
	return &sc
}

func modelartServicesServicesScheduleSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"duration": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Value mapping a time unit.`,
			},
			"time_unit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Scheduling time unit. Possible values are DAYS, HOURS, and MINUTES.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Scheduling type. Only the value **stop** is supported.`,
			},
		},
	}
	return &sc
}

func modelartServicesServicesAdditionalPropertySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"smn_notification": {
				Type:     schema.TypeList,
				Elem:     modelartServicesAdditionalPropertySmnNotificationSchema(),
				Computed: true,
			},
			"log_report_channels": {
				Type:     schema.TypeList,
				Elem:     modelartServicesAdditionalPropertyLogReportChannelSchema(),
				Computed: true,
			},
		},
	}
	return &sc
}

func modelartServicesAdditionalPropertySmnNotificationSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"topic_urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `URN of an SMN topic.`,
			},
			"events": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Computed:    true,
				Description: `Event ID.`,
			},
		},
	}
	return &sc
}

func modelartServicesAdditionalPropertyLogReportChannelSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of log report channel. The valid value is **LTS**.`,
			},
		},
	}
	return &sc
}

func resourceModelartServicesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listServices: Query the list of ModelArts services
	var (
		listServicesHttpUrl = "v1/{project_id}/services"
		listServicesProduct = "modelarts"
	)
	listServicesClient, err := cfg.NewServiceClient(listServicesProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	listServicesPath := listServicesClient.Endpoint + listServicesHttpUrl
	listServicesPath = strings.ReplaceAll(listServicesPath, "{project_id}", listServicesClient.ProjectID)

	listServicesqueryParams := buildListServicesQueryParams(d)
	listServicesPath += listServicesqueryParams

	listServicesResp, err := pagination.ListAllItems(
		listServicesClient,
		"offset",
		listServicesPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Modelart services")
	}

	listServicesRespJson, err := json.Marshal(listServicesResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listServicesRespBody interface{}
	err = json.Unmarshal(listServicesRespJson, &listServicesRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("services", flattenListServicesservices(listServicesRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListServicesservices(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("services", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                          utils.PathSearch("service_id", v, nil),
			"name":                        utils.PathSearch("service_name", v, nil),
			"workspace_id":                utils.PathSearch("workspace_id", v, nil),
			"description":                 utils.PathSearch("description", v, nil),
			"status":                      utils.PathSearch("status", v, nil),
			"infer_type":                  utils.PathSearch("infer_type", v, nil),
			"is_free":                     utils.PathSearch("is_free", v, nil),
			"schedule":                    flattenServicesSchedule(v),
			"additional_properties":       flattenServicesAdditionalProperty(v),
			"invocation_times":            utils.PathSearch("invocation_times", v, nil),
			"failed_times":                utils.PathSearch("failed_times", v, nil),
			"is_shared":                   utils.PathSearch("is_shared", v, nil),
			"shared_count":                utils.PathSearch("shared_count", v, nil),
			"is_opened_sample_collection": utils.PathSearch("is_opened_sample_collection", v, nil),
			"owner":                       utils.PathSearch("owner", v, nil),
		})
	}
	return rst
}

func flattenServicesSchedule(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("schedule", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"duration":  utils.PathSearch("duration", v, nil),
			"time_unit": utils.PathSearch("time_unit", v, nil),
			"type":      utils.PathSearch("type", v, nil),
		})
	}
	return rst
}

func flattenServicesAdditionalProperty(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("additional_properties", resp, nil)
	if curJson == nil {
		log.Printf("[ERROR] error parsing additional_properties from response")
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"smn_notification":    flattenServicesAdditionalPropertySmnNotification(curJson),
			"log_report_channels": flattenServicesAdditionalPropertyLogReportChannels(curJson),
		},
	}
	return rst
}

func flattenServicesAdditionalPropertySmnNotification(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("smn_notification", resp, nil)
	if curJson == nil {
		log.Printf("[ERROR] error parsing additional_properties from response")
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"topic_urn": utils.PathSearch("topic_urn", curJson, nil),
			"events":    utils.PathSearch("events", curJson, nil),
		},
	}
	return rst
}

func flattenServicesAdditionalPropertyLogReportChannels(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("log_report_channels", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"type": utils.PathSearch("type", v, nil),
		})
	}
	return rst
}

func buildListServicesQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("service_id"); ok {
		res = fmt.Sprintf("%s&service_id=%v", res, v)
	}

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&service_name=%v", res, v)
	}

	if v, ok := d.GetOk("model_id"); ok {
		res = fmt.Sprintf("%s&env=%v", res, v)
	}

	if v, ok := d.GetOk("workspace_id"); ok {
		res = fmt.Sprintf("%s&workspace_id=%v", res, v)
	}

	if v, ok := d.GetOk("infer_type"); ok {
		res = fmt.Sprintf("%s&infer_type=%v", res, v)
	}

	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
