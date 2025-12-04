package cdn

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN GET /v1.0/cdn/configuration/tml-apply-records
func DataSourceDomainTemplateApplyRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDomainTemplateApplyRecordsRead,

		Schema: map[string]*schema.Schema{
			// Optional parameters.
			"template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the domain template.`,
			},
			"template_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the domain template.`,
			},
			"operator_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The operation ID of the domain template.`,
			},

			// Attributes.
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of domain template apply records that matched the filter parameters.`,
				Elem:        domainTemplateApplyRecordSchema(),
			},
		},
	}
}

func domainTemplateApplyRecordSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"operator_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The operation ID.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The result of applying the template.`,
			},
			"template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the domain template.`,
			},
			"template_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the domain template.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the domain template.`,
			},
			"apply_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the domain template was applied, in RFC3339 format.`,
			},
			"type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The type of the domain template.`,
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The account ID.`,
			},
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of resources to which the template was applied.`,
				Elem:        domainTemplateApplyRecordResourceSchema(),
			},
			"configs": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The configuration of the domain template, in JSON format.`,
			},
		},
	}
}

func domainTemplateApplyRecordResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of applying the template to the domain.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain name.`,
			},
			"error_msg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The error message.`,
			},
		},
	}
}

func buildDomainTemplateApplyRecordsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("template_id"); ok {
		res = fmt.Sprintf("%s&tml_id=%v", res, v)
	}
	if v, ok := d.GetOk("template_name"); ok {
		res = fmt.Sprintf("%s&tml_name=%v", res, v)
	}
	if v, ok := d.GetOk("operator_id"); ok {
		res = fmt.Sprintf("%s&operator_id=%v", res, v)
	}

	return res
}

func listDomainTemplateApplyRecords(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1.0/cdn/configuration/tml-apply-records?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildDomainTemplateApplyRecordsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		elements := utils.PathSearch("elements", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, elements...)
		if len(elements) < limit {
			break
		}
		offset += len(elements)
	}

	return result, nil
}

func flattenDomainTemplateApplyRecordResources(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"status":      utils.PathSearch("status", item, nil),
			"domain_name": utils.PathSearch("domain_name", item, nil),
			"error_msg":   utils.PathSearch("error_msg", item, nil),
		})
	}

	return result
}

func flattenDomainTemplateApplyRecords(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"operator_id":   utils.PathSearch("operator_id", item, nil),
			"status":        utils.PathSearch("status", item, nil),
			"template_id":   utils.PathSearch("tml_id", item, nil),
			"template_name": utils.PathSearch("tml_name", item, nil),
			"description":   utils.PathSearch("remark", item, nil),
			"apply_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("apply_time", item,
				float64(0)).(float64))/1000, false),
			"type":       utils.PathSearch("type", item, nil),
			"account_id": utils.PathSearch("account_id", item, nil),
			"resources": flattenDomainTemplateApplyRecordResources(utils.PathSearch("resources", item,
				make([]interface{}, 0)).([]interface{})),
			"configs": utils.JsonToString(utils.PathSearch("configs", item, nil)),
		})
	}

	return result
}

func dataSourceDomainTemplateApplyRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	resp, err := listDomainTemplateApplyRecords(client, d)
	if err != nil {
		return diag.Errorf("error querying CDN domain template apply records: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("records", flattenDomainTemplateApplyRecords(resp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
