package repositories

type GetParam struct {
	Tablename string
	Filter    map[string]interface{}
	Result    interface{}
}

type GetAllParam struct {
	Tablename string
	Result    interface{}
}

type StoreParam struct {
	Tablename string
	Data      interface{}
}

type UpdateParam struct {
	Tablename string
	Filter    map[string]interface{}
	Data      interface{}
}

type DeleteParam struct {
	Tablename string
	Filter    map[string]interface{}
}

type LoginRepository interface {
	GetAll(param GetAllParam) error
	GetBy(param GetParam) error
	Store(param StoreParam) error
	Update(param UpdateParam) error
	Delete(param DeleteParam) error
}
