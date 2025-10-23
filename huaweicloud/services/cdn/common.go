package cdn

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func waitingForCdnDomainStatusOnline(ctx context.Context, client *golangsdk.ServiceClient, domainName string,
	epsID string, timeout time.Duration) error {
	unexpectedStatus := []string{"offline", "configure_failed", "check_failed", "deleting"}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			domainResp, err := ReadCdnDomainDetail(client, domainName, epsID)
			if err != nil {
				return nil, "ERROR", err
			}

			domainStatus := utils.PathSearch("domain.domain_status", domainResp, "").(string)
			if domainStatus == "" {
				return nil, "ERROR", errors.New("error retrieving CDN domain: domain_status is not found in API response")
			}

			if domainStatus == "online" {
				return domainResp, "COMPLETED", nil
			}

			if utils.StrSliceContains(unexpectedStatus, domainStatus) {
				return domainResp, domainStatus, nil
			}
			return domainResp, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitingForCdnDomainListStatusOnline(ctx context.Context, client *golangsdk.ServiceClient, domainNameList []string,
	epsID string, timeout time.Duration) error {
	mErr := multierror.Append(nil)
	for _, domainName := range domainNameList {
		if err := waitingForCdnDomainStatusOnline(ctx, client, domainName, epsID, timeout); err != nil {
			mErr = multierror.Append(fmt.Errorf("error waiting for domain (%s) to become online: %s", domainName, err))
		}
	}
	return mErr.ErrorOrNil()
}
