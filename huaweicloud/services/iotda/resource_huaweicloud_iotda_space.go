package iotda

import (
	"context"
	"log"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func ResourceSpace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSpaceCreate,
		ReadContext:   resourceSpaceRead,
		DeleteContext: resourceSpaceDelete,
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`[A-Za-z-_0-9?'#().,&%@!]{1,64}`),
					"The name contains a maximum of 64 characters. Only letters, digits, "+
						"hyphens (-), underscore (_) and the following specail characters are allowed: ?'#().,&%@!"),
			},

			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceSpaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	createOpts := model.AddApplicationRequest{
		Body: &model.AddApplication{
			AppName: d.Get("name").(string),
		},
	}
	log.Printf("[DEBUG] Create IoTDA space params: %#v", createOpts)

	resp, err := client.AddApplication(&createOpts)
	if err != nil {
		return diag.Errorf("error creating IoTDA space: %s", err)
	}

	if resp.AppId == nil {
		return diag.Errorf("error creating IoTDA space: id is not found in API response")
	}

	d.SetId(*resp.AppId)
	return resourceSpaceRead(ctx, d, meta)
}

func resourceSpaceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	detail, err := client.ShowApplication(&model.ShowApplicationRequest{AppId: d.Id()})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA space")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", detail.AppName),
		d.Set("is_default", detail.DefaultApp),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSpaceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	deleteOpts := &model.DeleteApplicationRequest{
		AppId: d.Id(),
	}
	_, err = client.DeleteApplication(deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting IoTDA space: %s", err)
	}

	return nil
}
