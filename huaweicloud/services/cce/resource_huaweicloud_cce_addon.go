package cce

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cce/v3/addons"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE DELETE /api/v3/addons/{id}
// @API CCE GET /api/v3/addons/{id}
// @API CCE PUT /api/v3/addons/{id}
// @API CCE POST /api/v3/addons
func ResourceAddon() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAddonCreate,
		ReadContext:   resourceAddonRead,
		UpdateContext: resourceAddonUpdate,
		DeleteContext: resourceAddonDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAddonImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ // request and response parameters
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"values": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"basic": {
							Type:         schema.TypeMap,
							Optional:     true,
							Elem:         &schema.Schema{Type: schema.TypeString},
							ExactlyOneOf: []string{"values.0.basic", "values.0.basic_json"},
						},
						"basic_json": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								equal, _ := utils.CompareJsonTemplateAreEquivalent(old, new)
								return equal
							},
							ExactlyOneOf: []string{"values.0.basic", "values.0.basic_json"},
						},
						"custom": {
							Type:          schema.TypeMap,
							Optional:      true,
							Elem:          &schema.Schema{Type: schema.TypeString},
							ConflictsWith: []string{"values.0.custom_json"},
						},
						"custom_json": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								equal, _ := utils.CompareJsonTemplateAreEquivalent(old, new)
								return equal
							},
							ConflictsWith: []string{"values.0.custom"},
						},
						"flavor": {
							Type:          schema.TypeMap,
							Optional:      true,
							Elem:          &schema.Schema{Type: schema.TypeString},
							ConflictsWith: []string{"values.0.flavor_json"},
						},
						"flavor_json": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								equal, _ := utils.CompareJsonTemplateAreEquivalent(old, new)
								return equal
							},
							ConflictsWith: []string{"values.0.flavor"},
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildAddonValues(d *schema.ResourceData) (basic, custom, flavor map[string]interface{}, err error) {
	values := d.Get("values").([]interface{})
	if len(values) == 0 || values[0] == nil {
		basic = make(map[string]interface{})
		return
	}

	valuesMap := values[0].(map[string]interface{})
	if basicRaw := valuesMap["basic"].(map[string]interface{}); len(basicRaw) != 0 {
		basic = basicRaw
	} else if basicJsonRaw := valuesMap["basic_json"].(string); basicJsonRaw != "" {
		err = json.Unmarshal([]byte(basicJsonRaw), &basic)
		if err != nil {
			err = fmt.Errorf("error unmarshalling basic json: %s", err)
			return
		}
	}

	if customRaw := valuesMap["custom"].(map[string]interface{}); len(customRaw) != 0 {
		custom = customRaw
	} else if customJsonRaw := valuesMap["custom_json"].(string); customJsonRaw != "" {
		err = json.Unmarshal([]byte(customJsonRaw), &custom)
		if err != nil {
			err = fmt.Errorf("error unmarshalling custom json: %s", err)
			return
		}
	}

	if flavorRaw := valuesMap["flavor"].(map[string]interface{}); len(flavorRaw) != 0 {
		flavor = flavorRaw
	} else if flavorJsonRaw := valuesMap["flavor_json"].(string); flavorJsonRaw != "" {
		err = json.Unmarshal([]byte(flavorJsonRaw), &flavor)
		if err != nil {
			err = fmt.Errorf("error unmarshalling flavor json %s", err)
			return
		}
	}

	return
}

func resourceAddonCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	cceClient, err := cfg.CceAddonV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE v3 Client (without project): %s", err)
	}

	clusterID := d.Get("cluster_id").(string)
	basic, custom, flavor, err := buildAddonValues(d)
	if err != nil {
		return diag.Errorf("error getting values for CCE add-on: %s", err)
	}

	createOpts := addons.CreateOpts{
		Kind:       "Addon",
		ApiVersion: "v3",
		Metadata: addons.CreateMetadata{
			Anno: addons.CreateAnnotations{
				AddonInstallType: "install",
			},
		},
		Spec: addons.RequestSpec{
			Version:           d.Get("version").(string),
			ClusterID:         clusterID,
			AddonTemplateName: d.Get("template_name").(string),
			Values: addons.Values{
				Basic:  basic,
				Custom: custom,
				Flavor: flavor,
			},
		},
	}

	create, err := addons.Create(cceClient, createOpts, clusterID).Extract()
	if err != nil {
		return diag.Errorf("error creating CCE add-on: %s", err)
	}

	d.SetId(create.Metadata.Id)

	log.Printf("[DEBUG] Waiting for CCE add-on (%s) to become available", create.Metadata.Id)
	stateConf := &resource.StateChangeConf{
		// The statuses of pending phase includes "installing" and "abnormal".
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      addonStateRefreshFunc(cceClient, create.Metadata.Id, clusterID, []string{"running", "available"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for CCE add-on (%s) to become available: %s", create.Metadata.Id, err)
	}

	return resourceAddonRead(ctx, d, meta)
}

func resourceAddonRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	cceClient, err := cfg.CceAddonV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE v3 Client (without project): %s", err)
	}

	clusterID := d.Get("cluster_id").(string)
	addonID := d.Id()
	n, err := addons.Get(cceClient, addonID, clusterID).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CCE add-on")
	}

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("cluster_id", n.Spec.ClusterID),
		d.Set("version", n.Spec.Version),
		d.Set("template_name", n.Spec.AddonTemplateName),
		d.Set("status", n.Status.Status),
		d.Set("description", n.Spec.Description),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting CCE add-on (%s) fields: %s", addonID, err)
	}

	return nil
}

func resourceAddonUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	cceClient, err := cfg.CceAddonV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE v3 Client (without project): %s", err)
	}

	clusterID := d.Get("cluster_id").(string)
	addonID := d.Id()
	basic, custom, flavor, err := buildAddonValues(d)
	if err != nil {
		return diag.Errorf("error getting values for CCE add-on: %s", err)
	}

	updateOpts := addons.UpdateOpts{
		Kind:       "Addon",
		ApiVersion: "v3",
		Metadata: addons.UpdateMetadata{
			Anno: addons.UpdateAnnotations{
				AddonUpgradeType: "upgrade",
			},
		},
		Spec: addons.RequestSpec{
			Version:           d.Get("version").(string),
			ClusterID:         clusterID,
			AddonTemplateName: d.Get("template_name").(string),
			Values: addons.Values{
				Basic:  basic,
				Custom: custom,
				Flavor: flavor,
			},
		},
	}

	_, err = addons.Update(cceClient, updateOpts, addonID, clusterID).Extract()
	if err != nil {
		return diag.Errorf("error updating CCE add-on (%s): %s", addonID, err)
	}

	stateConf := &resource.StateChangeConf{
		// The statuses of pending phase includes "installing" and "abnormal".
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      addonStateRefreshFunc(cceClient, addonID, clusterID, []string{"running", "available"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for CCE add-on (%s) to become available: %s", addonID, err)
	}

	return resourceAddonRead(ctx, d, meta)
}

func resourceAddonDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	cceClient, err := cfg.CceAddonV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE v3 Client (without project): %s", err)
	}

	clusterID := d.Get("cluster_id").(string)
	addonID := d.Id()
	err = addons.Delete(cceClient, addonID, clusterID).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting CCE add-on: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		// The statuses of pending phase includes "Deleting", "Available" and "Unavailable".
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      addonStateRefreshFunc(cceClient, addonID, clusterID, nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for CCE add-on (%s) to become deleted: %s", addonID, err)
	}
	return nil
}

func addonStateRefreshFunc(cceClient *golangsdk.ServiceClient, addonId, clusterId string,
	targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Expect the status of CCE add-on to be any one of the status list: %v.", targets)
		resp, err := addons.Get(cceClient, addonId, clusterId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] The add-on (%s) has been deleted", addonId)
				return resp, "COMPLETED", nil
			}
			return nil, "ERROR", err
		}

		invalidStatuses := []string{"installFailed", "upgradeFailed", "deleteFailed", "rollbackFailed", "unknown"}
		if utils.IsStrContainsSliceElement(resp.Status.Status, invalidStatuses, true, true) {
			return resp, "", fmt.Errorf("unexpect status (%s)", resp.Status.Status)
		}

		if utils.StrSliceContains(targets, resp.Status.Status) {
			return resp, "COMPLETED", nil
		}
		return resp, "PENDING", nil
	}
}

func resourceAddonImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()
	parts := strings.Split(importId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for CCE add-on, want '<cluster_id>/<id>', but got '%s'", importId)
	}

	clusterID := parts[0]
	addonID := parts[1]

	d.SetId(addonID)
	err := d.Set("cluster_id", clusterID)
	return []*schema.ResourceData{d}, err
}
