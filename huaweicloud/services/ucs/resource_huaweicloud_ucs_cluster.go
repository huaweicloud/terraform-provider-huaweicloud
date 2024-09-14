// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product UCS
// ---------------------------------------------------------------

package ucs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API UCS POST /v1/clusters
// @API UCS GET /v1/clusters/{id}
// @API UCS PUT /v1/clusters/{id}
// @API UCS POST /v1/clusters/{id}/unjoin
// @API UCS POST /v1/clusters/{id}/join
// @API UCS DELETE /v1/clusters/{id}
func ResourceCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterCreate,
		UpdateContext: resourceClusterUpdate,
		ReadContext:   resourceClusterRead,
		DeleteContext: resourceClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"category": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the category of the cloud.`,
			},
			"cluster_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the cluster type.`,
			},
			"fleet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies ID of the fleet to add the cluster into.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the cluster id.`,
			},
			"cluster_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the cluster region.`,
				RequiredWith: []string{
					"cluster_id",
				},
			},
			"cluster_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the cluster project ID.`,
				RequiredWith: []string{
					"cluster_id",
				},
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the name of the cluster to register.`,
			},
			"cluster_labels": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the labels of the cluster to register.`,
			},
			"annotations": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the name of the cluster to register.`,
			},
			"service_provider": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the cloud service provider.`,
			},
			"country": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the country name.`,
			},
			"city": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the city name.`,
			},
			"manage_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the manage type.`,
			},
		},
	}
}

func resourceClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createCluster: Create a UCS Cluster.
	var (
		createClusterHttpUrl = "v1/clusters"
		createClusterProduct = "ucs"
	)
	createClusterClient, err := cfg.NewServiceClient(createClusterProduct, region)
	if err != nil {
		return diag.Errorf("error creating UCS Client: %s", err)
	}

	createClusterPath := createClusterClient.Endpoint + createClusterHttpUrl

	createClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}

	createClusterOpt.JSONBody = utils.RemoveNil(buildCreateClusterBodyParams(d))
	createClusterResp, err := createClusterClient.Request("POST", createClusterPath, &createClusterOpt)
	if err != nil {
		return diag.Errorf("error creating Cluster: %s", err)
	}

	createClusterRespBody, err := utils.FlattenResponse(createClusterResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("uid", createClusterRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating Cluster: ID is not found in API response")
	}
	d.SetId(id)

	return resourceClusterRead(ctx, d, meta)
}

func buildCreateClusterBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"kind":       "Cluster",
		"apiVersion": "v1",
		"metadata":   buildCreateClusterMetadataBodyParams(d),
		"spec":       buildCreateClusterSpecBodyParams(d),
	}
	return bodyParams
}

func buildCreateClusterMetadataBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"uid":         utils.ValueIgnoreEmpty(d.Get("cluster_id")),
		"name":        utils.ValueIgnoreEmpty(d.Get("cluster_name")),
		"labels":      utils.ValueIgnoreEmpty(d.Get("cluster_labels")),
		"annotations": utils.ValueIgnoreEmpty(d.Get("annotations")),
	}
	return bodyParams
}

func buildCreateClusterSpecBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"category":  utils.ValueIgnoreEmpty(d.Get("category")),
		"type":      utils.ValueIgnoreEmpty(d.Get("cluster_type")),
		"region":    utils.ValueIgnoreEmpty(d.Get("cluster_region")),
		"projectID": utils.ValueIgnoreEmpty(d.Get("cluster_project_id")),
		"provider":  utils.ValueIgnoreEmpty(d.Get("service_provider")),
		"country":   utils.ValueIgnoreEmpty(d.Get("country")),
		"city":      utils.ValueIgnoreEmpty(d.Get("city")),
	}

	if v, ok := d.GetOk("fleet_id"); ok {
		bodyParams["clusterGroupID"] = v
		bodyParams["manageType"] = "grouped"
	} else {
		bodyParams["manageType"] = "discrete"
	}
	return bodyParams
}

func resourceClusterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getCluster: Query the UCS Cluster detail
	var (
		getClusterHttpUrl = "v1/clusters/{id}"
		getClusterProduct = "ucs"
	)
	getClusterClient, err := cfg.NewServiceClient(getClusterProduct, region)
	if err != nil {
		return diag.Errorf("error creating UCS Client: %s", err)
	}

	getClusterPath := getClusterClient.Endpoint + getClusterHttpUrl
	getClusterPath = strings.ReplaceAll(getClusterPath, "{id}", d.Id())

	getClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getClusterResp, err := getClusterClient.Request("GET", getClusterPath, &getClusterOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Cluster")
	}

	getClusterRespBody, err := utils.FlattenResponse(getClusterResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// annotations not returned by API
	mErr = multierror.Append(
		mErr,
		d.Set("cluster_id", utils.PathSearch("metadata.uid", getClusterRespBody, nil)),
		d.Set("cluster_region", utils.PathSearch("spec.region", getClusterRespBody, nil)),
		d.Set("cluster_project_id", utils.PathSearch("spec.projectID", getClusterRespBody, nil)),
		d.Set("cluster_name", utils.PathSearch("metadata.name", getClusterRespBody, nil)),
		d.Set("cluster_labels", utils.PathSearch("metadata.labels", getClusterRespBody, nil)),
		d.Set("fleet_id", utils.PathSearch("spec.clusterGroupID", getClusterRespBody, nil)),
		d.Set("service_provider", utils.PathSearch("spec.provider", getClusterRespBody, nil)),
		d.Set("country", utils.PathSearch("spec.country", getClusterRespBody, nil)),
		d.Set("city", utils.PathSearch("spec.city", getClusterRespBody, nil)),
		d.Set("category", utils.PathSearch("spec.category", getClusterRespBody, nil)),
		d.Set("manage_type", utils.PathSearch("spec.manageType", getClusterRespBody, nil)),
		d.Set("cluster_type", utils.PathSearch("spec.type", getClusterRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updateClusterProduct = "ucs"
	)
	updateClusterClient, err := cfg.NewServiceClient(updateClusterProduct, region)
	if err != nil {
		return diag.Errorf("error creating UCS Client: %s", err)
	}

	updateClusterChanges := []string{
		"country",
		"city",
	}

	if d.HasChanges(updateClusterChanges...) {
		// updateCluster: Update the UCS Cluster
		var (
			updateClusterHttpUrl = "v1/clusters/{id}"
		)

		updateClusterPath := updateClusterClient.Endpoint + updateClusterHttpUrl
		updateClusterPath = strings.ReplaceAll(updateClusterPath, "{id}", d.Id())

		updateClusterOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}

		updateClusterOpt.JSONBody = utils.RemoveNil(buildUpdateClusterBodyParams(d))
		_, err = updateClusterClient.Request("PUT", updateClusterPath, &updateClusterOpt)
		if err != nil {
			return diag.Errorf("error updating Cluster: %s", err)
		}
	}

	if d.HasChanges("fleet_id") {
		oldFleetID, newFleetID := d.GetChange("fleet_id")

		if oldFleetID.(string) != "" {
			var (
				removeClusterFormFleetHttpUrl = "v1/clusters/{id}/unjoin"
			)

			removeClusterFormFleetPath := updateClusterClient.Endpoint + removeClusterFormFleetHttpUrl
			removeClusterFormFleetPath = strings.ReplaceAll(removeClusterFormFleetPath, "{id}", d.Id())

			removeClusterFormFleetOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			_, err = updateClusterClient.Request("POST", removeClusterFormFleetPath, &removeClusterFormFleetOpt)
			if err != nil {
				return diag.Errorf("error updating Cluster: %s", err)
			}
		}

		if newFleetID.(string) != "" {
			var (
				addClusterToFleetHttpUrl = "v1/clusters/{id}/join"
			)

			addClusterToFleetPath := updateClusterClient.Endpoint + addClusterToFleetHttpUrl
			addClusterToFleetPath = strings.ReplaceAll(addClusterToFleetPath, "{id}", d.Id())

			addClusterToFleetOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}

			addClusterToFleetOpt.JSONBody = map[string]interface{}{
				"clusterGroupID": newFleetID,
			}
			_, err = updateClusterClient.Request("POST", addClusterToFleetPath, &addClusterToFleetOpt)
			if err != nil {
				return diag.Errorf("error updating Cluster: %s", err)
			}
		}
	}
	return resourceClusterRead(ctx, d, meta)
}

func buildUpdateClusterBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"kind":       "Cluster",
		"apiVersion": "v1",
		"spec":       buildUpdateClusterSpecBodyParams(d),
	}
	return bodyParams
}

func buildUpdateClusterSpecBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"country": utils.ValueIgnoreEmpty(d.Get("country")),
		"city":    utils.ValueIgnoreEmpty(d.Get("city")),
	}
	return bodyParams
}

func resourceClusterDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteCluster: Delete an existing UCS Cluster
	var (
		deleteClusterHttpUrl = "v1/clusters/{id}"
		deleteClusterProduct = "ucs"
	)
	deleteClusterClient, err := cfg.NewServiceClient(deleteClusterProduct, region)
	if err != nil {
		return diag.Errorf("error creating UCS Client: %s", err)
	}

	deleteClusterPath := deleteClusterClient.Endpoint + deleteClusterHttpUrl
	deleteClusterPath = strings.ReplaceAll(deleteClusterPath, "{id}", d.Id())

	deleteClusterOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	_, err = deleteClusterClient.Request("DELETE", deleteClusterPath, &deleteClusterOpt)
	if err != nil {
		return diag.Errorf("error deleting Cluster: %s", err)
	}

	return nil
}
