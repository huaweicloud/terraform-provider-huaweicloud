package cse

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/cse/dedicated/v4/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	internalPropertyKeys  = []string{"engineID", "engineName"}
	instanceNotFoundCodes = []string{
		"400017",
	}
)

// @API CSE GET /v2/{project_id}/registry/microservices/{serviceId}/instances/{instanceId}
// @API CSE DELETE /v2/{project_id}/registry/microservices/{serviceId}/instances/{instanceId}
// @API CSE POST /v2/{project_id}/registry/microservices/{serviceId}/instances
func ResourceMicroserviceInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMicroserviceInstanceCreate,
		ReadContext:   resourceMicroserviceInstanceRead,
		DeleteContext: resourceMicroserviceInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceMicroserviceInstanceImportState,
		},

		Schema: map[string]*schema.Schema{
			// Authentication and request parameters.
			"auth_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Description: utils.SchemaDesc(
					`The address that used to request the access token.`,
					utils.SchemaDescInput{
						Required: true,
					}),
			},
			"connect_address": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The address that used to send requests and manage configuration.`,
			},
			"admin_user": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The user name that used to pass the RBAC control.",
			},
			"admin_pass": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ForceNew:     true,
				RequiredWith: []string{"admin_user"},
				Description:  "The user password that used to pass the RBAC control.",
			},
			// Resource parameters.
			"microservice_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoints": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"properties": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"health_check": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"interval": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"max_retries": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"data_center": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"region": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildHealthCheckStructure(healthChecks []interface{}) *instances.HealthCheck {
	if len(healthChecks) < 1 {
		return nil
	}

	healthCheck := healthChecks[0].(map[string]interface{})

	return &instances.HealthCheck{
		Mode:     healthCheck["mode"].(string),
		Interval: healthCheck["interval"].(int),
		Times:    healthCheck["max_retries"].(int),
		Port:     healthCheck["port"].(int),
	}
}

func buildDataCenterStructure(dataCenters []interface{}) *instances.DataCenter {
	if len(dataCenters) < 1 {
		return nil
	}

	dataCenter := dataCenters[0].(map[string]interface{})

	return &instances.DataCenter{
		Name:          dataCenter["name"].(string),
		Region:        dataCenter["region"].(string),
		AvailableZone: dataCenter["availability_zone"].(string),
	}
}

func buildCustomProperties(properties map[string]interface{}) map[string]interface{} {
	if len(properties) < 1 {
		return nil
	}

	result := make(map[string]interface{})
	for k, v := range properties {
		if !utils.StrSliceContains(internalPropertyKeys, k) {
			result[k] = v
		}
	}

	return result
}

func buildInstanceCreateOpts(d *schema.ResourceData) instances.CreateOpts {
	return instances.CreateOpts{
		HostName:       d.Get("host_name").(string),
		Endpoints:      utils.ExpandToStringList(d.Get("endpoints").([]interface{})),
		Version:        d.Get("version").(string),
		Status:         d.Get("status").(string),
		Properties:     buildCustomProperties(d.Get("properties").(map[string]interface{})),
		HealthCheck:    buildHealthCheckStructure(d.Get("health_check").([]interface{})),
		DataCenterInfo: buildDataCenterStructure(d.Get("data_center").([]interface{})),
	}
}

func resourceMicroserviceInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token, err := GetAuthorizationToken(getAuthAddress(d), d.Get("admin_user").(string),
		d.Get("admin_pass").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	client := common.NewCustomClient(true, d.Get("connect_address").(string), "v4", "default")
	createOpts := buildInstanceCreateOpts(d)
	log.Printf("[DEBUG] The createOpts of the Microservice instance is: %v", createOpts)
	resp, err := instances.Create(client, createOpts, d.Get("microservice_id").(string), token)
	if err != nil {
		return diag.Errorf("error creating microservice instance: %s", err)
	}
	d.SetId(resp.ID)

	return resourceMicroserviceInstanceRead(ctx, d, meta)
}

func flattenHealthCheck(healthCheck instances.HealthCheck) (result []map[string]interface{}) {
	if reflect.DeepEqual(healthCheck, instances.HealthCheck{}) {
		return nil
	}

	result = append(result, map[string]interface{}{
		"mode":        healthCheck.Mode,
		"interval":    healthCheck.Interval,
		"max_retries": healthCheck.Times,
		"port":        healthCheck.Port,
	})

	log.Printf("[DEBUG] The health check result is %#v", result)
	return
}

func flattenDataCenter(dataCenter instances.DataCenter) (result []map[string]interface{}) {
	if reflect.DeepEqual(dataCenter, instances.DataCenter{}) {
		return nil
	}

	result = append(result, map[string]interface{}{
		"name":              dataCenter.Name,
		"region":            dataCenter.Region,
		"availability_zone": dataCenter.AvailableZone,
	})

	log.Printf("[DEBUG] The data center result is %#v", result)
	return
}

func resourceMicroserviceInstanceRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	token, err := GetAuthorizationToken(getAuthAddress(d), d.Get("admin_user").(string),
		d.Get("admin_pass").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	client := common.NewCustomClient(true, d.Get("connect_address").(string), "v4", "default")
	resp, err := instances.Get(client, d.Get("microservice_id").(string), d.Id(), token)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "errorCode", instanceNotFoundCodes...),
			"error retrieving Microservice instance")
	}

	mErr := multierror.Append(nil,
		d.Set("host_name", resp.HostName),
		d.Set("endpoints", resp.Endpoints),
		d.Set("version", resp.Version),
		d.Set("properties", buildCustomProperties(resp.Properties)),
		d.Set("health_check", flattenHealthCheck(resp.HealthCheck)),
		d.Set("data_center", flattenDataCenter(resp.DataCenterInfo)),
		d.Set("status", resp.Status),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceMicroserviceInstanceDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	token, err := GetAuthorizationToken(getAuthAddress(d), d.Get("admin_user").(string),
		d.Get("admin_pass").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	client := common.NewCustomClient(true, d.Get("connect_address").(string), "v4", "default")
	err = instances.Delete(client, d.Get("microservice_id").(string), d.Id(), token)
	if err != nil {
		return diag.Errorf("error deleting dedicated microservice instance (%s): %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func resourceMicroserviceInstanceImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	var (
		authAddr, connectAddr, importedIdWithoutAddrs, microserviceId, instanceId, adminUser, adminPwd string
		mErr                                                                                           *multierror.Error

		importedId   = d.Id()
		addressRegex = `https://\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5}`
		re           = regexp.MustCompile(fmt.Sprintf(`^(%[1]s)?/?(%[1]s)/(.*)$`, addressRegex))
		formatErr    = fmt.Errorf("the imported microservice ID specifies an invalid format, want "+
			"'<auth_address>/<connect_address>/<microservice_id>/<id>' or "+
			"'<auth_address>/<connect_address>/<microservice_id>/<id>/<admin_user>/<admin_pass>', but got '%s'",
			importedId)
	)

	if !re.MatchString(importedId) {
		return nil, formatErr
	}
	resp := re.FindAllStringSubmatch(importedId, -1)
	// If the imported ID matches the address regular expression, the length of the response result must be greater than 1.
	switch len(resp[0]) {
	case 4:
		authAddr = resp[0][1]
		connectAddr = resp[0][2]
		importedIdWithoutAddrs = resp[0][3]
		if authAddr == "" {
			authAddr = connectAddr // Using the connect address as the auth address if the auth address input is omitted.
		}
	default:
		return nil, formatErr
	}

	mErr = multierror.Append(mErr,
		d.Set("auth_address", authAddr),
		d.Set("connect_address", connectAddr),
	)

	parts := strings.Split(importedIdWithoutAddrs, "/")
	switch len(parts) {
	case 2:
		microserviceId = parts[0]
		instanceId = parts[1]
	case 4:
		microserviceId = parts[0]
		instanceId = parts[1]
		adminUser = parts[2]
		adminPwd = parts[3]

		mErr = multierror.Append(mErr,
			d.Set("admin_user", adminUser),
			d.Set("admin_pass", adminPwd),
		)
	default:
		return nil, formatErr
	}

	mErr = multierror.Append(mErr, d.Set("microservice_id", microserviceId))
	d.SetId(instanceId)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
