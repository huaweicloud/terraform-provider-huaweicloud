package evs

import (
	"context"
	"errors"
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

// @API EVS DELETE /v2/{project_id}/cloudsnapshots/{snapshot_id}
// @API EVS GET /v2/{project_id}/cloudsnapshots/{snapshot_id}
// @API EVS PUT /v2/{project_id}/cloudsnapshots/{snapshot_id}
// @API EVS POST /v2/{project_id}/cloudsnapshots
func ResourceEvsSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvsSnapshotCreate,
		ReadContext:   resourceEvsSnapshotRead,
		UpdateContext: resourceEvsSnapshotUpdate,
		DeleteContext: resourceEvsSnapshotDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"metadata": {
				Type:             schema.TypeMap,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: utils.SuppressMapDiffs(),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
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
			"metadata_origin": {
				Type:             schema.TypeMap,
				Optional:         true,
				Computed:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: utils.SuppressDiffAll,
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
 the new value next time the change is made. The corresponding parameter name is 'metadata'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildCreateEvsSnapshotBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"volume_id":   d.Get("volume_id"),
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"metadata":    utils.ValueIgnoreEmpty(utils.ExpandToStringMap(d.Get("metadata").(map[string]interface{}))),
	}

	if d.Get("force").(bool) {
		bodyParams["force"] = true
	}

	return map[string]interface{}{
		"snapshot": bodyParams,
	}
}

func GetSnapshotDetail(client *golangsdk.ServiceClient, snapshotID string) (interface{}, error) {
	requestPath := client.Endpoint + "v2/{project_id}/cloudsnapshots/{snapshot_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{snapshot_id}", snapshotID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func waitingForSnapshotStatusAvailable(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := GetSnapshotDetail(client, d.Id())
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("snapshot.status", respBody, "").(string)
			if status == "" {
				return respBody, "ERROR", errors.New("status is not found in API response")
			}

			if status == "available" {
				return respBody, "COMPLETED", nil
			}

			if status == "error" {
				return respBody, status, nil
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceEvsSnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v2/{project_id}/cloudsnapshots"
		product      = "evs"
		mapParamKeys = []string{
			"metadata",
		}
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateEvsSnapshotBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating EVS snapshot: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	snapshotID := utils.PathSearch("snapshot.id", respBody, "").(string)
	if snapshotID == "" {
		return diag.Errorf("error creating EVS snapshot: ID is not found in API response")
	}
	d.SetId(snapshotID)

	if err := waitingForSnapshotStatusAvailable(ctx, client, d, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for EVS snapshot (%s) creation to available: %s", d.Id(), err)
	}

	// If the request is successful, obtain the values of all JSON|object parameters first and save them to the
	// corresponding '_origin' attributes for subsequent determination and construction of the request body during
	// next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	err = utils.RefreshObjectParamOriginValues(d, mapParamKeys)
	if err != nil {
		return diag.Errorf("unable to refresh the origin values: %s", err)
	}

	return resourceEvsSnapshotRead(ctx, d, meta)
}

func resourceEvsSnapshotRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	respBody, err := GetSnapshotDetail(client, d.Id())
	if err != nil {
		// When the resource does not exist, calling the query API will return a `404` status code.
		return common.CheckDeletedDiag(d, err, "error retrieving EVS snapshot")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("volume_id", utils.PathSearch("snapshot.volume_id", respBody, nil)),
		d.Set("name", utils.PathSearch("snapshot.name", respBody, nil)),
		d.Set("metadata", utils.PathSearch("snapshot.metadata", respBody, nil)),
		d.Set("description", utils.PathSearch("snapshot.description", respBody, nil)),
		d.Set("status", utils.PathSearch("snapshot.status", respBody, nil)),
		d.Set("size", utils.PathSearch("snapshot.size", respBody, nil)),
		d.Set("created_at", utils.PathSearch("snapshot.created_at", respBody, nil)),
		d.Set("updated_at", utils.PathSearch("snapshot.updated_at", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateSnapshotBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}

	return map[string]interface{}{
		"snapshot": bodyParams,
	}
}

func resourceEvsSnapshotUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v2/{project_id}/cloudsnapshots/{snapshot_id}"
		product      = "evs"
		mapParamKeys = []string{
			"metadata",
		}
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{snapshot_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateSnapshotBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating EVS snapshot: %s", err)
	}

	// If the request is successful, obtain the values of all JSON|object parameters first and save them to the
	// corresponding '_origin' attributes for subsequent determination and construction of the request body during
	// next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	err = utils.RefreshObjectParamOriginValues(d, mapParamKeys)
	if err != nil {
		return diag.Errorf("unable to refresh the origin values: %s", err)
	}

	return resourceEvsSnapshotRead(ctx, d, meta)
}

func waitingForSnapshotDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := GetSnapshotDetail(client, d.Id())
			if err != nil {
				var errDefault404 golangsdk.ErrDefault404
				if errors.As(err, &errDefault404) {
					return "success deleted", "COMPLETED", nil
				}
				return nil, "ERROR", err
			}

			status := utils.PathSearch("snapshot.status", respBody, "").(string)
			if status == "" {
				return respBody, "ERROR", errors.New("status is not found in API response")
			}

			if status == "error_deleting" {
				return respBody, status, errors.New("an error occurred while deleting the EVS snapshot. " +
					"Please check with your cloud admin or check the API logs " +
					"to see why this error occurred")
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceEvsSnapshotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/cloudsnapshots/{snapshot_id}"
		product = "evs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{snapshot_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		// When the resource does not exist, calling the delete API will return a `404` status code.
		return common.CheckDeletedDiag(d, err, "error deleting EVS snapshot")
	}

	if err := waitingForSnapshotDeleted(ctx, client, d, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for EVS snapshot (%s) deleted: %s", d.Id(), err)
	}

	return nil
}
