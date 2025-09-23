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

var secretNonUpdatableParams = []string{"name", "namespace"}

// @API CCI POST /apis/cci/v2/namespaces/{namespace}/secrets
// @API CCI GET /apis/cci/v2/namespaces/{namespace}/secrets/{name}
// @API CCI PUT /apis/cci/v2/namespaces/{namespace}/secrets/{name}
// @API CCI DELETE /apis/cci/v2/namespaces/{namespace}/secrets/{name}
func ResourceV2Secret() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2SecretCreate,
		UpdateContext: resourceV2SecretUpdate,
		ReadContext:   resourceV2SecretRead,
		DeleteContext: resourceV2SecretDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV2SecretImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(secretNonUpdatableParams),

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
				Description: `Specifies the name of the CCI Secret.`,
			},
			"string_data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies string data of the CCI Secret.`,
			},
			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the data of the CCI Secret.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the used to facilitate programmatic handling of secret data.`,
			},
			"immutable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Specifies the immutable of the CCI Secret.`,
			},
			"annotations": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The annotations of the CCI Secret.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The labels of the CCI Secret.`,
			},
			"creation_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation timestamp of the CCI Secret.`,
			},
			"resource_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource version of the CCI Secret.`,
			},
			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uid of the CCI Secret.`,
			},
			"api_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API version of the CCI Secret.`,
			},
			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The kind of the CCI Secret.`,
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

func resourceV2SecretCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	createSecretHttpUrl := "apis/cci/v2/namespaces/{namespace}/secrets"
	createConfigPath := client.Endpoint + createSecretHttpUrl
	createConfigPath = strings.ReplaceAll(createConfigPath, "{namespace}", d.Get("namespace").(string))
	createSecretOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createSecretOpt.JSONBody = utils.RemoveNil(buildCreateSecretParams(d))

	resp, err := client.Request("POST", createConfigPath, &createSecretOpt)
	if err != nil {
		return diag.Errorf("error creating CCI v2 Secret: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	ns := utils.PathSearch("metadata.namespace", respBody, "").(string)
	name := utils.PathSearch("metadata.name", respBody, "").(string)
	if ns == "" || name == "" {
		return diag.Errorf("unable to find namespace or CCI v2 Secret name from API response")
	}
	d.SetId(ns + "/" + name)

	return resourceV2SecretRead(ctx, d, meta)
}

func buildCreateSecretParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name": d.Get("name"),
		},
		"stringData": utils.ValueIgnoreEmpty(d.Get("binary_data")),
		"data":       utils.ValueIgnoreEmpty(d.Get("data")),
		"type":       utils.ValueIgnoreEmpty(d.Get("type")),
	}

	return bodyParams
}

func resourceV2SecretRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)
	getSecretHttpUrl := "apis/cci/v2/namespaces/{namespace}/secrets/{name}"
	getSecretPath := client.Endpoint + getSecretHttpUrl
	getSecretPath = strings.ReplaceAll(getSecretPath, "{namespace}", ns)
	getSecretPath = strings.ReplaceAll(getSecretPath, "{name}", name)
	getSecretOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getSecretResp, err := client.Request("GET", getSecretPath, &getSecretOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting the specifies Secret form server")
	}

	resp, err := utils.FlattenResponse(getSecretResp)
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
		d.Set("string_data", utils.PathSearch("stringData", resp, nil)),
		d.Set("data", utils.PathSearch("data", resp, nil)),
		d.Set("type", utils.PathSearch("type", resp, nil)),
		d.Set("immutable", utils.PathSearch("immutable", resp, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceV2SecretUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	updateSecretHttpUrl := "apis/cci/v2/namespaces/{namespace}/secrets/{name}"
	updateSecretPath := client.Endpoint + updateSecretHttpUrl
	updateSecretPath = strings.ReplaceAll(updateSecretPath, "{namespace}", d.Get("namespace").(string))
	updateSecretPath = strings.ReplaceAll(updateSecretPath, "{name}", d.Get("name").(string))
	updateNetworkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updateNetworkOpt.JSONBody = utils.RemoveNil(buildUpdateV2SecretParams(d))

	_, err = client.Request("PUT", updateSecretPath, &updateNetworkOpt)
	if err != nil {
		return diag.Errorf("error updating CCI v2 Secret: %s", err)
	}
	return resourceV2SecretRead(ctx, d, meta)
}

func buildUpdateV2SecretParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name": d.Get("name"),
		},
		"stringData": utils.ValueIgnoreEmpty(d.Get("string_data")),
		"data":       utils.ValueIgnoreEmpty(d.Get("data")),
		"type":       utils.ValueIgnoreEmpty(d.Get("type")),
	}

	return bodyParams
}

func resourceV2SecretDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}

	namespace := d.Get("namespace").(string)
	name := d.Get("name").(string)

	deleteNamespaceHttpUrl := "apis/cci/v2/namespaces/{namespace}/secrets/{name}"
	deleteNamespacePath := client.Endpoint + deleteNamespaceHttpUrl
	deleteNamespacePath = strings.ReplaceAll(deleteNamespacePath, "{namespace}", namespace)
	deleteNamespacePath = strings.ReplaceAll(deleteNamespacePath, "{name}", name)
	deleteNamespaceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deleteNamespacePath, &deleteNamespaceOpt)
	if err != nil {
		return diag.Errorf("error deleting the specifies CCI v2 Secret: %s", err)
	}

	return nil
}

func resourceV2SecretImportState(_ context.Context, d *schema.ResourceData,
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
