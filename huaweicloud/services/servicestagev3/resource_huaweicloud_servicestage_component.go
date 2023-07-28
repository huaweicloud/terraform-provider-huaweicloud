package servicestagev3

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/servicestage/v3/jobs"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/servicestage/v3/components"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func ResourceComponent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComponentCreate,
		ReadContext:   resourceComponentRead,
		UpdateContext: resourceComponentUpdate,
		DeleteContext: resourceComponentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceComponentImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[a-z]([a-z0-9-]*[a-z0-9])?$`),
						"The name must start with a lowercase letter and end with a lowercase letter or digit, and "+
							"can only contain lowercase letters, digits and hyphens (-)."),
					validation.StringLenBetween(2, 64),
				),
			},
			"workload_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workload_kind": {
				Type: schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"deployment", "statefulset",
				}, false),
			},
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"replica": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"limit_cpu": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"limit_memory": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"request_cpu": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"request_memory": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"enable_sermant_injection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"timezone": {
				Type: schema.TypeString,
				Optional: true,
			},
			"jvm_opts": {
				Type: schema.TypeString,
				Optional: true,
			},
			"runtime_stack": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"deploy_mode": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"labels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
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
			"pod_labels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
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
			"source": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"GitHub", "GitLab", "Gitee", "Bitbucket", "package", "DevCloud",
							}, false),
						},
						"url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"storage": {
							Type:     schema.TypeString,
							Required: true,
						},
						"auth": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"repo_auth": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"repo_namespace": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"repo_ref": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"web_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"repo_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"build": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameters": {
							Type:     schema.TypeSet,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"artifact_namespace": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"build_cmd": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"cluster_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"dockerfile_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"environment_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"node_label_selector": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"envs": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.All(
								validation.StringMatch(regexp.MustCompile(`^[A-Za-z-_.]([\w-.]*)?$`),
									"The name can only contain letters, digits, underscores (_), "+
										"hyphens (-) and dots (.), and cannot start with a digit."),
								validation.StringLenBetween(1, 64),
							),
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value_from": {
							Type: schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"reference_type": {
										Type:     schema.TypeString,
										Optional: true,
										// ValidateFunc: validation.StringInSlice([]string{
										// 	"configMapKey", "secretKey",
										// }, false),
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"optional": {
										Type: schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"storage": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"HostPath", "EmptyDir", "ConfigMap", "Secret", "PersistentVolumeClaim",
							}, false),
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"parameters": {
							Type:     schema.TypeSet,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"default_mode": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"medium": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"mounts": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:     schema.TypeString,
										Required: true,
									},
									"read_only": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"sub_path": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"deploy_strategy": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"OneBatchRelease", "RollingRelease", "GrayRelease",
							}, false),
						},
						"rolling_release": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"batches": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"termination_seconds": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"fail_strategy": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"continue", "stop",
										}, false),
									},
								},
							},
						},
						"gray_release": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"weight", "content",
										}, false),
									},
									"first_batch_weight": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"first_batch_replica": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"remaining_batch": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"deployment_mode": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"replica_surge_mode": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"mirror", "mirror", "no_surge",
										}, false),
									},
									"rule_match_mode": {
										Type:     schema.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											"all", "any",
										}, false),
									},
									"rules": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Optional: true,
													ValidateFunc: validation.StringInSlice([]string{
														"header", "query_param", "custom", "method", "cookie",
													}, false),
												},
												"key": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"value": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"condition": {
													Type:     schema.TypeString,
													Optional: true,
													ValidateFunc: validation.StringInSlice([]string{
														"equal", "equal", "in",
													}, false),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"command": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"command": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"args": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"post_start": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"http", "command",
							}, false),
						},
						"scheme": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"HTTP", "HTTPS",
							}, false),
						},
						"hosts": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"command": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"pre_stop": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"http", "command",
							}, false),
						},
						"scheme": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"HTTP", "HTTPS",
							}, false),
						},
						"hosts": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"path": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"command": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"mesher": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"tomcat_opts": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_xml": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"host_aliases": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type: schema.TypeString,
							Optional: true,
						},
						"hostnames": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"dns_policy": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Default", "ClusterFirst", "ClusterFirstWithHostNet", "None",
				}, false),
			},
			"dns_config": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nameservers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"searches": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"options": {
							Type: schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type: schema.TypeString,
										Optional: true,
									},
									"value": {
										Type: schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"security_context": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"run_as_user": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"run_as_group": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"capabilities": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"add": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"drop": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"logs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"rotate": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"host_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"host_extend_path": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"None", "PodUID", "PodName", "PodUID/ContainerName", "PodName/ContainerName",
							}, false),
						},
					},
				},
			},
			"custom_metric": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"dimensions": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"affinity": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"node", "pod",
							}, false),
						},
						"condition": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"required", "preferred",
							}, false),
						},
						"weight": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"match_expressions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type: schema.TypeString,
										Optional: true,
									},
									"value": {
										Type: schema.TypeString,
										Optional: true,
									},
									"operation": {
										Type: schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"anti_affinity": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"node", "pod",
							}, false),
						},
						"condition": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"required", "preferred",
							}, false),
						},
						"weight": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"match_expressions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type: schema.TypeString,
										Optional: true,
									},
									"value": {
										Type: schema.TypeString,
										Optional: true,
									},
									"operation": {
										Type: schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"liveness_probe": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"http", "tcp", "command",
							}, false),
						},
						"delay": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"scheme": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"HTTP", "HTTPS",
							}, false),
						},
						"period_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"success_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"failure_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"host": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"command": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"readiness_probe": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"http", "tcp", "command",
							}, false),
						},
						"delay": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"scheme": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"HTTP", "HTTPS",
							}, false),
						},
						"period_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"success_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"failure_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"host": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"command": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"refer_resources": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"parameters": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem:     &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type: schema.TypeString,
										Optional: true,
									},
									"namespace": {
										Type: schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildRepoBuildStructure(build interface{}) *components.Build {
	buildSet := build.(*schema.Set)
	if buildSet.Len() != 1 {
		return &components.Build{}
	}

	buildMap := buildSet.List()[0].(map[string]interface{})
	paramSet := buildMap["parameters"].(*schema.Set)
	if paramSet.Len() != 1 {
		return &components.Build{Parameter: components.Parameter{}}
	}
	paramMap := paramSet.List()[0].(map[string]interface{})

	return &components.Build{
		Parameter: components.Parameter{
			BuildCmd:          paramMap["build_cmd"].(string),
			ArtifactNamespace: paramMap["artifact_namespace"].(string),
			ClusterId:         paramMap["cluster_id"].(string),
			DockerfilePath:    paramMap["dockerfile_path"].(string),
			NodeLabelSelector: paramMap["node_label_selector"].(map[string]interface{}),
		},
	}
}

func resourceComponentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.ServiceStageV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ServiceStage V3 client: %s", err)
	}

	appId := d.Get("application_id").(string)
	opt := components.CreateOpts{
		Name:                   d.Get("name").(string),
		WorkloadName:           d.Get("workload_name").(string),
		Description:            d.Get("description").(string),
		Labels:                 buildCompLabels(d.Get("labels").([]interface{})),
		PodLabels:              buildCompLabels(d.Get("pod_labels").([]interface{})),
		Version:                d.Get("version").(string),
		EnvironmentID:          d.Get("environment_id").(string),
		ApplicationID:          d.Get("application_id").(string),
		EnterpriseProjectId:    d.Get("enterprise_project_id").(string),
		LimitCpu:               d.Get("limit_cpu").(float64),
		LimitMemory:            d.Get("limit_memory").(float64),
		RequestCpu:             d.Get("request_cpu").(float64),
		RequestMemory:          d.Get("request_memory").(float64),
		Replica:                d.Get("replica").(int),
		EnableSermantInjection: d.Get("enable_sermant_injection").(bool),
		Timezone:               d.Get("timezone").(string),
		JvmOpts:                d.Get("jvm_opts").(string),
		WorkloadKind:           d.Get("workload_kind").(string),
		RuntimeStack:           buildCompRuntimeStack(d.Get("runtime_stack").(interface{})),
		Build:                  buildRepoBuildStructure(d.Get("build").(interface{})),
		Source:                 buildRepoSourceStructure(d.Get("source").(interface{})),
		Envs:                   buildEnvsStructure(d.Get("envs").([]interface{})),
		Storages:               buildStoragesStructure(d.Get("storage").([]interface{})),
		DeployStrategy:         buildDeployStrategyStructure(d.Get("deploy_strategy").(interface{})),
		Command:                buildCommandStructure(d.Get("command").(interface{})),
		PostStart:              buildComponentLifecycleStructure(d.Get("post_start").(interface{})),
		PreStop:                buildComponentLifecycleStructure(d.Get("pre_stop").(interface{})),
		Mesher:                 buildMesherStructure(d.Get("mesher").(interface{})),
		TomcatOpts:             buildTomcatOptStructure(d.Get("tomcat_opts").(interface{})),
		HostAliases:            buildHostAliasesStructure(d.Get("host_aliases").([]interface{})),
		DnsPolicy:              d.Get("dns_policy").(string),
		DnsConfig:              buildDnsConfigStructure(d.Get("dns_config").(interface{})),
		SecurityContext:        buildSecurityContextStructure(d.Get("security_context").(interface{})),
		Logs:                   buildLogsStructure(d.Get("logs").([]interface{})),
		CustomMetric:           buildCustomMetricStructure(d.Get("custom_metric").(interface{})),
		Affinity:               buildComponentAffinityStructure(d.Get("affinity").(interface{})),
		AntiAffinity:           buildComponentAffinityStructure(d.Get("anti_affinity").(interface{})),
		LivenessProbe:          buildComponentProbeStructure(d.Get("liveness_probe").(interface{})),
		ReadinessProbe:         buildComponentProbeStructure(d.Get("readiness_probe").(interface{})),
		ReferResources:         buildReferResourcesStructure(d.Get("refer_resources").([]interface{})),
	}
	resp, err := components.Create(client, appId, opt)
	if err != nil {
		return diag.Errorf("error creating ServiceStage component: %s", err)
	}

	d.SetId(resp.ComponentId)

	log.Printf("[DEBUG] Waiting for the component instance to become running, the instance ID is %s.", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"RUNNING"},
		Target:       []string{"SUCCEEDED"},
		Refresh:      componentInstanceRefreshFunc(client, resp.JobId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the creation of component (%s) to complete: %s",
			d.Id(), err)
	}

	return resourceComponentRead(ctx, d, meta)
}

func buildReferResourcesStructure(referResources []interface{}) []*components.Resource {
	if len(referResources) < 1 {
		return []*components.Resource{}
	}

	var result []*components.Resource
	for _, r := range referResources {
		var referResource components.Resource
		l := r.(map[string]interface{})
		referResource.ID = l["id"].(string)
		referResource.Type = l["type"].(string)

		parameterSet := l["parameters"].(*schema.Set)
		if parameterSet.Len() == 1 {
			var resourceParameter components.ResourceParameters
			parameters := parameterSet.List()[0].(map[string]interface{})
			if parameters != nil {
				resourceParameter.Type = parameters["type"].(string)
				resourceParameter.NameSpace = parameters["namespace"].(string)
			}
			referResource.Parameters = &resourceParameter
		}
		result = append(result, &referResource)
	}

	return result
}

func buildRepoSourceStructure(sources interface{}) *components.Source {
	sourcesSet := sources.(*schema.Set)

	if sourcesSet.Len() != 1 {
		return nil
	}

	source := sourcesSet.List()[0].(map[string]interface{})

	codeArtProjectId := ""
	codeArtProjectIdTemp := source["codearts_project_id"]
	if codeArtProjectIdTemp != nil {
		codeArtProjectId = codeArtProjectIdTemp.(string)
	}

	return &components.Source{
		Kind:              source["kind"].(string),
		Url:               source["url"].(string),
		Version:           source["version"].(string),
		Storage:           source["storage"].(string),
		CodeartsProjectId: codeArtProjectId,
	}
}

func buildEnvsStructure(envs []interface{}) []*components.Env {
	if len(envs) < 1 {
		return nil
	}

	var result []*components.Env
	for _, env := range envs {
		var environmentLabel components.Env
		e := env.(map[string]interface{})
		environmentLabel.Name = e["name"].(string)
		environmentLabel.Value = e["value"].(string)


		valueFromSet := e["value_from"].(*schema.Set)
		if valueFromSet.Len() == 1 {
			var envValueFrom components.EnvValueFrom
			valueFrom := valueFromSet.List()[0].(map[string]interface{})
			if valueFrom != nil {
				envValueFrom.ReferenceType = valueFrom["reference_type"].(string)
				envValueFrom.Name = valueFrom["name"].(string)
				envValueFrom.Key = valueFrom["key"].(string)
				envValueFrom.Optional = valueFrom["optional"].(bool)
			}
			environmentLabel.EnvValueFrom = &envValueFrom
		}
		result = append(result, &environmentLabel)
	}

	return result
}

func buildStoragesStructure(storages []interface{}) []*components.Storage {
	if len(storages) < 1 {
		return nil
	}
	var result []*components.Storage
	for _, storage := range storages {
		var environmentLabel components.Storage
		s := storage.(map[string]interface{})
		environmentLabel.Name = s["name"].(string)
		environmentLabel.Type = s["type"].(string)
		// parameters := s["parameters"].(map[string]interface{})

		parameterSet := s["parameters"].(*schema.Set)
		if parameterSet.Len() == 1 {
			var parameter components.StorageParameter
			parameters := parameterSet.List()[0].(map[string]interface{})
			if parameters != nil {
				parameter.Path = parameters["path"].(string)
				parameter.Name = parameters["name"].(string)
				parameter.DefaultMode = parameters["default_mode"].(int)
				parameter.Medium = parameters["medium"].(string)
			}
			environmentLabel.Parameters = &parameter
		}
		environmentLabel.Mounts = buildStorageMountsStructure(s["mounts"].([]interface{}))
		result = append(result, &environmentLabel)
	}

	return result
}

func buildStorageMountsStructure(mounts []interface{}) []*components.StorageMounts {
	if len(mounts) < 1 {
		return nil
	}

	var result []*components.StorageMounts
	for _, mount := range mounts {
		var environmentLabel components.StorageMounts
		m := mount.(map[string]interface{})
		environmentLabel.Path = m["path"].(string)
		environmentLabel.SubPath = m["sub_path"].(string)
		environmentLabel.Readonly = m["read_only"].(bool)
		result = append(result, &environmentLabel)
	}

	return result
}

func buildDeployStrategyStructure(deployStrategy interface{}) *components.DeployStrategy {
	deployStrategySet := deployStrategy.(*schema.Set)

	if deployStrategySet.Len() != 1 {
		return nil
	}

	var deploy = &components.DeployStrategy{}
	d := deployStrategySet.List()[0].(map[string]interface{})

	deploy.Type = d["type"].(string)
	// rollingRelease := d["rolling_release"].(map[string]interface{})

	rollingReleaseSet := d["rolling_release"].(*schema.Set)
	if rollingReleaseSet.Len() == 1 {
		var rollingReleaseVar components.RollingRelease
		rollingRelease := rollingReleaseSet.List()[0].(map[string]interface{})
		if rollingRelease != nil {
			rollingReleaseVar.Batches = rollingRelease["batches"].(int)
			rollingReleaseVar.TerminationSeconds = rollingRelease["termination_seconds"].(int)
			rollingReleaseVar.FailStrategy = rollingRelease["fail_strategy"].(string)
		}
		deploy.RollingRelease = &rollingReleaseVar
	}

	// grayRelease := d["gray_release"].(map[string]interface{})
	grayReleaseSet := d["gray_release"].(*schema.Set)
	if grayReleaseSet.Len() == 1 {
		var grayReleaseVar components.GrayRelease
		grayRelease := grayReleaseSet.List()[0].(map[string]interface{})
		if grayRelease != nil {
			grayReleaseVar.Type = grayRelease["type"].(string)
			grayReleaseVar.FirstBatchWeight = grayRelease["first_batch_weight"].(int)
			grayReleaseVar.FirstBatchReplica = grayRelease["first_batch_replica"].(int)
			grayReleaseVar.RemainingBatch = grayRelease["remaining_batch"].(int)
			grayReleaseVar.DeploymentMode = grayRelease["deployment_mode"].(int)
			grayReleaseVar.ReplicaSurgeMode = grayRelease["replica_surge_mode"].(string)
			grayReleaseVar.RuleMatchMode = grayRelease["rule_match_mode"].(string)
			grayReleaseVar.Rules = buildGrayRulesStructure(grayRelease["rules"].([]interface{}))
		}
		deploy.GrayRelease = &grayReleaseVar
	}

	return deploy
}

func buildGrayRulesStructure(rules []interface{}) []*components.GrayReleaseRule {
	if len(rules) < 1 {
		return nil
	}
	var result []*components.GrayReleaseRule
	for _, rule := range rules {
		var environmentLabel components.GrayReleaseRule
		r := rule.(map[string]interface{})
		environmentLabel.Key = r["key"].(string)
		environmentLabel.Type = r["type"].(string)
		environmentLabel.Value = r["value"].(string)
		environmentLabel.Condition = r["condition"].(string)
		result = append(result, &environmentLabel)
	}

	return result
}

func buildCommandStructure(commands interface{}) *components.Command {
	commandSet := commands.(*schema.Set)

	if commandSet.Len() != 1 {
		return nil
	}
	command := commandSet.List()[0].(map[string]interface{})
	var environmentLabel = &components.Command{}

	commandList := command["command"].([]interface{})
	commandValue := make([]string, len(commandList))
	for i, raw := range commandList {
		commandValue[i] = raw.(string)
	}
	environmentLabel.Command = commandValue

	argsList := command["args"].([]interface{})
	argsValue := make([]string, len(argsList))
	for i, raw := range argsList {
		argsValue[i] = raw.(string)
	}
	environmentLabel.Args = argsValue

	return environmentLabel
}

func buildComponentLifecycleStructure(componentLifecycle interface{}) *components.K8sLifeCycle {
	componentLifecycleSet := componentLifecycle.(*schema.Set)

	if componentLifecycleSet.Len() != 1 {
		return nil
	}
	lifecycle := componentLifecycleSet.List()[0].(map[string]interface{})

	var k8slifecycle = &components.K8sLifeCycle{}
	k8slifecycle.Type = lifecycle["type"].(string)
	if v, ok := lifecycle["scheme"].(string); ok {
		k8slifecycle.Scheme = v
	}
	if v, ok := lifecycle["host"].(string); ok {
		k8slifecycle.Host = v
	}
	if v, ok := lifecycle["path"].(string); ok {
		k8slifecycle.Path = v
	}
	if v, ok := lifecycle["port"].(int); ok {
		k8slifecycle.Port = v
	}
	// k8slifecycle.Port = lifecycle["port"].(int)

	commandList := lifecycle["command"].([]interface{})
	commandValue := make([]string, len(commandList))
	for i, raw := range commandList {
		commandValue[i] = raw.(string)
	}
	k8slifecycle.Command = commandValue

	return k8slifecycle
}

func buildMesherStructure(mesher interface{}) *components.Mesher {
	mesherSet := mesher.(*schema.Set)

	if mesherSet.Len() != 1 {
		return nil
	}
	m := mesherSet.List()[0].(map[string]interface{})

	var comMesher = &components.Mesher{}
	comMesher.Port = m["port"].(int)

	return comMesher
}

func buildTomcatOptStructure(tomcatOpts interface{}) *components.TomcatOpts {
	tomcatOptsSet := tomcatOpts.(*schema.Set)

	if tomcatOptsSet.Len() != 1 {
		return nil
	}
	t := tomcatOptsSet.List()[0].(map[string]interface{})

	var tomcatOpt = &components.TomcatOpts{}
	tomcatOpt.ServerXml = t["server_xml"].(string)

	return tomcatOpt
}

func buildHostAliasesStructure(hostAliases []interface{}) []*components.HostAlias {
	if len(hostAliases) < 1 {
		return nil
	}
	var result []*components.HostAlias
	for _, hostAlias := range hostAliases {
		var environmentLabel components.HostAlias
		h := hostAlias.(map[string]interface{})
		environmentLabel.IP = h["ip"].(string)

		hostnamesList := h["hostnames"].([]interface{})
		hostnamesValue := make([]string, len(hostnamesList))
		for i, raw := range hostnamesList {
			hostnamesValue[i] = raw.(string)
		}
		environmentLabel.HostNames = hostnamesValue
		result = append(result, &environmentLabel)
	}

	return result
}

func buildDnsConfigStructure(dnsconfig interface{}) *components.DnsConfig {
	dnsconfigSet := dnsconfig.(*schema.Set)

	if dnsconfigSet.Len() != 1 {
		return nil
	}
	d := dnsconfigSet.List()[0].(map[string]interface{})

	var conf = &components.DnsConfig{}

	nameserversList := d["nameservers"].([]interface{})
	nameserversValue := make([]string, len(nameserversList))
	for i, raw := range nameserversList {
		nameserversValue[i] = raw.(string)
	}
	conf.Nameservers = nameserversValue

	searchesList := d["nameservers"].([]interface{})
	searchesValue := make([]string, len(searchesList))
	for i, raw := range searchesList {
		searchesValue[i] = raw.(string)
	}
	conf.Searches = searchesValue
	conf.Options = buildDnsConfigOptionsStructure(d["options"].([]interface{}))

	return conf
}

func buildDnsConfigOptionsStructure(dnsConfigOptions []interface{}) []*components.NameValue {
	if len(dnsConfigOptions) < 1 {
		return nil
	}
	var result []*components.NameValue
	for _, options := range dnsConfigOptions {
		var environmentLabel components.NameValue
		o := options.(map[string]interface{})
		environmentLabel.Name = o["name"].(string)
		environmentLabel.Value = o["value"].(string)
		result = append(result, &environmentLabel)
	}

	return result
}

func buildSecurityContextStructure(securityContext interface{}) *components.SecurityContext {
	securityContextSet := securityContext.(*schema.Set)

	if securityContextSet.Len() != 1 {
		return nil
	}
	d := securityContextSet.List()[0].(map[string]interface{})

	var conf = &components.SecurityContext{}
	conf.RunAsUser = d["run_as_user"].(int)
	conf.RunAsGroup = d["run_as_group"].(int)
	capabilities := d["capabilities"].(interface{})

	conf.Capabilities = buildSecurityContextCapabilitiesStructure(capabilities)

	return conf
}

func buildSecurityContextCapabilitiesStructure(capabilities interface{}) *components.Capabilities {
	capabilitiesSet := capabilities.(*schema.Set)
	if capabilitiesSet.Len() != 1 {
		return nil
	}

	c := capabilitiesSet.List()[0].(map[string]interface{})

	var capability components.Capabilities
	addList := c["add"].([]interface{})
	addValue := make([]string, len(addList))
	for i, raw := range addList {
		addValue[i] = raw.(string)
	}
	capability.Add = addValue

	dropList := c["drop"].([]interface{})
	dropValue := make([]string, len(dropList))
	for i, raw := range dropList {
		dropValue[i] = raw.(string)
	}
	capability.Drop = dropValue

	return &capability
}

func buildLogsStructure(logs []interface{}) []*components.Log {
	if len(logs) < 1 {
		return nil
	}
	var result []*components.Log
	for _, log := range logs {
		var environmentLabel components.Log
		l := log.(map[string]interface{})
		environmentLabel.LogPath = l["log_path"].(string)
		environmentLabel.Rotate = l["rotate"].(string)
		environmentLabel.HostPath = l["host_path"].(string)
		environmentLabel.HostExtendPath = l["host_extend_path"].(string)
		result = append(result, &environmentLabel)
	}

	return result
}

func buildCustomMetricStructure(customMetric interface{}) *components.CustomMetric {
	customMetricSet := customMetric.(*schema.Set)

	if customMetricSet.Len() != 1 {
		return nil
	}
	m := customMetricSet.List()[0].(map[string]interface{})

	var conf = &components.CustomMetric{}
	conf.Path = m["path"].(string)
	conf.Port = m["port"].(int)
	conf.Dimensions = m["dimensions"].(string)

	return conf
}

func buildComponentAffinityStructure(affinity interface{}) *components.Affinity {
	affinitySet := affinity.(*schema.Set)

	if affinitySet.Len() != 1 {
		return nil
	}
	a := affinitySet.List()[0].(map[string]interface{})

	var conf = &components.Affinity{}
	conf.Kind = a["kind"].(string)
	conf.Condition = a["condition"].(string)
	conf.Weight = a["weight"].(int)
	conf.MatchExpressions = buildAffinityExpressionStructure(a["match_expressions"].([]interface{}))

	return conf
}

func buildAffinityExpressionStructure(matchExpressions []interface{}) []*components.MatchExpression {
	if len(matchExpressions) < 1 {
		return nil
	}
	var result []*components.MatchExpression
	for _, expression := range matchExpressions {
		var environmentLabel components.MatchExpression
		e := expression.(map[string]interface{})
		environmentLabel.Key = e["key"].(string)
		environmentLabel.Value = e["value"].(string)
		environmentLabel.Operation = e["operation"].(string)
		result = append(result, &environmentLabel)
	}

	return result
}

func buildComponentProbeStructure(componentProbe interface{}) *components.K8sProbe {
	componentLifecycleSet := componentProbe.(*schema.Set)

	if componentLifecycleSet.Len() != 1 {
		return nil
	}
	probe := componentLifecycleSet.List()[0].(map[string]interface{})

	var k8Probe = &components.K8sProbe{}
	k8Probe.Type = probe["type"].(string)
	k8Probe.Delay = probe["delay"].(int)
	k8Probe.Timeout = probe["timeout"].(int)
	if v, ok := probe["period_seconds"].(int); ok {
		k8Probe.PeriodSeconds = v
	}
	if v, ok := probe["success_Threshold"].(int); ok {
		k8Probe.SuccessThreshold = v
	}
	if v, ok := probe["failure_threshold"].(int); ok {
		k8Probe.FailureThreshold = v
	}
	// k8Probe.PeriodSeconds = probe["period_seconds"].(int)
	// k8Probe.SuccessThreshold = probe["success_Threshold"].(int)
	// k8Probe.FailureThreshold = probe["failure_threshold"].(int)

	if v, ok := probe["scheme"].(string); ok {
		k8Probe.Scheme = v
	}
	if v, ok := probe["host"].(string); ok {
		k8Probe.Host = v
	}
	if v, ok := probe["path"].(string); ok {
		k8Probe.Path = v
	}
	if v, ok := probe["port"].(int); ok {
		k8Probe.Port = v
	}

	commandList := probe["command"].([]interface{})
	commandValue := make([]string, len(commandList))
	for i, raw := range commandList {
		commandValue[i] = raw.(string)
	}
	k8Probe.Command = commandValue

	return k8Probe
}

func buildCompRuntimeStack(runtimeStacks interface{}) components.RuntimeStack {
	runtimeStacksSet := runtimeStacks.(*schema.Set)

	if runtimeStacksSet.Len() != 1 {
		return components.RuntimeStack{}
	}

	runtimeStack := runtimeStacksSet.List()[0].(map[string]interface{})
	return components.RuntimeStack{
		Name:       runtimeStack["name"].(string),
		Version:    runtimeStack["version"].(string),
		Type:       runtimeStack["type"].(string),
		DeployMode: runtimeStack["deploy_mode"].(string),
	}
}

func buildCompLabels(labels []interface{}) []*components.KeyValue {
	if len(labels) < 1 {
		return nil
	}
	var result []*components.KeyValue
	for _, label := range labels {
		var environmentLabel components.KeyValue
		l := label.(map[string]interface{})
		environmentLabel.Key = l["key"].(string)
		environmentLabel.Value = l["value"].(string)
		result = append(result, &environmentLabel)
	}

	return result
}

// func flattenRepoBuilder(builder components.Builder) (result []map[string]interface{}) {
//	defer func() {
//		if r := recover(); r != nil {
//			log.Printf("[ERROR] Recover panic when flattening builder structure: %#v", r)
//		}
//	}()
//
//	if !reflect.DeepEqual(builder, components.Builder{}) {
//		result = append(result, map[string]interface{}{
//			"cmd":                builder.Parameter.BuildCmd,
//			"organization":       builder.Parameter.ArtifactNamespace,
//			"cluster_id":         builder.Parameter.ClusterId,
//			"cluster_name":       builder.Parameter.ClusterName,
//			"cluster_type":       builder.Parameter.ClusterType,
//			"dockerfile_path":    builder.Parameter.DockerfilePath,
//			"use_public_cluster": builder.Parameter.UsePublicCluster,
//			"node_label":         builder.Parameter.NodeLabelSelector,
//		})
//	}
//
//	return
// }
//
// func flattenRepoSource(source components.Source) (result []map[string]interface{}) {
//	defer func() {
//		if r := recover(); r != nil {
//			log.Printf("[ERROR] Recover panic when flattening source structure: %#v", r)
//		}
//	}()
//
//	if (source != components.Source{}) {
//		if source.Spec.Type == "package" {
//			result = append(result, map[string]interface{}{
//				"type":         source.Spec.Type,
//				"storage_type": source.Spec.Storage,
//				"url":          source.Spec.Url,
//			})
//		} else if source.Spec.RepoType == "GitHub" || source.Spec.Type == "GitLab" ||
//			source.Spec.Type == "Gitee" || source.Spec.Type == "Bitbucket" || source.Spec.Type == "DevCloud" {
//			result = append(result, map[string]interface{}{
//				"type":           source.Spec.RepoType,
//				"authorization":  source.Spec.RepoAuth,
//				"url":            source.Spec.RepoUrl,
//				"repo_ref":       source.Spec.RepoRef,
//				"repo_namespace": source.Spec.RepoNamespace,
//			})
//		}
//	}
//
//	return
// }

func componentInstanceRefreshFunc(c *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opt := jobs.ListOpts{
			Limit: 50,
		}
		resp, err := jobs.List(c, jobId, opt)
		if err != nil {
			return resp, "ERROR", err
		}
		rl := len(resp)
		if rl < 1 {
			return resp, "NO TASK", nil
		}
		return resp, resp[rl-1].Status, nil
	}
}

func resourceComponentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.ServiceStageV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage V3 client: %s", err)
	}

	appId := d.Get("application_id").(string)
	resp, err := components.Get(client, appId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ServiceStage component")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", resp.Name),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceComponentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.ServiceStageV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ServiceStage V3 client: %s", err)
	}

	appId := d.Get("application_id").(string)
	// In normal changes, there is no situation in which source and builder are empty, so these two empty values are
	// ignored.
	opt := components.UpdateOpts{
		Name:    d.Get("name").(string),
		Builder: buildRepoBuildStructure(d.Get("build").(interface{})),
		Source:  buildRepoSourceStructure(d.Get("source").([]interface{})),
	}
	_, err = components.Update(client, appId, d.Id(), opt)
	if err != nil {
		return diag.Errorf("error updating ServiceStage component (%s): %s", d.Id(), err)
	}

	return resourceComponentRead(ctx, d, meta)
}

func resourceComponentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.ServiceStageV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ServiceStage V3 client: %s", err)
	}

	appId := d.Get("application_id").(string)
	err = components.Delete(client, appId, d.Id())
	if err != nil {
		return diag.Errorf("error deleting ServiceStage component: %s", err)
	}
	return nil
}

func resourceComponentImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid format specified for import id, must be <application_id>/<component_id>")
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("application_id", parts[0])
}
