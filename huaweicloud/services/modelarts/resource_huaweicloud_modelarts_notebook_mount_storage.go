package modelarts

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/modelarts/v1/notebook"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func ResourceNotebookMountStorage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotebookMountStorageCreate,
		ReadContext:   resourceNotebookMountStorageRead,
		DeleteContext: resourceNotebookMountStorageDelete,
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
			"notebook_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"local_mount_directory": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^/data/[a-zA-Z0-9_-]*/$"),
					"Only the sub directory of `/data/`can be mounted. e.g. /data/dir1/"),
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mount_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNotebookMountStorageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.ModelArtsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating ModelArts v1 client, err=%s", err)
	}

	notebookId := d.Get("notebook_id").(string)

	//check the notebook instance
	notebookDetail, err := notebook.Get(client, notebookId)
	if err != nil {
		err = parseModlArtsErrorToError404(err)
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			return fmtp.DiagErrorf("The nodebook id=%s is not exist.", notebookId)
		}
		return fmtp.DiagErrorf("Error retrieving ModelArts notebook id=%s", notebookId)
	}

	if notebookDetail.Status != notebook.StatusRunning {
		return fmtp.DiagErrorf("The notebook id=%s must be running, now is =%s", notebookId, notebookDetail.Status)
	}

	// mount storage
	opts := notebook.MountStorageOpts{
		Category:  "OBS",
		MountPath: d.Get("local_mount_directory").(string),
		Uri:       d.Get("storage_path").(string),
	}

	rs, err := notebook.Mount(client, notebookId, opts)
	if err != nil {
		return fmtp.DiagErrorf("Error mounting storage to ModelArts notebook id=%s: %s", notebookId, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", notebookId, rs.Id))
	d.Set("mount_id", rs.Id)

	err = waitingNotebookForMount(ctx, client, notebookId, rs.Id, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNotebookMountStorageRead(ctx, d, meta)
}

func resourceNotebookMountStorageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.ModelArtsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating ModelArts v1 client, err=%s", err)
	}

	notebookId, mountId, err := ParseMountFromId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	detail, err := notebook.GetMount(client, notebookId, mountId)
	if err != nil {
		return common.CheckDeletedDiag(d, parseModlArtsErrorToError404(err), "ModelArts notebook's mount storage")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("notebook_id", notebookId),
		d.Set("mount_id", mountId),
		d.Set("storage_path", detail.Uri),
		d.Set("local_mount_directory", detail.MountPath),
		d.Set("type", detail.Category),
		d.Set("status", detail.Status),
	)

	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("Error setting ModelArts notebook's mount storage fields: %s", mErr)
	}

	return nil
}

func resourceNotebookMountStorageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.ModelArtsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating ModelArts v1 client, err=%s", err)
	}

	notebookId, mountId, err := ParseMountFromId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = notebook.DeleteMount(client, notebookId, mountId)
	if err != nil {
		return fmtp.DiagErrorf("delete ModelArts notebook's mount storage failed. %q:%s", d.Id(), err)
	}

	err = waitingNotebookMountForDeleted(ctx, client, notebookId, mountId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

func waitingNotebookForMount(ctx context.Context, client *golangsdk.ServiceClient, notebookId, mountId string,
	timeout time.Duration) error {
	createStateConf := &resource.StateChangeConf{
		Pending: []string{"MOUNTING"},
		Target:  []string{"MOUNTED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := notebook.GetMount(client, notebookId, mountId)
			if err != nil {
				return nil, "failed", err
			}
			if resp.Status == "MOUNT_FAILED" {
				return nil, "failed", fmtp.Errorf("failed to mount the storage,%s", err)
			}
			return resp, resp.Status, err
		},
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := createStateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for mounting storage to ModelArts notebook (%s): %s", notebookId, err)
	}
	return nil
}

func waitingNotebookMountForDeleted(ctx context.Context, client *golangsdk.ServiceClient, notebookId, mountId string,
	timeout time.Duration) error {
	createStateConf := &resource.StateChangeConf{
		Pending: []string{"UNMOUNTING"},
		Target:  []string{"UNMOUNTED"},
		Refresh: func() (interface{}, string, error) {
			resp, err := notebook.GetMount(client, notebookId, mountId)
			if err != nil {
				err = parseModlArtsErrorToError404(err)
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return resp, "UNMOUNTED", nil
				}
				return nil, "failed", err
			}
			if resp.Status == "UNMOUNT_FAILED" {
				return nil, "failed", fmtp.Errorf("failed to unmount storage, %s", err)
			}
			return resp, resp.Status, err
		},
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := createStateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for ModelArts notebook storage (%s/%s) to be unmounted: %s",
			notebookId, mountId, err)
	}
	return nil
}

func ParseMountFromId(id string) (notebookId, mountId string, err error) {
	idArrays := strings.SplitN(id, "/", 2)
	if len(idArrays) != 2 {
		err = fmtp.Errorf("Invalid format specified for ID. Format must be <notebook_id>/<mount_id>")
		return
	}
	notebookId = idArrays[0]
	mountId = idArrays[1]
	return
}
