package obs

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/obs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API OBS PUT ?mirrorBackToSource
// @API OBS DELETE ?mirrorBackToSource
// @API OBS GET ?mirrorBackToSource
func ResourceObsBucketMirrorBackToSource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceObsBucketMirrorBackToSourceCreate,
		ReadContext:   resourceObsBucketMirrorBackToSourceRead,
		UpdateContext: resourceObsBucketMirrorBackToSourceUpdate,
		DeleteContext: resourceObsBucketMirrorBackToSourceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceObsBucketMirrorBackToSourceImport,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{"bucket"}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					}),
			},
		},
	}
}

func buildMirrorBackToSourceRules(rulesStr string) string {
	rule := utils.StringToJson(rulesStr)
	if rule == nil {
		return rulesStr
	}

	apiRules := map[string]interface{}{
		"rules": []interface{}{rule},
	}

	return utils.JsonToString(apiRules)
}

func resourceObsBucketMirrorBackToSourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		bucket = d.Get("bucket").(string)
		ruleID = utils.PathSearch("id", utils.StringToJson(d.Get("rule").(string)), "").(string)
	)

	obsClient, err := cfg.ObjectStorageClientWithSignature(region)
	if err != nil {
		return diag.Errorf("error creating OBS Client: %s", err)
	}

	opts := &obs.SetBucketMirrorBackToSourceInput{
		Bucket: bucket,
		Rules:  buildMirrorBackToSourceRules(d.Get("rule").(string)),
	}

	_, err = obsClient.SetBucketMirrorBackToSource(opts)
	if err != nil {
		return diag.FromErr(getObsError("error setting mirror back to source of OBS bucket", bucket, err))
	}

	d.SetId(ruleID)

	return resourceObsBucketMirrorBackToSourceRead(ctx, d, meta)
}

func resourceObsBucketMirrorBackToSourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		bucket = d.Get("bucket").(string)
		ruleID = utils.PathSearch("id", utils.StringToJson(d.Get("rule").(string)), "").(string)
	)

	obsClient, err := cfg.ObjectStorageClientWithSignature(region)
	if err != nil {
		return diag.Errorf("error creating OBS Client: %s", err)
	}

	opts := &obs.SetBucketMirrorBackToSourceInput{
		Bucket: bucket,
		Rules:  buildMirrorBackToSourceRules(d.Get("rule").(string)),
	}

	_, err = obsClient.SetBucketMirrorBackToSource(opts)
	if err != nil {
		return diag.FromErr(getObsError("Error updating mirror back to source of OBS bucket", bucket, err))
	}

	d.SetId(ruleID)

	return resourceObsBucketMirrorBackToSourceRead(ctx, d, meta)
}

func parseBucketAndRuleID(id string) (bucket, ruleId string, err error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return "", "", errors.New("invalid ID format, must be <bucket>/<id>")
	}

	bucket = parts[0]
	ruleId = parts[1]

	return bucket, ruleId, nil
}

func resourceObsBucketMirrorBackToSourceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		bucket = d.Get("bucket").(string)
		ruleID = d.Id()
	)

	obsClient, err := cfg.ObjectStorageClientWithSignature(region)
	if err != nil {
		return diag.Errorf("error creating OBS Client: %s", err)
	}

	output, err := obsClient.GetBucketMirrorBackToSource(bucket)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving OBS bucket mirror back to source")
	}

	ruleMap := utils.PathSearch(fmt.Sprintf("rules[?id=='%s']|[0]", ruleID), utils.StringToJson(output.Rules), nil)
	if ruleMap == nil {
		return diag.FromErr(golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/{bucket}?mirrorBackToSource",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the mirror back to source rule '%s' has been removed", ruleID)),
			},
		})
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("bucket", bucket),
		d.Set("rule", utils.JsonToString(ruleMap)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceObsBucketMirrorBackToSourceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		bucket = d.Get("bucket").(string)
	)

	obsClient, err := cfg.ObjectStorageClientWithSignature(region)
	if err != nil {
		return diag.Errorf("error creating OBS Client: %s", err)
	}

	_, err = obsClient.DeleteBucketMirrorBackToSource(bucket)
	if err != nil {
		return diag.FromErr(getObsError("Error deleting mirror back to source of OBS bucket", bucket, err))
	}

	return nil
}

func resourceObsBucketMirrorBackToSourceImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	bucket, ruleID, err := parseBucketAndRuleID(d.Id())
	if err != nil {
		return nil, err
	}

	d.SetId(ruleID)

	return []*schema.ResourceData{d}, multierror.Append(nil,
		d.Set("bucket", bucket),
	).ErrorOrNil()
}
