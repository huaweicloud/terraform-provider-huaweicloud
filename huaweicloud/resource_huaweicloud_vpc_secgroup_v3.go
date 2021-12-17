package huaweicloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	vpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3/model"
	region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3/region"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"time"
)

func ResourceVpcSecGroupV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcSecGroupV3Create,
		ReadContext:   resourceVpcSecGroupV3Read,
		UpdateContext: resourceVpcSecGroupV3Update,
		DeleteContext: resourceVpcSecGroupV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
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
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceVpcSecGroupV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	ak := config.AccessKey
	sk := config.SecretKey
	auth := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		Build()

	client := vpc.NewVpcClient(
		vpc.VpcClientBuilder().
			WithRegion(region.ValueOf(d.Get("region").(string))).
			WithCredential(auth).
			Build())

	request := &model.CreateSecurityGroupRequest{}
	descriptionSecurityGroupCreateSecurityGroupOption := d.Get("description").(string)
	enterpriseProjectIdSecurityGroupCreateSecurityGroupOption := d.Get("enterprise_project_id").(string)
	securityGroupbody := &model.CreateSecurityGroupOption{
		Name:                d.Get("name").(string),
		Description:         &descriptionSecurityGroupCreateSecurityGroupOption,
		EnterpriseProjectId: &enterpriseProjectIdSecurityGroupCreateSecurityGroupOption,
	}
	dryRunCreateSecurityGroupRequestBody := d.Get("dry_run").(bool)
	request.Body = &model.CreateSecurityGroupRequestBody{
		SecurityGroup: securityGroupbody,
		DryRun:        &dryRunCreateSecurityGroupRequestBody,
	}
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err := client.CreateSecurityGroup(request)
		if err != nil {
			if response.HttpStatusCode >= 500 {
				time.Sleep(3 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	response, err := client.CreateSecurityGroup(request)
	if err != nil {
		return fmtp.DiagErrorf("Error creationg security group %s: %s", d.Id(), err)
	}
	d.SetId(response.SecurityGroup.Id)
	return resourceVpcSecGroupV3Read(nil, d, meta)
}

func resourceVpcSecGroupV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	ak := config.AccessKey
	sk := config.SecretKey

	auth := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		Build()

	client := vpc.NewVpcClient(
		vpc.VpcClientBuilder().
			WithRegion(region.ValueOf(d.Get("region").(string))).
			WithCredential(auth).
			Build())

	request := &model.ShowSecurityGroupRequest{}
	request.SecurityGroupId = d.Id()
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err := client.ShowSecurityGroup(request)
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
		return fmtp.DiagErrorf("Error reading security group %s", err)
	}
	return nil
}

func resourceVpcSecGroupV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	ak := config.AccessKey
	sk := config.SecretKey

	auth := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		Build()

	client := vpc.NewVpcClient(
		vpc.VpcClientBuilder().
			WithRegion(region.ValueOf(d.Get("region").(string))).
			WithCredential(auth).
			Build())

	request := &model.UpdateSecurityGroupRequest{}
	request.SecurityGroupId = d.Id()
	descriptionSecurityGroupCreateSecurityGroupOption := d.Get("description").(string)
	nameSecurityGroupUpdateSecurityGroupOption := d.Get("name").(string)
	securityGroupbody := &model.UpdateSecurityGroupOption{
		Name:        &nameSecurityGroupUpdateSecurityGroupOption,
		Description: &descriptionSecurityGroupCreateSecurityGroupOption,
	}
	dryRunUpdateSecurityGroupRequestBody := d.Get("dry_run").(bool)
	request.Body = &model.UpdateSecurityGroupRequestBody{
		SecurityGroup: securityGroupbody,
		DryRun:        &dryRunUpdateSecurityGroupRequestBody,
	}
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err := client.UpdateSecurityGroup(request)
		if err != nil {
			if response.HttpStatusCode >= 500 {
				time.Sleep(3 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	response, err := client.UpdateSecurityGroup(request)
	if err != nil {
		return fmtp.DiagErrorf("Error updating security group %s: %s", err)
	}
	d.SetId(response.SecurityGroup.Id)
	return resourceVpcSecGroupV3Read(nil, d, meta)
}
func resourceVpcSecGroupV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	ak := config.AccessKey
	sk := config.SecretKey

	auth := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		Build()

	client := vpc.NewVpcClient(
		vpc.VpcClientBuilder().
			WithRegion(region.ValueOf(d.Get("region").(string))).
			WithCredential(auth).
			Build())

	request := &model.DeleteSecurityGroupRequest{}
	request.SecurityGroupId = d.Id()
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err := client.DeleteSecurityGroup(request)
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
		return fmtp.DiagErrorf("Error deleting security group %s", err)
	}
	return nil
}
