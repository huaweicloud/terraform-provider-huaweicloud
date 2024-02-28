package iam

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/identity/v3/groups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IAM POST /v3/groups
// @API IAM GET /v3/groups/{group_id}
// @API IAM PATCH /v3/groups/{group_id}
// @API IAM DELETE /v3/groups/{group_id}
func ResourceIdentityGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityGroupCreate,
		ReadContext:   resourceIdentityGroupRead,
		UpdateContext: resourceIdentityGroupUpdate,
		DeleteContext: resourceIdentityGroupDelete,
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
		},
	}
}

func resourceIdentityGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	createOpts := groups.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	group, err := groups.Create(identityClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating IAM group: %s", err)
	}

	d.SetId(group.ID)
	return resourceIdentityGroupRead(ctx, d, meta)
}

func resourceIdentityGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	group, err := groups.Get(identityClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "IAM group")
	}

	log.Printf("[DEBUG] Retrieved IAM group: %#v", group)
	mErr := multierror.Append(nil,
		d.Set("name", group.Name),
		d.Set("description", group.Description),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting IAM group fields: %s", err)
	}

	return nil
}

func resourceIdentityGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	var hasChange bool
	var updateOpts groups.UpdateOpts

	if d.HasChange("description") {
		hasChange = true
		updateOpts.Description = d.Get("description").(string)
	}

	if d.HasChange("name") {
		hasChange = true
		updateOpts.Name = d.Get("name").(string)
	}

	if hasChange {
		log.Printf("[DEBUG] Update Options: %#v", updateOpts)

		_, err := groups.Update(identityClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating IAM group: %s", err)
		}
	}

	return resourceIdentityGroupRead(ctx, d, meta)
}

func resourceIdentityGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	err = groups.Delete(identityClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting IAM group: %s", err)
	}

	return nil
}
