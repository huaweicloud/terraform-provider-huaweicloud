package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS GET /v3/{project_id}/instances/{instance_id}/db_user/detail
func DataSourceRdsMysqlAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsMysqlAccountsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the RDS instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the username of the DB account.`,
			},
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the IP address that is allowed to access your DB instance.`,
			},
			"users": {
				Type:        schema.TypeList,
				Elem:        accountsSchema(),
				Computed:    true,
				Description: `Indicates the list of users.`,
			},
		},
	}
}

func accountsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the username of the DB account.`,
			},
			"hosts": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `Indicates the IP addresses that are allowed to access your DB instance.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates remarks of the database account.`,
			},
		},
	}
	return &sc
}

func dataSourceRdsMysqlAccountsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		listMysqlAccountsHttpUrl = "v3/{project_id}/instances/{instance_id}/db_user/detail?page=1&limit=100"
		listMysqAccountsProduct  = "rds"
	)
	listMysqlAccountsClient, err := cfg.NewServiceClient(listMysqAccountsProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	listMysqlAccountsPath := listMysqlAccountsClient.Endpoint + listMysqlAccountsHttpUrl
	listMysqlAccountsPath = strings.ReplaceAll(listMysqlAccountsPath, "{project_id}", listMysqlAccountsClient.ProjectID)
	listMysqlAccountsPath = strings.ReplaceAll(listMysqlAccountsPath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	listMysqlAccountsResp, err := pagination.ListAllItems(
		listMysqlAccountsClient,
		"page",
		listMysqlAccountsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS Mysql accounts")
	}

	listMysqlAccountsRespJson, err := json.Marshal(listMysqlAccountsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listMysqlAccountsRespBody interface{}
	err = json.Unmarshal(listMysqlAccountsRespJson, &listMysqlAccountsRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("users", flattenListAccountsBody(filterListAccounts(d, listMysqlAccountsRespBody))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListAccountsBody(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	rst := make([]interface{}, 0, len(resp))

	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"name":        utils.PathSearch("name", v, nil),
			"hosts":       utils.PathSearch("hosts", v, nil),
			"description": utils.PathSearch("comment", v, nil),
		})
	}
	return rst
}

func filterListAccounts(d *schema.ResourceData, resp interface{}) []interface{} {
	accountsJson := utils.PathSearch("users", resp, make([]interface{}, 0))
	accountArray := accountsJson.([]interface{})
	if len(accountArray) < 1 {
		return nil
	}
	result := make([]interface{}, 0, len(accountArray))

	rawName, rawNameOK := d.GetOk("name")
	rawhost, rawhostOk := d.GetOk("host")

	for _, account := range accountArray {
		name := utils.PathSearch("name", account, nil)
		hosts := utils.ExpandToStringList(utils.PathSearch("hosts", account, make([]interface{}, 0)).([]interface{}))
		if rawNameOK && rawName != name {
			continue
		}
		if rawhostOk && !utils.StrSliceContains(hosts, rawhost.(string)) {
			continue
		}
		result = append(result, account)
	}

	return result
}
