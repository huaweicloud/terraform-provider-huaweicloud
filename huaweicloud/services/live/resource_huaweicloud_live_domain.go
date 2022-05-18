package live

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	v1 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/live/v1"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/live/v1/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	domainTypePull = "pull"
	domainTypePush = "push"

	domainStatusOn  = "on"
	domainStatusOff = "off"
)

func ResourceDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCreate,
		UpdateContext: resourceDomainUpdate,
		DeleteContext: resourceDomainDelete,
		ReadContext:   resourceDomainRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
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
			},

			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{domainTypePull, domainTypePush}, false),
			},

			"ingest_domain_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{domainStatusOn, domainStatusOff}, false),
				Computed:     true,
			},

			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcLiveV1Client(region)
	if err != nil {
		return diag.Errorf("error creating Live v1 client: %s", err)
	}

	createOpts, err := buildCreateParams(d, region)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Create Live domain params : %#v", createOpts)

	_, err = client.CreateDomain(createOpts)
	if err != nil {
		return diag.Errorf("error creating Live domain: %s", err)
	}

	d.SetId(createOpts.Body.Domain)

	err = waitingForDomainStatus(ctx, client, d.Id(), model.GetDecoupledLiveDomainInfoStatusEnum().ON,
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	// associate the streaming domain name with an ingest domain
	err = associatingDomain(d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	// off the domain
	if d.Get("status").(string) == domainStatusOff {
		err = updateStatus(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDomainRead(ctx, d, meta)
}

func resourceDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcLiveV1Client(region)
	if err != nil {
		return diag.Errorf("error creating Live v1 client: %s", err)
	}

	domain := d.Id()
	response, err := client.ShowDomain(&model.ShowDomainRequest{Domain: &domain})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Live domain")
	}

	if response.DomainInfo == nil || len(*response.DomainInfo) != 1 {
		return diag.Errorf("error retrieving Live domain")
	}
	r := *response.DomainInfo
	detail := r[0]

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", detail.Domain),
		d.Set("type", utils.MarshalValue(detail.DomainType)),
		d.Set("status", utils.MarshalValue(detail.Status)),
		d.Set("ingest_domain_name", detail.RelatedDomain),
		d.Set("cname", detail.DomainCname),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcLiveV1Client(region)
	if err != nil {
		return diag.Errorf("error creating Live v1 client: %s", err)
	}

	domainName := d.Get("name").(string)

	// associate the streaming domain name with an ingest domain Or delete association
	if d.HasChange("ingest_domain_name") {
		ingetstDomainNameOld, ingetstDomainName := d.GetChange("ingest_domain_name")

		if ingetstDomainName == "" {
			err = deleteAssociation(client, domainName, ingetstDomainNameOld.(string))
		} else {
			err = associatingDomain(d, client)
		}
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// update the domain status
	if d.HasChange("status") {
		err = updateStatus(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDomainRead(ctx, d, meta)
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcLiveV1Client(region)
	if err != nil {
		return diag.Errorf("error creating Live v1 client: %s", err)
	}

	// 1. off the domain
	reqStatus := model.GetLiveDomainModifyReqStatusEnum().OFF
	_, err = client.UpdateDomain(&model.UpdateDomainRequest{Body: &model.LiveDomainModifyReq{
		Domain: d.Get("name").(string),
		Status: &reqStatus,
	}})
	if err != nil {
		return diag.Errorf("error disabling Live domain: %s", err)
	}

	err = waitingForDomainStatus(ctx, client, d.Id(), model.GetDecoupledLiveDomainInfoStatusEnum().OFF,
		d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	// 2. delete the domain
	deleteOpts := &model.DeleteDomainRequest{
		Domain: d.Get("name").(string),
	}
	_, err = client.DeleteDomain(deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting Live domain: %s", err)
	}

	return nil
}

func associatingDomain(d *schema.ResourceData, client *v1.LiveClient) error {
	if v, ok := d.GetOk("ingest_domain_name"); ok {
		domainType := d.Get("type").(string)
		if domainType == domainTypePush {
			return fmt.Errorf("the ingest domain cannot associate with an ingest domain")
		}

		_, err := client.CreateDomainMapping(&model.CreateDomainMappingRequest{
			Body: &model.DomainMapping{
				PullDomain: d.Get("name").(string),
				PushDomain: v.(string),
			},
		})
		if err != nil {
			return fmt.Errorf("error associating the streaming domain name with an ingest domain: %s", err)
		}
	}

	return nil

}

func deleteAssociation(client *v1.LiveClient, pullDomain, pushDomain string) error {
	_, err := client.DeleteDomainMapping(&model.DeleteDomainMappingRequest{
		PullDomain: pullDomain,
		PushDomain: pushDomain,
	})
	if err != nil {
		return fmt.Errorf("error deleting the association between the streaming domain and ingest domain: %s", err)
	}

	return nil

}

func buildCreateParams(d *schema.ResourceData, region string) (*model.CreateDomainRequest, error) {
	var domainType model.LiveDomainCreateReqDomainType
	if d.Get("type").(string) == domainTypePull {
		domainType = model.GetLiveDomainCreateReqDomainTypeEnum().PULL
	} else {
		domainType = model.GetLiveDomainCreateReqDomainTypeEnum().PUSH
	}

	req := model.CreateDomainRequest{
		Body: &model.LiveDomainCreateReq{
			Domain:     d.Get("name").(string),
			DomainType: domainType,
			Region:     region,
		},
	}
	return &req, nil

}

func waitingForDomainStatus(ctx context.Context, client *v1.LiveClient, name string,
	status model.DecoupledLiveDomainInfoStatus, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Done"},
		Refresh: func() (interface{}, string, error) {
			resp, err := client.ShowDomain(&model.ShowDomainRequest{Domain: &name})
			if err != nil {
				return nil, "failed", err
			}

			if resp.DomainInfo == nil || len(*resp.DomainInfo) != 1 {
				return nil, "failed", fmt.Errorf("error retrieving Live domain")
			}
			r := *resp.DomainInfo
			detail := r[0]

			if detail.Status == status {
				return resp, "Done", nil
			}
			return resp, "Pending", nil
		},
		Timeout:      timeout,
		PollInterval: 5 * time.Second,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for Live domain (%s) status to become %v: %s", name, status, err)
	}
	return nil
}

func updateStatus(ctx context.Context, d *schema.ResourceData, client *v1.LiveClient) error {
	var reqStatus model.LiveDomainModifyReqStatus
	var respStatus model.DecoupledLiveDomainInfoStatus

	if d.Get("status").(string) == domainStatusOff {
		reqStatus = model.GetLiveDomainModifyReqStatusEnum().OFF
		respStatus = model.GetDecoupledLiveDomainInfoStatusEnum().OFF
	} else {
		reqStatus = model.GetLiveDomainModifyReqStatusEnum().ON
		respStatus = model.GetDecoupledLiveDomainInfoStatusEnum().ON
	}

	_, err := client.UpdateDomain(&model.UpdateDomainRequest{
		Body: &model.LiveDomainModifyReq{
			Domain: d.Get("name").(string),
			Status: &reqStatus,
		},
	})

	if err != nil {
		return fmt.Errorf("error updating Live domain: %s", err)
	}

	return waitingForDomainStatus(ctx, client, d.Id(), respStatus, d.Timeout(schema.TimeoutUpdate))
}
