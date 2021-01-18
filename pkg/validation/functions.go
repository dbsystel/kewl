package validation

// EnsureNonNil adds a validation failure to the validation.ResultCollector, in case the value is nil
func EnsureNonNil(value interface{}, results ResultCollector) {
	if value == nil {
		results.AddFailure("not set")
	}
}
