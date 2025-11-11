package keypair

func (uc *UseCase) GenerateAndSaveKeyPair() error {
	if uc.service.IsExistsPublicKey() && uc.service.IsExistsPrivateKey() {
		return nil
	}

	privateKeyObj, privatePEM, err := uc.service.GeneratePrivateKey()
	if err != nil {
		return err
	}

	publicPEM, err := uc.service.GeneratePublicKey(privateKeyObj)
	if err != nil {
		return err
	}

	err = uc.service.SavePrivateKey(privatePEM)
	if err != nil {
		return err
	}

	err = uc.service.SavePublicKey(publicPEM)
	if err != nil {
		return err
	}

	return nil
}
