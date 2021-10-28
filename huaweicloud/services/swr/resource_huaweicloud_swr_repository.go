package swr

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/chnsz/golangsdk/openstack/swr/v2/repositories"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func ResourceSWRRepository() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSWRRepositoryCreate,
		ReadContext:   resourceSWRRepositoryRead,
		UpdateContext: resourceSWRRepositoryUpdate,
		DeleteContext: resourceSWRRepositoryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSWRRepositoryImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		//request and response parameters
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 128),
					validation.StringMatch(
						regexp.MustCompile(`^[a-z0-9][a-z0-9._-]+[a-z0-9]+$`),
						"Only lowercase letters, digits, periods (.), underscores (_), and hyphens (-) are allowed.",
					),
					validation.StringDoesNotMatch(
						regexp.MustCompile(`_{3,}?|\.{2,}?|-{2,}?`),
						"Periods, underscores, and hyphens cannot be placed next to each other. A maximum of two consecutive underscores are allowed.",
					),
				),
			},
			"is_public": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"app_server", "linux", "framework_app", "database", "lang", "other", "windows", "arm"},
					false,
				),
			},
			"repository_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internal_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"num_images": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceSWRRepositoryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Unable to create HuaweiCloud SWR client : %s", err)
	}

	name := d.Get("name").(string)
	opts := repositories.CreateOpts{
		Repository:  name,
		Category:    d.Get("category").(string),
		Description: d.Get("description").(string),
		IsPublic:    d.Get("is_public").(bool),
	}

	organization := d.Get("organization").(string)

	err = repositories.Create(client, organization, opts).ExtractErr()
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud SWR Repository: %s", err)
	}
	d.SetId(name)

	return resourceSWRRepositoryRead(ctx, d, meta)
}

func resourceSWRRepositoryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Unable to create HuaweiCloud SWR client : %s", err)
	}

	organization := d.Get("organization").(string)

	repo, err := repositories.Get(client, organization, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error retrieving HuaweiCloud SWR Repository")
	}

	mErr := multierror.Append(
		d.Set("region", config.GetRegion(d)),
		d.Set("name", repo.Name),
		d.Set("repository_id", repo.ID),
		d.Set("description", repo.Description),
		d.Set("category", repo.Category),
		d.Set("is_public", repo.IsPublic),
		d.Set("path", repo.Path),
		d.Set("internal_path", repo.InternalPath),
		d.Set("num_images", repo.NumImages),
		d.Set("size", repo.Size),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error setting HuaweiCloud SWR Repository fields: %s", err)
	}

	return nil
}

func resourceSWRRepositoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Unable to create HuaweiCloud SWR client : %s", err)
	}

	opts := repositories.UpdateOpts{
		Category:    d.Get("category").(string),
		Description: d.Get("description").(string),
		IsPublic:    d.Get("is_public").(bool),
	}

	organization := d.Get("organization").(string)

	err = repositories.Update(client, organization, d.Id(), opts).ExtractErr()
	if err != nil {
		return fmtp.DiagErrorf("Error updating HuaweiCloud SWR Repository: %s", err)
	}

	return resourceSWRRepositoryRead(ctx, d, meta)
}

func resourceSWRRepositoryDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Unable to create HuaweiCloud SWR client : %s", err)
	}

	organization := d.Get("organization").(string)
	err = repositories.Delete(client, organization, d.Id()).ExtractErr()
	if err != nil {
		fmtp.DiagErrorf("error deleting HuaweiCloud SWR Repository: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceSWRRepositoryImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("invalid format specified for SWR repository import: format must be <organization>/<repository>")
		return nil, err
	}
	org := parts[0]
	repo := parts[1]
	d.SetId(repo)
	if err := d.Set("organization", org); err != nil {
		return nil, err
	}
	return schema.ImportStatePassthroughContext(ctx, d, meta)
}
