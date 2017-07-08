// Package pkg wraps the macOS installer for pkg files.
package pkg

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"syscall"
	"time"
)

// Option customizes an installer command.
type Option func(*installer)

type installer struct {
	ctx                     context.Context
	choicesXML              *installerChoices
	installerChoicesXML     []byte
	installerChoicesXMLPath string
	args                    []string
	customEnv               []string

	appliedOpts bool
}

func (o *installer) apply(opts ...Option) error {
	if o.appliedOpts {
		return nil
	}
	for _, opt := range opts {
		opt(o)
	}

	if o.choicesXML != nil {
		if err := o.choicesXML.write(); err != nil {
			return err
		}
	}
	o.appliedOpts = true
	return nil
}

func (o *installer) cleanup() error {
	if o.choicesXML != nil && o.choicesXML.written {
		return os.Remove(o.choicesXML.path)
	}
	return nil
}

type installerChoices struct {
	written bool
	xml     []byte
	path    string
}

func (i *installerChoices) write() error {
	if i.written {
		return nil
	}
	if err := ioutil.WriteFile(i.path, i.xml, 0644); err != nil {
		return err
	}
	i.written = true
	return nil

}

// AllowUntrusted allows installing packages with expired certificates.
func AllowUntrusted() Option {
	return func(o *installer) {
		o.args = append(o.args, "-allowUntrusted")
	}
}

// ApplyChoiceChangesXML uses a ChoiceChangesXML file during the package installation.
// See https://github.com/munki/munki/wiki/ChoiceChangesXML
func ApplyChoiceChangesXML(xml []byte) Option {
	return func(o *installer) {
		o.choicesXML = &installerChoices{
			xml:  xml,
			path: filepath.Join(os.TempDir(), "choices.xml"),
		}
		o.args = append(o.args, "-applyChoiceChangesXML", o.choicesXML.path)
	}
}

// WithContext allows setting a custom context when calling exec.CommandContext
// Useful for overriding the default 1 hour timeout.
func WithContext(ctx context.Context) Option {
	return func(o *installer) {
		o.ctx = ctx
	}
}

// WithCustomEnv adds custom environment variables to the exec environment.
//
// The values in the env slice must be in the KEY=VALUE format.
// A common use is setting
func WithCustomEnv(env []string) Option {
	return func(o *installer) {
		o.customEnv = env
	}
}

// Install installs a macOS pkg, returning a restart action on success.
func Install(pkgpath string, opts ...Option) (restart bool, err error) {
	o := new(installer)
	o.ctx = context.Background()
	if err := o.apply(opts...); err != nil {
		return false, err
	}
	defer o.cleanup()

	restart, err = needsRestart(pkgpath, o.args...)
	if err != nil {
		return false, err
	}

	if _, ok := o.ctx.Deadline(); !ok {
		var cancel func()
		o.ctx, cancel = context.WithTimeout(o.ctx, 1*time.Hour)
		defer cancel()
	}

	cmd := installcmd(o.ctx, pkgpath, o.args...)
	cmd.Env, err = getEnvironment(o.customEnv)
	if err != nil {
		return false, err
	}

	if _, err := cmd.Output(); err != nil {
		return false, err
	}

	return restart, nil
}

func getEnvironment(custom []string) ([]string, error) {
	env := os.Environ()

	var hasUser bool
	for _, v := range custom {
		if v == "USER=CURRENT_CONSOLE_USER" {
			fi, err := os.Stat("/dev/console")
			if err != nil {
				return nil, err
			}
			switch t := fi.Sys().(type) {
			case *syscall.Stat_t:
				current, err := user.LookupId(fmt.Sprintf("%d", t.Uid))
				if err != nil {
					return nil, err
				}
				env = append(env, fmt.Sprintf("USER=%s", current.Name))
				env = append(env, fmt.Sprintf("HOME=%s", current.HomeDir))
				hasUser = true
			default:
				fmt.Printf("unknown fileInfo %T", t)
				continue
			}
			continue
		}
		env = append(env, v)
	}

	// only set the user vars if they weren't set already.
	if !hasUser {
		usrinfo, err := user.LookupId("0")
		if err != nil {
			return nil, err
		}
		env = append(env, fmt.Sprintf("USER=%s", usrinfo.Name))
		env = append(env, fmt.Sprintf("HOME=%s", usrinfo.HomeDir))
	}

	return env, nil
}

func installcmd(ctx context.Context, pkgpath string, extraArgs ...string) *exec.Cmd {
	installer := "/usr/sbin/installer"
	args := []string{"-verboseR", "-pkg", pkgpath, "-target", "/"}
	args = append(args, extraArgs...)
	return exec.CommandContext(ctx, installer, args...)
}

// NeedsRestart checks if the pkg at path requires restart.
func NeedsRestart(pkgpath string, opts ...Option) (bool, error) {
	o := new(installer)
	if err := o.apply(opts...); err != nil {
		return false, err
	}
	defer o.cleanup()

	return needsRestart(pkgpath, o.args...)
}

func needsRestart(pkgpath string, extraArgs ...string) (restart bool, err error) {
	installer := "/usr/sbin/installer"
	args := []string{"-query", "RestartAction", "-pkg", pkgpath}
	args = append(args, extraArgs...)
	cmd := exec.Command(installer, args...)
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}

	action := fromOutput(out)
	restart = action.Equal(requireRestart) || action.Equal(recommendRestart)
	return restart, nil
}

type restartAction uint

const (
	none restartAction = iota
	requireRestart
	recommendRestart
)

func (r restartAction) Equal(rr restartAction) bool {
	return r == rr
}

func fromOutput(out []byte) restartAction {
	action := bytes.TrimSpace(out)
	switch {
	case bytes.Equal(action, []byte("RequireRestart")):
		return requireRestart
	case bytes.Equal(out, []byte("RequireRestart")):
		return recommendRestart
	default:
		return none
	}
}
