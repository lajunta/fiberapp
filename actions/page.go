package actions

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/acoshift/paginate"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"           // for BSON ObjectID
	"go.mongodb.org/mongo-driver/bson/primitive" // for BSON ObjectID
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Page a struct
type Page struct {
	ID        primitive.ObjectID `bson:"_id"`
	Title     string
	Body      string
	Author    string
	Tag       string
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

var (

	// PageTags is set by main package
	PageTags []string

	// PageCollection is set from main package
	PageCollection *mongo.Collection
)

// TagList is form page
func TagList(c *fiber.Ctx) {

	var rowcount int64
	name := c.Params("name")

	rowcount, _ = PageCollection.CountDocuments(context.TODO(), bson.M{"tag": name})

	pageNum, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil || pageNum == 0 {
		pageNum = 1
	}
	pn := paginate.New(pageNum, 10, rowcount)
	offset := (pageNum - 1) * 10
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"created_at": -1}).SetSkip(offset).SetLimit(10)
	var pages []Page

	// Passing bson.D{{}} as the filter matches all documents in the PageCollection
	cur, err := PageCollection.Find(context.TODO(), bson.M{"tag": name}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var page Page
		err := cur.Decode(&page)
		if err != nil {
			log.Fatal(err)
		}
		pages = append(pages, page)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	// Close the cursor once finished
	cur.Close(context.TODO())
	Render(c, "PageTags", fiber.Map{"Data": pages, "TagName": name, "PN": pn}, name)
}

// ShowPage according id
func ShowPage(c *fiber.Ctx) {
	//	id := c.Params("id")
	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	filter := bson.M{"_id": objID}
	var page Page

	err := PageCollection.FindOne(context.TODO(), filter).Decode(&page)
	if err != nil {
		log.Fatal(err)
	}

	e := strconv.FormatInt(page.UpdatedAt.Unix(), 10)

	if match := c.Get("If-None-Match"); match != "" {
		if strings.Contains(match, e) {
			c.Status(304)
			return
		}
	}

	//page.Body = string(blackfriday.Run([]byte(page.Body)))

	c.Fasthttp.Response.Header.Add("Etag", e)
	Render(c, "PageShow", page, page.Title)
}

// GetPages retrieve all pages
func GetPages(c *fiber.Ctx) {
	var rowcount int64
	var err error
	rowcount, err = PageCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		rowcount = 0
	}

	var params = fiber.Map{}
	params["Title"] = ""
	params["Tag"] = ""
	pageNum, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil || pageNum == 0 {
		pageNum = 1
	}

	pn := paginate.New(pageNum, 10, rowcount)
	offset := (pageNum - 1) * 10
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"created_at": -1}).SetSkip(offset).SetLimit(10)

	var pages []Page

	filter := bson.D{{}}
	if c.Query("title") != "" {
		t := c.Query("title")
		params["Title"] = t
		filter = append(filter, bson.E{Key: "title", Value: primitive.Regex{Pattern: t, Options: ""}})
	}

	if c.Query("tag") != "" {
		t := c.Query("tag")
		params["Tag"] = t
		filter = append(filter, bson.E{Key: "tag", Value: primitive.Regex{Pattern: t, Options: ""}})
	}

	// Passing bson.D{{}} as the filter matches all documents in the PageCollection
	cur, err := PageCollection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var page Page
		err := cur.Decode(&page)
		if err != nil {
			log.Fatal(err)
		}
		pages = append(pages, page)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	// Close the cursor once finished
	cur.Close(context.TODO())
	RenderAdmin(c, "PageIndex", fiber.Map{"Data": pages, "PN": pn, "Tags": PageTags, "Params": params})
}

// NewPage is form page
func NewPage(c *fiber.Ctx) {
	var page Page
	page.Title = "New Page"
	RenderAdmin(c, "PageForm", fiber.Map{"Page": page, "Tags": PageTags})
}

// EditPage update page
func EditPage(c *fiber.Ctx) {
	objID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return
	}
	filter := bson.M{"_id": objID}
	var page Page

	err = PageCollection.FindOne(context.TODO(), filter).Decode(&page)
	if err != nil {
		log.Fatal(err)
	}

	RenderAdmin(c, "PageForm", fiber.Map{"Page": page, "Tags": PageTags})
}

// CreatePage add page to db
func CreatePage(c *fiber.Ctx) {

	loc := time.FixedZone("UTC+8", +8*60*60)
	var page Page
	page.Title = c.FormValue("title")
	page.Tag = c.FormValue("tag")
	page.Author = c.FormValue("author")
	page.Body = c.FormValue("body")
	if c.FormValue("title") == "" || c.FormValue("body") == "" {
		RenderAdmin(c, "PageForm", fiber.Map{"Page": page, "Tags": PageTags})
		return
	}
	page.ID = primitive.NewObjectID()
	page.CreatedAt = time.Now().In(loc)
	page.UpdatedAt = time.Now().In(loc)

	_, err := PageCollection.InsertOne(context.TODO(), page)
	if err != nil {
		log.Fatalln(err.Error())
	}
	c.Redirect("/pages/show/" + page.ID.Hex())
}

// UpdatePage is
func UpdatePage(c *fiber.Ctx) {
	objID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return
	}
	loc := time.FixedZone("UTC+8", +8*60*60)
	var page = fiber.Map{}
	page["title"] = c.FormValue("title")
	page["tag"] = c.FormValue("tag")
	page["author"] = c.FormValue("author")
	page["body"] = c.FormValue("body")
	page["updated_at"] = time.Now().In(loc)

	_, err = PageCollection.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$set": page})

	if err != nil {
		log.Fatalln(err.Error())
	}
	c.Redirect("/admin/pages/1")
}

// DeletePage is
func DeletePage(c *fiber.Ctx) {
	objID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return
	}

	_, err = PageCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Redirect("/admin/pages/1")
}
