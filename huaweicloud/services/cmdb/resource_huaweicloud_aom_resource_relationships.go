package cmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
)

func ResourceCiRelationships() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceResourceCiRelationshipsCreate,
		ReadContext:   ResourceResourceCiRelationshipsRead,
		DeleteContext: ResourceResourceCiRelationshipsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: cmdbResourceImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"rf_resource_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"env_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enterprise_project_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ci_relationships": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"cmdb_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func buildResourceOpts(d *schema.ResourceData, conf *config.Config) []entity.ResourceImportDetailParam {
	projectId := d.Get("project_id").(string)
	resourceRegion := d.Get("resource_region").(string)
	if resourceRegion == "" {
		resourceRegion = conf.GetRegion(d)
	}
	if projectId == "" {
		projectId = conf.GetProjectID(conf.GetRegion(d))
	}
	return []entity.ResourceImportDetailParam{
		{
			ResourceId:     d.Get("resource_id").(string),
			ResourceName:   d.Get("resource_name").(string),
			ResourceRegion: resourceRegion,
			ProjectId:      projectId,
			EpsId:          d.Get("enterprise_project_id").(string),
			EpsName:        d.Get("enterprise_project_name").(string),
		},
	}
}

func buildDeleteResourceOpts(d *schema.ResourceData) []entity.UnbindResourceParam {
	return []entity.UnbindResourceParam{{
		Id:     d.Get("cmdb_id").(string),
		EnvIds: []string{d.Get("env_id").(string)},
	}}
}

func ResourceResourceCiRelationshipsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := httpclient_go.NewHttpClientGo(conf, "cmdb", conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("err creating Client: %s", err)
	}

	opts := entity.ResourceImportParam{
		EnvId:     d.Get("env_id").(string),
		Resources: buildResourceOpts(d, conf),
	}

	client.WithMethod(httpclient_go.MethodPut).WithUrl("v1/resource/" + d.Get("rf_resource_type").(string) +
		"/type/" + d.Get("type").(string) + "/ci-relationships").WithBody(opts)
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error Associate Resource field %s: client do error : %s", opts, err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error convert data %s, %s", string(body), err)
	}
	if response.StatusCode == 200 {
		rlt := &entity.CreateResourceResponse{}
		err = json.Unmarshal(body, rlt)
		if err != nil {
			return diag.Errorf("error convert data %s, %s", string(body), err)
		}
		if len(rlt.ResourceDetail) == 0 {
			return nil
		}
		d.SetId(rlt.ResourceDetail[0].ResourceId)
		return ResourceResourceCiRelationshipsRead(ctx, d, meta)
	}
	return diag.Errorf("error Associate Resource %v. error: %s", opts, string(body))
}

func ResourceResourceCiRelationshipsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := httpclient_go.NewHttpClientGo(cfg, "cmdb", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("err creating Client: %s", err)
	}

	opts := entity.PageResourceListParam{
		CiId:            d.Get("env_id").(string),
		CiType:          "environment",
		CiRelationships: d.Get("ci_relationships").(bool),
		Keywords:        map[string]string{"RESOURCE_ID": d.Get("resource_id").(string)},
	}

	client.WithMethod(httpclient_go.MethodPost).WithUrl("v1/resource/" + d.Get("rf_resource_type").(string) +
		"/type/" + d.Get("type").(string) + "/ci-relationships").WithBody(opts)
	response, err := client.Do()

	body, diags := client.CheckDeletedDiag(d, err, response, "error retrieving Resource")
	if body == nil {
		return diags
	}
	if err != nil {
		return diag.Errorf("error retrieving Resource %s", d.Id())
	}
	if response.StatusCode == 200 {
		rlt := &entity.ReadResourceResponse{}
		err = json.Unmarshal(body, rlt)
		if err != nil {
			return diag.Errorf("error convert data %s, %s", string(body), err)
		}
		if len(rlt.ResourceDetail) == 0 {
			d.SetId("")
		} else {
			d.SetId(rlt.ResourceDetail[0].ResourceId)
			d.Set("cmdb_id", rlt.ResourceDetail[0].Id)
			d.Set("resource_id", rlt.ResourceDetail[0].ResourceId)
			d.Set("resource_name", rlt.ResourceDetail[0].ResourceName)
			d.Set("resource_region", rlt.ResourceDetail[0].ResourceRegion)
		}
		return nil
	}
	return diag.Errorf("error Read Resource fields %v : %s", opts, string(body))
}

func ResourceResourceCiRelationshipsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := httpclient_go.NewHttpClientGo(cfg, "cmdb", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("err creating Client: %s", err)
	}

	opts := entity.DeleteResourceParam{
		Data: buildDeleteResourceOpts(d),
	}

	client.WithMethod(httpclient_go.MethodDelete).WithUrl("v1/resource/" + d.Get("rf_resource_type").(string) +
		"/type/" + d.Get("type").(string) + "/ci-relationships").WithBody(opts)

	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error delete Resource fields %s:client to error %s", opts, err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error delete Resource convert data %s to %v: %s", string(body), opts, err)
	}

	if response.StatusCode == 200 && !strings.Contains(string(body), "error_msg") {
		d.SetId("")
		return nil
	}

	return diag.Errorf("error delete Resource %s:  %s", d.Id(), string(body))
}

func cmdbResourceImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 4)
	if len(parts) != 4 {
		return nil, fmt.Errorf("Invalid format specified for import id, must be " +
			"<rf_resource_type>/<type>/<env_id>/<resource_id>")
	}

	d.SetId(parts[3])
	mErr := multierror.Append(nil,
		d.Set("rf_resource_type", parts[0]),
		d.Set("type", parts[1]),
		d.Set("env_id", parts[2]),
		d.Set("resource_id", parts[3]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
