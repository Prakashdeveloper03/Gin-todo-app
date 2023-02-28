package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
	"net/http"
)

var (
	db  *gorm.DB
	err error
)

type Todo struct {
	ID       uint `gorm:"primaryKey"`
	Title    string
	Complete bool
}

func main() {
	r := gin.Default()
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Todo{})
	r.SetHTMLTemplate(template.Must(template.ParseFiles("templates/base.go.tmpl")))
	r.GET("/", home)
	r.POST("/add", add)
	r.GET("/update/:id", update)
	r.GET("/delete/:id", delete)
	r.Run(":8080")
}

func home(c *gin.Context) {
	var todoList []Todo
	db.Find(&todoList)
	c.HTML(http.StatusOK, "base.go.tmpl", gin.H{
		"todoList": todoList,
	})
}

func add(c *gin.Context) {
	title := c.PostForm("title")
	newTodo := Todo{Title: title, Complete: false}
	db.Create(&newTodo)
	c.Redirect(http.StatusSeeOther, "/")
}

func update(c *gin.Context) {
	var todo Todo
	id := c.Param("id")
	db.First(&todo, id)
	todo.Complete = !todo.Complete
	db.Save(&todo)
	c.Redirect(http.StatusSeeOther, "/")
}

func delete(c *gin.Context) {
	var todo Todo
	id := c.Param("id")
	db.Delete(&todo, id)
	c.Redirect(http.StatusSeeOther, "/")
}
