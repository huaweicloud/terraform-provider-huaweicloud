package cci

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var poolBindingNonUpdatableParams = []string{
	"name", "namespace", "annotations", "labels", "finalizers", "generate_name",
	"owner_references",
	"pool_ref", "pool_ref.*.id",
	"target_ref", "target_ref.*.group", "target_ref.*.kind,name", "target_ref.*.namespace", "target_ref.*.port",
	"api_version", "kind",
}

// @API CCI POST /apis/loadbalancer.networking.openvessel.io/v1/namespaces/{namespace}/poolbindings
// @API CCI GET /apis/loadbalancer.networking.openvessel.io/v1/namespaces/{namespace}/poolbindings/{name}
// @API CCI DELETE /apis/loadbalancer.networking.openvessel.io/v1/namespaces/{namespace}/poolbindings/{name}
func ResourceV2PoolBinding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2PoolBindingCreate,
		ReadContext:   resourceV2PoolBindingRead,
		UpdateContext: resourceV2PoolBindingUpdate,
		DeleteContext: resourceV2PoolBindingDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV2PoolBindingImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(poolBindingNonUpdatableParams),

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
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"annotations": {
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"finalizers": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"generate_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"owner_references": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Type:     schema.TypeString,
							Required: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"uid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"block_owner_deletion": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"controller": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"pool_ref": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"target_ref": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": {
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
							Required: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"api_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "loadbalancer.networking.openvessel.io/v1",
			},
			"kind": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "PoolBinding",
			},
			"generation": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"creation_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deletion_grace_period_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"deletion_timestamp": {
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
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceV2PoolBindingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	createPoolBindingHttpUrl := "apis/loadbalancer.networking.openvessel.io/v1/namespaces/{namespace}/poolbindings"
	createConfigPath := client.Endpoint + createPoolBindingHttpUrl
	createConfigPath = strings.ReplaceAll(createConfigPath, "{namespace}", d.Get("namespace").(string))
	createPoolBindingOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createPoolBindingOpt.JSONBody = utils.RemoveNil(buildCreatePoolBindingParams(d))

	resp, err := client.Request("POST", createConfigPath, &createPoolBindingOpt)
	if err != nil {
		return diag.Errorf("error creating CCI pool binding: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	ns := utils.PathSearch("metadata.namespace", respBody, "").(string)
	name := utils.PathSearch("metadata.name", respBody, "").(string)
	if ns == "" || name == "" {
		return diag.Errorf("unable to find namespace or CCI pool binding name from API response")
	}
	d.SetId(ns + "/" + name)

	return resourceV2PoolBindingRead(ctx, d, meta)
}

func buildCreatePoolBindingParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"kind":       d.Get("kind"),
		"apiVersion": d.Get("api_version"),
		"metadata": map[string]interface{}{
			"name":            d.Get("name"),
			"namespace":       d.Get("namespace"),
			"labels":          utils.ValueIgnoreEmpty(d.Get("labels")),
			"finalizers":      d.Get("finalizers").(*schema.Set).List(),
			"generateName":    utils.ValueIgnoreEmpty(d.Get("generate_name")),
			"ownerReferences": buildCreatePoolBindingOwnerReferencesParams(d.Get("owner_references").(*schema.Set).List()),
		},
		"spec": map[string]interface{}{
			"poolRef":   buildCreatePoolBindingPoolRefParams(d.Get("pool_ref.0")),
			"targetRef": buildCreatePoolBindingTargetRefParams(d.Get("target_ref.0")),
		},
	}

	return bodyParams
}

func buildCreatePoolBindingPoolRefParams(poolRef interface{}) map[string]interface{} {
	if poolRef == nil {
		return nil
	}

	return map[string]interface{}{
		"id": utils.ValueIgnoreEmpty(utils.PathSearch("id", poolRef, nil)),
	}
}

func buildCreatePoolBindingTargetRefParams(targetRef interface{}) map[string]interface{} {
	if targetRef == nil {
		return nil
	}

	return map[string]interface{}{
		"group":     utils.ValueIgnoreEmpty(utils.PathSearch("group", targetRef, nil)),
		"kind":      utils.ValueIgnoreEmpty(utils.PathSearch("kind", targetRef, nil)),
		"name":      utils.PathSearch("name", targetRef, nil),
		"namespace": utils.ValueIgnoreEmpty(utils.PathSearch("namespace", targetRef, nil)),
		"port":      utils.ValueIgnoreEmpty(utils.PathSearch("port", targetRef, nil)),
	}
}

func buildCreatePoolBindingOwnerReferencesParams(ownerReferences []interface{}) []interface{} {
	if len(ownerReferences) == 0 {
		return nil
	}

	params := make([]interface{}, len(ownerReferences))
	for i, v := range ownerReferences {
		params[i] = map[string]interface{}{
			"apiVersion":         utils.ValueIgnoreEmpty(utils.PathSearch("api_version", v, nil)),
			"kind":               utils.ValueIgnoreEmpty(utils.PathSearch("kind", v, nil)),
			"name":               utils.ValueIgnoreEmpty(utils.PathSearch("name", v, nil)),
			"uid":                utils.ValueIgnoreEmpty(utils.PathSearch("uid", v, nil)),
			"blockOwnerDeletion": utils.ValueIgnoreEmpty(utils.PathSearch("block_owner_deletion", v, nil)),
			"controller":         utils.ValueIgnoreEmpty(utils.PathSearch("controller", v, nil)),
		}
	}
	return params
}

func resourceV2PoolBindingRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)
	getPoolBindingHttpUrl := "apis/loadbalancer.networking.openvessel.io/v1/namespaces/{namespace}/poolbindings/{name}"
	getPoolBindingPath := client.Endpoint + getPoolBindingHttpUrl
	getPoolBindingPath = strings.ReplaceAll(getPoolBindingPath, "{namespace}", ns)
	getPoolBindingPath = strings.ReplaceAll(getPoolBindingPath, "{name}", name)
	getPoolBindingOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getPoolBindingResp, err := client.Request("GET", getPoolBindingPath, &getPoolBindingOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting the specifies pool binding form server")
	}

	resp, err := utils.FlattenResponse(getPoolBindingResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("api_version", utils.PathSearch("apiVersion", resp, nil)),
		d.Set("kind", utils.PathSearch("kind", resp, nil)),
		d.Set("namespace", utils.PathSearch("metadata.namespace", resp, nil)),
		d.Set("name", utils.PathSearch("metadata.name", resp, nil)),
		d.Set("annotations", utils.PathSearch("metadata.annotations", resp, nil)),
		d.Set("labels", utils.PathSearch("metadata.labels", resp, nil)),
		d.Set("creation_timestamp", utils.PathSearch("metadata.creationTimestamp", resp, nil)),
		d.Set("resource_version", utils.PathSearch("metadata.resourceVersion", resp, nil)),
		d.Set("uid", utils.PathSearch("metadata.uid", resp, nil)),
		d.Set("finalizers", utils.PathSearch("metadata.finalizers", resp, nil)),
		d.Set("generate_name", utils.PathSearch("metadata.generateName", resp, nil)),
		d.Set("generation", utils.PathSearch("metadata.generation", resp, float64(0)).(float64)),
		d.Set("deletion_grace_period_seconds", utils.PathSearch("metadata.deletionGracePeriodSeconds", resp, float64(0)).(float64)),
		d.Set("deletion_timestamp", utils.PathSearch("metadata.deletionTimestamp", resp, nil)),
		d.Set("owner_references", flattenPoolBindingOwnerReferences(
			utils.PathSearch("metadata.ownerReferences", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("pool_ref", flattenPoolBindingPoolRef(utils.PathSearch("spec.poolRef", resp, nil))),
		d.Set("target_ref", flattenPoolBindingTargetRef(utils.PathSearch("spec.targetRef", resp, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPoolBindingOwnerReferences(ownerReferences []interface{}) []interface{} {
	if len(ownerReferences) == 0 {
		return nil
	}

	rst := make([]interface{}, len(ownerReferences))
	for i, v := range ownerReferences {
		rst[i] = map[string]interface{}{
			"api_version":          utils.PathSearch("apiVersion", v, nil),
			"kind":                 utils.PathSearch("kind", v, nil),
			"name":                 utils.PathSearch("name", v, nil),
			"uid":                  utils.PathSearch("uid", v, nil),
			"block_owner_deletion": utils.PathSearch("blockOwnerDeletion", v, nil),
			"controller":           utils.PathSearch("controller", v, nil),
		}
	}
	return rst
}

func flattenPoolBindingPoolRef(poolRef interface{}) []map[string]interface{} {
	if poolRef == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"id": utils.PathSearch("id", poolRef, nil),
		},
	}
}

func flattenPoolBindingTargetRef(targetRef interface{}) []map[string]interface{} {
	if targetRef == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"group":     utils.PathSearch("group", targetRef, nil),
			"kind":      utils.PathSearch("kind", targetRef, nil),
			"name":      utils.PathSearch("name", targetRef, nil),
			"namespace": utils.PathSearch("namespace", targetRef, nil),
			"port":      utils.PathSearch("port", targetRef, nil),
		},
	}
}

func resourceV2PoolBindingUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV2PoolBindingDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	namespace := d.Get("namespace").(string)
	name := d.Get("name").(string)

	deletePoolBindingHttpUrl := "apis/loadbalancer.networking.openvessel.io/v1/namespaces/{namespace}/poolbindings/{name}"
	deletePoolBindingPath := client.Endpoint + deletePoolBindingHttpUrl
	deletePoolBindingPath = strings.ReplaceAll(deletePoolBindingPath, "{namespace}", namespace)
	deletePoolBindingPath = strings.ReplaceAll(deletePoolBindingPath, "{name}", name)
	deletePoolBindingOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deletePoolBindingPath, &deletePoolBindingOpt)
	if err != nil {
		return diag.Errorf("error deleting the CCI pool binding: %s", err)
	}

	return nil
}

func resourceV2PoolBindingImportState(_ context.Context, d *schema.ResourceData,
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
