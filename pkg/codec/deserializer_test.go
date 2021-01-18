package codec_test

import (
	"encoding/json"
	"reflect"

	"github.com/dbsystel/kewl/pkg/panicutils"

	"github.com/dbsystel/kewl/pkg/codec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ codec.SchemeExtension = &corev1Extension{}

type corev1Extension struct{}

func (c *corev1Extension) AddToScheme(scheme *runtime.Scheme) error {
	return corev1.AddToScheme(scheme)
}

func gvk(gv schema.GroupVersion, obj interface{}) schema.GroupVersionKind {
	return schema.GroupVersionKind{Group: gv.Group, Version: gv.Version, Kind: reflect.TypeOf(obj).Name()}
}

func marshalJSON(obj runtime.Object) []byte {
	bytes, err := json.Marshal(obj)
	panicutils.PanicIfError(err)
	return bytes
}

var _ = Describe("Deserializer", func() {
	var sut codec.Deserializer
	BeforeEach(func() {
		sut = codec.NewDeserializer(runtime.NewScheme())
	})
	It("should throw an error if object is unknown", func() {
		deserialize, err := sut.Deserialize(schema.GroupVersionKind{Group: "meh", Version: "meh", Kind: "meh"}, nil)
		Expect(deserialize).To(BeNil())
		Expect(err).To(HaveOccurred())
	})
	It("should throw an error if the object is known, but cannot be deserialized", func() {
		Expect(sut.Register(&corev1Extension{})).To(Not(HaveOccurred()))
		deserialize, err := sut.Deserialize(gvk(corev1.SchemeGroupVersion, corev1.Pod{}), []byte("meh"))
		Expect(deserialize).To(BeNil())
		Expect(err).To(HaveOccurred())
	})
	It("should throw an error when the scheme extension is nil", func() {
		Expect(sut.Register(nil)).To(HaveOccurred())
	})
	It("should deserialize an object correctly", func() {
		Expect(sut.Register(&corev1Extension{})).To(Not(HaveOccurred()))
		Expect(sut.Deserialize(gvk(corev1.SchemeGroupVersion, corev1.Pod{}), marshalJSON(&corev1.Pod{}))).Should(Not(BeNil()))
	})
	It("should propagate not registered correctly", func() {
		Expect(sut.Register(&corev1Extension{})).To(Not(HaveOccurred()))
		_, err := sut.Deserialize(schema.GroupVersionKind{}, marshalJSON(&corev1.Pod{}))
		Expect(err).To(HaveOccurred())
		Expect(runtime.IsNotRegisteredError(err)).To(BeTrue())
	})
})
