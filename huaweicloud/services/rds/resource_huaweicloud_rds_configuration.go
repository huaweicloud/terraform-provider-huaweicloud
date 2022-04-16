package rds

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/rds/v3/configurations"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// ResourceRdsConfiguration is the impl for huaweicloud_rds_parametergroup resource
func ResourceRdsConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsConfigurationCreate,
		ReadContext:   resourceRdsConfigurationRead,
		UpdateContext: resourceRdsConfigurationUpdate,
		DeleteContext: resourceRdsConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"values": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"datastore": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: common.CaseInsensitiveFunc(),
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"configuration_parameters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"restart_required": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"readonly": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"value_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
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

func buildValues(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("values").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func buildDatastore(d *schema.ResourceData) configurations.DataStore {
	datastoreRaw := d.Get("datastore").([]interface{})
	rawMap := datastoreRaw[0].(map[string]interface{})

	datastore := configurations.DataStore{
		Type:    rawMap["type"].(string),
		Version: rawMap["version"].(string),
	}

	log.Printf("[DEBUG] Datastore: %#v", datastore)
	return datastore
}

func resourceRdsConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	rdsClient, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RDS Client: %s", err)
	}

	createOpts := configurations.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Values:      buildValues(d),
		DataStore:   buildDatastore(d),
	}

	log.Printf("[DEBUG] CreateOpts: %#v", createOpts)
	configuration, err := configurations.Create(rdsClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating RDS configuration: %s", err)
	}

	log.Printf("[DEBUG] RDS configuration created: %#v", configuration)
	d.SetId(configuration.Id)

	return resourceRdsConfigurationRead(ctx, d, meta)
}

func resourceRdsConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	rdsClient, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	n, err := configurations.Get(rdsClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS configuration")
	}

	d.SetId(n.Id)
	d.Set("name", n.Name)
	d.Set("description", n.Description)

	datastore := []map[string]string{
		{
			"type":    n.DatastoreName,
			"version": n.DatastoreVersionName,
		},
	}
	d.Set("datastore", datastore)

	parameters := make([]map[string]interface{}, len(n.Parameters))
	for i, parameter := range n.Parameters {
		parameters[i] = make(map[string]interface{})
		parameters[i]["name"] = parameter.Name
		parameters[i]["value"] = parameter.Value
		parameters[i]["restart_required"] = parameter.RestartRequired
		parameters[i]["readonly"] = parameter.ReadOnly
		parameters[i]["value_range"] = parameter.ValueRange
		parameters[i]["type"] = parameter.Type
		parameters[i]["description"] = parameter.Description
	}
	d.Set("configuration_parameters", parameters)
	return nil
}

func resourceRdsConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	rdsClient, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RDS Client: %s", err)
	}

	var updateOpts configurations.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}
	if d.HasChange("values") {
		updateOpts.Values = buildValues(d)
	}

	log.Printf("[DEBUG] updateOpts: %#v", updateOpts)
	err = configurations.Update(rdsClient, d.Id(), updateOpts).ExtractErr()
	if err != nil {
		return diag.Errorf("error updating RDS configuration: %s", err)
	}
	return resourceRdsConfigurationRead(ctx, d, meta)
}

func resourceRdsConfigurationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	rdsClient, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	err = configurations.Delete(rdsClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting RDS configuration: %s", err)
	}

	d.SetId("")
	return nil
}
