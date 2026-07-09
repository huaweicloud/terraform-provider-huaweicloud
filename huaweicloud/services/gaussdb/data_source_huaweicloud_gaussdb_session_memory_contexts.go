package gaussdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/session/memory-context
func DataSourceGaussDBSessionMemoryContexts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDBSessionMemoryContextsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"session_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"memory_context_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     sessionMemoryContextInfoSchema(),
			},
		},
	}
}

func sessionMemoryContextInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"context_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"amount": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
		},
	}
}

func dataSourceGaussDBSessionMemoryContextsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/session/memory-context"
	basePath := client.Endpoint + httpUrl
	basePath = strings.ReplaceAll(basePath, "{project_id}", client.ProjectID)
	basePath = strings.ReplaceAll(basePath, "{instance_id}", d.Get("instance_id").(string))
	basePath += buildGetSessionMemoryContextsQueryParams(d)

	var allMemoryContexts []interface{}
	offset := 0
	limit := 100

	for {
		requestPath := fmt.Sprintf("%s&offset=%d&limit=%d", basePath, offset, limit)

		opts := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
				"X-Language":   "en-us",
			},
		}

		requestResp, err := client.Request("GET", requestPath, &opts)
		if err != nil {
			return diag.Errorf("error retrieving GaussDB session memory contexts: %s", err)
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return diag.FromErr(err)
		}

		curJson := utils.PathSearch("memory_context_info", respBody, make([]interface{}, 0))
		curArray := curJson.([]interface{})
		if len(curArray) == 0 {
			break
		}
		allMemoryContexts = append(allMemoryContexts, curArray...)
		offset += limit
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("memory_context_info", flattenGetSessionMemoryContextInfoBody(allMemoryContexts)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetSessionMemoryContextsQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?node_id=%s&session_id=%s", d.Get("node_id").(string), d.Get("session_id").(string))
	return res
}

func flattenGetSessionMemoryContextInfoBody(resp interface{}) []interface{} {
	curArray, ok := resp.([]interface{})
	if !ok {
		return nil
	}
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"context_name": utils.PathSearch("context_name", v, nil),
			"amount":       utils.PathSearch("amount", v, nil),
			"size":         utils.PathSearch("size", v, nil),
		})
	}
	return res
}
