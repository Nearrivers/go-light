package app_errors

type HueApiError struct {
	Description string `json:"description"` // Description de l'erreur lisible par un humain
}

func (hae HueApiError) Error() string {
	return hae.Description
}
