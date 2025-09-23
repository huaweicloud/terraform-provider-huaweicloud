package cce

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cce/v1/namespaces"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

// @API CCE DELETE /api/v1/namespaces/{name}
// @API CCE GET /api/v1/namespaces/{name}
// @API CCE POST /api/v1/namespaces
func ResourceCCENamespaceV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCCENamespaceV1Create,
		ReadContext:   resourceCCENamespaceV1Read,
		DeleteContext: resourceCCENamespaceV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCCENamespaceV1Import,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(5 * time.Minute),
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
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"name", "prefix"},
			},
			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"annotations": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func resourceCCENamespaceV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.CceV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud CCE client: %s", err)
	}

	clusterID := d.Get("cluster_id").(string)
	createOpts := namespaces.CreateOpts{
		Kind:       "Namespace",
		ApiVersion: "v1",
		Metadata: namespaces.Metadata{
			Name:         d.Get("name").(string),
			GenerateName: d.Get("prefix").(string),
			Labels:       d.Get("labels").(map[string]interface{}),
			Annotations:  d.Get("annotations").(map[string]interface{}),
		},
	}

	namespace, err := namespaces.Create(client, clusterID, createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud CCE namespace: %s", err)
	}
	d.SetId(namespace.Metadata.UID)
	d.Set("name", namespace.Metadata.Name)

	return resourceCCENamespaceV1Read(ctx, d, meta)
}

func resourceCCENamespaceV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.CceV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud CCE client: %s", err)
	}

	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	namespace, err := namespaces.Get(client, clusterId, name).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error retrieving HuaweiCloud cce namespace")
	}

	labels := map[string]interface{}{}
	for key, val := range namespace.Metadata.Labels {
		labels[key] = val
	}

	annotations := map[string]interface{}{}
	for key, val := range namespace.Metadata.Annotations {
		annotations[key] = val
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("prefix", namespace.Metadata.GenerateName),
		d.Set("labels", labels),
		d.Set("annotations", annotations),
		d.Set("creation_timestamp", namespace.Metadata.CreationTimestamp),
		d.Set("status", namespace.Status.Phase),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("Error setting cce namespace fields: %s", err)
	}
	return nil
}

func resourceCCENamespaceV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.CceV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud CCE Client: %s", err)
	}

	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	namespace, err := namespaces.Delete(client, clusterID, name).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error deleting the specifies CCE namespace (%s): %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Active", "Terminating"},
		Target:       []string{"DELETED"},
		Refresh:      waitForNamepaceDelete(client, clusterID, name),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return fmtp.DiagErrorf(
			"Error waiting for cce namespace (%s) to become DELETED: %s",
			namespace.Metadata.UID, stateErr)
	}

	d.SetId("")
	return nil
}

func waitForNamepaceDelete(client *golangsdk.ServiceClient, clusterID, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		namespace, err := namespaces.Get(client, clusterID, name).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[INFO] Successfully deleted Huaweicloud CCE namespace %s", namespace.Metadata.UID)
				return namespace, "DELETED", nil
			}
			return namespace, "ACTIVE", err
		}
		return namespace, namespace.Status.Phase, nil
	}
}

func resourceCCENamespaceV1Import(context context.Context, d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*config.Config)
	client, err := config.CceV1Client(config.GetRegion(d))
	if err != nil {
		return []*schema.ResourceData{d}, fmtp.Errorf("Error creating HuaweiCloud CCE v1 client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmtp.Errorf("Invalid format specified for import id, must be <cluster_id>/<name>")
	}
	clsuterId := parts[0]
	name := parts[1]
	resp, err := namespaces.Get(client, clsuterId, name).Extract()
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	d.SetId(resp.Metadata.UID)
	d.Set("cluster_id", parts[0])
	d.Set("name", parts[1])

	return []*schema.ResourceData{d}, nil
}
