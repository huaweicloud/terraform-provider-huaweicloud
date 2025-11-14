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

// @API FunctionGraph GET /v2/{project_id}/fgs/triggers/{function_urn}
// @API FunctionGraph GET /v2/{project_id}/fgs/triggers
func DataSourceFunctionTriggers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFunctionTriggersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the triggers are located.`,
			},
			"function_urn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The URN of the function URN to which the triggers belong.`,
			},
			"trigger_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the function trigger.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the function trigger.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The status of the function trigger.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The start time of creation time of the function trigger.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The end time of creation time of the function trigger.`,
			},
			"triggers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the function trigger.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the function trigger.`,
						},
						"event_data": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The detailed configuration of the function trigger, in JSON format.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The current status of the function trigger.`,
						},
						"function_urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The URN of the function URN to which the triggers belong.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the function trigger, in RFC3339 format.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the function trigger, in RFC3339 format.`,
						},
					},
				},
				Description: `All triggers that match the filter parameters.`,
			},
		},
	}
}

func listTriggersWithFunctionUrn(client *golangsdk.ServiceClient, functionUrn string) ([]interface{}, error) {
	httpUrl := "v2/{project_id}/fgs/triggers/{function_urn}"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{function_urn}", functionUrn)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	resp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}
	return utils.PathSearch("[]", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func listTriggers(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/fgs/triggers"
		marker  = ""
		limit   = 500
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = fmt.Sprintf("%s?maxitems=%v", listPath, limit)

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

func dataSourceFunctionTriggersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		functionUrn = d.Get("function_urn").(string)
	)
	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	var triggers []interface{}
	if functionUrn != "" {
		triggers, err = listTriggersWithFunctionUrn(client, functionUrn)
		if err != nil {
			return diag.Errorf("error retrieving triggers under specified function (%s): %s", functionUrn, err)
		}
	} else {
		triggers, err = listTriggers(client)
		if err != nil {
			return diag.Errorf("error querying triggers: %s", err)
		}
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("triggers", filterTriggers(flattenTriggers(triggers), d)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func filterTriggers(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))

	for _, v := range all {
		if param, ok := d.GetOk("trigger_id"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("id", v, nil)) {
			continue
		}

		if param, ok := d.GetOk("type"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("type", v, nil)) {
			continue
		}

		if param, ok := d.GetOk("status"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("status", v, nil)) {
			continue
		}

		createdAt := utils.PathSearch("created_at", v, "").(string)
		// Some triggers do not return creation time, such as: "SMN".
		if createdAt == "" {
			continue
		}

		createdAtTimestamp := utils.ConvertTimeStrToNanoTimestamp(createdAt)
		if param, ok := d.GetOk("start_time"); ok &&
			utils.ConvertTimeStrToNanoTimestamp(param.(string)) > createdAtTimestamp {
			continue
		}

		if param, ok := d.GetOk("end_time"); ok &&
			utils.ConvertTimeStrToNanoTimestamp(param.(string)) < createdAtTimestamp {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func flattenTriggers(triggers []interface{}) []interface{} {
	if len(triggers) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(triggers))
	for _, trigger := range triggers {
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("trigger_id", trigger, nil),
			"type":         utils.PathSearch("trigger_type_code", trigger, nil),
			"status":       utils.PathSearch("trigger_status", trigger, nil),
			"event_data":   utils.JsonToString(utils.PathSearch("event_data", trigger, nil)),
			"function_urn": utils.PathSearch("func_urn", trigger, nil),
			"created_at":   convertTimeToRFC339(utils.PathSearch("created_time", trigger, "").(string)),
			"updated_at":   convertTimeToRFC339(utils.PathSearch("last_updated_time", trigger, "").(string)),
		})
	}
	return result
}

// The timeStr format is "yyyy-MM-ddTHH:mm:ss+08:00".
// Formats time according to the local computer time.
func convertTimeToRFC339(timeStr string) string {
	return utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(timeStr)/1000, false)
}
