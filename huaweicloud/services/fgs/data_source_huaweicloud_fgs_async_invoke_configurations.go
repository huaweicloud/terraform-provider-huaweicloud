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

// @API FunctionGraph GET /v2/{project_id}/fgs/functions/{function_urn}/async-invoke-configs
func DataSourceAsyncInvokeConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAsyncInvokeConfigurationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the async invoke configurations are located.`,
			},

			// Required parameters.
			"function_urn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The function URN to query async invoke configurations.`,
			},

			// Attributes.
			"configurations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of queried async invoke configurations.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"func_urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The function URN.`,
						},
						"max_async_event_age_in_seconds": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum validity period of a message.`,
						},
						"max_async_retry_attempts": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The maximum number of retry attempts to be made if asynchronous invocation fails.`,
						},
						"destination_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The destination configuration for async invoke.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"on_success": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The target to be invoked when a function is successfully executed.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"destination": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The target type.`,
												},
												"param": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The parameters corresponding to the target service, in JSON format.`,
												},
											},
										},
									},
									"on_failure": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The target to be invoked when a function fails to be executed.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"destination": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The target type.`,
												},
												"param": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: `The parameters corresponding to the target service, in JSON format.`,
												},
											},
										},
									},
								},
							},
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the async invoke configuration.`,
						},
						"last_modified": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The last modification time of the async invoke configuration.`,
						},
						"enable_async_status_log": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether async invoke status persistence is enabled.`,
						},
					},
				},
			},
		},
	}
}

func listAsyncInvokeConfigurations(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/fgs/functions/{function_urn}/async-invoke-configs?limit={limit}"
		limit   = 100
		marker  = "0"
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{function_urn}", d.Get("function_urn").(string))
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithMarker := fmt.Sprintf("%s&marker=%s", listPath, marker)
		requestResp, err := client.Request("GET", listPathWithMarker, &listOpt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		configs := utils.PathSearch("async_invoke_configs", respBody, make([]interface{}, 0)).([]interface{})
		if len(configs) < 1 {
			break
		}
		result = append(result, configs...)
		// Check if we've reached the last page
		nextMarker := utils.PathSearch("page_info.next_marker", respBody, float64(0)).(float64)
		if nextMarker == 0 {
			break
		}
		marker = strconv.FormatInt(int64(nextMarker), 10)
	}

	return result, nil
}

func flattenDestinationConfig(destConfig map[string]interface{}) []map[string]interface{} {
	if len(destConfig) < 1 {
		return nil
	}

	return []map[string]interface{}{
		utils.RemoveNil(map[string]interface{}{
			"on_success": flattenDestinationItem(utils.PathSearch("on_success", destConfig,
				make(map[string]interface{})).(map[string]interface{})),
			"on_failure": flattenDestinationItem(utils.PathSearch("on_failure", destConfig,
				make(map[string]interface{})).(map[string]interface{})),
		}),
	}
}

func flattenAsyncInvokeConfigurations(configurations []interface{}) []map[string]interface{} {
	if len(configurations) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(configurations))
	for _, item := range configurations {
		result = append(result, map[string]interface{}{
			"func_urn":                       utils.PathSearch("func_urn", item, nil),
			"max_async_event_age_in_seconds": utils.PathSearch("max_async_event_age_in_seconds", item, nil),
			"max_async_retry_attempts":       utils.PathSearch("max_async_retry_attempts", item, nil),
			"destination_config": flattenDestinationConfig(utils.PathSearch("destination_config", item,
				make(map[string]interface{})).(map[string]interface{})),
			"created_time":            utils.PathSearch("created_time", item, nil),
			"last_modified":           utils.PathSearch("last_modified", item, nil),
			"enable_async_status_log": utils.PathSearch("enable_async_status_log", item, nil),
		})
	}

	return result
}

func dataSourceAsyncInvokeConfigurationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	resp, err := listAsyncInvokeConfigurations(client, d)
	if err != nil {
		return diag.Errorf("error querying FunctionGraph async invoke configurations: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("configurations", flattenAsyncInvokeConfigurations(resp)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
