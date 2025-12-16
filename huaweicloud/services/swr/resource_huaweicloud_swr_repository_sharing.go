package swr

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/swr/v2/domains"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

var (
	permissions = map[string]string{
		"pull": "read",
		"push": "write",
	}

	permissionsReverse = map[string]string{
		"read":  "pull",
		"write": "push",
	}
)

// @API SWR GET /v2/manage/namespaces/{namespace}/repositories/{repository}/access-domains/{domain}
// @API SWR PATCH /v2/manage/namespaces/{namespace}/repositories/{repository}/access-domains/{domain}
// @API SWR DELETE /v2/manage/namespaces/{namespace}/repositories/{repository}/access-domains/{domain}
// @API SWR POST /v2/manage/namespaces/{namespace}/repositories/{repository}/access-domains
func ResourceSWRRepositorySharing() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSWRRepositorySharingCreate,
		ReadContext:   resourceSWRRepositorySharingRead,
		UpdateContext: resourceSWRRepositorySharingUpdate,
		DeleteContext: resourceSWRRepositorySharingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSWRRepositorySharingImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"organization": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sharing_account": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"deadline": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\d{4}-\d{2}-\d{2}|forever`),
					"The deadline should be forever or a date in format of YYYY-MM-DD"),
			},
			"permission": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "pull",
				ValidateFunc: validation.StringInSlice([]string{
					"pull",
				}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceSWRRepositorySharingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.SwrV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("unable to create SWR client: %s", err)
	}

	permission := d.Get("permission").(string)

	permit, ok := permissions[permission]
	if !ok {
		return diag.Errorf("the permission type (%s) is not available", permission)
	}

	domain := d.Get("sharing_account").(string)
	opts := domains.CreateOpts{
		AccessDomain: domain,
		Permit:       permit,
		Deadline:     deadlineTrans(d.Get("deadline").(string)),
		Description:  d.Get("description").(string),
	}

	organization := d.Get("organization").(string)
	repository := strings.ReplaceAll(d.Get("repository").(string), "/", "$")

	err = domains.Create(client, organization, repository, opts).ExtractErr()
	if err != nil {
		return diag.Errorf("error creating SWR repository sharing: %s", err)
	}
	d.SetId(domain)

	return resourceSWRRepositorySharingRead(ctx, d, meta)
}

func resourceSWRRepositorySharingRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.SwrV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("unable to create SWR client: %s", err)
	}

	organization := d.Get("organization").(string)
	repository := strings.ReplaceAll(d.Get("repository").(string), "/", "$")

	domain, err := domains.Get(client, organization, repository, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "errorCode", "SVCSTG.SWR.4001090"),
			"error retrieving SWR repository sharing")
	}

	mErr := multierror.Append(
		d.Set("region", cfg.GetRegion(d)),
		d.Set("sharing_account", domain.AccessDomain),
		d.Set("repository", domain.Repository),
		d.Set("organization", domain.Organization),
		d.Set("description", domain.Description),
		d.Set("deadline", deadlineTransReverse(domain.Deadline)),
		d.Set("status", domain.Status),
	)

	if permission, ok := permissionsReverse[domain.Permit]; ok {
		mErr = multierror.Append(mErr, d.Set("permission", permission))
	}
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting SWR repository sharing fields: %v", err)
	}

	return nil
}

func resourceSWRRepositorySharingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.SwrV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("unable to create SWR client: %s", err)
	}

	description := d.Get("description").(string)
	permission := d.Get("permission").(string)

	permit, ok := permissions[permission]
	if !ok {
		return diag.Errorf("the permission type (%s) is not available", permission)
	}

	opts := domains.UpdateOpts{
		Permit:      permit,
		Deadline:    deadlineTrans(d.Get("deadline").(string)),
		Description: &description,
	}

	organization := d.Get("organization").(string)
	repository := strings.ReplaceAll(d.Get("repository").(string), "/", "$")

	err = domains.Update(client, organization, repository, d.Id(), opts).ExtractErr()
	if err != nil {
		return diag.Errorf("error updating SWR repository sharing: %v", err)
	}

	return resourceSWRRepositorySharingRead(ctx, d, meta)
}

func resourceSWRRepositorySharingDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.SwrV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("unable to create SWR client: %s", err)
	}

	organization := d.Get("organization").(string)
	repository := strings.ReplaceAll(d.Get("repository").(string), "/", "$")

	err = domains.Delete(client, organization, repository, d.Id()).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SWR repository sharing")
	}

	return nil
}

func resourceSWRRepositorySharingImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ",")
	if len(parts) != 3 {
		parts = strings.Split(d.Id(), "/")
		if len(parts) != 3 {
			return nil, errors.New("invalid id format, must be <organization>/<repository>/<sharing_account> or " +
				"<organization>,<repository>,<sharing_account>")
		}
	}
	org := parts[0]
	repo := parts[1]
	domain := parts[2]
	d.SetId(domain)
	mErr := multierror.Append(
		d.Set("organization", org),
		d.Set("repository", repo),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, err
	}

	return schema.ImportStatePassthroughContext(ctx, d, meta)
}

func deadlineTrans(deadline string) string {
	if deadline == "forever" {
		return deadline
	}

	return fmt.Sprintf("%sT00:00:00Z", deadline)
}

func deadlineTransReverse(deadline string) string {
	if deadline == "forever" {
		return deadline
	}

	return deadline[:10]
}
