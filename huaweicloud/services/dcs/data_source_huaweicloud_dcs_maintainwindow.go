package dcs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DCS GET /v2/instances/maintain-windows
func DataSourceDcsMaintainWindow() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsMaintainWindowRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"seq": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"begin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"end": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceDcsMaintainWindowRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v2/instances/maintain-windows"
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DCS maintain windows: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	filerMaintainWindow := filterMaintainWindow(d, getRespBody)
	if len(filerMaintainWindow) < 1 {
		return diag.Errorf("your query returned no results. Please change your search criteria and try again.")
	}

	mErr = multierror.Append(nil,
		d.Set("region", region),
		d.Set("seq", utils.PathSearch("[0].seq", filerMaintainWindow, nil)),
		d.Set("begin", utils.PathSearch("[0].begin", filerMaintainWindow, nil)),
		d.Set("end", utils.PathSearch("[0].end", filerMaintainWindow, nil)),
		d.Set("default", utils.PathSearch("[0].default", filerMaintainWindow, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterMaintainWindow(d *schema.ResourceData, getRespBody interface{}) []interface{} {
	maintainWindows := utils.PathSearch("maintain_windows", getRespBody, make([]interface{}, 0)).([]interface{})
	if len(maintainWindows) < 1 {
		return nil
	}

	result := make([]interface{}, 0)

	rawSeq, rawSeqOK := d.GetOk("seq")
	rawBegin, rawBeginOK := d.GetOk("begin")
	rawEnd, rawEndOk := d.GetOk("end")
	rawDefault, rawDefaultOK := d.GetOk("default")

	for _, backupRecord := range maintainWindows {
		seq := utils.PathSearch("seq", backupRecord, float64(0)).(float64)
		begin := utils.PathSearch("begin", backupRecord, nil)
		end := utils.PathSearch("end", backupRecord, nil)
		isDefault := utils.PathSearch("default", backupRecord, nil)
		if rawSeqOK && rawSeq.(int) != int(seq) {
			continue
		}
		if rawBeginOK && rawBegin != begin {
			continue
		}
		if rawEndOk && rawEnd != end {
			continue
		}
		if rawDefaultOK && rawDefault != isDefault {
			continue
		}
		result = append(result, backupRecord)
	}

	return result
}
