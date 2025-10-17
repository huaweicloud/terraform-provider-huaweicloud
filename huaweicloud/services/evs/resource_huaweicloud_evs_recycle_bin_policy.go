package evs

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

var recycleBinPolicyNonUpdatableParams = []string{
	"switch",
	"threshold_time",
	"keep_time",
}

// @API EVS PUT /v3/{project_id}/recycle-bin-volumes/policy
func ResourceRecycleBinPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRecycleBinPolicyCreate,
		ReadContext:   resourceRecycleBinPolicyRead,
		UpdateContext: resourceRecycleBinPolicyUpdate,
		DeleteContext: resourceRecycleBinPolicyDelete,

		CustomizeDiff: config.FlexibleForceNew(recycleBinPolicyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"switch": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"threshold_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"keep_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildRecycleBinPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"switch":         d.Get("switch"),
		"threshold_time": utils.ValueIgnoreEmpty(d.Get("threshold_time")),
		"keep_time":      utils.ValueIgnoreEmpty(d.Get("keep_time")),
	}

	return bodyParams
}

func resourceRecycleBinPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/recycle-bin-volumes/policy"
		product = "evs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	createRecycleBinPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildRecycleBinPolicyBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &createRecycleBinPolicyOpt)
	if err != nil {
		return diag.Errorf("error creating EVS recycle bin policy: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	return resourceRecycleBinPolicyRead(ctx, d, meta)
}

func resourceRecycleBinPolicyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceRecycleBinPolicyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceRecycleBinPolicyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to update EVS recycle bin policy. Deleting this
    resource will not clear the corresponding request record, but will only remove the resource information from the tf
    state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
