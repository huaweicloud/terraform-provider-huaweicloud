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

// @API IAM POST /v5/oidc-providers
// @API IAM GET /v5/oidc-providers
// @API IAM GET /v5/oidc-providers/{provider_id}
// @API IAM PUT /v5/oidc-providers/{provider_id}
// @API IAM DELETE /v5/oidc-providers/{provider_id}
// @API IAM POST /v5/oidc-providers/{provider_id}/add-client-id
// @API IAM POST /v5/oidc-providers/{provider_id}/remove-client-id
// @API IAM PUT /v5/oidc-providers/{provider_id}/thumbprint
func ResourceV5OidcProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV5OidcProviderCreate,
		ReadContext:   resourceV5OidcProviderRead,
		UpdateContext: resourceV5OidcProviderUpdate,
		DeleteContext: resourceV5OidcProviderDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the OIDC identity provider. The name contains 1 to 64 characters and can only contain letters, digits, underscores (_) and hyphens (-).`,
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the URL of the OIDC identity provider. The value corresponds to the issuer field in the OpenID Connect metadata document of the identity provider, and contains 1 to 255 characters.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the OIDC identity provider. The description contains a maximum of 255 characters and cannot contain the following characters: @, #, %, &, <, >, ", $, ^ and *.`,
			},
			"client_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of client IDs registered in the OIDC identity provider. A maximum of 100 client IDs can be configured, and each client ID contains 1 to 255 characters.`,
			},
			"thumbprints": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of server certificate fingerprints of the OIDC identity provider. If this parameter is not specified, IAM automatically obtains the CA certificate fingerprint of the JWKS endpoint.`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uniform resource name of the OIDC identity provider.`,
			},
		},
	}
}

func resourceV5OidcProviderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	createHttpUrl := "v5/oidc-providers"
	createPath := client.Endpoint + createHttpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateV5OidcProviderBodyParams(d),
	}
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating OIDC identity provider: %s", err)
	}

	createBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	providerId := getV5OidcProviderIdFromResponse(createBody)
	if providerId == "" {
		// If the create response does not return the provider ID, query it by the name.
		providerId, err = getV5OidcProviderIdByName(client, d.Get("name").(string))
		if err != nil {
			return diag.Errorf("error retrieving the OIDC identity provider ID by name: %s", err)
		}
	}
	if providerId == "" {
		return diag.Errorf("unable to find the provider ID from the API response")
	}

	d.SetId(providerId)

	// After creation, the provider may not be immediately queryable.
	_, err = PollRequest(ctx, func() (interface{}, error) {
		return GetV5OidcProviderById(client, providerId)
	})
	if err != nil {
		return diag.Errorf("error waiting for the OIDC identity provider (%s) to be created: %s", providerId, err)
	}

	return resourceV5OidcProviderRead(ctx, d, meta)
}

func buildCreateV5OidcProviderBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":       d.Get("name").(string),
		"url":        d.Get("url").(string),
		"client_ids": d.Get("client_ids").(*schema.Set).List(),
	}
	if description := d.Get("description").(string); description != "" {
		bodyParams["description"] = description
	}
	if thumbprints := d.Get("thumbprints").(*schema.Set).List(); len(thumbprints) > 0 {
		bodyParams["thumbprints"] = thumbprints
	}
	return bodyParams
}

func getV5OidcProviderIdFromResponse(body interface{}) string {
	paths := []string{"provider_id", "oidc_provider.provider_id"}
	for _, path := range paths {
		if id, ok := utils.PathSearch(path, body, "").(string); ok && id != "" {
			return id
		}
	}
	return ""
}

func GetV5OidcProviderById(client *golangsdk.ServiceClient, providerId string) (interface{}, error) {
	getHttpUrl := "v5/oidc-providers/{provider_id}"
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

	// Compatible with both the wrapped ({"oidc_provider": {...}}) and unwrapped responses.
	if v, ok := utils.PathSearch("provider_id", respBody, "").(string); ok && v != "" {
		return respBody, nil
	}
	return utils.PathSearch("oidc_provider", respBody, respBody), nil
}

func listV5OidcProviders(client *golangsdk.ServiceClient) ([]interface{}, error) {
	result := make([]interface{}, 0)
	marker := ""
	for {
		listHttpUrl := "v5/oidc-providers?limit=200"
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

		providers := utils.PathSearch("oidc_providers", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, providers...)

		nextMarker := utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if nextMarker == "" {
			break
		}
		marker = nextMarker
	}
	return result, nil
}

func getV5OidcProviderIdByName(client *golangsdk.ServiceClient, name string) (string, error) {
	providers, err := listV5OidcProviders(client)
	if err != nil {
		return "", err
	}

	matched := utils.PathSearch(fmt.Sprintf("[?name=='%s']", name), providers, make([]interface{}, 0)).([]interface{})
	if len(matched) == 0 {
		return "", nil
	}
	return utils.PathSearch("provider_id", matched[0], "").(string), nil
}

func resourceV5OidcProviderRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		providerId = d.Id()
	)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	provider, err := GetV5OidcProviderById(client, providerId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error getting OIDC identity provider (%s)", providerId))
	}

	mErr := multierror.Append(
		d.Set("name", utils.PathSearch("name", provider, nil)),
		d.Set("url", utils.PathSearch("url", provider, nil)),
		d.Set("description", utils.PathSearch("description", provider, nil)),
		d.Set("urn", utils.PathSearch("urn", provider, nil)),
		d.Set("client_ids", utils.PathSearch("client_ids", provider, make([]interface{}, 0)).([]interface{})),
		d.Set("thumbprints", utils.PathSearch("thumbprints", provider, make([]interface{}, 0)).([]interface{})),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting fields for the OIDC identity provider (%s) resource: %s", providerId, err)
	}

	return nil
}

func resourceV5OidcProviderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		providerId = d.Id()
	)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	if d.HasChange("description") {
		updateHttpUrl := "v5/oidc-providers/{provider_id}"
		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{provider_id}", providerId)
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes:          []int{200, 201, 202, 204},
			JSONBody: map[string]interface{}{
				"description": d.Get("description").(string),
			},
		}
		_, err := client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating the description of the OIDC identity provider (%s): %s", providerId, err)
		}
	}

	if d.HasChange("client_ids") {
		o, n := d.GetChange("client_ids")
		err := updateV5OidcProviderClientIds(client, providerId, o.(*schema.Set), n.(*schema.Set))
		if err != nil {
			return diag.Errorf("error updating the client IDs of the OIDC identity provider (%s): %s", providerId, err)
		}
	}

	if d.HasChange("thumbprints") {
		thumbprintHttpUrl := "v5/oidc-providers/{provider_id}/thumbprint"
		thumbprintPath := client.Endpoint + thumbprintHttpUrl
		thumbprintPath = strings.ReplaceAll(thumbprintPath, "{provider_id}", providerId)
		thumbprintOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"thumbprints": d.Get("thumbprints").(*schema.Set).List(),
			},
		}
		_, err := client.Request("PUT", thumbprintPath, &thumbprintOpt)
		if err != nil {
			return diag.Errorf("error updating the thumbprints of the OIDC identity provider (%s): %s", providerId, err)
		}
	}

	return resourceV5OidcProviderRead(ctx, d, meta)
}

func updateV5OidcProviderClientIds(client *golangsdk.ServiceClient, providerId string, oldSet, newSet *schema.Set) error {
	// Add the newly added client IDs.
	for _, v := range newSet.Difference(oldSet).List() {
		addHttpUrl := "v5/oidc-providers/{provider_id}/add-client-id"
		addPath := client.Endpoint + addHttpUrl
		addPath = strings.ReplaceAll(addPath, "{provider_id}", providerId)
		addOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes:          []int{200, 201, 202, 204},
			JSONBody: map[string]interface{}{
				"client_id": v.(string),
			},
		}
		if _, err := client.Request("POST", addPath, &addOpt); err != nil {
			return fmt.Errorf("error adding client ID (%v): %s", v, err)
		}
	}

	// Remove the deleted client IDs.
	for _, v := range oldSet.Difference(newSet).List() {
		removeHttpUrl := "v5/oidc-providers/{provider_id}/remove-client-id"
		removePath := client.Endpoint + removeHttpUrl
		removePath = strings.ReplaceAll(removePath, "{provider_id}", providerId)
		removeOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes:          []int{200, 201, 202, 204},
			JSONBody: map[string]interface{}{
				"client_id": v.(string),
			},
		}
		if _, err := client.Request("POST", removePath, &removeOpt); err != nil {
			return fmt.Errorf("error removing client ID (%v): %s", v, err)
		}
	}

	return nil
}

func resourceV5OidcProviderDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		providerId = d.Id()
	)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	deleteHttpUrl := "v5/oidc-providers/{provider_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{provider_id}", providerId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200, 201, 202, 204},
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting OIDC identity provider (%s)", providerId))
	}

	return nil
}
