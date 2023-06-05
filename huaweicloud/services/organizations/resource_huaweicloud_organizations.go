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

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

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

func ResourceOrganizations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationsCreate,
		ReadContext:   resourceOrganizationsRead,
		UpdateContext: resourceOrganizationsUpdate,
		DeleteContext: resourceOrganizationsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"root_tags": common.TagsSchema(),
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the uniform resource name of the organization.`,
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the unique ID of the organization's management account.`,
			},
			"account_name": {
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

func resourceOrganizationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createOrganizations: create Organizations
	var (
		createOrganizationsHttpUrl = "v1/organizations"
		createOrganizationsProduct = "organizations"
	)
	createOrganizationsClient, err := cfg.NewServiceClient(createOrganizationsProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	createOrganizationsPath := createOrganizationsClient.Endpoint + createOrganizationsHttpUrl

	createOrganizationsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	createOrganizationsResp, err := createOrganizationsClient.Request("POST", createOrganizationsPath,
		&createOrganizationsOpt)
	if err != nil {
		return diag.Errorf("error creating Organizations: %s", err)
	}

	createOrganizationsRespBody, err := utils.FlattenResponse(createOrganizationsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("organization.id", createOrganizationsRespBody)
	if err != nil {
		return diag.Errorf("error creating Organizations: ID is not found in API response")
	}
	d.SetId(id.(string))

	if v, ok := d.GetOk("root_tags"); ok {
		getOrganizationsRootRespBody, diagErr := getOrganizationsRoot(d, createOrganizationsClient)
		if diagErr != nil {
			return diagErr
		}
		rootId := utils.PathSearch("roots|[0].id", getOrganizationsRootRespBody, "").(string)
		tagList := utils.ExpandResourceTags(v.(map[string]interface{}))
		err = addTags(createOrganizationsClient, rootType, rootId, tagList)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceOrganizationsRead(ctx, d, meta)
}

func resourceOrganizationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getOrganizations: Query Organizations
	var (
		getOrganizationsProduct = "organizations"
	)
	getOrganizationsClient, err := cfg.NewServiceClient(getOrganizationsProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	getOrganizationsRespBody, diagErr := getOrganizations(d, getOrganizationsClient)
	if diagErr != nil {
		return diagErr
	}

	getOrganizationsRootRespBody, diagErr := getOrganizationsRoot(d, getOrganizationsClient)
	if diagErr != nil {
		return diagErr
	}

	rootId := utils.PathSearch("roots|[0].id", getOrganizationsRootRespBody, "").(string)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("urn", utils.PathSearch("organization.urn", getOrganizationsRespBody, nil)),
		d.Set("account_id", utils.PathSearch("organization.management_account_id",
			getOrganizationsRespBody, nil)),
		d.Set("account_name", utils.PathSearch("organization.management_account_name",
			getOrganizationsRespBody, nil)),
		d.Set("created_at", utils.PathSearch("organization.created_at",
			getOrganizationsRespBody, nil)),
		d.Set("root_id", rootId),
		d.Set("root_name", utils.PathSearch("roots|[0].name",
			getOrganizationsRootRespBody, nil)),
		d.Set("root_urn", utils.PathSearch("roots|[0].urn", getOrganizationsRootRespBody, nil)),
	)

	tagMap, err := getTags(getOrganizationsClient, rootType, rootId)
	if err != nil {
		log.Printf("[WARN] error fetching tags of Organizations root (%s): %s", rootId, err)
	} else {
		mErr = multierror.Append(mErr, d.Set("root_tags", tagMap))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func getOrganizations(d *schema.ResourceData, client *golangsdk.ServiceClient) (interface{}, diag.Diagnostics) {
	var (
		getOrganizationsHttpUrl = "v1/organizations"
	)

	getOrganizationsPath := client.Endpoint + getOrganizationsHttpUrl

	getOrganizationsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getOrganizationsResp, err := client.Request("GET", getOrganizationsPath, &getOrganizationsOpt)

	if err != nil {
		return nil, common.CheckDeletedDiag(d, err, "error retrieving Organizations")
	}

	getOrganizationsRespBody, err := utils.FlattenResponse(getOrganizationsResp)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return getOrganizationsRespBody, nil
}

func getOrganizationsRoot(d *schema.ResourceData, client *golangsdk.ServiceClient) (interface{}, diag.Diagnostics) {
	var (
		getOrganizationsRootHttpUrl = "v1/organizations/roots"
	)

	getOrganizationsRootPath := client.Endpoint + getOrganizationsRootHttpUrl

	getOrganizationsRootOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getOrganizationsRootResp, err := client.Request("GET", getOrganizationsRootPath, &getOrganizationsRootOpt)

	if err != nil {
		return nil, common.CheckDeletedDiag(d, err, "error retrieving Organizations root")
	}

	getOrganizationsRootRespBody, err := utils.FlattenResponse(getOrganizationsRootResp)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return getOrganizationsRootRespBody, nil
}

func resourceOrganizationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updateOrganizations: update Organizations
	var (
		updateOrganizationsProduct = "organizations"
	)
	updateOrganizationsClient, err := cfg.NewServiceClient(updateOrganizationsProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	if d.HasChange("root_tags") {
		getOrganizationsRootRespBody, diagErr := getOrganizationsRoot(d, updateOrganizationsClient)
		if diagErr != nil {
			return diagErr
		}
		rootId := utils.PathSearch("roots|[0].id", getOrganizationsRootRespBody, "").(string)
		err = updateTags(d, updateOrganizationsClient, rootType, rootId, "root_tags")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceOrganizationsRead(ctx, d, meta)
}

func resourceOrganizationsDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteOrganizations: Delete Organizations
	var (
		deleteOrganizationsHttpUrl = "v1/organizations"
		deleteOrganizationsProduct = "organizations"
	)
	deleteOrganizationsClient, err := cfg.NewServiceClient(deleteOrganizationsProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	deleteOrganizationsPath := deleteOrganizationsClient.Endpoint + deleteOrganizationsHttpUrl

	deleteOrganizationsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = deleteOrganizationsClient.Request("DELETE", deleteOrganizationsPath, &deleteOrganizationsOpt)
	if err != nil {
		return diag.Errorf("error deleting Organizations: %s", err)
	}

	return nil
}

func addTags(client *golangsdk.ServiceClient, resourceType, resourceId string,
	tagList []tags.ResourceTag) error {
	var (
		addTagsToResourceHttpUrl = "v1/organizations/{resource_type}/{resource_id}/tags/create"
	)

	addTagsToResourcePath := client.Endpoint + addTagsToResourceHttpUrl
	addTagsToResourcePath = strings.ReplaceAll(addTagsToResourcePath, "{resource_type}", resourceType)
	addTagsToResourcePath = strings.ReplaceAll(addTagsToResourcePath, "{resource_id}", resourceId)

	addTagsToResourceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	addTagsToResourceOpt.JSONBody = utils.RemoveNil(buildTagsBodyParams(tagList))
	_, err := client.Request("POST", addTagsToResourcePath, &addTagsToResourceOpt)
	if err != nil {
		return fmt.Errorf("error creating Organizations tags: %s", err)
	}

	return nil
}

func deleteTags(client *golangsdk.ServiceClient, resourceType, resourceId string,
	tagList []tags.ResourceTag) error {
	var (
		addTagsToResourceHttpUrl = "v1/organizations/{resource_type}/{resource_id}/tags/delete"
	)

	addTagsToResourcePath := client.Endpoint + addTagsToResourceHttpUrl
	addTagsToResourcePath = strings.ReplaceAll(addTagsToResourcePath, "{resource_type}", resourceType)
	addTagsToResourcePath = strings.ReplaceAll(addTagsToResourcePath, "{resource_id}", resourceId)

	addTagsToResourceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	addTagsToResourceOpt.JSONBody = utils.RemoveNil(buildTagsBodyParams(tagList))
	_, err := client.Request("POST", addTagsToResourcePath, &addTagsToResourceOpt)
	if err != nil {
		return fmt.Errorf("error deleting Organizations tags: %s", err)
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
		return nil, fmt.Errorf("error get Organizations tags:: %s", err)
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
