package translator

type Translator interface {
	Translate(contentTH string) (string, error)
}
