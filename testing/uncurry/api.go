package uncurry

// Error2 extracts the second argument as error
func Error2(_ interface{}, err error) error {
	return err
}
