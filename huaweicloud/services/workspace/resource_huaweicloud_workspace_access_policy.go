package workspace

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/workspace/v2/accesspolicies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API Workspace GET /v2/{project_id}/access-policy/{access_policy_id}/objects
// @API Workspace PUT /v2/{project_id}/access-policy/{access_policy_id}/objects
// @API Workspace DELETE /v2/{project_id}/access-policy
// @API Workspace GET /v2/{project_id}/access-policy
// @API Workspace POST /v2/{project_id}/access-policy
func ResourceAccessPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccessPolicyCreate,
		ReadContext:   resourceAccessPolicyRead,
		UpdateContext: resourceAccessPolicyUpdate,
		DeleteContext: resourceAccessPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAccessPolicyImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the access policy is located.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the access policy.",
			},
			"blacklist_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of access policy blacklist.",
			},
			"blacklist": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The object type.",
						},
						"object_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The object ID.",
						},
						"object_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The object name.",
						},
					},
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the access policy.",
			},
		},
	}
}

func buildAccessPolicyObjects(objects *schema.Set) []accesspolicies.AccessPolicyObjectInfo {
	if objects.Len() < 1 {
		return nil
	}

	result := make([]accesspolicies.AccessPolicyObjectInfo, objects.Len())
	for i, val := range objects.List() {
		object := val.(map[string]interface{})
		result[i] = accesspolicies.AccessPolicyObjectInfo{
			ObjectType: object["object_type"].(string),
			ObjectId:   object["object_id"].(string),
		}
	}
	return result
}

func GetAccessPolicyByPolicyName(client *golangsdk.ServiceClient, policyName string) (*accesspolicies.AccessPolicyDetailInfo, error) {
	policies, err := accesspolicies.List(client)
	if err != nil {
		return nil, err
	}
	if len(policies) < 1 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/access-policy",
				RequestId: "NONE",
				Body:      []byte("all access policies have been deleted"),
			},
		}
	}

	for _, policy := range policies {
		if policy.PolicyName == policyName {
			return &policy, nil
		}
	}
	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Method:    "GET",
			URL:       "/v2/{project_id}/access-policy",
			RequestId: "NONE",
			Body:      []byte(fmt.Sprintf("no access policy matched the name '%s'", policyName)),
		},
	}
}

func resourceAccessPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WorkspaceV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	var (
		policyName = d.Get("name").(string)
		opts       = accesspolicies.CreateOpts{
			Policy: accesspolicies.AccessPolicy{
				PolicyName:    policyName,
				BlacklistType: d.Get("blacklist_type").(string),
			},
			PolicyObjectsList: buildAccessPolicyObjects(d.Get("blacklist").(*schema.Set)),
		}
	)
	err = accesspolicies.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating Workspace access policy: %s", err)
	}
	expectedPolicy, err := GetAccessPolicyByPolicyName(client, policyName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Workspace access policy")
	}
	d.SetId(expectedPolicy.PolicyId)

	return resourceAccessPolicyRead(ctx, d, meta)
}

func flattenAccessPolicyObjects(objects []accesspolicies.AccessPolicyObject) []interface{} {
	if len(objects) < 1 {
		return nil
	}
	result := make([]interface{}, len(objects))
	for i, object := range objects {
		result[i] = map[string]interface{}{
			"object_type": object.ObjectType,
			"object_id":   object.ObjectId,
			"object_name": object.ObjectName,
		}
	}
	return result
}

func resourceAccessPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.WorkspaceV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	expectedPolicy, err := GetAccessPolicyByPolicyName(client, d.Get("name").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Workspace access policy")
	}

	var (
		policyId = d.Id()
		opts     = accesspolicies.ListObjectsOpts{
			PolicyId: policyId,
		}
	)
	objects, err := accesspolicies.ListObjects(client, opts)
	if err != nil {
		return diag.Errorf("error querying object list of specified access policy (%s): %s", policyId, err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", expectedPolicy.PolicyName),
		d.Set("blacklist_type", expectedPolicy.BlacklistType),
		d.Set("blacklist", flattenAccessPolicyObjects(objects)),
		d.Set("created_at", expectedPolicy.CreateTime),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAccessPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WorkspaceV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	var (
		policyId = d.Id()
		opts     = accesspolicies.UpdateOpts{
			PolicyId:          policyId,
			PolicyObjectsList: buildAccessPolicyObjects(d.Get("blacklist").(*schema.Set)),
		}
	)
	err = accesspolicies.UpdateObjects(client, opts)
	if err != nil {
		return diag.Errorf("error updating Workspace access policy (%s): %s", policyId, err)
	}
	return resourceAccessPolicyRead(ctx, d, meta)
}

func resourceAccessPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WorkspaceV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	opts := accesspolicies.DeleteOpts{
		PolicyIdList: []string{d.Id()},
	}
	err = accesspolicies.Delete(client, opts)
	if err != nil {
		return diag.Errorf("error deleting Workspace access policy (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceAccessPolicyImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conf := meta.(*config.Config)
	client, err := conf.WorkspaceV2Client(conf.GetRegion(d))
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating Workspace v2 client: %s", err)
	}
	importedPolicyName := d.Id()
	expectedPolicy, err := GetAccessPolicyByPolicyName(client, importedPolicyName)
	if err != nil {
		return []*schema.ResourceData{d}, err
	}
	d.SetId(expectedPolicy.PolicyId)

	return []*schema.ResourceData{d}, d.Set("name", importedPolicyName)
}
