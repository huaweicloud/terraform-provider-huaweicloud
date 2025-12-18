package waf

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

// @API WAF POST /v1/{project_id}/premium-waf/agency
// @API WAF DELETE /v1/{project_id}/premium-waf/agency
// @API WAF GET /v1/{project_id}/premium-waf/agency
func ResourceWafDedicatedAgency() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafDedicatedAgencyCreate,
		ReadContext:   resourceWafDedicatedAgencyRead,
		UpdateContext: resourceWafDedicatedAgencyUpdate,
		DeleteContext: resourceWafDedicatedAgencyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `Specifies the region in which to create the resource.`,
			},
			// Field `role_name_list` is optional in API, make it required here.
			// Field `role_name_list` is not exist in detail API response.
			"role_name_list": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of dedicated engine agent policy names to create.`,
			},
			// Field `purged` only used in delete operation and not exist in API response.
			"purged": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to delete delegates synchronously.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The agency name.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version.`,
			},
			"duration": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The agent existence time period.`,
			},
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain ID.`,
			},
			"is_valid": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the agency is legal.`,
			},
			// The `role_list` field does not exist in the API documentation, but it appears in the response body during
			// actual testing. Use this field to determine if a `404` error has occurred.
			"role_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description.`,
						},
						"catalog": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The catalog.`,
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The role ID.`,
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The display name.`,
						},
						"is_granted": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether it is granted.`,
						},
					},
				},
			},
		},
	}
}

func configDedicatedAgency(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/premium-waf/agency"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: map[string]interface{}{
			"role_name_list": d.Get("role_name_list"),
		},
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceWafDedicatedAgencyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	respBody, err := configDedicatedAgency(client, d)
	if err != nil {
		return diag.Errorf("error creating WAF dedicated agency: %s", err)
	}

	id := utils.PathSearch("id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating WAF dedicated agency: ID is not found in API response")
	}
	d.SetId(id)

	return resourceWafDedicatedAgencyRead(ctx, d, meta)
}

func QueryDedicatedAgency(client *golangsdk.ServiceClient) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/premium-waf/agency"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	roleList := utils.PathSearch("role_list", respBody, make([]interface{}, 0)).([]interface{})
	if len(roleList) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return respBody, nil
}

func resourceWafDedicatedAgencyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	respBody, err := QueryDedicatedAgency(client)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving WAF dedicated agency")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("version", utils.PathSearch("version", respBody, nil)),
		d.Set("duration", utils.PathSearch("duration", respBody, nil)),
		d.Set("domain_id", utils.PathSearch("domain_id", respBody, nil)),
		d.Set("is_valid", utils.PathSearch("is_valid", respBody, nil)),
		d.Set("role_list", flattenRoleListAttribute(utils.PathSearch("role_list", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRoleListAttribute(respArray []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"description":  utils.PathSearch("description", v, nil),
			"catalog":      utils.PathSearch("catalog", v, nil),
			"id":           utils.PathSearch("id", v, nil),
			"display_name": utils.PathSearch("display_name", v, nil),
			"is_granted":   utils.PathSearch("is_granted", v, nil),
		})
	}

	return rst
}

func resourceWafDedicatedAgencyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	if d.HasChange("role_name_list") {
		_, err = configDedicatedAgency(client, d)
		if err != nil {
			return diag.Errorf("error updating WAF dedicated agency: %s", err)
		}
	}

	return resourceWafDedicatedAgencyRead(ctx, d, meta)
}

func buildDeleteDedicatedAgencyQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?role_id_list=%s&purged=%v", d.Id(), d.Get("purged").(bool))
}

func resourceWafDedicatedAgencyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/premium-waf/agency"
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildDeleteDedicatedAgencyQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting WAF dedicated agency: %s", err)
	}

	return nil
}
