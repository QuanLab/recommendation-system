package algorithm

import (
	"app/model/response"
	"app/model/news"
	"app/model/user"
	"strconv"
	"encoding/json"
	"log"
	"fmt"
)

const (
	TRENDING_FACTOR    = 0.2
	SIMILARITY_FACTOR  = 0.4
	INTERESTING_FACTOR = 0.4
)

func GetRecommendNews(domain string, boxid int, guid string, itemid int64) string {
	responseData := response.ResponseData{Algorithm: 1}
	fmt.Println(domain, boxid, guid, itemid)
	guiId, err := strconv.ParseInt(guid, 10, 64)
	if err != nil {
		log.Println(err)
		return "{}"
	}

	arrNews, err := news.GetNewsSimilarity(itemid)
	if err != nil {
		log.Println(err)
	}

	trendingNews, err := news.GetTrendingNews()
	if err != nil {
		log.Println(err)
	}

	mapUserNewsInterest := make(map[int64]float64, 0)
	recommends := make([]response.Post, 0)
	for _, news := range arrNews {
		mapUserNewsInterest[news.CatId] = 0
	}

	userProfile, err := user.GetUserProfileCached(guiId)
	if err != nil {
		log.Println(err)
	}

	for catId, _ := range mapUserNewsInterest {
		var dividend float64 = 10;
		var totalClick int64 = 10; // Init with G = 10 (virtual click)
		for _, period := range userProfile.PeriodTime {
			totalClick = totalClick + period.TotalClick
			for k, v := range period.MapPTCat {
				if k == catId {
					dividend = dividend + float64(period.TotalClick)*v
				}
			}
		}
		userNewsInterest := (news.GetProbCatIdEqualCi(catId) * dividend) / float64(totalClick)
		mapUserNewsInterest[catId] = userNewsInterest
		log.Println(catId, userNewsInterest)
		dividend = 10;
		totalClick = 10
	}

	for i, news := range arrNews {
		arrNews[i].SimilarityScore = SIMILARITY_FACTOR * arrNews[i].SimilarityScore

		for _, tNews := range trendingNews {
			if tNews.NewsId == news.NewsId {
				arrNews[i].SimilarityScore = arrNews[i].SimilarityScore + TRENDING_FACTOR *tNews.TrendingScore
			}
		}

		if prob, ok := mapUserNewsInterest[news.CatId]; ok {
			arrNews[i].SimilarityScore = arrNews[i].SimilarityScore + INTERESTING_FACTOR*prob
		}
	}

	arrNews = rankingScore(arrNews)
	for _, news := range arrNews {
		log.Println(news.NewsId, news.CatId, news.SimilarityScore)
		recommends = append(recommends, response.Post{ID: news.NewsId})
	}

	fmt.Println(recommends)
	responseData.Recommends = recommends
	jsonsbytes, _ := json.Marshal(responseData)
	return string(jsonsbytes)
}

func rankingScore(items []news.News) []news.News {
	var n = len(items)
	for i := 0; i < n; i++ {
		var minIdx = i
		for j := i; j < n; j++ {
			if items[j].SimilarityScore > items[minIdx].SimilarityScore {
				minIdx = j
			}
		}
		items[i], items[minIdx] = items[minIdx], items[i]
	}
	return items
}
