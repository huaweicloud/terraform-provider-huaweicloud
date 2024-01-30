// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DDM
// ---------------------------------------------------------------

package ddm

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDM GET /v2/{project_id}/engines
func DataSourceDdmEngines() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDdmEnginesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the engine version.`,
			},
			"engines": {
				Type:        schema.TypeList,
				Elem:        EnginesEngineSchema(),
				Computed:    true,
				Description: `Indicates the list of DDM engine.`,
			},
		},
	}
}

func EnginesEngineSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the engine.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the engine version.`,
			},
		},
	}
	return &sc
}

func resourceDdmEnginesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDdmEngines: Query the List of DDM engines
	var (
		getDdmEnginesHttpUrl = "v2/{project_id}/engines"
		getDdmEnginesProduct = "ddm"
	)
	getDdmEnginesClient, err := cfg.NewServiceClient(getDdmEnginesProduct, region)
	if err != nil {
		return diag.Errorf("error creating DdmEngines Client: %s", err)
	}

	getDdmEnginesPath := getDdmEnginesClient.Endpoint + getDdmEnginesHttpUrl
	getDdmEnginesPath = strings.ReplaceAll(getDdmEnginesPath, "{project_id}", getDdmEnginesClient.ProjectID)

	getDdmEnginesResp, err := pagination.ListAllItems(
		getDdmEnginesClient,
		"offset",
		getDdmEnginesPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DdmEngines")
	}

	getDdmEnginesRespJson, err := json.Marshal(getDdmEnginesResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getDdmEnginesRespBody interface{}
	err = json.Unmarshal(getDdmEnginesRespJson, &getDdmEnginesRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	version := d.Get("version").(string)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("engines", flattenGetEnginesResponseBodyEngine(getDdmEnginesRespBody, version)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetEnginesResponseBodyEngine(resp interface{}, version string) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("engineGroups", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		engineVersion := utils.PathSearch("version", v, nil)
		if version == "" || version == engineVersion {
			rst = append(rst, map[string]interface{}{
				"id":      utils.PathSearch("id", v, nil),
				"version": engineVersion,
			})
		}
	}
	return rst
}
