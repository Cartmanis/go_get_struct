package generator

import (
	"fmt"
	"github.com/cartmanis/go_get_struct/engine"
	"github.com/cartmanis/go_get_struct/node"
	"os"
	"path/filepath"
	"strings"
)

const (
	a = string('\u0060')
)

var mapList = make(map[string]bool, 0)

func CreateStruct(n *node.NodeXml, xmlPath string) error {
	if len(xmlPath) < 3 {
		return fmt.Errorf("Не достаточная длинна файла %v. Возможно файл без расширения", filepath.Base(xmlPath))
	}
	goPath := xmlPath[:len(xmlPath)-3] + "go"
	s := goNode(n)
	fmt.Println(s)
	goFile, err := os.Create(goPath)
	if err != nil {
		return err
	}
	defer goFile.Close()
	_, err = goFile.WriteString(s)
	if err != nil {
		return err
	}
	fmt.Printf("Структура go успешно создана и записана в файл %v", goPath)
	return nil
}

func goNode(n *node.NodeXml) string {
	if n == nil { //|| n.Childern == nil || len(n.Childern) <= 0 {
		return ""
	}
	if mapList[n.Namespace] {
		return ""
	}
	//if n.Namespace == "city" {
	//	fmt.Println(n)
	//}
	mapList[n.Namespace] = true
	temp := getNameStruct(n)
	temp += goAttribute(n)
	if n.Childern != nil && len(n.Childern) >= 0 {
		temp += goChildern(n.Childern)
	}

	for _, v := range n.Childern {
		temp += goNode(v)
	}
	return temp
}

func getNameStruct(n *node.NodeXml) string {
	if (n.Childern == nil || len(n.Childern) <= 0) && (n.Attr == nil || len(n.Attr) <= 0) {
		return ""
	}
	return fmt.Sprintf("\n"+`type %v struct {`, engine.GetCamelCase(n.Namespace)) + "\n\t" +
		fmt.Sprintf(`XMLName xml.Name %vxml:"%v"%v`+"\n", a, n.Namespace, a)
}

func goAttribute(n *node.NodeXml) string {
	if n == nil || n.Attr == nil || len(n.Attr) <= 0 {
		return ""
	}
	var temp string
	for _, v := range n.Attr {
		name := fmt.Sprintf("Attr%v", engine.GetCamelCase(v.Name.Local))
		temp += fmt.Sprintf("\t"+`%v string %vxml:"%v,attr"%v`+"\n", name, a, v.Name.Local, a)
	}
	if n.Childern == nil || len(n.Childern) <= 0 {
		temp += "}\n"
	}
	return temp
}

func goChildern(listChild []*node.NodeXml) string {
	var temp string
	mapList := make(map[string]bool)
	for _, v := range listChild {
		if mapList[v.Namespace] {
			continue
		}
		name := engine.GetCamelCase(v.Namespace)
		t := getType(v, isArray(v, listChild))
		temp += fmt.Sprintf("\t"+`%v %v %vxml:"%v"%v`+"\n",
			name, t, a, v.Namespace, a)
		mapList[v.Namespace] = true
	}
	temp += "}" + "\n"
	return temp
}

func isArray(current *node.NodeXml, listChild []*node.NodeXml) bool {
	if listChild == nil || len(listChild) <= 0 || current == nil {
		return false
	}
	var count int
	mapList := make(map[string]bool, 0)
	mapList[current.Namespace] = true
	for _, v := range listChild {
		if mapList[v.Namespace] {
			count++
		}
		if count > 1 {
			return true
		}
	}
	return false
}

func getType(n *node.NodeXml, isArray bool) string {
	if n.Childern == nil && isArray {
		return "[]string"
	}
	if n.Childern == nil && !isArray && (n.Attr == nil || len(n.Attr) <= 0) {
		return "string"
	}
	if isArray {
		return "[]*" + strings.Title(n.Namespace)
	}
	return "*" + strings.Title(n.Namespace)
}
