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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDM GET /v3/{project_id}/configurations
func DataSourceDdmConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDdmConfigurationsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"configurations": {
				Type:     schema.TypeList,
				Elem:     parameterConfigurationsSchema(),
				Computed: true,
			},
		},
	}
}

func parameterConfigurationsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datastore_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_defined": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceDdmConfigurationsRead(_ context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDdmConfigurations: Query the List of DDM configuration
	var (
		getDdmConfigurationsHttpUrl = "v3/{project_id}/configurations"
		getDdmConfigurationsProduct = "ddm"
	)
	getDdmConfigurationsClient, err := cfg.NewServiceClient(getDdmConfigurationsProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDM Client: %s", err)
	}

	getDdmConfigurationsPath := getDdmConfigurationsClient.Endpoint + getDdmConfigurationsHttpUrl
	getDdmConfigurationsPath = strings.ReplaceAll(getDdmConfigurationsPath, "{project_id}", getDdmConfigurationsClient.ProjectID)

	getDdmConfigurationsResp, err := pagination.ListAllItems(
		getDdmConfigurationsClient,
		"offset",
		getDdmConfigurationsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		diag.Errorf("error retrieving DDM configurations: %s", err)
	}

	getDdmConfigurationsRespJson, err := json.Marshal(getDdmConfigurationsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var getDdmConfigurationsRespBody any
	err = json.Unmarshal(getDdmConfigurationsRespJson, &getDdmConfigurationsRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("configurations", flattenGetConfigurationsResponseBody(getDdmConfigurationsRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetConfigurationsResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("configurations", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))

	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":             utils.PathSearch("id", v, nil),
			"name":           utils.PathSearch("name", v, nil),
			"description":    utils.PathSearch("description", v, nil),
			"datastore_name": utils.PathSearch("datastore_name", v, nil),
			"created":        utils.PathSearch("created", v, nil),
			"updated":        utils.PathSearch("updated", v, nil),
			"user_defined":   utils.PathSearch("user_defined", v, nil),
		})
	}
	return rst
}
