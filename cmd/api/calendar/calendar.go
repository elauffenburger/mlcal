package calendar

import "github.com/elauffenburger/musical-literacy-cal/pkg/mlcal"

type Getter interface {
	Get() (*mlcal.Calendar, error)
	GetICS() (string, error)
}

type Setter interface {
	Set(*mlcal.Calendar) error
}

type Cache interface {
	Getter
	Setter
}
