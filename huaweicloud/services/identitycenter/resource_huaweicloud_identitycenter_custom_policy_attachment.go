// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product IdentityCenter
// ---------------------------------------------------------------

package identitycenter

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IdentityCenter DELETE /v1/instances/{instance_id}/permission-sets/{permission_set_id}/custom-policy
// @API IdentityCenter GET /v1/instances/{instance_id}/permission-sets/{permission_set_id}/custom-policy
// @API IdentityCenter PUT /v1/instances/{instance_id}/permission-sets/{permission_set_id}/custom-policy
func ResourceCustomPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomPolicyAttachmentCreateOrUpdate,
		UpdateContext: resourceCustomPolicyAttachmentCreateOrUpdate,
		ReadContext:   resourceCustomPolicyAttachmentRead,
		DeleteContext: resourceCustomPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCustomPolicyAttachmentImport,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the IAM Identity Center instance.`,
			},
			"permission_set_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the IAM Identity Center permission set.`,
			},
			"custom_policy": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the custom policy to attach to a permission set.`,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := utils.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
		},
	}
}

func resourceCustomPolicyAttachmentCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// customPolicyAttachment: attach custom policy to permission set
	var (
		customPolicyAttachmentHttpUrl = "v1/instances/{instance_id}/permission-sets/{permission_set_id}/custom-policy"
		customPolicyAttachmentProduct = "identitycenter"
	)
	customPolicyAttachmentClient, err := cfg.NewServiceClient(customPolicyAttachmentProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	psID := d.Get("permission_set_id").(string)

	customPolicyAttachmentPath := customPolicyAttachmentClient.Endpoint + customPolicyAttachmentHttpUrl
	customPolicyAttachmentPath = strings.ReplaceAll(customPolicyAttachmentPath, "{instance_id}",
		d.Get("instance_id").(string))
	customPolicyAttachmentPath = strings.ReplaceAll(customPolicyAttachmentPath, "{permission_set_id}",
		d.Get("permission_set_id").(string))

	customPolicyAttachmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	customPolicyAttachmentOpt.JSONBody = utils.RemoveNil(buildCustomPolicyAttachmentBodyParams(d))
	_, err = customPolicyAttachmentClient.Request("PUT", customPolicyAttachmentPath, &customPolicyAttachmentOpt)
	if err != nil {
		return diag.Errorf("error creating/updating Identity Center custom policy attachment: %s", err)
	}

	d.SetId(psID)

	if diagErr := provisionPermissionSet(customPolicyAttachmentClient, instanceID, psID); diagErr != nil {
		return diagErr
	}

	return resourceCustomPolicyAttachmentRead(ctx, d, meta)
}

func buildCustomPolicyAttachmentBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"custom_policy": d.Get("custom_policy"),
	}
	return bodyParams
}

func resourceCustomPolicyAttachmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getCustomPolicyAttachment: query custom policy of the permission set
	var (
		getCustomPolicyAttachmentHttpUrl = "v1/instances/{instance_id}/permission-sets/{permission_set_id}/custom-policy"
		getCustomPolicyAttachmentProduct = "identitycenter"
	)
	getCustomPolicyAttachmentClient, err := cfg.NewServiceClient(getCustomPolicyAttachmentProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	getCustomPolicyAttachmentPath := getCustomPolicyAttachmentClient.Endpoint + getCustomPolicyAttachmentHttpUrl
	getCustomPolicyAttachmentPath = strings.ReplaceAll(getCustomPolicyAttachmentPath, "{instance_id}",
		d.Get("instance_id").(string))
	getCustomPolicyAttachmentPath = strings.ReplaceAll(getCustomPolicyAttachmentPath, "{permission_set_id}",
		d.Get("permission_set_id").(string))

	getCustomPolicyAttachmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getCustomPolicyAttachmentResp, err := getCustomPolicyAttachmentClient.Request("GET",
		getCustomPolicyAttachmentPath, &getCustomPolicyAttachmentOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center custom policy attachment")
	}

	getCustomPolicyAttachmentRespBody, err := utils.FlattenResponse(getCustomPolicyAttachmentResp)
	if err != nil {
		return diag.FromErr(err)
	}

	customPolicy := utils.PathSearch("custom_policy", getCustomPolicyAttachmentRespBody, "").(string)
	if customPolicy == "" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("custom_policy", customPolicy),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCustomPolicyAttachmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteCustomPolicyAttachment: delete custom policy of the permission set
	var (
		deleteCustomPolicyAttachmentHttpUrl = "v1/instances/{instance_id}/permission-sets/{permission_set_id}/custom-policy"
		deleteCustomPolicyAttachmentProduct = "identitycenter"
	)
	deleteCustomPolicyAttachmentClient, err := cfg.NewServiceClient(deleteCustomPolicyAttachmentProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	deleteCustomPolicyAttachmentPath := deleteCustomPolicyAttachmentClient.Endpoint + deleteCustomPolicyAttachmentHttpUrl
	deleteCustomPolicyAttachmentPath = strings.ReplaceAll(deleteCustomPolicyAttachmentPath, "{instance_id}",
		d.Get("instance_id").(string))
	deleteCustomPolicyAttachmentPath = strings.ReplaceAll(deleteCustomPolicyAttachmentPath, "{permission_set_id}",
		d.Get("permission_set_id").(string))

	deleteCustomPolicyAttachmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = deleteCustomPolicyAttachmentClient.Request("DELETE", deleteCustomPolicyAttachmentPath,
		&deleteCustomPolicyAttachmentOpt)
	if err != nil {
		return diag.Errorf("error deleting Identity Center custom policy attachment: %s", err)
	}

	return nil
}

func resourceCustomPolicyAttachmentImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		err := fmt.Errorf("invalid format: the format must be <instance_id>/<permission_set_id>")
		return nil, err
	}

	instanceID := parts[0]
	psID := parts[1]

	d.SetId(psID)
	d.Set("instance_id", instanceID)
	d.Set("permission_set_id", psID)

	return []*schema.ResourceData{d}, nil
}
