package iam

import (
	"context"
	"encoding/json"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/policies"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceIdentityRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityRoleCreate,
		ReadContext:   resourceIdentityRoleRead,
		UpdateContext: resourceIdentityRoleUpdate,
		DeleteContext: resourceIdentityRoleDelete,
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
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := utils.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"references": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceIdentityRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IAMV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud identity client: %s", err)
	}

	policy := policies.Policy{}
	policyDoc := d.Get("policy").(string)
	err = json.Unmarshal([]byte(policyDoc), &policy)
	if err != nil {
		return fmtp.DiagErrorf("Error unmarshalling policy, please check the format of the policy document: %s", err)
	}
	createOpts := policies.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        d.Get("type").(string),
		Policy:      policy,
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)

	role, err := policies.Create(identityClient, createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud Role: %s", err)
	}

	d.SetId(role.ID)

	return resourceIdentityRoleRead(ctx, d, meta)
}

func resourceIdentityRoleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IAMV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud identity client: %s", err)
	}

	role, err := policies.Get(identityClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "role")
	}

	logp.Printf("[DEBUG] Retrieved HuaweiCloud Role: %#v", role)

	policy, err := json.Marshal(role.Policy)
	if err != nil {
		return fmtp.DiagErrorf("Error marshalling policy: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("name", role.Name),
		d.Set("description", role.Description),
		d.Set("type", role.Type),
		d.Set("references", role.References),
		d.Set("policy", string(policy)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error setting identity role fields: %s", err)
	}

	return nil
}

func resourceIdentityRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IAMV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud identity client: %s", err)
	}

	policy := policies.Policy{}
	policyDoc := d.Get("policy").(string)
	err = json.Unmarshal([]byte(policyDoc), &policy)
	if err != nil {
		return fmtp.DiagErrorf("Error unmarshalling policy, please check the format of the policy document: %s", err)
	}
	createOpts := policies.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        d.Get("type").(string),
		Policy:      policy,
	}

	logp.Printf("[DEBUG] Update Options: %#v", createOpts)

	_, err = policies.Update(identityClient, d.Id(), createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error updating HuaweiCloud Role: %s", err)
	}

	return resourceIdentityRoleRead(ctx, d, meta)
}

func resourceIdentityRoleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IAMV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud identity client: %s", err)
	}

	err = policies.Delete(identityClient, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.DiagErrorf("Error deleting HuaweiCloud Role: %s", err)
	}

	return nil
}
