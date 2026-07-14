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

// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/slow-sql/optimization/stack
func DataSourceSqlStackInformation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSqlStackInformationRead,

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
			"pid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"comp_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gs_stack": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSqlStackInformationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/slow-sql/optimization/stack"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath += buildGetSqlStackInformationQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB SQL stack information: %s", err)
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
		d.Set("gs_stack", utils.PathSearch("gs_stack", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetSqlStackInformationQueryParams(d *schema.ResourceData) string {
	res := ""
	res = fmt.Sprintf("%s&pid=%v", res, d.Get("pid"))
	res = fmt.Sprintf("%s&node_id=%v", res, d.Get("node_id"))
	if v, ok := d.GetOk("comp_id"); ok {
		res = fmt.Sprintf("%s&comp_id=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
