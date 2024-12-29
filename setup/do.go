package setup

func Do() error {
	err := Commands()
	if err != nil {
		return err
	}

	err = Log()
	if err != nil {
		return err
	}

	err = Aumid()
	if err != nil {
		return err
	}

	return nil
}
