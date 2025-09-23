// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RDS
// ---------------------------------------------------------------

package rds

import (
	"context"
	"fmt"
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

var pgHbaNonUpdatableParams = []string{"instance_id"}

// @API RDS POST /v3/{project_id}/instances/{instance_id}/hba-info
// @API RDS GET /v3/{project_id}/instances/{instance_id}/hba-info
func ResourcePgHba() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePgHbaCreateOrUpdate,
		UpdateContext: resourcePgHbaCreateOrUpdate,
		ReadContext:   resourcePgHbaRead,
		DeleteContext: resourcePgHbaDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(pgHbaNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the RDS PostgreSQL instance.`,
			},
			"host_based_authentications": {
				Type:        schema.TypeList,
				Elem:        pgHbaHostBasedAuthenticationSchema(),
				Required:    true,
				Description: `Specifies the list of host based authentications.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func pgHbaHostBasedAuthenticationSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the connection type.`,
			},
			"database": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the database name.`,
			},
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the Name of a user.`,
			},
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the client IP address.`,
			},
			"method": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the authentication mode.`,
			},
			"mask": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the subnet mask.`,
			},
		},
	}
	return &sc
}

func resourcePgHbaCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := setPgHba(d, meta, buildPgHbaBodyParams(d.Get("host_based_authentications")))
	if err != nil {
		return diag.Errorf("error setting RDS PostgreSQL hba: %s", err)
	}

	d.SetId(d.Get("instance_id").(string))

	return resourcePgHbaRead(ctx, d, meta)
}

func buildPgHbaBodyParams(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"priority": i,
				"type":     raw["type"],
				"database": raw["database"],
				"user":     raw["user"],
				"address":  raw["address"],
				"method":   raw["method"],
				"mask":     utils.ValueIgnoreEmpty(raw["mask"]),
			}
		}
		return rst
	}
	return nil
}

func resourcePgHbaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getPgHba: query RDS PostgreSQL hba
	var (
		getPgHbaHttpUrl = "v3/{project_id}/instances/{instance_id}/hba-info"
		getPgHbaProduct = "rds"
	)
	getPgHbaClient, err := cfg.NewServiceClient(getPgHbaProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	getPgHbaPath := getPgHbaClient.Endpoint + getPgHbaHttpUrl
	getPgHbaPath = strings.ReplaceAll(getPgHbaPath, "{project_id}", getPgHbaClient.ProjectID)
	getPgHbaPath = strings.ReplaceAll(getPgHbaPath, "{instance_id}", d.Id())

	getPgHbaOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getPgHbaResp, err := getPgHbaClient.Request("GET", getPgHbaPath, &getPgHbaOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code||errCode", "DBS.280343"),
			"error retrieving RDS PostgreSQL hba")
	}

	getPgHbaRespBody, err := utils.FlattenResponse(getPgHbaResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", d.Id()),
		d.Set("host_based_authentications", flattenPgHbaRequestBodyHostBasedAuthentication(getPgHbaRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPgHbaRequestBodyHostBasedAuthentication(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curArray := resp.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"type":     utils.PathSearch("type", v, nil),
			"database": utils.PathSearch("database", v, nil),
			"user":     utils.PathSearch("user", v, nil),
			"address":  utils.PathSearch("address", v, nil),
			"method":   utils.PathSearch("method", v, nil),
			"mask":     utils.PathSearch("mask", v, nil),
		})
	}
	return rst
}

func resourcePgHbaDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := setPgHba(d, meta, make([]interface{}, 0))
	if err != nil {
		return diag.Errorf("error deleting RDS PostgreSQL hba: %s", err)
	}

	return nil
}

func setPgHba(d *schema.ResourceData, meta interface{}, requestOpt interface{}) error {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// setPgHba: create RDS PostgreSQL hba.
	var (
		pgHbaHttpUrl = "v3/{project_id}/instances/{instance_id}/hba-info"
		pgHbaProduct = "rds"
	)
	pgHbaClient, err := cfg.NewServiceClient(pgHbaProduct, region)
	if err != nil {
		return fmt.Errorf("error creating RDS client: %s", err)
	}

	pgHbaPath := pgHbaClient.Endpoint + pgHbaHttpUrl
	pgHbaPath = strings.ReplaceAll(pgHbaPath, "{project_id}", pgHbaClient.ProjectID)
	pgHbaPath = strings.ReplaceAll(pgHbaPath, "{instance_id}", d.Get("instance_id").(string))

	pgHbaOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	pgHbaOpt.JSONBody = requestOpt
	pgHbaResp, err := pgHbaClient.Request("POST", pgHbaPath, &pgHbaOpt)
	if err != nil {
		return err
	}

	pgHbaRespBody, err := utils.FlattenResponse(pgHbaResp)
	if err != nil {
		return err
	}
	message := utils.PathSearch("message", pgHbaRespBody, "")
	if message != "" {
		return fmt.Errorf("%s", message)
	}
	return nil
}
