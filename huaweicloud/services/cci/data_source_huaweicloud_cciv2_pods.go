package cci

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCI GET /apis/cci/v2/namespaces/{namespace}/pods
func DataSourceV2Pods() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2PodsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pods": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"annotations": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"labels": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"creation_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"active_deadline_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"affinity": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"node_affinity": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     podsNodeAffinitySchema(),
									},
									"pod_anti_affinity": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     podsAntiAffinitySchema(),
									},
								},
							},
						},
						"containers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     podsContainersSchema(),
						},
						"dns_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     podsDNSConfigSchema(),
						},
						"dns_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ephemeral_containers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     podsContainersSchema(),
						},
						"host_aliases": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostnames": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"hostname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_pull_secrets": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"init_containers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     podsContainersSchema(),
						},
						"node_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"overhead": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"readiness_gates": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"condition_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"restart_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scheduler_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_context": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     podsSecurityContextSchema(),
						},
						"set_hostname_as_fqdn": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"share_process_namespace": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"termination_grace_period_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volumes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     podsVolumesSchema(),
						},
						"api_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"finalizers": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"status": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     podsStatusSchema(),
						},
					},
				},
			},
		},
	}
}

func podsStatusSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"observed_generation": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"conditions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_transition_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}

	return &sc
}

func podsVolumesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"projected": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesProjectedSchema(),
			},
			"config_map": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesConfigMapSchema(),
			},
			"nfs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesNfsSchema(),
			},
			"persistent_volume_claim": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesPersistentVolumeClaimSchema(),
			},
			"secret": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesSecretSchema(),
			},
		},
	}

	return &sc
}

func podsVolumesSecretSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"default_mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesKeyToPathSchema(),
			},
			"optional": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"secret_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func podsVolumesKeyToPathSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func podsVolumesPersistentVolumeClaimSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"claim_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"read_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}

	return &sc
}

func podsVolumesNfsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"read_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}

	return &sc
}

func podsVolumesConfigMapSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"default_mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesKeyToPathSchema(),
			},
			"optional": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func podsVolumesProjectedSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"sources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesProjectedSourcesSchema(),
			},
			"default_mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}

	return &sc
}

func podsVolumesProjectedSourcesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"config_map": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsolumesProjectedSourcesConfigMapSchema(),
			},
			"downward_api": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesProjectedSourcesDownwardAPISchema(),
			},
			"secret": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesProjectedSourcesSecretSchema(),
			},
		},
	}

	return &sc
}

func podsVolumesProjectedSourcesSecretSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesKeyToPathSchema(),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"optional": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}

	return &sc
}

func podsVolumesProjectedSourcesDownwardAPISchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsDownwardAPIFileSchema(),
			},
		},
	}

	return &sc
}

func podsDownwardAPIFileSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"field_ref": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"field_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_file_ref": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}

	return &sc
}

func podsolumesProjectedSourcesConfigMapSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesKeyToPathSchema(),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"optional": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}

	return &sc
}

func podsSecurityContextSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"fs_group": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"fs_group_change_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"run_as_group": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"run_as_non_root": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"run_as_user": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"supplemental_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"sysctls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}

	return &sc
}

func podsDNSConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"nameservers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"options": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"searches": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

	return &sc
}

func podsNodeAffinitySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"required_during_scheduling_ignored_during_execution": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsNodeSelectorSchema(),
			},
		},
	}
	return &sc
}

func podsNodeSelectorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_selector_terms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsNodeSelectorTermSchema(),
			},
		},
	}

	return &sc
}

func podsNodeSelectorTermSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"match_expressions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsNodeSelectorRequirementSchema(),
			},
		},
	}

	return &sc
}

func podsNodeSelectorRequirementSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

	return &sc
}

func podsAntiAffinitySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"preferred_during_scheduling_ignored_during_execution": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsWeightedPodAffinityTermSchema(),
			},
			"required_during_scheduling_ignored_during_execution": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsAffinityTermSchema(),
			},
		},
	}
	return &sc
}

func podsAffinityTermSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"label_selector": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsLabelSelectorSchema(),
			},
			"namespaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"topology_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func podsLabelSelectorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"match_labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"match_expressions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsMatchExpressionsSchema(),
			},
		},
	}
	return &sc
}

func podsMatchExpressionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

	return &sc
}

func podsWeightedPodAffinityTermSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"pod_affinity_term": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsAffinityTermSchema(),
			},
			"weight": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &sc
}

func podsContainersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"args": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"command": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"env": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"env_from": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersEnvFromSchema(),
			},
			"image": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lifecycle": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"post_start": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     podsContainersLifecycleHandlerSchema(),
						},
						"pre_stop": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     podsContainersLifecycleHandlerSchema(),
						},
					},
				},
			},
			"liveness_probe": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersProbeSchema(),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"readiness_probe": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersProbeSchema(),
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limits": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"requests": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"security_context": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersSecurityContextSchema(),
			},
			"startup_probe": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersProbeSchema(),
			},
			"stdin": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"stdin_once": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"termination_message_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"termination_message_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tty": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"working_dir": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func podsContainersSecurityContextSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"capabilities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"add": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"drop": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"proc_mount": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"read_only_root_file_system": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"run_as_group": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"run_as_non_root": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"run_as_user": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}

	return &sc
}

func podsContainersProbeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"exec": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersLifecycleHandlerExecSchema(),
			},
			"failure_threshold": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"http_get": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersLifecycleHandlerHttpGetActionSchema(),
			},
			"initial_delay_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"period_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"success_threshold": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"termination_grace_period_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}

	return &sc
}

func podsContainersLifecycleHandlerSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"exec": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersLifecycleHandlerExecSchema(),
			},
			"http_get": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersLifecycleHandlerHttpGetActionSchema(),
			},
		},
	}

	return &sc
}

func podsContainersLifecycleHandlerHttpGetActionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_headers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scheme": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func podsContainersLifecycleHandlerExecSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"command": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

	return &sc
}

func podsContainersEnvFromSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"config_map_ref": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersEnvSourceSchema(),
			},
			"prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_ref": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersEnvSourceSchema(),
			},
		},
	}

	return &sc
}

func podsContainersEnvSourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"optional": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}

	return &sc
}

func dataSourceV2PodsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}
	listPodsHttpUrl := "apis/cci/v2/namespaces/{namespace}/pods"
	listPodsPath := client.Endpoint + listPodsHttpUrl
	listPodsPath = strings.ReplaceAll(listPodsPath, "{namespace}", d.Get("namespace").(string))
	listPodsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listPodsResp, err := client.Request("GET", listPodsPath, &listPodsOpt)
	if err != nil {
		return diag.Errorf("error getting CCI pods list: %s", err)
	}

	listPodsRespBody, err := utils.FlattenResponse(listPodsResp)
	if err != nil {
		return diag.Errorf("error retrieving CCI pods: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	pods := utils.PathSearch("items", listPodsRespBody, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("pods", flattenPods(pods)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPods(pods []interface{}) []interface{} {
	if len(pods) == 0 {
		return nil
	}

	rst := make([]interface{}, len(pods))
	for i, v := range pods {
		containers := utils.PathSearch("spec.containers", v, make([]interface{}, 0)).([]interface{})
		ephemeralContainers := utils.PathSearch("spec.ephemeralContainers", v, make([]interface{}, 0)).([]interface{})
		hostAliases := utils.PathSearch("spec.hostAliases", v, make([]interface{}, 0)).([]interface{})
		imagePullSecrets := utils.PathSearch("spec.imagePullSecrets", v, make([]interface{}, 0)).([]interface{})
		initContainers := utils.PathSearch("spec.initContainers", v, make([]interface{}, 0)).([]interface{})
		readinessGates := utils.PathSearch("spec.readinessGates", v, make([]interface{}, 0)).([]interface{})
		volumes := utils.PathSearch("spec.volumes", v, make([]interface{}, 0)).([]interface{})
		rst[i] = map[string]interface{}{
			"name":                             utils.PathSearch("metadata.name", v, nil),
			"namespace":                        utils.PathSearch("metadata.namespace", v, nil),
			"annotations":                      utils.PathSearch("metadata.annotations", v, nil),
			"labels":                           utils.PathSearch("metadata.labels", v, nil),
			"creation_timestamp":               utils.PathSearch("metadata.creationTimestamp", v, nil),
			"resource_version":                 utils.PathSearch("metadata.resourceVersion", v, nil),
			"uid":                              utils.PathSearch("metadata.uid", v, nil),
			"finalizers":                       utils.PathSearch("metadata.finalizers", v, nil),
			"active_deadline_seconds":          utils.PathSearch("spec.activeDeadlineSeconds", v, nil),
			"affinity":                         flattenPodAffinity(utils.PathSearch("spec.affinity", v, nil)),
			"containers":                       flattenPodContainers(containers),
			"dns_config":                       flattenPodDNSconfig(utils.PathSearch("spec.dnsConfig", v, nil)),
			"dns_policy":                       utils.PathSearch("spec.dnsPolicy", v, nil),
			"ephemeral_containers":             flattenPodContainers(ephemeralContainers),
			"host_aliases":                     flattenPodHostAliases(hostAliases),
			"hostname":                         utils.PathSearch("spec.hostname", v, nil),
			"image_pull_secrets":               flattenImagePullSecrets(imagePullSecrets),
			"node_name":                        utils.PathSearch("spec.nodeName", v, nil),
			"init_containers":                  flattenPodContainers(initContainers),
			"overhead":                         utils.PathSearch("spec.overhead", v, nil),
			"readiness_gates":                  flattenPodReadinessGates(readinessGates),
			"restart_policy":                   utils.PathSearch("spec.restartPolicy", v, nil),
			"scheduler_name":                   utils.PathSearch("spec.schedulerName", v, nil),
			"termination_grace_period_seconds": int(utils.PathSearch("spec.terminationGracePeriodSeconds", v, float64(0)).(float64)),
			"security_context":                 flattenPodseCurityContext(utils.PathSearch("spec.securityContext", v, nil)),
			"volumes":                          flattenPodVolumes(volumes),
			"status":                           flattenPodStatus(utils.PathSearch("status", v, nil)),
		}
	}
	return rst
}
