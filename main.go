package main

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type SqlLogger struct {
	logger.Interface
}

func (l SqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	fmt.Printf("%v\n=========================\n", sql)

}

var db *gorm.DB

func main() {
	dsn := "root:1234@tcp(localhost:3306)/go_pek?parseTime=true"
	dial := mysql.Open(dsn)

	var err error
	db, err = gorm.Open(dial, &gorm.Config{
		Logger: &SqlLogger{},
		DryRun: false,
	})
	if err != nil {
		panic(err)
	}

	// db.AutoMigrate(Gender{}, Test{}, Customer{})

	// CreateGender("None")
	// GetGenders()
	// GetGendersById(4)
	// UpdateGenderById(4, "test")
	// Update2GenderById(4, "test222")
	// DeleteGenderById(4)
	// CreateTest(1, "Code0001")
	// CreateTest(2, "Code0002")
	// CreateTest(3, "Code0003")

	// DeleteTest(1)
	// GetTestById(1)
	// GetTest()

	// CreateCustomer("Jame", 1)
	// CreateCustomer("Katty", 2)

	GetCustomers()
}

func GetCustomers() {
	customers := []Customer{}
	tx := db.Preload(clause.Associations).Find(&customers)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}

	for _, c := range customers {
		fmt.Println(c.ID, "|", c.Name, c.Gender.Name)
	}
}
func CreateCustomer(name string, genderID uint) {
	customer := Customer{
		Name:     name,
		GenderID: genderID}

	tx := db.Create(&customer)
	if tx.Error != nil {
		return
	}

	fmt.Println("Customer created %v", customer)
}
func GetGenders() {
	genders := []Gender{}
	tx := db.Order("id desc").Find(&genders)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(genders)

}

func GetGendersById(id uint) {
	gender := Gender{}
	tx := db.Where("id = ?", id).Find(&gender)
	if tx.Error != nil {
		fmt.Println("=============")
		fmt.Println(tx.Error)
		return
	}
	fmt.Println("****** %v", gender)

}

func UpdateGenderById(id uint, name string) {
	gender := Gender{}
	tx := db.Where("id = ?", id).Find(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println("Before Update: %v", gender)

	gender.Name = name
	tx2 := db.Save(&gender)
	if tx2.Error != nil {
		fmt.Println(tx.Error)
		return
	}

	fmt.Println("After Update: %v", gender)

}

func Update2GenderById(id uint, name string) {
	gender := Gender{Name: name}
	tx := db.Model(&Gender{}).Where("id = ?", id).Updates(gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	GetGendersById(id)

}

func DeleteGenderById(id uint) {

	tx := db.Delete(&Gender{}, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	GetGendersById(id)

}
func CreateGender(name string) {
	gender := Gender{
		Name: name,
	}
	tx := db.Create(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
}

func CreateTest(code uint, name string) {
	test := Test{
		Code: code,
		Name: name,
	}

	tx := db.Create(&test)

	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
}

func GetTest() {
	tests := []Test{}
	tx := db.Find(&tests)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}

	for _, test := range tests {
		fmt.Println(test)
	}
}

func GetTestById(id uint) {
	test := &Test{}
	tx := db.Find(&test, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}

	fmt.Println(test)
}

func DeleteTest(id uint) {
	tx := db.Unscoped().Delete(&Test{}, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
	}

}

type Gender struct {
	ID   uint
	Name string `gorm:"unique;size(10)"`
}

type Test struct {
	gorm.Model
	Code uint
	Name string `gorm:"column:gender_name;size:10;default:Hello;not null"`
}

func (t Test) TableName() string {
	return "my_test"
}

type Customer struct {
	ID       uint
	Name     string
	Gender   Gender
	GenderID uint
}
