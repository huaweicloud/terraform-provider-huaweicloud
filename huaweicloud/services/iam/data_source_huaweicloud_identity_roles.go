package iam

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v3/roles
func DataSourceV3Roles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV3RolesRead,

		Schema: map[string]*schema.Schema{
			// Optional parameters.
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The display name of the role to be queried.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the role to be queried.`,
			},
			"catalog": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The catalog of the role to be queried.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The display mode of the role to be queried.`,
			},
			"permission_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   `The permission type of the role to be queried.`,
				ConflictsWith: []string{"domain_id"},
			},
			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The domain ID to be queried.`,
			},

			// Attributes.
			"roles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the roles.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the role.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the role.`,
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The display name of the role.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the role.`,
						},
						"description_cn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the role in Chinese.`,
						},
						"catalog": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The catalog of the role.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The display mode of the role.`,
						},
						"flag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The flag of the role.`,
						},
						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain ID of the role.`,
						},
						"policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The content of the role, in JSON format.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the role, in RFC3339 format.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The last update time of the role, in RFC3339 format.`,
						},
						"links": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The links of the role.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"self": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The self link of the role.`,
									},
									"previous": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The previous link of the role.`,
									},
									"next": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The next link of the role.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildV3RolesQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("display_name"); ok {
		res = fmt.Sprintf("%s&display_name=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("catalog"); ok {
		res = fmt.Sprintf("%s&catalog=%v", res, v)
	}
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}
	if v, ok := d.GetOk("permission_type"); ok {
		res = fmt.Sprintf("%s&permission_type=%v", res, v)
	}
	if v, ok := d.GetOk("domain_id"); ok {
		res = fmt.Sprintf("%s&domain_id=%v", res, v)
	}

	return res
}

func flattenV3RoleLinks(links map[string]interface{}) []map[string]interface{} {
	if len(links) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"self":     utils.PathSearch("self", links, nil),
			"previous": utils.PathSearch("previous", links, nil),
			"next":     utils.PathSearch("next", links, nil),
		},
	}
}

func flattenV3Roles(roles []interface{}) []map[string]interface{} {
	if len(roles) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(roles))
	for _, role := range roles {
		createdAt, _ := strconv.ParseInt(utils.PathSearch("created_time", role, "").(string), 10, 64)
		updatedAt, _ := strconv.ParseInt(utils.PathSearch("updated_time", role, "").(string), 10, 64)

		result = append(result, map[string]interface{}{
			"id":             utils.PathSearch("id", role, nil),
			"name":           utils.PathSearch("name", role, nil),
			"display_name":   utils.PathSearch("display_name", role, nil),
			"description":    utils.PathSearch("description", role, nil),
			"description_cn": utils.PathSearch("description_cn", role, nil),
			"catalog":        utils.PathSearch("catalog", role, nil),
			"type":           utils.PathSearch("type", role, nil),
			"flag":           utils.PathSearch("flag", role, nil),
			"domain_id":      utils.PathSearch("domain_id", role, nil),
			"policy":         utils.JsonToString(utils.PathSearch("policy", role, nil)),
			"created_at":     utils.FormatTimeStampRFC3339(createdAt/1000, false),
			"updated_at":     utils.FormatTimeStampRFC3339(updatedAt/1000, false),
			"links": flattenV3RoleLinks(
				utils.PathSearch("links", role, make(map[string]interface{})).(map[string]interface{})),
		})
	}

	return result
}

func listV3Roles(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v3/roles?per_page={per_page}"
		result  = make([]interface{}, 0)
		page    = 1
		perPage = 300
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{per_page}", strconv.Itoa(perPage))
	listPath += buildV3RolesQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithPage := fmt.Sprintf("%s&page=%s", listPath, strconv.Itoa(page))

		resp, err := client.Request("GET", listPathWithPage, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		roles := utils.PathSearch("roles", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, roles...)
		if len(roles) < perPage {
			break
		}
		page++
	}

	return result, nil
}

func dataSourceV3RolesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	roles, err := listV3Roles(client, d)
	if err != nil {
		return diag.Errorf("error querying roles: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("roles", flattenV3Roles(roles)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
