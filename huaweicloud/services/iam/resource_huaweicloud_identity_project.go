package iam

import (
	"context"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3/projects"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceIdentityProjectV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityProjectV3Create,
		ReadContext:   resourceIdentityProjectV3Read,
		UpdateContext: resourceIdentityProjectV3Update,
		DeleteContext: resourceIdentityProjectV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceIdentityProjectV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud identity client: %s", err)
	}

	createOpts := projects.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	project, err := projects.Create(identityClient, createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud project: %s", err)
	}

	d.SetId(project.ID)

	return resourceIdentityProjectV3Read(ctx, d, meta)
}

func resourceIdentityProjectV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud identity client: %s", err)
	}

	project, err := projects.Get(identityClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "IAM project")
	}

	logp.Printf("[DEBUG] Retrieved Huaweicloud project: %#v", project)

	mErr := multierror.Append(nil,
		d.Set("name", project.Name),
		d.Set("description", project.Description),
		d.Set("parent_id", project.ParentID),
		d.Set("enabled", project.Enabled),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error setting identity project fields: %s", err)
	}

	return nil
}

func resourceIdentityProjectV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud identity client: %s", err)
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
			return fmtp.DiagErrorf("Error updating Huaweicloud project: %s", err)
		}
	}

	return resourceIdentityProjectV3Read(ctx, d, meta)
}

func resourceIdentityProjectV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating Huaweicloud identity client: %s", err)
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

		return fmtp.DiagErrorf("Error deleting IAM project: %s", err)
	}

	return nil
}
