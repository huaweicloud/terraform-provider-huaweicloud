package iotda

import (
	"context"
	"log"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	stringRegxp     = `^[\x{4E00}-\x{9FFC}A-Za-z-_0-9?'#().,&%@!]*$`
	stringFormatMsg = "Only letters, Chinese characters, digits, hyphens (-), underscores (_) and" +
		" the following specail characters are allowed: ?'#().,&%@!"
)

func ResourceProduct() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceProductCreate,
		UpdateContext: ResourceProductUpdate,
		DeleteContext: ResourceProductDelete,
		ReadContext:   ResourceProductRead,
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
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 64),
					validation.StringMatch(regexp.MustCompile(stringRegxp), stringFormatMsg),
				),
			},

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{"MQTT", "CoAP", "HTTP", "HTTPS", "Modbus", "ONVIF",
					"OPC-UA", "OPC-DA", "Other"}, false),
			},

			"data_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"json", "binary"}, false),
			},

			"device_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 32),
					validation.StringMatch(regexp.MustCompile(stringRegxp), stringFormatMsg),
				),
			},

			"services": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 500,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(1, 64),
								validation.StringMatch(regexp.MustCompile(stringRegxp), stringFormatMsg),
							),
						},

						"type": { // keep same with console
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(0, 64),
								validation.StringMatch(regexp.MustCompile(stringRegxp), stringFormatMsg),
							),
						},

						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(0, 128),
								validation.StringMatch(regexp.MustCompile(stringRegxp), stringFormatMsg),
							),
						},

						"properties": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 500,
							Elem:     propertySchema("services.properties"),
						},

						"commands": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 500,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.All(
											validation.StringLenBetween(1, 64),
											validation.StringMatch(regexp.MustCompile(stringRegxp), stringFormatMsg),
										),
									},

									"paras": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 500,
										Elem:     propertySchema("services.commands.paras"),
									},

									"responses": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 500,
										Elem:     propertySchema("services.commands.responses"),
									},
								},
							},
						},
					},
				},
			},

			"manufacturer_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 32),
					validation.StringMatch(regexp.MustCompile(stringRegxp), stringFormatMsg),
				),
			},

			"industry": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 64),
					validation.StringMatch(regexp.MustCompile(stringRegxp), stringFormatMsg),
				),
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 128),
					validation.StringMatch(regexp.MustCompile(stringRegxp), stringFormatMsg),
				),
			},

			"space_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"product_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 32),
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z-_0-9]*$`),
						"Only letters, digits, underscores (_) and hyphens (-) are allowed."),
				),
			},
		},
	}
}

// propertySchema get the schema define for services.properties; services.commands.paras; services.commands.responses
func propertySchema(category string) *schema.Resource {
	common := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 64),
					validation.StringMatch(regexp.MustCompile(stringRegxp), stringFormatMsg),
				),
			},

			"type": { // keep same with console
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{"int", "decimal", "string", "DateTime",
					"jsonObject", "string list"}, false),
			},

			"enum_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"min": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0",
			},

			"max": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "65535",
			},

			"max_length": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 2147483647),
			},

			"step": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"unit": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 16),
				),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 128),
					validation.StringMatch(regexp.MustCompile(stringRegxp), stringFormatMsg),
				),
			},
		},
	}

	if category == "services.properties" {
		common.Schema["method"] = &schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"RW", "W", "R"}, false),
		}
	}
	return &common
}

func ResourceProductCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	createOpts := buildProductCreateParams(d)
	log.Printf("[DEBUG] Create IoTDA product params: %#v", createOpts)

	resp, err := client.CreateProduct(createOpts)
	if err != nil {
		return diag.Errorf("error creating IoTDA product: %s", err)
	}

	if resp.ProductId == nil {
		return diag.Errorf("error creating IoTDA product: id is not found in API response")
	}

	d.SetId(*resp.ProductId)
	return ResourceProductRead(ctx, d, meta)
}

func ResourceProductRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	response, err := client.ShowProduct(&model.ShowProductRequest{ProductId: d.Id()})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA product")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("product_id", response.ProductId),
		d.Set("name", response.Name),
		d.Set("device_type", response.DeviceType),
		d.Set("protocol", response.ProtocolType),
		d.Set("data_type", response.DataFormat),
		d.Set("manufacturer_name", response.ManufacturerName),
		d.Set("industry", response.Industry),
		d.Set("description", response.Description),
		d.Set("space_id", response.AppId),
		d.Set("services", flattenServices(response.ServiceCapabilities)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func ResourceProductUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	updateOpts := buildProductUpdateParams(d)
	_, err = client.UpdateProduct(updateOpts)

	if err != nil {
		return diag.Errorf("error updating IoTDA product: %s", err)
	}

	return ResourceProductRead(ctx, d, meta)
}

func ResourceProductDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	deleteOpts := &model.DeleteProductRequest{ProductId: d.Id()}
	_, err = client.DeleteProduct(deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting IoTDA product: %s", err)
	}

	return nil
}

func buildProductCreateParams(d *schema.ResourceData) *model.CreateProductRequest {
	req := model.CreateProductRequest{
		Body: &model.AddProduct{
			ProductId:           utils.StringIgnoreEmpty(d.Get("product_id").(string)),
			Name:                d.Get("name").(string),
			DeviceType:          d.Get("device_type").(string),
			ProtocolType:        d.Get("protocol").(string),
			DataFormat:          d.Get("data_type").(string),
			ManufacturerName:    utils.StringIgnoreEmpty(d.Get("manufacturer_name").(string)),
			Industry:            utils.StringIgnoreEmpty(d.Get("industry").(string)),
			Description:         utils.StringIgnoreEmpty(d.Get("description").(string)),
			AppId:               utils.StringIgnoreEmpty(d.Get("space_id").(string)),
			ServiceCapabilities: *buildServices(d.Get("services").([]interface{})),
		},
	}
	return &req
}

func buildProductUpdateParams(d *schema.ResourceData) *model.UpdateProductRequest {
	req := model.UpdateProductRequest{
		ProductId: d.Id(),
		Body: &model.UpdateProduct{
			Name:                utils.String(d.Get("name").(string)),
			DeviceType:          utils.String(d.Get("device_type").(string)),
			ProtocolType:        utils.String(d.Get("protocol").(string)),
			DataFormat:          utils.String(d.Get("data_type").(string)),
			ManufacturerName:    utils.StringIgnoreEmpty(d.Get("manufacturer_name").(string)),
			Industry:            utils.StringIgnoreEmpty(d.Get("industry").(string)),
			Description:         utils.StringIgnoreEmpty(d.Get("description").(string)),
			AppId:               utils.StringIgnoreEmpty(d.Get("space_id").(string)),
			ServiceCapabilities: buildServices(d.Get("services").([]interface{})),
		},
	}
	return &req
}

func buildServices(raw []interface{}) *[]model.ServiceCapability {
	rst := make([]model.ServiceCapability, len(raw))
	for i, v := range raw {
		service := v.(map[string]interface{})
		rst[i] = model.ServiceCapability{
			ServiceId:   service["id"].(string),
			ServiceType: service["type"].(string),
			Description: utils.StringIgnoreEmpty(service["description"].(string)),
			Properties:  buildServiceProperties(service["properties"]),
			Commands:    buildServiceCommands(service["commands"]),
		}
	}

	return &rst
}

func buildServiceCommands(commandsRaw interface{}) *[]model.ServiceCommand {
	properties := commandsRaw.([]interface{})
	rst := make([]model.ServiceCommand, len(properties))
	for i, v := range properties {
		s := v.(map[string]interface{})
		rst[i] = model.ServiceCommand{
			CommandName: s["name"].(string),
			Paras:       buildServiceCommandParas(s["paras"]),
			Responses:   buildServiceCommandResponses(s["responses"]),
		}
	}

	return &rst
}

func buildServiceCommandResponses(raw interface{}) *[]model.ServiceCommandResponse {
	rst := model.ServiceCommandResponse{
		ResponseName: "cmdResponses",
		Paras:        buildServiceCommandParas(raw),
	}
	return &[]model.ServiceCommandResponse{rst}
}

func buildServiceCommandParas(raw interface{}) *[]model.ServiceCommandPara {
	properties := raw.([]interface{})
	rst := make([]model.ServiceCommandPara, len(properties))
	for i, v := range properties {
		s := v.(map[string]interface{})
		rst[i] = model.ServiceCommandPara{
			ParaName:    s["name"].(string),
			DataType:    s["type"].(string),
			EnumList:    utils.ExpandToStringListPointer(s["enum_list"].([]interface{})),
			Min:         utils.String(s["min"].(string)),
			Max:         utils.String(s["max"].(string)),
			MaxLength:   utils.Int32(int32(s["max_length"].(int))),
			Step:        utils.Float64(s["step"].(float64)),
			Unit:        utils.StringIgnoreEmpty(s["unit"].(string)),
			Description: utils.StringIgnoreEmpty(s["description"].(string)),
		}
	}

	return &rst
}

func buildServiceProperties(propertiesRaw interface{}) *[]model.ServiceProperty {
	properties := propertiesRaw.([]interface{})
	rst := make([]model.ServiceProperty, len(properties))
	for i, v := range properties {
		s := v.(map[string]interface{})
		rst[i] = model.ServiceProperty{
			PropertyName: s["name"].(string),
			DataType:     s["type"].(string),
			EnumList:     utils.ExpandToStringListPointer(s["enum_list"].([]interface{})),
			Min:          utils.String(s["min"].(string)),
			Max:          utils.String(s["max"].(string)),
			MaxLength:    utils.Int32(int32(s["max_length"].(int))),
			Step:         utils.Float64(s["step"].(float64)),
			Unit:         utils.StringIgnoreEmpty(s["unit"].(string)),
			Description:  utils.StringIgnoreEmpty(s["description"].(string)),
			Method:       s["method"].(string),
		}
	}

	return &rst
}

func flattenServices(s *[]model.ServiceCapability) []interface{} {
	if s != nil {
		rst := make([]interface{}, len(*s))
		for i, v := range *s {
			rst[i] = map[string]interface{}{
				"id":          v.ServiceId,
				"type":        v.ServiceType,
				"description": v.Description,
				"properties":  flattenServiceProperties(v.Properties),
				"commands":    flattenServiceCommands(v.Commands),
			}
		}

		return rst
	}

	return make([]interface{}, 0)
}

func flattenServiceProperties(s *[]model.ServiceProperty) []interface{} {
	if s != nil {
		rst := make([]interface{}, len(*s))
		for i, v := range *s {
			rst[i] = map[string]interface{}{
				"name":        v.PropertyName,
				"type":        v.DataType,
				"enum_list":   v.EnumList,
				"min":         v.Min,
				"max":         v.Max,
				"max_length":  v.MaxLength,
				"step":        v.Step,
				"unit":        v.Unit,
				"description": v.Description,
				"method":      v.Method,
			}
		}

		return rst
	}

	return make([]interface{}, 0)
}

func flattenServiceCommands(s *[]model.ServiceCommand) []interface{} {
	if s != nil {
		rst := make([]interface{}, len(*s))
		for i, v := range *s {
			rst[i] = map[string]interface{}{
				"name":      v.CommandName,
				"paras":     flattenServiceCommandParas(v.Paras),
				"responses": flattenServiceCommandResponses(v.Responses),
			}
		}

		return rst
	}

	return make([]interface{}, 0)
}

func flattenServiceCommandResponses(s *[]model.ServiceCommandResponse) []interface{} {
	if s != nil {
		responses := *s
		if len(responses) > 0 {
			return flattenServiceCommandParas(responses[0].Paras)
		}
	}

	return make([]interface{}, 0)
}

func flattenServiceCommandParas(s *[]model.ServiceCommandPara) []interface{} {
	if s != nil {
		rst := make([]interface{}, len(*s))
		for i, v := range *s {
			rst[i] = map[string]interface{}{
				"name":        v.ParaName,
				"type":        v.DataType,
				"enum_list":   v.EnumList,
				"min":         v.Min,
				"max":         v.Max,
				"max_length":  v.MaxLength,
				"step":        v.Step,
				"unit":        v.Unit,
				"description": v.Description,
			}
		}
		return rst
	}

	return make([]interface{}, 0)
}
