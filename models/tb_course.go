package models

import (
	"fmt"

	"github.com/luonannet/playground-backend/util"
	orm "gorm.io/gorm"
)

type TbCourse struct {
	Id            int    `orm:"column(id);pk"`
	Label         string `orm:"column(label);size(245);null"`
	Desc          string `orm:"column(desc);null"`
	Status        string `orm:"column(status);size(41);null"`
	TbCommunityId int    `orm:"column(tb_community_id)"`
	MaxPods       int    `orm:"column(max_pods);null"`
	MinPods       int    `orm:"column(min_pods);null"`
	EnvPath       string `orm:"column(env_path);null"`
	CoursePath    string `orm:"column(course_path);null"`
	TbClusterId   int    `orm:"column(tb_cluster_id)"`
}

func (t *TbCourse) TableName() string {
	return "tb_course"
}

func init() {
	orm.RegisterModel(new(TbCourse))
}

// AddTbCourse insert a new TbCourse into database and returns
// last inserted Id on success.
func AddTbCourse(m *TbCourse) (id int, err error) {
	o := util.GetDB()
	result := o.Create(m)
	return m.Id, result.Error
}

// GetTbCourseById retrieves TbCourse by Id. Returns error if
// Id doesn't exist
func GetTbCourseById(id int) (*TbCourse, error) {
	o := util.GetDB()
	v := &TbCourse{Id: id}
	if tx := o.First(v); tx.Error == nil {
		return v, nil
	}
	return nil, tx.Error
}

// GetAllTbCourse retrieves all TbCourse matches certain condition. Returns empty list if
// no records exist
func GetAllTbCourse(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) total int64, ml []*TbCourse, err error) {
	o := util.GetDB() 
	tx := o.Model(new(TbCourse)).Where("user_id", userid)
	tx.Count(&total)
	tx.Limit(limit).Offset(offset).Order("id desc").Scan(&ml)
	return nil, err
}

// UpdateTbCourse updates TbCourse by Id and returns error if
// the record to be updated doesn't exist
func UpdateTbCourseById(m *TbCourse) (err error) {
	o := orm.NewOrm()
	v := TbCourse{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteTbCourse deletes TbCourse by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTbCourse(id int) (err error) {
	o := orm.NewOrm()
	v := TbCourse{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&TbCourse{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
