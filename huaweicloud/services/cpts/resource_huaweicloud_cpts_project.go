package cpts

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cpts/v1/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,
		ReadContext:   resourceProjectRead,
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

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 42),
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 50),
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

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcCptsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CPTS v1 client: %s", err)
	}

	createOpts := &model.CreateProjectRequest{
		Body: &model.CreateProjectRequestBody{
			Name:        d.Get("name").(string),
			Description: utils.String(d.Get("description").(string)),
		},
	}
	response, err := client.CreateProject(createOpts)
	if err != nil {
		return diag.Errorf("error creating CPTS project: %s", err)
	}

	if response.ProjectId == nil {
		return diag.Errorf("error creating CPTS project: id not found in api response")
	}

	d.SetId(strconv.Itoa(int(*response.ProjectId)))
	return resourceProjectRead(ctx, d, meta)
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcCptsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CPTS v1 client: %s", err)
	}

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.Errorf("the project ID must be integer: %s", err)
	}

	response, err := client.ShowProject(&model.ShowProjectRequest{
		TestSuiteId: int32(id),
	})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "the project is not found")
	}

	layout := "2006-01-02T15:04:05-07:00"
	createTime, err := time.Parse(layout, *response.Project.CreateTime)
	if err != nil {
		return diag.Errorf("error parsing the time: %s", err)
	}
	updateTime, err := time.Parse(layout, *response.Project.UpdateTime)
	if err != nil {
		return diag.Errorf("error parsing the time: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", response.Project.Name),
		d.Set("description", response.Project.Description),
		d.Set("created_at", utils.FormatTimeStampUTC(createTime.Unix())),
		d.Set("updated_at", utils.FormatTimeStampUTC(updateTime.Unix())),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcCptsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CPTS v1 client: %s", err)
	}

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.Errorf("the project ID must be integer: %s", err)
	}

	_, err = client.UpdateProject(&model.UpdateProjectRequest{
		TestSuiteId: int32(id),
		Body: &model.UpdateProjectRequestBody{
			Id:          int32(id),
			Name:        d.Get("name").(string),
			Description: utils.String(d.Get("description").(string)),
		},
	})

	if err != nil {
		return diag.Errorf("error updating the project %q: %s", id, err)
	}

	return resourceProjectRead(ctx, d, meta)
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcCptsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CPTS v1 client: %s", err)
	}

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.Errorf("the project ID must be integer: %s", err)
	}

	deleteOpts := &model.DeleteProjectRequest{
		TestSuiteId: int32(id),
	}

	_, err = client.DeleteProject(deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting CPTS project %q: %s", id, err)
	}

	return nil
}
