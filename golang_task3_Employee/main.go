package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// Employee 结构体，映射到employees表
type Employee struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	Department string `db:"department"`
	Salary     int    `db:"salary"`
}

type EmployeeService struct {
	db *sqlx.DB
}

// NewEmployeeService 创建EmployeeService实例
func NewEmployeeService(dataSourceName string) (*EmployeeService, error) {
	db, err := sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %v", err)
	}
	
	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库ping失败: %v", err)
	}
	
	return &EmployeeService{db: db}, nil
}

// GetTechDepartmentEmployees 查询所有部门为"技术部"的员工信息
func (es *EmployeeService) GetTechDepartmentEmployees() ([]Employee, error) {
	var employees []Employee
	
	// 使用sqlx的Select方法，自动映射到结构体切片
	query := "SELECT id, name, department, salary FROM employees WHERE department = ?"
	err := es.db.Select(&employees, query, "技术部")
	if err != nil {
		return nil, fmt.Errorf("查询技术部员工失败: %v", err)
	}
	
	return employees, nil
}

// GetHighestPaidEmployee 查询工资最高的员工信息
func (es *EmployeeService) GetHighestPaidEmployee() (*Employee, error) {
	var employee Employee
	
	// 查询工资最高的员工
	query := "SELECT id, name, department, salary FROM employees ORDER BY salary DESC LIMIT 1"
	err := es.db.Get(&employee, query)
	if err != nil {
		return nil, fmt.Errorf("查询最高工资员工失败: %v", err)
	}
	
	return &employee, nil
}

// Close 关闭数据库连接
func (es *EmployeeService) Close() error {
	return es.db.Close()
}

// CreateSampleData 创建测试数据
func (es *EmployeeService) CreateSampleData() error {
	// 创建employees表
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS employees (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			department VARCHAR(50) NOT NULL,
			salary INT NOT NULL
		)`
	
	_, err := es.db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("创建表失败: %v", err)
	}
	
	// 插入测试数据
	employees := []Employee{
		{Name: "张三", Department: "技术部", Salary: 15000},
		{Name: "李四", Department: "技术部", Salary: 18000},
		{Name: "王五", Department: "技术部", Salary: 20000},
		{Name: "赵六", Department: "销售部", Salary: 12000},
		{Name: "钱七", Department: "市场部", Salary: 25000}, // 最高工资
		{Name: "孙八", Department: "技术部", Salary: 16000},
	}
	
	for _, emp := range employees {
		_, err := es.db.NamedExec(`
			INSERT INTO employees (name, department, salary) 
			VALUES (:name, :department, :salary)`,
			emp)
		if err != nil {
			return fmt.Errorf("插入数据失败: %v", err)
		}
	}
	
	return nil
}

// 使用sqlx的事务示例 - 给所有技术部员工加薪
func (es *EmployeeService) GiveRaiseToTechDepartment(raiseAmount int) error {
	tx, err := es.db.Beginx()
	if err != nil {
		return fmt.Errorf("开始事务失败: %v", err)
	}
	
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()
	
	// 更新技术部员工工资
	result, err := tx.Exec(
		"UPDATE employees SET salary = salary + ? WHERE department = ?",
		raiseAmount, "技术部",
	)
	if err != nil {
		return fmt.Errorf("更新工资失败: %v", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %v", err)
	}
	
	fmt.Printf("成功为 %d 名技术部员工加薪 %d 元\n", rowsAffected, raiseAmount)
	
	return tx.Commit()
}

func main() {
	// 数据库连接配置
	dsn := "username:password@tcp(localhost:3306)/company_db?parseTime=true"

	// 创建EmployeeService实例
	employeeService, err := NewEmployeeService(dsn)
	if err != nil {
		log.Fatal("初始化EmployeeService失败:", err)
	}
	defer employeeService.Close()

	// 创建测试数据
	fmt.Println("创建测试数据...")
	err = employeeService.CreateSampleData()
	if err != nil {
		log.Printf("创建测试数据失败: %v", err)
	} else {
		fmt.Println("测试数据创建成功!")
	}

	// 1. 查询所有技术部员工
	fmt.Println("\n=== 技术部员工列表 ===")
	techEmployees, err := employeeService.GetTechDepartmentEmployees()
	if err != nil {
		log.Printf("查询技术部员工失败: %v", err)
	} else {
		if len(techEmployees) == 0 {
			fmt.Println("没有找到技术部员工")
		} else {
			for i, emp := range techEmployees {
				fmt.Printf("%d. ID: %d, 姓名: %s, 部门: %s, 工资: %d\n", 
					i+1, emp.ID, emp.Name, emp.Department, emp.Salary)
			}
		}
	}

	// 4. 给技术部员工加薪示例
	fmt.Println("\n=== 给技术部员工加薪 ===")
	err = employeeService.GiveRaiseToTechDepartment(1000)
	if err != nil {
		log.Printf("加薪操作失败: %v", err)
	} else {
		// 重新查询技术部员工查看加薪结果
		updatedTechEmployees, err := employeeService.GetTechDepartmentEmployees()
		if err != nil {
			log.Printf("查询更新后的技术部员工失败: %v", err)
		} else {
			fmt.Println("\n加薪后的技术部员工:")
			for i, emp := range updatedTechEmployees {
				fmt.Printf("%d. ID: %d, 姓名: %s, 部门: %s, 工资: %d\n", 
					i+1, emp.ID, emp.Name, emp.Department, emp.Salary)
			}
		}
	}
}

