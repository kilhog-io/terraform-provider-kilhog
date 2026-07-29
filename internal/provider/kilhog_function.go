// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var (
	_ function.Function = KilhogFunction{}
)

func NewKilhogFunction() function.Function {
	return KilhogFunction{}
}

type KilhogFunction struct{}

func (r KilhogFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "echo"
}

func (r KilhogFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Echo function",
		MarkdownDescription: "Echoes given argument as result",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "input",
				MarkdownDescription: "String to echo",
			},
		},
		Return: function.StringReturn{},
	}
}

func (r KilhogFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var data string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &data))

	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, data))
}
