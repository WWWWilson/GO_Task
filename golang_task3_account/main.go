package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type Account struct {
	ID      int     `json:"id"`
	Balance float64 `json:"balance"`
}

type Transaction struct {
	ID            int       `json:"id"`
	FromAccountID int       `json:"from_account_id"`
	ToAccountID   int       `json:"to_account_id"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type BankService struct {
	db *sql.DB
}

func NewBankService(dataSourceName string) (*BankService, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &BankService{db: db}, nil
}

// Transfer 执行转账操作（使用事务保证数据一致性）
func (bs *BankService) Transfer(fromAccountID, toAccountID int, amount float64) error {
	// 开启事务
	tx, err := bs.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// 检查转账金额是否有效
	if amount <= 0 {
		return errors.New("转账金额必须大于0")
	}

	// 检查转出账户是否存在且余额充足
	var fromBalance float64
	err = tx.QueryRow("SELECT balance FROM accounts WHERE id = ? FOR UPDATE", fromAccountID).Scan(&fromBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("转出账户不存在")
		}
		return err
	}

	if fromBalance < amount {
		return errors.New("账户余额不足")
	}

	// 检查转入账户是否存在
	var toAccountExists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM accounts WHERE id = ?)", toAccountID).Scan(&toAccountExists)
	if err != nil {
		return err
	}
	if !toAccountExists {
		return errors.New("转入账户不存在")
	}

	// 创建交易记录（状态为pending）
	result, err := tx.Exec(
		"INSERT INTO transactions (from_account_id, to_account_id, amount, status) VALUES (?, ?, ?, 'pending')",
		fromAccountID, toAccountID, amount,
	)
	if err != nil {
		return err
	}

	transactionID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// 更新转出账户余额
	_, err = tx.Exec(
		"UPDATE accounts SET balance = balance - ? WHERE id = ?",
		amount, fromAccountID,
	)
	if err != nil {
		return err
	}

	// 更新转入账户余额
	_, err = tx.Exec(
		"UPDATE accounts SET balance = balance + ? WHERE id = ?",
		amount, toAccountID,
	)
	if err != nil {
		return err
	}

	// 更新交易记录状态为completed
	_, err = tx.Exec(
		"UPDATE transactions SET status = 'completed' WHERE id = ?",
		transactionID,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetAccountBalance 获取账户余额
func (bs *BankService) GetAccountBalance(accountID int) (float64, error) {
	var balance float64
	err := bs.db.QueryRow("SELECT balance FROM accounts WHERE id = ?", accountID).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("账户不存在")
		}
		return 0, err
	}
	return balance, nil
}

// GetTransactionHistory 获取账户交易记录
func (bs *BankService) GetTransactionHistory(accountID int) ([]Transaction, error) {
	rows, err := bs.db.Query(`
		SELECT id, from_account_id, to_account_id, amount, status, created_at 
		FROM transactions 
		WHERE from_account_id = ? OR to_account_id = ?
		ORDER BY created_at DESC
	`, accountID, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.FromAccountID, &t.ToAccountID, &t.Amount, &t.Status, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

// CreateAccount 创建新账户
func (bs *BankService) CreateAccount(initialBalance float64) (int, error) {
	result, err := bs.db.Exec(
		"INSERT INTO accounts (balance) VALUES (?)",
		initialBalance,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func main() {
	dsn := "username:password@tcp(localhost:3306)/bank_db?parseTime=true"
	bankService, err := NewBankService(dsn)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer bankService.db.Close()

	// 示例使用
	// 1. 创建账户
	account1, _ := bankService.CreateAccount(1000.00)
	account2, _ := bankService.CreateAccount(500.00)
	fmt.Printf("创建账户: ID=%d, 初始余额=%.2f\n", account1, 1000.00)
	fmt.Printf("创建账户: ID=%d, 初始余额=%.2f\n", account2, 500.00)

	// 2. 执行转账
	amount := 100.00
	fmt.Printf("执行转账: 从账户 %d 向账户 %d 转账 %.2f\n", account1, account2, amount)

	err = bankService.Transfer(account1, account2, amount)
	if err != nil {
		log.Printf("转账失败: %v\n", err)
	} else {
		fmt.Println("转账成功!")
	}

	// 3. 查询余额
	balance1, _ := bankService.GetAccountBalance(account1)
	balance2, _ := bankService.GetAccountBalance(account2)
	fmt.Printf("账户 %d 余额: %.2f\n", account1, balance1)
	fmt.Printf("账户 %d 余额: %.2f\n", account2, balance2)

	// 4. 查询交易记录
	transactions, _ := bankService.GetTransactionHistory(account1)
	fmt.Printf("账户 %d 的交易记录: %+v\n", account1, transactions)

}