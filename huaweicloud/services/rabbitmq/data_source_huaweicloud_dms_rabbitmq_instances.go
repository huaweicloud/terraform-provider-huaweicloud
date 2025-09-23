package rabbitmq

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RabbitMQ GET /v2/{project_id}/instances
func DataSourceDmsRabbitMQInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDmsRabbitMQInstancesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the RabbitMQ instance.`,
			},
			"engine_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the version of the RabbitMQ engine.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the enterprise project ID to which the RabbitMQ instance belongs.`,
			},
			"exact_match_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies whether to search for the instance that precisely matches a specified instance name.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the flavor ID of the RabbitMQ instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the RabbitMQ instance.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the status of the RabbitMQ instance.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the RabbitMQ instance type.`,
			},
			"instances": {
				Type:        schema.TypeList,
				Elem:        rabbitmqInstanceSchema(),
				Computed:    true,
				Description: `Indicates the list of RabbitMQ instances.`,
			},
		},
	}
}

func rabbitmqInstanceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the RabbitMQ instance.`,
			},
			"access_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the user accessing the RabbitMQ instance.`,
			},
			"availability_zones": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `Indicates the list of the availability zone names.`,
			},
			"broker_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the number of the brokers.`,
			},
			"charging_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the billing mode.`,
			},
			"connect_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the IP address of the RabbitMQ instance.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time of the RabbitMQ instance.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the description of the RabbitMQ instance.`,
			},
			"engine": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the message engine type.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the enterprise project ID to which the RabbitMQ instance belongs.`,
			},
			"engine_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the version of the RabbitMQ engine.`,
			},
			"extend_times": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the number of disk expansion times. If the value exceeds 20, disk expansion is no longer allowed.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the flavor ID of the RabbitMQ instance.`,
			},
			"is_logical_volume": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: `Indicates whether the instance is a new instance. This parameter is used to distinguish old instances
				from new instances during instance capacity expansion.`,
			},
			"maintain_begin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time at which the maintenance window starts. The format is HH:mm:ss.`,
			},
			"maintain_end": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time at which the maintenance window ends. The format is HH:mm:ss.`,
			},
			"management_connect_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the management address of the RabbitMQ instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the RabbitMQ instance.`,
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the port of the RabbitMQ instance.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of a security group.`,
			},
			"security_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of a security group.`,
			},
			"specification": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the instance specification.`,
			},
			"ssl_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the RabbitMQ SASL_SSL is enabled.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the RabbitMQ instance.`,
			},
			"storage_resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the storage resource.`,
			},
			"storage_space": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the message storage space in GB.`,
			},
			"storage_spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the storage I/O specification.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of a subnet.`,
			},
			"tags": common.TagsComputedSchema(),
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the RabbitMQ instance type.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the user creating the RabbitMQ instance.`,
			},
			"used_storage_space": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the used message storage space in GB.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of a VPC.`,
			},
			"vpc_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of a VPC.`,
			},
		},
	}
	return &sc
}

func resourceDmsRabbitMQInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getRabbitmqInstances: Query the list of RabbitMQ instances
	var (
		getRabbitmqInstancesHttpUrl = "v2/{project_id}/instances"
		getRabbitmqInstancesProduct = "dms"
	)
	getRabbitmqInstancesClient, err := cfg.NewServiceClient(getRabbitmqInstancesProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS Client: %s", err)
	}

	getRabbitmqInstancesBasePath := getRabbitmqInstancesClient.Endpoint + getRabbitmqInstancesHttpUrl
	getRabbitmqInstancesBasePath = strings.ReplaceAll(getRabbitmqInstancesBasePath, "{project_id}", getRabbitmqInstancesClient.ProjectID)

	getRabbitmqInstancesOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	offset := 0
	totalResults := make([]map[string]interface{}, 0)
	var getRabbitmqInstancesPath string
	for {
		getRabbitmqInstancesQueryParams := buildGetRabbitmqInstancesQueryParams(d, offset)
		getRabbitmqInstancesPath = getRabbitmqInstancesBasePath + getRabbitmqInstancesQueryParams
		getRabbitmqInstancesResp, err := getRabbitmqInstancesClient.Request("GET", getRabbitmqInstancesPath, &getRabbitmqInstancesOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving RabbitMQ instances")
		}

		getRabbitmqInstancesRespBody, err := utils.FlattenResponse(getRabbitmqInstancesResp)
		if err != nil {
			return diag.FromErr(err)
		}

		instances := utils.PathSearch("instances", getRabbitmqInstancesRespBody, make([]interface{}, 0)).([]interface{})
		instanceNum := utils.PathSearch("instance_num", getRabbitmqInstancesRespBody, 0)
		filterResults := filterRabbitmqInstances(d, instances)
		for _, v := range filterResults {
			chargingMode := "postPaid"
			if utils.PathSearch("charging_mode", v, float64(0)).(float64) == 0 {
				chargingMode = "prePaid"
			}
			createdAt, _ := strconv.ParseInt(utils.PathSearch("created_at", v, "").(string), 10, 64)
			totalResults = append(totalResults, map[string]interface{}{
				"id":                         utils.PathSearch("instance_id", v, nil),
				"access_user":                utils.PathSearch("access_user", v, nil),
				"availability_zones":         utils.PathSearch("available_zone_names", v, nil),
				"broker_num":                 utils.PathSearch("broker_num", v, nil),
				"charging_mode":              chargingMode,
				"connect_address":            utils.PathSearch("connect_address", v, nil),
				"created_at":                 utils.FormatTimeStampRFC3339(createdAt/1000, false),
				"description":                utils.PathSearch("description", v, nil),
				"engine":                     utils.PathSearch("engine", v, nil),
				"engine_version":             utils.PathSearch("engine_version", v, nil),
				"extend_times":               utils.PathSearch("extend_times", v, nil),
				"enterprise_project_id":      utils.PathSearch("enterprise_project_id", v, nil),
				"flavor_id":                  utils.PathSearch("product_id", v, nil),
				"is_logical_volume":          utils.PathSearch("is_logical_volume", v, nil),
				"maintain_begin":             utils.PathSearch("maintain_begin", v, nil),
				"maintain_end":               utils.PathSearch("maintain_end", v, nil),
				"management_connect_address": utils.PathSearch("management_connect_address", v, nil),
				"name":                       utils.PathSearch("name", v, nil),
				"port":                       utils.PathSearch("port", v, nil),
				"security_group_id":          utils.PathSearch("security_group_id", v, nil),
				"security_group_name":        utils.PathSearch("security_group_name", v, nil),
				"specification":              utils.PathSearch("specification", v, nil),
				"ssl_enable":                 utils.PathSearch("ssl_enable", v, nil),
				"status":                     utils.PathSearch("status", v, nil),
				"storage_resource_id":        utils.PathSearch("storage_resource_id", v, nil),
				"storage_space":              utils.PathSearch("total_storage_space", v, nil),
				"storage_spec_code":          utils.PathSearch("storage_spec_code", v, nil),
				"subnet_id":                  utils.PathSearch("subnet_id", v, nil),
				"tags":                       utils.FlattenTagsToMap(utils.PathSearch("tags", v, make([]interface{}, 0))),
				"type":                       utils.PathSearch("type", v, nil),
				"user_name":                  utils.PathSearch("user_name", v, nil),
				"used_storage_space":         utils.PathSearch("used_storage_space", v, nil),
				"vpc_id":                     utils.PathSearch("vpc_id", v, nil),
				"vpc_name":                   utils.PathSearch("vpc_name", v, nil),
			})
		}

		offset += len(instances)
		if offset >= int(instanceNum.(float64)) {
			break
		}
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instances", totalResults),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func filterRabbitmqInstances(d *schema.ResourceData, instanceArray []interface{}) []interface{} {
	if len(instanceArray) < 1 {
		return nil
	}
	result := make([]interface{}, 0, len(instanceArray))

	rawEngineVersion, rawEngineVersionOK := d.GetOk("engine_version")
	rawFlavorId, rawFlavorIdOK := d.GetOk("flavor_id")
	rawType, rawTypeOK := d.GetOk("type")

	for _, instance := range instanceArray {
		flavorId := utils.PathSearch("product_id", instance, nil)
		instanceType := utils.PathSearch("type", instance, nil)
		engineVersion := utils.PathSearch("engine_version", instance, nil)

		if rawFlavorIdOK && rawFlavorId != flavorId {
			continue
		}

		if rawTypeOK && rawType != instanceType {
			continue
		}

		if rawEngineVersionOK && rawEngineVersion != engineVersion {
			continue
		}
		result = append(result, instance)
	}
	return result
}

func buildGetRabbitmqInstancesQueryParams(d *schema.ResourceData, offset int) string {
	res := fmt.Sprintf("?engine=%v&limit=50&offset=%v", "rabbitmq", offset)

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}

	if v, ok := d.GetOk("exact_match_name"); ok {
		res = fmt.Sprintf("%s&exact_match_name=%v", res, v)
	}

	if v, ok := d.GetOk("instance_id"); ok {
		res = fmt.Sprintf("%s&instance_id=%v", res, v)
	}

	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}

	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, v)
	}

	return res
}
