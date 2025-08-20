package dns

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DNS GET /v2.1/ptrs
func DataSourceDNSV21PtrRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSV21PtrRecordsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The enterprise project ID corresponding to the PTR record.`,
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The key/value pairs to associate with the PTR record.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The status of the PTR record.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The resource type.`,
			},
			"ptrrecords": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the PTR records list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the PTR record.`,
						},
						"names": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The domain names of the PTR record.`,
						},
						"publicip_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the EIP.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the PTR record.`,
						},
						"ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The time to live (TTL) of the record set (in seconds).`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The enterprise project ID of the PTR record.`,
						},
						"tags": common.TagsComputedSchema(`The key/value pairs to associate with the PTR record.`),
						"address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The address of the EIP.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the PTR record.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceDNSV21PtrRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dns_region", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	httpUrl := "v2.1/ptrs?limit=500"
	getPath := client.Endpoint + httpUrl
	getPath += buildListDNSV21PtrRecordsQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	offset := 0
	rst := make([]map[string]interface{}, 0)
	for {
		currentPath := getPath + fmt.Sprintf("&offset=%d", offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error getting PTR Records: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flattening response: %s", err)
		}

		records := utils.PathSearch("floatingips", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(records) == 0 {
			break
		}

		for _, record := range records {
			rst = append(rst, map[string]interface{}{
				"id":                    utils.PathSearch("id", record, nil),
				"names":                 utils.PathSearch("ptrdnames", record, nil),
				"description":           utils.PathSearch("description", record, nil),
				"publicip_id":           utils.PathSearch("publicip.id", record, nil),
				"ttl":                   utils.PathSearch("ttl", record, nil),
				"address":               utils.PathSearch("publicip.address", record, nil),
				"enterprise_project_id": utils.PathSearch("enterprise_project_id", record, nil),
				"status":                utils.PathSearch("status", record, nil),
				"tags":                  utils.FlattenTagsToMap(utils.PathSearch("tags", record, make([]interface{}, 0))),
			})
		}

		offset += len(records)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("ptrrecords", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListDNSV21PtrRecordsQueryParams(d *schema.ResourceData) string {
	queryParam := ""

	if v, ok := d.GetOk("tags"); ok {
		tags := v.(map[string]interface{})
		tagsList := make([]string, 0, len(tags))
		for k, v := range tags {
			tagsList = append(tagsList, k+","+v.(string))
		}
		queryParam = fmt.Sprintf("%s&tags=%v", queryParam, strings.Join(tagsList, "|"))
	}

	if v, ok := d.GetOk("status"); ok {
		queryParam = fmt.Sprintf("%s&status=%v", queryParam, v)
	}

	if v, ok := d.GetOk("resource_type"); ok {
		queryParam = fmt.Sprintf("%s&resource_type=%v", queryParam, v)
	}

	if v, ok := d.GetOk("enterprise_project_id"); ok {
		queryParam = fmt.Sprintf("%s&enterprise_project_id=%v", queryParam, v)
	}

	return queryParam
}
