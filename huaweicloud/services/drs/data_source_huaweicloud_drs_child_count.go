package drs

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS GET /v5/{project_id}/jobs/child-count
func DataSourceDrsChildCount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsChildCountRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the data source. If omitted, the provider-level region will be used.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the instance ID of the DDM or GaussDBv5 database.",
			},
			"db_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"gaussdbv5", "ddm"}, false),
				Description:  "Specifies the type of the database instance.",
			},
			"child_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of child tasks.",
			},
		},
	}
}

func dataSourceDrsChildCountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v5/{project_id}/jobs/child-count"
		product = "drs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildListChildCountQueryParams(d, region)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DRS child count: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("child_count", utils.PathSearch("count", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListChildCountQueryParams(d *schema.ResourceData, region string) string {
	res := fmt.Sprintf("?instance_id=%v", d.Get("instance_id"))
	res = fmt.Sprintf("%s&db_type=%v", res, d.Get("db_type"))
	res = fmt.Sprintf("%s&region=%v", res, region)

	return res
}
