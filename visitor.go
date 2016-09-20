package filter

type Visitor interface {
	Visit(node Node)
}

type SQLVisitor struct {
	sql string
}

func (v *SQLVisitor) Visit(node Node) {

}

func (v *SQLVisitor) Sql() string {
	return v.sql
}
