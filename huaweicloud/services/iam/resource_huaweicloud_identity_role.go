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
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
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
	role, err := policies.Create(identityClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating IAM custom policy: %s", err)
	}

	d.SetId(role.ID)
	return resourceIdentityRoleRead(ctx, d, meta)
}

func resourceIdentityRoleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	role, err := policies.Get(identityClient, d.Id()).Extract()
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

func resourceIdentityRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
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
	_, err = policies.Update(identityClient, d.Id(), createOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating IAM custom policy: %s", err)
	}

	return resourceIdentityRoleRead(ctx, d, meta)
}

func resourceIdentityRoleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	err = policies.Delete(identityClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting IAM custom policy: %s", err)
	}

	return nil
}
