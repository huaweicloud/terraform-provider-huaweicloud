package ecs

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ECS POST /v2.1/{project_id}/servers/{server_id}/os-interface
// @API ECS GET /v2.1/{project_id}/servers/{server_id}/os-interface
// @API ECS PUT /v1/{project_id}/cloudservers/{server_id}/os-interface/{port_id}
// @API ECS DELETE /v2.1/{project_id}/servers/{server_id}/os-interface/{port_id}
// @API IAM POST /v3/auth/tokens
func ResourceComputePortAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputePortAttachCreate,
		ReadContext:   resourceComputePortAttachRead,
		DeleteContext: resourceComputePortAttachDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceComputeInstancePortAttachImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"port_id", "net_id"},
			},
			"net_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"port_id", "net_id"},
			},
			"api_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "v2.1",
			},
			"assume_role": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "This is a beta feature and requires application.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agency_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"domain_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"port_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

type ComputePortAttachWrapper struct {
	*schemas.ResourceDataWrapper
	Config     *config.Config
	ApiVersion string
}

func newComputePortAttachWrapper(d *schema.ResourceData, meta interface{}) *ComputePortAttachWrapper {
	return &ComputePortAttachWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
		ApiVersion:          d.Get("api_version").(string),
	}
}

func resourceComputePortAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newComputePortAttachWrapper(d, meta)
	rst, err := wrapper.AttachPort()
	instanceId := d.Get("instance_id").(string)
	if err != nil {
		log.Printf("[ERROR] failed to attach port to instance[%s], error: %s", instanceId, err)
		return diag.Errorf("failed to attach port: %s", err)
	}

	portID := rst.Get("interfaceAttachment").Get("port_id").String()
	d.SetId(fmt.Sprintf("%s/%s", instanceId, portID))
	return resourceComputePortAttachRead(ctx, d, meta)
}

func resourceComputePortAttachRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newComputePortAttachWrapper(d, meta)

	instanceId := d.Get("instance_id").(string)
	portId := d.Get("port_id").(string)
	rst, err := wrapper.GetPort(instanceId, portId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error get ECS port[%s/%s] details", instanceId, portId))
	}

	var aErr error
	if _, ok := d.GetOk("assume_role"); ok {
		assumeRole := make([]map[string]any, 0)
		assumeRole = append(assumeRole, map[string]any{
			"agency_name": d.Get("assume_role.0.agency_name").(string),
			"domain_name": d.Get("assume_role.0.domain_name").(string),
			"domain_id":   d.Get("assume_role.0.domain_id").(string),
		})
		aErr = d.Set("assume_role", assumeRole)
	}

	interfaceAttachment := rst.Get("interfaceAttachment")
	mErr := multierror.Append(aErr,
		d.Set("region", wrapper.Config.GetRegion(d)),
		d.Set("port_id", interfaceAttachment.Get("port_id").String()),
		d.Set("net_id", interfaceAttachment.Get("net_id").String()),
		d.Set("port_state", interfaceAttachment.Get("port_state").String()),
		d.Set("api_version", d.Get("api_version")),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceComputePortAttachDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newComputePortAttachWrapper(d, meta)
	_, err := wrapper.DetachPort(ctx)
	if err != nil {
		return common.CheckDeletedDiag(d,
			err,
			fmt.Sprintf("error deleting ECS port[%s/%s]", d.Get("instance_id").(string), d.Get("port_id").(string)),
		)
	}
	d.SetId("")
	return nil
}

func resourceComputeInstancePortAttachImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<port_id>', but got '%s'", importedId)
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("port_id", parts[1]),
		d.Set("api_version", "v2.1"),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

func (w *ComputePortAttachWrapper) AttachPort() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "ecs")
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("/%s/{project_id}/servers/%s/os-interface", w.ApiVersion, w.Get("instance_id"))

	interfaceAttachment := map[string]any{
		"port_id": w.Get("port_id"),
		"net_id":  w.Get("net_id"),
	}
	if _, ok := w.ResourceData.GetOk("assume_role"); ok {
		verifyToken, err := w.getAssumeToken()
		if err != nil {
			return nil, fmt.Errorf("failed to get assume token: %s", err)
		}
		interfaceAttachment["verify_token_for_port"] = verifyToken
	}
	params := map[string]any{
		"interfaceAttachment": interfaceAttachment,
	}
	params = utils.RemoveNil(params)
	hc := httphelper.New(client)
	rst, err := hc.Method("POST").
		URI(uri).
		Body(params).
		Request().
		Result()

	if hc.Response.StatusCode > 400 {
		err = fmt.Errorf("statusCode: %v", hc.Response.StatusCode)
		if rst != nil {
			err = fmt.Errorf("%s", rst.String())
		}
		return nil, err
	}
	return rst, err
}

func (w *ComputePortAttachWrapper) GetPort(instanceId, portId string) (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "ecs")
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("/%s/{project_id}/servers/%s/os-interface/%s", w.ApiVersion, instanceId, portId)
	hc := httphelper.New(client)
	rst, err := hc.Method("GET").
		URI(uri).
		Request().
		Result()

	if hc.Response.StatusCode == 404 {
		return rst, golangsdk.ErrDefault404{}
	}

	if hc.Response.StatusCode > 400 {
		err = fmt.Errorf("failed to query port details, statusCode: %v", hc.Response.StatusCode)
		if rst != nil {
			err = fmt.Errorf("failed to query port details: %s", rst.String())
		}
		return nil, err
	}
	if err == nil {
		return rst, nil
	}
	return nil, err
}

func (w *ComputePortAttachWrapper) DetachPort(ctx context.Context) (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "ecs")
	if err != nil {
		return nil, err
	}
	instanceId := w.Get("instance_id").(string)
	portId := w.Get("port_id").(string)
	uri := fmt.Sprintf("/%s/{project_id}/servers/%s/os-interface/%s", w.ApiVersion, instanceId, portId)
	hc := httphelper.New(client)

	rst, err := hc.Method("DELETE").URI(uri).Send()
	if err != nil {
		return nil, err
	}

	if hc.Response.StatusCode == 404 {
		w.SetId("")
		return nil, golangsdk.ErrDefault404{}
	}
	if hc.Response.StatusCode > 400 {
		err = fmt.Errorf("failed to detach port, statusCode: %v", hc.Response.StatusCode)
		if rst != nil {
			err = fmt.Errorf("failed to detach port: %s", rst.String())
		}
		return nil, err
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"ACTIVE"},
		Target:       []string{"COMPLETED"},
		Refresh:      w.waitToBeDeleted(instanceId, portId),
		Timeout:      w.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}
	w.SetId("")
	return nil, nil
}

func (w *ComputePortAttachWrapper) waitToBeDeleted(instanceId, portId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := w.GetPort(instanceId, portId)
		log.Printf("[DEBUG] Attempting to delete ECS Port %#v", err)
		if err == nil {
			return nil, "ACTIVE", nil
		}

		var err404 golangsdk.ErrDefault404
		if errors.As(err, &err404) {
			log.Printf("[DEBUG] The ECS Port (%s) has been deleted.", portId)
			return r, "COMPLETED", nil
		}

		if err != nil {
			return nil, "ERROR", err
		}
		return r, "ACTIVE", nil
	}
}
func (w *ComputePortAttachWrapper) getAssumeToken() (*string, error) {
	client, err := w.NewClient(w.Config, "iam")
	if err != nil {
		return nil, err
	}

	uri := "/v3/auth/tokens"
	params := map[string]any{
		"auth": map[string]any{
			"identity": map[string]any{
				"methods": []string{"assume_role"},
				"assume_role": map[string]any{
					"agency_name": w.Get("assume_role.0.agency_name"),
					"domain_name": w.Get("assume_role.0.domain_name"),
					"domain_id":   w.Get("assume_role.0.domain_id"),
				},
			},
			"scope": map[string]any{
				"project": map[string]any{
					"name": w.Config.GetRegion(w.ResourceData),
				},
			},
		},
	}
	params = utils.RemoveNil(params)
	hp := httphelper.New(client).
		Method("POST").
		URI(uri).
		Body(params).
		Request()
	if hp.ErrorOrNil() != nil {
		return nil, hp.ErrorOrNil()
	}

	token := hp.Response.Header.Get("X-Subject-Token")
	return &token, nil
}
