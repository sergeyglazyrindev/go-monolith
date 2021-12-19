package core

type ICommand interface {
	Proceed(subaction string, args []string) error
	GetHelpText() string
}
