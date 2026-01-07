package dns

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

var resolverAccessLogNonUpdatableParams = []string{"lts_group_id", "lts_topic_id"}

// @API DNS POST /v2/resolver/queryloggingconfig
// @API DNS GET /v2/resolver/queryloggingconfig/{id}
// @API DNS POST /v2/resolver/queryloggingconfig/{id}/associatevpc
// @API DNS POST /v2/resolver/queryloggingconfig/{id}/disassociatevpc
// @API DNS DELETE /v2/resolver/queryloggingconfig/{id}
func ResourceResolverAccessLog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResolverAccessLogCreate,
		ReadContext:   resourceResolverAccessLogRead,
		UpdateContext: resourceResolverAccessLogUpdate,
		DeleteContext: resourceResolverAccessLogDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(resolverAccessLogNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the resolver access log is located.`,
			},
			"lts_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the log group.`,
			},
			"lts_topic_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the log stream.`,
			},
			"vpc_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of VPC IDs associated with the resolver access log.`,
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

func buildResolverAccessLogVpc(vpcId string) map[string]interface{} {
	return map[string]interface{}{
		"vpc_id": vpcId,
	}
}

func buildResolverAccessLogBodyParams(d *schema.ResourceData, vpcId string) map[string]interface{} {
	return map[string]interface{}{
		"lts_group_id": d.Get("lts_group_id"),
		"lts_topic_id": d.Get("lts_topic_id"),
		"vpc":          buildResolverAccessLogVpc(vpcId),
	}
}

func resourceResolverAccessLogCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v2/resolver/queryloggingconfig"
		vpcIds  = d.Get("vpc_ids").(*schema.Set)
	)
	client, err := cfg.NewServiceClient("dns_region", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildResolverAccessLogBodyParams(d, vpcIds.List()[0].(string))),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating resolver access log: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	accessLogId := utils.PathSearch("resolver_query_log_config.id", respBody, "").(string)
	if accessLogId == "" {
		return diag.Errorf("unable to find resolver access log ID from API response")
	}

	d.SetId(accessLogId)

	// Because the first item is the VPC associated when creating, it has been associated successfully.
	// So start from the second item, associate the VPC.
	if vpcIds.Len() > 1 {
		if err = associateResolverAccessLogVpcIds(client, accessLogId, vpcIds.List()[1:]); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceResolverAccessLogRead(ctx, d, meta)
}

func GetResolverAccessLog(client *golangsdk.ServiceClient, accessLogId string) (interface{}, error) {
	httpUrl := "v2/resolver/queryloggingconfig/{id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{id}", accessLogId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("resolver_query_log_config", respBody, nil), nil
}

func resourceResolverAccessLogRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dns_region", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	accessLogId := d.Id()
	respBody, err := GetResolverAccessLog(client, accessLogId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving resolver access log (%s): %s", accessLogId, err))
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("lts_group_id", utils.PathSearch("lts_group_id", respBody, nil)),
		d.Set("lts_topic_id", utils.PathSearch("lts_topic_id", respBody, nil)),
		d.Set("vpc_ids", utils.PathSearch("vpc_ids", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateResolverAccessLogVpcBodyParams(vpcId string) golangsdk.RequestOpts {
	return golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"vpc": buildResolverAccessLogVpc(vpcId),
		},
	}
}

func associateResolverAccessLogVpcIds(client *golangsdk.ServiceClient, accessLogId string, vpcIds []interface{}) error {
	httpUrl := "v2/resolver/queryloggingconfig/{id}/associatevpc"
	associatePath := client.Endpoint + httpUrl
	associatePath = strings.ReplaceAll(associatePath, "{id}", accessLogId)
	for _, vpcId := range vpcIds {
		associateOpt := buildUpdateResolverAccessLogVpcBodyParams(vpcId.(string))
		if _, err := client.Request("POST", associatePath, &associateOpt); err != nil {
			return fmt.Errorf("error associating VPC (%v) to resolver access log (%s): %s", vpcId, accessLogId, err)
		}
	}

	return nil
}

func disassociateResolverAccessLogVpcIds(client *golangsdk.ServiceClient, accessLogId string, vpcIds []interface{}) error {
	httpUrl := "v2/resolver/queryloggingconfig/{id}/disassociatevpc"
	associatePath := client.Endpoint + httpUrl
	associatePath = strings.ReplaceAll(associatePath, "{id}", accessLogId)
	for _, vpcId := range vpcIds {
		associateOpt := buildUpdateResolverAccessLogVpcBodyParams(vpcId.(string))
		_, err := client.Request("POST", associatePath, &associateOpt)
		if err != nil {
			return fmt.Errorf("error disassociating VPC (%v) from resolver access log (%s): %s", vpcId, accessLogId, err)
		}
	}

	return nil
}

func resourceResolverAccessLogUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dns_region", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	accessLogId := d.Id()
	oldVal, newVal := d.GetChange("vpc_ids")
	rmRaw := oldVal.(*schema.Set).Difference(newVal.(*schema.Set))
	addRaw := newVal.(*schema.Set).Difference(oldVal.(*schema.Set))
	if rmRaw.Len() > 0 {
		if err = disassociateResolverAccessLogVpcIds(client, accessLogId, rmRaw.List()); err != nil {
			return diag.FromErr(err)
		}
	}

	if addRaw.Len() > 0 {
		if err = associateResolverAccessLogVpcIds(client, accessLogId, addRaw.List()); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceResolverAccessLogRead(ctx, d, meta)
}

func resourceResolverAccessLogDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		accessLogId = d.Id()
	)
	client, err := cfg.NewServiceClient("dns_region", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	respBody, err := GetResolverAccessLog(client, accessLogId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error getting vpc ids associated with resolver access log (%s): %s",
			accessLogId, err))
	}

	// If the resolver access log is associated with more than one VPC, the remaining VPCs must be unbound before the resource can be deleted.
	if vpcIds := utils.PathSearch("vpc_ids", respBody, make([]interface{}, 0)).([]interface{}); len(vpcIds) > 1 {
		if err = disassociateResolverAccessLogVpcIds(client, accessLogId, vpcIds[1:]); err != nil {
			return diag.FromErr(err)
		}
	}

	err = deleteResolverAccessLog(client, accessLogId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting DNS resolver access log (%s): %s", accessLogId, err))
	}
	return nil
}

func deleteResolverAccessLog(client *golangsdk.ServiceClient, accessLogId string) error {
	httpUrl := "v2/resolver/queryloggingconfig/{id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{id}", accessLogId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	return err
}
