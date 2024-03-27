package ddm

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API DDM PUT /v2/{project_id}/instances/{instance_id}/action/read-write-strategy
func ResourceDdmInstanceReadStrategy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDdmReadStrategyCreateOrUpdate,
		ReadContext:   resourceDdmReadStrategyRead,
		UpdateContext: resourceDdmReadStrategyCreateOrUpdate,
		DeleteContext: resourceDdmReadStrategyDelete,
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
				ForceNew:    true,
				Description: "Specifies the ID of the DDM instance.",
			},
			"read_weights": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the ID of the DB instance associated with the DDM schema.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Specifies read weight of the DB instance associated with the DDM schema.",
						},
					},
				},
				Required:    true,
				Description: `Specifies the list of read weights of the primary DB instance and its read replicas.`,
			},
		},
	}
}

func resourceDdmReadStrategyCreateOrUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		dbReadStrategyHttpUrl = "v2/{project_id}/instances/{instance_id}/action/read-write-strategy"
		dbReadStrategyProduct = "ddm"
	)
	dbReadStrategyClient, err := cfg.NewServiceClient(dbReadStrategyProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDM client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	dbReadStrategyPath := dbReadStrategyClient.Endpoint + dbReadStrategyHttpUrl
	dbReadStrategyPath = strings.ReplaceAll(dbReadStrategyPath, "{project_id}", dbReadStrategyClient.ProjectID)
	dbReadStrategyPath = strings.ReplaceAll(dbReadStrategyPath, "{instance_id}", fmt.Sprintf("%v", instanceID))

	dbReadStrategyOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	dbReadStrategyOpt.JSONBody = buildDbReadStrategyBodyParams(d)
	_, err = dbReadStrategyClient.Request("PUT", dbReadStrategyPath, &dbReadStrategyOpt)
	if err != nil {
		return diag.Errorf("error setting read strategy of the DDM instance: %s", err)
	}
	d.SetId(instanceID)

	return nil
}

func buildDbReadStrategyBodyParams(d *schema.ResourceData) map[string]interface{} {
	readWeightMap := make(map[string]interface{})
	readWeightList := d.Get("read_weights").(*schema.Set).List()
	for _, readWeight := range readWeightList {
		variable := readWeight.(map[string]interface{})
		readWeightMap[variable["db_id"].(string)] = variable["weight"]
	}

	bodyParams := map[string]interface{}{
		"read_weight": readWeightMap,
	}
	return bodyParams
}

func resourceDdmReadStrategyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDdmReadStrategyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting read strategy is not supported. The read strategy is only removed from the state," +
		" but it remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
