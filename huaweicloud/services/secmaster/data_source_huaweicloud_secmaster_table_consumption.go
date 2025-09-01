package secmaster

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Secmaster GET /v2/{project_id}/workspaces/{workspace_id}/siem/tables/{table_id}/consumption
func DataSourceTableConsumption() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTableConsumptionRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"table_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"table_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dataspace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pipe_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pipe_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_point": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// This field is misspelled in the API documentation. It is consistent here with the API documentation.
			"subscirption": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTableConsumptionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v2/{project_id}/workspaces/{workspace_id}/siem/tables/{table_id}/consumption"
		workspaceID = d.Get("workspace_id").(string)
		tableID     = d.Get("table_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceID)
	requestPath = strings.ReplaceAll(requestPath, "{table_id}", tableID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving SecMaster table consumption configuration: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
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
		d.Set("table_name", utils.PathSearch("table_name", respBody, nil)),
		d.Set("dataspace_id", utils.PathSearch("dataspace_id", respBody, nil)),
		d.Set("pipe_id", utils.PathSearch("pipe_id", respBody, nil)),
		d.Set("pipe_name", utils.PathSearch("pipe_name", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("type", utils.PathSearch("type", respBody, nil)),
		d.Set("access_point", utils.PathSearch("access_point", respBody, nil)),
		d.Set("subscirption", utils.PathSearch("subscirption", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
