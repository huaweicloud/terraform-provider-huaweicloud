package cce

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE POST /cce/cam/v3/clusters/{cluster_id}/releases
// @API CCE GET /cce/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}
// @API CCE PUT /cce/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}
// @API CCE DELETE /cce/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}

var releaseNonUpdatableParams = []string{"cluster_id", "chart_id", "name", "namespace", "version", "description"}

func ResourceRelease() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceReleaseCreate,
		ReadContext:   resourceReleaseRead,
		UpdateContext: resourceReleaseUpdate,
		DeleteContext: resourceReleaseDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceReleaseImport,
		},

		CustomizeDiff: config.FlexibleForceNew(releaseNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"chart_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"values": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_pull_policy": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"image_tag": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parameters": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dry_run": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"name_template": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"no_hooks": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"replace": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"recreate": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"reset_values": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"release_version": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"include_hooks": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"chart_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"chart_public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"chart_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateReleaseBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"chart_id":    d.Get("chart_id"),
		"name":        d.Get("name"),
		"namespace":   d.Get("namespace"),
		"version":     d.Get("version"),
		"values":      buildReleaseValuesParams(d),
		"description": d.Get("description"),
		"parameters":  buildReleaseParametersParams(d),
	}

	return bodyParams
}

func buildReleaseValuesParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"imagePullPolicy": utils.PathSearch("[0].image_pull_policy", d.Get("values"), nil),
		"imageTag":        utils.PathSearch("[0].image_tag", d.Get("values"), nil),
	}

	return bodyParams
}

func buildReleaseParametersParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"dry_run":         utils.PathSearch("[0].dry_run", d.Get("parameters"), nil),
		"name_template":   utils.PathSearch("[0].name_template", d.Get("parameters"), nil),
		"no_hooks":        utils.PathSearch("[0].no_hooks", d.Get("parameters"), nil),
		"replace":         utils.PathSearch("[0].replace", d.Get("parameters"), nil),
		"recreate":        utils.PathSearch("[0].recreate", d.Get("parameters"), nil),
		"reset_values":    utils.PathSearch("[0].reset_values", d.Get("parameters"), nil),
		"release_version": utils.PathSearch("[0].release_version", d.Get("parameters"), nil),
		"include_hooks":   utils.PathSearch("[0].include_hooks", d.Get("parameters"), nil),
	}

	return bodyParams
}

func resourceReleaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createReleaseHttpUrl = "cce/cam/v3/clusters/{cluster_id}/releases"
		createReleaseProduct = "cce"
	)
	createReleaseClient, err := cfg.NewServiceClient(createReleaseProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	createReleasePath := createReleaseClient.Endpoint + createReleaseHttpUrl
	createReleasePath = strings.ReplaceAll(createReleasePath, "{cluster_id}", d.Get("cluster_id").(string))

	createReleaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createReleaseOpt.JSONBody = utils.RemoveNil(buildCreateReleaseBodyParams(d))
	_, err = createReleaseClient.Request("POST", createReleasePath, &createReleaseOpt)
	if err != nil {
		return diag.Errorf("error creating CCE release: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceReleaseRead(ctx, d, meta)
}

func resourceReleaseRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		getReleaseHttpUrl = "cce/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}"
		getReleaseProduct = "cce"
	)
	getReleaseClient, err := cfg.NewServiceClient(getReleaseProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	getReleaseHttpPath := getReleaseClient.Endpoint + getReleaseHttpUrl
	getReleaseHttpPath = strings.ReplaceAll(getReleaseHttpPath, "{cluster_id}", d.Get("cluster_id").(string))
	getReleaseHttpPath = strings.ReplaceAll(getReleaseHttpPath, "{namespace}", d.Get("namespace").(string))
	getReleaseHttpPath = strings.ReplaceAll(getReleaseHttpPath, "{name}", d.Id())

	getReleaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getReleaseResp, err := getReleaseClient.Request("GET", getReleaseHttpPath, &getReleaseOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CCE release")
	}

	getReleaseRespBody, err := utils.FlattenResponse(getReleaseResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// version, values, chart_id, description, parameters, action are not returned in GET API
	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("cluster_id", utils.PathSearch("cluster_id", getReleaseRespBody, nil)),
		d.Set("name", utils.PathSearch("name", getReleaseRespBody, nil)),
		d.Set("namespace", utils.PathSearch("namespace", getReleaseRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getReleaseRespBody, nil)),
		d.Set("status_description", utils.PathSearch("status_description", getReleaseRespBody, nil)),
		d.Set("cluster_name", utils.PathSearch("cluster_name", getReleaseRespBody, nil)),
		d.Set("chart_name", utils.PathSearch("chart_name", getReleaseRespBody, nil)),
		d.Set("chart_public", utils.PathSearch("chart_public", getReleaseRespBody, nil)),
		d.Set("chart_version", utils.PathSearch("chart_version", getReleaseRespBody, nil)),
		d.Set("created_at", utils.PathSearch("create_at", getReleaseRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("update_at", getReleaseRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateReleaseBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"chart_id":   d.Get("chart_id"),
		"action":     d.Get("action"),
		"values":     buildReleaseValuesParams(d),
		"parameters": buildReleaseParametersParams(d),
	}

	return bodyParams
}

func resourceReleaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updateReleaseHttpUrl = "cce/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}"
		updateReleaseProduct = "cce"
	)
	updateReleaseClient, err := cfg.NewServiceClient(updateReleaseProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	updateReleasePath := updateReleaseClient.Endpoint + updateReleaseHttpUrl
	updateReleasePath = strings.ReplaceAll(updateReleasePath, "{cluster_id}", d.Get("cluster_id").(string))
	updateReleasePath = strings.ReplaceAll(updateReleasePath, "{namespace}", d.Get("namespace").(string))
	updateReleasePath = strings.ReplaceAll(updateReleasePath, "{name}", d.Id())

	updateReleaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	updateReleaseOpt.JSONBody = utils.RemoveNil(buildUpdateReleaseBodyParams(d))
	_, err = updateReleaseClient.Request("PUT", updateReleasePath, &updateReleaseOpt)
	if err != nil {
		return diag.Errorf("error updating CCE release: %s", err)
	}

	return resourceReleaseRead(ctx, d, meta)
}

func resourceReleaseDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteReleaseHttpUrl = "cce/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}"
		deleteReleaseProduct = "cce"
	)
	deleteReleaseClient, err := cfg.NewServiceClient(deleteReleaseProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	deleteReleaseHttpPath := deleteReleaseClient.Endpoint + deleteReleaseHttpUrl
	deleteReleaseHttpPath = strings.ReplaceAll(deleteReleaseHttpPath, "{cluster_id}", d.Get("cluster_id").(string))
	deleteReleaseHttpPath = strings.ReplaceAll(deleteReleaseHttpPath, "{namespace}", d.Get("namespace").(string))
	deleteReleaseHttpPath = strings.ReplaceAll(deleteReleaseHttpPath, "{name}", d.Id())

	deleteReleaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = deleteReleaseClient.Request("DELETE", deleteReleaseHttpPath, &deleteReleaseOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CCE release")
	}

	return nil
}

func resourceReleaseImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		err := errors.New("invalid format specified for CCE Node. Format must be <cluster id>/<namespace>/<name>")
		return nil, err
	}

	clusterID := parts[0]
	namespace := parts[1]
	name := parts[2]

	d.SetId(name)
	d.Set("name", name)
	d.Set("cluster_id", clusterID)
	d.Set("namespace", namespace)

	return []*schema.ResourceData{d}, nil
}
