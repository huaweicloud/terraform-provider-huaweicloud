package identitycenter

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IdentityCenter GET /v1/identity-stores/{identity_store_id}/group-memberships
// @API IdentityCenter GET /v1/identity-stores/{identity_store_id}/group-memberships-for-member
func DataSourceIdentityCenterGroupMemberships() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityCenterGroupMembershipsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"identity_store_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_memberships": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"member_id": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"identity_store_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"membership_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildListGroupMembershipsParams(marker string, d *schema.ResourceData) string {
	res := "?limit=100"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	if v, ok := d.GetOk("group_id"); ok {
		res = fmt.Sprintf("%s&group_id=%v", res, v)
	}

	return res
}

func listGroupMemberships(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/identity-stores/{identity_store_id}/group-memberships"
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{identity_store_id}", d.Get("identity_store_id").(string))

	queryParams := buildListGroupMembershipsParams("", d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestResp, err := client.Request("GET", listPath+queryParams, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving Identity Center group memberships: %s", err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		memberships := utils.PathSearch("group_memberships", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, memberships...)

		marker := utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}

		queryParams = buildListGroupMembershipsParams(marker, d)
	}
	return result, nil
}

func buildListGroupMembershipsForMemberParams(marker string, d *schema.ResourceData) string {
	res := "?limit=100"

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	if v, ok := d.GetOk("user_id"); ok {
		res = fmt.Sprintf("%s&user_id=%v", res, v)
	}

	return res
}

func listGroupMembershipsForMember(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/identity-stores/{identity_store_id}/group-memberships-for-member"
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{identity_store_id}", d.Get("identity_store_id").(string))

	queryParams := buildListGroupMembershipsForMemberParams("", d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestResp, err := client.Request("GET", listPath+queryParams, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving Identity Center group memberships for member: %s", err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		memberships := utils.PathSearch("group_memberships", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, memberships...)

		marker := utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}

		queryParams = buildListGroupMembershipsForMemberParams(marker, d)
	}
	return result, nil
}

func dataSourceIdentityCenterGroupMembershipsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("identitystore", region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	_, groupIDOk := d.GetOk("group_id")
	_, userIDOk := d.GetOk("user_id")

	if !groupIDOk && !userIDOk {
		return diag.Errorf("Exactly one of group_id or user_id must be set")
	}
	if groupIDOk && userIDOk {
		return diag.Errorf("Only one of group_id or user_id can be set")
	}

	var memberships []interface{}

	if groupIDOk {
		memberships, err = listGroupMemberships(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if userIDOk {
		memberships, err = listGroupMembershipsForMember(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("group_memberships", flattenGroupMemberships(memberships)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGroupMemberships(memberships []interface{}) []interface{} {
	if len(memberships) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(memberships))
	for _, membership := range memberships {
		result = append(result, map[string]interface{}{
			"membership_id":     utils.PathSearch("membership_id", membership, nil),
			"group_id":          utils.PathSearch("group_id", membership, nil),
			"identity_store_id": utils.PathSearch("identity_store_id", membership, nil),
			"member_id":         flattenMemberId(utils.PathSearch("member_id", membership, nil)),
		})
	}
	return result
}

func flattenMemberId(memberId interface{}) []map[string]interface{} {
	if memberId == nil || len(memberId.(map[string]interface{})) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"user_id": utils.PathSearch("user_id", memberId, nil),
		},
	}
}
