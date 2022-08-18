package cse

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cse/dedicated/v2/engines"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var DefaultVersion = "CSE2"

func ResourceMicroserviceEngine() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMicroserviceEngineCreate,
		ReadContext:   resourceMicroserviceEngineRead,
		DeleteContext: resourceMicroserviceEngineDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
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
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z]([A-Za-z0-9-]*[A-Za-z0-9])?$`),
						"The name must start a letter and cannot end with a hyphen (-), and can only contain "+
							"letters, digits and hyphens (-)."),
					validation.StringLenBetween(3, 24),
				),
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"auth_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"RBAC", "NONE",
				}, false),
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  DefaultVersion,
			},
			"admin_pass": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"eip_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"extend_params": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"service_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"service_registry_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"private": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"config_center_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"private": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceMicroserviceEngineCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.CseV2Client(region)
	if err != nil {
		return diag.Errorf("error creating CSE v2 client: %s", err)
	}

	networkId := d.Get("network_id").(string)
	subnetResp, err := vpc.GetVpcSubnetById(conf, region, networkId)
	if err != nil {
		return diag.FromErr(err)
	}
	vpcResp, err := vpc.GetVpcById(conf, region, subnetResp.VPC_ID)
	if err != nil {
		return diag.FromErr(err)
	}

	authType := d.Get("auth_type").(string)
	createOpts := engines.CreateOpts{
		Payment:             "1",
		SpecType:            d.Get("version").(string),
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Flavor:              d.Get("flavor").(string),
		AvailabilityZones:   utils.ExpandToStringList(d.Get("availability_zones").([]interface{})),
		AuthType:            authType,
		VpcName:             vpcResp.Name,
		VpcId:               vpcResp.ID,
		NetworkId:           networkId,
		SubnetCidr:          subnetResp.CIDR,
		PublicIpId:          d.Get("eip_id").(string),
		Inputs:              d.Get("extend_params").(map[string]interface{}),
		EnterpriseProjectId: common.GetEnterpriseProjectID(d, conf),
	}

	if authType == "RBAC" {
		createOpts.AuthCred = &engines.AuthCred{
			Password: d.Get("admin_pass").(string),
		}
	}

	resp, err := engines.Create(client, createOpts)
	if err != nil {
		return diag.Errorf("error creating Microservice engine: %s", err)
	}
	d.SetId(resp.ID)

	log.Printf("[DEBUG] Waiting for the Microservice engine to become running, the engine ID is %s.", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Init", "Executing"},
		Target:       []string{"Finished"},
		Refresh:      MicroserviceJobRefreshFunc(client, d.Id(), strconv.Itoa(resp.JobId)),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        180 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the creation of Microservice engine (%s) to complete: %s", d.Id(), err)
	}

	return resourceMicroserviceEngineRead(ctx, d, meta)
}

func flattenServiceRegistryAddresses(entrypoint engines.ExternalEntrypoint) (result []map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Recover panic when flattening service registry center structure: %#v", r)
		}
	}()

	entrypoints := map[string]interface{}{
		"private": entrypoint.ServiceEndpoint.ServiceCenter.MasterEntrypoint,
	}
	if !reflect.DeepEqual(entrypoint.PublicServiceEndpoint, engines.ServiceEndpoint{}) {
		entrypoints["public"] = entrypoint.PublicServiceEndpoint.ServiceCenter.MasterEntrypoint
	}

	return append(result, entrypoints)
}

func flattenConfigAddresses(entrypoint engines.ExternalEntrypoint) (result []map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Recover panic when flattening config center structure: %#v", r)
		}
	}()

	entrypoints := map[string]interface{}{
		"private": entrypoint.ServiceEndpoint.ConfigCenter.MasterEntrypoint,
	}
	if !reflect.DeepEqual(entrypoint.PublicServiceEndpoint, engines.ServiceEndpoint{}) {
		entrypoints["public"] = entrypoint.PublicServiceEndpoint.ConfigCenter.MasterEntrypoint
	}

	return append(result, entrypoints)
}

func resourceMicroserviceEngineRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.CseV2Client(region)
	if err != nil {
		return diag.Errorf("error creating CSE v2 client: %s", err)
	}

	resp, err := engines.Get(client, d.Id(), common.GetEnterpriseProjectID(d, conf))
	if err != nil {
		return common.CheckDeletedDiag(d, parseEngineJobError(err), "error retrieving Microservice engine")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", resp.Name),
		d.Set("flavor", resp.Flavor),
		d.Set("availability_zones", resp.Reference.AzList),
		d.Set("auth_type", resp.AuthType),
		d.Set("version", resp.SpecType),
		d.Set("enterprise_project_id", resp.EnterpriseProjectId),
		d.Set("network_id", resp.Reference.NetworkId),
		d.Set("description", resp.Description),
		d.Set("eip_id", resp.Reference.PublicIpId),
		d.Set("extend_params", resp.Reference.Inputs),
		d.Set("service_registry_addresses", flattenServiceRegistryAddresses(resp.ExternalEntrypoint)),
		d.Set("config_center_addresses", flattenConfigAddresses(resp.ExternalEntrypoint)),
	)

	diagErr := make([]diag.Diagnostic, 0, 3)
	// Attributes
	if resp.Reference.ServiceLimit != "" {
		limit, err := strconv.Atoi(resp.Reference.ServiceLimit)
		if err != nil {
			// Record and continue.
			diagErr = append(diagErr, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Wrong format",
				Detail:   fmt.Sprintf("Unable to parse the service limit (%#v).", resp.Reference.ServiceLimit),
			})
		} else {
			mErr = multierror.Append(mErr, d.Set("service_limit", limit))
		}
	}
	if resp.Reference.InstanceLimit != "" {
		limit, err := strconv.Atoi(resp.Reference.InstanceLimit)
		if err != nil {
			// Record and continue.
			diagErr = append(diagErr, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Wrong format",
				Detail:   fmt.Sprintf("Unable to parse the instance limit (%#v).", resp.Reference.InstanceLimit),
			})
		} else {
			mErr = multierror.Append(mErr, d.Set("instance_limit", limit))
		}
	}

	diagErr = append(diagErr, diag.FromErr(mErr.ErrorOrNil())...)
	return diagErr
}

func resourceMicroserviceEngineDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.CseV2Client(region)
	if err != nil {
		return diag.Errorf("error creating CSE v2 client: %s", err)
	}

	resp, err := engines.Delete(client, d.Id(), common.GetEnterpriseProjectID(d, conf))
	if err != nil {
		return diag.Errorf("error getting Microservice engine: %s", err)
	}

	log.Printf("[DEBUG] Waiting for the Microservice engine delete complete, the engine ID is %s.", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Init", "Executing"},
		Target:       []string{"Deleted"},
		Refresh:      MicroserviceJobRefreshFunc(client, d.Id(), strconv.Itoa(resp.JobId)),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        120 * time.Second,
		PollInterval: 15 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting the Microservice engine (%s): %s", d.Id(), err)
	}

	d.SetId("")

	return nil
}

func parseEngineJobError(respErr error) error {
	var apiErr engines.ErrorResponse
	if errCode, ok := respErr.(golangsdk.ErrDefault400); ok {
		pErr := json.Unmarshal(errCode.Body, &apiErr)
		if pErr == nil && (apiErr.ErrCode == "SVCSTG.00501116") {
			return golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Body: []byte("the Microservice engine has been deleted"),
				},
			}
		}
	}
	return respErr
}

func MicroserviceJobRefreshFunc(c *golangsdk.ServiceClient, engineId, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := engines.GetJob(c, engineId, jobId)
		if newErr := parseEngineJobError(err); newErr != nil {
			if _, ok := newErr.(golangsdk.ErrDefault404); ok {
				return resp, "Deleted", nil
			}
			return resp, "ERROR", newErr
		}
		return resp, resp.Status, nil
	}
}
