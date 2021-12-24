package huaweicloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	swr "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/swr/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/swr/v2/model"
	region "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/swr/v2/region"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func ResourceSwrV2Namespaces() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceSwrV2NamespacesCreate,
		ReadContext:   ResourceSwrV2NamespacesRead,
		DeleteContext: ResourceSwrV2NamespacesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Second),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func ResourceSwrV2NamespacesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	ak := config.AccessKey
	sk := config.SecretKey
	auth := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		Build()

	client := swr.NewSwrClient(
		swr.SwrClientBuilder().
			WithRegion(region.ValueOf("cn-north-4")).
			WithCredential(auth).
			Build())

	request := &model.CreateNamespaceRequest{}
	request.Body = &model.CreateNamespaceRequestBody{
		Namespace: d.Get("namespace").(string),
	}

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err := client.CreateNamespace(request)
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
		return fmtp.DiagErrorf("[Error] Error create Huaweicloud Swr_Namespace:%+v\n", err)
	} else {
		d.SetId(d.Get("namespace").(string))
	}
	return nil
}

func ResourceSwrV2NamespacesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	ak := config.AccessKey
	sk := config.SecretKey
	auth := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		Build()

	client := swr.NewSwrClient(
		swr.SwrClientBuilder().
			WithRegion(region.ValueOf("cn-north-4")).
			WithCredential(auth).
			Build())

	request := &model.ShowNamespaceRequest{}
	request.Namespace = d.Get("namespace").(string)

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err := client.ShowNamespace(request)
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
		return fmtp.DiagErrorf("[Error]Error Read Huaweicloud Swr_Namespace:%+v\n", err)
	}
	return nil
}

func ResourceSwrV2NamespacesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	ak := config.AccessKey
	sk := config.SecretKey
	auth := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		Build()

	client := swr.NewSwrClient(
		swr.SwrClientBuilder().
			WithRegion(region.ValueOf("cn-north-4")).
			WithCredential(auth).
			Build())

	request := &model.DeleteNamespacesRequest{}
	request.Namespace = d.Get("namespace").(string)

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err := client.DeleteNamespaces(request)
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
		return fmtp.DiagErrorf("[Error]Error Delete Swr_Namespace:%+v\n", err)
	} else {
		d.SetId("")
	}
	return nil
}
