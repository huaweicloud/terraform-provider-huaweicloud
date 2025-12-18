package cceautopilot

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE POST /autopilot/v2/charts
// @API CCE GET /autopilot/v2/charts/{chart_id}
// @API CCE PUT /autopilot/v2/charts/{chart_id}
// @API CCE DELETE /autopilot/v2/charts/{chart_id}

func ResourceAutopilotChart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAutopilotChartCreate,
		ReadContext:   resourceAutopilotChartRead,
		UpdateContext: resourceAutopilotChartUpdate,
		DeleteContext: resourceAutopilotChartDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parameters": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"translate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instruction": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"chart_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAutopilotChartCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	file, err := os.Open(d.Get("content").(string))
	if err != nil {
		return diag.Errorf("error opening autopilot chart file: %s", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("content", file.Name())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return diag.FromErr(err)
	}

	if v, ok := d.GetOk("parameters"); ok {
		err = writer.WriteField("parameters", v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	writer.Close()

	var (
		createAutopilotChartHttpUrl = "autopilot/v2/charts"
		createAutopilotChartProduct = "cce"
	)
	createAutopilotChartClient, err := cfg.NewServiceClient(createAutopilotChartProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	createAutopilotChartHttpPath := createAutopilotChartClient.Endpoint + createAutopilotChartHttpUrl

	createAutopilotChartOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": writer.FormDataContentType(),
		},
		RawBody: body,
	}

	createAutopilotChartResp, err := createAutopilotChartClient.Request("POST", createAutopilotChartHttpPath, &createAutopilotChartOpt)
	if err != nil {
		return diag.Errorf("error uploading autopilot chart: %s", err)
	}

	createAutopilotChartRespBody, err := utils.FlattenResponse(createAutopilotChartResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", createAutopilotChartRespBody, nil)
	if id == nil {
		return diag.Errorf("error uploading autopilot chart: ID is not found in API response")
	}

	d.SetId(id.(string))

	return resourceAutopilotChartRead(ctx, d, meta)
}

func resourceAutopilotChartRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		getAutopilotChartHttpUrl = "autopilot/v2/charts/{chart_id}"
		getAutopilotChartProduct = "cce"
	)
	getAutopilotChartClient, err := cfg.NewServiceClient(getAutopilotChartProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	getAutopilotChartHttpPath := getAutopilotChartClient.Endpoint + getAutopilotChartHttpUrl
	getAutopilotChartHttpPath = strings.ReplaceAll(getAutopilotChartHttpPath, "{chart_id}", d.Id())

	getAutopilotChartOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getAutopilotChartResp, err := getAutopilotChartClient.Request("GET", getAutopilotChartHttpPath, &getAutopilotChartOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CCE autopilot chart")
	}

	getAutopilotChartRespBody, err := utils.FlattenResponse(getAutopilotChartResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("name", utils.PathSearch("name", getAutopilotChartRespBody, nil)),
		d.Set("value", utils.PathSearch("values", getAutopilotChartRespBody, nil)),
		d.Set("translate", utils.PathSearch("translate", getAutopilotChartRespBody, nil)),
		d.Set("instruction", utils.PathSearch("instruction", getAutopilotChartRespBody, nil)),
		d.Set("version", utils.PathSearch("version", getAutopilotChartRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getAutopilotChartRespBody, nil)),
		d.Set("source", utils.PathSearch("source", getAutopilotChartRespBody, nil)),
		d.Set("public", utils.PathSearch("public", getAutopilotChartRespBody, nil)),
		d.Set("chart_url", utils.PathSearch("chart_url", getAutopilotChartRespBody, nil)),
		d.Set("created_at", utils.PathSearch("create_at", getAutopilotChartRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("update_at", getAutopilotChartRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAutopilotChartUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	file, err := os.Open(d.Get("content").(string))
	if err != nil {
		return diag.Errorf("error opening autopilot chart file: %s", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("content", file.Name())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return diag.FromErr(err)
	}

	if v, ok := d.GetOk("parameters"); ok {
		err = writer.WriteField("parameters", v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	writer.Close()

	var (
		updateAutopilotChartHttpUrl = "autopilot/v2/charts/{chart_id}"
		updateAutopilotChartProduct = "cce"
	)
	updateAutopilotChartClient, err := cfg.NewServiceClient(updateAutopilotChartProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	updateAutopilotChartHttpPath := updateAutopilotChartClient.Endpoint + updateAutopilotChartHttpUrl
	updateAutopilotChartHttpPath = strings.ReplaceAll(updateAutopilotChartHttpPath, "{chart_id}", d.Id())

	updateAutopilotChartOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": writer.FormDataContentType(),
		},
		RawBody: body,
	}

	_, err = updateAutopilotChartClient.Request("PUT", updateAutopilotChartHttpPath, &updateAutopilotChartOpt)
	if err != nil {
		return diag.Errorf("error updating CCE autopilot chart: %s", err)
	}

	return resourceAutopilotChartRead(ctx, d, meta)
}

func resourceAutopilotChartDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteAutopilotChartHttpUrl = "autopilot/v2/charts/{chart_id}"
		deleteAutopilotChartProduct = "cce"
	)
	deleteAutopilotChartClient, err := cfg.NewServiceClient(deleteAutopilotChartProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	deleteAutopilotChartHttpPath := deleteAutopilotChartClient.Endpoint + deleteAutopilotChartHttpUrl
	deleteAutopilotChartHttpPath = strings.ReplaceAll(deleteAutopilotChartHttpPath, "{chart_id}", d.Id())

	deleteAutopilotChartOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = deleteAutopilotChartClient.Request("DELETE", deleteAutopilotChartHttpPath, &deleteAutopilotChartOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CCE autopilot chart")
	}

	return nil
}
