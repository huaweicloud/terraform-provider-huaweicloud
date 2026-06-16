package geminidb

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GeminiDB GET /v3/{project_id}/dbcache/mappings
func DataSourceMemoryMappings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMemoryMappingsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mapping_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dbcache_mappings": {
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
						"source_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
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
						"rule_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildMemoryMappingsQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("mapping_id"); ok {
		queryParams = fmt.Sprintf("%s&id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("source_instance_id"); ok {
		queryParams = fmt.Sprintf("%s&source_instance_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("source_instance_name"); ok {
		queryParams = fmt.Sprintf("%s&source_instance_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("target_instance_id"); ok {
		queryParams = fmt.Sprintf("%s&target_instance_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("target_instance_name"); ok {
		queryParams = fmt.Sprintf("%s&target_instance_name=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceMemoryMappingsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/dbcache/mappings?limit=100"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildMemoryMappingsQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving the memory acceleration mappings: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		memoryMappings := utils.PathSearch("dbcache_mappings", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(memoryMappings) == 0 {
			break
		}

		result = append(result, memoryMappings...)
		offset += len(memoryMappings)
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("dbcache_mappings", flattenMemoryMappings(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMemoryMappings(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"id":                   utils.PathSearch("id", v, nil),
			"name":                 utils.PathSearch("name", v, nil),
			"source_instance_id":   utils.PathSearch("source_instance_id", v, nil),
			"source_instance_name": utils.PathSearch("source_instance_name", v, nil),
			"target_instance_id":   utils.PathSearch("target_instance_id", v, nil),
			"target_instance_name": utils.PathSearch("target_instance_name", v, nil),
			"status":               utils.PathSearch("status", v, nil),
			"created":              utils.PathSearch("created", v, nil),
			"updated":              utils.PathSearch("updated", v, nil),
			"rule_count":           utils.PathSearch("rule_count", v, nil),
		})
	}

	return result
}
