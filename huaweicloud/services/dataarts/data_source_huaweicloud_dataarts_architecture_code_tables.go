package dataarts

import (
	"context"
	"fmt"
	"net/url"
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

// @API DataArtsStudio GET /v2/{project_id}/design/code-tables
func DataSourceArchitectureCodeTables() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceArchitectureCodeTablesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the code tables are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the code tables belong.`,
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The Chinese name of the code table to be exactly queried.`,
			},
			"code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The English name of the code table to be exactly queried.`,
			},
			"create_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The creator name of the code table to be queried.`,
			},
			"directory_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The directory ID of the code table to be queried.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The status of the code table to be queried.`,
			},
			"begin_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The start time of the code table to be queried, in RFC3339 format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The end time of the code table to be queried, in RFC3339 format.`,
			},

			// Attributes.
			"code_tables": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataArchitectureCodeTable(),
				Description: `The list of code tables that matched filter parameters.`,
			},
		},
	}
}

func dataArchitectureCodeTable() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the code table.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the code table.`,
			},
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The code of the code table.`,
			},
			"directory_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The directory ID of the code table.`,
			},
			"fields": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataArchitectureCodeTableField(),
				Description: `The fields information of the code table.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the code table.`,
			},
			"directory_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The directory path of the code table.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user who created the code table.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the code table was created.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the code table was updated.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the code table.`,
			},
		},
	}
}

func dataArchitectureCodeTableField() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the field.`,
			},
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The code of the field.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the field.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the field.`,
			},
			"ordinal": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The ordinal of the field.`,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the field.`,
			},
		},
	}
}

func buildArchitectureCodeTablesQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name_ch=%v", res, v)
	}
	if v, ok := d.GetOk("code"); ok {
		res = fmt.Sprintf("%s&name_en=%v", res, v)
	}
	if v, ok := d.GetOk("create_by"); ok {
		res = fmt.Sprintf("%s&create_by=%v", res, v)
	}
	if v, ok := d.GetOk("approver"); ok {
		res = fmt.Sprintf("%s&approver=%v", res, v)
	}
	if v, ok := d.GetOk("directory_id"); ok {
		res = fmt.Sprintf("%s&directory_id=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("begin_time"); ok {
		res = fmt.Sprintf("%s&begin_time=%s", res, url.QueryEscape(v.(string)))
	}
	if v, ok := d.GetOk("end_time"); ok {
		res = fmt.Sprintf("%s&end_time=%s", res, url.QueryEscape(v.(string)))
	}

	return res
}

func listArchitectureCodeTables(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/design/code-tables?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildArchitectureCodeTablesQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(d.Get("workspace_id").(string)),
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		records := utils.PathSearch("data.value.records", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, records...)

		if len(records) < limit {
			break
		}
		offset += len(records)
	}

	return result, nil
}

func flattenArchitectureCodeTables(codeTables []interface{}) []map[string]interface{} {
	if len(codeTables) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(codeTables))
	for _, codeTable := range codeTables {
		result = append(result, map[string]interface{}{
			"id":             utils.PathSearch("id", codeTable, nil),
			"name":           utils.PathSearch("name_ch", codeTable, nil),
			"code":           utils.PathSearch("name_en", codeTable, nil),
			"directory_id":   utils.PathSearch("directory_id", codeTable, nil),
			"directory_path": utils.PathSearch("directory_path", codeTable, nil),
			"description":    utils.PathSearch("description", codeTable, nil),
			"created_by":     utils.PathSearch("create_by", codeTable, nil),
			"created_at":     utils.PathSearch("create_time", codeTable, nil),
			"updated_at":     utils.PathSearch("update_time", codeTable, nil),
			"status":         utils.PathSearch("status", codeTable, nil),
			"fields":         flattenCodeTableFields(codeTable),
		})
	}
	return result
}

func dataSourceArchitectureCodeTablesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	codeTables, err := listArchitectureCodeTables(client, d)
	if err != nil {
		return diag.Errorf("error querying DataArts Architecture code tables: %s", err)
	}

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("code_tables", flattenArchitectureCodeTables(codeTables)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
