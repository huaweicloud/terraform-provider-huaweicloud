package iam

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3/projects"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IAM POST /v3/projects
// @API IAM PATCH /v3/projects/{project_id}
// @API IAM PUT /v3-ext/projects/{project_id}
// @API IAM GET /v3-ext/projects/{project_id}
// @API IAM DELETE /v3/projects/{project_id}
func ResourceIdentityProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityProjectCreate,
		ReadContext:   resourceIdentityProjectRead,
		UpdateContext: resourceIdentityProjectUpdate,
		DeleteContext: resourceIdentityProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceIdentityProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	identityExtClient, err := cfg.IdentityV3ExtClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM extension client: %s", err)
	}

	createOpts := projects.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	project, err := projects.Create(identityClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating IAM project: %s", err)
	}

	d.SetId(project.ID)

	if d.Get("status").(string) == "suspended" {
		updateProjectStatusOpts := projects.UpdateStatusOpts{
			Status: d.Get("status").(string),
		}
		err = projects.UpdateStatus(identityExtClient, d.Id(), updateProjectStatusOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error updating IAM project status: %s", err)
		}
	}

	return resourceIdentityProjectRead(ctx, d, meta)
}

func resourceIdentityProjectRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityExtClient, err := cfg.IdentityV3ExtClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM extension client: %s", err)
	}

	project, err := projects.Get(identityExtClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "IAM project")
	}

	log.Printf("[DEBUG] Retrieved IAM project: %#v", project)

	mErr := multierror.Append(nil,
		d.Set("name", project.Name),
		d.Set("description", project.Description),
		d.Set("parent_id", project.ParentID),
		d.Set("enabled", project.Enabled),
		d.Set("status", project.Status),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting IAM project fields: %s", err)
	}

	return nil
}

func resourceIdentityProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	identityExtClient, err := cfg.IdentityV3ExtClient(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM extension client: %s", err)
	}

	var hasChange bool
	var updateOpts projects.UpdateOpts

	if d.HasChange("name") {
		hasChange = true
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateOpts.Description = description
	}

	if hasChange {
		_, err := projects.Update(identityClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating IAM project: %s", err)
		}
	}

	if d.HasChange("status") {
		updateProjectStatusOpts := projects.UpdateStatusOpts{
			Status: d.Get("status").(string),
		}
		err = projects.UpdateStatus(identityExtClient, d.Id(), updateProjectStatusOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error updating IAM project status: %s", err)
		}
	}

	return resourceIdentityProjectRead(ctx, d, meta)
}

func resourceIdentityProjectDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	err = projects.Delete(identityClient, d.Id()).ExtractErr()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			errorMsg := "Deleting projects is not supported. The project is only removed from the state, but it remains in the cloud."
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  errorMsg,
				},
			}
		}

		return diag.Errorf("error deleting IAM project: %s", err)
	}

	return nil
}
