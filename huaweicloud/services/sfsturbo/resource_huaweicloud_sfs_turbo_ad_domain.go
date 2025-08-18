package sfsturbo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var adDomainNonUpdatableParams = []string{"share_id"}

// @API SFSTurbo POST /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/active-directory-domain
// @API SFSTurbo GET /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/active-directory-domain
// @API SFSTurbo PUT /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/active-directory-domain
// @API SFSTurbo DELETE /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/active-directory-domain
// @API SFSTurbo GET /v1/{project_id}/sfs-turbo/jobs/{job_id}
func ResourceSFSTurboAdDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSFSTurboAdDomainCreate,
		ReadContext:   resourceSFSTurboAdDomainRead,
		UpdateContext: resourceSFSTurboAdDomainUpdate,
		DeleteContext: resourceSFSTurboAdDomainDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(adDomainNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"share_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_account": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"system_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dns_server": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"overwrite_same_account": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"organization_unit": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateOrUpdateSFSTurboAdDomainBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"service_account":        d.Get("service_account"),
		"password":               d.Get("password"),
		"domain_name":            d.Get("domain_name"),
		"system_name":            d.Get("system_name"),
		"dns_server":             d.Get("dns_server"),
		"overwrite_same_account": d.Get("overwrite_same_account"),
		"organization_unit":      utils.ValueIgnoreEmpty(d.Get("organization_unit")),
		"vpc_id":                 utils.ValueIgnoreEmpty(d.Get("vpc_id")),
	}

	return bodyParams
}

func buildDeleteSFSTurboAdDomainBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"service_account": d.Get("service_account"),
		"password":        d.Get("password"),
	}

	return bodyParams
}

func resourceSFSTurboAdDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		shareId = d.Get("share_id").(string)
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/active-directory-domain"
		product = "sfs-turbo"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{share_id}", shareId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateSFSTurboAdDomainBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo AD domain: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	// The jobId field is inconsistent with the API doc cause there is a problem with the API doc.
	jobID := utils.PathSearch("jobId", respBody, "").(string)
	if jobID == "" {
		return diag.Errorf("error creating SFS Turbo AD domain: jobId is not found in API response")
	}

	if err := waitForSFSTurboAdDomainJobSuccess(ctx, client, jobID, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for SFS Turbo AD domain job creation to succeed: %s", err)
	}

	d.SetId(shareId)

	return resourceSFSTurboAdDomainRead(ctx, d, meta)
}

func resourceSFSTurboAdDomainRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		shareId = d.Id()
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/active-directory-domain"
		product = "sfs-turbo"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{share_id}", shareId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		// When the resource does not exist, the API returns 403 and errCode is SFS.TURBO.9000, which needs to be
		// handled as 404 to be handled by CheckDeletedDiag.
		convertedErr := common.ConvertExpected403ErrInto404Err(err, "errCode", "SFS.TURBO.9000")
		return common.CheckDeletedDiag(d, convertedErr, "error retrieving SFS Turbo AD domain")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("domain_name", utils.PathSearch("domain_name", respBody, nil)),
		d.Set("system_name", utils.PathSearch("system_name", respBody, nil)),
		d.Set("dns_server", utils.PathSearch("dns_server", respBody, nil)),
		d.Set("organization_unit", utils.PathSearch("organization_unit", respBody, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSFSTurboAdDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		shareId = d.Id()
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/active-directory-domain"
		product = "sfs-turbo"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{share_id}", shareId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateSFSTurboAdDomainBodyParams(d)),
	}

	resp, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating SFS Turbo AD domain: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	// The jobId field is inconsistent with the API doc cause there is a problem with the API doc.
	jobID := utils.PathSearch("jobId", respBody, "").(string)
	if jobID == "" {
		return diag.Errorf("error updating SFS Turbo AD domain: jobId is not found in API response")
	}

	if err := waitForSFSTurboAdDomainJobSuccess(ctx, client, jobID, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return diag.Errorf("error waiting for SFS Turbo AD domain job update to succeed: %s", err)
	}

	return resourceSFSTurboAdDomainRead(ctx, d, meta)
}

func resourceSFSTurboAdDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		shareId = d.Id()
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/active-directory-domain"
		product = "sfs-turbo"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{share_id}", shareId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDeleteSFSTurboAdDomainBodyParams(d),
	}

	resp, err := client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		// When the resource does not exist, delete the resource and the API returns 404.
		return common.CheckDeletedDiag(d, err, "error deleting SFS Turbo AD domain")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	// The jobId field is inconsistent with the API doc cause there is a problem with the API doc.
	jobID := utils.PathSearch("jobId", respBody, "").(string)
	if jobID == "" {
		return diag.Errorf("error deleting SFS Turbo AD domain: jobId is not found in API response")
	}

	if err := waitForSFSTurboAdDomainJobSuccess(ctx, client, jobID, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for SFS Turbo AD domain job deletion to succeed: %s", err)
	}

	return nil
}

func waitForSFSTurboAdDomainJobSuccess(ctx context.Context, client *golangsdk.ServiceClient, jobID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"SUCCESS"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := getSFSTurboJobDetail(client, jobID)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("status", respBody, "").(string)
			if status == "" {
				return respBody, "ERROR", fmt.Errorf("status is not found in SFS Turbo job (%s) detail API response", jobID)
			}

			if status == "success" {
				return respBody, "SUCCESS", nil
			}

			if status == "failed" {
				return respBody, status, fmt.Errorf("the SFS Turbo job (%s) status is FAIL, the fail reason is: %s",
					jobID, utils.PathSearch("fail_reason", respBody, "").(string))
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func getSFSTurboJobDetail(client *golangsdk.ServiceClient, jobID string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/sfs-turbo/jobs/{job_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error querying SFS Turbo job detail: %s", err)
	}

	return utils.FlattenResponse(resp)
}
