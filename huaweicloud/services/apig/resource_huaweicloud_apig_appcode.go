package apig

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/applications"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes/{app_code_id}
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes/{app_code_id}
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes
func ResourceAppcode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppcodeCreate,
		ReadContext:   resourceAppcodeRead,
		DeleteContext: resourceAppcodeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAppcodeImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the application and APPCODE are located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the dedicated instance to which the application and APPCODE belong.",
			},
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the application to which the APPCODE belongs.",
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The APPCODE value (content).",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the APPCODE.",
			},
		},
	}
}

func resourceAppcodeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		resp *applications.AppCode

		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)

		instanceId = d.Get("instance_id").(string)
		appId      = d.Get("application_id").(string)
		appCode    = d.Get("value").(string)
	)
	client, err := cfg.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	if appCode != "" {
		opt := applications.AppCodeOpts{
			AppCode: appCode,
		}
		resp, err = applications.CreateAppCode(client, instanceId, appId, opt).Extract()
		if err != nil {
			return diag.Errorf("generating APPCODE failed: %s", err)
		}
	} else {
		resp, err = applications.AutoGenerateAppCode(client, instanceId, appId).Extract()
		if err != nil {
			return diag.Errorf("auto generating APPCODE failed: %s", err)
		}
	}

	d.SetId(resp.Id)

	return resourceAppcodeRead(ctx, d, meta)
}

func resourceAppcodeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)

		instanceId = d.Get("instance_id").(string)
		appId      = d.Get("application_id").(string)
		appCodeId  = d.Id()
	)
	client, err := cfg.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	resp, err := applications.GetAppCode(client, instanceId, appId, appCodeId).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error querying APPCODE (%s) from specified application (%s) under dedicated instance (%s)",
				appCodeId, appId, instanceId))
	}
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("value", resp.Code),
		d.Set("created_at", resp.CreateTime),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving APPCODE resource fields: %s", err)
	}
	return nil
}

func resourceAppcodeDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)

		instanceId = d.Get("instance_id").(string)
		appId      = d.Get("application_id").(string)
		appCodeId  = d.Id()
	)
	client, err := cfg.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	err = applications.RemoveAppCode(client, instanceId, appId, appCodeId).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error deleting APPCODE (%s) from specified application (%s) under dedicated instance (%s)",
				appCodeId, appId, instanceId))
	}
	return nil
}

func resourceAppcodeImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<application_id>/<id>', "+
			"but got '%s'", importedId)
	}
	d.SetId(parts[2])
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("application_id", parts[1]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error saving APPCODE resource fields during import: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
