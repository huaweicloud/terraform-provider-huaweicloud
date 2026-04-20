package rfs

import (
	"context"
	"fmt"
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

// @API RFS POST /v1/private-providers/{provider_name}/versions
// @API RFS GET /v1/private-providers/{provider_name}/versions/{provider_version}/metadata
// @API RFS DELETE /v1/private-providers/{provider_name}/versions/{provider_version}
func ResourcePrivateProviderVersion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateProviderVersionCreate,
		ReadContext:   resourcePrivateProviderVersionRead,
		UpdateContext: resourcePrivateProviderVersionUpdate,
		DeleteContext: resourcePrivateProviderVersionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePrivateProviderVersionImportState,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"provider_name",
			"provider_version",
			"function_graph_urn",
			"provider_id",
			"version_description",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"function_graph_urn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_id": {
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
			"provider_source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreatePrivateProviderVersionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"provider_version":    d.Get("provider_version"),
		"function_graph_urn":  d.Get("function_graph_urn"),
		"provider_id":         utils.ValueIgnoreEmpty(d.Get("provider_id")),
		"version_description": utils.ValueIgnoreEmpty(d.Get("version_description")),
	}
}

func resourcePrivateProviderVersionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		product      = "rfs"
		httpUrl      = "v1/private-providers/{provider_name}/versions"
		providerName = d.Get("provider_name").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	reqUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate RFS request UUID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{provider_name}", providerName)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Client-Request-Id": reqUUID},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreatePrivateProviderVersionBodyParams(d)),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating RFS private provider version: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", providerName, d.Get("provider_version").(string)))

	return resourcePrivateProviderVersionRead(ctx, d, meta)
}

func resourcePrivateProviderVersionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "rfs"
		httpUrl = "v1/private-providers/{provider_name}/versions/{provider_version}/metadata"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	reqUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate RFS request UUID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{provider_name}", d.Get("provider_name").(string))
	requestPath = strings.ReplaceAll(requestPath, "{provider_version}", d.Get("provider_version").(string))
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Client-Request-Id": reqUUID},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		// If the resource does not exist, the response HTTP status code of the details API is `404`.
		return common.CheckDeletedDiag(d, err, "error retrieving RFS private provider version")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("provider_name", utils.PathSearch("provider_name", respBody, nil)),
		d.Set("provider_version", utils.PathSearch("provider_version", respBody, nil)),
		d.Set("function_graph_urn", utils.PathSearch("function_graph_urn", respBody, nil)),
		d.Set("provider_id", utils.PathSearch("provider_id", respBody, nil)),
		d.Set("version_description", utils.PathSearch("version_description", respBody, nil)),
		d.Set("provider_source", utils.PathSearch("provider_source", respBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePrivateProviderVersionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePrivateProviderVersionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "rfs"
		httpUrl = "v1/private-providers/{provider_name}/versions/{provider_version}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RFS client: %s", err)
	}

	reqUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate RFS request UUID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{provider_name}", d.Get("provider_name").(string))
	requestPath = strings.ReplaceAll(requestPath, "{provider_version}", d.Get("provider_version").(string))
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Client-Request-Id": reqUUID},
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting RFS private provider version: %s", err)
	}

	return nil
}

func resourcePrivateProviderVersionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID,"+
			" want '<provider_name>/<provider_version>', but got '%s'", importedId)
	}

	mErr := multierror.Append(
		d.Set("provider_name", parts[0]),
		d.Set("provider_version", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
