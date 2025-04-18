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
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the metadata annotations of the CCI Pod.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The labels of the namespace.`,
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
			"volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Specifies the type of the CCI Pod strategy.`,
						},
						"projected": {
							Type:        schema.TypeList,
							Computed:    true,
							MaxItems:    1,
							Elem:        podVolumesProjectedSchema(),
							Description: `Specifies the rolling update config of the CCI Pod strategy.`,
						},
					},
				},
				Description: `Specifies the volumes of the CCI Pod.`,
			},
			"containers": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        podContainersSchema(),
				Description: `Specifies the container of the CCI Pod.`,
			},
			"restart_policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The restart policy for all containers within the pod.`,
			},
			"active_deadline_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The active deadline seconds the pod.`,
			},
			"dns_policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The restart policy for all containers within the pod.`,
			},
			"host_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The restart policy for all containers within the pod.`,
			},
			"node_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The restart policy for all containers within the pod.`,
			},
			"overhead": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The overhead.`,
			},
			"termination_grace_period_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The restart policy for all containers within the pod.`,
			},
			"scheduler_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The restart policy for all containers within the pod.`,
			},
			"set_hostname_as_fqdn": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `whether the pod hostname is configured as the pod FQDN.`,
			},
			"share_process_namespace": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to share a single process namespace between all of containers in a pod.`,
			},
			"dns_config": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The restart policy for all containers within the pod.`,
			},
			// d.Set("security_context", utils.PathSearch("spec.securityContext", resp, nil)),
			// d.Set("image_pull_secrets", flattenSpecStrategy(utils.PathSearch("spec.imagePullSecrets", resp, nil))),
			"progress_deadline_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The progress deadline seconds of the CCI Pod.`,
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

func podVolumesProjectedSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"sources_config_map": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Elem:        podVolumesProjectedSourcesSchema(),
				Description: `Specifies the type of the CCI Pod strategy.`,
			},
			"sources_downward_api": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Elem:        podVolumesProjectedSourcesSchema(),
				Description: `Specifies the type of the CCI Pod strategy.`,
			},
			"sources_secret": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Elem:        podVolumesProjectedSourcesSchema(),
				Description: `Specifies the type of the CCI Pod strategy.`,
			},
			"default_mode": {
				Type:        schema.TypeInt,
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
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies the name of the CCI Pod container.`,
			},
			"image": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the image name of the CCI Pod container.`,
			},
			"env": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the subnet ID of the CCI network.`,
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the subnet ID of the CCI network.`,
						},
					},
				},
				Description: `Specifies the environment variables of the CCI Pod container.`,
			},
			"resources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limits": {
							Type:        schema.TypeMap,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the subnet ID of the CCI network.`,
						},
						"requests": {
							Type:        schema.TypeMap,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the subnet ID of the CCI network.`,
						},
					},
				},
				Description: `Specifies the container of the CCI Pod.`,
			},
			"args": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the arguments to the entrypoint of the CCI Pod container.`,
			},
			"command": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the entrypoint array of the CCI Pod container.`,
			},
			"env_from": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The sources to populate environment variables of the CCI Pod container.`,
			},
			"stdin": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether this container should allocate a buffer for stdin in the container runtime.`,
			},
			"stdin_once": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether this container runtime should close the stdin channel.`,
			},
			"termination_message_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The termination message path of the CCI Pod container.`,
			},
			"termination_message_policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The termination message policy of the CCI Pod container.`,
			},
			"tty": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether this container should allocate a TTY for itself.`,
			},
			"working_dir": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The working directory of the CCI Pod container.`,
			},
			"lifecycle": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"post_start": {
							Type:        schema.TypeList,
							Computed:    true,
							MaxItems:    1,
							Elem:        podContainersLifecycleHandlerSchema(),
							Description: `The lifecycle post start of the CCI Pod container.`,
						},
						"pre_stop": {
							Type:        schema.TypeList,
							Computed:    true,
							MaxItems:    1,
							Elem:        podContainersLifecycleHandlerSchema(),
							Description: `The lifecycle pre stop of the CCI Pod container.`,
						},
					},
				},
				Description: `Specifies the container of the CCI Pod.`,
			},
			"liveness_probe": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Elem:        podContainersProbeSchema(),
				Description: `The liveness probe of the CCI Pod container.`,
			},
			"ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of port to expose on the IP address of pod.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The port name.`,
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The protocol for port.`,
						},
					},
				},
				Description: `The container of the CCI Pod.`,
			},
			"readiness_probe": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Elem:        podContainersProbeSchema(),
				Description: `The readiness probe of the CCI Pod container.`,
			},
			"startup_probe": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Elem:        podContainersProbeSchema(),
				Description: `The startup probe of the CCI Pod container.`,
			},
			"security_context": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capabilities": {
							Type:        schema.TypeList,
							Computed:    true,
							MaxItems:    1,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the subnet ID of the CCI network.`,
						},
						"proc_mount": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The denotes the type of proc mount to use for the containers.`,
						},
						"read_only_root_file_system": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether this container has a read-only root file system.`,
						},
						"run_as_group": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `The GID TO run the entrypoint of the container process.`,
						},
						"run_as_non_root": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `The container must run as a non-root user.`,
						},
						"run_as_user": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `The UID to run the entrypoint of the container process.`,
						},
					},
				},
				Description: `Specifies the container of the CCI Pod.`,
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
				Elem:        podContainersLifecycleHandlerExecSchema(),
				Description: `The lifecycle post start of the CCI Pod container.`,
			},
			"http_get": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        podContainersLifecycleHandlerHttpGetActionSchema(),
				Description: `The lifecycle pre stop of the CCI Pod container.`,
			},
			"failure_threshold": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The minimum consecutive failures for the probe to be considered failed after having succeeded.`,
			},
			"initial_delay_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of seconds after the container has started before liveness probes are initialed.`,
			},
			"period_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `How often to perform the probe.`,
			},
			"success_threshold": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The minimum consecutive successes for the probe to be considered failed after having succeeded.`,
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
				Elem:        podContainersLifecycleHandlerExecSchema(),
				Description: `The lifecycle post start of the CCI Pod container.`,
			},
			"http_get": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        podContainersLifecycleHandlerHttpGetActionSchema(),
				Description: `The lifecycle pre stop of the CCI Pod container.`,
			},
		},
	}

	return &sc
}

func podContainersLifecycleHandlerExecSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"command": {
				Type:        schema.TypeList,
				Optional:    true,
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
				Computed:    true,
				Description: `The host name.`,
			},
			"http_headers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the custom HTTP headers.`,
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The value of the custom HTTP headers.`,
						},
					},
				},
				Description: `The custom headers to set in the request.`,
			},
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The path to access on the HTTP server.`,
			},
			"port": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The port to access on the HTTP server.`,
			},
			"scheme": {
				Type:        schema.TypeString,
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
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Specifies the type of the CCI Pod strategy.`,
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Specifies the rolling update config of the CCI Pod strategy.`,
						},
						"mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Specifies the rolling update config of the CCI Pod strategy.`,
						},
					},
				},
				Description: `Specifies the type of the CCI Pod strategy.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the rolling update config of the CCI Pod strategy.`,
			},
			"optional": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Specifies the rolling update config of the CCI Pod strategy.`,
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

	createNetworkHttpUrl := "apis/yangtse/v2/namespaces/{namespace}/pods"
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
		"kind":       "pod",
		"apiVersion": "cci/v2",
		"metadata": map[string]interface{}{
			"name":      d.Get("name"),
			"namespace": d.Get("namespace"),
		},
		"spec": map[string]interface{}{
			"replicas": d.Get("replicas"),
			"selector": buildV2PodSelectorParams(d.Get("selector")),
			"template": buildV2PodTemplateParams(d),
			// "strategy":     buildV2PodStrategyParams(d.Get("strategy")),
		},
	}

	return bodyParams
}

func buildV2PodTemplateParams(d *schema.ResourceData) map[string]interface{} {
	template := map[string]interface{}{
		"metadata": buildV2PodTemplateMetadataParams(d),
		"spec":     buildV2PodTemplateSpecParams(d),
	}

	return template
}

func buildV2PodTemplateSpecParams(d *schema.ResourceData) map[string]interface{} {
	metadata := map[string]interface{}{
		"containers":       buildV2PodTemplateSpecContainersParams(d),
		"dnsPolicy":        d.Get("template.0.spec_dns_policy"),
		"imagePullSecrets": d.Get("template.0.spec_image_pull_secrets"),
		"affinity":         d.Get("template.0.affinity"),
	}

	return metadata
}

func buildV2PodTemplateSpecContainersParams(d *schema.ResourceData) []interface{} {
	containers := d.Get("template.0.spec_containers").([]interface{})
	if len(containers) == 0 {
		return nil
	}
	containersParams := make([]interface{}, 0, len(containers))
	for i, v := range containers {
		containersParams[i] = map[string]interface{}{
			"name":      utils.PathSearch("name", v, nil),
			"image":     utils.PathSearch("image", v, nil),
			"env":       utils.PathSearch("env", v, nil),
			"resources": utils.PathSearch("resources", v, nil),
		}
	}

	return containersParams
}

func buildV2PodTemplateMetadataParams(d *schema.ResourceData) map[string]interface{} {
	metadata := map[string]interface{}{
		"labels":      d.Get("template.0.metadata_labels"),
		"annotations": d.Get("template.0.metadata_annotations"),
	}

	return metadata
}

func buildV2PodSelectorParams(selector interface{}) interface{} {
	if selector == nil {
		return nil
	}

	params := map[string]interface{}{
		"matchLabels": utils.PathSearch("[0].match_labels", selector, nil),
	}
	return params
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
		d.Set("volumes", flattenSpecSelector(utils.PathSearch("spec.volumes", resp, nil))),
		// d.Set("containers", flattenSpecTemplate(utils.PathSearch("spec.containers", resp, nil))),
		d.Set("restart_policy", utils.PathSearch("spec.restartPolicy", resp, nil)),
		d.Set("termination_grace_period_seconds", int(utils.PathSearch("spec.terminationGracePeriodSeconds", resp, float64(0)).(float64))),
		d.Set("dns_policy", utils.PathSearch("spec.dnsPolicy", resp, nil)),
		d.Set("node_name", utils.PathSearch("spec.nodeName", resp, nil)),
		d.Set("security_context", utils.PathSearch("spec.securityContext", resp, nil)),
		d.Set("image_pull_secrets", utils.PathSearch("spec.imagePullSecrets", resp, nil)),
		d.Set("host_name", utils.PathSearch("spec.hostname", resp, nil)),
		// d.Set("scheduler_name", utils.PathSearch("spec.schedulerName", resp, nil)),
		d.Set("dns_config", utils.PathSearch("spec.dnsConfig", resp, nil)),
		d.Set("status", flattenPodStatus(utils.PathSearch("status", resp, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSpecSelector(selector interface{}) []map[string]interface{} {
	if selector == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"match_labels": utils.PathSearch("matchLabels", selector, nil),
		},
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

	updateNetworkHttpUrl := "apis/yangtse/v2/namespaces/{namespace}/pods/{name}"
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
	deleteNetworkHttpUrl := "apis/yangtse/v2/namespaces/{namespace}/pods/{name}"
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
		Target:       []string{"Ready"},
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
		status := utils.PathSearch("status.status", resp, "").(string)
		if status != "Ready" {
			return resp, "Pending", nil
		}

		return resp, status, nil
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
			log.Printf("[DEBUG] successfully deleted CCI network: %s", name)
			return "", "Deleted", nil
		}
		return resp, "Pending", nil
	}
}

func resourceV2PodImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
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
