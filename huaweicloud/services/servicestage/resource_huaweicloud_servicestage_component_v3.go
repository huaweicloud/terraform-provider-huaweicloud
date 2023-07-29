package servicestage

import (
    "context"
    "fmt"
    "log"
    "reflect"
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

func ResourceComponentV3() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceComponentV3Create,
        ReadContext:   resourceComponentV3Read,
        UpdateContext: resourceComponentV3Update,
        DeleteContext: resourceComponentV3Delete,

        Importer: &schema.ResourceImporter{
            StateContext: resourceComponentV3ImportState,
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
                Type:     schema.TypeString,
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
                Type:     schema.TypeString,
                Optional: true,
            },
            "jvm_opts": {
                Type:     schema.TypeString,
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
                            Type:     schema.TypeSet,
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
                                        Type:     schema.TypeBool,
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
                            Type:     schema.TypeString,
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
                            Type:     schema.TypeList,
                            Optional: true,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "name": {
                                        Type:     schema.TypeString,
                                        Optional: true,
                                    },
                                    "value": {
                                        Type:     schema.TypeString,
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
                        "az": {
                            Type:     schema.TypeList,
                            Optional: true,
                            Elem:     &schema.Schema{Type: schema.TypeString},
                        },
                        "node": {
                            Type:     schema.TypeList,
                            Optional: true,
                            Elem:     &schema.Schema{Type: schema.TypeString},
                        },
                        "component": {
                            Type:     schema.TypeList,
                            Optional: true,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "display_name": {
                                        Type:     schema.TypeString,
                                        Optional: true,
                                    },
                                    "name": {
                                        Type:     schema.TypeString,
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
                        "az": {
                            Type:     schema.TypeList,
                            Optional: true,
                            Elem:     &schema.Schema{Type: schema.TypeString},
                        },
                        "node": {
                            Type:     schema.TypeList,
                            Optional: true,
                            Elem:     &schema.Schema{Type: schema.TypeString},
                        },
                        "component": {
                            Type:     schema.TypeList,
                            Optional: true,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "display_name": {
                                        Type:     schema.TypeString,
                                        Optional: true,
                                    },
                                    "name": {
                                        Type:     schema.TypeString,
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
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "type": {
                                        Type:     schema.TypeString,
                                        Optional: true,
                                    },
                                    "namespace": {
                                        Type:     schema.TypeString,
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

func resourceComponentV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
        Source:                 buildRepoSourceV3Structure(d.Get("source").(interface{})),
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
        Refresh:      componentInstanceV3RefreshFunc(client, resp.JobId),
        Timeout:      d.Timeout(schema.TimeoutCreate),
        Delay:        5 * time.Second,
        PollInterval: 5 * time.Second,
    }
    _, err = stateConf.WaitForStateContext(ctx)
    if err != nil {
        return diag.Errorf("error waiting for the creation of component (%s) to complete: %s",
            d.Id(), err)
    }

    return resourceComponentV3Read(ctx, d, meta)
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

func buildRepoSourceV3Structure(sources interface{}) *components.Source {
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

    azList := a["az"].([]interface{})
    azValue := make([]string, len(azList))
    for i, raw := range azList {
        azValue[i] = raw.(string)
    }
    conf.AZ = azValue

    nodeList := a["az"].([]interface{})
    nodeValue := make([]string, len(nodeList))
    for i, raw := range nodeList {
        nodeValue[i] = raw.(string)
    }
    conf.Node = nodeValue
    conf.Component = buildAffinityAppInnerParameters(a["component"].([]interface{}))

    return conf
}

func buildAffinityAppInnerParameters(appInnerParameters []interface{}) []*components.AppInnerParameters {
    if len(appInnerParameters) < 1 {
        return nil
    }
    var result []*components.AppInnerParameters
    for _, parameter := range appInnerParameters {
        var environmentLabel components.AppInnerParameters
        p := parameter.(map[string]interface{})
        environmentLabel.DisplayName = p["display_name"].(string)
        environmentLabel.Name = p["name"].(string)
        result = append(result, &environmentLabel)
    }

    return result
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

func componentInstanceV3RefreshFunc(c *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
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

func resourceComponentV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
        d.Set("name", resp.Name),
        d.Set("workload_name", resp.WorkloadName),
        d.Set("description", resp.Description),
        d.Set("labels", flattenComponentLabels(resp.Labels)),
        d.Set("pod_labels", flattenComponentLabels(resp.PodLabels)),
        d.Set("version", resp.Version),
        d.Set("description", resp.Description),
        d.Set("environment_id", resp.EnvironmentID),
        d.Set("application_id", resp.ApplicationID),
        d.Set("enterprise_project_id", resp.EnterpriseProjectId),
        d.Set("limit_cpu", resp.LimitCpu),
        d.Set("limit_memory", resp.LimitMemory),
        d.Set("request_cpu", resp.RequestCpu),
        d.Set("request_memory", resp.RequestMemory),
        d.Set("replica", resp.Replica),
        d.Set("timezone", resp.Timezone),
        d.Set("jvm_opts", resp.JvmOpts),
        d.Set("workload_kind", resp.WorkloadKind),
        d.Set("runtime_stack", flattenComponentRuntimeStack(resp.RuntimeStack)),
        d.Set("build", flattenComponentBuild(resp.Build)),
        d.Set("source", flattenComponentSource(resp.Source)),
        d.Set("envs", flattenComponentEnvs(resp.Envs)),
        d.Set("storage", flattenComponentStorages(resp.Storage)),
        d.Set("command", flattenComponentCommand(resp.Command)),
        d.Set("post_start", flattenComponentK8sLifeCycle(resp.PostStart)),
        d.Set("pre_stop", flattenComponentK8sLifeCycle(resp.PreStop)),
        d.Set("mesher", flattenComponentMesher(resp.Mesher)),
        d.Set("tomcat_opts", flattenComponentTomcatOpt(resp.TomcatOpts)),
        d.Set("dns_policy", resp.DnsPolicy),
        d.Set("dns_config", flattenComponentDnsConfig(resp.DnsConfig)),
        d.Set("security_context", flattenComponentSecurityContext(resp.SecurityContext)),
        d.Set("logs", flattenComponentLogs(resp.Logs)),
        d.Set("custom_metric", flattenComponentCustomMetric(resp.CustomMetric)),
        d.Set("affinity", flattenComponentAffinity(resp.Affinity)),
        d.Set("anti_affinity", flattenComponentAffinity(resp.AntiAffinity)),
        d.Set("liveness_probe", flattenComponentK8sProbe(resp.LivenessProbe)),
        d.Set("readiness_probe", flattenComponentK8sProbe(resp.ReadinessProbe)),
        d.Set("refer_resources", flattenComponentReferResource(resp.ReferResources)),

    )

    return diag.FromErr(mErr.ErrorOrNil())
}

func flattenComponentLabels(labels []*components.KeyValue) (result []map[string]interface{}) {
    if len(labels) < 1 {
        return nil
    }

    for _, val := range labels {
        s := map[string]interface{}{
            "key":   val.Key,
            "value": val.Value,
        }
        result = append(result, s)
    }

    log.Printf("[DEBUG] The labels result is %#v", result)
    return result
}

func flattenComponentRuntimeStack(runtimeStack *components.RuntimeStack) (result []map[string]interface{}) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("[ERROR] Recover panic when flattening scheduler structure: %#v", r)
        }
    }()

    if reflect.DeepEqual(runtimeStack, components.RuntimeStack{}) {
        return nil
    }

    result = []map[string]interface{}{
        {
            "name":        runtimeStack.Name,
            "type":        runtimeStack.Type,
            "version":     runtimeStack.Version,
            "deploy_mode": runtimeStack.DeployMode,
        },
    }
    log.Printf("[DEBUG] The runtimeStack result is %#v", result)
    return result

}

func flattenComponentBuild(build *components.Build) (result []map[string]interface{}) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("[ERROR] Recover panic when flattening scheduler structure: %#v", r)
        }
    }()

    if reflect.DeepEqual(build, components.Build{}) {
        return nil
    }

    result = []map[string]interface{}{
        {
            "parameters": flattenComponentBuildParameters(build.Parameter),
        },
    }
    log.Printf("[DEBUG] The build result is %#v", result)
    return result

}

func flattenComponentBuildParameters(parameters components.Parameter) (result []map[string]interface{}) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("[ERROR] Recover panic when flattening scheduler structure: %#v", r)
        }
    }()

    if reflect.DeepEqual(parameters, components.Parameter{}) {
        return nil
    }

    result = []map[string]interface{}{
        {
            "build_cmd":           parameters.BuildCmd,
            "dockerfile_path":     parameters.DockerfilePath,
            "artifact_namespace":  parameters.ArtifactNamespace,
            "cluster_id":          parameters.ClusterId,
            "environment_id":      parameters.EnvironmentId,
            "node_label_selector": parameters.NodeLabelSelector,
        },
    }
    log.Printf("[DEBUG] The build.parameter result is %#v", result)
    return result

}

func flattenComponentSource(source *components.Source) (result []map[string]interface{}) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("[ERROR] Recover panic when flattening scheduler structure: %#v", r)
        }
    }()

    if reflect.DeepEqual(source, components.Source{}) {
        return nil
    }

    result = []map[string]interface{}{
        {
            "kind":           source.Kind,
            "version":        source.Version,
            "url":            source.Url,
            "storage":        source.Storage,
            "auth":           source.Auth,
            "repo_auth":      source.RepoAuth,
            "repo_namespace": source.RepoNamespace,
            "repo_ref":       source.RepoRef,
            "web_url":        source.WebUrl,
            "repo_url":       source.RepoUrl,
        },
    }
    log.Printf("[DEBUG] The source result is %#v", result)
    return result

}

func flattenComponentEnvs(envs []*components.Env) (result []map[string]interface{}) {
    if len(envs) < 1 {
        return nil
    }

    for _, val := range envs {
        s := map[string]interface{}{
            "name":       val.Name,
            "value":      val.Value,
            "value_from": val.EnvValueFrom,
        }
        result = append(result, s)
    }

    log.Printf("[DEBUG] The envs result is %#v", result)
    return result
}

func flattenComponentStorages(storages []*components.Storage) (result []map[string]interface{}) {
    if len(storages) < 1 {
        return nil
    }

    for _, val := range storages {
        s := map[string]interface{}{
            "name":       val.Name,
            "type":       val.Type,
            "parameters": flattenComponentStoragesParameters(val.Parameters),
            "mounts":     flattenComponentStoragesMounts(val.Mounts),
        }
        result = append(result, s)
    }

    log.Printf("[DEBUG] The storages result is %#v", result)
    return result
}

func flattenComponentStoragesParameters(parameters *components.StorageParameter) (result []map[string]interface{}) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("[ERROR] Recover panic when flattening scheduler structure: %#v", r)
        }
    }()

    if reflect.DeepEqual(parameters, components.StorageParameter{}) {
        return nil
    }

    result = []map[string]interface{}{
        {
            "path":         parameters.Path,
            "name":         parameters.Name,
            "default_mode": parameters.DefaultMode,
            "medium":       parameters.Medium,
        },
    }
    log.Printf("[DEBUG] The StoragesParameters result is %#v", result)
    return result

}

func flattenComponentStoragesMounts(mounts []*components.StorageMounts) (result []map[string]interface{}) {
    if len(mounts) < 1 {
        return nil
    }

    for _, val := range mounts {
        s := map[string]interface{}{
            "path":      val.Path,
            "sub_path":  val.SubPath,
            "read_only": val.Readonly,
        }
        result = append(result, s)
    }

    log.Printf("[DEBUG] The StoragesMounts result is %#v", result)
    return result
}

func flattenComponentCommand(command *components.Command) (result []map[string]interface{}) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("[ERROR] Recover panic when flattening scheduler structure: %#v", r)
        }
    }()

    if reflect.DeepEqual(command, components.Command{}) {
        return nil
    }

    result = []map[string]interface{}{
        {
            "command": command.Command,
            "args":    command.Args,
        },
    }
    log.Printf("[DEBUG] The command result is %#v", result)
    return result

}

func flattenComponentK8sLifeCycle(k8sLifeCycle *components.K8sLifeCycle) (result []map[string]interface{}) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("[ERROR] Recover panic when flattening scheduler structure: %#v", r)
        }
    }()

    if reflect.DeepEqual(k8sLifeCycle, components.K8sLifeCycle{}) {
        return nil
    }

    result = []map[string]interface{}{
        {
            "command": k8sLifeCycle.Command,
            "path":    k8sLifeCycle.Path,
            "port":    k8sLifeCycle.Port,
            "scheme":  k8sLifeCycle.Scheme,
            "type":    k8sLifeCycle.Type,
        },
    }
    log.Printf("[DEBUG] The k8sLifeCycle result is %#v", result)
    return result
}

func flattenComponentMesher(mesher *components.Mesher) (result []map[string]interface{}) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("[ERROR] Recover panic when flattening scheduler structure: %#v", r)
        }
    }()

    if reflect.DeepEqual(mesher, components.Mesher{}) {
        return nil
    }

    result = []map[string]interface{}{
        {
            "port": mesher.Port,
        },
    }
    log.Printf("[DEBUG] The k8sLifeCycle result is %#v", result)
    return result
}

func flattenComponentTomcatOpt(tomcatOpts *components.TomcatOpts) (result []map[string]interface{}) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("[ERROR] Recover panic when flattening scheduler structure: %#v", r)
        }
    }()

    if reflect.DeepEqual(tomcatOpts, components.TomcatOpts{}) {
        return nil
    }

    result = []map[string]interface{}{
        {
            "server_xml": tomcatOpts.ServerXml,
        },
    }
    log.Printf("[DEBUG] The tomcatOpts result is %#v", result)
    return result
}

func flattenComponentHostAliases(hostAliases []*components.HostAlias) (result []map[string]interface{}) {
    if len(hostAliases) < 1 {
        return nil
    }

    for _, val := range hostAliases {
        s := map[string]interface{}{
            "ip":       val.IP,
            "hostname": val.HostNames,
        }
        result = append(result, s)
    }

    log.Printf("[DEBUG] The hostAliases result is %#v", result)
    return result
}

func flattenComponentDnsConfig(dnsConfig *components.DnsConfig) (result []map[string]interface{}) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("[ERROR] Recover panic when flattening scheduler structure: %#v", r)
        }
    }()

    if reflect.DeepEqual(dnsConfig, components.DnsConfig{}) {
        return nil
    }

    result = []map[string]interface{}{
        {
            "nameservers": dnsConfig.Nameservers,
            "searches":    dnsConfig.Searches,
            "options":     flattenComponentDnsConfigOption(dnsConfig.Options),
        },
    }
    log.Printf("[DEBUG] The tomcatOpts result is %#v", result)
    return result
}

func flattenComponentDnsConfigOption(dnsConfigOption []*components.NameValue) (result []map[string]interface{}) {
    if len(dnsConfigOption) < 1 {
        return nil
    }

    for _, val := range dnsConfigOption {
        s := map[string]interface{}{
            "name":  val.Name,
            "value": val.Value,
        }
        result = append(result, s)
    }

    log.Printf("[DEBUG] The dnsConfigOption result is %#v", result)
    return result
}

func flattenComponentSecurityContext(securityContext *components.SecurityContext) (result []map[string]interface{}) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("[ERROR] Recover panic when flattening scheduler structure: %#v", r)
        }
    }()

    if reflect.DeepEqual(securityContext, components.SecurityContext{}) {
        return nil
    }

    result = []map[string]interface{}{
        {
            "run_as_user":  securityContext.RunAsUser,
            "run_as_group": securityContext.RunAsGroup,
            "capabilities": flattenComponentCapabilities(securityContext.Capabilities),
        },
    }
    log.Printf("[DEBUG] The securityContext result is %#v", result)
    return result
}

func flattenComponentCapabilities(capabilities *components.Capabilities) (result []map[string]interface{}) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("[ERROR] Recover panic when flattening scheduler structure: %#v", r)
        }
    }()

    if reflect.DeepEqual(capabilities, components.Capabilities{}) {
        return nil
    }

    result = []map[string]interface{}{
        {
            "add":  capabilities.Add,
            "drop": capabilities.Drop,
        },
    }
    log.Printf("[DEBUG] The capabilities result is %#v", result)
    return result
}

func flattenComponentLogs(Logs []*components.Log) (result []map[string]interface{}) {
    if len(Logs) < 1 {
        return nil
    }

    for _, val := range Logs {
        s := map[string]interface{}{
            "log_path":         val.LogPath,
            "rotate":           val.Rotate,
            "host_path":        val.HostPath,
            "host_extend_path": val.HostExtendPath,
        }
        result = append(result, s)
    }

    log.Printf("[DEBUG] The logs result is %#v", result)
    return result
}

func flattenComponentCustomMetric(customMetric *components.CustomMetric) (result []map[string]interface{}) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("[ERROR] Recover panic when flattening scheduler structure: %#v", r)
        }
    }()

    if reflect.DeepEqual(customMetric, components.CustomMetric{}) {
        return nil
    }

    result = []map[string]interface{}{
        {
            "path":       customMetric.Path,
            "port":       customMetric.Port,
            "dimensions": customMetric.Dimensions,
        },
    }
    log.Printf("[DEBUG] The customMetric result is %#v", result)
    return result
}

func flattenComponentAffinity(affinity *components.Affinity) (result []map[string]interface{}) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("[ERROR] Recover panic when flattening scheduler structure: %#v", r)
        }
    }()

    if reflect.DeepEqual(affinity, components.Affinity{}) {
        return nil
    }

    result = []map[string]interface{}{
        {
            "az":        affinity.AZ,
            "node":      affinity.Node,
            "component": flattenComponentAffinityComponent(affinity.Component),
        },
    }
    log.Printf("[DEBUG] The affinity result is %#v", result)
    return result
}

func flattenComponentAffinityComponent(affinityComponent []*components.AppInnerParameters) (result []map[string]interface{}) {
    if len(affinityComponent) < 1 {
        return nil
    }

    for _, val := range affinityComponent {
        s := map[string]interface{}{
            "display_name": val.DisplayName,
            "name":         val.Name,
        }
        result = append(result, s)
    }

    log.Printf("[DEBUG] The affinityComponent result is %#v", result)
    return result
}

func flattenComponentK8sProbe(k8sProbe *components.K8sProbe) (result []map[string]interface{}) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("[ERROR] Recover panic when flattening scheduler structure: %#v", r)
        }
    }()

    if reflect.DeepEqual(k8sProbe, components.K8sProbe{}) {
        return nil
    }

    result = []map[string]interface{}{
        {
            "command":           k8sProbe.Command,
            "path":              k8sProbe.Path,
            "port":              k8sProbe.Port,
            "host":              k8sProbe.Host,
            "scheme":            k8sProbe.Scheme,
            "type":              k8sProbe.Type,
            "delay":             k8sProbe.Delay,
            "timeout":           k8sProbe.Timeout,
            "period_seconds":    k8sProbe.PeriodSeconds,
            "success_threshold": k8sProbe.SuccessThreshold,
            "failure_threshold": k8sProbe.FailureThreshold,
        },
    }
    log.Printf("[DEBUG] The k8sProbe result is %#v", result)
    return result
}

func flattenComponentReferResource(referResource []*components.Resource) (result []map[string]interface{}) {
    if len(referResource) < 1 {
        return nil
    }

    for _, val := range referResource {
        s := map[string]interface{}{
            "id":         val.ID,
            "type":       val.Type,
            "parameters": flattenComponentReferResourceParameter(val.Parameters),
        }
        result = append(result, s)
    }

    log.Printf("[DEBUG] The referResource result is %#v", result)
    return result
}

func flattenComponentReferResourceParameter(referResourceParameter *components.ResourceParameters) (result []map[string]interface{}) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("[ERROR] Recover panic when flattening scheduler structure: %#v", r)
        }
    }()

    if reflect.DeepEqual(referResourceParameter, components.ResourceParameters{}) {
        return nil
    }

    result = []map[string]interface{}{
        {
            "type":      referResourceParameter.Type,
            "namespace": referResourceParameter.NameSpace,
        },
    }
    log.Printf("[DEBUG] The referResourceParameter result is %#v", result)
    return result
}

func resourceComponentV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
    conf := meta.(*config.Config)
    client, err := conf.ServiceStageV3Client(conf.GetRegion(d))
    if err != nil {
        return diag.Errorf("error creating ServiceStage V3 client: %s", err)
    }

    appId := d.Get("application_id").(string)
    // In normal changes, there is no situation in which source and builder are empty, so these two empty values are
    // ignored.
    opt := components.UpdateOpts{
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
        Source:                 buildRepoSourceV3Structure(d.Get("source").(interface{})),
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
    _, err = components.Update(client, appId, d.Id(), opt)
    if err != nil {
        return diag.Errorf("error updating ServiceStage component (%s): %s", d.Id(), err)
    }

    return resourceComponentV3Read(ctx, d, meta)
}

func resourceComponentV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceComponentV3ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
    parts := strings.SplitN(d.Id(), "/", 2)
    if len(parts) != 2 {
        return nil, fmt.Errorf("Invalid format specified for import id, must be <application_id>/<component_id>")
    }

    d.SetId(parts[1])
    return []*schema.ResourceData{d}, d.Set("application_id", parts[0])
}
