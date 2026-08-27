package cceautopilot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

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
			"release_version": {
				Type:     schema.TypeInt,
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
	client, err := cfg.NewServiceClient("cce", region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	createPath := client.Endpoint + "autopilot/cam/v3/clusters/{cluster_id}/releases"
	createPath = strings.ReplaceAll(createPath, "{cluster_id}", d.Get("cluster_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildCreateAutopilotReleaseBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CCE autopilot release: %s", err)
	}
	d.SetId(d.Get("name").(string))

	err = waitingForReleaseJobCompleted(ctx, client, d, d.Timeout(schema.TimeoutCreate), []string{"DEPLOYED"})
	if err != nil {
		return diag.Errorf("error creating CCE autopilot release: %s", err)
	}

	return resourceAutopilotReleaseRead(ctx, d, meta)
}

func resourceAutopilotReleaseRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cce", region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	getAutopilotReleaseRespBody, err := getAutopilotRelease(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CCE autopilot release")
	}

	// values, chart_id, description, parameters, action these set by user are not returned in GET API
	// version returned in the response is the version of release
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
		d.Set("release_version", utils.PathSearch("version", getAutopilotReleaseRespBody, nil)),
		d.Set("created_at", utils.PathSearch("create_at", getAutopilotReleaseRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("update_at", getAutopilotReleaseRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getAutopilotRelease(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	getPath := client.Endpoint + "autopilot/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}"
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", d.Get("cluster_id").(string))
	getPath = strings.ReplaceAll(getPath, "{namespace}", d.Get("namespace").(string))
	getPath = strings.ReplaceAll(getPath, "{name}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
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
	client, err := cfg.NewServiceClient("cce", region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	if d.HasChangeExcept("enable_force_new") {
		updatePath := client.Endpoint + "autopilot/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}"
		updatePath = strings.ReplaceAll(updatePath, "{cluster_id}", d.Get("cluster_id").(string))
		updatePath = strings.ReplaceAll(updatePath, "{namespace}", d.Get("namespace").(string))
		updatePath = strings.ReplaceAll(updatePath, "{name}", d.Id())

		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			JSONBody: utils.RemoveNil(buildUpdateAutopilotReleaseBodyParams(d)),
		}

		retryFunc := func() (interface{}, bool, error) {
			res, err := client.Request("PUT", updatePath, &updateOpt)
			if err == nil {
				return res, false, nil
			}
			shouldRetry := strings.Contains(err.Error(), "Update release is forbidden")
			return nil, shouldRetry, err
		}

		_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     refreshAutoPilotReleaseStatus(client, d, []string{"DEPLOYED"}),
			WaitTarget:   []string{"COMPLETED"},
			WaitPending:  []string{"PENDING"},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			DelayTimeout: 10 * time.Second,
			PollInterval: 5 * time.Second,
		})
		if err != nil {
			return diag.Errorf("error updating CCE autopilot release: %s", err)
		}

		err = waitingForReleaseJobCompleted(ctx, client, d, d.Timeout(schema.TimeoutUpdate), []string{"DEPLOYED"})
		if err != nil {
			return diag.Errorf("error updating CCE autopilot release: %s", err)
		}
	}
	return resourceAutopilotReleaseRead(ctx, d, meta)
}

func resourceAutopilotReleaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cce", region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	deletePath := client.Endpoint + "autopilot/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}"
	deletePath = strings.ReplaceAll(deletePath, "{cluster_id}", d.Get("cluster_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{namespace}", d.Get("namespace").(string))
	deletePath = strings.ReplaceAll(deletePath, "{name}", d.Id())

	deleteAutopilotReleaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteAutopilotReleaseOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CCE autopilot release")
	}

	stateConf := &retry.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"DELETED"},
		Refresh:      refreshAutoPilotReleaseStatus(client, d, []string{"DELETED"}),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting CCE autopilot release: %s", err)
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
	mErr := multierror.Append(nil,
		d.Set("name", name),
		d.Set("cluster_id", clusterID),
		d.Set("namespace", namespace),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

func waitingForReleaseJobCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	t time.Duration, targets []string) error {
	stateConf := &retry.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshAutoPilotReleaseStatus(client, d, targets),
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func refreshAutoPilotReleaseStatus(client *golangsdk.ServiceClient, d *schema.ResourceData, targets []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Expect the status of CCE autopilot release to be any one of the status list: %v.", targets)
		releaseRespBody, err := getAutopilotRelease(client, d)
		if err != nil {
			var errDefault404 golangsdk.ErrDefault404
			if errors.As(err, &errDefault404) {
				return "Resource Not Found", "DELETED", nil
			}
			return nil, "ERROR", fmt.Errorf("error retrieving CCE autopilot release: %s", err)
		}

		status := utils.PathSearch("status", releaseRespBody, nil)
		if status == nil {
			return nil, "ERROR", fmt.Errorf("error parsing status from response body")
		}

		statusStr := status.(string)
		invalidStatuses := []string{"FAILED", "UNKNOWN"}
		if utils.IsStrContainsSliceElement(statusStr, invalidStatuses, true, true) {
			if statusStr == "UNKNOWN" {
				return nil, "ERROR", fmt.Errorf("the release status is unknown, please try to delete and reinstall manually")
			}
			return nil, "ERROR", fmt.Errorf("the release job failed, status: %s", statusStr)
		}

		if utils.StrSliceContains(targets, statusStr) {
			return releaseRespBody, "COMPLETED", nil
		}
		return releaseRespBody, "PENDING", nil
	}
}
