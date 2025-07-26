// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product IdentityCenter
// ---------------------------------------------------------------

package identitycenter

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IdentityCenter POST /v1/instances/{instance_id}/permission-sets
// @API IdentityCenter DELETE /v1/instances/{instance_id}/permission-sets/{id}
// @API IdentityCenter GET /v1/instances/{instance_id}/permission-sets/{id}
// @API IdentityCenter PUT /v1/instances/{instance_id}/permission-sets/{id}
// @API IdentityCenter POST /v1/instances/{resource_type}/{resource_id}/tags/create
// @API IdentityCenter POST /v1/instances/{resource_type}/{resource_id}/tags/delete
// @API IdentityCenter GET /v1/instances/{resource_type}/{resource_id}/tags
func ResourcePermissionSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePermissionSetCreate,
		UpdateContext: resourcePermissionSetUpdate,
		ReadContext:   resourcePermissionSetRead,
		DeleteContext: resourcePermissionSetDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePermissionSetImport,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Description: "schema: Internal",
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"session_duration": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"relay_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": common.TagsSchema(),
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourcePermissionSetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createPermissionSetHttpUrl = "v1/instances/{instance_id}/permission-sets"
		createPermissionSetProduct = "identitycenter"
	)
	createPermissionSetClient, err := cfg.NewServiceClient(createPermissionSetProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	createPermissionSetPath := createPermissionSetClient.Endpoint + createPermissionSetHttpUrl
	createPermissionSetPath = strings.ReplaceAll(createPermissionSetPath, "{instance_id}", d.Get("instance_id").(string))

	createPermissionSetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createPermissionSetOpt.JSONBody = utils.RemoveNil(buildCreatePermissionSetBodyParams(d))
	createPermissionSetResp, err := createPermissionSetClient.Request("POST", createPermissionSetPath, &createPermissionSetOpt)
	if err != nil {
		return diag.Errorf("error creating permission set: %s", err)
	}

	createPermissionSetRespBody, err := utils.FlattenResponse(createPermissionSetResp)
	if err != nil {
		return diag.FromErr(err)
	}

	permissionSetId := utils.PathSearch("permission_set.permission_set_id", createPermissionSetRespBody, "").(string)
	if permissionSetId == "" {
		return diag.Errorf("unable to find the Identity Center permission set ID from the API response")
	}

	d.SetId(permissionSetId)

	if _, ok := d.GetOk("tags"); ok {
		if err := updateTags(createPermissionSetClient, d, "identitycenter:permissionset", d.Id()); err != nil {
			return diag.Errorf("error creating tags of Identity Center permission set %s: %s", d.Id(), err)
		}
	}

	return resourcePermissionSetRead(ctx, d, meta)
}

func buildCreatePermissionSetBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":             utils.ValueIgnoreEmpty(d.Get("name")),
		"session_duration": utils.ValueIgnoreEmpty(d.Get("session_duration")),
		"description":      utils.ValueIgnoreEmpty(d.Get("description")),
		"relay_state":      utils.ValueIgnoreEmpty(d.Get("relay_state")),
	}
	return bodyParams
}

func resourcePermissionSetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		getPermissionSetHttpUrl = "v1/instances/{instance_id}/permission-sets/{id}"
		getPermissionSetProduct = "identitycenter"
	)
	getPermissionSetClient, err := cfg.NewServiceClient(getPermissionSetProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	psID := d.Id()
	getPermissionSetPath := getPermissionSetClient.Endpoint + getPermissionSetHttpUrl
	getPermissionSetPath = strings.ReplaceAll(getPermissionSetPath, "{instance_id}", instanceID)
	getPermissionSetPath = strings.ReplaceAll(getPermissionSetPath, "{id}", psID)

	getPermissionSetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getPermissionSetResp, err := getPermissionSetClient.Request("GET", getPermissionSetPath, &getPermissionSetOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving permission set")
	}

	getPermissionSetRespBody, err := utils.FlattenResponse(getPermissionSetResp)
	if err != nil {
		return diag.FromErr(err)
	}

	accountIDs, err := getAssignededAccounts(getPermissionSetClient, instanceID, psID)
	if err != nil {
		log.Printf("[WARN] failed to get accounts assigned to the permission set %s: %s", psID, err)
	}

	timeStamp := utils.PathSearch("permission_set.created_date", getPermissionSetRespBody, float64(0)).(float64)
	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("permission_set.name", getPermissionSetRespBody, nil)),
		d.Set("session_duration", utils.PathSearch("permission_set.session_duration", getPermissionSetRespBody, nil)),
		d.Set("relay_state", utils.PathSearch("permission_set.relay_state", getPermissionSetRespBody, nil)),
		d.Set("description", utils.PathSearch("permission_set.description", getPermissionSetRespBody, nil)),
		d.Set("urn", utils.PathSearch("permission_set.permission_urn", getPermissionSetRespBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(timeStamp)/1000, false)),
		d.Set("account_ids", accountIDs),
	)

	tags, err := getPermissionSetTags(getPermissionSetClient, d.Id())
	if err != nil {
		log.Printf("[WARN] error fetching tags of permission set: %s", err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("tags", tags),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getAssignededAccounts(client *golangsdk.ServiceClient, instanceID, psID string) ([]string, error) {
	requestURI := fmt.Sprintf("v1/instances/%s/permission-sets/%s/accounts",
		instanceID, psID)
	requestPath := client.Endpoint + requestURI

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	response, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return nil, err
	}

	accountsRaw := utils.PathSearch("account_ids", respBody, make([]interface{}, 0))
	return utils.ExpandToStringList(accountsRaw.([]interface{})), nil
}

func resourcePermissionSetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("identitycenter", region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	updatePermissionSetChanges := []string{
		"session_duration",
		"description",
		"relay_state",
	}

	if d.HasChanges(updatePermissionSetChanges...) {
		updatePermissionSetHttpUrl := "v1/instances/{instance_id}/permission-sets/{id}"

		updatePermissionSetPath := client.Endpoint + updatePermissionSetHttpUrl
		updatePermissionSetPath = strings.ReplaceAll(updatePermissionSetPath, "{instance_id}", d.Get("instance_id").(string))
		updatePermissionSetPath = strings.ReplaceAll(updatePermissionSetPath, "{id}", d.Id())

		updatePermissionSetOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		updatePermissionSetOpt.JSONBody = utils.RemoveNil(buildUpdatePermissionSetBodyParams(d))
		_, err = client.Request("PUT", updatePermissionSetPath, &updatePermissionSetOpt)
		if err != nil {
			return diag.Errorf("error updating permission set: %s", err)
		}
	}

	if d.HasChange("tags") {
		if err := updateTags(client, d, "identitycenter:permissionset", d.Id()); err != nil {
			return diag.Errorf("error updating tags of Identitycenter permission set %s: %s", d.Id(), err)
		}
	}

	return resourcePermissionSetRead(ctx, d, meta)
}

func buildUpdatePermissionSetBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"session_duration": utils.ValueIgnoreEmpty(d.Get("session_duration")),
		"relay_state":      utils.ValueIgnoreEmpty(d.Get("relay_state")),
		// the description parameter can be cleared
		"description": d.Get("description"),
	}
	return bodyParams
}

func resourcePermissionSetDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deletePermissionSetHttpUrl = "v1/instances/{instance_id}/permission-sets/{id}"
		deletePermissionSetProduct = "identitycenter"
	)
	deletePermissionSetClient, err := cfg.NewServiceClient(deletePermissionSetProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	deletePermissionSetPath := deletePermissionSetClient.Endpoint + deletePermissionSetHttpUrl
	deletePermissionSetPath = strings.ReplaceAll(deletePermissionSetPath, "{instance_id}", d.Get("instance_id").(string))
	deletePermissionSetPath = strings.ReplaceAll(deletePermissionSetPath, "{id}", d.Id())

	deletePermissionSetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	_, err = deletePermissionSetClient.Request("DELETE", deletePermissionSetPath, &deletePermissionSetOpt)
	if err != nil {
		return diag.Errorf("error deleting permission set: %s", err)
	}

	return nil
}

func resourcePermissionSetImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		err := fmt.Errorf("invalid format: the format must be <instance id>/<permission set id>")
		return nil, err
	}

	instanceID := parts[0]
	psID := parts[1]

	d.SetId(psID)
	d.Set("instance_id", instanceID)

	return []*schema.ResourceData{d}, nil
}

func updateTags(client *golangsdk.ServiceClient, d *schema.ResourceData, tagsType string, id string) error {
	oRaw, nRaw := d.GetChange("tags")
	oMap := oRaw.(map[string]interface{})
	nMap := nRaw.(map[string]interface{})

	manageTagsHttpUrl := "v1/instances/{resource_type}/{resource_id}/tags/{action}"
	manageTagsPath := client.Endpoint + manageTagsHttpUrl
	manageTagsPath = strings.ReplaceAll(manageTagsPath, "{resource_type}", tagsType)
	manageTagsPath = strings.ReplaceAll(manageTagsPath, "{resource_id}", id)
	manageTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// remove old tags
	if len(oMap) > 0 {
		manageDeleteTagsPath := strings.ReplaceAll(manageTagsPath, "{action}", "delete")
		manageTagsOpt.JSONBody = map[string]interface{}{
			"tags": utils.ExpandResourceTags(oMap),
		}
		_, err := client.Request("POST", manageDeleteTagsPath, &manageTagsOpt)
		if err != nil {
			return err
		}
	}

	// set new tags
	if len(nMap) > 0 {
		manageCreateTagsPath := strings.ReplaceAll(manageTagsPath, "{action}", "create")
		manageTagsOpt.JSONBody = map[string]interface{}{
			"tags": utils.ExpandResourceTags(nMap),
		}
		_, err := client.Request("POST", manageCreateTagsPath, &manageTagsOpt)
		if err != nil {
			return err
		}
	}

	return nil
}

func getPermissionSetTags(client *golangsdk.ServiceClient, id string) (interface{}, error) {
	getPermissionSetTagsHttpUrl := "v1/instances/{resource_type}/{resource_id}/tags?limit=10"
	getPermissionSetTagsPath := client.Endpoint + getPermissionSetTagsHttpUrl
	getPermissionSetTagsPath = strings.ReplaceAll(getPermissionSetTagsPath, "{resource_type}", "identitycenter:permissionset")
	getPermissionSetTagsPath = strings.ReplaceAll(getPermissionSetTagsPath, "{resource_id}", id)

	getPermissionSetTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	path := getPermissionSetTagsPath
	var permissionSetTags []interface{}
	for {
		getPermissionSetTagsResp, err := client.Request("GET", path, &getPermissionSetTagsOpt)
		if err != nil {
			return nil, err
		}
		getPermissionSetTagsRespBody, err := utils.FlattenResponse(getPermissionSetTagsResp)
		if err != nil {
			return nil, err
		}
		tags := utils.PathSearch("tags", getPermissionSetTagsRespBody, make([]interface{}, 0)).([]interface{})
		permissionSetTags = append(permissionSetTags, tags...)

		marker := utils.PathSearch("page_info.next_marker", getPermissionSetTagsRespBody, nil)
		if marker == nil {
			break
		}
		path = fmt.Sprintf("%s&marker=%s", getPermissionSetTagsPath, marker)
	}

	result := make(map[string]interface{})
	for _, val := range permissionSetTags {
		valMap := val.(map[string]interface{})
		result[valMap["key"].(string)] = valMap["value"]
	}

	return result, nil
}
