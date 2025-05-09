package lts

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LTS POST /v3/{project_id}/lts/access-config-list
func DataSourceCceAccesses() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCceAccessesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query CCE access configurations.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the CCE access.`,
			},
			"log_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the log group to which the access configurations and log streams belong.`,
			},
			"log_stream_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the log stream to which the access configurations belong.`,
			},
			"tags": common.TagsSchema(),
			"accesses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the CCE access.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the CCE access.`,
						},
						"log_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the log group.`,
						},
						"log_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the log group.`,
						},
						"log_stream_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the log stream.`,
						},
						"log_stream_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the log stream.`,
						},
						"access_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        dataCceAccessConfigDeatilSchema(),
							Description: `The configuration of the CCE access.`,
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the cluster corresponding to CCE access.`,
						},
						"host_group_ids": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: `The ID list of the log access host groups.`,
						},
						"tags": common.TagsComputedSchema(),
						"binary_collect": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether collect in binary format.`,
						},
						"log_split": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether to split log.`,
						},
						"access_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the log access.`,
						},
						"processor_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the ICAgent structuring parsing.`,
						},
						"processors": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The type of the parser.`,
									},
									"detail": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The configuration of the parser, in JSON format.`,
									},
								},
							},
							Description: `The list of the ICAgent structuring parsing rules.`,
						},
						"demo_log": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The example log of the ICAgent structuring parsing.`,
						},
						"demo_fields": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the parsed field.`,
									},
									"field_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The value of the parsed field.`,
									},
								},
							},
							Description: `The list of the parsed fields of the example log`,
						},
						"encoding_format": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The encoding format log file.`,
						},
						"incremental_collect": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether to collect logs incrementally.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the CCE access, in RFC3339 format.`,
						},
					},
				},
				Description: `The list of CCE access configurations.`,
			},
		},
	}
}

func dataCceAccessConfigDeatilSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"path_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the CCE access.`,
			},
			"paths": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The collection paths.`,
			},
			"black_paths": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The collection path blacklist.`,
			},
			"windows_log_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"categorys": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The collection path blacklist.`,
						},
						"event_level": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of Windows event levels.`,
						},
						"time_offset": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The collection time offset unit.`,
						},
						"time_offset_unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The collection time offset.`,
						},
					},
				},
				Description: `The configuration of Windows event logs.`,
			},
			"single_log_format": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The mode of single-line log format.`,
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The value of single-line log format.`,
						},
					},
				},
				Description: `The configuration single-line logs.`,
			},
			"multi_log_format": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The mode of multi-line log format.`,
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The value of multi-line log format.`,
						},
					},
				},
				Description: `The configuration multi-line logs.`,
			},
			"stdout": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether output is standard.`,
			},
			"stderr": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether error output is standard.`,
			},
			"name_space_regex": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The regular expression matching of kubernetes namespaces.`,
			},
			"pod_name_regex": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The regular expression matching of kubernetes pods.`,
			},
			"container_name_regex": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The regular expression matching of kubernetes container names.`,
			},
			"log_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The container label log tag.`,
			},
			"include_labels_logical": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The logical relationship between multiple container label whitelists.`,
			},
			"include_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The container label whitelist.`,
			},
			"exclude_labels_logical": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The logical relationship between multiple container label blacklists.`,
			},
			"exclude_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The container label blacklist.`,
			},
			"log_envs": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The environment variable tag.`,
			},
			"include_envs_logical": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The logical relationship between multiple environment variable whitelists.`,
			},
			"include_envs": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The environment variable whitelist.`,
			},
			"exclude_envs_logical": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The logical relationship between multiple environment variable blacklists.`,
			},
			"exclude_envs": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The environment variable blacklist.`,
			},
			"log_k8s": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The kubernetes label log tag.`,
			},
			"include_k8s_labels_logical": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The logical relationship between multiple kubernetes label whitelists.`,
			},
			"include_k8s_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The kubernetes label whitelist.`,
			},
			"exclude_k8s_labels_logical": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The logical relationship between multiple kubernetes label blacklists.`,
			},
			"exclude_k8s_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The kubernetes label blacklist.`,
			},
			"repeat_collect": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to allow repeated file collection.`,
			},
			"custom_key_value": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The custom key/value pairs.`,
			},
			"system_fields": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of system built-in fields of the CCE access.`,
			},
		},
	}
	return &sc
}

func buildQueryCceAccessConfigsBodyParams(d *schema.ResourceData) map[string]interface{} {
	queryOpts := make(map[string]interface{})

	if accessConfigName, ok := d.GetOk("name"); ok {
		queryOpts["access_config_name_list"] = []interface{}{accessConfigName}
	}
	if logGroupName, ok := d.GetOk("log_group_name"); ok {
		queryOpts["log_group_name_list"] = []interface{}{logGroupName}
	}
	if logStreamName, ok := d.GetOk("log_stream_name"); ok {
		queryOpts["log_stream_name_list"] = []interface{}{logStreamName}
	}
	if tagsConfig, ok := d.GetOk("tags"); ok {
		queryOpts["access_config_tag_list"] = utils.ExpandResourceTags(tagsConfig.(map[string]interface{}))
	}

	return queryOpts
}

func dataSourceCceAccessesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                        = meta.(*config.Config)
		region                     = cfg.GetRegion(d)
		listCceAccessConfigHttpUrl = "v3/{project_id}/lts/access-config-list"
	)
	ltsClient, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	listCceAccessConfigPath := ltsClient.Endpoint + listCceAccessConfigHttpUrl
	listCceAccessConfigPath = strings.ReplaceAll(listCceAccessConfigPath, "{project_id}", ltsClient.ProjectID)

	listCceAccessConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildQueryCceAccessConfigsBodyParams(d),
	}

	requestResp, err := ltsClient.Request("POST", listCceAccessConfigPath, &listCceAccessConfigOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CCE access configs")
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	flattenedAccesses := flattenAccessConfigDetails(utils.PathSearch("result", respBody, make([]interface{}, 0)).([]interface{}))
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("accesses", flattenedAccesses),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAccessConfigDetails(accesses []interface{}) []interface{} {
	result := make([]interface{}, 0, len(accesses))
	for _, v := range accesses {
		result = append(result, map[string]interface{}{
			"id":                  utils.PathSearch("access_config_id", v, nil),
			"name":                utils.PathSearch("access_config_name", v, nil),
			"log_group_id":        utils.PathSearch("log_info.log_group_id", v, nil),
			"log_group_name":      utils.PathSearch("log_info.log_group_name", v, nil),
			"log_stream_id":       utils.PathSearch("log_info.log_stream_id", v, nil),
			"log_stream_name":     utils.PathSearch("log_info.log_stream_name", v, nil),
			"access_config":       flattenCceAccessConfigDetail(v),
			"cluster_id":          utils.PathSearch("cluster_id", v, nil),
			"host_group_ids":      utils.PathSearch("host_group_info.host_group_id_list", v, nil),
			"tags":                utils.FlattenTagsToMap(utils.PathSearch("access_config_tag", v, nil)),
			"binary_collect":      utils.PathSearch("binary_collect", v, nil),
			"log_split":           utils.PathSearch("log_split", v, nil),
			"access_type":         utils.PathSearch("access_config_type", v, nil),
			"processor_type":      utils.PathSearch("processor_type", v, nil),
			"processors":          flattenCceAccessProcessors(utils.PathSearch("processors", v, make([]interface{}, 0)).([]interface{})),
			"demo_log":            utils.PathSearch("demo_log", v, nil),
			"demo_fields":         flattenCceAccessDemoFields(utils.PathSearch("demo_fields", v, make([]interface{}, 0)).([]interface{})),
			"encoding_format":     utils.PathSearch("encoding_format", v, nil),
			"incremental_collect": utils.PathSearch("incremental_collect", v, nil),
			"created_at":          utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", v, float64(0)).(float64))/1000, false),
		})
	}
	return result
}

func flattenCceAccessProcessors(processors []interface{}) []map[string]interface{} {
	if len(processors) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(processors))
	for i, v := range processors {
		result[i] = map[string]interface{}{
			"type":   utils.PathSearch("type", v, nil),
			"detail": utils.JsonToString(utils.PathSearch("detail", v, nil)),
		}
	}
	return result
}
