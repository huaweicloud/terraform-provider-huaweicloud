package modelarts

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

var v2ServiceNonUpdatableParams = []string{
	"name",
	"type",
	"workspace_id",
	"deploy_type",
}

// @API ModelArts POST /v2/{project_id}/services
// @API ModelArts GET /v2/{project_id}/services/{service_id}
// @API ModelArts GET /v2/{project_id}/services/{service_id}/versions
// @API ModelArts POST /v2/{project_id}/services/{service_id}/versions/switch
// @API ModelArts PUT /v2/{project_id}/services/{service_id}
// @API ModelArts POST /v2/{project_id}/modelarts-service-v2/{service_id}/tags/create
// @API ModelArts DELETE /v2/{project_id}/modelarts-service-v2/{service_id}/tags/create
// @API ModelArts POST /v2/{project_id}/services/delete
func ResourceV2Service() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2ServiceCreate,
		ReadContext:   resourceV2ServiceRead,
		UpdateContext: resourceV2ServiceUpdate,
		DeleteContext: resourceV2ServiceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(v2ServiceNonUpdatableParams),
			config.MergeDefaultTags(),
		),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the service is located.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the service.`,
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The version of the service.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The reasoning method of the service.`,
			},
			"group_configs": {
				// The order of this parameter is determined by the user and the service will not reorder this result.
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Required parameters.
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
							Description:      `The name of the instance group.`,
						},
						"count": {
							Type:             schema.TypeInt,
							Required:         true,
							DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
							Description:      `The number of service instances in the deployment scenario.`,
						},
						"weight": {
							Type:             schema.TypeInt,
							Required:         true,
							DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
							Description:      `The weight percentage of the instance group.`,
						},
						"unit_configs": {
							Type:             schema.TypeList,
							Required:         true,
							Elem:             schemaServiceGroupConfigUnitConfig(),
							DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
							Description:      `The unit configurations of the instance group.`,
						},
						// Optional parameters.
						"pool_id": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
							Description:      `The ID of the dedicated resource pool for the instance group.`,
						},
						"framework": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
							Description:      `The algorithm framework.`,
						},
						// Attributes
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the instance group.`,
						},
					},
				},
				Description: `The instance group configurations of the service.`,
			},
			"runtime_config": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The configuration of the service runtime, in JSON format.`,
			},
			"upgrade_config": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The upgrade configuration of the service, in JSON format.`,
			},
			// Optional parameters.
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the service.`,
			},
			"workspace_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The workspace ID of the service.`,
			},
			"deploy_type": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The deploy type of the service.`,
			},
			"log_configs": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Required parameters.
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
							Description:      `The type of LTS configuration.`,
						},
						// Optional parameters.
						"log_group_id": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
							Description:      `The ID of the LTS group.`,
						},
						"log_stream_id": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
							Description:      `The ID of the LTS stream.`,
						},
					},
				},
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The log configurations of the service.`,
			},
			"tags": common.TagsSchema(`The key/value pairs to associate with the service.`),
			// Attributes.
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the service.`,
			},
			"predict_url": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of service access.`,
						},
						"urls": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The URLs of service access.`,
						},
					},
				},
				Description: `The access addresses of the service.`,
			},
			// Internal attributes.
			// Store all version information of the current service deployment, used to determine whether the version of
			// the current script configuration belongs to the historical version when the script is changed.
			// + Yes: execute the version switch, ignore the changes of other parameters, and only query the historical
			//        version configuration.
			// + No: deploy a new version and apply the current configuration to the new version.
			"history_versions": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: utils.SchemaDesc(
					`The deployed history information of the service versions.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func schemaServiceGroupConfigUnitConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required parameters.
			"image": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Required parameters.
						"source": {
							Type:     schema.TypeString,
							Required: true,
							DiffSuppressFunc: utils.ComposeAnySchemaDiffSuppressFunc(
								suppressV2ServiceUpgradeParamDiffs(),
								utils.SuppressCaseDiffs(),
							),
							Description: `The image type of the group unit.`,
						},
						"swr_path": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
							Description:      `The SWR storage path of the group unit.`,
						},
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The image ID of the group unit.`,
						},
					},
				},
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The image configuration of the group unit.`,
			},
			// Optional parameters.
			"role": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: utils.ComposeAnySchemaDiffSuppressFunc(
					suppressV2ServiceUpgradeParamDiffs(),
					utils.SuppressCaseDiffs(),
				),
				Description: `The role of the group unit.`,
			},
			"custom_spec": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Optional parameters.
						"gpu": {
							Type:     schema.TypeFloat,
							Optional: true,
							Computed: true,
							DiffSuppressFunc: utils.ComposeAnySchemaDiffSuppressFunc(
								suppressV2ServiceUpgradeParamDiffs(),
								suppressV2ServiceFloatValuesDiffs(),
							),
							Description: `The GPU number of the custom specification.`,
						},
						"memory": {
							Type:             schema.TypeInt,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
							Description:      `The memory size of the custom specification.`,
						},
						"cpu": {
							Type:     schema.TypeFloat,
							Optional: true,
							Computed: true,
							DiffSuppressFunc: utils.ComposeAnySchemaDiffSuppressFunc(
								suppressV2ServiceUpgradeParamDiffs(),
								suppressV2ServiceFloatValuesDiffs(),
							),
							Description: `The CPU number of the custom specification.`,
						},
						"ascend": {
							Type:             schema.TypeInt,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
							Description:      `The number of Ascend chips.`,
						},
					},
				},
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The configuration of the custom resource specification.`,
			},
			"flavor": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The instance flavor of the group unit.`,
			},
			"models": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
							Description:      `The source type of the model configuration.`,
						},
						"mount_path": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
							Description:      `The path to mount into the container.`,
						},
						// Optional parameters.
						"address": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
							Description:      `The source address of the model configuration.`,
						},
						"source_id": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
							Description:      `The source ID of the model configuration.`,
						},
					},
				},
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The model configuration of the group unit.`,
			},
			"codes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
							Description:      `The source type of the code configuration.`,
						},
						"mount_path": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
							Description:      `The path to mount into the container.`,
						},
						// Optional parameters.
						"address": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
							Description:      `The source address of the code configuration.`,
						},
						"source_id": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
							Description:      `The source ID of the code configuration.`,
						},
					},
				},
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The code configuration of the group unit.`,
			},
			"count": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The instance number of the group unit.`,
			},
			"cmd": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The startup commands of the group unit.`,
			},
			"envs": {
				Type:             schema.TypeMap,
				Optional:         true,
				Computed:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The environment variables of the group unit.`,
			},
			"readiness_health": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				MaxItems:         1,
				Elem:             schemaServiceUnitConfigHealth(),
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The configuration of the readiness health check.`,
			},
			"startup_health": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				MaxItems:         1,
				Elem:             schemaServiceUnitConfigHealth(),
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The configuration of the startup health check.`,
			},
			"liveness_health": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				MaxItems:         1,
				Elem:             schemaServiceUnitConfigHealth(),
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The configuration of the liveness health check.`,
			},
			"port": {
				Type:             schema.TypeInt,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The port of the group unit.`,
			},
			"recovery": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The recovery strategy of the group unit.`,
			},
			// Attributes.
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the unit configuration.`,
			},
		},
	}
}

func schemaServiceUnitConfigHealth() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required parameters.
			"initial_delay_seconds": {
				Type:             schema.TypeInt,
				Required:         true,
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The time to wait when performing the first probe.`,
			},
			"timeout_seconds": {
				Type:             schema.TypeInt,
				Required:         true,
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The timeout for executing the probe.`,
			},
			"period_seconds": {
				Type:             schema.TypeInt,
				Required:         true,
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The period time for performing health check.`,
			},
			"failure_threshold": {
				Type:             schema.TypeInt,
				Required:         true,
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The minimum number of consecutive detection failures.`,
			},
			"check_method": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The method of the health check.`,
			},
			// Optional parameters.
			"command": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The commands configuration of the health check.`,
			},
			"url": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressV2ServiceUpgradeParamDiffs(),
				Description:      `The request URL of the health check.`,
			},
		},
	}
}

func isHistoryServiceVersion(historyVersions map[string]interface{}, targetVersion string) bool {
	_, isExist := historyVersions[targetVersion]
	return isExist
}

func suppressV2ServiceUpgradeParamDiffs() schema.SchemaDiffSuppressFunc {
	return func(_, _, _ string, d *schema.ResourceData) bool {
		historyVersions := d.Get("history_versions").(map[string]interface{})
		targetVersion := d.Get("version").(string)
		return isHistoryServiceVersion(historyVersions, targetVersion)
	}
}

// Only available for the float type of the parameter values compare operation.
func suppressV2ServiceFloatValuesDiffs() schema.SchemaDiffSuppressFunc {
	return func(k, _, _ string, d *schema.ResourceData) bool {
		oldVal, newVal := d.GetChange(k)
		oldF, oldValIsTypeFloat := oldVal.(float64)
		newF, newValIsTypeFloat := newVal.(float64)
		if !oldValIsTypeFloat || !newValIsTypeFloat {
			return false
		}

		return utils.EqualFloat(utils.Round(oldF, 2), utils.Round(newF, 2))
	}
}

func buildV2ServiceUnitConfigCustomSpec(specConfigs []interface{}) map[string]interface{} {
	// All sub parameter behaviors are optional.
	if len(specConfigs) < 1 || specConfigs[0] == nil {
		return nil
	}

	specConfig := specConfigs[0]
	return map[string]interface{}{
		"gpu":    utils.ValueIgnoreEmpty(utils.PathSearch("gpu", specConfig, nil)),
		"memory": utils.ValueIgnoreEmpty(utils.PathSearch("memory", specConfig, nil)),
		"cpu":    utils.ValueIgnoreEmpty(utils.PathSearch("cpu", specConfig, nil)),
		"ascend": utils.ValueIgnoreEmpty(utils.PathSearch("ascend", specConfig, nil)),
	}
}

func buildV2ServiceUnitConfigImage(imageConfigs []interface{}) map[string]interface{} {
	if len(imageConfigs) < 1 {
		return nil
	}

	imageConfig := imageConfigs[0]
	return map[string]interface{}{
		"source":   utils.PathSearch("source", imageConfig, "").(string),
		"swr_path": utils.PathSearch("swr_path", imageConfig, nil),
		// Only 'IMAGE' source type of image requires this parameter input.
		"id": utils.ValueIgnoreEmpty(utils.PathSearch("source == 'IMAGE' && id || ''", imageConfig, nil)),
	}
}

func buildV2ServiceUnitConfigModels(modelConfigs []interface{}) []map[string]interface{} {
	if len(modelConfigs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(modelConfigs))
	for _, modelConfig := range modelConfigs {
		result = append(result, map[string]interface{}{
			"source":     utils.PathSearch("source", modelConfig, nil),
			"mount_path": utils.PathSearch("mount_path", modelConfig, nil),
			"address":    utils.ValueIgnoreEmpty(utils.PathSearch("address", modelConfig, nil)),
			"source_id":  utils.ValueIgnoreEmpty(utils.PathSearch("source_id", modelConfig, nil)),
		})
	}
	return result
}

func buildV2ServiceUnitConfigCodes(codeConfigs []interface{}) []map[string]interface{} {
	if len(codeConfigs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(codeConfigs))
	for _, codeConfig := range codeConfigs {
		result = append(result, map[string]interface{}{
			"source":     utils.PathSearch("source", codeConfig, nil),
			"mount_path": utils.PathSearch("mount_path", codeConfig, nil),
			"address":    utils.ValueIgnoreEmpty(utils.PathSearch("address", codeConfig, nil)),
			"source_id":  utils.ValueIgnoreEmpty(utils.PathSearch("source_id", codeConfig, nil)),
		})
	}
	return result
}

func buildV2ServiceUnitConfigHealthImage(imageConfigs []interface{}) map[string]interface{} {
	if len(imageConfigs) < 1 {
		return nil
	}

	imageConfig := imageConfigs[0]
	return map[string]interface{}{
		// Required parameters.
		"source":   utils.PathSearch("source", imageConfig, nil),
		"swr_path": utils.PathSearch("swr_path", imageConfig, nil),
	}
}

func buildV2ServiceUnitConfigHealth(healthConfigs []interface{}) map[string]interface{} {
	if len(healthConfigs) < 1 {
		return nil
	}

	healthConfig := healthConfigs[0]
	return map[string]interface{}{
		// Optional parameters.
		"initial_delay_seconds": utils.PathSearch("initial_delay_seconds", healthConfig, nil),
		"timeout_seconds":       utils.PathSearch("timeout_seconds", healthConfig, nil),
		"period_seconds":        utils.PathSearch("period_seconds", healthConfig, nil),
		"failure_threshold":     utils.PathSearch("failure_threshold", healthConfig, nil),
		// Optional parameters.
		"check_method": utils.ValueIgnoreEmpty(utils.PathSearch("check_method", healthConfig, nil)),
		"command":      utils.ValueIgnoreEmpty(utils.PathSearch("command", healthConfig, nil)),
		"url":          utils.ValueIgnoreEmpty(utils.PathSearch("url", healthConfig, nil)),
		"image": utils.ValueIgnoreEmpty(buildV2ServiceUnitConfigHealthImage(utils.PathSearch("image",
			healthConfig, make([]interface{}, 0)).([]interface{}))),
	}
}

func buildV2ServiceUnitConfigs(unitConfigs []interface{}) []map[string]interface{} {
	if len(unitConfigs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(unitConfigs))
	for _, groupConfig := range unitConfigs {
		result = append(result, map[string]interface{}{
			// Required parameters.
			"image": buildV2ServiceUnitConfigImage(utils.PathSearch("image", groupConfig, make([]interface{}, 0)).([]interface{})),
			// Optional parameters.
			"role": utils.ValueIgnoreEmpty(utils.PathSearch("role", groupConfig, nil)),
			"custom_spec": utils.ValueIgnoreEmpty(buildV2ServiceUnitConfigCustomSpec(utils.PathSearch("custom_spec",
				groupConfig, make([]interface{}, 0)).([]interface{}))),
			"flavor": utils.ValueIgnoreEmpty(utils.PathSearch("flavor", groupConfig, nil)),
			"models": utils.ValueIgnoreEmpty(buildV2ServiceUnitConfigModels(utils.PathSearch("models",
				groupConfig, make([]interface{}, 0)).([]interface{}))),
			"codes": utils.ValueIgnoreEmpty(buildV2ServiceUnitConfigCodes(utils.PathSearch("codes",
				groupConfig, make([]interface{}, 0)).([]interface{}))),
			"count": utils.ValueIgnoreEmpty(utils.PathSearch("count", groupConfig, nil)),
			"cmd":   utils.ValueIgnoreEmpty(utils.PathSearch("cmd", groupConfig, nil)),
			"envs":  utils.ValueIgnoreEmpty(utils.PathSearch("envs", groupConfig, make(map[string]interface{})).(map[string]interface{})),
			"readiness_health": utils.ValueIgnoreEmpty(buildV2ServiceUnitConfigHealth(utils.PathSearch("readiness_health",
				groupConfig, make([]interface{}, 0)).([]interface{}))),
			"startup_health": utils.ValueIgnoreEmpty(buildV2ServiceUnitConfigHealth(utils.PathSearch("startup_health",
				groupConfig, make([]interface{}, 0)).([]interface{}))),
			"liveness_health": utils.ValueIgnoreEmpty(buildV2ServiceUnitConfigHealth(utils.PathSearch("liveness_health",
				groupConfig, make([]interface{}, 0)).([]interface{}))),
			"port":     utils.ValueIgnoreEmpty(utils.PathSearch("port", groupConfig, nil)),
			"recovery": utils.ValueIgnoreEmpty(utils.PathSearch("recovery", groupConfig, nil)),
		})
	}

	return result
}

func buildV2ServiceGroupConfigs(groupConfigs []interface{}) []map[string]interface{} {
	if len(groupConfigs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(groupConfigs))
	for _, groupConfig := range groupConfigs {
		result = append(result, map[string]interface{}{
			// Required parameters.
			"name":  utils.PathSearch("name", groupConfig, nil),
			"count": utils.PathSearch("count", groupConfig, nil),
			// Only required if the value of parameter type is 'REAL_TIME'.
			"weight": utils.ValueIgnoreEmpty(utils.PathSearch("weight", groupConfig, nil)),
			// Optional parameters.
			"pool_id":   utils.ValueIgnoreEmpty(utils.PathSearch("pool_id", groupConfig, nil)),
			"framework": utils.ValueIgnoreEmpty(utils.PathSearch("framework", groupConfig, nil)),
			"unit_configs": utils.ValueIgnoreEmpty(buildV2ServiceUnitConfigs(utils.PathSearch("unit_configs",
				groupConfig, make([]interface{}, 0)).([]interface{}))),
		})
	}

	return result
}

func buildV2ServiceLogConfigs(logConfigs []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(logConfigs))

	for _, logConfig := range logConfigs {
		result = append(result, map[string]interface{}{
			// Required parameters.
			"weight": utils.PathSearch("weight", logConfig, nil),
			// Optional parameters.
			"log_group_id":  utils.PathSearch("log_group_id", logConfig, nil),
			"log_stream_id": utils.PathSearch("log_stream_id", logConfig, nil),
		})
	}

	return result
}

func buildV2ServiceCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Required parameters.
		"name":           d.Get("name"),
		"version":        d.Get("version"),
		"type":           d.Get("type"),
		"group_configs":  buildV2ServiceGroupConfigs(d.Get("group_configs").([]interface{})),
		"runtime_config": utils.StringToJson(d.Get("runtime_config").(string)),
		"upgrade_config": utils.StringToJson(d.Get("upgrade_config").(string)),
		// Optional parameters.
		"description":  utils.ValueIgnoreEmpty(d.Get("description")),
		"workspace_id": utils.ValueIgnoreEmpty(d.Get("workspace_id")),
		"deploy_type":  utils.ValueIgnoreEmpty(d.Get("deploy_type")),
		"log_configs":  utils.ValueIgnoreEmpty(buildV2ServiceLogConfigs(d.Get("log_configs").([]interface{}))),
		"tags":         utils.ValueIgnoreEmpty(utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{}), true)),
	}
}

func serviceStatusRefreshFunc(client *golangsdk.ServiceClient, serviceId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetServiceById(client, serviceId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return "Resource Not Found", "COMPLETED", nil
			}
			return respBody, "ERROR", err
		}

		statusResp := utils.PathSearch("status", respBody, "").(string)
		if utils.StrSliceContains([]string{"FAILED", "EERROR"}, statusResp) {
			return respBody, "ERROR", fmt.Errorf("unexpect status (%s)", statusResp)
		}

		if utils.StrSliceContains(targets, statusResp) {
			return respBody, "COMPLETED", nil
		}

		return "continue", "PENDING", nil
	}
}

func resourceV2ServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/services"
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildV2ServiceCreateBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating service: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	serviceId := utils.PathSearch("id", respBody, "").(string)
	if serviceId == "" {
		return diag.Errorf("unable to find the service ID from the API response")
	}
	d.SetId(serviceId)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      serviceStatusRefreshFunc(client, serviceId, []string{"STOPPED", "RUNNING"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the service create operation to become complete: %s", err)
	}

	return resourceV2ServiceRead(ctx, d, meta)
}

func GetServiceById(client *golangsdk.ServiceClient, serviceId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/services/{service_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{service_id}", serviceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func flattenServiceUnitConfigImage(imageConfig interface{}) []map[string]interface{} {
	if imageConfig == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"source":   utils.PathSearch("source", imageConfig, nil),
			"swr_path": utils.PathSearch("swr_path", imageConfig, nil),
			"id":       utils.PathSearch("id", imageConfig, nil),
		},
	}
}

func flattenServiceUnitConfigCustomSpec(customSpec interface{}) []map[string]interface{} {
	if customSpec == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"gpu":    utils.PathSearch("gpu", customSpec, nil),
			"memory": utils.PathSearch("memory", customSpec, nil),
			"cpu":    utils.PathSearch("cpu", customSpec, nil),
			"ascend": utils.PathSearch("ascend", customSpec, nil),
		},
	}
}

func flattenServiceUnitConfigModels(modelConfigs []interface{}) []map[string]interface{} {
	if len(modelConfigs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(modelConfigs))
	for _, modelConfig := range modelConfigs {
		result = append(result, map[string]interface{}{
			"source":     utils.PathSearch("source", modelConfig, nil),
			"mount_path": utils.PathSearch("mount_path", modelConfig, nil),
			"address":    utils.PathSearch("address", modelConfig, nil),
			"source_id":  utils.PathSearch("source_id", modelConfig, nil),
		})
	}

	return result
}

func flattenServiceUnitConfigCodes(codeConfigs []interface{}) []map[string]interface{} {
	if len(codeConfigs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(codeConfigs))
	for _, codeConfig := range codeConfigs {
		result = append(result, map[string]interface{}{
			"source":     utils.PathSearch("source", codeConfig, nil),
			"mount_path": utils.PathSearch("mount_path", codeConfig, nil),
			"address":    utils.PathSearch("address", codeConfig, nil),
			"source_id":  utils.PathSearch("source_id", codeConfig, nil),
		})
	}

	return result
}

func flattenServiceGroupConfigHealthCheck(healthCheck interface{}) []map[string]interface{} {
	if healthCheck == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"initial_delay_seconds": utils.PathSearch("initial_delay_seconds", healthCheck, nil),
			"timeout_seconds":       utils.PathSearch("timeout_seconds", healthCheck, nil),
			"period_seconds":        utils.PathSearch("period_seconds", healthCheck, nil),
			"failure_threshold":     utils.PathSearch("failure_threshold", healthCheck, nil),
			"check_method":          utils.PathSearch("check_method", healthCheck, nil),
			"command":               utils.PathSearch("command", healthCheck, nil),
			"url":                   utils.PathSearch("url", healthCheck, nil),
		},
	}
}

func flattenServiceGroupConfigUnitConfigs(unitConfigs []interface{}) []map[string]interface{} {
	if len(unitConfigs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(unitConfigs))
	for _, unitConfig := range unitConfigs {
		result = append(result, map[string]interface{}{
			"image":       flattenServiceUnitConfigImage(utils.PathSearch("image", unitConfig, nil)),
			"role":        utils.PathSearch("role", unitConfig, nil),
			"custom_spec": flattenServiceUnitConfigCustomSpec(utils.PathSearch("custom_spec", unitConfig, nil)),
			"flavor":      utils.PathSearch("flavor", unitConfig, nil),
			"models": flattenServiceUnitConfigModels(utils.PathSearch("models",
				unitConfig, make([]interface{}, 0)).([]interface{})),
			"codes": flattenServiceUnitConfigCodes(utils.PathSearch("codes",
				unitConfig, make([]interface{}, 0)).([]interface{})),
			"count": utils.PathSearch("count", unitConfig, nil),
			"cmd":   utils.PathSearch("cmd", unitConfig, nil),
			"envs":  utils.PathSearch("envs", unitConfig, nil),
			"readiness_health": flattenServiceGroupConfigHealthCheck(utils.PathSearch("readiness_health",
				unitConfig, nil)),
			"startup_health": flattenServiceGroupConfigHealthCheck(utils.PathSearch("startup_health",
				unitConfig, nil)),
			"liveness_health": flattenServiceGroupConfigHealthCheck(utils.PathSearch("liveness_health",
				unitConfig, nil)),
			"port":     utils.PathSearch("port", unitConfig, nil),
			"recovery": utils.PathSearch("recovery", unitConfig, nil),
			"id":       utils.PathSearch("id", unitConfig, nil),
		})
	}

	return result
}

func flattenServiceGroupConfigs(groupConfigs []interface{}) []map[string]interface{} {
	if len(groupConfigs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(groupConfigs))
	for _, groupConfig := range groupConfigs {
		result = append(result, map[string]interface{}{
			"name":      utils.PathSearch("name", groupConfig, nil),
			"weight":    utils.PathSearch("weight", groupConfig, nil),
			"pool_id":   utils.PathSearch("pool_id", groupConfig, nil),
			"framework": utils.PathSearch("framework", groupConfig, nil),
			"count":     utils.PathSearch("count", groupConfig, nil),
			"unit_configs": flattenServiceGroupConfigUnitConfigs(utils.PathSearch("unit_configs",
				groupConfig, make([]interface{}, 0)).([]interface{})),
			"id": utils.PathSearch("id", groupConfig, nil),
		})
	}

	return result
}

func flattenServiceLogConfigs(logConfigs []interface{}) []map[string]interface{} {
	if len(logConfigs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(logConfigs))

	for _, logConfig := range logConfigs {
		result = append(result, map[string]interface{}{
			"weight":        utils.PathSearch("weight", logConfig, nil),
			"log_group_id":  utils.PathSearch("log_group_id", logConfig, nil),
			"log_stream_id": utils.PathSearch("log_stream_id", logConfig, nil),
		})
	}

	return result
}

func listV2ServiceHistoryVersions(client *golangsdk.ServiceClient, serviceId, workspaceId string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/services/{service_id}/versions?limit={limit}"
		limit   = 500
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{service_id}", serviceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"workspace_id": workspaceId,
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			// If the service does not exist, a 404 error will be returned.
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		historyVersions := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		// Even if the empty list insert, the same result as the original list will be returned.
		result = append(result, historyVersions...)
		if len(historyVersions) < limit {
			break
		}
		offset += len(historyVersions)
	}

	return result, nil
}

func flattenV2ServiceHistoryVersions(historyVersions []interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for _, historyVersion := range historyVersions {
		versionNum := utils.PathSearch("version", historyVersion, "").(string)
		if versionNum == "" {
			log.Printf("[ERROR] The version number is empty in the version record: %#v", historyVersion)
			continue
		}
		versionId := utils.PathSearch("id", historyVersion, "").(string)
		if versionId == "" {
			log.Printf("[ERROR] The version ID is empty in the version record: %#v", historyVersion)
			continue
		}
		result[versionNum] = versionId
	}
	return result
}

func flattenPredictUrls(predictUrls []interface{}) []map[string]interface{} {
	if len(predictUrls) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(predictUrls))
	for _, predictUrl := range predictUrls {
		result = append(result, map[string]interface{}{
			"type": utils.PathSearch("type", predictUrl, nil),
			"urls": utils.PathSearch("urls", predictUrl, make([]interface{}, 0)),
		})
	}

	return result
}

func resourceV2ServiceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		serviceId = d.Id()
	)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	serviceInfo, err := GetServiceById(client, serviceId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error querying service (%s)", serviceId))
	}

	workspaceId := utils.PathSearch("workspace_id", serviceInfo, "").(string)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		// Requird parameters.
		d.Set("name", utils.PathSearch("name", serviceInfo, nil)),
		d.Set("version", utils.PathSearch("version.version", serviceInfo, nil)),
		d.Set("type", utils.PathSearch("type", serviceInfo, nil)),
		d.Set("group_configs", flattenServiceGroupConfigs(utils.PathSearch("version.instance_groups",
			serviceInfo, make([]interface{}, 0)).([]interface{}))),
		d.Set("runtime_config", utils.JsonToString(utils.PathSearch("version.runtime_config", serviceInfo, nil))),
		d.Set("upgrade_config", utils.JsonToString(utils.PathSearch("version.upgrade_config", serviceInfo, nil))),
		// Optional parameters.
		d.Set("description", utils.PathSearch("version.description", serviceInfo, nil)),
		d.Set("workspace_id", workspaceId),
		d.Set("deploy_type", utils.PathSearch("deploy_type", serviceInfo, nil)),
		d.Set("log_configs", flattenServiceLogConfigs(utils.PathSearch("version.log_configs",
			serviceInfo, make([]interface{}, 0)).([]interface{}))),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", serviceInfo, make([]interface{}, 0)))),
		// Attributes.
		d.Set("status", utils.PathSearch("status", serviceInfo, nil)),
		d.Set("predict_url", flattenPredictUrls(utils.PathSearch("predict_url", serviceInfo, make([]interface{}, 0)).([]interface{}))),
	)

	historyVersions, err := listV2ServiceHistoryVersions(client, serviceId, workspaceId)
	if err != nil {
		log.Printf("[ERROR] Unable to obtain the history version records for the specified version (%s): %s",
			serviceId, err)
	} else {
		mErr = multierror.Append(mErr, d.Set("history_versions", flattenV2ServiceHistoryVersions(historyVersions)))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving resource fields of the ModelArts service (%s): %s", serviceId, err)
	}
	return nil
}

func modifyV2ServiceTags(client *golangsdk.ServiceClient, method, action string, actionTags map[string]interface{}) error {
	modifyPath := strings.ReplaceAll(client.ResourceBase, "{action}", action)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"tags": utils.ExpandResourceTagsMap(actionTags),
		},
	}

	_, err := client.Request(method, modifyPath, &opt)
	return err
}

func updateV2ServiceTags(client *golangsdk.ServiceClient, d *schema.ResourceData) (err error) {
	var (
		httpUrl   = "v2/{project_id}/modelarts-service-v2/{service_id}/tags/{action}"
		serviceId = d.Id()
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{service_id}", serviceId)
	client.ResourceBase = updatePath

	oldTagsRaw, newTagsRaw := d.GetChange("tags")
	deleteTags := utils.TakeObjectsDifferent(oldTagsRaw.(map[string]interface{}), newTagsRaw.(map[string]interface{}))
	createTags := utils.TakeObjectsDifferent(newTagsRaw.(map[string]interface{}), oldTagsRaw.(map[string]interface{}))

	if len(deleteTags) > 0 {
		err = modifyV2ServiceTags(client, "DELETE", "delete", deleteTags)
		if err != nil {
			return
		}
	}

	if len(createTags) > 0 {
		err = modifyV2ServiceTags(client, "POST", "create", createTags)
		if err != nil {
			return
		}
	}

	return
}

func switchServiceVersion(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, targetVersion,
	targetVersionId string) (interface{}, error) {
	var (
		httpUrl     = "v2/{project_id}/services/{service_id}/versions/switch"
		serviceId   = d.Id()
		workspaceId = d.Get("workspace_id").(string)
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{service_id}", serviceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"id":                serviceId,
			"target_version_id": targetVersionId,
			"workspace_id":      workspaceId,
		},
	}

	_, err := client.Request("POST", updatePath, &opt)
	if err != nil {
		return nil, fmt.Errorf("error switching specified service version (%s): %s", targetVersion, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      serviceStatusRefreshFunc(client, serviceId, []string{"STOPPED", "RUNNING"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		PollInterval: 15 * time.Second,
	}
	historyServiceInfo, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error waiting for the service create operation to become complete: %s", err)
	}

	return historyServiceInfo, nil
}

func updateServiceDescription(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl   = "v2/{project_id}/services/{service_id}"
		serviceId = d.Id()
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{service_id}", serviceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"version":     d.Get("version"),
			"description": d.Get("description"),
			"id":          serviceId,
		},
	}

	_, err := client.Request("PUT", updatePath, &opt)
	if err != nil {
		return err
	}
	return nil
}

func buildV2ServiceUpgradeUnitConfigs(oldUnitConfigs, newUnitConfigs []interface{}) []map[string]interface{} {
	if len(newUnitConfigs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(newUnitConfigs))
	for i, newUnitConfig := range newUnitConfigs {
		unitConfiElem := map[string]interface{}{
			// Required parameters.
			"image": buildV2ServiceUnitConfigImage(utils.PathSearch("image", newUnitConfig, make([]interface{}, 0)).([]interface{})),
			// Optional parameters.
			"role": utils.ValueIgnoreEmpty(utils.PathSearch("role", newUnitConfig, nil)),
			"custom_spec": utils.ValueIgnoreEmpty(buildV2ServiceUnitConfigCustomSpec(utils.PathSearch("custom_spec",
				newUnitConfig, make([]interface{}, 0)).([]interface{}))),
			"flavor": utils.ValueIgnoreEmpty(utils.PathSearch("flavor", newUnitConfig, nil)),
			"models": utils.ValueIgnoreEmpty(buildV2ServiceUnitConfigModels(utils.PathSearch("models",
				newUnitConfig, make([]interface{}, 0)).([]interface{}))),
			"codes": utils.ValueIgnoreEmpty(buildV2ServiceUnitConfigCodes(utils.PathSearch("codes",
				newUnitConfig, make([]interface{}, 0)).([]interface{}))),
			"count": utils.ValueIgnoreEmpty(utils.PathSearch("count", newUnitConfig, nil)),
			"cmd":   utils.ValueIgnoreEmpty(utils.PathSearch("cmd", newUnitConfig, nil)),
			"envs":  utils.ValueIgnoreEmpty(utils.PathSearch("envs", newUnitConfig, make(map[string]interface{})).(map[string]interface{})),
			"readiness_health": utils.ValueIgnoreEmpty(buildV2ServiceUnitConfigHealth(utils.PathSearch("readiness_health",
				newUnitConfig, make([]interface{}, 0)).([]interface{}))),
			"startup_health": utils.ValueIgnoreEmpty(buildV2ServiceUnitConfigHealth(utils.PathSearch("startup_health",
				newUnitConfig, make([]interface{}, 0)).([]interface{}))),
			"liveness_health": utils.ValueIgnoreEmpty(buildV2ServiceUnitConfigHealth(utils.PathSearch("liveness_health",
				newUnitConfig, make([]interface{}, 0)).([]interface{}))),
			"port":     utils.ValueIgnoreEmpty(utils.PathSearch("port", newUnitConfig, nil)),
			"recovery": utils.ValueIgnoreEmpty(utils.PathSearch("recovery", newUnitConfig, nil)),
		}
		if i < len(oldUnitConfigs) {
			unitConfiElem["id"] = utils.PathSearch("id", oldUnitConfigs[i], nil)
		}

		result = append(result, unitConfiElem)
	}

	return result
}

func buildV2ServiceUpgradeGroupConfigs(oldConfigVal, newConfigVal interface{}) []map[string]interface{} {
	oldConfigs := oldConfigVal.([]interface{})
	newConfigs := newConfigVal.([]interface{})

	result := make([]map[string]interface{}, 0, len(newConfigs))
	for _, newConfig := range newConfigs {
		configName := utils.PathSearch("name", newConfig, "").(string)
		result = append(result, utils.RemoveNil(map[string]interface{}{
			"id":        utils.ValueIgnoreEmpty(utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0].id", configName), oldConfigs, "")),
			"name":      utils.ValueIgnoreEmpty(utils.PathSearch("name", newConfig, nil)),
			"weight":    utils.ValueIgnoreEmpty(utils.PathSearch("weight", newConfig, nil)),
			"pool_id":   utils.ValueIgnoreEmpty(utils.PathSearch("pool_id", newConfig, nil)),
			"framework": utils.ValueIgnoreEmpty(utils.PathSearch("framework", newConfig, nil)),
			"count":     utils.ValueIgnoreEmpty(utils.PathSearch("count", newConfig, nil)),
			"unit_configs": utils.ValueIgnoreEmpty(buildV2ServiceUpgradeUnitConfigs(
				utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0].unit_configs", configName), oldConfigs, make([]interface{}, 0)).([]interface{}),
				utils.PathSearch("unit_configs", newConfig, make([]interface{}, 0)).([]interface{}),
			)),
		}))
	}

	return result
}

func buildV2ServiceUpgradeBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"id":             d.Id(),
		"version":        d.Get("version"),
		"group_configs":  buildV2ServiceUpgradeGroupConfigs(d.GetChange("group_configs")),
		"runtime_config": utils.StringToJson(d.Get("runtime_config").(string)),
		"upgrade_config": utils.StringToJson(d.Get("upgrade_config").(string)),
		"description":    d.Get("description"),
		"workspace_id":   d.Get("workspace_id"),
		"deploy_type":    d.Get("deploy_type"),
		"log_configs":    buildV2ServiceLogConfigs(d.Get("log_configs").([]interface{})), // Cannot be omitted.
	}
}

func upgradeService(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl   = "v2/{project_id}/services/{service_id}"
		serviceId = d.Id()
	)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{service_id}", serviceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildV2ServiceUpgradeBodyParams(d),
	}

	_, err := client.Request("PUT", createPath, &opt)
	if err != nil {
		return fmt.Errorf("error upgrading service (%s): %s", serviceId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      serviceStatusRefreshFunc(client, serviceId, []string{"STOPPED", "RUNNING"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the switch version operation to become complete: %s", err)
	}
	return err
}

func resourceV2ServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		historyVersions = d.Get("history_versions").(map[string]interface{})
		targetVersion   = d.Get("version").(string)
		historyDesc     string
	)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	// Tags do not distinguish between versions, and all versions share the same set of tags.
	if d.HasChange("tags") {
		err = updateV2ServiceTags(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("version") {
		if !isHistoryServiceVersion(historyVersions, targetVersion) {
			err = upgradeService(ctx, client, d)
			if err != nil {
				return diag.FromErr(err)
			}
			// No other actions after upgrading the service.
			return resourceV2ServiceRead(ctx, d, meta)
		}

		// If the current version number belongs to a historical version, perform version switching.
		// Try to assert the version ID to prevent illegal version value ​​from entering subsequent request build.
		targetVersionId, ok := historyVersions[targetVersion].(string)
		if !ok || targetVersionId == "" {
			return diag.Errorf("target history version ID (corresponding version is: %s) does not found: %#v",
				targetVersion, historyVersions)
		}
		historyServiceInfo, err := switchServiceVersion(ctx, client, d, targetVersion, targetVersionId)
		if err != nil {
			return diag.FromErr(err)
		}
		historyDesc = utils.PathSearch("version.description", historyServiceInfo, "").(string)
	}

	// After the version switch is completed, it is allowed to change the description of the corresponding version.
	// If only change the version number but do not modify the description of the resource in the script, switch
	// version only in this change and show description change in the next apply (if it has).
	if d.HasChange("description") && historyDesc != d.Get("description").(string) {
		err = updateServiceDescription(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
		// Updating the service description will not trigger the status change (whatever the version is current or
		// history).
	}

	return resourceV2ServiceRead(ctx, d, meta)
}

func resourceV2ServiceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v2/{project_id}/services/delete"
		serviceId = d.Id()
	)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"service_ids": []interface{}{
				serviceId,
			},
		},
	}

	_, err = client.Request("POST", deletePath, &opt)
	if err != nil {
		return diag.Errorf("error deleting service (%s): %s", serviceId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      serviceStatusRefreshFunc(client, serviceId, nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        2 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the service delete operation to become complete: %s", err)
	}
	return diag.FromErr(err)
}
