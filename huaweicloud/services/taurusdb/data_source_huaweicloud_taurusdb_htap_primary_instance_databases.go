package taurusdb

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

// @API TaurusDB POST /v3/{project_id}/instances/{instance_id}/htap/databases
func DataSourceTaurusDBHtapPrimaryInstanceDatabases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBHtapPrimaryInstanceDatabasesRead,

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
			"databases": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"source_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"database_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceTaurusDBHtapPrimaryInstanceDatabasesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/htap/databases"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildGetHtapPrimaryInstanceDatabasesBodyParams(d),
	}

	offset := 0
	result := make([]interface{}, 0)
	for {
		currentPath := fmt.Sprintf("%s?limit=100&offset=%v", listPath, offset)

		listResp, err := client.Request("POST", currentPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving TaurusDB HTAP primary instance databases: %s", err)
		}

		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		databases := utils.PathSearch("databases", listRespBody, make([]interface{}, 0)).([]interface{})

		if len(databases) == 0 {
			break
		}
		result = append(result, databases...)

		totalCount := utils.PathSearch("total_count", listRespBody, float64(0)).(float64)
		if int(totalCount) == len(result) {
			break
		}

		offset += len(databases)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("database_names", utils.ExpandToStringList(result)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetHtapPrimaryInstanceDatabasesBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"databases":          utils.ExpandToStringList(d.Get("databases").([]interface{})),
		"source_instance_id": d.Get("source_instance_id").(string),
	}
	return bodyParams
}
