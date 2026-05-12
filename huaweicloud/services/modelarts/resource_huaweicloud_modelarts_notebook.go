package modelarts

import (
	"context"
	"errors"
	"fmt"
	"log"
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

var (
	notebookNonUpdatableParams = []string{
		"key_pair",
		"pool_id",
		"workspace_id",
		"volume.*.type",
		"volume.*.ownership",
		"volume.*.uri",
		"volume.*.id",
	}
	notebookNotFoundErrCodes = []string{
		"ModelArts.6309",
		"ModelArts.6357",
		"ModelArts.6404",
	}
)

// @API ModelArts POST /v1/{project_id}/notebooks/{id}/start
// @API ModelArts POST /v1/{project_id}/notebooks/{id}/stop
// @API ModelArts POST /v1/{project_id}/notebooks
// @API ModelArts GET /v1/{project_id}/notebooks/{id}
// @API ModelArts GET /v1/{project_id}/notebooks/{instance_id}/storage
// @API ModelArts PUT /v1/{project_id}/notebooks/{id}
// @API ModelArts DELETE /v1/{project_id}/notebooks/{id}
func ResourceNotebook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotebookCreate,
		ReadContext:   resourceNotebookRead,
		UpdateContext: resourceNotebookUpdate,
		DeleteContext: resourceNotebookDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(notebookNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the notebook is located.`,
			},

			// Required parameters.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the notebook.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The flavor ID of the notebook.`,
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The image ID of the notebook.`,
			},
			"volume": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The storage type.`,
						},
						"ownership": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "MANAGED",
							Description: `The storage ownership.`,
						},
						"size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `The storage size, in GB.`,
						},
						"uri": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The storage URL of the dedicated disk.`,
						},
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The storage ID of the dedicated disk.`,
						},
						"mount_path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							Description: utils.SchemaDesc(
								`The local mount path of the storage.`,
								utils.SchemaDescInput{
									Internal: true,
								},
							),
						},
					},
				},
				Description: `The volume configuration of the notebook.`,
			},

			// Optional parameters.
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the notebook.`,
			},
			"key_pair": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The key pair name for remote SSH access.`,
			},
			"allowed_access_ips": {
				Type:         schema.TypeList,
				Optional:     true,
				Elem:         &schema.Schema{Type: schema.TypeString},
				RequiredWith: []string{"key_pair"},
				Description:  `The public IP addresses that are allowed for remote SSH access.`,
			},
			"pool_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the dedicated resource pool which the notebook used.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The workspace ID to which the notebook belongs.`,
			},

			// Attributes.
			"auto_stop_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the notebook auto stop is enabled.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the notebook.`,
			},
			"image_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the image which the notebook used.`,
			},
			"image_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the image which the notebook used.`,
			},
			"image_swr_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The SWR repository path of the image which the notebook used.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the notebook, in UTC format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the notebook, in UTC format.`,
			},
			"pool_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the dedicated resource pool which the notebook used.`,
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The web URL of the notebook.`,
			},
			"ssh_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URL for remote SSH access of the notebook.`,
			},
			"mount_storages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the storage which is mounted to the notebook.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the storage which is mounted to the notebook.`,
						},
						"mount_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The local mount path of the storage which is mounted to the notebook.`,
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The source path of the storage which is mounted to the notebook.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the storage which is mounted to the notebook.`,
						},
					},
				},
				Description: `The storages which are mounted to the notebook.`,
			},
		},
	}
}

func buildNotebookVolume(volumes []interface{}) map[string]interface{} {
	if len(volumes) < 1 {
		return nil
	}

	category := utils.PathSearch("type", volumes[0], "").(string)
	ownership := utils.PathSearch("ownership", volumes[0], "").(string)
	result := map[string]interface{}{
		"category":  category,
		"ownership": ownership,
		"capacity":  utils.ValueIgnoreEmpty(utils.PathSearch("size", volumes[0], 0).(int)),
	}
	if category == "EFS" && ownership == "DEDICATED" {
		result["uri"] = utils.PathSearch("uri", volumes[0], nil)
		result["id"] = utils.PathSearch("id", volumes[0], nil)
	}

	return result
}

func buildCreateNotebookBodyParams(d *schema.ResourceData) map[string]interface{} {
	result := map[string]interface{}{
		// Fixed values.
		"feature":  "NOTEBOOK",
		"duration": -1,
		// Required parameters.
		"name":     d.Get("name"),
		"flavor":   d.Get("flavor_id"),
		"image_id": d.Get("image_id"),
		"volume":   buildNotebookVolume(d.Get("volume").([]interface{})),
		// Optional parameters.
		"description":  utils.ValueIgnoreEmpty(d.Get("description")),
		"pool_id":      utils.ValueIgnoreEmpty(d.Get("pool_id")),
		"workspace_id": utils.ValueIgnoreEmpty(d.Get("workspace_id")),
	}

	if v, ok := d.GetOk("key_pair"); ok {
		result["endpoints"] = []map[string]interface{}{
			{
				"service":            "SSH",
				"key_pair_names":     []string{v.(string)},
				"allowed_access_ips": d.Get("allowed_access_ips"),
			},
		}
	}

	return result
}

func createNotebook(client *golangsdk.ServiceClient, requestBody map[string]interface{}) (interface{}, error) {
	httpUrl := "v1/{project_id}/notebooks"

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(requestBody),
	}

	requestResp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func refreshNotebookStatus(client *golangsdk.ServiceClient, notebookId string, targetStatus []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetNotebookById(client, notebookId)
		if err != nil {
			parsedErr := common.ConvertExpected400ErrInto404Err(err, "error_code", notebookNotFoundErrCodes...)
			if parsedErr != nil {
				if _, ok := parsedErr.(golangsdk.ErrDefault404); ok && len(targetStatus) < 1 {
					return "Resource Not Found", "COMPLETED", nil
				}
				return respBody, "ERROR", parsedErr
			}
		}

		unexpectedStatus := []string{
			"CREATE_FAILED",
			"START_FAILED",
			"DELETE_FAILED",
			"ERROR",
		}
		if utils.StrSliceContains(unexpectedStatus, utils.PathSearch("status", respBody, "").(string)) {
			return respBody, "ERROR", fmt.Errorf("unexpected status (%s)", utils.PathSearch("status", respBody, "").(string))
		}

		if utils.StrSliceContains(targetStatus, utils.PathSearch("status", respBody, "").(string)) {
			evsVolumePendingStatuses := []string{"INITIALIZING", "RESIZING", "DELETING"}
			evsVolumeStatus := utils.PathSearch("[volume][?category=='EVS'].status|[0]", respBody, "").(string)
			if evsVolumeStatus != "" && utils.StrSliceContains(evsVolumePendingStatuses, evsVolumeStatus) {
				return respBody, "PENDING", nil
			}
			return respBody, "COMPLETED", nil
		}
		return respBody, "PENDING", nil
	}
}

func waitForNotebookStatusCompleted(ctx context.Context, client *golangsdk.ServiceClient, notebookId string, targetStatus []string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshNotebookStatus(client, notebookId, targetStatus),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for ModelArts notebook (%s) status to be completed: %s", notebookId, err)
	}
	return nil
}

func resourceNotebookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	respBody, err := createNotebook(client, buildCreateNotebookBodyParams(d))
	if err != nil {
		return diag.Errorf("error creating ModelArts notebook: %s", err)
	}

	resourceId := utils.PathSearch("id", respBody, "").(string)
	if resourceId == "" {
		return diag.Errorf("unable to find the ModelArts notebook ID from the API response")
	}
	d.SetId(resourceId)

	err = waitForNotebookStatusCompleted(ctx, client, resourceId, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNotebookRead(ctx, d, meta)
}

func GetNotebookById(client *golangsdk.ServiceClient, notebookId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/notebooks/{id}"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", notebookId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func flattenNotebookVolume(volume map[string]interface{}) []interface{} {
	if len(volume) < 1 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"type":       utils.PathSearch("category", volume, nil),
			"ownership":  utils.PathSearch("ownership", volume, nil),
			"size":       utils.PathSearch("capacity", volume, nil),
			"uri":        utils.PathSearch("uri", volume, nil),
			"id":         utils.PathSearch("id", volume, nil),
			"mount_path": utils.PathSearch("mount_path", volume, nil),
		},
	}
}

func listNotebookMountedStorages(client *golangsdk.ServiceClient, notebookId string) ([]interface{}, error) {
	httpUrl := "v1/{project_id}/notebooks/{instance_id}/storage"

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", notebookId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}
	return utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenNotebookMountStorages(mounts []interface{}) []interface{} {
	if len(mounts) < 1 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":         utils.PathSearch("id", mounts[0], nil),
			"type":       utils.PathSearch("category", mounts[0], nil),
			"mount_path": utils.PathSearch("mount_path", mounts[0], nil),
			"path":       utils.PathSearch("uri", mounts[0], nil),
			"status":     utils.PathSearch("status", mounts[0], nil),
		},
	}
}

func resourceNotebookRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		notebookId = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	respBody, err := GetNotebookById(client, notebookId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", notebookNotFoundErrCodes...),
			"error retrieving ModelArts notebook")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		// Required parameters.
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("flavor_id", utils.PathSearch("flavor", respBody, nil)),
		d.Set("image_id", utils.PathSearch("image.id", respBody, nil)),
		d.Set("volume", flattenNotebookVolume(utils.PathSearch("volume", respBody, make(map[string]interface{})).(map[string]interface{}))),
		// Optional parameters.
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("key_pair", utils.PathSearch("endpoints[?service=='SSH']|[0].key_pair_names|[0]", respBody, nil)),
		d.Set("allowed_access_ips", utils.PathSearch("endpoints[?service=='SSH']|[0].allowed_access_ips", respBody, nil)),
		d.Set("pool_id", utils.PathSearch("pool.id", respBody, nil)),
		d.Set("workspace_id", utils.PathSearch("workspace_id", respBody, nil)),
		// Attributes.
		d.Set("auto_stop_enabled", utils.PathSearch("lease.enable", respBody, false)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("image_type", utils.PathSearch("image.type", respBody, nil)),
		d.Set("image_name", utils.PathSearch("image.name", respBody, nil)),
		d.Set("image_swr_path", utils.PathSearch("image.swr_path", respBody, nil)),
		d.Set("created_at", time.Unix(int64(utils.PathSearch("create_at",
			respBody, float64(0)).(float64))/1000, 0).UTC().Format("2006-01-02 15:04:05 MST")),
		d.Set("updated_at", time.Unix(int64(utils.PathSearch("update_at",
			respBody, float64(0)).(float64))/1000, 0).UTC().Format("2006-01-02 15:04:05 MST")),
		d.Set("pool_name", utils.PathSearch("pool.name", respBody, nil)),
		d.Set("url", utils.PathSearch("url", respBody, nil)),
		d.Set("ssh_uri", utils.PathSearch("endpoints[?service=='SSH']|[0].uri", respBody, nil)),
	)

	mountedStorages, err := listNotebookMountedStorages(client, d.Id())
	if err != nil {
		log.Printf("[ERROR] Failed to query the mounted storages of ModelArts notebook (%s): %s", notebookId, err)
	} else {
		mErr = multierror.Append(mErr, d.Set("mount_storages", flattenNotebookMountStorages(mountedStorages)))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateNotebookByRunningStateBodyParams(d *schema.ResourceData) map[string]interface{} {
	result := map[string]interface{}{
		"name":        d.Get("name"),
		"description": d.Get("description"),
	}

	if d.HasChange("volume.0.size") {
		result["storage_new_size"] = d.Get("volume.0.size")
	}

	return result
}

func updateNotebook(client *golangsdk.ServiceClient, notebookId string, requestBody map[string]interface{}) error {
	httpUrl := "v1/{project_id}/notebooks/{id}"

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{id}", notebookId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(requestBody),
	}

	_, err := client.Request("PUT", updatePath, &opt)
	return err
}

func stopNotebook(ctx context.Context, client *golangsdk.ServiceClient, notebookId string, timeout time.Duration) error {
	httpUrl := "v1/{project_id}/notebooks/{id}/stop"

	stopPath := client.Endpoint + httpUrl
	stopPath = strings.ReplaceAll(stopPath, "{project_id}", client.ProjectID)
	stopPath = strings.ReplaceAll(stopPath, "{id}", notebookId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("POST", stopPath, &opt)
	if err != nil {
		return err
	}
	return waitForNotebookStatusCompleted(ctx, client, notebookId, []string{"STOPPED"}, timeout)
}

func buildUpdateNotebookByStoppedStateBodyParams(d *schema.ResourceData) map[string]interface{} {
	result := map[string]interface{}{
		"flavor":   d.Get("flavor_id"),
		"image_id": d.Get("image_id"),
	}

	if v, ok := d.GetOk("key_pair"); ok {
		result["endpoints"] = []map[string]interface{}{
			{
				"service":            "SSH",
				"key_pair_names":     []string{v.(string)},
				"allowed_access_ips": d.Get("allowed_access_ips"),
			},
		}
	}

	return result
}

func buildStartNotebookBodyParams(duration int, actionType string) string {
	res := "?"
	if duration != 0 {
		res += fmt.Sprintf("duration=%d", duration)
	}
	if actionType != "" {
		res += fmt.Sprintf("type=%s", actionType)
	}
	return res
}

func startNotebook(ctx context.Context, client *golangsdk.ServiceClient, notebookId string, timeout time.Duration) error {
	httpUrl := "v1/{project_id}/notebooks/{id}/start"

	startPath := client.Endpoint + httpUrl
	startPath = strings.ReplaceAll(startPath, "{project_id}", client.ProjectID)
	startPath = strings.ReplaceAll(startPath, "{id}", notebookId)
	startPath += buildStartNotebookBodyParams(-1, "")

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("POST", startPath, &opt)
	if err != nil {
		return err
	}
	return waitForNotebookStatusCompleted(ctx, client, notebookId, []string{"RUNNING"}, timeout)
}

func resourceNotebookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		notebookId = d.Id()
	)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	if d.HasChanges("name", "description", "volume.0.size") {
		notebook, err := GetNotebookById(client, notebookId)
		if err != nil {
			return diag.Errorf("error get ModelArts notebook (%s): %s", notebookId, err)
		}
		if utils.PathSearch("status", notebook, "").(string) == "STOPPED" {
			log.Printf("[DEBUG] Try to start notebook (%s) because the status is STOPPED now", notebookId)
			err = startNotebook(ctx, client, notebookId, d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.Errorf("error start ModelArts notebook (%s): %s", notebookId, err)
			}
		}

		err = updateNotebook(client, notebookId, buildUpdateNotebookByRunningStateBodyParams(d))
		if err != nil {
			return diag.Errorf("error update ModelArts notebook (%s): %s", notebookId, err)
		}
		err = waitForNotebookStatusCompleted(ctx, client, notebookId, []string{"RUNNING"}, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for ModelArts notebook (%s) status to be running: %s", notebookId, err)
		}
	}

	if d.HasChanges("flavor_id", "image_id", "allowed_access_ips") {
		// Stop notebook before updating the flavor, image and endpoints
		err := stopNotebook(ctx, client, notebookId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			if !errors.As(common.ConvertExpected400ErrInto404Err(err, "error_code", "ModelArts.6357"), &golangsdk.ErrDefault404{}) {
				return diag.Errorf("error stop ModelArts notebook (%s): %s", notebookId, err)
			}
			log.Printf("[DEBUG] The notebook (%s) is already stopped, so skip the stop operation", notebookId)
		}

		err = updateNotebook(client, notebookId, buildUpdateNotebookByStoppedStateBodyParams(d))
		if err != nil {
			return diag.Errorf("error update ModelArts notebook (%s): %s", notebookId, err)
		}

		// Start notebook after updating the flavor, image and endpoints
		err = startNotebook(ctx, client, notebookId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error start ModelArts notebook (%s): %s", notebookId, err)
		}
	}

	return resourceNotebookRead(ctx, d, meta)
}

func deleteNotebook(ctx context.Context, client *golangsdk.ServiceClient, notebookId string, timeout time.Duration) error {
	httpUrl := "v1/{project_id}/notebooks/{id}"

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{id}", notebookId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &opt)
	if err != nil {
		return err
	}
	return waitForNotebookStatusCompleted(ctx, client, notebookId, nil, timeout)
}

func resourceNotebookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		notebookId = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	err = deleteNotebook(ctx, client, notebookId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", notebookNotFoundErrCodes...),
			fmt.Sprintf("error delete ModelArts notebook (%s)", notebookId))
	}

	return nil
}
