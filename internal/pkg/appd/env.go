package appd

type Env map[string]interface{}

func (env *Env) Get(k string) interface{} {
	if env == nil {
		return nil
	}
	return (*env)[k]
}

func (env *Env) Set(k string, v interface{}) {
	if env == nil {
		return
	}
	(*env)[k] = v
}

func (env *Env) Clean() {
	if env == nil {
		return
	}
	*env = make(Env)
}

func (env *Env) Del(k string) {
	if env == nil {
		return
	}
	delete(*env, k)
}
