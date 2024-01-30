// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RDS
// ---------------------------------------------------------------

package rds

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS GET /v3/{project_id}/collations
func DataSourceSQLServerCollations() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceSQLServerCollationsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"char_sets": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `Indicates the character set information list.`,
			},
		},
	}
}

func resourceSQLServerCollationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listSQLServerCollations: Query the List of RDS SQLServer collations.
	var (
		listSQLServerCollationsHttpUrl = "v3/{project_id}/collations"
		listSQLServerCollationsProduct = "rds"
	)
	listSQLServerCollationsClient, err := cfg.NewServiceClient(listSQLServerCollationsProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	listSQLServerCollationsPath := listSQLServerCollationsClient.Endpoint + listSQLServerCollationsHttpUrl
	listSQLServerCollationsPath = strings.ReplaceAll(listSQLServerCollationsPath, "{project_id}",
		listSQLServerCollationsClient.ProjectID)

	listSQLServerCollationsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listSQLServerCollationsResp, err := listSQLServerCollationsClient.Request("GET", listSQLServerCollationsPath,
		&listSQLServerCollationsOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS SQLServer collations")
	}

	listSQLServerCollationsRespBody, err := utils.FlattenResponse(listSQLServerCollationsResp)
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
		d.Set("char_sets", utils.PathSearch("charSets", listSQLServerCollationsRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
