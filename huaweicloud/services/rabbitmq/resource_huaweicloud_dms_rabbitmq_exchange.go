package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RabbitMQ PUT /v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/exchanges
// @API RabbitMQ POST /v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/exchanges
// @API RabbitMQ GET /v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/exchanges
// @API RabbitMQ GET /v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/exchanges/{exchange}/binding
func ResourceDmsRabbitmqExchange() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRabbitmqExchangeCreate,
		ReadContext:   resourceDmsRabbitmqExchangeRead,
		DeleteContext: resourceDmsRabbitmqExchangeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceExchangeOrQueueImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vhost": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"auto_delete": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"durable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"internal": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"arguments": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsJSON,
				Description:  `The argument configuration of the exchange, in JSON format.`,
			},
			"bindings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"routing_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"properties_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceDmsRabbitmqExchangeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	vhost := d.Get("vhost").(string)
	name := d.Get("name").(string)

	createHttpUrl := "v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/exchanges"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)
	createPath = strings.ReplaceAll(createPath, "{vhost}", vhost)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildRabbitmqExchangeRequestBody(d)),
	}
	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating exchange: %s", err)
	}

	id := fmt.Sprintf("%s/%s/%s", instanceID, vhost, name)
	d.SetId(id)

	return resourceDmsRabbitmqExchangeRead(ctx, d, cfg)
}

func buildRabbitmqExchangeRequestBody(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"type":        d.Get("type"),
		"auto_delete": d.Get("auto_delete"),
		"durable":     utils.ValueIgnoreEmpty(d.Get("durable")),
		"internal":    utils.ValueIgnoreEmpty(d.Get("internal")),
		"arguments":   utils.StringToJson(d.Get("arguments").(string)),
	}
	return bodyParams
}

func resourceDmsRabbitmqExchangeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	vhost := d.Get("vhost").(string)
	name := d.Get("name").(string)

	exchange, err := GetRabbitmqExchange(client, instanceID, vhost, name)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving the exchange")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("type", utils.PathSearch("type", exchange, nil)),
		d.Set("auto_delete", utils.PathSearch("auto_delete", exchange, nil)),
		d.Set("durable", utils.PathSearch("durable", exchange, nil)),
		d.Set("internal", utils.PathSearch("internal", exchange, nil)),
		d.Set("arguments", flattenExchangeArguments(utils.PathSearch("arguments",
			exchange, make(map[string]interface{})).(map[string]interface{}))),
	)

	listRabbitmqExchangeBindingsResp, err := listRabbitmqExchangeBindings(client, d)
	if err != nil {
		log.Printf("[WARN] Error fetching bindings of exchange (%s): %s", d.Id(), err)
	} else {
		mErr = multierror.Append(
			d.Set("bindings", flattenExchangeBindings(
				utils.PathSearch("items", listRabbitmqExchangeBindingsResp, make([]interface{}, 0)).([]interface{}))),
		)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetRabbitmqExchange(client *golangsdk.ServiceClient, instanceID, vhost, name string) (interface{}, error) {
	listHttpUrl := "v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/exchanges"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceID)
	listPath = strings.ReplaceAll(listPath, "{vhost}", vhost)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// pageLimit is `10`
	listPath += fmt.Sprintf("?limit=%d", pageLimit)
	offset := 0
	for {
		currentPath := listPath + fmt.Sprintf("&offset=%d", offset)
		listResp, err := client.Request("GET", currentPath, &listOpt)
		if err != nil {
			return nil, err
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return nil, err
		}

		searchPath := fmt.Sprintf("items[?name=='%s']|[0]", name)
		result := utils.PathSearch(searchPath, listRespBody, nil)
		if result != nil {
			return result, nil
		}

		// `total` means the number of all `exchange`, and type is float64.
		offset += pageLimit
		total := utils.PathSearch("total", listRespBody, float64(0))
		if int(total.(float64)) <= offset {
			return nil, golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Method:    "GET",
					URL:       "/v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/exchanges",
					RequestId: "NONE",
					Body:      []byte(fmt.Sprintf("the exchange (%s) does not exist", name)),
				},
			}
		}
	}
}

func listRabbitmqExchangeBindings(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	listHttpUrl := "v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/exchanges/{exchange}/binding"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath = strings.ReplaceAll(listPath, "{vhost}", d.Get("vhost").(string))

	// exchange name may have /, % or |
	name := d.Get("name").(string)
	if strings.Contains(name, "/") {
		replacedName := strings.ReplaceAll(name, "/", "__F_SLASH__")
		listPath = strings.ReplaceAll(listPath, "{exchange}", url.PathEscape(replacedName))
	} else {
		listPath = strings.ReplaceAll(listPath, "{exchange}", url.PathEscape(name))
	}

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving the exchange bingdings infos: %s", err)
	}

	return utils.FlattenResponse(listResp)
}

func flattenExchangeBindings(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		rst = append(rst, map[string]interface{}{
			"destination_type": utils.PathSearch("destination_type", params, nil),
			"destination":      utils.PathSearch("destination", params, nil),
			"routing_key":      utils.PathSearch("routing_key", params, nil),
			"properties_key":   utils.PathSearch("properties_key", params, nil),
		})
	}
	return rst
}

func flattenExchangeArguments(arguments map[string]interface{}) interface{} {
	// If the arguments is empty object, return nil.
	if len(arguments) == 0 {
		return nil
	}

	return utils.JsonToString(arguments)
}

func resourceDmsRabbitmqExchangeDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	vhost := d.Get("vhost").(string)
	name := d.Get("name").(string)

	deleteHttpUrl := "v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/exchanges"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceID)
	deletePath = strings.ReplaceAll(deletePath, "{vhost}", vhost)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		JSONBody: map[string]interface{}{
			"name": []string{name},
		},
	}

	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting exchange")
	}

	return nil
}
