package workspace

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace GET /v2/{project_id}/assist-auth-config/apply-objects
func DataSourceAssistAuthConfigurationAppliedObjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssistAuthConfigurationAppliedObjectsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the applied objects of assist auth configuration are located.`,
			},

			// Optional parameters.
			"object_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the binding object.`,
			},
			"object_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the object.`,
			},

			// Attributes.
			"objects": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of assist auth configuration applied objects that matched filter parameters.`,
				Elem:        assistAuthConfigurationAppliedObjectSchema(),
			},
		},
	}
}

func assistAuthConfigurationAppliedObjectSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the user or user group.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the binding object.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the user or user group.`,
			},
		},
	}
}

func buildAssistAuthConfigurationAppliedObjectsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("object_type"); ok {
		res = fmt.Sprintf("%s&object_type=%v", res, v)
	}
	if v, ok := d.GetOk("object_name"); ok {
		res = fmt.Sprintf("%s&object_name=%v", res, v)
	}

	return res
}

func listAssistAuthConfigurationAppliedObjects(client *golangsdk.ServiceClient, d ...*schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/assist-auth-config/apply-objects?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	if len(d) > 0 {
		listPath += buildAssistAuthConfigurationAppliedObjectsQueryParams(d[0])
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		objects := utils.PathSearch("objects", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, objects...)
		if len(objects) < limit {
			break
		}
		offset += len(objects)
	}

	return result, nil
}

func dataSourceAssistAuthConfigurationAppliedObjectsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	resp, err := listAssistAuthConfigurationAppliedObjects(client, d)
	if err != nil {
		return diag.Errorf("error querying Workspace assist auth configuration applied objects: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("objects", flattenAssistAuthConfigurationAppliedObjects(resp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
