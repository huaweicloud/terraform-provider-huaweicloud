package drs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS GET /v5/{project_id}/job/{job_id}/object-column-info
func DataSourceDrsObjectColumnInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsObjectColumnInfoRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the data source. If omitted, the provider-level region will be used.",
			},
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the ID of the DRS job to query object column info.",
			},
			"object_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the object ID. If not specified, database-level objects are queried by default.",
			},
			"is_refresh": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether to force refresh the query result.",
			},
			"object_with_column_infos": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of the objects related to column info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the node.",
						},
						"parent_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The parent node ID.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The node type.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The node name.",
						},
						"alias_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The node alias name.",
						},
						"notices": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The prompt messages, such as too many tables under a database.",
						},
						"extend_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The extended information.",
						},
						"is_support_expand": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to support expand query.",
						},
						"has_column_info": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the object has column info.",
						},
						"is_preset": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the object is preset.",
						},
						"token_count": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The token count.",
						},
						"is_sent": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the object has been sent to node.",
						},
						"sent_alias_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The alias name sent to node.",
						},
					},
				},
			},
		},
	}
}

func dataSourceDrsObjectColumnInfoRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v5/{project_id}/job/{job_id}/object-column-info"
		product = "drs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{job_id}", d.Get("job_id").(string))
	listPath += buildListObjectColumnInfoQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving DRS object column info: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("object_with_column_infos", flattenObjectWithColumnInfos(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListObjectColumnInfoQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("object_id"); ok {
		res = fmt.Sprintf("%s&object_id=%v", res, v)
	}
	if v, ok := d.GetOk("is_refresh"); ok {
		res = fmt.Sprintf("%s&is_refresh=%v", res, v.(bool))
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenObjectWithColumnInfos(resp interface{}) []interface{} {
	curJson := utils.PathSearch("object_with_column_infos", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"parent_id":         utils.PathSearch("parent_id", v, nil),
			"type":              utils.PathSearch("type", v, nil),
			"name":              utils.PathSearch("name", v, nil),
			"alias_name":        utils.PathSearch("alias_name", v, nil),
			"notices":           utils.PathSearch("notices", v, make([]interface{}, 0)),
			"extend_info":       utils.PathSearch("extend_info", v, nil),
			"is_support_expand": utils.StringToBool(utils.PathSearch("is_support_expand", v, nil)),
			"has_column_info":   utils.StringToBool(utils.PathSearch("has_column_info", v, nil)),
			"is_preset":         utils.StringToBool(utils.PathSearch("is_preset", v, nil)),
			"token_count":       utils.PathSearch("token_count", v, nil),
			"is_sent":           utils.StringToBool(utils.PathSearch("is_sent", v, nil)),
			"sent_alias_name":   utils.PathSearch("sent_alias_name", v, nil),
		})
	}
	return res
}
