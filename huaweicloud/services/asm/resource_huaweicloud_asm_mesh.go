package asm

import (
	"context"
	"fmt"
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

// @API ASM POST /v1/{project_id}/meshes
// @API ASM GET /v1/{project_id}/meshes/{mesh_id}
// @API ASM DELETE /v1/{project_id}/meshes/{mesh_id}

var nonUpdatableParams = []string{
	"name", "type", "version", "annotations", "labels", "tags", "extend_params",
	"extend_params.*.clusters",
	"extend_params.*.clusters.*.cluster_id",
	"extend_params.*.clusters.*.installation",
	"extend_params.*.clusters.*.installation.*.nodes",
	"extend_params.*.clusters.*.installation.*.nodes.*.field_selector",
	"extend_params.*.clusters.*.installation.*.nodes.*.field_selector.*.key",
	"extend_params.*.clusters.*.installation.*.nodes.*.field_selector.*.operator",
	"extend_params.*.clusters.*.installation.*.nodes.*.field_selector.*.values",
	"extend_params.*.clusters.*.injection",
	"extend_params.*.clusters.*.injection.*.namespaces",
	"extend_params.*.clusters.*.injection.*.namespaces.*.field_selector",
	"extend_params.*.clusters.*.injection.*.namespaces.*.field_selector.*.key",
	"extend_params.*.clusters.*.injection.*.namespaces.*.field_selector.*.operator",
	"extend_params.*.clusters.*.injection.*.namespaces.*.field_selector.*.values",
}

func ResourceAsmMesh() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAsmMeshCreate,
		ReadContext:   resourceAsmMeshRead,
		UpdateContext: resourceAsmMeshUpdate,
		DeleteContext: resourceAsmMeshDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies mesh name.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the mesh type.`,
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the mesh version.`,
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the mesh annotations in key/value format.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the mesh labels in key/value format.`,
			},
			"extend_params": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `Specifies the extend parameters of the mesh.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"clusters": {
							Type:        schema.TypeList,
							Required:    true,
							Description: `Specifies the cluster informations in the mesh.`,
							Elem:        extendParamsClustersElem(),
						},
					},
				},
			},
			"tags": common.TagsSchema(),
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the mesh.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the mesh is created.`,
			},
		},
	}
}

func extendParamsClustersElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the cluster ID.`,
			},
			"installation": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `Specifies the mesh components installation configuration.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nodes": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: `Specifies the mesh components installation configuration.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_selector": {
										Type:        schema.TypeList,
										Required:    true,
										MaxItems:    1,
										Description: `Specifies the field selector.`,
										Elem:        selectorElem(),
									},
								},
							},
						},
					},
				},
			},
			"injection": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `Specifies the sidecar injection configuration.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespaces": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: `Specifies the namespace of the sidecar injection.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_selector": {
										Type:        schema.TypeList,
										Required:    true,
										MaxItems:    1,
										Description: `Specifies the field selector.`,
										Elem:        selectorElem(),
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

func selectorElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the key of the selector.`,
			},
			"operator": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the operator of the selector.`,
			},
			"values": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the value of the selector.`,
			},
		},
	}
}

func resourceAsmMeshCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createMeshHttpUrl = "v1/{project_id}/meshes"
		createMeshProduct = "asm"
	)
	createMeshClient, err := cfg.NewServiceClient(createMeshProduct, region)
	if err != nil {
		return diag.Errorf("error creating Mesh Client: %s", err)
	}

	createMeshPath := createMeshClient.Endpoint + createMeshHttpUrl
	createMeshPath = strings.ReplaceAll(createMeshPath, "{project_id}", createMeshClient.ProjectID)

	createMeshOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}

	createMeshOpt.JSONBody = utils.RemoveNil(buildCreateMeshBodyParams(d, cfg))
	createMeshResp, err := createMeshClient.Request("POST", createMeshPath, &createMeshOpt)
	if err != nil {
		return diag.Errorf("error creating Mesh: %s", err)
	}

	createMeshRespBody, err := utils.FlattenResponse(createMeshResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("metadata.uid", createMeshRespBody, nil)
	if id == nil {
		return diag.Errorf("error creating Mesh: ID is not found in API response")
	}
	d.SetId(id.(string))

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETE"},
		Refresh:      meshStateRefreshFunc(createMeshClient, d.Id(), []string{"Running"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        20 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the status of mesh (%s) to become running: %s", d.Id(), err)
	}

	return resourceAsmMeshRead(ctx, d, meta)
}

func buildCreateMeshBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "mesh",
		"metadata":   buildCreateMeshMetadataOpts(d),
		"spec":       buildCreateMeshSpecOpts(d, cfg),
	}
	return bodyParams
}

func buildCreateMeshMetadataOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"annotations": utils.ValueIgnoreEmpty(d.Get("annotations")),
		"labels":      utils.ValueIgnoreEmpty(d.Get("labels")),
	}
	return bodyParams
}

func buildCreateMeshSpecOpts(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"type":         d.Get("type"),
		"version":      d.Get("version"),
		"tags":         utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
		"extendParams": buildCreateMeshExtendParamsOpts(d, cfg),
	}
	return bodyParams
}

func buildCreateMeshExtendParamsOpts(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	extendParamsRaw := d.Get("extend_params").([]interface{})
	if len(extendParamsRaw) == 0 {
		return nil
	}
	extendParams := extendParamsRaw[0].(map[string]interface{})
	bodyParams := map[string]interface{}{
		"clusters": buildCreateMeshClustersOpts(extendParams["clusters"].([]interface{}), d, cfg),
	}
	return bodyParams
}

func buildCreateMeshClustersOpts(clustersRaw []interface{}, d *schema.ResourceData, cfg *config.Config) []map[string]interface{} {
	if len(clustersRaw) == 0 {
		return nil
	}

	region := cfg.GetRegion(d)

	bodyParams := make([]map[string]interface{}, len(clustersRaw))
	for i, clusterRaw := range clustersRaw {
		cluster := clusterRaw.(map[string]interface{})
		bodyParams[i] = map[string]interface{}{
			"clusterID":    cluster["cluster_id"],
			"projectID":    cfg.GetProjectID(region),
			"injection":    buildCreateMeshInjectionOpts(cluster["injection"].([]interface{})),
			"installation": buildCreateMeshInstallationOpts(cluster["installation"].([]interface{})),
		}
	}
	return bodyParams
}

func buildCreateMeshInjectionOpts(injectionRaw []interface{}) map[string]interface{} {
	if len(injectionRaw) == 0 {
		return nil
	}

	injection := injectionRaw[0].(map[string]interface{})
	bodyParams := map[string]interface{}{
		"namespaces": buildCreateMeshNamespacesOpts(injection["namespaces"].([]interface{})),
	}
	return bodyParams
}

func buildCreateMeshNamespacesOpts(namespacesRaw []interface{}) map[string]interface{} {
	if len(namespacesRaw) == 0 {
		return nil
	}

	namespaces := namespacesRaw[0].(map[string]interface{})
	bodyParams := map[string]interface{}{
		"fieldSelector": buildCreateMeshFieldSelectorOpts(namespaces["field_selector"].([]interface{})),
	}
	return bodyParams
}

func buildCreateMeshInstallationOpts(installationRaw []interface{}) map[string]interface{} {
	if len(installationRaw) == 0 {
		return nil
	}

	installation := installationRaw[0].(map[string]interface{})
	bodyParams := map[string]interface{}{
		"nodes": buildCreateMeshNodesOpts(installation["nodes"].([]interface{})),
	}
	return bodyParams
}

func buildCreateMeshNodesOpts(nodesRaw []interface{}) map[string]interface{} {
	if len(nodesRaw) == 0 {
		return nil
	}

	nodes := nodesRaw[0].(map[string]interface{})
	bodyParams := map[string]interface{}{
		"fieldSelector": buildCreateMeshFieldSelectorOpts(nodes["field_selector"].([]interface{})),
	}
	return bodyParams
}

func buildCreateMeshFieldSelectorOpts(fieldSelectorRaw []interface{}) map[string]interface{} {
	if len(fieldSelectorRaw) == 0 {
		return nil
	}

	fieldSelector := fieldSelectorRaw[0].(map[string]interface{})
	bodyParams := map[string]interface{}{
		"key":      fieldSelector["key"],
		"operator": fieldSelector["operator"],
		"values":   utils.ExpandToStringList(fieldSelector["values"].([]interface{})),
	}
	return bodyParams
}

func resourceAsmMeshRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getMeshHttpUrl = "v1/{project_id}/meshes/{mesh_id}"
		getMeshProduct = "asm"
	)
	getMeshClient, err := cfg.NewServiceClient(getMeshProduct, region)
	if err != nil {
		return diag.Errorf("error creating Mesh Client: %s", err)
	}

	getMeshPath := getMeshClient.Endpoint + getMeshHttpUrl
	getMeshPath = strings.ReplaceAll(getMeshPath, "{project_id}", getMeshClient.ProjectID)
	getMeshPath = strings.ReplaceAll(getMeshPath, "{mesh_id}", d.Id())

	getPotectionRulesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getMeshResp, err := getMeshClient.Request("GET", getMeshPath, &getPotectionRulesOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ASM mesh")
	}

	getMeshRespBody, err := utils.FlattenResponse(getMeshResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// extend_params not set, because extend_params.*.clusters.*.installation and
	// extend_params.*.clusters.*.injection not returned by API
	// annotations, labels, tags not set, because they are not returned by API
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("metadata.name", getMeshRespBody, nil)),
		d.Set("created_at", utils.PathSearch("metadata.creationTimestamp", getMeshRespBody, nil)),
		d.Set("type", utils.PathSearch("spec.type", getMeshRespBody, nil)),
		d.Set("version", utils.PathSearch("spec.version", getMeshRespBody, nil)),
		d.Set("status", utils.PathSearch("status.phase", getMeshRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAsmMeshUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAsmMeshDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteMeshHttpUrl = "v1/{project_id}/meshes/{mesh_id}"
		deleteMeshProduct = "asm"
	)
	deleteMeshClient, err := cfg.NewServiceClient(deleteMeshProduct, region)
	if err != nil {
		return diag.Errorf("error creating Mesh Client: %s", err)
	}

	deleteMeshPath := deleteMeshClient.Endpoint + deleteMeshHttpUrl
	deleteMeshPath = strings.ReplaceAll(deleteMeshPath, "{project_id}", cfg.GetProjectID(region))
	deleteMeshPath = strings.ReplaceAll(deleteMeshPath, "{mesh_id}", d.Id())

	deleteMeshOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	_, err = deleteMeshClient.Request("DELETE", deleteMeshPath, &deleteMeshOpt)
	if err != nil {
		return diag.Errorf("error deleting Mesh: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETE"},
		Refresh:      meshStateRefreshFunc(deleteMeshClient, d.Id(), []string{"Deleted"}),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        20 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the mesh (%s) to be deleted: %s", d.Id(), err)
	}

	return nil
}

func meshStateRefreshFunc(client *golangsdk.ServiceClient, meshId string, targetStatus []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			getMeshHttpUrl  = "v1/{project_id}/meshes/{mesh_id}"
			status          string
			getMeshRespBody interface{}
		)

		getMeshPath := client.Endpoint + getMeshHttpUrl
		getMeshPath = strings.ReplaceAll(getMeshPath, "{project_id}", client.ProjectID)
		getMeshPath = strings.ReplaceAll(getMeshPath, "{mesh_id}", meshId)

		getPotectionRulesOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		getMeshResp, err := client.Request("GET", getMeshPath, &getPotectionRulesOpt)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); !ok {
				return getMeshResp, "ERROR", fmt.Errorf("error retrieving ASM mesh: %s", err)
			}
			status = "Deleted"
		} else {
			getMeshRespBody, err = utils.FlattenResponse(getMeshResp)
			if err != nil {
				return nil, "ERROR", err
			}

			statusRaw := utils.PathSearch(`status.phase`, getMeshRespBody, nil)
			if statusRaw == nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `status.phase`)
			}

			status = fmt.Sprintf("%v", statusRaw)
		}

		unexpecetdStatus := []string{
			"CreateFailed", "DeleteFailed", "UpgradeFailed", "RollbackFailed",
		}
		if utils.StrSliceContains(unexpecetdStatus, status) {
			return getMeshResp, "ERROR", fmt.Errorf("unexpect status (%s)", status)
		}

		if utils.StrSliceContains(targetStatus, status) {
			return getMeshResp, "COMPLETE", nil
		}

		return getMeshRespBody, "PENDING", nil
	}
}
