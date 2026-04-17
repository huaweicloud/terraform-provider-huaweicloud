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

// @API SWR GET /v2/manage/namespaces/{namespace}/repos/{repository}/{tag}/references
func DataSourceSwrSignedImageAccossicatedImageTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSwrSignedImageAccossicatedImageTagsRead,

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
			"sig_tag": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the tag of the signed image.",
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of refenerce tags of the signed image.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceSwrSignedImageAccossicatedImageTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}
	tags, err := listAccossicatedTags(client, d)
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
		d.Set("tags", tags),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func listAccossicatedTags(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/{tag}/references"
		result  = make([]interface{}, 0)
		limit   = 10
		marker  = ""
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{namespace}", d.Get("namespace").(string))
	listPath = strings.ReplaceAll(listPath, "{repository}", d.Get("repository").(string))
	listPath = strings.ReplaceAll(listPath, "{tag}", d.Get("sig_tag").(string))

	reqOpt := &golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		listPathWithMarker := fmt.Sprintf("%s?limit=%d", listPath, limit)
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s&marker=%s", listPathWithMarker, marker)
		}

		resp, err := client.Request("GET", listPathWithMarker, reqOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		tags := utils.PathSearch("tags", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, tags...)
		if len(tags) < limit {
			break
		}

		marker = utils.PathSearch("next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}
