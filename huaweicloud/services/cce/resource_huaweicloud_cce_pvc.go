package cce

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cce/v1/persistentvolumeclaims"
	"github.com/chnsz/golangsdk/openstack/evs/v2/cloudvolumes"
	"github.com/chnsz/golangsdk/openstack/sfs/v2/shares"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

const (
	// EvsVolume is the type of the EVS disk.
	EvsVolume = "bs"
	// ObsVolume is the type of the OBS bucket.
	ObsVolume = "obs"
	// SfsVolume is the type of the SFS file system.
	SfsVolume = "nfs"
)

var storageClassType = map[string]string{
	"csi-disk": EvsVolume,
	"csi-obs":  ObsVolume,
	"csi-nas":  SfsVolume,
}

type StateRefresh struct {
	Pending      []string
	Target       []string
	Delay        time.Duration
	Timeout      time.Duration
	PollInterval time.Duration
}

func ResourceCcePersistentVolumeClaimsV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCcePersistentVolumeClaimV1Create,
		ReadContext:   resourceCcePersistentVolumeClaimV1Read,
		DeleteContext: resourceCcePersistentVolumeClaimV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCcePvcResourceImportState,
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
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
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
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[a-z0-9][a-z0-9-]{0,61}[a-z0-9]$"),
					"The name consists of 1 to 63 characters, including lowercase letters, digits and hyphens, "+
						"and must start and end with lowercase letters and digits"),
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  EvsVolume,
				ValidateFunc: validation.StringInSlice([]string{
					EvsVolume, ObsVolume, SfsVolume,
				}, false),
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"access_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_timestamp": {
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

func buildCcePersistentVolumeClaimCreateOpts(d *schema.ResourceData,
	config *config.Config) (persistentvolumeclaims.CreateOpts, error) {
	var volumeType = d.Get("volume_type").(string)
	var accessMode string
	switch volumeType {
	case EvsVolume:
		accessMode = "ReadWriteOnce"
	case ObsVolume, SfsVolume:
		accessMode = "ReadWriteMany"
	default:
		return persistentvolumeclaims.CreateOpts{}, fmtp.Errorf("Volumes does not support this type: %s", volumeType)
	}
	createOpts := persistentvolumeclaims.CreateOpts{
		ApiVersion: "v1",                    // The value is fixed at v1.
		Kind:       "PersistentVolumeClaim", // The value is fixed at PersistentVolumeClaim.
		Metadata: persistentvolumeclaims.Metadata{
			Name: d.Get("name").(string),
			Labels: &persistentvolumeclaims.Labels{
				Region:           config.GetRegion(d),
				AvailabilityZone: d.Get("availability_zone").(string),
			},
		},
		Spec: persistentvolumeclaims.Spec{
			VolumeID:    d.Get("volume_id").(string),
			StorageType: volumeType,
			// Only the first value in all selected options is valid.
			AccessModes: []string{accessMode},
		},
	}

	return createOpts, nil
}

func resourceCcePersistentVolumeClaimV1Create(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	c, err := config.CceV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud CCE client: %s", err)
	}

	clusterId := d.Get("cluster_id").(string)
	ns := d.Get("namespace").(string)
	opts, err := buildCcePersistentVolumeClaimCreateOpts(d, config)
	if err != nil {
		return fmtp.DiagErrorf("Unable to build createOpts: %s", err)
	}
	namespace, err := persistentvolumeclaims.Create(c, clusterId, ns, opts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud CCE PVC: %s", err)
	}

	d.SetId(namespace.Metadata.UID)

	stateRef := StateRefresh{
		Pending:      []string{"Pending"},
		Target:       []string{"Bound"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	diagErr := waitForCcePersistentVolumeClaimtateRefresh(ctx,
		buildCcePvcStateRefreshStruct(c, clusterId, ns, d.Id(), stateRef))
	if diagErr != nil {
		return diagErr
	}

	return resourceCcePersistentVolumeClaimV1Read(ctx, d, meta)
}

func setCcePersistentVolumeClaimVolumeParams(d *schema.ResourceData, config *config.Config,
	spec persistentvolumeclaims.SpecResp) error {
	var volumeId, volumeType string

	switch spec.StorageClassName {
	case "csi-disk":
		c, err := config.BlockStorageV2Client(config.GetRegion(d))
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud EVS v2 client: %s", err)
		}
		volumeType = EvsVolume
		pages, err := cloudvolumes.List(c, cloudvolumes.ListOpts{Name: spec.VolumeName}).AllPages()
		if err != nil {
			return err
		}
		volumes, err := cloudvolumes.ExtractVolumes(pages)
		if err != nil {
			return err
		}
		if len(volumes) <= 0 {
			return fmtp.Errorf("EVS disk (%s) not found.", spec.VolumeName)
		}
		volumeId = volumes[0].ID
	case "csi-obs":
		volumeId = spec.VolumeName
		volumeType = ObsVolume
	case "csi-nas":
		c, err := config.SfsV2Client(config.GetRegion(d))
		if err != nil {
			return fmtp.Errorf("Error creating Huaweicloud File Share v2 Client: %s", err)
		}
		systems, err := shares.List(c, shares.ListOpts{Name: spec.VolumeName})
		if len(systems) <= 0 {
			return fmtp.Errorf("SFS file system (%s) not found.", spec.VolumeName)
		}
		volumeId = systems[0].ID
		volumeType = SfsVolume
	default:
		return fmtp.Errorf("Storage Classes does not support this type: %s", spec.StorageClassName)
	}
	return multierror.Append(nil,
		d.Set("volume_id", volumeId),
		d.Set("volume_type", volumeType),
	)
}

func setCcePersistentVolumeClaimAccessMode(d *schema.ResourceData, accessModes []string) error {
	if len(accessModes) > 0 {
		return d.Set("access_mode", accessModes[0])
	}
	return nil
}

func saveCcePersistentVolumeClaimState(d *schema.ResourceData, config *config.Config,
	resp *persistentvolumeclaims.PersistentVolumeClaim) error {
	metadata := &resp.Metadata

	mErr := multierror.Append(nil,
		d.Set("availability_zone", metadata.Labels.AvailabilityZone),
		d.Set("region", metadata.Labels.Region),
		d.Set("name", metadata.Name),
		setCcePersistentVolumeClaimVolumeParams(d, config, resp.Spec),
		setCcePersistentVolumeClaimAccessMode(d, resp.Spec.AccessModes),
		d.Set("creation_timestamp", metadata.CreationTimestamp),
		d.Set("status", &resp.Status.Phase),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}
	return nil
}

func GetCcePvcInfoById(c *golangsdk.ServiceClient, clusterId, namespace,
	id string) (*persistentvolumeclaims.PersistentVolumeClaim, error) {
	pages, err := persistentvolumeclaims.List(c, clusterId, namespace).AllPages()
	if err != nil {
		return nil, err
	}
	responses, err := persistentvolumeclaims.ExtractPersistentVolumeClaims(pages)
	if err != nil {
		return nil, err
	}
	for _, v := range responses {
		if v.Metadata.UID == id {
			return &v, nil
		}
	}
	// PVC has not exist.
	return nil, nil
}

func resourceCcePersistentVolumeClaimV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.CceV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud CCE client: %s", err)
	}

	clusterId := d.Get("cluster_id").(string)
	ns := d.Get("namespace").(string)
	resp, err := GetCcePvcInfoById(client, clusterId, ns, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "pvc")
	}
	if resp != nil {
		if err := saveCcePersistentVolumeClaimState(d, config, resp); err != nil {
			return fmtp.DiagErrorf("Error saving the specifies CCE PVC (%s) to state: %s", d.Id(), err)
		}
	}

	return nil
}

func resourceCcePersistentVolumeClaimV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	c, err := config.CceV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud CCE Client: %s", err)
	}

	clusterId := d.Get("cluster_id").(string)
	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)
	_, err = persistentvolumeclaims.Delete(c, clusterId, ns, name).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error deleting the specifies CCE PVC (%s): %s", d.Id(), err)
	}

	stateRef := StateRefresh{
		Pending:      []string{"Bound"},
		Target:       []string{"DELETED"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 3 * time.Second,
	}
	diagErr := waitForCcePersistentVolumeClaimtateRefresh(ctx, buildCcePvcStateRefreshStruct(c, clusterId, ns, d.Id(), stateRef))
	if diagErr != nil {
		return diagErr
	}

	d.SetId("")
	return nil
}

func buildCcePvcStateRefreshStruct(c *golangsdk.ServiceClient, clusterId, namespace, pvcId string,
	s StateRefresh) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:      s.Pending,
		Target:       s.Target,
		Refresh:      pvcStateRefreshFunc(c, clusterId, namespace, pvcId),
		Timeout:      s.Timeout,
		Delay:        s.Delay,
		PollInterval: s.PollInterval,
	}
}

func waitForCcePersistentVolumeClaimtateRefresh(ctx context.Context, conf *resource.StateChangeConf) diag.Diagnostics {
	_, err := conf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf("Timeout for waiting PVC status become ready: %s", err)
	}
	return nil
}

func pvcStateRefreshFunc(c *golangsdk.ServiceClient, clusterId, namespace, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := GetCcePvcInfoById(c, clusterId, namespace, id)
		if err != nil {
			return resp, "ERROR", nil
		}
		if resp != nil {
			return resp, resp.Status.Phase, nil
		}
		return resp, "DELETED", nil
	}
}

func resourceCcePvcResourceImportState(context context.Context, d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*config.Config)
	c, err := config.CceV1Client(config.GetRegion(d))
	if err != nil {
		return []*schema.ResourceData{d}, fmtp.Errorf("Error creating HuaweiCloud CCE v1 client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return nil, fmtp.Errorf("Invalid format specified for import id, must be <cluster_id>/<namespace>/<name>")
	}

	clsuterId := parts[0]
	namespace := parts[1]
	resp, err := GetCcePvcInfoById(c, clsuterId, namespace, parts[2])
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	d.SetId(resp.Metadata.UID)
	d.Set("cluster_id", parts[0])
	d.Set("namespace", parts[1])

	return []*schema.ResourceData{d}, nil
}
