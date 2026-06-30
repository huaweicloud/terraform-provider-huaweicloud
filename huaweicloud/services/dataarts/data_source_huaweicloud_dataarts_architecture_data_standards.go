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

// @API DataArtsStudio GET /v2/{project_id}/design/standards
func DataSourceArchitectureDataStandards() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceArchitectureDataStandardsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the data standards are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the data standards belong.`,
			},

			// Optional parameters.
			"name_ch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The Chinese name of the data standard to be exactly queried.`,
			},
			"name_en": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The English code of the data standard to be exactly queried.`,
			},
			"directory_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The directory ID of the data standard to be queried.`,
			},
			"begin_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The start time of the data standard to be queried, in RFC3339 format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The end time of the data standard to be queried, in RFC3339 format.`,
			},

			// Attributes.
			"data_standards": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataArchitectureDataStandard(),
				Description: `The list of data standards that matched filter parameters.`,
			},
		},
	}
}

func dataArchitectureDataStandard() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the data standard.`,
			},
			"directory_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The directory ID of the data standard.`,
			},
			"directory_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The directory path of the data standard.`,
			},
			"values": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataArchitectureDataStandardValue(),
				Description: `The value of data standard attributes.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the data standard.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of data standard creator.`,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of data standard updater.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the data standard, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the data standard, in RFC3339 format.`,
			},
			"new_biz": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataArchitectureDataStandardNewBiz(),
				Description: `The biz info of the data standard.`,
			},
		},
	}
}

func dataArchitectureDataStandardValue() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"fd_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the data standard attribute.`,
			},
			"fd_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the data standard attribute.`,
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the data standard attribute.`,
			},
			"fd_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the data standard attribute definition.`,
			},
			"row_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The row ID of the data standard attribute.`,
			},
			"directory_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The directory ID to which the attribute belongs.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the data standard attribute.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of attribute creator.`,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of attribute updater.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the data standard attribute, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the data standard attribute, in RFC3339 format.`,
			},
		},
	}
}

func dataArchitectureDataStandardNewBiz() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the new biz.`,
			},
			"biz_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the new biz.`,
			},
			"biz_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the new biz.`,
			},
			"biz_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The info of the new biz.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the new biz.`,
			},
			"biz_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version of the new biz.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the new biz was created, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the new biz was updated, in RFC3339 format.`,
			},
		},
	}
}

func buildArchitectureDataStandardsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("directory_id"); ok {
		res = fmt.Sprintf("%s&directory_id=%v", res, v)
	}
	if v, ok := d.GetOk("name_ch"); ok {
		res = fmt.Sprintf("%s&name_ch=%v", res, v)
	}
	if v, ok := d.GetOk("name_en"); ok {
		res = fmt.Sprintf("%s&name_en=%v", res, v)
	}
	if v, ok := d.GetOk("begin_time"); ok {
		res = fmt.Sprintf("%s&begin_time=%s", res, url.QueryEscape(v.(string)))
	}
	if v, ok := d.GetOk("end_time"); ok {
		res = fmt.Sprintf("%s&end_time=%s", res, url.QueryEscape(v.(string)))
	}

	return res
}

func listArchitectureDataStandards(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/design/standards?limit={limit}&need_path=true"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildArchitectureDataStandardsQueryParams(d)

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

func flattenArchitectureDataStandards(standards []interface{}) []map[string]interface{} {
	if len(standards) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(standards))
	for _, standard := range standards {
		result = append(result, map[string]interface{}{
			"id":             utils.PathSearch("id", standard, nil),
			"directory_id":   utils.PathSearch("directory_id", standard, nil),
			"directory_path": utils.PathSearch("directory_path", standard, nil),
			"status":         utils.PathSearch("status", standard, nil),
			"created_by":     utils.PathSearch("create_by", standard, nil),
			"updated_by":     utils.PathSearch("update_by", standard, nil),
			"created_at":     utils.PathSearch("create_time", standard, nil),
			"updated_at":     utils.PathSearch("update_time", standard, nil),
			"values":         flattenGetDataStandardResponseBodyValue(standard),
			"new_biz":        flattenGetDataStandardResponseBodyNewBiz(standard),
		})
	}
	return result
}

func dataSourceArchitectureDataStandardsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	standards, err := listArchitectureDataStandards(client, d)
	if err != nil {
		return diag.Errorf("error querying DataArts Architecture data standards: %s", err)
	}

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_standards", flattenArchitectureDataStandards(standards)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
