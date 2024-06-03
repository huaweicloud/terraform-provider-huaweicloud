// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CNAD
// ---------------------------------------------------------------

package cnad

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AAD POST /v1/cnad/policies/{policy_id}/bind
// @API AAD POST /v1/cnad/policies/{policy_id}/unbind
// @API AAD GET /v1/cnad/protected-ips
func ResourcePolicyAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyAssociateCreate,
		UpdateContext: resourcePolicyAssociateUpdate,
		ReadContext:   resourcePolicyAssociateRead,
		DeleteContext: resourcePolicyAssociateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePolicyAssociateImportState,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `Specifies the instance ID. This field must be the instance ID where the policy is
located.`,
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the CNAD advanced policy ID in which to associate protected objects.`,
			},
			"protected_object_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `Specifies the protected object IDs to associate.`,
			},
			"protected_objects": {
				Type:        schema.TypeList,
				Elem:        associateProtectedObjectSchema(),
				Computed:    true,
				Description: `The advanced protected objects.`,
			},
		},
	}
}

func associateProtectedObjectSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the protected object.`,
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the IP of the protected object.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the type of the protected object.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the protected object.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the instance ID of the protected object.`,
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the instance name of the protected object.`,
			},
			"instance_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the instance version of the protected object.`,
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the region to which the protected object belongs.`,
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the status of the protected object.`,
			},
			"block_threshold": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the block threshold of the protected object.`,
			},
			"clean_threshold": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the clean threshold of the protected object.`,
			},
			"policy_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the policy name of the protected object.`,
			},
		},
	}
	return &sc
}

func resourcePolicyAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aad", region)
	if err != nil {
		return diag.Errorf("error creating CNAD Client: %s", err)
	}

	ipList := utils.ExpandToStringList(d.Get("protected_object_ids").(*schema.Set).List())
	if err := bindObjectsToPolicy(client, d, ipList); err != nil {
		return diag.FromErr(err)
	}
	id := fmt.Sprintf("%s/%s", d.Get("policy_id"), d.Get("instance_id"))
	d.SetId(id)
	return resourcePolicyAssociateRead(ctx, d, meta)
}

func resourcePolicyAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                       = meta.(*config.Config)
		region                    = cfg.GetRegion(d)
		mErr                      *multierror.Error
		getProtectedObjectHttpUrl = "v1/cnad/protected-ips"
		getProtectedObjectProduct = "aad"
	)
	getProtectedObjectClient, err := cfg.NewServiceClient(getProtectedObjectProduct, region)
	if err != nil {
		return diag.Errorf("error creating CNAD Client: %s", err)
	}

	getProtectedObjectPath := getProtectedObjectClient.Endpoint + getProtectedObjectHttpUrl
	getProtectedObjectPath += buildGetObjectsPolicyQueryParams(d)

	getProtectedObjectsResp, err := pagination.ListAllItems(
		getProtectedObjectClient,
		"offset",
		getProtectedObjectPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		// here is no special error code
		return diag.Errorf("error retrieving policy binding protected objects, %s", err)
	}

	getProtectedObjectsRespJson, err := json.Marshal(getProtectedObjectsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getProtectedObjectsRespBody interface{}
	if err := json.Unmarshal(getProtectedObjectsRespJson, &getProtectedObjectsRespBody); err != nil {
		return diag.FromErr(err)
	}

	protectedObjects, protectedObjectIDs := flattenGetAssociateProtectedObjects(getProtectedObjectsRespBody)
	if len(protectedObjects) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{},
			"error retrieving policy binding protected objects")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("protected_object_ids", protectedObjectIDs),
		d.Set("protected_objects", protectedObjects),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetAssociateProtectedObjects(resp interface{}) ([]interface{}, []string) {
	if resp == nil {
		return nil, nil
	}
	curJson := utils.PathSearch("items", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	ids := make([]string, len(curArray))
	rst := make([]interface{}, len(curArray))
	for i, v := range curArray {
		id := utils.PathSearch("id", v, "")
		ids[i] = id.(string)
		rst[i] = map[string]interface{}{
			"id":               id,
			"ip_address":       utils.PathSearch("ip", v, nil),
			"type":             utils.PathSearch("type", v, nil),
			"name":             utils.PathSearch("name", v, nil),
			"instance_id":      utils.PathSearch("package_id", v, nil),
			"instance_name":    utils.PathSearch("package_name", v, nil),
			"instance_version": utils.PathSearch("package_version", v, nil),
			"region":           utils.PathSearch("region", v, nil),
			"status":           utils.PathSearch("status", v, nil),
			"block_threshold":  utils.PathSearch("block_threshold", v, nil),
			"clean_threshold":  utils.PathSearch("clean_threshold", v, nil),
			"policy_name":      utils.PathSearch("policy_name", v, nil),
		}
	}
	return rst, ids
}

func buildGetObjectsPolicyQueryParams(d *schema.ResourceData) string {
	policyID := d.Get("policy_id")
	instanceID := d.Get("instance_id")
	return fmt.Sprintf("?policy_id=%s&package_id=%s", policyID, instanceID)
}

func resourcePolicyAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aad", region)
	if err != nil {
		return diag.Errorf("error creating CNAD Client: %s", err)
	}

	oldRaws, newRaws := d.GetChange("protected_object_ids")
	unbindIdList := utils.ExpandToStringList(oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set)).List())
	bindIdList := utils.ExpandToStringList(newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set)).List())

	if len(unbindIdList) > 0 {
		if err := unbindObjectsToPolicy(client, d, unbindIdList); err != nil {
			return diag.FromErr(err)
		}
	}
	if len(bindIdList) > 0 {
		if err := bindObjectsToPolicy(client, d, bindIdList); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourcePolicyAssociateRead(ctx, d, meta)
}

func resourcePolicyAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aad", region)
	if err != nil {
		return diag.Errorf("error creating CNAD Client: %s", err)
	}

	ipList := utils.ExpandToStringList(d.Get("protected_object_ids").(*schema.Set).List())
	if err := unbindObjectsToPolicy(client, d, ipList); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func bindObjectsToPolicy(client *golangsdk.ServiceClient, d *schema.ResourceData, idList []string) error {
	policyID := d.Get("policy_id").(string)
	bindPath := client.Endpoint + "v1/cnad/policies/{policy_id}/bind"
	bindPath = strings.ReplaceAll(bindPath, "{policy_id}", policyID)

	bindOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"package_id": utils.ValueIgnoreEmpty(d.Get("instance_id")),
			"id_list":    idList,
		},
	}

	_, err := client.Request("POST", bindPath, &bindOpt)
	if err != nil {
		return fmt.Errorf("error binding protected objects to policy : %s", err)
	}
	return nil
}

func unbindObjectsToPolicy(client *golangsdk.ServiceClient, d *schema.ResourceData, idList []string) error {
	policyID := d.Get("policy_id").(string)
	unbindPath := client.Endpoint + "v1/cnad/policies/{policy_id}/unbind"
	unbindPath = strings.ReplaceAll(unbindPath, "{policy_id}", policyID)

	unbindOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"package_id": utils.ValueIgnoreEmpty(d.Get("instance_id")),
			"id_list":    idList,
		},
	}

	_, err := client.Request("POST", unbindPath, &unbindOpt)
	if err != nil {
		return fmt.Errorf("error unbinding protected objects to policy : %s", err)
	}
	return nil
}

func resourcePolicyAssociateImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	partLength := len(parts)

	if partLength == 2 {
		mErr := multierror.Append(nil,
			d.Set("policy_id", parts[0]),
			d.Set("instance_id", parts[1]),
		)
		if err := mErr.ErrorOrNil(); err != nil {
			return nil, fmt.Errorf("failed to set value to state when import policy associate, %s", err)
		}
		return []*schema.ResourceData{d}, nil
	}
	return nil, fmt.Errorf("invalid format specified for import id, must be <policy_id>/<instance_id>")
}
