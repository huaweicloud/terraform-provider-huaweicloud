package servicestage

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/servicestage/v2/components"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API ServiceStage POST /v2/{project_id}/cas/applications/{application_id}/components
// @API ServiceStage GET /v2/{project_id}/cas/applications/{application_id}/components/{component_id}
// @API ServiceStage PUT /v2/{project_id}/cas/applications/{application_id}/components/{component_id}
// @API ServiceStage DELETE /v2/{project_id}/cas/applications/{application_id}/components/{component_id}
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Webapp", "MicroService", "Common",
				}, false),
			},
			"runtime": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"framework": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
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
						"authorization": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ExactlyOneOf: []string{"source.0.storage_type"},
						},
						"repo_ref": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"source.0.storage_type"},
						},
						"repo_namespace": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"source.0.storage_type"},
						},
						"storage_type": {
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
									"endpoint": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"bucket": {
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
			"builder": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"organization": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"cmd": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"dockerfile_path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"use_public_cluster": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"node_label": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func buildPropertiesStructure(params []interface{}) components.Properties {
	if len(params) < 1 {
		return components.Properties{}
	}

	param := params[0].(map[string]interface{})
	return components.Properties{
		Endpoint: param["endpoint"].(string),
		Bucket:   param["bucket"].(string),
		Key:      param["key"].(string),
	}
}

func buildRepoBuilderStructure(params []interface{}) *components.Builder {
	if len(params) < 1 {
		return nil
	}

	param := params[0].(map[string]interface{})

	return &components.Builder{
		Parameter: components.Parameter{
			BuildCmd:          param["cmd"].(string),
			ArtifactNamespace: param["organization"].(string),
			ClusterId:         param["cluster_id"].(string),
			ClusterName:       param["cluster_name"].(string),
			ClusterType:       param["cluster_type"].(string),
			UsePublicCluster:  param["use_public_cluster"].(bool),
			DockerfilePath:    param["dockerfile_path"].(string),
			NodeLabelSelector: param["node_label"].(map[string]interface{}),
		},
	}
}

func buildRepoSourceStructure(sources []interface{}) *components.Source {
	if len(sources) < 1 {
		return nil
	}
	var result components.Source

	source := sources[0].(map[string]interface{})
	rType := source["type"].(string)
	switch rType {
	case "package":
		result = components.Source{
			Kind: "artifact",
			Spec: components.Spec{
				Type:       rType,
				Storage:    source["storage_type"].(string),
				Url:        source["url"].(string),
				Properties: buildPropertiesStructure(source["properties"].([]interface{})),
			},
		}
	case "GitHub", "GitLab", "Gitee", "Bitbucket", "DevCloud":
		result = components.Source{
			Kind: "code",
			Spec: components.Spec{
				RepoType:      rType,
				RepoAuth:      source["authorization"].(string),
				RepoUrl:       source["url"].(string),
				RepoRef:       source["repo_ref"].(string),
				RepoNamespace: source["repo_namespace"].(string),
			},
		}
	default:
	}
	return &result
}

func resourceComponentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	v2Client, err := conf.ServiceStageV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ServiceStage V2 client: %s", err)
	}

	appId := d.Get("application_id").(string)
	opt := components.CreateOpts{
		Name:     d.Get("name").(string),
		Runtime:  d.Get("runtime").(string),
		Type:     d.Get("type").(string),
		Framwork: d.Get("framework").(string),
		Builder:  buildRepoBuilderStructure(d.Get("builder").([]interface{})),
		Source:   buildRepoSourceStructure(d.Get("source").([]interface{})),
	}
	resp, err := components.Create(v2Client, appId, opt)
	if err != nil {
		return diag.Errorf("error creating ServiceStage component: %s", err)
	}

	d.SetId(resp.ID)

	return resourceComponentRead(ctx, d, meta)
}

func flattenRepoBuilder(builder components.Builder) (result []map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Recover panic when flattening builder structure: %#v", r)
		}
	}()

	if !reflect.DeepEqual(builder, components.Builder{}) {
		result = append(result, map[string]interface{}{
			"cmd":                builder.Parameter.BuildCmd,
			"organization":       builder.Parameter.ArtifactNamespace,
			"cluster_id":         builder.Parameter.ClusterId,
			"cluster_name":       builder.Parameter.ClusterName,
			"cluster_type":       builder.Parameter.ClusterType,
			"dockerfile_path":    builder.Parameter.DockerfilePath,
			"use_public_cluster": builder.Parameter.UsePublicCluster,
			"node_label":         builder.Parameter.NodeLabelSelector,
		})
	}

	return
}

func flattenRepoSource(source components.Source) (result []map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Recover panic when flattening source structure: %#v", r)
		}
	}()

	if (source != components.Source{}) {
		if source.Spec.Type == "package" {
			result = append(result, map[string]interface{}{
				"type":         source.Spec.Type,
				"storage_type": source.Spec.Storage,
				"url":          source.Spec.Url,
			})
		} else if source.Spec.RepoType == "GitHub" || source.Spec.Type == "GitLab" ||
			source.Spec.Type == "Gitee" || source.Spec.Type == "Bitbucket" || source.Spec.Type == "DevCloud" {
			result = append(result, map[string]interface{}{
				"type":           source.Spec.RepoType,
				"authorization":  source.Spec.RepoAuth,
				"url":            source.Spec.RepoUrl,
				"repo_ref":       source.Spec.RepoRef,
				"repo_namespace": source.Spec.RepoNamespace,
			})
		}
	}

	return
}

func resourceComponentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.ServiceStageV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage V2 client: %s", err)
	}

	appId := d.Get("application_id").(string)
	resp, err := components.Get(client, appId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ServiceStage component")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", resp.Name),
		d.Set("type", resp.Type),
		d.Set("runtime", resp.Runtime),
		d.Set("framework", resp.Framwork),
		d.Set("builder", flattenRepoBuilder(resp.Builder)),
		d.Set("source", flattenRepoSource(resp.Source)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceComponentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.ServiceStageV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ServiceStage V2 client: %s", err)
	}

	appId := d.Get("application_id").(string)
	// In normal changes, there is no situation in which source and builder are empty, so these two empty values are
	// ignored.
	opt := components.UpdateOpts{
		Name:    d.Get("name").(string),
		Builder: buildRepoBuilderStructure(d.Get("builder").([]interface{})),
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
	client, err := conf.ServiceStageV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ServiceStage V2 client: %s", err)
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
