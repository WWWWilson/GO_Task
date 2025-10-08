package main

import (
	"fmt"
	"golang_task3/repository"
	"golang_task3/service"
	"log"
)

func main() {
	// 数据库连接字符串
	dsn := "username:password@tcp(localhost:3306)/school_db?parseTime=true"
	// 初始化数据库连接
	db, err := repository.InitDB(dsn)
	if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
	defer db.Close()

	// 初始化仓库和服务
	studentRepo := repository.NewStudentRepository(db)
	studentService := service.NewStudentService(studentRepo)

	// 示例操作
    // 1. 创建学生
    student, err := studentService.CreateStudent("张三", 20, "三年级")
    if err != nil {
        log.Printf("Error creating student: %v", err)
    } else {
        fmt.Printf("Created student: %+v\n", student)
    }

	 // 2. 获取所有学生
    students, err := studentService.GetAllStudents()
    if err != nil {
        log.Printf("Error getting students: %v", err)
    } else {
        fmt.Println("All students:")
        for _, s := range students {
            fmt.Printf("  %+v\n", s)
        }
    }
	
    // 3. 根据年级查询
    gradeStudents, err := studentService.GetStudentsByGrade("高三")
    if err != nil {
        log.Printf("Error getting students by grade: %v", err)
    } else {
        fmt.Println("Grade 12 students:")
        for _, s := range gradeStudents {
            fmt.Printf("  %+v\n", s)
        }
    }

    // 4. students 表中姓名为 "张三" 的学生年级更新为 "四年级"
    rowsAffected,err := studentService.UpdateStudent("张三", 20, "四年级")
    if err != nil {
        log.Printf("Error getting students by grade: %v", err)
    } 

    if rowsAffected == 0 {
        fmt.Println("没有找到姓名为'张三'的学生")
    } else {
        fmt.Printf("成功更新了 %d 条记录\n", rowsAffected)
    }

     // 4. 删除 students 表中年龄小于 15 岁的学生记录
    rowsDeleted, err := studentService.DeleteStudent(15)
    if err != nil {
        log.Printf("删除失败: %v", err)
    } else {
        fmt.Printf("成功删除了 %d 条年龄小于15岁的学生记录\n", rowsDeleted)
    }

}