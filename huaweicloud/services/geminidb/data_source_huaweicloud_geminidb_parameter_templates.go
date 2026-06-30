package geminidb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GeminiDB GET /v3.1/{project_id}/configurations
func DataSourceGeminiDBParameterTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeminiDBParameterTemplatesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"datastore_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_defined": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"configurations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datastore_version_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datastore_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_defined": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGeminiDBParameterTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/configurations"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildParameterTemplatesQueryParams(d)

	resp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving GeminiDB parameter templates: %s", err)
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	var respBody interface{}
	err = json.Unmarshal(respJson, &respBody)
	if err != nil {
		return diag.FromErr(err)
	}

	configurations := flattenParameterTemplatesConfigurations(respBody)

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("configurations", configurations),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildParameterTemplatesQueryParams(d *schema.ResourceData) string {
	queryParams := ""
	if v, ok := d.GetOk("datastore_name"); ok {
		queryParams = fmt.Sprintf("%s&datastore_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("mode"); ok {
		queryParams = fmt.Sprintf("%s&mode=%v", queryParams, v)
	}
	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("user_defined"); ok {
		userDefined, err := strconv.ParseBool(v.(string))
		if err != nil {
			log.Printf("[ERROR] error parsing 'user_defined' field to Boolean: %s", err)
		}
		queryParams = fmt.Sprintf("%s&user_defined=%v", queryParams, userDefined)
	}
	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}
	return queryParams
}

func flattenParameterTemplatesConfigurations(resp interface{}) []interface{} {
	curJson := utils.PathSearch("configurations", resp, make([]interface{}, 0))
	curArray, ok := curJson.([]interface{})
	if !ok {
		return nil
	}
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":                     utils.PathSearch("id", v, nil),
			"name":                   utils.PathSearch("name", v, nil),
			"description":            utils.PathSearch("description", v, nil),
			"datastore_version_name": utils.PathSearch("datastore_version_name", v, nil),
			"datastore_name":         utils.PathSearch("datastore_name", v, nil),
			"created":                utils.PathSearch("created", v, nil),
			"updated":                utils.PathSearch("updated", v, nil),
			"mode":                   utils.PathSearch("mode", v, nil),
			"user_defined":           utils.PathSearch("user_defined", v, nil),
		})
	}
	return res
}
