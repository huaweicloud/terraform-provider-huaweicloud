package cceautopilot

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

// @API CCE POST /autopilot/cam/v3/clusters/{cluster_id}/releases
// @API CCE GET /autopilot/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}
// @API CCE PUT /autopilot/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}
// @API CCE DELETE /autopilot/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}

var releaseNonUpdatableParams = []string{"cluster_id", "chart_id", "name", "namespace", "version", "description"}

func ResourceAutopilotRelease() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAutopilotReleaseCreate,
		ReadContext:   resourceAutopilotReleaseRead,
		UpdateContext: resourceAutopilotReleaseUpdate,
		DeleteContext: resourceAutopilotReleaseDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAutopilotReleaseImport,
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

func buildCreateAutopilotReleaseBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"chart_id":    d.Get("chart_id"),
		"name":        d.Get("name"),
		"namespace":   d.Get("namespace"),
		"version":     d.Get("version"),
		"values":      buildAutopilotReleaseValuesParams(d),
		"description": d.Get("description"),
		"parameters":  buildAutopilotReleaseParametersParams(d),
	}

	return bodyParams
}

func buildAutopilotReleaseValuesParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"imagePullPolicy": utils.PathSearch("[0].image_pull_policy", d.Get("values"), nil),
		"imageTag":        utils.PathSearch("[0].image_tag", d.Get("values"), nil),
	}

	return bodyParams
}

func buildAutopilotReleaseParametersParams(d *schema.ResourceData) map[string]interface{} {
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

func resourceAutopilotReleaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createAutopilotReleaseHttpUrl = "autopilot/cam/v3/clusters/{cluster_id}/releases"
		createAutopilotReleaseProduct = "cce"
	)
	createAutopilotReleaseClient, err := cfg.NewServiceClient(createAutopilotReleaseProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	createAutopilotReleasePath := createAutopilotReleaseClient.Endpoint + createAutopilotReleaseHttpUrl
	createAutopilotReleasePath = strings.ReplaceAll(createAutopilotReleasePath, "{cluster_id}", d.Get("cluster_id").(string))

	createAutopilotReleaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createAutopilotReleaseOpt.JSONBody = utils.RemoveNil(buildCreateAutopilotReleaseBodyParams(d))
	_, err = createAutopilotReleaseClient.Request("POST", createAutopilotReleasePath, &createAutopilotReleaseOpt)
	if err != nil {
		return diag.Errorf("error creating CCE autopilot release: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceAutopilotReleaseRead(ctx, d, meta)
}

func resourceAutopilotReleaseRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		getAutopilotReleaseHttpUrl = "autopilot/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}"
		getAutopilotReleaseProduct = "cce"
	)
	getAutopilotReleaseClient, err := cfg.NewServiceClient(getAutopilotReleaseProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	getAutopilotReleaseHttpPath := getAutopilotReleaseClient.Endpoint + getAutopilotReleaseHttpUrl
	getAutopilotReleaseHttpPath = strings.ReplaceAll(getAutopilotReleaseHttpPath, "{cluster_id}", d.Get("cluster_id").(string))
	getAutopilotReleaseHttpPath = strings.ReplaceAll(getAutopilotReleaseHttpPath, "{namespace}", d.Get("namespace").(string))
	getAutopilotReleaseHttpPath = strings.ReplaceAll(getAutopilotReleaseHttpPath, "{name}", d.Id())

	getAutopilotReleaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getAutopilotReleaseResp, err := getAutopilotReleaseClient.Request("GET", getAutopilotReleaseHttpPath, &getAutopilotReleaseOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CCE autopilot release")
	}

	getAutopilotReleaseRespBody, err := utils.FlattenResponse(getAutopilotReleaseResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// version, values, chart_id, description, parameters, action are not returned in GET API
	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("cluster_id", utils.PathSearch("cluster_id", getAutopilotReleaseRespBody, nil)),
		d.Set("name", utils.PathSearch("name", getAutopilotReleaseRespBody, nil)),
		d.Set("namespace", utils.PathSearch("namespace", getAutopilotReleaseRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getAutopilotReleaseRespBody, nil)),
		d.Set("status_description", utils.PathSearch("status_description", getAutopilotReleaseRespBody, nil)),
		d.Set("cluster_name", utils.PathSearch("cluster_name", getAutopilotReleaseRespBody, nil)),
		d.Set("chart_name", utils.PathSearch("chart_name", getAutopilotReleaseRespBody, nil)),
		d.Set("chart_public", utils.PathSearch("chart_public", getAutopilotReleaseRespBody, nil)),
		d.Set("chart_version", utils.PathSearch("chart_version", getAutopilotReleaseRespBody, nil)),
		d.Set("created_at", utils.PathSearch("create_at", getAutopilotReleaseRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("update_at", getAutopilotReleaseRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateAutopilotReleaseBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"chart_id":   d.Get("chart_id"),
		"action":     d.Get("action"),
		"values":     buildAutopilotReleaseValuesParams(d),
		"parameters": buildAutopilotReleaseParametersParams(d),
	}

	return bodyParams
}

func resourceAutopilotReleaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updateAutopilotReleaseHttpUrl = "autopilot/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}"
		updateAutopilotReleaseProduct = "cce"
	)
	updateAutopilotReleaseClient, err := cfg.NewServiceClient(updateAutopilotReleaseProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	updateAutopilotReleasePath := updateAutopilotReleaseClient.Endpoint + updateAutopilotReleaseHttpUrl
	updateAutopilotReleasePath = strings.ReplaceAll(updateAutopilotReleasePath, "{cluster_id}", d.Get("cluster_id").(string))
	updateAutopilotReleasePath = strings.ReplaceAll(updateAutopilotReleasePath, "{namespace}", d.Get("namespace").(string))
	updateAutopilotReleasePath = strings.ReplaceAll(updateAutopilotReleasePath, "{name}", d.Id())

	updateAutopilotReleaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	updateAutopilotReleaseOpt.JSONBody = utils.RemoveNil(buildUpdateAutopilotReleaseBodyParams(d))
	_, err = updateAutopilotReleaseClient.Request("PUT", updateAutopilotReleasePath, &updateAutopilotReleaseOpt)
	if err != nil {
		return diag.Errorf("error updating CCE autopilot release: %s", err)
	}

	return resourceAutopilotReleaseRead(ctx, d, meta)
}

func resourceAutopilotReleaseDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteAutopilotReleaseHttpUrl = "autopilot/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}"
		deleteAutopilotReleaseProduct = "cce"
	)
	deleteAutopilotReleaseClient, err := cfg.NewServiceClient(deleteAutopilotReleaseProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	deleteAutopilotReleaseHttpPath := deleteAutopilotReleaseClient.Endpoint + deleteAutopilotReleaseHttpUrl
	deleteAutopilotReleaseHttpPath = strings.ReplaceAll(deleteAutopilotReleaseHttpPath, "{cluster_id}", d.Get("cluster_id").(string))
	deleteAutopilotReleaseHttpPath = strings.ReplaceAll(deleteAutopilotReleaseHttpPath, "{namespace}", d.Get("namespace").(string))
	deleteAutopilotReleaseHttpPath = strings.ReplaceAll(deleteAutopilotReleaseHttpPath, "{name}", d.Id())

	deleteAutopilotReleaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = deleteAutopilotReleaseClient.Request("DELETE", deleteAutopilotReleaseHttpPath, &deleteAutopilotReleaseOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CCE autopilot release")
	}

	return nil
}

func resourceAutopilotReleaseImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
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
