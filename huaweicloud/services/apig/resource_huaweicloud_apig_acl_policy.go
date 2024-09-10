package apig

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/acls"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// ResourceAclPolicy is a provider resource of the APIG ACL policy.
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/acls/{acl_id}
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/acls/{acl_id}
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/acls/{acl_id}
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/acls
func ResourceAclPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAclPolicyCreate,
		ReadContext:   resourceAclPolicyRead,
		UpdateContext: resourceAclPolicyUpdate,
		DeleteContext: resourceAclPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAclPolicyImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the ACL policy is located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the dedicated instance to which the ACL policy belongs.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the ACL policy.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the ACL policy.",
			},
			"entity_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The entity type of the ACL policy.",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "One or more objects from which the access will be controlled.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the ACL policy.",
			},
		},
	}
}

func resourceAclPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	// Build ACL policy create options according to schema configuration.
	opts := acls.CreateOpts{
		Name:       d.Get("name").(string),
		Type:       d.Get("type").(string),
		EntityType: d.Get("entity_type").(string),
		Value:      d.Get("value").(string),
	}
	resp, err := acls.Create(client, instanceId, opts)
	if err != nil {
		return diag.Errorf("error creating ACL policy: %s", err)
	}
	d.SetId(resp.ID)

	return resourceAclPolicyRead(ctx, d, meta)
}

func resourceAclPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		policyId   = d.Id()
	)
	resp, err := acls.Get(client, instanceId, policyId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "ACL policy")
	}
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("type", resp.Type),
		d.Set("name", resp.Name),
		d.Set("entity_type", resp.EntityType),
		d.Set("value", resp.Value),
		d.Set("updated_at", resp.UpdatedAt),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving ACL policy (%s) fields: %s", policyId, err)
	}
	return nil
}

func resourceAclPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		policyId   = d.Id()
		// Build ACL policy update options according to schema configuration.
		opts = acls.UpdateOpts{
			Name:       d.Get("name").(string),
			Type:       d.Get("type").(string),
			EntityType: d.Get("entity_type").(string),
			Value:      d.Get("value").(string),
		}
	)
	_, err = acls.Update(client, instanceId, policyId, opts)
	if err != nil {
		return diag.Errorf("error updating ACL policy: %s", err)
	}

	return resourceAclPolicyRead(ctx, d, meta)
}

func resourceAclPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		policyId   = d.Id()
	)
	if err = acls.Delete(client, instanceId, policyId); err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("unable to delete the ACL policy (%s)", policyId))
	}

	return nil
}

func resourceAclPolicyImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but '%s'", importedId)
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("instance_id", parts[0])
}
