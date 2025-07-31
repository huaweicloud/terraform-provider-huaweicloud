package servicestage

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	v3ComponentNotFoundCodes = []string{
		"SVCSTG.00100401",
	}

	componentJsonParamKeys = []string{
		"source",
		"build",
		"command",
		"tomcat_opts",
		"deploy_strategy.0.rolling_release",
		"deploy_strategy.0.gray_release",
		"update_strategy",
	}

	componentNonUpdatableParams = []string{
		"application_id",
		"environment_id",
		"name",
		"runtime_stack",
		"runtime_stack.*.name",
		"runtime_stack.*.type",
		"runtime_stack.*.deploy_mode",
		"runtime_stack.*.version",
		"replica",
	}
)

// @API ServiceStage POST /v3/{project_id}/cas/applications/{application_id}/components
// @API ServiceStage GET /v3/{project_id}/cas/jobs/{job_id}
// @API ServiceStage GET /v3/{project_id}/cas/applications/{application_id}/components/{component_id}
// @API ServiceStage PUT /v3/{project_id}/cas/applications/{application_id}/components/{component_id}
// @API ServiceStage DELETE /v3/{project_id}/cas/applications/{application_id}/components/{component_id}
func ResourceV3Component() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3ComponentCreate,
		ReadContext:   resourceV3ComponentRead,
		UpdateContext: resourceV3ComponentUpdate,
		DeleteContext: resourceV3ComponentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV3ComponentImportState,
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(componentNonUpdatableParams),
			config.MergeDefaultTags(),
		),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the component is located.`,
			},
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The application ID to which the component belongs.`,
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The environment ID where the component is deployed.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the component.`,
			},
			"runtime_stack": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The stack name.`,
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The stack type.`,
						},
						"deploy_mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The deploy mode of the stack.`,
						},
						"version": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The stack version.`,
						},
					},
				},
				Description: "The configuration of the runtime stack.",
			},
			"source": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: utils.SuppressObjectDiffs(),
				Description:      `The source configuration of the component, in JSON format.`,
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The version of the component.`,
			},
			"refer_resources": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The resource ID.`,
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The resource type.`,
						},
						"parameters": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							Description:  `The resource parameters, in JSON format.`,
						},
					},
				},
				Description: `The configuration of the reference resources.`,
			},
			"config_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The configuration mode of the component.`,
			},
			"workload_content": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The workload content of the component, in JSON format.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description of the component.`,
			},
			"build": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: utils.SuppressObjectDiffs(),
				Description:      `The build configuration of the component, in JSON format.`,
			},
			"replica": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The replica number of the component.`,
			},
			"limit_cpu": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Computed:    true,
				Description: `The maximum number of the CPU limit.`,
			},
			"limit_memory": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Computed:    true,
				Description: `The maximum number of the memory limit.`,
			},
			"request_cpu": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Computed:    true,
				Description: `The number of the CPU request resources.`,
			},
			"request_memory": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Computed:    true,
				Description: `The number of the memory request resources.`,
			},
			"envs": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the environment variable.`,
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The value of the environment variable.`,
						},
					},
				},
				Description: "The configuration of the environment variables.",
			},
			"storages": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The type of the data storage.`,
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the disk where the data is stored.`,
						},
						"parameters": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsJSON,
							Description:  `The information corresponding to the specific types of data storage, in JSON format.`,
						},
						"mounts": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The mount path.`,
									},
									"sub_path": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The sub mount path.`,
									},
									"read_only": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: `Whether the disk mount is read-only.`,
									},
								},
							},
							Description: `The configuration of the disk mounts.`,
						},
					},
				},
				Description: "The storage configuration.",
			},
			"deploy_strategy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The deploy type.`,
						},
						"rolling_release": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ValidateFunc:     validation.StringIsJSON,
							DiffSuppressFunc: utils.SuppressObjectDiffs(),
							Description:      `The rolling release parameters, in JSON format.`,
						},
						"gray_release": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ValidateFunc:     validation.StringIsJSON,
							DiffSuppressFunc: utils.SuppressObjectDiffs(),
							Description:      `The gray release parameters, in JSON format.`,
						},
						"rolling_release_origin": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The script configuration value of this change is also the original value used for comparison with
the new value next time the change is made. The corresponding parameter name is 'deploy_strategy.0.rolling_release'.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
						"gray_release_origin": {
							Type:     schema.TypeString,
							Computed: true,
							Description: utils.SchemaDesc(
								`The script configuration value of this change is also the original value used for comparison with
the new value next time the change is made. The corresponding parameter name is 'deploy_strategy.0.gray_release'.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
					},
				},
				Description: `The configuration of the deploy strategy.`,
			},
			// Most of the strategy configuration inputs for component deployment and upgrades have been changed from
			// deploy_strategy to update_strategy now.
			"update_strategy": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringIsJSON,
				Description:      `The configuration of the update strategy, in JSON format.`,
				DiffSuppressFunc: utils.SuppressObjectDiffs(),
			},
			"command": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: utils.SuppressObjectDiffs(),
				Description:      `The start commands of the component, in JSON format.`,
			},
			"post_start": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        componentLifecycleSchema(),
				Description: `The post start configuration.`,
			},
			"pre_stop": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        componentLifecycleSchema(),
				Description: `The pre stop configuration.`,
			},
			"mesher": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The process listening port.`,
						},
					},
				},
				Description: `The configuration of the access mesher.`,
			},
			"timezone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The time zone in which the component runs.`,
			},
			"jvm_opts": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The JVM parameters of the component.`,
			},
			"tomcat_opts": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: utils.SuppressObjectDiffs(),
				Description:      `The configuration of the tomcat server.`,
			},
			"logs": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The log path of the container.`,
						},
						"rotate": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The interval for dumping logs.`,
						},
						"host_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The mounted host path.`,
						},
						"host_extend_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The extension path of the host.`,
						},
					},
				},
				Description: `The configuration of the logs collection.`,
			},
			"custom_metric": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The collection path.`,
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The collection port.`,
						},
						"dimensions": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The monitoring dimension.`,
						},
					},
				},
				Description: `The configuration of the monitor metric.`,
			},
			"affinity": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        componentAffinitySchema(),
				Description: `The affinity configuration of the component.`,
			},
			"anti_affinity": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        componentAffinitySchema(),
				Description: `The anti-affinity configuration of the component.`,
			},
			"liveness_probe": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        componentProbeSchema(),
				Description: "The liveness probe configuration of the component.",
			},
			"readiness_probe": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        componentProbeSchema(),
				Description: "The readiness probe configuration of the component.",
			},
			"external_accesses": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The protocol of the external access.`,
						},
						"address": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The address of the external access.`,
						},
						"forward_port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `The forward port of the external access.`,
						},
					},
				},
				Description: "The configuration of the external accesses.",
			},
			"tags": common.TagsSchema(
				`The key/value pairs to associate with the component.`,
			),
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the component.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the component, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the component, in RFC3339 format.`,
			},
			// Internal parameters/attributes.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"source_origin": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
 the new value next time the change is made. The corresponding parameter name is 'source'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"build_origin": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
 the new value next time the change is made. The corresponding parameter name is 'build'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"update_strategy_origin": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
 the new value next time the change is made. The corresponding parameter name is 'update_strategy'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"command_origin": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
 the new value next time the change is made. The corresponding parameter name is 'command'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"tomcat_opts_origin": {
				Type:     schema.TypeString,
				Computed: true,
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
 the new value next time the change is made. The corresponding parameter name is 'tomcat_opts'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func componentLifecycleSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The processing method.`,
			},
			"scheme": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The HTTP request type.`,
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The host (IP) of the lifecycle configuration.`,
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The port number of the lifecycle configuration.`,
			},
			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The request path of the lifecycle configuration.`,
			},
			"command": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The command list of the lifecycle configuration.`,
			},
		},
	}
	return &sc
}

func componentAffinitySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"condition": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The condition type of the (anti) affinity rule.`,
			},
			"kind": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The kind of the (anti) affinity rule.`,
			},
			"match_expressions": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The key of the match rule.`,
						},
						"operation": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The operation of the match rule.`,
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The value of the match rule.`,
						},
					},
				},
				Description: "The list of the match rules for (anti) affinity.",
			},
			"weight": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The weight of the (anti) affinity rule.`,
			},
		},
	}
	return &sc
}

func componentProbeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the probe.`,
			},
			"delay": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The delay time of the probe.`,
			},
			"timeout": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The timeout of the probe.`,
			},
			"scheme": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The scheme type of the probe.`,
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The host of the probe.`,
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The port of the probe.`,
			},
			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The path of the probe.`,
			},
			"command": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The command list of the probe.`,
			},
		},
	}
	return &sc
}

func buildV3ComponentRuntimeStackConfig(runtimeStacks []interface{}) map[string]interface{} {
	if len(runtimeStacks) < 1 {
		return nil
	}

	runtimeStack := runtimeStacks[0]
	return map[string]interface{}{
		"name":        utils.PathSearch("name", runtimeStack, nil),
		"type":        utils.PathSearch("type", runtimeStack, nil),
		"deploy_mode": utils.PathSearch("deploy_mode", runtimeStack, nil),
		"version":     utils.PathSearch("version", runtimeStack, nil),
	}
}

func buildV3ComponentEnvVariables(variables *schema.Set) []interface{} {
	if variables.Len() < 1 {
		return nil
	}

	result := make([]interface{}, 0, variables.Len())
	for _, variable := range variables.List() {
		result = append(result, map[string]interface{}{
			"name":  utils.PathSearch("name", variable, nil),
			"value": utils.PathSearch("value", variable, nil),
		})
	}

	return result
}

func buildV3ComponentStorageMounts(mounts *schema.Set) []interface{} {
	if mounts.Len() < 1 {
		return nil
	}

	result := make([]interface{}, 0, mounts.Len())
	for _, mount := range mounts.List() {
		result = append(result, utils.RemoveNil(map[string]interface{}{
			"path":      utils.PathSearch("path", mount, nil),
			"sub_path":  utils.PathSearch("sub_path", mount, nil),
			"read_only": utils.PathSearch("read_only", mount, nil),
		}))
	}

	return result
}

func buildV3ComponentStorages(storages *schema.Set) []interface{} {
	if storages.Len() < 1 {
		return nil
	}

	result := make([]interface{}, 0, storages.Len())
	for _, storage := range storages.List() {
		result = append(result, map[string]interface{}{
			"type":       utils.PathSearch("type", storage, nil),
			"name":       utils.PathSearch("name", storage, nil),
			"parameters": utils.StringToJson(utils.PathSearch("parameters", storage, "").(string)),
			"mounts":     buildV3ComponentStorageMounts(utils.PathSearch("mounts", storage, schema.NewSet(schema.HashString, nil)).(*schema.Set)),
		})
	}

	return result
}

func buildV3ComponentDeployStrategy(strategies []interface{}) map[string]interface{} {
	if len(strategies) < 1 {
		return nil
	}

	strategy := strategies[0]
	return map[string]interface{}{
		"type":            utils.PathSearch("type", strategy, nil),
		"rolling_release": utils.StringToJson(utils.PathSearch("rolling_release", strategy, "").(string)),
		"gray_release":    utils.StringToJson(utils.PathSearch("gray_release", strategy, "").(string)),
	}
}

func buildV3ComponentLifecycle(lifecycles []interface{}) map[string]interface{} {
	if len(lifecycles) < 1 {
		return nil
	}

	lifecycle := lifecycles[0]
	return map[string]interface{}{
		"type":   utils.PathSearch("type", lifecycle, nil),
		"scheme": utils.ValueIgnoreEmpty(utils.PathSearch("scheme", lifecycle, nil)),
		"host":   utils.ValueIgnoreEmpty(utils.PathSearch("host", lifecycle, nil)),
		"port":   utils.ValueIgnoreEmpty(utils.PathSearch("port", lifecycle, nil)),
		"path":   utils.ValueIgnoreEmpty(utils.PathSearch("path", lifecycle, nil)),
		"command": utils.ValueIgnoreEmpty(utils.ExpandToStringListBySet(utils.PathSearch("command", lifecycle,
			schema.NewSet(schema.HashString, nil)).(*schema.Set))),
	}
}

func buildV3ComponentMesher(meshers []interface{}) map[string]interface{} {
	if len(meshers) < 1 {
		return nil
	}

	mesher := meshers[0]
	return map[string]interface{}{
		"port": utils.PathSearch("port", mesher, nil),
	}
}

func buildV3ComponentLogs(logs *schema.Set) []interface{} {
	if logs.Len() < 1 {
		return nil
	}

	result := make([]interface{}, 0, logs.Len())
	for _, v := range logs.List() {
		result = append(result, map[string]interface{}{
			"log_path":         utils.PathSearch("log_path", v, nil),
			"rotate":           utils.PathSearch("rotate", v, nil),
			"host_path":        utils.PathSearch("host_path", v, nil),
			"host_extend_path": utils.PathSearch("host_extend_path", v, nil),
		})
	}

	return result
}

func buildV3ComponentCustomMetric(metrics []interface{}) map[string]interface{} {
	if len(metrics) < 1 {
		return nil
	}

	metric := metrics[0]
	return map[string]interface{}{
		"path":       utils.ValueIgnoreEmpty(utils.PathSearch("path", metric, nil)),
		"port":       utils.ValueIgnoreEmpty(utils.PathSearch("port", metric, nil)),
		"dimensions": utils.ValueIgnoreEmpty(utils.PathSearch("dimensions", metric, nil)),
	}
}

func buildV3ComponentAffinityMatchExpressions(matchRules *schema.Set) []interface{} {
	if matchRules.Len() < 1 {
		return nil
	}

	result := make([]interface{}, 0, matchRules.Len())
	for _, rule := range matchRules.List() {
		result = append(result, map[string]interface{}{
			"key":       utils.PathSearch("key", rule, nil),
			"operation": utils.PathSearch("operation", rule, nil),
			"value":     utils.PathSearch("value", rule, nil),
		})
	}

	return result
}

func buildV3ComponentAffinity(affinityRules *schema.Set) []interface{} {
	if affinityRules.Len() < 1 {
		return nil
	}

	result := make([]interface{}, 0, affinityRules.Len())
	for _, rule := range affinityRules.List() {
		result = append(result, map[string]interface{}{
			"condition": utils.PathSearch("condition", rule, nil),
			"kind":      utils.PathSearch("kind", rule, nil),
			"match_expressions": utils.ValueIgnoreEmpty(buildV3ComponentAffinityMatchExpressions(utils.PathSearch("match_expressions", rule,
				schema.NewSet(schema.HashString, nil)).(*schema.Set))),
			"weight": utils.PathSearch("weight", rule, nil),
		})
	}

	return result
}

func buildV3ComponentProbeConfiguration(probeConfigs []interface{}) map[string]interface{} {
	if len(probeConfigs) < 1 {
		return nil
	}

	probeConfig := probeConfigs[0]
	return map[string]interface{}{
		"type":    utils.ValueIgnoreEmpty(utils.PathSearch("type", probeConfig, nil)),
		"delay":   utils.ValueIgnoreEmpty(utils.PathSearch("delay", probeConfig, nil)),
		"timeout": utils.ValueIgnoreEmpty(utils.PathSearch("timeout", probeConfig, nil)),
		"scheme":  utils.ValueIgnoreEmpty(utils.PathSearch("scheme", probeConfig, nil)),
		"host":    utils.ValueIgnoreEmpty(utils.PathSearch("host", probeConfig, nil)),
		"port":    utils.ValueIgnoreEmpty(utils.PathSearch("port", probeConfig, nil)),
		"path":    utils.ValueIgnoreEmpty(utils.PathSearch("path", probeConfig, nil)),
		"command": utils.ValueIgnoreEmpty(utils.ExpandToStringListBySet(utils.PathSearch("command",
			probeConfig, schema.NewSet(schema.HashString, nil)).(*schema.Set))),
	}
}

func buildV3ComponentReferResources(refResources *schema.Set) []interface{} {
	if refResources.Len() < 1 {
		return nil
	}

	result := make([]interface{}, 0, refResources.Len())
	for _, refRsource := range refResources.List() {
		result = append(result, utils.RemoveNil(map[string]interface{}{
			"id":         utils.ValueIgnoreEmpty(utils.PathSearch("id", refRsource, nil)),
			"type":       utils.ValueIgnoreEmpty(utils.PathSearch("type", refRsource, nil)),
			"parameters": utils.StringToJson(utils.PathSearch("parameters", refRsource, "").(string)),
		}))
	}

	return result
}

func buildV3ComponentExternalAccesses(accesses *schema.Set) []interface{} {
	if accesses.Len() < 1 {
		return nil
	}

	result := make([]interface{}, 0, accesses.Len())
	for _, access := range accesses.List() {
		result = append(result, utils.RemoveNil(map[string]interface{}{
			"protocol":     utils.PathSearch("protocol", access, nil),
			"address":      utils.ValueIgnoreEmpty(utils.PathSearch("address", access, nil)),
			"forward_port": utils.ValueIgnoreEmpty(utils.PathSearch("forward_port", access, nil)),
		}))
	}

	return result
}

func buildV3ComponentCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Required parameters.
		"name":            d.Get("name").(string),
		"runtime_stack":   utils.ValueIgnoreEmpty(buildV3ComponentRuntimeStackConfig(d.Get("runtime_stack").([]interface{}))),
		"source":          utils.StringToJson(d.Get("source").(string)),
		"version":         d.Get("version").(string),
		"refer_resources": utils.ValueIgnoreEmpty(buildV3ComponentReferResources(d.Get("refer_resources").(*schema.Set))),
		// Optional parameters.
		"environment_id":    d.Get("environment_id").(string),
		"config_mode":       utils.ValueIgnoreEmpty(d.Get("config_mode")),
		"workload_content":  utils.ValueIgnoreEmpty(d.Get("workload_content")),
		"description":       utils.ValueIgnoreEmpty(d.Get("description")),
		"build":             utils.StringToJson(d.Get("build").(string)),
		"replica":           utils.ValueIgnoreEmpty(d.Get("replica").(int)),
		"limit_cpu":         utils.ValueIgnoreEmpty(d.Get("limit_cpu").(float64)),
		"limit_memory":      utils.ValueIgnoreEmpty(d.Get("limit_memory").(float64)),
		"request_cpu":       utils.ValueIgnoreEmpty(d.Get("request_cpu").(float64)),
		"request_memory":    utils.ValueIgnoreEmpty(d.Get("request_memory").(float64)),
		"envs":              utils.ValueIgnoreEmpty(buildV3ComponentEnvVariables(d.Get("envs").(*schema.Set))),
		"storages":          utils.ValueIgnoreEmpty(buildV3ComponentStorages(d.Get("storages").(*schema.Set))),
		"deploy_strategy":   utils.ValueIgnoreEmpty(buildV3ComponentDeployStrategy(d.Get("deploy_strategy").([]interface{}))),
		"update_strategy":   utils.StringToJson(d.Get("update_strategy").(string)),
		"command":           utils.StringToJson(d.Get("command").(string)),
		"post_start":        utils.ValueIgnoreEmpty(buildV3ComponentLifecycle(d.Get("post_start").([]interface{}))),
		"pre_stop":          utils.ValueIgnoreEmpty(buildV3ComponentLifecycle(d.Get("pre_stop").([]interface{}))),
		"mesher":            utils.ValueIgnoreEmpty(buildV3ComponentMesher(d.Get("mesher").([]interface{}))),
		"timezone":          utils.ValueIgnoreEmpty(d.Get("timezone").(string)),
		"jvm_opts":          utils.ValueIgnoreEmpty(d.Get("jvm_opts").(string)),
		"tomcat_opts":       utils.StringToJson(d.Get("tomcat_opts").(string)),
		"logs":              utils.ValueIgnoreEmpty(buildV3ComponentLogs(d.Get("logs").(*schema.Set))),
		"custom_metric":     utils.ValueIgnoreEmpty(buildV3ComponentCustomMetric(d.Get("custom_metric").([]interface{}))),
		"affinity":          utils.ValueIgnoreEmpty(buildV3ComponentAffinity(d.Get("affinity").(*schema.Set))),
		"anti_affinity":     utils.ValueIgnoreEmpty(buildV3ComponentAffinity(d.Get("anti_affinity").(*schema.Set))),
		"liveness_probe":    utils.ValueIgnoreEmpty(buildV3ComponentProbeConfiguration(d.Get("liveness_probe").([]interface{}))),
		"readiness_probe":   utils.ValueIgnoreEmpty(buildV3ComponentProbeConfiguration(d.Get("readiness_probe").([]interface{}))),
		"external_accesses": utils.ValueIgnoreEmpty(buildV3ComponentExternalAccesses(d.Get("external_accesses").(*schema.Set))),
		"labels":            utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
	}
}

func queryV3Job(client *golangsdk.ServiceClient, jobId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/cas/jobs/{job_id}"

	queryPath := client.Endpoint + httpUrl
	queryPath = strings.ReplaceAll(queryPath, "{project_id}", client.ProjectID)
	queryPath = strings.ReplaceAll(queryPath, "{job_id}", jobId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", queryPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func jobStatusRefreshFunc(client *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := queryV3Job(client, jobId)
		if err != nil {
			return respBody, "ERROR", err
		}

		jobStatus := utils.PathSearch("job.execution_status", respBody, "").(string)
		if utils.StrSliceContains([]string{"FAILED", "UNKNOWN"}, jobStatus) {
			return nil, "ERROR", fmt.Errorf("unexpected status: %s", jobStatus)
		}
		if utils.StrSliceContains([]string{"SUCCEEDED"}, jobStatus) {
			return respBody, "COMPLETED", nil
		}
		return "continue", "PENDING", nil
	}
}

func waitV3JobCompleted(ctx context.Context, client *golangsdk.ServiceClient, jobId string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      jobStatusRefreshFunc(client, jobId),
		Timeout:      timeout,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the job (in ServiceStage service) to complete: %s", err)
	}
	return nil
}

func resourceV3ComponentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/cas/applications/{application_id}/components"
		appId   = d.Get("application_id").(string)
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{application_id}", appId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildV3ComponentCreateBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating component: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	componentId := utils.PathSearch("component_id", respBody, "").(string)
	if componentId == "" {
		return diag.Errorf("unable to find the component ID from the API response")
	}
	d.SetId(componentId)

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID of the component creating operation from the API response")
	}
	err = waitV3JobCompleted(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	// If the request is successful, obtain the values ​​of all JSON parameters first and save them to the
	// corresponding '_origin' attributes for subsequent determination and construction of the request body during
	// next updates.
	err = utils.RefreshObjectParamOriginValues(d, componentJsonParamKeys)
	if err != nil {
		return diag.Errorf("unable to refresh the origin values: %s", err)
	}

	return resourceV3ComponentRead(ctx, d, meta)
}

func QueryV3Component(client *golangsdk.ServiceClient, applicationId, componentId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/cas/applications/{application_id}/components/{component_id}"

	queryPath := client.Endpoint + httpUrl
	queryPath = strings.ReplaceAll(queryPath, "{project_id}", client.ProjectID)
	queryPath = strings.ReplaceAll(queryPath, "{application_id}", applicationId)
	queryPath = strings.ReplaceAll(queryPath, "{component_id}", componentId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", queryPath, &opt)
	if err != nil {
		return nil, common.ConvertExpected401ErrInto404Err(err, "error_code", v3ComponentNotFoundCodes...)
	}

	return utils.FlattenResponse(requestResp)
}

func flattenV3ComponentRuntimeStackConfig(runtimeStack map[string]interface{}) []map[string]interface{} {
	if len(runtimeStack) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":        utils.PathSearch("name", runtimeStack, nil),
			"type":        utils.PathSearch("type", runtimeStack, nil),
			"version":     utils.PathSearch("version", runtimeStack, nil),
			"deploy_mode": utils.PathSearch("deploy_mode", runtimeStack, nil),
		},
	}
}

func flattenV3ComponentEnvVariables(variables []interface{}) []interface{} {
	if len(variables) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(variables))
	for _, variable := range variables {
		envName := utils.PathSearch("name", variable, "").(string)
		if utils.IsStrContainsSliceElement(envName, []string{"TZ"}, true, true) {
			continue
		}
		result = append(result, map[string]interface{}{
			"name":  envName,
			"value": utils.PathSearch("value", variable, nil),
		})
	}

	return result
}

func flattenV3ComponentStorageMounts(mounts []interface{}) []interface{} {
	if len(mounts) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(mounts))
	for _, mount := range mounts {
		result = append(result, map[string]interface{}{
			"path":      utils.PathSearch("path", mount, nil),
			"sub_path":  utils.PathSearch("sub_path", mount, nil),
			"read_only": utils.PathSearch("read_only", mount, nil),
		})
	}

	return result
}

func flattenV3ComponentStorages(storages []interface{}) []interface{} {
	if len(storages) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(storages))
	for _, storage := range storages {
		result = append(result, map[string]interface{}{
			"type":       utils.PathSearch("type", storage, nil),
			"name":       utils.PathSearch("name", storage, nil),
			"parameters": utils.JsonToString(utils.PathSearch("parameters", storage, nil)),
			"mounts":     flattenV3ComponentStorageMounts(utils.PathSearch("mounts", storage, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenV3ComponentDeployStrategy(d *schema.ResourceData, strategy map[string]interface{}) []map[string]interface{} {
	if len(strategy) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"type":                   utils.PathSearch("type", strategy, nil),
			"rolling_release":        utils.JsonToString(utils.PathSearch("rolling_release", strategy, nil)),
			"gray_release":           utils.JsonToString(utils.PathSearch("gray_release", strategy, nil)),
			"rolling_release_origin": d.Get("deploy_strategy.0.rolling_release_origin"),
			"gray_release_origin":    d.Get("deploy_strategy.0.gray_release_origin"),
		},
	}
}

func flattenV3ComponentLifecycle(lifecycle map[string]interface{}) []map[string]interface{} {
	if len(lifecycle) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"type":    utils.PathSearch("type", lifecycle, nil),
			"scheme":  utils.PathSearch("scheme", lifecycle, nil),
			"host":    utils.PathSearch("host", lifecycle, nil),
			"port":    utils.PathSearch("port", lifecycle, nil),
			"path":    utils.PathSearch("path", lifecycle, nil),
			"command": utils.PathSearch("command", lifecycle, nil),
		},
	}
}

func flattenV3ComponentMesher(mesher map[string]interface{}) []map[string]interface{} {
	if len(mesher) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"port": utils.PathSearch("port", mesher, nil),
		},
	}
}

func flattenV3ComponentLogs(logList []interface{}) []interface{} {
	if len(logList) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(logList))
	for _, val := range logList {
		result = append(result, map[string]interface{}{
			"log_path":         utils.PathSearch("log_path", val, nil),
			"rotate":           utils.PathSearch("rotate", val, nil),
			"host_path":        utils.PathSearch("host_path", val, nil),
			"host_extend_path": utils.PathSearch("host_extend_path", val, nil),
		})
	}
	return result
}

func flattenV3ComponentCustomMetric(customMetric map[string]interface{}) []map[string]interface{} {
	if len(customMetric) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"path":       utils.PathSearch("path", customMetric, nil),
			"port":       utils.PathSearch("port", customMetric, nil),
			"dimensions": utils.PathSearch("dimensions", customMetric, nil),
		},
	}
}

func flattenV3ComponentAffinityMatchExpressions(matchRules []interface{}) []map[string]interface{} {
	if len(matchRules) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(matchRules))
	for _, rule := range matchRules {
		result = append(result, map[string]interface{}{
			"key":       utils.PathSearch("key", rule, nil),
			"operation": utils.PathSearch("operation", rule, nil),
			"value":     utils.PathSearch("value", rule, nil),
		})
	}

	return result
}

func flattenV3ComponentAffinity(affinityRules []interface{}) []map[string]interface{} {
	if len(affinityRules) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(affinityRules))
	for _, rule := range affinityRules {
		result = append(result, map[string]interface{}{
			"condition": utils.PathSearch("condition", rule, nil),
			"kind":      utils.PathSearch("kind", rule, nil),
			"match_expressions": flattenV3ComponentAffinityMatchExpressions(utils.PathSearch("match_expressions",
				rule, make([]interface{}, 0)).([]interface{})),
			"weight": int(utils.PathSearch("weight", rule, float64(0)).(float64)),
		})
	}

	return result
}

func flattenV3ComponentProbe(probe map[string]interface{}) []map[string]interface{} {
	if len(probe) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"type":    utils.PathSearch("type", probe, nil),
			"delay":   utils.PathSearch("delay", probe, nil),
			"timeout": utils.PathSearch("timeout", probe, nil),
			"scheme":  utils.PathSearch("scheme", probe, nil),
			"host":    utils.PathSearch("host", probe, nil),
			"port":    utils.PathSearch("port", probe, nil),
			"path":    utils.PathSearch("path", probe, nil),
			"command": utils.PathSearch("command", probe, make([]interface{}, 0)),
		},
	}
}

func flattenV3ComponentReferResources(refResources []interface{}) []map[string]interface{} {
	if len(refResources) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(refResources))
	for _, refRsource := range refResources {
		result = append(result, map[string]interface{}{
			"id":         utils.PathSearch("id", refRsource, nil),
			"type":       utils.PathSearch("type", refRsource, nil),
			"parameters": utils.JsonToString(utils.PathSearch("parameters", refRsource, nil)),
		})
	}

	return result
}

func flattenV3ExternalAccesses(accesses []interface{}) []map[string]interface{} {
	if len(accesses) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(accesses))
	for _, access := range accesses {
		protocol := utils.PathSearch("protocol", access, "").(string)
		lowercaseProtocol := strings.ToLower(protocol)
		// Only external accesses of protocol http or https can be defined manually.
		if lowercaseProtocol != "http" && lowercaseProtocol != "https" {
			continue
		}
		result = append(result, map[string]interface{}{
			"protocol":     utils.PathSearch("protocol", access, nil),
			"address":      utils.PathSearch("address", access, nil),
			"forward_port": utils.PathSearch("forward_port", access, nil),
		})
	}

	return result
}

func resourceV3ComponentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		appId       = d.Get("application_id").(string)
		componentId = d.Id()
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	respBody, err := QueryV3Component(client, appId, componentId)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error getting component (%s)", componentId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("runtime_stack", flattenV3ComponentRuntimeStackConfig(utils.PathSearch("runtime_stack", respBody,
			make(map[string]interface{})).(map[string]interface{}))),
		d.Set("environment_id", utils.PathSearch("environment_id", respBody, nil)),
		d.Set("config_mode", utils.PathSearch("config_mode", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("source", utils.JsonToString(utils.PathSearch("source", respBody, nil))),
		d.Set("build", utils.JsonToString(utils.PathSearch("build", respBody, nil))),
		d.Set("limit_cpu", utils.PathSearch("limit_cpu", respBody, nil)),
		d.Set("limit_memory", utils.PathSearch("limit_memory", respBody, nil)),
		d.Set("request_cpu", utils.PathSearch("request_cpu", respBody, nil)),
		d.Set("request_memory", utils.PathSearch("request_memory", respBody, nil)),
		d.Set("version", utils.PathSearch("version", respBody, nil)),
		d.Set("envs", flattenV3ComponentEnvVariables(utils.PathSearch("envs[?!inner]", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("replica", utils.PathSearch("replica", respBody, nil)),
		d.Set("storages", flattenV3ComponentStorages(utils.PathSearch("storages", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("deploy_strategy", flattenV3ComponentDeployStrategy(d, utils.PathSearch("deploy_strategy", respBody,
			make(map[string]interface{})).(map[string]interface{}))),
		d.Set("update_strategy", utils.JsonToString(utils.PathSearch("update_strategy", respBody, nil))),
		d.Set("command", utils.JsonToString(utils.PathSearch("command", respBody, nil))),
		d.Set("post_start", flattenV3ComponentLifecycle(utils.PathSearch("post_start", respBody,
			make(map[string]interface{})).(map[string]interface{}))),
		d.Set("pre_stop", flattenV3ComponentLifecycle(utils.PathSearch("pre_stop", respBody,
			make(map[string]interface{})).(map[string]interface{}))),
		d.Set("mesher", flattenV3ComponentMesher(utils.PathSearch("mesher", respBody,
			make(map[string]interface{})).(map[string]interface{}))),
		d.Set("timezone", utils.PathSearch("timezone", respBody, nil)),
		d.Set("jvm_opts", utils.PathSearch("jvm_opts", respBody, nil)),
		d.Set("tomcat_opts", utils.JsonToString(utils.PathSearch("tomcat_opts", respBody, nil))),
		d.Set("logs", flattenV3ComponentLogs(utils.PathSearch("logs", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("custom_metric", flattenV3ComponentCustomMetric(utils.PathSearch("custom_metric", respBody,
			make(map[string]interface{})).(map[string]interface{}))),
		d.Set("affinity", flattenV3ComponentAffinity(utils.PathSearch("affinity", respBody,
			make([]interface{}, 0)).([]interface{}))),
		d.Set("anti_affinity", flattenV3ComponentAffinity(utils.PathSearch("anti_affinity", respBody,
			make([]interface{}, 0)).([]interface{}))),
		d.Set("liveness_probe", flattenV3ComponentProbe(utils.PathSearch("liveness_probe", respBody,
			make(map[string]interface{})).(map[string]interface{}))),
		d.Set("readiness_probe", flattenV3ComponentProbe(utils.PathSearch("readiness_probe", respBody,
			make(map[string]interface{})).(map[string]interface{}))),
		d.Set("refer_resources", flattenV3ComponentReferResources(utils.PathSearch("refer_resources", respBody,
			make([]interface{}, 0)).([]interface{}))),
		d.Set("external_accesses", flattenV3ExternalAccesses(utils.PathSearch("external_accesses", respBody,
			make([]interface{}, 0)).([]interface{}))),
		d.Set("status", utils.PathSearch("status.component_status", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("status.create_time", respBody,
			float64(0)).(float64))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("status.update_time", respBody,
			float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildV3ComponentUpdteBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Cannot be updated but the request body needs them.
		"name":          d.Get("name").(string),
		"runtime_stack": utils.ValueIgnoreEmpty(buildV3ComponentRuntimeStackConfig(d.Get("runtime_stack").([]interface{}))),
		// Required parameters
		"source":          utils.StringToJson(d.Get("source").(string)),
		"version":         d.Get("version").(string),
		"refer_resources": utils.ValueIgnoreEmpty(buildV3ComponentReferResources(d.Get("refer_resources").(*schema.Set))),
		// Optional parameters.
		"description":       d.Get("description").(string),
		"build":             utils.StringToJson(d.Get("build").(string)),
		"replica":           d.Get("replica").(int),
		"limit_cpu":         d.Get("limit_cpu").(float64),
		"limit_memory":      d.Get("limit_memory").(float64),
		"request_cpu":       d.Get("request_cpu").(float64),
		"request_memory":    d.Get("request_memory").(float64),
		"envs":              utils.ValueIgnoreEmpty(buildV3ComponentEnvVariables(d.Get("envs").(*schema.Set))),
		"storages":          utils.ValueIgnoreEmpty(buildV3ComponentStorages(d.Get("storages").(*schema.Set))),
		"deploy_strategy":   utils.ValueIgnoreEmpty(buildV3ComponentDeployStrategy(d.Get("deploy_strategy").([]interface{}))),
		"update_strategy":   utils.StringToJson(d.Get("update_strategy").(string)),
		"command":           utils.StringToJson(d.Get("command").(string)),
		"post_start":        utils.ValueIgnoreEmpty(buildV3ComponentLifecycle(d.Get("post_start").([]interface{}))),
		"pre_stop":          utils.ValueIgnoreEmpty(buildV3ComponentLifecycle(d.Get("pre_stop").([]interface{}))),
		"mesher":            utils.ValueIgnoreEmpty(buildV3ComponentMesher(d.Get("mesher").([]interface{}))),
		"timezone":          utils.ValueIgnoreEmpty(d.Get("timezone").(string)),
		"jvm_opts":          utils.ValueIgnoreEmpty(d.Get("jvm_opts").(string)),
		"tomcat_opts":       utils.StringToJson(d.Get("tomcat_opts").(string)),
		"logs":              utils.ValueIgnoreEmpty(buildV3ComponentLogs(d.Get("logs").(*schema.Set))),
		"custom_metric":     utils.ValueIgnoreEmpty(buildV3ComponentCustomMetric(d.Get("custom_metric").([]interface{}))),
		"affinity":          utils.ValueIgnoreEmpty(buildV3ComponentAffinity(d.Get("affinity").(*schema.Set))),
		"anti_affinity":     utils.ValueIgnoreEmpty(buildV3ComponentAffinity(d.Get("anti_affinity").(*schema.Set))),
		"liveness_probe":    utils.ValueIgnoreEmpty(buildV3ComponentProbeConfiguration(d.Get("liveness_probe").([]interface{}))),
		"readiness_probe":   utils.ValueIgnoreEmpty(buildV3ComponentProbeConfiguration(d.Get("readiness_probe").([]interface{}))),
		"external_accesses": utils.ValueIgnoreEmpty(buildV3ComponentExternalAccesses(d.Get("external_accesses").(*schema.Set))),
		"labels":            utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
	}
}

func componentStatusRefreshFunc(client *golangsdk.ServiceClient, appId, commponetId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := QueryV3Component(client, appId, commponetId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				log.Printf("[DEBUG] The component (%s) does not exist", commponetId)
				return "Resource Not Found", "COMPLETED", nil
			}
			return respBody, "ERROR", err
		}

		componentStatus := utils.PathSearch("status.component_status", respBody, "").(string)
		if utils.IsStrContainsSliceElement(componentStatus, []string{"FAILED", "UNKNOWN", "PARTIALLY_FAILED"}, true, true) {
			return nil, "ERROR", fmt.Errorf("unexpected status: %s", componentStatus)
		}
		if utils.IsStrContainsSliceElement(componentStatus, targets, true, true) {
			return respBody, "COMPLETED", nil
		}
		return "continue", "PENDING", nil
	}
}

func waitV3ComponentUpdateCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	expectedStatuses := []string{
		"PENDING",
		"RUNNING",
		"GRAYING", // Upgrade the component by gray release and waiting for manual continuation.
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      componentStatusRefreshFunc(client, d.Get("application_id").(string), d.Id(), expectedStatuses),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the update operation to complete: %s", err)
	}
	return nil
}

func resourceV3ComponentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v3/{project_id}/cas/applications/{application_id}/components/{component_id}"
		appId       = d.Get("application_id").(string)
		componentId = d.Id()
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{application_id}", appId)
	createPath = strings.ReplaceAll(createPath, "{component_id}", componentId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildV3ComponentUpdteBodyParams(d)),
	}

	_, err = client.Request("PUT", createPath, &opt)
	if err != nil {
		return diag.Errorf("error updating component (%s): %s", componentId, err)
	}

	err = waitV3ComponentUpdateCompleted(ctx, client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	// If the request is successful, obtain the values ​​of all JSON parameters first and save them to the
	// corresponding '_origin' attributes for subsequent determination and construction of the request body during
	// next updates.
	err = utils.RefreshObjectParamOriginValues(d, componentJsonParamKeys)
	if err != nil {
		return diag.Errorf("unable to refresh the origin values: %s", err)
	}

	return resourceV3ComponentRead(ctx, d, meta)
}

func waitV3ComponentDeleteCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      componentStatusRefreshFunc(client, d.Get("application_id").(string), d.Id(), nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the delete operation to complete: %s", err)
	}
	return nil
}

func resourceV3ComponentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v3/{project_id}/cas/applications/{application_id}/components/{component_id}"
		appId       = d.Get("application_id").(string)
		componentId = d.Id()
	)

	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{application_id}", appId)
	deletePath = strings.ReplaceAll(deletePath, "{component_id}", componentId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// Returns the state code 200 and structure (with format '{"component_id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"}').
	_, err = client.Request("DELETE", deletePath, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected401ErrInto404Err(err, "error_code", v3ComponentNotFoundCodes...),
			fmt.Sprintf("error deleting component (%s)", componentId))
	}

	return diag.FromErr(waitV3ComponentDeleteCompleted(ctx, client, d))
}

func resourceV3ComponentImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<application_id>/<id>', but got '%s'", importedId)
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("application_id", parts[0])
}
