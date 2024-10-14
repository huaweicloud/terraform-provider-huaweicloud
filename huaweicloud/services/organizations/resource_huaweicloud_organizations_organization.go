// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product Organizations
// ---------------------------------------------------------------

package organizations

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	rootType     = "organizations:roots"
	unitType     = "organizations:ous"
	accountsType = "organizations:accounts"
	policiesType = "organizations:policies"
)

// @API Organizations POST /v1/organizations
// @API Organizations GET /v1/organizations/roots
// @API Organizations POST /v1/organizations/{resource_type}/{resource_id}/tags/create
// @API Organizations POST /v1/organizations/{resource_type}/{resource_id}/tags/delete
// @API Organizations POST /v1/organizations/policies/enable
// @API Organizations GET /v1/organizations
// @API Organizations GET /v1/organizations/{resource_type}/{resource_id}/tags
// @API Organizations POST /v1/organizations/policies/disable
// @API Organizations DELETE /v1/organizations
func ResourceOrganization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationCreate,
		UpdateContext: resourceOrganizationUpdate,
		ReadContext:   resourceOrganizationRead,
		DeleteContext: resourceOrganizationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"enabled_policy_types": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of Organizations policy types to enable in the Organization Root.`,
			},
			"root_tags": common.TagsSchema(),
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the uniform resource name of the organization.`,
			},
			"master_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the unique ID of the organization's management account.`,
			},
			"master_account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the organization's management account.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time when the organization was created.`,
			},
			"root_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the root.`,
			},
			"root_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the root.`,
			},
			"root_urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the urn of the root.`,
			},
		},
	}
}

func resourceOrganizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createOrganization: create Organizations organization
	var (
		createOrganizationHttpUrl = "v1/organizations"
		createOrganizationProduct = "organizations"
	)
	createOrganizationClient, err := cfg.NewServiceClient(createOrganizationProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	createOrganizationPath := createOrganizationClient.Endpoint + createOrganizationHttpUrl

	createOrganizationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	createOrganizationResp, err := createOrganizationClient.Request("POST", createOrganizationPath,
		&createOrganizationOpt)
	if err != nil {
		return diag.Errorf("error creating Organizations organization: %s", err)
	}

	createOrganizationRespBody, err := utils.FlattenResponse(createOrganizationResp)
	if err != nil {
		return diag.FromErr(err)
	}

	organizationId := utils.PathSearch("organization.id", createOrganizationRespBody, "").(string)
	if organizationId == "" {
		return diag.Errorf("unable to find the organization ID from the API response")
	}
	d.SetId(organizationId)

	getRootRespBody, err := getRoot(createOrganizationClient)
	if err != nil {
		return diag.FromErr(err)
	}
	rootId := utils.PathSearch("roots|[0].id", getRootRespBody, "").(string)

	if v, ok := d.GetOk("root_tags"); ok {
		tagList := utils.ExpandResourceTags(v.(map[string]interface{}))
		err = addTags(createOrganizationClient, rootType, rootId, tagList)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := d.GetOk("enabled_policy_types"); ok {
		enabledPolicyTypes := v.(*schema.Set).List()
		for _, enabledPolicyType := range enabledPolicyTypes {
			if err = enablePolicy(ctx, d.Timeout(schema.TimeoutCreate), createOrganizationClient,
				enabledPolicyType.(string), rootId); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceOrganizationRead(ctx, d, meta)
}

func resourceOrganizationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getOrganization: Query Organizations organization
	var (
		getOrganizationProduct = "organizations"
	)
	getOrganizationClient, err := cfg.NewServiceClient(getOrganizationProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	getOrganizationRespBody, err := getOrganization(getOrganizationClient)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Organizations organization")
	}

	getRootRespBody, diagErr := getRoot(getOrganizationClient)
	if diagErr != nil {
		return common.CheckDeletedDiag(d, diagErr, "error retrieving Organizations root")
	}

	rootId := utils.PathSearch("roots|[0].id", getRootRespBody, "").(string)

	policyTypes := utils.PathSearch("roots|[0].policy_types[?status=='enabled'].type", getRootRespBody,
		make([]interface{}, 0)).([]interface{})

	mErr = multierror.Append(
		mErr,
		d.Set("urn", utils.PathSearch("organization.urn", getOrganizationRespBody, nil)),
		d.Set("master_account_id", utils.PathSearch("organization.management_account_id",
			getOrganizationRespBody, nil)),
		d.Set("master_account_name", utils.PathSearch("organization.management_account_name",
			getOrganizationRespBody, nil)),
		d.Set("created_at", utils.PathSearch("organization.created_at",
			getOrganizationRespBody, nil)),
		d.Set("root_id", rootId),
		d.Set("root_name", utils.PathSearch("roots|[0].name",
			getRootRespBody, nil)),
		d.Set("root_urn", utils.PathSearch("roots|[0].urn", getRootRespBody, nil)),
		d.Set("enabled_policy_types", policyTypes),
	)

	tagMap, err := getTags(getOrganizationClient, rootType, rootId)
	if err != nil {
		log.Printf("[WARN] error fetching Organizations tags of root (%s): %s", rootId, err)
	} else {
		mErr = multierror.Append(mErr, d.Set("root_tags", tagMap))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func getOrganization(client *golangsdk.ServiceClient) (interface{}, error) {
	var (
		getOrganizationHttpUrl = "v1/organizations"
	)

	getOrganizationPath := client.Endpoint + getOrganizationHttpUrl

	getOrganizationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getOrganizationResp, err := client.Request("GET", getOrganizationPath, &getOrganizationOpt)

	if err != nil {
		return nil, err
	}

	getOrganizationRespBody, err := utils.FlattenResponse(getOrganizationResp)
	if err != nil {
		return nil, err
	}
	return getOrganizationRespBody, nil
}

func getRoot(client *golangsdk.ServiceClient) (interface{}, error) {
	var (
		getRootHttpUrl = "v1/organizations/roots"
	)

	getRootPath := client.Endpoint + getRootHttpUrl

	getRootOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRootResp, err := client.Request("GET", getRootPath, &getRootOpt)

	if err != nil {
		return nil, err
	}

	getRootRespBody, err := utils.FlattenResponse(getRootResp)
	if err != nil {
		return nil, err
	}
	return getRootRespBody, nil
}

func enablePolicy(ctx context.Context, timeout time.Duration, client *golangsdk.ServiceClient, policyType,
	rootId string) error {
	err := requestRootPolicy(client, policyType, rootId, "enable")
	if err != nil {
		return err
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"pending_enable"},
		Target:       []string{"enabled"},
		Refresh:      policyStateRefreshFunc(client, policyType),
		Timeout:      timeout,
		Delay:        1 * time.Second,
		PollInterval: 1 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for Organizations policy type (%s) to be enabled: %s", policyType, err)
	}
	return nil
}

func disablePolicy(ctx context.Context, timeout time.Duration, client *golangsdk.ServiceClient, policyType,
	rootId string) error {
	err := requestRootPolicy(client, policyType, rootId, "disable")
	if err != nil {
		return err
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"pending_disable"},
		Target:       []string{"disabled"},
		Refresh:      policyStateRefreshFunc(client, policyType),
		Timeout:      timeout,
		Delay:        1 * time.Second,
		PollInterval: 1 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for Organizations policy type (%s) to be disabled: %s", policyType, err)
	}
	return nil
}

func requestRootPolicy(client *golangsdk.ServiceClient, policyType, rootId, action string) error {
	requestRootPolicyHttpUrl := fmt.Sprintf("v1/organizations/policies/%s", action)

	requestRootPolicyPath := client.Endpoint + requestRootPolicyHttpUrl

	requestRootPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestRootPolicyOpt.JSONBody = utils.RemoveNil(buildRequestRootPolicyBodyParams(policyType, rootId))
	_, err := client.Request("POST", requestRootPolicyPath, &requestRootPolicyOpt)

	if err != nil {
		return err
	}
	return nil
}

func buildRequestRootPolicyBodyParams(policyType, rootId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"policy_type": policyType,
		"root_id":     rootId,
	}
	return bodyParams
}

func policyStateRefreshFunc(client *golangsdk.ServiceClient, policyType string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getRootRespBody, err := getRoot(client)
		if err != nil {
			return nil, "", err
		}

		enabled := utils.PathSearch(fmt.Sprintf("roots|[0].policy_types[?type=='%s'].status|[0]",
			policyType), getRootRespBody, "").(string)
		if err != nil {
			return nil, "", err
		}

		return getRootRespBody, enabled, nil
	}
}

func resourceOrganizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updateOrganization: update Organizations organization
	var (
		updateOrganizationProduct = "organizations"
	)
	updateOrganizationClient, err := cfg.NewServiceClient(updateOrganizationProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	getRootRespBody, err := getRoot(updateOrganizationClient)
	if err != nil {
		return diag.FromErr(err)
	}
	rootId := utils.PathSearch("roots|[0].id", getRootRespBody, "").(string)

	if d.HasChange("root_tags") {
		if err = updateTags(d, updateOrganizationClient, rootType, rootId, "root_tags"); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("enabled_policy_types") {
		oldRaw, newRaw := d.GetChange("enabled_policy_types")
		enabledPolicyTypes := newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set))
		disabledPolicyTypes := oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set))

		for _, enabledPolicyType := range enabledPolicyTypes.List() {
			if err = enablePolicy(ctx, d.Timeout(schema.TimeoutUpdate), updateOrganizationClient,
				enabledPolicyType.(string), rootId); err != nil {
				return diag.FromErr(err)
			}
		}

		for _, disabledPolicyType := range disabledPolicyTypes.List() {
			if err = disablePolicy(ctx, d.Timeout(schema.TimeoutUpdate), updateOrganizationClient,
				disabledPolicyType.(string), rootId); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceOrganizationRead(ctx, d, meta)
}

func resourceOrganizationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteOrganization: Delete Organizations organization
	var (
		deleteOrganizationHttpUrl = "v1/organizations"
		deleteOrganizationProduct = "organizations"
	)
	deleteOrganizationClient, err := cfg.NewServiceClient(deleteOrganizationProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	deleteOrganizationPath := deleteOrganizationClient.Endpoint + deleteOrganizationHttpUrl

	deleteOrganizationOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = deleteOrganizationClient.Request("DELETE", deleteOrganizationPath, &deleteOrganizationOpt)
	if err != nil {
		return diag.Errorf("error deleting Organizations organization: %s", err)
	}

	return nil
}

func addTags(client *golangsdk.ServiceClient, resourceType, resourceId string, tagList []tags.ResourceTag) error {
	var (
		addTagsToHttpUrl = "v1/organizations/{resource_type}/{resource_id}/tags/create"
	)

	addTagsToPath := client.Endpoint + addTagsToHttpUrl
	addTagsToPath = strings.ReplaceAll(addTagsToPath, "{resource_type}", resourceType)
	addTagsToPath = strings.ReplaceAll(addTagsToPath, "{resource_id}", resourceId)

	addTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	addTagsOpt.JSONBody = utils.RemoveNil(buildTagsBodyParams(tagList))
	_, err := client.Request("POST", addTagsToPath, &addTagsOpt)
	if err != nil {
		return fmt.Errorf("error creating tags of resourceType (%s): %s", resourceType, err)
	}

	return nil
}

func deleteTags(client *golangsdk.ServiceClient, resourceType, resourceId string, tagList []tags.ResourceTag) error {
	var (
		addTagsHttpUrl = "v1/organizations/{resource_type}/{resource_id}/tags/delete"
	)

	addTagsPath := client.Endpoint + addTagsHttpUrl
	addTagsPath = strings.ReplaceAll(addTagsPath, "{resource_type}", resourceType)
	addTagsPath = strings.ReplaceAll(addTagsPath, "{resource_id}", resourceId)

	addTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	addTagsOpt.JSONBody = utils.RemoveNil(buildTagsBodyParams(tagList))
	_, err := client.Request("POST", addTagsPath, &addTagsOpt)
	if err != nil {
		return fmt.Errorf("error deleting tags of resourceType (%s): %s", resourceType, err)
	}

	return nil
}

func buildTagsBodyParams(tagList []tags.ResourceTag) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"tags": tagList,
	}
	return bodyParams
}

func updateTags(d *schema.ResourceData, client *golangsdk.ServiceClient, resourceType, resourceId, tagsName string) error {
	oRaw, nRaw := d.GetChange(tagsName)
	oMap := oRaw.(map[string]interface{})
	nMap := nRaw.(map[string]interface{})
	if len(oMap) > 0 {
		tagList := utils.ExpandResourceTags(oMap)
		err := deleteTags(client, resourceType, resourceId, tagList)
		if err != nil {
			return err
		}
	}
	if len(nMap) > 0 {
		tagList := utils.ExpandResourceTags(nMap)
		err := addTags(client, resourceType, resourceId, tagList)
		if err != nil {
			return err
		}
	}

	return nil
}

func getTags(client *golangsdk.ServiceClient, resourceType, resourceId string) (map[string]string, error) {
	var (
		getTagsHttpUrl = "v1/organizations/{resource_type}/{resource_id}/tags"
	)

	getTagsPath := client.Endpoint + getTagsHttpUrl
	getTagsPath = strings.ReplaceAll(getTagsPath, "{resource_type}", resourceType)
	getTagsPath = strings.ReplaceAll(getTagsPath, "{resource_id}", resourceId)
	getTagsPath += buildQueryTagsParam()

	getTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getTagsResp, err := client.Request("GET", getTagsPath, &getTagsOpt)
	if err != nil {
		return nil, fmt.Errorf("error get Organizations tags: %s", err)
	}

	getTagsRespBody, err := utils.FlattenResponse(getTagsResp)
	if err != nil {
		return nil, err
	}

	curJson := utils.PathSearch("tags", getTagsRespBody, nil)
	if curJson == nil {
		return nil, fmt.Errorf("error get tags by resourceId (%s) and resourceType (%s)", resourceId, resourceType)
	}

	result := make(map[string]string)
	for _, v := range curJson.([]interface{}) {
		key := utils.PathSearch("key", v, "").(string)
		value := utils.PathSearch("value", v, "").(string)
		result[key] = value
	}

	return result, nil
}

func buildQueryTagsParam() string {
	return fmt.Sprintf("?limit=%d", 2000)
}
