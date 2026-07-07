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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	trainingJobNonUpdatableParams = []string{
		"kind",
		"metadata.*.name",
		"metadata.*.workspace_id",
		"metadata.*.training_experiment_reference",
		"metadata.*.training_experiment_reference.*.id",
		"metadata.*.annotations",
		"algorithm",
		"algorithm.*.id",
		"algorithm.*.subscription_id",
		"algorithm.*.item_version_id",
		"algorithm.*.code_dir",
		"algorithm.*.boot_file",
		"algorithm.*.autosearch_config_path",
		"algorithm.*.autosearch_framework_path",
		"algorithm.*.command",
		"algorithm.*.local_code_dir",
		"algorithm.*.working_dir",
		"algorithm.*.parameters",
		"algorithm.*.parameters.*.name",
		"algorithm.*.parameters.*.value",
		"algorithm.*.parameters.*.description",
		"algorithm.*.parameters.*.constraint",
		"algorithm.*.parameters.*.constraint.*.type",
		"algorithm.*.parameters.*.constraint.*.editable",
		"algorithm.*.parameters.*.constraint.*.required",
		"algorithm.*.parameters.*.constraint.*.valid_type",
		"algorithm.*.parameters.*.constraint.*.valid_range",
		"algorithm.*.engine",
		"algorithm.*.engine.*.engine_id",
		"algorithm.*.engine.*.engine_name",
		"algorithm.*.engine.*.engine_version",
		"algorithm.*.engine.*.image_url",
		"algorithm.*.engine.*.install_sys_packages",
		"algorithm.*.environments",
		"algorithm.*.inputs",
		"algorithm.*.inputs.*.name",
		"algorithm.*.inputs.*.description",
		"algorithm.*.inputs.*.local_dir",
		"algorithm.*.inputs.*.access_method",
		"algorithm.*.inputs.*.remote",
		"algorithm.*.inputs.*.remote.*.dataset",
		"algorithm.*.inputs.*.remote.*.dataset.*.id",
		"algorithm.*.inputs.*.remote.*.dataset.*.version_id",
		"algorithm.*.inputs.*.remote.*.dataset.*.service_type",
		"algorithm.*.inputs.*.remote.*.dataset.*.name",
		"algorithm.*.inputs.*.remote.*.dataset.*.dataset_proportion",
		"algorithm.*.inputs.*.remote.*.obs",
		"algorithm.*.inputs.*.remote.*.obs.*.obs_url",
		"algorithm.*.outputs",
		"algorithm.*.outputs.*.name",
		"algorithm.*.outputs.*.description",
		"algorithm.*.outputs.*.local_dir",
		"algorithm.*.outputs.*.access_method",
		"algorithm.*.outputs.*.remote",
		"algorithm.*.outputs.*.remote.*.obs",
		"algorithm.*.outputs.*.remote.*.obs.*.obs_url",
		"spec",
		"spec.*.resource",
		"spec.*.resource.*.flavor_id",
		"spec.*.resource.*.node_count",
		"spec.*.resource.*.pool_id",
		"spec.*.resource.*.pool_group_id",
		"spec.*.resource.*.main_container_customized_flavor",
		"spec.*.resource.*.main_container_customized_flavor.*.cpu_core_num",
		"spec.*.resource.*.main_container_customized_flavor.*.mem_size",
		"spec.*.resource.*.main_container_customized_flavor.*.accelerator_num",
		"spec.*.runtime_type",
		"spec.*.log_export_path",
		"spec.*.log_export_path.*.obs_url",
		"spec.*.log_export_path.*.host_path",
		"spec.*.auto_stop",
		"spec.*.auto_stop.*.time_unit",
		"spec.*.auto_stop.*.duration",
		"spec.*.schedule_policy",
		"spec.*.schedule_policy.*.priority",
		"spec.*.schedule_policy.*.preemptible",
		"spec.*.schedule_policy.*.required_affinity",
		"spec.*.schedule_policy.*.required_affinity.*.affinity_type",
		"spec.*.schedule_policy.*.required_affinity.*.affinity_group_size",
		"spec.*.schedule_policy.*.required_affinity.*.node_affinity",
		"spec.*.schedule_policy.*.required_affinity.*.node_affinity.*.node_selector_terms",
		"spec.*.schedule_policy.*.required_affinity.*.node_affinity.*.node_selector_terms.*.match_expressions",
		"spec.*.schedule_policy.*.required_affinity.*.node_affinity.*.node_selector_terms.*.match_expressions.*.key",
		"spec.*.schedule_policy.*.required_affinity.*.node_affinity.*.node_selector_terms.*.match_expressions.*.operator",
		"spec.*.schedule_policy.*.required_affinity.*.node_affinity.*.node_selector_terms.*.match_expressions.*.values",
		"spec.*.schedule_policy.*.required_affinity.*.node_affinity.*.node_selector_terms.*.match_fields",
		"spec.*.schedule_policy.*.required_affinity.*.node_affinity.*.node_selector_terms.*.match_fields.*.key",
		"spec.*.schedule_policy.*.required_affinity.*.node_affinity.*.node_selector_terms.*.match_fields.*.operator",
		"spec.*.schedule_policy.*.required_affinity.*.node_affinity.*.node_selector_terms.*.match_fields.*.values",
		"spec.*.schedule_policy.*.preferred_affinity",
		"spec.*.schedule_policy.*.preferred_affinity.*.node_affinity",
		"spec.*.schedule_policy.*.preferred_affinity.*.node_affinity.*.weight",
		"spec.*.schedule_policy.*.preferred_affinity.*.node_affinity.*.preference",
		"spec.*.schedule_policy.*.preferred_affinity.*.node_affinity.*.preference.*.match_expressions",
		"spec.*.schedule_policy.*.preferred_affinity.*.node_affinity.*.preference.*.match_expressions.*.key",
		"spec.*.schedule_policy.*.preferred_affinity.*.node_affinity.*.preference.*.match_expressions.*.operator",
		"spec.*.schedule_policy.*.preferred_affinity.*.node_affinity.*.preference.*.match_expressions.*.values",
		"spec.*.schedule_policy.*.preferred_affinity.*.node_affinity.*.preference.*.match_fields",
		"spec.*.schedule_policy.*.preferred_affinity.*.node_affinity.*.preference.*.match_fields.*.key",
		"spec.*.schedule_policy.*.preferred_affinity.*.node_affinity.*.preference.*.match_fields.*.operator",
		"spec.*.schedule_policy.*.preferred_affinity.*.node_affinity.*.preference.*.match_fields.*.values",
		"spec.*.log_export_config",
		"spec.*.log_export_config.*.version",
		"spec.*.log_export_config.*.rotation_enabled",
		"spec.*.notification",
		"spec.*.notification.*.topic_urn",
		"spec.*.notification.*.events",
		"spec.*.custom_metrics",
		"spec.*.custom_metrics.*.exec",
		"spec.*.custom_metrics.*.exec.*.command",
		"spec.*.custom_metrics.*.http_get",
		"spec.*.custom_metrics.*.http_get.*.path",
		"spec.*.custom_metrics.*.http_get.*.port",
		"spec.*.output_model",
		"spec.*.output_model.*.obs",
		"spec.*.output_model.*.obs.*.obs_path",
		"spec.*.output_model.*.obs.*.local_path",
		"spec.*.asset_model",
		"spec.*.asset_model.*.name",
		"spec.*.asset_model.*.code",
		"spec.*.asset_model.*.version",
		"spec.*.asset_model.*.desc",
		"spec.*.asset_model.*.series",
		"spec.*.asset_model.*.type",
		"spec.*.asset_id",
		"spec.*.volumes",
		"spec.*.volumes.*.nfs",
		"spec.*.volumes.*.nfs.*.nfs_server_path",
		"spec.*.volumes.*.nfs.*.local_path",
		"spec.*.volumes.*.nfs.*.read_only",
		"spec.*.volumes.*.pfs",
		"spec.*.volumes.*.pfs.*.pfs_path",
		"spec.*.volumes.*.pfs.*.local_path",
		"spec.*.volumes.*.obs",
		"spec.*.volumes.*.obs.*.obs_path",
		"spec.*.volumes.*.obs.*.local_path",
		"train_type",
		"ftjob_config",
		"ftjob_config.*.ft_job_uuid",
		"ftjob_config.*.ft_train_type",
		"ftjob_config.*.model_type",
		"ftjob_config.*.train_output_path",
		"ftjob_config.*.train_process",
		"ftjob_config.*.checkpoint_id",
		"ftjob_config.*.task_env",
		"ftjob_config.*.task_env.*.envs",
		"ftjob_config.*.task_env.*.envs.*.label",
		"ftjob_config.*.task_env.*.envs.*.des",
		"ftjob_config.*.task_env.*.envs.*.env_name",
		"ftjob_config.*.task_env.*.envs.*.env_type",
		"ftjob_config.*.task_env.*.envs.*.value",
		"ftjob_config.*.task_env.*.envs.*.modifiable",
		"ftjob_config.*.task_env.*.envs.*.displayable",
		"ftjob_config.*.task_env.*.envs.*.used_steps",
		"ftjob_config.*.checkpoint_config",
		"ftjob_config.*.checkpoint_config.*.checkpoint_id",
		"ftjob_config.*.checkpoint_config.*.save_checkpoints_max",
		"ftjob_config.*.checkpoint_config.*.skipped_steps",
		"ftjob_config.*.checkpoint_config.*.restore_training",
		"endpoints",
		"endpoints.*.ssh",
		"endpoints.*.ssh.*.key_pair_names",
	}
	trainingJobNotFoundErrCodes = []string{
		"ModelArts.2775", // The resource does not exist.
	}
)

// @API ModelArts POST /v2/{project_id}/training-jobs
// @API ModelArts GET /v2/{project_id}/training-jobs/{training_job_id}
// @API ModelArts PUT /v2/{project_id}/training-jobs/{training_job_id}
// @API ModelArts DELETE /v2/{project_id}/training-jobs/{training_job_id}
// @API ModelArts POST /v2/{project_id}/modelarts-training-job/{training_job_id}/tags/create
// @API ModelArts DELETE /v2/{project_id}/modelarts-training-job/{training_job_id}/tags/delete
// @API ModelArts GET /v2/{project_id}/modelarts-training-job/{training_job_id}/tags
func ResourceTrainingJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTrainingJobCreate,
		ReadContext:   resourceTrainingJobRead,
		UpdateContext: resourceTrainingJobUpdate,
		DeleteContext: resourceTrainingJobDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(trainingJobNonUpdatableParams),
			config.MergeDefaultTags(),
		),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the training job is located.`,
			},

			// Required parameters.
			"kind": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the training job.`,
			},
			"metadata": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the training job.`,
						},
						"workspace_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The ID of the workspace to which the training job belongs.`,
						},
						"training_experiment_reference": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The ID of the experiment.`,
									},
								},
							},
							Description: `The experiment configuration to be associated with the training job.`,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The description of the training job.`,
						},
						"annotations": {
							Type:        schema.TypeMap,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The advanced feature configuration of the training job.`,
						},
					},
				},
				Description: `The metadata configuration of the training job.`,
			},
			"algorithm": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The ID of the algorithm.`,
						},
						"subscription_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The subscription ID of the subscribed algorithm.`,
						},
						"item_version_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The version ID of the subscribed algorithm.`,
						},
						"code_dir": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The code directory of the training job.`,
						},
						"boot_file": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The boot file of the training job code.`,
						},
						"autosearch_config_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The YAML configuration path of the auto search job.`,
						},
						"autosearch_framework_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The framework code directory of the auto search job.`,
						},
						"command": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The container startup command for custom image scenarios.`,
						},
						"local_code_dir": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The local path where the algorithm code is downloaded in the training container.`,
						},
						"working_dir": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The working directory when running the algorithm.`,
						},
						"engine": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Elem:        trainingJobAlgorithmEngineSchema(),
							Description: `The engine configuration of the training job.`,
						},
						"inputs": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        trainingJobAlgorithmInputSchema(),
							Description: `The data input list of the training job.`,
						},
						"outputs": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        trainingJobAlgorithmOutputSchema(),
							Description: `The data output list of the training job.`,
						},
						"parameters": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        trainingJobAlgorithmParameterSchema(),
							Description: `The runtime parameter list of the training job.`,
						},
						"environments": {
							Type:        schema.TypeMap,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The environment variables of the training job.`,
						},
					},
				},
				Description: `The algorithm configuration of the training job.`,
			},
			"spec": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Elem:        trainingJobSpecResourceSchema(),
							Description: `The resource specification of the training job.`,
						},
						"runtime_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The runtime type of the training job.`,
						},
						"log_export_path": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"obs_url": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `The OBS path where training job logs are saved.`,
									},
									"host_path": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `The host path where training job logs are saved.`,
									},
								},
							},
							Description: `The log export path configuration of the training job.`,
						},
						"log_export_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The log version.`,
									},
									"rotation_enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `Whether to enable log rotation download.`,
									},
								},
							},
							Description: `The log export configuration of the training job.`,
						},
						"auto_stop": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time_unit": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The time unit of the auto stop duration.`,
									},
									"duration": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: `The running duration before the training job is automatically stopped.`,
									},
								},
							},
							Description: `The auto stop configuration of the training job.`,
						},
						"schedule_policy": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Elem:        trainingJobSpecSchedulePolicySchema(),
							Description: `The schedule policy configuration of the training job.`,
						},
						"notification": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"topic_urn": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The URN of the SMN topic for training event notifications.`,
									},
									"events": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `The training events that trigger notifications.`,
									},
								},
							},
							Description: `The notification configuration of the training job.`,
						},
						"custom_metrics": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        trainingJobSpecCustomMetricsSchema(),
							Description: `The custom metrics collection configuration of the training job.`,
						},
						"output_model": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Elem:        trainingJobSpecOutputModelSchema(),
							Description: `The output model configuration of the training job.`,
						},
						"asset_model": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The model name.`,
									},
									"version": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The model version.`,
									},
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The model type.`,
									},
									"code": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The model code.`,
									},
									"desc": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The model description.`,
									},
									"series": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The model series.`,
									},
								},
							},
							Description: `The asset model configuration of the training job.`,
						},
						"asset_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The asset model ID for fine-tuning training jobs.`,
						},
						"volumes": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        trainingJobSpecVolumeSchema(),
							Description: `The volume mount configuration of the training job.`,
						},
					},
				},
				Description: `The specification configuration of the training job.`,
			},

			// Optional parameters.
			"endpoints": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        trainingJobEndpointsSchema(),
				Description: `The remote access configuration of the training job.`,
			},
			"train_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The training type of the fine-tuning job.`,
			},
			"ftjob_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        trainingJobFtjobConfigSchema(),
				Description: `The fine-tuning training job configuration.`,
			},
			"tags": common.TagsSchema(`The key/value pairs to associate with the training job.`),

			// Attributes.
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the training job, in RFC3339 format.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current status of the training job.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{Internal: true},
				),
			},
		},
	}
}

func trainingJobAlgorithmParameterSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The parameter name.`,
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The parameter value.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The parameter description.`,
			},
			"constraint": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The parameter type.`,
						},
						"editable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether the parameter is editable.`,
						},
						"required": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether the parameter is required.`,
						},
						"valid_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The valid type of the parameter.`,
						},
						"valid_range": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The valid range of the parameter.`,
						},
					},
				},
				Description: `The parameter constraint configuration.`,
			},
		},
	}
}

func trainingJobAlgorithmEngineSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"engine_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The engine specification ID.`,
			},
			"engine_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The engine specification name.`,
			},
			"engine_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The engine specification version.`,
			},
			"image_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The custom image URL obtained from SWR.`,
			},
			"install_sys_packages": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to install the moxing version specified by the training platform.`,
			},
		},
	}
}

func trainingJobAlgorithmInputSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"remote": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dataset": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The dataset ID.`,
									},
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The dataset name.`,
									},
									"version_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The dataset version ID.`,
									},
									"service_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The dataset service type.`,
									},
									"dataset_proportion": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: `The dataset proportion used for fine-tuning training jobs.`,
									},
								},
							},
							Description: `The dataset input configuration.`,
						},
						"obs": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"obs_url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The OBS path URL of the dataset.`,
									},
								},
							},
							Description: `The OBS input configuration.`,
						},
					},
				},
				Description: `The actual input data configuration.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The input channel name.`,
			},
			"local_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The local path mapped by the input channel in the container.`,
			},
			"access_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The delivery method of the input channel path.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The input channel description.`,
			},
		},
	}
}

func trainingJobAlgorithmOutputSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The output channel name.`,
			},
			"remote": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"obs": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"obs_url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The OBS path URL where the data is output.`,
									},
								},
							},
							Description: `The OBS output configuration.`,
						},
					},
				},
				Description: `The actual output data configuration.`,
			},
			"local_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The local path mapped by the output channel in the container.`,
			},
			"access_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The delivery method of the output channel path.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The output channel description.`,
			},
		},
	}
}

func trainingJobSpecResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The number of resource replicas used by the training job.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The resource flavor ID of the training job.`,
			},
			"pool_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The dedicated resource pool ID.`,
			},
			"pool_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The federated resource pool ID.`,
			},
			"main_container_customized_flavor": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu_core_num": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Description: `The number of CPU cores.`,
						},
						"mem_size": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Description: `The memory size.`,
						},
						"accelerator_num": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Description: `The number of accelerator cards.`,
						},
					},
				},
				Description: `The customized flavor configuration of the main container.`,
			},
		},
	}
}

func trainingJobSpecSchedulePolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The priority of the training job.`,
			},
			"preemptible": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether the training job can be preempted.`,
			},
			"required_affinity": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"affinity_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The affinity scheduling policy type.`,
						},
						"affinity_group_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `The affinity group size.`,
						},
						"node_affinity": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"node_selector_terms": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        trainingJobSpecSchedulePolicyNodeAffinityNodeSelectorTermsSchema(),
										Description: `The node selector term list.`,
									},
								},
							},
							Description: `The node affinity configuration.`,
						},
					},
				},
				Description: `The required affinity configuration of the training job.`,
			},
			"preferred_affinity": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_affinity": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"weight": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: `The weight associated with the preferred node selector term.`,
									},
									"preference": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Elem:        trainingJobSpecSchedulePolicyNodeAffinityNodeSelectorTermsSchema(),
										Description: `The preferred node selector term.`,
									},
								},
							},
							Description: `The preferred node affinity terms.`,
						},
					},
				},
				Description: `The preferred affinity configuration of the training job.`,
			},
		},
	}
}

func trainingJobSpecSchedulePolicyNodeAffinityNodeSelectorTermsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"match_expressions": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        trainingJobNodeSelectorRequirementSchema(),
				Description: `The node selector requirements based on node labels.`,
			},
			"match_fields": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        trainingJobNodeSelectorRequirementSchema(),
				Description: `The node selector requirements based on node fields.`,
			},
		},
	}
}

func trainingJobNodeSelectorRequirementSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The label key used by the node selector requirement.`,
			},
			"operator": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The operator used by the node selector requirement.`,
			},
			"values": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The label values used by the node selector requirement.`,
			},
		},
	}
}

func trainingJobSpecCustomMetricsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"exec": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"command": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The command used to collect metrics.`,
						},
					},
				},
				Description: `The command-based metrics collection configuration.`,
			},
			"http_get": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The HTTP path used to collect metrics.`,
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The HTTP port used to collect metrics.`,
						},
					},
				},
				Description: `The HTTP-based metrics collection configuration.`,
			},
		},
	}
}

func trainingJobSpecOutputModelSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"obs": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"obs_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The OBS path where the output model is saved.`,
						},
						"local_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The local path where the output model is saved.`,
						},
					},
				},
				Description: `The OBS output configuration of the model.`,
			},
		},
	}
}

func trainingJobSpecVolumeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"nfs": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nfs_server_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The NFS server path.`,
						},
						"local_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The local mount path in the training container.`,
						},
						"read_only": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether the NFS volume is read-only in the container.`,
						},
					},
				},
				Description: `The NFS volume mount configuration.`,
			},
			"pfs": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pfs_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The OBSFS path.`,
						},
						"local_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The local mount path in the training container.`,
						},
					},
				},
				Description: `The PFS volume mount configuration.`,
			},
			"obs": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"obs_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The OBS path to be mounted.`,
						},
						"local_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The local mount path in the training container.`,
						},
					},
				},
				Description: `The OBS volume mount configuration.`,
			},
		},
	}
}

func trainingJobEndpointsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ssh": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_pair_names": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The SSH key pair names.`,
						},
					},
				},
				Description: `The SSH connection configuration.`,
			},
		},
	}
}

func trainingJobFtjobConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ft_job_uuid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The model ID.`,
			},
			"ft_train_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The model training type.`,
			},
			"model_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The model type.`,
			},
			"train_output_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The output path of the training job.`,
			},
			"train_process": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: `The training process progress.`,
			},
			"checkpoint_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The checkpoint ID.`,
			},
			"task_env": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        trainingJobFtjobConfigTaskEnvSchema(),
				Description: `The fine-tuning training job environment parameters.`,
			},
			"checkpoint_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"checkpoint_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The checkpoint ID.`,
						},
						"save_checkpoints_max": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `The maximum number of checkpoints to save.`,
						},
						"skipped_steps": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `The number of steps to skip.`,
						},
						"restore_training": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `Whether to restore training from a checkpoint.`,
						},
					},
				},
				Description: `The checkpoint configuration.`,
			},
		},
	}
}

func trainingJobFtjobConfigTaskEnvSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"envs": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The label of the environment variable.`,
						},
						"des": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The description of the environment variable.`,
						},
						"env_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The name of the environment variable.`,
						},
						"env_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The type of the environment variable.`,
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The value of the environment variable.`,
						},
						"modifiable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether the environment variable is modifiable.`,
						},
						"displayable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether the environment variable is displayable.`,
						},
						"used_steps": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The steps where the environment variable is used.`,
						},
					},
				},
				Description: `The fine-tuning training environment variables.`,
			},
		},
	}
}

func buildTrainingJobCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	preemptible := utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "spec.0.schedule_policy.0.preemptible")
	return map[string]interface{}{
		"kind":         d.Get("kind"),
		"metadata":     buildTrainingJobMetadata(d.Get("metadata.0")),
		"algorithm":    buildTrainingJobAlgorithm(d.Get("algorithm.0")),
		"spec":         buildTrainingJobSpec(d.Get("spec.0"), preemptible),
		"endpoints":    buildTrainingJobEndpoints(d.Get("endpoints").([]interface{})),
		"train_type":   utils.ValueIgnoreEmpty(d.Get("train_type")),
		"ftjob_config": buildTrainingJobFtjobConfig(d.Get("ftjob_config").([]interface{})),
	}
}

func buildTrainingJobEndpoints(endpoints []interface{}) map[string]interface{} {
	if len(endpoints) == 0 || endpoints[0] == nil {
		return nil
	}

	return map[string]interface{}{
		"ssh": buildTrainingJobEndpointsSSH(utils.PathSearch("ssh[0]",
			endpoints[0], make(map[string]interface{}, 0)).(map[string]interface{})),
	}
}

func buildTrainingJobEndpointsSSH(ssh map[string]interface{}) map[string]interface{} {
	if len(ssh) < 1 {
		return nil
	}

	return map[string]interface{}{
		"key_pair_names": utils.ExpandToStringList(utils.PathSearch("key_pair_names", ssh, make([]interface{}, 0)).([]interface{})),
	}
}

func buildTrainingJobFtjobConfig(ftjobConfig []interface{}) map[string]interface{} {
	if len(ftjobConfig) == 0 || ftjobConfig[0] == nil {
		return nil
	}

	jobConfig := ftjobConfig[0]
	return map[string]interface{}{
		"ft_job_uuid":       utils.ValueIgnoreEmpty(utils.PathSearch("ft_job_uuid", jobConfig, nil)),
		"ft_train_type":     utils.ValueIgnoreEmpty(utils.PathSearch("ft_train_type", jobConfig, nil)),
		"model_type":        utils.ValueIgnoreEmpty(utils.PathSearch("model_type", jobConfig, nil)),
		"train_output_path": utils.ValueIgnoreEmpty(utils.PathSearch("train_output_path", jobConfig, nil)),
		"train_process":     utils.ValueIgnoreEmpty(utils.PathSearch("train_process", jobConfig, nil)),
		"checkpoint_id":     utils.ValueIgnoreEmpty(utils.PathSearch("checkpoint_id", jobConfig, nil)),
		"task_env": buildTrainingJobFtjobConfigTaskEnv(utils.PathSearch("task_env[0]",
			jobConfig, make(map[string]interface{}, 0)).(map[string]interface{})),
		"checkpoint_config": buildTrainingJobFtjobConfigCheckpointConfig(utils.PathSearch("checkpoint_config[0]",
			jobConfig, make(map[string]interface{}, 0)).(map[string]interface{})),
	}
}

func buildTrainingJobFtjobConfigTaskEnv(taskEnv map[string]interface{}) map[string]interface{} {
	if len(taskEnv) < 1 {
		return nil
	}

	return map[string]interface{}{
		"envs": buildTrainingJobFtjobConfigEnvs(utils.PathSearch("envs", taskEnv, make([]interface{}, 0)).([]interface{})),
	}
}

func buildTrainingJobFtjobConfigEnvs(envs []interface{}) []map[string]interface{} {
	if len(envs) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(envs))
	for _, env := range envs {
		result = append(result, map[string]interface{}{
			"label":       utils.ValueIgnoreEmpty(utils.PathSearch("label", env, nil)),
			"des":         utils.ValueIgnoreEmpty(utils.PathSearch("des", env, nil)),
			"env_name":    utils.ValueIgnoreEmpty(utils.PathSearch("env_name", env, nil)),
			"env_type":    utils.ValueIgnoreEmpty(utils.PathSearch("env_type", env, nil)),
			"value":       utils.ValueIgnoreEmpty(utils.PathSearch("value", env, nil)),
			"modifiable":  utils.PathSearch("modifiable", env, nil),
			"displayable": utils.PathSearch("displayable", env, nil),
			"used_steps": utils.ValueIgnoreEmpty(utils.ExpandToStringList(utils.PathSearch("used_steps",
				env, make([]interface{}, 0)).([]interface{}))),
		})
	}

	return result
}

func buildTrainingJobFtjobConfigCheckpointConfig(checkpointConfig map[string]interface{}) map[string]interface{} {
	if len(checkpointConfig) < 1 {
		return nil
	}

	return map[string]interface{}{
		"checkpoint_id": utils.ValueIgnoreEmpty(utils.PathSearch("checkpoint_id", checkpointConfig, nil)),
		// The following fields support filling 0.
		"save_checkpoints_max": utils.PathSearch("save_checkpoints_max", checkpointConfig, nil),
		"skipped_steps":        utils.PathSearch("skipped_steps", checkpointConfig, nil),
		"restore_training":     utils.PathSearch("restore_training", checkpointConfig, nil),
	}
}

func buildTrainingJobMetadata(metadata interface{}) map[string]interface{} {
	return map[string]interface{}{
		"name":         utils.PathSearch("name", metadata, nil),
		"workspace_id": utils.ValueIgnoreEmpty(utils.PathSearch("workspace_id", metadata, nil)),
		"training_experiment_reference": buildTrainingJobMetadataTrainingExperimentReference(
			utils.PathSearch("training_experiment_reference[0]", metadata, nil)),
		"description": utils.ValueIgnoreEmpty(utils.PathSearch("description", metadata, nil)),
		"annotations": utils.ValueIgnoreEmpty(utils.PathSearch("annotations", metadata, nil)),
	}
}

func buildTrainingJobMetadataTrainingExperimentReference(experimentReference interface{}) map[string]interface{} {
	if experimentReference == nil {
		return nil
	}

	return map[string]interface{}{
		"id": utils.PathSearch("id", experimentReference, nil),
	}
}

func buildTrainingJobAlgorithm(algorithm interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":                        utils.ValueIgnoreEmpty(utils.PathSearch("id", algorithm, nil)),
		"subscription_id":           utils.ValueIgnoreEmpty(utils.PathSearch("subscription_id", algorithm, nil)),
		"item_version_id":           utils.ValueIgnoreEmpty(utils.PathSearch("item_version_id", algorithm, nil)),
		"code_dir":                  utils.ValueIgnoreEmpty(utils.PathSearch("code_dir", algorithm, nil)),
		"boot_file":                 utils.ValueIgnoreEmpty(utils.PathSearch("boot_file", algorithm, nil)),
		"autosearch_config_path":    utils.ValueIgnoreEmpty(utils.PathSearch("autosearch_config_path", algorithm, nil)),
		"autosearch_framework_path": utils.ValueIgnoreEmpty(utils.PathSearch("autosearch_framework_path", algorithm, nil)),
		"command":                   utils.ValueIgnoreEmpty(utils.PathSearch("command", algorithm, nil)),
		"local_code_dir":            utils.ValueIgnoreEmpty(utils.PathSearch("local_code_dir", algorithm, nil)),
		"working_dir":               utils.ValueIgnoreEmpty(utils.PathSearch("working_dir", algorithm, nil)),
		"engine": buildTrainingJobAlgorithmEngine(utils.PathSearch("engine[0]",
			algorithm, make(map[string]interface{}, 0)).(map[string]interface{})),
		"inputs": buildTrainingJobAlgorithmInputs(utils.PathSearch("inputs",
			algorithm, make([]interface{}, 0)).([]interface{})),
		"outputs": buildTrainingJobAlgorithmOutputs(utils.PathSearch("outputs",
			algorithm, make([]interface{}, 0)).([]interface{})),
		"parameters": buildTrainingJobAlgorithmParameters(utils.PathSearch("parameters",
			algorithm, make([]interface{}, 0)).([]interface{})),
		"environments": utils.PathSearch("environments", algorithm, nil),
	}
}

func buildTrainingJobAlgorithmEngine(engine map[string]interface{}) map[string]interface{} {
	if len(engine) < 1 {
		return nil
	}

	return map[string]interface{}{
		"engine_id":            utils.ValueIgnoreEmpty(utils.PathSearch("engine_id", engine, nil)),
		"engine_name":          utils.ValueIgnoreEmpty(utils.PathSearch("engine_name", engine, nil)),
		"engine_version":       utils.ValueIgnoreEmpty(utils.PathSearch("engine_version", engine, nil)),
		"image_url":            utils.ValueIgnoreEmpty(utils.PathSearch("image_url", engine, nil)),
		"install_sys_packages": utils.PathSearch("install_sys_packages", engine, nil),
	}
}

func buildTrainingJobAlgorithmInputs(inputs []interface{}) []map[string]interface{} {
	if len(inputs) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(inputs))
	for _, input := range inputs {
		result = append(result, map[string]interface{}{
			"remote": buildTrainingJobAlgorithmInputRemote(utils.PathSearch("remote[0]",
				input, make(map[string]interface{}, 0)).(map[string]interface{})),
			"name":          utils.ValueIgnoreEmpty(utils.PathSearch("name", input, nil)),
			"local_dir":     utils.ValueIgnoreEmpty(utils.PathSearch("local_dir", input, nil)),
			"access_method": utils.ValueIgnoreEmpty(utils.PathSearch("access_method", input, nil)),
			"description":   utils.ValueIgnoreEmpty(utils.PathSearch("description", input, nil)),
		})
	}

	return result
}

func buildTrainingJobAlgorithmInputRemote(remote map[string]interface{}) map[string]interface{} {
	if len(remote) < 1 {
		return nil
	}

	return map[string]interface{}{
		"dataset": buildTrainingJobAlgorithmInputDataset(utils.PathSearch("dataset[0]",
			remote, make(map[string]interface{}, 0)).(map[string]interface{})),
		"obs": buildTrainingJobAlgorithmInputObs(utils.PathSearch("obs[0]",
			remote, make(map[string]interface{}, 0)).(map[string]interface{})),
	}
}

func buildTrainingJobAlgorithmInputDataset(dataset map[string]interface{}) map[string]interface{} {
	if len(dataset) < 1 {
		return nil
	}

	return map[string]interface{}{
		"id":                 utils.PathSearch("id", dataset, nil),
		"name":               utils.ValueIgnoreEmpty(utils.PathSearch("name", dataset, nil)),
		"version_id":         utils.ValueIgnoreEmpty(utils.PathSearch("version_id", dataset, nil)),
		"service_type":       utils.ValueIgnoreEmpty(utils.PathSearch("service_type", dataset, nil)),
		"dataset_proportion": utils.ValueIgnoreEmpty(utils.PathSearch("dataset_proportion", dataset, nil)),
	}
}

func buildTrainingJobAlgorithmInputObs(obs map[string]interface{}) map[string]interface{} {
	if len(obs) < 1 {
		return nil
	}

	return map[string]interface{}{
		"obs_url": utils.PathSearch("obs_url", obs, nil),
	}
}

func buildTrainingJobAlgorithmOutputs(outputs []interface{}) []map[string]interface{} {
	if len(outputs) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(outputs))
	for _, output := range outputs {
		result = append(result, map[string]interface{}{
			"name": utils.PathSearch("name", output, nil),
			"remote": buildTrainingJobAlgorithmOutputRemote(utils.PathSearch("remote[0]",
				output, make(map[string]interface{}, 0)).(map[string]interface{})),
			"local_dir":     utils.ValueIgnoreEmpty(utils.PathSearch("local_dir", output, nil)),
			"access_method": utils.ValueIgnoreEmpty(utils.PathSearch("access_method", output, nil)),
			"description":   utils.ValueIgnoreEmpty(utils.PathSearch("description", output, nil)),
		})
	}

	return result
}

func buildTrainingJobAlgorithmOutputRemote(remote map[string]interface{}) map[string]interface{} {
	if len(remote) < 1 {
		return nil
	}

	return map[string]interface{}{
		"obs": map[string]interface{}{
			"obs_url": utils.PathSearch("obs[0].obs_url", remote, nil),
		},
	}
}

func buildTrainingJobAlgorithmParameters(parameters []interface{}) []map[string]interface{} {
	if len(parameters) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(parameters))
	for _, parameter := range parameters {
		result = append(result, map[string]interface{}{
			"name":        utils.ValueIgnoreEmpty(utils.PathSearch("name", parameter, nil)),
			"value":       utils.ValueIgnoreEmpty(utils.PathSearch("value", parameter, nil)),
			"description": utils.ValueIgnoreEmpty(utils.PathSearch("description", parameter, nil)),
			"constraint": buildTrainingJobAlgorithmParameterConstraint(utils.PathSearch("constraint[0]",
				parameter, make(map[string]interface{}, 0)).(map[string]interface{})),
		})
	}

	return result
}

func buildTrainingJobAlgorithmParameterConstraint(constraint map[string]interface{}) map[string]interface{} {
	if len(constraint) < 1 {
		return nil
	}

	return map[string]interface{}{
		"type":       utils.PathSearch("type", constraint, nil),
		"editable":   utils.PathSearch("editable", constraint, nil),
		"required":   utils.PathSearch("required", constraint, nil),
		"valid_type": utils.ValueIgnoreEmpty(utils.PathSearch("valid_type", constraint, nil)),
		"valid_range": utils.ValueIgnoreEmpty(utils.ExpandToStringList(utils.PathSearch("valid_range",
			constraint, make([]interface{}, 0)).([]interface{}))),
	}
}

func buildTrainingJobSpec(spec, preemptible interface{}) map[string]interface{} {
	return map[string]interface{}{
		"resource": buildTrainingJobSpecResource(utils.PathSearch("resource[0]",
			spec, make(map[string]interface{}, 0)).(map[string]interface{})),
		"runtime_type": utils.ValueIgnoreEmpty(utils.PathSearch("runtime_type", spec, nil)),
		"log_export_path": buildTrainingJobSpecLogExportPath(utils.PathSearch("log_export_path[0]",
			spec, make(map[string]interface{}, 0)).(map[string]interface{})),
		"log_export_config": buildTrainingJobSpecLogExportConfig(utils.PathSearch("log_export_config[0]",
			spec, make(map[string]interface{}, 0)).(map[string]interface{})),
		"auto_stop": buildTrainingJobSpecAutoStop(utils.PathSearch("auto_stop[0]",
			spec, make(map[string]interface{}, 0)).(map[string]interface{})),
		"schedule_policy": buildTrainingJobSpecSchedulePolicy(utils.PathSearch("schedule_policy[0]",
			spec, make(map[string]interface{}, 0)).(map[string]interface{}), preemptible),
		"notification": buildTrainingJobSpecNotification(utils.PathSearch("notification[0]",
			spec, make(map[string]interface{}, 0)).(map[string]interface{})),
		"custom_metrics": buildTrainingJobSpecCustomMetrics(utils.PathSearch("custom_metrics",
			spec, make([]interface{}, 0)).([]interface{})),
		"output_model": buildTrainingJobSpecOutputModel(utils.PathSearch("output_model[0]",
			spec, make(map[string]interface{}, 0)).(map[string]interface{})),
		"asset_model": buildTrainingJobSpecAssetModel(utils.PathSearch("asset_model[0]",
			spec, make(map[string]interface{}, 0)).(map[string]interface{})),
		"asset_id": utils.ValueIgnoreEmpty(utils.PathSearch("asset_id", spec, nil)),
		"volumes": buildTrainingJobSpecVolumes(utils.PathSearch("volumes",
			spec, make([]interface{}, 0)).([]interface{})),
	}
}

func buildTrainingJobSpecResource(specResource map[string]interface{}) map[string]interface{} {
	if len(specResource) < 1 {
		return nil
	}

	return map[string]interface{}{
		"node_count":    utils.PathSearch("node_count", specResource, nil),
		"flavor_id":     utils.ValueIgnoreEmpty(utils.PathSearch("flavor_id", specResource, nil)),
		"pool_id":       utils.ValueIgnoreEmpty(utils.PathSearch("pool_id", specResource, nil)),
		"pool_group_id": utils.ValueIgnoreEmpty(utils.PathSearch("pool_group_id", specResource, nil)),
		"main_container_customized_flavor": buildTrainingJobSpecMainContainerCustomizedFlavor(
			utils.PathSearch("main_container_customized_flavor[0]",
				specResource, make(map[string]interface{}, 0)).(map[string]interface{})),
	}
}

func buildTrainingJobSpecMainContainerCustomizedFlavor(flavor map[string]interface{}) map[string]interface{} {
	if len(flavor) < 1 {
		return nil
	}

	return map[string]interface{}{
		"cpu_core_num":    utils.ValueIgnoreEmpty(utils.PathSearch("cpu_core_num", flavor, nil)),
		"mem_size":        utils.ValueIgnoreEmpty(utils.PathSearch("mem_size", flavor, nil)),
		"accelerator_num": utils.ValueIgnoreEmpty(utils.PathSearch("accelerator_num", flavor, nil)),
	}
}

func buildTrainingJobSpecVolumes(volumes []interface{}) []map[string]interface{} {
	if len(volumes) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(volumes))
	for _, volume := range volumes {
		result = append(result, map[string]interface{}{
			"nfs": buildTrainingJobSpecVolumeNfs(utils.PathSearch("nfs[0]",
				volume, make(map[string]interface{}, 0)).(map[string]interface{})),
			"pfs": buildTrainingJobSpecVolumePfs(utils.PathSearch("pfs[0]",
				volume, make(map[string]interface{}, 0)).(map[string]interface{})),
			"obs": buildTrainingJobSpecVolumeObs(utils.PathSearch("obs[0]",
				volume, make(map[string]interface{}, 0)).(map[string]interface{})),
		})
	}

	return result
}

func buildTrainingJobSpecVolumeNfs(nfs map[string]interface{}) map[string]interface{} {
	if len(nfs) < 1 {
		return nil
	}

	return map[string]interface{}{
		"nfs_server_path": utils.ValueIgnoreEmpty(utils.PathSearch("nfs_server_path", nfs, nil)),
		"local_path":      utils.ValueIgnoreEmpty(utils.PathSearch("local_path", nfs, nil)),
		"read_only":       utils.PathSearch("read_only", nfs, nil),
	}
}

func buildTrainingJobSpecVolumePfs(pfs map[string]interface{}) map[string]interface{} {
	if len(pfs) < 1 {
		return nil
	}

	return map[string]interface{}{
		"pfs_path":   utils.ValueIgnoreEmpty(utils.PathSearch("pfs_path", pfs, nil)),
		"local_path": utils.ValueIgnoreEmpty(utils.PathSearch("local_path", pfs, nil)),
	}
}

func buildTrainingJobSpecVolumeObs(obs map[string]interface{}) map[string]interface{} {
	if len(obs) < 1 {
		return nil
	}

	return map[string]interface{}{
		"obs_path":   utils.ValueIgnoreEmpty(utils.PathSearch("obs_path", obs, nil)),
		"local_path": utils.ValueIgnoreEmpty(utils.PathSearch("local_path", obs, nil)),
	}
}

func buildTrainingJobSpecLogExportPath(logExportPath map[string]interface{}) map[string]interface{} {
	if len(logExportPath) < 1 {
		return nil
	}

	return utils.RemoveNil(map[string]interface{}{
		"obs_url":   utils.ValueIgnoreEmpty(utils.PathSearch("obs_url", logExportPath, nil)),
		"host_path": utils.ValueIgnoreEmpty(utils.PathSearch("host_path", logExportPath, nil)),
	})
}

func buildTrainingJobSpecAutoStop(autoStop map[string]interface{}) map[string]interface{} {
	if len(autoStop) < 1 {
		return nil
	}

	return map[string]interface{}{
		"time_unit": utils.PathSearch("time_unit", autoStop, nil),
		"duration":  utils.PathSearch("duration", autoStop, nil),
	}
}

func buildTrainingJobSpecSchedulePolicy(schedulePolicy map[string]interface{}, preemptible interface{}) map[string]interface{} {
	if len(schedulePolicy) < 1 {
		return nil
	}

	// For public resource pool, preemptible is not specified.
	return utils.RemoveNil(map[string]interface{}{
		"priority":    utils.ValueIgnoreEmpty(utils.PathSearch("priority", schedulePolicy, nil)),
		"preemptible": preemptible,
		"required_affinity": buildTrainingJobSpecSchedulePolicyRequiredAffinity(utils.PathSearch("required_affinity[0]",
			schedulePolicy, make(map[string]interface{}, 0)).(map[string]interface{})),
		"preferred_affinity": buildTrainingJobSpecSchedulePolicyPreferredAffinity(utils.PathSearch("preferred_affinity[0]",
			schedulePolicy, make(map[string]interface{}, 0)).(map[string]interface{})),
	})
}

func buildTrainingJobSpecSchedulePolicyRequiredAffinity(requiredAffinity map[string]interface{}) map[string]interface{} {
	if len(requiredAffinity) < 1 {
		return nil
	}

	return map[string]interface{}{
		"affinity_type":       utils.ValueIgnoreEmpty(utils.PathSearch("affinity_type", requiredAffinity, nil)),
		"affinity_group_size": utils.ValueIgnoreEmpty(utils.PathSearch("affinity_group_size", requiredAffinity, nil)),
		"node_affinity": buildTrainingJobSpecSchedulePolicyRequiredAffinityNodeAffinity(utils.PathSearch("node_affinity[0]",
			requiredAffinity, make(map[string]interface{}, 0)).(map[string]interface{})),
	}
}

func buildTrainingJobSpecSchedulePolicyRequiredAffinityNodeAffinity(nodeAffinity map[string]interface{}) map[string]interface{} {
	if len(nodeAffinity) == 0 {
		return nil
	}

	return map[string]interface{}{
		"nodeSelectorTerms": buildTrainingJobSpecSchedulePolicyNodeSelectorTerms(utils.PathSearch("node_selector_terms",
			nodeAffinity, make([]interface{}, 0)).([]interface{})),
	}
}

func buildTrainingJobSpecSchedulePolicyNodeSelectorTerms(nodeSelectorTerms []interface{}) []map[string]interface{} {
	if len(nodeSelectorTerms) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(nodeSelectorTerms))
	for _, term := range nodeSelectorTerms {
		result = append(result, map[string]interface{}{
			"matchExpressions": buildTrainingJobSchedulePolicyNodeSelectorRequirements(utils.PathSearch("match_expressions",
				term, make([]interface{}, 0)).([]interface{})),
			"matchFields": buildTrainingJobSchedulePolicyNodeSelectorRequirements(utils.PathSearch("match_fields",
				term, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func buildTrainingJobSpecSchedulePolicyPreferredAffinity(preferredAffinity map[string]interface{}) map[string]interface{} {
	if len(preferredAffinity) < 1 {
		return nil
	}

	return map[string]interface{}{
		"node_affinity": buildTrainingJobSpecSchedulePolicyPreferredAffinityNodeAffinity(utils.PathSearch("node_affinity",
			preferredAffinity, make([]interface{}, 0)).([]interface{})),
	}
}

func buildTrainingJobSpecSchedulePolicyPreferredAffinityNodeAffinity(affinities []interface{}) []map[string]interface{} {
	if len(affinities) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(affinities))
	for _, term := range affinities {
		result = append(result, map[string]interface{}{
			"weight": term.(map[string]interface{})["weight"],
			"preference": buildTrainingJobSpecSchedulePolicyPreferredAffinityNodeAffinityPreference(utils.PathSearch("preference[0]",
				term, make(map[string]interface{}, 0)).(map[string]interface{})),
		})
	}

	return result
}

func buildTrainingJobSpecSchedulePolicyPreferredAffinityNodeAffinityPreference(preference map[string]interface{}) map[string]interface{} {
	if len(preference) == 0 {
		return nil
	}

	return map[string]interface{}{
		"matchExpressions": buildTrainingJobSchedulePolicyNodeSelectorRequirements(utils.PathSearch("match_expressions",
			preference, make([]interface{}, 0)).([]interface{})),
		"matchFields": buildTrainingJobSchedulePolicyNodeSelectorRequirements(utils.PathSearch("match_fields",
			preference, make([]interface{}, 0)).([]interface{})),
	}
}

func buildTrainingJobSchedulePolicyNodeSelectorRequirements(requirements []interface{}) []map[string]interface{} {
	if len(requirements) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(requirements))
	for _, requirement := range requirements {
		result = append(result, map[string]interface{}{
			"key":      utils.PathSearch("key", requirement, nil),
			"operator": utils.PathSearch("operator", requirement, nil),
			"values": utils.ValueIgnoreEmpty(utils.ExpandToStringList(utils.PathSearch("values",
				requirement, make([]interface{}, 0)).([]interface{}))),
		})
	}

	return result
}

func buildTrainingJobSpecLogExportConfig(logExportConfig map[string]interface{}) map[string]interface{} {
	if len(logExportConfig) < 1 {
		return nil
	}

	return map[string]interface{}{
		"version":          utils.ValueIgnoreEmpty(utils.PathSearch("version", logExportConfig, nil)),
		"rotation_enabled": utils.PathSearch("rotation_enabled", logExportConfig, nil),
	}
}

func buildTrainingJobSpecNotification(notification map[string]interface{}) map[string]interface{} {
	if len(notification) < 1 {
		return nil
	}

	return map[string]interface{}{
		"topic_urn": utils.PathSearch("topic_urn", notification, nil),
		"events": utils.ValueIgnoreEmpty(utils.ExpandToStringList(utils.PathSearch("events",
			notification, make([]interface{}, 0)).([]interface{}))),
	}
}

func buildTrainingJobSpecCustomMetrics(customMetrics []interface{}) []map[string]interface{} {
	if len(customMetrics) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(customMetrics))
	for _, metric := range customMetrics {
		result = append(result, map[string]interface{}{
			"exec": buildTrainingJobSpecCustomMetricsExec(utils.PathSearch("exec[0]",
				metric, make(map[string]interface{}, 0)).(map[string]interface{})),
			"http_get": buildTrainingJobSpecCustomMetricsHttpGet(utils.PathSearch("http_get[0]",
				metric, make(map[string]interface{}, 0)).(map[string]interface{})),
		})
	}

	return result
}

func buildTrainingJobSpecCustomMetricsExec(exec map[string]interface{}) map[string]interface{} {
	if len(exec) < 1 {
		return nil
	}

	return map[string]interface{}{
		"command": utils.ExpandToStringList(utils.PathSearch("command",
			exec, make([]interface{}, 0)).([]interface{})),
	}
}

func buildTrainingJobSpecCustomMetricsHttpGet(httpGet map[string]interface{}) map[string]interface{} {
	if len(httpGet) < 1 {
		return nil
	}

	return map[string]interface{}{
		"path": utils.PathSearch("path", httpGet, nil),
		"port": utils.PathSearch("port", httpGet, nil),
	}
}

func buildTrainingJobSpecOutputModel(outputModel map[string]interface{}) map[string]interface{} {
	if len(outputModel) < 1 {
		return nil
	}

	return map[string]interface{}{
		"obs": buildTrainingJobSpecOutputModelObs(utils.PathSearch("obs[0]",
			outputModel, make(map[string]interface{}, 0)).(map[string]interface{})),
	}
}

func buildTrainingJobSpecOutputModelObs(obs map[string]interface{}) map[string]interface{} {
	if len(obs) < 1 {
		return nil
	}

	return map[string]interface{}{
		"obs_path":   utils.PathSearch("obs_path", obs, nil),
		"local_path": utils.ValueIgnoreEmpty(utils.PathSearch("local_path", obs, nil)),
	}
}

func buildTrainingJobSpecAssetModel(assetModel map[string]interface{}) map[string]interface{} {
	if len(assetModel) < 1 {
		return nil
	}

	return map[string]interface{}{
		"name":    utils.PathSearch("name", assetModel, nil),
		"version": utils.PathSearch("version", assetModel, nil),
		"type":    utils.PathSearch("type", assetModel, nil),
		"code":    utils.ValueIgnoreEmpty(utils.PathSearch("code", assetModel, nil)),
		"desc":    utils.ValueIgnoreEmpty(utils.PathSearch("desc", assetModel, nil)),
		"series":  utils.ValueIgnoreEmpty(utils.PathSearch("series", assetModel, nil)),
	}
}

func createTrainingJob(client *golangsdk.ServiceClient, params map[string]interface{}) (interface{}, error) {
	httpUrl := "v2/{project_id}/training-jobs"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(params),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func resourceTrainingJobCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("modelarts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	resp, err := createTrainingJob(client, buildTrainingJobCreateBodyParams(d))
	if err != nil {
		return diag.Errorf("error creating training job: %s", err)
	}

	trainingJobId := utils.PathSearch("metadata.id", resp, "").(string)
	if trainingJobId == "" {
		return diag.Errorf("unable to find the ID of the training job from the API response")
	}

	d.SetId(trainingJobId)

	stateConf := &retry.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshTrainingJobCreateStatus(client, trainingJobId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for training job (%s) to be created: %s", trainingJobId, err)
	}

	if rawTags, ok := d.GetOk("tags"); ok && len(rawTags.(map[string]interface{})) > 0 {
		err = addTrainingJobTags(client, trainingJobId, rawTags.(map[string]interface{}))
		if err != nil {
			return diag.Errorf("error creating training job tags: %s", err)
		}
	}

	return resourceTrainingJobRead(ctx, d, meta)
}

func refreshTrainingJobCreateStatus(client *golangsdk.ServiceClient, trainingJobId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := GetTrainingJobById(client, trainingJobId)
		if err != nil {
			return nil, "ERROR", err
		}

		status := utils.PathSearch("status.phase", res, "").(string)
		// Pending: The job has been created successfully but is waiting to run.
		// Failed: The job has been created successfully but failed to run.
		// Abnormal: The job has been created successfully but encountered an error during runtime.
		if utils.StrSliceContains([]string{"Running", "Completed", "Pending", "Failed", "Abnormal"}, status) {
			return res, "COMPLETED", nil
		}

		return res, "PENDING", nil
	}
}

func GetTrainingJobById(client *golangsdk.ServiceClient, trainingJobId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/training-jobs/{training_job_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{training_job_id}", trainingJobId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", trainingJobNotFoundErrCodes...)
	}

	return utils.FlattenResponse(requestResp)
}

func resourceTrainingJobRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		trainingJobId = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	resp, err := GetTrainingJobById(client, trainingJobId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving training job (%s)", trainingJobId))
	}

	tags, err := getTrainingJobTags(client, trainingJobId)
	if err != nil {
		log.Printf("[ERROR] error retrieving training job (%s) tags: %s", trainingJobId, err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("kind", utils.PathSearch("kind", resp, nil)),
		d.Set("metadata", flattenTrainingJobMetadata(utils.PathSearch("metadata", resp, nil), d.Get("metadata.0.annotations"))),
		d.Set("algorithm", flattenTrainingJobAlgorithm(utils.PathSearch("algorithm", resp, nil), d.Get("algorithm.0"))),
		d.Set("spec", flattenTrainingJobSpec(utils.PathSearch("spec", resp, nil), d.Get("spec.0"))),
		d.Set("endpoints", flattenTrainingJobEndpoints(utils.PathSearch("endpoints", resp, nil))),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", tags, nil))),
		d.Set("create_time", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("metadata.create_time",
			resp, float64(0)).(float64))/1000, false)),
		d.Set("status", utils.PathSearch("status.phase", resp, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getTrainingJobTags(client *golangsdk.ServiceClient, trainingJobId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/modelarts-training-job/{training_job_id}/tags"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{training_job_id}", trainingJobId)

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

func flattenTrainingJobMetadata(resp interface{}, annotations interface{}) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"name":         utils.PathSearch("name", resp, nil),
			"workspace_id": utils.PathSearch("workspace_id", resp, nil),
			"training_experiment_reference": flattenTrainingJobMetadataTrainingExperimentReference(
				utils.PathSearch("training_experiment_reference.id", resp, "").(string)),
			"description": utils.PathSearch("description", resp, nil),
			"annotations": annotations,
		},
	}
}

func flattenTrainingJobMetadataTrainingExperimentReference(experimentId string) []interface{} {
	if experimentId == "" {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id": experimentId,
		},
	}
}

func flattenTrainingJobAlgorithm(algorithm interface{}, scriptAlgorithm interface{}) []interface{} {
	if algorithm == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":                        utils.PathSearch("id", algorithm, nil),
			"subscription_id":           utils.PathSearch("subscription_id", algorithm, nil),
			"item_version_id":           utils.PathSearch("item_version_id", algorithm, nil),
			"code_dir":                  utils.PathSearch("code_dir", algorithm, nil),
			"boot_file":                 utils.PathSearch("boot_file", algorithm, nil),
			"autosearch_config_path":    utils.PathSearch("autosearch_config_path", algorithm, nil),
			"autosearch_framework_path": utils.PathSearch("autosearch_framework_path", algorithm, nil),
			"command":                   utils.PathSearch("command", algorithm, nil),
			"local_code_dir":            utils.PathSearch("local_code_dir", algorithm, nil),
			"working_dir":               utils.PathSearch("working_dir", algorithm, nil),
			"engine":                    flattenTrainingJobAlgorithmEngine(utils.PathSearch("engine", algorithm, nil)),
			// For fine-tuning job, `inputs` is not returned.
			"inputs": utils.PathSearch("inputs", scriptAlgorithm, nil),
			"outputs": flattenTrainingJobAlgorithmOutputs(utils.PathSearch("outputs",
				algorithm, make([]interface{}, 0)).([]interface{})),
			"parameters": flattenTrainingJobAlgorithmParameters(utils.PathSearch("parameters",
				algorithm, make([]interface{}, 0)).([]interface{})),
			"environments": utils.PathSearch("environments", algorithm, nil),
		},
	}
}

func flattenTrainingJobAlgorithmEngine(engine interface{}) []interface{} {
	if engine == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"engine_id":            utils.PathSearch("engine_id", engine, nil),
			"engine_name":          utils.PathSearch("engine_name", engine, nil),
			"engine_version":       utils.PathSearch("engine_version", engine, nil),
			"image_url":            utils.PathSearch("image_url", engine, nil),
			"install_sys_packages": utils.PathSearch("install_sys_packages", engine, nil),
		},
	}
}

func flattenTrainingJobAlgorithmOutputs(outputs []interface{}) []interface{} {
	if len(outputs) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(outputs))
	for _, output := range outputs {
		result = append(result, map[string]interface{}{
			"name":          utils.PathSearch("name", output, nil),
			"remote":        flattenTrainingJobAlgorithmOutputRemote(utils.PathSearch("remote", output, nil)),
			"local_dir":     utils.PathSearch("local_dir", output, nil),
			"access_method": utils.PathSearch("access_method", output, nil),
			"description":   utils.PathSearch("description", output, nil),
		})
	}

	return result
}

func flattenTrainingJobAlgorithmOutputRemote(remote interface{}) []interface{} {
	if remote == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"obs": flattenTrainingJobAlgorithmOutputObs(utils.PathSearch("obs", remote, nil)),
		},
	}
}

func flattenTrainingJobAlgorithmOutputObs(obs interface{}) []interface{} {
	if obs == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"obs_url": utils.PathSearch("obs_url", obs, nil),
		},
	}
}

func flattenTrainingJobAlgorithmParameters(parameters []interface{}) []interface{} {
	if len(parameters) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(parameters))
	for _, parameter := range parameters {
		result = append(result, map[string]interface{}{
			"name":        utils.PathSearch("name", parameter, nil),
			"value":       utils.PathSearch("value", parameter, nil),
			"description": utils.PathSearch("description", parameter, nil),
			"constraint": flattenTrainingJobAlgorithmParameterConstraint(utils.PathSearch("constraint",
				parameter, nil)),
		})
	}

	return result
}

func flattenTrainingJobAlgorithmParameterConstraint(constraint interface{}) []interface{} {
	if constraint == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"type":        utils.PathSearch("type", constraint, nil),
			"editable":    utils.PathSearch("editable", constraint, nil),
			"required":    utils.PathSearch("required", constraint, nil),
			"valid_type":  utils.PathSearch("valid_type", constraint, nil),
			"valid_range": utils.PathSearch("valid_range", constraint, nil),
		},
	}
}

func flattenTrainingJobEndpoints(endpoints interface{}) []interface{} {
	if endpoints == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"ssh": flattenTrainingJobEndpointsSSH(utils.PathSearch("ssh", endpoints, nil)),
		},
	}
}

func flattenTrainingJobEndpointsSSH(ssh interface{}) []interface{} {
	if ssh == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"key_pair_names": utils.PathSearch("key_pair_names", ssh, nil),
		},
	}
}

func flattenTrainingJobSpec(spec interface{}, scriptSpec interface{}) []interface{} {
	if spec == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"resource":        flattenTrainingJobSpecResource(utils.PathSearch("resource", spec, nil)),
			"runtime_type":    utils.PathSearch("runtime_type", spec, nil),
			"log_export_path": flattenTrainingJobSpecLogExportPath(utils.PathSearch("log_export_path", spec, nil)),
			"log_export_config": flattenTrainingJobSpecLogExportConfig(utils.PathSearch("log_export_config",
				spec, nil)),
			"auto_stop": flattenTrainingJobSpecAutoStop(utils.PathSearch("auto_stop", spec, nil)),
			"schedule_policy": flattenTrainingJobSpecSchedulePolicy(utils.PathSearch("schedule_policy",
				spec, nil)),
			"notification": flattenTrainingJobSpecNotification(utils.PathSearch("notification", spec, nil)),
			"custom_metrics": flattenTrainingJobSpecCustomMetrics(utils.PathSearch("custom_metrics",
				spec, make([]interface{}, 0)).([]interface{})),
			// For fine-tuning job, `output_model`, `asset_model`, `asset_id` are not returned.
			"output_model": utils.PathSearch("output_model", scriptSpec, nil),
			"asset_model":  utils.PathSearch("asset_model", scriptSpec, nil),
			"asset_id":     utils.PathSearch("asset_id", scriptSpec, nil),
			"volumes": flattenTrainingJobSpecVolumes(utils.PathSearch("volumes",
				spec, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenTrainingJobSpecResource(specResource interface{}) []interface{} {
	if specResource == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"node_count":    utils.PathSearch("node_count", specResource, nil),
			"flavor_id":     utils.PathSearch("flavor_id", specResource, nil),
			"pool_id":       utils.PathSearch("pool_id", specResource, nil),
			"pool_group_id": utils.PathSearch("pool_group_id", specResource, nil),
			"main_container_customized_flavor": flattenTrainingJobSpecMainContainerCustomizedFlavor(
				utils.PathSearch("main_container_customized_flavor", specResource, nil)),
		},
	}
}

func flattenTrainingJobSpecMainContainerCustomizedFlavor(mainContainerCustomizedFlavor interface{}) []interface{} {
	if mainContainerCustomizedFlavor == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"cpu_core_num":    utils.PathSearch("cpu_core_num", mainContainerCustomizedFlavor, nil),
			"mem_size":        utils.PathSearch("mem_size", mainContainerCustomizedFlavor, nil),
			"accelerator_num": utils.PathSearch("accelerator_num", mainContainerCustomizedFlavor, nil),
		},
	}
}

func flattenTrainingJobSpecLogExportPath(logExportPath interface{}) []interface{} {
	if logExportPath == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"obs_url":   utils.PathSearch("obs_url", logExportPath, nil),
			"host_path": utils.PathSearch("host_path", logExportPath, nil),
		},
	}
}

func flattenTrainingJobSpecAutoStop(autoStop interface{}) []interface{} {
	if autoStop == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"time_unit": utils.PathSearch("time_unit", autoStop, nil),
			"duration":  utils.PathSearch("duration", autoStop, nil),
		},
	}
}

func flattenTrainingJobSpecSchedulePolicy(schedulePolicy interface{}) []interface{} {
	if schedulePolicy == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"priority":    utils.PathSearch("priority", schedulePolicy, nil),
			"preemptible": utils.PathSearch("preemptible", schedulePolicy, nil),
			"required_affinity": flattenTrainingJobSpecRequiredAffinity(utils.PathSearch("required_affinity",
				schedulePolicy, nil)),
			"preferred_affinity": flattenTrainingJobSpecPreferredAffinity(utils.PathSearch("preferred_affinity",
				schedulePolicy, nil)),
		},
	}
}

func flattenTrainingJobSpecRequiredAffinity(requiredAffinity interface{}) []interface{} {
	if requiredAffinity == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"affinity_type":       utils.PathSearch("affinity_type", requiredAffinity, nil),
			"affinity_group_size": utils.PathSearch("affinity_group_size", requiredAffinity, nil),
			"node_affinity": flattenTrainingJobSpecRequiredAffinityNodeAffinity(utils.PathSearch("node_affinity",
				requiredAffinity, nil)),
		},
	}
}

func flattenTrainingJobSpecRequiredAffinityNodeAffinity(nodeAffinity interface{}) []interface{} {
	if nodeAffinity == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"node_selector_terms": flattenTrainingJobSpecNodeAffinityNodeSelectorTerms(utils.PathSearch("nodeSelectorTerms",
				nodeAffinity, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenTrainingJobSpecNodeAffinityNodeSelectorTerms(nodeSelectorTerms []interface{}) []interface{} {
	if len(nodeSelectorTerms) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(nodeSelectorTerms))
	for _, term := range nodeSelectorTerms {
		result = append(result, map[string]interface{}{
			"match_expressions": flattenTrainingJobNodeSelectorRequirements(utils.PathSearch("matchExpressions",
				term, make([]interface{}, 0)).([]interface{})),
			"match_fields": flattenTrainingJobNodeSelectorRequirements(utils.PathSearch("matchFields",
				term, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenTrainingJobSpecPreferredAffinity(preferredAffinity interface{}) []interface{} {
	if preferredAffinity == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"node_affinity": flattenTrainingJobPreferredSchedulingTerms(utils.PathSearch("node_affinity",
				preferredAffinity, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenTrainingJobPreferredSchedulingTerms(terms []interface{}) []interface{} {
	if len(terms) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(terms))
	for _, term := range terms {
		result = append(result, map[string]interface{}{
			"weight":     utils.PathSearch("weight", term, nil),
			"preference": flattenTrainingJobNodeSelectorTerm(utils.PathSearch("preference", term, nil)),
		})
	}

	return result
}

func flattenTrainingJobNodeSelectorTerm(term interface{}) []interface{} {
	if term == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"match_expressions": flattenTrainingJobNodeSelectorRequirements(utils.PathSearch("matchExpressions",
				term, make([]interface{}, 0)).([]interface{})),
			"match_fields": flattenTrainingJobNodeSelectorRequirements(utils.PathSearch("matchFields",
				term, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenTrainingJobNodeSelectorRequirements(requirements []interface{}) []interface{} {
	if len(requirements) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(requirements))
	for _, requirement := range requirements {
		result = append(result, map[string]interface{}{
			"key":      utils.PathSearch("key", requirement, nil),
			"operator": utils.PathSearch("operator", requirement, nil),
			"values":   utils.PathSearch("values", requirement, nil),
		})
	}

	return result
}

func flattenTrainingJobSpecLogExportConfig(logExportConfig interface{}) []interface{} {
	if logExportConfig == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"version":          utils.PathSearch("version", logExportConfig, nil),
			"rotation_enabled": utils.PathSearch("rotation_enabled", logExportConfig, nil),
		},
	}
}

func flattenTrainingJobSpecNotification(notification interface{}) []interface{} {
	if notification == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"topic_urn": utils.PathSearch("topic_urn", notification, nil),
			"events":    utils.PathSearch("events", notification, nil),
		},
	}
}

func flattenTrainingJobSpecCustomMetrics(customMetrics []interface{}) []interface{} {
	if len(customMetrics) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(customMetrics))
	for _, rawMetric := range customMetrics {
		result = append(result, map[string]interface{}{
			"exec":     flattenTrainingJobSpecCustomMetricsExec(utils.PathSearch("exec", rawMetric, nil)),
			"http_get": flattenTrainingJobSpecCustomMetricsHttpGet(utils.PathSearch("http_get", rawMetric, nil)),
		})
	}

	return result
}

func flattenTrainingJobSpecCustomMetricsExec(exec interface{}) []interface{} {
	if exec == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"command": utils.PathSearch("command", exec, nil),
		},
	}
}

func flattenTrainingJobSpecCustomMetricsHttpGet(httpGet interface{}) []interface{} {
	if httpGet == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"path": utils.PathSearch("path", httpGet, nil),
			"port": utils.PathSearch("port", httpGet, nil),
		},
	}
}

func flattenTrainingJobSpecVolumes(volumes []interface{}) []interface{} {
	if len(volumes) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(volumes))
	for _, volume := range volumes {
		result = append(result, map[string]interface{}{
			"nfs": flattenTrainingJobSpecVolumeNfs(utils.PathSearch("nfs", volume, nil)),
			"pfs": flattenTrainingJobSpecVolumePfs(utils.PathSearch("pfs", volume, nil)),
			"obs": flattenTrainingJobSpecVolumeObs(utils.PathSearch("obs", volume, nil)),
		})
	}

	return result
}

func flattenTrainingJobSpecVolumeNfs(nfs interface{}) []interface{} {
	if nfs == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"nfs_server_path": utils.PathSearch("nfs_server_path", nfs, nil),
			"local_path":      utils.PathSearch("local_path", nfs, nil),
			"read_only":       utils.PathSearch("read_only", nfs, nil),
		},
	}
}

func flattenTrainingJobSpecVolumePfs(pfs interface{}) []interface{} {
	if pfs == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"pfs_path":   utils.PathSearch("pfs_path", pfs, nil),
			"local_path": utils.PathSearch("local_path", pfs, nil),
		},
	}
}

func flattenTrainingJobSpecVolumeObs(obs interface{}) []interface{} {
	if obs == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"obs_path":   utils.PathSearch("obs_path", obs, nil),
			"local_path": utils.PathSearch("local_path", obs, nil),
		},
	}
}

func updateTrainingJob(client *golangsdk.ServiceClient, trainingJobId string, description interface{}) error {
	httpUrl := "v2/{project_id}/training-jobs/{training_job_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{training_job_id}", trainingJobId)

	opt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"description": description,
		},
	}

	_, err := client.Request("PUT", updatePath, &opt)
	return err
}

func resourceTrainingJobUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		trainingJobId = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	if d.HasChange("metadata.0.description") {
		err = updateTrainingJob(client, trainingJobId, d.Get("metadata.0.description"))
		if err != nil {
			return diag.Errorf("error updating training job (%s): %s", trainingJobId, err)
		}
	}

	if d.HasChange("tags") {
		err = updateTrainingJobTags(client, trainingJobId, d)
		if err != nil {
			return diag.Errorf("error updating training job tags: %s", err)
		}
	}

	return resourceTrainingJobRead(ctx, d, meta)
}

func addTrainingJobTags(client *golangsdk.ServiceClient, trainingJobId string, tags map[string]interface{}) error {
	httpUrl := "v2/{project_id}/modelarts-training-job/{training_job_id}/tags/create"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{training_job_id}", trainingJobId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"tags": utils.ExpandResourceTagsMap(tags),
		},
		OkCodes: []int{204},
	}

	_, err := client.Request("POST", createPath, &opt)
	return err
}

func deleteTrainingJobTags(client *golangsdk.ServiceClient, trainingJobId string, tags map[string]interface{}) error {
	httpUrl := "v2/{project_id}/modelarts-training-job/{training_job_id}/tags/delete"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{training_job_id}", trainingJobId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"tags": utils.ExpandResourceTagsMap(tags),
		},
		OkCodes: []int{204},
	}

	_, err := client.Request("DELETE", deletePath, &opt)
	return err
}

func updateTrainingJobTags(client *golangsdk.ServiceClient, trainingJobId string, d *schema.ResourceData) error {
	var (
		oldRaw, newRaw = d.GetChange("tags")
		removeTags     = oldRaw.(map[string]interface{})
		addTags        = newRaw.(map[string]interface{})
	)

	if len(removeTags) > 0 {
		if err := deleteTrainingJobTags(client, trainingJobId, removeTags); err != nil {
			return err
		}
	}

	if len(addTags) > 0 {
		if err := addTrainingJobTags(client, trainingJobId, addTags); err != nil {
			return err
		}
	}

	return nil
}

func deleteTrainingJob(client *golangsdk.ServiceClient, trainingJobId string) error {
	httpUrl := "v2/{project_id}/training-jobs/{training_job_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{training_job_id}", trainingJobId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &opt)
	return err
}

func resourceTrainingJobDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		trainingJobId = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	err = deleteTrainingJob(client, trainingJobId)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", trainingJobNotFoundErrCodes...),
			fmt.Sprintf("error deleting training job (%s)", trainingJobId),
		)
	}

	err = waitingForTrainingJobDeleteCompleted(ctx, client, trainingJobId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for training job (%s) to be deleted: %s", trainingJobId, err)
	}

	return nil
}

func waitingForTrainingJobDeleteCompleted(ctx context.Context, client *golangsdk.ServiceClient, trainingJobId string, t time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			res, err := GetTrainingJobById(client, trainingJobId)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					res = map[string]string{"code": "COMPLETED"}
					return res, "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			return res, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
