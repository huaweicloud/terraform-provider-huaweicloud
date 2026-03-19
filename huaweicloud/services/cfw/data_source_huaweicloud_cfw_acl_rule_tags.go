package cfw

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CFW GET /v2/{project_id}/cfw-acl/tags
func DataSourceAclRuleTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAclRuleTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildQueryAclRuleTagsQueryParams(cfg *config.Config, d *schema.ResourceData) string {
	queryParam := fmt.Sprintf("?fw_instance_id=%s&limit=1024", d.Get("fw_instance_id").(string))

	if v := cfg.GetEnterpriseProjectID(d); v != "" {
		queryParam += fmt.Sprintf("&enterprise_project_id=%s", v)
	}

	return queryParam
}

func dataSourceAclRuleTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		product    = "cfw"
		httpUrl    = "v2/{project_id}/cfw-acl/tags"
		offset     = 0
		allResults = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildQueryAclRuleTagsQueryParams(cfg, d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%d", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving CFW acl rule tags: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		records := utils.PathSearch("data.records", respBody, make([]interface{}, 0)).([]interface{})
		if len(records) == 0 {
			break
		}

		allResults = append(allResults, records...)
		offset += len(records)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("records", flattenAclRuleTagsResponse(allResults)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAclRuleTagsResponse(respArray []interface{}) []map[string]interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"tag_id":    utils.PathSearch("tag_id", v, nil),
			"tag_key":   utils.PathSearch("tag_key", v, nil),
			"tag_value": utils.PathSearch("tag_value", v, nil),
		})
	}

	return rst
}
