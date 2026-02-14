package oms

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API OMS POST /v2/{project_id}/objectstorage/buckets/objects
func DataSourceBucketObjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBucketObjectsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cloud_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ak": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"sk": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"file_path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"json_auth_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildListBucketObjectsBody(d *schema.ResourceData, region string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"cloud_type":        d.Get("cloud_type"),
		"ak":                d.Get("ak"),
		"sk":                d.Get("sk"),
		"bucket_name":       d.Get("bucket_name"),
		"file_path":         d.Get("file_path"),
		"json_auth_file":    utils.ValueIgnoreEmpty(d.Get("json_auth_file")),
		"connection_string": utils.ValueIgnoreEmpty(d.Get("connection_string")),
		"app_id":            utils.ValueIgnoreEmpty(d.Get("app_id")),
		"data_center":       region,
		"page_size":         100,
	}

	return bodyParams
}

func listBucketObjects(client *golangsdk.ServiceClient, d *schema.ResourceData, region string) ([]interface{}, error) {
	var (
		httpUrl     = "v2/{project_id}/objectstorage/buckets/objects"
		nextMarker  = ""
		result      = make([]interface{}, 0)
		requestBody = utils.RemoveNil(buildListBucketObjectsBody(d, region))
	)

	requestBody["behind_filename"] = nextMarker
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         requestBody,
	}

	for {
		requestResp, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return nil, err
		}

		getRespBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		records := utils.PathSearch("records", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(records) == 0 {
			break
		}

		result = append(result, records...)

		nextMarker = utils.PathSearch("next_marker", getRespBody, "").(string)
		if nextMarker == "" {
			break
		}

		requestBody["behind_filename"] = nextMarker
	}

	return result, nil
}

func dataSourceBucketObjectsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("oms", region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	respbody, err := listBucketObjects(client, d, region)
	if err != nil {
		return diag.Errorf("error retrieving bucket objects: %s", err)
	}

	datasourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(datasourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("records", flattenBucketObjects(respbody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBucketObjects(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"name": utils.PathSearch("name", v, nil),
			"size": utils.PathSearch("size", v, nil),
		})
	}
	return rst
}
