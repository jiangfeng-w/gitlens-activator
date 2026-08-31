package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const v18Snippet = `function i9(e,t,i){let r,s,n,o={id:e.user.id,name:e.user.name,email:e.user.email};return o}
changeSubscription(e,t,i){e??={plan:{actual:(0,ty.le)("community",!1,0,void 0),effective:(0,ty.le)("community",!1,0,void 0)},account:void 0,state:iJ.z.Community},(0,ty.Jc)(e)}`

const v19Snippet = `function iL(e,t,i){let r,s,n,o={id:e.user.id,name:e.user.name,email:e.user.email};return o}
changeSubscription(e,t,i){e??={plan:{actual:(0,tw.le)("community",!1,0,void 0),effective:(0,tw.le)("community",!1,0,void 0)},account:void 0,state:iS.z.Community},(0,tw.Jc)(e)}`

func TestActivateForVersion16AppliesBothPatches(t *testing.T) {
	patched, changed, err := activateForVersion16(v18Snippet)
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("expected changed")
	}
	if !strings.Contains(patched, `licenses:{paidLicenses:{},effectiveLicenses:{"gitlens-pro"`) {
		t.Fatal("license injection not applied")
	}
	if !strings.Contains(patched, `e={plan:{actual:(0,ty.le)("pro",!1,0,void 0)`) {
		t.Fatal("no-login patch not applied")
	}
	if !strings.Contains(patched, `account:{id:"88888888-8888-8888-8888-888888888888"`) {
		t.Fatal("mock account not injected")
	}

	// 幂等：再次激活不应有任何改动
	again, changed2, err := activateForVersion16(patched)
	if err != nil {
		t.Fatal(err)
	}
	if changed2 || again != patched {
		t.Fatal("expected idempotent")
	}
}

func TestActivateForVersion16V19Module(t *testing.T) {
	patched, changed, err := activateForVersion16(v19Snippet)
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("expected changed")
	}
	if !strings.Contains(patched, `e={plan:{actual:(0,tw.le)("pro",!1,0,void 0)`) {
		t.Fatal("v19 no-login patch not applied")
	}
}

func TestActivateForVersion16NoLicenseInjectionPoint(t *testing.T) {
	// 只有 changeSubscription 注入点，没有 let xxx={id:e.user.id... 结构
	content := `changeSubscription(e,t,i){e??={plan:{actual:(0,ty.le)("community",!1,0,void 0),effective:(0,ty.le)("community",!1,0,void 0)},account:void 0,state:iJ.z.Community}}`
	patched, changed, err := activateForVersion16(content)
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("expected no-login patch to apply")
	}
	if !hasNoLoginPatch(patched) {
		t.Fatal("no-login patch marker missing")
	}
	if hasLicenseInjection(patched) {
		t.Fatal("license injection should not be present")
	}
}

func TestActivateForVersion16NoPatchPoints(t *testing.T) {
	content := `function foo(){return 1}`
	_, changed, err := activateForVersion16(content)
	if err != nil {
		t.Fatal(err)
	}
	if changed {
		t.Fatal("expected unchanged")
	}
}

const v16Snippet = `function i9(e,t,i){let r,s,n,o={id:e.user.id,name:e.user.name,email:e.user.email};return o}`

func TestActivateRestoreDir16Idempotent(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "eamodio.gitlens-16.10.0")
	if err := os.MkdirAll(filepath.Join(dir, "dist"), 0o755); err != nil {
		t.Fatal(err)
	}
	js := filepath.Join(dir, "dist", "gitlens.js")
	if err := os.WriteFile(js, []byte(v16Snippet), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := activateExtensionDir(dir); err != nil {
		t.Fatal("activate:", err)
	}
	// 16.x 只有 license 注入点，重复激活必须幂等成功
	if err := activateExtensionDir(dir); err != nil {
		t.Fatal("second activate should be idempotent:", err)
	}

	if err := restoreExtensionDir(dir); err != nil {
		t.Fatal("restore:", err)
	}
	restored, err := os.ReadFile(js)
	if err != nil {
		t.Fatal(err)
	}
	if string(restored) != v16Snippet {
		t.Fatal("restored content mismatch")
	}
}

func TestActivateRestoreDir18(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "eamodio.gitlens-18.3.0")
	if err := os.MkdirAll(filepath.Join(dir, "dist"), 0o755); err != nil {
		t.Fatal(err)
	}
	js := filepath.Join(dir, "dist", "gitlens.js")
	if err := os.WriteFile(js, []byte(v18Snippet), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := activateExtensionDir(dir); err != nil {
		t.Fatal("activate:", err)
	}
	data, err := os.ReadFile(js)
	if err != nil {
		t.Fatal(err)
	}
	if !hasNoLoginPatch(string(data)) {
		t.Fatal("no-login patch missing after activate")
	}
	if _, err := os.Stat(js + ".backup"); err != nil {
		t.Fatal("backup not created:", err)
	}

	// 幂等激活
	if err := activateExtensionDir(dir); err != nil {
		t.Fatal("second activate:", err)
	}

	if err := restoreExtensionDir(dir); err != nil {
		t.Fatal("restore:", err)
	}
	restored, err := os.ReadFile(js)
	if err != nil {
		t.Fatal(err)
	}
	if string(restored) != v18Snippet {
		t.Fatal("restored content mismatch")
	}
}
