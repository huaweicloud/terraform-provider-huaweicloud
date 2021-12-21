package huaweicloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	vpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3/model"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3/region"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"time"
)

func ResourceVpcAddressGroupV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcAddressGroupCreate,
		UpdateContext: resourceVpcAddressGroupUpdate,
		DeleteContext: resourceVpcAddressGroupDelete,
		ReadContext:   resourceVpcAddressGroupRead,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_set": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceVpcAddressGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	auth := basic.NewCredentialsBuilder().
		WithAk(config.AccessKey).
		WithSk(config.SecretKey).
		Build()

	client := vpc.NewVpcClient(
		vpc.VpcClientBuilder().
			WithRegion(region.ValueOf(config.GetRegion(d))).
			WithCredential(auth).
			Build())

	request := &model.CreateAddressGroupRequest{}
	var finallyIpSet []string
	if listIpSetAddressGroup, ok := d.Get("ip_set").([]interface{}); ok {
		for _, value := range listIpSetAddressGroup {
			finallyValue := value.(string)
			finallyIpSet = append(finallyIpSet, finallyValue)
		}
	}
	descriptionAddressGroupCreateAddressGroupOption := d.Get("description").(string)
	addressGroupBody := &model.CreateAddressGroupOption{
		Name:        d.Get("name").(string),
		Description: &descriptionAddressGroupCreateAddressGroupOption,
		IpVersion:   int32(d.Get("ip_version").(int)),
		IpSet:       &finallyIpSet,
	}
	dryRunCreateAddressGroupRequestBody := d.Get("dry_run").(bool)
	request.Body = &model.CreateAddressGroupRequestBody{
		AddressGroup: addressGroupBody,
		DryRun:       &dryRunCreateAddressGroupRequestBody,
	}
	var response *model.CreateAddressGroupResponse
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error
		response, err = client.CreateAddressGroup(request)
		if err != nil {
			if response.HttpStatusCode >= 500 {
				time.Sleep(3 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmtp.DiagErrorf("Error creating VPC addressGroup: %s", err)
	}
	d.SetId(response.AddressGroup.Id)
	return resourceVpcAddressGroupRead(ctx, d, meta)
}

func resourceVpcAddressGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	auth := basic.NewCredentialsBuilder().
		WithAk(config.AccessKey).
		WithSk(config.SecretKey).
		Build()

	client := vpc.NewVpcClient(
		vpc.VpcClientBuilder().
			WithRegion(region.ValueOf(config.GetRegion(d))).
			WithCredential(auth).
			Build())

	request := &model.ShowAddressGroupRequest{}
	request.AddressGroupId = d.Id()
	var response *model.ShowAddressGroupResponse
	err := resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		var err error
		response, err = client.ShowAddressGroup(request)
		if err != nil {
			if response.HttpStatusCode >= 500 {
				time.Sleep(3 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmtp.DiagErrorf("Error query VPC addressGroup: %s", err)
	}
	return nil
}

func resourceVpcAddressGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	auth := basic.NewCredentialsBuilder().
		WithAk(config.AccessKey).
		WithSk(config.SecretKey).
		Build()

	client := vpc.NewVpcClient(
		vpc.VpcClientBuilder().
			WithRegion(region.ValueOf(config.GetRegion(d))).
			WithCredential(auth).
			Build())

	request := &model.DeleteAddressGroupRequest{}
	request.AddressGroupId = d.Id()
	var response *model.DeleteAddressGroupResponse
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		var err error
		response, err = client.DeleteAddressGroup(request)
		if err != nil {
			if response.HttpStatusCode >= 500 {
				time.Sleep(3 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmtp.DiagErrorf("Error delete VPC addressGroup: %s", err)
	}
	return nil
}

func resourceVpcAddressGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	auth := basic.NewCredentialsBuilder().
		WithAk(config.AccessKey).
		WithSk(config.SecretKey).
		Build()

	client := vpc.NewVpcClient(
		vpc.VpcClientBuilder().
			WithRegion(region.ValueOf(config.GetRegion(d))).
			WithCredential(auth).
			Build())

	request := &model.UpdateAddressGroupRequest{}
	request.AddressGroupId = d.Id()
	var finallyIpSet []string
	if listIpSetAddressGroup, ok := d.Get("ip_set").([]interface{}); ok {
		for _, value := range listIpSetAddressGroup {
			finallyValue := value.(string)
			finallyIpSet = append(finallyIpSet, finallyValue)
		}
	}
	nameAddressGroupUpdateAddressGroupOption := d.Get("name").(string)
	descriptionAddressGroupUpdateAddressGroupOption := d.Get("description").(string)
	addressGroupBody := &model.UpdateAddressGroupOption{
		Name:        &nameAddressGroupUpdateAddressGroupOption,
		Description: &descriptionAddressGroupUpdateAddressGroupOption,
		IpSet:       &finallyIpSet,
	}
	dryRunUpdateAddressGroupRequestBody := d.Get("dry_run").(bool)
	request.Body = &model.UpdateAddressGroupRequestBody{
		AddressGroup: addressGroupBody,
		DryRun:       &dryRunUpdateAddressGroupRequestBody,
	}
	var response *model.UpdateAddressGroupResponse
	err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		var err error
		response, err = client.UpdateAddressGroup(request)
		if err != nil {
			if response.HttpStatusCode >= 500 {
				time.Sleep(3 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmtp.DiagErrorf("Error update VPC addressGroup: %s", err)
	}
	d.SetId(response.AddressGroup.Id)
	return resourceVpcAddressGroupRead(ctx, d, meta)
}
