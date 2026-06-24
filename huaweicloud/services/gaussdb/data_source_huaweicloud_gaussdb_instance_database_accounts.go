package gaussdb

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/db-users
func DataSourceGaussDBInstanceDatabaseAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDBInstanceDatabaseAccountsRead,

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
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     databaseUserSchema(),
			},
		},
	}
}

func databaseUserSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attribute": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     userAttributeSchema(),
			},
			"memberof": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lock_status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func userAttributeSchema() *schema.Resource {
	return &schema.Resource{
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
			"rolpassworddeadline": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGaussDBInstanceDatabaseAccountsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/db-users"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving GaussDB instance database accounts: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("users", flattenGetDatabaseUsersBody(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetDatabaseUsersBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("users", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"name":        utils.PathSearch("name", v, nil),
			"attribute":   flattenGetUserAttributeBody(v),
			"memberof":    utils.PathSearch("memberof", v, nil),
			"lock_status": utils.PathSearch("lock_status", v, nil),
		})
	}
	return res
}

func flattenGetUserAttributeBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("attribute", resp, nil)
	if curJson == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"rolsuper":            utils.PathSearch("rolsuper", curJson, nil),
			"rolinherit":          utils.PathSearch("rolinherit", curJson, nil),
			"rolcreaterole":       utils.PathSearch("rolcreaterole", curJson, nil),
			"rolcreatedb":         utils.PathSearch("rolcreatedb", curJson, nil),
			"rolcanlogin":         utils.PathSearch("rolcanlogin", curJson, nil),
			"rolconnlimit":        utils.PathSearch("rolconnlimit", curJson, nil),
			"rolreplication":      utils.PathSearch("rolreplication", curJson, nil),
			"rolbypassrls":        utils.PathSearch("rolbypassrls", curJson, nil),
			"rolpassworddeadline": utils.PathSearch("rolpassworddeadline", curJson, nil),
		},
	}
}
