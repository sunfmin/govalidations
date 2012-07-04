package govalidations

type GateKeeper struct {
	Validators []Validator
}

func NewGateKeeper() *GateKeeper {
	return &GateKeeper{}
}

func (gk *GateKeeper) Add(vd Validator) (r *GateKeeper) {
	gk.Validators = append(gk.Validators, vd)
	return gk
}

func (gk *GateKeeper) Join(newgk *GateKeeper) (r *GateKeeper) {
	gk.Validators = append(gk.Validators, newgk.Validators...)
	return gk
}

func (gk *GateKeeper) Validate(object interface{}) (r *Validated) {
	r = &Validated{}
	r.Object = object
	for _, vd := range gk.Validators {
		r.Errors = append(r.Errors, vd(object)...)
	}
	return
}
