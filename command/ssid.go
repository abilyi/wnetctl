package command

func GetSsidCommand(argv []string) Command {
	return nil
}

type SsidCommand struct {
	GenericCommand
	name string
}
