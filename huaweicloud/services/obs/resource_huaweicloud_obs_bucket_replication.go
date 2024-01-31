package obs

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/obs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API OBS PUT ?replication
// @API OBS DELETE ?replication
// @API OBS GET ?replication
func ResourceObsBucketReplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceObsBucketReplicationCreate,
		UpdateContext: resourceObsBucketReplicationCreate,
		ReadContext:   resourceObsBucketReplicationRead,
		DeleteContext: resourceObsBucketReplicationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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
				ForceNew: true,
			},
			"destination_bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"agency": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"storage_class": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"STANDARD", "WARM", "COLD",
							}, false),
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"history_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildReplicationRuleFromRawMap(rawMap map[string]interface{}, destBucket, prefix string) obs.ReplicationRule {
	replicationRule := obs.ReplicationRule{
		Prefix:            prefix,
		DestinationBucket: destBucket,
	}

	if enabled, ok := rawMap["enabled"].(bool); ok && enabled {
		replicationRule.Status = obs.RuleStatusEnabled
	} else {
		replicationRule.Status = obs.RuleStatusDisabled
	}

	if historyEnabled, ok := rawMap["history_enabled"].(bool); ok && historyEnabled {
		replicationRule.HistoricalObjectReplication = obs.Enabled
	} else {
		replicationRule.HistoricalObjectReplication = obs.Disabled
	}

	if val, ok := rawMap["storage_class"].(string); ok {
		replicationRule.StorageClass = obs.ParseStringToStorageClassType(val)
	}
	return replicationRule
}

func buildReplicationRules(d *schema.ResourceData) ([]obs.ReplicationRule, error) {
	destBucket := d.Get("destination_bucket").(string)
	rawArray := d.Get("rule").([]interface{})

	replicationRules := make([]obs.ReplicationRule, 0, len(rawArray))
	for _, raw := range rawArray {
		if rawMap, rawOk := raw.(map[string]interface{}); rawOk {
			prefix := rawMap["prefix"].(string)
			if prefix == "" && len(rawArray) > 1 {
				return nil, fmt.Errorf("to apply a rule to all objects, delete all rules that take effect" +
					" by prefixes first")
			}
			replicationRules = append(replicationRules, buildReplicationRuleFromRawMap(rawMap, destBucket, prefix))
		}
	}

	if len(replicationRules) == 0 {
		replicationRules = append(replicationRules, obs.ReplicationRule{
			Status:            obs.RuleStatusEnabled,
			DestinationBucket: destBucket,
		})
	}
	return replicationRules, nil
}

func resourceObsBucketReplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	obsClient, err := cfg.ObjectStorageClientWithSignature(region)
	if err != nil {
		return diag.Errorf("error creating OBS Client: %s", err)
	}

	bucket := d.Get("bucket").(string)
	replicationRules, err := buildReplicationRules(d)
	if err != nil {
		return diag.FromErr(err)
	}
	opts := &obs.SetBucketReplicationInput{
		Bucket: bucket,
		BucketReplicationConfiguration: obs.BucketReplicationConfiguration{
			Agency:           d.Get("agency").(string),
			ReplicationRules: replicationRules,
		},
	}
	_, err = obsClient.SetBucketReplication(opts)
	if err != nil {
		return diag.FromErr(getObsError("Error creating cross-region replication of OBS bucket", bucket, err))
	}

	// Assign the source bucket name as the resource ID
	d.SetId(bucket)
	return resourceObsBucketReplicationRead(ctx, d, meta)
}

func flattenDestinationBucket(output *obs.GetBucketReplicationOutput) string {
	rules := output.ReplicationRules
	if len(rules) == 0 {
		return ""
	}
	return rules[0].DestinationBucket
}

func flattenDestinationRules(output *obs.GetBucketReplicationOutput) []map[string]interface{} {
	rules := output.ReplicationRules
	if len(rules) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(rules))
	for i, rule := range rules {
		result[i] = map[string]interface{}{
			"prefix":          rule.Prefix,
			"storage_class":   rule.StorageClass,
			"enabled":         rule.Status == obs.RuleStatusEnabled,
			"history_enabled": rule.HistoricalObjectReplication == obs.Enabled,
			"id":              rule.ID,
		}
	}
	return result
}

func resourceObsBucketReplicationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	obsClient, err := cfg.ObjectStorageClientWithSignature(region)
	if err != nil {
		return diag.Errorf("error creating OBS Client: %s", err)
	}

	output, err := obsClient.GetBucketReplication(d.Id())
	if err != nil {
		if obsError, ok := err.(obs.ObsError); ok && obsError.Code == "ReplicationConfigurationNotFoundError" {
			// The bucket does not have replication configurations
			d.SetId("")
			return nil
		}
		return diag.FromErr(getObsError("Error retrieving OBS bucket replication", d.Id(), err))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("bucket", d.Id()),
		d.Set("destination_bucket", flattenDestinationBucket(output)),
		d.Set("agency", output.Agency),
		d.Set("rule", flattenDestinationRules(output)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting OBS bucket replication fields: %s", err)
	}
	return nil
}

func resourceObsBucketReplicationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	obsClient, err := cfg.ObjectStorageClientWithSignature(region)
	if err != nil {
		return diag.Errorf("error creating OBS Client: %s", err)
	}

	bucket := d.Id()
	log.Printf("[DEBUG] delete cross-region replication configuration of OBS bucket %s", bucket)
	_, err = obsClient.DeleteBucketReplication(bucket)
	if err != nil {
		return diag.FromErr(getObsError("Error deleting cross-region replication "+
			"configuration of OBS bucket", bucket, err))
	}

	return nil
}
