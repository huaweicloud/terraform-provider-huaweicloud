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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var pvNonUpdatableParams = []string{"name", "csi.*.driver", "csi.*.volume_handle"}

// @API CCI POST /apis/cci/v2/persistentvolumes
// @API CCI GET /apis/cci/v2/persistentvolumes/{name}
// @API CCI PUT /apis/cci/v2/persistentvolumes/{name}
// @API CCI DELETE /apis/cci/v2/persistentvolumes/{name}
func ResourceV2PersistentVolume() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2PersistentVolumeCreate,
		ReadContext:   resourceV2PersistentVolumeRead,
		UpdateContext: resourceV2PersistentVolumeUpdate,
		DeleteContext: resourceV2PersistentVolumeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(pvNonUpdatableParams),

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
			"annotations": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"access_modes": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"capacity": {
				Type:     schema.TypeMap,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"claim_ref": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"field_path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"resource_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"uid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"csi": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     csiSchema(),
			},
			"mount_options": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"node_affinity": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     volumeNodeAffinitySchema(),
			},
			"reclaim_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"storage_class_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"volume_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"finalizers": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"phase": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reason": {
							Type:     schema.TypeString,
							Computed: true,
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
		},
	}
}

func csiSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"driver": {
				Type:     schema.TypeString,
				Required: true,
			},
			"volume_handle": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fs_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"read_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"volume_attributes": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"controller_expand_secret_ref": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     secretReferenceSchema(),
			},
			"controller_publish_secret_ref": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     secretReferenceSchema(),
			},
			"node_expand_secret_ref": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     secretReferenceSchema(),
			},
			"node_publish_secret_ref": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     secretReferenceSchema(),
			},
			"node_stage_secret_ref": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     secretReferenceSchema(),
			},
		},
	}

	return &sc
}

func secretReferenceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}

	return &sc
}

func volumeNodeAffinitySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"required": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     pvNodeSelectorSchema(),
			},
		},
	}

	return &sc
}

func pvNodeSelectorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_selector_terms": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     pvNodeSelectorTermSchema(),
			},
		},
	}

	return &sc
}

func pvNodeSelectorTermSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"match_expressions": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     pvNodeSelectorRequirementSchema(),
			},
		},
	}

	return &sc
}

func pvNodeSelectorRequirementSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"operator": {
				Type:     schema.TypeString,
				Required: true,
			},
			"values": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

	return &sc
}

func resourceV2PersistentVolumeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	createPersistentVolumeHttpUrl := "apis/cci/v2/persistentvolumes"
	createPersistentVolumePath := client.Endpoint + createPersistentVolumeHttpUrl
	createPersistentVolumeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createPersistentVolumeOpt.JSONBody = utils.RemoveNil(buildCreateV2PersistentVolumeParams(d))

	resp, err := client.Request("POST", createPersistentVolumePath, &createPersistentVolumeOpt)
	if err != nil {
		return diag.Errorf("error creating CCI persistent volume: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	name := utils.PathSearch("metadata.name", respBody, "").(string)
	if name == "" {
		return diag.Errorf("unable to find CCI persistent volume name from API response")
	}
	d.SetId(name)

	err = waitForCreateOrUpdatePvStatus(ctx, client, name, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceV2PersistentVolumeRead(ctx, d, meta)
}

func buildCreateV2PersistentVolumeParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name":        d.Get("name"),
			"annotations": utils.ValueIgnoreEmpty(d.Get("annotations")),
			"labels":      utils.ValueIgnoreEmpty(d.Get("labels")),
		},
		"spec": map[string]interface{}{
			"accessModes":                   d.Get("access_modes").(*schema.Set).List(),
			"capacity":                      d.Get("capacity"),
			"claimRef":                      buildClainRefParams(d.Get("claim_ref.0")),
			"csi":                           buildCSIParams(d.Get("csi.0")),
			"mountOptions":                  utils.ValueIgnoreEmpty(d.Get("mount_options").(*schema.Set).List()),
			"nodeAffinity":                  buildPvNodeAffinityParams(d.Get("node_affinity.0")),
			"PersistentVolumeReclaimPolicy": utils.ValueIgnoreEmpty(d.Get("reclaim_policy")),
			"storageClassName":              utils.ValueIgnoreEmpty(d.Get("storage_class_name")),
			"volumeMode":                    utils.ValueIgnoreEmpty(d.Get("volume_mode")),
		},
	}

	return bodyParams
}

func buildClainRefParams(claimRef interface{}) interface{} {
	if claimRef == nil || len(claimRef.(map[string]interface{})) == 0 {
		return nil
	}

	rst := map[string]interface{}{
		"apiVersion":      utils.ValueIgnoreEmpty(utils.PathSearch("api_version", claimRef, nil)),
		"fieldPath":       utils.ValueIgnoreEmpty(utils.PathSearch("field_path", claimRef, nil)),
		"kind":            utils.ValueIgnoreEmpty(utils.PathSearch("kind", claimRef, nil)),
		"name":            utils.ValueIgnoreEmpty(utils.PathSearch("name", claimRef, nil)),
		"namespace":       utils.ValueIgnoreEmpty(utils.PathSearch("namespace", claimRef, nil)),
		"resourceVersion": utils.ValueIgnoreEmpty(utils.PathSearch("resource_version", claimRef, nil)),
		"uid":             utils.ValueIgnoreEmpty(utils.PathSearch("uid", claimRef, nil)),
	}

	return rst
}

func buildCSIParams(csi interface{}) interface{} {
	if csi == nil {
		return nil
	}

	rst := map[string]interface{}{
		"driver":                     utils.PathSearch("driver", csi, nil),
		"volumeHandle":               utils.PathSearch("volume_handle", csi, nil),
		"fsType":                     utils.ValueIgnoreEmpty(utils.PathSearch("fs_type", csi, nil)),
		"readOnly":                   utils.ValueIgnoreEmpty(utils.PathSearch("read_only", csi, nil)),
		"volumeAttributes":           utils.ValueIgnoreEmpty(utils.PathSearch("volume_attributes", csi, nil)),
		"controllerExpandSecretRef":  buildSecretReferenceParams(utils.PathSearch("controller_expand_secret_ref|[0]", csi, nil)),
		"controllerPublishSecretRef": buildSecretReferenceParams(utils.PathSearch("controller_publish_secret_ref|[0]", csi, nil)),
		"nodeExpandSecretRef":        buildSecretReferenceParams(utils.PathSearch("node_expand_secret_ref|[0]", csi, nil)),
		"nodePublishSecretRef":       buildSecretReferenceParams(utils.PathSearch("node_publish_secret_ref|[0]", csi, nil)),
		"nodeStageSecretRef":         buildSecretReferenceParams(utils.PathSearch("node_stage_secret_ref|[0]", csi, nil)),
	}

	return rst
}

func buildSecretReferenceParams(sr interface{}) interface{} {
	if sr == nil {
		return nil
	}

	rst := map[string]interface{}{
		"name":      utils.ValueIgnoreEmpty(utils.PathSearch("name", sr, nil)),
		"namespace": utils.ValueIgnoreEmpty(utils.PathSearch("namespace", sr, nil)),
	}

	return rst
}

func buildPvNodeAffinityParams(nodeAffinity interface{}) interface{} {
	if nodeAffinity == nil || len(nodeAffinity.(map[string]interface{})) == 0 {
		return nil
	}

	rst := map[string]interface{}{
		"required": buildPvNodeSelectorParams(utils.PathSearch("required|[0]", nodeAffinity, nil)),
	}
	return rst
}

func buildPvNodeSelectorParams(nodeSelector interface{}) interface{} {
	if nodeSelector == nil {
		return nil
	}
	nodeSelectorTerms := utils.PathSearch("node_selector_terms", nodeSelector, &schema.Set{}).(*schema.Set).List()
	rst := map[string]interface{}{
		"nodeSelectorTerms": buildPvNodeSelectorTermParams(nodeSelectorTerms),
	}
	return rst
}

func buildPvNodeSelectorTermParams(nodeSelectorTerms []interface{}) []interface{} {
	if len(nodeSelectorTerms) == 0 {
		return nil
	}

	rst := make([]interface{}, len(nodeSelectorTerms))
	for i, v := range nodeSelectorTerms {
		matchExpressions := utils.PathSearch("match_expressions", v, &schema.Set{}).(*schema.Set).List()
		rst[i] = map[string]interface{}{
			"matchExpressions": buildPvMatchExpressionsParams(matchExpressions),
		}
	}
	return rst
}

func buildPvMatchExpressionsParams(matchExpressions []interface{}) []interface{} {
	if len(matchExpressions) == 0 {
		return nil
	}

	rst := make([]interface{}, len(matchExpressions))
	for i, v := range matchExpressions {
		rst[i] = map[string]interface{}{
			"key":      utils.ValueIgnoreEmpty(utils.PathSearch("key", v, nil)),
			"operator": utils.ValueIgnoreEmpty(utils.PathSearch("operator", v, nil)),
			"values":   utils.ValueIgnoreEmpty(utils.PathSearch("values", v, &schema.Set{}).(*schema.Set).List()),
		}
	}
	return rst
}

func resourceV2PersistentVolumeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	resp, err := GetPersistentVolume(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying CCI persistent volume")
	}

	mErr := multierror.Append(
		d.Set("name", utils.PathSearch("metadata.name", resp, nil)),
		d.Set("kind", utils.PathSearch("kind", resp, nil)),
		d.Set("api_version", utils.PathSearch("apiVersion", resp, nil)),
		d.Set("annotations", utils.PathSearch("metadata.annotations", resp, nil)),
		d.Set("labels", utils.PathSearch("metadata.labels", resp, nil)),
		d.Set("creation_timestamp", utils.PathSearch("metadata.creationTimestamp", resp, nil)),
		d.Set("finalizers", utils.PathSearch("metadata.finalizers", resp, nil)),
		d.Set("resource_version", utils.PathSearch("metadata.resourceVersion", resp, nil)),
		d.Set("uid", utils.PathSearch("metadata.uid", resp, nil)),
		d.Set("access_modes", utils.PathSearch("spec.accessModes", resp, nil)),
		d.Set("capacity", utils.PathSearch("spec.capacity", resp, nil)),
		d.Set("claim_ref", flattenClaimRef(utils.PathSearch("spec.claimRef", resp, nil))),
		d.Set("csi", flattenCSI(utils.PathSearch("spec.csi", resp, nil))),
		d.Set("mount_options", utils.PathSearch("spec.mountOptions", resp, nil)),
		d.Set("node_affinity", flattenPvNodeAffinity(utils.PathSearch("spec.nodeAffinity", resp, nil))),
		d.Set("reclaim_policy", utils.PathSearch("spec.reclaimPolicy", resp, nil)),
		d.Set("storage_class_name", utils.PathSearch("spec.storageClassName", resp, nil)),
		d.Set("volume_mode", utils.PathSearch("spec.volumeMode", resp, nil)),
		d.Set("status", flattenPersistentVolumeStatus(utils.PathSearch("status", resp, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPersistentVolumeStatus(status interface{}) []map[string]interface{} {
	if status == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"message": utils.PathSearch("message", status, nil),
			"phase":   utils.PathSearch("phase", status, nil),
			"reason":  utils.PathSearch("reason", status, nil),
		},
	}

	return rst
}

func flattenClaimRef(claimRef interface{}) []map[string]interface{} {
	if claimRef == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"api_version":      utils.PathSearch("apiVersion", claimRef, nil),
			"field_path":       utils.PathSearch("fieldPath", claimRef, nil),
			"kind":             utils.PathSearch("kind", claimRef, nil),
			"name":             utils.PathSearch("name", claimRef, nil),
			"namespace":        utils.PathSearch("namespace", claimRef, nil),
			"resource_version": utils.PathSearch("resourceVersion", claimRef, nil),
			"uid":              utils.PathSearch("uid", claimRef, nil),
		},
	}

	return rst
}

func flattenCSI(csi interface{}) []map[string]interface{} {
	if csi == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"driver":                        utils.PathSearch("driver", csi, nil),
			"volume_handle":                 utils.PathSearch("volumeHandle", csi, nil),
			"fs_type":                       utils.PathSearch("fsType", csi, nil),
			"read_only":                     utils.PathSearch("readOnly", csi, nil),
			"volume_attributes":             utils.PathSearch("volumeAttributes", csi, nil),
			"controller_expand_secret_ref":  flattenCSISecret(utils.PathSearch("controllerExpandSecretRef", csi, nil)),
			"controller_publish_secret_ref": flattenCSISecret(utils.PathSearch("controllerPublishSecretRef", csi, nil)),
			"node_expand_secret_ref":        flattenCSISecret(utils.PathSearch("nodeExpandSecretRef", csi, nil)),
			"node_publish_secret_ref":       flattenCSISecret(utils.PathSearch("nodePublish_secretRef", csi, nil)),
			"node_stage_secret_ref":         flattenCSISecret(utils.PathSearch("nodeStageSecretRef", csi, nil)),
		},
	}

	return rst
}

func flattenCSISecret(sr interface{}) []map[string]interface{} {
	if sr == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"name":      utils.PathSearch("name", sr, nil),
			"namespace": utils.PathSearch("namespace", sr, nil),
		},
	}

	return rst
}

func flattenPvNodeAffinity(nodeAffinity interface{}) []map[string]interface{} {
	if nodeAffinity == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"required": flattenPvNodeAffinityRequired(utils.PathSearch("required", nodeAffinity, nil)),
		},
	}

	return rst
}

func flattenPvNodeAffinityRequired(required interface{}) []map[string]interface{} {
	if required == nil {
		return nil
	}

	nodeSelectorTerms := utils.PathSearch("nodeSelectorTerms", required, make([]interface{}, 0)).([]interface{})
	rst := []map[string]interface{}{
		{
			"node_selector_terms": flattenPvNodeSelectorTerms(nodeSelectorTerms),
		},
	}

	return rst
}

func flattenPvNodeSelectorTerms(nodeSelectorTerms []interface{}) []interface{} {
	if len(nodeSelectorTerms) == 0 {
		return nil
	}

	rst := make([]interface{}, len(nodeSelectorTerms))
	for i, v := range nodeSelectorTerms {
		matchExpressions := utils.PathSearch("matchExpressions", v, make([]interface{}, 0)).([]interface{})
		rst[i] = map[string]interface{}{
			"match_expressions": flattenPvMatchExpressions(matchExpressions),
		}
	}
	return rst
}

func flattenPvMatchExpressions(matchExpressions []interface{}) []interface{} {
	if len(matchExpressions) == 0 {
		return nil
	}

	rst := make([]interface{}, len(matchExpressions))
	for i, v := range matchExpressions {
		rst[i] = map[string]interface{}{
			"key":      utils.PathSearch("key", v, nil),
			"operator": utils.PathSearch("operator", v, nil),
			"values":   utils.PathSearch("values", v, nil),
		}
	}
	return rst
}

func resourceV2PersistentVolumeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	name := d.Get("name").(string)
	updatePersistentVolumeHttpUrl := "apis/cci/v2/persistentvolumes/{name}"
	updatePersistentVolumePath := client.Endpoint + updatePersistentVolumeHttpUrl
	updatePersistentVolumePath = strings.ReplaceAll(updatePersistentVolumePath, "{name}", name)
	updatePersistentVolumeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updatePersistentVolumeOpt.JSONBody = utils.RemoveNil(buildUpdateV2PersistentVolumeParams(d))

	_, err = client.Request("PUT", updatePersistentVolumePath, &updatePersistentVolumeOpt)
	if err != nil {
		return diag.Errorf("error updating CCI persistent volume: %s", err)
	}
	err = waitForCreateOrUpdatePvStatus(ctx, client, name, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceV2PersistentVolumeRead(ctx, d, meta)
}

func buildUpdateV2PersistentVolumeParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"kind":       d.Get("kind"),
		"apiVersion": d.Get("api_version"),
		"metadata": map[string]interface{}{
			"name":              d.Get("name"),
			"uid":               d.Get("uid"),
			"resourceVersion":   d.Get("resource_version"),
			"creationTimestamp": d.Get("creation_timestamp"),
			"annotations":       utils.ValueIgnoreEmpty(d.Get("annotations")),
			"labels":            utils.ValueIgnoreEmpty(d.Get("labels")),
			"finalizers":        utils.ValueIgnoreEmpty(d.Get("finalizers")),
		},
		"spec": map[string]interface{}{
			"accessModes":                   d.Get("access_modes").(*schema.Set).List(),
			"capacity":                      d.Get("capacity"),
			"claimRef":                      buildClainRefParams(d.Get("claim_ref.0")),
			"csi":                           buildCSIParams(d.Get("csi.0")),
			"mountOptions":                  utils.ValueIgnoreEmpty(d.Get("mount_options").(*schema.Set).List()),
			"nodeAffinity":                  buildPvNodeAffinityParams(d.Get("node_affinity.0")),
			"PersistentVolumeReclaimPolicy": utils.ValueIgnoreEmpty(d.Get("reclaim_policy")),
			"storageClassName":              utils.ValueIgnoreEmpty(d.Get("storage_class_name")),
			"volumeMode":                    utils.ValueIgnoreEmpty(d.Get("volume_mode")),
		},
	}

	return bodyParams
}

func resourceV2PersistentVolumeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	name := d.Get("name").(string)
	deletePersistentVolumeHttpUrl := "apis/cci/v2/persistentvolumes/{name}"
	deletePersistentVolumePath := client.Endpoint + deletePersistentVolumeHttpUrl
	deletePersistentVolumePath = strings.ReplaceAll(deletePersistentVolumePath, "{name}", d.Get("name").(string))
	deletePersistentVolumeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deletePersistentVolumePath, &deletePersistentVolumeOpt)
	if err != nil {
		return diag.Errorf("error deleting CCI v2 persistent volume: %s", err)
	}

	err = waitForDeletePvStatus(ctx, client, name, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForCreateOrUpdatePvStatus(ctx context.Context, client *golangsdk.ServiceClient, name string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Available"},
		Refresh: func() (interface{}, string, error) {
			resp, err := GetPersistentVolume(client, name)
			if err != nil {
				return nil, "failed", err
			}
			return resp, utils.PathSearch("status.phase", resp, "").(string), nil
		},
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for the status of the persistent volume to complete: %s", err)
	}
	return nil
}

func waitForDeletePvStatus(ctx context.Context, client *golangsdk.ServiceClient, name string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Deleted"},
		Refresh: func() (interface{}, string, error) {
			resp, err := GetPersistentVolume(client, name)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return "", "Deleted", nil
				}
				return nil, "ERROR", err
			}
			return resp, "Pending", nil
		},
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for the status of the CCI persistent volume to delete: %s", err)
	}
	return nil
}

func GetPersistentVolume(client *golangsdk.ServiceClient, name string) (interface{}, error) {
	getPersistentVolumeHttpUrl := "apis/cci/v2/persistentvolumes/{name}"
	getPersistentVolumePath := client.Endpoint + getPersistentVolumeHttpUrl
	getPersistentVolumePath = strings.ReplaceAll(getPersistentVolumePath, "{name}", name)
	getPersistentVolumeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getPersistentVolumeResp, err := client.Request("GET", getPersistentVolumePath, &getPersistentVolumeOpt)
	if err != nil {
		return getPersistentVolumeResp, err
	}

	return utils.FlattenResponse(getPersistentVolumeResp)
}
