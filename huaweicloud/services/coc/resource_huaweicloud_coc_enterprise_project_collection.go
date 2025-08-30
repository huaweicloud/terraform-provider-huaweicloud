package coc

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API COC PUT /v1/enterprise-project-collect
// @API COC GET /v1/enterprise-project-collect
func ResourceEnterpriseProjectCollection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnterpriseProjectCollectionCreate,
		ReadContext:   resourceEnterpriseProjectCollectionRead,
		UpdateContext: resourceEnterpriseProjectCollectionUpdate,
		DeleteContext: resourceEnterpriseProjectCollectionDelete,

		Schema: map[string]*schema.Schema{
			"ep_id_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildEnterpriseProjectCollectionOpts(epIdList []interface{}) map[string]interface{} {
	if epIdList == nil {
		epIdList = make([]interface{}, 0)
		epIdList = append(epIdList, "")
	}

	bodyParams := map[string]interface{}{
		"ep_id_list": epIdList,
	}

	return bodyParams
}

func createOrUpdateOrDeleteEnterpriseProjectCollection(client *golangsdk.ServiceClient, bodyParams map[string]interface{}) error {
	httpUrl := "v1/enterprise-project-collect"
	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(bodyParams),
	}

	_, err := client.Request("PUT", createPath, &createOpt)
	return err
}

func resourceEnterpriseProjectCollectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("coc", region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	err = createOrUpdateOrDeleteEnterpriseProjectCollection(client,
		buildEnterpriseProjectCollectionOpts(d.Get("ep_id_list").(*schema.Set).List()))
	if err != nil {
		return diag.Errorf("error creating the COC enterprise project collection: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	return resourceEnterpriseProjectCollectionRead(ctx, d, meta)
}

func getEnterpriseProjectCollection(client *golangsdk.ServiceClient) (interface{}, error) {
	httpUrl := "v1/enterprise-project-collect?limit=100"
	getPath := client.Endpoint + httpUrl
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)

	if err != nil {
		return nil, fmt.Errorf("error retrieving COC enterprise project collections: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	enterpriseProjectCollection := utils.PathSearch("data[0]", getRespBody, nil)
	if enterpriseProjectCollection == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return enterpriseProjectCollection, nil
}

func resourceEnterpriseProjectCollectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("coc", region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	enterpriseProjectCollection, err := getEnterpriseProjectCollection(client)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving the collection list of enterprise project ID")
	}

	mErr := multierror.Append(
		d.Set("ep_id_list", utils.PathSearch("ep_id_list", enterpriseProjectCollection,
			make([]interface{}, 0)).([]interface{})),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEnterpriseProjectCollectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("coc", region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	err = createOrUpdateOrDeleteEnterpriseProjectCollection(client,
		buildEnterpriseProjectCollectionOpts(d.Get("ep_id_list").(*schema.Set).List()))
	if err != nil {
		return diag.Errorf("error updating the COC enterprise project collection: %s", err)
	}

	return resourceEnterpriseProjectCollectionRead(ctx, d, meta)
}

func resourceEnterpriseProjectCollectionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("coc", region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	err = createOrUpdateOrDeleteEnterpriseProjectCollection(client, buildEnterpriseProjectCollectionOpts(nil))
	if err != nil {
		return diag.Errorf("error deleting the COC enterprise project collection: %s", err)
	}

	return nil
}
