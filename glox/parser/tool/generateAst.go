package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Types from other packages must be added here
// Format is of
// Type : package
// Importing of package must be done manually
var (
	otherPackageTypes map[string]string = map[string]string{
		"Token":  "token",
		"Object": "token",
	}
)

func main() {
	generateExpr()
	generateStmt()
}

func generateStmt() {
	filename := "stmt.go"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	paths := strings.Split(dir, "/")
	parserPath := strings.Join(paths[:len(paths)-1], "/")

	// Remove file if present and make a new one
	_ = os.Remove(parserPath + "/" + filename)
	file, err := os.OpenFile(parserPath+"/"+filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	defineAst(file, "Stmt", []string{
		"Block : []Stmt statements",
		"Expression : Expr expression",
		"Function : Token name, []Token params, []Stmt body",
		"If : Expr condition, Stmt thenBranch, Stmt elseBranch",
		"Print : Expr expression",
		"Var : Token name, Expr initializer",
		"While : Expr condition, Stmt body",
	})
}

func generateExpr() {
	filename := "expr.go"
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	paths := strings.Split(dir, "/")
	parserPath := strings.Join(paths[:len(paths)-1], "/")

	// Remove file if present and make a new one
	_ = os.Remove(parserPath + "/" + filename)
	file, err := os.OpenFile(parserPath+"/"+filename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	defineAst(file, "Expr", []string{
		"Assign : Token name, Expr value",
		"Binary : Expr left, Token operator, Expr right",
		"Call : Expr callee, Token paren, []Expr arguments",
		"Grouping : Expr expression",
		"Literal : Object value",
		"Logical: Expr left, Token operator, Expr right",
		"Unary : Token operator, Expr right",
		"Variable : Token name",
	})
}

func defineAst(file *os.File, basename string, types []string) {
	builder := strings.Builder{}
	_, err := builder.WriteString(`
package parser
import (
	"github.com/madraceee/interpreters/glox/token"
)`)
	if err != nil {
		log.Fatal(err)
	}

	// Add interface for all rules to implement
	// Allows all rules to return a  string of what they hold
	builder.WriteString("\ntype  " + basename + " interface{\n")
	builder.WriteString("\tVisit(Visit" + basename + ") (interface{}, error)\n}\n\n")

	_, err = file.WriteString(builder.String())
	if err != nil {
		log.Fatal(err)
	}

	defineToBeImplementedInterface(file, basename, types)

	for _, _type := range types {
		split := strings.Split(_type, ":")
		classname := strings.Trim(split[0], " ")
		fields := strings.Trim(split[1], " ")
		defineTypes(file, basename, classname, fields)
	}
}

func defineTypes(file *os.File, basename, classname, fieldlist string) {
	builder := strings.Builder{}

	fieldsWithType := make([]string, 0)
	fields := make([]string, 0)

	unformatedFields := strings.Split(fieldlist, ",")
	for _, field := range unformatedFields {
		vals := strings.Split(strings.Trim(field, " "), " ")
		fieldsWithType = append(fieldsWithType, vals[1]+" "+replaceOtherPackageTypes(vals[0]))
		fields = append(fields, vals[1])
	}

	// Struct for the class along with the fields
	_, err := builder.WriteString("type " + classname + " struct {\n")
	if err != nil {
		log.Fatal(err)
	}
	for _, field := range fieldsWithType {
		builder.WriteString(toUpperFirstChar(field) + "\n")
	}
	_, err = builder.WriteString("\n}\n")
	if err != nil {
		log.Fatal(err)
	}

	// Create a builder function for the rule
	builder.WriteString("func New" + classname + "(")
	for i, field := range fieldsWithType {
		if i != 0 {
			builder.WriteString(",")
		}
		builder.WriteString(field)
	}
	builder.WriteString(")" + basename + "{\n")
	builder.WriteString("\t return &" + classname + "{\n")
	for _, field := range fields {
		builder.WriteString(toUpperFirstChar(field) + " : " + field + ",\n")
	}
	builder.WriteString("\t}\n")
	builder.WriteString("}\n")

	_, err = file.WriteString(builder.String())
	if err != nil {
		log.Fatal(err)
	}

	defineVisitFunc(file, basename, classname)
}

func defineVisitFunc(file *os.File, basename, classname string) {
	builder := strings.Builder{}
	builder.WriteString("\nfunc (expr *" + classname + ") Visit(visitor Visit" + basename + ") (interface{}, error){\n")
	builder.WriteString("\treturn visitor.Visit" + classname + basename + "(expr)\n")
	builder.WriteString("}\n")

	_, err := file.WriteString(builder.String())
	if err != nil {
		log.Fatal(err)
	}
}

func defineToBeImplementedInterface(file *os.File, basename string, types []string) {
	builder := strings.Builder{}
	builder.WriteString("\ntype  Visit" + basename + " interface{\n")
	for _, _type := range types {
		split := strings.Split(_type, ":")
		classname := strings.Trim(split[0], " ")
		builder.WriteString("\tVisit" + classname + basename + "(*" + classname + ") (interface{}, error)\n")
	}
	builder.WriteString("}\n\n")

	_, err := file.WriteString(builder.String())
	if err != nil {
		log.Fatal(err)
	}
}

func replaceOtherPackageTypes(otherPackageType string) string {
	isArr := false
	if otherPackageType[0:2] == "[]" {
		isArr = true
		otherPackageType = otherPackageType[2:]
	}
	if packageName, ok := otherPackageTypes[otherPackageType]; ok {
		if isArr {
			packageName = "[]" + packageName
		}
		return packageName + "." + otherPackageType
	}
	if isArr {
		otherPackageType = "[]" + otherPackageType
	}
	return otherPackageType
}

func toUpperFirstChar(str string) string {
	runes := []rune(str)
	if runes[0] >= 97 {
		runes[0] = runes[0] - 32
	}

	return string(runes)
}
