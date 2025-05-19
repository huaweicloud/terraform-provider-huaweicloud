package apig

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

var apiCheckNonUpdatableParams = []string{
	"instance_id",
	"type",
	"api_id",
	"name",
	"group_id",
	"req_method",
	"req_uri",
	"match_mode",
}

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/apis/check
func ResourceApiCheck() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApiCheckCreate,
		ReadContext:   resourceApiCheckRead,
		UpdateContext: resourceApiCheckUpdate,
		DeleteContext: resourceApiCheckDelete,

		CustomizeDiff: config.FlexibleForceNew(apiCheckNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the dedicated instance to which the API belongs.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the API to be checked.",
			},
			"api_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the API to be excluded from the check.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the API.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the group to which the API belongs.",
			},
			"req_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The request method of the API.",
			},
			"req_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The request path of the API.",
			},
			"match_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The matching mode of the API.",
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

func buildApiCheckBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"type":       d.Get("type"),
		"api_id":     utils.ValueIgnoreEmpty(d.Get("api_id")),
		"name":       utils.ValueIgnoreEmpty(d.Get("name")),
		"group_id":   utils.ValueIgnoreEmpty(d.Get("group_id")),
		"req_method": utils.ValueIgnoreEmpty(d.Get("req_method")),
		"req_uri":    utils.ValueIgnoreEmpty(d.Get("req_uri")),
		"match_mode": utils.ValueIgnoreEmpty(d.Get("match_mode")),
	}
}

func resourceApiCheckCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/apis/check"
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildApiCheckBodyParams(d)),
		OkCodes:          []int{204},
	}
	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error checking API under APIG instance (%s): %s", instanceId, err)
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomId)

	return nil
}

func resourceApiCheckRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceApiCheckUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceApiCheckDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time resource for checking the API definition. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
