package kafka

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var smartConnectorValidateNonUpdatableParams = []string{
	"instance_id",
	"type",
	"task",
	"task.*.current_cluster_name",
	"task.*.cluster_name",
	"task.*.user_name",
	"task.*.password",
	"task.*.sasl_mechanism",
	"task.*.instance_id",
	"task.*.bootstrap_servers",
	"task.*.security_protocol",
	"task.*.direction",
	"task.*.sync_consumer_offsets_enabled",
	"task.*.replication_factor",
	"task.*.task_num",
	"task.*.rename_topic_enabled",
	"task.*.provenance_header_enabled",
	"task.*.consumer_strategy",
	"task.*.compression_type",
	"task.*.topics_mapping",
	"task.*.topics_mapping",
	"task.task_num",
	"task.*.rename_topic_enabled",
	"task.*.provenance_header_enabled",
	"task.*.consumer_strategy",
	"task.compression_type",
}

// @API Kafka POST /v2/{project_id}/instances/{instance_id}/connector/validate
func ResourceSmartConnectorValidate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSmartConnectorValidateCreate,
		ReadContext:   resourceSmartConnectorValidateRead,
		UpdateContext: resourceSmartConnectorValidateUpdate,
		DeleteContext: resourceSmartConnectorValidateDelete,

		CustomizeDiff: config.FlexibleForceNew(smartConnectorValidateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the Smart Connect is located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Kafka instance to which the Smart Connect belongs.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the Smart Connect task.`,
			},
			"task": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `The configuration of the Smart Connect task.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"current_cluster_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The alias of the current instance.`,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The ID of the peer instance.`,
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The alias of the peer instance.`,
						},
						"user_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The username of the peer instance.`,
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: `The password of the peer instance.`,
						},
						"sasl_mechanism": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The authentication mechanism of the peer instance.`,
						},
						"bootstrap_servers": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The address of the peer instance.`,
						},
						"security_protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The authentication method of the peer instance.`,
						},
						"direction": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The synchronization direction of the Smart Connect task.`,
						},
						"sync_consumer_offsets_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether to synchronize consumption progress.`,
						},
						"replication_factor": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `The number of replicas of the Smart Connect task.`,
						},
						"task_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `The number of tasks of the data replication.`,
						},
						"rename_topic_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether to rename topic.`,
						},
						"provenance_header_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether to add source header.`,
						},
						"consumer_strategy": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The startup offset of the smart connect task.`,
						},
						"compression_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The compression algorithm of the smart connect task.`,
						},
						"topics_mapping": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The topics mapping of the smart connect task.`,
						},
					},
				},
			},
			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildValidateConnectorValidateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"type": utils.ValueIgnoreEmpty(d.Get("type")),
		"task": buildValidateConnectorValidateTaskBodyParams(d.Get("task").([]interface{})),
	}

	return bodyParams
}

func buildValidateConnectorValidateTaskBodyParams(tasks []interface{}) map[string]interface{} {
	if len(tasks) == 0 || tasks[0] == nil {
		return nil
	}

	task := tasks[0]
	return map[string]interface{}{
		"current_cluster_name":          utils.ValueIgnoreEmpty(utils.PathSearch("current_cluster_name", task, nil)),
		"instance_id":                   utils.ValueIgnoreEmpty(utils.PathSearch("instance_id", task, nil)),
		"cluster_name":                  utils.ValueIgnoreEmpty(utils.PathSearch("cluster_name", task, nil)),
		"user_name":                     utils.ValueIgnoreEmpty(utils.PathSearch("user_name", task, nil)),
		"password":                      utils.ValueIgnoreEmpty(utils.PathSearch("password", task, nil)),
		"sasl_mechanism":                utils.ValueIgnoreEmpty(utils.PathSearch("sasl_mechanism", task, nil)),
		"bootstrap_servers":             utils.ValueIgnoreEmpty(utils.PathSearch("bootstrap_servers", task, nil)),
		"security_protocol":             utils.ValueIgnoreEmpty(utils.PathSearch("security_protocol", task, nil)),
		"direction":                     utils.ValueIgnoreEmpty(utils.PathSearch("direction", task, nil)),
		"sync_consumer_offsets_enabled": utils.ValueIgnoreEmpty(utils.PathSearch("sync_consumer_offsets_enabled", task, nil)),
		"replication_factor":            utils.ValueIgnoreEmpty(utils.PathSearch("replication_factor", task, nil)),
		"task_num":                      utils.ValueIgnoreEmpty(utils.PathSearch("task_num", task, nil)),
		"rename_topic_enabled":          utils.ValueIgnoreEmpty(utils.PathSearch("rename_topic_enabled", task, nil)),
		"provenance_header_enabled":     utils.ValueIgnoreEmpty(utils.PathSearch("provenance_header_enabled", task, nil)),
		"consumer_strategy":             utils.ValueIgnoreEmpty(utils.PathSearch("consumer_strategy", task, nil)),
		"compression_type":              utils.ValueIgnoreEmpty(utils.PathSearch("compression_type", task, nil)),
		"topics_mapping":                utils.ValueIgnoreEmpty(utils.PathSearch("topics_mapping", task, nil)),
	}
}

func resourceSmartConnectorValidateCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		httpUrl    = "v2/{project_id}/instances/{instance_id}/connector/validate"
		instanceId = d.Get("instance_id").(string)
	)
	client, err := cfg.NewServiceClient("dmsv2", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS Client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildValidateConnectorValidateBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error validating Smart Connect connectivity: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	return nil
}

func resourceSmartConnectorValidateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSmartConnectorValidateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSmartConnectorValidateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for validating Kafka instances connectivity. Deleting this
resource will not clear the corresponding request record, but will only remove the resource information from the tfstate
file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
