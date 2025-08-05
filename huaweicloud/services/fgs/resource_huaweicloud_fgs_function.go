package fgs

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API FunctionGraph POST /v2/{project_id}/fgs/functions
// @API FunctionGraph PUT /v2/{project_id}/fgs/functions/{function_urn}/config
// @API FunctionGraph PUT /v2/{project_id}/fgs/functions/{function_urn}/config-max-instance
// @API FunctionGraph GET /v2/{project_id}/fgs/functions/{function_urn}/config
// @API FunctionGraph GET /v2/{project_id}/fgs/functions/{function_urn}/versions
// @API FunctionGraph GET /v2/{project_id}/fgs/functions/{function_urn}/lts-log-detail
// @API FunctionGraph POST /v2/{project_id}/fgs/functions/{function_urn}/tags/create
// @API FunctionGraph DELETE /v2/{project_id}/fgs/functions/{function_urn}/tags/delete
// @API FunctionGraph PUT /v2/{project_id}/fgs/functions/{function_urn}/code
// @API FunctionGraph POST /v2/{project_id}/fgs/functions/{function_urn}/versions
// @API FunctionGraph GET /v2/{project_id}/fgs/functions/{function_urn}/aliases
// @API FunctionGraph POST /v2/{project_id}/fgs/functions/{function_urn}/aliases
// @API FunctionGraph DELETE /v2/{project_id}/fgs/functions/{function_urn}/aliases/{alias_name}
// @API FunctionGraph DELETE /v2/{project_id}/fgs/functions/{function_urn}
// @API FunctionGraph PUT /v2/{project_id}/fgs/functions/{function_urn}/reservedinstances
// @API FunctionGraph GET /v2/{project_id}/fgs/functions/reservedinstanceconfigs
func ResourceFgsFunction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFunctionCreate,
		ReadContext:   resourceFunctionRead,
		UpdateContext: resourceFunctionUpdate,
		DeleteContext: resourceFunctionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the function is located.`,
			},

			// Required parameters.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the function.`,
			},
			"memory_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The memory size allocated to the function, in MByte (MB).`,
			},
			"runtime": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The environment for executing the function.`,
			},
			"timeout": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The timeout interval of the function, in seconds.`,
			},

			// Optional parameters but required in documentation.
			"app": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"package"},
				Description: utils.SchemaDesc(
					`The group to which the function belongs.`,
					utils.SchemaDescInput{
						Required: true,
					},
				),
			},
			"code_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`The code type of the function.`,
					utils.SchemaDescInput{
						Required: true,
					},
				),
			},
			"handler": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`The entry point of the function.`,
					utils.SchemaDescInput{
						Required: true,
					},
				),
			},

			// Optional parameters.
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the function.`,
			},
			"functiongraph_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"v1", "v2",
				}, false), // The current default value is v1, which may be adjusted in the future.
				Description: `The description of the function.`,
			},
			"func_code": {
				Type:        schema.TypeString,
				Optional:    true,
				StateFunc:   utils.DecodeHashAndHexEncode,
				Description: `The function code.`,
			},
			"code_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The URL where the function code is stored in OBS.`,
			},
			"code_filename": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of the function file.`,
			},
			"depend_list": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The ID list of the dependencies.`,
			},
			"user_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The key/value information defined for the function.`,
			},
			"encrypted_user_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `The key/value information defined to be encrypted for the function.`,
			},
			"agency": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"xrole"},
				Description:   `The agency configuration of the function.`,
			},
			"app_agency": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The execution agency enables you to obtain a token or an AK/SK for accessing other cloud services.`,
			},
			"initializer_handler": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The initializer of the function.`,
			},
			"initializer_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The maximum duration the function can be initialized.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the enterprise project to which the function belongs.`,
			},
			"vpc_id": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"network_id"},
				Description:  `The ID of the VPC to which the function belongs.`,
			},
			"network_id": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"vpc_id"},
				Description:  `The network ID of subnet.`,
			},
			"dns_list": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"vpc_id"},
				Description:  `The private DNS configuration of the function network.`,
			},
			"mount_user_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The mount user ID.`,
			},
			"mount_user_group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The mount user group ID.`,
			},
			"func_mounts": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The mount type.`,
						},
						"mount_resource": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The ID of the mounted resource (corresponding cloud service).`,
						},
						"mount_share_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The remote mount path.`,
						},
						"local_mount_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The function access path.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The mount status.`,
						},
					},
				},
				Description: `The list of function mount configuration.`,
			},
			"custom_image": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The URL of SWR image.`,
						},
						"command": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The startup commands of the SWR image.`,
						},
						"args": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The command line arguments used to start the SWR image.`,
						},
						"working_dir": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The working directory of the SWR image.`,
						},
						"user_id": {
							Type:     schema.TypeString,
							Optional: true,
							Description: utils.SchemaDesc(
								`The user ID that used to run SWR image.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"user_group_id": {
							Type:     schema.TypeString,
							Optional: true,
							Description: utils.SchemaDesc(
								`The user group ID that used to run SWR image.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
					},
				},
				ConflictsWith: []string{
					"code_type",
				},
				Description: `The custom image configuration of the function.`,
			},
			"max_instance_num": {
				// The original type of this parameter is int, but its zero value is meaningful.
				// So, the following types of parameter passing are realized through the logic of terraform's implicit
				// conversion of int:
				//   + -1: the number of instances is unlimited.
				//   + 0: this function is disabled.
				//   + (0, +1000]: Specific value (2023.06.26).
				//   + empty: keep the default (latest updated) value.
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The maximum number of instances of the function.`,
			},
			"versions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The version name.`,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The description of the version.`,
						},
						"aliases": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The name of the version alias.`,
									},
									"description": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The description of the version alias.`,
									},
									"additional_version_weights": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsJSON,
										Description:  `The percentage grayscale configuration of the version alias.`,
									},
									"additional_version_strategy": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsJSON,
										Description:  `The description of the version alias.`,
									},
								},
							},
							Description: `The aliases management for specified version.`,
						},
					},
				},
				Description: `The versions management of the function.`,
			},
			"tags": common.TagsSchema(
				`The key/value pairs to associate with the function.`,
			),
			"log_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The LTS group ID for collecting logs.`,
			},
			"log_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The LTS group name for collecting logs.`,
			},
			"log_stream_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The LTS stream ID for collecting logs.`,
			},
			"log_stream_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"log_group_id"},
				Description:  `The LTS stream name for collecting logs.`,
			},
			"reserved_instances": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"qualifier_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The qualifier type of reserved instance.`,
						},
						"qualifier_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The version name or alias name.`,
						},
						"count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The number of reserved instance.`,
						},
						"idle_mode": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether to enable the idle mode.`,
						},
						"tactics_config": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Elem:        tracticsConfigsSchema(),
							Description: `The auto scaling policies for reserved instance.`,
						},
					},
				},
				Description: `The reserved instance policies of the function.`,
			},
			// The value in the api document is -1 to 1000, After confirmation, when the parameter set to -1 or 0,
			// the actual number of concurrent requests is 1, so the value range is set to 1 to 1000, and the document
			// will be modified later (2024.02.29).
			"concurrency_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The number of concurrent requests of the function.`,
			},
			"gpu_type": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"gpu_memory"},
				Description:  `The GPU type of the function.`,
			},
			"gpu_memory": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"gpu_type"},
				Description:  `The GPU memory size allocated to the function, in MByte (MB).`,
			},
			// Currently, the "pre_stop_timeout" and "pre_stop_timeout" are not visible on the page,
			// so they are temporarily used as internal parameters.
			"pre_stop_handler": {
				Type:     schema.TypeString,
				Optional: true,
				Description: utils.SchemaDesc(
					`The pre-stop handler of a function.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"pre_stop_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: utils.SchemaDesc(
					`The maximum duration that the function can be initialized.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"enable_dynamic_memory": {
				Type:     schema.TypeBool,
				Optional: true,
				// The dynamic memory function can be closed, so computed behavior cannot be supported.
				Description: `Whether the dynamic memory configuration is enabled.`,
			},
			"is_stateful_function": {
				Type:     schema.TypeBool,
				Optional: true,
				// The stateful function can be closed, so computed behavior cannot be supported.
				Description: `Whether the function is a stateful function.`,
			},
			"network_controller": {
				Type:     schema.TypeList,
				Optional: true,
				// The network controller can be closed, so computed behavior cannot be supported.
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trigger_access_vpcs": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `The ID of the VPC that can trigger the function.`,
									},
									"vpc_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `The ID of the VPC that can trigger the function.`,
									},
								},
							},
							Description: `The configuration of the VPCs that can trigger the function.`,
						},
						"disable_public_network": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether to disable the public network access.`,
						},
					},
				},
				Description: `The network configuration of the function.`,
			},
			"peering_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				// The perring CIDR can be canceled, so computed behavior cannot be supported.
				Description: `The VPC CIDR blocks used in the function code to detect whether it conflicts with the VPC
CIDR blocks used by the service.`,
			},
			"enable_auth_in_header": {
				Type:     schema.TypeBool,
				Optional: true,
				// The auth function can be closed, so computed behavior cannot be supported.
				// And the default value (in the service API) is false.
				Description: `Whether the authentication in the request header is enabled.`,
			},
			"enable_class_isolation": {
				Type:     schema.TypeBool,
				Optional: true,
				// The isolation function can be closed, so computed behavior cannot be supported.
				// And the default value (in the service API) is false.
				Description: `Whether the class isolation is enabled for the JAVA runtime functions.`,
			},
			"ephemeral_storage": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The size of the function ephemeral storage.`,
			},
			"heartbeat_handler": {
				Type:     schema.TypeString,
				Optional: true,
				// The handler of heartbeat can be omitted, and if this parameter is not configured, the default value
				// will not be returned. So, the computed behavior cannot be supported.
				Description: `The heartbeat handler of the function.`,
			},
			"restore_hook_handler": {
				Type:     schema.TypeString,
				Optional: true,
				// The handler of restore hook can be omitted, and if this parameter is not configured, the default
				// value will not be returned. So, the computed behavior cannot be supported.
				Description: `The restore hook handler of the function.`,
			},
			"restore_hook_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				// The timeout of restore hook can be omitted, and if this parameter is not configured, the default
				// value will not be returned. So, the computed behavior cannot be supported.
				Description: `The timeout of the function restore hook.`,
			},
			"lts_custom_tag": {
				Type:             schema.TypeMap,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: utils.SuppressMapDiffs(),
				// The custom tags can be set to empty, so computed behavior cannot be supported.
				Description: `The custom tags configuration that used to filter the LTS logs.`,
			},
			"enable_lts_log": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to enable the LTS log.`,
			},
			"user_data_encrypt_kms_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The KMS key ID for encrypting the user data.`,
			},
			"code_encrypt_kms_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The KMS key ID for encrypting the function code.`,
			},

			// Deprecated parameters.
			"package": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"app"},
				Deprecated:    `use app instead`,
			},
			"xrole": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"agency"},
				Deprecated:    `use agency instead`,
			},

			// Attributes.
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URN (Uniform Resource Name) of the function.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version of the function.`,
			},
			"lts_custom_tag_origin": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
 the new value next time the change is made. The corresponding parameter name is 'lts_custom_tag'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func tracticsConfigsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cron_configs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of scheduled policy configuration.`,
						},
						"cron": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The cron expression.`,
						},
						"count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The number of reserved instance to which the policy belongs.`,
						},
						"start_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The effective timestamp of policy.`,
						},
						"expired_time": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The expiration timestamp of the policy.`,
						},
					},
				},
				Description: `The list of scheduled policy configurations.`,
			},
			"metric_configs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of metric policy.`,
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The type of metric policy.`,
						},
						"threshold": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The metric policy threshold.`,
						},
						"min": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The minimun of traffic.`,
						},
					},
				},
				Description: `The list of metric policy configurations.`,
			},
		},
	}
}

func buildFunctionCustomImage(imageConfigs []interface{}) map[string]interface{} {
	if len(imageConfigs) < 1 {
		return nil
	}

	imageConfig := imageConfigs[0]
	return map[string]interface{}{
		"enabled":     true,
		"image":       utils.ValueIgnoreEmpty(utils.PathSearch("url", imageConfig, nil)),
		"command":     utils.ValueIgnoreEmpty(utils.PathSearch("command", imageConfig, nil)),
		"args":        utils.ValueIgnoreEmpty(utils.PathSearch("args", imageConfig, nil)),
		"working_dir": utils.ValueIgnoreEmpty(utils.PathSearch("working_dir", imageConfig, nil)),
		"uid":         utils.ValueIgnoreEmpty(utils.PathSearch("user_id", imageConfig, nil)),
		"gid":         utils.ValueIgnoreEmpty(utils.PathSearch("user_group_id", imageConfig, nil)),
	}
}

func buildFunctionCodeConfig(funcCode string) map[string]interface{} {
	if funcCode == "" {
		return nil
	}

	return map[string]interface{}{
		"file": utils.TryBase64EncodeString(funcCode),
	}
}

func buildFunctionLogConfig(rawConfig cty.Value) interface{} {
	params := utils.RemoveNil(map[string]interface{}{
		"group_id":    utils.ValueIgnoreEmpty(utils.GetNestedObjectFromRawConfig(rawConfig, "log_group_id")),
		"group_name":  utils.ValueIgnoreEmpty(utils.GetNestedObjectFromRawConfig(rawConfig, "log_group_name")),
		"stream_id":   utils.ValueIgnoreEmpty(utils.GetNestedObjectFromRawConfig(rawConfig, "log_stream_id")),
		"stream_name": utils.ValueIgnoreEmpty(utils.GetNestedObjectFromRawConfig(rawConfig, "log_stream_name")),
	})

	// If the value of `enable_lts_log` parameter is `true`, the corresponding LTS log parameters be configured.
	if len(params) == 0 {
		return nil
	}

	return params
}

func buildNetworkControllerTriggerAccessVpcs(triggerAccessVpcs []interface{}) []map[string]interface{} {
	if len(triggerAccessVpcs) < 1 || triggerAccessVpcs[0] == nil {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(triggerAccessVpcs))
	for _, triggerAccessVpc := range triggerAccessVpcs {
		result = append(result, map[string]interface{}{
			"vpc_id":   utils.PathSearch("vpc_id", triggerAccessVpc, nil),
			"vpc_name": utils.PathSearch("vpc_name", triggerAccessVpc, nil),
		})
	}
	return result
}

func buildFunctionNetworkController(networkControlers []interface{}) map[string]interface{} {
	if len(networkControlers) < 1 || networkControlers[0] == nil {
		return nil
	}

	networkControler := networkControlers[0]
	return map[string]interface{}{
		"trigger_access_vpcs": buildNetworkControllerTriggerAccessVpcs(utils.PathSearch("trigger_access_vpcs",
			networkControler, schema.NewSet(schema.HashString, nil)).(*schema.Set).List()),
		"disable_public_network": utils.PathSearch("disable_public_network", networkControler, nil),
	}
}

func buildCreateFunctionBodyParams(cfg *config.Config, d *schema.ResourceData) map[string]interface{} {
	// Parameter app is recommended to replace parameter package.
	pkg, ok := d.GetOk("app")
	if !ok {
		pkg = d.Get("package")
	}

	// Parameter agency is recommended to replace parameter xrole.
	agency, ok := d.GetOk("agency")
	if !ok {
		agency = d.Get("xrole")
	}

	rawConfig := d.GetRawConfig()
	return map[string]interface{}{
		// Required parameters.
		"func_name":   d.Get("name"),
		"runtime":     d.Get("runtime"),
		"timeout":     d.Get("timeout"),
		"memory_size": d.Get("memory_size"),
		// Optional parameters but required in documentation.
		"package":   utils.ValueIgnoreEmpty(pkg),
		"handler":   utils.ValueIgnoreEmpty(d.Get("handler")),
		"code_type": utils.ValueIgnoreEmpty(d.Get("code_type")),
		// Optional parameters.
		"description":                  utils.ValueIgnoreEmpty(d.Get("description")),
		"type":                         utils.ValueIgnoreEmpty(d.Get("functiongraph_version")),
		"code_url":                     utils.ValueIgnoreEmpty(d.Get("code_url")),
		"code_filename":                utils.ValueIgnoreEmpty(d.Get("code_filename")),
		"user_data":                    utils.ValueIgnoreEmpty(d.Get("user_data")),
		"encrypted_user_data":          utils.ValueIgnoreEmpty(d.Get("encrypted_user_data")),
		"xrole":                        utils.ValueIgnoreEmpty(agency),
		"enterprise_project_id":        cfg.GetEnterpriseProjectID(d),
		"custom_image":                 buildFunctionCustomImage(d.Get("custom_image").([]interface{})),
		"gpu_memory":                   utils.ValueIgnoreEmpty(d.Get("gpu_memory")),
		"gpu_type":                     utils.ValueIgnoreEmpty(d.Get("gpu_type")),
		"pre_stop_handler":             utils.ValueIgnoreEmpty(d.Get("pre_stop_handler")),
		"pre_stop_timeout":             utils.ValueIgnoreEmpty(d.Get("pre_stop_timeout")),
		"func_code":                    buildFunctionCodeConfig(d.Get("func_code").(string)),
		"log_config":                   buildFunctionLogConfig(rawConfig),
		"enable_dynamic_memory":        d.Get("enable_dynamic_memory"),
		"is_stateful_function":         d.Get("is_stateful_function"),
		"network_controller":           buildFunctionNetworkController(d.Get("network_controller").([]interface{})),
		"lts_custom_tag":               utils.ValueIgnoreEmpty(d.Get("lts_custom_tag")),
		"enable_lts_log":               utils.GetNestedObjectFromRawConfig(rawConfig, "enable_lts_log"),
		"user_data_encrypt_kms_key_id": utils.ValueIgnoreEmpty(d.Get("user_data_encrypt_kms_key_id")),
		"code_encrypt_kms_key_id":      utils.ValueIgnoreEmpty(d.Get("code_encrypt_kms_key_id")),
	}
}

func createFunction(cfg *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData) (string, error) {
	httpUrl := "v2/{project_id}/fgs/functions"

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildCreateFunctionBodyParams(cfg, d)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return "", err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return "", err
	}

	return utils.PathSearch("func_urn", respBody, "").(string), nil
}

func buildFunctionVpcConfig(d *schema.ResourceData) map[string]interface{} {
	vpcId, ok := d.GetOk("vpc_id")
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"vpc_id":    vpcId,
		"subnet_id": d.Get("network_id"),
	}
}

func parseFunctionMountId(mountId int) int {
	if mountId < 1 {
		return -1
	}
	return mountId
}

func buildFunctionMountConfig(mounts []interface{}, mountUserId, mountGroupId int) map[string]interface{} {
	if len(mounts) < 1 {
		return nil
	}

	parsedMounts := make([]interface{}, 0, len(mounts))
	for _, mount := range mounts {
		parsedMounts = append(parsedMounts, map[string]interface{}{
			"mount_type":       utils.PathSearch("mount_type", mount, nil),
			"mount_resource":   utils.PathSearch("mount_resource", mount, nil),
			"mount_share_path": utils.PathSearch("mount_share_path", mount, nil),
			"local_mount_path": utils.PathSearch("local_mount_path", mount, nil),
		})
	}

	return map[string]interface{}{
		"mount_config": parsedMounts,
		"mount_user": map[string]interface{}{
			"mount_user_id":       parseFunctionMountId(mountUserId),
			"mount_user_group_id": parseFunctionMountId(mountGroupId),
		},
	}
}

func buildFunctionStrategyConfig(concurrencyNum int) map[string]interface{} {
	if concurrencyNum < 1 {
		return nil
	}

	return map[string]interface{}{
		"concurrent_num": concurrencyNum,
	}
}

func buildUpdateFunctionMetadataBodyParams(cfg *config.Config, d *schema.ResourceData) map[string]interface{} {
	// Parameter app is recommended to replace parameter package.
	pkg, ok := d.GetOk("app")
	if !ok {
		pkg = d.Get("package")
	}

	// Parameter agency is recommended to replace parameter xrole.
	agency, ok := d.GetOk("agency")
	if !ok {
		pkg = d.Get("xrole")
	}

	return map[string]interface{}{
		// Required parameters.
		"runtime":     d.Get("runtime"),
		"timeout":     d.Get("timeout"),
		"memory_size": d.Get("memory_size"),
		// Optional parameters but required in documentation.
		"package": utils.ValueIgnoreEmpty(pkg),
		"handler": utils.ValueIgnoreEmpty(d.Get("handler")),
		// Optional parameters.
		"description":         d.Get("description"),
		"user_data":           utils.ValueIgnoreEmpty(d.Get("user_data")),
		"encrypted_user_data": utils.ValueIgnoreEmpty(d.Get("encrypted_user_data")),
		"xrole":               utils.ValueIgnoreEmpty(agency),
		"app_xrole":           utils.ValueIgnoreEmpty(d.Get("app_agency")),
		"custom_image":        buildFunctionCustomImage(d.Get("custom_image").([]interface{})),
		"gpu_memory":          utils.ValueIgnoreEmpty(d.Get("gpu_memory")),
		"gpu_type":            utils.ValueIgnoreEmpty(d.Get("gpu_type")),
		"initializer_handler": utils.ValueIgnoreEmpty(d.Get("initializer_handler")),
		"initializer_timeout": utils.ValueIgnoreEmpty(d.Get("initializer_timeout")),
		"pre_stop_handler":    utils.ValueIgnoreEmpty(d.Get("pre_stop_handler")),
		"pre_stop_timeout":    utils.ValueIgnoreEmpty(d.Get("pre_stop_timeout")),
		"domain_names":        utils.ValueIgnoreEmpty(d.Get("dns_list")),
		"func_vpc":            buildFunctionVpcConfig(d),
		"func_mounts": buildFunctionMountConfig(d.Get("func_mounts").([]interface{}),
			d.Get("mount_user_id").(int), d.Get("mount_user_group_id").(int)),
		"strategy_config":              buildFunctionStrategyConfig(d.Get("concurrency_num").(int)),
		"enable_dynamic_memory":        d.Get("enable_dynamic_memory"),
		"is_stateful_function":         d.Get("is_stateful_function"),
		"network_controller":           buildFunctionNetworkController(d.Get("network_controller").([]interface{})),
		"enterprise_project_id":        cfg.GetEnterpriseProjectID(d),
		"peering_cidr":                 d.Get("peering_cidr"),
		"enable_auth_in_header":        d.Get("enable_auth_in_header"),
		"enable_class_isolation":       d.Get("enable_class_isolation"),
		"ephemeral_storage":            utils.ValueIgnoreEmpty(d.Get("ephemeral_storage")),
		"heartbeat_handler":            d.Get("heartbeat_handler"),
		"restore_hook_handler":         d.Get("restore_hook_handler"),
		"restore_hook_timeout":         d.Get("restore_hook_timeout"),
		"lts_custom_tag":               utils.ValueIgnoreEmpty(d.Get("lts_custom_tag")),
		"user_data_encrypt_kms_key_id": utils.ValueIgnoreEmpty(d.Get("user_data_encrypt_kms_key_id")),
	}
}

func updateFunctionMetadata(client *golangsdk.ServiceClient, functionUrn string, params map[string]interface{}) error {
	httpUrl := "v2/{project_id}/fgs/functions/{function_urn}/config"

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{function_urn}", functionUrn)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(params),
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("failed to update the function metadata: %s", err)
	}
	return nil
}

func buildUpdateFunctionCodeBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"code_type":           d.Get("code_type"),
		"code_url":            d.Get("code_url"),
		"code_filename":       d.Get("code_filename"),
		"depend_version_list": d.Get("depend_list").(*schema.Set).List(),
		"func_code":           buildFunctionCodeConfig(d.Get("func_code").(string)),
	}
}

func updateFunctionCode(client *golangsdk.ServiceClient, d *schema.ResourceData, functionUrn string) error {
	httpUrl := "v2/{project_id}/fgs/functions/{function_urn}/code"

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{function_urn}", functionUrn)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildUpdateFunctionCodeBodyParams(d)),
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("failed to update the function code: %s", err)
	}
	return nil
}

func updateFunctionMaxInstanceNum(client *golangsdk.ServiceClient, functionUrn string, maxInstanceNum int) error {
	httpUrl := "v2/{project_id}/fgs/functions/{function_urn}/config-max-instance"

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{function_urn}", functionUrn)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"max_instance_num": maxInstanceNum,
		},
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("failed to update the function max instance number: %s", err)
	}
	return nil
}

func buildFunctionTagsBodyParams(tags map[string]interface{}) map[string]interface{} {
	tagsList := make([]interface{}, 0, len(tags))

	for k, v := range tags {
		tagsList = append(tagsList, map[string]interface{}{
			"key":   k,
			"value": v,
		})
	}

	return map[string]interface{}{
		"tags": tagsList,
	}
}

func createFunctionTags(client *golangsdk.ServiceClient, functionUrn string, tags map[string]interface{}) error {
	if len(tags) < 1 {
		return nil
	}

	httpUrl := "v2/{project_id}/functions/{function_urn}/tags/create"

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{function_urn}", functionUrn)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		OkCodes:  []int{204},
		JSONBody: buildFunctionTagsBodyParams(tags),
	}
	_, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return fmt.Errorf("failed to create the function tags: %s", err)
	}
	return nil
}

func getFunctionVersionUrn(client *golangsdk.ServiceClient, functionUrn string, qualifierName string) (string, error) {
	versions, err := getFunctionVersions(client, functionUrn)
	if err != nil {
		return "", err
	}

	return utils.PathSearch(fmt.Sprintf("[?version=='%s']|[0].func_urn", qualifierName), versions, "").(string), nil
}

func getFunctionAliasUrn(client *golangsdk.ServiceClient, functionUrn string, qualifierName string) (string, error) {
	aliases, err := getFunctionAliases(client, functionUrn)
	if err != nil {
		return "", err
	}

	// There will be no aliases with the same name even between different versions.
	return utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0].alias_urn", qualifierName), aliases, "").(string), nil
}

func getReservedInstanceUrn(client *golangsdk.ServiceClient, functionUrn, qualifierType, qualifierName string) (string, error) {
	if qualifierType == "version" {
		return getFunctionVersionUrn(client, functionUrn, qualifierName)
	}

	return getFunctionAliasUrn(client, functionUrn, qualifierName)
}

func deleteFunctionReservedInstances(client *golangsdk.ServiceClient, functionUrn string, policies []interface{}) error {
	httpUrl := "v2/{project_id}/fgs/functions/{function_urn}/reservedinstances"

	basePath := client.Endpoint + httpUrl
	basePath = strings.ReplaceAll(basePath, "{project_id}", client.ProjectID)

	for _, policy := range policies {
		qualifierType := utils.PathSearch("qualifier_type", policy, "").(string)
		qualifierName := utils.PathSearch("qualifier_name", policy, "").(string)

		urn, err := getReservedInstanceUrn(client, functionUrn, qualifierType, qualifierName)
		if err != nil {
			return err
		}
		// Deleting the alias will also delete the corresponding reserved instance.
		if urn == "" {
			return nil
		}

		deletePath := basePath
		deletePath = strings.ReplaceAll(deletePath, "{function_urn}", urn)
		deleteOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			JSONBody: map[string]interface{}{
				"count":     0,
				"idle_mode": false,
			},
		}

		_, err = client.Request("PUT", deletePath, &deleteOpt)
		if err != nil {
			return fmt.Errorf("failed to remove the function reversed instance: %s", err)
		}
	}

	return nil
}

func buildTracticsConfigCronConfigs(cronConfigs []interface{}) []map[string]interface{} {
	if len(cronConfigs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(cronConfigs))
	for _, cronConfig := range cronConfigs {
		result = append(result, map[string]interface{}{
			"name":         utils.ValueIgnoreEmpty(utils.PathSearch("name", cronConfig, nil)),
			"cron":         utils.ValueIgnoreEmpty(utils.PathSearch("cron", cronConfig, nil)),
			"count":        utils.ValueIgnoreEmpty(utils.PathSearch("count", cronConfig, nil)),
			"start_time":   utils.ValueIgnoreEmpty(utils.PathSearch("start_time", cronConfig, nil)),
			"expired_time": utils.ValueIgnoreEmpty(utils.PathSearch("expired_time", cronConfig, nil)),
		})
	}

	return result
}

func buildTracticsConfigMetricConfigs(metricConfigs []interface{}) []map[string]interface{} {
	if len(metricConfigs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(metricConfigs))
	for _, metricConfig := range metricConfigs {
		result = append(result, map[string]interface{}{
			"name":      utils.ValueIgnoreEmpty(utils.PathSearch("name", metricConfig, nil)),
			"type":      utils.ValueIgnoreEmpty(utils.PathSearch("type", metricConfig, nil)),
			"threshold": utils.ValueIgnoreEmpty(utils.PathSearch("threshold", metricConfig, nil)),
			"min":       utils.ValueIgnoreEmpty(utils.PathSearch("min", metricConfig, nil)),
		})
	}

	return result
}

func buildReservedInstanceTracticsConfigs(tacticsConfigs []interface{}) map[string]interface{} {
	if len(tacticsConfigs) < 1 {
		return nil
	}

	return map[string]interface{}{
		"cron_configs": buildTracticsConfigCronConfigs(utils.PathSearch("cron_configs",
			tacticsConfigs[0], make([]interface{}, 0)).([]interface{})),
		"metric_configs": buildTracticsConfigMetricConfigs(utils.PathSearch("metric_configs",
			tacticsConfigs[0], make([]interface{}, 0)).([]interface{})),
	}
}

func createFunctionReservedInstances(client *golangsdk.ServiceClient, functionUrn string, policies []interface{}) error {
	httpUrl := "v2/{project_id}/fgs/functions/{function_urn}/reservedinstances"

	basePath := client.Endpoint + httpUrl
	basePath = strings.ReplaceAll(basePath, "{project_id}", client.ProjectID)

	for _, policy := range policies {
		qualifierType := utils.PathSearch("qualifier_type", policy, "").(string)
		qualifierName := utils.PathSearch("qualifier_name", policy, "").(string)

		urn, err := getReservedInstanceUrn(client, functionUrn, qualifierType, qualifierName)
		if err != nil {
			return err
		}

		createPath := basePath
		createPath = strings.ReplaceAll(createPath, "{function_urn}", urn)
		createOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			JSONBody: utils.RemoveNil(map[string]interface{}{
				"count":     utils.PathSearch("count", policy, nil),
				"idle_mode": utils.PathSearch("idle_mode", policy, nil),
				"tactics_config": buildReservedInstanceTracticsConfigs(utils.PathSearch("tactics_config",
					policy, make([]interface{}, 0)).([]interface{})),
			}),
		}

		_, err = client.Request("PUT", createPath, &createOpt)
		if err != nil {
			return fmt.Errorf("failed to configure the function reversed instance: %s", err)
		}
	}

	return nil
}

func updateFunctionReservedInstances(client *golangsdk.ServiceClient, d *schema.ResourceData, functionUrn string) error {
	oldVal, newVal := d.GetChange("reserved_instances")

	oldRaws := oldVal.(*schema.Set).Difference(newVal.(*schema.Set))
	newRaws := newVal.(*schema.Set).Difference(oldVal.(*schema.Set))

	if oldRaws.Len() > 0 {
		if err := deleteFunctionReservedInstances(client, functionUrn, oldRaws.List()); err != nil {
			return err
		}
	}

	if newRaws.Len() > 0 {
		if err := createFunctionReservedInstances(client, functionUrn, newRaws.List()); err != nil {
			return err
		}
	}

	return nil
}

func resourceFunctionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                             = meta.(*config.Config)
		region                          = cfg.GetRegion(d)
		functionMetadataObjectParamKeys = []string{
			"lts_custom_tag",
		}
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	funcUrn, err := createFunction(cfg, client, d)
	if err != nil {
		return diag.Errorf("error creating function: %s", err)
	}
	if funcUrn == "" {
		return diag.Errorf("unable to find the function URN from the API response")
	}
	d.SetId(funcUrn)
	funcUrnWithoutVersion := parseFunctionUrnWithoutVersion(funcUrn)

	// lintignore:R019
	if d.HasChanges("vpc_id", "network_id", "func_mounts", "app_agency", "initializer_handler", "initializer_timeout",
		"concurrency_num", "peering_cidr", "enable_auth_in_header", "enable_class_isolation", "ephemeral_storage",
		"heartbeat_handler", "restore_hook_handler", "restore_hook_timeout", "lts_custom_tag") {
		params := buildUpdateFunctionMetadataBodyParams(cfg, d)
		err = updateFunctionMetadata(client, funcUrnWithoutVersion, params)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// If the request is successful, obtain the values of all JSON|object parameters first and save them to the
	// corresponding '_origin' attributes for subsequent determination and construction of the request body during
	// next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	err = utils.RefreshObjectParamOriginValues(d, functionMetadataObjectParamKeys)
	if err != nil {
		return diag.Errorf("unable to refresh the origin values: %s", err)
	}

	if d.HasChange("depend_list") {
		err := updateFunctionCode(client, d, funcUrnWithoutVersion)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if strNum, ok := d.GetOk("max_instance_num"); ok {
		// If the maximum number of instances is omitted (after type conversion, the value is zero), means this feature
		// is disabled.
		maxInstanceNum, _ := strconv.Atoi(strNum.(string))
		err = updateFunctionMaxInstanceNum(client, funcUrnWithoutVersion, maxInstanceNum)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if tags, ok := d.GetOk("tags"); ok {
		if err := createFunctionTags(client, funcUrn, tags.(map[string]interface{})); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("versions") {
		if err = updateFunctionVersions(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("reserved_instances") {
		if err = updateFunctionReservedInstances(client, d, funcUrnWithoutVersion); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceFunctionRead(ctx, d, meta)
}

func GetFunctionMetadata(client *golangsdk.ServiceClient, functionUrn string) (interface{}, error) {
	httpUrl := "v2/{project_id}/fgs/functions/{function_urn}/config"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{function_urn}", functionUrn)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func flattenFgsCustomImage(imageConfig map[string]interface{}) []map[string]interface{} {
	if len(imageConfig) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"url":           utils.PathSearch("image", imageConfig, nil),
			"command":       utils.PathSearch("command", imageConfig, nil),
			"args":          utils.PathSearch("args", imageConfig, nil),
			"working_dir":   utils.PathSearch("working_dir", imageConfig, nil),
			"user_id":       utils.PathSearch("uid", imageConfig, nil),
			"user_group_id": utils.PathSearch("gid", imageConfig, nil),
		},
	}
}

func flattenFuncionMounts(mounts []interface{}) []map[string]interface{} {
	if len(mounts) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(mounts))
	for _, mount := range mounts {
		result = append(result, map[string]interface{}{
			"mount_type":       utils.PathSearch("mount_type", mount, nil),
			"mount_resource":   utils.PathSearch("mount_resource", mount, nil),
			"mount_share_path": utils.PathSearch("mount_share_path", mount, nil),
			"local_mount_path": utils.PathSearch("local_mount_path", mount, nil),
			"status":           utils.PathSearch("status", mount, nil),
		})
	}

	return result
}

func getFunctionVersions(client *golangsdk.ServiceClient, functionUrn string) ([]interface{}, error) {
	httpUrl := "v2/{project_id}/fgs/functions/{function_urn}/versions"

	// The query parameter 'marker' and 'maxitems' are not available.
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{function_urn}", functionUrn)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("versions", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func getFunctionAliases(client *golangsdk.ServiceClient, functionUrn string) ([]interface{}, error) {
	httpUrl := "v2/{project_id}/fgs/functions/{function_urn}/aliases"

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{function_urn}", functionUrn)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("[]", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenFunctionVersionAliases(aliases []interface{}) []map[string]interface{} {
	if len(aliases) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(aliases))
	for _, alias := range aliases {
		result = append(result, map[string]interface{}{
			"name":                        utils.PathSearch("name", alias, nil),
			"description":                 utils.PathSearch("description", alias, nil),
			"additional_version_weights":  utils.JsonToString(utils.PathSearch("additional_version_weights", alias, nil)),
			"additional_version_strategy": utils.JsonToString(utils.PathSearch("additional_version_strategy", alias, nil)),
		})
	}

	return result
}

func flattenNetworkControllerTriggerAccessVpcs(triggerAccessVpcs []interface{}) []map[string]interface{} {
	if len(triggerAccessVpcs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(triggerAccessVpcs))
	for _, triggerAccessVpc := range triggerAccessVpcs {
		result = append(result, map[string]interface{}{
			"vpc_id":   utils.PathSearch("vpc_id", triggerAccessVpc, nil),
			"vpc_name": utils.PathSearch("vpc_name", triggerAccessVpc, nil),
		})
	}

	return result
}

func flattenFunctionNetworkController(networkController interface{}) []map[string]interface{} {
	if networkController == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"trigger_access_vpcs": flattenNetworkControllerTriggerAccessVpcs(utils.PathSearch("trigger_access_vpcs",
				networkController, make([]interface{}, 0)).([]interface{})),
			"disable_public_network": utils.PathSearch("disable_public_network", networkController, nil),
		},
	}
}

func flattenFunctionVersions(client *golangsdk.ServiceClient, functionUrn string) ([]map[string]interface{}, error) {
	versionList, err := getFunctionVersions(client, functionUrn)
	if err != nil {
		return nil, err
	}
	aliasesConfig, err := getFunctionAliases(client, functionUrn)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(versionList))
	for _, version := range versionList {
		versionName := utils.PathSearch("version", version, "").(string)
		if versionName == "" {
			log.Printf("[DEBUG] The version name is not found from the API response: %v", version)
			continue
		}
		versionDesc := utils.PathSearch("version_description", version, "").(string)
		aliases := utils.PathSearch(fmt.Sprintf("[?version=='%s']", versionName),
			aliasesConfig, make([]interface{}, 0)).([]interface{})
		if versionName == "latest" {
			if len(aliases) < 1 {
				continue
			}
			// The description of the latest version is configured through the 'description' parameter of the function resource,
			// not the 'versions.description' parameter.
			versionDesc = ""
		}
		result = append(result, utils.RemoveNil(map[string]interface{}{
			"name":        versionName,
			"aliases":     flattenFunctionVersionAliases(aliases),
			"description": utils.ValueIgnoreEmpty(versionDesc),
		}))
	}

	return result, nil
}

func getFunctionReservedInstances(client *golangsdk.ServiceClient, functionUrn string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/fgs/functions/reservedinstanceconfigs?function_urn={function_urn}&limit=100"
		marker  = 0
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{function_urn}", functionUrn)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	result := make([]interface{}, 0)
	for {
		listPathWithMarker := fmt.Sprintf("%s&marker=%d", listPath, marker)
		requestResp, err := client.Request("GET", listPathWithMarker, &listOpt)
		if err != nil {
			return nil, fmt.Errorf("error querying function reserved instances: %s", err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		reservedInstances := utils.PathSearch("reserved_instances", respBody, make([]interface{}, 0)).([]interface{})
		if len(reservedInstances) < 1 {
			break
		}
		result = append(result, reservedInstances...)
		marker += len(reservedInstances)
	}

	return result, nil
}

func flattenReservedInstanceCronConfig(cronConfigs []interface{}) []map[string]interface{} {
	if len(cronConfigs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(cronConfigs))
	for _, cronConfig := range cronConfigs {
		result = append(result, map[string]interface{}{
			"name":         utils.PathSearch("name", cronConfig, nil),
			"cron":         utils.PathSearch("cron", cronConfig, nil),
			"count":        utils.PathSearch("count", cronConfig, nil),
			"start_time":   utils.PathSearch("start_time", cronConfig, nil),
			"expired_time": utils.PathSearch("expired_time", cronConfig, nil),
		})
	}

	return result
}

func flattenReservedInstanceMetricConfig(metricConfigs []interface{}) []map[string]interface{} {
	if len(metricConfigs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(metricConfigs))
	for _, metricConfig := range metricConfigs {
		result = append(result, map[string]interface{}{
			"name":      utils.PathSearch("name", metricConfig, nil),
			"type":      utils.PathSearch("type", metricConfig, nil),
			"threshold": utils.PathSearch("threshold", metricConfig, nil),
			"min":       utils.PathSearch("min", metricConfig, nil),
		})
	}

	return result
}

func flattenReservedInstanceTracticsConfig(tracticsConfig map[string]interface{}) []map[string]interface{} {
	if len(tracticsConfig) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"cron_configs": flattenReservedInstanceCronConfig(utils.PathSearch("cron_configs",
				tracticsConfig, make([]interface{}, 0)).([]interface{})),
			"metric_configs": flattenReservedInstanceMetricConfig(utils.PathSearch("metric_configs",
				tracticsConfig, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenReservedInstances(reservedInstances []interface{}) []map[string]interface{} {
	if len(reservedInstances) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(reservedInstances))
	for _, reservedInstance := range reservedInstances {
		result = append(result, map[string]interface{}{
			"count":          utils.PathSearch("min_count", reservedInstance, nil),
			"idle_mode":      utils.PathSearch("idle_mode", reservedInstance, nil),
			"qualifier_name": utils.PathSearch("qualifier_name", reservedInstance, nil),
			"qualifier_type": utils.PathSearch("qualifier_type", reservedInstance, nil),
			"tactics_config": flattenReservedInstanceTracticsConfig(utils.PathSearch("tactics_config",
				reservedInstance, make(map[string]interface{})).(map[string]interface{})),
		})
	}
	return result
}

func setFunctionFieldApp(d *schema.ResourceData, app string) error {
	// If the deprecated parameter package is not set, saving value to the parameter app.
	if _, ok := d.GetOk("package"); !ok {
		return d.Set("app", app)
	}
	return d.Set("package", app)
}

func setFunctionFieldAgency(d *schema.ResourceData, agency string) error {
	// If the deprecated parameter xrole is not set, saving value to the parameter agency.
	if _, ok := d.GetOk("xrole"); !ok {
		return d.Set("agency", agency)
	}
	return d.Set("xrole", agency)
}

func resourceFunctionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                   = meta.(*config.Config)
		region                = cfg.GetRegion(d)
		funcUrn               = d.Id()
		funcUrnWithoutVersion = parseFunctionUrnWithoutVersion(funcUrn)
	)
	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	function, err := GetFunctionMetadata(client, funcUrnWithoutVersion)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying function (%s) metadata", funcUrnWithoutVersion))
	}

	log.Printf("[DEBUG] Retrieved Function %s: %+v", funcUrnWithoutVersion, function)
	mErr := multierror.Append(
		// Required parameters.
		d.Set("name", utils.PathSearch("func_name", function, nil)),
		d.Set("runtime", utils.PathSearch("runtime", function, nil)),
		d.Set("timeout", utils.PathSearch("timeout", function, nil)),
		d.Set("memory_size", utils.PathSearch("memory_size", function, nil)),
		// Optional parameters but required in documentation.
		setFunctionFieldApp(d, utils.PathSearch("package", function, "").(string)),
		d.Set("handler", utils.PathSearch("handler", function, nil)),
		d.Set("code_type", utils.PathSearch("code_type", function, nil)),
		// Optional parameters.
		d.Set("description", utils.PathSearch("description", function, nil)),
		d.Set("functiongraph_version", utils.PathSearch("type", function, nil)),
		d.Set("code_url", utils.PathSearch("code_url", function, nil)),
		d.Set("code_filename", utils.PathSearch("code_filename", function, nil)),
		d.Set("user_data", utils.PathSearch("user_data", function, nil)),
		setFunctionFieldAgency(d, utils.PathSearch("xrole", function, "").(string)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", function, nil)),
		d.Set("custom_image", flattenFgsCustomImage(utils.PathSearch("custom_image",
			function, make(map[string]interface{})).(map[string]interface{}))),
		d.Set("gpu_type", utils.PathSearch("gpu_type", function, nil)),
		d.Set("gpu_memory", utils.PathSearch("gpu_memory", function, nil)),
		d.Set("lts_custom_tag", utils.PathSearch("lts_custom_tag", function, nil)),
		d.Set("pre_stop_handler", utils.PathSearch("pre_stop_handler", function, nil)),
		d.Set("pre_stop_timeout", utils.PathSearch("pre_stop_timeout", function, nil)),
		d.Set("app_agency", utils.PathSearch("app_xrole", function, nil)),
		d.Set("depend_list", utils.PathSearch("depend_version_list", function, nil)),
		d.Set("initializer_handler", utils.PathSearch("initializer_handler", function, nil)),
		d.Set("initializer_timeout", utils.PathSearch("initializer_timeout", function, nil)),
		d.Set("max_instance_num", strconv.Itoa(int(utils.PathSearch("strategy_config.concurrency",
			function, float64(0)).(float64)))),
		d.Set("concurrency_num", int(utils.PathSearch("strategy_config.concurrent_num",
			function, float64(0)).(float64))),
		d.Set("dns_list", utils.PathSearch("domain_names", function, nil)),
		d.Set("vpc_id", utils.PathSearch("func_vpc.vpc_id", function, nil)),
		d.Set("network_id", utils.PathSearch("func_vpc.subnet_id", function, nil)),
		d.Set("mount_user_id", utils.PathSearch("mount_config.mount_user.user_id", function, nil)),
		d.Set("mount_user_group_id", utils.PathSearch("mount_config.mount_user.user_group_id", function, nil)),
		d.Set("func_mounts", flattenFuncionMounts(utils.PathSearch("mount_config.func_mounts",
			function, make([]interface{}, 0)).([]interface{}))),
		d.Set("enable_dynamic_memory", utils.PathSearch("enable_dynamic_memory", function, nil)),
		d.Set("is_stateful_function", utils.PathSearch("is_stateful_function", function, nil)),
		d.Set("network_controller", flattenFunctionNetworkController(utils.PathSearch("network_controller", function, nil))),
		d.Set("peering_cidr", utils.PathSearch("peering_cidr", function, nil)),
		d.Set("enable_auth_in_header", utils.PathSearch("enable_auth_in_header", function, nil)),
		d.Set("enable_class_isolation", utils.PathSearch("enable_class_isolation", function, nil)),
		d.Set("ephemeral_storage", utils.PathSearch("ephemeral_storage", function, nil)),
		d.Set("heartbeat_handler", utils.PathSearch("heartbeat_handler", function, nil)),
		d.Set("restore_hook_handler", utils.PathSearch("restore_hook_handler", function, nil)),
		d.Set("restore_hook_timeout", utils.PathSearch("restore_hook_timeout", function, nil)),
		d.Set("enable_lts_log", utils.PathSearch("enable_lts_log", function, nil)),
		d.Set("user_data_encrypt_kms_key_id", utils.PathSearch("user_data_encrypt_kms_key_id", function, nil)),
		d.Set("code_encrypt_kms_key_id", utils.PathSearch("code_encrypt_kms_key_id", function, nil)),
		d.Set("tags", d.Get("tags")),
		// Attributes.
		d.Set("urn", utils.PathSearch("func_urn", function, nil)),
		d.Set("version", utils.PathSearch("version", function, nil)),
	)

	// The metadata API for obtaining function does not return the names of the log group and log stream.
	logConfiguration, err := getFunctionLogConfiguration(client, funcUrn)
	if err != nil {
		log.Printf("[ERROR] Unable to get log configuration: %s", err)
	}
	mErr = multierror.Append(mErr,
		d.Set("log_group_name", utils.PathSearch("group_name", logConfiguration, nil)),
		d.Set("log_stream_name", utils.PathSearch("stream_name", logConfiguration, nil)),
		d.Set("log_group_id", utils.PathSearch("group_id", logConfiguration, nil)),
		d.Set("log_stream_id", utils.PathSearch("stream_id", logConfiguration, nil)),
	)

	versionConfig, err := flattenFunctionVersions(client, funcUrnWithoutVersion)
	if err != nil {
		// Not all regions support the version related API calls.
		log.Printf("[ERROR] Unable to parsing the function versions: %s", err)
	}
	mErr = multierror.Append(mErr, d.Set("versions", versionConfig))

	reservedInstances, err := getFunctionReservedInstances(client, funcUrn)
	if err != nil {
		return diag.Errorf("error retrieving function reserved instance: %s", err)
	}
	mErr = multierror.Append(mErr, d.Set("reserved_instances", flattenReservedInstances(reservedInstances)))

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting function fields: %s", err)
	}

	return nil
}

func getFunctionLogConfiguration(client *golangsdk.ServiceClient, functionUrn string) (interface{}, error) {
	httpUrl := "v2/{project_id}/fgs/functions/{function_urn}/lts-log-detail"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{function_urn}", functionUrn)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func resourceFunctionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                             = meta.(*config.Config)
		region                          = cfg.GetRegion(d)
		funcUrnWithoutVersion           = parseFunctionUrnWithoutVersion(d.Id())
		functionMetadataObjectParamKeys = []string{
			"lts_custom_tag",
		}
		rawConfig = d.GetRawConfig()
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	// lintignore:R019
	if d.HasChanges("code_type", "code_url", "code_filename", "depend_list", "func_code") {
		err := updateFunctionCode(client, d, funcUrnWithoutVersion)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// lintignore:R019
	if d.HasChanges("app", "handler", "memory_size", "timeout", "encrypted_user_data",
		"user_data", "agency", "app_agency", "description", "initializer_handler", "initializer_timeout",
		"vpc_id", "network_id", "dns_list", "mount_user_id", "mount_user_group_id", "func_mounts", "custom_image",
		"log_group_id", "log_stream_id", "log_group_name", "log_stream_name", "concurrency_num", "gpu_memory", "gpu_type",
		"enable_dynamic_memory", "is_stateful_function", "network_controller", "enterprise_project_id", "peering_cidr",
		"enable_auth_in_header", "enable_class_isolation", "ephemeral_storage", "heartbeat_handler", "restore_hook_handler",
		"restore_hook_timeout", "lts_custom_tag", "enable_lts_log", "user_data_encrypt_kms_key_id") {
		params := buildUpdateFunctionMetadataBodyParams(cfg, d)
		if d.HasChanges("log_group_id", "log_stream_id", "log_group_name", "log_stream_name", "enable_lts_log") {
			params["enable_lts_log"] = utils.GetNestedObjectFromRawConfig(rawConfig, "enable_lts_log")
			params["log_config"] = buildFunctionLogConfig(rawConfig)
		}
		err := updateFunctionMetadata(client, funcUrnWithoutVersion, params)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// If the request is successful, obtain the values ​​of all JSON|object parameters first and save them to the
	// corresponding '_origin' attributes for subsequent determination and construction of the request body during
	// next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	err = utils.RefreshObjectParamOriginValues(d, functionMetadataObjectParamKeys)
	if err != nil {
		return diag.Errorf("unable to refresh the origin values: %s", err)
	}

	if d.HasChange("max_instance_num") {
		// If the maximum number of instances is omitted (after type conversion, the value is zero), means this feature
		// is disabled.
		maxInstanceNum, _ := strconv.Atoi(d.Get("max_instance_num").(string))
		err = updateFunctionMaxInstanceNum(client, funcUrnWithoutVersion, maxInstanceNum)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		if err = updateFunctionTags(client, d, funcUrnWithoutVersion); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("versions") {
		if err = updateFunctionVersions(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("reserved_instances") {
		if err = updateFunctionReservedInstances(client, d, funcUrnWithoutVersion); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceFunctionRead(ctx, d, meta)
}

func deleteFunctionTags(client *golangsdk.ServiceClient, functionUrn string, tags map[string]interface{}) error {
	if len(tags) < 1 {
		return nil
	}

	httpUrl := "v2/{project_id}/functions/{function_urn}/tags/delete"

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{function_urn}", functionUrn)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		OkCodes:  []int{204},
		JSONBody: buildFunctionTagsBodyParams(tags),
	}
	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return fmt.Errorf("failed to delete the function tags: %s", err)
	}
	return nil
}

func updateFunctionTags(client *golangsdk.ServiceClient, d *schema.ResourceData, functionUrn string) error {
	oldVal, newVal := d.GetChange("tags")
	oldTags := oldVal.(map[string]interface{})
	newTags := newVal.(map[string]interface{})

	if len(oldTags) > 0 {
		if err := deleteFunctionTags(client, functionUrn, oldTags); err != nil {
			return err
		}
	}

	if len(newTags) > 0 {
		if err := createFunctionTags(client, functionUrn, newTags); err != nil {
			return err
		}
	}
	return nil
}

// Does not allow to delete the latest version of the function (URN with the latest version).
// To delete the entire function (including all versions), the URN should be without any version number/alias,
// such as 'urn:fss:cn-north-4:xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx:function:default:xxxxxx'
func deleteFunctionOrVersion(client *golangsdk.ServiceClient, functionUrn string) error {
	httpUrl := "v2/{project_id}/fgs/functions/{function_urn}"

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{function_urn}", functionUrn)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return err
	}
	return nil
}

func deleteFunctionVersionAlias(client *golangsdk.ServiceClient, functionUrn, aliasName string) error {
	httpUrl := "v2/{project_id}/fgs/functions/{function_urn}/aliases/{alias_name}"

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{function_urn}", functionUrn)
	deletePath = strings.ReplaceAll(deletePath, "{alias_name}", aliasName)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return fmt.Errorf("error deleting function version alias (%s): %s", aliasName, err)
	}
	return nil
}

func createFunctionVersion(client *golangsdk.ServiceClient, functionUrn, versionName, versionDesc string) error {
	httpUrl := "v2/{project_id}/fgs/functions/{function_urn}/versions"

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{function_urn}", functionUrn)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"version":     versionName,
			"description": versionDesc,
		},
	}

	_, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return fmt.Errorf("failed to create the function version: %s", err)
	}
	return nil
}

func createFunctionVersionAlias(client *golangsdk.ServiceClient, functionUrn, versionName string, aliasCfg interface{}) error {
	if aliasCfg == nil {
		return nil
	}

	httpUrl := "v2/{project_id}/fgs/functions/{function_urn}/aliases"

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{function_urn}", functionUrn)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(map[string]interface{}{
			"version":     versionName,
			"name":        utils.ValueIgnoreEmpty(utils.PathSearch("name", aliasCfg, nil)),
			"description": utils.ValueIgnoreEmpty(utils.PathSearch("description", aliasCfg, nil)),
			"additional_version_weights": utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("additional_version_weights",
				aliasCfg, "").(string))),
			"additional_version_strategy": utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("additional_version_strategy",
				aliasCfg, "").(string))),
		}),
	}

	_, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return fmt.Errorf("failed to create the function version alias: %s", err)
	}
	return nil
}

func updateFunctionVersions(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		functionUrn = parseFunctionUrnWithoutVersion(d.Id())

		oldVal, newVal  = d.GetChange("versions")
		oldVersions     = oldVal.(*schema.Set).List()
		newVersions     = newVal.(*schema.Set).List()
		newVersionNames = utils.PathSearch("[*].name", newVersions, make([]interface{}, 0)).([]interface{})
		oldVersionNames = utils.PathSearch("[*].name", oldVersions, make([]interface{}, 0)).([]interface{})

		err error
	)

	// Version         -> null:                Remove version
	// Version + Alias -> null:                Remove version
	// Version + Alias -> Version:             Remove alias
	// Version + Alias -> Version + New Alias: Remove alias before new alias create
	// Do not use the difference function of the type schema set, and use the custom compare function as follows.
	for _, oldVersion := range oldVersions {
		versionName := utils.PathSearch("name", oldVersion, nil)
		versionDesc := utils.PathSearch("description", oldVersion, nil)
		// Check if the version (by name and description) has changed and decide whether to delete it and the latest
		// version is check in particular, because the latest version cannot be deleted, so, if the version structure
		// has been removed (delete alias), skip the version delete logic and just only delete the corresponding alias.
		// For versions with other names, deleting the version will also delete the alias, regardless of whether it has
		// an alias.
		if fmt.Sprint(versionName) != "latest" && (!utils.SliceContains(newVersionNames, versionName) ||
			fmt.Sprint(versionDesc) != utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0].description", versionName), newVersions, "").(string)) {
			// Delete a specified function version.
			err = deleteFunctionOrVersion(client, fmt.Sprintf("%s:%s", functionUrn, versionName))
			if err != nil {
				return err
			}
			continue
		}

		aliases := utils.PathSearch("aliases", oldVersion, make([]interface{}, 0)).([]interface{})
		// If the version info does not change, check if the alias information is removed or updated and decide whether
		// to delete current alias, because these is no API to update the alias, just have APIs for create and delete.
		// Any configuration update needs to be implemented by delete function (if the alias configuration has been
		// configured in the history list) before new configuration create.
		if len(aliases) > 0 && !reflect.DeepEqual(aliases, utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0].aliases", versionName),
			newVersions, make([]interface{}, 0)).([]interface{})) {
			// Delete an alias from a specified function version.
			err = deleteFunctionVersionAlias(client, functionUrn, utils.PathSearch("[0].name", aliases, "").(string))
			if err != nil {
				return err
			}
		}
	}

	// null            -> Version:             Create version
	// null            -> Version + Alias:     Create version + Create alias
	// Version         -> Version + Alias:     Create alias
	// Version + Alias -> Version + New Alias: Create alias after old alias removed
	for _, newVersion := range newVersions {
		versionName := utils.PathSearch("name", newVersion, nil)
		versionDesc := utils.PathSearch("description", newVersion, "").(string)
		// Check if the new version (by name and description) is supported and decide whether to create a new version
		// and the latest version is check in particular, because the create version name cannot be 'latest', so, if the
		// version name has been updated (ignore alias change), new a new version first.
		isCreateNewVersion := versionName != "latest" && (!utils.SliceContains(oldVersionNames, versionName) ||
			fmt.Sprint(versionDesc) != utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0].description", versionName), oldVersions, "").(string))
		if isCreateNewVersion {
			// Create a new function version.
			err := createFunctionVersion(client, functionUrn, fmt.Sprint(versionName), versionDesc)
			if err != nil {
				return err
			}
		}

		aliases := utils.PathSearch("aliases", newVersion, make([]interface{}, 0)).([]interface{})
		// If the version is supported or alias is changed (also can be supported), it means that a new alias needs to
		// be created according to the alias configuration in the current script (the latter will delete the old alias
		// configuration first).
		if isCreateNewVersion || (len(aliases) > 0 && !reflect.DeepEqual(aliases,
			utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0].aliases", versionName), oldVersions, make([]interface{}, 0)).([]interface{}))) {
			// Create a new alias under a specified function version.
			err := createFunctionVersionAlias(client, functionUrn, fmt.Sprint(versionName), utils.PathSearch("[0]", aliases, nil))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceFunctionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                   = meta.(*config.Config)
		region                = cfg.GetRegion(d)
		funcUrnWithoutVersion = parseFunctionUrnWithoutVersion(d.Id())
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	// Using function URN delete function without the version means delete the function and all its versions (aliases).
	err = deleteFunctionOrVersion(client, funcUrnWithoutVersion)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting function (%s)", funcUrnWithoutVersion))
	}
	return nil
}

// The function ID consists of the function URN and the current version.
// Some requests require the URN information without the version, so this function is used to extract it.
func parseFunctionUrnWithoutVersion(urn string) string {
	// The format of the function URN is: 'urn:fss:cn-north-4:xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx:function:default:xxxxxx'.
	index := strings.LastIndex(urn, ":")
	if index != -1 {
		urn = urn[0:index]
	}
	return urn
}
