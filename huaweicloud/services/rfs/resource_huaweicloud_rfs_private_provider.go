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

// @API RFS POST /v1/private-providers
// @API RFS GET /v1/private-providers/{provider_name}/metadata
// @API RFS PATCH /v1/private-providers/{provider_name}/metadata
// @API RFS DELETE /v1/private-providers/{provider_name}
func ResourcePrivateProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateProviderCreate,
		ReadContext:   resourcePrivateProviderRead,
		UpdateContext: resourcePrivateProviderUpdate,
		DeleteContext: resourcePrivateProviderDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"provider_name",
			"provider_version",
			"version_description",
			"function_graph_urn",
			"provider_agency_urn",
			"provider_agency_name",
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
			// The `function_graph_urn` field is optional in the API documentation, but it is actually required.
			"function_graph_urn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// The `provider_version` field was not returned.
			"provider_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// The `version_description` field was not returned.
			"version_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"provider_agency_urn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"provider_agency_name": {
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
			"provider_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_source": {
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

func buildCreatePrivateProviderBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"provider_name":        d.Get("provider_name"),
		"provider_description": utils.ValueIgnoreEmpty(d.Get("provider_description")),
		"provider_version":     utils.ValueIgnoreEmpty(d.Get("provider_version")),
		"version_description":  utils.ValueIgnoreEmpty(d.Get("version_description")),
		"function_graph_urn":   utils.ValueIgnoreEmpty(d.Get("function_graph_urn")),
		"provider_agency_urn":  utils.ValueIgnoreEmpty(d.Get("provider_agency_urn")),
		"provider_agency_name": utils.ValueIgnoreEmpty(d.Get("provider_agency_name")),
	}
}

func resourcePrivateProviderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		product      = "rfs"
		httpUrl      = "v1/private-providers"
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
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Client-Request-Id": reqUUID},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreatePrivateProviderBodyParams(d)),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating RFS private provider: %s", err)
	}

	d.SetId(providerName)

	return resourcePrivateProviderRead(ctx, d, meta)
}

func resourcePrivateProviderRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "rfs"
		httpUrl = "v1/private-providers/{provider_name}/metadata"
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
	requestPath = strings.ReplaceAll(requestPath, "{provider_name}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Client-Request-Id": reqUUID},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		// If the resource does not exist, the response HTTP status code of the details API is `404`.
		return common.CheckDeletedDiag(d, err, "error retrieving RFS private provider")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("provider_name", utils.PathSearch("provider_name", respBody, nil)),
		d.Set("provider_description", utils.PathSearch("provider_description", respBody, nil)),
		d.Set("provider_agency_urn", utils.PathSearch("provider_agency_urn", respBody, nil)),
		d.Set("provider_agency_name", utils.PathSearch("provider_agency_name", respBody, nil)),
		d.Set("provider_id", utils.PathSearch("provider_id", respBody, nil)),
		d.Set("provider_source", utils.PathSearch("provider_source", respBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", respBody, nil)),
		d.Set("update_time", utils.PathSearch("update_time", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePrivateProviderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "rfs"
		httpUrl = "v1/private-providers/{provider_name}/metadata"
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
	requestPath = strings.ReplaceAll(requestPath, "{provider_name}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Client-Request-Id": reqUUID},
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"provider_description": d.Get("provider_description")},
	}

	_, err = client.Request("PATCH", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating RFS private provider: %s", err)
	}

	return resourcePrivateProviderRead(ctx, d, meta)
}

func resourcePrivateProviderDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "rfs"
		httpUrl = "v1/private-providers/{provider_name}"
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
	requestPath = strings.ReplaceAll(requestPath, "{provider_name}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Client-Request-Id": reqUUID,
		},
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting RFS private provider: %s", err)
	}

	return nil
}
