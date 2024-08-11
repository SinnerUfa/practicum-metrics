package config

func Load(v any, args []string) error {
	if err := LoadFlags(v, args); err != nil {
		return err
	}
	if err := LoadEnv(v); err != nil {
		return err
	}
	return nil
}
