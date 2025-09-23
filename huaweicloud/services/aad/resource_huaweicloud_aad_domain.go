package aad

import (
	"context"
	"fmt"
	"log"
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

var domainNonUpdatableParams = []string{"domain_name", "enterprise_project_id", "vips", "instance_ids", "port_http", "port_https"}

// @API AAD POST /v1/{project_id}/aad/external/domains
// @API AAD GET /v1/aad/protected-domains
// @API AAD GET /v1/aad/protected-domains/{domain_id}
// @API AAD PUT /v1/aad/protected-domains/{domain_id}
// @API AAD DELETE /v2/aad/domains
func ResourceDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainRead,
		UpdateContext: resourceDomainUpdate,
		DeleteContext: resourceDomainDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(domainNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the domain name to be protected by AAD instance.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the enterprise project ID.",
			},
			"real_server_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Specifies the origin server type.",
			},
			"real_server": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the value of the origin server.",
			},
			"vips": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifies the list of AAD instance IP addresses.",
			},
			"instance_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifies the list of AAD instance IDs.",
			},
			"port_http": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Specifies the port when forwarding protocol is HTTP.",
			},
			"port_https": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Specifies the port when forwarding protocol is HTTPS.",
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"cname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cname of domain.",
			},
			"protocol": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The protocol of the domain.",
			},
			"waf_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The protect status of WAF server.",
			},
		},
	}
}

func buildCreateDomainBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"domain_name":           d.Get("domain_name"),
		"enterprise_project_id": d.Get("enterprise_project_id"),
		"real_server_type":      d.Get("real_server_type"),
		"real_server":           d.Get("real_server"),
		"vips":                  utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("vips").([]interface{}))),
		"instance_ids":          utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("instance_ids").([]interface{}))),
		"port_http":             utils.ValueIgnoreEmpty(utils.ExpandToIntList(d.Get("port_http").([]interface{}))),
		"port_https":            utils.ValueIgnoreEmpty(utils.ExpandToIntList(d.Get("port_https").([]interface{}))),
	}
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		product    = "aad"
		httpUrl    = "v1/{project_id}/aad/external/domains"
		domainName = d.Get("domain_name")
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildCreateDomainBodyParams(d)),
	}
	reap, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating the AAD domain(%s): %s", domainName, err)
	}

	respBody, err := utils.FlattenResponse(reap)
	if err != nil {
		return diag.FromErr(err)
	}

	domainId := utils.PathSearch("domainId", respBody, "").(string)
	if domainId == "" {
		return diag.Errorf("error creating AAD domain: domain ID is not found in API response")
	}
	d.SetId(domainId)

	return resourceDomainRead(ctx, d, meta)
}

func setDomainInstanceIDs(d *schema.ResourceData, client *golangsdk.ServiceClient, id string) error {
	var (
		getInstanceIdsHttpUrl = "v1/aad/protected-domains/{domain_id}"
	)

	getInstanceIdsPath := client.Endpoint + getInstanceIdsHttpUrl
	getInstanceIdsPath = strings.ReplaceAll(getInstanceIdsPath, "{domain_id}", id)
	getInstanceIdsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getInstanceIdsResp, err := client.Request("GET", getInstanceIdsPath, &getInstanceIdsOpt)
	if err != nil {
		return nil
	}

	getInstanceIdsRespBody, err := utils.FlattenResponse(getInstanceIdsResp)
	if err != nil {
		log.Printf("[WARN] failed to flatten AAD domain `instance_ids` field: %s", err)
		return nil
	}

	instanceIds := utils.PathSearch("instance_ids", getInstanceIdsRespBody, make([]interface{}, 0)).([]interface{})
	return d.Set("instance_ids", instanceIds)
}

func resourceDomainRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg               = meta.(*config.Config)
		region            = cfg.GetRegion(d)
		product           = "aad"
		getDomainsHttpUrl = "v1/aad/protected-domains"
		mErr              *multierror.Error
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	getDomainsPath := client.Endpoint + getDomainsHttpUrl

	getDomainsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getDomainsResp, err := client.Request("GET", getDomainsPath, &getDomainsOpt)

	if err != nil {
		return diag.Errorf("error retrieving AAD domains: %s", err)
	}

	getDomainsRespBody, err := utils.FlattenResponse(getDomainsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("items[?domain_id=='%s']|[0]", d.Id())
	domain := utils.PathSearch(jsonPath, getDomainsRespBody, nil)
	if domain == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("domain_name", utils.PathSearch("domain_name", domain, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", domain, nil)),
		d.Set("real_server_type", utils.PathSearch("real_server_type", domain, nil)),
		d.Set("real_server", utils.PathSearch("real_servers", domain, nil)),
		d.Set("cname", utils.PathSearch("cname", domain, nil)),
		d.Set("protocol", utils.PathSearch("protocol", domain, nil)),
		d.Set("waf_status", utils.PathSearch("waf_status", domain, 0)),
		setDomainInstanceIDs(d, client, d.Id()),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateDomainsBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"real_server_type": d.Get("real_server_type"),
		"real_servers":     d.Get("real_server"),
	}
}

func resourceDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + "v1/aad/protected-domains/{domain_id}"
	requestPath = strings.ReplaceAll(requestPath, "{domain_id}", d.Id())

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildUpdateDomainsBodyParams(d),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating AAD domain: %s", err)
	}

	return resourceDomainRead(ctx, d, meta)
}

func resourceDomainDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/aad/domains"
		product = "aad"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"domain_id": []string{d.Id()},
		},
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting AAD domain: %s", err)
	}

	return nil
}
