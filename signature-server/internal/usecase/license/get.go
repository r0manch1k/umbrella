package license

func (uc *UseCase) GetPublicKey() ([]byte, error) {
	return uc.signer.GetPublicKey()
}
