package iam

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM POST /v5/saml-providers
// @API IAM GET /v5/saml-providers
// @API IAM GET /v5/saml-providers/{provider_id}
// @API IAM PUT /v5/saml-providers/{provider_id}
// @API IAM DELETE /v5/saml-providers/{provider_id}
func ResourceV5SamlProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV5SamlProviderCreate,
		ReadContext:   resourceV5SamlProviderRead,
		UpdateContext: resourceV5SamlProviderUpdate,
		DeleteContext: resourceV5SamlProviderDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the SAML identity provider. The name contains 1 to 64 characters and can only contain letters, digits, underscores (_) and hyphens (-).`,
			},
			"saml_metadata_document": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Specifies the SAML metadata document obtained from the identity provider. The document is an XML string that contains the issuer name, ` +
					`the validity period and the signing certificate used to verify SAML assertions.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the SAML identity provider. The description contains a maximum of 255 characters and cannot contain the following characters: @, #, %, &, <, >, ", $, ^ and *.`,
			},
			"provider_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the SAML identity provider.`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uniform resource name of the SAML identity provider.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the SAML identity provider.`,
			},
			"expire_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The expiration time of the SAML identity provider.`,
			},
		},
	}
}

func resourceV5SamlProviderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	createHttpUrl := "v5/saml-providers"
	createPath := client.Endpoint + createHttpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200, 201},
		JSONBody:         buildCreateV5SamlProviderBodyParams(d),
	}
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SAML identity provider: %s", err)
	}

	createBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	providerId := getV5SamlProviderIdFromResponse(createBody)
	if providerId == "" {
		// If the create response does not return the provider ID, query it by the name.
		providerId, err = getV5SamlProviderIdByName(client, d.Get("name").(string))
		if err != nil {
			return diag.Errorf("error retrieving the SAML identity provider ID by name: %s", err)
		}
	}
	if providerId == "" {
		return diag.Errorf("unable to find the provider ID from the API response")
	}

	d.SetId(providerId)

	// After creation, the provider may not be immediately queryable.
	_, err = PollRequest(ctx, func() (interface{}, error) {
		return GetV5SamlProviderById(client, providerId)
	})
	if err != nil {
		return diag.Errorf("error waiting for the SAML identity provider (%s) to be created: %s", providerId, err)
	}

	return resourceV5SamlProviderRead(ctx, d, meta)
}

func buildCreateV5SamlProviderBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                   d.Get("name").(string),
		"saml_metadata_document": d.Get("saml_metadata_document").(string),
	}
	if description := d.Get("description").(string); description != "" {
		bodyParams["description"] = description
	}
	return bodyParams
}

func getV5SamlProviderIdFromResponse(body interface{}) string {
	paths := []string{"provider_id", "saml_provider.provider_id"}
	for _, path := range paths {
		if id, ok := utils.PathSearch(path, body, "").(string); ok && id != "" {
			return id
		}
	}
	return ""
}

func GetV5SamlProviderById(client *golangsdk.ServiceClient, providerId string) (interface{}, error) {
	getHttpUrl := "v5/saml-providers/{provider_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{provider_id}", providerId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	// Compatible with both the wrapped ({"saml_provider": {...}}) and unwrapped responses.
	if v, ok := utils.PathSearch("provider_id", respBody, "").(string); ok && v != "" {
		return respBody, nil
	}
	return utils.PathSearch("saml_provider", respBody, respBody), nil
}

func listV5SamlProviders(client *golangsdk.ServiceClient) ([]interface{}, error) {
	result := make([]interface{}, 0)
	marker := ""
	for {
		listHttpUrl := "v5/saml-providers?limit=200"
		if marker != "" {
			listHttpUrl += "&marker=" + marker
		}
		listPath := client.Endpoint + listHttpUrl
		listOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		resp, err := client.Request("GET", listPath, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		providers := utils.PathSearch("saml_providers", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, providers...)

		nextMarker := utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if nextMarker == "" {
			break
		}
		marker = nextMarker
	}
	return result, nil
}

func getV5SamlProviderIdByName(client *golangsdk.ServiceClient, name string) (string, error) {
	providers, err := listV5SamlProviders(client)
	if err != nil {
		return "", err
	}

	matched := utils.PathSearch(fmt.Sprintf("[?name=='%s']", name), providers, make([]interface{}, 0)).([]interface{})
	if len(matched) == 0 {
		return "", nil
	}
	return utils.PathSearch("provider_id", matched[0], "").(string), nil
}

func resourceV5SamlProviderRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		providerId = d.Id()
	)
	client, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	provider, err := GetV5SamlProviderById(client, providerId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error getting SAML identity provider (%s)", providerId))
	}

	mErr := multierror.Append(
		d.Set("name", utils.PathSearch("name", provider, nil)),
		d.Set("description", utils.PathSearch("description", provider, nil)),
		d.Set("saml_metadata_document", utils.PathSearch("saml_metadata_document", provider, nil)),
		d.Set("provider_id", utils.PathSearch("provider_id", provider, nil)),
		d.Set("urn", utils.PathSearch("urn", provider, nil)),
		d.Set("create_time", utils.PathSearch("create_time", provider, nil)),
		d.Set("expire_time", utils.PathSearch("expire_time", provider, nil)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting fields for the SAML identity provider (%s) resource: %s", providerId, err)
	}

	return nil
}

func resourceV5SamlProviderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		providerId = d.Id()
	)
	client, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	if d.HasChanges("description", "saml_metadata_document") {
		updateHttpUrl := "v5/saml-providers/{provider_id}"
		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{provider_id}", providerId)
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes:          []int{200, 201, 202, 204},
			JSONBody: map[string]interface{}{
				"saml_metadata_document": d.Get("saml_metadata_document").(string),
				"description":            d.Get("description").(string),
			},
		}
		_, err := client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating the SAML identity provider (%s): %s", providerId, err)
		}
	}

	return resourceV5SamlProviderRead(ctx, d, meta)
}

func resourceV5SamlProviderDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		providerId = d.Id()
	)
	client, err := cfg.IAMNoVersionClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	deleteHttpUrl := "v5/saml-providers/{provider_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{provider_id}", providerId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200, 201, 202, 204},
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting SAML identity provider (%s)", providerId))
	}

	return nil
}
