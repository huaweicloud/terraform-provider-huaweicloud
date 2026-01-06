package aad

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/aad/v1/rules"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// ResourceForwardRule is the imple of huaweicloud_aad_forward_rule
// @API AAD POST /v1/aad/instances/{instance_id}/{ip}/rules/batch-create
// @API AAD POST /v1/aad/instances/{instance_id}/{ip}/rules/batch-delete
// @API AAD PUT /v1/aad/instances/{instance_id}/{ip}/rules/{rule_id}
// @API AAD GET /v1/aad/instances/{instance_id}/{ip}/rules
func ResourceForwardRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceForwardRuleCreate,
		ReadContext:   resourceForwardRuleRead,
		UpdateContext: resourceForwardRuleUpdate,
		DeleteContext: resourceForwardRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceForwardRuleImportState,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"forward_protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"forward_port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"source_port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"source_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"lb_method": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "schema: Deprecated",
			},
		},
	}
}

func resourceForwardRuleCreate(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.AadV1Client("")
	if err != nil {
		return diag.Errorf("error creating Advanced Anti-DDoS V1 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	advancedIp := d.Get("ip").(string)
	protocol := d.Get("forward_protocol").(string)
	forwardPort := d.Get("forward_port").(int)
	opts := rules.BatchCreateOpts{
		Rules: []rules.RuleOpts{
			{
				ForwardProtocol: protocol,
				ForwardPort:     forwardPort,
				SourcePort:      d.Get("source_port").(int),
				SourceIp:        d.Get("source_ip").(string),
			},
		},
	}
	_, err = rules.BatchCreate(client, instanceId, advancedIp, opts)
	if err != nil {
		return diag.Errorf("error creating Advanced Anti-DDoS forward rule: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/%d", instanceId, advancedIp, protocol, forwardPort))
	return resourceForwardRuleRead(ctx, d, meta)
}

func filterForwardRuleFromList(ruleList []rules.Rule, protocol string, port int) *rules.Rule {
	for _, rule := range ruleList {
		if rule.ForwardPort == port && rule.ForwardProtocol == protocol {
			return &rule
		}
	}
	return nil
}

func GetForwardRuleFromServer(client *golangsdk.ServiceClient, instanceId, advancedIp, protocol string,
	port int) (*rules.Rule, error) {
	resp, err := rules.List(client, instanceId, advancedIp)
	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/aad/instances/{instance_id}/{ip}/rules",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the Advanced Anti-DDoS forward rule (%s/%s) does not exist", instanceId, advancedIp)),
			},
		}
	}
	if result := filterForwardRuleFromList(resp, protocol, port); result != nil {
		return result, nil
	}
	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Method:    "GET",
			URL:       "/v1/aad/instances/{instance_id}/{ip}/rules",
			RequestId: "NONE",
			Body:      []byte(fmt.Sprintf("no Advanced Anti-DDoS forward rule matched the given protocol (%s) and (or) port (%d)", protocol, port)),
		},
	}
}

func resourceForwardRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.AadV1Client("")
	if err != nil {
		return diag.Errorf("error creating Advanced Anti-DDoS V1 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	advancedIp := d.Get("ip").(string)
	protocol := d.Get("forward_protocol").(string)
	forwardPort := d.Get("forward_port").(int)
	resp, err := GetForwardRuleFromServer(client, instanceId, advancedIp, protocol, forwardPort)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Advanced Anti-DDoS forward rule")
	}
	log.Printf("[DEBUG] Retrieved Advanced Anti-DDos forward rule: %#v", resp)

	mErr := multierror.Append(nil,
		d.Set("forward_protocol", resp.ForwardProtocol),
		d.Set("forward_port", resp.ForwardPort),
		d.Set("source_port", resp.SourcePort),
		d.Set("source_ip", resp.SourceIp),
		d.Set("status", resp.Status),
		d.Set("lb_method", resp.LbMethod),
		d.Set("rule_id", resp.ID),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting resource: %s", mErr)
	}
	return nil
}

func resourceForwardRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.AadV1Client("")
	if err != nil {
		return diag.Errorf("error creating Advanced Anti-DDoS V1 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	advancedIp := d.Get("ip").(string)
	ruleId := d.Get("rule_id").(string)

	opts := rules.UpdateOpts{
		ForwardProtocol: d.Get("forward_protocol").(string),
		ForwardPort:     d.Get("forward_port").(int),
		SourcePort:      d.Get("source_port").(int),
		SourceIp:        d.Get("source_ip").(string),
	}
	err = rules.Update(client, instanceId, advancedIp, ruleId, opts)
	if err != nil {
		return diag.Errorf("error updating Advanced Anti-DDoS forward rule (%s): %v", ruleId, err)
	}
	return resourceForwardRuleRead(ctx, d, meta)
}

func resourceForwardRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.AadV1Client("")
	if err != nil {
		return diag.Errorf("error creating Advanced Anti-DDoS V1 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	advancedIp := d.Get("ip").(string)
	ruleId := d.Get("rule_id").(string)
	opts := rules.BatchDeleteOpts{
		RuleIds: []string{
			ruleId,
		},
	}
	_, err = rules.BatchDelete(client, instanceId, advancedIp, opts)
	if err != nil {
		return diag.Errorf("error deleting Advanced Anti-DDoS forward rule (%s): %v", ruleId, err)
	}
	return nil
}

func resourceForwardRuleImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.SplitN(d.Id(), "/", 4)
	if len(parts) != 4 {
		return nil, fmt.Errorf("Invalid format specified for import id, must be " +
			"<instance_id>/<ip>/<forward_protocol>/<forward_port>")
	}

	portNum, err := strconv.Atoi(parts[3])
	if err != nil {
		return nil, err
	}
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("ip", parts[1]),
		d.Set("forward_protocol", parts[2]),
		d.Set("forward_port", portNum),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
