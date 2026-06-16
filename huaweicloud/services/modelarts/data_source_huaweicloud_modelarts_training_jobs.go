package modelarts

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts POST /v2/{project_id}/training-job-searches
func DataSourceTrainingJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTrainingJobsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the training jobs are located.`,
			},

			// Optional parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The workspace ID of the training jobs to be queried.`,
			},
			"sort_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The metric used to sort the training jobs.`,
			},
			"order": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The sort order of the training jobs.`,
			},
			"unified_jobs": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to query custom jobs and fine-tuning jobs together.`,
			},
			"train_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The training job type to be queried when unified_jobs is enabled.`,
			},
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 20,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The filter key.`,
						},
						"operator": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The filter operator.`,
						},
						"value": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    10,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The filter values.`,
						},
					},
				},
				Description: `The filter conditions used to query training jobs.`,
			},

			// Attributes.
			"jobs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        trainingJobsSchema(),
				Description: `The list of training jobs that match the filter parameters.`,
			},
		},
	}
}

func trainingJobsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the training job.`,
			},
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the training job.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the training job.`,
						},
						"workspace_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The workspace ID to which the training job belongs.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the training job.`,
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The user name that created the training job.`,
						},
						"annotations": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The advanced feature configuration of the training job.`,
						},
					},
				},
				Description: `The metadata of the training job.`,
			},
			"status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"phase": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The primary status of the training job.`,
						},
						"secondary_phase": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The secondary status of the training job.`,
						},
						"duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The running duration of the training job, in milliseconds.`,
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The start time of the training job, in RFC3339 format.`,
						},
						"tasks": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The subtask names of the training job.`,
						},
						"node_count_metrics": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The node count metrics of the training job, in JSON format.`,
						},
						"task_statuses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"task": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The subtask name.`,
									},
									"exit_code": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The exit code of the subtask.`,
									},
									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The error message of the subtask.`,
									},
								},
							},
							Description: `The status of the first failed subtask.`,
						},
						"running_records": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        trainingJobsRunningRecordSchema(),
							Description: `The running and fault recovery records of the training job.`,
						},
					},
				},
				Description: `The status of the training job.`,
			},
			"algorithm": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        trainingJobsAlgorithmSchema(),
				Description: `The algorithm configuration of the training job.`,
			},
			"spec": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        trainingJobsSpecSchema(),
				Description: `The specification of the training job.`,
			},
			"endpoints": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        trainingJobsEndpointsSchema(),
				Description: `The remote access endpoints of the training job.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the training job, in RFC3339 format.`,
			},
		},
	}
}

func trainingJobsAlgorithmSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The algorithm ID of the training job.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The algorithm name of the training job.`,
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subscription ID of the subscribed algorithm.`,
			},
			"item_version_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version ID of the subscribed algorithm.`,
			},
			"code_dir": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The code directory of the training job.`,
			},
			"boot_file": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The boot file of the training job.`,
			},
			"autosearch_config_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The YAML configuration path of the auto search job.`,
			},
			"autosearch_framework_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The framework code directory of the auto search job.`,
			},
			"command": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The startup command of the custom image training job.`,
			},
			"local_code_dir": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The local code directory in the training container.`,
			},
			"working_dir": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The working directory when running the algorithm.`,
			},
			"parameters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The parameter name.`,
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The parameter value.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The parameter description.`,
						},
						"constraint": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The parameter constraint type.`,
									},
									"editable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether the parameter is editable.`,
									},
									"required": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether the parameter is required.`,
									},
									"valid_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The parameter valid type.`,
									},
									"valid_range": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `The parameter valid values.`,
									},
								},
							},
							Description: `The parameter constraint.`,
						},
					},
				},
				Description: `The runtime parameters of the training job.`,
			},
			"inputs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        trainingJobsAlgorithmInputSchema(),
				Description: `The input channels of the training job.`,
			},
			"outputs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        trainingJobsAlgorithmOutputSchema(),
				Description: `The output channels of the training job.`,
			},
			"engine": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        trainingJobsAlgorithmEngineSchema(),
				Description: `The engine configuration of the training job.`,
			},
			"environments": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The environment variables of the training job.`,
			},
		},
	}
}

func trainingJobsAlgorithmInputSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the input channel.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the input channel.`,
			},
			"local_dir": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The local directory mapped by the input channel.`,
			},
			"access_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The delivery method of the input channel path.`,
			},
			"remote": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dataset": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The dataset ID.`,
									},
									"version_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The dataset version ID.`,
									},
									"obs_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The OBS URL of the dataset.`,
									},
									"service_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The dataset service type.`,
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The dataset name.`,
									},
								},
							},
							Description: `The dataset input information.`,
						},
						"obs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"obs_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The OBS URL of the input data.`,
									},
								},
							},
							Description: `The OBS input information.`,
						},
					},
				},
				Description: `The actual input information.`,
			},
			"remote_constraint": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The data type of the remote constraint.`,
						},
						"attributes": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The attributes of the remote constraint.`,
						},
					},
				},
				Description: `The data input constraint.`,
			},
		},
	}
}

func trainingJobsAlgorithmOutputSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the output channel.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the output channel.`,
			},
			"local_dir": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The local directory mapped by the output channel.`,
			},
			"access_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The delivery method of the output channel path.`,
			},
			"remote": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"obs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"obs_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The OBS URL of the output data.`,
									},
								},
							},
							Description: `The OBS output information.`,
						},
					},
				},
				Description: `The remote output information.`,
			},
		},
	}
}

func trainingJobsAlgorithmEngineSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"engine_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The engine specification ID of the training job.`,
			},
			"engine_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The engine specification name of the training job.`,
			},
			"engine_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The engine specification version of the training job.`,
			},
			"image_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The custom image URL of the training job.`,
			},
			"install_sys_packages": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to install the moxing version specified by the training platform.`,
			},
		},
	}
}

func trainingJobsRecoverRecordSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"recover_start_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The start time of the fault tolerance strategy, in RFC3339 format.`,
			},
			"recover_end_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The end time of the fault tolerance strategy, in RFC3339 format.`,
			},
			"recover": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The fault tolerance strategy.`,
			},
			"fault_scenario": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The fault scenario.`,
			},
			"reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The fault reason.`,
			},
			"related_task": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The task worker ID that triggered the fault.`,
			},
			"recover_result": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The execution result of the fault tolerance strategy.`,
			},
		},
	}
}

func trainingJobsRunningRecordSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"start_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The start time of the run, in RFC3339 format.`,
			},
			"end_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The end time of the run, in RFC3339 format.`,
			},
			"xpu_start_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The accelerator start time of the run, in RFC3339 format.`,
			},
			"start_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The start type of the run.`,
			},
			"end_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The end reason of the run.`,
			},
			"end_related_task": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The task worker ID that caused the run to end.`,
			},
			"end_recover": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The final fault tolerance strategy when the run ended abnormally.`,
			},
			"end_recover_before_downgrade": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The fault tolerance strategy before downgrade when the run ended abnormally.`,
			},
			"recover_records": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        trainingJobsRecoverRecordSchema(),
				Description: `The fault tolerance strategy details when the run ended abnormally.`,
			},
		},
	}
}

func trainingJobsSpecSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        trainingJobsSpecResourceSchema(),
				Description: `The resource specification of the training job.`,
			},
			"runtime_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The runtime type of the training job.`,
			},
			"volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nfs": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        trainingJobsSpecVolumeNfsSchema(),
							Description: `The NFS volume configuration.`,
						},
					},
				},
				Description: `The mounted volumes of the training job.`,
			},
			"log_export_path": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"obs_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The OBS path where the training job logs are exported.`,
						},
					},
				},
				Description: `The log export path of the training job.`,
			},
			"schedule_policy": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        trainingJobsSpecSchedulePolicySchema(),
				Description: `The scheduling policy of the training job.`,
			},
			"custom_metrics": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        trainingJobsSpecCustomMetricsSchema(),
				Description: `The custom metrics configuration of the training job.`,
			},
		},
	}
}

func trainingJobsSpecResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource specification mode of the training job.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The flavor ID of the training job.`,
			},
			"flavor_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The flavor name of the training job.`,
			},
			"node_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of resource replicas selected by the training job.`,
			},
			"pool_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource pool ID selected by the training job.`,
			},
			"pool_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The federated resource pool ID selected by the training job.`,
			},
			"main_container_allocated_resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        trainingJobsSpecResourceMainContainerAllocatedResourcesSchema(),
				Description: `The allocated resources of the main container.`,
			},
			"main_container_customized_flavor": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        trainingJobsSpecResourceMainContainerCustomizedFlavorSchema(),
				Description: `The customized flavor of the main container.`,
			},
		},
	}
}

func trainingJobsSpecResourceMainContainerAllocatedResourcesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cpu_arch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The CPU architecture.`,
			},
			"cpu_core_num": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The number of CPU cores.`,
			},
			"mem_size": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The memory size.`,
			},
			"accelerator_num": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The number of accelerator cards.`,
			},
			"accelerator_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The accelerator type.`,
			},
		},
	}
}

func trainingJobsSpecResourceMainContainerCustomizedFlavorSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cpu_core_num": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The number of CPU cores.`,
			},
			"mem_size": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The memory size.`,
			},
			"accelerator_num": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The number of accelerator cards.`,
			},
		},
	}
}

func trainingJobsSpecVolumeNfsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"nfs_server_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The NFS server path.`,
			},
			"local_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The path for attaching volumes to the training container.`,
			},
			"read_only": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the disks attached in NFS mode are read-only.`,
			},
		},
	}
}

func trainingJobsSpecSchedulePolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"priority": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The scheduling priority of the training job.`,
			},
			"preemptible": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the training job can be preempted.`,
			},
			"required_affinity": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        trainingJobsSpecSchedulePolicyRequiredAffinitySchema(),
				Description: `The required affinity policy of the training job.`,
			},
			"preferred_affinity": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        trainingJobsSpecSchedulePolicyPreferredAffinitySchema(),
				Description: `The preferred affinity configuration of the training job.`,
			},
		},
	}
}

func trainingJobsSpecSchedulePolicyRequiredAffinitySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"affinity_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The affinity scheduling policy type.`,
			},
			"affinity_group_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The affinity group size.`,
			},
			"node_affinity": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        trainingJobsSpecSchedulePolicyRequiredAffinityNodeAffinitySchema(),
				Description: `The node affinity configuration.`,
			},
		},
	}
}

func trainingJobsSpecSchedulePolicyRequiredAffinityNodeAffinitySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_selector_terms": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        trainingJobsSpecSchedulePolicyNodeSelectorTermsSchema(),
				Description: `The node selector term list.`,
			},
		},
	}
}

func trainingJobsSpecSchedulePolicyPreferredAffinitySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_affinity": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        trainingJobsSpecSchedulePolicyPreferredAffinityNodeAffinitySchema(),
				Description: `The preferred node affinity terms.`,
			},
		},
	}
}

func trainingJobsSpecSchedulePolicyPreferredAffinityNodeAffinitySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"weight": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The weight associated with the preferred node selector term.`,
			},
			"preference": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        trainingJobsSpecSchedulePolicyNodeSelectorTermsSchema(),
				Description: `The preferred node selector term.`,
			},
		},
	}
}

func trainingJobsSpecSchedulePolicyNodeSelectorTermsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"match_expressions": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        trainingJobsNodeSelectorRequirementSchema(),
				Description: `The node selector requirements based on node labels.`,
			},
			"match_fields": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        trainingJobsNodeSelectorRequirementSchema(),
				Description: `The node selector requirements based on node fields.`,
			},
		},
	}
}

func trainingJobsNodeSelectorRequirementSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The label key used by the node selector requirement.`,
			},
			"operator": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The operator used by the node selector requirement.`,
			},
			"values": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The label values used by the node selector requirement.`,
			},
		},
	}
}

func trainingJobsSpecCustomMetricsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"exec": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"command": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The command for metrics collection.`,
						},
					},
				},
				Description: `The command-based metrics collection configuration.`,
			},
			"http_get": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The URL path for HTTP metrics collection.`,
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The port for HTTP metrics collection.`,
						},
					},
				},
				Description: `The HTTP-based metrics collection configuration.`,
			},
		},
	}
}

func trainingJobsEndpointsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ssh": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_pair_names": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The SSH key pair names.`,
						},
						"task_urls": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"task": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The task ID.`,
									},
									"url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The SSH connection URL.`,
									},
								},
							},
							Description: `The SSH connection URLs.`,
						},
					},
				},
				Description: `The SSH connection information.`,
			},
			"jupyter_lab": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The JupyterLab URL.`,
						},
						"token": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The JupyterLab token.`,
						},
					},
				},
				Description: `The JupyterLab connection information.`,
			},
			"tensorboard": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The Tensorboard URL.`,
						},
						"token": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The Tensorboard token.`,
						},
					},
				},
				Description: `The Tensorboard connection information.`,
			},
			"mindstudio_insight": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The MindStudio Insight URL.`,
						},
						"token": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The MindStudio Insight token.`,
						},
					},
				},
				Description: `The MindStudio Insight connection information.`,
			},
		},
	}
}

func buildTrainingJobsRequestBody(d *schema.ResourceData, offset, limit int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"offset":       offset,
		"limit":        limit,
		"workspace_id": utils.ValueIgnoreEmpty(d.Get("workspace_id")),
		"sort_by":      utils.ValueIgnoreEmpty(d.Get("sort_by")),
		"order":        utils.ValueIgnoreEmpty(d.Get("order")),
		"unified_jobs": utils.ValueIgnoreEmpty(d.Get("unified_jobs")),
		"train_type":   utils.ValueIgnoreEmpty(d.Get("train_type")),
		"filters":      buildTrainingJobsFiltersBodyParams(d.Get("filters").([]interface{})),
	}

	return bodyParams
}

func buildTrainingJobsFiltersBodyParams(filtersInput []interface{}) []map[string]interface{} {
	if len(filtersInput) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(filtersInput))
	for _, item := range filtersInput {
		result = append(result, map[string]interface{}{
			"key":      utils.ValueIgnoreEmpty(utils.PathSearch("key", item, nil)),
			"operator": utils.ValueIgnoreEmpty(utils.PathSearch("operator", item, nil)),
			"value":    utils.ValueIgnoreEmpty(utils.PathSearch("value", item, nil)),
		})
	}

	return result
}

func listTrainingJobs(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/training-job-searches"
		// Maximum is 50.
		limit    = 50
		pageSize = 0
		result   = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildTrainingJobsRequestBody(d, pageSize, limit)),
	}

	for {
		listOpt.JSONBody.(map[string]interface{})["offset"] = pageSize
		requestResp, err := client.Request("POST", listPath, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		jobs := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, jobs...)

		if len(jobs) < limit {
			break
		}

		pageSize++
	}

	return result, nil
}

func dataSourceTrainingJobsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	jobs, err := listTrainingJobs(client, d)
	if err != nil {
		return diag.Errorf("error querying training jobs: %s", err)
	}

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("jobs", flattenTrainingJobs(jobs)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTrainingJobs(jobs []interface{}) []map[string]interface{} {
	if len(jobs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(jobs))
	for _, job := range jobs {
		result = append(result, map[string]interface{}{
			"kind":      utils.PathSearch("kind", job, nil),
			"metadata":  flattenTrainingJobsMetadata(utils.PathSearch("metadata", job, nil)),
			"status":    flattenTrainingJobsStatus(utils.PathSearch("status", job, nil)),
			"algorithm": flattenTrainingJobsAlgorithm(utils.PathSearch("algorithm", job, nil)),
			"spec":      flattenTrainingJobsSpec(utils.PathSearch("spec", job, nil)),
			"endpoints": flattenTrainingJobsEndpoints(utils.PathSearch("endpoints", job, nil)),
			"create_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("metadata.create_time",
				job, float64(0)).(float64))/1000, false),
		})
	}
	return result
}

func flattenTrainingJobsMetadata(metadata interface{}) []map[string]interface{} {
	if metadata == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"id":           utils.PathSearch("id", metadata, nil),
			"name":         utils.PathSearch("name", metadata, nil),
			"workspace_id": utils.PathSearch("workspace_id", metadata, nil),
			"description":  utils.PathSearch("description", metadata, nil),
			"user_name":    utils.PathSearch("user_name", metadata, nil),
			"annotations":  utils.PathSearch("annotations", metadata, nil),
		},
	}
}

func flattenTrainingJobsStatus(status interface{}) []map[string]interface{} {
	if status == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"phase":           utils.PathSearch("phase", status, nil),
			"secondary_phase": utils.PathSearch("secondary_phase", status, nil),
			"duration":        utils.PathSearch("duration", status, nil),
			"start_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("start_time",
				status, float64(0)).(float64))/1000, false),
			"tasks":              utils.PathSearch("tasks", status, make([]interface{}, 0)).([]interface{}),
			"node_count_metrics": utils.JsonToString(utils.PathSearch("node_count_metrics", status, nil)),
			"task_statuses": flattenTrainingJobsTaskStatuses(utils.PathSearch("task_statuses", status,
				make([]interface{}, 0)).([]interface{})),
			"running_records": flattenTrainingJobsRunningRecords(utils.PathSearch("running_records", status,
				make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenTrainingJobsTaskStatuses(taskStatuses []interface{}) []map[string]interface{} {
	if len(taskStatuses) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(taskStatuses))
	for _, taskStatus := range taskStatuses {
		result = append(result, map[string]interface{}{
			"task":      utils.PathSearch("task", taskStatus, nil),
			"exit_code": utils.PathSearch("exit_code", taskStatus, nil),
			"message":   utils.PathSearch("message", taskStatus, nil),
		})
	}

	return result
}

func flattenTrainingJobsRunningRecords(runningRecords []interface{}) []map[string]interface{} {
	if len(runningRecords) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(runningRecords))
	for _, runningRecord := range runningRecords {
		result = append(result, map[string]interface{}{
			"start_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("start_at",
				runningRecord, float64(0)).(float64)), false),
			"end_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("end_at",
				runningRecord, float64(0)).(float64)), false),
			"xpu_start_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("xpu_start_at",
				runningRecord, float64(0)).(float64)), false),
			"start_type":                   utils.PathSearch("start_type", runningRecord, nil),
			"end_reason":                   utils.PathSearch("end_reason", runningRecord, nil),
			"end_related_task":             utils.PathSearch("end_related_task", runningRecord, nil),
			"end_recover":                  utils.PathSearch("end_recover", runningRecord, nil),
			"end_recover_before_downgrade": utils.PathSearch("end_recover_before_downgrade", runningRecord, nil),
			"recover_records": flattenTrainingJobsRecoverRecords(utils.PathSearch("recover_records", runningRecord,
				make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func flattenTrainingJobsRecoverRecords(records []interface{}) []map[string]interface{} {
	if len(records) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(records))
	for _, record := range records {
		result = append(result, map[string]interface{}{
			"recover_start_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("recover_start_at",
				record, float64(0)).(float64)), false),
			"recover_end_at": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("recover_end_at",
				record, float64(0)).(float64)), false),
			"recover":        utils.PathSearch("recover", record, nil),
			"fault_scenario": utils.PathSearch("fault_scenario", record, nil),
			"reason":         utils.PathSearch("reason", record, nil),
			"related_task":   utils.PathSearch("related_task", record, nil),
			"recover_result": utils.PathSearch("recover_result", record, nil),
		})
	}
	return result
}

func flattenTrainingJobsAlgorithm(algorithm interface{}) []map[string]interface{} {
	if algorithm == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"id":                        utils.PathSearch("id", algorithm, nil),
			"name":                      utils.PathSearch("name", algorithm, nil),
			"subscription_id":           utils.PathSearch("subscription_id", algorithm, nil),
			"item_version_id":           utils.PathSearch("item_version_id", algorithm, nil),
			"code_dir":                  utils.PathSearch("code_dir", algorithm, nil),
			"boot_file":                 utils.PathSearch("boot_file", algorithm, nil),
			"autosearch_config_path":    utils.PathSearch("autosearch_config_path", algorithm, nil),
			"autosearch_framework_path": utils.PathSearch("autosearch_framework_path", algorithm, nil),
			"command":                   utils.PathSearch("command", algorithm, nil),
			"local_code_dir":            utils.PathSearch("local_code_dir", algorithm, nil),
			"working_dir":               utils.PathSearch("working_dir", algorithm, nil),
			"parameters": flattenTrainingJobsAlgorithmParameters(utils.PathSearch("parameters", algorithm,
				make([]interface{}, 0)).([]interface{})),
			"inputs": flattenTrainingJobsAlgorithmInputs(utils.PathSearch("inputs", algorithm,
				make([]interface{}, 0)).([]interface{})),
			"outputs": flattenTrainingJobsAlgorithmOutputs(utils.PathSearch("outputs", algorithm,
				make([]interface{}, 0)).([]interface{})),
			"engine":       flattenTrainingJobsAlgorithmEngine(utils.PathSearch("engine", algorithm, nil)),
			"environments": utils.PathSearch("environments", algorithm, nil),
		},
	}
}

func flattenTrainingJobsAlgorithmParameters(parameters []interface{}) []map[string]interface{} {
	if len(parameters) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(parameters))
	for _, parameter := range parameters {
		result = append(result, map[string]interface{}{
			"name":        utils.PathSearch("name", parameter, nil),
			"value":       utils.PathSearch("value", parameter, nil),
			"description": utils.PathSearch("description", parameter, nil),
			"constraint": flattenTrainingJobsAlgorithmParameterConstraint(utils.PathSearch("constraint",
				parameter, nil)),
		})
	}
	return result
}

func flattenTrainingJobsAlgorithmParameterConstraint(constraint interface{}) []map[string]interface{} {
	if constraint == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"type":        utils.PathSearch("type", constraint, nil),
			"editable":    utils.PathSearch("editable", constraint, nil),
			"required":    utils.PathSearch("required", constraint, nil),
			"valid_type":  utils.PathSearch("valid_type", constraint, nil),
			"valid_range": utils.PathSearch("valid_range", constraint, make([]interface{}, 0)),
		},
	}
}

func flattenTrainingJobsAlgorithmInputs(inputs []interface{}) []map[string]interface{} {
	if len(inputs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(inputs))
	for _, input := range inputs {
		result = append(result, map[string]interface{}{
			"name":          utils.PathSearch("name", input, nil),
			"description":   utils.PathSearch("description", input, nil),
			"local_dir":     utils.PathSearch("local_dir", input, nil),
			"access_method": utils.PathSearch("access_method", input, nil),
			"remote":        flattenTrainingJobsAlgorithmInputRemote(utils.PathSearch("remote", input, nil)),
			"remote_constraint": flattenTrainingJobsAlgorithmInputRemoteConstraints(utils.PathSearch("remote_constraint",
				input, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func flattenTrainingJobsAlgorithmInputRemote(remote interface{}) []map[string]interface{} {
	if remote == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"dataset": flattenTrainingJobsAlgorithmInputRemoteDataset(utils.PathSearch("dataset", remote, nil)),
			"obs":     flattenTrainingJobsAlgorithmInputRemoteObs(utils.PathSearch("obs", remote, nil)),
		},
	}
}

func flattenTrainingJobsAlgorithmInputRemoteDataset(dataset interface{}) []map[string]interface{} {
	if dataset == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"id":           utils.PathSearch("id", dataset, nil),
			"version_id":   utils.PathSearch("version_id", dataset, nil),
			"obs_url":      utils.PathSearch("obs_url", dataset, nil),
			"service_type": utils.PathSearch("service_type", dataset, nil),
			"name":         utils.PathSearch("name", dataset, nil),
		},
	}
}

func flattenTrainingJobsAlgorithmInputRemoteObs(obs interface{}) []map[string]interface{} {
	if obs == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"obs_url": utils.PathSearch("obs_url", obs, nil),
		},
	}
}

func flattenTrainingJobsAlgorithmInputRemoteConstraints(constraints []interface{}) []map[string]interface{} {
	if len(constraints) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(constraints))
	for _, constraint := range constraints {
		result = append(result, map[string]interface{}{
			"data_type":  utils.PathSearch("data_type", constraint, nil),
			"attributes": utils.PathSearch("attributes", constraint, nil),
		})
	}
	return result
}

func flattenTrainingJobsAlgorithmOutputs(outputs []interface{}) []map[string]interface{} {
	if len(outputs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(outputs))
	for _, output := range outputs {
		result = append(result, map[string]interface{}{
			"name":          utils.PathSearch("name", output, nil),
			"description":   utils.PathSearch("description", output, nil),
			"local_dir":     utils.PathSearch("local_dir", output, nil),
			"access_method": utils.PathSearch("access_method", output, nil),
			"remote":        flattenTrainingJobsAlgorithmOutputRemote(utils.PathSearch("remote", output, nil)),
		})
	}
	return result
}

func flattenTrainingJobsAlgorithmOutputRemote(remote interface{}) []map[string]interface{} {
	if remote == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"obs": flattenTrainingJobsAlgorithmOutputRemoteObs(utils.PathSearch("obs", remote, nil)),
		},
	}
}

func flattenTrainingJobsAlgorithmOutputRemoteObs(obs interface{}) []map[string]interface{} {
	if obs == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"obs_url": utils.PathSearch("obs_url", obs, nil),
		},
	}
}

func flattenTrainingJobsAlgorithmEngine(engine interface{}) []map[string]interface{} {
	if engine == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"engine_id":            utils.PathSearch("engine_id", engine, nil),
			"engine_name":          utils.PathSearch("engine_name", engine, nil),
			"engine_version":       utils.PathSearch("engine_version", engine, nil),
			"image_url":            utils.PathSearch("image_url", engine, nil),
			"install_sys_packages": utils.PathSearch("install_sys_packages", engine, nil),
		},
	}
}

func flattenTrainingJobsSpec(spec interface{}) []map[string]interface{} {
	if spec == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"resource":     flattenTrainingJobsSpecResource(utils.PathSearch("resource", spec, nil)),
			"runtime_type": utils.PathSearch("runtime_type", spec, nil),
			"volumes": flattenTrainingJobsSpecVolumes(utils.PathSearch("volumes", spec,
				make([]interface{}, 0)).([]interface{})),
			"log_export_path": flattenTrainingJobsLogExportPath(utils.PathSearch("log_export_path",
				spec, make(map[string]interface{})).(map[string]interface{})),
			"schedule_policy": flattenTrainingJobsSchedulePolicy(utils.PathSearch("schedule_policy", spec, nil)),
			"custom_metrics": flattenTrainingJobsSpecCustomMetrics(utils.PathSearch("custom_metrics", spec,
				make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenTrainingJobsSpecResource(resource interface{}) []map[string]interface{} {
	if resource == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"policy":        utils.PathSearch("policy", resource, nil),
			"flavor_id":     utils.PathSearch("flavor_id", resource, nil),
			"flavor_name":   utils.PathSearch("flavor_name", resource, nil),
			"node_count":    utils.PathSearch("node_count", resource, nil),
			"pool_id":       utils.PathSearch("pool_id", resource, nil),
			"pool_group_id": utils.PathSearch("pool_group_id", resource, nil),
			"main_container_allocated_resources": flattenTrainingJobsSpecResourceMainContainerAllocatedResources(
				utils.PathSearch("main_container_allocated_resources", resource, nil)),
			"main_container_customized_flavor": flattenTrainingJobsSpecResourceMainContainerCustomizedFlavor(
				utils.PathSearch("main_container_customized_flavor", resource, nil)),
		},
	}
}

func flattenTrainingJobsSpecResourceMainContainerAllocatedResources(allocatedResources interface{}) []map[string]interface{} {
	if allocatedResources == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"cpu_arch":         utils.PathSearch("cpu_arch", allocatedResources, nil),
			"cpu_core_num":     utils.PathSearch("cpu_core_num", allocatedResources, nil),
			"mem_size":         utils.PathSearch("mem_size", allocatedResources, nil),
			"accelerator_num":  utils.PathSearch("accelerator_num", allocatedResources, nil),
			"accelerator_type": utils.PathSearch("accelerator_type", allocatedResources, nil),
		},
	}
}

func flattenTrainingJobsSpecResourceMainContainerCustomizedFlavor(customizedFlavor interface{}) []map[string]interface{} {
	if customizedFlavor == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"cpu_core_num":    utils.PathSearch("cpu_core_num", customizedFlavor, nil),
			"mem_size":        utils.PathSearch("mem_size", customizedFlavor, nil),
			"accelerator_num": utils.PathSearch("accelerator_num", customizedFlavor, nil),
		},
	}
}

func flattenTrainingJobsSpecVolumes(volumes []interface{}) []map[string]interface{} {
	if len(volumes) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(volumes))
	for _, volume := range volumes {
		result = append(result, map[string]interface{}{
			"nfs": flattenTrainingJobsSpecVolumeNfs(utils.PathSearch("nfs", volume, nil)),
		})
	}
	return result
}

func flattenTrainingJobsSpecVolumeNfs(nfs interface{}) []map[string]interface{} {
	if nfs == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"nfs_server_path": utils.PathSearch("nfs_server_path", nfs, nil),
			"local_path":      utils.PathSearch("local_path", nfs, nil),
			"read_only":       utils.PathSearch("read_only", nfs, nil),
		},
	}
}

func flattenTrainingJobsLogExportPath(logExportPath map[string]interface{}) []map[string]interface{} {
	if len(logExportPath) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"obs_url": utils.PathSearch("obs_url", logExportPath, nil),
		},
	}
}

func flattenTrainingJobsSchedulePolicy(schedulePolicy interface{}) []map[string]interface{} {
	if schedulePolicy == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"priority":    utils.PathSearch("priority", schedulePolicy, nil),
			"preemptible": utils.PathSearch("preemptible", schedulePolicy, nil),
			"required_affinity": flattenTrainingJobsSpecSchedulePolicyRequiredAffinity(utils.PathSearch("required_affinity",
				schedulePolicy, nil)),
			"preferred_affinity": flattenTrainingJobsSpecSchedulePolicyPreferredAffinity(utils.PathSearch("preferred_affinity",
				schedulePolicy, nil)),
		},
	}
}

func flattenTrainingJobsSpecSchedulePolicyRequiredAffinity(requiredAffinity interface{}) []map[string]interface{} {
	if requiredAffinity == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"affinity_type":       utils.PathSearch("affinity_type", requiredAffinity, nil),
			"affinity_group_size": utils.PathSearch("affinity_group_size", requiredAffinity, nil),
			"node_affinity": flattenTrainingJobsSpecSchedulePolicyRequiredAffinityNodeAffinity(utils.PathSearch("node_affinity",
				requiredAffinity, nil)),
		},
	}
}

func flattenTrainingJobsSpecSchedulePolicyRequiredAffinityNodeAffinity(nodeAffinity interface{}) []map[string]interface{} {
	if nodeAffinity == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"node_selector_terms": flattenTrainingJobsSpecSchedulePolicyNodeSelectorTerms(utils.PathSearch("nodeSelectorTerms",
				nodeAffinity, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenTrainingJobsSpecSchedulePolicyNodeSelectorTerms(nodeSelectorTerms []interface{}) []map[string]interface{} {
	if len(nodeSelectorTerms) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(nodeSelectorTerms))
	for _, term := range nodeSelectorTerms {
		result = append(result, map[string]interface{}{
			"match_expressions": flattenTrainingJobsNodeSelectorRequirements(utils.PathSearch("matchExpressions",
				term, make([]interface{}, 0)).([]interface{})),
			"match_fields": flattenTrainingJobsNodeSelectorRequirements(utils.PathSearch("matchFields",
				term, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func flattenTrainingJobsNodeSelectorRequirements(requirements []interface{}) []map[string]interface{} {
	if len(requirements) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(requirements))
	for _, requirement := range requirements {
		result = append(result, map[string]interface{}{
			"key":      utils.PathSearch("key", requirement, nil),
			"operator": utils.PathSearch("operator", requirement, nil),
			"values":   utils.PathSearch("values", requirement, nil),
		})
	}
	return result
}

func flattenTrainingJobsSpecSchedulePolicyPreferredAffinity(preferredAffinity interface{}) []map[string]interface{} {
	if preferredAffinity == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"node_affinity": flattenTrainingJobsPreferredSchedulingTerms(utils.PathSearch("node_affinity",
				preferredAffinity, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenTrainingJobsPreferredSchedulingTerms(terms []interface{}) []map[string]interface{} {
	if len(terms) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(terms))
	for _, term := range terms {
		result = append(result, map[string]interface{}{
			"weight":     utils.PathSearch("weight", term, nil),
			"preference": flattenTrainingJobsNodeSelectorTerm(utils.PathSearch("preference", term, nil)),
		})
	}
	return result
}

func flattenTrainingJobsNodeSelectorTerm(term interface{}) []map[string]interface{} {
	if term == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"match_expressions": flattenTrainingJobsNodeSelectorRequirements(utils.PathSearch("matchExpressions",
				term, make([]interface{}, 0)).([]interface{})),
			"match_fields": flattenTrainingJobsNodeSelectorRequirements(utils.PathSearch("matchFields",
				term, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenTrainingJobsSpecCustomMetrics(customMetrics []interface{}) []map[string]interface{} {
	if len(customMetrics) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(customMetrics))
	for _, metric := range customMetrics {
		result = append(result, map[string]interface{}{
			"exec":     flattenTrainingJobsSpecCustomMetricsExec(utils.PathSearch("exec", metric, nil)),
			"http_get": flattenTrainingJobsSpecCustomMetricsHttpGet(utils.PathSearch("http_get", metric, nil)),
		})
	}
	return result
}

func flattenTrainingJobsSpecCustomMetricsExec(exec interface{}) []map[string]interface{} {
	if exec == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"command": utils.PathSearch("command", exec, nil),
		},
	}
}

func flattenTrainingJobsSpecCustomMetricsHttpGet(httpGet interface{}) []map[string]interface{} {
	if httpGet == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"path": utils.PathSearch("path", httpGet, nil),
			"port": utils.PathSearch("port", httpGet, nil),
		},
	}
}

func flattenTrainingJobsEndpoints(endpoints interface{}) []map[string]interface{} {
	if endpoints == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"ssh":         flattenTrainingJobsEndpointsSSH(utils.PathSearch("ssh", endpoints, nil)),
			"jupyter_lab": flattenTrainingJobsEndpointsAccess(utils.PathSearch("jupyter_lab", endpoints, nil)),
			"tensorboard": flattenTrainingJobsEndpointsAccess(utils.PathSearch("tensorboard", endpoints, nil)),
			"mindstudio_insight": flattenTrainingJobsEndpointsAccess(utils.PathSearch("mindstudio_insight",
				endpoints, nil)),
		},
	}
}

func flattenTrainingJobsEndpointsSSH(ssh interface{}) []map[string]interface{} {
	if ssh == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"key_pair_names": utils.PathSearch("key_pair_names", ssh, nil),
			"task_urls": flattenTrainingJobsEndpointsSSHTaskURLs(utils.PathSearch("task_urls",
				ssh, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenTrainingJobsEndpointsSSHTaskURLs(taskURLs []interface{}) []map[string]interface{} {
	if len(taskURLs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(taskURLs))
	for _, taskURL := range taskURLs {
		result = append(result, map[string]interface{}{
			"task": utils.PathSearch("task", taskURL, nil),
			"url":  utils.PathSearch("url", taskURL, nil),
		})
	}
	return result
}

func flattenTrainingJobsEndpointsAccess(endpoint interface{}) []map[string]interface{} {
	if endpoint == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"url":   utils.PathSearch("url", endpoint, nil),
			"token": utils.PathSearch("token", endpoint, nil),
		},
	}
}
