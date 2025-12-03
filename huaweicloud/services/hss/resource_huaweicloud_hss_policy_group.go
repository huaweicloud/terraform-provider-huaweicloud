package hss

import (
	"context"
	"errors"
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

// @API HSS PUT /v5/{project_id}/policy/group
// @API HSS DELETE /v5/{project_id}/policy/group
// @API HSS POST /v5/{project_id}/policy/group
// @API HSS GET /v5/{project_id}/policy/groups
func ResourcePolicyGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyGroupCreate,
		ReadContext:   resourcePolicyGroupRead,
		UpdateContext: resourcePolicyGroupUpdate,
		DeleteContext: resourcePolicyGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePolicyGroupImportState,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"group_id",
			"name",
			"description",
			"enterprise_project_id",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			// Field `group_id` not exist in API response.
			// Field `group_id` cannot be updated.
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Field `name` cannot be updated.
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Field `description` cannot be updated.
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Field `enterprise_project_id` not exist in API response.
			// Field `enterprise_project_id` cannot be updated.
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protect_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"host_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"default_group": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"deletable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"support_os": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"support_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateGroupQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}
	return ""
}

func buildCreateGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"group_id":    d.Get("group_id"),
		"name":        d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func buildQueryGroupByNameQueryParams(epsId string, offset int) string {
	rst := ""

	if offset > 0 {
		rst += fmt.Sprintf("&offset=%d", offset)
	}
	if epsId != "" {
		rst += fmt.Sprintf("&enterprise_project_id=%s", epsId)
	}

	if rst != "" {
		rst = "?" + rst[1:]
	}

	return rst
}

func queryPolicyGroupByName(client *golangsdk.ServiceClient, name, epsId string) (interface{}, error) {
	requestPath := client.Endpoint + "v5/{project_id}/policy/groups"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var (
		allPolicyGroups []interface{}
		offset          int
	)

	for {
		requestPathWithOffset := requestPath + buildQueryGroupByNameQueryParams(epsId, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		dataList := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataList) == 0 {
			break
		}

		allPolicyGroups = append(allPolicyGroups, dataList...)
		offset += len(dataList)
	}

	policyGroup := utils.PathSearch(fmt.Sprintf("[?group_name == '%s']|[0]", name), allPolicyGroups, nil)
	if policyGroup == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return policyGroup, nil
}

func buildUpdateGroupQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}
	return ""
}

func updatePolicyGroup(client *golangsdk.ServiceClient, d *schema.ResourceData, epsId string) error {
	requestPath := client.Endpoint + "v5/{project_id}/policy/group"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildUpdateGroupQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"group_id":     d.Id(),
			"protect_mode": d.Get("protect_mode"),
		},
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	return err
}

func resourcePolicyGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/policy/group"
		name    = d.Get("name").(string)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath += buildCreateGroupQueryParams(epsId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateGroupBodyParams(d)),
	}

	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error adding HSS policy group: %s", err)
	}

	policyGroup, err := queryPolicyGroupByName(client, name, epsId)
	if err != nil {
		return diag.Errorf("error querying HSS policy group in create operation: %s", err)
	}

	groupId := utils.PathSearch("group_id", policyGroup, "").(string)
	if groupId == "" {
		return diag.Errorf("error adding HSS policy group: group ID is not found in API response")
	}

	d.SetId(groupId)

	if _, ok := d.GetOk("protect_mode"); ok {
		if err := updatePolicyGroup(client, d, epsId); err != nil {
			return diag.Errorf("error updating HSS policy group in create operation: %s", err)
		}
	}

	return resourcePolicyGroupRead(ctx, d, meta)
}

func buildQueryGroupByIdQueryParams(policyGroupId, epsId string) string {
	rst := fmt.Sprintf("?group_id=%s", policyGroupId)

	if epsId != "" {
		rst += fmt.Sprintf("&enterprise_project_id=%s", epsId)
	}
	return rst
}

func QueryPolicyGroupById(client *golangsdk.ServiceClient, policyGroupId, epsId string) (interface{}, error) {
	requestPath := client.Endpoint + "v5/{project_id}/policy/groups"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildQueryGroupByIdQueryParams(policyGroupId, epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	policyGroup := utils.PathSearch("data_list|[0]", respBody, nil)
	if policyGroup == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return policyGroup, nil
}

func resourcePolicyGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	policyGroup, err := QueryPolicyGroupById(client, d.Id(), epsId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying HSS policy group")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("group_name", policyGroup, nil)),
		d.Set("description", utils.PathSearch("description", policyGroup, nil)),
		d.Set("protect_mode", utils.PathSearch("protect_mode", policyGroup, nil)),
		d.Set("host_num", utils.PathSearch("host_num", policyGroup, nil)),
		d.Set("default_group", utils.PathSearch("default_group", policyGroup, nil)),
		d.Set("deletable", utils.PathSearch("deletable", policyGroup, nil)),
		d.Set("support_os", utils.PathSearch("support_os", policyGroup, nil)),
		d.Set("support_version", utils.PathSearch("support_version", policyGroup, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePolicyGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	if err := updatePolicyGroup(client, d, epsId); err != nil {
		return diag.Errorf("error updating HSS policy group in update operation: %s", err)
	}

	return resourcePolicyGroupRead(ctx, d, meta)
}

func buildDeleteGroupQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}
	return ""
}

func resourcePolicyGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/policy/group"
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildDeleteGroupQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"id_list": []string{d.Id()},
		},
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting HSS policy group: %s", err)
	}

	return nil
}

func resourcePolicyGroupImportState(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format of import ID, must be <enterprise_project_id>/<id>")
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("enterprise_project_id", parts[0])
}
