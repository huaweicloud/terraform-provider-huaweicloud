package rfs

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RFS POST /v1/private-modules/{module_name}/versions
// @API RFS GET /v1/private-modules/{module_name}/versions
// @API RFS DELETE /v1/private-modules/{module_name}/versions/{module_version}
func ResourcePrivateModuleVersion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateModuleVersionCreate,
		UpdateContext: resourcePrivateModuleVersionUpdate,
		ReadContext:   resourcePrivateModuleVersionRead,
		DeleteContext: resourcePrivateModuleVersionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePrivateModuleVersionImportState,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"module_name",
			"module_version",
			"module_uri",
			"module_id",
			"version_description",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"module_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"module_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"module_uri": {
				Type:     schema.TypeString,
				Required: true,
			},
			"module_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"version_description": {
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
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreatePrivateModuleVersionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"module_version":      d.Get("module_version"),
		"module_uri":          d.Get("module_uri"),
		"module_id":           utils.ValueIgnoreEmpty(d.Get("module_id")),
		"version_description": utils.ValueIgnoreEmpty(d.Get("version_description")),
	}

	return bodyParams
}

func resourcePrivateModuleVersionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v1/private-modules/{module_name}/versions"
		moduleName = d.Get("module_name").(string)
	)
	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	requestId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate RFS request ID: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{module_name}", moduleName)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Client-Request-Id": requestId.String(),
		},
		JSONBody: utils.RemoveNil(buildCreatePrivateModuleVersionBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating RFS private module version: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", moduleName, d.Get("module_version").(string)))

	return resourcePrivateModuleVersionRead(ctx, d, meta)
}

func QueryPrivateModuleVersion(client *golangsdk.ServiceClient, moduleName, moduleVersion, uuid string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/private-modules/{module_name}/versions"
	requestPath = strings.ReplaceAll(requestPath, "{module_name}", moduleName)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Client-Request-Id": uuid,
		},
	}

	nextMarker := ""
	for {
		requestPathWithQueryParams := requestPath + buildListPrivateModuleVersionsQuery(nextMarker)
		resp, err := client.Request("GET", requestPathWithQueryParams, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		targetVersion := utils.PathSearch(fmt.Sprintf("versions[?module_version=='%s'] | [0]", moduleVersion), respBody, nil)
		if targetVersion != nil {
			return targetVersion, nil
		}

		nextMarker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if nextMarker == "" {
			break
		}
	}

	return nil, golangsdk.ErrDefault404{}
}

func resourcePrivateModuleVersionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		moduleName    = d.Get("module_name").(string)
		moduleVersion = d.Get("module_version").(string)
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	requestId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate RFS request ID: %s", err)
	}

	respBody, err := QueryPrivateModuleVersion(client, moduleName, moduleVersion, requestId.String())

	if err != nil {
		// If the resource does not exist, the response HTTP status code of the details API is `404`.
		return common.CheckDeletedDiag(d, err, "error retrieving RFS private module version")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("module_name", utils.PathSearch("module_name", respBody, "")),
		d.Set("module_id", utils.PathSearch("module_id", respBody, "")),
		d.Set("module_version", utils.PathSearch("module_version", respBody, "")),
		d.Set("create_time", utils.PathSearch("create_time", respBody, "")),
		d.Set("version_description", utils.PathSearch("version_description", respBody, "")),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePrivateModuleVersionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePrivateModuleVersionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/private-modules/{module_name}/versions/{module_version}"
	)

	client, err := cfg.NewServiceClient("rfs", region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	requestId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate RFS request ID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{module_name}", d.Get("module_name").(string))
	requestPath = strings.ReplaceAll(requestPath, "{module_version}", d.Get("module_version").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Client-Request-Id": requestId.String(),
		},
	}

	_, err = client.Request("DELETE", requestPath, &opt)
	if err != nil {
		return diag.Errorf("error deleting RFS private module version: %s", err)
	}

	return nil
}

func resourcePrivateModuleVersionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import ID format, expected '<module_name>/<module_version>', but got: %s", d.Id())
	}

	moduleName := parts[0]
	moduleVersion := parts[1]

	mErr := multierror.Append(
		d.Set("module_name", moduleName),
		d.Set("module_version", moduleVersion),
	)

	if mErr.ErrorOrNil() != nil {
		return nil, fmt.Errorf("error setting attributes during import: %s", mErr.ErrorOrNil())
	}

	return []*schema.ResourceData{d}, nil
}

func buildListPrivateModuleVersionsQuery(marker string) string {
	if marker == "" {
		return ""
	}

	return fmt.Sprintf("?marker=%s", marker)
}
