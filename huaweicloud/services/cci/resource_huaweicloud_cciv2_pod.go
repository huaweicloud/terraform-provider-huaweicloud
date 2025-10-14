package cci

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var podNonUpdatableParams = []string{"namespace", "name"}

// @API CCI POST /apis/cci/v2/namespaces/{namespace}/pods
// @API CCI GET /apis/cci/v2/namespaces/{namespace}/pods/{name}
// @API CCI PUT /apis/cci/v2/namespaces/{namespace}/pods/{name}
// @API CCI DELETE /apis/cci/v2/namespaces/{namespace}/pods/{name}
func ResourceV2Pod() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2PodCreate,
		ReadContext:   resourceV2PodRead,
		UpdateContext: resourceV2PodUpdate,
		DeleteContext: resourceV2PodDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV2PodImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(podNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the namespace of the CCI.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the CCI Pod.`,
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the annotations of the CCI Pod.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The labels of the CCI Pod.`,
			},
			"active_deadline_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The active deadline seconds the pod.`,
			},
			"affinity": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_affinity": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem:     nodeAffinitySchema(),
						},
						"pod_anti_affinity": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem:     podAntiAffinitySchema(),
						},
					},
				},
			},
			"containers": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        podContainersSchema(),
				Description: `Specifies the container of the CCI Pod.`,
			},
			"dns_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        podDNSConfigSchema(),
				Description: `Specifies The DNS config of the pod.`,
			},
			"dns_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the DNS policy of the pod.`,
			},
			"ephemeral_containers": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        podContainersSchema(),
				Description: `Specifies the ephemeral container of the CCI Pod.`,
			},
			"host_aliases": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostnames": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
				Description: `Specifies the host aliases of the CCI Pod.`,
			},
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the host name of the pod.`,
			},
			"image_pull_secrets": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"init_containers": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        podContainersSchema(),
				Description: `Specifies the init container of the CCI Pod.`,
			},
			"node_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the node name of the CCI Pod.`,
			},
			"overhead": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the overhead of the CCI Pod.`,
			},
			"readiness_gates": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"condition_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
				Description: `Specifies the readiness gates of the CCI Pod.`,
			},
			"restart_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The restart policy for all containers within the pod.`,
			},
			"scheduler_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The restart policy for all containers within the pod.`,
			},
			"security_context": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     podSecurityContextSchema(),
			},
			"set_hostname_as_fqdn": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `whether the pod hostname is configured as the pod FQDN.`,
			},
			"share_process_namespace": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to share a single process namespace between all of containers in a pod.`,
			},
			"termination_grace_period_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The restart policy for all containers within the pod.`,
			},
			"volumes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        podVolumesSchema(),
				Description: `Specifies the volumes of the CCI Pod.`,
			},
			"api_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API version of the CCI Pod.`,
			},
			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The kind of the CCI Pod.`,
			},
			"creation_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation timestamp of the CCI Pod.`,
			},
			"resource_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource version of the CCI Pod.`,
			},
			"finalizers": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The finalizers of the namespace.`,
			},
			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uid of the CCI Pod.`,
			},
			"status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"observed_generation": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The observed generation of the CCI Pod.`,
						},
						"conditions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The type of the CCI Pod conditions.`,
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Tthe status of the CCI Pod conditions.`,
									},
									"last_update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The last update time of the CCI Pod conditions.`,
									},
									"last_transition_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The last transition time of the CCI Pod conditions.`,
									},
									"reason": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The reason of the CCI Pod conditions.`,
									},
									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The message of the CCI Pod conditions.`,
									},
								},
							},
							Description: `Tthe conditions of the CCI Pod.`,
						},
					},
				},
				Description: `The status of the CCI Pod.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func podDNSConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"nameservers": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the name servers of the DNS config.`,
			},
			"options": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the name of the options.`,
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the value of the options.`,
						},
					},
				},
				Description: `Specifies the options of the DNS config.`,
			},
			"searches": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the searches of the DNS config.`,
			},
		},
	}

	return &sc
}

func podVolumesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"projected": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     podVolumesProjectedSchema(),
			},
			"config_map": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     podVolumesConfigMapSchema(),
			},
			"nfs": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     podVolumesNfsSchema(),
			},
			"persistent_volume_claim": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     podVolumesPersistentVolumeClaimSchema(),
			},
			"secret": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     podVolumesSecretSchema(),
			},
		},
	}

	return &sc
}

func podVolumesConfigMapSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"default_mode": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"items": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     podVolumesKeyToPathSchema(),
			},
			"optional": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}

	return &sc
}

func podVolumesNfsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"server": {
				Type:     schema.TypeString,
				Required: true,
			},
			"read_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}

	return &sc
}

func podVolumesPersistentVolumeClaimSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"claim_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"read_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}

	return &sc
}

func podVolumesSecretSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"default_mode": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"items": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     podVolumesKeyToPathSchema(),
			},
			"optional": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"secret_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}

	return &sc
}

func podVolumesKeyToPathSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}

	return &sc
}

func podVolumesProjectedSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"sources": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        podVolumesProjectedSourcesSchema(),
				Description: `Specifies the type of the CCI Pod strategy.`,
			},
			"default_mode": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the rolling update config of the CCI Pod strategy.`,
			},
		},
	}

	return &sc
}

func podContainersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"args": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the arguments to the entrypoint of the container.`,
			},
			"command": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the command of the container.`,
			},
			"env": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"env_from": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        podContainersEnvFromSchema(),
				Description: `The sources to populate environment variables of the container.`,
			},
			"image": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the image name of the CCI Pod container.`,
			},
			"lifecycle": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"post_start": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Elem:        podContainersLifecycleHandlerSchema(),
							Description: `The lifecycle post start of the CCI Pod container.`,
						},
						"pre_stop": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Elem:        podContainersLifecycleHandlerSchema(),
							Description: `The lifecycle pre stop of the CCI Pod container.`,
						},
					},
				},
				Description: `Specifies the lifecycle of the container.`,
			},
			"liveness_probe": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        podContainersProbeSchema(),
				Description: `Specifies the liveness probe of the container.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the container.`,
			},
			"ports": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container_port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `Specifies the number of port to expose on the IP address of pod.`,
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the port name of the container.`,
						},
						"protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the protocol for container port.`,
						},
					},
				},
				Description: `Specifies the ports of the container.`,
			},
			"readiness_probe": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        podContainersProbeSchema(),
				Description: `Specifies the readiness probe of the container.`,
			},
			"resources": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limits": {
							Type:        schema.TypeMap,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the limits of resource.`,
						},
						"requests": {
							Type:        schema.TypeMap,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the requests of the resource.`,
						},
					},
				},
				Description: `Specifies the resources of the container.`,
			},
			"security_context": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        podContainersSecurityContextSchema(),
				Description: `Specifies the security context of the container.`,
			},
			"startup_probe": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        podContainersProbeSchema(),
				Description: `Specifies the startup probe of the container.`,
			},
			"stdin": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether this container should allocate a buffer for stdin in the container runtime.`,
			},
			"stdin_once": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether this container runtime should close the stdin channel.`,
			},
			"termination_message_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the termination message path of the CCI Pod container.`,
			},
			"termination_message_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the termination message policy of the CCI Pod container.`,
			},
			"tty": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether this container should allocate a TTY for itself.`,
			},
			"working_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the working directory of the CCI Pod container.`,
			},
			"volume_mounts": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"extend_path_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"mount_path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"read_only": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"sub_path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"sub_path_expr": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}

	return &sc
}

func podContainersProbeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"exec": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        podContainersLifecycleHandlerExecSchema(),
				Description: `Specifies the exec.`,
			},
			"failure_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the minimum consecutive failures for the probe to be considered failed after having succeeded.`,
			},
			"http_get": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MinItems:    1,
				Elem:        podContainersLifecycleHandlerHttpGetActionSchema(),
				Description: `Specifies the HTTP get.`,
			},
			"initial_delay_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The number of seconds after the container has started before liveness probes are initialed.`,
			},
			"period_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies how often to perform the probe.`,
			},
			"success_threshold": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The minimum consecutive successes for the probe to be considered failed after having succeeded.`,
			},
			"termination_grace_period_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}

	return &sc
}

func podContainersSecurityContextSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"capabilities": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"add": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"drop": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				Description: `Specifies the capabilities of the security context.`,
			},
			"proc_mount": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the denotes the type of proc mount to use for the containers.`,
			},
			"read_only_root_file_system": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether this container has a read-only root file system.`,
			},
			"run_as_group": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The GID TO run the entrypoint of the container process.`,
			},
			"run_as_non_root": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `The container must run as a non-root user.`,
			},
			"run_as_user": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The UID to run the entrypoint of the container process.`,
			},
		},
	}

	return &sc
}

func podSecurityContextSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"fs_group": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"fs_group_change_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"run_as_group": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"run_as_non_root": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"run_as_user": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"supplemental_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"sysctls": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}

	return &sc
}

func podContainersLifecycleHandlerSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"exec": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        podContainersLifecycleHandlerExecSchema(),
				Description: `The lifecycle post start of the CCI Pod container.`,
			},
			"http_get": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        podContainersLifecycleHandlerHttpGetActionSchema(),
				Description: `The lifecycle pre stop of the CCI Pod container.`,
			},
		},
	}

	return &sc
}

func podContainersEnvFromSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"config_map_ref": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        podContainersEnvSourceSchema(),
				Description: `Specifies the config map.`,
			},
			"prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the prefix.`,
			},
			"secret_ref": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        podContainersEnvSourceSchema(),
				Description: `Specifies the secret.`,
			},
		},
	}

	return &sc
}

func podContainersEnvSourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the name.`,
			},
			"optional": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to be defined.`,
			},
		},
	}

	return &sc
}

func podContainersLifecycleHandlerExecSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"command": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The command line to execute inside the container.`,
			},
		},
	}

	return &sc
}

func podContainersLifecycleHandlerHttpGetActionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The host name.`,
			},
			"http_headers": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The name of the custom HTTP headers.`,
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The value of the custom HTTP headers.`,
						},
					},
				},
				Description: `The custom headers to set in the request.`,
			},
			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The path to access on the HTTP server.`,
			},
			"port": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The port to access on the HTTP server.`,
			},
			"scheme": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The scheme to use for connecting to the host.`,
			},
		},
	}

	return &sc
}

func podVolumesProjectedSourcesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"config_map": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     podVolumesProjectedSourcesConfigMapSchema(),
			},
			"downward_api": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     podVolumesProjectedSourcesDownwardAPISchema(),
			},
			"secret": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     podVolumesProjectedSourcesSecretSchema(),
			},
		},
	}

	return &sc
}

func podVolumesProjectedSourcesSecretSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"items": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     podVolumesKeyToPathSchema(),
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"optional": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}

	return &sc
}

func podVolumesProjectedSourcesDownwardAPISchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"items": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     downwardAPIFileSchema(),
			},
		},
	}

	return &sc
}

func downwardAPIFileSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"field_ref": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"field_path": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"mode": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_file_ref": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"resource": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}

	return &sc
}

func podVolumesProjectedSourcesConfigMapSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"items": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     podVolumesKeyToPathSchema(),
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"optional": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}

	return &sc
}

func resourceV2PodCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	createNetworkHttpUrl := "apis/cci/v2/namespaces/{namespace}/pods"
	createNetworkPath := client.Endpoint + createNetworkHttpUrl
	createNetworkPath = strings.ReplaceAll(createNetworkPath, "{namespace}", d.Get("namespace").(string))
	createNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createNetworkOpt.JSONBody = utils.RemoveNil(buildCreateV2PodParams(d))

	resp, err := client.Request("POST", createNetworkPath, &createNetworkOpt)
	if err != nil {
		return diag.Errorf("error creating CCI Network: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	ns := utils.PathSearch("metadata.namespace", respBody, "").(string)
	name := utils.PathSearch("metadata.name", respBody, "").(string)
	if ns == "" || name == "" {
		return diag.Errorf("unable to find CCI Pod name or namespace from API response")
	}
	d.SetId(ns + "/" + name)

	err = waitForCreateV2PodStatus(ctx, client, ns, name, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceV2PodRead(ctx, d, meta)
}

func buildCreateV2PodParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"kind":       "Pod",
		"apiVersion": "cci/v2",
		"metadata": map[string]interface{}{
			"name":        d.Get("name"),
			"namespace":   d.Get("namespace"),
			"annotations": utils.ValueIgnoreEmpty(d.Get("annotations")),
			"labels":      utils.ValueIgnoreEmpty(d.Get("labels")),
		},
		"spec": map[string]interface{}{
			"activeDeadlineSeconds":         utils.ValueIgnoreEmpty(d.Get("active_deadline_seconds")),
			"affinity":                      buildV2PodAffinityParams(d.Get("affinity.0")),
			"containers":                    buildV2PodContainersParams(d.Get("containers").(*schema.Set).List()),
			"dnsConfig":                     buildV2PodDNSconfigParams(d.Get("dns_config.0")),
			"dnsPolicy":                     utils.ValueIgnoreEmpty(d.Get("dns_policy")),
			"ephemeralContainers":           buildV2PodContainersParams(d.Get("ephemeral_containers").(*schema.Set).List()),
			"hostAliases":                   buildV2PodHostAliasesParams(d.Get("host_aliases").(*schema.Set).List()),
			"hostname":                      utils.ValueIgnoreEmpty(d.Get("hostname")),
			"imagePullSecrets":              buildImagePullSecretsParams(d.Get("image_pull_secrets").(*schema.Set).List()),
			"initContainers":                buildV2PodContainersParams(d.Get("init_containers").(*schema.Set).List()),
			"nodeName":                      utils.ValueIgnoreEmpty(d.Get("node_name")),
			"overhead":                      utils.ValueIgnoreEmpty(d.Get("overhead")),
			"readinessGates":                buildV2PodReadinessGatesParams(d.Get("readiness_gates").(*schema.Set).List()),
			"restartPolicy":                 utils.ValueIgnoreEmpty(d.Get("restart_policy")),
			"schedulerName":                 utils.ValueIgnoreEmpty(d.Get("scheduler_name")),
			"securityContext":               buildV2PodseCurityContextParams(d.Get("security_context.0")),
			"setHostnameAsPQDN":             utils.ValueIgnoreEmpty(d.Get("set_hostname_as_fqdn")),
			"shareProcessNamespace":         utils.ValueIgnoreEmpty(d.Get("share_process_namespace")),
			"terminationGracePeriodSeconds": utils.ValueIgnoreEmpty(d.Get("termination_grace_period_seconds")),
			"volumes":                       buildV2PodVolumesParams(d.Get("volumes").(*schema.Set).List()),
		},
	}

	return bodyParams
}

func buildV2PodDNSconfigParams(dnsConfig interface{}) map[string]interface{} {
	if dnsConfig == nil {
		return nil
	}
	options := utils.PathSearch("options", dnsConfig, &schema.Set{}).(*schema.Set).List()
	return map[string]interface{}{
		"nameservers": utils.ValueIgnoreEmpty(utils.PathSearch("nameservers", dnsConfig, &schema.Set{}).(*schema.Set).List()),
		"options":     utils.ValueIgnoreEmpty(buildV2PodDNSconfigOptionsParams(options)),
		"searches":    utils.ValueIgnoreEmpty(utils.PathSearch("searches", dnsConfig, &schema.Set{}).(*schema.Set).List()),
	}
}

func buildV2PodDNSconfigOptionsParams(options []interface{}) []interface{} {
	if len(options) == 0 {
		return nil
	}
	params := make([]interface{}, len(options))
	for i, v := range options {
		params[i] = map[string]interface{}{
			"name":  utils.PathSearch("name", v, nil),
			"value": utils.PathSearch("value", v, nil),
		}
	}

	return params
}

func buildV2PodHostAliasesParams(hostAliases []interface{}) []interface{} {
	if len(hostAliases) == 0 {
		return nil
	}
	params := make([]interface{}, len(hostAliases))
	for i, v := range hostAliases {
		params[i] = map[string]interface{}{
			"hostnames": utils.PathSearch("hostnames", v, &schema.Set{}).(*schema.Set).List(),
			"ip":        utils.PathSearch("ip", v, nil),
		}
	}

	return params
}

func buildV2PodReadinessGatesParams(readinessGates []interface{}) []interface{} {
	if len(readinessGates) == 0 {
		return nil
	}
	params := make([]interface{}, len(readinessGates))
	for i, v := range readinessGates {
		params[i] = map[string]interface{}{
			"conditionType": utils.PathSearch("condition_type", v, nil),
		}
	}

	return params
}

func buildV2PodseCurityContextParams(sc interface{}) map[string]interface{} {
	if sc == nil {
		return nil
	}
	sysctls := utils.PathSearch("sysctls", sc, &schema.Set{}).(*schema.Set).List()
	return map[string]interface{}{
		"fsGroup":             utils.PathSearch("fs_group", sc, nil),
		"fsGroupChangePolicy": utils.PathSearch("fs_group_change_policy", sc, nil),
		"runAsGroup":          utils.PathSearch("run_as_group", sc, nil),
		"runAsNonRoot":        utils.PathSearch("run_as_non_root", sc, nil),
		"runAsUser":           utils.PathSearch("run_as_user", sc, nil),
		"supplementalGroups":  utils.PathSearch("supplemental_groups", sc, &schema.Set{}).(*schema.Set).List(),
		"sysctls":             buildV2PodseCurityContextSysctlsParams(sysctls),
	}
}

func buildV2PodseCurityContextSysctlsParams(sysctls []interface{}) []interface{} {
	if len(sysctls) == 0 {
		return nil
	}
	params := make([]interface{}, len(sysctls))
	for i, v := range sysctls {
		params[i] = map[string]interface{}{
			"name":  utils.PathSearch("name", v, nil),
			"value": utils.PathSearch("value", v, nil),
		}
	}

	return params
}

func buildV2PodVolumesParams(volumes []interface{}) []interface{} {
	if len(volumes) == 0 {
		return nil
	}
	params := make([]interface{}, len(volumes))
	for i, v := range volumes {
		params[i] = map[string]interface{}{
			"name":                  utils.PathSearch("name", v, nil),
			"projected":             buildV2PodVolumesProjectedParams(utils.PathSearch("projected|[0]", v, nil)),
			"configMap":             buildV2PodVolumesConfigMapParams(utils.PathSearch("config_map|[0]", v, nil)),
			"nfs":                   buildV2PodVolumesNfsParams(utils.PathSearch("nfs|[0]", v, nil)),
			"persistentVolumeClaim": buildV2PodVolumesPvcParams(utils.PathSearch("persistent_volume_claim|[0]", v, nil)),
			"secret":                buildV2PodVolumesSecretParams(utils.PathSearch("secret|[0]", v, nil)),
		}
	}

	return params
}

func buildV2PodVolumesSecretParams(secret interface{}) map[string]interface{} {
	if secret == nil {
		return nil
	}
	items := utils.PathSearch("items", secret, &schema.Set{}).(*schema.Set).List()
	return map[string]interface{}{
		"defaultMode": utils.PathSearch("default_mode", secret, nil),
		"items":       buildPodVolumesKeyToPathParams(items),
		"optional":    utils.ValueIgnoreEmpty(utils.PathSearch("optional", secret, nil)),
		"secretName":  utils.ValueIgnoreEmpty(utils.PathSearch("secret_name", secret, nil)),
	}
}

func buildV2PodVolumesPvcParams(pvc interface{}) map[string]interface{} {
	if pvc == nil {
		return nil
	}

	return map[string]interface{}{
		"claimName": utils.PathSearch("claim_name", pvc, nil),
		"readOnly":  utils.ValueIgnoreEmpty(utils.PathSearch("read_only", pvc, nil)),
	}
}

func buildV2PodVolumesNfsParams(nfs interface{}) map[string]interface{} {
	if nfs == nil {
		return nil
	}

	return map[string]interface{}{
		"path":     utils.PathSearch("path", nfs, nil),
		"server":   utils.PathSearch("server", nfs, nil),
		"readOnly": utils.ValueIgnoreEmpty(utils.PathSearch("read_only", nfs, nil)),
	}
}

func buildV2PodVolumesConfigMapParams(configMap interface{}) map[string]interface{} {
	if configMap == nil {
		return nil
	}
	items := utils.PathSearch("items", configMap, &schema.Set{}).(*schema.Set).List()
	return map[string]interface{}{
		"defaultMode": utils.PathSearch("default_mode", configMap, nil),
		"items":       buildPodVolumesKeyToPathParams(items),
		"name":        utils.ValueIgnoreEmpty(utils.PathSearch("name", configMap, nil)),
		"optional":    utils.ValueIgnoreEmpty(utils.PathSearch("optional", configMap, nil)),
	}
}

func buildV2PodVolumesProjectedParams(projected interface{}) map[string]interface{} {
	if projected == nil {
		return nil
	}
	sources := utils.PathSearch("sources", projected, &schema.Set{}).(*schema.Set).List()
	return map[string]interface{}{
		"sources":     buildPodVolumesProjectedSourcesParams(sources),
		"defaultMode": utils.PathSearch("default_mode", projected, nil),
	}
}

func buildPodVolumesProjectedSourcesParams(sources []interface{}) []interface{} {
	if len(sources) == 0 {
		return nil
	}
	params := make([]interface{}, len(sources))
	for i, v := range sources {
		params[i] = map[string]interface{}{
			"configMap":   buildV2PodVolumesProjectedSourcesConfigMapParams(utils.PathSearch("config_map|[0]", v, nil)),
			"downwardAPI": buildV2PodVolumesProjectedSourcesDownwardAPIarams(utils.PathSearch("downward_api|[0]", v, nil)),
			"secret":      buildV2PodVolumesProjectedSourcesSecretParams(utils.PathSearch("secret|[0]", v, nil)),
		}
	}

	return params
}

func buildV2PodVolumesProjectedSourcesSecretParams(secret interface{}) map[string]interface{} {
	if secret == nil {
		return nil
	}
	items := utils.PathSearch("items", secret, &schema.Set{}).(*schema.Set).List()
	return map[string]interface{}{
		"items":    buildPodVolumesKeyToPathParams(items),
		"name":     utils.ValueIgnoreEmpty(utils.PathSearch("name", secret, nil)),
		"optional": utils.ValueIgnoreEmpty(utils.PathSearch("optional", secret, nil)),
	}
}

func buildV2PodVolumesProjectedSourcesDownwardAPIarams(downwardAPI interface{}) map[string]interface{} {
	if downwardAPI == nil {
		return nil
	}
	items := utils.PathSearch("items", downwardAPI, &schema.Set{}).(*schema.Set).List()
	return map[string]interface{}{
		"items": buildPodDownwardAPIFileParams(items),
	}
}

func buildPodDownwardAPIFileParams(items []interface{}) []interface{} {
	if len(items) == 0 {
		return nil
	}
	params := make([]interface{}, len(items))
	for i, v := range items {
		params[i] = map[string]interface{}{
			"fieldRef":        buildPodDownwardAPIFileFieldRefParams(utils.PathSearch("field_ref|[0]", v, nil)),
			"mode":            utils.PathSearch("mode", v, nil),
			"path":            utils.PathSearch("path", v, nil),
			"resourceFileRef": buildPodDownwardAPIFileResourceFileRefParams(utils.PathSearch("resource_file_ref|[0]", v, nil)),
		}
	}

	return params
}

func buildPodDownwardAPIFileResourceFileRefParams(resourceFileRef interface{}) map[string]interface{} {
	if resourceFileRef == nil {
		return nil
	}
	return map[string]interface{}{
		"containerName": utils.ValueIgnoreEmpty(utils.PathSearch("container_name", resourceFileRef, nil)),
		"resource":      utils.PathSearch("resource", resourceFileRef, nil),
	}
}

func buildPodDownwardAPIFileFieldRefParams(fieldRef interface{}) map[string]interface{} {
	if fieldRef == nil {
		return nil
	}
	return map[string]interface{}{
		"apiVersion": utils.ValueIgnoreEmpty(utils.PathSearch("api_version", fieldRef, nil)),
		"fieldPath":  utils.PathSearch("field_path", fieldRef, nil),
	}
}

func buildV2PodVolumesProjectedSourcesConfigMapParams(configMap interface{}) map[string]interface{} {
	if configMap == nil {
		return nil
	}
	items := utils.PathSearch("items", configMap, &schema.Set{}).(*schema.Set).List()
	return map[string]interface{}{
		"items":    buildPodVolumesKeyToPathParams(items),
		"name":     utils.ValueIgnoreEmpty(utils.PathSearch("name", configMap, nil)),
		"optional": utils.ValueIgnoreEmpty(utils.PathSearch("optional", configMap, nil)),
	}
}

func buildPodVolumesKeyToPathParams(items []interface{}) []interface{} {
	if len(items) == 0 {
		return nil
	}
	params := make([]interface{}, len(items))
	for i, v := range items {
		params[i] = map[string]interface{}{
			"key":  utils.PathSearch("key", v, nil),
			"mode": utils.PathSearch("mode", v, nil),
			"path": utils.PathSearch("path", v, nil),
		}
	}

	return params
}

func buildV2PodAffinityParams(affinity interface{}) map[string]interface{} {
	if affinity == nil {
		return nil
	}
	return map[string]interface{}{
		"nodeAffinity":    buildNodeAffinityParams(utils.PathSearch("node_affinity|[0]", affinity, nil)),
		"podAntiAffinity": buildPodAntiAffinityParams(utils.PathSearch("pod_anti_affinity|[0]", affinity, nil)),
	}
}

func buildV2PodContainersParams(containers []interface{}) []interface{} {
	if len(containers) == 0 {
		return nil
	}

	containersParams := make([]interface{}, len(containers))
	for i, v := range containers {
		container := utils.RemoveNil(map[string]interface{}{
			"args":                     utils.ValueIgnoreEmpty(utils.PathSearch("args", v, &schema.Set{}).(*schema.Set).List()),
			"command":                  utils.ValueIgnoreEmpty(utils.PathSearch("command", v, &schema.Set{}).(*schema.Set).List()),
			"name":                     utils.ValueIgnoreEmpty(utils.PathSearch("name", v, nil)),
			"image":                    utils.ValueIgnoreEmpty(utils.PathSearch("image", v, nil)),
			"stdin":                    utils.ValueIgnoreEmpty(utils.PathSearch("stdin", v, nil)),
			"stdinOnce":                utils.ValueIgnoreEmpty(utils.PathSearch("stdin_once", v, nil)),
			"terminationMessagePath":   utils.ValueIgnoreEmpty(utils.PathSearch("termination_message_path", v, nil)),
			"terminationMessagePolicy": utils.ValueIgnoreEmpty(utils.PathSearch("termination_message_policy", v, nil)),
			"tty":                      utils.ValueIgnoreEmpty(utils.PathSearch("tty", v, nil)),
			"workingDir":               utils.ValueIgnoreEmpty(utils.PathSearch("working_dir", v, nil)),
			"resources":                buildContainerResourceParams(utils.PathSearch("resources|[0]", v, nil)),
			"lifecycle":                buildPodLifecycleParams(utils.PathSearch("lifecycle|[0]", v, nil)),
			"livenessProbe":            buildPodContainerProbeParams(utils.PathSearch("liveness_probe|[0]", v, nil)),
			"readinessProbe":           buildPodContainerProbeParams(utils.PathSearch("readiness_probe|[0]", v, nil)),
			"startupProbe":             buildPodContainerProbeParams(utils.PathSearch("startup_probe|[0]", v, nil)),
			"securityContext":          buildPodContainersSecurityContextParams(utils.PathSearch("security_context|[0]", v, nil)),
		})

		env := utils.PathSearch("env", v, &schema.Set{}).(*schema.Set).List()
		if len(env) > 0 {
			container["env"] = buildContainerEnvParams(env)
		}

		envFrom := utils.PathSearch("env_from", v, &schema.Set{}).(*schema.Set).List()
		if len(envFrom) > 0 {
			container["envFrom"] = buildPodContainersEnvFromParams(envFrom)
		}

		ports := utils.PathSearch("ports", v, &schema.Set{}).(*schema.Set).List()
		if len(ports) > 0 {
			container["ports"] = buildPodContainersPortsParams(ports)
		}

		volumeMounts := utils.PathSearch("volume_mounts", v, &schema.Set{}).(*schema.Set).List()
		if len(volumeMounts) > 0 {
			container["volumeMounts"] = buildPodContainersVolumeMountsParams(volumeMounts)
		}

		containersParams[i] = container
	}

	return containersParams
}

func buildPodContainersVolumeMountsParams(volumeMounts []interface{}) []interface{} {
	params := make([]interface{}, len(volumeMounts))
	for i, v := range volumeMounts {
		params[i] = utils.RemoveNil(map[string]interface{}{
			"extendPathMode": utils.ValueIgnoreEmpty(utils.PathSearch("extend_path_mode", v, nil)),
			"mountPath":      utils.PathSearch("mount_path", v, nil),
			"name":           utils.PathSearch("name", v, nil),
			"readOnly":       utils.ValueIgnoreEmpty(utils.PathSearch("read_only", v, nil)),
			"subPath":        utils.ValueIgnoreEmpty(utils.PathSearch("sub_path", v, nil)),
			"subPathExpr":    utils.ValueIgnoreEmpty(utils.PathSearch("sub_path_expr", v, nil)),
		})
	}

	return params
}

func buildPodContainersSecurityContextParams(sc interface{}) map[string]interface{} {
	if sc == nil {
		return nil
	}
	return map[string]interface{}{
		"capabilities":           buildPodContainersSecurityContextCapParams(utils.PathSearch("capabilities|[0]", sc, nil)),
		"procMount":              utils.PathSearch("proc_mount", sc, nil),
		"readOnlyRootFileSystem": utils.PathSearch("read_only_root_file_system", sc, nil),
		"runAsGroup":             utils.PathSearch("run_as_group", sc, nil),
		"runAsNonRoot":           utils.PathSearch("run_as_non_root", sc, nil),
		"runAsUser":              utils.PathSearch("run_as_user", sc, nil),
	}
}

func buildPodContainersSecurityContextCapParams(capabilities interface{}) map[string]interface{} {
	if capabilities == nil {
		return nil
	}
	return map[string]interface{}{
		"add":  utils.PathSearch("add", capabilities, &schema.Set{}).(*schema.Set).List(),
		"drop": utils.PathSearch("drop", capabilities, &schema.Set{}).(*schema.Set).List(),
	}
}

func buildPodContainersPortsParams(ports []interface{}) []interface{} {
	params := make([]interface{}, len(ports))
	for i, v := range ports {
		params[i] = map[string]interface{}{
			"containerPort": utils.PathSearch("container_port", v, nil),
			"name":          utils.ValueIgnoreEmpty(utils.PathSearch("name", v, nil)),
			"protocol":      utils.ValueIgnoreEmpty(utils.PathSearch("protocol", v, nil)),
		}
	}

	return params
}

func buildPodContainerProbeParams(probe interface{}) map[string]interface{} {
	if probe == nil {
		return nil
	}
	return map[string]interface{}{
		"failureThreshold":              utils.PathSearch("failure_threshold", probe, nil),
		"initialDelaySeconds":           utils.PathSearch("initial_delay_seconds", probe, nil),
		"periodSeconds":                 utils.PathSearch("period_seconds", probe, nil),
		"successThreshold":              utils.PathSearch("success_threshold", probe, nil),
		"terminationGracePeriodSeconds": utils.PathSearch("termination_grace_period_seconds", probe, nil),
		"exec":                          buildPodContainersLifecycleHandlerExecParams(utils.PathSearch("exec|[0]", probe, nil)),
		"httpGet":                       buildPodContainersLifecycleHandlerHttpGetParams(utils.PathSearch("http_get|[0]", probe, nil)),
	}
}

func buildPodLifecycleParams(lifecycle interface{}) map[string]interface{} {
	if lifecycle == nil {
		return nil
	}
	return map[string]interface{}{
		"postStart": buildPodContainersLifecycleHandlerParams(utils.PathSearch("post_start|[0]", lifecycle, nil)),
		"preStop":   buildPodContainersLifecycleHandlerParams(utils.PathSearch("pre_stop|[0]", lifecycle, nil)),
	}
}

func buildPodContainersLifecycleHandlerParams(lifecycle interface{}) map[string]interface{} {
	if lifecycle == nil {
		return nil
	}
	return map[string]interface{}{
		"exec":    buildPodContainersLifecycleHandlerExecParams(utils.PathSearch("exec|[0]", lifecycle, nil)),
		"httpGet": buildPodContainersLifecycleHandlerHttpGetParams(utils.PathSearch("http_get|[0]", lifecycle, nil)),
	}
}

func buildPodContainersLifecycleHandlerHttpGetParams(httpGet interface{}) map[string]interface{} {
	if httpGet == nil {
		return nil
	}
	httpHeaders := utils.PathSearch("http_headers", httpGet, &schema.Set{}).(*schema.Set).List()
	return map[string]interface{}{
		"host":        utils.PathSearch("host", httpGet, nil),
		"httpHeaders": buildPodContainersLifecycleHandlerHttpHeadersParams(httpHeaders),
		"path":        utils.PathSearch("path", httpGet, nil),
		"port":        utils.PathSearch("port", httpGet, nil),
		"scheme":      utils.PathSearch("scheme", httpGet, nil),
	}
}

func buildPodContainersLifecycleHandlerHttpHeadersParams(httpHeaders []interface{}) []interface{} {
	if len(httpHeaders) == 0 {
		return nil
	}
	params := make([]interface{}, len(httpHeaders))
	for i, v := range httpHeaders {
		params[i] = map[string]interface{}{
			"name":  utils.PathSearch("name", v, nil),
			"value": utils.PathSearch("value", v, nil),
		}
	}

	return params
}

func buildPodContainersLifecycleHandlerExecParams(exec interface{}) map[string]interface{} {
	if exec == nil {
		return nil
	}
	return map[string]interface{}{
		"command": utils.PathSearch("command", exec, &schema.Set{}).(*schema.Set).List(),
	}
}

func buildPodContainersEnvFromParams(envFrom []interface{}) []interface{} {
	params := make([]interface{}, len(envFrom))
	for i, v := range envFrom {
		params[i] = map[string]interface{}{
			"configMapRef": buildPodContainersEnvSourceParams(utils.PathSearch("config_map_ref|[0]", v, nil)),
			"prefix":       utils.PathSearch("prefix", v, nil),
			"secretRef":    buildPodContainersEnvSourceParams(utils.PathSearch("secret_ref|[0]", v, nil)),
		}
	}

	return params
}

func buildPodContainersEnvSourceParams(envSource interface{}) map[string]interface{} {
	if envSource == nil {
		return nil
	}
	return map[string]interface{}{
		"name":     utils.PathSearch("name", envSource, nil),
		"optional": utils.PathSearch("optional", envSource, nil),
	}
}

func resourceV2PodRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)
	resp, err := GetV2Pod(client, ns, name)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying CCI v2 network")
	}

	mErr := multierror.Append(
		d.Set("name", utils.PathSearch("metadata.name", resp, nil)),
		d.Set("namespace", utils.PathSearch("metadata.namespace", resp, nil)),
		d.Set("uid", utils.PathSearch("metadata.uid", resp, nil)),
		d.Set("resource_version", utils.PathSearch("metadata.resourceVersion", resp, nil)),
		d.Set("creation_timestamp", utils.PathSearch("metadata.creationTimestamp", resp, nil)),
		d.Set("kind", utils.PathSearch("kind", resp, nil)),
		d.Set("api_version", utils.PathSearch("apiVersion", resp, nil)),
		d.Set("annotations", utils.PathSearch("metadata.annotations", resp, nil)),
		d.Set("labels", utils.PathSearch("metadata.labels", resp, nil)),
		d.Set("finalizers", utils.PathSearch("metadata.finalizers", resp, nil)),
		d.Set("active_deadline_seconds", utils.PathSearch("spec.activeDeadlineSeconds", resp, nil)),
		d.Set("affinity", flattenPodAffinity(utils.PathSearch("spec.affinity", resp, nil))),
		d.Set("containers", flattenPodContainers(
			utils.PathSearch("spec.containers", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("dns_config", flattenPodDNSconfig(utils.PathSearch("spec.dnsConfig", resp, nil))),
		d.Set("dns_policy", utils.PathSearch("spec.dnsPolicy", resp, nil)),
		d.Set("ephemeral_containers", flattenPodContainers(
			utils.PathSearch("spec.ephemeralContainers", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("host_aliases", flattenPodHostAliases(
			utils.PathSearch("spec.hostAliases", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("hostname", utils.PathSearch("spec.hostname", resp, nil)),
		d.Set("image_pull_secrets", flattenImagePullSecrets(
			utils.PathSearch("spec.imagePullSecrets", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("node_name", utils.PathSearch("spec.nodeName", resp, nil)),
		d.Set("init_containers", flattenPodContainers(
			utils.PathSearch("spec.initContainers", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("overhead", utils.PathSearch("spec.overhead", resp, nil)),
		d.Set("readiness_gates", flattenPodReadinessGates(
			utils.PathSearch("spec.readinessGates", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("restart_policy", utils.PathSearch("spec.restartPolicy", resp, nil)),
		d.Set("scheduler_name", utils.PathSearch("spec.schedulerName", resp, nil)),
		d.Set("termination_grace_period_seconds", int(utils.PathSearch("spec.terminationGracePeriodSeconds", resp, float64(0)).(float64))),
		d.Set("security_context", flattenPodseCurityContext(utils.PathSearch("spec.securityContext", resp, nil))),
		d.Set("volumes", flattenPodVolumes(utils.PathSearch("spec.volumes", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("status", flattenPodStatus(utils.PathSearch("status", resp, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPodVolumes(volumes []interface{}) []interface{} {
	if len(volumes) == 0 {
		return nil
	}
	rst := make([]interface{}, len(volumes))
	for i, v := range volumes {
		rst[i] = map[string]interface{}{
			"name":                    utils.PathSearch("name", v, nil),
			"projected":               flattenPodVolumesProjected(utils.PathSearch("projected", v, nil)),
			"config_map":              flattenPodVolumesConfigMap(utils.PathSearch("configMap", v, nil)),
			"nfs":                     flattenPodVolumesNfs(utils.PathSearch("nfs", v, nil)),
			"persistent_volume_claim": flattenPodVolumesPvc(utils.PathSearch("persistentVolumeClaim", v, nil)),
			"secret":                  flattenPodVolumesSecret(utils.PathSearch("secret", v, nil)),
		}
	}

	return rst
}

func flattenPodVolumesSecret(secret interface{}) []map[string]interface{} {
	if secret == nil || len(secret.(map[string]interface{})) == 0 {
		return nil
	}
	items := utils.PathSearch("items", secret, &schema.Set{}).(*schema.Set).List()
	return []map[string]interface{}{
		{
			"defaultMode": utils.PathSearch("default_mode", secret, nil),
			"items":       flattenPodVolumesKeyToPath(items),
			"optional":    utils.PathSearch("optional", secret, nil),
			"secretName":  utils.PathSearch("secret_name", secret, nil),
		},
	}
}

func flattenPodVolumesPvc(pvc interface{}) []map[string]interface{} {
	if pvc == nil || len(pvc.(map[string]interface{})) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"claim_name": utils.PathSearch("claimName", pvc, nil),
			"read_only":  utils.PathSearch("readOnly", pvc, nil),
		},
	}
}

func flattenPodVolumesConfigMap(configMap interface{}) []map[string]interface{} {
	if configMap == nil || len(configMap.(map[string]interface{})) == 0 {
		return nil
	}
	items := utils.PathSearch("items", configMap, &schema.Set{}).(*schema.Set).List()
	return []map[string]interface{}{
		{
			"default_mode": utils.PathSearch("defaultMode", configMap, nil),
			"items":        flattenPodVolumesKeyToPath(items),
			"name":         utils.PathSearch("name", configMap, nil),
			"optional":     utils.PathSearch("optional", configMap, nil),
		},
	}
}

func flattenPodVolumesNfs(nfs interface{}) []map[string]interface{} {
	if nfs == nil || len(nfs.(map[string]interface{})) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"path":     utils.PathSearch("path", nfs, nil),
			"server":   utils.PathSearch("server", nfs, nil),
			"readOnly": utils.PathSearch("read_only", nfs, nil),
		},
	}
}

func flattenPodVolumesProjected(projected interface{}) []map[string]interface{} {
	if projected == nil || len(projected.(map[string]interface{})) == 0 {
		return nil
	}
	sources := utils.PathSearch("sources", projected, make([]interface{}, 0)).([]interface{})
	return []map[string]interface{}{
		{
			"sources":      flattenPodVolumesProjectedSources(sources),
			"default_mode": utils.PathSearch("defaultMode", projected, nil),
		},
	}
}

func flattenPodVolumesProjectedSources(sources []interface{}) []interface{} {
	if len(sources) == 0 {
		return nil
	}
	rst := make([]interface{}, len(sources))
	for i, v := range sources {
		rst[i] = map[string]interface{}{
			"config_map":   flattenPodVolumesProjectedSourcesConfigMap(utils.PathSearch("configMap", v, nil)),
			"downward_api": flattenPodVolumesProjectedSourcesDownwardAPI(utils.PathSearch("downwardAPI", v, nil)),
			"secret":       flattenPodVolumesProjectedSourcesSecret(utils.PathSearch("secret", v, nil)),
		}
	}

	return rst
}

func flattenPodVolumesProjectedSourcesSecret(secret interface{}) []map[string]interface{} {
	if secret == nil || len(secret.(map[string]interface{})) == 0 {
		return nil
	}
	items := utils.PathSearch("items", secret, make([]interface{}, 0)).([]interface{})
	return []map[string]interface{}{
		{
			"items":    flattenPodVolumesKeyToPath(items),
			"name":     utils.ValueIgnoreEmpty(utils.PathSearch("name", secret, nil)),
			"optional": utils.ValueIgnoreEmpty(utils.PathSearch("optional", secret, nil)),
		},
	}
}

func flattenPodVolumesKeyToPath(items []interface{}) []interface{} {
	if len(items) == 0 {
		return nil
	}
	rst := make([]interface{}, len(items))
	for i, v := range items {
		rst[i] = map[string]interface{}{
			"key":  utils.PathSearch("key", v, nil),
			"mode": utils.PathSearch("mode", v, nil),
			"path": utils.PathSearch("path", v, nil),
		}
	}

	return rst
}

func flattenPodVolumesProjectedSourcesDownwardAPI(downwardAPI interface{}) []map[string]interface{} {
	if downwardAPI == nil || len(downwardAPI.(map[string]interface{})) == 0 {
		return nil
	}
	items := utils.PathSearch("items", downwardAPI, make([]interface{}, 0)).([]interface{})
	return []map[string]interface{}{
		{
			"items": flattenPodDownwardAPIFile(items),
		},
	}
}

func flattenPodDownwardAPIFile(items []interface{}) []interface{} {
	if len(items) == 0 {
		return nil
	}
	rst := make([]interface{}, len(items))
	for i, v := range items {
		rst[i] = map[string]interface{}{
			"field_ref":         flattenPodDownwardAPIFileFieldRef(utils.PathSearch("fieldRef", v, nil)),
			"mode":              utils.PathSearch("mode", v, nil),
			"path":              utils.PathSearch("path", v, nil),
			"resource_file_ref": flattenPodDownwardAPIFileResourceFileRef(utils.PathSearch("resourceFileRef", v, nil)),
		}
	}

	return rst
}

func flattenPodDownwardAPIFileResourceFileRef(resourceFileRef interface{}) []map[string]interface{} {
	if resourceFileRef == nil || len(resourceFileRef.(map[string]interface{})) == 0 {
		return nil
	}
	return []map[string]interface{}{
		{
			"container_name": utils.PathSearch("containerName", resourceFileRef, nil),
			"resource":       utils.PathSearch("resource", resourceFileRef, nil),
		},
	}
}

func flattenPodDownwardAPIFileFieldRef(fieldRef interface{}) []map[string]interface{} {
	if fieldRef == nil || len(fieldRef.(map[string]interface{})) == 0 {
		return nil
	}
	return []map[string]interface{}{
		{
			"api_version": utils.ValueIgnoreEmpty(utils.PathSearch("apiVersion", fieldRef, nil)),
			"field_path":  utils.PathSearch("fieldPath", fieldRef, nil),
		},
	}
}

func flattenPodVolumesProjectedSourcesConfigMap(configMap interface{}) []map[string]interface{} {
	if configMap == nil || len(configMap.(map[string]interface{})) == 0 {
		return nil
	}
	items := utils.PathSearch("items", configMap, make([]interface{}, 0)).([]interface{})
	return []map[string]interface{}{
		{
			"items":    flattenPodVolumesKeyToPath(items),
			"name":     utils.ValueIgnoreEmpty(utils.PathSearch("name", configMap, nil)),
			"optional": utils.ValueIgnoreEmpty(utils.PathSearch("optional", configMap, nil)),
		},
	}
}

func flattenPodDNSconfig(dnsConfig interface{}) []map[string]interface{} {
	if dnsConfig == nil || len(dnsConfig.(map[string]interface{})) == 0 {
		return nil
	}
	options := utils.PathSearch("options", dnsConfig, make([]interface{}, 0)).([]interface{})
	return []map[string]interface{}{
		{
			"nameservers": utils.PathSearch("nameservers", dnsConfig, make([]interface{}, 0)).([]interface{}),
			"options":     flattenPodDNSconfigOptions(options),
			"searches":    utils.PathSearch("searches", dnsConfig, make([]interface{}, 0)).([]interface{}),
		},
	}
}

func flattenPodDNSconfigOptions(options []interface{}) []interface{} {
	if len(options) == 0 {
		return nil
	}
	rst := make([]interface{}, len(options))
	for i, v := range options {
		rst[i] = map[string]interface{}{
			"name":  utils.PathSearch("name", v, nil),
			"value": utils.PathSearch("value", v, nil),
		}
	}

	return rst
}

func flattenPodAffinity(affinity interface{}) []map[string]interface{} {
	if affinity == nil || len(affinity.(map[string]interface{})) == 0 {
		return nil
	}
	return []map[string]interface{}{
		{
			"node_affinity":     flattenNodeAffinity(utils.PathSearch("nodeAffinity", affinity, nil)),
			"pod_anti_affinity": flattenPodAntiAffinity(utils.PathSearch("podAntiAffinity", affinity, nil)),
		},
	}
}

func flattenPodseCurityContext(sc interface{}) []map[string]interface{} {
	if sc == nil || len(sc.(map[string]interface{})) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"fs_group":               utils.PathSearch("fsGroup", sc, nil),
			"fs_group_change_policy": utils.PathSearch("fsGroupChangePolicy", sc, nil),
			"run_as_group":           utils.PathSearch("runAsGroup", sc, nil),
			"run_as_non_root":        utils.PathSearch("runAsNonRoot", sc, nil),
			"run_as_user":            utils.PathSearch("runAsUser", sc, nil),
			"supplemental_groups":    utils.PathSearch("supplementalGroups", sc, &schema.Set{}).(*schema.Set).List(),
			"sysctls":                flattenPodseCurityContextSysctls(utils.PathSearch("sysctls", sc, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenPodseCurityContextSysctls(sysctls []interface{}) []interface{} {
	if len(sysctls) == 0 {
		return nil
	}
	params := make([]interface{}, len(sysctls))
	for i, v := range sysctls {
		params[i] = map[string]interface{}{
			"name":  utils.PathSearch("name", v, nil),
			"value": utils.PathSearch("value", v, nil),
		}
	}

	return params
}

func flattenPodReadinessGates(readinessGates []interface{}) []interface{} {
	if len(readinessGates) == 0 {
		return nil
	}
	params := make([]interface{}, len(readinessGates))
	for i, v := range readinessGates {
		params[i] = map[string]interface{}{
			"condition_type": utils.PathSearch("conditionType", v, nil),
		}
	}

	return params
}

func flattenImagePullSecrets(imagePullSecrets []interface{}) []interface{} {
	if len(imagePullSecrets) == 0 {
		return nil
	}

	rst := make([]interface{}, len(imagePullSecrets))
	for i, v := range imagePullSecrets {
		rst[i] = map[string]interface{}{
			"name": utils.PathSearch("name", v, nil),
		}
	}
	return rst
}

func flattenPodHostAliases(hostAliases []interface{}) []interface{} {
	if len(hostAliases) == 0 {
		return nil
	}
	params := make([]interface{}, len(hostAliases))
	for i, v := range hostAliases {
		params[i] = map[string]interface{}{
			"hostnames": utils.PathSearch("hostnames", v, make([]interface{}, 0)).([]interface{}),
			"ip":        utils.PathSearch("ip", v, nil),
		}
	}

	return params
}

func flattenPodContainers(containers []interface{}) []interface{} {
	if len(containers) == 0 {
		return nil
	}

	rst := make([]interface{}, len(containers))
	for i, v := range containers {
		volumeMounts := utils.PathSearch("volumeMounts", v, make([]interface{}, 0)).([]interface{})
		rst[i] = utils.RemoveNil(map[string]interface{}{
			"args":                       utils.PathSearch("args", v, nil),
			"command":                    utils.PathSearch("command", v, nil),
			"name":                       utils.PathSearch("name", v, nil),
			"image":                      utils.ValueIgnoreEmpty(utils.PathSearch("image", v, nil)),
			"stdin":                      utils.PathSearch("stdin", v, nil),
			"stdin_once":                 utils.PathSearch("stdin_once", v, nil),
			"termination_message_path":   utils.PathSearch("terminationMessagePath", v, nil),
			"termination_message_policy": utils.PathSearch("terminationMessagePolicy", v, nil),
			"tty":                        utils.PathSearch("tty", v, nil),
			"working_dir":                utils.PathSearch("workingDir", v, nil),
			"resources":                  flattenContainerResource(utils.PathSearch("resources", v, nil)),
			"lifecycle":                  flattenPodLifecycle(utils.PathSearch("lifecycle", v, nil)),
			"liveness_probe":             flattenPodContainerProbe(utils.PathSearch("livenessProbe", v, nil)),
			"readiness_probe":            flattenPodContainerProbe(utils.PathSearch("readinessProbe", v, nil)),
			"startup_probe":              flattenPodContainerProbe(utils.PathSearch("startupProbe", v, nil)),
			"security_context":           flattenPodContainersSecurityContext(utils.PathSearch("securityContext", v, nil)),
			"env":                        flattenContainerEnv(utils.PathSearch("env", v, make([]interface{}, 0)).([]interface{})),
			"env_from":                   flattenPodContainersEnvFrom(utils.PathSearch("envFrom", v, make([]interface{}, 0)).([]interface{})),
			"ports":                      flattenPodContainersPorts(utils.PathSearch("ports", v, make([]interface{}, 0)).([]interface{})),
			"volume_mounts":              flattenPodContainersVolumeMounts(volumeMounts),
		})
	}

	return rst
}

func flattenPodContainersVolumeMounts(volumeMounts []interface{}) []interface{} {
	params := make([]interface{}, len(volumeMounts))
	for i, v := range volumeMounts {
		params[i] = map[string]interface{}{
			"extend_path_mode": utils.ValueIgnoreEmpty(utils.PathSearch("extendPathMode", v, nil)),
			"mount_path":       utils.PathSearch("mountPath", v, nil),
			"name":             utils.PathSearch("name", v, nil),
			"read_only":        utils.ValueIgnoreEmpty(utils.PathSearch("readOnly", v, nil)),
			"sub_path":         utils.ValueIgnoreEmpty(utils.PathSearch("subPath", v, nil)),
			"sub_path_expr":    utils.ValueIgnoreEmpty(utils.PathSearch("subPathExpr", v, nil)),
		}
	}

	return params
}

func flattenPodContainersPorts(ports []interface{}) []interface{} {
	params := make([]interface{}, len(ports))
	for i, v := range ports {
		params[i] = map[string]interface{}{
			"container_port": utils.PathSearch("containerPort", v, nil),
			"name":           utils.ValueIgnoreEmpty(utils.PathSearch("name", v, nil)),
			"protocol":       utils.ValueIgnoreEmpty(utils.PathSearch("protocol", v, nil)),
		}
	}

	return params
}

func flattenPodContainersEnvFrom(envFrom []interface{}) []interface{} {
	if len(envFrom) == 0 {
		return nil
	}
	params := make([]interface{}, len(envFrom))
	for i, v := range envFrom {
		params[i] = map[string]interface{}{
			"config_map_ref": flattenPodContainersEnvSource(utils.PathSearch("configMapRef", v, nil)),
			"prefix":         utils.PathSearch("prefix", v, nil),
			"secret_ref":     flattenPodContainersEnvSource(utils.PathSearch("secretRef", v, nil)),
		}
	}

	return params
}

func flattenPodContainersEnvSource(envSource interface{}) []map[string]interface{} {
	if envSource == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"name":     utils.PathSearch("name", envSource, nil),
			"optional": utils.PathSearch("optional", envSource, nil),
		},
	}
}

func flattenPodContainersSecurityContext(sc interface{}) []map[string]interface{} {
	if sc == nil || len(sc.(map[string]interface{})) == 0 {
		return nil
	}
	return []map[string]interface{}{
		{
			"capabilities":               flattenPodContainersSecurityContextCap(utils.PathSearch("capabilities", sc, nil)),
			"proc_mount":                 utils.PathSearch("procMount", sc, nil),
			"read_only_root_file_system": utils.PathSearch("readOnlyRootFileSystem", sc, nil),
			"run_as_group":               utils.PathSearch("runAsGroup", sc, nil),
			"run_as_non_root":            utils.PathSearch("runAsNonRoot", sc, nil),
			"run_as_user":                utils.PathSearch("runAsUser", sc, nil),
		},
	}
}

func flattenPodContainersSecurityContextCap(capabilities interface{}) []map[string]interface{} {
	if capabilities == nil || len(capabilities.(map[string]interface{})) == 0 {
		return nil
	}
	return []map[string]interface{}{
		{
			"add":  utils.PathSearch("add", capabilities, make([]interface{}, 0)).([]interface{}),
			"drop": utils.PathSearch("drop", capabilities, make([]interface{}, 0)).([]interface{}),
		},
	}
}

func flattenPodContainerProbe(probe interface{}) []map[string]interface{} {
	if probe == nil || len(probe.(map[string]interface{})) == 0 {
		return nil
	}
	return []map[string]interface{}{
		{
			"failure_threshold":                utils.PathSearch("failureThreshold", probe, nil),
			"initial_delay_seconds":            utils.PathSearch("initialDelaySeconds", probe, nil),
			"period_seconds":                   utils.PathSearch("periodSeconds", probe, nil),
			"success_threshold":                utils.PathSearch("successThreshold", probe, nil),
			"termination_grace_period_seconds": utils.PathSearch("terminationGracePeriodSeconds", probe, nil),
			"exec":                             flattenPodContainersLifecycleHandlerExec(utils.PathSearch("exec", probe, nil)),
			"http_get": flattenPodContainersLifecycleHandlerHttpGet(
				utils.PathSearch("httpGet", probe, nil)),
		},
	}
}

func flattenPodLifecycle(lifecycle interface{}) []map[string]interface{} {
	if lifecycle == nil || len(lifecycle.(map[string]interface{})) == 0 {
		return nil
	}
	return []map[string]interface{}{
		{
			"post_start": flattenPodContainersLifecycleHandler(utils.PathSearch("postStart", lifecycle, nil)),
			"pre_stop":   flattenPodContainersLifecycleHandler(utils.PathSearch("preStop", lifecycle, nil)),
		},
	}
}

func flattenPodContainersLifecycleHandler(handler interface{}) map[string]interface{} {
	if handler == nil || len(handler.(map[string]interface{})) == 0 {
		return nil
	}
	return map[string]interface{}{
		"exec":     flattenPodContainersLifecycleHandlerExec(utils.PathSearch("exec", handler, nil)),
		"http_get": flattenPodContainersLifecycleHandlerHttpGet(utils.PathSearch("httpGet", handler, make([]interface{}, 0)).([]interface{})),
	}
}

func flattenPodContainersLifecycleHandlerExec(exec interface{}) []map[string]interface{} {
	if exec == nil || len(exec.(map[string]interface{})) == 0 {
		return nil
	}
	return []map[string]interface{}{
		{
			"command": utils.PathSearch("command", exec, make([]interface{}, 0)).([]interface{}),
		},
	}
}

func flattenPodContainersLifecycleHandlerHttpGet(httpGet interface{}) []map[string]interface{} {
	if httpGet == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"host": utils.PathSearch("host", httpGet, nil),
			"http_headers": flattenPodContainersLifecycleHandlerHttpHeaders(
				utils.PathSearch("httpHeaders", httpGet, make([]interface{}, 0)).([]interface{})),
			"path":   utils.PathSearch("path", httpGet, nil),
			"port":   utils.PathSearch("port", httpGet, nil),
			"scheme": utils.PathSearch("scheme", httpGet, nil),
		},
	}
}

func flattenPodContainersLifecycleHandlerHttpHeaders(httpHeaders []interface{}) []interface{} {
	if len(httpHeaders) == 0 {
		return nil
	}
	rst := make([]interface{}, len(httpHeaders))
	for i, v := range httpHeaders {
		rst[i] = map[string]interface{}{
			"name":  utils.PathSearch("name", v, nil),
			"value": utils.PathSearch("value", v, nil),
		}
	}

	return rst
}

func flattenPodStatus(status interface{}) []map[string]interface{} {
	if status == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"observed_generation": utils.PathSearch("observedGeneration", status, nil),
			"conditions":          flattenPodStatusConditions(utils.PathSearch("conditions", status, make([]interface{}, 0)).([]interface{})),
		},
	}

	return rst
}

func flattenPodStatusConditions(conditions []interface{}) []interface{} {
	if len(conditions) == 0 {
		return nil
	}

	rst := make([]interface{}, len(conditions))
	for i, v := range conditions {
		rst[i] = map[string]interface{}{
			"type":                 utils.PathSearch("type", v, nil),
			"status":               utils.PathSearch("status", v, nil),
			"last_update_time":     utils.PathSearch("lastUpdateTime", v, nil),
			"last_transition_time": utils.PathSearch("lastTransitionTime", v, nil),
			"reason":               utils.PathSearch("reason", v, nil),
			"message":              utils.PathSearch("message", v, nil),
		}
	}

	return rst
}

func resourceV2PodUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	updateNetworkHttpUrl := "apis/cci/v2/namespaces/{namespace}/pods//{name}"
	updateNetworkPath := client.Endpoint + updateNetworkHttpUrl
	updateNetworkPath = strings.ReplaceAll(updateNetworkPath, "{namespace}", d.Get("namespace").(string))
	updateNetworkPath = strings.ReplaceAll(updateNetworkPath, "{name}", d.Get("name").(string))
	updateNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updateNetworkOpt.JSONBody = utils.RemoveNil(buildCreateV2PodParams(d))

	_, err = client.Request("PUT", updateNetworkPath, &updateNetworkOpt)
	if err != nil {
		return diag.Errorf("error updating CCI v2 Network: %s", err)
	}
	return resourceV2PodRead(ctx, d, meta)
}

func resourceV2PodDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)
	deleteNetworkHttpUrl := "apis/cci/v2/namespaces/{namespace}/pods/{name}"
	deleteNetworkPath := client.Endpoint + deleteNetworkHttpUrl
	deleteNetworkPath = strings.ReplaceAll(deleteNetworkPath, "{namespace}", ns)
	deleteNetworkPath = strings.ReplaceAll(deleteNetworkPath, "{name}", name)
	deleteNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deleteNetworkPath, &deleteNetworkOpt)
	if err != nil {
		return diag.Errorf("error deleting CCI v2 network: %s", err)
	}

	err = waitForDeleteV2PodStatus(ctx, client, ns, name, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForCreateV2PodStatus(ctx context.Context, client *golangsdk.ServiceClient, ns, name string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending"},
		Target:       []string{"Completed"},
		Refresh:      refreshCreateV2PodStatus(client, ns, name),
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the status of the CCI network to complete: %s", err)
	}
	return nil
}

func refreshCreateV2PodStatus(client *golangsdk.ServiceClient, ns, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := GetV2Pod(client, ns, name)
		if err != nil {
			return nil, "ERROR", err
		}
		status := utils.PathSearch("status.phase", resp, "").(string)
		if status == "Running" {
			return resp, "Completed", nil
		}

		return resp, "Pending", nil
	}
}

func waitForDeleteV2PodStatus(ctx context.Context, client *golangsdk.ServiceClient, ns, name string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending"},
		Target:       []string{"Deleted"},
		Refresh:      refreshDeleteV2PodStatus(client, ns, name),
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the status of the CCI v2 network to complete: %s", err)
	}
	return nil
}

func refreshDeleteV2PodStatus(client *golangsdk.ServiceClient, ns, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := GetV2Pod(client, ns, name)
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[DEBUG] successfully deleted CCI pod: %s", name)
			return "", "Deleted", nil
		}
		return resp, "Pending", nil
	}
}

func resourceV2PodImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<namespace>/<name>', but '%s'", importedId)
	}

	d.Set("namespace", parts[0])
	d.Set("name", parts[1])

	return []*schema.ResourceData{d}, nil
}

func GetV2Pod(client *golangsdk.ServiceClient, namespace, name string) (interface{}, error) {
	getV2PodHttpUrl := "apis/cci/v2/namespaces/{namespace}/pods/{name}"
	getV2PodPath := client.Endpoint + getV2PodHttpUrl
	getV2PodPath = strings.ReplaceAll(getV2PodPath, "{namespace}", namespace)
	getV2PodPath = strings.ReplaceAll(getV2PodPath, "{name}", name)
	getV2PodOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getV2PodResp, err := client.Request("GET", getV2PodPath, &getV2PodOpt)
	if err != nil {
		return getV2PodResp, err
	}

	return utils.FlattenResponse(getV2PodResp)
}
