package cdn

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var CacheSharingGroupNonUpdatableParams = []string{"name", "primary_domain"}

// @API CDN POST /v1.0/cdn/configuration/share-cache-groups
// @API CDN GET /v1.0/cdn/configuration/share-cache-groups
// @API CDN PUT /v1.0/cdn/configuration/share-cache-groups/{id}
// @API CDN DELETE /v1.0/cdn/configuration/share-cache-groups/{id}
func ResourceCacheSharingGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCacheSharingGroupCreate,
		ReadContext:   resourceCacheSharingGroupRead,
		UpdateContext: resourceCacheSharingGroupUpdate,
		DeleteContext: resourceCacheSharingGroupDelete,

		CustomizeDiff: config.FlexibleForceNew(CacheSharingGroupNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceCacheSharingGroupImportState,
		},

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the cache sharing group.`,
			},
			"primary_domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The primary domain name of the cache sharing group.`,
			},
			"share_cache_records": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The associated domain name.`,
						},
					},
				},
				Description: `The list of associated domain names of the cache sharing group.`,
			},

			// Attributes.
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the cache sharing group, in RFC3339 format.`,
			},

			// Internal parameter.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildShareCacheRecordsBodyParams(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"domain_name": utils.PathSearch("domain_name", item, nil),
		})
	}

	return result
}

func buildCacheSharingGroupCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"group_name":          d.Get("name").(string),
		"primary_domain":      d.Get("primary_domain").(string),
		"share_cache_records": buildShareCacheRecordsBodyParams(d.Get("share_cache_records").(*schema.Set).List()),
	}
}

func createCacheSharingGroup(client *golangsdk.ServiceClient, bodyParams map[string]interface{}) error {
	httpUrl := "v1.0/cdn/configuration/share-cache-groups"
	createPath := client.Endpoint + httpUrl

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(bodyParams),
		OkCodes:          []int{204},
	}

	_, err := client.Request("POST", createPath, &createOpt)
	return err
}

func listCacheSharingGroups(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v1.0/cdn/configuration/share-cache-groups?limit={limit}"
		limit   = 1000
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		groups := utils.PathSearch("share_cache_groups", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, groups...)
		if len(groups) < limit {
			break
		}
		offset += len(groups)
	}

	return result, nil
}

func GetCacheSharingGroupByName(client *golangsdk.ServiceClient, groupName string) (interface{}, error) {
	groups, err := listCacheSharingGroups(client)
	if err != nil {
		return nil, err
	}

	group := utils.PathSearch(fmt.Sprintf("[?group_name == '%s']|[0]", groupName), groups, nil)
	if group == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1.0/cdn/configuration/share-cache-groups",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the cache sharing group with name '%s' has been removed", groupName)),
			},
		}
	}
	return group, nil
}

func resourceCacheSharingGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	bodyParams := buildCacheSharingGroupCreateBodyParams(d)
	if err := createCacheSharingGroup(client, bodyParams); err != nil {
		return diag.Errorf("error creating CDN cache sharing group: %s", err)
	}

	groupName := d.Get("name").(string)
	group, err := GetCacheSharingGroupByName(client, groupName)
	if err != nil {
		return diag.Errorf("error querying CDN cache sharing groups: %s", err)
	}

	d.SetId(utils.PathSearch("id", group, "").(string))

	return resourceCacheSharingGroupRead(ctx, d, meta)
}

func GetCacheSharingGroupById(client *golangsdk.ServiceClient, groupId string) (interface{}, error) {
	groups, err := listCacheSharingGroups(client)
	if err != nil {
		return nil, err
	}

	group := utils.PathSearch(fmt.Sprintf("[?id == '%s']|[0]", groupId), groups, nil)
	if group == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1.0/cdn/configuration/share-cache-groups",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the cache sharing group with ID '%s' has been removed", groupId)),
			},
		}
	}
	return group, nil
}

func flattenShareCacheRecords(records []interface{}) []map[string]interface{} {
	if len(records) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(records))
	for _, record := range records {
		result = append(result, map[string]interface{}{
			"domain_name": utils.PathSearch("domain_name", record, nil),
		})
	}

	return result
}

func resourceCacheSharingGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		groupId = d.Id()
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	group, err := GetCacheSharingGroupById(client, groupId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error getting cache sharing group (%s)", groupId))
	}

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("group_name", group, nil)),
		d.Set("primary_domain", utils.PathSearch("primary_domain", group, nil)),
		d.Set("share_cache_records", flattenShareCacheRecords(utils.PathSearch("share_cache_records", group,
			make([]interface{}, 0)).([]interface{}))),
		d.Set("create_time", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", group, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCacheSharingGroupUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"share_cache_records": buildShareCacheRecordsBodyParams(d.Get("share_cache_records").(*schema.Set).List()),
	}
}

func updateCacheSharingGroup(client *golangsdk.ServiceClient, groupId string, bodyParams map[string]interface{}) error {
	httpUrl := "v1.0/cdn/configuration/share-cache-groups/{id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{id}", groupId)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(bodyParams),
		OkCodes:          []int{204},
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func resourceCacheSharingGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		groupId = d.Id()
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	bodyParams := buildCacheSharingGroupUpdateBodyParams(d)
	if err := updateCacheSharingGroup(client, groupId, bodyParams); err != nil {
		return diag.Errorf("error updating cache sharing group (%s): %s", groupId, err)
	}

	return resourceCacheSharingGroupRead(ctx, d, meta)
}

func deleteCacheSharingGroup(client *golangsdk.ServiceClient, groupId string) error {
	httpUrl := "v1.0/cdn/configuration/share-cache-groups/{id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{id}", groupId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	return err
}

func resourceCacheSharingGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                     = meta.(*config.Config)
		groupId                 = d.Id()
		cacheGroupNotFoundCodes = []string{
			"CDN.0001", // The share cache group does not exist.
		}
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	// Before deleting the share cache group, clean up the share cache records.
	if d.Get("share_cache_records").(*schema.Set).Len() > 0 {
		bodyParams := map[string]interface{}{
			"share_cache_records": make([]interface{}, 0),
		}
		if err := updateCacheSharingGroup(client, groupId, bodyParams); err != nil {
			return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error.error_code", cacheGroupNotFoundCodes...),
				fmt.Sprintf("error cleaning up share cache records for cache sharing group (%s)", groupId))
		}
	}

	err = deleteCacheSharingGroup(client, groupId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error.error_code", cacheGroupNotFoundCodes...),
			fmt.Sprintf("error deleting cache sharing group (%s)", groupId))
	}

	return nil
}

func resourceCacheSharingGroupImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()

	if utils.IsUUID(importedId) {
		d.SetId(importedId)
		return []*schema.ResourceData{d}, nil
	}

	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return nil, fmt.Errorf("error creating CDN client: %s", err)
	}

	group, err := GetCacheSharingGroupByName(client, importedId)
	if err != nil {
		return nil, err
	}
	d.SetId(utils.PathSearch("id", group, "").(string))

	return []*schema.ResourceData{d}, nil
}
