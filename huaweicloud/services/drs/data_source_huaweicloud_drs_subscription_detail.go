package drs

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

// @API DRS GET /v5/{project_id}/subscriptions/{job_id}
func DataSourceDrsSubscriptionDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsSubscriptionDetailRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the data source. If omitted, the provider-level region will be used.",
			},
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the ID of the DRS subscription task.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the subscription task.",
			},
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The intranet IP address.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The enterprise project ID.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the subscription task.",
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the subscription task, in timestamp format.",
			},
			"begin_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The start time of the subscription task, in timestamp format.",
			},
			"now_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current time, in timestamp format.",
			},
			"engine_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The link type of the subscription task.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the subscription task.",
			},
			"subscription_data_type": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The subscribed data types.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_dml_subscribed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether DML is subscribed.",
						},
						"is_ddl_subscribed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether DDL is subscribed.",
						},
					},
				},
			},
			"source_endpoint": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The source database instance information of the subscription.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The database instance ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The database name.",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The database intranet IP address.",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The database port.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The database type.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The database user name.",
						},
					},
				},
			},
		},
	}
}

func dataSourceDrsSubscriptionDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v5/{project_id}/subscriptions/{job_id}"
		product = "drs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{job_id}", d.Get("job_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DRS subscription detail: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("ip", utils.PathSearch("ip", getRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", getRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
		d.Set("created_time", utils.PathSearch("created_time", getRespBody, nil)),
		d.Set("begin_time", utils.PathSearch("begin_time", getRespBody, nil)),
		d.Set("now_time", utils.PathSearch("now_time", getRespBody, nil)),
		d.Set("engine_type", utils.PathSearch("engine_type", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("subscription_data_type", flattenSubscriptionDataType(getRespBody)),
		d.Set("source_endpoint", flattenSubscriptionSourceEndpoint(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSubscriptionDataType(resp interface{}) []interface{} {
	curJson := utils.PathSearch("subscription_data_type", resp, nil)
	if curJson == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"is_dml_subscribed": utils.PathSearch("is_dml_subscribed", curJson, nil),
			"is_ddl_subscribed": utils.PathSearch("is_ddl_subscribed", curJson, nil),
		},
	}
}

func flattenSubscriptionSourceEndpoint(resp interface{}) []interface{} {
	curJson := utils.PathSearch("source_endpoint", resp, nil)
	if curJson == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"db_instance_id": utils.PathSearch("db_instance_id", curJson, nil),
			"name":           utils.PathSearch("name", curJson, nil),
			"ip":             utils.PathSearch("ip", curJson, nil),
			"port":           utils.PathSearch("port", curJson, nil),
			"type":           utils.PathSearch("type", curJson, nil),
			"user_name":      utils.PathSearch("user_name", curJson, nil),
		},
	}
}
