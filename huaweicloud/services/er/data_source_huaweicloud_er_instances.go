package er

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/er/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ER GET /v3/{project_id}/enterprise-router/instances
func DataSourceInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The region where the ER instances are located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID used to query specified instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name used to filter the instances.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The enterprise project ID of the instances to be queried.`,
			},
			"owned_by_self": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether resources belong to the current renant.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The status used to filter the instances.`,
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The key/value pairs used to filter the instances.`,
			},
			// Attributes
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The instance ID.`,
						},
						"asn": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The BGP AS number of the ER instance.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the instance.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the instance.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The current status of the instance.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of enterprise project to which the instance belongs.`,
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The key/value pairs to associate with the instance.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the instance.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The last update time of the instance.`,
						},
						"enable_default_propagation": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether to enable the propagation of the default route table.`,
						},
						"enable_default_association": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether to enable the association of the default route table.`,
						},
						"auto_accept_shared_attachments": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether to automatically accept the creation of shared attachment.`,
						},
						"default_propagation_route_table_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the default propagation route table.`,
						},
						"default_association_route_table_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the default association route table.`,
						},
						"availability_zones": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The availability zone list where the ER instance is located.`,
						},
					},
				},
				Description: `All instances that match the filter parameters.`,
			},
		},
	}
}

// Filter instances by name and tags.
func filterInstances(d *schema.ResourceData, all []instances.Instance) ([]instances.Instance, error) {
	filter := map[string]interface{}{}
	if name, ok := d.GetOk("name"); ok {
		filter["Name"] = name
	}
	filterResult, err := utils.FilterSliceWithField(all, filter)
	if err != nil {
		return nil, fmt.Errorf("error filting instance list: %s", err)
	}

	result := make([]instances.Instance, 0, len(filterResult))
	tagFilter := d.Get("tags").(map[string]interface{})
	for _, val := range filterResult {
		item := val.(instances.Instance)
		tagmap := utils.TagsToMap(item.Tags)

		// Filter instances list by tags, if the filter is nil, skip and return all fileterResult elements.
		if !utils.HasMapContains(tagmap, tagFilter) {
			continue
		}
		result = append(result, item)
	}
	return result, nil
}

func flattenInstances(all []instances.Instance) []map[string]interface{} {
	if len(all) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(all))
	for i, instance := range all {
		result[i] = map[string]interface{}{
			"id":                    instance.ID,
			"asn":                   instance.ASN,
			"name":                  instance.Name,
			"description":           instance.Description,
			"status":                instance.Status,
			"enterprise_project_id": instance.EnterpriseProjectId,
			"tags":                  utils.TagsToMap(instance.Tags),
			// The time results are not the time in RF3339 format without milliseconds.
			"created_at":                         utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(instance.CreatedAt)/1000, false),
			"updated_at":                         utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(instance.UpdatedAt)/1000, false),
			"enable_default_propagation":         instance.EnableDefaultPropagation,
			"enable_default_association":         instance.EnableDefaultAssociation,
			"auto_accept_shared_attachments":     instance.AutoAcceptSharedAttachments,
			"default_propagation_route_table_id": instance.DefaultPropagationRouteTableId,
			"default_association_route_table_id": instance.DefaultAssociationRouteTableId,
			"availability_zones":                 instance.AvailabilityZoneIds,
		}
	}
	return result
}

func buildSliceIgnoreEmptyElement(e string) []string {
	if e != "" {
		return []string{e}
	}
	return nil
}

func buildInstanceListOpts(d *schema.ResourceData) instances.ListOpts {
	return instances.ListOpts{
		EnterpriseProjectIds: buildSliceIgnoreEmptyElement(d.Get("enterprise_project_id").(string)),
		Statuses:             buildSliceIgnoreEmptyElement(d.Get("status").(string)),
		IDs:                  buildSliceIgnoreEmptyElement(d.Get("instance_id").(string)),
		OwnedBySelf:          d.Get("owned_by_self").(bool),
		SortKey:              []string{"name"},
	}
}

func dataSourceInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ErV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ER v3 client: %s", err)
	}

	resp, err := instances.List(client, buildInstanceListOpts(d))
	if err != nil {
		return diag.Errorf("error retrieving instances: %s", err)
	}
	if resp, err = filterInstances(d, resp); err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instances", flattenInstances(resp)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving instance data source fields: %s", mErr)
	}
	return nil
}
