package cae

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

// @API CAE POST /v1/{project_id}/cae/domains
// @API CAE GET /v1/{project_id}/cae/domains
// @API CAE DELETE /v1/{project_id}/cae/domains/{domain_id}
func ResourceDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainRead,
		DeleteContext: resourceDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDomainImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the CAE environment.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The domain name to be associated with the CAE environment.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The ID of the enterprise project to which the domain name belongs.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the domain name is associated, in RFC3339 format.`,
			},
		},
	}
}

func buildCreateDomainBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"api_version": "v1",
		"kind":        "Domain",
		"metadata": map[string]interface{}{
			"name": d.Get("name"),
		},
	}
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		httpUrl       = "v1/{project_id}/cae/domains"
		environmentId = d.Get("environment_id").(string)
	)

	client, err := cfg.NewServiceClient("cae", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(environmentId, cfg.GetEnterpriseProjectID(d)),
		JSONBody:         buildCreateDomainBodyParams(d),
	}
	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("unable to associate domain name for specified environment (%s): %s", environmentId, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	domainId := utils.PathSearch("items[0].metadata.id", respBody, "").(string)
	if domainId == "" {
		return diag.Errorf("unable to find the domain name ID from the API response")
	}

	d.SetId(domainId)

	return resourceDomainRead(ctx, d, meta)
}

func resourceDomainRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("cae", region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	domainInfo, err := GetDomainById(client, d.Get("environment_id").(string), d.Id(), cfg.GetEnterpriseProjectID(d))
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving domain name (%s)", d.Get("name").(string)))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", domainInfo, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("created_at",
			domainInfo, "").(string))/1000, false)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func GetDomainById(client *golangsdk.ServiceClient, environmentId, domainId, epsId string) (interface{}, error) {
	domains, err := getDomains(client, environmentId, epsId)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", envResourceNotFoundCodes...)
	}

	domainInfo := utils.PathSearch(fmt.Sprintf("items[?metadata.id=='%s']|[0].metadata", domainId), domains, nil)
	if domainInfo == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return domainInfo, nil
}

func resourceDomainDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1/{project_id}/cae/domains/{domain_id}"
	)

	client, err := cfg.NewServiceClient("cae", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{domain_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(d.Get("environment_id").(string), cfg.GetEnterpriseProjectID(d)),
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// CAE.01500005: The environment or resource not found.
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "CAE.01500005"),
			fmt.Sprintf("error deleting associated domain name (%s)", d.Get("name").(string)))
	}
	return nil
}

func getDomains(client *golangsdk.ServiceClient, environmentId, epsId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/cae/domains"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(environmentId, epsId),
	}
	resp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

// Since the ID cannot be found on the console, so we need to import by the domain name.
func resourceDomainImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var (
		cfg        = meta.(*config.Config)
		importedId = d.Id()
		parts      = strings.Split(importedId, "/")
	)

	if len(parts) != 2 && len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<environment_id>/<name>' or "+
			"'<environment_id>/<name>/<enterprise_project_id>', but got '%s'",
			importedId)
	}

	var (
		environmentId = parts[0]
		domainName    = parts[1]
	)

	mErr := multierror.Append(
		d.Set("environment_id", environmentId),
		d.Set("name", domainName),
	)

	if len(parts) == 3 {
		mErr = multierror.Append(mErr, d.Set("enterprise_project_id", parts[2]))
	}

	if mErr.ErrorOrNil() != nil {
		return nil, mErr
	}

	client, err := cfg.NewServiceClient("cae", cfg.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating CAE client: %s", err)
	}

	domains, err := getDomains(client, environmentId, cfg.GetEnterpriseProjectID(d))
	if err != nil {
		return nil, fmt.Errorf("error retrieving domains: %s", err)
	}

	domainId := utils.PathSearch(fmt.Sprintf("items[?metadata.name=='%s']|[0].metadata.id", domainName), domains, "").(string)
	if domainId == "" {
		return nil, fmt.Errorf("unable to find ID of the domain name (%s) from API response : %s", domainName, err)
	}

	d.SetId(domainId)
	return []*schema.ResourceData{d}, nil
}
