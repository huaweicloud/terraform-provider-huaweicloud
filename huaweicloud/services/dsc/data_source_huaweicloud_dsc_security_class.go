package dsc

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC GET /v1/{project_id}/scan-templates/{template_id}/classifications
func DataSourceSecurityClass() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityClassRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"classification_trees": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     templateClassificationTreeSchema(),
			},
			"template_category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func templateClassificationTreeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"children": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     templateClassificationTreeChildrenSchema(),
			},
			"children_nums": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"depth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_nums": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}

	return &sc
}

func templateClassificationTreeChildrenSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"children_nums": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"depth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rule_nums": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}

	return &sc
}

func dataSourceSecurityClassRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/scan-templates/{template_id}/classifications"
		product = "dsc"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{template_id}", d.Get("template_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC security class: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("classification_trees", flattenTemplateClassificationTrees(
			utils.PathSearch("classification_trees", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("template_category", utils.PathSearch("template_category", respBody, nil)),
		d.Set("template_name", utils.PathSearch("template_name", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTemplateClassificationTrees(trees []interface{}) []interface{} {
	if len(trees) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(trees))
	for _, v := range trees {
		rst = append(rst, map[string]interface{}{
			"children": flattenTemplateClassificationTreesChildren(
				utils.PathSearch("children", v, make([]interface{}, 0)).([]interface{})),
			"children_nums": utils.PathSearch("children_nums", v, nil),
			"create_time":   utils.PathSearch("create_time", v, nil),
			"depth":         utils.PathSearch("depth", v, nil),
			"id":            utils.PathSearch("id", v, nil),
			"name":          utils.PathSearch("name", v, nil),
			"parent_id":     utils.PathSearch("parent_id", v, nil),
			"project_id":    utils.PathSearch("project_id", v, nil),
			"root_id":       utils.PathSearch("root_id", v, nil),
			"rule_nums":     utils.PathSearch("rule_nums", v, nil),
			"update_time":   utils.PathSearch("update_time", v, nil),
		})
	}

	return rst
}

func flattenTemplateClassificationTreesChildren(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"children_nums": utils.PathSearch("children_nums", v, nil),
			"create_time":   utils.PathSearch("create_time", v, nil),
			"depth":         utils.PathSearch("depth", v, nil),
			"id":            utils.PathSearch("id", v, nil),
			"name":          utils.PathSearch("name", v, nil),
			"parent_id":     utils.PathSearch("parent_id", v, nil),
			"project_id":    utils.PathSearch("project_id", v, nil),
			"root_id":       utils.PathSearch("root_id", v, nil),
			"rule_nums":     utils.PathSearch("rule_nums", v, nil),
			"update_time":   utils.PathSearch("update_time", v, nil),
		})
	}

	return rst
}
