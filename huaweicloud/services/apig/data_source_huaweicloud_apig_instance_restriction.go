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

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/restriction
func DataSourceInstanceRestriction() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceRestrictionRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the dedicated instance is located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to be queried.`,
			},

			// Attributes.
			"restrict_cidrs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of restricted IP CIDR blocks.`,
			},
			"resource_subnet_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The CIDR block of the resource subnet.`,
			},
		},
	}
}

func getInstanceRestriction(client *golangsdk.ServiceClient, instanceId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/restriction"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func dataSourceInstanceRestrictionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
	)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	resp, err := getInstanceRestriction(client, instanceId)
	if err != nil {
		return diag.Errorf("error querying the restricted information of the dedicated instance (%s): %s", instanceId, err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("restrict_cidrs", utils.PathSearch("restrict_cidrs", resp, make([]interface{}, 0))),
		d.Set("resource_subnet_cidr", utils.PathSearch("resource_subnet_cidr", resp, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
