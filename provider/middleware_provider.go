package provider

type MiddlewareProvider interface {
}

type middlewareProvider struct {
}

func NewMiddlewareProvider(servicesProvider ServicesProvider) MiddlewareProvider {
	return &middlewareProvider{}
}
