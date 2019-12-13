package models

type UserWithdraw struct {
	Id            int     `json:"id" orm:"column(id);auto"`
	Uid           int64   `json:"uid" orm:"column(uid)"`
	Amount        float64 `json:"amount" orm:"column(amount);digits(9);decimals(2)"`
	Type          int8    `json:"type" orm:"column(uid)"`
	Pay_method    int8     `json:"id" orm:"column(pay_method);auto"`
	Realname      int64   `json:"uid" orm:"column(uid)"`
	Bank_account  float64 `json:"amount" orm:"digits(9);decimals(2)"`
	Title         int8    `json:"amount" orm:"digits(9);decimals(2)"`
	Status_remark int     `json:"id" orm:"column(id);auto"`
	Biz_order_id  int64   `json:"uid" orm:"column(uid)"`
	Out_order_id  float64 `json:"amount" orm:"digits(9);decimals(2)"`
	Pay_status    int8    `json:"amount" orm:"digits(9);decimals(2)"`
	Pay_status_re int     `json:"id" orm:"column(id);auto"`
	Balance_id    int64   `json:"uid" orm:"column(uid)"`
	Admin_uid     int     `json:"id" orm:"column(id);auto"`
	Ymd           int64   `json:"uid" orm:"column(uid)"`
	Status        float64 `json:"amount" orm:"digits(9);decimals(2)"`
	Created_at    int8    `json:"amount" orm:"digits(9);decimals(2)"`
	Updated_at    int     `json:"id" orm:"column(id);auto"`
}

