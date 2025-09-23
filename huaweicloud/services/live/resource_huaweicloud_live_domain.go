package live

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Live POST /v1/{project_id}/domain
// @API Live GET /v1/{project_id}/domain
// @API Live PUT /v1/{project_id}/domain
// @API Live DELETE /v1/{project_id}/domain
// @API Live PUT /v1/{project_id}/domains_mapping
// @API Live DELETE /v1/{project_id}/domains_mapping
// @API Live PUT /v1/{project_id}/domain/ipv6-switch
func ResourceDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainRead,
		UpdateContext: resourceDomainUpdate,
		DeleteContext: resourceDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_area": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"ingest_domain_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/domain"
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateDomainBodyParams(cfg, d, region)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating Live domain: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	domainName := utils.PathSearch("domain", respBody, "").(string)
	if domainName == "" {
		return diag.Errorf("error creating Live domain: domain name is not found in API response")
	}

	d.SetId(domainName)

	err = waitingForDomainStateCompleted(ctx, client, d.Timeout(schema.TimeoutCreate), domainName, "on")
	if err != nil {
		return diag.Errorf("error waiting for the Live domain (%s) creation to complete: %s", domainName, err)
	}

	// Associate the streaming domain name with an ingest domain.
	domainType := d.Get("type").(string)
	if ingestDomain, ok := d.GetOk("ingest_domain_name"); ok && domainType == "pull" {
		err = associateIngestDomain(client, domainName, ingestDomain.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Disable the domain.
	status := d.Get("status").(string)
	if status == "off" {
		err = updateDomainStatus(ctx, cfg, d, client, d.Timeout(schema.TimeoutCreate), status)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Enable the IPv6.
	if d.Get("is_ipv6").(bool) {
		if err := updateIPv6Switch(client, d, domainName); err != nil {
			return diag.Errorf("error updating Live domain IPv6 switch in creation operation: %s", err)
		}
	}

	return resourceDomainRead(ctx, d, meta)
}

func buildCreateDomainBodyParams(cfg *config.Config, d *schema.ResourceData, region string) map[string]interface{} {
	params := map[string]interface{}{
		"domain":                d.Get("name"),
		"domain_type":           d.Get("type"),
		"region":                region,
		"service_area":          utils.ValueIgnoreEmpty(d.Get("service_area")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
	}

	return params
}

func resourceDomainRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	respBody, err := GetDomain(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Live domain")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("domain", respBody, nil)),
		d.Set("type", utils.PathSearch("domain_type", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("ingest_domain_name", utils.PathSearch("related_domain", respBody, nil)),
		d.Set("cname", utils.PathSearch("domain_cname", respBody, nil)),
		d.Set("service_area", utils.PathSearch("service_area", respBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", respBody, nil)),
		d.Set("is_ipv6", utils.PathSearch("is_ipv6", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetDomain(client *golangsdk.ServiceClient, domainName string) (interface{}, error) {
	httpUrl := "v1/{project_id}/domain"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = fmt.Sprintf("%s?domain=%s", getPath, domainName)
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

	domainInfo := utils.PathSearch("domain_info|[0]", respBody, nil)
	if domainInfo == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return domainInfo, nil
}

func resourceDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		domainType = d.Get("type").(string)
		domainName = d.Id()
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	// Associate the streaming domain name with an ingest domain or delete association.
	if d.HasChange("ingest_domain_name") && domainType == "pull" {
		oldIngetstDomain, newIngetstDomain := d.GetChange("ingest_domain_name")

		if newIngetstDomain == "" {
			err = disassociateIngestDomain(client, domainName, oldIngetstDomain.(string))
		} else {
			err = associateIngestDomain(client, domainName, newIngetstDomain.(string))
		}

		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Update the domain status.
	if d.HasChange("status") {
		status := d.Get("status").(string)
		err = updateDomainStatus(ctx, cfg, d, client, d.Timeout(schema.TimeoutUpdate), status)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Enable or disable IPv6
	if d.HasChange("is_ipv6") {
		if err := updateIPv6Switch(client, d, domainName); err != nil {
			return diag.Errorf("error updating Live domain IPv6 switch in update operation: %s", err)
		}
	}

	return resourceDomainRead(ctx, d, meta)
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/domain"
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	// Disable the domain.
	// Only the status is `off`, the domain can be deleted.
	status := d.Get("status").(string)
	if status != "off" {
		err = updateDomainStatus(ctx, cfg, d, client, d.Timeout(schema.TimeoutDelete), "off")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Call the deletion API, deleting the domain.
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = fmt.Sprintf("%s?domain=%s", deletePath, d.Id())
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "LIVE.103011019"),
			"error deleting Live domain")
	}

	return nil
}

func associateIngestDomain(client *golangsdk.ServiceClient, pullDomain, pushDomain string) error {
	httpUrl := "v1/{project_id}/domains_mapping"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildAssociateDomainBodyParams(pullDomain, pushDomain)),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error associating the ingest domain for the streaming domain: %s", err)
	}

	return nil
}

func buildAssociateDomainBodyParams(pullDomain, pushDomain string) map[string]interface{} {
	params := map[string]interface{}{
		"pull_domain": pullDomain,
		"push_domain": pushDomain,
	}

	return params
}

func disassociateIngestDomain(client *golangsdk.ServiceClient, pullDomain, pushDomain string) error {
	httpUrl := "v1/{project_id}/domains_mapping"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = fmt.Sprintf("%s?pull_domain=%s&push_domain=%s", requestPath, pullDomain, pushDomain)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("DELETE", requestPath, &requestOpts)
	if err != nil {
		return fmt.Errorf("error disassociating the ingest domain from the streaming domain: %s", err)
	}

	return nil
}

func updateDomainStatus(ctx context.Context, cfg *config.Config, d *schema.ResourceData, client *golangsdk.ServiceClient,
	t time.Duration, status string) error {
	httpUrl := "v1/{project_id}/domain"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateaDomainStatusBodyParams(cfg, d, status)),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error updating the Live domain status: %s", err)
	}

	err = waitingForDomainStateCompleted(ctx, client, t, d.Id(), status)
	if err != nil {
		return fmt.Errorf("error waiting for the Live domain status update to complete: %s", err)
	}

	return nil
}

func buildUpdateaDomainStatusBodyParams(cfg *config.Config, d *schema.ResourceData,
	status string) map[string]interface{} {
	params := map[string]interface{}{
		"domain":                d.Get("name"),
		"status":                status,
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
	}

	return params
}

func updateIPv6Switch(client *golangsdk.ServiceClient, d *schema.ResourceData, domainName string) error {
	httpUrl := "v1/{project_id}/domain/ipv6-switch"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateIPv6SwitchBodyParams(d, domainName)),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error enabling or disabling the IPv6: %s", err)
	}

	return nil
}

func buildUpdateIPv6SwitchBodyParams(d *schema.ResourceData, domainName string) map[string]interface{} {
	params := map[string]interface{}{
		"domain":  domainName,
		"is_ipv6": utils.ValueIgnoreEmpty(d.Get("is_ipv6")),
	}

	return params
}

func waitingForDomainStateCompleted(ctx context.Context, client *golangsdk.ServiceClient, t time.Duration, domainName, status string) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      waitDomainStatusRefreshFunc(client, domainName, status),
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitDomainStatusRefreshFunc(client *golangsdk.ServiceClient, domainName, status string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetDomain(client, domainName)
		if err != nil {
			return nil, "ERROR", err
		}

		domainStatus := utils.PathSearch("status", respBody, "").(string)

		if domainStatus == status {
			return respBody, "COMPLETED", nil
		}

		return respBody, "PENDING", nil
	}
}
