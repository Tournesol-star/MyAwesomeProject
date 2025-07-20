package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"net/http"
	"strconv"

	"log"
	"database/sql"
)

type Student struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    Age       int    `json:"age"`
    Gender    string `json:"gender"`
    CreatedAt string `json:"created_at"` 
    UpdatedAt string `json:"updated_at"` 
}
var db *sql.DB

func main() {
	// 初始化数据库连接
	dsn := "root:@Hu123456@tcp(localhost:3306)/student_db"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 测试数据库连接
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// 创建Gin引擎
	r := gin.Default()

	// 定义API路由
	api := r.Group("/api")
	{
		api.GET("/students", ListStudents)
		api.POST("/students", CreateStudent)
		api.GET("/students/:id", GetStudent)
		api.PUT("/students/:id", UpdateStudent)
		api.DELETE("/students/:id", DeleteStudent)
	}

	fmt.Println("Server is running on :8080")
	r.Run(":8080")
}


// 获取所有学生
func ListStudents(c *gin.Context) {
    rows, err := db.Query("SELECT * FROM students")
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var students []Student
    for rows.Next() {
        var student Student
        // 按表中列的顺序，扫描到对应的结构体字段
        err := rows.Scan(&student.ID, &student.Name, &student.Age, &student.Gender, &student.CreatedAt, &student.UpdatedAt) 
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }
        students = append(students, student)
    }

    c.JSON(200, students)
}


// 创建学生
func CreateStudent(c *gin.Context) {
	var student Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO students (name, age, gender) VALUES (?, ?, ?)", student.Name, student.Age, student.Gender)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	student.ID = int(id)
	c.JSON(201, student)
}

// 根据ID查询学生
func GetStudent(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var student Student
    // 按表中列的顺序，扫描到对应的结构体字段
    err = db.QueryRow("SELECT * FROM students WHERE id=?", id).Scan(&student.ID, &student.Name, &student.Age, &student.Gender, &student.CreatedAt, &student.UpdatedAt) 
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(404, gin.H{"error": "student not found"})
            return
        }
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, student)
}

// 根据ID更新学生
func UpdateStudent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var student Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 根据ID更新学生
	result, err := db.Exec("UPDATE students SET name=?, age=?, gender=? WHERE id=?", student.Name, student.Age, student.Gender, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "student not found"})
		return
	}

	c.JSON(200, gin.H{"message": "student updated successfully"})
}

// 根据ID删除学生
func DeleteStudent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("DELETE FROM students WHERE id=?", id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "student not found"})
		return
	}

	c.JSON(200, gin.H{"message": "student deleted successfully"})
}
