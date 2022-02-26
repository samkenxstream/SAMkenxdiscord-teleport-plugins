// Code generated by _gen/main.go DO NOT EDIT
/*
Copyright 2015-2021 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/gravitational/teleport-plugins/terraform/tfschema"
	apitypes "github.com/gravitational/teleport/api/types"
)

// resourceTeleportAppType is the resource metadata type
type resourceTeleportAppType struct{}

// resourceTeleportApp is the resource
type resourceTeleportApp struct {
	p Provider
}

// GetSchema returns the resource schema
func (r resourceTeleportAppType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfschema.GenSchemaAppV3(ctx)
}

// NewResource creates the empty resource
func (r resourceTeleportAppType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceTeleportApp{
		p: *(p.(*Provider)),
	}, nil
}

// Create creates the provision token
func (r resourceTeleportApp) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.IsConfigured(resp.Diagnostics) {
		return
	}

	var plan types.Object
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	app := &apitypes.AppV3{}
	diags = tfschema.CopyAppV3FromTerraform(ctx, plan, app)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := app.CheckAndSetDefaults()
	if err != nil {
		resp.Diagnostics.AddError("Error setting App defaults", err.Error())
		return
	}

	err = r.p.Client.CreateApp(ctx, app)
	if err != nil {
		resp.Diagnostics.AddError("Error creating App", err.Error())
		return
	}

	id := app.Metadata.Name
	appI, err := r.p.Client.GetApp(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading App", err.Error())
		return
	}

	app, ok := appI.(*apitypes.AppV3)
	if !ok {
		resp.Diagnostics.AddError("Error reading App", fmt.Sprintf("Can not convert %T to AppV3", appI))
		return
	}

	diags = tfschema.CopyAppV3ToTerraform(ctx, *app, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read reads teleport App
func (r resourceTeleportApp) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state types.Object
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id types.String
	diags = req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("metadata").WithAttributeName("name"), &id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	appI, err := r.p.Client.GetApp(ctx, id.Value)
	if err != nil {
		resp.Diagnostics.AddError("Error reading App", err.Error())
		return
	}

	app := appI.(*apitypes.AppV3)
	diags = tfschema.CopyAppV3ToTerraform(ctx, *app, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates teleport App
func (r resourceTeleportApp) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	if !r.p.IsConfigured(resp.Diagnostics) {
		return
	}

	var plan types.Object
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	app := &apitypes.AppV3{}
	diags = tfschema.CopyAppV3FromTerraform(ctx, plan, app)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := app.Metadata.Name

	err := app.CheckAndSetDefaults()
	if err != nil {
		resp.Diagnostics.AddError("Error updating App", err.Error())
		return
	}

	err = r.p.Client.UpdateApp(ctx, app)
	if err != nil {
		resp.Diagnostics.AddError("Error updating App", err.Error())
		return
	}

	appI, err := r.p.Client.GetApp(ctx, name)
	if err != nil {
		resp.Diagnostics.AddError("Error reading App", err.Error())
		return
	}

	app = appI.(*apitypes.AppV3)
	diags = tfschema.CopyAppV3ToTerraform(ctx, *app, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes Teleport App
func (r resourceTeleportApp) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var id types.String
	diags := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("metadata").WithAttributeName("name"), &id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.p.Client.DeleteApp(ctx, id.Value)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting AppV3", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports App state
func (r resourceTeleportApp) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	appI, err := r.p.Client.GetApp(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading App", err.Error())
		return
	}

	app := appI.(*apitypes.AppV3)

	var state types.Object

	diags := resp.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = tfschema.CopyAppV3ToTerraform(ctx, *app, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Attrs["id"] = types.String{Value: app.Metadata.Name}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
