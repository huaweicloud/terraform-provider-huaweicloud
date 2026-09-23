package drs

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var drsAssociateSmnTopicNonUpdatableParams = []string{
	"job_id",
	"topic_urn",
	"is_alarm_to_user",
	"increment_delay_threshold",
	"rto_delay_threshold",
	"rpo_delay_threshold",
}

// @API DRS POST /v5/{project_id}/jobs/{job_id}/associate-smn-topic
func ResourceDrsAssociateSmnTopic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDrsAssociateSmnTopicCreate,
		ReadContext:   resourceDrsAssociateSmnTopicRead,
		UpdateContext: resourceDrsAssociateSmnTopicUpdate,
		DeleteContext: resourceDrsAssociateSmnTopicDelete,

		CustomizeDiff: config.FlexibleForceNew(drsAssociateSmnTopicNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the ID of the DRS job to associate the SMN topic.",
			},
			"topic_urn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the topic URN of the SMN topic.",
			},
			"is_alarm_to_user": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies whether to send alarm notifications to the user.",
			},
			"increment_delay_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the incremental delay threshold in seconds.",
			},
			"rto_delay_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the RTO delay threshold in seconds, only available in disaster recovery scenarios.",
			},
			"rpo_delay_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the RPO delay threshold in seconds, only available in disaster recovery scenarios.",
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceDrsAssociateSmnTopicCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v5/{project_id}/jobs/{job_id}/associate-smn-topic"
		product = "drs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{job_id}", d.Get("job_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateDrsAssociateSmnTopicBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error associating DRS SMN topic: %s", err)
	}

	_, err = utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// The API returns an empty response body, use job_id as the resource ID
	d.SetId(d.Get("job_id").(string))

	return resourceDrsAssociateSmnTopicRead(ctx, d, meta)
}

func buildCreateDrsAssociateSmnTopicBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"topic_urn":                 d.Get("topic_urn"),
		"is_alarm_to_user":          d.Get("is_alarm_to_user"),
		"increment_delay_threshold": utils.ValueIgnoreEmpty(d.Get("increment_delay_threshold")),
		"rto_delay_threshold":       utils.ValueIgnoreEmpty(d.Get("rto_delay_threshold")),
		"rpo_delay_threshold":       utils.ValueIgnoreEmpty(d.Get("rpo_delay_threshold")),
	}
	return bodyParams
}

func resourceDrsAssociateSmnTopicRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDrsAssociateSmnTopicUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDrsAssociateSmnTopicDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DRS associate SMN topic resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
