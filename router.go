package hemera

type PatternSet struct {
	fields  []string
	payload interface{}
}

type Router struct {
	items []PatternSet
}

func NewRouter() Router {
	items := make([]PatternSet, 100)
	return Router{items: items}
}

func Add(r *Router, set PatternSet) error {
	r.items = append(r.items, set)
	return nil
}
