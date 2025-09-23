package sfsturbo

import (
	"context"
	"errors"
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

// Due to the lack of a test LDAP test environment, this resource has not been fully tested and the code coverage is
// only 35%, so the document is temporarily placed in the incubating directory.

// @API SFSTurbo POST /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/ldap
// @API SFSTurbo GET /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/ldap
// @API SFSTurbo DELETE /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/ldap
// @API SFSTurbo PUT /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/ldap
// @API SFSTurbo GET /v1/{project_id}/sfs-turbo/jobs/{job_id}
func ResourceLdapConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLdapConfigCreate,
		ReadContext:   resourceLdapConfigRead,
		UpdateContext: resourceLdapConfigUpdate,
		DeleteContext: resourceLdapConfigDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew([]string{"share_id"}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Description: `Specifies the region in which to create the resource. If omitted, the provider-level
				region will be used.`,
			},
			"share_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the file system ID.`,
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the URL of the LDAP server.`,
			},
			"base_dn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the base DN.`,
			},
			"user_dn": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the bind DN.`,
			},
			// This field is not returned by the API
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `Specifies the LDAP authentication password.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the ID of the VPC that the specified LDAP server can connect to.`,
			},
			"backup_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the URL of the standby LDAP server.`,
			},
			"schema": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the LDAP schema.`,
			},
			"search_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the LDAP search timeout interval, in seconds. `,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildLdapConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"url":            d.Get("url"),
		"base_dn":        d.Get("base_dn"),
		"user_dn":        utils.ValueIgnoreEmpty(d.Get("user_dn")),
		"password":       utils.ValueIgnoreEmpty(d.Get("password")),
		"vpc_id":         utils.ValueIgnoreEmpty(d.Get("vpc_id")),
		"backup_url":     utils.ValueIgnoreEmpty(d.Get("backup_url")),
		"schema":         utils.ValueIgnoreEmpty(d.Get("schema")),
		"search_timeout": utils.ValueIgnoreEmpty(d.Get("search_timeout")),
	}

	return bodyParams
}

func getTurboJobDetail(client *golangsdk.ServiceClient, jobID string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/sfs-turbo/jobs/{job_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func waitingForTurboJobSuccess(ctx context.Context, client *golangsdk.ServiceClient, jobID string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			jobDetail, err := getTurboJobDetail(client, jobID)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("status", jobDetail, "").(string)
			if status == "" {
				return nil, "ERROR", errors.New("status is not found in job API response")
			}

			if status == "failed" {
				return jobDetail, status, nil
			}

			if status == "success" {
				return jobDetail, "COMPLETED", nil
			}

			return jobDetail, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceLdapConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		shareID = d.Get("share_id").(string)
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/ldap"
	)

	client, err := cfg.NewServiceClient("sfs-turbo", region)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{share_id}", shareID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildLdapConfigBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo LDAP configuration: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobID := utils.PathSearch("jobId", respBody, "").(string)
	if jobID == "" {
		return diag.Errorf("error creating SFS Turbo LDAP configuration: job ID is empty")
	}

	if err := waitingForTurboJobSuccess(ctx, client, jobID, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for SFS Turbo LDAP configuration creation to complete: %s", err)
	}

	d.SetId(shareID)

	return resourceLdapConfigRead(ctx, d, meta)
}

func ReadLdapConfig(client *golangsdk.ServiceClient, shareID string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/ldap"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{share_id}", shareID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceLdapConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		shareID = d.Get("share_id").(string)
	)

	client, err := cfg.NewServiceClient("sfs-turbo", region)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo client: %s", err)
	}

	respBody, err := ReadLdapConfig(client, shareID)
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected403ErrInto404Err(err, "errCode", "SFS.TURBO.9000"),
			"error retrieving SFS Turbo LDAP configuration")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("url", utils.PathSearch("url", respBody, nil)),
		d.Set("base_dn", utils.PathSearch("base_dn", respBody, nil)),
		d.Set("user_dn", utils.PathSearch("user_dn", respBody, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", respBody, nil)),
		d.Set("backup_url", utils.PathSearch("backup_url", respBody, nil)),
		d.Set("schema", utils.PathSearch("schema", respBody, nil)),
		d.Set("search_timeout", utils.PathSearch("search_timeout", respBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceLdapConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		shareID = d.Get("share_id").(string)
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/ldap"
	)

	client, err := cfg.NewServiceClient("sfs-turbo", region)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{share_id}", shareID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildLdapConfigBodyParams(d)),
	}

	resp, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating SFS Turbo LDAP configuration: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobID := utils.PathSearch("jobId", respBody, "").(string)
	if jobID == "" {
		return diag.Errorf("error updating SFS Turbo LDAP configuration: job ID is empty")
	}

	if err := waitingForTurboJobSuccess(ctx, client, jobID, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return diag.Errorf("error waiting for SFS Turbo LDAP configuration update to complete: %s", err)
	}

	return resourceLdapConfigRead(ctx, d, meta)
}

func resourceLdapConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		shareID = d.Get("share_id").(string)
		httpUrl = "v1/{project_id}/sfs-turbo/shares/{share_id}/fs/ldap"
	)

	client, err := cfg.NewServiceClient("sfs-turbo", region)
	if err != nil {
		return diag.Errorf("error creating SFS Turbo client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{share_id}", shareID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting SFS Turbo LDAP configuration: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobID := utils.PathSearch("jobId", respBody, "").(string)
	if jobID == "" {
		return diag.Errorf("error deleting SFS Turbo LDAP configuration: job ID is empty")
	}

	if err := waitingForTurboJobSuccess(ctx, client, jobID, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for SFS Turbo LDAP configuration deletion to complete: %s", err)
	}

	return nil
}
