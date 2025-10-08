package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// Book 结构体，映射到books表
type Book struct {
	ID     int     `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
}

type BookService struct {
	db *sqlx.DB
}

// NewBookService 创建BookService实例
func NewBookService(dataSourceName string) (*BookService, error) {
	db, err := sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %v", err)
	}
	
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库ping失败: %v", err)
	}
	
	return &BookService{db: db}, nil
}

// Close 关闭数据库连接
func (bs *BookService) Close() error {
	return bs.db.Close()
}

// CreateBooksTable 创建books表
func (bs *BookService) CreateBooksTable() error {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS books (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(200) NOT NULL,
			author VARCHAR(100) NOT NULL,
			price DECIMAL(10,2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_price (price),
			INDEX idx_author (author)
		)`
	
	_, err := bs.db.Exec(createTableSQL)
	return err
}

// InsertSampleData 插入测试数据
func (bs *BookService) InsertSampleData() error {
	books := []Book{
		{Title: "Go语言编程", Author: "张三", Price: 65.50},
		{Title: "数据库系统概念", Author: "李四", Price: 89.00},
		{Title: "算法导论", Author: "王五", Price: 128.00},
		{Title: "设计模式", Author: "赵六", Price: 45.00},
		{Title: "计算机网络", Author: "钱七", Price: 75.50},
		{Title: "操作系统", Author: "孙八", Price: 58.00},
		{Title: "软件工程", Author: "周九", Price: 42.50},
		{Title: "机器学习", Author: "吴十", Price: 99.00},
		{Title: "深度学习", Author: "郑十一", Price: 120.00},
		{Title: "Python编程", Author: "王十二", Price: 55.00},
	}
	
	for _, book := range books {
		_, err := bs.db.NamedExec(`
			INSERT INTO books (title, author, price) 
			VALUES (:title, :author, :price)`, book)
		if err != nil {
			return err
		}
	}
	
	return nil
}

// GetBooksByPrice 查询价格大于指定值的书籍
func (bs *BookService) GetBooksByPrice(minPrice float64) ([]Book, error) {
	var books []Book
	
	query := "SELECT id, title, author, price FROM books WHERE price > ? ORDER BY price DESC"
	err := bs.db.Select(&books, query, minPrice)
	if err != nil {
		return nil, fmt.Errorf("查询价格大于 %.2f 的书籍失败: %v", minPrice, err)
	}
	
	return books, nil
}

// GetBooksByPriceRange 查询价格在指定范围内的书籍（复杂查询示例）
func (bs *BookService) GetBooksByPriceRange(minPrice, maxPrice float64) ([]Book, error) {
	var books []Book
	
	query := `
		SELECT id, title, author, price 
		FROM books 
		WHERE price > ? AND price < ? 
		ORDER BY price DESC, title ASC`
	
	err := bs.db.Select(&books, query, minPrice, maxPrice)
	if err != nil {
		return nil, fmt.Errorf("查询价格在 %.2f 到 %.2f 之间的书籍失败: %v", minPrice, maxPrice, err)
	}
	
	return books, nil
}

// GetBooksByAuthorAndPrice 根据作者和价格查询书籍
func (bs *BookService) GetBooksByAuthorAndPrice(author string, minPrice float64) ([]Book, error) {
	var books []Book
	
	query := `
		SELECT id, title, author, price 
		FROM books 
		WHERE author LIKE ? AND price > ? 
		ORDER BY price DESC`
	
	err := bs.db.Select(&books, query, "%"+author+"%", minPrice)
	if err != nil {
		return nil, fmt.Errorf("查询作者包含 '%s' 且价格大于 %.2f 的书籍失败: %v", author, minPrice, err)
	}
	
	return books, nil
}

// GetExpensiveBooksByAuthor 查询指定作者的价格最高的书籍
func (bs *BookService) GetExpensiveBooksByAuthor(author string) (*Book, error) {
	var book Book
	
	query := `
		SELECT id, title, author, price 
		FROM books 
		WHERE author LIKE ? 
		ORDER BY price DESC 
		LIMIT 1`
	
	err := bs.db.Get(&book, query, "%"+author+"%")
	if err != nil {
		return nil, fmt.Errorf("查询作者 '%s' 的最贵书籍失败: %v", author, err)
	}
	
	return &book, nil
}

// GetBooksStatistics 获取书籍价格统计信息（复杂查询）
func (bs *BookService) GetBooksStatistics() (struct {
	TotalBooks    int     `db:"total_books"`
	AveragePrice  float64 `db:"avg_price"`
	MaxPrice      float64 `db:"max_price"`
	MinPrice      float64 `db:"min_price"`
	ExpensiveCount int    `db:"expensive_count"`
}, error) {
	var stats struct {
		TotalBooks    int     `db:"total_books"`
		AveragePrice  float64 `db:"avg_price"`
		MaxPrice      float64 `db:"max_price"`
		MinPrice      float64 `db:"min_price"`
		ExpensiveCount int    `db:"expensive_count"`
	}
	
	query := `
		SELECT 
			COUNT(*) as total_books,
			AVG(price) as avg_price,
			MAX(price) as max_price,
			MIN(price) as min_price,
			COUNT(CASE WHEN price > 50 THEN 1 END) as expensive_count
		FROM books`
	
	err := bs.db.Get(&stats, query)
	if err != nil {
		return stats, fmt.Errorf("获取书籍统计信息失败: %v", err)
	}
	
	return stats, nil
}

// UpdateBookPrice 更新书籍价格（事务示例）
func (bs *BookService) UpdateBookPrice(bookID int, newPrice float64) error {
	tx, err := bs.db.Beginx()
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
	
	// 更新书籍价格
	result, err := tx.Exec(
		"UPDATE books SET price = ? WHERE id = ?",
		newPrice, bookID,
	)
	if err != nil {
		return fmt.Errorf("更新书籍价格失败: %v", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %v", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("未找到ID为 %d 的书籍", bookID)
	}
	
	fmt.Printf("成功更新书籍ID %d 的价格为 %.2f\n", bookID, newPrice)
	
	return tx.Commit()
}


func main(){
	// 数据库连接配置
	dsn := "username:password@tcp(localhost:3306)/bookstore?parseTime=true&charset=utf8mb4"
	
	// 创建BookService实例
	bookService, err := NewBookService(dsn)
	if err != nil {
		log.Fatal("初始化BookService失败:", err)
	}	
	defer bookService.Close()

	// 创建表
	fmt.Println("创建books表...")
	err = bookService.CreateBooksTable()
	if err != nil {
		log.Fatal("创建表失败:", err)
	}

	// 插入测试数据
	fmt.Println("插入测试数据...")
	err = bookService.InsertSampleData()
	if err != nil {
		log.Printf("插入测试数据失败: %v", err)
	} else {
		fmt.Println("测试数据插入成功!")
	}

	// 1. 查询价格大于50元的书籍
	fmt.Println("\n=== 价格大于50元的书籍 ===")
	expensiveBooks, err := bookService.GetBooksByPrice(50.0)
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("找到 %d 本价格大于50元的书籍:\n", len(expensiveBooks))
		for i, book := range expensiveBooks {
			fmt.Printf("%d. 《%s》- %s - ￥%.2f\n", 
				i+1, book.Title, book.Author, book.Price)
		}
	}
	// 2. 查询价格在指定范围内的书籍
	fmt.Println("\n=== 价格在60-100元之间的书籍 ===")
	midRangeBooks, err := bookService.GetBooksByPriceRange(60.0, 100.0)
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("找到 %d 本价格在60-100元之间的书籍:\n", len(midRangeBooks))
		for i, book := range midRangeBooks {
			fmt.Printf("%d. 《%s》- %s - ￥%.2f\n", 
				i+1, book.Title, book.Author, book.Price)
		}
	}
	
	// 3. 根据作者和价格查询
	fmt.Println("\n=== 作者包含'王'且价格大于50元的书籍 ===")
	authorBooks, err := bookService.GetBooksByAuthorAndPrice("王", 50.0)
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("找到 %d 本符合条件的书籍:\n", len(authorBooks))
		for i, book := range authorBooks {
			fmt.Printf("%d. 《%s》- %s - ￥%.2f\n", 
				i+1, book.Title, book.Author, book.Price)
		}
	}
	
	// 4. 获取最贵的书籍
	fmt.Println("\n=== 最贵的书籍 ===")
	mostExpensive, err := bookService.GetExpensiveBooksByAuthor("")
	if err != nil {
		log.Printf("查询失败: %v", err)
	} else {
		fmt.Printf("最贵的书籍: 《%s》- %s - ￥%.2f\n", 
			mostExpensive.Title, mostExpensive.Author, mostExpensive.Price)
	}
	
	// 5. 获取统计信息
	fmt.Println("\n=== 书籍统计信息 ===")
	stats, err := bookService.GetBooksStatistics()
	if err != nil {
		log.Printf("获取统计信息失败: %v", err)
	} else {
		fmt.Printf("总书籍数量: %d\n", stats.TotalBooks)
		fmt.Printf("平均价格: ￥%.2f\n", stats.AveragePrice)
		fmt.Printf("最高价格: ￥%.2f\n", stats.MaxPrice)
		fmt.Printf("最低价格: ￥%.2f\n", stats.MinPrice)
		fmt.Printf("价格大于50元的书籍数量: %d\n", stats.ExpensiveCount)
	}
	
	// 6. 更新书籍价格示例
	fmt.Println("\n=== 更新书籍价格 ===")
	// 假设我们要更新第一本书的价格
	if len(expensiveBooks) > 0 {
		err = bookService.UpdateBookPrice(expensiveBooks[0].ID, 80.00)
		if err != nil {
			log.Printf("更新价格失败: %v", err)
		}
	}

}