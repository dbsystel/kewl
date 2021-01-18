package validation

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ ResultCollector = &resultCollectorImpl{}

// statusCauseSliceRef is a helper struct for storing the actual pointer to a slice of strings
type statusCauseSliceRef struct {
	causes []v1.StatusCause
}

// resultCollectorImpl is the default implementation of ResultCollector
type resultCollectorImpl struct {
	origin      string
	messagesRef *statusCauseSliceRef
}

func (r *resultCollectorImpl) AppendField(suffix string) ResultCollector {
	return &resultCollectorImpl{origin: r.origin + suffix, messagesRef: r.messagesRef}
}

func (r *resultCollectorImpl) Failures() []v1.StatusCause {
	return r.messagesRef.causes
}

func (r *resultCollectorImpl) AddFailure(failure string) {
	r.messagesRef.causes = append(
		r.messagesRef.causes,
		v1.StatusCause{Type: "invalid", Message: failure, Field: r.origin},
	)
}
