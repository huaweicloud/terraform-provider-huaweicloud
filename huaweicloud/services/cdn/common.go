package cdn

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"

	"github.com/chnsz/golangsdk"
)

func waitingForCdnDomainListStatusOnline(ctx context.Context, client *golangsdk.ServiceClient, domainNameList []string,
	epsID string, timeout time.Duration) error {
	mErr := multierror.Append(nil)
	for _, domainName := range domainNameList {
		if err := waitForDomainStatusAvailable(ctx, client, domainName, epsID, []string{"online"}, timeout); err != nil {
			mErr = multierror.Append(fmt.Errorf("error waiting for domain (%s) to become online: %s", domainName, err))
		}
	}
	return mErr.ErrorOrNil()
}
