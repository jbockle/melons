package ioc

type ServiceLifetime uint8

const (
	Singleton ServiceLifetime = iota
	Transient
)

var lifetimeString [2]string = [2]string{
	"Singleton",
	"Transient",
}

func (sl ServiceLifetime) String() string {
	return lifetimeString[sl]
}
