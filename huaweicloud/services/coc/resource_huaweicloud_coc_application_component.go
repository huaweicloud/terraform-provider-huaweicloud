package coc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var applicationComponentNonUpdatableParams = []string{"application_id"}

// @API COC POST /v1/components
// @API COC PUT /v1/components/{id}
// @API COC DELETE /v1/components/{id}
// @API COC GET /v1/components
func ResourceApplicationComponent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationComponentCreate,
		ReadContext:   resourceApplicationComponentRead,
		UpdateContext: resourceApplicationComponentUpdate,
		DeleteContext: resourceApplicationComponentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(applicationComponentNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ep_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceApplicationComponentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	createHttpUrl := "v1/components"
	createPath := client.Endpoint + createHttpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateApplicationComponentBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating COC application component: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening response: %s", err)
	}

	id := utils.PathSearch("data.id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the COC application component ID from the API response")
	}

	d.SetId(id)

	return resourceApplicationComponentRead(ctx, d, meta)
}

func buildCreateApplicationComponentBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"application_id": d.Get("application_id"),
		"name":           d.Get("name"),
	}

	return bodyParams
}

func resourceApplicationComponentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	applicationID := ""
	if v, ok := d.GetOk("application_id"); ok {
		applicationID = v.(string)
	}
	applicationComponent, err := GetApplicationComponent(client, applicationID, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving application component")
	}

	mErr = multierror.Append(mErr,
		d.Set("name", utils.PathSearch("name", applicationComponent, nil)),
		d.Set("code", utils.PathSearch("code", applicationComponent, nil)),
		d.Set("application_id", utils.PathSearch("application_id", applicationComponent, nil)),
		d.Set("ep_id", utils.PathSearch("ep_id", applicationComponent, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetApplicationComponent(client *golangsdk.ServiceClient, applicationID string, componentID string) (interface{}, error) {
	getHttpUrl := "v1/components?limit=100"
	basePath := client.Endpoint + getHttpUrl
	if applicationID != "" {
		basePath = fmt.Sprintf("%s&application_id=%v", basePath, applicationID)
	}
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	var marker string
	for {
		getPath := basePath + buildGetApplicationComponentParams(marker)
		getResp, err := client.Request("GET", getPath, &getOpt)

		if err != nil {
			return nil, err
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, err
		}
		component, hasNext, nextMarker := flattenCocGetApplicationComponent(getRespBody, componentID)
		if hasNext {
			marker = nextMarker
			continue
		}
		if component != nil {
			return component, nil
		}
		return nil, golangsdk.ErrDefault404{}
	}
}

func buildGetApplicationComponentParams(marker string) string {
	res := ""
	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	return res
}

func flattenCocGetApplicationComponent(resp interface{}, componentID string) (map[string]interface{}, bool, string) {
	if resp == nil {
		return nil, false, ""
	}
	componentsJson := utils.PathSearch("data", resp, make([]interface{}, 0))
	componentsArray := componentsJson.([]interface{})
	if len(componentsArray) == 0 {
		return nil, false, ""
	}

	marker := ""
	var res map[string]interface{}
	for _, component := range componentsArray {
		res = map[string]interface{}{
			"id":             utils.PathSearch("id", component, nil),
			"name":           utils.PathSearch("name", component, nil),
			"code":           utils.PathSearch("code", component, nil),
			"application_id": utils.PathSearch("application_id", component, nil),
			"ep_id":          utils.PathSearch("ep_id", component, nil),
		}
		marker = utils.PathSearch("id", component, "").(string)
		if marker == componentID {
			return res, false, marker
		}
	}
	return nil, true, marker
}

func resourceApplicationComponentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	updateHttpUrl := "v1/components/{id}"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{id}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateApplicationComponentBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating application component: %s", err)
	}

	return resourceApplicationComponentRead(ctx, d, meta)
}

func buildUpdateApplicationComponentBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name": d.Get("name"),
	}

	return bodyParams
}

func resourceApplicationComponentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	deleteHttpUrl := "v1/components/{id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code",
			"common.00000400"), "error deleting COC application component")
	}

	return nil
}
