package database

import (
	"fmt"
	"gorm.io/gorm"
	"testing"
)

func TestMysql_Init(t *testing.T) {

}

func TestMysql_ListAllInfo(t *testing.T) {
	fmt.Println("Test")
	m := new(Mysql)
	m.connect("192.168.52.153", "3306", "apikiller", "root", "123456")
	m.addHttpItem("123123")
}

func TestMysql_addHttpItem(t *testing.T) {
	type fields struct {
		db       *gorm.DB
		MaxCount int
	}
	type args struct {
		item string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
		{
			name:   "",
			fields: fields{},
			args:   args{},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mysql{
				db:       tt.fields.db,
				MaxCount: tt.fields.MaxCount,
			}
			if got := m.addHttpItem(tt.args.item); got != tt.want {
				t.Errorf("addHttpItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
