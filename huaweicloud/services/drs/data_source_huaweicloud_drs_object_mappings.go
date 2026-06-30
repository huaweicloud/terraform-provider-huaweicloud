package drs

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS POST /v5/{project_id}/jobs/{job_id}/object-mappings
func DataSourceDrsObjectMappings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsObjectMappingsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schema_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"table_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"has_column_info": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"object_mapping_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_db_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_schema_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_table_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_db_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_schema_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_table_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"has_column_info": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildObjectMappingsRequestBody(d *schema.ResourceData, offset int) map[string]interface{} {
	bodyParams := map[string]interface{}{}

	if offset > 0 {
		bodyParams["offset"] = offset
	}

	if v, ok := d.GetOk("db_name"); ok {
		bodyParams["db_name"] = v.(string)
	}
	if v, ok := d.GetOk("schema_name"); ok {
		bodyParams["schema_name"] = v.(string)
	}
	if v, ok := d.GetOk("table_name"); ok {
		bodyParams["table_name"] = v.(string)
	}
	if v, ok := d.GetOk("has_column_info"); ok {
		hasColumnInfo, err := strconv.ParseBool(v.(string))
		if err != nil {
			log.Printf("[ERROR] error parsing 'has_column_info' field to Boolean: %s", err)
		}
		bodyParams["has_column_info"] = hasColumnInfo
	}

	return utils.RemoveNil(bodyParams)
}

func dataSourceDrsObjectMappingsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v5/{project_id}/jobs/{job_id}/object-mappings"
		result  = make([]interface{}, 0)
		offset  = 0
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Get("job_id").(string))

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		reqOpt.JSONBody = buildObjectMappingsRequestBody(d, offset)

		listResp, err := client.Request("POST", requestPath, &reqOpt)

		if err != nil {
			return diag.Errorf("error retrieving DRS job object mappings data: %s", err)
		}

		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		objectMappings := utils.PathSearch("object_mapping_list", listRespBody, make([]interface{}, 0)).([]interface{})

		if len(objectMappings) == 0 {
			break
		}

		result = append(result, objectMappings...)

		offset += len(objectMappings)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("object_mapping_list", flattenObjectMappings(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenObjectMappings(objectMappingsResp []interface{}) []interface{} {
	if len(objectMappingsResp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(objectMappingsResp))
	for _, v := range objectMappingsResp {
		rst = append(rst, map[string]interface{}{
			"source_db_name":     utils.PathSearch("source_db_name", v, nil),
			"source_schema_name": utils.PathSearch("source_schema_name", v, nil),
			"source_table_name":  utils.PathSearch("source_table_name", v, nil),
			"target_db_name":     utils.PathSearch("target_db_name", v, nil),
			"target_schema_name": utils.PathSearch("target_schema_name", v, nil),
			"target_table_name":  utils.PathSearch("target_table_name", v, nil),
			"has_column_info":    utils.PathSearch("has_column_info", v, nil),
		})
	}
	return rst
}
