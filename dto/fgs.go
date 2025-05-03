package dto

type InsertFGS struct {
	FGGId  int `json:"fgs_id" db:"fgs_id"`
	AreaId int `json:"area_id" db:"area_id"`
	Count  int `json:"count" db:"count"`
}

type StaffFGS struct {
	StaffId  int    `json:"staff_id" db:"staff_id"`
	Fullname string `json:"fullname" db:"fullname"`
	AreaId   int    `json:"id" db:"id"`
	AreaName string `json:"area" db:"area"`
	FGGId    *int    `json:"fgs_id" db:"fgs_id"`
	Count    *int    `json:"count" db:"count"`
}
