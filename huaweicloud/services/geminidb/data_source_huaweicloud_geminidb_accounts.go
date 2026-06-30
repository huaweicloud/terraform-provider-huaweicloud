package geminidb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GeminiDB GET /v3/{project_id}/redis/instances/{instance_id}/db-users
func DataSourceGeminiDbAccounts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeminiDbAccountsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"privilege": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"databases": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceGeminiDbAccountsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	getPath := client.Endpoint + "v3/{project_id}/redis/instances/{instance_id}/db-users"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	if v, ok := d.GetOk("name"); ok {
		getPath += fmt.Sprintf("?name=%v", v)
	}
	resp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving GeminiDB accounts: %s", err)
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		return diag.Errorf("error retrieving GeminiDB accounts: %s", err)
	}
	var respBody interface{}
	err = json.Unmarshal(respJson, &respBody)
	if err != nil {
		return diag.Errorf("error retrieving GeminiDB accounts: %s", err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("users", flattenListGeminiDbAccountsResponseBody(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListGeminiDbAccountsResponseBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("users", resp, make([]interface{}, 0))

	curArray := curJson.([]interface{})

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":      utils.PathSearch("name", v, nil),
			"type":      utils.PathSearch("type", v, nil),
			"privilege": utils.PathSearch("privilege", v, nil),
			"databases": utils.PathSearch("databases", v, nil),
		})
	}
	return rst
}
