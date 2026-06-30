package geminidb

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

// @API GeminiDB GET /v3/{project_id}/configurations/{config_id}/applicable-instances
func DataSourceGeminiDBPtApplicableInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeminiDBPtApplicableInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"config_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     geminiDBPtApplicableInstanceSchema(),
			},
		},
	}
}

func geminiDBPtApplicableInstanceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGeminiDBPtApplicableInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/configurations/{config_id}/applicable-instances"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{config_id}", d.Get("config_id").(string))
	getPath += buildListGeminiDBPtApplicableInstancesQueryParams(d)

	getResp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving GeminiDB applicable instances: %s", err)
	}
	getRespJson, err := json.Marshal(getResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getRespBody interface{}
	err = json.Unmarshal(getRespJson, &getRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("instances", flattenListGeminiDBPtApplicableInstances(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListGeminiDBPtApplicableInstancesQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("instance_name"); ok {
		queryParams = fmt.Sprintf("%s&instance_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("instance_id"); ok {
		queryParams = fmt.Sprintf("%s&instance_id=%v", queryParams, v)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func flattenListGeminiDBPtApplicableInstances(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	instancesRaw := utils.PathSearch("instances", resp, nil)
	if instancesRaw == nil {
		return nil
	}

	instancesSlice, ok := instancesRaw.([]interface{})
	if !ok {
		return nil
	}

	instances := make([]map[string]interface{}, 0, len(instancesSlice))
	for _, instanceRaw := range instancesSlice {
		instanceMap := map[string]interface{}{
			"id":   utils.PathSearch("id", instanceRaw, nil),
			"name": utils.PathSearch("name", instanceRaw, nil),
		}
		instances = append(instances, instanceMap)
	}

	return instances
}
