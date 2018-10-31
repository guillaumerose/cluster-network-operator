package operator

import (
	"testing"

	. "github.com/onsi/gomega"
	uns "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// TestReconcileNamespace makes sure that namespace
// annotations are merged, and everything else is overwritten
// Namespaces use the "generic" logic; deployments and services
// have custom logic
func TestMergeNamespace(t *testing.T) {
	g := NewGomegaWithT(t)

	cur := UnstructuredFromYaml(t, `
apiVersion: v1
kind: Namespace
metadata:
  name: ns1
  labels:
    a: cur
    b: cur
  annotations:
    a: cur
    b: cur`)

	upd := UnstructuredFromYaml(t, `
apiVersion: v1
kind: Namespace
metadata:
  name: ns1
  labels:
    a: upd
    c: upd
  annotations:
    a: upd
    c: upd`)

	// this mutates updated
	err := MergeObjectForUpdate(cur, upd)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(upd.GetLabels()).To(Equal(map[string]string{
		"a": "upd",
		"b": "cur",
		"c": "upd",
	}))

	g.Expect(upd.GetAnnotations()).To(Equal(map[string]string{
		"a": "upd",
		"b": "cur",
		"c": "upd",
	}))
}

func TestMergeDeployment(t *testing.T) {
	g := NewGomegaWithT(t)

	cur := UnstructuredFromYaml(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: d1
  labels:
    a: cur
    b: cur
  annotations:
    deployment.kubernetes.io/revision: cur
    a: cur
    b: cur`)

	upd := UnstructuredFromYaml(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: d1
  labels:
    a: upd
    c: upd
  annotations:
    deployment.kubernetes.io/revision: upd
    a: upd
    c: upd`)

	// this mutates updated
	err := MergeObjectForUpdate(cur, upd)
	g.Expect(err).NotTo(HaveOccurred())

	// labels are not merged
	g.Expect(upd.GetLabels()).To(Equal(map[string]string{
		"a": "upd",
		"b": "cur",
		"c": "upd",
	}))

	// annotations are merged
	g.Expect(upd.GetAnnotations()).To(Equal(map[string]string{
		"a": "upd",
		"b": "cur",
		"c": "upd",
		"deployment.kubernetes.io/revision": "cur",
	}))
}

func TestMergeNilCur(t *testing.T) {
	g := NewGomegaWithT(t)

	cur := UnstructuredFromYaml(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: d1`)

	upd := UnstructuredFromYaml(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: d1
  labels:
    a: upd
    c: upd
  annotations:
    a: upd
    c: upd`)

	// this mutates updated
	err := MergeObjectForUpdate(cur, upd)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(upd.GetLabels()).To(Equal(map[string]string{
		"a": "upd",
		"c": "upd",
	}))

	g.Expect(upd.GetAnnotations()).To(Equal(map[string]string{
		"a": "upd",
		"c": "upd",
	}))
}

func TestMergeNilMeta(t *testing.T) {
	g := NewGomegaWithT(t)

	cur := UnstructuredFromYaml(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: d1`)

	upd := UnstructuredFromYaml(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: d1`)

	// this mutates updated
	err := MergeObjectForUpdate(cur, upd)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(upd.GetLabels()).To(BeEmpty())
}

func TestMergeNilUpd(t *testing.T) {
	g := NewGomegaWithT(t)

	cur := UnstructuredFromYaml(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: d1
  labels:
    a: cur
    b: cur
  annotations:
    a: cur
    b: cur`)

	upd := UnstructuredFromYaml(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: d1`)

	// this mutates updated
	err := MergeObjectForUpdate(cur, upd)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(upd.GetLabels()).To(Equal(map[string]string{
		"a": "cur",
		"b": "cur",
	}))

	g.Expect(upd.GetAnnotations()).To(Equal(map[string]string{
		"a": "cur",
		"b": "cur",
	}))
}

func TestMergeService(t *testing.T) {
	g := NewGomegaWithT(t)

	cur := UnstructuredFromYaml(t, `
apiVersion: v1
kind: Service
metadata:
  name: d1
spec:
  clusterIP: cur`)

	upd := UnstructuredFromYaml(t, `
apiVersion: v1
kind: Service
metadata:
  name: d1
spec:
  clusterIP: upd`)

	err := MergeObjectForUpdate(cur, upd)
	g.Expect(err).NotTo(HaveOccurred())

	ip, _, err := uns.NestedString(upd.Object, "spec", "clusterIP")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(ip).To(Equal("cur"))
}