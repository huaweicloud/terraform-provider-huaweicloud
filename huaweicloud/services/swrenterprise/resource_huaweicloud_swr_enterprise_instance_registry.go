package swrenterprise

import (
	"context"
	"strconv"
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

var enterpriseInstanceRegistryNonUpdatableParams = []string{
	"instance_id", "type",
}

// @API SWR POST /v2/{project_id}/instances/{instance_id}/registries
// @API SWR GET /v2/{project_id}/instances/{instance_id}/registries/{registry_id}
// @API SWR PUT /v2/{project_id}/instances/{instance_id}/registries/{registry_id}
// @API SWR DELETE /v2/{project_id}/instances/{instance_id}/registries/{registry_id}
func ResourceSwrEnterpriseInstanceRegistry() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrEnterpriseInstanceRegistryCreate,
		UpdateContext: resourceSwrEnterpriseInstanceRegistryUpdate,
		ReadContext:   resourceSwrEnterpriseInstanceRegistryRead,
		DeleteContext: resourceSwrEnterpriseInstanceRegistryDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(enterpriseInstanceRegistryNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the enterprise instance ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the registry name.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the registry type.`,
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the registry url.`,
			},
			"insecure": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Specifies whether the registry is insecure.`,
			},
			"credential": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `Specifies the credential infos.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the credential type.`,
						},
						"access_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the access key.`,
						},
						"access_secret": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the access secret.`,
						},
					},
				},
			},
			// when `type` value is swr-pro-internal, the replication can be cross-region
			"region_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the region ID of the target registry.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the project ID of the target registry.`,
			},
			"target_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the target enterprise instance ID.`,
			},
			"dns_conf": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Specifies the DNS configuration.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hosts": {
							Type:        schema.TypeMap,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the hosts map.`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the registry description.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"registry_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the registry ID.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the last update time.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the registry status.`,
			},
		},
	}
}

func resourceSwrEnterpriseInstanceRegistryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	createHttpUrl := "v2/{project_id}/instances/{instance_id}/registries"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateSwrEnterpriseInstanceRegistryBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SWR registry: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := int(utils.PathSearch("id", createRespBody, float64(-1)).(float64))
	if id == -1 {
		return diag.Errorf("unable to find SWR instance ID from the API response")
	}

	d.SetId(d.Get("instance_id").(string) + "/" + strconv.Itoa(id))

	return resourceSwrEnterpriseInstanceRegistryRead(ctx, d, meta)
}

func buildCreateOrUpdateSwrEnterpriseInstanceRegistryBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"type":        d.Get("type"),
		"url":         d.Get("url"),
		"insecure":    d.Get("insecure"),
		"credential":  buildSwrEnterpriseInstanceRegistryCredentialBodyParams(d),
		"description": d.Get("description"),
		"region_id":   utils.ValueIgnoreEmpty(d.Get("region_id")),
		"project_id":  utils.ValueIgnoreEmpty(d.Get("project_id")),
		"instance_id": utils.ValueIgnoreEmpty(d.Get("target_instance_id")),
	}

	if d.Get("type").(string) == "swr-pro" {
		bodyParams["dns_conf"] = buildSwrEnterpriseInstanceRegistryDNSConfBodyParams(d)
	}

	return bodyParams
}

func buildSwrEnterpriseInstanceRegistryCredentialBodyParams(d *schema.ResourceData) map[string]interface{} {
	paramsList := d.Get("credential").([]interface{})
	if params, ok := paramsList[0].(map[string]interface{}); ok {
		return map[string]interface{}{
			"type":          params["type"],
			"access_key":    params["access_key"],
			"access_secret": params["access_secret"],
		}
	}

	return nil
}

func buildSwrEnterpriseInstanceRegistryDNSConfBodyParams(d *schema.ResourceData) map[string]interface{} {
	if paramsList, ok := d.Get("dns_conf").([]interface{}); ok && len(paramsList) > 0 {
		if params, ok := paramsList[0].(map[string]interface{}); ok {
			return map[string]interface{}{
				"hosts": params["hosts"],
			}
		}
	}

	return nil
}

func resourceSwrEnterpriseInstanceRegistryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid ID format, want '<instance_id>/<registry_id>', but got '%s'", d.Id())
	}
	instanceId := parts[0]
	id := parts[1]

	getHttpUrl := "v2/{project_id}/instances/{instance_id}/registries/{registry_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{registry_id}", id)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SWR registry")
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getRespBody, nil)),
		d.Set("url", utils.PathSearch("url", getRespBody, nil)),
		d.Set("insecure", utils.PathSearch("insecure", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("credential", flattenSwrEnterpriseInstanceRegistryCredential(getRespBody, d)),
		d.Set("dns_conf", flattenSwrEnterpriseInstanceRegistryDNSConf(getRespBody)),
		d.Set("region_id", utils.PathSearch("region_id", getRespBody, nil)),
		d.Set("project_id", utils.PathSearch("project_id", getRespBody, nil)),
		d.Set("target_instance_id", utils.PathSearch("instance_id", getRespBody, nil)),
		d.Set("registry_id", utils.PathSearch("id", getRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
		d.Set("created_at", utils.PathSearch("created_at", getRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSwrEnterpriseInstanceRegistryCredential(resp interface{}, d *schema.ResourceData) []interface{} {
	params := utils.PathSearch("credential", resp, nil)
	if params == nil {
		return nil
	}

	rst := map[string]interface{}{
		"type":          utils.PathSearch("type", params, nil),
		"access_key":    utils.PathSearch("access_key", params, nil),
		"access_secret": d.Get("credential.0.access_secret"),
	}

	return []interface{}{rst}
}

func flattenSwrEnterpriseInstanceRegistryDNSConf(resp interface{}) []interface{} {
	params := utils.PathSearch("dns_conf", resp, nil)
	if params == nil {
		return nil
	}

	rst := map[string]interface{}{
		"hosts": utils.PathSearch("hosts", params, nil),
	}

	return []interface{}{rst}
}

func resourceSwrEnterpriseInstanceRegistryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	changes := []string{
		"name",
		"url",
		"insecure",
		"credential",
		"region_id",
		"project_id",
		"target_instance_id",
		"dns_conf",
		"description",
	}

	if d.HasChanges(changes...) {
		updateHttpUrl := "v2/{project_id}/instances/{instance_id}/registries/{registry_id}"
		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
		updatePath = strings.ReplaceAll(updatePath, "{registry_id}", strconv.Itoa(d.Get("registry_id").(int)))
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         buildCreateOrUpdateSwrEnterpriseInstanceRegistryBodyParams(d),
		}

		_, err := client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating SWR instance registry: %s", err)
		}
	}

	return resourceSwrEnterpriseInstanceRegistryRead(ctx, d, meta)
}

func resourceSwrEnterpriseInstanceRegistryDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	deleteHttpUrl := "v2/{project_id}/instances/{instance_id}/registries/{registry_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{registry_id}", strconv.Itoa(d.Get("registry_id").(int)))
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SWR registry")
	}

	return nil
}
