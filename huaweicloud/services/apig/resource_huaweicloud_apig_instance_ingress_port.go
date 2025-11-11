package apig

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var instanceIngressPortNonUpdatableParams = []string{"instance_id", "protocol", "port"}

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/custom-ingress-ports
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/custom-ingress-ports
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/custom-ingress-ports/{ingress_port_id}
func ResourceInstanceIngressPort() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceIngressPortCreate,
		ReadContext:   resourceInstanceIngressPortRead,
		UpdateContext: resourceInstanceIngressPortUpdate,
		DeleteContext: resourceInstanceIngressPortDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceInstanceIngressPortImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(instanceIngressPortNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the ingress port is located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the custom ingress port belongs.`,
			},
			"protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The protocol of the custom ingress port.`,
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The port number of the custom ingress port.`,
			},

			// Attributes.
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the custom ingress port.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildInstanceIngressPortBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"protocol":     d.Get("protocol"),
		"ingress_port": d.Get("port"),
	}
}

func resourceInstanceIngressPortCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	// Lock the resource to prevent concurrent creations (error APIC.9224 will be returned if multiple requests are sent
	// concurrently)
	config.MutexKV.Lock(instanceId)
	defer config.MutexKV.Unlock(instanceId)

	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/custom-ingress-ports"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildInstanceIngressPortBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating APIG instance ingress port: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	ingressPortId := utils.PathSearch("ingress_port_id", respBody, "").(string)
	if ingressPortId == "" {
		return diag.Errorf("unable to find the ingress port ID from the API response")
	}
	d.SetId(ingressPortId)

	return resourceInstanceIngressPortRead(ctx, d, meta)
}

func buildInstanceIngressPortsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("protocol"); ok {
		res = fmt.Sprintf("%s&protocol=%v", res, v)
	}
	if v, ok := d.GetOk("port"); ok {
		res = fmt.Sprintf("%s&ingress_port=%v", res, v)
	}

	return res
}

func listInstanceIngressPorts(client *golangsdk.ServiceClient, instanceId string, d ...*schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/apigw/instances/{instance_id}/custom-ingress-ports?limit={limit}"
		result  = make([]interface{}, 0)
		limit   = 500
		offset  = 0
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	if len(d) > 0 {
		listPath += buildInstanceIngressPortsQueryParams(d[0])
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		ingressPortInfos := utils.PathSearch("ingress_port_infos", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, ingressPortInfos...)
		if len(ingressPortInfos) < limit {
			break
		}
		offset += len(ingressPortInfos)
	}

	return result, nil
}

func GetInstanceIngressPortById(client *golangsdk.ServiceClient, instanceId, ingressPortId string, d ...*schema.ResourceData) (interface{}, error) {
	ingressPorts, err := listInstanceIngressPorts(client, instanceId, d...)
	if err != nil {
		return nil, err
	}

	for _, ingressPort := range ingressPorts {
		portId := utils.PathSearch("ingress_port_id", ingressPort, "").(string)
		if portId == ingressPortId {
			return ingressPort, nil
		}
	}

	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Method:    "GET",
			URL:       "/v2/{project_id}/apigw/instances/{instance_id}/custom-ingress-ports",
			RequestId: "NONE",
			Body:      []byte(fmt.Sprintf("ingress port (%s) is not found", ingressPortId)),
		},
	}
}

func resourceInstanceIngressPortRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	ingressPort, err := GetInstanceIngressPortById(client, instanceId, d.Id(), d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving APIG instance ingress port (%s)", d.Id()))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("protocol", utils.PathSearch("protocol", ingressPort, nil)),
		d.Set("port", utils.PathSearch("ingress_port", ingressPort, nil)),
		d.Set("status", utils.PathSearch("status", ingressPort, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceInstanceIngressPortUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInstanceIngressPortDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	var (
		httpUrl       = "v2/{project_id}/apigw/instances/{instance_id}/custom-ingress-ports/{ingress_port_id}"
		instanceId    = d.Get("instance_id").(string)
		ingressPortId = d.Id()
	)

	// Lock the resource to prevent concurrent deletions (error APIC.9224 will be returned if multiple requests are sent
	// concurrently)
	config.MutexKV.Lock(instanceId)
	defer config.MutexKV.Unlock(instanceId)

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)
	deletePath = strings.ReplaceAll(deletePath, "{ingress_port_id}", ingressPortId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "APIC.7200"),
			fmt.Sprintf("error deleting APIG instance ingress port (%s)", ingressPortId))
	}

	return nil
}

func resourceInstanceIngressPortImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")

	switch len(parts) {
	case 2:
		if !utils.IsUUID(parts[1]) {
			return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but got '%s'", importedId)
		}
		d.SetId(parts[1])
		return []*schema.ResourceData{d}, d.Set("instance_id", parts[0])
	case 3:
		if utils.IsUUID(parts[1]) {
			return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<protocol>/<port>', but got '%s'", importedId)
		}
		portNum, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid port number, want '<instance_id>/<protocol>/<port>', but got '%s'", importedId)
		}

		// Store the instance ID as the temporary resource ID first, then update it once the ingress port ID is obtained.
		d.SetId(parts[0])
		mErr := multierror.Append(nil,
			d.Set("instance_id", parts[0]),
			d.Set("protocol", parts[1]),
			d.Set("port", portNum),
		)
		if err := mErr.ErrorOrNil(); err != nil {
			return nil, fmt.Errorf("failed to set values in import state, %s", err)
		}

		cfg := meta.(*config.Config)
		region := cfg.GetRegion(d)
		client, err := cfg.NewServiceClient("apig", region)
		if err != nil {
			return nil, fmt.Errorf("error creating APIG client: %s", err)
		}
		respBody, err := listInstanceIngressPorts(client, parts[0], d)
		if err != nil {
			return nil, err
		}
		d.SetId(utils.PathSearch("[0].ingress_port_id", respBody, "").(string))
		return []*schema.ResourceData{d}, nil
	}

	return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>' or "+
		"'<instance_id>/<protocol>/<port>', but got '%s'", importedId)
}
