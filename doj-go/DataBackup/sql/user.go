/**
 * @ Author: ClearDewy
 * @ Desc:
 **/
package sql

type UserInfoType struct {
	Uid                 String `json:"uid"`
	Username            String `json:"username"`
	Password            String `json:"password"`
	School              String `json:"school"`
	Major               String `json:"major"`
	Number              String `json:"number"`
	Name                String `json:"name"`
	Gender              int    `json:"gender"`
	Cf_username         String `json:"cf_username"`
	Email               String `json:"email"`
	Avatar              String `json:"avatar"`
	Signature           String `json:"signature"`
	Title_name          String `json:"title_name"`
	Title_color         String `json:"title_color"`
	System_auth         int    `json:"system_auth"`
	User_auth           int    `json:"user_auth"`
	Problem_auth        int    `json:"problem_auth"`
	Context_auth        int    `json:"context_auth"`
	Train_auth          int    `json:"train_auth"`
	Submit_auth         int    `json:"submit_auth"`
	Context_attend_auth int    `json:"context_attend_auth"`
	Train_attend_auth   int    `json:"train_attend_auth"`
}

func GetUserInfoByUsername(username string) (userInfo UserInfoType, err error) {
	row := db.QueryRow("SELECT `uid`,`username`,`password`,`school`,`major`,`number`,`name`,`gender`,`cf_username`,`email`,`avatar`,IFNULL(`signature`,''),`title_name`,`title_color`,`system_auth`,`user_auth`,`problem_auth`,`context_auth`,`train_auth`,`submit_auth`,`context_attend_auth`,`train_attend_auth` FROM user_info WHERE `username`= ? ", username)
	err = row.Scan(&userInfo.Uid, &userInfo.Username, &userInfo.Password, &userInfo.School, &userInfo.Major,
		&userInfo.Number, &userInfo.Name, &userInfo.Gender, &userInfo.Cf_username, &userInfo.Email,
		&userInfo.Avatar, &userInfo.Signature, &userInfo.Title_name, &userInfo.Title_color, &userInfo.System_auth,
		&userInfo.User_auth, &userInfo.Problem_auth, &userInfo.Context_auth, &userInfo.Train_auth, &userInfo.Submit_auth, &userInfo.Context_attend_auth, &userInfo.Train_attend_auth)
	return
}

func GetUserInfoByUid(uid string) (userInfo UserInfoType, err error) {
	row := db.QueryRow("SELECT `uid`,`username`,`password`,`school`,`major`,`number`,`name`,`gender`,`cf_username`,`email`,`avatar`,IFNULL(`signature`,''),`title_name`,`title_color`,`system_auth`,`user_auth`,`problem_auth`,`context_auth`,`train_auth`,`submit_auth`,`context_attend_auth`,`train_attend_auth` FROM user_info WHERE `uid`= ? ", uid)
	err = row.Scan(&userInfo.Uid, &userInfo.Username, &userInfo.Password, &userInfo.School, &userInfo.Major,
		&userInfo.Number, &userInfo.Name, &userInfo.Gender, &userInfo.Cf_username, &userInfo.Email,
		&userInfo.Avatar, &userInfo.Signature, &userInfo.Title_name, &userInfo.Title_color, &userInfo.System_auth,
		&userInfo.User_auth, &userInfo.Problem_auth, &userInfo.Context_auth, &userInfo.Train_auth, &userInfo.Submit_auth, &userInfo.Context_attend_auth, &userInfo.Train_attend_auth)
	return
}
