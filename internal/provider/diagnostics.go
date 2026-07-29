// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"errors"
	"fmt"

	kilhogsdk "github.com/kilhog-io/kilhog/pkg/kilhog"
)

func apiErrorSummary(err error, action string) string {
	var apiErr *kilhogsdk.APIError
	if errors.As(err, &apiErr) {
		return fmt.Sprintf("%s failed (HTTP %d): %s", action, apiErr.StatusCode, apiErr.Message)
	}

	return fmt.Sprintf("%s failed: %s", action, err.Error())
}
