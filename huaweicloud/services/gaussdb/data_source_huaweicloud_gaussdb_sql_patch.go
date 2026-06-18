package gaussdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/sql-patch
func DataSourceGaussDbSqlPatch() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDbSqlPatchRead,

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
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sql_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"patch_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"patch_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGaussDbSqlPatchRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/sql-patch"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath += buildGetGaussDbSqlPatchQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB SQL patch: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
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
		d.Set("patch_name", utils.PathSearch("patch_name", getRespBody, nil)),
		d.Set("hint", utils.PathSearch("hint", getRespBody, nil)),
		d.Set("patch_status", utils.PathSearch("patch_status", getRespBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetGaussDbSqlPatchQueryParams(d *schema.ResourceData) string {
	res := ""
	res = fmt.Sprintf("%s&node_id=%v", res, d.Get("node_id").(string))
	res = fmt.Sprintf("%s&sql_id=%v", res, d.Get("sql_id").(string))

	if v, ok := d.GetOk("database_name"); ok {
		res = fmt.Sprintf("%s&database_name=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
