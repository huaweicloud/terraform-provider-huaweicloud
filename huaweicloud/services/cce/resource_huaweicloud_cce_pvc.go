package cce

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cce/v1/persistentvolumeclaims"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

type StateRefresh struct {
	Pending      []string
	Target       []string
	Delay        time.Duration
	Timeout      time.Duration
	PollInterval time.Duration
}

// @API CCE DELETE /api/v1/namespaces/{ns}/persistentvolumeclaims/{name}
// @API CCE POST /api/v1/namespaces/{ns}/persistentvolumeclaims
// @API CCE GET /api/v1/namespaces/{ns}/persistentvolumeclaims
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
			"annotations": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"access_modes": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"ReadWriteOnce",
						"ReadOnlyMany",
						"ReadWriteMany",
					}, false),
				},
				Set: schema.HashString,
			},
			"storage": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_class_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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

func resourcePvcLabels(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("labels").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func resourcePvcAnnotations(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("annotations").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func resourcePvcAccessMode(d *schema.ResourceData) []string {
	rawAccessModes := d.Get("access_modes").(*schema.Set).List()
	accessModes := make([]string, len(rawAccessModes))
	for i, raw := range rawAccessModes {
		accessModes[i] = raw.(string)
	}
	return accessModes
}

func buildCcePersistentVolumeClaimCreateOpts(d *schema.ResourceData,
	config *config.Config) (persistentvolumeclaims.CreateOpts, error) {
	createOpts := persistentvolumeclaims.CreateOpts{
		ApiVersion: "v1",                    // The value is fixed at v1.
		Kind:       "PersistentVolumeClaim", // The value is fixed at PersistentVolumeClaim.
		Metadata: persistentvolumeclaims.Metadata{
			Name:        d.Get("name").(string),
			Namespace:   d.Get("namespace").(string),
			Labels:      resourcePvcLabels(d),
			Annotations: resourcePvcAnnotations(d),
		},
		Spec: persistentvolumeclaims.Spec{
			AccessModes:      resourcePvcAccessMode(d),
			StorageClassName: d.Get("storage_class_name").(string),
			// resources
			Resources: persistentvolumeclaims.ResourceRequest{
				Requests: persistentvolumeclaims.CapacityReq{
					Storage: d.Get("storage").(string),
				},
			},
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

func saveCcePersistentVolumeClaimState(d *schema.ResourceData, config *config.Config,
	resp *persistentvolumeclaims.PersistentVolumeClaim) error {
	metadata := &resp.Metadata

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", metadata.Name),
		d.Set("namespace", metadata.Namespace),
		d.Set("access_modes", resp.Spec.AccessModes),
		d.Set("storage_class_name", resp.Spec.StorageClassName),
		d.Set("storage", resp.Spec.Resources.Requests.Storage),
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
	return nil, golangsdk.ErrDefault404{}
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
			if _, ok := err.(golangsdk.ErrDefault404); !ok {
				return resp, "ERROR", nil
			}
		}
		if resp != nil {
			return resp, resp.Status.Phase, nil
		}
		return resp, "DELETED", nil
	}
}

func resourceCcePvcResourceImportState(context context.Context, d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return nil, fmtp.Errorf("Invalid format specified for import id, must be <cluster_id>/<namespace>/<id>")
	}

	clsuterId := parts[0]
	namespace := parts[1]
	id := parts[2]
	d.SetId(id)
	d.Set("cluster_id", clsuterId)
	d.Set("namespace", namespace)

	return []*schema.ResourceData{d}, nil
}
