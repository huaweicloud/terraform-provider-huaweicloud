package modelarts

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/modelarts/v1/notebook"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts POST /v1/{project_id}/notebooks/{id}/start
// @API ModelArts POST /v1/{project_id}/notebooks/{id}/stop
// @API ModelArts GET /v1/{project_id}/notebooks/{id}/storage
// @API ModelArts DELETE /v1/{project_id}/notebooks/{id}
// @API ModelArts GET /v1/{project_id}/notebooks/{id}
// @API ModelArts PUT /v1/{project_id}/notebooks/{id}
// @API ModelArts POST /v1/{project_id}/notebooks
func ResourceNotebook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotebookCreate,
		ReadContext:   resourceNotebookRead,
		UpdateContext: resourceNotebookUpdate,
		DeleteContext: resourceNotebookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"volume": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"EFS", "EVS"}, false),
						},
						"ownership": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "MANAGED",
							ValidateFunc: validation.StringInSlice([]string{"MANAGED", "DEDICATED"}, false),
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"uri": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"mount_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "schema: Computed",
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_pair": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"allowed_access_ips"},
			},
			"allowed_access_ips": {
				Type:         schema.TypeList,
				Optional:     true,
				Elem:         &schema.Schema{Type: schema.TypeString},
				RequiredWith: []string{"key_pair"},
			},
			"pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"auto_stop_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_swr_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_type": {
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
			"pool_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssh_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mount_storages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mount_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceNotebookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ModelArtsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating ModelArts v1 client, err=%s", err)
	}

	volume, err := buildVolumeParamter(d)
	if err != nil {
		return diag.FromErr(err)
	}
	leaseHour := -1
	opts := notebook.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Feature:     "NOTEBOOK",
		Flavor:      d.Get("flavor_id").(string),
		ImageId:     d.Get("image_id").(string),
		Duration:    &leaseHour,
		PoolId:      d.Get("pool_id").(string),
		WorkspaceId: d.Get("workspace_id").(string),
		Volume:      *volume,
		Endpoints:   buildEndpointsParamter(d),
	}

	rs, err := notebook.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating ModelArts notebook: %s", err)
	}

	d.SetId(rs.Id)

	err = waitingNotebookForRunning(ctx, client, rs.Id, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNotebookRead(ctx, d, meta)
}

func resourceNotebookRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ModelArtsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating ModelArts v1 client, err=%s", err)
	}

	detail, err := notebook.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, parseModlArtsErrorToError404(err), "error retrieving ModelArts notebook")
	}

	keyPair, uri, ips := parseEndpoints(detail.Endpoints)
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", detail.Name),
		d.Set("flavor_id", detail.Flavor),
		d.Set("image_id", detail.Image.Id),
		d.Set("description", detail.Description),
		d.Set("pool_id", detail.Pool.Id),
		d.Set("workspace_id", detail.WorkspaceId),
		d.Set("status", detail.Status),
		d.Set("image_name", detail.Image.Name),
		d.Set("image_type", detail.Image.Type),
		d.Set("image_swr_path", detail.Image.SwrPath),
		d.Set("created_at", time.Unix(int64(detail.CreateAt)/1000, 0).UTC().Format("2006-01-02 15:04:05 MST")),
		d.Set("updated_at", time.Unix(int64(detail.UpdateAt)/1000, 0).UTC().Format("2006-01-02 15:04:05 MST")),
		d.Set("auto_stop_enabled", detail.Lease.Enable),
		d.Set("pool_name", detail.Pool.Name),
		d.Set("url", detail.Url),
		d.Set("key_pair", keyPair),
		d.Set("allowed_access_ips", ips),
		d.Set("ssh_uri", uri),
		setVolumeToState(d, detail.Volume),
		setMountStoragesToState(d, client),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceNotebookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ModelArtsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating ModelArts v1 client, err=%s", err)
	}

	if d.HasChanges("name", "description", "endpoints") {
		desc := d.Get("description").(string)
		opts := notebook.UpdateOpts{
			Name:        d.Get("name").(string),
			Description: &desc,
			Endpoints:   buildEndpointsParamter(d),
		}

		_, err := notebook.Update(client, d.Id(), opts)
		if err != nil {
			return diag.Errorf("error update ModelArts notebook: %s", err)
		}
	}

	if d.HasChanges("flavor_id", "image_id", "volume.0.size") {
		// stop
		status := d.Get("status").(string)
		if status != notebook.StatusStopped {
			_, err := notebook.Stop(client, d.Id())
			if err != nil {
				return diag.FromErr(err)
			}

			err = waitingNotebookForStopped(ctx, client, d.Id(), d.Timeout(schema.TimeoutUpdate))
			if err != nil {
				return diag.FromErr(err)
			}
		}

		// change
		storageSize := d.Get("volume.0.size").(int)
		opts := notebook.UpdateOpts{
			Flavor:         d.Get("flavor_id").(string),
			ImageId:        d.Get("image_id").(string),
			StorageNewSize: &storageSize,
		}

		_, err := notebook.Update(client, d.Id(), opts)
		if err != nil {
			return diag.Errorf("error update ModelArts notebook: %s", err)
		}

		// start the instance
		_, err = notebook.Start(client, d.Id(), -1)
		if err != nil {
			return diag.FromErr(err)
		}

		err = waitingNotebookForRunning(ctx, client, d.Id(), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceNotebookRead(ctx, d, meta)
}

func resourceNotebookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ModelArtsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating ModelArts v1 client, err=%s", err)
	}

	_, err = notebook.Delete(client, d.Id())
	if err != nil {
		return diag.Errorf("delete ModelArts notebook failed. %q:%s", d.Id(), err)
	}

	err = waitingNotebookForDeleted(ctx, client, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildVolumeParamter(d *schema.ResourceData) (*notebook.VolumeReq, error) {
	rst := notebook.VolumeReq{
		Category:  d.Get("volume.0.type").(string),
		Ownership: d.Get("volume.0.ownership").(string),
	}

	if rst.Category == "EFS" && rst.Ownership == "DEDICATED" {
		v, ok := d.GetOk("volume.0.uri")
		if !ok {
			return nil, fmt.Errorf("the parameter 'uri' is mandatory if the storage type is EFS and ownership is DEDICATED")
		}
		rst.Uri = v.(string)
		v, ok = d.GetOk("volume.0.id")
		if !ok {
			return nil, fmt.Errorf("the parameter 'id' is mandatory if the storage type is EFS and ownership is DEDICATED")
		}
		rst.ID = v.(string)
	}

	if v, ok := d.GetOk("volume.0.size"); ok {
		capacity := v.(int)
		rst.Capacity = &capacity
	}

	return &rst, nil
}

func buildEndpointsParamter(d *schema.ResourceData) []notebook.EndpointsReq {
	if v, ok := d.GetOk("key_pair"); ok {
		endpoint := notebook.EndpointsReq{
			Service:          "SSH",
			AllowedAccessIps: utils.ExpandToStringList(d.Get("allowed_access_ips").([]interface{})),
			KeyPairNames:     []string{v.(string)},
		}
		return []notebook.EndpointsReq{endpoint}
	}
	return nil
}

func setVolumeToState(d *schema.ResourceData, volume notebook.VolumeRes) error {
	result := make(map[string]interface{})
	result["type"] = volume.Category
	result["ownership"] = volume.Ownership
	result["size"] = volume.Capacity
	result["uri"] = volume.URI
	result["id"] = volume.ID
	result["mount_path"] = volume.MountPath
	return d.Set("volume", []map[string]interface{}{result})
}

func parseEndpoints(configs []notebook.Endpoints) (keyPair, uri string, ips []string) {
	for _, v := range configs {
		if v.Service == "SSH" {
			return v.KeyPairNames[0], v.Uri, v.AllowedAccessIps
		}
	}
	return
}

func setMountStoragesToState(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	resp, err := notebook.ListMounts(client, d.Id())
	if err != nil {
		log.Printf("[ERROR] Failed to query the mount storage of ModelArts notebook instance=%s", d.Id())
		return nil
	}
	rst := make([]map[string]interface{}, len(resp.Data))
	for i, v := range resp.Data {
		storage := make(map[string]interface{})
		storage["id"] = v.Id
		storage["type"] = v.Category
		storage["mount_path"] = v.MountPath
		storage["path"] = v.Uri
		storage["status"] = v.Status
		rst[i] = storage
	}
	return d.Set("mount_storages", rst)
}

func waitingNotebookForRunning(ctx context.Context, client *golangsdk.ServiceClient, id string, timeout time.Duration) error {
	createStateConf := &resource.StateChangeConf{
		Pending: []string{notebook.StatusInit, notebook.StatusCreating, notebook.StatusStarting, notebook.StatusSnapshotting},
		Target:  []string{notebook.StatusRunning},
		Refresh: func() (interface{}, string, error) {
			resp, err := notebook.Get(client, id)
			if err != nil {
				return nil, "failed", err
			}
			if resp.Status == notebook.StatusCreateFailed || resp.Status == notebook.StatusError {
				return nil, "failed", fmt.Errorf("error_code: %s, error_msg: %s", resp.Status, resp.FailReason)
			}
			return resp, resp.Status, err
		},
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := createStateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for ModelArts notebook (%s) to be created: %s", id, err)
	}
	return nil
}

func waitingNotebookForStopped(ctx context.Context, client *golangsdk.ServiceClient, id string, timeout time.Duration) error {
	createStateConf := &resource.StateChangeConf{
		Pending: []string{notebook.StatusRunning, notebook.StatusStopping},
		Target:  []string{notebook.StatusStopped},
		Refresh: func() (interface{}, string, error) {
			resp, err := notebook.Get(client, id)
			if err != nil {
				return nil, "failed", err
			}
			if resp.Status == notebook.StatusError {
				return nil, "failed", fmt.Errorf("error_code: %s, error_msg: %s", resp.Status, resp.FailReason)
			}
			return resp, resp.Status, err
		},
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := createStateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for ModelArts notebook (%s) to be stopped: %s", id, err)
	}
	return nil
}

func waitingNotebookForDeleted(ctx context.Context, client *golangsdk.ServiceClient, id string, timeout time.Duration) error {
	createStateConf := &resource.StateChangeConf{
		Pending: []string{notebook.StatusRunning, notebook.StatusDeleting},
		Target:  []string{notebook.StatusDeleted},
		Refresh: func() (interface{}, string, error) {
			resp, err := notebook.Get(client, id)
			if err != nil {
				err = parseModlArtsErrorToError404(err)
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return resp, notebook.StatusDeleted, nil
				}
				return nil, "failed", err
			}
			if resp.Status == notebook.StatusError || resp.Status == notebook.StatusDeleteFailed {
				return nil, "failed", fmt.Errorf("error_code: %s, error_msg: %s", resp.Status, resp.FailReason)
			}
			return resp, resp.Status, err
		},
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := createStateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for ModelArts notebook (%s) to be deleted: %s", id, err)
	}
	return nil
}

func parseModlArtsErrorToError404(respErr error) error {
	var apiError notebook.ModelartsError
	if errCode, ok := respErr.(golangsdk.ErrDefault400); ok {
		pErr := json.Unmarshal(errCode.Body, &apiError)
		if pErr == nil && (apiError.ErrorCode == "ModelArts.6309" || apiError.ErrorCode == "ModelArts.6404") {
			return golangsdk.ErrDefault404(errCode)
		}
	}
	return respErr
}
