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
	"github.com/gravitational/trace"
)

// resourceTeleportOIDCConnectorType is the resource metadata type
type resourceTeleportOIDCConnectorType struct{}

// resourceTeleportOIDCConnector is the resource
type resourceTeleportOIDCConnector struct {
	p Provider
}

// GetSchema returns the resource schema
func (r resourceTeleportOIDCConnectorType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfschema.GenSchemaOIDCConnectorV3(ctx)
}

// NewResource creates the empty resource
func (r resourceTeleportOIDCConnectorType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceTeleportOIDCConnector{
		p: *(p.(*Provider)),
	}, nil
}

// Create creates the provision token
func (r resourceTeleportOIDCConnector) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.IsConfigured(resp.Diagnostics) {
		return
	}

	var plan types.Object
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	oidcConnector := &apitypes.OIDCConnectorV3{}
	diags = tfschema.CopyOIDCConnectorV3FromTerraform(ctx, plan, oidcConnector)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	

	_, err := r.p.Client.GetOIDCConnector(ctx, oidcConnector.Metadata.Name, true)
	if !trace.IsNotFound(err) {
		if err == nil {
			n := oidcConnector.Metadata.Name
			existErr := fmt.Sprintf("OIDCConnector exists in Teleport. Either remove it (tctl rm oidc/%v)"+
				" or import it to the existing state (terraform import teleport_app.%v %v)", n, n, n)

			resp.Diagnostics.Append(diagFromErr("OIDCConnector exists in Teleport", trace.Errorf(existErr)))
			return
		}

		resp.Diagnostics.Append(diagFromWrappedErr("Error reading OIDCConnector", trace.Wrap(err), "oidc"))
		return
	}

	err = oidcConnector.CheckAndSetDefaults()
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error setting OIDCConnector defaults", trace.Wrap(err), "oidc"))
		return
	}

	err = r.p.Client.UpsertOIDCConnector(ctx, oidcConnector)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error creating OIDCConnector", trace.Wrap(err), "oidc"))
		return
	}

	id := oidcConnector.Metadata.Name
	oidcConnectorI, err := r.p.Client.GetOIDCConnector(ctx, id, true)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading OIDCConnector", trace.Wrap(err), "oidc"))
		return
	}

	oidcConnector, ok := oidcConnectorI.(*apitypes.OIDCConnectorV3)
	if !ok {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading OIDCConnector", trace.Errorf("Can not convert %T to OIDCConnectorV3", oidcConnectorI), "oidc"))
		return
	}

	diags = tfschema.CopyOIDCConnectorV3ToTerraform(ctx, *oidcConnector, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Attrs["id"] = types.String{Value: oidcConnector.Metadata.Name}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read reads teleport OIDCConnector
func (r resourceTeleportOIDCConnector) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
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

	oidcConnectorI, err := r.p.Client.GetOIDCConnector(ctx, id.Value, true)
	if trace.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading OIDCConnector", trace.Wrap(err), "oidc"))
		return
	}

	oidcConnector := oidcConnectorI.(*apitypes.OIDCConnectorV3)
	diags = tfschema.CopyOIDCConnectorV3ToTerraform(ctx, *oidcConnector, &state)
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

// Update updates teleport OIDCConnector
func (r resourceTeleportOIDCConnector) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	if !r.p.IsConfigured(resp.Diagnostics) {
		return
	}

	var plan types.Object
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	oidcConnector := &apitypes.OIDCConnectorV3{}
	diags = tfschema.CopyOIDCConnectorV3FromTerraform(ctx, plan, oidcConnector)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := oidcConnector.Metadata.Name

	err := oidcConnector.CheckAndSetDefaults()
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error updating OIDCConnector", err, "oidc"))
		return
	}

	err = r.p.Client.UpsertOIDCConnector(ctx, oidcConnector)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error updating OIDCConnector", err, "oidc"))
		return
	}

	oidcConnectorI, err := r.p.Client.GetOIDCConnector(ctx, name, true)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading OIDCConnector", err, "oidc"))
		return
	}

	oidcConnector = oidcConnectorI.(*apitypes.OIDCConnectorV3)
	diags = tfschema.CopyOIDCConnectorV3ToTerraform(ctx, *oidcConnector, &plan)
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

// Delete deletes Teleport OIDCConnector
func (r resourceTeleportOIDCConnector) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var id types.String
	diags := req.State.GetAttribute(ctx, tftypes.NewAttributePath().WithAttributeName("metadata").WithAttributeName("name"), &id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.p.Client.DeleteOIDCConnector(ctx, id.Value)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error deleting OIDCConnectorV3", trace.Wrap(err), "oidc"))
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports OIDCConnector state
func (r resourceTeleportOIDCConnector) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	oidcConnectorI, err := r.p.Client.GetOIDCConnector(ctx, req.ID, true)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading OIDCConnector", trace.Wrap(err), "oidc"))
		return
	}

	oidcConnector := oidcConnectorI.(*apitypes.OIDCConnectorV3)

	var state types.Object

	diags := resp.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = tfschema.CopyOIDCConnectorV3ToTerraform(ctx, *oidcConnector, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Attrs["id"] = types.String{Value: oidcConnector.Metadata.Name}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
