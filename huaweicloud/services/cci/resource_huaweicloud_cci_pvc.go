package cci

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
	"github.com/chnsz/golangsdk/openstack/cci/v1/persistentvolumeclaims"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

var (
	fsType = map[string]string{
		"sas":             "ext4",
		"ssd":             "ext4",
		"sata":            "ext4",
		"nfs-rw":          "nfs",
		"efs-performance": "nfs",
		"efs-standard":    "nfs",
		"obs":             "obs",
	}
	volumeTypeForList = map[string]string{
		"sas":             "bs",
		"ssd":             "bs",
		"sata":            "bs",
		"obs":             "obs",
		"nfs-rw":          "nfs",
		"efs-performance": "efs",
		"efs-standard":    "efs",
	}
)

type StateRefresh struct {
	Pending      []string
	Target       []string
	Delay        time.Duration
	Timeout      time.Duration
	PollInterval time.Duration
}

// @API CCI GET /api/v1/namespaces/{ns}/extended-persistentvolumeclaims
// @API CCI POST /api/v1/namespaces/{ns}/extended-persistentvolumeclaims
// @API CCI DELETE /api/v1/namespaces/{ns}/persistentvolumeclaims/{name}
func ResourcePersistentVolumeClaimV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePersistentVolumeClaimV1Create,
		ReadContext:   resourcePersistentVolumeClaimV1Read,
		DeleteContext: resourcePersistentVolumeClaimV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePersistentVolumeClaimV1ImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"device_mount_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "sas",
			},
			"access_modes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func buildPersistentVolumeClaimV1CreateParams(d *schema.ResourceData) (persistentvolumeclaims.CreateOpts, error) {
	createOpts := persistentvolumeclaims.CreateOpts{
		Kind:       "PersistentVolumeClaim",
		ApiVersion: "v1",
	}
	volumeType := d.Get("volume_type").(string)
	fsType, ok := fsType[volumeType]
	if !ok {
		return createOpts, fmt.Errorf("the volume type (%s) is not available", volumeType)
	}
	createOpts.Metadata = persistentvolumeclaims.Metadata{
		Namespace: d.Get("namespace").(string),
		Name:      d.Get("name").(string),
		Annotations: &persistentvolumeclaims.Annotations{
			FsType:          fsType,
			VolumeID:        d.Get("volume_id").(string),
			DeviceMountPath: d.Get("device_mount_path").(string),
		},
	}
	createOpts.Spec = persistentvolumeclaims.Spec{
		StorageClassName: volumeType,
		Resources: persistentvolumeclaims.ResourceRequirement{
			Requests: &persistentvolumeclaims.ResourceName{
				// At present, due to design defects of the CCI service, the storage has no practical meaning.
				Storage: "1Gi",
			},
		},
	}

	return createOpts, nil
}

func resourcePersistentVolumeClaimV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.CciV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCI v1 client: %s", err)
	}
	createOpts, err := buildPersistentVolumeClaimV1CreateParams(d)
	if err != nil {
		return diag.Errorf("unable to build createOpts of the PVC: %s", err)
	}
	namespace := d.Get("namespace").(string)
	create, err := persistentvolumeclaims.Create(client, createOpts, namespace).Extract()
	if err != nil {
		return diag.Errorf("error creating CCI PVC: %s", err)
	}
	d.SetId(create.Metadata.UID)
	stateRef := StateRefresh{
		Pending:      []string{"Pending"},
		Target:       []string{"Bound"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	if err := waitForPersistentVolumeClaimStateRefresh(ctx, d, client, namespace, stateRef); err != nil {
		return diag.Errorf("create the specifies PVC (%s) timed out: %s", d.Id(), err)
	}

	return resourcePersistentVolumeClaimV1Read(ctx, d, meta)
}

func savePersistentVolumeClaimV1State(d *schema.ResourceData, resp *persistentvolumeclaims.ListResp) error {
	spec := &resp.PersistentVolume.Spec
	metadata := &resp.PersistentVolumeClaim.Metadata
	mErr := multierror.Append(nil,
		d.Set("namespace", metadata.Namespace),
		d.Set("name", metadata.Name),
		d.Set("volume_id", spec.FlexVolume.Options.VolumeID),
		d.Set("volume_type", spec.StorageClassName),
		d.Set("device_mount_path", spec.FlexVolume.Options.DeviceMountPath),
		d.Set("access_modes", spec.AccessModes),
		d.Set("status", resp.PersistentVolumeClaim.Status.Phase),
		d.Set("creation_timestamp", metadata.CreationTimestamp),
		d.Set("enable", metadata.Enable),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}

	return nil
}

func resourcePersistentVolumeClaimV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.CciV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCI v1 client: %s", err)
	}

	namespace := d.Get("namespace").(string)
	volumeType := d.Get("volume_type").(string)
	id := d.Id()

	response, err := GetPvcInfoById(client, namespace, volumeType, id)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting the specifies PVC form server")
	}
	if response != nil {
		d.Set("region", region)
		if err := savePersistentVolumeClaimV1State(d, response); err != nil {
			return diag.Errorf("error saving the specifies PVC (%s) to state: %s", id, err)
		}
	}

	return nil
}

func resourcePersistentVolumeClaimV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.CciV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCI v1 Client: %s", err)
	}

	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)
	_, err = persistentvolumeclaims.Delete(client, namespace, name).Extract()
	if err != nil {
		return diag.Errorf("error deleting the specifies PVC (%s): %s", d.Id(), err)
	}

	stateRef := StateRefresh{
		Pending:      []string{"Bound"},
		Target:       []string{"DELETED"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        3 * time.Second,
		PollInterval: 2 * time.Second,
	}
	if err := waitForPersistentVolumeClaimStateRefresh(ctx, d, client, namespace, stateRef); err != nil {
		return diag.Errorf("delete the specifies PVC (%s) timed out: %s", d.Id(), err)
	}

	return nil
}

func waitForPersistentVolumeClaimStateRefresh(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	ns string, s StateRefresh) error {
	stateConf := &resource.StateChangeConf{
		Pending:      s.Pending,
		Target:       s.Target,
		Refresh:      pvcStateRefreshFunc(d, client, ns),
		Timeout:      s.Timeout,
		Delay:        s.Delay,
		PollInterval: s.PollInterval,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for the status of the PVC (%s) to complete timeout: %s", d.Id(), err)
	}
	return nil
}

func pvcStateRefreshFunc(d *schema.ResourceData, client *golangsdk.ServiceClient,
	ns string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		volumeType := d.Get("volume_type").(string)
		response, err := GetPvcInfoById(client, ns, volumeType, d.Id())
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return response, "DELETED", nil
			}
			return response, "ERROR", nil
		}
		if response != nil {
			return response, response.PersistentVolumeClaim.Status.Phase, nil
		}
		return response, "ERROR", nil
	}
}

func GetPvcInfoById(client *golangsdk.ServiceClient, ns, volumeType,
	id string) (*persistentvolumeclaims.ListResp, error) {
	// If the storage of listOpts is not set, the list method will search for all PVCs of evs type.
	storageType, ok := volumeTypeForList[volumeType]
	if !ok {
		return nil, golangsdk.ErrDefault400{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte(fmt.Sprintf("the volume type (%s) is not available", volumeType)),
			},
		}
	}

	listOpts := persistentvolumeclaims.ListOpts{
		StorageType: storageType,
	}
	pages, err := persistentvolumeclaims.List(client, listOpts, ns).AllPages()
	if err != nil {
		return nil, err
	}
	responses, err := persistentvolumeclaims.ExtractPersistentVolumeClaims(pages)
	if err != nil {
		return nil, err
	}
	for _, v := range responses {
		if v.PersistentVolumeClaim.Metadata.UID == id {
			return &v, nil
		}
	}
	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Method:    "GET",
			URL:       "/api/v1/namespaces/{ns}/extended-persistentvolumeclaims",
			RequestId: "NONE",
			Body:      []byte(fmt.Sprintf("the PVC (%s) does not exist", id)),
		},
	}
}

func resourcePersistentVolumeClaimV1ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for CCI PVC, must be <namespace>/<volume type>/<pvc id>")
	}
	d.SetId(parts[2])
	mErr := multierror.Append(nil,
		d.Set("namespace", parts[0]),
		d.Set("volume_type", parts[1]),
	)
	if mErr.ErrorOrNil() != nil {
		return []*schema.ResourceData{d}, mErr
	}
	return []*schema.ResourceData{d}, nil
}
