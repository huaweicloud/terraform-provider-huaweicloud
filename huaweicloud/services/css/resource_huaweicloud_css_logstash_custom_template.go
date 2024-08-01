package css

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CSS POST /v1.0/{project_id}/clusters/{cluster_id}/lgsconf/favorite
// @API CSS GET /v1.0/{project_id}/lgsconf/template
// @API CSS DELETE /v1.0/{project_id}/lgsconf/deletetemplate
func ResourceLogstashCustomTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLogstashCustomTemplateCreate,
		ReadContext:   resourceLogstashCustomTemplateRead,
		DeleteContext: resourceLogstashCustomTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"configuration_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"conf_content": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLogstashCustomTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	clusterID := d.Get("cluster_id").(string)
	cssV1Client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	addCustomTemplateHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/lgsconf/favorite"
	addCustomTemplatePath := cssV1Client.Endpoint + addCustomTemplateHttpUrl
	addCustomTemplatePath = strings.ReplaceAll(addCustomTemplatePath, "{project_id}", cssV1Client.ProjectID)
	addCustomTemplatePath = strings.ReplaceAll(addCustomTemplatePath, "{cluster_id}", clusterID)

	addCustomTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	addCustomTemplateOpt.JSONBody = utils.RemoveNil(map[string]interface{}{
		"name": d.Get("configuration_name").(string),
		"template": map[string]interface{}{
			"templateName": d.Get("name").(string),
			"desc":         utils.ValueIgnoreEmpty(d.Get("description")),
		},
	})
	_, err = cssV1Client.Request("POST", addCustomTemplatePath, &addCustomTemplateOpt)
	if err != nil {
		return diag.Errorf("error adding CSS logstash cluster custom template: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceLogstashCustomTemplateRead(ctx, d, meta)
}

func resourceLogstashCustomTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssV1Client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	getCustomTemplateHttpUrl := "v1.0/{project_id}/lgsconf/template"
	getCustomTemplatePath := cssV1Client.Endpoint + getCustomTemplateHttpUrl
	getCustomTemplatePath = strings.ReplaceAll(getCustomTemplatePath, "{project_id}", cssV1Client.ProjectID)

	getCustomTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getCustomTemplateResp, err := cssV1Client.Request("GET", getCustomTemplatePath, &getCustomTemplateOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "CSS logstash cluster custom template")
	}

	getCustomTemplateRespBody, err := utils.FlattenResponse(getCustomTemplateResp)
	if err != nil {
		return diag.Errorf("erorr retrieving CSS logstash cluster custom template: %s", err)
	}

	getCustomTemplateExp := fmt.Sprintf("customTemplates[?name=='%s']|[0]", d.Id())
	customTemplate := utils.PathSearch(getCustomTemplateExp, getCustomTemplateRespBody, nil)
	if customTemplate == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "CSS logstash cluster custom template")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("template_id", utils.PathSearch("id", customTemplate, nil)),
		d.Set("name", utils.PathSearch("name", customTemplate, nil)),
		d.Set("description", utils.PathSearch("desc", customTemplate, nil)),
		d.Set("conf_content", utils.PathSearch("confContent", customTemplate, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceLogstashCustomTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssV1Client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	deleteCustomTemplateUrl := "v1.0/{project_id}/lgsconf/deletetemplate"
	deleteCustomTemplatePath := cssV1Client.Endpoint + deleteCustomTemplateUrl
	deleteCustomTemplatePath = strings.ReplaceAll(deleteCustomTemplatePath, "{project_id}", cssV1Client.ProjectID)

	deleteCustomTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         map[string]interface{}{"name": d.Id()},
	}

	_, err = cssV1Client.Request("DELETE", deleteCustomTemplatePath, &deleteCustomTemplateOpt)
	if err != nil {
		err = common.ConvertExpected400ErrInto404Err(err, "errCode", "CSS.0001")
		return common.CheckDeletedDiag(d, err, "error deleting CSS logstash cluster custom template")
	}

	return nil
}
