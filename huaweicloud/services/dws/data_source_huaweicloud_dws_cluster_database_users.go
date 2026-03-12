package dws

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

// @API DWS GET /v1/{project_id}/clusters/{cluster_id}/db-manager/users
func DataSourceClusterDatabaseUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterDatabaseUsersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the database users are located.`,
			},

			// Required parameters.
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the cluster to be queried.`,
			},

			// Optional parameters.
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The object type to be queried.`,
			},
			"user_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The user type to be queried.`,
			},

			// Attributes.
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        clusterDatabaseUserSchema(),
				Description: `The list of the database users that matched filter parameters.`,
			},
		},
	}
}

func clusterDatabaseUserSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the database user or role.`,
			},
			"login": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the database user can login.`,
			},
			"user_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the database user.`,
			},
		},
	}
}

func buildClusterDatabaseUsersQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}
	if v, ok := d.GetOk("user_type"); ok {
		res = fmt.Sprintf("%s&user_type=%v", res, v)
	}

	return res
}

func listClusterDatabaseUsers(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl   = "v1/{project_id}/clusters/{cluster_id}/db-manager/users?limit={limit}"
		clusterId = d.Get("cluster_id").(string)
		limit     = 1000
		offset    = 0
		result    = make([]interface{}, 0)
	)

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{cluster_id}", clusterId)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))
	listPathWithLimit += buildClusterDatabaseUsersQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPathWithLimit, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		users := utils.PathSearch("users", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, users...)
		if len(users) < limit {
			break
		}
		offset += len(users)
	}

	return result, nil
}

func flattenClusterDatabaseUsers(users []interface{}) []map[string]interface{} {
	if len(users) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(users))
	for _, user := range users {
		result = append(result, map[string]interface{}{
			"name":      utils.PathSearch("name", user, nil),
			"login":     utils.PathSearch("login", user, nil),
			"user_type": utils.PathSearch("user_type", user, nil),
		})
	}

	return result
}

func dataSourceClusterDatabaseUsersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	users, err := listClusterDatabaseUsers(client, d)
	if err != nil {
		return diag.Errorf("error querying cluster database users: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("users", flattenClusterDatabaseUsers(users)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
