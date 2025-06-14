package dew

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DEW POST /v1.0/{project_id}/kms/resource_instances/action
func DataSourceKmsCustomKeysByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKmsCustomKeysByTagsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the resource.",
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the operation type.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        tagsSchema(),
				Description: "Specifies the tag list.",
			},
			"matches": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        matchesSchema(),
				Description: "Specifies the field to be matched.",
			},
			"sequence": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the sequence number of the request message.",
			},
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        resourcesSchema(),
				Description: "The list of key resources.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of records.",
			},
		},
	}
}

func tagsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the tag key.",
			},
			"values": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifies the tag value set.",
			},
		},
	}
}

func matchesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the field to be matched.",
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the value to be matched.",
			},
		},
	}
}

func resourcesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource ID.",
			},
			"resource_detail": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        resourceDetailSchema(),
				Description: "The key details.",
			},
			"resource_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource name.",
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        resourceTagsSchema(),
				Description: "The tag list.",
			},
		},
	}
}

func resourceDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CMK ID.",
			},
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user domain ID.",
			},
			"key_alias": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The key alias.",
			},
			"realm": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The key realm.",
			},
			"key_spec": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The key generation algorithm.",
			},
			"key_usage": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CMK usage.",
			},
			"key_description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The key description.",
			},
			"creation_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the key was created.",
			},
			"scheduled_deletion_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the key was scheduled to be deleted.",
			},
			"key_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The key status.",
			},
			"default_key_flag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The master key identifier.",
			},
			"key_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The key type.",
			},
			"expiration_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the key material expires.",
			},
			"origin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The key source.",
			},
			"key_rotation_enabled": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The key rotation status.",
			},
			"sys_enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The enterprise project ID.",
			},
			"keystore_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The keystore ID.",
			},
			"key_label": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The key label in the encryption machine.",
			},
			"partition_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The partition type.",
			},
		},
	}
}

func resourceTagsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The tag key.",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The tag value.",
			},
		},
	}
}

func dataSourceKmsCustomKeysByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		mErr   *multierror.Error
	)

	client, err := cfg.NewServiceClient("kms", region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	requestPath := client.Endpoint + "v1.0/{project_id}/kms/resource_instances/action"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	allResources := make([]interface{}, 0)
	allCount := 0
	offset := 0

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	bodyParams := buildKmsCustomKeysByTagsParams(d)
	for {
		if d.Get("action").(string) == "filter" {
			bodyParams["limit"] = "1000"
			bodyParams["offset"] = offset
		}
		listOpt.JSONBody = utils.RemoveNil(bodyParams)
		resp, err := client.Request("POST", requestPath, &listOpt)

		if err != nil {
			return diag.Errorf("error retrieving KMS custom keys by tags: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		resources := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})

		// When the action is count, resources is null, but total_count is not null.
		if d.Get("action").(string) == "count" {
			allCount = int(utils.PathSearch("total_count", respBody, float64(0)).(float64))
			break
		}

		if len(resources) == 0 {
			break
		}
		allResources = append(allResources, resources...)
		allCount += len(resources)
		offset += len(resources)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("resources", flattenKmsCustomKeysResources(allResources)),
		d.Set("total_count", allCount),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildKmsCustomKeysByTagsParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"action": d.Get("action"),
	}

	if v, ok := d.GetOk("tags"); ok {
		params["tags"] = buildKmsTags(v.([]interface{}))
	}

	if v, ok := d.GetOk("matches"); ok {
		params["matches"] = buildKmsMatches(v.([]interface{}))
	}

	if v, ok := d.GetOk("sequence"); ok {
		params["sequence"] = v.(string)
	}

	return params
}

func buildKmsTags(tags []interface{}) []map[string]interface{} {
	if len(tags) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(tags))

	for i, tag := range tags {
		tagMap := tag.(map[string]interface{})
		result[i] = map[string]interface{}{
			"key":    tagMap["key"],
			"values": tagMap["values"],
		}
	}

	return result
}

func buildKmsMatches(matches []interface{}) []map[string]interface{} {
	if len(matches) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(matches))

	for i, match := range matches {
		matchMap := match.(map[string]interface{})
		result[i] = map[string]interface{}{
			"key":   matchMap["key"],
			"value": matchMap["value"],
		}
	}

	return result
}

func flattenKmsCustomKeysResources(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))

	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"resource_id":     utils.PathSearch("resource_id", v, nil),
			"resource_name":   utils.PathSearch("resource_name", v, nil),
			"resource_detail": flattenKmsResourceDetail(utils.PathSearch("resource_detail", v, nil)),
			"tags":            flattenKmsResourceTags(utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})),
		})
	}
	return rst
}

func flattenKmsResourceDetail(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	resourceDetail := map[string]interface{}{
		"key_id":                    utils.PathSearch("key_id", resp, nil),
		"domain_id":                 utils.PathSearch("domain_id", resp, nil),
		"key_alias":                 utils.PathSearch("key_alias", resp, nil),
		"realm":                     utils.PathSearch("realm", resp, nil),
		"key_spec":                  utils.PathSearch("key_spec", resp, nil),
		"key_usage":                 utils.PathSearch("key_usage", resp, nil),
		"key_description":           utils.PathSearch("key_description", resp, nil),
		"creation_date":             utils.PathSearch("creation_date", resp, nil),
		"scheduled_deletion_date":   utils.PathSearch("scheduled_deletion_date", resp, nil),
		"key_state":                 utils.PathSearch("key_state", resp, nil),
		"default_key_flag":          utils.PathSearch("default_key_flag", resp, nil),
		"key_type":                  utils.PathSearch("key_type", resp, nil),
		"expiration_time":           utils.PathSearch("expiration_time", resp, nil),
		"origin":                    utils.PathSearch("origin", resp, nil),
		"key_rotation_enabled":      utils.PathSearch("key_rotation_enabled", resp, nil),
		"sys_enterprise_project_id": utils.PathSearch("sys_enterprise_project_id", resp, nil),
		"keystore_id":               utils.PathSearch("keystore_id", resp, nil),
		"key_label":                 utils.PathSearch("key_label", resp, nil),
		"partition_type":            utils.PathSearch("partition_type", resp, nil),
	}
	return []interface{}{resourceDetail}
}

func flattenKmsResourceTags(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))

	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}
	return rst
}
