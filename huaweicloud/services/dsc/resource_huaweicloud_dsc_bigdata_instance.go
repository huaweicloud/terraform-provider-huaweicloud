package dsc

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var bigdataInstanceNonUpdatableParams = []string{
	"asset_name",
	"create_time",
	"ds_address",
	"ds_name",
	"ds_password",
	"ds_port",
	"ds_type",
	"ds_user",
	"ds_version",
	"ins_id",
	"ins_name",
	"ins_type",
	"lts_group_id",
	"lts_group_name",
	"lts_stream_id",
	"lts_stream_name",
	"queue_name",
	"region",
	"scan_metadata",
	"security_group_id",
	"subnet_id",
	"vpc_id",
}

// @API DSC POST /v1/{project_id}/asset-center/bigdata/instances
func ResourceBigdataInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigdataInstanceCreate,
		ReadContext:   resourceBigdataInstanceRead,
		UpdateContext: resourceBigdataInstanceUpdate,
		DeleteContext: resourceBigdataInstanceDelete,

		CustomizeDiff: config.FlexibleForceNew(bigdataInstanceNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"asset_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ds_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ds_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ds_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"ds_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ds_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ds_user": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ds_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ins_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ins_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ins_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lts_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lts_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lts_stream_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lts_stream_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"queue_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scan_metadata": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildBigdataInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"asset_name":        utils.ValueIgnoreEmpty(d.Get("asset_name")),
		"create_time":       utils.ValueIgnoreEmpty(d.Get("create_time")),
		"ds_address":        utils.ValueIgnoreEmpty(d.Get("ds_address")),
		"ds_name":           utils.ValueIgnoreEmpty(d.Get("ds_name")),
		"ds_password":       utils.ValueIgnoreEmpty(d.Get("ds_password")),
		"ds_port":           utils.ValueIgnoreEmpty(d.Get("ds_port")),
		"ds_type":           utils.ValueIgnoreEmpty(d.Get("ds_type")),
		"ds_user":           utils.ValueIgnoreEmpty(d.Get("ds_user")),
		"ds_version":        utils.ValueIgnoreEmpty(d.Get("ds_version")),
		"ins_id":            utils.ValueIgnoreEmpty(d.Get("ins_id")),
		"ins_name":          utils.ValueIgnoreEmpty(d.Get("ins_name")),
		"ins_type":          utils.ValueIgnoreEmpty(d.Get("ins_type")),
		"lts_group_id":      utils.ValueIgnoreEmpty(d.Get("lts_group_id")),
		"lts_group_name":    utils.ValueIgnoreEmpty(d.Get("lts_group_name")),
		"lts_stream_id":     utils.ValueIgnoreEmpty(d.Get("lts_stream_id")),
		"lts_stream_name":   utils.ValueIgnoreEmpty(d.Get("lts_stream_name")),
		"queue_name":        utils.ValueIgnoreEmpty(d.Get("queue_name")),
		"region":            utils.ValueIgnoreEmpty(d.Get("region_name")),
		"scan_metadata":     utils.ValueIgnoreEmpty(d.Get("scan_metadata")),
		"security_group_id": utils.ValueIgnoreEmpty(d.Get("security_group_id")),
		"subnet_id":         utils.ValueIgnoreEmpty(d.Get("subnet_id")),
		"vpc_id":            utils.ValueIgnoreEmpty(d.Get("vpc_id")),
	}

	return utils.RemoveNil(bodyParams)
}

func resourceBigdataInstanceCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/asset-center/bigdata/instances"
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildBigdataInstanceBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error adding DSC bigdata instance: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId.String())

	mErr := multierror.Append(nil,
		d.Set("msg", utils.PathSearch("msg", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceBigdataInstanceRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceBigdataInstanceUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceBigdataInstanceDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to add a bigdata instance to the DSC asset center.
Deleting this resource will not remove the added bigdata instance or undo the add action, but will only remove
the resource information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
