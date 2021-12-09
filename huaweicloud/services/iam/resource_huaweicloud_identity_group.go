package iam

import (
	"context"

	"github.com/chnsz/golangsdk/openstack/identity/v3/groups"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceIdentityGroupV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityGroupV3Create,
		ReadContext:   resourceIdentityGroupV3Read,
		UpdateContext: resourceIdentityGroupV3Update,
		DeleteContext: resourceIdentityGroupV3Delete,
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

func resourceIdentityGroupV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud identity client: %s", err)
	}

	createOpts := groups.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)

	group, err := groups.Create(identityClient, createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud Group: %s", err)
	}

	d.SetId(group.ID)

	return resourceIdentityGroupV3Read(ctx, d, meta)
}

func resourceIdentityGroupV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud identity client: %s", err)
	}

	group, err := groups.Get(identityClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "group")
	}

	logp.Printf("[DEBUG] Retrieved HuaweiCloud Group: %#v", group)

	mErr := multierror.Append(nil,
		d.Set("name", group.Name),
		d.Set("description", group.Description),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error setting identity group key fields: %s", err)
	}

	return nil
}

func resourceIdentityGroupV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud identity client: %s", err)
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
		logp.Printf("[DEBUG] Update Options: %#v", updateOpts)
	}

	if hasChange {
		_, err := groups.Update(identityClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmtp.DiagErrorf("Error updating HuaweiCloud group: %s", err)
		}
	}

	return resourceIdentityGroupV3Read(ctx, d, meta)
}

func resourceIdentityGroupV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud identity client: %s", err)
	}

	err = groups.Delete(identityClient, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.DiagErrorf("Error deleting HuaweiCloud group: %s", err)
	}

	return nil
}
