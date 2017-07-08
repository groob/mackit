package pkg

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"testing"
)

func Test_installcmd(t *testing.T) {
	want := []string{"/usr/sbin/installer", "-verboseR", "-pkg", "/tmp/pkgpath", "-target", "/", "-foo", "-bar"}
	cmd := installcmd(context.Background(), "/tmp/pkgpath", "-foo", "-bar")
	have := cmd.Args
	for i, v := range want {
		if have[i] != want[i] {
			t.Errorf("have %s, want %s", v, want[i])
		}
	}
}

func TestOption(t *testing.T) {
	envUser := "USER=CURRENT_CONSOLE_USER"
	opts := []Option{
		AllowUntrusted(),
		ApplyChoiceChangesXML([]byte("hello")),
		WithCustomEnv([]string{envUser}),
		WithContext(context.WithValue(context.Background(), "fooKey", "foo")),
	}

	i := new(installer)
	// apply opts multiple times
	if err := i.apply(opts...); err != nil {
		t.Fatal(err)
	}
	if err := i.apply(opts...); err != nil {
		t.Fatal(err)
	}

	argMap := make(map[string]struct{}, len(i.args))
	for _, v := range i.args {
		argMap[v] = struct{}{}
	}

	if _, ok := argMap["-allowUntrusted"]; !ok {
		t.Error("-allowUntrusted flag missing")
	}
	if _, ok := argMap["-applyChoiceChangesXML"]; !ok {
		t.Error("-applyChoiceChangesXML flag missing")
	}

	choiceFile := i.args[len(i.args)-1]
	if have, want := filepath.Base(choiceFile), "choices.xml"; have != want {
		t.Errorf("have %s, want %s", have, want)
	}

	env, err := getEnvironment(i.customEnv)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("check user env", func(t *testing.T) {
		u, err := user.Current()
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range env {
			if v == fmt.Sprintf("USER=%s", u.Name) {
				return
			}
		}
		t.Errorf("%s environment variable not set", envUser)
	})

	if _, ok := i.ctx.Value("fooKey").(string); !ok {
		t.Errorf("expected fooKey in ctx")
	}

	// test cleanup
	if err := i.cleanup(); err != nil {
		t.Fatal(err)
	}
}

func TestInstall(t *testing.T) {
	if !isRoot(t) {
		t.Skipf("skip pkg installation, needs root")
	}
	restart, err := needsRestart("/Users/victor/Downloads/munkitools-2.8.2.2855.pkg")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(restart)

	pkg := createPkg(t)
	defer os.Remove(pkg)
	restart, err = Install(pkg,
		WithCustomEnv([]string{"USER=CURRENT_CONSOLE_USER"}),
		AllowUntrusted(),
	)
	if err != nil {
		t.Fatal(err)
	}

}

func Test_suppressBundleRelocation(t *testing.T) {
	pkg := createPkg(t)
	if err := suppressBundleRelocation(pkg); err != nil {
		t.Fatal(err)
	}
}

// create package in tmp, returning path
func createPkg(t *testing.T) string {
	pkgroot := filepath.Join(os.TempDir(), "pkgroot")
	tp := filepath.Join(pkgroot, "/tmp/")
	os.MkdirAll(tp, 0755)
	if err := ioutil.WriteFile(filepath.Join(tp, "testfile"), []byte("testfile"), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(pkgroot)

	path := filepath.Join(os.TempDir(), "mackit-test-package.pkg")
	cmd := exec.Command(
		"/usr/bin/pkgbuild",
		"--root", pkgroot,
		"--identifier", "test.pkg",
		"--version", "1.2.3",
		path,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		t.Fatal(err)
	}
	return path
}

func isRoot(t *testing.T) bool {
	u, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	return u.Uid == "0"
}
