package model

import "time"

type Account struct {
	ID                    int64      `gorm:"column:id;type:bigint;primaryKey"`
	AccountID             int64      `json:"account_id,omitempty" gorm:"column:account_id;type:bigint;"`
	AppID                 *int64     `json:"app_id,omitempty" gorm:"column:app_id;type:bigint;default:null"` //
	Channel               *string    `json:"channel,omitempty" gorm:"column:channel;type:varchar(512);default:null"`
	InviterID             *int64     `json:"inviter_id,omitempty" gorm:"column:inviter_id;type:bigint;default:null"`
	AccountCatalog        int        `json:"account_catalog,omitempty" gorm:"column:account_catalog;type:integer;"`
	AccountStatus         int        `json:"account_status,omitempty" gorm:"column:account_status;type:integer;"`
	StatusStart           int64      `json:"status_start,omitempty" gorm:"column:status_start;type:bigint"`
	Name                  *string    `json:"name,omitempty" gorm:"column:name;type:varchar(64);default:null"` //
	Avatar                *string    `json:"avatar,omitempty" gorm:"column:avatar;type:varchar(1024);default:null"`
	Email                 *string    `json:"email,omitempty" gorm:"column:email;type:varchar(128);default:null"`
	PhoneArea             *string    `json:"phone_area,omitempty" gorm:"column:phone_area;type:varchar(32);default:null"`
	PhoneNumber           *string    `json:"phone_number,omitempty" gorm:"column:phone_number;type:varchar(32);default:null"`
	Password              string     `json:"password,omitempty" gorm:"column:password;type:varchar(255)"`
	RegisterType          int        `json:"register_type,omitempty" gorm:"column:register_type;type:integer"`
	RegisterIP            string     `json:"register_ip,omitempty" gorm:"column:register_ip;type:varchar(64)"`
	RegisterPlace         *string    `json:"register_place,omitempty" gorm:"column:register_place;type:varchar(128);default:null"`
	RegisterDevice        *int       `json:"register_device,omitempty" gorm:"column:register_device;type:integer;default:null"`
	LatestLoginIP         string     `json:"latest_login_ip,omitempty" gorm:"column:latest_login_ip;type:varchar(64)"`
	LatestLoginPlace      *string    `json:"latest_login_place,omitempty" gorm:"column:latest_login_place;type:varchar(128);default:null"`
	LatestLoginDevice     *int       `json:"latest_login_device,omitempty" gorm:"column:latest_login_device;type:integer;default:null"`
	LatestLoginAt         *time.Time `json:"latest_login_at,omitempty" gorm:"column:latest_login_at;type:timestamp(6)"`
	AssetSecurityStrategy int        `json:"asset_security_strategy" gorm:"column:asset_security_strategy;type:integer;default:-2" `
	AssetPassword         *string    `json:"asset_password,omitempty" gorm:"column:asset_password;type:varchar(64);default:null"`
	AssetPasswordStart    *int64     `json:"asset_password_start,omitempty" gorm:"column:asset_password_start;type:bigint;default:null"`
	GAKey                 *string    `json:"ga_key,omitempty" gorm:"column:ga_key;type:varchar(64);default:null"`
	AntiFishingText       *string    `json:"anti_fishing_text,omitempty" gorm:"column:anti_fishing_text;type:varchar(256);default:null"`
	KYCDataID             *int64     `json:"kyc_data_id,omitempty" gorm:"column:kyc_data_id;type:bigint;default:null"`
	CreatedAt             *time.Time `gorm:"column:created_at"`
	UpdatedAt             *time.Time `gorm:"column:updated_at"`
	Remark                string     `gorm:"column:remark;size:256"`
}
