package das

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DAS GET /v3/{project_id}/instances/{instance_id}/db-users
func DataSourceDatabaseUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatabaseUsersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the database users are located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the instance to which the database user belongs.`,
			},

			// Optional parameters.
			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the database user.",
			},
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the database user.",
			},

			// Attributes.
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        databaseUsersElem(),
				Description: `The list of users that matched filter parameters.`,
			},
		},
	}
}

func databaseUsersElem() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the database user, in UUID format.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the database user.",
			},
		},
	}
	return &sc
}

func buildDatabaseUsersQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("user_id"); ok {
		res = fmt.Sprintf("%s&db_user_id=%v", res, v)
	}
	if v, ok := d.GetOk("user_name"); ok {
		res = fmt.Sprintf("%s&db_username=%v", res, v)
	}

	return res
}

func listDatabaseUsers(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/db-users?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildDatabaseUsersQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		users := utils.PathSearch("db_users", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, users...)
		if len(users) < limit {
			break
		}
		offset += len(users)
	}

	return result, nil
}

func flattenDatabaseUsers(users []interface{}) []map[string]interface{} {
	if len(users) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(users))
	for _, user := range users {
		result = append(result, map[string]interface{}{
			"id":   utils.PathSearch("db_user_id", user, nil),
			"name": utils.PathSearch("db_username", user, nil),
		})
	}
	return result
}

func dataSourceDatabaseUsersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	users, err := listDatabaseUsers(client, d)
	if err != nil {
		return diag.Errorf("error querying DAS Database users: %s", err)
	}

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("users", flattenDatabaseUsers(users)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
