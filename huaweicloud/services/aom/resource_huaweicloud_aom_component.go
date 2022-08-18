package aom

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

func ResourceAomComponent() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceAomComponentCreate,
		ReadContext:   ResourceAomComponentRead,
		UpdateContext: ResourceAomComponentUpdate,
		DeleteContext: ResourceAomComponentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"model_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"model_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"aom_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"app_id": {
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
			"modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"register_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sub_app_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func ResourceAomComponentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(cfg)
	if diaErr != nil {
		return diaErr
	}
	opts := entity2.ComponentParam{
		Description: d.Get("description").(string),
		ModelType:   d.Get("model_type").(string),
		ModelId:     d.Get("model_id").(string),
		Name:        d.Get("name").(string),
	}
	client.WithMethod(httpclient_go.MethodPost).
		WithUrlWithoutEndpoint(cfg, "aom", cfg.GetRegion(d), "v1/components").WithBody(opts)
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error create Component fields %s: %s", opts, err)
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
		return ResourceAomComponentRead(ctx, d, meta)
	}
	return diag.Errorf("error create Component %v. error: %s", opts, string(body))
}

func ResourceAomComponentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(cfg)
	if diaErr != nil {
		return diaErr
	}

	client.WithMethod(httpclient_go.MethodGet).
		WithUrlWithoutEndpoint(cfg, "aom", cfg.GetRegion(d), "v1/components/"+d.Id())
	response, err := client.Do()

	body, diags := client.CheckDeletedDiag(d, err, response, "error retrieving Component")
	if body == nil {
		return diags
	}

	rlt := &entity2.ComponentVo{}
	err = json.Unmarshal(body, rlt)
	if err != nil {
		return diag.Errorf("error retrieving Component %s", d.Id())
	}

	mErr := multierror.Append(nil,
		d.Set("aom_id", rlt.AomId),
		d.Set("app_id", rlt.AppId),
		d.Set("create_time", rlt.CreateTime),
		d.Set("creator", rlt.Creator),
		d.Set("description", rlt.Description),
		d.Set("modified_time", rlt.ModifiedTime),
		d.Set("modifier", rlt.Modifier),
		d.Set("name", rlt.Name),
		d.Set("register_type", rlt.RegisterType),
		d.Set("sub_app_id", rlt.SubAppId),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Component fields: %s", err)
	}

	return nil
}

func ResourceAomComponentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(cfg)
	if diaErr != nil {
		return diaErr
	}
	opts := entity2.ComponentParam{
		Description: d.Get("description").(string),
		ModelType:   d.Get("model_type").(string),
		ModelId:     d.Get("model_id").(string),
		Name:        d.Get("name").(string),
	}

	client.WithMethod(httpclient_go.MethodPut).
		WithUrlWithoutEndpoint(cfg, "aom", cfg.GetRegion(d), "v1/components/"+d.Id()).WithBody(opts)
	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error update Component fields %s: %s", opts, err)
	}

	if response.StatusCode == 200 {
		return nil
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error update Component %s: %s", string(body), err)
	}
	return diag.Errorf("error update Component %s:  %s", opts, string(body))
}

func ResourceAomComponentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(cfg)
	if diaErr != nil {
		return diaErr
	}

	client.WithMethod(httpclient_go.MethodDelete).
		WithUrlWithoutEndpoint(cfg, "aom", cfg.GetRegion(d), "v1/components/"+d.Id())

	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error delete Component %s: %s", d.Id(), err)
	}

	if response.StatusCode == 200 {
		return nil
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error delete Component %s: %s", d.Id(), err)
	}
	return diag.Errorf("error delete Component %s:  %s", d.Id(), string(body))
}
