package workspace

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace GET /v1/{project_id}/app-center/apps
func DataSourceApplications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the applications are located.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The application name to be queried and supports fuzzy matching.`,
			},
			"applications": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        applicationSchema(),
				Description: `The list of applications that match the filter parameters.`,
			},
		},
	}
}

func applicationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the application.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the application.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the application.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version of the application.`,
			},
			"authorization_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The authorization type of the application.`,
			},
			"application_file_store": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        applicationFileStoreSchema(),
				Description: `The file store configuration of the application.`,
			},
			"application_icon_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The icon URL of the application.`,
			},
			"install_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The installation type of the application.`,
			},
			"install_command": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The installation command of the application.`,
			},
			"uninstall_command": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uninstallation command of the application.`,
			},
			"support_os": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The supported operating system of the application.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the application.`,
			},
			"application_source": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The source of the application.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the application, in UTC format.`,
			},
			"catalog_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The catalog ID of the application.`,
			},
			"catalog": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The catalog name of the application.`,
			},
			"install_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The installation information of the application.`,
			},
		},
	}
}

func applicationFileStoreSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"store_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The store type of the application file.`,
			},
			"bucket_store": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        applicationBucketStoreSchema(),
				Description: `The OBS bucket store configuration.`,
			},
			"file_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The external file link.`,
			},
		},
	}
}

func applicationBucketStoreSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the OBS bucket.`,
			},
			"bucket_file_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The file path in the OBS bucket.`,
			},
		},
	}
}

func flattenApplications(applications []interface{}) []interface{} {
	if len(applications) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(applications))
	for _, item := range applications {
		result = append(result, map[string]interface{}{
			"id":                     utils.PathSearch("id", item, nil),
			"name":                   utils.PathSearch("name", item, nil),
			"description":            utils.PathSearch("description", item, nil),
			"version":                utils.PathSearch("version", item, nil),
			"authorization_type":     utils.PathSearch("authorization_type", item, nil),
			"application_file_store": flattenApplicationFileStore(utils.PathSearch("app_file_store", item, nil)),
			"application_icon_url":   utils.PathSearch("app_icon_url", item, nil),
			"install_type":           utils.PathSearch("install_type", item, nil),
			"install_command":        utils.PathSearch("install_command", item, nil),
			"uninstall_command":      utils.PathSearch("uninstall_command", item, nil),
			"support_os":             utils.PathSearch("support_os", item, nil),
			"status":                 utils.PathSearch("status", item, nil),
			"application_source":     utils.PathSearch("application_source", item, nil),
			"create_time":            utils.PathSearch("create_time", item, nil),
			"catalog_id":             utils.PathSearch("catalog_id", item, nil),
			"catalog":                utils.PathSearch("catalog", item, nil),
			"install_info":           utils.PathSearch("install_info", item, nil),
		})
	}
	return result
}

func dataSourceApplicationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	applications, err := listApplications(client, d)
	if err != nil {
		return diag.Errorf("error querying Workspace applications: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("applications", flattenApplications(applications)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
