package cci

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var pvcNonUpdatableParams = []string{"namespace", "name", "annotations", "labels", "access_modes",
	"storage_class_name", "volume_mode", "valume_name",
}

// @API CCI POST /apis/cci/v2/namespaces/{namespace}/persistentvolumeclaims
// @API CCI GET /apis/cci/v2/namespaces/{namespace}/persistentvolumeclaims/{name}
// @API CCI PUT /apis/cci/v2/namespaces/{namespace}/persistentvolumeclaims/{name}
// @API CCI DELETE /apis/cci/v2/namespaces/{namespace}/persistentvolumeclaims/{name}
func ResourceV2PersistentVolumeClaim() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2PersistentVolumeClaimCreate,
		UpdateContext: resourceV2PersistentVolumeClaimUpdate,
		ReadContext:   resourceV2PersistentVolumeClaimRead,
		DeleteContext: resourceV2PersistentVolumeClaimDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV2PersistentVolumeClaimImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(pvcNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the namespace.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the persistent volume claim in the namespace.`,
			},
			"annotations": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The annotations of the persistent volume claim.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The labels of the persistent volume claim.`,
			},
			"access_modes": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The access modes of the persistent volume claim.`,
			},
			"resources": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limits": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"requests": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"selector": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match_expressions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"operator": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"values": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"match_labels": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"storage_class_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volume_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"valume_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"creation_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation timestamp of the persistent volume claim.`,
			},
			"resource_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource version of the persistent volume claim.`,
			},
			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uid of the persistent volume claim.`,
			},
			"api_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API version of the persistent volume claim.`,
			},
			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The kind of the persistent volume claim.`,
			},
			"finalizers": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The finalizers of the persistent volume claim.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the persistent volume claim.`,
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

func resourceV2PersistentVolumeClaimCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	createV2PersistentVolumeClaimHttpUrl := "apis/cci/v2/namespaces/{namespace}/persistentvolumeclaims"
	createV2PersistentVolumeClaimPath := client.Endpoint + createV2PersistentVolumeClaimHttpUrl
	createV2PersistentVolumeClaimPath = strings.ReplaceAll(createV2PersistentVolumeClaimPath, "{namespace}", d.Get("namespace").(string))
	createV2PersistentVolumeClaimOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createV2PersistentVolumeClaimOpt.JSONBody = buildCreateV2PersistentVolumeClaimParams(d)

	resp, err := client.Request("POST", createV2PersistentVolumeClaimPath, &createV2PersistentVolumeClaimOpt)
	if err != nil {
		return diag.Errorf("error creating CCI namespace: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	name := utils.PathSearch("metadata.name", respBody, "").(string)
	ns := utils.PathSearch("metadata.namespace", respBody, "").(string)
	if ns == "" {
		return diag.Errorf("unable to find V2PersistentVolumeClaim name from API response")
	}
	d.SetId(ns + "/" + name)

	return resourceV2PersistentVolumeClaimRead(ctx, d, meta)
}

func buildCreateV2PersistentVolumeClaimParams(d *schema.ResourceData) map[string]interface{} {

	bodyParams := map[string]interface{}{
		"metadata": map[string]interface{}{
			"namespace":   d.Get("namespace"),
			"name":        d.Get("name"),
			"annotations": utils.ValueIgnoreEmpty(d.Get("annotations")),
			"labels":      utils.ValueIgnoreEmpty(d.Get("labels")),
		},
		"spec": map[string]interface{}{
			"accessModes":      d.Get("access_modes"),
			"resources":        buildCreateV2PVCResources(d),
			"selector":         buildCreateV2PVCSelector(d),
			"storageClassName": d.Get("storage_class_name"),
			"volumeMode":       d.Get("volume_mode"),
			"valumeName":       d.Get("valume_name"),
		},
	}

	return bodyParams
}

func buildCreateV2PVCResources(d *schema.ResourceData) map[string]interface{} {
	resourcs := d.Get("resources").([]interface{})
	if len(resourcs) == 0 {
		return nil
	}

	res := resourcs[0]
	bodyParams := map[string]interface{}{
		"limits":   utils.PathSearch("limits", res, nil),
		"requests": utils.PathSearch("requests", res, nil),
	}

	return bodyParams
}

func buildCreateV2PVCSelector(d *schema.ResourceData) map[string]interface{} {
	selector := d.Get("selector").([]interface{})
	if len(selector) == 0 {
		return nil
	}

	res := selector[0]
	bodyParams := map[string]interface{}{
		"matchExpressions": buildCreateV2PVCMatchExpressions(res),
		"matchLabels":      utils.PathSearch("match_labels", res, nil),
	}

	return bodyParams
}

func buildCreateV2PVCMatchExpressions(selector interface{}) []map[string]interface{} {
	expressions := utils.PathSearch("match_expressions", selector, nil)
	if expressions == nil {
		return nil
	}

	bodyParams := make([]map[string]interface{}, len(expressions.([]interface{})))

	for i, v := range expressions.([]interface{}) {
		bodyParams[i] = map[string]interface{}{
			"key":      utils.PathSearch("key", v, nil),
			"operator": utils.PathSearch("operator", v, nil),
			"values":   utils.PathSearch("values", v, nil),
		}
	}

	return bodyParams
}

func resourceV2PersistentVolumeClaimRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	namespace := d.Get("namespace").(string)
	name := d.Get("name").(string)
	resp, err := GetV2PersistentVolumeClaimDetail(client, namespace, name)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting the specifies namespace form server")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("api_version", utils.PathSearch("apiVersion", resp, nil)),
		d.Set("kind", utils.PathSearch("kind", resp, nil)),
		d.Set("name", utils.PathSearch("metadata.name", resp, nil)),
		d.Set("creation_timestamp", utils.PathSearch("metadata.creationTimestamp", resp, nil)),
		d.Set("resource_version", utils.PathSearch("metadata.resourceVersion", resp, nil)),
		d.Set("uid", utils.PathSearch("metadata.uid", resp, nil)),
		d.Set("finalizers", utils.PathSearch("spec.finalizers", resp, nil)),
		d.Set("status", utils.PathSearch("status.phase", resp, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceV2PersistentVolumeClaimUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV2PersistentVolumeClaimDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}
	namespace := d.Get("namespace").(string)
	name := d.Get("name").(string)

	deleteV2PersistentVolumeClaimHttpUrl := "apis/cci/v2/namespaces/{namespace}/persistentvolumeclaims/{name}"
	deleteV2PersistentVolumeClaimPath := client.Endpoint + deleteV2PersistentVolumeClaimHttpUrl
	deleteV2PersistentVolumeClaimPath = strings.ReplaceAll(deleteV2PersistentVolumeClaimPath, "{namespace}", namespace)
	deleteV2PersistentVolumeClaimPath = strings.ReplaceAll(deleteV2PersistentVolumeClaimPath, "{name}", name)
	deleteV2PersistentVolumeClaimOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deleteV2PersistentVolumeClaimPath, &deleteV2PersistentVolumeClaimOpt)
	if err != nil {
		return diag.Errorf("error deleting the specifies namespace (%s): %s", namespace, err)
	}

	return nil
}

func GetV2PersistentVolumeClaimDetail(client *golangsdk.ServiceClient, namespace, name string) (interface{}, error) {
	getV2PersistentVolumeClaimDetailHttpUrl := "apis/cci/v2/namespaces/{namespace}/persistentvolumeclaims/{name}"
	getV2PersistentVolumeClaimDetailPath := client.Endpoint + getV2PersistentVolumeClaimDetailHttpUrl
	getV2PersistentVolumeClaimDetailPath = strings.ReplaceAll(getV2PersistentVolumeClaimDetailPath, "{namespace}", namespace)
	getV2PersistentVolumeClaimDetailPath = strings.ReplaceAll(getV2PersistentVolumeClaimDetailPath, "{name}", name)
	getV2PersistentVolumeClaimDetailOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getV2PersistentVolumeClaimDetailResp, err := client.Request("GET", getV2PersistentVolumeClaimDetailPath, &getV2PersistentVolumeClaimDetailOpt)
	if err != nil {
		return getV2PersistentVolumeClaimDetailResp, err
	}

	return utils.FlattenResponse(getV2PersistentVolumeClaimDetailResp)
}

func resourceV2PersistentVolumeClaimImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<namespace>/<name>', but '%s'", importedId)
	}

	d.Set("namespace", parts[0])
	d.Set("name", parts[1])

	return []*schema.ResourceData{d}, nil
}
