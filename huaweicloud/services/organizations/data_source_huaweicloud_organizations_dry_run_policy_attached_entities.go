package organizations

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Organizations GET /v1/organizations/dry-run-policies/{policy_id}/attached-entities
func DataSourceDryRunPolicyAttachedEntities() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDryRunPolicyAttachedEntitiesRead,
		Schema: map[string]*schema.Schema{
			// Required parameters.
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dry run policy.`,
			},

			// Attributes.
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The unique ID of the entity.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the entity.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the entity.`,
						},
					},
				},
				Description: `The entities that are attached to the specified dry run policy.`,
			},
		},
	}
}

func dataSourceDryRunPolicyAttachedEntitiesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("organizations", region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	policyId := d.Get("policy_id").(string)
	attachedEntities, err := listAttachedEntitiesForDryRunPolicy(client, policyId)
	if err != nil {
		return diag.Errorf("error retrieving attached entities for dry run policy (%s): %s", policyId, err)
	}

	dataSourceID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceID)

	return diag.FromErr(d.Set("entities", flattenDryRunPolicyAttachedEntities(attachedEntities)))
}

func flattenDryRunPolicyAttachedEntities(attachedEntities []interface{}) []interface{} {
	if len(attachedEntities) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(attachedEntities))
	for _, entity := range attachedEntities {
		rst = append(rst, map[string]interface{}{
			"id":   utils.PathSearch("id", entity, nil),
			"type": utils.PathSearch("type", entity, nil),
			"name": utils.PathSearch("name", entity, nil),
		})
	}

	return rst
}
