// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RDS
// ---------------------------------------------------------------

package rds

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

// @API RDS GET /v3/{project_id}/configurations
func DataSourceParametergroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceParametergroupsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datastore_version_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datastore_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_defined": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"configurations": {
				Type:     schema.TypeList,
				Elem:     ConfigurationsSchema(),
				Computed: true,
			},
		},
	}
}

func ConfigurationsSchema() *schema.Resource {
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
			"datastore_version_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datastore_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_defined": {
				Type:     schema.TypeBool,
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
	return &sc
}

func DataSourceParametergroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listConfigurationsHttpUrl = "v3/{project_id}/configurations"
		listConfigurationsProduct = "rds"
	)
	listConfigurationsClient, err := cfg.NewServiceClient(listConfigurationsProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS Client: %s", err)
	}

	listConfigurationsPath := listConfigurationsClient.Endpoint + listConfigurationsHttpUrl
	listConfigurationsPath = strings.ReplaceAll(listConfigurationsPath, "{project_id}", listConfigurationsClient.ProjectID)
	listConfigurationsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	listConfigurationsResp, err := listConfigurationsClient.Request("GET", listConfigurationsPath, &listConfigurationsOpt)
	if err != nil {
		return diag.Errorf("error retrieving RDS configurations")
	}
	listConfigurationsRespBody, err := utils.FlattenResponse(listConfigurationsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("configurations", flattenConfigurationsBody(listConfigurationsRespBody, d)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenConfigurationsBody(resp interface{}, d *schema.ResourceData) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("configurations", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	configurationsName, nameOk := d.GetOk("name")
	configurationsDatastoreVersionName, datastoreVersionNameOk := d.GetOk("datastore_version_name")
	configurationsDatastoreName, datastoreNameOk := d.GetOk("datastore_name")
	configurationsUserDefined, _ := d.GetOk("user_defined")
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		name := utils.PathSearch("name", v, "")
		datastoreVersionName := utils.PathSearch("datastore_version_name", v, "")
		datastoreName := utils.PathSearch("datastore_name", v, "")
		userDefined := utils.PathSearch("user_defined", v, false).(bool)
		if nameOk && configurationsName.(string) != name.(string) {
			continue
		}
		if datastoreVersionNameOk && configurationsDatastoreVersionName.(string) != datastoreVersionName.(string) {
			continue
		}
		if datastoreNameOk && configurationsDatastoreName.(string) != datastoreName.(string) {
			continue
		}
		if configurationsUserDefined.(bool) != userDefined {
			continue
		}
		rst = append(rst, map[string]interface{}{
			"id":                     utils.PathSearch("id", v, nil),
			"name":                   name,
			"description":            utils.PathSearch("description", v, nil),
			"datastore_version_name": datastoreVersionName,
			"datastore_name":         datastoreName,
			"user_defined":           userDefined,
			"created_at":             utils.PathSearch("created", v, nil),
			"updated_at":             utils.PathSearch("updated", v, nil),
		})
	}
	return rst
}
