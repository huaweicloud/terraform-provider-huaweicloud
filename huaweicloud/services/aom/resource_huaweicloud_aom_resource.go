package aom

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
	"io/ioutil"
	"strings"
	"time"
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
				Required: true,
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
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_region": {
				Type:     schema.TypeString,
				Required: true,
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
			"maker": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"keywords": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ci_relationships": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"ci_region": {
				Type:     schema.TypeString,
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

func buildResourceOpts(d *schema.ResourceData) []entity.ResourceImportDetailParam {
	return []entity.ResourceImportDetailParam{
		{
			ResourceId:     d.Get("resource_id").(string),
			ResourceName:   d.Get("resource_name").(string),
			ResourceRegion: d.Get("resource_rgion").(string),
			ProjectId:      d.Get("project_id").(string),
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
	cfg := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(cfg)
	if diaErr != nil {
		return diaErr
	}

	opts := entity.ResourceImportParam{
		EnvId:     d.Get("env_id").(string),
		Resources: buildResourceOpts(d),
	}

	client.WithMethod(httpclient_go.MethodPut).
		WithUrlWithoutEndpoint(cfg, "aom", cfg.GetRegion(d), "v1/resource/"+d.Get("rf_resource_type").(string)+
			"/type/"+d.Get("type").(string)+"/ci-relationships").WithBody(opts)
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
	client, diaErr := httpclient_go.NewHttpClientGo(cfg)
	if diaErr != nil {
		return diaErr
	}

	opts := entity.PageResourceListParam{
		CiId:     d.Get("env_id").(string),
		CiType:   "environment",
		Keywords: map[string]string{"RESOURCE_ID": d.Get("resource_id").(string)},
	}

	client.WithMethod(httpclient_go.MethodPost).
		WithUrlWithoutEndpoint(cfg, "aom", cfg.GetRegion(d), "v1/resource/"+d.Get("rf_resource_type").(string)+
			"/type/"+d.Get("type").(string)+"/ci-relationships").WithBody(opts)
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
	client, diaErr := httpclient_go.NewHttpClientGo(cfg)
	if diaErr != nil {
		return diaErr
	}

	opts := entity.DeleteResourceParam{
		Data: buildDeleteResourceOpts(d),
	}

	client.WithMethod(httpclient_go.MethodPost).
		WithUrlWithoutEndpoint(cfg, "aom", cfg.GetRegion(d), "v1/resource/"+d.Get("rf_resource_type").(string)+
			"/type/"+d.Get("type").(string)+"/ci-relationships").WithBody(opts)

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
