package cmdb

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	entity2 "github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
	"io/ioutil"
	"time"
)

func ResourceAomEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceAomEnvironmentCreate,
		ReadContext:   ResourceAomEnvironmentRead,
		UpdateContext: ResourceAomEnvironmentUpdate,
		DeleteContext: ResourceAomEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Computed: true,
				ForceNew: true,
			},
			"env_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"env_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"component_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"register_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"aom_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"env_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"env_tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"tag_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceAomEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(cfg)
	if diaErr != nil {
		return diaErr
	}

	opts := entity2.EnvParam{
		ComponentId:  d.Get("component_id").(string),
		Description:  d.Get("description").(string),
		EnvName:      d.Get("env_name").(string),
		EnvType:      d.Get("env_type").(string),
		OsType:       d.Get("os_type").(string),
		Region:       cfg.GetRegion(d),
		RegisterType: d.Get("register_type").(string),
	}
	client.WithMethod(httpclient_go.MethodPost).
		WithUrlWithoutEndpoint(cfg, "aom", cfg.GetRegion(d), "v1/environments").WithBody(opts)
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error create Environment fields %s: %s", opts, err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error convert data %s, %s", string(body), err)
	}
	if response.StatusCode == 200 {
		rlt := &entity2.CreateModelVo{}
		err = json.Unmarshal(body, rlt)
		if err != nil {
			return diag.Errorf("error convert data %s, %s", string(body), err)
		}
		d.SetId(rlt.Id)
		return ResourceAomEnvironmentRead(ctx, d, meta)
	}
	return diag.Errorf("error create Environment %v. error: %s", opts, string(body))
}

func ResourceAomEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(cfg)
	if diaErr != nil {
		return diaErr
	}

	client.WithMethod(httpclient_go.MethodGet).
		WithUrlWithoutEndpoint(cfg, "aom", cfg.GetRegion(d), "v1/environments/"+d.Id())
	response, err := client.Do()

	body, diags := client.CheckDeletedDiag(d, err, response, "error retrieving Environment")
	if body == nil {
		return diags
	}

	rlt := &entity2.EnvVo{}
	err = json.Unmarshal(body, rlt)
	if err != nil {
		return diag.Errorf("error retrieving Environment %s", d.Id())
	}

	mErr := multierror.Append(nil,
		d.Set("aom_id", rlt.AomId),
		d.Set("component_id", rlt.ComponentId),
		d.Set("create_time", rlt.CreateTime),
		d.Set("creator", rlt.Creator),
		d.Set("description", rlt.Description),
		d.Set("env_id", rlt.EnvId),
		d.Set("env_name", rlt.EnvName),
		d.Set("env_type", rlt.EnvType),
		d.Set("enterprise_project_id", rlt.EpsId),
		d.Set("modified_time", rlt.ModifiedTime),
		d.Set("modifier", rlt.Modifier),
		d.Set("os_type", rlt.OsType),
		d.Set("region", rlt.Region),
		d.Set("register_type", rlt.RegisterType),
	)
	var envTags []map[string]interface{}
	for _, obj := range rlt.EnvTags {
		envTag := make(map[string]string)
		envTag["tag_id"] = obj.TagId
		envTag["tag_name"] = obj.TagName
	}
	mErr = multierror.Append(mErr, d.Set("env_tags", envTags))

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Environment fields: %s", err)
	}

	return nil
}

func ResourceAomEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(cfg)
	if diaErr != nil {
		return diaErr
	}
	opts := entity2.EnvParam{
		ComponentId:  d.Get("component_id").(string),
		Description:  d.Get("description").(string),
		EnvName:      d.Get("env_name").(string),
		EnvType:      d.Get("env_type").(string),
		OsType:       d.Get("os_type").(string),
		Region:       cfg.GetRegion(d),
		RegisterType: d.Get("register_type").(string),
	}

	client.WithMethod(httpclient_go.MethodPut).
		WithUrlWithoutEndpoint(cfg, "aom", cfg.GetRegion(d), "v1/environments/"+d.Id()).WithBody(opts)
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error update Environment fields %s: %s", opts, err)
	}

	if response.StatusCode == 200 {
		return nil
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error update Environment %s: %s", string(body), err)
	}
	return diag.Errorf("error update Environment %s:  %s", opts, string(body))
}

func ResourceAomEnvironmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(cfg)
	if diaErr != nil {
		return diaErr
	}

	client.WithMethod(httpclient_go.MethodDelete).
		WithUrlWithoutEndpoint(cfg, "aom", cfg.GetRegion(d), "v1/environments/"+d.Id())
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error delete Environment %s: %s", d.Id(), err)
	}

	if response.StatusCode == 200 {
		return nil
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error delete Environment %s: %s", d.Id(), err)
	}
	return diag.Errorf("error delete Environment %s:  %s", d.Id(), string(body))
}
