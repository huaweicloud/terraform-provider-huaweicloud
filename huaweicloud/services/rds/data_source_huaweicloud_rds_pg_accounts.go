package rds

import (
	"context"
	"encoding/json"
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
func DataSourcePgAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePgAccountsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"users": {
				Type:     schema.TypeList,
				Elem:     pgAccountsSchema(),
				Computed: true,
			},
		},
	}
}

func pgAccountsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attributes": {
				Type:     schema.TypeList,
				Elem:     attributesSchema(),
				Computed: true,
			},
			"memberof": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func attributesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"rolsuper": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"rolinherit": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"rolcreaterole": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"rolcreatedb": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"rolcanlogin": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"rolconnlimit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rolreplication": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"rolbypassrls": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourcePgAccountsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listPgAccountsHttpUrl = "v3/{project_id}/instances/{instance_id}/db_user/detail?page=1&limit=100"
		listPgAccountsProduct = "rds"
	)
	listPgAccountsClient, err := cfg.NewServiceClient(listPgAccountsProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	listPgAccountsPath := listPgAccountsClient.Endpoint + listPgAccountsHttpUrl
	listPgAccountsPath = strings.ReplaceAll(listPgAccountsPath, "{project_id}", listPgAccountsClient.ProjectID)
	listPgAccountsPath = strings.ReplaceAll(listPgAccountsPath, "{instance_id}", d.Get("instance_id").(string))

	listPgAccountsResp, err := pagination.ListAllItems(
		listPgAccountsClient,
		"page",
		listPgAccountsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS PostgreSQL accounts")
	}

	listPgAccountsRespJson, err := json.Marshal(listPgAccountsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listPgAccountsRespBody interface{}
	err = json.Unmarshal(listPgAccountsRespJson, &listPgAccountsRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("users", flattenPgAccountsBody(filterPgAccounts(d, listPgAccountsRespBody))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPgAccountsBody(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	rst := make([]interface{}, 0, len(resp))

	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"name":        utils.PathSearch("name", v, nil),
			"attributes":  flattenAttributes(v),
			"memberof":    utils.PathSearch("memberof", v, nil),
			"description": utils.PathSearch("comment", v, nil),
		})
	}
	return rst
}

func filterPgAccounts(d *schema.ResourceData, resp interface{}) []interface{} {
	accountsJson := utils.PathSearch("users", resp, make([]interface{}, 0))
	accountArray := accountsJson.([]interface{})
	if len(accountArray) < 1 {
		return nil
	}
	result := make([]interface{}, 0)

	rawUserName, rawUserNameOK := d.GetOk("user_name")

	for _, account := range accountArray {
		name := utils.PathSearch("name", account, nil)
		if rawUserNameOK && rawUserName != name {
			continue
		}
		result = append(result, account)
	}
	return result
}

func flattenAttributes(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("attributes", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"rolsuper":       utils.PathSearch("rolsuper", curJson, nil),
			"rolinherit":     utils.PathSearch("rolinherit", curJson, nil),
			"rolcreaterole":  utils.PathSearch("rolcreaterole", curJson, nil),
			"rolcreatedb":    utils.PathSearch("rolcreatedb", curJson, nil),
			"rolcanlogin":    utils.PathSearch("rolcanlogin", curJson, nil),
			"rolconnlimit":   utils.PathSearch("rolconnlimit", curJson, nil),
			"rolreplication": utils.PathSearch("rolreplication", curJson, nil),
			"rolbypassrls":   utils.PathSearch("rolbypassrls", curJson, nil),
		},
	}
	return rst
}
