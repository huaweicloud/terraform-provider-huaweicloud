package servicestage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/servicestage/v2/instances"
	"github.com/chnsz/golangsdk/openstack/servicestage/v2/jobs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func lifecycleProcessSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"command", "http",
				}, false),
			},
			"parameters": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"commands": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"host": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func affinitySchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"availability_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"private_ips": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instance_names": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func probeDetailSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"command", "http", "tcp",
				}, false),
			},
			"command_param": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"commands": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"http_param": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scheme": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"HTTP", "HTTPS",
							}, false),
						},
						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"host": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"tcp_param": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"delay": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

// ResourceComponentInstance is the imple of huaweicloud_servicestage_component_instance
// @API ServiceStage POST /v2/{project_id}/cas/applications/{application_id}/components/{component_id}/instances
// @API ServiceStage GET /v2/{project_id}/cas/jobs/{job_id}
// @API ServiceStage GET /v2/{project_id}/cas/applications/{application_id}/components/{component_id}/instances/{instance_id}
// @API ServiceStage PUT /v2/{project_id}/cas/applications/{application_id}/components/{component_id}/instances/{instance_id}
// @API ServiceStage DELETE /v2/{project_id}/cas/applications/{application_id}/components/{component_id}/instances/{instance_id}
func ResourceComponentInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComponentInstanceCreate,
		ReadContext:   resourceComponentInstanceRead,
		UpdateContext: resourceComponentInstanceUpdate,
		DeleteContext: resourceComponentInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceComponentInstanceImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"component_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"refer_resource": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"alias": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"parameters": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"artifact": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"package", "image",
							}, false),
						},
						"storage": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"swr", "obs",
							}, false), // The devcloud does not support yet.
						},
						"url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"auth_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "iam",
							ValidateFunc: validation.StringInSlice([]string{
								"iam", "none",
							}, false),
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"properties": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"endpoint": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"key": {
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"env_variable": {
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
						"storage": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"HostPath", "EmptyDir", "ConfigMap", "Secret", "PersistentVolumeClaim",
										}, false),
									},
									"parameter": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"path": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"claim_name": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"secret_name": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"mount": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"path": {
													Type:     schema.TypeString,
													Required: true,
												},
												"readonly": {
													Type:     schema.TypeBool,
													Required: true,
												},
												"subpath": {
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
						"strategy": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"upgrade": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "RollingUpdate",
										ValidateFunc: validation.StringInSlice([]string{
											"RollingUpdate", "Recreate",
										}, false),
									},
								},
							},
						},
						"lifecycle": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"entrypoint": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"commands": {
													Type:     schema.TypeList,
													Required: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"args": {
													Type:     schema.TypeList,
													Required: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"post_start": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem:     lifecycleProcessSchemaResource(),
									},
									"pre_stop": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem:     lifecycleProcessSchemaResource(),
									},
								},
							},
						},
						"log_collection_policy": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"container_mounting": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"path": {
													Type:     schema.TypeString,
													Required: true,
												},
												"host_extend_path": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"aging_period": {
													Type:     schema.TypeString,
													Optional: true,
													Default:  "Hourly",
													ValidateFunc: validation.StringInSlice([]string{
														"Hourly", "Daily", "Weekly",
													}, false),
												},
											},
										},
									},
									"host_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"scheduler": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"affinity": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem:     affinitySchemaResource(),
									},
									"anti_affinity": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem:     affinitySchemaResource(),
									},
								},
							},
						},
						"probe": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"liveness": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem:     probeDetailSchemaResource(),
									},
									"readiness": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem:     probeDetailSchemaResource(),
									},
								},
							},
						},
					},
				},
			},
			"external_access": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"HTTP", "HTTPS",
							}, false),
						},
						"address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildArtifactStructure(artifacts *schema.Set) map[string]instances.Artifact {
	if artifacts.Len() < 1 {
		return nil
	}

	result := make(map[string]instances.Artifact)
	for _, val := range artifacts.List() {
		artifact := val.(map[string]interface{})
		s := instances.Artifact{
			Type:    artifact["type"].(string),
			Storage: artifact["storage"].(string),
			URL:     artifact["url"].(string),
			Auth:    artifact["auth_type"].(string),
			Version: artifact["version"].(string),
		}
		properties := artifact["properties"].([]interface{})
		if len(properties) > 0 {
			property := properties[0].(map[string]interface{})
			s.Properties = map[string]interface{}{
				"bucket":   property["bucket"].(string),
				"endpoint": property["endpoint"].(string),
				"key":      property["key"].(string),
			}
		}
		result[artifact["name"].(string)] = s
	}

	return result
}

func buildReferResourcesList(resources *schema.Set) ([]instances.ReferResource, error) {
	if resources.Len() < 1 {
		return nil, nil
	}

	result := make([]instances.ReferResource, resources.Len())
	for i, val := range resources.List() {
		res := val.(map[string]interface{})
		refer := instances.ReferResource{
			Type:       res["type"].(string),
			ID:         res["id"].(string),
			ReferAlias: res["alias"].(string),
		}
		pResult := make(map[string]interface{})
		if param, ok := res["parameters"]; ok {
			log.Printf("[DEBUG] The parameters is %#v", param)
			p := param.(map[string]interface{})
			for k, v := range p {
				if k == "hosts" {
					var r []string
					err := json.Unmarshal([]byte(v.(string)), &r)
					if err != nil {
						return nil, fmt.Errorf("the format of the host value is not right: %#v", v)
					}
					pResult[k] = &r
					continue
				}
				pResult[k] = v
			}

			refer.Parameters = pResult
		}

		log.Printf("[DEBUG] The parameter map is %v", pResult)
		result[i] = refer
	}

	return result, nil
}

func buildEnvVariables(variables *schema.Set) []instances.Variable {
	if variables.Len() < 1 {
		return nil
	}

	result := make([]instances.Variable, 0, variables.Len())
	for _, val := range variables.List() {
		variable := val.(map[string]interface{})
		result = append(result, instances.Variable{
			Name:  variable["name"].(string),
			Value: variable["value"].(string),
		})
	}

	return result
}

func buildMountsList(mounts *schema.Set) []instances.Mount {
	if mounts.Len() < 1 {
		return nil
	}

	result := make([]instances.Mount, mounts.Len())
	for i, val := range mounts.List() {
		mount := val.(map[string]interface{})
		result[i] = instances.Mount{
			Path:     mount["path"].(string),
			SubPath:  mount["subpath"].(string),
			ReadOnly: utils.Bool(mount["readonly"].(bool)),
		}
	}

	return result
}

func buildStoragesList(storages *schema.Set) []instances.Storage {
	if storages.Len() < 1 {
		return nil
	}

	result := make([]instances.Storage, storages.Len())
	for i, val := range storages.List() {
		storage := val.(map[string]interface{})
		var parameters instances.StorageParams
		if paramVal, ok := storage["parameter"]; ok {
			parameter := paramVal.([]interface{})[0].(map[string]interface{})
			parameters.Path = parameter["path"].(string)
			parameters.Name = parameter["name"].(string)
			parameters.ClaimName = parameter["claim_name"].(string)
			parameters.SecretName = parameter["secret_name"].(string)
		}

		result[i] = instances.Storage{
			Type:       storage["type"].(string),
			Parameters: &parameters,
			Mounts:     buildMountsList(storage["mount"].(*schema.Set)),
		}
	}

	return result
}

func buildStrategyStructure(strategies []interface{}) *instances.Strategy {
	if len(strategies) < 1 {
		return nil
	}

	strategy := strategies[0].(map[string]interface{})

	return &instances.Strategy{
		Upgrade: strategy["upgrade"].(string),
	}
}

func buildLifecycleProcess(processes []interface{}) *instances.Process {
	if len(processes) < 1 {
		return nil
	}

	process := processes[0].(map[string]interface{})
	// The configuration structure of the process parameters is required.
	parameters := process["parameters"].([]interface{})
	param := parameters[0].(map[string]interface{})

	return &instances.Process{
		Type: process["type"].(string),
		Parameters: &instances.ProcessParams{
			Commands: utils.ExpandToStringList(param["commands"].([]interface{})),
			Port:     param["port"].(int),
			Path:     param["path"].(string),
			Host:     param["host"].(string),
		},
	}
}

func buildLifecycleStructure(lifecycles []interface{}) *instances.Lifecycle {
	if len(lifecycles) < 1 {
		return nil
	}

	lifecycle := lifecycles[0].(map[string]interface{})
	result := instances.Lifecycle{
		PostStart: buildLifecycleProcess(lifecycle["post_start"].([]interface{})),
		PreStop:   buildLifecycleProcess(lifecycle["pre_stop"].([]interface{})),
	}

	if val, ok := lifecycle["entrypoint"]; ok && len(val.([]interface{})) > 0 {
		entrypoint := val.([]interface{})[0].(map[string]interface{})
		result.Entrypoint = &instances.Entrypoint{
			Commands: utils.ExpandToStringList(entrypoint["commands"].([]interface{})),
			Args:     utils.ExpandToStringList(entrypoint["args"].([]interface{})),
		}
	}

	return &result
}

func buildLogCollectionPoliciesStructure(policies *schema.Set) []instances.LogCollectionPolicy {
	if policies.Len() < 1 {
		return nil
	}

	result := make([]instances.LogCollectionPolicy, 0, policies.Len())
	for _, val := range policies.List() {
		policy := val.(map[string]interface{})
		hostPath := policy["host_path"].(string)
		cmSet := policy["container_mounting"].(*schema.Set)
		for _, val := range cmSet.List() {
			cm := val.(map[string]interface{})
			result = append(result, instances.LogCollectionPolicy{
				LogPath:        cm["path"].(string),
				HostExtendPath: cm["host_extend_path"].(string),
				AgingPeriod:    cm["aging_period"].(string),
				HostPath:       hostPath,
			})
		}
	}

	return result
}

func buildAffinityStructure(affinities []interface{}) *instances.Affinity {
	if len(affinities) < 1 {
		return nil
	}

	result := instances.Affinity{}

	affinity := affinities[0].(map[string]interface{})
	if val, ok := affinity["availability_zones"]; ok && len(val.([]interface{})) > 0 {
		result.AvailabilityZones = utils.ExpandToStringList(val.([]interface{}))
	}
	if val, ok := affinity["private_ips"]; ok && len(val.([]interface{})) > 0 {
		result.Nodes = utils.ExpandToStringList(val.([]interface{}))
	}
	if val, ok := affinity["instance_names"]; ok && len(val.([]interface{})) > 0 {
		result.Applications = utils.ExpandToStringList(val.([]interface{}))
	}

	return &result
}

func buildSchedulerStructure(schedulers []interface{}) *instances.Scheduler {
	if len(schedulers) < 1 {
		return nil
	}

	scheduler := schedulers[0].(map[string]interface{})
	return &instances.Scheduler{
		Affinity:     buildAffinityStructure(scheduler["affinity"].([]interface{})),
		AntiAffinity: buildAffinityStructure(scheduler["anti_affinity"].([]interface{})),
	}
}

func buildProbeDetailStructure(details []interface{}) (*instances.ProbeDetail, error) {
	if len(details) < 1 {
		return nil, nil
	}

	detail := details[0].(map[string]interface{})
	pType := detail["type"].(string)
	result := instances.ProbeDetail{
		Type:    pType,
		Delay:   detail["delay"].(int),
		Timeout: detail["timeout"].(int),
	}

	params := make(map[string]interface{})
	switch pType {
	case "command":
		cmdParams := detail["command_param"].([]interface{})
		if len(cmdParams) < 1 {
			return nil, fmt.Errorf("The command parameters must be set if the probe type is 'command'.")
		}
		cmdParam := cmdParams[0].(map[string]interface{})
		params["command"] = utils.ExpandToStringList(cmdParam["commands"].([]interface{}))
	case "http":
		httpParams := detail["http_param"].([]interface{})
		if len(httpParams) < 1 {
			return nil, fmt.Errorf("The http parameters must be set if the probe type is 'http'.")
		}
		httpParam := httpParams[0].(map[string]interface{})
		params["scheme"] = httpParam["scheme"]
		params["port"] = httpParam["port"]
		params["path"] = httpParam["path"]
		params["host"] = httpParam["host"]
	case "tcp":
		tcpParams := detail["tcp_param"].([]interface{})
		if len(tcpParams) < 1 {
			return nil, fmt.Errorf("The tcp parameters must be set if the probe type is 'tcp'.")
		}
		tcpParam := tcpParams[0].(map[string]interface{})
		params["port"] = tcpParam["port"]
	default:
		return nil, fmt.Errorf("One of the following parameters must be set: command, http or tcp.")
	}

	result.Parameters = params
	return &result, nil
}

func buildProbeStructure(probes []interface{}) (*instances.Probe, error) {
	if len(probes) < 1 {
		return nil, nil
	}

	probe := probes[0].(map[string]interface{})
	liveness, err := buildProbeDetailStructure(probe["liveness"].([]interface{}))
	if err != nil {
		return nil, err
	}
	readiness, err := buildProbeDetailStructure(probe["readiness"].([]interface{}))
	if err != nil {
		return nil, err
	}

	return &instances.Probe{
		LivenessProbe:  liveness,
		ReadinessProbe: readiness,
	}, nil
}

func buildConfigurationStructure(configs []interface{}) (instances.Configuration, error) {
	if len(configs) < 1 {
		return instances.Configuration{}, nil
	}

	config := configs[0].(map[string]interface{})
	probe, err := buildProbeStructure(config["probe"].([]interface{}))
	if err != nil {
		return instances.Configuration{}, err
	}

	return instances.Configuration{
		EnvVariables:          buildEnvVariables(config["env_variable"].(*schema.Set)),
		Storages:              buildStoragesList(config["storage"].(*schema.Set)),
		Strategy:              buildStrategyStructure(config["strategy"].([]interface{})),
		Lifecycle:             buildLifecycleStructure(config["lifecycle"].([]interface{})),
		LogCollectionPolicies: buildLogCollectionPoliciesStructure(config["log_collection_policy"].(*schema.Set)),
		Scheduler:             buildSchedulerStructure(config["scheduler"].([]interface{})),
		Probe:                 probe,
	}, nil
}

func buildExternalAccessList(accesses *schema.Set) []instances.ExternalAccess {
	if accesses.Len() < 1 {
		return nil
	}

	result := make([]instances.ExternalAccess, accesses.Len())
	for i, val := range accesses.List() {
		access := val.(map[string]interface{})
		result[i] = instances.ExternalAccess{
			Protocol:    access["protocol"].(string),
			Address:     access["address"].(string),
			ForwardPort: access["port"].(int),
		}
	}

	return result
}

func buildInstanceCreateOpts(d *schema.ResourceData) (instances.CreateOpts, error) {
	result := instances.CreateOpts{
		EnvId:       d.Get("environment_id").(string),
		Name:        d.Get("name").(string),
		Version:     d.Get("version").(string),
		Replica:     d.Get("replica").(int),
		FlavorId:    d.Get("flavor_id").(string),
		Description: d.Get("description").(string),
		Artifacts:   buildArtifactStructure(d.Get("artifact").(*schema.Set)),

		ExternalAccesses: buildExternalAccessList(d.Get("external_access").(*schema.Set)),
	}
	referRes, err := buildReferResourcesList(d.Get("refer_resource").(*schema.Set))
	if err != nil {
		return result, err
	}
	result.ReferResources = referRes

	conf, err := buildConfigurationStructure(d.Get("configuration").([]interface{}))
	if err != nil {
		return result, err
	}

	result.Configuration = conf
	return result, nil
}

func resourceComponentInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.ServiceStageV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ServiceStage v2 client: %s", err)
	}

	appId := d.Get("application_id").(string)
	componentId := d.Get("component_id").(string)
	opt, err := buildInstanceCreateOpts(d)
	if err != nil {
		return diag.Errorf("error building the CreateOpts of the component instance: %s", err)
	}
	log.Printf("[DEBUG] The instance create option of ServiceStage component is: %v", opt)

	resp, err := instances.Create(client, appId, componentId, opt)
	if err != nil {
		return diag.Errorf("error creating ServiceStage component instance: %s", err)
	}

	d.SetId(resp.InstanceId)

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
		return diag.Errorf("error waiting for the creation of component instance (%s) to complete: %s",
			d.Id(), err)
	}

	return resourceComponentInstanceRead(ctx, d, meta)
}

func flattenProperties(properties map[string]interface{}) []map[string]interface{} {
	result := make(map[string]interface{})
	if bucket, ok := properties["bucket"]; ok {
		result["bucket"] = bucket
	}
	if endpoint, ok := properties["endpoint"]; ok {
		result["endpoint"] = endpoint
	}
	if key, ok := properties["key"]; ok {
		result["key"] = key
	}
	if len(result) < 1 {
		return nil
	}
	return []map[string]interface{}{result}
}

func flattenArtifact(artifacts map[string]instances.Artifact) []map[string]interface{} {
	if len(artifacts) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(artifacts))
	for key, val := range artifacts {
		s := map[string]interface{}{
			"name":      key,
			"type":      val.Type,
			"storage":   val.Storage,
			"url":       val.URL,
			"auth_type": val.Auth,
			"version":   val.Version,
		}
		if p := flattenProperties(val.Properties); p != nil {
			s["properties"] = p
		}
		result = append(result, s)
	}

	log.Printf("[DEBUG] The artifacts result is %#v", result)
	return result
}

func flattenReferResources(resources []instances.ReferResource) (result []map[string]interface{}) {
	if len(resources) < 1 {
		return nil
	}

	for _, val := range resources {
		s := map[string]interface{}{
			"type":  val.Type,
			"id":    val.ID,
			"alias": val.ReferAlias,
		}
		params := make(map[string]interface{})
		for k, v := range val.Parameters {
			if _, ok := v.([]interface{}); ok {
				jsonByte, _ := json.Marshal(v.([]interface{}))
				params[k] = string(jsonByte)
				continue
			}
			params[k] = v
		}
		if len(params) > 0 {
			s["parameters"] = params
		}
		result = append(result, s)
	}

	log.Printf("[DEBUG] The resources result is %#v", result)
	return
}

func flattenEnvVariables(variables []instances.VariableResp) (result []map[string]interface{}) {
	if len(variables) < 1 {
		return nil
	}

	for _, val := range variables {
		// After the instance is created, the system will automatically add a series of environment variables to it.
		// These variables will mark as internal.
		if val.Internal {
			continue
		}
		result = append(result, map[string]interface{}{
			"name":  val.Name,
			"value": val.Value,
		})
	}

	log.Printf("[DEBUG] The environment variables result is %#v", result)
	return
}

func flattenMounts(mounts []instances.Mount) (result []map[string]interface{}) {
	if len(mounts) < 1 {
		return nil
	}

	for _, val := range mounts {
		result = append(result, map[string]interface{}{
			"path":     val.Path,
			"readonly": val.ReadOnly,
			"subpath":  val.SubPath,
		})
	}

	log.Printf("[DEBUG] The mounts result is %#v", result)
	return
}

func flattenStorages(storages []instances.StorageResp) (result []map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Recover panic when flattening storage structure: %#v", r)
		}
	}()

	if len(storages) < 1 {
		return nil
	}

	for _, val := range storages {
		result = append(result, map[string]interface{}{
			"type": val.Type,
			"parameter": []map[string]interface{}{
				{
					"path":        val.Parameters.Path,
					"name":        val.Parameters.Name,
					"claim_name":  val.Parameters.ClaimName,
					"secret_name": val.Parameters.SecretName,
				},
			},
			"mount": flattenMounts(val.Mounts),
		})
	}

	log.Printf("[DEBUG] The storages result is %#v", result)
	return result
}

func flattenStrategy(strategy instances.StrategyResp) (result []map[string]interface{}) {
	if reflect.DeepEqual(strategy, instances.StrategyResp{}) {
		return nil
	}

	result = append(result, map[string]interface{}{
		"upgrade": strategy.Upgrade,
	})

	log.Printf("[DEBUG] The strategy result is %#v", result)
	return
}

func flattenProcess(process instances.ProcessResp) (result []map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Recover panic when flattening process structure: %#v", r)
		}
	}()

	if reflect.DeepEqual(process, instances.ProcessResp{}) {
		return nil
	}

	result = append(result, map[string]interface{}{
		"type": process.Type,
		"parameters": []map[string]interface{}{
			{
				"commands": process.Parameters.Commands,
				"port":     process.Parameters.Port,
				"path":     process.Parameters.Path,
				"host":     process.Parameters.Host,
			},
		},
	})

	log.Printf("[DEBUG] The process result is %#v", result)
	return
}

func flattenLifecycle(lifecycle instances.LifecycleResp) (result []map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Recover panic when flattening lifecycle structure: %#v", r)
		}
	}()

	if reflect.DeepEqual(lifecycle, instances.LifecycleResp{}) {
		return nil
	}

	result = append(result, map[string]interface{}{
		"entrypoint": []map[string]interface{}{
			{
				"commands": lifecycle.Entrypoint.Commands,
				"args":     lifecycle.Entrypoint.Args,
			},
		},
		"post_start": flattenProcess(lifecycle.PostStart),
		"pre_stop":   flattenProcess(lifecycle.PreStop),
	})

	log.Printf("[DEBUG] The lifecycle result is %#v", result)
	return
}

func flattenLogCollectionPolicies(policies []instances.LogCollectionPolicyResp) (result []map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Recover panic when flattening policies structure: %#v", r)
		}
	}()

	if len(policies) < 1 {
		return nil
	}

	policiesMap := make(map[string][]interface{})
	for _, val := range policies {
		policiesMap[val.HostPath] = append(policiesMap[val.HostPath], map[string]interface{}{
			"path":             val.LogPath,
			"host_extend_path": val.HostExtendPath,
			"aging_period":     val.AgingPeriod,
		})
	}

	for k, v := range policiesMap {
		result = append(result, map[string]interface{}{
			"host_path":          k,
			"container_mounting": v,
		})
	}

	log.Printf("[DEBUG] The collection policies result is %#v", result)
	return
}

func flattenScheduler(scheduler instances.SchedulerResp) (result []map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Recover panic when flattening scheduler structure: %#v", r)
		}
	}()

	if reflect.DeepEqual(scheduler, instances.SchedulerResp{}) {
		return nil
	}

	result = append(result, map[string]interface{}{
		"affinity": []map[string]interface{}{
			{
				"availability_zones": scheduler.Affinity.AvailabilityZones,
				"private_ips":        scheduler.Affinity.Nodes,
				"instance_names":     scheduler.Affinity.Applications,
			},
		},
		"anti_affinity": []map[string]interface{}{
			{
				"availability_zones": scheduler.AntiAffinity.AvailabilityZones,
				"private_ips":        scheduler.AntiAffinity.Nodes,
				"instance_names":     scheduler.AntiAffinity.Applications,
			},
		},
	})

	log.Printf("[DEBUG] The scheduler result is %#v", result)
	return
}

func flattenProbeDetail(detail instances.ProbeDetail) (result []map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Recover panic when flattening probe detail structure: %#v", r)
		}
	}()

	if reflect.DeepEqual(detail, instances.ProbeDetail{}) {
		return nil
	}

	pType := detail.Type
	structure := map[string]interface{}{
		"type":    pType,
		"delay":   detail.Delay,
		"timeout": detail.Timeout,
	}
	switch pType {
	case "command":
		structure["command_param"] = []map[string]interface{}{
			{
				"commands": detail.Parameters["command"],
			},
		}
	case "http":
		structure["http_param"] = []map[string]interface{}{
			{
				"scheme": detail.Parameters["scheme"],
				"port":   detail.Parameters["port"],
				"path":   detail.Parameters["path"],
				"host":   detail.Parameters["host"],
			},
		}
	case "tcp":
		structure["tcp_param"] = []map[string]interface{}{
			{
				"port": detail.Parameters["port"],
			},
		}
	default:
	}

	result = append(result, structure)
	log.Printf("[DEBUG] The probe detail result is %#v", result)
	return
}

func flattenProbe(probe instances.ProbeResp) []map[string]interface{} {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Recover panic when flattening probe structure: %#v", r)
		}
	}()

	if reflect.DeepEqual(probe, instances.ProbeResp{}) {
		return nil
	}

	result := []map[string]interface{}{
		{
			"liveness":  flattenProbeDetail(probe.LivenessProbe),
			"readiness": flattenProbeDetail(probe.ReadinessProbe),
		},
	}
	log.Printf("[DEBUG] The probe result is %#v", result)
	return result
}

func flattenConfiguration(configuration instances.ConfigurationResp) (result []map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Recover panic when flattening configuration structure: %#v", r)
		}
	}()

	if reflect.DeepEqual(configuration, instances.ConfigurationResp{}) {
		return nil
	}

	result = []map[string]interface{}{
		{
			"env_variable":          flattenEnvVariables(configuration.EnvVariables),
			"storage":               flattenStorages(configuration.Storages),
			"strategy":              flattenStrategy(configuration.Strategy),
			"lifecycle":             flattenLifecycle(configuration.Lifecycle),
			"log_collection_policy": flattenLogCollectionPolicies(configuration.LogCollectionPolicy),
			"scheduler":             flattenScheduler(configuration.Scheduler),
			"probe":                 flattenProbe(configuration.Probe),
		},
	}
	log.Printf("[DEBUG] The configuration result is %#v", result)
	return result
}

func flattenExternalAccesses(accesses []instances.ExternalAccessResp) []map[string]interface{} {
	result := make([]map[string]interface{}, len(accesses))

	if len(accesses) < 1 {
		return nil
	}

	for i, val := range accesses {
		result[i] = map[string]interface{}{
			"protocol": val.Protocol,
			"address":  val.Address,
			"port":     val.ForwardPort,
		}
	}
	log.Printf("[DEBUG] The accesses result is %#v", result)
	return result
}

func resourceComponentInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.ServiceStageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage v2 client: %s", err)
	}

	appId := d.Get("application_id").(string)
	componentId := d.Get("component_id").(string)
	resp, err := instances.Get(client, appId, componentId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ServiceStage component instance")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("environment_id", resp.EnvironmentId),
		d.Set("name", resp.Name),
		d.Set("version", resp.Version),
		d.Set("replica", resp.StatusDetail.Replica),
		d.Set("flavor_id", resp.FlavorId),
		d.Set("description", resp.Description),
		d.Set("artifact", flattenArtifact(resp.Artifacts)),
		d.Set("refer_resource", flattenReferResources(resp.ReferResources)),
		d.Set("configuration", flattenConfiguration(resp.Configuration)),
		d.Set("external_access", flattenExternalAccesses(resp.ExternalAccesses)),
		// Attributes
		d.Set("status", resp.StatusDetail.Status),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildInstanceUpdateOpts(d *schema.ResourceData) (instances.UpdateOpts, error) {
	desc := d.Get("description").(string)
	result := instances.UpdateOpts{
		Version:          d.Get("version").(string),
		FlavorId:         d.Get("flavor_id").(string),
		Description:      &desc,
		Artifacts:        buildArtifactStructure(d.Get("artifact").(*schema.Set)),
		ExternalAccesses: buildExternalAccessList(d.Get("external_access").(*schema.Set)),
	}

	referRes, err := buildReferResourcesList(d.Get("refer_resource").(*schema.Set))
	if err != nil {
		return result, err
	}
	result.ReferResources = referRes

	conf, err := buildConfigurationStructure(d.Get("configuration").([]interface{}))
	if err != nil {
		return result, err
	}

	result.Configuration = conf
	return result, nil
}

func resourceComponentInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.ServiceStageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage v2 client: %s", err)
	}

	appId := d.Get("application_id").(string)
	componentId := d.Get("component_id").(string)
	opt, err := buildInstanceUpdateOpts(d)
	if err != nil {
		return diag.Errorf("error building the UpdateOpts of the component instance: %s", err)
	}
	log.Printf("[DEBUG] The instance update option of ServiceStage component is: %v", opt)

	resp, err := instances.Update(client, appId, componentId, d.Id(), opt)
	if err != nil {
		return diag.Errorf("error updating component instance: %s", err)
	}

	log.Printf("[DEBUG] Waiting for the component instance to become running, the instance ID is %s.", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"RUNNING"},
		Target:       []string{"SUCCEEDED"},
		Refresh:      componentInstanceRefreshFunc(client, resp.JobId),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the updation of component instance (%s) to complete: %s",
			d.Id(), err)
	}

	return resourceComponentInstanceRead(ctx, d, meta)
}

func resourceComponentInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.ServiceStageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage v2 client: %s", err)
	}

	appId := d.Get("application_id").(string)
	componentId := d.Get("component_id").(string)
	resp, err := instances.Delete(client, appId, componentId, d.Id())
	if err != nil {
		return diag.Errorf("error deleting ServiceStage component instance (%s): %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Waiting for the component instance to become deleted, the instance ID is %s.", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"RUNNING"},
		Target:       []string{"SUCCEEDED"},
		Refresh:      componentInstanceRefreshFunc(client, resp.JobId),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the delete of component instance (%s) to complete: %s",
			d.Id(), err)
	}

	return nil
}

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

func resourceComponentInstanceImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid format specified for import id, must be " +
			"<application_id>/<component_id>/<instance_id>")
	}

	d.SetId(parts[2])
	mErr := multierror.Append(nil,
		d.Set("application_id", parts[0]),
		d.Set("component_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
