// Generated by PMS #3
package dds

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceDdsDatabaseVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDdsDatabaseVersionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"datastore_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the database name.`,
			},
			"versions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Indicates the database version.`,
			},
		},
	}
}

type DatabaseVersionsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newDatabaseVersionsDSWrapper(d *schema.ResourceData, meta interface{}) *DatabaseVersionsDSWrapper {
	return &DatabaseVersionsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDdsDatabaseVersionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newDatabaseVersionsDSWrapper(d, meta)
	rst, err := wrapper.ListDatastoreVersions()
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := uuid.GenerateUUID()
	d.SetId(id)

	err = wrapper.listDatastoreVersionsToSchema(rst)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// @API DDS GET /v3/{project_id}/datastores/{datastore_name}/versions
func (w *DatabaseVersionsDSWrapper) ListDatastoreVersions() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "dds")
	if err != nil {
		return nil, err
	}

	d := w.ResourceData
	uri := "/v3/{project_id}/datastores/{datastore_name}/versions"
	uri = strings.ReplaceAll(uri, "{datastore_name}", d.Get("datastore_name").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *DatabaseVersionsDSWrapper) listDatastoreVersionsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("versions", body.Get("versions").Value()),
	)
	return mErr.ErrorOrNil()
}
