package as

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AS GET /v2/{project_id}/scaling-groups/{scaling_group_id}/warm-pool-instances
func DataSourceWarmPoolInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAsWarmPoolInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"warm_pool_instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAsWarmPoolInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("autoscaling", region)
	if err != nil {
		return diag.Errorf("error creating AS client: %s", err)
	}

	getPath := client.Endpoint + "v2/{project_id}/scaling-groups/{scaling_group_id}/warm-pool-instances"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{scaling_group_id}", d.Get("scaling_group_id").(string))
	resp, err := pagination.ListAllItems(
		client,
		"marker",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving warm pool instances: %s", err)
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		return diag.Errorf("error retrieving warm pool instances: %s", err)
	}
	var respBody interface{}
	err = json.Unmarshal(respJson, &respBody)
	if err != nil {
		return diag.Errorf("error retrieving warm pool instances: %s", err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("warm_pool_instances", flattenListWarmPoolInstancesResponseBody(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListWarmPoolInstancesResponseBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("warm_pool_instances", resp, make([]interface{}, 0))

	curArray := curJson.([]interface{})

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"instance_id": utils.PathSearch("instance_id", v, nil),
			"name":        utils.PathSearch("name", v, nil),
			"status":      utils.PathSearch("status", v, nil),
		})
	}
	return rst
}
