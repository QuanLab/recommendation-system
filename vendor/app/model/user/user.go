package user

import (
	"github.com/coocood/freecache"
	"app/shared/database/mysql"
	"app/shared/util"
	"encoding/json"
	"log"
	"fmt"
)

const (
	CACHE_SIZE  = 100 * 1024 * 1024
	EXPIRE_TIME = 60 * 86400 //60 days
)

var (
	cache *freecache.Cache
)

func init() {
	cache = freecache.NewCache(CACHE_SIZE)
}

type UserProfile struct {
	Guid       int64                 `json:"guid,omitempty"`
	PeriodTime map[string]PeriodTime `json:"periods,omitempty"`
}

type PeriodTime struct {
	TotalClick int64             `json:"TotalClick,omitempty"`
	MapPTCat   map[int64]float64 `json:"MapPTCat,omitempty"`
}

func GetUserProfile(guid int64) (UserProfile, error) {
	var profile string
	query := "SELECT profile FROM user_profile WHERE guid=?";
	rows, err := mysql.DB.Query(query, guid)
	defer rows.Close()

	var userProfile = UserProfile{}
	if err != nil {
		log.Println(err)
		return userProfile, err
	}
	if rows.Next() {
		rows.Scan(&profile)
		log.Println(profile)
		err := json.Unmarshal([]byte(profile), &userProfile)
		if err != nil {
			log.Println("Unable to Unmarshal user profile")
			return userProfile, err
		}
		userProfile.Guid = guid
		return userProfile, nil
	}
	return userProfile, err
}

func GetUserProfileCached(guid int64) (UserProfile, error) {
	var userProfile UserProfile
	value, err := cache.Get([]byte(string(guid)))
	if err!=nil {
		fmt.Println("Not found user ", guid)
		return userProfile, err
	}
	e := json.Unmarshal(value, &userProfile)
	if e != nil {
		log.Println("Unable to Unmarshal user profile")
		return userProfile, err
	}
	return userProfile, nil
}

func LoadUserProfileToCache() error {
	log.Println("Load user profile to cache")
	query := "SELECT guid, profile FROM user_profile";
	rows, err := mysql.DB.Query(query)
	if err != nil {
		log.Println(err)
		return err
	}
	defer rows.Close()

	lastMonths := util.GetListYearMonthFromCurrent(3)
	for rows.Next() {
		var userProfile UserProfile
		var profile string
		var guid int64
		rows.Scan(&guid, &profile)
		err := json.Unmarshal([]byte(profile), &userProfile)
		if err!=nil {
			log.Println(err)
			continue
		}
		var isDelete bool = true

		for _, month := range lastMonths {
			if _, ok := userProfile.PeriodTime[month]; ok {
				isDelete = false
				break
			}
		}
		if isDelete {
			DeleteUserProfile(guid)
		}
		cache.Set([]byte(string(guid)), []byte(profile), EXPIRE_TIME)
	}
	log.Println("Load user profile to cache done.")
	return nil
}

func DeleteUserProfile(guid int64) (error) {
	query := "DELETE FROM user_profile WHERE guid=?";
	_, err := mysql.DB.Exec(query, guid)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
