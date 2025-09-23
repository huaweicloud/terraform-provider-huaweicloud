package huaweicloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/identity/v3/domains"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
)

// @API IAM GET /v3/auth/domains
func DataSourceAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAccountRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating IAM client: %s", err)
	}

	// ResourceBase: https://iam.{CLOUD}/v3/auth/
	identityClient.ResourceBase += "auth/"
	allPages, err := domains.List(identityClient, nil).AllPages()
	if err != nil {
		return diag.Errorf("failed to query account: %s", err)
	}

	accounts, err := domains.ExtractDomains(allPages)
	if err != nil {
		return diag.Errorf("failed to extract details of account: %s", err)
	}

	if len(accounts) == 0 {
		return diag.Errorf("failed to query account: you are not authorized to perform the action")
	}

	result := accounts[0]
	mErr := multierror.Append(
		d.Set("name", result.Name),
		d.Set("current_project_id", identityClient.ProjectID),
	)

	var user *UserDetail
	if cfg.UserID != "" {
		user, _ = queryIamUser(cfg)
	} else {
		user, _ = queryCallerUser(cfg)
	}

	if user != nil {
		mErr = multierror.Append(mErr,
			d.Set("username", user.Name),
			d.Set("user_id", user.ID),
		)
	}

	d.SetId(result.ID)
	return diag.FromErr(mErr.ErrorOrNil())
}

type UserDetail struct {
	ID       string
	Name     string
	DomainID string
}

func queryIamUser(c *config.Config) (*UserDetail, error) {
	client, err := c.IdentityV3Client(c.Region)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	uri := "/v3.0/OS-USER/users/" + c.UserID
	data, err := httphelper.New(client).
		Method("GET").
		URI(uri).
		OkCode(200).
		Request().
		Result()
	if err != nil {
		return nil, err
	}

	return &UserDetail{
		ID:       data.Get("user.id").String(),
		Name:     data.Get("user.name").String(),
		DomainID: data.Get("user.domain_id").String(),
	}, nil
}

func queryCallerUser(c *config.Config) (*UserDetail, error) {
	client, err := c.StsClient(c.Region)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	uri := "/v5/caller-identity"
	data, err := httphelper.New(client).
		Method("GET").
		URI(uri).
		OkCode(200).
		Request().
		Result()
	if err != nil {
		return nil, err
	}

	urn := data.Get("principal_urn").String()
	userName := ""
	if urn != "" && strings.Contains(urn, ":") {
		userName = urn[strings.LastIndex(urn, ":")+1:]
	}

	return &UserDetail{
		ID:       data.Get("principal_id").String(),
		Name:     userName,
		DomainID: data.Get("account_id").String(),
	}, nil
}
