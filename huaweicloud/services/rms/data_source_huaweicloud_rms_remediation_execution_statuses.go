package rms

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

// @API CONFIG POST /v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/remediation-execution-statuses/summary
func DataSourceRemediationExecutionStatuses() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRemediationExecutionStatusesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"policy_assignment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The policy assignment ID.`,
			},
			"resource_keys": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `The list of query criteria required to collect remediation results.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The resource type.`,
						},
						"resource_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The resource ID.`,
						},
						"resource_provider": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The cloud service name.`,
						},
					},
				},
			},
			"value": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The compliance rule remediation execution results.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_key": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The query criteria required to collect remediation results.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The resource type.`,
									},
									"resource_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The resource ID.`,
									},
									"resource_provider": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The cloud service name.`,
									},
								},
							},
						},
						"invocation_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The start time of remediation.`,
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The execution result of remediation.`,
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The information of remediation execution.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceRemediationExecutionStatusesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	product := "rms"
	var mErr *multierror.Error

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	statuses, err := getRemediationExecutionStatuses(client, d, cfg.DomainID)
	if err != nil {
		return diag.Errorf("error retrieving RMS remediation execution statuses: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(mErr,
		d.Set("region", region),
		d.Set("value", flattenValue(statuses)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func getRemediationExecutionStatuses(client *golangsdk.ServiceClient, d *schema.ResourceData, domainId string) (interface{}, error) {
	httpUrl := "v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/remediation-execution-statuses/summary"
	getStatusesPath := client.Endpoint + httpUrl
	getStatusesPath = strings.ReplaceAll(getStatusesPath, "{domain_id}", domainId)
	getStatusesPath = strings.ReplaceAll(getStatusesPath, "{policy_assignment_id}", d.Get("policy_assignment_id").(string))

	getStatusesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	if resourceKeys, ok := d.GetOk("resource_keys"); ok {
		getStatusesOpt.JSONBody = buildRemediationExecutionStatusesBodyParams(resourceKeys)
		return fetchSpecificStatuses(client, getStatusesPath, getStatusesOpt)
	}

	return fetchAllStatuses(client, getStatusesPath, getStatusesOpt)
}

func fetchSpecificStatuses(client *golangsdk.ServiceClient, path string, opts golangsdk.RequestOpts) (interface{}, error) {
	// API does not support pagination display after specifying query conditions.
	resp, err := client.Request("POST", path, &opts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("value", respBody, make([]interface{}, 0)), nil
}

func fetchAllStatuses(client *golangsdk.ServiceClient, basePath string, opts golangsdk.RequestOpts) (interface{}, error) {
	results := make([]interface{}, 0)
	opts.JSONBody = map[string]interface{}{
		"resource_keys": []interface{}{},
	}
	basePath += "?limit=100"
	path := basePath
	for {
		resp, err := client.Request("POST", path, &opts)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		statusInfo := utils.PathSearch("value[*]", respBody, make([]interface{}, 0))
		results = append(results, statusInfo.([]interface{})...)

		marker := utils.PathSearch("page_info.next_marker", respBody, "")
		if marker == "" {
			break
		}
		path = fmt.Sprintf("%s&marker=%s", basePath, marker)
	}

	return results, nil
}

func buildRemediationExecutionStatusesBodyParams(param interface{}) map[string]interface{} {
	if param == nil {
		return nil
	}

	rawResourceKeys := param.([]interface{})
	resourceKeys := make([]interface{}, 0, len(rawResourceKeys))
	for _, rawResourceKey := range rawResourceKeys {
		rawResourceKeyMap := rawResourceKey.(map[string]interface{})
		resourceKey := map[string]interface{}{
			"resource_type":     rawResourceKeyMap["resource_type"],
			"resource_id":       rawResourceKeyMap["resource_id"],
			"resource_provider": rawResourceKeyMap["resource_provider"],
		}
		resourceKeys = append(resourceKeys, resourceKey)
	}
	return map[string]interface{}{
		"resource_keys": resourceKeys,
	}
}

func flattenValue(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	rst := make([]interface{}, 0)
	curArray := resp.([]interface{})
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"resource_key":    flattenResourceKey(v),
			"invocation_time": utils.PathSearch("invocation_time", v, nil),
			"state":           utils.PathSearch("state", v, nil),
			"message":         utils.PathSearch("message", v, nil),
		})
	}
	return rst
}

func flattenResourceKey(resp interface{}) []interface{} {
	rawResourceKey := utils.PathSearch("resource_key", resp, nil)
	if rawResourceKey == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"resource_type":     utils.PathSearch("resource_type", rawResourceKey, nil),
			"resource_id":       utils.PathSearch("resource_id", rawResourceKey, nil),
			"resource_provider": utils.PathSearch("resource_provider", rawResourceKey, nil),
		},
	}
}
