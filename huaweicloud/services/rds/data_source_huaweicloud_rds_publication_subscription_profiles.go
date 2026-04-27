package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS GET /v3/{project_id}/instances/{instance_id}/replication/profiles
func DataSourceRdsPublicationSubscriptionProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsPublicationSubscriptionProfilesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"agent_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"profiles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     publicationSubscriptionProfilesSchema(),
			},
		},
	}
}

func publicationSubscriptionProfilesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"profile_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"profile_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"agent_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_def_profile": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceRdsPublicationSubscriptionProfilesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/replication/profiles"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath += buildListPublicationSubscriptionProfilesQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving RDS publication and subscription profiles: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("profiles", flattenListPublicationSubscriptionProfilesBody(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListPublicationSubscriptionProfilesQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=100"
	if v, ok := d.GetOk("agent_type"); ok {
		queryParams += fmt.Sprintf("&agent_type=%v", v)
	}
	return queryParams
}

func flattenListPublicationSubscriptionProfilesBody(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	profiles := utils.PathSearch("profiles", resp, make([]interface{}, 0)).([]interface{})
	if len(profiles) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(profiles))
	for i, profile := range profiles {
		result[i] = map[string]interface{}{
			"profile_id":     utils.PathSearch("profile_id", profile, nil),
			"profile_name":   utils.PathSearch("profile_name", profile, nil),
			"agent_type":     utils.PathSearch("agent_type", profile, nil),
			"description":    utils.PathSearch("description", profile, nil),
			"is_def_profile": utils.PathSearch("is_def_profile", profile, nil),
		}
	}

	return result
}
