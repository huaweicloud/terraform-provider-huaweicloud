// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product EG
// ---------------------------------------------------------------

package eg

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EG POST /v1/{project_id}/connections
// @API EG PUT /v1/{project_id}/connections/{connection_id}
// @API EG GET /v1/{project_id}/connections/{connection_id}
// @API EG DELETE /v1/{project_id}/connections/{connection_id}
func ResourceConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectionCreate,
		UpdateContext: resourceConnectionUpdate,
		ReadContext:   resourceConnectionRead,
		DeleteContext: resourceConnectionDelete,
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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the connection.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the VPC to which the connection belongs.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the subnet to which the connection belongs.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the connection.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the type of the connection.`,
			},
			"kafka_detail": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        connectionKafkaDetailSchema(),
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the configuration details of the kafka intance.`,
			},

			// Attributes
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the connection.`,
			},
			"agency": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user-delegated name used for private network target connection.`,
			},
			"flavor": {
				Type:        schema.TypeList,
				Elem:        connectionKafkaFlavorSchema(),
				Computed:    true,
				Description: `The configuration details of the kafka instance.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the connection.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The last update time of the connection.`,
			},
		},
	}
}

func connectionKafkaDetailSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the kafka instance.`,
			},
			"connect_address": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the IP address of the kafka instance.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the user name of the kafka instance.`,
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: `Specifies the password of the kafka instance.`,
			},
			"acks": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Description: `Specifies the number of confirmation signals the procuder
                    needs to receive to consider the message sent successfully.`,
			},
			"security_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the security protocol of the kafka instance.`,
			},
		},
	}
	return &sc
}

func connectionKafkaFlavorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"bandwidth_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The bandwidth type of the kafka instance.`,
			},
			"concurrency": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The concurrency number of the kafka instance..`,
			},
			"concurrency_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The concurrency type of the kafka instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the kafka instance.`,
			},
		},
	}
	return &sc
}

func resourceConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createConnection: Create an EG Connection.
	var (
		createConnectionHttpUrl = "v1/{project_id}/connections"
		createConnectionProduct = "eg"
	)
	createConnectionClient, err := cfg.NewServiceClient(createConnectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	createConnectionPath := createConnectionClient.Endpoint + createConnectionHttpUrl
	createConnectionPath = strings.ReplaceAll(createConnectionPath, "{project_id}", createConnectionClient.ProjectID)

	createConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	createConnectionOpt.JSONBody = utils.RemoveNil(buildCreateConnectionBodyParams(d))
	createConnectionResp, err := createConnectionClient.Request("POST", createConnectionPath, &createConnectionOpt)
	if err != nil {
		return diag.Errorf("error creating Connection: %s", err)
	}

	createConnectionRespBody, err := utils.FlattenResponse(createConnectionResp)
	if err != nil {
		return diag.FromErr(err)
	}

	connectionId := utils.PathSearch("id", createConnectionRespBody, "").(string)
	if connectionId == "" {
		return diag.Errorf("unable to find the connection ID from the API response")
	}
	d.SetId(connectionId)

	return resourceConnectionRead(ctx, d, meta)
}

func buildCreateConnectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":         d.Get("name"),
		"vpc_id":       d.Get("vpc_id"),
		"subnet_id":    d.Get("subnet_id"),
		"description":  utils.ValueIgnoreEmpty(d.Get("description")),
		"type":         utils.ValueIgnoreEmpty(d.Get("type")),
		"kafka_detail": buildCreateConnectionRequestBodyKafkaDetail(d.Get("kafka_detail")),
	}
	return bodyParams
}

func buildCreateConnectionRequestBodyKafkaDetail(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"instance_id":       utils.ValueIgnoreEmpty(raw["instance_id"]),
			"addr":              utils.ValueIgnoreEmpty(raw["connect_address"]),
			"acks":              utils.ValueIgnoreEmpty(raw["acks"]),
			"security_protocol": utils.ValueIgnoreEmpty(raw["security_protocol"]),
		}

		if raw["user_name"].(string) != "" {
			params["sasl_ssl"] = true
			params["username"] = raw["user_name"]
			params["password"] = raw["password"]
		}
		return params
	}
	return nil
}

func resourceConnectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getConnection: Query the EG Connection detail
	var (
		getConnectionHttpUrl = "v1/{project_id}/connections/{id}"
		getConnectionProduct = "eg"
	)
	getConnectionClient, err := cfg.NewServiceClient(getConnectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	getConnectionPath := getConnectionClient.Endpoint + getConnectionHttpUrl
	getConnectionPath = strings.ReplaceAll(getConnectionPath, "{project_id}", getConnectionClient.ProjectID)
	getConnectionPath = strings.ReplaceAll(getConnectionPath, "{id}", d.Id())

	getConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getConnectionResp, err := getConnectionClient.Request("GET", getConnectionPath, &getConnectionOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Connection")
	}

	getConnectionRespBody, err := utils.FlattenResponse(getConnectionResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getConnectionRespBody, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", getConnectionRespBody, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", getConnectionRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getConnectionRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getConnectionRespBody, nil)),
		d.Set("kafka_detail", flattenGetConnectionResponseBodyKafkaDetail(getConnectionRespBody)),
		d.Set("flavor", flattenGetConnectionResponseBodyKafkaFlavor(getConnectionRespBody)),
		d.Set("status", utils.PathSearch("status", getConnectionRespBody, nil)),
		d.Set("agency", utils.PathSearch("agency", getConnectionRespBody, nil)),
		d.Set("created_at", utils.PathSearch("created_time", getConnectionRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("updated_time", getConnectionRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetConnectionResponseBodyKafkaDetail(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("kafka_detail", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"instance_id":       utils.PathSearch("instance_id", curJson, nil),
			"connect_address":   utils.PathSearch("addr", curJson, nil),
			"user_name":         utils.PathSearch("username", curJson, nil),
			"password":          utils.PathSearch("password", curJson, nil),
			"acks":              utils.PathSearch("acks", curJson, nil),
			"security_protocol": utils.PathSearch("security_protocol", curJson, nil),
		},
	}
	return rst
}

func flattenGetConnectionResponseBodyKafkaFlavor(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("flavor", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"bandwidth_type":   utils.PathSearch("bandwidth_type", curJson, nil),
			"concurrency":      utils.PathSearch("concurrency", curJson, nil),
			"concurrency_type": utils.PathSearch("concurrency_type", curJson, nil),
			"name":             utils.PathSearch("name", curJson, nil),
		},
	}
	return rst
}

func resourceConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateConnectionChanges := []string{
		"description",
	}

	if d.HasChanges(updateConnectionChanges...) {
		// updateConnection: Update the EG Connection.
		var (
			updateConnectionHttpUrl = "v1/{project_id}/connections/{id}"
			updateConnectionProduct = "eg"
		)
		updateConnectionClient, err := cfg.NewServiceClient(updateConnectionProduct, region)
		if err != nil {
			return diag.Errorf("error creating EG client: %s", err)
		}

		updateConnectionPath := updateConnectionClient.Endpoint + updateConnectionHttpUrl
		updateConnectionPath = strings.ReplaceAll(updateConnectionPath, "{project_id}", updateConnectionClient.ProjectID)
		updateConnectionPath = strings.ReplaceAll(updateConnectionPath, "{id}", d.Id())

		updateConnectionOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{"Content-Type": "application/json"},
		}

		updateConnectionOpt.JSONBody = utils.RemoveNil(buildUpdateConnectionBodyParams(d))
		_, err = updateConnectionClient.Request("PUT", updateConnectionPath, &updateConnectionOpt)
		if err != nil {
			return diag.Errorf("error updating Connection: %s", err)
		}
	}
	return resourceConnectionRead(ctx, d, meta)
}

func buildUpdateConnectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"description": d.Get("description"),
	}
	return bodyParams
}

func resourceConnectionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteConnection: Delete an existing EG Connection
	var (
		deleteConnectionHttpUrl = "v1/{project_id}/connections/{id}"
		deleteConnectionProduct = "eg"
	)
	deleteConnectionClient, err := cfg.NewServiceClient(deleteConnectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	deleteConnectionPath := deleteConnectionClient.Endpoint + deleteConnectionHttpUrl
	deleteConnectionPath = strings.ReplaceAll(deleteConnectionPath, "{project_id}", deleteConnectionClient.ProjectID)
	deleteConnectionPath = strings.ReplaceAll(deleteConnectionPath, "{id}", d.Id())

	deleteConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	_, err = deleteConnectionClient.Request("DELETE", deleteConnectionPath, &deleteConnectionOpt)
	if err != nil {
		return diag.Errorf("error deleting Connection: %s", err)
	}

	return nil
}
