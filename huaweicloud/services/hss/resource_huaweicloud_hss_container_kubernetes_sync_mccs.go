package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var syncClusterStatusNonUpdatableParams = []string{
	"total_num", "data_list", "data_list.*.cluster_id", "enterprise_project_id",
}

// @API HSS POST /v5/{project_id}/container/kubernetes/multi-cloud/clusters/status-synchronize
func ResourceContainerKubernetesSyncMccs() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContainerKubernetesSyncMccsCreate,
		ReadContext:   resourceContainerKubernetesSyncMccsRead,
		UpdateContext: resourceContainerKubernetesSyncMccsUpdate,
		DeleteContext: resourceContainerKubernetesSyncMccsDelete,

		CustomizeDiff: config.FlexibleForceNew(syncClusterStatusNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Specifies the region where the resource is located. If omitted, the provider-level region will be used.",
			},
			"total_num": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The total number of clusters to synchronize.",
			},
			// Fields "data_list" and "data_list.*.cluster_id" are optional in the API documentation, but are actually required.
			"data_list": {
				Type:     schema.TypeList,
				Optional: true,
				Description: utils.SchemaDesc("The list of cluster IDs to synchronize.",
					utils.SchemaDescInput{
						Required: true,
					}),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:     schema.TypeString,
							Required: true,
							Description: utils.SchemaDesc("The ID of the cluster to synchronize.",
								utils.SchemaDescInput{
									Required: true,
								}),
						},
					},
				},
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the enterprise project ID.",
			},
		},
	}
}

func buildContainerKubernetesSyncMccsQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	rst := ""
	if epsID := cfg.GetEnterpriseProjectID(d); epsID != "" {
		rst = fmt.Sprintf("?enterprise_project_id=%s", epsID)
	}
	return rst
}

func buildContainerKubernetesSyncMccsDataListBodyParams(d *schema.ResourceData) []map[string]interface{} {
	rawArray, ok := d.Get("data_list").([]interface{})
	if !ok {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, item := range rawArray {
		cluster, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		rst = append(rst, map[string]interface{}{
			"cluster_id": cluster["cluster_id"],
		})
	}
	return rst
}

func buildContainerKubernetesSyncMccsBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"total_num": d.Get("total_num"),
		"data_list": buildContainerKubernetesSyncMccsDataListBodyParams(d),
	}
}

func resourceContainerKubernetesSyncMccsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/container/kubernetes/multi-cloud/clusters/status-synchronize"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildContainerKubernetesSyncMccsQueryParams(d, cfg)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildContainerKubernetesSyncMccsBodyParams(d)),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error synchronizing multi-cloud cluster status for container kubernetes: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}
	d.SetId(id)

	return resourceContainerKubernetesSyncMccsRead(ctx, d, meta)
}

func resourceContainerKubernetesSyncMccsRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceContainerKubernetesSyncMccsUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceContainerKubernetesSyncMccsDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to synchronize multi-cloud cluster status. Deleting
	this resource will not affect the synchronization status, but will only remove the resource information from
	the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
