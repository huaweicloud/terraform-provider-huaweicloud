// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product IdentityCenter
// ---------------------------------------------------------------

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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IdentityCenter GET /v1/identity-stores/{identity_store_id}/users
func DataSourceIdentityCenterUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceIdentityCenterUsersRead,

		Description: "schema: Internal",
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"identity_store_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the identity store that associated with IAM Identity Center.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the user.`,
			},
			"family_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the family name of the user.`,
			},
			"given_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the given name of the user.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the display name of the user.`,
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the email of the user.`,
			},
			"users": {
				Type:        schema.TypeList,
				Elem:        identityCenterUsersUserSchema(),
				Computed:    true,
				Description: `Indicates the list of IdentityCenter user.`,
			},
		},
	}
}

func identityCenterUsersUserSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the user.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the user.`,
			},
			"family_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the family name of the user.`,
			},
			"given_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the given name of the user.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the display name of the user.`,
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the email of the user.`,
			},
			"phone_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The phone number of the user.`,
			},
			"user_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the user.`,
			},
			"title": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The title of the user.`,
			},
			"addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The addresses information of the user.`,
				Elem:        identityCenterUsersAddressesSchema(),
			},
			"enterprise": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The enterprise information of the user.`,
				Elem:        identityCenterUsersEnterpriseSchema(),
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the user.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the user.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time of the user.`,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The updater of the user.`,
			},
			"email_verified": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the email is verified.`,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the user is enabled.`,
			},
		},
	}
	return &sc
}

func identityCenterUsersAddressesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"country": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The country of the user.`,
			},
			"formatted": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `A string containing a formatted version of the address to be displayed.`,
			},
			"locality": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The locality of the user.`,
			},
			"postal_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The postal code of the user.`,
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The region of the user.`,
			},
			"street_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The street address of the user.`,
			},
		},
	}
	return &sc
}

func identityCenterUsersEnterpriseSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"cost_center": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The cost center of the enterprise.`,
			},
			"department": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The department of the enterprise.`,
			},
			"division": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The division of the enterprise.`,
			},
			"employee_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The employee number of the enterprise.`,
			},
			"organization": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The organization of the enterprise.`,
			},
			"manager": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The manager of the enterprise.`,
			},
		},
	}
	return &sc
}

func resourceIdentityCenterUsersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getIdentityCenterUsers: Query Identity Center users
	var (
		getIdentityCenterUsersHttpUrl = "v1/identity-stores/{identity_store_id}/users"
		getIdentityCenterUsersProduct = "identitystore"
	)
	getIdentityCenterUsersClient, err := cfg.NewServiceClient(getIdentityCenterUsersProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	getIdentityCenterUsersBasePath := getIdentityCenterUsersClient.Endpoint + getIdentityCenterUsersHttpUrl
	getIdentityCenterUsersBasePath = strings.ReplaceAll(getIdentityCenterUsersBasePath, "{identity_store_id}",
		d.Get("identity_store_id").(string))

	getIdentityCenterUsersOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var users []interface{}
	var marker string
	var getIdentityCenterUsersPath string
	for {
		getIdentityCenterUsersPath = getIdentityCenterUsersBasePath + buildGetIdentityCenterUsersQueryParams(d, marker)
		getIdentityCenterUsersResp, err := getIdentityCenterUsersClient.Request("GET",
			getIdentityCenterUsersPath, &getIdentityCenterUsersOpt)

		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving Identity Center users")
		}

		getIdentityCenterUsersRespBody, err := utils.FlattenResponse(getIdentityCenterUsersResp)
		if err != nil {
			return diag.FromErr(err)
		}
		users = append(users, flattenGetIdentityCenterUsersResponseBodyUser(d, getIdentityCenterUsersRespBody)...)
		marker = utils.PathSearch("page_info.next_marker", getIdentityCenterUsersRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("users", users),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetIdentityCenterUsersResponseBodyUser(d *schema.ResourceData, resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("users", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rawFamilyName := utils.PathSearch("name.family_name", v, "").(string)
		if familyName, ok := d.GetOk("family_name"); ok && rawFamilyName != familyName {
			continue
		}
		rawGivenName := utils.PathSearch("name.given_name", v, "").(string)
		if givenName, ok := d.GetOk("given_name"); ok && rawGivenName != givenName {
			continue
		}
		rawDisplayName := utils.PathSearch("display_name", v, "").(string)
		if displayName, ok := d.GetOk("display_name"); ok && rawDisplayName != displayName {
			continue
		}
		rawEmail := utils.PathSearch("emails|[0].value", v, "").(string)
		if email, ok := d.GetOk("email"); ok && rawEmail != email {
			continue
		}
		rst = append(rst, map[string]interface{}{
			"id":             utils.PathSearch("user_id", v, nil),
			"user_name":      utils.PathSearch("user_name", v, nil),
			"family_name":    rawFamilyName,
			"given_name":     rawGivenName,
			"display_name":   rawDisplayName,
			"email":          rawEmail,
			"phone_number":   utils.PathSearch("phone_numbers|[0].value", v, nil),
			"title":          utils.PathSearch("title", v, nil),
			"user_type":      utils.PathSearch("user_type", v, nil),
			"created_by":     utils.PathSearch("created_by", v, nil),
			"updated_by":     utils.PathSearch("updated_by", v, nil),
			"email_verified": utils.PathSearch("email_verified", v, false),
			"enabled":        utils.PathSearch("enabled", v, false),
			"created_at":     utils.FormatTimeStampRFC3339(int64(utils.PathSearch("created_at", v, float64(0)).(float64))/1000, false),
			"updated_at":     utils.FormatTimeStampRFC3339(int64(utils.PathSearch("updated_at", v, float64(0)).(float64))/1000, false),
			"addresses":      flattenIdentityCenterUserAddresses(utils.PathSearch("addresses|[0]", v, nil)),
			"enterprise":     flattenIdentityCenterUserEnterprise(utils.PathSearch("enterprise", v, nil)),
		})
	}
	return rst
}

func buildGetIdentityCenterUsersQueryParams(d *schema.ResourceData, marker string) string {
	res := "?limit=100"
	if v, ok := d.GetOk("user_name"); ok {
		res = fmt.Sprintf("%s&user_name=%v", res, v)
	}

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}
