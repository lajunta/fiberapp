package actions

import (
	"context"
	"log"
	"strconv"

	"github.com/acoshift/paginate"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Home page
func Home(c *fiber.Ctx) {
	var rowcount int64
	var err error
	rowcount, err = PageCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		rowcount = 0
	}

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

	Render(c, "Home", fiber.Map{"Data": pages, "Header": "Latest", "PN": pn}, "App Slogan")
}
