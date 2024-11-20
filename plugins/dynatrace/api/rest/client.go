package rest

type Client interface {
	Get(url string, expectedStatusCodes ...int) Request
}
