package lts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LTS POST /v2/{project_id}/lts/log2metric/rules
// @API LTS GET /v2/{project_id}/lts/log2metric/rules/{rule_id}
// @API LTS PUT /v2/{project_id}/lts/log2metric/rules/{rule_id}
// @API LTS DELETE /v2/{project_id}/lts/log2metric/rules/{rule_id}
func ResourceMetricRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMetricRuleCreate,
		ReadContext:   resourceMetricRuleRead,
		UpdateContext: resourceMetricRuleUpdate,
		DeleteContext: resourceMetricRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the log metric rule.`,
			},
			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The status of the log metric rule.`,
			},
			"log_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The log group ID to which the log metric rule belongs.`,
			},
			"log_stream_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The log stream ID to which the log metric rule belongs.`,
			},
			"sampler": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The type of the log sampling.`,
						},
						"ratio": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The sampling rate of the log.`,
						},
					},
				},
				Description: `The sampling configuration of the log.`,
			},
			"sinks": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The type of the stored object.`,
						},
						"metric_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the generated log metric.`,
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The name of the AOM Prometheus common instance.`,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The ID of the AOM Prometheus common instance.`,
						},
					},
				},
				Description: `The storage location of the generated metrics.`,
			},
			"aggregator": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The type of the log statistics.`,
						},
						"field": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The field of the log statistics.`,
						},
						"group_by": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of the group fields of the log statistics.`,
						},
						"keyword": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The keyword of the log statistics.`,
						},
					},
				},
				Description: `The configuration of log statistics mode.`,
			},
			"window_size": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The interval time for processing data windows.`,
			},
			"report": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to report data to sinks.`,
			},
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							RequiredWith: []string{"filter.0.filters"},
							Description:  `The filter type of the log.`,
						},
						"filters": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        metricRuleFiltersSchema(),
							Description: `The list of log filtering rule groups.`,
						},
					},
				},
				DiffSuppressFunc: func(_, _, _ string, d *schema.ResourceData) bool {
					oldRaw, newRaw := d.GetChange("filter")
					// If filter is set to {}, the DiffSuppress function is needed to prevent changes.
					return len(oldRaw.([]interface{})) == 0 && newRaw.([]interface{})[0] == nil
				},
				Description: `The configuration of log filtering rule.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the log metric rule.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the log metric rule, in RFC3339 format.`,
			},
		},
	}
}

func metricRuleFiltersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The filter type of the log.`,
			},
			"filters": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The filter field of the log.`,
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The filter conditions of the log.`,
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The value corresponding to the log filter field.`,
						},
						"lower": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The minimum value corresponding to the log filter field.`,
						},
						"upper": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The maximum value corresponding to the log filter field.`,
						},
					},
				},
				Description: `The list of the log filter rule associations.`,
			},
		},
	}
}

func buildMetricRuleBodyParams(d *schema.ResourceData, domainId, projectId string) map[string]interface{} {
	return map[string]interface{}{
		"domain_id":     domainId,
		"project_id":    projectId,
		"name":          d.Get("name"),
		"status":        d.Get("status"),
		"log_group_id":  d.Get("log_group_id"),
		"log_stream_id": d.Get("log_stream_id"),
		"sampler":       buildMetricRuleSampler(d.Get("sampler.0")),
		"report":        d.Get("report"),
		"sinks":         buildMetricRuleSinks(d.Get("sinks").(*schema.Set).List()),
		"aggregator":    buildMetricRuleAggregator(d.Get("aggregator.0")),
		"window_size":   d.Get("window_size"),
		"filter":        buildMetricRuleFilter(d.Get("filter").([]interface{})),
		"description":   d.Get("description"),
	}
}

func buildMetricRuleSampler(sampler interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":  utils.PathSearch("type", sampler, nil),
		"ratio": utils.PathSearch("ratio", sampler, nil),
	}
}

func buildMetricRuleSinks(sinks []interface{}) []map[string]interface{} {
	rest := make([]map[string]interface{}, len(sinks))
	for i, item := range sinks {
		rest[i] = map[string]interface{}{
			"type":        utils.PathSearch("type", item, nil),
			"metric_name": utils.PathSearch("metric_name", item, nil),
			"name":        utils.PathSearch("name", item, nil),
			"instance":    utils.PathSearch("instance_id", item, nil),
		}
	}
	return rest
}

func buildMetricRuleAggregator(aggregator interface{}) map[string]interface{} {
	params := map[string]interface{}{
		"type":     utils.PathSearch("type", aggregator, nil),
		"field":    utils.PathSearch("field", aggregator, nil),
		"group_by": utils.ValueIgnoreEmpty(utils.PathSearch("group_by", aggregator, nil)),
		"keyword":  utils.ValueIgnoreEmpty(utils.PathSearch("keyword", aggregator, nil)),
	}
	return utils.RemoveNil(params)
}

func buildMetricRuleFilter(filter []interface{}) map[string]interface{} {
	if len(filter) == 0 || filter[0] == nil {
		return nil
	}
	return map[string]interface{}{
		"type":    utils.PathSearch("type", filter[0], nil),
		"filters": buildMetricRuleFilterRules(utils.PathSearch("filters", filter[0], schema.NewSet(schema.HashString, nil)).(*schema.Set)),
	}
}

func buildMetricRuleFilterRules(filterRules *schema.Set) []map[string]interface{} {
	rst := make([]map[string]interface{}, 0)
	for _, v := range filterRules.List() {
		filters := utils.PathSearch("filters", v, schema.NewSet(schema.HashString, nil)).(*schema.Set)
		if filters.Len() > 0 {
			rst = append(rst, map[string]interface{}{
				"type":    utils.PathSearch("type", v, nil),
				"filters": buildMetricRuleFilters(utils.PathSearch("filters", v, schema.NewSet(schema.HashString, nil)).(*schema.Set)),
			})
		}
	}
	return rst
}

func buildMetricRuleFilters(filters *schema.Set) []map[string]interface{} {
	if filters.Len() == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, filters.Len())
	for i, v := range filters.List() {
		params := map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"type":  utils.PathSearch("type", v, nil),
			"value": utils.ValueIgnoreEmpty(utils.PathSearch("value", v, nil)),
			"lower": utils.ValueIgnoreEmpty(utils.PathSearch("lower", v, nil)),
			"upper": utils.ValueIgnoreEmpty(utils.PathSearch("upper", v, nil)),
		}
		rst[i] = utils.RemoveNil(params)
	}
	return rst
}

func resourceMetricRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/lts/log2metric/rules"
	)

	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}
	projectId := client.ProjectID
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", projectId)

	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildMetricRuleBodyParams(d, cfg.DomainID, projectId),
	}

	requestResp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("error creating log metric rule: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleId := utils.PathSearch("rule_id", respBody, "").(string)
	if ruleId == "" {
		return diag.Errorf("unable to find the log metric rule ID from the API response")
	}

	d.SetId(ruleId)

	return resourceMetricRuleRead(ctx, d, meta)
}

// GetMetricRuleById is a method used to get metric rule detail by rule ID.
func GetMetricRuleById(client *golangsdk.ServiceClient, ruleId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/lts/log2metric/rules/{rule_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{rule_id}", ruleId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func resourceMetricRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		ruleId = d.Id()
	)

	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	respBody, err := GetMetricRuleById(client, ruleId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving log metric rule (%s)", ruleId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("log_group_id", utils.PathSearch("log_group_id", respBody, nil)),
		d.Set("log_stream_id", utils.PathSearch("log_stream_id", respBody, nil)),
		d.Set("sampler", flattenMetricRuleSampler(utils.PathSearch("sampler", respBody, nil))),
		d.Set("report", utils.PathSearch("report", respBody, nil)),
		d.Set("sinks", flattenMetricRuleSinks(utils.PathSearch("sinks", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("aggregator", flattenMetricRuleAggregator(utils.PathSearch("aggregator", respBody, nil))),
		d.Set("window_size", utils.PathSearch("window_size", respBody, nil)),
		d.Set("filter", flattenMetricRuleFilter(utils.PathSearch("filter", respBody, nil))),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", respBody, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMetricRuleSampler(sampler interface{}) []map[string]interface{} {
	if sampler == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"type":  utils.PathSearch("type", sampler, nil),
			"ratio": utils.PathSearch("ratio", sampler, nil),
		},
	}
}

func flattenMetricRuleSinks(sinks []interface{}) []map[string]interface{} {
	if len(sinks) == 0 {
		return nil
	}

	rest := make([]map[string]interface{}, len(sinks))
	for i, item := range sinks {
		rest[i] = map[string]interface{}{
			"type":        utils.PathSearch("type", item, nil),
			"metric_name": utils.PathSearch("metric_name", item, nil),
			"name":        utils.PathSearch("name", item, nil),
			"instance_id": utils.PathSearch("instance", item, nil),
		}
	}
	return rest
}

func flattenMetricRuleAggregator(aggregator interface{}) []map[string]interface{} {
	if aggregator == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"type":     utils.PathSearch("type", aggregator, nil),
			"field":    utils.PathSearch("field", aggregator, nil),
			"group_by": utils.PathSearch("group_by", aggregator, nil),
			"keyword":  utils.PathSearch("keyword", aggregator, nil),
		},
	}
}

func flattenMetricRuleFilter(filter interface{}) []map[string]interface{} {
	if filter == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"type":    utils.PathSearch("type", filter, nil),
			"filters": flattenmetricRuleFliter(utils.PathSearch("filters", filter, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenmetricRuleFliter(filters []interface{}) []map[string]interface{} {
	if len(filters) == 0 {
		// The outermost layer is an object, and the `filter.filters` parameter and sub parameters are optional. so a default value
		// needs to be set to prevent change.
		return []map[string]interface{}{
			{
				"type":    nil,
				"filters": nil,
			},
		}
	}

	rest := make([]map[string]interface{}, len(filters))
	for i, v := range filters {
		rest[i] = map[string]interface{}{
			"type":    utils.PathSearch("type", v, nil),
			"filters": flattenmetricRuleAssociatedFliter(utils.PathSearch("filters", v, make([]interface{}, 0)).([]interface{})),
		}
	}
	return rest
}

func flattenmetricRuleAssociatedFliter(associatedFilters []interface{}) []map[string]interface{} {
	if len(associatedFilters) == 0 {
		return nil
	}
	rest := make([]map[string]interface{}, len(associatedFilters))
	for i, v := range associatedFilters {
		rest[i] = map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"type":  utils.PathSearch("type", v, nil),
			"value": utils.PathSearch("value", v, nil),
			"lower": utils.PathSearch("lower", v, nil),
			"upper": utils.PathSearch("upper", v, nil),
		}
	}
	return rest
}

func resourceMetricRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/lts/log2metric/rules/{rule_id}"
		ruleId  = d.Id()
	)

	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	projectId := client.ProjectID
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{rule_id}", ruleId)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildMetricRuleBodyParams(d, cfg.DomainID, projectId),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating log metric rule (%s): %s", ruleId, err)
	}

	return resourceMetricRuleRead(ctx, d, meta)
}

func resourceMetricRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/lts/log2metric/rules/{rule_id}"
		ruleId  = d.Id()
	)

	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{rule_id}", ruleId)

	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting log metric rule (%s)", ruleId))
	}
	return nil
}
