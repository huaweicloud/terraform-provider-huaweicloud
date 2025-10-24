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

// @API IdentityCenter GET /v1/instances/{instance_id}/permission-sets/{permission_set_id}/custom-role
// @API IdentityCenter PUT /v1/instances/{instance_id}/permission-sets/{permission_set_id}/custom-role
// @API IdentityCenter DELETE /v1/instances/{instance_id}/permission-sets/{permission_set_id}/custom-role
func ResourceCustomRoleAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomRoleAttachmentCreateOrUpdate,
		UpdateContext: resourceCustomRoleAttachmentCreateOrUpdate,
		ReadContext:   resourceCustomRoleAttachmentRead,
		DeleteContext: resourceCustomRoleAttachmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCustomRoleAttachmentImport,
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
			"custom_role": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the custom role to attach to a permission set.`,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := utils.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
		},
	}
}

func resourceCustomRoleAttachmentCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// addCustomRoleAttachment: attach custom role to permission set
	var (
		customRoleAttachmentHttpUrl = "v1/instances/{instance_id}/permission-sets/{permission_set_id}/custom-role"
		customRoleAttachmentProduct = "identitycenter"
	)
	customRoleAttachmentClient, err := cfg.NewServiceClient(customRoleAttachmentProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	psID := d.Get("permission_set_id").(string)

	customRoleAttachmentPath := customRoleAttachmentClient.Endpoint + customRoleAttachmentHttpUrl
	customRoleAttachmentPath = strings.ReplaceAll(customRoleAttachmentPath, "{instance_id}", instanceID)
	customRoleAttachmentPath = strings.ReplaceAll(customRoleAttachmentPath, "{permission_set_id}", psID)

	customRoleAttachmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	customRoleAttachmentOpt.JSONBody = utils.RemoveNil(buildAddCustomRoleAttachmentBodyParams(d))
	_, err = customRoleAttachmentClient.Request("PUT", customRoleAttachmentPath, &customRoleAttachmentOpt)
	if err != nil {
		return diag.Errorf("error creating/updating Identity Center custom role attachment: %s", err)
	}

	d.SetId(psID)

	if diagErr := provisionPermissionSet(customRoleAttachmentClient, instanceID, psID); diagErr != nil {
		return diagErr
	}

	return resourceCustomRoleAttachmentRead(ctx, d, meta)
}

func buildAddCustomRoleAttachmentBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"custom_role": d.Get("custom_role"),
	}
	return bodyParams
}

func resourceCustomRoleAttachmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getCustomRoleAttachment: query custom role of the permission set
	var (
		getCustomRoleAttachmentHttpUrl = "v1/instances/{instance_id}/permission-sets/{permission_set_id}/custom-role"
		getCustomRoleAttachmentProduct = "identitycenter"
	)
	getCustomRoleAttachmentClient, err := cfg.NewServiceClient(getCustomRoleAttachmentProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	getCustomRoleAttachmentPath := getCustomRoleAttachmentClient.Endpoint + getCustomRoleAttachmentHttpUrl
	getCustomRoleAttachmentPath = strings.ReplaceAll(getCustomRoleAttachmentPath, "{instance_id}",
		d.Get("instance_id").(string))
	getCustomRoleAttachmentPath = strings.ReplaceAll(getCustomRoleAttachmentPath, "{permission_set_id}",
		d.Get("permission_set_id").(string))

	getCustomRoleAttachmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getCustomRoleAttachmentResp, err := getCustomRoleAttachmentClient.Request("GET", getCustomRoleAttachmentPath,
		&getCustomRoleAttachmentOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center custom role attachment")
	}

	getCustomRoleAttachmentRespBody, err := utils.FlattenResponse(getCustomRoleAttachmentResp)
	if err != nil {
		return diag.FromErr(err)
	}

	customRole := utils.PathSearch("custom_role", getCustomRoleAttachmentRespBody, "").(string)
	if customRole == "" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("custom_role", utils.PathSearch("custom_role", getCustomRoleAttachmentRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCustomRoleAttachmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteCustomRoleAttachment: delete custom role of the permission set
	var (
		deleteCustomRoleAttachmentHttpUrl = "v1/instances/{instance_id}/permission-sets/{permission_set_id}/custom-role"
		deleteCustomRoleAttachmentProduct = "identitycenter"
	)
	deleteCustomRoleAttachmentClient, err := cfg.NewServiceClient(deleteCustomRoleAttachmentProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	deleteCustomRoleAttachmentPath := deleteCustomRoleAttachmentClient.Endpoint + deleteCustomRoleAttachmentHttpUrl
	deleteCustomRoleAttachmentPath = strings.ReplaceAll(deleteCustomRoleAttachmentPath, "{instance_id}",
		d.Get("instance_id").(string))
	deleteCustomRoleAttachmentPath = strings.ReplaceAll(deleteCustomRoleAttachmentPath, "{permission_set_id}",
		d.Get("permission_set_id").(string))

	deleteCustomRoleAttachmentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = deleteCustomRoleAttachmentClient.Request("DELETE", deleteCustomRoleAttachmentPath,
		&deleteCustomRoleAttachmentOpt)
	if err != nil {
		return diag.Errorf("error deleting Identity Center custom role attachment: %s", err)
	}

	return nil
}

func resourceCustomRoleAttachmentImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
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
