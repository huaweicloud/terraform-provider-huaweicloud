package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var eventDeleteIsolatedFileNonUpdatableParams = []string{"data_list", "data_list.*.host_id", "data_list.*.file_hash",
	"data_list.*.file_path", "data_list.*.file_attr", "enterprise_project_id"}

// This resource is used to delete HSS isolated files.
// Due to the lack of test scenarios, this code is not tested and is not documented externally.
// Documentation is only stored in docs/incubating.

// @API HSS DELETE /v5/{project_id}/event/isolated-file
func ResourceEventDeleteIsolatedFile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEventDeleteIsolatedFileCreate,
		ReadContext:   resourceEventDeleteIsolatedFileRead,
		UpdateContext: resourceEventDeleteIsolatedFileUpdate,
		DeleteContext: resourceEventDeleteIsolatedFileDelete,

		CustomizeDiff: config.FlexibleForceNew(eventDeleteIsolatedFileNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"file_hash": {
							Type:     schema.TypeString,
							Required: true,
						},
						"file_path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"file_attr": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"enterprise_project_id": {
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

func buildCreateEventDeleteIsolatedFileBodyParams(d *schema.ResourceData) map[string]interface{} {
	dataList := d.Get("data_list").([]interface{})
	dataListRequestBody := make([]interface{}, 0, len(dataList))
	for _, v := range dataList {
		dataListRequestBody = append(dataListRequestBody, map[string]interface{}{
			"host_id":   utils.PathSearch("host_id", v, nil),
			"file_hash": utils.PathSearch("file_hash", v, nil),
			"file_path": utils.PathSearch("file_path", v, nil),
			"file_attr": utils.PathSearch("file_attr", v, nil),
		})
	}

	return map[string]interface{}{
		"data_list": dataListRequestBody,
	}
}

func resourceEventDeleteIsolatedFileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/event/isolated-file"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	if epsId != "" {
		requestPath += fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateEventDeleteIsolatedFileBodyParams(d),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting HSS isolated file: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	return resourceEventDeleteIsolatedFileRead(ctx, d, meta)
}

func resourceEventDeleteIsolatedFileRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceEventDeleteIsolatedFileUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceEventDeleteIsolatedFileDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to delete HSS isolated file. Deleting this resource
    will not clear the corresponding isolated file deletion record, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
