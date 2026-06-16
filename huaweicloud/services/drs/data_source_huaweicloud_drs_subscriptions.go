package drs

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS POST /v5/{project_id}/subscriptions
func DataSourceDrsSubscriptions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsSubscriptionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"engine_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"net_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instance_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"begin_at": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_billing": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subscriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     subscriptionJobSchema(),
			},
		},
	}
}

func subscriptionJobSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"begin_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"now_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"job_action": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"available_actions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"unavailable_actions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"current_action": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildSubscriptionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"job_type":              "subscription",
		"engine_type":           d.Get("engine_type"),
		"net_type":              d.Get("net_type"),
		"name":                  d.Get("name"),
		"status":                d.Get("status"),
		"description":           d.Get("description"),
		"enterprise_project_id": d.Get("enterprise_project_id"),
		"instance_ids":          utils.ExpandToStringList(d.Get("instance_ids").([]interface{})),
		"instance_ip":           d.Get("instance_ip"),
		"sort_key":              d.Get("sort_key"),
		"sort_dir":              d.Get("sort_dir"),
		"service_name":          d.Get("service_name"),
		"begin_at":              d.Get("begin_at"),
		"tags":                  d.Get("tags"),
		"is_billing":            d.Get("is_billing"),
	}
	return utils.RemoveNil(bodyParams)
}

func dataSourceDrsSubscriptionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v5/{project_id}/subscriptions"
		result  = make([]interface{}, 0)
		limit   = 1000
		offset  = 0
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildSubscriptionBodyParams(d),
	}

	for {
		currentPath := fmt.Sprintf("%s?limit=%d&offset=%d", listPath, limit, offset)
		listResp, err := client.Request("POST", currentPath, &reqOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		jobs := utils.PathSearch("jobs", listRespBody, make([]interface{}, 0)).([]interface{})
		result = append(result, jobs...)
		if len(jobs) < limit {
			break
		}

		offset += len(jobs)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("subscriptions", flattenSubscriptions(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenJobAction(jobAction interface{}) []map[string]interface{} {
	if jobAction == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"available_actions":   utils.PathSearch("available_actions", jobAction, make([]interface{}, 0)),
			"unavailable_actions": utils.PathSearch("unavailable_actions", jobAction, make([]interface{}, 0)),
			"current_action":      utils.PathSearch("current_action", jobAction, nil),
		},
	}
}

func flattenSubscriptions(subscriptions []interface{}) []interface{} {
	if len(subscriptions) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		result = append(result, map[string]interface{}{
			"id":                    utils.PathSearch("id", subscription, nil),
			"name":                  utils.PathSearch("name", subscription, nil),
			"status":                utils.PathSearch("status", subscription, nil),
			"created_time":          utils.PathSearch("created_time", subscription, nil),
			"begin_time":            utils.PathSearch("begin_time", subscription, nil),
			"now_time":              utils.PathSearch("now_time", subscription, nil),
			"description":           utils.PathSearch("description", subscription, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", subscription, nil),
			"job_action":            flattenJobAction(utils.PathSearch("job_action", subscription, nil)),
		})
	}

	return result
}
