package obs

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/obs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type GrantType string

const (
	GrantOwner           GrantType = "OWNER"
	GrantPublic          GrantType = "PUBLIC"
	GrantLogDeliveryUser GrantType = "LOG_DELIVERY_USER"
	GrantAccount         GrantType = "ACCOUNT"
)

func ResourceOBSBucketAcl() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOBSBucketAclCreate,
		ReadContext:   resourceOBSBucketAclRead,
		UpdateContext: resourceOBSBucketAclCreate,
		DeleteContext: resourceOBSBucketAclDelete,
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
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the bucket to which to set the acl.`,
			},
			"owner_permission": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        permissionSchema(),
				Optional:    true,
				Computed:    true,
				Description: `Specifies the bucket owner permission.`,
			},
			"public_permission": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        permissionSchema(),
				Optional:    true,
				Description: `Specifies the public permission.`,
			},
			"log_delivery_user_permission": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        permissionSchema(),
				Optional:    true,
				Description: `Specifies the log delivery user permission.`,
			},
			"account_permission": {
				Type:        schema.TypeList,
				Elem:        accountPermissionSchema(),
				Optional:    true,
				Description: `Specifies the account permissions.`,
			},
		},
	}
}

func permissionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"access_to_bucket": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the access to bucket. Valid values are **READ** and **WRITE**.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"READ", "WRITE",
					}, false),
				},
			},
			"access_to_acl": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the access to acl. Valid values are **READ_ACP** and **WRITE_ACP**.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"READ_ACP", "WRITE_ACP",
					}, false),
				},
			},
		},
	}
	return &sc
}

func accountPermissionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Specifies the account id to authorize. The account id cannot be the bucket owner, 
and must be unique.`,
			},
			"access_to_bucket": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the access to bucket. Valid values are **READ** and **WRITE**.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"READ", "WRITE",
					}, false),
				},
			},
			"access_to_acl": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the access to acl. Valid values are **READ_ACP** and **WRITE_ACP**.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"READ_ACP", "WRITE_ACP",
					}, false),
				},
			},
		},
	}
	return &sc
}

func buildPermissionTypesFromRawMap(rawMap map[string]interface{}) []obs.PermissionType {
	var accesses []string
	if accessArray, ok := rawMap["access_to_bucket"].([]interface{}); ok {
		accesses = append(accesses, utils.ExpandToStringList(accessArray)...)
	}
	if accessArray, ok := rawMap["access_to_acl"].([]interface{}); ok {
		accesses = append(accesses, utils.ExpandToStringList(accessArray)...)
	}

	permissionTypes := make([]obs.PermissionType, len(accesses))
	for i, access := range accesses {
		switch access {
		case "READ":
			permissionTypes[i] = obs.PermissionRead
		case "WRITE":
			permissionTypes[i] = obs.PermissionWrite
		case "READ_ACP":
			permissionTypes[i] = obs.PermissionReadAcp
		default:
			permissionTypes[i] = obs.PermissionWriteAcp
		}
	}
	return permissionTypes
}

func buildUserTypeGrants(rawMap map[string]interface{}, grants []obs.Grant, domainID string) []obs.Grant {
	permissionTypes := buildPermissionTypesFromRawMap(rawMap)
	log.Printf("[DEBUG] The grantee user permission types: %v.", permissionTypes)
	for _, permissionType := range permissionTypes {
		grants = append(grants, obs.Grant{
			Permission: permissionType,
			Grantee: obs.Grantee{
				Type: obs.GranteeUser,
				ID:   domainID,
			},
		})
	}
	return grants
}

func buildGroupTypeGrants(rawMap map[string]interface{}, grants []obs.Grant, uriType obs.GroupUriType) []obs.Grant {
	permissionTypes := buildPermissionTypesFromRawMap(rawMap)
	log.Printf("[DEBUG] The grantee group permission types: %v.", permissionTypes)
	for _, permissionType := range permissionTypes {
		grants = append(grants, obs.Grant{
			Permission: permissionType,
			Grantee: obs.Grantee{
				Type: obs.GranteeGroup,
				URI:  uriType,
			},
		})
	}
	return grants
}

func findCurrentOwnerGrant(obsClient *obs.ObsClient, d *schema.ResourceData, cfg *config.Config) (*obs.Grant, error) {
	bucket := d.Get("bucket").(string)
	output, err := obsClient.GetBucketAcl(bucket)
	if err != nil {
		return nil, getObsError("Error retrieving OBS bucket current acl", bucket, err)
	}
	for _, grant := range output.Grants {
		grantee := grant.Grantee
		if grantee.Type == obs.GranteeUser && grantee.ID == cfg.DomainID {
			return &grant, nil
		}
	}
	return nil, getObsError("Cannot find owner grant from current grants", bucket, err)
}

func buildObsBucketAclGrantsParam(obsClient *obs.ObsClient, d *schema.ResourceData,
	cfg *config.Config) ([]obs.Grant, error) {
	var grants []obs.Grant
	if rawArray, ok := d.Get("owner_permission").([]interface{}); ok && len(rawArray) > 0 {
		if rawMap, rawOk := rawArray[0].(map[string]interface{}); rawOk {
			grants = buildUserTypeGrants(rawMap, grants, cfg.DomainID)
		}
	}

	if len(grants) == 0 {
		// the owner permission is empty, read the current owner permission. make the owner has permissions
		grant, err := findCurrentOwnerGrant(obsClient, d, cfg)
		if err != nil {
			return nil, err
		}
		grants = append(grants, *grant)
	}

	if rawArray, ok := d.Get("account_permission").([]interface{}); ok && len(rawArray) > 0 {
		for _, raw := range rawArray {
			if rawMap, rawOk := raw.(map[string]interface{}); rawOk {
				grants = buildUserTypeGrants(rawMap, grants, rawMap["account_id"].(string))
			}
		}
	}

	if rawArray, ok := d.Get("public_permission").([]interface{}); ok && len(rawArray) > 0 {
		if rawMap, rawOk := rawArray[0].(map[string]interface{}); rawOk {
			grants = buildGroupTypeGrants(rawMap, grants, obs.GroupAllUsers)
		}
	}

	if rawArray, ok := d.Get("log_delivery_user_permission").([]interface{}); ok && len(rawArray) > 0 {
		if rawMap, rawOk := rawArray[0].(map[string]interface{}); rawOk {
			grants = buildGroupTypeGrants(rawMap, grants, obs.GroupLogDelivery)
		}
	}
	return grants, nil
}

func resourceOBSBucketAclCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	obsClient, err := cfg.ObjectStorageClient(region)
	if err != nil {
		return diag.Errorf("error creating OBS Client: %s", err)
	}

	bucket := d.Get("bucket").(string)
	grantParam, err := buildObsBucketAclGrantsParam(obsClient, d, cfg)
	if err != nil {
		return diag.FromErr(err)
	}

	params := &obs.SetBucketAclInput{
		Bucket: bucket,
	}
	params.Owner.ID = cfg.DomainID
	params.Grants = grantParam
	_, err = obsClient.SetBucketAcl(params)
	if err != nil {
		return diag.FromErr(getObsError("Error creating OBS bucket acl", bucket, err))
	}

	d.SetId(bucket)
	return resourceOBSBucketAclRead(ctx, d, meta)
}

func flattenGrantsByPermissionType(output *obs.GetBucketAclOutput, cfg *config.Config) map[GrantType][]obs.Grant {
	var ownerGrants []obs.Grant
	var accountGrants []obs.Grant
	var logDeliveryUserGrants []obs.Grant
	var publicGrants []obs.Grant

	for _, grant := range output.Grants {
		grantee := grant.Grantee
		if grantee.Type == obs.GranteeUser && grantee.ID == cfg.DomainID {
			// owner grants
			ownerGrants = append(ownerGrants, grant)
			continue
		}
		if grantee.Type == obs.GranteeUser && grantee.ID != cfg.DomainID {
			// account grants
			accountGrants = append(accountGrants, grant)
			continue
		}

		granteeURI := obs.GroupUriType(parseGranteeURI(grantee.URI))
		if grantee.Type == obs.GranteeGroup && granteeURI == obs.GroupLogDelivery {
			// log delivery user grants
			logDeliveryUserGrants = append(logDeliveryUserGrants, grant)
			continue
		}
		if grantee.Type == obs.GranteeGroup && granteeURI == obs.GroupAllUsers {
			// public grants
			publicGrants = append(publicGrants, grant)
		}
	}
	return map[GrantType][]obs.Grant{
		GrantOwner:           ownerGrants,
		GrantAccount:         accountGrants,
		GrantLogDeliveryUser: logDeliveryUserGrants,
		GrantPublic:          publicGrants,
	}
}

// parseGranteeURI use to parse uri. For example: http://acs.amazonaws.com/groups/global/AllUsers
func parseGranteeURI(granteeURI obs.GroupUriType) string {
	uri := string(granteeURI)
	if len(uri) == 0 {
		return ""
	}
	uriArray := strings.Split(uri, "/")
	if len(uriArray) == 0 {
		return ""
	}
	return uriArray[len(uriArray)-1]
}

func buildAccessesFromGrants(grant obs.Grant, accessToBucket []string, accessToAcl []string) (bucketAccesses []string,
	aclAccesses []string) {
	switch grant.Permission {
	case obs.PermissionRead:
		accessToBucket = append(accessToBucket, "READ")
	case obs.PermissionWrite:
		accessToBucket = append(accessToBucket, "WRITE")
	case obs.PermissionReadAcp:
		accessToAcl = append(accessToAcl, "READ_ACP")
	case obs.PermissionWriteAcp:
		accessToAcl = append(accessToAcl, "WRITE_ACP")
	case obs.PermissionFullControl:
		// permission type `PermissionFullControl` means domain id has all permission
		accessToBucket = append(accessToBucket, "READ", "WRITE")
		accessToAcl = append(accessToAcl, "READ_ACP", "WRITE_ACP")
	default:
		log.Printf("[WARN] The grant permission: %s not support.", grant.Permission)
	}
	return accessToBucket, accessToAcl
}

func flattenPermission(grants []obs.Grant) []map[string]interface{} {
	if len(grants) == 0 {
		return nil
	}
	var accessToBucket []string
	var accessToAcl []string
	for _, grant := range grants {
		accessToBucket, accessToAcl = buildAccessesFromGrants(grant, accessToBucket, accessToAcl)
	}
	if len(accessToBucket) == 0 && len(accessToAcl) == 0 {
		return nil
	}
	ownerPermissionMap := map[string]interface{}{
		"access_to_bucket": accessToBucket,
		"access_to_acl":    accessToAcl,
	}
	return []map[string]interface{}{ownerPermissionMap}
}

func flattenAccountPermission(grants []obs.Grant) []map[string]interface{} {
	if len(grants) == 0 {
		return nil
	}
	accountIDSet := make(map[string]bool)
	accessToBucketMap := make(map[string][]string)
	accessToAclMap := make(map[string][]string)
	for _, grant := range grants {
		granteeID := grant.Grantee.ID
		var accessToBucket []string
		var accessToAcl []string
		if v, ok := accessToBucketMap[granteeID]; ok {
			accessToBucket = v
		}
		if v, ok := accessToAclMap[granteeID]; ok {
			accessToAcl = v
		}
		accessToBucket, accessToAcl = buildAccessesFromGrants(grant, accessToBucket, accessToAcl)
		accessToBucketMap[granteeID] = accessToBucket
		accessToAclMap[granteeID] = accessToAcl
		accountIDSet[granteeID] = true
	}

	var m []map[string]interface{}
	for id := range accountIDSet {
		m = append(m, map[string]interface{}{
			"access_to_bucket": accessToBucketMap[id],
			"access_to_acl":    accessToAclMap[id],
			"account_id":       id,
		})
	}
	return m
}

func resourceOBSBucketAclRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	obsClient, err := cfg.ObjectStorageClient(region)
	if err != nil {
		return diag.Errorf("error creating OBS Client: %s", err)
	}

	output, err := obsClient.GetBucketAcl(d.Id())
	if err != nil {
		return diag.FromErr(getObsError("Error retrieving OBS bucket acl", d.Id(), err))
	}

	permissionTypeMap := flattenGrantsByPermissionType(output, cfg)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("bucket", d.Id()),
		d.Set("owner_permission", flattenPermission(permissionTypeMap[GrantOwner])),
		d.Set("public_permission", flattenPermission(permissionTypeMap[GrantPublic])),
		d.Set("log_delivery_user_permission", flattenPermission(permissionTypeMap[GrantLogDeliveryUser])),
		d.Set("account_permission", flattenAccountPermission(permissionTypeMap[GrantAccount])),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting OBS bucket acl fields: %s", err)
	}
	return nil
}

func resourceOBSBucketAclDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	obsClient, err := cfg.ObjectStorageClient(region)
	if err != nil {
		return diag.Errorf("error creating OBS Client: %s", err)
	}

	// the operation to delete acl is to set the owner full permissions
	params := &obs.SetBucketAclInput{
		Bucket: d.Id(),
	}
	params.Owner.ID = cfg.DomainID
	params.Grants = []obs.Grant{
		{
			Grantee: obs.Grantee{
				Type: obs.GranteeUser,
				ID:   cfg.DomainID,
			},
			Permission: obs.PermissionFullControl,
		},
	}
	_, err = obsClient.SetBucketAcl(params)
	if err != nil {
		return diag.FromErr(getObsError("Error deleting OBS bucket acl", d.Id(), err))
	}
	return nil
}
