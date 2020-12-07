package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/csbs/v1/policies"
)

func dataSourceCSBSBackupPolicyV1() *schema.Resource {
	return &schema.Resource{
		Read:               dataSourceCSBSBackupPolicyV1Read,
		DeprecationMessage: "It has been deprecated.",
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"common": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"scheduled_operation": {
				Type:     schema.TypeSet,
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
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"max_backups": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"retention_duration_days": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"permanent": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"trigger_pattern": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operation_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCSBSBackupPolicyV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	policyClient, err := config.CsbsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating csbs client: %s", err)
	}

	listOpts := policies.ListOpts{
		ID:     d.Get("id").(string),
		Name:   d.Get("name").(string),
		Status: d.Get("status").(string),
	}

	refinedPolicies, err := policies.List(policyClient, listOpts)

	if err != nil {
		return fmt.Errorf("Unable to retrieve backup policies: %s", err)
	}

	if len(refinedPolicies) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedPolicies) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	backupPolicy := refinedPolicies[0]

	log.Printf("[INFO] Retrieved backup policy %s using given filter", backupPolicy.ID)

	d.SetId(backupPolicy.ID)

	if err := d.Set("resource", flattenCSBSPolicyResources(backupPolicy)); err != nil {
		return err
	}

	if err := d.Set("scheduled_operation", flattenCSBSScheduledOperations(backupPolicy)); err != nil {
		return err
	}

	d.Set("name", backupPolicy.Name)
	d.Set("id", backupPolicy.ID)
	d.Set("common", backupPolicy.Parameters.Common)
	d.Set("status", backupPolicy.Status)
	d.Set("description", backupPolicy.Description)
	d.Set("provider_id", backupPolicy.ProviderId)

	d.Set("region", GetRegion(d, config))

	return nil
}
