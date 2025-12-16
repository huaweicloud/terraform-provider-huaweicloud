// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product SWR
// ---------------------------------------------------------------

package swr

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cce/v3/clusters"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SWR POST /v2/manage/namespaces/{namespace}/repos/{repository}/triggers
// @API SWR DELETE /v2/manage/namespaces/{namespace}/repos/{repository}/triggers/{trigger}
// @API SWR GET /v2/manage/namespaces/{namespace}/repos/{repository}/triggers/{trigger}
// @API SWR PATCH /v2/manage/namespaces/{namespace}/repos/{repository}/triggers/{trigger}
func ResourceSwrImageTrigger() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrImageTriggerCreate,
		UpdateContext: resourceSwrImageTriggerUpdate,
		ReadContext:   resourceSwrImageTriggerRead,
		DeleteContext: resourceSwrImageTriggerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSwrImageTriggerImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the organization.`,
			},
			"repository": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the repository.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the trigger name.`,
			},
			"workload_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the type of the application.`,
			},
			"workload_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the application.`,
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the namespace where the application is located.`,
			},
			"condition_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the trigger condition type.`,
			},
			"condition_value": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the trigger condition value.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the cluster.`,
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the name of the cluster.`,
			},
			"container": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the name of the container to be updated.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the trigger type.`,
			},
			"enabled": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to enable the trigger.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time.`,
			},
			"creator_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creator name of the trigger.`,
			},
		},
	}
}

func resourceSwrImageTriggerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createSwrImageTrigger: create SWR image trigger
	var (
		createSwrImageTriggerHttpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/triggers"
		createSwrImageTriggerProduct = "swr"
	)
	createSwrImageTriggerClient, err := cfg.NewServiceClient(createSwrImageTriggerProduct, region)
	if err != nil {
		return diag.Errorf("error creating SWR Client: %s", err)
	}

	organization := d.Get("organization").(string)
	repository := d.Get("repository").(string)
	name := d.Get("name").(string)
	createSwrImageTriggerPath := createSwrImageTriggerClient.Endpoint + createSwrImageTriggerHttpUrl
	createSwrImageTriggerPath = strings.ReplaceAll(createSwrImageTriggerPath, "{namespace}",
		organization)
	createSwrImageTriggerPath = strings.ReplaceAll(createSwrImageTriggerPath, "{repository}",
		strings.ReplaceAll(repository, "/", "$"))

	createSwrImageTriggerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	clusterName := d.Get("cluster_name").(string)
	clusterId := d.Get("cluster_id").(string)
	if clusterName == "" && clusterId != "" {
		clusterName, err = getClusterNameByClusterId(d, cfg, clusterId)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	createSwrImageTriggerOpt.JSONBody = utils.RemoveNil(buildCreateSwrImageTriggerBodyParams(d, clusterName))
	_, err = createSwrImageTriggerClient.Request("POST", createSwrImageTriggerPath, &createSwrImageTriggerOpt)
	if err != nil {
		return diag.Errorf("error creating SWR image trigger: %s", err)
	}

	d.SetId(organization + "/" + repository + "/" + name)

	return resourceSwrImageTriggerRead(ctx, d, meta)
}

func getClusterNameByClusterId(d *schema.ResourceData, cfg *config.Config, clusterId string) (string, error) {
	cceClient, err := cfg.CceV3Client(cfg.GetRegion(d))
	if err != nil {
		return "", fmt.Errorf("error creating CCE client: %s", err)
	}
	ns, err := clusters.Get(cceClient, clusterId).Extract()
	if err != nil {
		return "", fmt.Errorf("error retrieving CCE cluster, err: %s", err)
	}
	return ns.Metadata.Name, nil
}

func buildCreateSwrImageTriggerBodyParams(d *schema.ResourceData, clusterName string) map[string]interface{} {
	enabled := d.Get("enabled")
	if enabled == "" {
		enabled = "true"
	}
	bodyParams := map[string]interface{}{
		"action":       "update",
		"name":         utils.ValueIgnoreEmpty(d.Get("name")),
		"app_type":     utils.ValueIgnoreEmpty(d.Get("workload_type")),
		"application":  utils.ValueIgnoreEmpty(d.Get("workload_name")),
		"cluster_id":   utils.ValueIgnoreEmpty(d.Get("cluster_id")),
		"cluster_name": utils.ValueIgnoreEmpty(clusterName),
		"cluster_ns":   utils.ValueIgnoreEmpty(d.Get("namespace")),
		"trigger_type": utils.ValueIgnoreEmpty(d.Get("condition_type")),
		"condition":    utils.ValueIgnoreEmpty(d.Get("condition_value")),
		"container":    utils.ValueIgnoreEmpty(d.Get("container")),
		"trigger_mode": utils.ValueIgnoreEmpty(d.Get("type")),
		"enable":       enabled,
	}
	return bodyParams
}

func resourceSwrImageTriggerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getSwrImageTrigger: Query SWR image trigger
	var (
		getSwrImageTriggerHttpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/triggers/{trigger}"
		getSwrImageTriggerProduct = "swr"
	)
	getSwrImageTriggerClient, err := cfg.NewServiceClient(getSwrImageTriggerProduct, region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	organization := d.Get("organization").(string)
	repository := d.Get("repository").(string)
	trigger := d.Get("name").(string)

	getSwrImageTriggerPath := getSwrImageTriggerClient.Endpoint + getSwrImageTriggerHttpUrl
	getSwrImageTriggerPath = strings.ReplaceAll(getSwrImageTriggerPath, "{namespace}", organization)
	getSwrImageTriggerPath = strings.ReplaceAll(getSwrImageTriggerPath, "{repository}", strings.ReplaceAll(repository, "/", "$"))
	getSwrImageTriggerPath = strings.ReplaceAll(getSwrImageTriggerPath, "{trigger}", trigger)

	getSwrImageTriggerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getSwrImageTriggerResp, err := getSwrImageTriggerClient.Request("GET",
		getSwrImageTriggerPath, &getSwrImageTriggerOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SWR image trigger")
	}

	getSwrImageTriggerRespBody, err := utils.FlattenResponse(getSwrImageTriggerResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("organization", organization),
		d.Set("repository", repository),
		d.Set("workload_type", utils.PathSearch("app_type", getSwrImageTriggerRespBody, nil)),
		d.Set("workload_name", utils.PathSearch("application", getSwrImageTriggerRespBody, nil)),
		d.Set("cluster_id", utils.PathSearch("cluster_id", getSwrImageTriggerRespBody, nil)),
		d.Set("cluster_name", utils.PathSearch("cluster_name", getSwrImageTriggerRespBody, nil)),
		d.Set("namespace", utils.PathSearch("cluster_ns", getSwrImageTriggerRespBody, nil)),
		d.Set("condition_value", utils.PathSearch("condition", getSwrImageTriggerRespBody, nil)),
		d.Set("container", utils.PathSearch("container", getSwrImageTriggerRespBody, nil)),
		d.Set("enabled", utils.PathSearch("enable", getSwrImageTriggerRespBody, nil)),
		d.Set("name", utils.PathSearch("name", getSwrImageTriggerRespBody, nil)),
		d.Set("type", utils.PathSearch("trigger_mode", getSwrImageTriggerRespBody, nil)),
		d.Set("condition_type", utils.PathSearch("trigger_type", getSwrImageTriggerRespBody, nil)),
		d.Set("created_at", utils.PathSearch("created_at", getSwrImageTriggerRespBody, nil)),
		d.Set("creator_name", utils.PathSearch("creator_name", getSwrImageTriggerRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSwrImageTriggerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateSwrImageTriggerHasChanges := []string{
		"enabled",
	}

	if d.HasChanges(updateSwrImageTriggerHasChanges...) {
		// updateSwrImageTrigger: update SWR image trigger
		var (
			updateSwrImageTriggerHttpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/triggers/{trigger}"
			updateSwrImageTriggerProduct = "swr"
		)
		updateSwrImageTriggerClient, err := cfg.NewServiceClient(updateSwrImageTriggerProduct, region)
		if err != nil {
			return diag.Errorf("error creating SWR client: %s", err)
		}

		updateSwrImageTriggerPath := updateSwrImageTriggerClient.Endpoint + updateSwrImageTriggerHttpUrl
		updateSwrImageTriggerPath = strings.ReplaceAll(updateSwrImageTriggerPath, "{namespace}",
			fmt.Sprintf("%v", d.Get("organization")))
		updateSwrImageTriggerPath = strings.ReplaceAll(updateSwrImageTriggerPath, "{repository}",
			strings.ReplaceAll(d.Get("repository").(string), "/", "$"))
		updateSwrImageTriggerPath = strings.ReplaceAll(updateSwrImageTriggerPath, "{trigger}",
			fmt.Sprintf("%v", d.Get("name")))

		updateSwrImageTriggerOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				201,
			},
		}
		updateSwrImageTriggerOpt.JSONBody = utils.RemoveNil(buildUpdateSwrImageTriggerBodyParams(d))
		_, err = updateSwrImageTriggerClient.Request("PATCH", updateSwrImageTriggerPath,
			&updateSwrImageTriggerOpt)
		if err != nil {
			return diag.Errorf("error updating SWR image trigger: %s", err)
		}
	}
	return resourceSwrImageTriggerRead(ctx, d, meta)
}

func buildUpdateSwrImageTriggerBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"enable": utils.ValueIgnoreEmpty(d.Get("enabled")),
	}
	return bodyParams
}

func resourceSwrImageTriggerDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteSwrImageTrigger: Delete SWR image trigger
	var (
		deleteSwrImageTriggerHttpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/triggers/{trigger}"
		deleteSwrImageTriggerProduct = "swr"
	)
	deleteSwrImageTriggerClient, err := cfg.NewServiceClient(deleteSwrImageTriggerProduct, region)
	if err != nil {
		return diag.Errorf("error creating SWR image trigger client: %s", err)
	}

	deleteSwrImageTriggerPath := deleteSwrImageTriggerClient.Endpoint + deleteSwrImageTriggerHttpUrl
	deleteSwrImageTriggerPath = strings.ReplaceAll(deleteSwrImageTriggerPath, "{namespace}",
		fmt.Sprintf("%v", d.Get("organization")))
	deleteSwrImageTriggerPath = strings.ReplaceAll(deleteSwrImageTriggerPath, "{repository}",
		strings.ReplaceAll(d.Get("repository").(string), "/", "$"))
	deleteSwrImageTriggerPath = strings.ReplaceAll(deleteSwrImageTriggerPath, "{trigger}",
		fmt.Sprintf("%v", d.Get("name")))

	deleteSwrImageTriggerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = deleteSwrImageTriggerClient.Request("DELETE", deleteSwrImageTriggerPath, &deleteSwrImageTriggerOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SWR image trigger")
	}

	return nil
}

func resourceSwrImageTriggerImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ",")
	if len(parts) != 3 {
		parts = strings.Split(d.Id(), "/")
		if len(parts) != 3 {
			return nil, errors.New("invalid id format, must be <organization_name>/<repository_name>/<trigger_name> or " +
				"<organization_name>,<repository_name>,<trigger_name>")
		}
	} else {
		// reform ID to be separated by slashes
		id := fmt.Sprintf("%s/%s/%s", parts[0], parts[1], parts[2])
		d.SetId(id)
	}

	d.Set("organization", parts[0])
	d.Set("repository", parts[1])
	d.Set("name", parts[2])

	return []*schema.ResourceData{d}, nil
}
