// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type stringOneOfValidator struct {
	values []string
}

func stringOneOf(values ...string) validator.String {
	return &stringOneOfValidator{values: values}
}

func (v *stringOneOfValidator) Description(_ context.Context) string {
	return "value must be one of the allowed strings"
}

func (v *stringOneOfValidator) MarkdownDescription(_ context.Context) string {
	return "value must be one of the allowed strings"
}

func (v *stringOneOfValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	for _, allowed := range v.values {
		if value == allowed {
			return
		}
	}

	resp.Diagnostics.AddAttributeError(req.Path, "Invalid String Value", fmt.Sprintf("value must be one of: %s", strings.Join(v.values, ", ")))
}

var _ validator.String = (*stringOneOfValidator)(nil)
