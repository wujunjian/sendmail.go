package util

import (
	config "../config"
	"strings"
	"fmt"
    "strconv"
    "time"
)

type themeinfo struct {
	tid int
	tp int
	theme_name string
	custom_name string
	mtime string
	uid uint64
}

type userinfo struct {
    country string
	uid uint64
	name string
	nickname string
	email string
}

func getthemeinfo(themepkgname string) (themeinfo,bool) {
	var tmpthemeinfo themeinfo

	var sqlstr string = "select tid,type,theme_name, custom_name, mtime,uid from theme where custom_name = '" +
	themepkgname + "' limit 1"

	//var tid,tp,uid int
    //var theme_name,custom_name,mtime string
    var mtime int64
	err := QueryOneRow(sqlstr, &tmpthemeinfo.tid,
							   &tmpthemeinfo.tp, 
							   &tmpthemeinfo.theme_name, 
							   &tmpthemeinfo.custom_name, 
							   &mtime,
							   &tmpthemeinfo.uid)
	if err != nil {
		fmt.Println("row.Scan err themepkgname = ", themepkgname, err)
		return tmpthemeinfo, false
	} else {
        //fmt.Println(tmpthemeinfo.tid, tmpthemeinfo.tp, tmpthemeinfo.theme_name, tmpthemeinfo.custom_name, tmpthemeinfo.mtime, tmpthemeinfo.uid)
        tmpthemeinfo.mtime = time.Unix(mtime, 0).String()
		return tmpthemeinfo, true
	}
}

func getuserinfo(uid uint64) (userinfo,bool) {
	var tmpuserinfo userinfo 

	var sqlstr string = "select country, uid, name, nickname, email from user where uid = '" +
	strconv.FormatUint(uid, 10) + "' limit 1"

	//var tid,tp,uid int
	//var theme_name,custom_name,mtime string
    err := QueryOneRow(sqlstr, &tmpuserinfo.country, 
                               &tmpuserinfo.uid, 
							   &tmpuserinfo.name, 
							   &tmpuserinfo.nickname, 
							   &tmpuserinfo.email)
	if err != nil {  
		fmt.Println("row.Scan err , uid = ", uid,  err)
		return tmpuserinfo, false
	} else {
		///fmt.Println(tmpuserinfo.uid, 
		///	        tmpuserinfo.name, 
		///	        tmpuserinfo.nickname, 
		///	        tmpuserinfo.email)
		return tmpuserinfo, true
	}
}

func Sendcrashmail(errcount, themepkgname, datestr string) {

	var designerurl string 
	tmpthemeinfo, HaveDesignerInfo := getthemeinfo(themepkgname)
	if !HaveDesignerInfo {
		tmpthemeinfo.mtime = config.Gf.DefaultNames.Mtime
		tmpthemeinfo.theme_name = config.Gf.DefaultNames.ThemeName
        designerurl = config.Gf.DefaultNames.DesignerUrl
	} else {
		tid := strconv.Itoa(tmpthemeinfo.tid)
		switch tmpthemeinfo.tp {
		case 0:    designerurl = ""
		case 2:    designerurl = ""
		case 3:    designerurl = ""
		case 8:    designerurl = ""
		case 10:   designerurl = ""
		case 5:    designerurl = ""
		case 7:    designerurl = ""
		case 9:    designerurl = ""
		default:   designerurl = ""
		}
    }

	tmpuserinfo, isuserok := getuserinfo(tmpthemeinfo.uid)
	if !isuserok {
		tmpuserinfo.email    = config.Gf.DefaultNames.Email
		tmpuserinfo.name     = config.Gf.DefaultNames.AccountName
		tmpuserinfo.nickname = config.Gf.DefaultNames.NickName
	}

    var tmp string
    var isCn bool
    if tmpuserinfo.country == "CN" || !isuserok {
        tmp = config.CnTemplet
        isCn = true
    } else {
        tmp = config.EnTemplet
        isCn = false
    }

    var tmpSubJect string 
    if isCn {
        tmpSubJect = config.Gf.Subject
    } else {
        tmpSubJect = config.Gf.EnSubject
    }
    
    if HaveDesignerInfo {
        if isCn {
            tmp = strings.Replace(tmp, "%%MESSAGE%%", config.Gf.DefaultNames.HaveDesignerInfo, -1)
        } else {
            tmp = strings.Replace(tmp, "%%MESSAGE%%", config.Gf.DefaultNames.EnHaveDesignerInfo, -1)
        }
        tmpSubJect = strings.Replace(tmpSubJect, "%%THEMENAME%%", tmpthemeinfo.theme_name, -1)
    } else {
        if isCn {
            tmp = strings.Replace(tmp, "%%MESSAGE%%", config.Gf.DefaultNames.NoDesignerInfo, -1)
        } else {
            tmp = strings.Replace(tmp, "%%MESSAGE%%", config.Gf.DefaultNames.EnNoDesignerInfo, -1)
        }

        tmp = strings.Replace(tmp, "%%THEMENAME%%", "%%THEMEPKGNAME%%", -1)
        tmpSubJect = strings.Replace(tmpSubJect, "%%THEMENAME%%", themepkgname, -1)
    }

	tmp = strings.Replace(tmp, "%%THEMEPKGNAME%%", themepkgname, -1)          //
	tmp = strings.Replace(tmp, "%%COUNT%%", errcount, -1)                     //
    tmp = strings.Replace(tmp, "%%DATE%%", datestr, -1)                       //
    tmp = strings.Replace(tmp, "%%THEMENAME%%", tmpthemeinfo.theme_name, -1)  //
	tmp = strings.Replace(tmp, "%%NICKNAME%%", tmpuserinfo.nickname, -1)      //
	
	if len(tmpuserinfo.email) > 0 && isuserok {
		tmp = strings.Replace(tmp, "%%ACCOUNTNAME%%", tmpuserinfo.email, -1)  //
	} else {
		tmp = strings.Replace(tmp, "%%ACCOUNTNAME%%", tmpuserinfo.name, -1)   //	
	}

	tmp = strings.Replace(tmp, "%%LASTMODIFYDATE%%", tmpthemeinfo.mtime, -1)  //最后修改日期
	tmp = strings.Replace(tmp, "%%GPURL%%", "https://play.google.com/store/apps/details?id="+themepkgname, -1)  //

	tmp = strings.Replace(tmp, "%%CMDESIGNER%%", designerurl, -1)             //

    var ToEmail string
    //如果默认邮箱为 wujunjian@cmcm.com 则为测试
    if config.Gf.DefaultNames.Email == "wujunjian@cmcm.com" || len(tmpuserinfo.email) == 0 {
        ToEmail = config.Gf.DefaultNames.Email
    } else {
        ToEmail = tmpuserinfo.email
    }

	msg := []byte(
		"From: Theme Designer Center<@.>\r\n" +
		"To: "+ToEmail+"\r\n" +
		tmpSubJect+"\r\n" +
		"Content-Type: text/html;\r\n\tcharset=\"utf-8\"\r\n" +
		"\r\n" +
		tmp+"\r\n")

	//
	SendMail(msg, ToEmail)
}