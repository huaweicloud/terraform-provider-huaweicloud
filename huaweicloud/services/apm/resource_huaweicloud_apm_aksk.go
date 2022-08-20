package apm

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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"io/ioutil"
	"time"
)

func ResourceApmAkSk() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceApmAkSkCreate,
		ReadContext:   ResourceApmAkSkRead,
		UpdateContext: ResourceApmAkSkUpdate,
		DeleteContext: ResourceApmAkSkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"access_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"secret_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func ResourceApmAkSkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(cfg)
	if diaErr != nil {
		return diaErr
	}

	opts := entity.CreateAkSkParam{
		Descp: d.Get("description").(string),
	}

	client.WithMethod(httpclient_go.MethodPost).
		WithUrlWithoutEndpoint(cfg, "apm", cfg.GetRegion(d), "v1/apm2/access-keys").WithBody(opts)

	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error creating aksk fields %s: client do error: %s", opts, err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return diag.Errorf("error convert data %s to %v : %s", string(body), opts, err)
	}
	if response.StatusCode == 200 {
		rlt := &entity.AkSkResultVO{}

		err = json.Unmarshal(body, rlt)
		if err != nil {
			return diag.Errorf("error convert data %s to %v : %s", string(body), opts, err)
		}
		d.SetId("success")
		d.Set("access_key", rlt.Ak)
		d.Set("secret_key", rlt.Sk)
		return nil
	}
	return diag.Errorf("error creating aksk fields %s: %s", opts, string(body))
}

func ResourceApmAkSkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(cfg)
	if diaErr != nil {
		return diaErr
	}

	client.WithMethod(httpclient_go.MethodGet).
		WithUrlWithoutEndpoint(cfg, "apm", cfg.GetRegion(d), "v1/apm2/access-keys")

	response, err := client.Do()

	body, diags := client.CheckDeletedDiag(d, err, response, "error to query akSks")
	if diags != nil {
		return diags
	}

	rlt := &entity.GetAkSkListVO{}
	err = json.Unmarshal(body, &rlt)
	if err != nil {
		return diag.Errorf("error convert data %s: %s", string(body), err)
	}
	for _, item := range rlt.AccessAkSkModels {
		if item.Ak == d.Get("access_key").(string) {
			return nil
		}
	}

	resourceID := d.Id()
	d.SetId("")
	d.Set("access_key", "")
	d.Set("secret_key", "")
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Resource not found",
			Detail:   fmt.Sprintf("the resource %s is goneand will be removed in Teraform state.", resourceID),
		},
	}
}

func ResourceApmAkSkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func ResourceApmAkSkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, diaErr := httpclient_go.NewHttpClientGo(cfg)
	if diaErr != nil {
		return diaErr
	}

	client.WithMethod(httpclient_go.MethodDelete).
		WithUrlWithoutEndpoint(cfg, "apm", cfg.GetRegion(d), "v1/apm2/access-keys/"+d.Get("access_key").(string))

	response, err := client.Do()
	if err != nil {
		return diag.Errorf("error delete aksk %s: %s", d.Get("application_id"), err)
	}

	if response.StatusCode == 200 {
		return nil
	}
	mErr := &multierror.Error{}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		mErr = multierror.Append(mErr, err)
	}

	rlt := &entity.ErrorResp{}
	err = json.Unmarshal(body, rlt)

	if err != nil {
		mErr = multierror.Append(mErr, err)
	}
	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error create A fields: %w", err)
	}

	return nil
}
