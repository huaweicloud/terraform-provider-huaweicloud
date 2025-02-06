package iotda

import (
	"context"
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

// @API IoTDA PUT /v5/iot/{project_id}/products/{product_id}
// @API IoTDA DELETE /v5/iot/{project_id}/products/{product_id}
// @API IoTDA GET /v5/iot/{project_id}/products/{product_id}
// @API IoTDA POST /v5/iot/{project_id}/products
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
						},
						"type": { // keep same with console
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"option": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
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
			},
			"industry": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
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
			},
		},
	}
}

// propertySchema get the schema define for services.properties; services.commands.paras; services.commands.responses
func propertySchema(category string) *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": { // keep same with console
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{"int", "decimal", "string", "DateTime",
					"jsonObject", "string list"}, false),
			},
			"required": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
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
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"default_value": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}

	if category == "services.properties" {
		sc.Schema["method"] = &schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"RW", "W", "R"}, false),
		}
	}
	return &sc
}

func buildServiceCapabilitiesPropertiesBodyParams(propertiesRaw interface{}) []map[string]interface{} {
	properties := propertiesRaw.([]interface{})
	rst := make([]map[string]interface{}, len(properties))
	for i, v := range properties {
		s := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"property_name": s["name"],
			"data_type":     s["type"],
			"required":      s["required"],
			"enum_list":     s["enum_list"],
			"min":           s["min"],
			"max":           s["max"],
			"max_length":    s["max_length"],
			"step":          s["step"],
			"unit":          utils.ValueIgnoreEmpty(s["unit"]),
			"description":   utils.ValueIgnoreEmpty(s["description"]),
			"method":        s["method"],
			"default_value": utils.StringToJson(s["default_value"].(string)),
		}
	}

	return rst
}

func buildCommandsParasBodyParams(rawParas interface{}) []map[string]interface{} {
	paras := rawParas.([]interface{})
	rst := make([]map[string]interface{}, len(paras))
	for i, v := range paras {
		s := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"para_name":   s["name"],
			"data_type":   s["type"],
			"required":    s["required"],
			"enum_list":   s["enum_list"],
			"min":         s["min"],
			"max":         s["max"],
			"max_length":  s["max_length"],
			"step":        s["step"],
			"unit":        utils.ValueIgnoreEmpty(s["unit"]),
			"description": utils.ValueIgnoreEmpty(s["description"]),
		}
	}

	return rst
}

func buildCommandsResponsesParasBodyParams(rawParas interface{}) []map[string]interface{} {
	paras := rawParas.([]interface{})
	rst := make([]map[string]interface{}, len(paras))
	for i, v := range paras {
		s := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"para_name":   s["name"],
			"data_type":   s["type"],
			"required":    s["required"],
			"enum_list":   s["enum_list"],
			"min":         s["min"],
			"max":         s["max"],
			"max_length":  s["max_length"],
			"step":        s["step"],
			"unit":        utils.ValueIgnoreEmpty(s["unit"]),
			"description": utils.ValueIgnoreEmpty(s["description"]),
		}
	}

	return rst
}

func buildCommandsResponsesBodyParams(responsesRaw interface{}) []map[string]interface{} {
	rst := map[string]interface{}{
		"response_name": "cmdResponses",
		"paras":         buildCommandsResponsesParasBodyParams(responsesRaw),
	}

	return []map[string]interface{}{rst}
}

func buildServiceCapabilitiesCommandsBodyParams(commandsRaw interface{}) []map[string]interface{} {
	commands := commandsRaw.([]interface{})
	rst := make([]map[string]interface{}, len(commands))
	for i, v := range commands {
		s := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"command_name": s["name"],
			"paras":        buildCommandsParasBodyParams(s["paras"]),
			"responses":    buildCommandsResponsesBodyParams(s["responses"]),
		}
	}

	return rst
}

func buildProductServiceCapabilitiesBodyParams(raw []interface{}) []map[string]interface{} {
	rst := make([]map[string]interface{}, len(raw))
	for i, v := range raw {
		service := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"service_id":   service["id"],
			"service_type": service["type"],
			"description":  utils.ValueIgnoreEmpty(service["description"]),
			"option":       utils.ValueIgnoreEmpty(service["option"]),
			"properties":   buildServiceCapabilitiesPropertiesBodyParams(service["properties"]),
			"commands":     buildServiceCapabilitiesCommandsBodyParams(service["commands"]),
		}
	}

	return rst
}

func buildCreateProductBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"product_id":           utils.ValueIgnoreEmpty(d.Get("product_id")),
		"name":                 d.Get("name"),
		"device_type":          d.Get("device_type"),
		"protocol_type":        d.Get("protocol"),
		"data_format":          d.Get("data_type"),
		"manufacturer_name":    utils.ValueIgnoreEmpty(d.Get("manufacturer_name")),
		"industry":             utils.ValueIgnoreEmpty(d.Get("industry")),
		"description":          utils.ValueIgnoreEmpty(d.Get("description")),
		"app_id":               utils.ValueIgnoreEmpty(d.Get("space_id")),
		"service_capabilities": buildProductServiceCapabilitiesBodyParams(d.Get("services").([]interface{})),
	}
}

func ResourceProductCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/iot/{project_id}/products"
		product = "iotda"
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateProductBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating IoTDA product: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	productID := utils.PathSearch("product_id", respBody, "").(string)
	if productID == "" {
		return diag.Errorf("error creating IoTDA product: ID is not found in API response")
	}

	d.SetId(productID)
	return ResourceProductRead(ctx, d, meta)
}

func flattenServicesPropertiesAttribute(respBody interface{}) []interface{} {
	properties := utils.PathSearch("properties", respBody, make([]interface{}, 0)).([]interface{})
	rst := make([]interface{}, len(properties))
	for i, v := range properties {
		rst[i] = map[string]interface{}{
			"name":          utils.PathSearch("property_name", v, nil),
			"type":          utils.PathSearch("data_type", v, nil),
			"required":      utils.PathSearch("required", v, nil),
			"enum_list":     utils.PathSearch("enum_list", v, nil),
			"min":           utils.PathSearch("min", v, nil),
			"max":           utils.PathSearch("max", v, nil),
			"max_length":    utils.PathSearch("max_length", v, nil),
			"step":          utils.PathSearch("step", v, nil),
			"unit":          utils.PathSearch("unit", v, nil),
			"description":   utils.PathSearch("description", v, nil),
			"method":        utils.PathSearch("method", v, nil),
			"default_value": utils.JsonToString(utils.PathSearch("default_value", v, nil)),
		}
	}

	return rst
}

func flattenCommandsParasAttribute(respBody interface{}) []interface{} {
	rawArray := utils.PathSearch("paras", respBody, make([]interface{}, 0)).([]interface{})
	rst := make([]interface{}, len(rawArray))
	for i, v := range rawArray {
		rst[i] = map[string]interface{}{
			"name":        utils.PathSearch("para_name", v, nil),
			"type":        utils.PathSearch("data_type", v, nil),
			"required":    utils.PathSearch("required", v, nil),
			"enum_list":   utils.PathSearch("enum_list", v, nil),
			"min":         utils.PathSearch("min", v, nil),
			"max":         utils.PathSearch("max", v, nil),
			"max_length":  utils.PathSearch("max_length", v, nil),
			"step":        utils.PathSearch("step", v, nil),
			"unit":        utils.PathSearch("unit", v, nil),
			"description": utils.PathSearch("description", v, nil),
		}
	}

	return rst
}

func flattenCommandsResponsesAttribute(respBody interface{}) []interface{} {
	rawArray := utils.PathSearch("responses|[0].paras", respBody, make([]interface{}, 0)).([]interface{})
	rst := make([]interface{}, len(rawArray))
	for i, v := range rawArray {
		rst[i] = map[string]interface{}{
			"name":        utils.PathSearch("para_name", v, nil),
			"type":        utils.PathSearch("data_type", v, nil),
			"required":    utils.PathSearch("required", v, nil),
			"enum_list":   utils.PathSearch("enum_list", v, nil),
			"min":         utils.PathSearch("min", v, nil),
			"max":         utils.PathSearch("max", v, nil),
			"max_length":  utils.PathSearch("max_length", v, nil),
			"step":        utils.PathSearch("step", v, nil),
			"unit":        utils.PathSearch("unit", v, nil),
			"description": utils.PathSearch("description", v, nil),
		}
	}

	return rst
}

func flattenServicesCommandsAttribute(respBody interface{}) []interface{} {
	rawArray := utils.PathSearch("commands", respBody, make([]interface{}, 0)).([]interface{})
	rst := make([]interface{}, len(rawArray))
	for i, v := range rawArray {
		rst[i] = map[string]interface{}{
			"name":      utils.PathSearch("command_name", v, nil),
			"paras":     flattenCommandsParasAttribute(v),
			"responses": flattenCommandsResponsesAttribute(v),
		}
	}

	return rst
}

func flattenServicesAttribute(respBody interface{}) []interface{} {
	rawArray := utils.PathSearch("service_capabilities", respBody, make([]interface{}, 0)).([]interface{})
	rst := make([]interface{}, len(rawArray))
	for i, v := range rawArray {
		rst[i] = map[string]interface{}{
			"id":          utils.PathSearch("service_id", v, nil),
			"type":        utils.PathSearch("service_type", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"option":      utils.PathSearch("option", v, nil),
			"properties":  flattenServicesPropertiesAttribute(v),
			"commands":    flattenServicesCommandsAttribute(v),
		}
	}

	return rst
}

func ResourceProductRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/iot/{project_id}/products/{product_id}"
		product = "iotda"
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{product_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		// When the resource does not exist, the query API will return a `404` status code.
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA product")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("product_id", utils.PathSearch("product_id", respBody, nil)),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("device_type", utils.PathSearch("device_type", respBody, nil)),
		d.Set("protocol", utils.PathSearch("protocol_type", respBody, nil)),
		d.Set("data_type", utils.PathSearch("data_format", respBody, nil)),
		d.Set("manufacturer_name", utils.PathSearch("manufacturer_name", respBody, nil)),
		d.Set("industry", utils.PathSearch("industry", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("space_id", utils.PathSearch("app_id", respBody, nil)),
		d.Set("services", flattenServicesAttribute(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdatePropertiesBodyParams(propertiesRaw interface{}) []interface{} {
	properties := propertiesRaw.([]interface{})
	rst := make([]interface{}, len(properties))
	for i, v := range properties {
		s := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"property_name": s["name"],
			"data_type":     s["type"],
			"required":      s["required"],
			"enum_list":     s["enum_list"],
			"min":           s["min"],
			"max":           s["max"],
			"max_length":    s["max_length"],
			"step":          s["step"],
			"unit":          utils.ValueIgnoreEmpty(s["unit"]),
			"description":   utils.ValueIgnoreEmpty(s["description"]),
			"method":        s["method"],
			"default_value": utils.StringToJson(s["default_value"].(string)),
		}
	}

	return rst
}

func buildUpdateCommandsParasBodyParams(rawParas interface{}) []interface{} {
	paras := rawParas.([]interface{})
	rst := make([]interface{}, len(paras))
	for i, v := range paras {
		s := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"para_name":   s["name"],
			"data_type":   s["type"],
			"required":    s["required"],
			"enum_list":   s["enum_list"],
			"min":         s["min"],
			"max":         s["max"],
			"max_length":  s["max_length"],
			"step":        s["step"],
			"unit":        utils.ValueIgnoreEmpty(s["unit"]),
			"description": utils.ValueIgnoreEmpty(s["description"]),
		}
	}

	return rst
}

func buildUpdateCommandsResponsesParasBodyParams(rawResponses interface{}) []interface{} {
	paras := rawResponses.([]interface{})
	rst := make([]interface{}, len(paras))
	for i, v := range paras {
		s := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"para_name":   s["name"],
			"data_type":   s["type"],
			"required":    s["required"],
			"enum_list":   s["enum_list"],
			"min":         s["min"],
			"max":         s["max"],
			"max_length":  s["max_length"],
			"step":        s["step"],
			"unit":        utils.ValueIgnoreEmpty(s["unit"]),
			"description": utils.ValueIgnoreEmpty(s["description"]),
		}
	}

	return rst
}

func buildUpdateCommandsResponsesBodyParams(rawResponses interface{}) []interface{} {
	rst := map[string]interface{}{
		"response_name": "cmdResponses",
		"paras":         buildUpdateCommandsResponsesParasBodyParams(rawResponses),
	}

	return []interface{}{rst}
}

func buildUpdateCommandsBodyParams(commandsRaw interface{}) []interface{} {
	commands := commandsRaw.([]interface{})
	rst := make([]interface{}, len(commands))
	for i, v := range commands {
		s := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"command_name": s["name"],
			"paras":        buildUpdateCommandsParasBodyParams(s["paras"]),
			"responses":    buildUpdateCommandsResponsesBodyParams(s["responses"]),
		}
	}

	return rst
}

func buildUpdateServiceCapabilitiesBodyParams(servicesRaw []interface{}) []interface{} {
	rst := make([]interface{}, len(servicesRaw))
	for i, v := range servicesRaw {
		service := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"service_id":   service["id"],
			"service_type": service["type"],
			"description":  utils.ValueIgnoreEmpty(service["description"]),
			"option":       utils.ValueIgnoreEmpty(service["option"]),
			"properties":   buildUpdatePropertiesBodyParams(service["properties"]),
			"commands":     buildUpdateCommandsBodyParams(service["commands"]),
		}
	}

	return rst
}

func buildUpdateProductBodyParams(d *schema.ResourceData) map[string]interface{} {
	rst := map[string]interface{}{
		"name":                 d.Get("name"),
		"device_type":          d.Get("device_type"),
		"protocol_type":        d.Get("protocol"),
		"data_format":          d.Get("data_type"),
		"manufacturer_name":    utils.ValueIgnoreEmpty(d.Get("manufacturer_name")),
		"industry":             utils.ValueIgnoreEmpty(d.Get("industry")),
		"description":          utils.ValueIgnoreEmpty(d.Get("description")),
		"app_id":               utils.ValueIgnoreEmpty(d.Get("space_id")),
		"service_capabilities": buildUpdateServiceCapabilitiesBodyParams(d.Get("services").([]interface{})),
	}

	return rst
}

func ResourceProductUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/iot/{project_id}/products/{product_id}"
		product = "iotda"
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{product_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateProductBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating IoTDA product: %s", err)
	}

	return ResourceProductRead(ctx, d, meta)
}

func ResourceProductDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/iot/{project_id}/products/{product_id}"
		product = "iotda"
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{product_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		// When the resource does not exist, the deletion API will return a `404` status code.
		return common.CheckDeletedDiag(d, err, "error deleting IoTDA product")
	}

	return nil
}
