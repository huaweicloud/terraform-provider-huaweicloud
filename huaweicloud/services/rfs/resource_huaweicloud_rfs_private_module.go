package rfs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RFS POST /v1/private-modules
// @API RFS DELETE /v1/private-modules/{module_name}
// @API RFS GET /v1/private-modules/{module_name}/metadata
// @API RFS PATCH /v1/private-modules/{module_name}/metadata
func ResourcePrivateModule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateModuleCreate,
		ReadContext:   resourcePrivateModuleRead,
		UpdateContext: resourcePrivateModuleUpdate,
		DeleteContext: resourcePrivateModuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"module_name",
			"module_version",
			"module_uri",
			"version_description",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"module_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			// The `module_version` field was not returned.
			"module_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The `module_uri` field was not returned.
			"module_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The `version_description` field was not returned.
			"version_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The `module_description` cannot be updated to empty.
			"module_description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"module_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreatePrivateModuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"module_name":         d.Get("module_name"),
		"module_version":      utils.ValueIgnoreEmpty(d.Get("module_version")),
		"module_uri":          utils.ValueIgnoreEmpty(d.Get("module_uri")),
		"version_description": utils.ValueIgnoreEmpty(d.Get("version_description")),
		"module_description":  utils.ValueIgnoreEmpty(d.Get("module_description")),
	}
}

func resourcePrivateModuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v1/private-modules"
		moduleName = d.Get("module_name").(string)
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate UUID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Client-Request-Id": uuid,
		},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreatePrivateModuleBodyParams(d)),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating RFS private module: %s", err)
	}

	d.SetId(moduleName)

	return resourcePrivateModuleRead(ctx, d, meta)
}

func QueryPrivateModule(client *golangsdk.ServiceClient, moduleName, uuid string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/private-modules/{module_name}/metadata"
	requestPath = strings.ReplaceAll(requestPath, "{module_name}", moduleName)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Client-Request-Id": uuid,
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourcePrivateModuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate UUID: %s", err)
	}

	respBody, err := QueryPrivateModule(client, d.Id(), uuid)
	if err != nil {
		// If the resource does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving RFS private module")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("module_name", utils.PathSearch("module_name", respBody, nil)),
		d.Set("module_id", utils.PathSearch("module_id", respBody, nil)),
		d.Set("module_description", utils.PathSearch("module_description", respBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", respBody, nil)),
		d.Set("update_time", utils.PathSearch("update_time", respBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePrivateModuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/private-modules/{module_name}/metadata"
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate UUID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{module_name}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Client-Request-Id": uuid,
		},
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"module_description": d.Get("module_description"),
		},
	}

	_, err = client.Request("PATCH", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating RFS private module description: %s", err)
	}

	return resourcePrivateModuleRead(ctx, d, meta)
}

func resourcePrivateModuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/private-modules/{module_name}"
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate UUID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{module_name}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Client-Request-Id": uuid,
		},
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting RFS private module: %s", err)
	}

	return nil
}
