package workspace

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace POST /v1/{project_id}/image-servers
// @API Workspace GET /v1/{project_id}/image-servers
// @API Workspace PATCH /v1/{project_id}/image-servers/{server_id}
// @API Workspace PATCH /v1/{project_id}/image-servers/actions/batch-delete
func ResourceAppImageServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppImageServerCreate,
		ReadContext:   resourceAppImageServerRead,
		UpdateContext: resourceAppImageServerUpdate,
		DeleteContext: resourceAppImageServerDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the image server.",
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The flavor ID of the image server.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The VPC ID to which the image server belongs.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The subnet ID to which the image server belongs.",
			},
			"root_volume": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The disk type of the image server.",
						},
						"size": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    true,
							Description: "The disk size of the image server, in GB.",
						},
					},
				},
				Description: "The system disk configuration of the image server.",
			},
			"authorize_accounts": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The name of the account.",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The type of the account.",
						},
						"domain": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The domain name of the Workspace service.",
						},
					},
				},
				Description: "The list of the management accounts for creating the image.",
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The basic image ID of the image server.",
			},
			"image_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The basic image type of the image server.",
			},
			"spec_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The specification code of the basic image to which the image server belongs.",
			},
			"image_source_product_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The basic image product ID of the image server.",
			},
			"is_vdi": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "The session mode of the image server.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The availability zone of the image server.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the image server.",
			},
			"ou_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The OU name corresponding to the AD server.",
			},
			"extra_session_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The additional session type.",
			},
			"extra_session_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The number of additional sessions for a single server.",
			},
			"route_policy": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_session": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "The number of session connections of the server.",
						},
						"cpu_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "The CPU usage of the server.",
						},
						"mem_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "The memory usage of the server.",
						},
					},
				},
				Description: "The session scheduling policy of the server associated with the image server.",
			},
			"tags": common.TagsForceNewSchema("The key/value pairs to associate with the image server."),
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ID of the enterprise project to which the image server belong.",
			},
			"scheduler_hints": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "The configuration of the dedicate host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dedicated_host_id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The ID of the dedicate host.",
						},
						"tenancy": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The type of the dedicate host.",
						},
					},
				},
			},
			"is_delete_associated_resources": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to delete resources associated with this image server after deleting it.",
			},
			"attach_apps": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: utils.SchemaDesc("The list of the warehouse apps.", utils.SchemaDescInput{Internal: true}),
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the image server, in RFC3339 format.",
			},
		},
	}
}

func buildCreateAppImageServerOpts(d *schema.ResourceData, epsId string) map[string]interface{} {
	return map[string]interface{}{
		"name":               d.Get("name"),
		"product_id":         d.Get("flavor_id"),
		"vpc_id":             d.Get("vpc_id"),
		"subnet_id":          d.Get("subnet_id"),
		"root_volume":        buildAppServerRootVolume(d.Get("root_volume").([]interface{})),
		"authorize_accounts": buildImageServerAccounts(d.Get("authorize_accounts").(*schema.Set).List()),
		"image_ref": map[string]interface{}{
			"id":         d.Get("image_id"),
			"image_type": d.Get("image_type"),
			"spce_code":  utils.ValueIgnoreEmpty(d.Get("spec_code")),
			"product_id": utils.ValueIgnoreEmpty(d.Get("image_source_product_id")),
		},
		"is_vdi":                d.Get("is_vdi"),
		"availability_zone":     utils.ValueIgnoreEmpty(d.Get("availability_zone")),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"ou_name":               utils.ValueIgnoreEmpty(d.Get("ou_name")),
		"extra_session_type":    utils.ValueIgnoreEmpty(d.Get("extra_session_type")),
		"extra_session_size":    utils.ValueIgnoreEmpty(d.Get("extra_session_size")),
		"route_policy":          buildAppServerGroupRoutePolicy(d.Get("route_policy").([]interface{})),
		"tags":                  utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		"enterprise_project_id": utils.ValueIgnoreEmpty(epsId),
		"scheduler_hints":       buildAppServerSchedulerHints(d.Get("scheduler_hints").([]interface{})),
		"attach_apps":           utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("attach_apps").([]interface{}))),
	}
}

func buildImageServerAccounts(accounts []interface{}) []map[string]interface{} {
	if len(accounts) == 0 {
		return nil
	}

	res := make([]map[string]interface{}, len(accounts))
	for i, v := range accounts {
		res[i] = map[string]interface{}{
			"account":      utils.PathSearch("account", v, ""),
			"account_type": utils.PathSearch("type", v, ""),
			"domain":       utils.PathSearch("domain", v, ""),
		}
	}
	return res
}

func resourceAppImageServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/image-servers"
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateAppImageServerOpts(d, cfg.GetEnterpriseProjectID(d))),
	}
	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating image server of Workspace APP: %s", err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find job ID from API response")
	}

	serverResp, err := waitForImageServerJobCompleted(ctx, client, d.Timeout(schema.TimeoutCreate), jobId)
	if err != nil {
		return diag.Errorf("error waiting for creating image server job (%s) completed: %s", jobId, err)
	}

	imageServerId := utils.PathSearch("sub_jobs|[0].job_resource_info.resource_id", serverResp, "").(string)
	if imageServerId == "" {
		return diag.Errorf("unable to find image server ID from API response")
	}

	d.SetId(imageServerId)

	return resourceAppImageServerRead(ctx, d, meta)
}

func waitForImageServerJobCompleted(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration, jobId string) (interface{},
	error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"RUNNING"},
		Target:       []string{"SUCCESS"},
		Refresh:      refreshImageServerJobStatusFunc(client, jobId),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 30 * time.Second,
	}

	serverResp, err := stateConf.WaitForStateContext(ctx)
	return serverResp, err
}

func refreshImageServerJobStatusFunc(client *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		httpUrl := "v1/{project_id}/image-server-jobs/{job_id}"
		getJobPath := client.Endpoint + httpUrl
		getJobPath = strings.ReplaceAll(getJobPath, "{project_id}", client.ProjectID)
		getJobPath = strings.ReplaceAll(getJobPath, "{job_id}", jobId)
		getOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		resp, err := client.Request("GET", getJobPath, &getOpt)
		if err != nil {
			return resp, "ERROR", err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return resp, "ERROR", err
		}

		return respBody, utils.PathSearch("status", respBody, nil).(string), nil
	}
}

func resourceAppImageServerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		imageServerId = d.Id()
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	imageServer, err := GetAppImageServerById(client, imageServerId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Workspace APP image server")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", imageServer, nil)),
		d.Set("authorize_accounts",
			flattenAuthorizAccounts(utils.PathSearch("authorize_accounts", imageServer, make([]interface{}, 0)).([]interface{}))),
		d.Set("image_id", utils.PathSearch("image_ref.id", imageServer, nil)),
		d.Set("image_type", utils.PathSearch("image_ref.image_type", imageServer, nil)),
		d.Set("spec_code", utils.PathSearch("image_ref.spce_code", imageServer, nil)),
		d.Set("description", utils.PathSearch("description", imageServer, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", imageServer, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("create_time",
			imageServer, "").(string))/1000, false)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func flattenAuthorizAccounts(accounts []interface{}) []map[string]interface{} {
	if len(accounts) == 0 {
		return nil
	}

	rest := make([]map[string]interface{}, len(accounts))
	for i, v := range accounts {
		rest[i] = map[string]interface{}{
			"account": utils.PathSearch("account", v, nil),
			"type":    utils.PathSearch("account_type", v, nil),
			"domain":  utils.PathSearch("domain", v, nil),
		}
	}

	return rest
}

func GetAppImageServerById(client *golangsdk.ServiceClient, imageServerId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/image-servers"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = fmt.Sprintf("%s?server_id=%s", getPath, imageServerId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	// In any case, the response status code is 200.
	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving image server (%s): %s", imageServerId, err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	imageServer := utils.PathSearch("items|[0]", respBody, nil)
	if imageServer == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return imageServer, nil
}

func resourceAppImageServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// To ensure that the update logic is skipped when the `is_delete_associated_resources` parameters is updated.
	if d.HasChangesExcept("name", "description") {
		return nil
	}

	var (
		cfg           = meta.(*config.Config)
		httpUrl       = "v1/{project_id}/image-servers/{server_id}"
		imageServerId = d.Id()
	)

	client, err := cfg.NewServiceClient("appstream", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{server_id}", imageServerId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"name":        d.Get("name"),
			"description": d.Get("description"),
		},
	}

	_, err = client.Request("PATCH", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating image server (%s) of Workspace APP: %s", imageServerId, err)
	}

	return resourceAppImageServerRead(ctx, d, meta)
}

func resourceAppImageServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		httpUrl       = "v1/{project_id}/image-servers/actions/batch-delete"
		imageServerId = d.Id()
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"items":     []string{imageServerId},
			"recursive": d.Get("is_delete_associated_resources"),
		},
	}

	resp, err := client.Request("PATCH", deletePath, &deleteOpt)
	if err != nil {
		// Although the deletion result of the main region shows that the interface returns a 200 status code when
		// deleting a non-existent image server, in order to avoid the possible return of a 404 status code in the
		// future, the CheckDeleted design is retained here.
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting image server (%s)", imageServerId))
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find job ID from API response")
	}

	_, err = waitForImageServerJobCompleted(ctx, client, d.Timeout(schema.TimeoutDelete), jobId)
	if err != nil {
		return diag.Errorf("error waiting for deleting image server job (%s) completed: %s", jobId, err)
	}
	return nil
}
