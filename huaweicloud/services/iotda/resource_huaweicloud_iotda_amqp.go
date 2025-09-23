package iotda

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA DELETE /v5/iot/{project_id}/amqp-queues/{queue_id}
// @API IoTDA GET /v5/iot/{project_id}/amqp-queues/{queue_id}
// @API IoTDA POST /v5/iot/{project_id}/amqp-queues
func ResourceAmqp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAmqpCreate,
		ReadContext:   resourceAmqpRead,
		DeleteContext: resourceAmqpDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateAmqpBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"queue_name": d.Get("name").(string),
	}
}

func resourceAmqpCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		product   = "iotda"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	createPath := client.Endpoint + "v5/iot/{project_id}/amqp-queues"
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateAmqpBodyParams(d),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating IoTDA AMQP queue: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	queueId := utils.PathSearch("queue_id", createRespBody, "").(string)
	if queueId == "" {
		return diag.Errorf("error creating IoTDA AMQP queue: ID is not found in API response")
	}

	d.SetId(queueId)

	return resourceAmqpRead(ctx, d, meta)
}

func ReadAmqpById(client *golangsdk.ServiceClient, queueId string) (interface{}, error) {
	getPath := client.Endpoint + "v5/iot/{project_id}/amqp-queues/{queue_id}"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{queue_id}", queueId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IoTDA AMQP queue: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	// When the resource does not exist, the query API returns a `200` status code and `queue_id` returns null.
	// For this situation, a `404` status code needs to be returned for subsequent checkDeleted logic processing.
	queueIdResp := utils.PathSearch("queue_id", getRespBody, "").(string)
	if queueIdResp == "" {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func resourceAmqpRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		product   = "iotda"
		queueId   = d.Id()
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	queueResp, err := ReadAmqpById(client, queueId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA AMQP queue")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("queue_name", queueResp, nil)),
		d.Set("created_at", utils.PathSearch("create_time", queueResp, nil)),
		d.Set("updated_at", utils.PathSearch("last_modify_time", queueResp, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAmqpDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		product   = "iotda"
		queueId   = d.Id()
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	deletePath := client.Endpoint + "v5/iot/{project_id}/amqp-queues/{queue_id}"
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{queue_id}", queueId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	// When deleting non-existent resource, the API returns a `404` status code.
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting IoTDA AMQP queue")
	}

	return nil
}
