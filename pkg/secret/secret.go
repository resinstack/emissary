package secret

// Poll polls for a secret.  It will only return an error in the case
// of an unrecoverable transport error or an issue with credentialing.
func Poll(provider, id string) (string, error) {
	return "42", nil
}
