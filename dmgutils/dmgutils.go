// Package dmgutils wraps the macOS hdiutil command
package dmgutils

import (
	"bytes"
	"context"
	"log"
	"os/exec"
	"strings"
	"time"

	plist "github.com/groob/plist"
)

// Option customizes an hdiutil command.
type Option func(*hdiutil)

type hdiutil struct {
	ctx              context.Context
	args             []string
	randomMountPoint bool
	useShadow        bool
	appliedOpts      bool
}

func (o *hdiutil) apply(opts ...Option) error {
	if o.appliedOpts {
		return nil
	}
	for _, opt := range opts {
		opt(o)
	}

	o.appliedOpts = true
	return nil
}

// WithContext allows setting a custom context when calling exec.CommandContext
// Useful for overriding the default 1 hour timeout.
func WithContext(ctx context.Context) Option {
	return func(o *hdiutil) {
		o.ctx = ctx
	}
}

func mountcmd(ctx context.Context, dmgpath string, extraArgs ...string) *exec.Cmd {
	hdiutil := "/usr/bin/hdiutil"
	args := []string{"attach", dmgpath, "-nobrowse", "-plist"}
	return exec.CommandContext(ctx, hdiutil, args...)
}

type dmgAttachHeader struct {
	SystemEntities []systemEntities `plist:"system-entities"`
}

type systemEntities struct {
	mounts *mounts
}

type mounts struct {
	ContentHint         string `plist:"content-hint"`
	DevEntry            string `plist:"dev-entry"`
	MountPoint          string `plist:"mount-point"`
	UnmappedContentHint string `plist:"unmapped-content-hint"`
	VolumeKind          string `plist:"volume-kind"`
}

func (p *systemEntities) UnmarshalPlist(f func(i interface{}) error) error {
	var mounts mounts
	if err := f(&mounts); err != nil {
		return err
	}
	p.mounts = &mounts

	return nil
}

// MountDMG mounts a macOS dmg, returning a mount path on success.
func MountDMG(pkgpath string, opts ...Option) (mountedpaths []string, err error) {
	o := new(hdiutil)
	mountpoints := []string{}
	o.ctx = context.Background()
	if err := o.apply(opts...); err != nil {
		return nil, err
	}

	if _, ok := o.ctx.Deadline(); !ok {
		var cancel func()
		o.ctx, cancel = context.WithTimeout(o.ctx, 30*time.Minute)
		defer cancel()
	}

	cmd := mountcmd(o.ctx, pkgpath, o.args...)

	out, _ := cmd.Output()
	buf := bytes.NewBuffer(out)

	var p dmgAttachHeader
	if err := plist.NewDecoder(buf).Decode(&p); err != nil {
		log.Fatal(err)
	}

	for _, element := range p.SystemEntities {
		mountpoints = append(mountpoints, strings.TrimSpace(element.mounts.MountPoint))
	}

	return mountpoints, nil
}