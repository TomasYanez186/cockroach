// Copyright 2018 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License included
// in the file licenses/BSL.txt and at www.mariadb.com/bsl11.
//
// Change Date: 2022-10-01
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by the Apache License, Version 2.0,
// included in the file licenses/APL.txt and at
// https://www.apache.org/licenses/LICENSE-2.0

package flagstub

import (
	"time"

	"github.com/cockroachdb/cockroach/pkg/cmd/roachprod/vm"
	"github.com/pkg/errors"
)

// New wraps a delegate vm.Provider to only return its name and
// flags.  This allows roachprod to provide a consistent tooling
// experience. Operations that can be reasonably stubbed out are
// implemented as no-op or no-value. All other operations will
// return the provided error.
func New(delegate vm.Provider, unimplemented string) vm.Provider {
	return &provider{delegate: delegate, unimplemented: errors.New(unimplemented)}
}

type provider struct {
	delegate      vm.Provider
	unimplemented error
}

// CleanSSH implements vm.Provider and is a no-op.
func (p *provider) CleanSSH() error {
	return nil
}

// ConfigSSH implements vm.Provider and is a no-op.
func (p *provider) ConfigSSH() error {
	return nil
}

// Create implements vm.Provider and returns Unimplemented.
func (p *provider) Create(names []string, opts vm.CreateOpts) error {
	return p.unimplemented
}

// Delete implements vm.Provider and returns Unimplemented.
func (p *provider) Delete(vms vm.List) error {
	return p.unimplemented
}

// Extend implements vm.Provider and returns Unimplemented.
func (p *provider) Extend(vms vm.List, lifetime time.Duration) error {
	return p.unimplemented
}

// FindActiveAccount implements vm.Provider and returns an empty account.
func (p *provider) FindActiveAccount() (string, error) {
	return "", nil
}

// Flags implements vm.Provider and returns the delegate's name.
func (p *provider) Flags() vm.ProviderFlags {
	return p.delegate.Flags()
}

// List implements vm.Provider and returns an empty list.
func (p *provider) List() (vm.List, error) {
	return nil, nil
}

// Name implements vm.Provider and returns the delegate's name.
func (p *provider) Name() string {
	return p.delegate.Name()
}

// Active is part of the vm.Provider interface.
func (p *provider) Active() bool {
	return false
}
