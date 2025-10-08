package repository

import (
	"database/sql"
	"fmt"
	"golang_task3/models"
)

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository{
	return &StudentRepository{db:db}
}

//初始化数据库连接
func InitDB(dataSourceName string)(*sql.DB,error){
	db,err := sql.Open("mysql",dataSourceName)
	if err != nil {
        return nil, err
    }
	 if err = db.Ping(); err != nil {
        return nil, err
    }
	return db, nil
}

// 创建学生
func (r *StudentRepository) CreateStudent(student *models.Student) error{
	query := "INSERT INTO students (name, age, grade) VALUES (?, ?, ?)"
	result, err := r.db.Exec(query, student.Name, student.Age, student.Grade)
	if err != nil {
        return err
    }
	id, err := result.LastInsertId()
    if err != nil {
        return err
    }
	student.ID = int(id)
	return nil
}

// 根据ID获取学生
func (r *StudentRepository) GetStudentByID(id int) (*models.Student, error) {
    query := "SELECT id, name, age, grade, created_at FROM students WHERE id = ?"
    row := r.db.QueryRow(query, id)
    
    student := &models.Student{}
    err := row.Scan(&student.ID, &student.Name, &student.Age, &student.Grade, &student.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("student not found")
        }
        return nil, err
    }
    
    return student, nil
}

// 查询 students 表中所有年龄大于 18 岁的学生信息
func (r *StudentRepository) GetStudentThan18(id int) ([]*models.Student, error) {
    query := "SELECT id, name, age, grade FROM students WHERE age > 18 ORDER BY age DESC"
    row,err := r.db.Query(query)
     if err = row.Err(); err != nil {
        return nil, err
    }
    defer row.Close()

    var students []*models.Student
    for row.Next() {
        var student *models.Student
        err := row.Scan(&student.ID, &student.Name, &student.Age, &student.Grade)
        if err != nil {
            return nil, err
        }
        students = append(students, student)
    }
     if err = row.Err(); err != nil {
        return nil, err
    }
    
    return students, nil
}


// 获取所有学生
func (r *StudentRepository) GetAllStudents() ([]*models.Student, error) {
    query := "SELECT id, name, age, grade, created_at FROM students ORDER BY id"
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var students []*models.Student
    for rows.Next() {
        student := &models.Student{}
        err := rows.Scan(&student.ID, &student.Name, &student.Age, &student.Grade, &student.CreatedAt)
        if err != nil {
            return nil, err
        }
        students = append(students, student)
    }
    
    return students, nil
}

// 更新学生信息
func (r *StudentRepository) UpdateStudent(student *models.Student) (int64,error) {
    query := "UPDATE students SET name = ?, age = ?, grade = ? WHERE id = ?"
    result, err := r.db.Exec(query, student.Name, student.Age, student.Grade, student.ID)
    if err != nil {
        return 0,err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return 0,err
    }
    return rowsAffected,nil
}

// 删除学生
func (r *StudentRepository) DeleteStudent(age int) (int64, error) {
     var count int
    countQuery := "SELECT COUNT(*) FROM students WHERE age < ?"
     err := r.db.QueryRow(countQuery, age).Scan(&count)
    if err != nil {
        return 0, err
    }
    if count == 0 {
        return 0, nil // 没有需要删除的记录
    }

     // 执行删除
    deleteQuery := "DELETE FROM students WHERE age < ?"
    result, err := r.db.Exec(deleteQuery, age)
    if err != nil {
        return 0, err
    }
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return 0, err
    }
    return rowsAffected, nil
}

// 根据年级查询学生
func (r *StudentRepository) GetStudentsByGrade(grade string) ([]*models.Student, error) {
    query := "SELECT id, name, age, grade, created_at FROM students WHERE grade = ? ORDER BY id"
    rows, err := r.db.Query(query, grade)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var students []*models.Student
    for rows.Next() {
        student := &models.Student{}
        err := rows.Scan(&student.ID, &student.Name, &student.Age, &student.Grade, &student.CreatedAt)
        if err != nil {
            return nil, err
        }
        students = append(students, student)
    }
    
    return students, nil
}