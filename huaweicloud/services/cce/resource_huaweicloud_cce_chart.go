package cce

import (
	"context"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/def"
	cce "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cce/v3/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API CCE POST /v2/charts
// @API CCE GET /v2/charts/{chart_id}
// @API CCE PUT /v2/charts/{chart_id}
// @API CCE DELETE /v2/charts/{chart_id}

func ResourceChart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChartCreate,
		ReadContext:   resourceChartRead,
		UpdateContext: resourceChartUpdate,
		DeleteContext: resourceChartDelete,

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

func resourceChartCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.HcCceV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE v3 client: %s", err)
	}

	file, err := os.Open(d.Get("content").(string))
	if err != nil {
		return diag.Errorf("error opening chart file: %s", err)
	}
	defer file.Close()

	createOpts := cce.UploadChartRequestBody{
		Parameters: &def.MultiPart{
			Content: d.Get("parameters").(string),
		},
		Content: &def.FilePart{
			Content: file,
		},
	}

	req := cce.UploadChartRequest{
		Body: &createOpts,
	}
	resp, err := client.UploadChart(&req)
	if err != nil {
		return diag.Errorf("error uploading CCE chart: %s", err)
	}

	if resp == nil || resp.Id == nil {
		return diag.Errorf("unable to find resource ID in the response: %v", resp)
	}

	d.SetId(*resp.Id)

	return resourceChartRead(ctx, d, meta)
}

func resourceChartRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.HcCceV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE v3 client: %s", err)
	}

	req := cce.ShowChartRequest{
		ChartId: d.Id(),
	}

	resp, err := client.ShowChart(&req)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CCE chart")
	}

	if resp == nil {
		return diag.Errorf("unable to find the response: %v", resp)
	}

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("name", resp.Name),
		d.Set("value", resp.Values),
		d.Set("translate", resp.Translate),
		d.Set("instruction", resp.Instruction),
		d.Set("version", resp.Version),
		d.Set("description", resp.Description),
		d.Set("source", resp.Source),
		d.Set("public", resp.Public),
		d.Set("chart_url", resp.ChartUrl),
		d.Set("created_at", resp.CreateAt),
		d.Set("updated_at", resp.UpdateAt),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceChartUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.HcCceV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE v3 client: %s", err)
	}

	file, err := os.Open(d.Get("content").(string))
	if err != nil {
		return diag.Errorf("error opening chart file: %s", err)
	}
	defer file.Close()

	updateOpts := cce.UpdateChartRequestBody{
		Parameters: &def.MultiPart{
			Content: d.Get("parameters").(string),
		},
		Content: &def.FilePart{
			Content: file,
		},
	}

	req := cce.UpdateChartRequest{
		ChartId: d.Id(),
		Body:    &updateOpts,
	}
	_, err = client.UpdateChart(&req)
	if err != nil {
		return diag.Errorf("error updating CCE chart: %s", err)
	}

	return resourceChartRead(ctx, d, meta)
}

func resourceChartDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.HcCceV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE v3 client: %s", err)
	}

	req := cce.DeleteChartRequest{
		ChartId: d.Id(),
	}

	_, err = client.DeleteChart(&req)
	if err != nil {
		return diag.Errorf("error deleting CCE chart: %s", err)
	}
	return nil
}
