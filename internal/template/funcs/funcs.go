package funcs

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"path/filepath"
	"reflect"
	"sort"
	"github.com/jiangjiax/stars/internal/asset"
	"strings"
	"time"
)

// 创建函数映射组
var (
	stringFuncs = template.FuncMap{
		"lower":      strings.ToLower,
		"upper":      strings.ToUpper,
		"title":      strings.Title,
		"urlize":     urlize,
		"truncate":   truncate,
		"trimPrefix": strings.TrimPrefix,
		"jsonify":    jsonify,
		"AssetPath":  getAssetPath,
		"urlquery":   url.QueryEscape,
		"hasPrefix":  strings.HasPrefix,
		"hasSuffix":  strings.HasSuffix,
		"contains":   strings.Contains,
		"replace":    strings.Replace,
		"safeurl": func(s string) template.URL {
			return template.URL(s)
		},
		"pathEscape": func(s string) string {
			return url.PathEscape(s)
		},
	}

	dateFuncs = template.FuncMap{
		"formatDate": formatDate,
		"dateFormat": dateFormat,
		"now":        time.Now,
		"toDate":     toDate,
	}

	comparisonFuncs = template.FuncMap{
		"eq": reflect.DeepEqual,
		"ne": ne,
		"lt": lt,
		"gt": gt,
	}

	filterFuncs = template.FuncMap{
		"where": where,
		"sort":  sortBy,
	}

	sliceFuncs = template.FuncMap{
		"slice":   slice,
		"append":  appendSlice,
		"uniq":    uniq,
		"delimit": delimit,
		"first":   first,
		"reverse": reverse,
	}

	// HTML 相关的函数
	htmlFuncs = template.FuncMap{
		"safeHTML": safeHTML,
	}

	// 添加 len 函数映射
	utilFuncs = template.FuncMap{
		"len": length,
	}

	// 添加数学运算函数映射组
	mathFuncs = template.FuncMap{
		"add":      add,
		"sub":      sub,
		"div":      div,
		"mul":      mul,
		"seq":      seq,
		"ceil":     ceil,
		"sequence": sequence,
	}

	// 添加一些有用的路径处理函数
	pathFuncs = template.FuncMap{
		"basename": filepath.Base,
		"dirname":  filepath.Dir,
		"ext":      filepath.Ext,
	}

	// 合并所有函数映射
	DefaultFuncs = mergeFuncMaps(
		stringFuncs,
		dateFuncs,
		comparisonFuncs,
		filterFuncs,
		sliceFuncs,
		htmlFuncs,
		utilFuncs,
		mathFuncs,
		pathFuncs,
	)

	assetPipeline *asset.Pipeline
)

// 字符串处理函数
func urlize(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// 日期处理函数
func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func dateFormat(layout string, t time.Time) string {
	return t.Format(layout)
}

func toDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

// 比较函数
func ne(a, b interface{}) bool {
	return !reflect.DeepEqual(a, b)
}

func lt(a, b interface{}) bool {
	switch v := a.(type) {
	case int:
		if bv, ok := b.(int); ok {
			return v < bv
		}
	case string:
		if bv, ok := b.(string); ok {
			return v < bv
		}
	case time.Time:
		if bv, ok := b.(time.Time); ok {
			return v.Before(bv)
		}
	}
	return false
}

func gt(a, b interface{}) bool {
	switch v := a.(type) {
	case int:
		if bv, ok := b.(int); ok {
			return v > bv
		}
	case string:
		if bv, ok := b.(string); ok {
			return v > bv
		}
	case time.Time:
		if bv, ok := b.(time.Time); ok {
			return v.After(bv)
		}
	}
	return false
}

// 合并多个函数映射
func mergeFuncMaps(maps ...template.FuncMap) template.FuncMap {
	merged := make(template.FuncMap)
	for _, m := range maps {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
}

// GetFunc 获取指定的模板函数
func GetFunc(name string) interface{} {
	return DefaultFuncs[name]
}

// AddFunc 添加自定义模板函数
func AddFunc(name string, fn interface{}) {
	DefaultFuncs[name] = fn
}

// where 函数用于过滤切片
func where(arr interface{}, key string, value interface{}) (interface{}, error) {
	// 获取数组的反射值
	arrValue := reflect.ValueOf(arr)
	if arrValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("first argument must be a slice")
	}

	// 创建结果切片
	result := reflect.MakeSlice(arrValue.Type(), 0, arrValue.Len())

	// 遍历切片
	for i := 0; i < arrValue.Len(); i++ {
		item := arrValue.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}

		// 获取字段值
		field := item.FieldByName(key)
		if !field.IsValid() {
			continue
		}

		// 比较值
		if reflect.DeepEqual(field.Interface(), value) {
			result = reflect.Append(result, arrValue.Index(i))
		}
	}

	return result.Interface(), nil
}

// sortBy 函数用于排序切片
func sortBy(arr interface{}, key string, order string) (interface{}, error) {
	// 获取数组的反射值
	arrValue := reflect.ValueOf(arr)
	if arrValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("first argument must be a slice")
	}

	// 创建一个可以排序的切片
	length := arrValue.Len()
	sorted := make([]reflect.Value, length)
	for i := 0; i < length; i++ {
		sorted[i] = arrValue.Index(i)
	}

	// 定义排序函数
	less := func(i, j int) bool {
		// 获取字段值
		a := sorted[i]
		b := sorted[j]

		// 处理接口类型
		if a.Kind() == reflect.Interface {
			a = a.Elem()
		}
		if b.Kind() == reflect.Interface {
			b = b.Elem()
		}

		// 处理指针类型
		if a.Kind() == reflect.Ptr {
			a = a.Elem()
		}
		if b.Kind() == reflect.Ptr {
			b = b.Elem()
		}

		fieldA := a.FieldByName(key)
		fieldB := b.FieldByName(key)

		if !fieldA.IsValid() || !fieldB.IsValid() {
			return false
		}

		// 根据字段类型进行比较
		switch fieldA.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if order == "desc" {
				return fieldA.Int() > fieldB.Int()
			}
			return fieldA.Int() < fieldB.Int()
		case reflect.String:
			if order == "desc" {
				return fieldA.String() > fieldB.String()
			}
			return fieldA.String() < fieldB.String()
			// 可以添加其他类型的比较
		}
		return false
	}

	// 排序
	sort.SliceStable(sorted, less)

	// 创建结果切片
	result := reflect.MakeSlice(arrValue.Type(), length, length)
	for i := 0; i < length; i++ {
		result.Index(i).Set(sorted[i])
	}

	return result.Interface(), nil
}

// uniq 函数用于去重切片
func uniq(slice interface{}) (interface{}, error) {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("uniq: input must be a slice")
	}

	seen := make(map[interface{}]bool)
	result := reflect.MakeSlice(v.Type(), 0, v.Len())

	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()
		if !seen[item] {
			seen[item] = true
			result = reflect.Append(result, v.Index(i))
		}
	}

	return result.Interface(), nil
}

// delimit 函数用于连接切片
func delimit(slice interface{}, delimiter string) (string, error) {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return "", fmt.Errorf("delimit: input must be a slice")
	}

	var items []string
	for i := 0; i < v.Len(); i++ {
		items = append(items, fmt.Sprint(v.Index(i).Interface()))
	}

	return strings.Join(items, delimiter), nil
}

// slice 函数用于创建新的片
func slice() interface{} {
	return make([]interface{}, 0)
}

// appendSlice 函数用于向切片添加元素
func appendSlice(slice interface{}, elements ...interface{}) (interface{}, error) {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("first argument must be a slice")
	}

	result := reflect.MakeSlice(v.Type(), v.Len(), v.Len()+len(elements))
	reflect.Copy(result, v)

	for _, elem := range elements {
		result = reflect.Append(result, reflect.ValueOf(elem))
	}

	return result.Interface(), nil
}

// safeHTML 函数用于安全地渲染 HTML 内容
func safeHTML(s interface{}) template.HTML {
	switch v := s.(type) {
	case string:
		return template.HTML(v)
	case template.HTML:
		return v
	default:
		return template.HTML(fmt.Sprint(v))
	}
}

// length 函数用于获取切片、数组、字符串或映射的长度
func length(v interface{}) int {
	if v == nil {
		return 0
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Array, reflect.Slice, reflect.String, reflect.Map:
		return rv.Len()
	default:
		return 0
	}
}

// add 函数用于数字相加
func add(a, b int) int {
	return a + b
}

// sub 函数用于数字相减
func sub(a, b int) int {
	return a - b
}

// first 函数用于获取切片的前 n 个元素
func first(n int, items interface{}) (interface{}, error) {
	v := reflect.ValueOf(items)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("first: argument must be a slice")
	}

	if n > v.Len() {
		n = v.Len()
	}

	result := reflect.MakeSlice(v.Type(), n, n)
	for i := 0; i < n; i++ {
		result.Index(i).Set(v.Index(i))
	}

	return result.Interface(), nil
}

// reverse 函数用于反转切片
func reverse(items interface{}) (interface{}, error) {
	v := reflect.ValueOf(items)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("reverse: argument must be a slice")
	}

	length := v.Len()
	result := reflect.MakeSlice(v.Type(), length, length)

	for i := 0; i < length; i++ {
		result.Index(i).Set(v.Index(length - 1 - i))
	}

	return result.Interface(), nil
}

// jsonify 将数据转换为 JSON 字符串
func jsonify(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return string(b)
}

// 设置资源管理器实例
func SetAssetPipeline(p *asset.Pipeline) {
	assetPipeline = p
}

// 获取资源文件的哈希路径
func getAssetPath(name string) string {
	if assetPipeline == nil {
		return name
	}
	return assetPipeline.GetAssetPath(name)
}

// 数学运算函数
func div(a, b int) float64 {
	if b == 0 {
		return 0
	}
	return float64(a) / float64(b)
}

func mul(a, b int) int {
	return a * b
}

// seq 生成一个从1到n的序列
func seq(n int) []int {
	if n <= 0 {
		return []int{}
	}
	result := make([]int, n)
	for i := range result {
		result[i] = i + 1
	}
	return result
}

// 添加 ceil 函数
func ceil(x float64) int {
	return int((x + 0.99999999))
}

// sequence 函数生成从0到n-1的序列
func sequence(n int) []int {
	seq := make([]int, n)
	for i := range seq {
		seq[i] = i
	}
	return seq
}
