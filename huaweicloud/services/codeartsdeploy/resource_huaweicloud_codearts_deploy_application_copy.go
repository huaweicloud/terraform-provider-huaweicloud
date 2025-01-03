package codeartsdeploy

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CodeArtsDeploy POST /v1/applications/{app_id}/duplicate
func ResourceDeployApplicationCopy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCodeArtsDeployApplicationCopyCreate,
		ReadContext:   resourceCodeArtsDeployApplicationCopyRead,
		DeleteContext: resourceCodeArtsDeployApplicationCopyDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the application ID.`,
			},
		},
	}
}

func resourceCodeArtsDeployApplicationCopyCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_deploy", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	httpUrl := "v1/applications/{app_id}/duplicate"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{app_id}", d.Get("app_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error copying CodeArts deploy application: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	copyId := utils.PathSearch("result.id", createRespBody, "").(string)
	if copyId == "" {
		return diag.Errorf("unable to find the CodeArts deploy new application ID from the API response")
	}

	d.SetId(copyId)

	return nil
}

func resourceCodeArtsDeployApplicationCopyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCodeArtsDeployApplicationCopyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting application copy resource is not supported. The resource is only removed from the" +
		"state, the application remain in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
