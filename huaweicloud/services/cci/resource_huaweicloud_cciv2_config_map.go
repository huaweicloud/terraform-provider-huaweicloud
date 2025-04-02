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

var configMapNonUpdatableParams = []string{"name", "namespace"}

// @API CCI POST /apis/cci/v2/namespaces/{namespace}/configmaps
// @API CCI GET /apis/cci/v2/namespaces/{namespace}/configmaps/{name}
// @API CCI PUT /apis/cci/v2/namespaces/{namespace}/configmaps/{name}
// @API CCI DELETE /apis/cci/v2/namespaces/{namespace}/configmaps/{name}
func ResourceV2ConfigMap() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2ConfigMapCreate,
		UpdateContext: resourceV2ConfigMapUpdate,
		ReadContext:   resourceV2ConfigMapRead,
		DeleteContext: resourceV2ConfigMapDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV2ConfigMapImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(configMapNonUpdatableParams),

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
				Description: `Specifies the namespace.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the CCI ConfigMap.`,
			},
			"binary_data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the binary data of the CCI ConfigMap.`,
			},
			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the data of the CCI ConfigMap.`,
			},
			"immutable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `The immutable of the CCI ConfigMap.`,
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The annotations of the CCI ConfigMap.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The labels of the CCI ConfigMap.`,
			},
			"creation_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation timestamp of the CCI ConfigMap.`,
			},
			"resource_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource version of the CCI ConfigMap.`,
			},
			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uid of the CCI ConfigMap.`,
			},
			"api_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API version of the CCI ConfigMap.`,
			},
			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The kind of the CCI ConfigMap.`,
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

func resourceV2ConfigMapCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	createConfigMapHttpUrl := "apis/cci/v2/namespaces/{namespace}/configmaps"
	createConfigPath := client.Endpoint + createConfigMapHttpUrl
	createConfigPath = strings.ReplaceAll(createConfigPath, "{namespace}", d.Get("namespace").(string))
	createConfigMapOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createConfigMapOpt.JSONBody = utils.RemoveNil(buildCreateConfigMapParams(d))

	resp, err := client.Request("POST", createConfigPath, &createConfigMapOpt)
	if err != nil {
		return diag.Errorf("error creating CCI v2 ConfigMap: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	ns := utils.PathSearch("metadata.namespace", respBody, "").(string)
	name := utils.PathSearch("metadata.name", respBody, "").(string)
	if ns == "" || name == "" {
		return diag.Errorf("unable to find namespace or CCI v2 ConfigMap name from API response")
	}
	d.SetId(ns + "/" + name)

	return resourceV2ConfigMapRead(ctx, d, meta)
}

func buildCreateConfigMapParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name": d.Get("name"),
		},
		"data": utils.ValueIgnoreEmpty(d.Get("data")),
	}

	return bodyParams
}

func resourceV2ConfigMapRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)
	getConfigMapHttpUrl := "apis/cci/v2/namespaces/{namespace}/configmaps/{name}"
	getConfigMapPath := client.Endpoint + getConfigMapHttpUrl
	getConfigMapPath = strings.ReplaceAll(getConfigMapPath, "{namespace}", ns)
	getConfigMapPath = strings.ReplaceAll(getConfigMapPath, "{name}", name)
	getConfigMapOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getConfigMapResp, err := client.Request("GET", getConfigMapPath, &getConfigMapOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting the specifies ConfigMap form server")
	}

	resp, err := utils.FlattenResponse(getConfigMapResp)
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
		d.Set("binary_data", utils.PathSearch("binaryData", resp, nil)),
		d.Set("data", utils.PathSearch("data", resp, nil)),
		d.Set("immutable", utils.PathSearch("immutable", resp, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceV2ConfigMapUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	updateConfigMapHttpUrl := "apis/cci/v2/namespaces/{namespace}/configmaps/{name}"
	updateConfigMapPath := client.Endpoint + updateConfigMapHttpUrl
	updateConfigMapPath = strings.ReplaceAll(updateConfigMapPath, "{namespace}", d.Get("namespace").(string))
	updateConfigMapPath = strings.ReplaceAll(updateConfigMapPath, "{name}", d.Get("name").(string))
	updateNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updateNetworkOpt.JSONBody = utils.RemoveNil(buildUpdateV2ConfigMapParams(d))

	_, err = client.Request("PUT", updateConfigMapPath, &updateNetworkOpt)
	if err != nil {
		return diag.Errorf("error updating CCI v2 ConfigMap: %s", err)
	}
	return resourceV2ConfigMapRead(ctx, d, meta)
}

func buildUpdateV2ConfigMapParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name": d.Get("name"),
		},
		"binaryData": utils.ValueIgnoreEmpty(d.Get("binary_data")),
		"data":       utils.ValueIgnoreEmpty(d.Get("data")),
	}

	return bodyParams
}

func resourceV2ConfigMapDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	namespace := d.Get("namespace").(string)
	name := d.Get("name").(string)
	deleteNamespaceHttpUrl := "apis/cci/v2/namespaces/{namespace}/configmaps/{name}"
	deleteNamespacePath := client.Endpoint + deleteNamespaceHttpUrl
	deleteNamespacePath = strings.ReplaceAll(deleteNamespacePath, "{namespace}", namespace)
	deleteNamespacePath = strings.ReplaceAll(deleteNamespacePath, "{name}", name)
	deleteNamespaceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deleteNamespacePath, &deleteNamespaceOpt)
	if err != nil {
		return diag.Errorf("error deleting the specifies CCI v2 ConfigMap (%s): %s", namespace, err)
	}

	return nil
}

func resourceV2ConfigMapImportState(_ context.Context, d *schema.ResourceData,
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
