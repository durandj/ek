package configuration

import "fmt"

type HTTPConfiguration struct {
	Interface string `default:"0.0.0.0"`
	Port      int    `default:"8080"`
}

func (httpConfig HTTPConfiguration) Address() string {
	return fmt.Sprintf("%s:%d", httpConfig.Interface, httpConfig.Port)
}
