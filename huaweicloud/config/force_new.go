package config

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FlexibleForceNew make the ForceNew of parameters configurable
// this func accepts a list of non-updatable parameters
// when non-updatable parameters are changed
// if ForceNew is enabled, the resource will be recreated
// if ForceNew is not enabled, an error will be raise
func FlexibleForceNew(keys []string) schema.CustomizeDiffFunc {
	return func(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
		cfg := meta.(*Config)
		var err error
		forceNew := cfg.GetForceNew(d)
		keysExpand := expandKeys(keys, d)
		if forceNew {
			for _, k := range keysExpand {
				if err := d.ForceNew(k); err != nil {
					log.Printf("[WARN] unable to require attribute replacement of %s: %s", k, err)
				}
			}
		} else {
			for _, k := range keysExpand {
				if d.Id() != "" && d.HasChange(k) {
					oldValue, newValue := d.GetChange(k)
					err = multierror.Append(err, fmt.Errorf("%s can't be updated, %v -> %v", k, oldValue, newValue))
				}
			}
		}

		return err
	}
}

func expandKeys(keys []string, d *schema.ResourceDiff) []string {
	res := []string{}
	for _, k := range keys {
		if strings.Contains(k, "*") {
			parts := strings.SplitN(k, ".*.", 2)
			l := len(d.Get(parts[0]).([]interface{}))
			i := 0
			var tempKeys []string
			for i < l {
				tempKeys = append(tempKeys, strings.Join([]string{parts[0], parts[1]}, fmt.Sprintf(".%s.", strconv.Itoa(i))))
				i++
			}
			res = append(res, expandKeys(tempKeys, d)...)
		} else {
			res = append(res, k)
		}
	}
	return res
}
