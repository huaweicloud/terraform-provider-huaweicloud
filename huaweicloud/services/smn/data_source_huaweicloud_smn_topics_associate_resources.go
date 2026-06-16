package smn

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SMN GET /v2/{project_id}/notifications/topics-with-associated-resources
func DataSourceSmnTopicsAssociateResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSmnTopicsAssociateResourcesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"topic_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fuzzy_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fuzzy_display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"topics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topicsAssociateResourcesSchema(),
			},
		},
	}
}

func topicsAssociateResourcesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"topic_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"push_policy": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"topic_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topicsTagsSchema(),
			},
			"attributes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topicsAccessPolicySchema(),
			},
			"logtanks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topicsLogtanksSchema(),
			},
		},
	}
}

func topicsTagsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func topicsAccessPolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"access_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func topicsLogtanksSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_stream_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSmnTopicsAssociateResourcesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
	)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("smn", region)
	if err != nil {
		return diag.Errorf("error creating SMN client: %s", err)
	}

	httpUrl := "v2/{project_id}/notifications/topics-with-associated-resources"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildTopicsAssociateResourcesQueryParams(d, epsId)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving SMN topics with associated resources: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("topics", flattenTopicsAssociateResources(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildTopicsAssociateResourcesQueryParams(d *schema.ResourceData, epsId string) string {
	res := ""
	if v, ok := d.GetOk("topic_id"); ok {
		res = fmt.Sprintf("%s&topic_id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("fuzzy_name"); ok {
		res = fmt.Sprintf("%s&fuzzy_name=%v", res, v)
	}
	if v, ok := d.GetOk("fuzzy_display_name"); ok {
		res = fmt.Sprintf("%s&fuzzy_display_name=%v", res, v)
	}
	if epsId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsId)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenTopicsAssociateResources(resp interface{}) []interface{} {
	curJson := utils.PathSearch("topics", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"topic_urn":             utils.PathSearch("topic_urn", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"display_name":          utils.PathSearch("display_name", v, nil),
			"push_policy":           utils.PathSearch("push_policy", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"topic_id":              utils.PathSearch("topic_id", v, nil),
			"create_time":           utils.PathSearch("create_time", v, nil),
			"update_time":           utils.PathSearch("update_time", v, nil),
			"tags":                  flattenTopicsTags(v),
			"attributes":            flattenTopicsAttributes(v),
			"logtanks":              flattenTopicsLogtanks(v),
		})
	}
	return res
}

func flattenTopicsTags(resp interface{}) []interface{} {
	curJson := utils.PathSearch("tags", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}
	return res
}

func flattenTopicsAttributes(resp interface{}) []interface{} {
	curJson := utils.PathSearch("attributes", resp, nil)
	if curJson == nil {
		return nil
	}

	res := []interface{}{
		map[string]interface{}{
			"access_policy": utils.PathSearch("access_policy", curJson, nil),
			"create_time":   utils.PathSearch("create_time", curJson, nil),
			"update_time":   utils.PathSearch("update_time", curJson, nil),
		},
	}
	return res
}

func flattenTopicsLogtanks(resp interface{}) []interface{} {
	curJson := utils.PathSearch("logtanks", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":            utils.PathSearch("id", v, nil),
			"log_group_id":  utils.PathSearch("log_group_id", v, nil),
			"log_stream_id": utils.PathSearch("log_stream_id", v, nil),
			"create_time":   utils.PathSearch("create_time", v, nil),
			"update_time":   utils.PathSearch("update_time", v, nil),
		})
	}
	return res
}
