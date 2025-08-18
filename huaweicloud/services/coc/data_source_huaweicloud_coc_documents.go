package coc

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API COC GET /v1/documents
func DataSourceCocDocuments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCocDocumentsRead,

		Schema: map[string]*schema.Schema{
			"name_like": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"document_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"document_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modifier": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCocDocumentsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("coc", region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	documents, err := queryDocuments(client, d)
	if err != nil {
		return diag.Errorf("error querying documents: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		nil,
		d.Set("data", flattenCocGetDocuments(documents)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func queryDocuments(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/documents"
		offset  = 0
		result  = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath += buildGetDocumentsParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		documents := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		if len(documents) < 1 {
			break
		}
		result = append(result, documents...)
		offset += len(documents)
	}

	return result, nil
}

func buildGetDocumentsParams(d *schema.ResourceData) string {
	res := "?limit=100"
	if v, ok := d.GetOk("name_like"); ok {
		res = fmt.Sprintf("%s&name_like=%v", res, v)
	}
	if v, ok := d.GetOk("creator"); ok {
		res = fmt.Sprintf("%s&creator=%v", res, v)
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, v)
	}
	if v, ok := d.GetOk("document_type"); ok {
		res = fmt.Sprintf("%s&document_type=%v", res, v)
	}

	return res
}

func flattenCocGetDocuments(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	result := make([]interface{}, 0, len(paramsList))

	for _, params := range paramsList {
		result = append(result, map[string]interface{}{
			"document_id":           utils.PathSearch("document_id", params, nil),
			"name":                  utils.PathSearch("name", params, nil),
			"create_time":           utils.PathSearch("create_time", params, nil),
			"update_time":           utils.PathSearch("update_time", params, nil),
			"version":               utils.PathSearch("version", params, nil),
			"creator":               utils.PathSearch("creator", params, nil),
			"modifier":              utils.PathSearch("modifier", params, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", params, nil),
		})
	}

	return result
}
