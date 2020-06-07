package actions

import (
	"bytes"
	"html/template"
	"time"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	isLogined = false
	tmpl      *template.Template
	funcs     = template.FuncMap{
		"to_date":  toDate,
		"to_date1": toDate1,
		"to_date2": toDate2,
		"htmlSafe": htmlSafe,
		"hex":      hex,
		"page_num": pageNum,
	}
)

// Render a template with a container
func Render(c *fiber.Ctx, tname string, data interface{}, title string) {
	tmpl = template.Must(template.New(tname).Funcs(funcs).ParseGlob("./views/*.html"))
	var w bytes.Buffer
	tmpl.ExecuteTemplate(&w, tname, data)
	c.Render("layout", fiber.Map{"Logined": isLogined, "Data": template.HTML(w.Bytes()), "Title": title})
}

// RenderAdmin a template with a container
func RenderAdmin(c *fiber.Ctx, tname string, data interface{}) {
	tmpl = template.Must(template.New(tname).Funcs(funcs).ParseGlob("./views/*.html"))
	var w bytes.Buffer
	tmpl.ExecuteTemplate(&w, tname, data)
	store := sess.Get(c)
	islogined := store.Get("islogined")
	if islogined != nil {
		isLogined = islogined.(bool)
	}
	c.Render("admin", fiber.Map{"Logined": isLogined, "Data": template.HTML(w.Bytes())})
}

func pageNum(page int64, num int) int64 {
	return (page-1)*10 + 1 + int64(num)
}

func toDate(tm time.Time) string {
	loc := time.FixedZone("UTC+8", +8*60*60)
	return tm.In(loc).Format("2006-01-02")
}

func toDate1(tm time.Time) string {
	loc := time.FixedZone("UTC+8", +8*60*60)
	return tm.In(loc).Format("200601021504")
}

func toDate2(tm time.Time) string {
	loc := time.FixedZone("UTC+8", +8*60*60)
	return tm.In(loc).Format("2006-01-02 15:04")
}

func htmlSafe(html string) template.HTML {
	return template.HTML(html)
}

func hex(id primitive.ObjectID) string {
	return id.Hex()
}

// Keys get keys  array
func Keys(myMap map[string][]string) []string {
	keys := make([]string, 0, len(myMap))
	for k := range myMap {
		keys = append(keys, k)
	}
	return keys
}

// Values get map values array
func Values(myMap map[string][]string) [][]string {
	keys := make([][]string, 0, len(myMap))
	for _, v := range myMap {
		keys = append(keys, v)
	}
	return keys
}
