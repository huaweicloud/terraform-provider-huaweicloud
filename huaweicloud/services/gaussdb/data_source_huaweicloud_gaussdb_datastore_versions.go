package gaussdb

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB GET /v3.2/{project_id}/datastore/versions
func DataSourceDatastoreVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatastoreVersionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"database_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"software_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hotfixes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"create_time": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"deploy_modes": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"default_upgrade": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"update_time": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"mode": {
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
				},
			},
		},
	}
}

func listDatastoreVersions(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v3.2/{project_id}/datastore/versions?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", listPath, offset)
		listResp, err := client.Request("GET", currentPath, &listOpt)
		if err != nil {
			return nil, err
		}

		resp, err := utils.FlattenResponse(listResp)
		if err != nil {
			return nil, err
		}

		versions := utils.PathSearch("database_versions", resp, make([]interface{}, 0)).([]interface{})
		result = append(result, versions...)
		if len(versions) < limit {
			break
		}

		offset += len(versions)
	}

	return result, nil
}

func dataSourceDatastoreVersionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	databaseVersions, err := listDatastoreVersions(client)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB datastore versions: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("database_versions", flattenDatastoreVersions(databaseVersions)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDatastoreVersions(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"software_version": utils.PathSearch("software_version", v, nil),
			"hotfixes": flattenDatastoreVersionsHotfixes(
				utils.PathSearch("hotfixes", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenDatastoreVersionsHotfixes(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"version":     utils.PathSearch("version", v, nil),
			"create_time": utils.PathSearch("create_time", v, nil),
			"deploy_modes": flattenDatastoreVersionsDeployModes(
				utils.PathSearch("deploy_modes", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenDatastoreVersionsDeployModes(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"default_upgrade": utils.PathSearch("default_upgrade", v, nil),
			"update_time":     utils.PathSearch("update_time", v, nil),
			"mode":            utils.PathSearch("mode", v, nil),
		})
	}

	return result
}
