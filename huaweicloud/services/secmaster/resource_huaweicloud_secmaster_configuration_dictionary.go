package secmaster

import (
	"context"
	"fmt"
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

var nonUpdatableParamsConfigurationDictionary = []string{
	"dict_id",
	"dict_key",
	"language",
	"version",
	"scope",
	"is_built_in",
}

// @API SecMaster POST /v1/{project_id}/configurations/dictionaries
// @API SecMaster GET /v1/{project_id}/configurations/dictionaries
// @API SecMaster PUT /v1/{project_id}/configurations/dictionaries
// @API SecMaster DELETE /v1/{project_id}/configurations/dictionaries
func ResourceConfigurationDictionary() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConfigurationDictionaryCreate,
		ReadContext:   resourceConfigurationDictionaryRead,
		UpdateContext: resourceConfigurationDictionaryUpdate,
		DeleteContext: resourceConfigurationDictionaryDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParamsConfigurationDictionary),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"dict_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dict_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dict_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dict_val": {
				Type:     schema.TypeString,
				Required: true,
			},
			"language": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dict_pkey": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dict_pcode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"extend_field": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"is_built_in": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publish_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateConfigurationDictionaryBodyParams(d *schema.ResourceData) map[string]interface{} {
	dictItem := map[string]interface{}{
		"dict_id":   d.Get("dict_id"),
		"dict_key":  d.Get("dict_key"),
		"dict_code": d.Get("dict_code"),
		"dict_val":  d.Get("dict_val"),
		"language":  d.Get("language"),
	}

	if v, ok := d.GetOk("version"); ok {
		dictItem["version"] = v
	}
	if v, ok := d.GetOk("dict_pkey"); ok {
		dictItem["dict_pkey"] = v
	}
	if v, ok := d.GetOk("dict_pcode"); ok {
		dictItem["dict_pcode"] = v
	}
	if v, ok := d.GetOk("scope"); ok {
		dictItem["scope"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		dictItem["description"] = v
	}
	if _, ok := d.GetOk("extend_field"); ok {
		dictItem["extend_field"] = utils.ValueIgnoreEmpty(d.Get("extend_field"))
	}

	bodyParams := map[string]interface{}{
		"dict_list": []interface{}{dictItem},
	}

	if v, ok := d.GetOk("is_built_in"); ok {
		bodyParams["is_built_in"] = v
	}

	return bodyParams
}

func resourceConfigurationDictionaryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		createHttpUrl = "v1/{project_id}/configurations/dictionaries"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         buildCreateConfigurationDictionaryBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster configuration dictionary: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	// The create API returns success_list and failed_list, get the id from success_list
	dictUUID := utils.PathSearch("success_list|[0].id", respBody, "").(string)
	if dictUUID == "" {
		return diag.Errorf("error creating SecMaster configuration dictionary: unable to find dictionary ID")
	}

	d.SetId(dictUUID)

	return resourceConfigurationDictionaryRead(ctx, d, meta)
}

func GetConfigurationDictionaryInfo(client *golangsdk.ServiceClient, dictUUID string) (interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/configurations/dictionaries"
		limit   = 100
		offset  = 0
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpts := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
	}

	for {
		getPathWithParams := fmt.Sprintf("%s?limit=%d&offset=%d", getPath, limit, offset)
		resp, err := client.Request("GET", getPathWithParams, &getOpts)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		dictItem := utils.PathSearch(fmt.Sprintf("success_list[?id=='%s']|[0]", dictUUID), respBody, nil)
		if dictItem != nil {
			return dictItem, nil
		}

		successList := utils.PathSearch("success_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(successList) < limit {
			break
		}
		offset += limit
	}

	return nil, golangsdk.ErrDefault404{}
}

func resourceConfigurationDictionaryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	dictInfo, err := GetConfigurationDictionaryInfo(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SecMaster configuration dictionary")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("dict_id", utils.PathSearch("dict_id", dictInfo, nil)),
		d.Set("dict_key", utils.PathSearch("dict_key", dictInfo, nil)),
		d.Set("dict_code", utils.PathSearch("dict_code", dictInfo, nil)),
		d.Set("dict_val", utils.PathSearch("dict_val", dictInfo, nil)),
		d.Set("language", utils.PathSearch("language", dictInfo, nil)),
		d.Set("version", utils.PathSearch("version", dictInfo, nil)),
		d.Set("dict_pkey", utils.PathSearch("dict_pkey", dictInfo, nil)),
		d.Set("dict_pcode", utils.PathSearch("dict_pcode", dictInfo, nil)),
		d.Set("scope", utils.PathSearch("scope", dictInfo, nil)),
		d.Set("description", utils.PathSearch("description", dictInfo, nil)),
		d.Set("create_time", utils.PathSearch("create_time", dictInfo, nil)),
		d.Set("update_time", utils.PathSearch("update_time", dictInfo, nil)),
		d.Set("publish_time", utils.PathSearch("publish_time", dictInfo, nil)),
		d.Set("extend_field", utils.PathSearch("extension_field", dictInfo, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateConfigurationDictionaryBodyParams(d *schema.ResourceData) map[string]interface{} {
	dictItem := map[string]interface{}{
		"dict_id":   d.Get("dict_id"),
		"dict_key":  d.Get("dict_key"),
		"dict_code": d.Get("dict_code"),
		"dict_val":  d.Get("dict_val"),
		"language":  d.Get("language"),
	}

	if v, ok := d.GetOk("version"); ok {
		dictItem["version"] = v
	}
	if v, ok := d.GetOk("dict_pkey"); ok {
		dictItem["dict_pkey"] = v
	}
	if v, ok := d.GetOk("dict_pcode"); ok {
		dictItem["dict_pcode"] = v
	}
	if v, ok := d.GetOk("scope"); ok {
		dictItem["scope"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		dictItem["description"] = v
	}
	if _, ok := d.GetOk("extend_field"); ok {
		dictItem["extend_field"] = utils.ValueIgnoreEmpty(d.Get("extend_field"))
	}

	bodyParams := map[string]interface{}{
		"dict_list": []interface{}{dictItem},
	}

	if v, ok := d.GetOk("is_built_in"); ok {
		bodyParams["is_built_in"] = v
	}

	return bodyParams
}

func resourceConfigurationDictionaryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		updateHttpUrl = "v1/{project_id}/configurations/dictionaries"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         buildUpdateConfigurationDictionaryBodyParams(d),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating SecMaster configuration dictionary: %s", err)
	}

	return resourceConfigurationDictionaryRead(ctx, d, meta)
}

func resourceConfigurationDictionaryDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		deleteHttpUrl = "v1/{project_id}/configurations/dictionaries"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)

	deleteDictItem := map[string]interface{}{
		"dict_id":  d.Get("dict_id"),
		"dict_key": d.Get("dict_key"),
		"language": d.Get("language"),
	}

	deleteBodyParams := map[string]interface{}{
		"dict_list": []interface{}{deleteDictItem},
	}

	if v, ok := d.GetOk("is_built_in"); ok {
		deleteBodyParams["is_built_in"] = v
	}

	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         deleteBodyParams,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting SecMaster configuration dictionary: %s", err)
	}

	return nil
}
