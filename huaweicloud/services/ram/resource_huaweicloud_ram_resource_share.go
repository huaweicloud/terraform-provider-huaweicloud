// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RAM
// ---------------------------------------------------------------

package ram

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

// @API RAM POST /v1/resource-shares
// @API RAM PUT /v1/resource-shares/{resource_share_id}
// @API RAM POST /v1/resource-shares/{resource_share_id}/tags/create
// @API RAM POST /v1/resource-shares/{resource_share_id}/tags/delete
// @API RAM POST /v1/resource-shares/{resource_share_id}/disassociate
// @API RAM POST /v1/resource-shares/{resource_share_id}/associate
// @API RAM POST /v1/resource-shares/search
// @API RAM POST /v1/resource-share-associations/search
// @API RAM GET /v1/resource-shares/{resource_share_id}/associated-permissions
// @API RAM DELETE /v1/resource-shares/{resource_share_id}
func ResourceRAMShare() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRAMShareCreate,
		UpdateContext: resourceRAMShareUpdate,
		ReadContext:   resourceRAMShareRead,
		DeleteContext: resourceRAMShareDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the resource share.`,
			},
			"principals": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `Specifies the list of one or more principals associated with the resource share.`,
			},
			"resource_urns": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `Specifies the list of URNs of one or more resources associated with the resource share.`,
			},
			"permission_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `Specifies the list of RAM permissions associated with the resource share.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the description of the resource share.`,
			},
			"allow_external_principals": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Specifies whether resources can be shared with any accounts outside the organization.`,
			},
			"tags": common.TagsSchema(),
			"owning_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The owning account ID of the RAM share.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the RAM share.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the RAM share.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the RAM share.`,
			},
			"associated_permissions": {
				Type:     schema.TypeList,
				Elem:     associatedPermissionsSchema(),
				Computed: true,
			},
		},
	}
}

func associatedPermissionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"permission_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The permission ID.`,
			},
			"permission_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The permission name.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource type of the permission.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the permission.`,
			},
		},
	}
	return &sc
}

func resourceRAMShareCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createRAMShare: create a RAM share.
	var (
		createRAMShareHttpUrl = "v1/resource-shares"
		createRAMShareProduct = "ram"
	)
	createRAMShareClient, err := cfg.NewServiceClient(createRAMShareProduct, region)
	if err != nil {
		return diag.Errorf("error creating RAM Client: %s", err)
	}

	createRAMSharePath := createRAMShareClient.Endpoint + createRAMShareHttpUrl
	createRAMShareOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createRAMShareOpt.JSONBody = utils.RemoveNil(buildCreateRAMShareBodyParams(d))
	createRAMShareResp, err := createRAMShareClient.Request("POST", createRAMSharePath, &createRAMShareOpt)
	if err != nil {
		return diag.Errorf("error creating RAM share: %s", err)
	}

	createRAMShareRespBody, err := utils.FlattenResponse(createRAMShareResp)
	if err != nil {
		return diag.FromErr(err)
	}

	shareId := utils.PathSearch("resource_share.id", createRAMShareRespBody, "").(string)
	if shareId == "" {
		return diag.Errorf("unable to find the RAM share ID from the API response")
	}
	d.SetId(shareId)
	return resourceRAMShareRead(ctx, d, meta)
}

func buildCreateRAMShareBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":           utils.ValueIgnoreEmpty(d.Get("name")),
		"description":    utils.ValueIgnoreEmpty(d.Get("description")),
		"permission_ids": utils.ValueIgnoreEmpty(d.Get("permission_ids").(*schema.Set).List()),
		"principals":     utils.ValueIgnoreEmpty(d.Get("principals").(*schema.Set).List()),
		"resource_urns":  utils.ValueIgnoreEmpty(d.Get("resource_urns").(*schema.Set).List()),
		"tags":           utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
	}

	// The default value of this field when created is true. We do not actively add a default value.
	if !d.Get("allow_external_principals").(bool) {
		bodyParams["allow_external_principals"] = false
	}

	return bodyParams
}

func resourceRAMShareRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	shareClient, err := cfg.NewServiceClient("ram", region)
	if err != nil {
		return diag.Errorf("error creating RAM Client: %s", err)
	}
	// Search share instance information
	if err = setRAMShareInstance(shareClient, d); err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RAM share")
	}

	// Search share association resource urns
	if err = setRAMShareAssociations("resource", shareClient, d); err != nil {
		return diag.FromErr(err)
	}

	// Search share association principals
	if err = setRAMShareAssociations("principal", shareClient, d); err != nil {
		return diag.FromErr(err)
	}

	// Search share permissions
	return diag.FromErr(setRAMSharePermissions(shareClient, d))
}

func setRAMShareInstance(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	getRAMShareHttpUrl := "v1/resource-shares/search"
	getRAMSharePath := client.Endpoint + getRAMShareHttpUrl
	getRAMShareOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"resource_share_ids": []string{d.Id()},
			"resource_owner":     "self",
		},
	}

	getRAMShareResp, err := client.Request("POST", getRAMSharePath, &getRAMShareOpt)
	if err != nil {
		// There is no special error code.
		return err
	}

	getRAMShareRespBody, err := utils.FlattenResponse(getRAMShareResp)
	if err != nil {
		return err
	}

	curJson := utils.PathSearch("resource_shares", getRAMShareRespBody, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return golangsdk.ErrDefault404{}
	}
	if len(curArray) > 1 {
		return fmt.Errorf("except retrieving one RAM share, but got %d", len(curArray))
	}

	resourceShare := curArray[0]
	status := utils.PathSearch("status", resourceShare, "")
	if status == "deleted" {
		// The deleted resource share will exist 48 hours with "deleted" status. And will be removed after 48 hours.
		return golangsdk.ErrDefault404{}
	}

	mErr := multierror.Append(
		nil,
		d.Set("name", utils.PathSearch("name", resourceShare, nil)),
		d.Set("description", utils.PathSearch("description", resourceShare, nil)),
		d.Set("owning_account_id", utils.PathSearch("owning_account_id", resourceShare, nil)),
		d.Set("status", status),
		d.Set("created_at", utils.PathSearch("created_at", resourceShare, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", resourceShare, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", resourceShare, nil))),
		d.Set("allow_external_principals", utils.PathSearch("allow_external_principals", resourceShare, nil)),
	)
	return mErr.ErrorOrNil()
}

// buildAssociationsBodyParams The default limit value for paging query is `200`, so the limit value is not configured
// here.
func buildAssociationsBodyParams(d *schema.ResourceData, associationType, nextMarker string) interface{} {
	requestParams := map[string]interface{}{
		"resource_share_ids": []string{d.Id()},
		"association_type":   associationType,
	}

	if nextMarker != "" {
		requestParams["marker"] = nextMarker
	}
	return requestParams
}

// setRAMShareAssociations associationType has two valid values: "resource" and "principal"
func setRAMShareAssociations(associationType string, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl             = "v1/resource-share-associations/search"
		associationJsonPath = "resource_share_associations[?status=='associating' || status=='associated'].associated_entity"
		nextMarker          string
		totalAssociations   []interface{}
	)

	requestPath := client.Endpoint + httpUrl
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestOpt.JSONBody = buildAssociationsBodyParams(d, associationType, nextMarker)
		resp, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			// There is no special error code.
			return err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return err
		}

		associatedEntities := utils.PathSearch(associationJsonPath, respBody, make([]interface{}, 0)).([]interface{})
		if len(associatedEntities) > 0 {
			totalAssociations = append(totalAssociations, associatedEntities...)
		}

		nextMarker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if nextMarker == "" {
			break
		}
	}

	if associationType == "resource" {
		return d.Set("resource_urns", totalAssociations)
	}

	if associationType == "principal" {
		return d.Set("principals", totalAssociations)
	}

	return fmt.Errorf("got an invalid association type: %s when search share associations", associationType)
}

// buildPermissionPathWithQueryParams The default limit value for paging query is `200`, so the limit value is not
// configured here.
func buildPermissionPathWithQueryParams(path, nextMarker string) string {
	if nextMarker == "" {
		return path
	}
	return fmt.Sprintf("%s?marker=%s", path, nextMarker)
}

func setRAMSharePermissions(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl            = "v1/resource-shares/{resource_share_id}/associated-permissions"
		permissionJsonPath = "associated_permissions[?status=='associating' || status=='associated']"
		nextMarker         string
		totalPermissions   []interface{}
	)

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{resource_share_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithMarker := buildPermissionPathWithQueryParams(requestPath, nextMarker)
		resp, err := client.Request("GET", requestPathWithMarker, &requestOpt)
		if err != nil {
			// There is no special error code.
			return err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return err
		}

		permissions := utils.PathSearch(permissionJsonPath, respBody, make([]interface{}, 0)).([]interface{})
		if len(permissions) > 0 {
			totalPermissions = append(totalPermissions, permissions...)
		}

		nextMarker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if nextMarker == "" {
			break
		}
	}

	mErr := multierror.Append(
		nil,
		d.Set("associated_permissions", flattenAssociatedPermissions(totalPermissions)),
	)
	return mErr.ErrorOrNil()
}

func flattenAssociatedPermissions(permissions []interface{}) []interface{} {
	if len(permissions) == 0 {
		return nil
	}

	rst := make([]interface{}, len(permissions))
	for i, v := range permissions {
		rst[i] = map[string]interface{}{
			"permission_id":   utils.PathSearch("permission_id", v, nil),
			"permission_name": utils.PathSearch("permission_name", v, nil),
			"resource_type":   utils.PathSearch("resource_type", v, nil),
			"status":          utils.PathSearch("status", v, nil),
		}
	}
	return rst
}

func resourceRAMShareUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateRAMShareClient, err := cfg.NewServiceClient("ram", region)
	if err != nil {
		return diag.Errorf("error creating RAM Client: %s", err)
	}

	updateRAMShareChanges := []string{
		"name",
		"description",
		"allow_external_principals",
	}

	if d.HasChanges(updateRAMShareChanges...) {
		// updateRAMShare: update the RAM share.
		updateRAMShareHttpUrl := "v1/resource-shares/{resource_share_id}"
		updateRAMSharePath := updateRAMShareClient.Endpoint + updateRAMShareHttpUrl
		updateRAMSharePath = strings.ReplaceAll(updateRAMSharePath, "{resource_share_id}", d.Id())

		updateRAMShareOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		updateRAMShareOpt.JSONBody = utils.RemoveNil(buildUpdateRAMShareBodyParams(d))
		_, err = updateRAMShareClient.Request("PUT", updateRAMSharePath, &updateRAMShareOpt)
		if err != nil {
			return diag.Errorf("error updating RAM share: %s", err)
		}
	}

	if d.HasChange("tags") {
		err = updateRAMShareTags(updateRAMShareClient, d)
		if err != nil {
			return diag.Errorf("error updating RAM share tags: %s", err)
		}
	}

	if d.HasChanges("principals", "resource_urns") {
		oldPrincipalsRaws, newPrincipalsRaws := d.GetChange("principals")
		oldPrincipals := oldPrincipalsRaws.(*schema.Set).Difference(newPrincipalsRaws.(*schema.Set)).List()
		newPrincipals := newPrincipalsRaws.(*schema.Set).Difference(oldPrincipalsRaws.(*schema.Set)).List()

		oldResourceUrnsRaws, newResourceUrnsRaws := d.GetChange("resource_urns")
		oldResourceUrns := oldResourceUrnsRaws.(*schema.Set).Difference(newResourceUrnsRaws.(*schema.Set)).List()
		newResourceUrns := newResourceUrnsRaws.(*schema.Set).Difference(oldResourceUrnsRaws.(*schema.Set)).List()

		if len(oldPrincipals) > 0 || len(oldResourceUrns) > 0 {
			err = disassociatePrincipalsAndResourceUrns(updateRAMShareClient, d.Id(), oldPrincipals, oldResourceUrns)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		if len(newPrincipals) > 0 || len(newResourceUrns) > 0 {
			err = associatePrincipalsAndResourceUrns(updateRAMShareClient, d.Id(), newPrincipals, newResourceUrns)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceRAMShareRead(ctx, d, meta)
}

func updateRAMShareTags(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	ramShareTagsHttpUrl := "v1/resource-shares/{resource_share_id}/tags"
	ramShareTagsPath := client.Endpoint + ramShareTagsHttpUrl
	ramShareTagsPath = strings.ReplaceAll(ramShareTagsPath, "{resource_share_id}", d.Id())

	oRaw, nRaw := d.GetChange("tags")
	oMap := oRaw.(map[string]interface{})
	nMap := nRaw.(map[string]interface{})

	// remove old tags
	if len(oMap) > 0 {
		deleteTagsPath := fmt.Sprintf("%s/%s", ramShareTagsPath, "delete")
		deleteTagsOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"tags": utils.ExpandResourceTagsMap(oMap),
			},
		}
		_, err := client.Request("POST", deleteTagsPath, &deleteTagsOpt)
		if err != nil {
			return fmt.Errorf("error delete RAM share tags: %s", err)
		}
	}

	// set new tags
	if len(nMap) > 0 {
		createTagsPath := fmt.Sprintf("%s/%s", ramShareTagsPath, "create")
		createTagsOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"tags": utils.ExpandResourceTagsMap(nMap),
			},
		}
		_, err := client.Request("POST", createTagsPath, &createTagsOpt)
		if err != nil {
			return fmt.Errorf("error create RAM share tags: %s", err)
		}
	}
	return nil
}

func buildUpdateRAMShareBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                      utils.ValueIgnoreEmpty(d.Get("name")),
		"description":               utils.ValueIgnoreEmpty(d.Get("description")),
		"allow_external_principals": d.Get("allow_external_principals"),
	}
	return bodyParams
}

func resourceRAMShareDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteRAMShare: delete RAM share
	var (
		deleteRAMShareHttpUrl = "v1/resource-shares/{resource_share_id}"
		deleteRAMShareProduct = "ram"
	)
	deleteRAMShareClient, err := cfg.NewServiceClient(deleteRAMShareProduct, region)
	if err != nil {
		return diag.Errorf("error creating RAM Client: %s", err)
	}

	deleteRAMSharePath := deleteRAMShareClient.Endpoint + deleteRAMShareHttpUrl
	deleteRAMSharePath = strings.ReplaceAll(deleteRAMSharePath, "{resource_share_id}", d.Id())
	deleteRAMShareOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteRAMShareClient.Request("DELETE", deleteRAMSharePath, &deleteRAMShareOpt)
	if err != nil {
		if errCode, ok := err.(golangsdk.ErrDefault400); ok {
			if resp, pErr := common.ParseErrorMsg(errCode.Body); pErr == nil && resp.ErrorCode == "ram.1102" {
				// There are resources in use in the resource share. Do disassociate before delete share
				principals := d.Get("principals").(*schema.Set).List()
				resourceUrns := d.Get("resource_urns").(*schema.Set).List()
				err = disassociatePrincipalsAndResourceUrns(deleteRAMShareClient, d.Id(), principals, resourceUrns)
				if err == nil {
					// retry delete
					return resourceRAMShareDelete(ctx, d, meta)
				}
			}
		}
		return diag.Errorf("error deleting RAM share: %s", err)
	}
	return nil
}

func associatePrincipalsAndResourceUrns(client *golangsdk.ServiceClient, resourceId string, principals,
	resourceUrns []interface{}) error {
	associateHttpUrl := "v1/resource-shares/{resource_share_id}/associate"
	associatePath := client.Endpoint + associateHttpUrl
	associatePath = strings.ReplaceAll(associatePath, "{resource_share_id}", resourceId)
	associateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"principals":    principals,
			"resource_urns": resourceUrns,
		},
	}
	_, err := client.Request("POST", associatePath, &associateOpt)
	if err != nil {
		return fmt.Errorf("error associate RAM share principals and resource urns, %s", err)
	}
	return nil
}

func disassociatePrincipalsAndResourceUrns(client *golangsdk.ServiceClient, resourceId string, principals,
	resourceUrns []interface{}) error {
	disassociateHttpUrl := "v1/resource-shares/{resource_share_id}/disassociate"
	disassociatePath := client.Endpoint + disassociateHttpUrl
	disassociatePath = strings.ReplaceAll(disassociatePath, "{resource_share_id}", resourceId)
	disassociateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"principals":    principals,
			"resource_urns": resourceUrns,
		},
	}
	_, err := client.Request("POST", disassociatePath, &disassociateOpt)
	if err != nil {
		return fmt.Errorf("error disassociate RAM share principals and resource urns, %s", err)
	}
	return nil
}
