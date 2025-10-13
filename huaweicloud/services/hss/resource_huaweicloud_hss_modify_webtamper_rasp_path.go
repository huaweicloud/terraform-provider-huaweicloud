package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var modifyWebtamperRaspPathNonUpdatableParams = []string{"host_id", "rasp_path", "enterprise_project_id", "host_name"}

// @API HSS PUT /v5/{project_id}/wtp/{host_id}/rasp-path
func ResourceModifyWebtamperRaspPath() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceModifyWebtamperRaspPathCreate,
		ReadContext:   resourceModifyWebtamperRaspPathRead,
		UpdateContext: resourceModifyWebtamperRaspPathUpdate,
		DeleteContext: resourceModifyWebtamperRaspPathDelete,

		CustomizeDiff: config.FlexibleForceNew(modifyWebtamperRaspPathNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rasp_path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// This parameter does not tack effect
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
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

func buildModifyWebtamperRaspPathBodyParams(d *schema.ResourceData, epsId string) string {
	queryParams := ""

	if v, ok := d.GetOk("host_name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func resourceModifyWebtamperRaspPathCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/wtp/{host_id}/rasp-path"
		epsId   = cfg.GetEnterpriseProjectID(d)
		hostId  = d.Get("host_id").(string)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{host_id}", hostId)
	requestPath += buildModifyWebtamperRaspPathBodyParams(d, epsId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"rasp_path": d.Get("rasp_path"),
		},
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating the dynamic web tamper protection Tomcat bin directory: %s", err)
	}

	d.SetId(hostId)

	return resourceModifyWebtamperRaspPathRead(ctx, d, meta)
}

func resourceModifyWebtamperRaspPathRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceModifyWebtamperRaspPathUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceModifyWebtamperRaspPathDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to updating the dynamic web tamper protection Tomcat bin directory.
	  Deleting this resource will not clear the corresponding close records, but will only remove the resource information from
	  the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
