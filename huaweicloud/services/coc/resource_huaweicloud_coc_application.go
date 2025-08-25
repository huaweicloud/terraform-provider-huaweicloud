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

var applicationNonUpdatableParams = []string{"parent_id"}

// @API COC POST /v1/applications
// @API COC PUT /v1/applications/{id}
// @API COC DELETE /v1/applications/{id}
// @API COC GET /v1/applications
func ResourceApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationCreate,
		ReadContext:   resourceApplicationRead,
		UpdateContext: resourceApplicationUpdate,
		DeleteContext: resourceApplicationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(applicationNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_collection": {
				Type:     schema.TypeBool,
				Optional: true,
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
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceApplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	createHttpUrl := "v1/applications"
	createPath := client.Endpoint + createHttpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateApplicationBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating COC application: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening response: %s", err)
	}

	id := utils.PathSearch("data.id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the COC application ID from the API response")
	}

	d.SetId(id)

	if _, ok := d.GetOk("is_collection"); ok {
		err = updateApplication(client, d)
		if err != nil {
			return diag.Errorf("error updating COC application `is_collection`: %s", err)
		}
	}

	return resourceApplicationRead(ctx, d, meta)
}

func buildCreateApplicationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"parent_id":   utils.ValueIgnoreEmpty(d.Get("parent_id")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}

	return bodyParams
}

func resourceApplicationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	application, err := GetApplication(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving application")
	}

	mErr = multierror.Append(mErr,
		d.Set("name", utils.PathSearch("name", application, nil)),
		d.Set("code", utils.PathSearch("code", application, nil)),
		d.Set("parent_id", utils.PathSearch("parent_id", application, nil)),
		d.Set("description", utils.PathSearch("description", application, nil)),
		d.Set("path", utils.PathSearch("path", application, nil)),
		d.Set("is_collection", utils.PathSearch("is_collection", application, nil)),
		d.Set("create_time", utils.PathSearch("create_time", application, nil)),
		d.Set("update_time", utils.PathSearch("update_time", application, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetApplication(client *golangsdk.ServiceClient, applicationID string) (interface{}, error) {
	getHttpUrl := "v1/applications?id_list={application_id}&limit=1"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{application_id}", applicationID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening application: %s", err)
	}

	application := utils.PathSearch("data[0]", getRespBody, nil)
	if application == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return application, nil
}

func resourceApplicationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	changeList := []string{
		"name", "description", "is_collection",
	}
	if d.HasChanges(changeList...) {
		err = updateApplication(client, d)
		if err != nil {
			return diag.Errorf("error updating application: %s", err)
		}
	}

	return resourceApplicationRead(ctx, d, meta)
}

func updateApplication(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updateHttpUrl := "v1/applications/{id}"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{id}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateApplicationBodyParams(d)),
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func buildUpdateApplicationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":          d.Get("name"),
		"description":   d.Get("description"),
		"is_collection": utils.ValueIgnoreEmpty(d.Get("is_collection")),
	}

	return bodyParams
}

func resourceApplicationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	deleteHttpUrl := "v1/applications/{id}"
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
			"common.00000400"), "error deleting COC application")
	}

	return nil
}
