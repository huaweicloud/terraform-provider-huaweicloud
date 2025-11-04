package dew

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var keyReplicateNonUpdatableParams = []string{
	"key_id",
	"replica_region",
	"key_alias",
	"replica_project_id",
	"key_description",
	"enterprise_project_id",
	"tags",
}

// @API DEW POST /v2/{project_id}/kms/keys/{key_id}/replicate
func ResourceKmsKeyReplicate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKmsKeyReplicateCreate,
		ReadContext:   resourceKmsKeyReplicateRead,
		UpdateContext: resourceKmsKeyReplicateUpdate,
		DeleteContext: resourceKmsKeyReplicateDelete,

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(keyReplicateNonUpdatableParams),
			config.MergeDefaultTags(),
		),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"key_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"replica_region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_alias": {
				Type:     schema.TypeString,
				Required: true,
			},
			"replica_project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildKmsKeyReplicateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"key_id":                d.Get("key_id"),
		"replica_region":        d.Get("replica_region"),
		"key_alias":             d.Get("key_alias"),
		"replica_project_id":    d.Get("replica_project_id"),
		"key_description":       utils.ValueIgnoreEmpty(d.Get("key_description")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(d.Get("enterprise_project_id")),
		"tags":                  utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
	}
}

func resourceKmsKeyReplicateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/kms/keys/{key_id}/replicate"
	)

	client, err := cfg.NewServiceClient("kms", region)
	if err != nil {
		return diag.Errorf("error creating DEW KMS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{key_id}", d.Get("key_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildKmsKeyReplicateBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error replicating DEW KMS key: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	replicaKeyId := utils.PathSearch("key_id", respBody, "").(string)
	if replicaKeyId == "" {
		return diag.Errorf("error replicating DEW KMS key: ID is not found in API response")
	}

	d.SetId(replicaKeyId)
	return resourceKmsKeyReplicateRead(ctx, d, meta)
}

func resourceKmsKeyReplicateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceKmsKeyReplicateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceKmsKeyReplicateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to replicate a DEW KMS key.
Deleting this resource will not change the current DEW KMS key, but will only remove the resource information from the
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
