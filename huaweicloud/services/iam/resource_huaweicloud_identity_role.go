package iam

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/policies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM POST /v3.0/OS-ROLE/roles
// @API IAM GET /v3.0/OS-ROLE/roles/{role_id}
// @API IAM PATCH /v3.0/OS-ROLE/roles/{role_id}
// @API IAM DELETE /v3.0/OS-ROLE/roles/{role_id}
func ResourceV3Role() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3IdentityRoleCreate,
		ReadContext:   resourceV3IdentityRoleRead,
		UpdateContext: resourceV3IdentityRoleUpdate,
		DeleteContext: resourceV3IdentityRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the custom policy.`,
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The description of the custom policy.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The display mode of the custom policy.`,
			},
			"policy": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  `The content of the custom policy, in JSON format.`,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := utils.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},

			// Attribute
			"references": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of references.`,
			},
		},
	}
}

func resourceV3IdentityRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	policy := policies.Policy{}
	policyDoc := d.Get("policy").(string)
	err = json.Unmarshal([]byte(policyDoc), &policy)
	if err != nil {
		return diag.Errorf("error unmarshalling policy, please check the format of the policy document: %s", err)
	}

	createOpts := policies.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        d.Get("type").(string),
		Policy:      policy,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	role, err := policies.Create(client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating IAM custom policy: %s", err)
	}

	d.SetId(role.ID)
	return resourceV3IdentityRoleRead(ctx, d, meta)
}

func resourceV3IdentityRoleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	role, err := policies.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "IAM custom policy")
	}

	log.Printf("[DEBUG] Retrieved IAM custom policy: %#v", role)
	policy, err := json.Marshal(role.Policy)
	if err != nil {
		return diag.Errorf("error marshaling policy: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("name", role.Name),
		d.Set("description", role.Description),
		d.Set("type", role.Type),
		d.Set("references", role.References),
		d.Set("policy", string(policy)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting IAM custom policy fields: %s", err)
	}

	return nil
}

func resourceV3IdentityRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	policy := policies.Policy{}
	policyDoc := d.Get("policy").(string)
	err = json.Unmarshal([]byte(policyDoc), &policy)
	if err != nil {
		return diag.Errorf("error unmarshalling policy, please check the format of the policy document: %s", err)
	}

	createOpts := policies.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        d.Get("type").(string),
		Policy:      policy,
	}

	log.Printf("[DEBUG] Update Options: %#v", createOpts)
	_, err = policies.Update(client, d.Id(), createOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating IAM custom policy: %s", err)
	}

	return resourceV3IdentityRoleRead(ctx, d, meta)
}

func resourceV3IdentityRoleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	err = policies.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting IAM custom policy: %s", err)
	}

	return nil
}
