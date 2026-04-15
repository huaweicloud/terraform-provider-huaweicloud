package swr

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SWR GET /v2/manage/namespaces/{namespace}/repos/{repository}/{tag}/accessories
func DataSourceSwrSignedImageAttachments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSwrSignedImageAttachmentsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the data source.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the namespace (organization) name of the signed image belongs.",
			},
			"repository": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the repository name of the signed image.",
			},
			"tag": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the tag of the signed image.",
			},
			"accessories": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of attachments of the signed image.",
				Elem:        repositoryAccessoriesSchema(),
			},
		},
	}
}

func repositoryAccessoriesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the attachment.",
			},
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the tenant the attachment belongs to.",
			},
			"namespace_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the organization (namespace) the attachment belongs to.",
			},
			"repo_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the repository the attachment belongs to.",
			},
			"sig_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The signature tag of the attachment.",
			},
			"sig_digest": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The hash value of the attachment.",
			},
			"target_digest": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The hash value of the signed image associated with the attachment.",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the attachment.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the attachment.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the attachment.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The update time of the attachment.",
			},
		},
	}
}

func dataSourceSwrSignedImageAttachmentsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}
	accessories, err := queryAccessories(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("accessories", flattenAccessories(accessories)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func queryAccessories(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/{tag}/accessories"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{namespace}", d.Get("namespace").(string))
	listPath = strings.ReplaceAll(listPath, "{repository}", d.Get("repository").(string))
	listPath = strings.ReplaceAll(listPath, "{tag}", d.Get("tag").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s?limit=100&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving SWR repository tag accessories: %s", err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		accessories := utils.PathSearch("accessories", respBody, make([]interface{}, 0)).([]interface{})
		if len(accessories) < 1 {
			break
		}
		result = append(result, accessories...)
		offset += len(accessories)
	}
	return result, nil
}

func flattenAccessories(accessories []interface{}) []interface{} {
	if len(accessories) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(accessories))
	for _, v := range accessories {
		result = append(result, map[string]interface{}{
			"id":             utils.PathSearch("id", v, nil),
			"domain_id":      utils.PathSearch("domain_id", v, nil),
			"namespace_name": utils.PathSearch("namespace_name", v, nil),
			"repo_name":      utils.PathSearch("repo_name", v, nil),
			"sig_tag":        utils.PathSearch("sig_tag", v, nil),
			"sig_digest":     utils.PathSearch("sig_digest", v, nil),
			"target_digest":  utils.PathSearch("target_digest", v, nil),
			"size":           utils.PathSearch("size", v, nil),
			"type":           utils.PathSearch("type", v, nil),
			"created_at":     utils.PathSearch("created_at", v, nil),
			"updated_at":     utils.PathSearch("updated_at", v, nil),
		})
	}
	return result
}
