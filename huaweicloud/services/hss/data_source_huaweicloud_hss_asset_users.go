package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS GET /v5/{project_id}/asset/users
func DataSourceAssetUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssetUsersRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the resource.",
			},
			"host_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the host ID.",
			},
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the account name.",
			},
			"host_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the host name.",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the private IP of the server.",
			},
			"login_permission": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether the user has the login permission.",
			},
			"root_permission": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether the user has root permissions.",
			},
			"user_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the server user group.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the enterprise project ID.",
			},
			"category": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the type. The default value is host.",
			},
			"part_match": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether to use fuzzy matching.",
			},
			"data_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        userSchema(),
				Description: "The list of account information.",
			},
		},
	}
}

func userSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"agent_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The agent ID.",
			},
			"host_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The host ID.",
			},
			"host_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The host name.",
			},
			"host_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The host IP.",
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user name.",
			},
			"login_permission": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the user has the login permission.",
			},
			"root_permission": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the user has root permissions.",
			},
			"user_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user group name.",
			},
			"user_home_dir": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user home directory.",
			},
			"shell": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user startup shell.",
			},
			"recent_scan_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The latest scan time.",
			},
			"container_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The container ID.",
			},
			"container_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The container name.",
			},
		},
	}
	return &sc
}

func buildAssetUsersQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	epsId := cfg.GetEnterpriseProjectID(d)
	queryParams := "?limit=10"
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("host_id"); ok {
		queryParams = fmt.Sprintf("%s&host_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("user_name"); ok {
		queryParams = fmt.Sprintf("%s&user_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("private_ip"); ok {
		queryParams = fmt.Sprintf("%s&private_ip=%v", queryParams, v)
	}
	if v, ok := d.GetOk("login_permission"); ok {
		queryParams = fmt.Sprintf("%s&login_permission=%v", queryParams, v)
	}
	if v, ok := d.GetOk("root_permission"); ok {
		queryParams = fmt.Sprintf("%s&root_permission=%v", queryParams, v)
	}
	if v, ok := d.GetOk("user_group"); ok {
		queryParams = fmt.Sprintf("%s&user_group=%v", queryParams, v)
	}
	if v, ok := d.GetOk("category"); ok {
		queryParams = fmt.Sprintf("%s&category=%v", queryParams, v)
	}
	if v, ok := d.GetOk("part_match"); ok {
		queryParams = fmt.Sprintf("%s&part_match=%v", queryParams, v)
	}
	return queryParams
}

func flattenAssetUsers(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"agent_id":         utils.PathSearch("agent_id", v, nil),
			"host_id":          utils.PathSearch("host_id", v, nil),
			"host_name":        utils.PathSearch("host_name", v, nil),
			"host_ip":          utils.PathSearch("host_ip", v, nil),
			"user_name":        utils.PathSearch("user_name", v, nil),
			"login_permission": utils.PathSearch("login_permission", v, nil),
			"root_permission":  utils.PathSearch("root_permission", v, nil),
			"user_group_name":  utils.PathSearch("user_group_name", v, nil),
			"user_home_dir":    utils.PathSearch("user_home_dir", v, nil),
			"shell":            utils.PathSearch("shell", v, nil),
			"recent_scan_time": utils.PathSearch("recent_scan_time", v, nil),
			"container_id":     utils.PathSearch("container_id", v, nil),
			"container_name":   utils.PathSearch("container_name", v, nil),
		})
	}
	return rst
}

func dataSourceAssetUsersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/asset/users"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAssetUsersQueryParams(d, cfg)
	allUsers := make([]interface{}, 0)
	offset := 0

	listUsersOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", currentPath, &listUsersOpt)
		if err != nil {
			return diag.Errorf("error retrieving HSS asset users: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		usersResp := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(usersResp) == 0 {
			break
		}
		allUsers = append(allUsers, usersResp...)
		offset += len(usersResp)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("data_list", flattenAssetUsers(allUsers)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
