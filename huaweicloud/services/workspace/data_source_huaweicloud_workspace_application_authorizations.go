package workspace

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

// @API Workspace GET /v1/{project_id}/app-center/apps/{app_id}/authorizations
func DataSourceApplicationAuthorizations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationAuthorizationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the application authorizations are located.`,
			},

			// Required parameters.
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the application.`,
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The username or user group name.`,
			},
			"target_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the target.`,
			},

			// Attributes.
			"authorizations": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        applicationAuthorizationSchema(),
				Description: `The list of application authorizations that matched filter parameters.`,
			},
		},
	}
}

func applicationAuthorizationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"account_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The account type.`,
			},
			"account": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The account information.`,
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain name.`,
			},
			"platform_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The platform type.`,
			},
		},
	}
}

func buildApplicationAuthorizationsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("target_type"); ok {
		res = fmt.Sprintf("%s&target_type=%v", res, v)
	}

	return res
}

func listApplicationAuthorizations(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/app-center/apps/{app_id}/authorizations?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	appId := d.Get("app_id").(string)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{app_id}", appId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildApplicationAuthorizationsQueryParams(d)

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
		authorizations := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, authorizations...)
		if len(authorizations) < limit {
			break
		}
		offset += len(authorizations)
	}

	return result, nil
}

func flattenApplicationAuthorizations(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"account":       utils.PathSearch("account", item, nil),
			"domain":        utils.PathSearch("domain", item, nil),
			"account_type":  utils.PathSearch("account_type", item, nil),
			"platform_type": utils.PathSearch("platform_type", item, nil),
		})
	}

	return result
}

func dataSourceApplicationAuthorizationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	authorizations, err := listApplicationAuthorizations(client, d)
	if err != nil {
		return diag.Errorf("error querying Workspace application authorizations: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("authorizations", flattenApplicationAuthorizations(authorizations)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
