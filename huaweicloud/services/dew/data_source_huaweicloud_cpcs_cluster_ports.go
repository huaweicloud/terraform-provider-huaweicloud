package dew

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

// @API DEW GET /v1/{project_id}/dew/cpcs/cluster/{cluster_id}/port
func DataSourceClusterPorts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterPortsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"result": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"elb_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"elb_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"listener_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_group_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"server_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"validate_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"wrong": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"wrong_msg": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceClusterPortsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v1/{project_id}/dew/cpcs/cluster/{cluster_id}/port"
		clusterId = d.Get("cluster_id").(string)
	)

	client, err := cfg.NewServiceClient("kms", region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", clusterId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving cluster ports: %s", err)
	}

	respBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterPorts := utils.PathSearch("result", respBody, make([]interface{}, 0)).([]interface{})

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("result", flattenClusterPorts(clusterPorts)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenClusterPorts(results []interface{}) []interface{} {
	if len(results) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(results))
	for _, association := range results {
		result = append(result, map[string]interface{}{
			"id":                utils.PathSearch("id", association, nil),
			"cluster_id":        utils.PathSearch("cluster_id", association, nil),
			"elb_id":            utils.PathSearch("elb_id", association, nil),
			"elb_ip":            utils.PathSearch("elb_ip", association, nil),
			"mode":              utils.PathSearch("mode", association, nil),
			"listener_port":     utils.PathSearch("listener_port", association, nil),
			"listener_id":       utils.PathSearch("listener_id", association, nil),
			"server_group_port": utils.PathSearch("server_group_port", association, nil),
			"server_group_id":   utils.PathSearch("server_group_id", association, nil),
			"project_id":        utils.PathSearch("project_id", association, nil),
			"validate_time":     utils.PathSearch("validate_time", association, nil),
			"wrong":             utils.PathSearch("wrong", association, nil),
			"wrong_msg":         utils.PathSearch("wrong_msg", association, nil),
		})
	}

	return result
}
