package dsc

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC POST /v1/{project_id}/logtrace/data-links
func DataSourceLinks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLinksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dbss_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ecs_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"egress_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"internet_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"labels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"oem_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// In the API documentation, it is of type `int`, but here it is modified to `string` to support
			// scenarios with `0`.
			"sensitive_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     linkSchema(),
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     nodeSchema(),
			},
		},
	}
}

func linkSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"access_times": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"encrypt_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_access_times": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"target_node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func nodeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"encrypt_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fixed_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"floating_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sensitive_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     sensitiveInfoSchema(),
			},
		},
	}
}

func sensitiveInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"label_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sensitive_level": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildLinksBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := make(map[string]interface{})
	if v, ok := d.GetOk("db_name"); ok {
		bodyParams["db_name"] = v
	}

	if v, ok := d.GetOk("dbss_instance_id"); ok {
		bodyParams["dbss_instance_id"] = v
	}

	if v, ok := d.GetOk("ecs_name"); ok {
		bodyParams["ecs_name"] = v
	}

	if v, ok := d.GetOk("egress_type"); ok {
		bodyParams["egress_type"] = v
	}

	if v, ok := d.GetOk("end_time"); ok {
		bodyParams["end_time"] = v
	}

	if v, ok := d.GetOk("internet_ip"); ok {
		bodyParams["internet_ip"] = v
	}

	labelsInput := d.Get("labels").([]interface{})
	if len(labelsInput) > 0 {
		bodyParams["labels"] = utils.ExpandToStringList(labelsInput)
	}

	if v, ok := d.GetOk("oem_instance_id"); ok {
		bodyParams["oem_instance_id"] = v
	}

	if v, ok := d.GetOk("sensitive_level"); ok {
		bodyParams["sensitive_level"] = v
	}

	if v, ok := d.GetOk("start_time"); ok {
		bodyParams["start_time"] = v
	}

	return bodyParams
}

func dataSourceLinksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/logtrace/data-links"
		product = "dsc"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
		JSONBody: buildLinksBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC links: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("links", flattenLinks(
			utils.PathSearch("links", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("nodes", flattenNodes(
			utils.PathSearch("nodes", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenLinks(links []interface{}) []interface{} {
	if len(links) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(links))
	for _, v := range links {
		rst = append(rst, map[string]interface{}{
			"access_times":     utils.PathSearch("access_times", v, nil),
			"encrypt_status":   utils.PathSearch("encrypt_status", v, nil),
			"source_node_id":   utils.PathSearch("source_node_id", v, nil),
			"ssl_access_times": utils.PathSearch("ssl_access_times", v, nil),
			"target_node_id":   utils.PathSearch("target_node_id", v, nil),
		})
	}

	return rst
}

func flattenNodes(nodes []interface{}) []interface{} {
	if len(nodes) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(nodes))
	for _, v := range nodes {
		rst = append(rst, map[string]interface{}{
			"encrypt_status": utils.PathSearch("encrypt_status", v, nil),
			"fixed_ip":       utils.PathSearch("fixed_ip", v, nil),
			"floating_ip":    utils.PathSearch("floating_ip", v, nil),
			"id":             utils.PathSearch("id", v, nil),
			"name":           utils.PathSearch("name", v, nil),
			"node_type":      utils.PathSearch("node_type", v, nil),
			"sensitive_infos": flattenSensitiveInfos(
				utils.PathSearch("sensitive_infos", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenSensitiveInfos(sensitiveInfos []interface{}) []interface{} {
	if len(sensitiveInfos) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(sensitiveInfos))
	for _, v := range sensitiveInfos {
		rst = append(rst, map[string]interface{}{
			"label_name":      utils.PathSearch("label_name", v, nil),
			"sensitive_level": utils.PathSearch("sensitive_level", v, nil),
		})
	}

	return rst
}
