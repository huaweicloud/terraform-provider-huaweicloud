package swr

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SWR POST /v2/manage/agency
// @API SWR GET /v2/manage/agency
func ResourceSwrAgency() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrAgencyCreate,
		ReadContext:   resourceSwrAgencyRead,
		DeleteContext: resourceSwrAgencyDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSwrAgencyCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	// precheck whether the agency is exist, because the agency can be created once
	// or a 409 HTTP response code will be returned
	isAgency, err := getSWRAgency(client)
	if err != nil {
		return diag.FromErr(err)
	}

	if isAgency {
		log.Println("[WARN] agency already exist")
	} else {
		createHttpUrl := "v2/manage/agency"
		createPath := client.Endpoint + createHttpUrl
		createOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				204,
			},
		}

		_, err = client.Request("POST", createPath, &createOpt)
		if err != nil {
			return diag.Errorf("error creating SWR agency: %s", err)
		}
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	return nil
}

func getSWRAgency(client *golangsdk.ServiceClient) (bool, error) {
	getHttpUrl := "v2/manage/agency"
	getPath := client.Endpoint + getHttpUrl
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return false, fmt.Errorf("error getting SWR agency: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return false, fmt.Errorf("error flattening SWR agency response: %s", err)
	}

	isAgency := utils.PathSearch("is_agency", getRespBody, nil)
	if isAgency == nil {
		return false, errors.New("error finding is_agency in API response")
	}

	return isAgency.(bool), nil
}

func resourceSwrAgencyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSwrAgencyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting SWR agency is not supported. The resource is only removed from the state " +
		"and the agency still remains in the cloud"
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
