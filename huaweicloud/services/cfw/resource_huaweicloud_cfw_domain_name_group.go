// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CFW
// ---------------------------------------------------------------

package cfw

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

// @API CFW POST /v1/{project_id}/domain-set
// @API CFW GET /v1/{project_id}/domain-sets
// @API CFW GET /v1/{project_id}/domain-set/domains/{id}
// @API CFW PUT /v1/{project_id}/domain-set/{id}
// @API CFW DELETE /v1/{project_id}/domain-set/domains/{id}
// @API CFW POST /v1/{project_id}/domain-set/domains/{id}
// @API CFW DELETE /v1/{project_id}/domain-set/{id}

func ResourceDomainNameGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainNameGroupCreate,
		ReadContext:   resourceDomainNameGroupRead,
		UpdateContext: resourceDomainNameGroupUpdate,
		DeleteContext: resourceDomainNameGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDomainNameGroupImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"fw_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the firewall instance ID.`,
			},
			"object_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the protected object ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the domain name group.`,
			},
			"type": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the type of the domain name group.`,
			},
			"domain_names": {
				Type:        schema.TypeList,
				Elem:        domainNameGroupDomainNamesSchema(),
				Optional:    true,
				Description: `Specifies the list of domain names.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the domain name group.`,
			},
			"config_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The config status of the domain name group.`,
			},
			"ref_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The reference count of the domain name group.`,
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The message of the domain name group.`,
			},
		},
	}
}

func domainNameGroupDomainNamesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the domain name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the description.`,
			},
			"domain_address_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain address ID.`,
			},
			"dns_ips": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The DNS IP list.`,
			},
		},
	}
	return &sc
}

func resourceDomainNameGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createDomainNameGroup: Create a CFW domain name group.
	var (
		createDomainNameGroupHttpUrl = "v1/{project_id}/domain-set"
		createDomainNameGroupProduct = "cfw"
	)
	createDomainNameGroupClient, err := cfg.NewServiceClient(createDomainNameGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	createDomainNameGroupPath := createDomainNameGroupClient.Endpoint + createDomainNameGroupHttpUrl
	createDomainNameGroupPath = strings.ReplaceAll(createDomainNameGroupPath, "{project_id}", createDomainNameGroupClient.ProjectID)
	createDomainNameGroupPath += fmt.Sprintf("?fw_instance_id=%v", d.Get("fw_instance_id"))

	createDomainNameGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	// domain_names is required in API but can be an empty list, so RemoveNil is not used here
	createDomainNameGroupOpt.JSONBody = buildCreateDomainNameGroupBodyParams(d)
	createDomainNameGroupResp, err := createDomainNameGroupClient.Request("POST", createDomainNameGroupPath, &createDomainNameGroupOpt)
	if err != nil {
		return diag.Errorf("error creating DomainNameGroup: %s", err)
	}

	createDomainNameGroupRespBody, err := utils.FlattenResponse(createDomainNameGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("data.id", createDomainNameGroupRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating DomainNameGroup: ID is not found in API response")
	}
	d.SetId(id)

	return resourceDomainNameGroupRead(ctx, d, meta)
}

func buildCreateDomainNameGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"fw_instance_id":  d.Get("fw_instance_id"),
		"object_id":       d.Get("object_id"),
		"name":            d.Get("name"),
		"domain_set_type": d.Get("type"),
		"domain_names":    buildCreateDomainNameGroupRequestBodyDomainNames(d.Get("domain_names")),
		"description":     utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func buildCreateDomainNameGroupRequestBodyDomainNames(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"domain_name": raw["domain_name"],
				"description": utils.ValueIgnoreEmpty(raw["description"]),
			}
		}
		return rst
	}
	return nil
}

func resourceDomainNameGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDomainNameGroup: Query the CFW domain name group detail
	getDomainNameGroupProduct := "cfw"
	client, err := cfg.NewServiceClient(getDomainNameGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	getDomainNameGroupRespBody, err := getDomainNameGroup(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DomainNameGroup")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getDomainNameGroupRespBody, nil)),
		d.Set("type", utils.PathSearch("domain_set_type", getDomainNameGroupRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getDomainNameGroupRespBody, nil)),
		d.Set("config_status", utils.PathSearch("config_status", getDomainNameGroupRespBody, nil)),
		d.Set("ref_count", utils.PathSearch("ref_count", getDomainNameGroupRespBody, nil)),
		d.Set("message", utils.PathSearch("message", getDomainNameGroupRespBody, nil)),
	)

	// getDomainNameList: Query the domain list of the CFW domain name group
	domainNameList, err := getDomainNames(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("domain_names", flattenDomainNames(domainNameList.([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getDomainNameGroup(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	getDomainNameGroupHttpUrl := "v1/{project_id}/domain-sets"
	getDomainNameGroupPath := client.Endpoint + getDomainNameGroupHttpUrl
	getDomainNameGroupPath = strings.ReplaceAll(getDomainNameGroupPath, "{project_id}", client.ProjectID)

	getDomainNameGroupQueryParams := buildGetDomainNameGroupQueryParams(d)
	getDomainNameGroupPath += getDomainNameGroupQueryParams

	getDomainNameGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getDomainNameGroupResp, err := client.Request("GET", getDomainNameGroupPath, &getDomainNameGroupOpt)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005")
	}

	getDomainNameGroupRespBody, err := utils.FlattenResponse(getDomainNameGroupResp)
	if err != nil {
		return nil, err
	}

	jsonPath := fmt.Sprintf("data.records[?set_id=='%s']|[0]", d.Id())
	getDomainNameGroupRespBody = utils.PathSearch(jsonPath, getDomainNameGroupRespBody, nil)
	if getDomainNameGroupRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return getDomainNameGroupRespBody, nil
}

func getDomainNames(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	getDomainNameListHttpUrl := "v1/{project_id}/domain-set/domains/{id}"
	getDomainNameListPath := client.Endpoint + getDomainNameListHttpUrl
	getDomainNameListPath = strings.ReplaceAll(getDomainNameListPath, "{project_id}", client.ProjectID)
	getDomainNameListPath = strings.ReplaceAll(getDomainNameListPath, "{id}", d.Id())
	getDomainNameListQueryParams := buildGetDomainNameListQueryParams(d)
	getDomainNameListPath += getDomainNameListQueryParams

	getDomainNameListOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getDomainNameListResp, err := client.Request("GET", getDomainNameListPath, &getDomainNameListOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving domain name list of domain group(%s): %s", d.Id(), err)
	}

	getDomainNameListRespBody, err := utils.FlattenResponse(getDomainNameListResp)
	if err != nil {
		return nil, err
	}

	domainNameList := utils.PathSearch("data.records", getDomainNameListRespBody, nil)
	if domainNameList == nil {
		return nil, fmt.Errorf("can not find records in response: %s", getDomainNameListRespBody)
	}
	return domainNameList, nil
}

func buildGetDomainNameGroupQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?limit=1024&offset=0&fw_instance_id=%v&object_id=%v",
		d.Get("fw_instance_id"), d.Get("object_id"))

	if v, ok := d.GetOk("name"); ok {
		res += fmt.Sprintf("&key_word=%v", v)
	}
	return res
}

func buildGetDomainNameListQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?limit=1024&offset=0&fw_instance_id=%v", d.Get("fw_instance_id"))

	return res
}

func flattenDomainNames(domainNameList []interface{}) []map[string]interface{} {
	res := make([]map[string]interface{}, len(domainNameList))
	for i, v := range domainNameList {
		raw := v.(map[string]interface{})
		res[i] = map[string]interface{}{
			"domain_name":       raw["domain_name"],
			"description":       raw["description"],
			"domain_address_id": raw["domain_address_id"],
			"dns_ips":           raw["dns_ips"],
		}
	}
	return res
}

func resourceDomainNameGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateDomainNameGroupProduct := "cfw"
	updateDomainNameGroupClient, err := cfg.NewServiceClient(updateDomainNameGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	updateDomainNameGroupChanges := []string{
		"name",
		"description",
	}

	if d.HasChanges(updateDomainNameGroupChanges...) {
		// updateDomainNameGroup: Update the configuration of CFW domain name group
		updateDomainNameGroupHttpUrl := "v1/{project_id}/domain-set/{id}"
		updateDomainNameGroupPath := updateDomainNameGroupClient.Endpoint + updateDomainNameGroupHttpUrl
		updateDomainNameGroupPath = strings.ReplaceAll(updateDomainNameGroupPath, "{project_id}", updateDomainNameGroupClient.ProjectID)
		updateDomainNameGroupPath = strings.ReplaceAll(updateDomainNameGroupPath, "{id}", d.Id())
		updateDomainNameGroupPath += fmt.Sprintf("?fw_instance_id=%v", d.Get("fw_instance_id"))

		updateDomainNameGroupOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		}

		updateDomainNameGroupOpt.JSONBody = utils.RemoveNil(buildUpdateDomainNameGroupBodyParams(d))
		_, err = updateDomainNameGroupClient.Request("PUT", updateDomainNameGroupPath, &updateDomainNameGroupOpt)
		if err != nil {
			return diag.Errorf("error updating DomainNameGroup: %s", err)
		}
	}

	if d.HasChange("domain_names") {
		// The interface requires that one domain name be retained during deletion.
		// A placeholder domain name is introduced here to facilitate the update of the domain name group.
		var placeholderDomainName interface{}
		oldRaws, newRaws := d.GetChange("domain_names")
		oldDomainNames := oldRaws.([]interface{})
		newDomainNames := newRaws.([]interface{})

		// Retain the old value of the first domain name.
		if len(oldDomainNames) > 1 {
			err := removeDomainNames(updateDomainNameGroupClient, d, oldDomainNames[1:])
			if err != nil {
				return diag.FromErr(err)
			}
		}
		// Add a placeholder domain name to ensure the old value of the first domain name can be deleted.
		if len(newDomainNames) != 0 && len(oldDomainNames) != 0 {
			placeholderDomainName, err = addPlaceholderDomainName(updateDomainNameGroupClient, d)
			if err != nil {
				return diag.FromErr(err)
			}
		}
		// Delete the old value of the first domain name.
		if len(oldDomainNames) != 0 {
			err := removeDomainNames(updateDomainNameGroupClient, d, oldDomainNames[:1])
			if err != nil {
				return diag.FromErr(err)
			}
		}
		// Add a real domain name to ensure the placeholder domain name can be deleted.
		if len(newDomainNames) != 0 {
			err := addDomainNames(updateDomainNameGroupClient, d, newDomainNames[:1])
			if err != nil {
				return diag.FromErr(err)
			}
		}

		// Delete the placeholder domain name.
		if len(newDomainNames) != 0 && len(oldDomainNames) != 0 {
			err := removeDomainNames(updateDomainNameGroupClient, d, []interface{}{placeholderDomainName})
			if err != nil {
				return diag.FromErr(err)
			}
		}

		// Complete to add the rest of the domain names.
		if len(newDomainNames) > 1 {
			err := addDomainNames(updateDomainNameGroupClient, d, newDomainNames[1:])
			if err != nil {
				return diag.FromErr(err)
			}
		}

	}
	return resourceDomainNameGroupRead(ctx, d, meta)
}

func buildUpdateDomainNameGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": d.Get("description"),
	}
	return bodyParams
}

func removeDomainNames(client *golangsdk.ServiceClient, d *schema.ResourceData, domainNameList []interface{}) error {
	if len(domainNameList) == 0 {
		return nil
	}

	id := d.Id()

	updateDomainNameGroupHttpUrl := "v1/{project_id}/domain-set/domains/{id}"
	removeDomainNamesPath := client.Endpoint + updateDomainNameGroupHttpUrl
	removeDomainNamesPath = strings.ReplaceAll(removeDomainNamesPath, "{project_id}", client.ProjectID)
	removeDomainNamesPath = strings.ReplaceAll(removeDomainNamesPath, "{id}", id)

	updateDomainNameGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	updateDomainNameGroupOpt.JSONBody = utils.RemoveNil(buildDemoveDomainNamesBodyParams(d, domainNameList))
	_, err := client.Request("DELETE", removeDomainNamesPath, &updateDomainNameGroupOpt)
	if err != nil {
		return fmt.Errorf("error removing domain names from domain name group(%s): %s", id, err)
	}

	return nil
}

func addDomainNames(client *golangsdk.ServiceClient, d *schema.ResourceData, domainNameList []interface{}) error {
	if len(domainNameList) == 0 {
		return nil
	}

	id := d.Id()

	updateDomainNameGroupHttpUrl := "v1/{project_id}/domain-set/domains/{id}"
	removeDomainNamesPath := client.Endpoint + updateDomainNameGroupHttpUrl
	removeDomainNamesPath = strings.ReplaceAll(removeDomainNamesPath, "{project_id}", client.ProjectID)
	removeDomainNamesPath = strings.ReplaceAll(removeDomainNamesPath, "{id}", id)

	updateDomainNameGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	updateDomainNameGroupOpt.JSONBody = utils.RemoveNil(buildAddDomainNamesBodyParams(d, domainNameList))
	_, err := client.Request("POST", removeDomainNamesPath, &updateDomainNameGroupOpt)
	if err != nil {
		return fmt.Errorf("error adding domain names to domain name group(%s): %s", id, err)
	}

	return nil
}

func addPlaceholderDomainName(client *golangsdk.ServiceClient, d *schema.ResourceData) (map[string]interface{}, error) {
	placeholderDomainName := map[string]interface{}{
		"domain_name": "placeholder.placeholder",
	}
	err := addDomainNames(client, d, []interface{}{placeholderDomainName})
	if err != nil {
		return nil, err
	}

	domainNameList, err := getDomainNames(client, d)
	if err != nil {
		return nil, err
	}

	for _, v := range domainNameList.([]interface{}) {
		raw := v.(map[string]interface{})
		if raw["domain_name"] == "placeholder.placeholder" {
			return raw, nil
		}
	}
	return nil, fmt.Errorf("error adding placeholder domain name")
}

func buildDemoveDomainNamesBodyParams(d *schema.ResourceData, domainNameList []interface{}) map[string]interface{} {
	if len(domainNameList) == 0 {
		return nil
	}

	domainAddressIds := make([]string, len(domainNameList))
	for i, v := range domainNameList {
		raw := v.(map[string]interface{})
		domainAddressIds[i] = raw["domain_address_id"].(string)
	}

	return map[string]interface{}{
		"object_id":          d.Get("object_id"),
		"domain_address_ids": domainAddressIds,
	}
}

func buildAddDomainNamesBodyParams(d *schema.ResourceData, domainNameList []interface{}) map[string]interface{} {
	if len(domainNameList) == 0 {
		return nil
	}

	domainNames := make([]map[string]interface{}, len(domainNameList))
	for i, v := range domainNameList {
		raw := v.(map[string]interface{})
		domainNames[i] = map[string]interface{}{
			"domain_name": raw["domain_name"],
			"description": raw["description"],
		}
	}

	return map[string]interface{}{
		"fw_instance_id": d.Get("fw_instance_id"),
		"object_id":      d.Get("object_id"),
		"domain_names":   domainNames,
	}
}

func resourceDomainNameGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteDomainNameGroup: Delete an existing CFW domain name group
	var (
		deleteDomainNameGroupHttpUrl = "v1/{project_id}/domain-set/{id}"
		deleteDomainNameGroupProduct = "cfw"
	)
	deleteDomainNameGroupClient, err := cfg.NewServiceClient(deleteDomainNameGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	deleteDomainNameGroupPath := deleteDomainNameGroupClient.Endpoint + deleteDomainNameGroupHttpUrl
	deleteDomainNameGroupPath = strings.ReplaceAll(deleteDomainNameGroupPath, "{project_id}", deleteDomainNameGroupClient.ProjectID)
	deleteDomainNameGroupPath = strings.ReplaceAll(deleteDomainNameGroupPath, "{id}", d.Id())
	deleteDomainNameGroupPath += fmt.Sprintf("?fw_instance_id=%v", d.Get("fw_instance_id"))

	deleteDomainNameGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = deleteDomainNameGroupClient.Request("DELETE", deleteDomainNameGroupPath, &deleteDomainNameGroupOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005"),
			"error deleting CFW domain name group",
		)
	}

	// Successful deletion API call does not guarantee that the resource is successfully deleted.
	// Call the details API to confirm that the resource has been successfully deleted.
	_, err = getDomainNameGroup(deleteDomainNameGroupClient, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CFW domain name group")
	}

	return diag.Errorf("error deleting CFW domain name group")
}

func resourceDomainNameGroupImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <fw_instance_id>/<object_id>/<id>")
	}

	mErr := multierror.Append(nil,
		d.Set("fw_instance_id", parts[0]),
		d.Set("object_id", parts[1]),
	)
	err := mErr.ErrorOrNil()
	if err != nil {
		return nil, err
	}
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}
