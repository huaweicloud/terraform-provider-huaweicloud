package apig

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

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-acl
func DataSourceApplicationAcl() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationAclRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the application and ACL rules are located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the dedicated instance to which the application belongs.",
			},
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the application to which the ACL rules belong.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ACL type.",
			},
			"values": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The ACL values.",
			},
		},
	}
}

func dataSourceApplicationAclRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-acl"
		instanceId = d.Get("instance_id").(string)
		appId      = d.Get("application_id").(string)
		mErr       *multierror.Error
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{app_id}", appId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return diag.Errorf("error retrieving the application ACL rules from the application (%s): %s", appId, err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(mErr,
		d.Set("region", region),
		d.Set("type", utils.PathSearch("app_acl_type", respBody, "").(string)),
		d.Set("values", utils.PathSearch("app_acl_values", respBody, make([]interface{}, 0)).([]interface{})),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving the fields of the application ACL rules: %s", err)
	}
	return nil
}
