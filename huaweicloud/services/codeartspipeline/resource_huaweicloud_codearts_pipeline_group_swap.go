package codeartspipeline

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

var pipelineGroupSwapNonUpdatableParams = []string{
	"project_id", "group_id1", "group_id2",
}

// @API CodeArtsPipeline POST /v5/{project_id}/api/pipeline-group/swap
func ResourceCodeArtsPipelineGroupSwap() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCodeArtsPipelineGroupSwapCreate,
		ReadContext:   resourceCodeArtsPipelineGroupSwapRead,
		UpdateContext: resourceCodeArtsPipelineGroupSwapUpdate,
		DeleteContext: resourceCodeArtsPipelineGroupSwapDelete,

		CustomizeDiff: config.FlexibleForceNew(pipelineGroupSwapNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the project ID for CodeArts service.`,
			},
			"group_id1": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the pipeline group ID1.`,
			},
			"group_id2": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the pipeline group ID2.`,
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

func resourceCodeArtsPipelineGroupSwapCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts pipeline client: %s", err)
	}

	createHttpUrl := "v5/{project_id}/api/pipeline-group/swap?groupId1={groupId1}&groupId2={groupId2}"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", d.Get("project_id").(string))
	createPath = strings.ReplaceAll(createPath, "{groupId1}", d.Get("group_id1").(string))
	createPath = strings.ReplaceAll(createPath, "{groupId2}", d.Get("group_id2").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error swaping pipeline group: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(createRespBody, ""); err != nil {
		return diag.Errorf("error swaping pipeline group: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	return nil
}

func resourceCodeArtsPipelineGroupSwapRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCodeArtsPipelineGroupSwapUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCodeArtsPipelineGroupSwapDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting pipeline group swap resource is not supported. The resource is only removed from the state," +
		" the pipeline group remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
