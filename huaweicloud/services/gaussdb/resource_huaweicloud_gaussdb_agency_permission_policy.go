package gaussdb

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var agencyPermissionPolicyNonUpdatableParams = []string{
	"agency_name",
}

// @API GaussDB PUT /v3/{project_id}/agency/{agency_name}/policy
// @API GaussDB GET /v3/{project_id}/agency/{agency_name}
func ResourceAgencyPermissionPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAgencyPermissionPolicyCreate,
		ReadContext:   resourceAgencyPermissionPolicyRead,
		UpdateContext: resourceAgencyPermissionPolicyUpdate,
		DeleteContext: resourceAgencyPermissionPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(agencyPermissionPolicyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"agency_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bind_role_names": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"unbind_role_names": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"is_existed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAgencyPermissionPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"bind_role_names":   utils.ExpandToStringListBySet(d.Get("bind_role_names").(*schema.Set)),
		"unbind_role_names": utils.ExpandToStringListBySet(d.Get("unbind_role_names").(*schema.Set)),
	}

	return bodyParams
}

func updateAgencyPermissionPolicyParameters(client *golangsdk.ServiceClient, d *schema.ResourceData, agencyName string) error {
	httpUrl := "v3/{project_id}/agency/{agency_name}/policy"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{agency_name}", agencyName)
	updateopt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildAgencyPermissionPolicyBodyParams(d)),
	}

	_, err := client.Request("PUT", updatePath, &updateopt)
	if err != nil {
		return err
	}

	return nil
}

func resourceAgencyPermissionPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		agencyName = d.Get("agency_name").(string)
	)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	err = updateAgencyPermissionPolicyParameters(client, d, agencyName)
	if err != nil {
		return diag.Errorf("error updating GaussDB agency permission policy: %s", err)
	}

	d.SetId(agencyName)

	return resourceAgencyPermissionPolicyRead(ctx, d, meta)
}

func resourceAgencyPermissionPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	agencyInfo, err := GetAgencyPermissionPolicyInfo(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "errCode", "DBS.280001"),
			"error retrieving GaussDB agency permission policy")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("agency_name", utils.PathSearch("name", agencyInfo, nil)),
		d.Set("is_existed", utils.PathSearch("is_existed", agencyInfo, nil)),
		d.Set("roles", flattenAgencyPermissionPolicyInfo(
			utils.PathSearch("roles", agencyInfo, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAgencyPermissionPolicyInfo(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"name":        utils.PathSearch("name", v, nil),
			"description": utils.PathSearch("description", v, nil),
		})
	}

	return rst
}

func GetAgencyPermissionPolicyInfo(client *golangsdk.ServiceClient, agencyName string) (interface{}, error) {
	httpUrl := "v3/{project_id}/agency/{agency_name}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{agency_name}", agencyName)
	getOpts := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceAgencyPermissionPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	if d.HasChangesExcept("enable_force_new") {
		err = updateAgencyPermissionPolicyParameters(client, d, d.Id())
		if err != nil {
			return diag.Errorf("error updating GaussDB agency permission policy: %s", err)
		}
	}

	return resourceAgencyPermissionPolicyRead(ctx, d, meta)
}

func resourceAgencyPermissionPolicyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GaussDB agency permission policy resource is not supported. The resource is only removed from the state," +
		"the instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
