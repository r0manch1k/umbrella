package keypair

func (uc *UseCase) GetPublicKey() ([]byte, error) {
	publicKey, err := uc.service.GetPublicKey()

	return publicKey, err
}
