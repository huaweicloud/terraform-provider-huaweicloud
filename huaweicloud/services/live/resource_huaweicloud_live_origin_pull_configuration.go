package live

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Live PUT /v1/{project_id}/domain/pull-sources
// @API Live GET /v1/{project_id}/domain/pull-sources
func ResourceOriginPullConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOriginPullConfigurationCreate,
		ReadContext:   resourceOriginPullConfigurationRead,
		UpdateContext: resourceOriginPullConfigurationUpdate,
		DeleteContext: resourceOriginPullConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceOriginPullConfigurationImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sources": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"sources_ip": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"source_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"scheme": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"additional_args": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildOriginPullConfigurationBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"play_domain":     d.Get("domain_name"),
		"source_type":     d.Get("source_type"),
		"sources":         utils.ExpandToStringList(d.Get("sources").([]interface{})),
		"sources_ip":      utils.ExpandToStringList(d.Get("sources_ip").([]interface{})),
		"source_port":     utils.ValueIgnoreEmpty(d.Get("source_port")),
		"scheme":          utils.ValueIgnoreEmpty(d.Get("scheme")),
		"additional_args": utils.ExpandToStringMap(d.Get("additional_args").(map[string]interface{})),
	}

	return params
}

func updateOriginPullConfiguration(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	validationHttpUrl := "v1/{project_id}/domain/pull-sources"
	validationPath := client.Endpoint + validationHttpUrl
	validationPath = strings.ReplaceAll(validationPath, "{project_id}", client.ProjectID)

	validationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildOriginPullConfigurationBodyParams(d)),
	}

	_, err := client.Request("PUT", validationPath, &validationOpt)
	return err
}

func resourceOriginPullConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	err = updateOriginPullConfiguration(client, d)
	if err != nil {
		return diag.Errorf("error creating Live origin pull configuration: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	return resourceOriginPullConfigurationRead(ctx, d, meta)
}

func ReadOriginPullConfiguration(client *golangsdk.ServiceClient, domainName string) (interface{}, error) {
	getPath := client.Endpoint + "v1/{project_id}/domain/pull-sources"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += fmt.Sprintf("?play_domain=%s", domainName)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		// When the `domain_name` does not exist, calling the query API will return a `400` status code.
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", domainNameNotExistsCode)
	}

	return utils.FlattenResponse(getResp)
}

func resourceOriginPullConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		domainName = d.Get("domain_name").(string)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	getRespBody, err := ReadOriginPullConfiguration(client, domainName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Live origin pull configuration")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("domain_name", domainName),
		d.Set("source_type", utils.PathSearch("source_type", getRespBody, nil)),
		d.Set("sources", flattenSourcesOrSourcesIp(utils.PathSearch("sources", getRespBody, nil))),
		d.Set("sources_ip", flattenSourcesOrSourcesIp(utils.PathSearch("sources_ip", getRespBody, nil))),
		d.Set("source_port", utils.PathSearch("source_port", getRespBody, float64(0))),
		d.Set("scheme", utils.PathSearch("scheme", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSourcesOrSourcesIp(sources interface{}) []string {
	if sources == nil {
		return nil
	}

	result := make([]string, 0)
	for _, v := range sources.([]interface{}) {
		result = append(result, v.(string))
	}

	return result
}

func resourceOriginPullConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	err = updateOriginPullConfiguration(client, d)
	if err != nil {
		return diag.Errorf("error updating Live origin pull configuration: %s", err)
	}

	return resourceOriginPullConfigurationRead(ctx, d, meta)
}

func resourceOriginPullConfigurationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceOriginPullConfigurationImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	if importedId == "" {
		return nil, fmt.Errorf("invalid format specified for import ID, `domain_name` is empty")
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	mErr := multierror.Append(nil,
		d.Set("domain_name", importedId),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
