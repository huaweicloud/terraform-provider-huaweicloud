package gaussdb

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB GET /v3/{project_id}/datastore/instances
func DataSourceDatastoreInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatastoreInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine_instance_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instances": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_name": {
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
									"solution": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hotfix_versions": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func listDatastoreInstances(client *golangsdk.ServiceClient) (interface{}, error) {
	httpUrl := "v3/{project_id}/datastore/instances"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func dataSourceDatastoreInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	resp, err := listDatastoreInstances(client)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB datastore instances: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("engine_instance_details", flattenDatastoreInstances(
			utils.PathSearch("engine_instance_details", resp, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDatastoreInstances(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"engine_version": utils.PathSearch("engine_version", v, nil),
			"instances": flattenInstancesInfo(
				utils.PathSearch("instances", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenInstancesInfo(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"instance_id":     utils.PathSearch("instance_id", v, nil),
			"instance_name":   utils.PathSearch("instance_name", v, nil),
			"status":          utils.PathSearch("status", v, nil),
			"type":            utils.PathSearch("type", v, nil),
			"solution":        utils.PathSearch("solution", v, nil),
			"hotfix_versions": utils.PathSearch("hotfix_versions", v, nil),
		})
	}

	return result
}
