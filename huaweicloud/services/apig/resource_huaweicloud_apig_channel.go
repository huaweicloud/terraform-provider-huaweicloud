package apig

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/channels"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/vpc-channels/{vpc_channel_id}
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/vpc-channels
func ResourceChannel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChannelCreate,
		ReadContext:   resourceChannelRead,
		UpdateContext: resourceChannelUpdate,
		DeleteContext: resourceChannelDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceChannelResourceImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the dedicated instance is located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the dedicated instance to which the channel belongs.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The channel name.",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The default port for health check in channel.",
			},
			"balance_strategy": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The distribution algorithm.",
			},
			"member_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The member type of the channel.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The type of the channel.",
				DiffSuppressFunc: func(_, o, n string, _ *schema.ResourceData) bool {
					return n == "2" && o == "builtin" || n == "3" && o == "microservice"
				},
			},
			"member_group": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the member group.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The description of the member group.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The weight of the current member group.",
						},
						"microservice_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The microservice version of the backend server group.",
						},
						"microservice_port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The microservice port of the backend server group.",
						},
						"microservice_labels": {
							Type:        schema.TypeMap,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The microservice tags of the backend server group.",
						},
						"reference_vpc_channel_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The ID of the reference load balance channel.",
						},
					},
				},
				Description: "The backend server groups of the channel.",
			},
			"member": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The IP address of the backend server.",
						},
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The ID of the backend server.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The name of the backend server.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The weight of current backend server.",
						},
						"is_backup": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether this member is the backup member.",
						},
						"group_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The group name of the backend server.",
						},
						"status": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The status of the backend server.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The port of the backend server.",
						},
					},
				},
				Description: "The backend servers of the channel.",
			},
			"health_check": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The rotocol for performing health check on backend servers.",
						},
						"threshold_normal": {
							Type:     schema.TypeInt,
							Required: true,
							Description: "The the healthy threshold, which refers to the number of consecutive successful " +
								"checks required for a backend server to be considered healthy.",
						},
						"threshold_abnormal": {
							Type:     schema.TypeInt,
							Required: true,
							Description: "The unhealthy threshold, which refers to the number of consecutive failed check " +
								"required for a backend server to be considered unhealthy.",
						},
						"interval": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The interval between consecutive check, in second.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The timeout for determining whether a health check fails, in second.",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The destination path for health check.",
						},
						"method": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The request method for health check.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The destination host port for health check.",
						},
						"http_codes": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The response codes for determining a successful HTTP response.",
						},
						"enable_client_ssl": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: `Whether to enable two-way authentication.`,
						},
						"status": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `The status of health check.`,
						},
					},
				},
				Description: "The health configuration of cloud servers associated with the load balance channel for " +
					"APIG regularly check.",
			},
			"microservice": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cce_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The ID of the CCE cluster.",
									},
									"namespace": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name of the CCE namespace.",
									},
									"workload_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The workload type.",
									},
									"workload_name": {
										Type:     schema.TypeString,
										Optional: true,
										Description: utils.SchemaDesc(
											`The workload name.`,
											utils.SchemaDescInput{
												Deprecated: true,
											}),
									},
									"label_key": {
										Type:     schema.TypeString,
										Optional: true,
										Description: utils.SchemaDesc(
											`The service label key.`,
											utils.SchemaDescInput{
												Required: true,
											}),
									},
									"label_value": {
										Type:     schema.TypeString,
										Optional: true,
										Description: utils.SchemaDesc(
											`The service label value.`,
											utils.SchemaDescInput{
												Required: true,
											}),
									},
								},
							},
							Description:  "The CCE microservice details.",
							ExactlyOneOf: []string{"microservice.0.cse_config"},
						},
						"cse_config": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"engine_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "schema:Internal; The microservice engine ID.",
									},
									"service_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "schema:Internal; The microservice ID.",
									},
								},
							},
							Description: "schema:Internal; The CSE microservice details.",
						},
					},
				},
				Description: "The configuration of the microservice.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the channel.",
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The current status of the channel.`,
			},
		},
	}
}

func buildMicroserviceLabels(labels map[string]interface{}) []channels.MicroserviceLabel {
	result := make([]channels.MicroserviceLabel, 0, len(labels))
	for k, v := range labels {
		result = append(result, channels.MicroserviceLabel{
			Name:  k,
			Value: v.(string),
		})
	}
	return result
}

func buildChannelMemberGroups(groups *schema.Set) []channels.MemberGroup {
	if groups.Len() < 1 {
		return nil
	}

	result := make([]channels.MemberGroup, groups.Len())
	for i, val := range groups.List() {
		group := val.(map[string]interface{})
		result[i] = channels.MemberGroup{
			Name:                  group["name"].(string),
			Description:           group["description"].(string),
			Weight:                group["weight"].(int),
			MicroserviceVersion:   group["microservice_version"].(string),
			MicroservicePort:      group["microservice_port"].(int),
			MicroserviceLabels:    buildMicroserviceLabels(group["microservice_labels"].(map[string]interface{})),
			ReferenceVpcChannelId: group["reference_vpc_channel_id"].(string),
		}
	}

	return result
}

func buildChannelMembers(members *schema.Set) []channels.MemberInfo {
	if members.Len() < 1 {
		return nil
	}

	result := make([]channels.MemberInfo, members.Len())
	for i, val := range members.List() {
		member := val.(map[string]interface{})
		result[i] = channels.MemberInfo{
			Host:      member["host"].(string),
			EcsId:     member["id"].(string),
			EcsName:   member["name"].(string),
			Weight:    utils.Int(member["weight"].(int)),
			IsBackup:  utils.Bool(member["is_backup"].(bool)),
			GroupName: member["group_name"].(string),
			Status:    member["status"].(int),
			Port:      utils.Int(member["port"].(int)),
		}
	}

	return result
}

func buildChannelHealthCheckConfig(healthConfigs []interface{}) *channels.VpcHealthConfig {
	if len(healthConfigs) < 1 {
		return nil
	}

	// The `health_check` is a `Computed` and `Optional` behavior. If `health_check.protocol` is empty, it means that
	// the `health_check` is not specified in the script, and `health_check.protocol` is defined as required in the SDK,
	// so `health_check` should be ignored if empty, otherwise, the SDK will report an error during updates.
	healthConfig := healthConfigs[0].(map[string]interface{})
	protocol := healthConfig["protocol"].(string)
	if protocol == "" {
		return nil
	}

	return &channels.VpcHealthConfig{
		Protocol:          protocol,
		ThresholdNormal:   healthConfig["threshold_normal"].(int),
		ThresholdAbnormal: healthConfig["threshold_abnormal"].(int),
		TimeInterval:      healthConfig["interval"].(int),
		Path:              healthConfig["path"].(string),
		Method:            healthConfig["method"].(string),
		Port:              healthConfig["port"].(int),
		HttpCodes:         healthConfig["http_codes"].(string),
		EnableClientSsl:   utils.Bool(healthConfig["enable_client_ssl"].(bool)),
		Status:            healthConfig["status"].(int),
		Timeout:           healthConfig["timeout"].(int),
	}
}

func buildChannelMicroserviceConfig(microserviceConfigs []interface{}) *channels.MicroserviceConfig {
	if len(microserviceConfigs) < 1 {
		return nil
	}

	micConfig := microserviceConfigs[0].(map[string]interface{})
	if cceConfig, ok := micConfig["cce_config"]; ok {
		log.Printf("[DEBUG] The CCE configuration of the microservice is: %#v", cceConfig)
		configs := cceConfig.([]interface{})
		if len(configs) > 0 {
			details := configs[0].(map[string]interface{})
			return &channels.MicroserviceConfig{
				ServiceType: "CCE",
				CceInfo: &channels.MicroserviceCceInfo{
					ClusterId:    details["cluster_id"].(string),
					Namespace:    details["namespace"].(string),
					WorkloadType: details["workload_type"].(string),
					AppName:      details["workload_name"].(string),
					LabelKey:     details["label_key"].(string),
					LabelValue:   details["label_value"].(string),
				},
			}
		}
	}
	if cseConfig, ok := micConfig["cse_config"]; ok {
		log.Printf("[DEBUG] The CSE configuration of the microservice is: %#v", cseConfig)
		configs := cseConfig.([]interface{})
		if len(configs) > 0 {
			details := configs[0].(map[string]interface{})
			return &channels.MicroserviceConfig{
				ServiceType: "CSE",
				CseInfo: &channels.MicroserviceCseInfo{
					EngineId:  details["engine_id"].(string),
					ServiceId: details["service_id"].(string),
				},
			}
		}
	}
	return nil
}

func buildChannelCreateOpts(d *schema.ResourceData) channels.ChannelOpts {
	result := channels.ChannelOpts{
		InstanceId:         d.Get("instance_id").(string),
		Name:               d.Get("name").(string),
		Port:               d.Get("port").(int),
		BalanceStrategy:    d.Get("balance_strategy").(int),
		MemberType:         d.Get("member_type").(string),
		MemberGroups:       buildChannelMemberGroups(d.Get("member_group").(*schema.Set)),
		Members:            buildChannelMembers(d.Get("member").(*schema.Set)),
		VpcHealthConfig:    buildChannelHealthCheckConfig(d.Get("health_check").([]interface{})),
		MicroserviceConfig: buildChannelMicroserviceConfig(d.Get("microservice").([]interface{})),
	}

	cType := d.Get("type").(string)
	switch cType {
	// Due the type conversion of the terraform provider, the number can be convert to the string without errors.
	case "2":
		result.Type = 2
	case "3":
		result.Type = 3
	default:
		result.VpcChannelType = cType
	}

	return result
}

func resourceChannelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	opts := buildChannelCreateOpts(d)
	v, err := channels.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating APIG channel: %s", err)
	}
	d.SetId(v.ID)
	return resourceChannelRead(ctx, d, meta)
}

func flattenMicroserviceLabels(labels []channels.MicroserviceLabel) map[string]interface{} {
	result := make(map[string]interface{})
	for _, label := range labels {
		result[label.Name] = label.Value
	}
	return result
}

func flattenChannelMemberGroups(groups []channels.MemberGroup) []map[string]interface{} {
	result := make([]map[string]interface{}, len(groups))
	for i, v := range groups {
		result[i] = map[string]interface{}{
			"name":                     v.Name,
			"description":              v.Description,
			"weight":                   v.Weight,
			"microservice_version":     v.MicroserviceVersion,
			"microservice_port":        v.MicroservicePort,
			"microservice_labels":      flattenMicroserviceLabels(v.MicroserviceLabels),
			"reference_vpc_channel_id": v.ReferenceVpcChannelId,
		}
	}
	return result
}

func flattenChannelMicroserivceCceConfig(cceConfig *channels.MicroserviceCceInfo) []map[string]interface{} {
	if cceConfig == nil {
		return nil
	}
	result := []map[string]interface{}{
		{
			"cluster_id":    cceConfig.ClusterId,
			"namespace":     cceConfig.Namespace,
			"workload_type": cceConfig.WorkloadType,
			"workload_name": cceConfig.AppName,
			"label_key":     cceConfig.LabelKey,
			"label_value":   cceConfig.LabelValue,
		},
	}
	return result
}

func flattenChannelMicroserivceCseConfig(cseConfig *channels.MicroserviceCseInfo) []map[string]interface{} {
	if cseConfig == nil {
		return nil
	}
	result := []map[string]interface{}{
		{
			"engine_id":  cseConfig.EngineId,
			"service_id": cseConfig.ServiceId,
		},
	}
	return result
}

func flattenHealthCheckConfig(healthConfig *channels.VpcHealthConfig) []map[string]interface{} {
	if healthConfig == nil {
		return nil
	}
	result := []map[string]interface{}{
		{
			"protocol":           strings.ToUpper(healthConfig.Protocol),
			"threshold_normal":   healthConfig.ThresholdNormal,
			"threshold_abnormal": healthConfig.ThresholdAbnormal,
			"interval":           healthConfig.TimeInterval,
			"timeout":            healthConfig.Timeout,
			"path":               healthConfig.Path,
			"method":             healthConfig.Method,
			"port":               healthConfig.Port,
			"http_codes":         healthConfig.HttpCodes,
			"enable_client_ssl":  healthConfig.EnableClientSsl,
			"status":             healthConfig.Status,
		},
	}
	return result
}

func flattenChannelMicroserivceConfig(microserviceConfig *channels.MicroserviceConfig) []map[string]interface{} {
	if microserviceConfig == nil {
		return nil
	}
	result := make([]map[string]interface{}, 0, 1)
	switch microserviceConfig.ServiceType {
	case "CCE":
		result = append(result, map[string]interface{}{
			"cce_config": flattenChannelMicroserivceCceConfig(microserviceConfig.CceInfo),
		})
	case "CSE":
		result = append(result, map[string]interface{}{
			"cse_config": flattenChannelMicroserivceCseConfig(microserviceConfig.CseInfo),
		})
	}
	return result
}

func flattenChannelMembers(members []channels.MemberInfo) []map[string]interface{} {
	if len(members) < 1 {
		return nil
	}
	result := make([]map[string]interface{}, len(members))
	for i, member := range members {
		result[i] = map[string]interface{}{
			"host":       member.Host,
			"id":         member.EcsId,
			"name":       member.EcsName,
			"weight":     member.Weight,
			"is_backup":  member.IsBackup,
			"group_name": member.GroupName,
			"status":     member.Status,
			"port":       member.Port,
		}
	}
	return result
}

func parseChannelType(resp *channels.Channel) string {
	if resp.VpcChannelType != "" {
		return resp.VpcChannelType
	}
	return strconv.Itoa(resp.Type)
}

func resourceChannelRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	channelId := d.Id()
	resp, err := channels.Get(client, instanceId, channelId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "channel")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", resp.Name),
		d.Set("port", resp.Port),
		d.Set("balance_strategy", resp.BalanceStrategy),
		d.Set("member_type", resp.MemberType),
		d.Set("type", parseChannelType(resp)),
		d.Set("member_group", flattenChannelMemberGroups(resp.MemberGroups)),
		d.Set("member", flattenChannelMembers(resp.Members)),
		d.Set("health_check", flattenHealthCheckConfig(resp.VpcHealthConfig)),
		d.Set("microservice", flattenChannelMicroserivceConfig(resp.MicroserviceConfig)),
		d.Set("created_at", resp.CreateTime),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving APIG channel (%s) fields: %s", channelId, mErr)
	}
	return nil
}

func buildChannelUpdateOpts(d *schema.ResourceData) channels.ChannelOpts {
	return buildChannelCreateOpts(d)
}

func resourceChannelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		channelId = d.Id()
		opts      = buildChannelUpdateOpts(d)
	)
	_, err = channels.Update(client, channelId, opts)
	if err != nil {
		return diag.Errorf("error updating APIG channel (%s): %s", channelId, err)
	}
	return resourceChannelRead(ctx, d, meta)
}

func resourceChannelDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	channelId := d.Id()
	if err = channels.Delete(client, instanceId, channelId); err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting APIG channel (%s)", channelId))
	}

	return nil
}

func resourceChannelResourceImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	channelId := d.Id()
	parts := strings.Split(channelId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but got '%s'",
			channelId)
	}

	d.SetId(parts[1])
	err := d.Set("instance_id", parts[0])
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error saving instance ID: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
