package fgs

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API FunctionGraph GET /v2/{project_id}/fgs/triggers
func DataSourceTriggers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTriggersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the triggers are located.`,
			},
			"trigger_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the trigger.`,
			},
			"triggers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trigger_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the trigger.`,
						},
						"trigger_type_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the trigger.`,
						},
						"trigger_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the trigger.`,
						},
						"event_data": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The detailed configuration of the trigger, in JSON format.`,
						},
						"func_urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The URN of the function to which the trigger belongs.`,
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the trigger.`,
						},
						"last_updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the trigger.`,
						},
					},
				},
				Description: `The list of triggers that match the filter parameters.`,
			},
		},
	}
}

func getTriggers(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/fgs/triggers"
		marker  = ""
		limit   = 500
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = fmt.Sprintf("%s?maxitems=%v", listPath, limit)
	if triggerType, ok := d.GetOk("trigger_type"); ok {
		listPath = fmt.Sprintf("%s&trigger_type=%v", listPath, triggerType)
	}

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s&marker=%s", listPath, marker)
		}

		requestResp, err := client.Request("GET", listPathWithMarker, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		triggers := utils.PathSearch("triggers", respBody, make([]interface{}, 0)).([]interface{})
		if len(triggers) < 1 {
			break
		}

		result = append(result, triggers...)
		if len(triggers) < limit {
			break
		}
		marker = strconv.Itoa(int(utils.PathSearch("next_marker", respBody, float64(0)).(float64)))
	}

	return result, nil
}

func dataSourceTriggersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	triggers, err := getTriggers(client, d)
	if err != nil {
		return diag.Errorf("error querying triggers: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("triggers", flattenTriggersList(triggers)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTriggersList(triggers []interface{}) []map[string]interface{} {
	if len(triggers) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(triggers))
	for _, trigger := range triggers {
		result = append(result, map[string]interface{}{
			"trigger_id":        utils.PathSearch("trigger_id", trigger, nil),
			"trigger_type_code": utils.PathSearch("trigger_type_code", trigger, nil),
			"trigger_status":    utils.PathSearch("trigger_status", trigger, nil),
			"event_data":        utils.JsonToString(utils.PathSearch("event_data", trigger, nil)),
			"func_urn":          utils.PathSearch("func_urn", trigger, nil),
			"created_time":      utils.PathSearch("created_time", trigger, nil),
			"last_updated_time": utils.PathSearch("last_updated_time", trigger, nil),
		})
	}
	return result
}
