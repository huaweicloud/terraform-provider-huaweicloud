package elb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ELB GET /v3/{project_id}/elb/system-security-policies
// @API ELB GET /v3/{project_id}/elb/security-policies
func DataSourceElbSecurityPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceElbSecurityPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cipher": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type: schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					"system", "custom",
				}, false),
				Optional: true,
			},
			"security_policies": {
				Type:     schema.TypeList,
				Elem:     SecurityPoliciesSchema(),
				Computed: true,
			},
		},
	}
}

func SecurityPoliciesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"listeners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"protocols": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ciphers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceElbSecurityPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		listSecurityPoliciesProduct = "elb"
	)
	listSecurityPoliciesClient, err := cfg.NewServiceClient(listSecurityPoliciesProduct, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}
	var sceurityPolicies []interface{}
	v := d.Get("type").(string)
	switch v {
	case "system":
		systemSecurityPolicies, err := systemElbSecurityPoliciesRead(d, listSecurityPoliciesClient)
		if err != nil {
			return diag.FromErr(err)
		}
		sceurityPolicies = systemSecurityPolicies
	case "custom":
		customSceurityPolicies, err := customElbSecurityPoliciesRead(d, listSecurityPoliciesClient)
		if err != nil {
			return diag.FromErr(err)
		}
		sceurityPolicies = customSceurityPolicies
	default:
		allSceurityPolicies, err := allElbSecurityPoliciesRead(d, listSecurityPoliciesClient)
		if err != nil {
			return diag.FromErr(err)
		}
		sceurityPolicies = allSceurityPolicies
	}
	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("security_policies", sceurityPolicies),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func systemElbSecurityPolicies(listSecurityPoliciesClient *golangsdk.ServiceClient) (interface{}, error) {
	var (
		listSystemSecurityPoliciesHttpUrl = "v3/{project_id}/elb/system-security-policies"
	)
	listSystemSecurityPoliciesPath := listSecurityPoliciesClient.Endpoint + listSystemSecurityPoliciesHttpUrl
	listSystemSecurityPoliciesPath = strings.ReplaceAll(listSystemSecurityPoliciesPath, "{project_id}",
		listSecurityPoliciesClient.ProjectID)
	getSystemSecurityPoliciesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	listSystemSecurityPoliciesResp, err := listSecurityPoliciesClient.Request("GET",
		listSystemSecurityPoliciesPath, &getSystemSecurityPoliciesOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving ELB system security policies")
	}
	listSystemSecurityPoliciesRespBody, err := utils.FlattenResponse(listSystemSecurityPoliciesResp)
	if err != nil {
		return nil, fmt.Errorf("error returns the api response body: %s", err)
	}
	return listSystemSecurityPoliciesRespBody, nil
}

func customElbSecurityPolicies(listSecurityPoliciesClient *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var (
		listCustomSecurityPoliciesHttpUrl = "v3/{project_id}/elb/security-policies"
	)
	listSecurityPoliciesPath := listSecurityPoliciesClient.Endpoint + listCustomSecurityPoliciesHttpUrl
	listSecurityPoliciesPath = strings.ReplaceAll(listSecurityPoliciesPath, "{project_id}",
		listSecurityPoliciesClient.ProjectID)
	listSecurityPoliciesQueryParams := buildListSecurityPoliciesQueryParams(d)
	listSecurityPoliciesPath += listSecurityPoliciesQueryParams
	listSecurityPoliciesResp, err := pagination.ListAllItems(
		listSecurityPoliciesClient,
		"marker",
		listSecurityPoliciesPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving ELB custom security policies")
	}

	listSecurityPoliciesRespJson, err := json.Marshal(listSecurityPoliciesResp)
	if err != nil {
		return nil, fmt.Errorf("error json marshal: %s", err)
	}
	var listSecurityPoliciesRespBody interface{}
	err = json.Unmarshal(listSecurityPoliciesRespJson, &listSecurityPoliciesRespBody)
	if err != nil {
		return nil, fmt.Errorf("error json unmarshal: %s", err)
	}
	return listSecurityPoliciesRespBody, nil
}

func systemElbSecurityPoliciesRead(d *schema.ResourceData, listSecurityPoliciesClient *golangsdk.ServiceClient) ([]interface{}, error) {
	listSystemSecurityPoliciesRespBody, err := systemElbSecurityPolicies(listSecurityPoliciesClient)
	if err != nil {
		return nil, err
	}
	securityPolicies := flattenSystemSecurityPoliciesBody(listSystemSecurityPoliciesRespBody, d)
	return securityPolicies, nil
}

func customElbSecurityPoliciesRead(d *schema.ResourceData, listSecurityPoliciesClient *golangsdk.ServiceClient) ([]interface{}, error) {
	listCustomSecurityPoliciesRespBody, err := customElbSecurityPolicies(listSecurityPoliciesClient, d)
	if err != nil {
		return nil, err
	}
	securityPolicies := flattenCustomSecurityPoliciesBody(listCustomSecurityPoliciesRespBody, d)
	return securityPolicies, nil
}

func allElbSecurityPoliciesRead(d *schema.ResourceData, listSecurityPoliciesClient *golangsdk.ServiceClient) ([]interface{}, error) {
	listSystemSecurityPoliciesRespBody, err := systemElbSecurityPolicies(listSecurityPoliciesClient)
	if err != nil {
		return nil, err
	}
	listCustomSecurityPoliciesRespBody, err := customElbSecurityPolicies(listSecurityPoliciesClient, d)
	if err != nil {
		return nil, err
	}
	securityPolicies := flattenSystemSecurityPoliciesBody(listSystemSecurityPoliciesRespBody, d)
	securityPolicies = append(securityPolicies, flattenCustomSecurityPoliciesBody(listCustomSecurityPoliciesRespBody, d)...)
	return securityPolicies, nil
}

func buildListSecurityPoliciesQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("security_policy_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("description"); ok {
		res = fmt.Sprintf("%s&description=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenCustomSecurityPoliciesBody(resp interface{}, d *schema.ResourceData) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("security_policies", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	securityPoliciesProtocol, protocolsOk := d.GetOk("protocol")
	securityPoliciesCipher, ciphersOk := d.GetOk("cipher")
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		protocols := utils.PathSearch("protocols", v, make([]interface{}, 0)).([]interface{})
		ciphers := utils.PathSearch("ciphers", v, make([]interface{}, 0)).([]interface{})
		if protocolsOk && !utils.StrSliceContains(utils.ExpandToStringList(protocols), securityPoliciesProtocol.(string)) {
			continue
		}
		if ciphersOk && !utils.StrSliceContains(utils.ExpandToStringList(ciphers), securityPoliciesCipher.(string)) {
			continue
		}
		rst = append(rst, map[string]interface{}{
			"type":        "custom",
			"id":          utils.PathSearch("id", v, nil),
			"name":        utils.PathSearch("name", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"listeners":   utils.PathSearch("listeners", v, nil),
			"protocols":   protocols,
			"ciphers":     ciphers,
			"created_at":  utils.PathSearch("created_at", v, nil),
			"updated_at":  utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}

func flattenSystemSecurityPoliciesBody(resp interface{}, d *schema.ResourceData) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("system_security_policies", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	securityPoliciesName, nameOk := d.GetOk("name")
	securityPoliciesProtocols, protocolsOk := d.GetOk("protocol")
	securityPoliciesCiphers, ciphersOk := d.GetOk("cipher")
	securityPoliciesId, idOk := d.GetOk("security_policy_id")
	securityPoliciesdescription, descriptionOk := d.GetOk("description")
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		name := utils.PathSearch("name", v, "")
		protocols := strings.Split(utils.PathSearch("protocols", v, "").(string), "\\s")
		ciphers := strings.Split(utils.PathSearch("ciphers", v, "").(string), ":")
		id := utils.PathSearch("id", v, "")
		description := utils.PathSearch("description", v, "")
		if nameOk && securityPoliciesName.(string) != name.(string) {
			continue
		}
		if idOk && securityPoliciesId.(string) != id.(string) {
			continue
		}
		if descriptionOk && securityPoliciesdescription.(string) != description.(string) {
			continue
		}
		if protocolsOk && !utils.StrSliceContains(protocols, securityPoliciesProtocols.(string)) {
			continue
		}
		if ciphersOk && !utils.StrSliceContains(ciphers, securityPoliciesCiphers.(string)) {
			continue
		}
		rst = append(rst, map[string]interface{}{
			"type":      "system",
			"id":        id,
			"name":      name,
			"protocols": protocols,
			"ciphers":   ciphers,
		})
	}
	return rst
}
