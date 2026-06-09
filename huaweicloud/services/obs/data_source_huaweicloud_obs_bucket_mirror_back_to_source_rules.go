package obs

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API OBS GET /?mirrorBackToSource
func DataSourceObsBucketMirrorBackToSourceRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceObsBucketMirrorBackToSourceRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the OBS bucket mirror back to source rules are located.`,
			},
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the bucket.`,
			},
			"rules": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The rules of the bucket mirror back to source, in JSON format.`,
			},
		},
	}
}

func dataSourceObsBucketMirrorBackToSourceRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.ObjectStorageClient(region)
	if err != nil {
		return diag.Errorf("error creating OBS client: %s", err)
	}

	bucketName := d.Get("bucket").(string)
	resp, err := client.GetBucketMirrorBackToSource(bucketName)
	if err != nil {
		return diag.Errorf("error querying OBS bucket mirror back to source rules: %s", err)
	}

	d.SetId(bucketName)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rules", resp.Rules),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
