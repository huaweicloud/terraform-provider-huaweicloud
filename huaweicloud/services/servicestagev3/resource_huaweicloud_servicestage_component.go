package servicestagev3

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z]([\w-]*[A-Za-z0-9])?$`),
						"The name can only contain letters, digits, underscores (_) and hyphens (-), and the name must"+
							" start with a letter and end with a letter or digit."),
					validation.StringLenBetween(2, 64),
				),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
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
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"limit_cpu": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit_memory": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"request_cpu": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"request_memory": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"tomcat_opts": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_xml": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"jvm_opts": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"refer_resources": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"replica": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"external_accesses": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
						"codearts_project_id": {
							Type:     schema.TypeString,
							Optional: true,
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
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Labels:              buildCompLabels(d.Get("labels").([]interface{})),
		Version:             d.Get("version").(string),
		EnvironmentID:       d.Get("environment_id").(string),
		ApplicationID:       d.Get("application_id").(string),
		EnterpriseProjectId: d.Get("enterprise_project_id").(string),
		RuntimeStack:        buildCompRuntimeStack(d.Get("runtime_stack").(interface{})),
		Build:               buildRepoBuildStructure(d.Get("build").(interface{})),
		Source:              buildRepoSourceStructure(d.Get("source").(interface{})),
		ReferResources:      buildReferResourcesStructure(d.Get("refer_resources").([]interface{})),
	}
	resp, err := components.Create(client, appId, opt)
	if err != nil {
		return diag.Errorf("error creating ServiceStage component: %s", err)
	}

	d.SetId(resp.ID)

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
		result = append(result, &referResource)
	}

	return result
}

func buildRepoSourceStructure(sources interface{}) *components.Source {
	sourcesSet := sources.(*schema.Set)

	if sourcesSet.Len() != 1 {
		return &components.Source{}
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

//func flattenRepoBuilder(builder components.Builder) (result []map[string]interface{}) {
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
//}
//
//func flattenRepoSource(source components.Source) (result []map[string]interface{}) {
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
//}

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
