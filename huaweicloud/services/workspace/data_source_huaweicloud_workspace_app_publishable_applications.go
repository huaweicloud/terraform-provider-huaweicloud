package workspace

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceAppPublishableApplications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppPublishableApplicationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},

			// Required parameters.
			"app_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the application group.`,
			},

			// Attributes.
			"group_images": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of image IDs under the server group.`,
			},
			"applications": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        schemaApplications(),
				Description: `The list of the publishable applications.`,
			},

			// Deprecated attributes.
			"apps": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     schemaApplications(),
				Description: utils.SchemaDesc(
					`The list of the publishable applications.`,
					utils.SchemaDescInput{
						Deprecated: true,
					},
				),
			},
		},
	}
}

func schemaApplications() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the the application.`,
			},
			"execute_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The execution path where the application is located.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version of the the application.`,
			},
			"publisher": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The publisher of the the application.`,
			},
			"work_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The work path of the the application.`,
			},
			"command_param": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The command line arguments used to start the application.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the the application.`,
			},
			"publishable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the application can be published.`,
			},
			"icon_index": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The index of the application icon.`,
			},
			"icon_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The path where the application icon is located.`,
			},
			"source_image_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of image IDs to which the application belongs.`,
			},
		},
	}
}

type AppPublishableApplicationsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newAppPublishableApplicationsDSWrapper(d *schema.ResourceData, meta interface{}) *AppPublishableApplicationsDSWrapper {
	return &AppPublishableApplicationsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceAppPublishableApplicationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newAppPublishableApplicationsDSWrapper(d, meta)
	shoPubAppRst, err := wrapper.ShowPublishableApplications()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.showPublishableApplicationsToSchema(shoPubAppRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API Workspace GET /v1/{project_id}/app-groups/{app_group_id}/publishable-app
func (w *AppPublishableApplicationsDSWrapper) ShowPublishableApplications() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "appstream")
	if err != nil {
		return nil, err
	}

	uri := "/v1/{project_id}/app-groups/{app_group_id}/publishable-app"
	uri = strings.ReplaceAll(uri, "{app_group_id}", w.Get("app_group_id").(string))
	resp, err := httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
		// When the status code is 404, the error message is empty in some cases. so it is uniformly processed as an empty list.
		// The scenarios with status code 404 are as follows:
		// a. The application group id does not exist.
		// b. The application group is not associated with a server group.
		// c. There is no server under the server group associated with the application group, and the API response error message is empty.
	if _, ok := err.(golangsdk.ErrDefault404); ok {
		return &gjson.Result{}, nil
	}
	return resp, err
}

func flattenPublishableApplications(body *gjson.Result) any {
	return schemas.SliceToList(body.Get("items"),
		func(apps gjson.Result) any {
			return map[string]any{
				"name":             apps.Get("name").Value(),
				"execute_path":     apps.Get("execute_path").Value(),
				"version":          apps.Get("version").Value(),
				"publisher":        apps.Get("publisher").Value(),
				"work_path":        apps.Get("work_path").Value(),
				"command_param":    apps.Get("command_param").Value(),
				"description":      apps.Get("description").Value(),
				"publishable":      apps.Get("publishable").Value(),
				"icon_index":       apps.Get("icon_index").Value(),
				"icon_path":        apps.Get("icon_path").Value(),
				"source_image_ids": schemas.SliceToStrList(apps.Get("source_image_ids")),
			}
		},
	)
}

func (w *AppPublishableApplicationsDSWrapper) showPublishableApplicationsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		// Attributes.
		d.Set("group_images", body.Get("group_images").Value()),
		d.Set("applications", flattenPublishableApplications(body)),
		// Deprecated attributes.
		d.Set("apps", flattenPublishableApplications(body)),
	)
	return mErr.ErrorOrNil()
}
