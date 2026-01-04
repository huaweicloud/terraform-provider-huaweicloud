package rabbitmq

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

//nolint:revive
// @API RabbitMQ DELETE /v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/exchanges/{exchange}/destination-type/{destination_type}/destination/{destination}/properties-key/{properties_key}/unbinding
// @API RabbitMQ POST /v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/exchanges/{exchange}/binding
// @API RabbitMQ GET /v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/exchanges/{exchange}/binding

func ResourceDmsRabbitmqExchangeAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRabbitmqExchangeAssociateCreate,
		ReadContext:   resourceDmsRabbitmqExchangeAssociateRead,
		DeleteContext: resourceDmsRabbitmqExchangeAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceExchangeAssociateImportState,
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
			"exchange": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"routing_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"properties_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDmsRabbitmqExchangeAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	vhost := d.Get("vhost").(string)
	exchange := d.Get("exchange").(string)
	destination := d.Get("destination").(string)
	destinationType := d.Get("destination_type").(string)

	createHttpUrl := "v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/exchanges/{exchange}/binding"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)
	createPath = strings.ReplaceAll(createPath, "{vhost}", vhost)
	createPath = strings.ReplaceAll(createPath, "{exchange}", exchange)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildRabbitmqExchangeAssociateRequestBody(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating exchange: %s", err)
	}

	id := fmt.Sprintf("%s/%s/%s/%s/%s", instanceID, vhost, exchange, destinationType, destination)
	if routingKey, ok := d.GetOk("routing_key"); ok {
		id += fmt.Sprintf("/%s", routingKey.(string))
	}
	d.SetId(id)

	return resourceDmsRabbitmqExchangeAssociateRead(ctx, d, cfg)
}

func buildRabbitmqExchangeAssociateRequestBody(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"destination":      d.Get("destination"),
		"destination_type": d.Get("destination_type"),
		// routing_key have to send empty string if it's empty
		"routing_key": d.Get("routing_key"),
	}
	return bodyParams
}

func resourceDmsRabbitmqExchangeAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	result, err := getRabbitmqExchangeAssociate(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving the exchange")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("routing_key", utils.PathSearch("routing_key", result, nil)),
		d.Set("properties_key", url.PathEscape(utils.PathSearch("properties_key", result, "").(string))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getRabbitmqExchangeAssociate(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	getHttpUrl := "v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/exchanges/{exchange}/binding"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath = strings.ReplaceAll(getPath, "{vhost}", d.Get("vhost").(string))
	getPath = strings.ReplaceAll(getPath, "{exchange}", d.Get("exchange").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	// Queue or queue are all available for destination_type when creating,
	// but the return will always be queue, so search it in lowercase
	searchPath := fmt.Sprintf("items[?destination_type=='%s']|[?destination=='%s']",
		strings.ToLower(d.Get("destination_type").(string)), d.Get("destination").(string))
	associations := utils.PathSearch(searchPath, getRespBody, make([]interface{}, 0)).([]interface{})

	routingKey := d.Get("routing_key").(string)

	// routing_key may contain `\`
	for _, association := range associations {
		if routingKey == utils.PathSearch("routing_key", association, "").(string) {
			return association, nil
		}
	}

	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Method:    "GET",
			URL:       "/v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/exchanges/{exchange}/binding",
			RequestId: "NONE",
			Body:      []byte(fmt.Sprintf("the exchange (%s) has no association", d.Get("exchange").(string))),
		},
	}
}

func resourceDmsRabbitmqExchangeAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	destination := d.Get("destination").(string)
	destination = strings.ReplaceAll(destination, "/", "__F_SLASH__")

	propertiesKey := d.Get("properties_key").(string)
	propertiesKey = strings.ReplaceAll(propertiesKey, "%2F", "__F_SLASH__")
	propertiesKey = strings.ReplaceAll(propertiesKey, "%5C", "__B_SLASH__")

	deleteHttpUrl := "v2/rabbitmq/{project_id}/instances/{instance_id}/vhosts/{vhost}/exchanges/{exchange}" +
		"/destination-type/{destination_type}/destination/{destination}/properties-key/{properties_key}/unbinding"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{vhost}", d.Get("vhost").(string))
	deletePath = strings.ReplaceAll(deletePath, "{exchange}", d.Get("exchange").(string))
	deletePath = strings.ReplaceAll(deletePath, "{destination_type}", d.Get("destination_type").(string))
	// destination may have % or |
	deletePath = strings.ReplaceAll(deletePath, "{destination}", url.PathEscape(destination))
	// properties_key is already encoded in READ
	deletePath = strings.ReplaceAll(deletePath, "{properties_key}", propertiesKey)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting exchange")
	}

	return nil
}

func resourceExchangeAssociateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ",")
	if !(len(parts) == 5 || len(parts) == 6) {
		return nil, fmt.Errorf("invalid ID format, must be <instance_id>,<vhost>,<exchange>,<destination_type>,<destination> or" +
			"<instance_id>,<vhost>,<exchange>,<destination_type>,<destination>,<routing_key>")
	}

	d.Set("instance_id", parts[0])
	d.Set("vhost", parts[1])
	d.Set("exchange", parts[2])
	d.Set("destination_type", parts[3])
	d.Set("destination", parts[4])

	id := fmt.Sprintf("%s/%s/%s/%s/%s", parts[0], parts[1], parts[2], parts[3], parts[4])

	if len(parts) == 6 {
		d.Set("routing_key", parts[5])
		id += fmt.Sprintf("/%s", parts[5])
	}

	// reset ID to be separated by slashes
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
