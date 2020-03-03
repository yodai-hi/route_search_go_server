package db_helper

import (
	"database/sql"
	"fmt"
	. "github.com/RyanCarrier/dijkstra"
	. "github.com/yodai-hi/pbl_signage/route_search_server/entities"
	"log"
	"os"
	"strconv"
)

//DBから取得するデータを保存する配列を初期化
var (
	db *sql.DB
	err error
)


// 初期化処理
func InitDBConnection() {
	//mysqlへ接続。ドライバ名（mysql）と、ユーザー名・データソースを指定。
	db, err = sql.Open("mysql", getConnectionString())
	//接続でエラーが発生した場合の処理
	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}
}


//func FinishDBConnection() {
//	db.Close()
//}


func getConnectionString() string {
	user := getParamString("MYSQL_USER", "root")
	pass := getParamString("MYSQL_PASSWORD", "")
	protocol := getParamString("MYSQL_PROTOCOL", "tcp")
	host := getParamString("MYSQL_DATABASE_HOST", "localhost")
	port := getParamString("MYSQL_PORT", "3306")
	dbname := getParamString("MYSQL_DATABASE", "mysql")

	return fmt.Sprintf("%s:%s@%s(%s:%s)/%s", user, pass, protocol, host, port, dbname)
}


func getParamString(param string, defaultValue string) string {
	env := os.Getenv(param)
	if env != "" {
		return env
	}
	return defaultValue
}


func GenerateGraph() Graph {
	var graph Graph

	var points [] GraphPoint
	pointRows, err := db.Query("SELECT id FROM points")
	if err != nil {
		panic(err.Error())
	}
	defer pointRows.Close()
	//レコード一件一件をあらかじめ用意しておいた構造体に当てはめていく。
	for pointRows.Next() {
		//構造体Path型の変数pathを定義
		var point GraphPoint
		err := pointRows.Scan(
			&point.ID,
		)

		if err != nil {
			panic(err.Error())
		}
		points = append(points, point)
		//fmt.Println(
		//	point.ID,
		//)
	}

	var paths [] GraphPath
	pathRows, err := db.Query("SELECT start_point_id, destination_point_id, cost FROM paths")
	if err != nil {
		panic(err.Error())
	}
	defer pathRows.Close()
	//レコード一件一件をあらかじめ用意しておいた構造体に当てはめていく。
	for pathRows.Next() {
		//構造体Path型の変数pathを定義
		var path GraphPath
		err := pathRows.Scan(
			&path.StartPointID,
			&path.DestinationPointID,
			&path.Cost,
		)

		if err != nil {
			panic(err.Error())
		}
		paths = append(paths, path)
	}

	//Add the vertexes
	for _, point := range points {
		graph.AddVertex(point.ID)
	}

	//Add the arcs
	for _, path := range paths {
		var cost int64
		if path.Cost == -1 {
			cost = 1500
		}else{
			cost = int64(path.Cost)
		}
		_ = graph.AddArc(path.StartPointID, path.DestinationPointID, cost)
	}

	return graph
}


func FetchAllPoi() []Point {
	var points []Point
	//データベースへクエリを送信。引っ張ってきたデータがrowsに入る。
	pointRows, err := db.Query("SELECT id, latitude, longitude, name_jp, address_jp, outline_jp, access, tel, business, off_day, fee, how_to_book, parking, remark, image_url_1, image_url_2, image_url_3, image_url_4, image_url_5, view_order, video_url  FROM points WHERE class='poi'")
	if err != nil {
		panic(err.Error())
	}
	defer pointRows.Close()

	//レコード一件一件をあらかじめ用意しておいた構造体に当てはめていく。
	for pointRows.Next() {
		//構造体Path型の変数pathを定義
		var point Point
		err := pointRows.Scan(
			&point.ID,
			&point.Latitude,
			&point.Longitude,
			&point.NameJP,
			&point.AddressJP,
			&point.OutlineJP,
			&point.Access,
			&point.Tel,
			&point.Business,
			&point.OffDay,
			&point.Fee,
			&point.HowToBook,
			&point.Parking,
			&point.Remark,
			&point.ImageUrl1,
			&point.ImageUrl2,
			&point.ImageUrl3,
			&point.ImageUrl4,
			&point.ImageUrl5,
			&point.ViewOrder,
			&point.SpotVideo,
		)

		if err != nil {
			panic(err.Error())
		}

		points = append(points, point)

	}

	return points
}


func FetchPathId(startPoint int, destPoint int ) int {
	//データベースへクエリを送信。引っ張ってきたデータがrowsに入る。
	pathRows, err := db.Query("SELECT id FROM paths WHERE start_point_id=? and destination_point_id=?", startPoint, destPoint)
	if err != nil {
		panic(err.Error())
	}
	defer pathRows.Close()

	var pathId int
	for pathRows.Next() {
		//構造体Path型の変数pathを定義
		err := pathRows.Scan(
			&pathId,
		)

		if err != nil {
			panic(err.Error())
		}

	}

	return pathId
}


func FetchBusVideoData(currentPathId int) []Video {
	//データベースへクエリを送信。引っ張ってきたデータがrowsに入る。
	var videos []Video
	videoRows, err := db.Query("SELECT videos.id, paths_videos.fragment_path_id, videos.video_type, videos.video_name, videos.video_url FROM videos INNER JOIN paths_videos ON paths_videos.video_id=videos.id WHERE paths_videos.whole_path_id=?  ORDER BY paths_videos.order", currentPathId)
	if err != nil {
		panic(err.Error())
	}

	for videoRows.Next() {
		//構造体Path型の変数pathを定義
		var video Video
		err := videoRows.Scan(
			&video.Id,
			&video.PathId,
			&video.VideoType,
			&video.VideoName,
			&video.VideoUrl,
		)
		fmt.Print(video.PathId)

		if err != nil {
			panic(err.Error())
		}

		videos = append(videos, video)
	}

	return videos
}


func FetchWalkVideoData(beforePathId int, currentPathId int, afterPathId int) []Video {
	//データベースへクエリを送信。引っ張ってきたデータがrowsに入る。
	var resultVideo []Video
	var video Video
	var connectionVideo Video

	var beforeCuration string
	var afterCuration string

	if beforePathId == -1 {
		beforeCuration = "s"
	}else{
		beforeConnection := FetchEdgeData(beforePathId, currentPathId).ConnectionStatus
		if beforeConnection=="s"{
			beforeCuration = "f"
		}else{
			beforeCuration = "s"
		}
	}
	if afterPathId == -1 {
		afterCuration  = "s"
	}else{
		afterConnection := FetchEdgeData(currentPathId, afterPathId).ConnectionStatus
		if afterConnection=="s"{
			afterCuration = "f"
		}else{
			afterCuration = "s"
			videoRows, err := db.Query("SELECT videos.id, videos.video_type, videos.video_name, videos.video_url FROM videos WHERE videos.video_type=?",  "c"+afterConnection)
			if err != nil {
				panic(err.Error())
			}
			for videoRows.Next() {
				//構造体Path型の変数pathを定義
				err := videoRows.Scan(
					&connectionVideo.Id,
					&connectionVideo.VideoType,
					&connectionVideo.VideoName,
					&connectionVideo.VideoUrl,
				)
				if err != nil {
					panic(err.Error())
				}
			}
			resultVideo = append(resultVideo, connectionVideo)
		}
	}
	fmt.Print(strconv.Itoa(beforePathId)+","+strconv.Itoa(afterPathId)+":"+beforeCuration+afterCuration+"/")

	videoRows, err := db.Query("SELECT videos.id, paths_videos.whole_path_id, videos.video_type, videos.video_name, videos.video_url FROM videos INNER JOIN paths_videos ON paths_videos.video_id=videos.id WHERE paths_videos.whole_path_id=? and videos.video_type=?", currentPathId,  beforeCuration+afterCuration)
	if err != nil {
		panic(err.Error())
	}

	for videoRows.Next() {
		//構造体Path型の変数pathを定義
		err := videoRows.Scan(
			&video.Id,
			&video.PathId,
			&video.VideoType,
			&video.VideoName,
			&video.VideoUrl,
		)

		if err != nil {
			panic(err.Error())
		}
	}
	resultVideo = append(resultVideo, video)

	return resultVideo
}

func FetchPointId(pathId int) []int {
	pathRows, err := db.Query("SELECT start_point_id, destination_point_id FROM paths WHERE id=?", pathId)
	if err != nil {
		panic(err.Error())
	}
	defer pathRows.Close()

	var path Path
	for pathRows.Next() {
		//構造体Path型の変数pathを定義
		err := pathRows.Scan(
			&path.StartPointID,
			&path.DestinationPointID,
		)

		if err != nil {
			panic(err.Error())
		}
	}

	var result = [] int{path.StartPointID, path.DestinationPointID}
	return result
}


func FetchPathInfo(pathId int) Path {
	//データベースへクエリを送信。引っ張ってきたデータがrowsに入る。
	pathRows, err := db.Query("SELECT id, start_point_id,  destination_point_id, cost, transport FROM paths WHERE id=?", pathId)
	if err != nil {
		panic(err.Error())
	}
	defer pathRows.Close()

	var path Path
	for pathRows.Next() {
		//構造体Path型の変数pathを定義
		err := pathRows.Scan(
			&path.ID,
			&path.StartPointID,
			&path.DestinationPointID,
			&path.Cost,
			&path.Transport,
		)

		if err != nil {
			panic(err.Error())
		}
	}

	return path
}


func FetchEdgeData(beforePathId int, afterPathId int) Edge {
	//データベースへクエリを送信。引っ張ってきたデータがrowsに入る。
	edgeRows, err := db.Query("SELECT id, before_path_id, after_path_id, angle, connection_status FROM edges WHERE before_path_id=?  and after_path_id=?", beforePathId, afterPathId)
	if err != nil {
		panic(err.Error())
	}
	defer edgeRows.Close()

	var edge Edge
	for edgeRows.Next() {
		//構造体Path型の変数pathを定義
		err := edgeRows.Scan(
			&edge.ID,
			&edge.BeforePathID,
			&edge.AfterPathID,
			&edge.Angle,
			&edge.ConnectionStatus,
		)

		if err != nil {
			panic(err.Error())
		}
	}

	return edge
}


func FetchPolyline(pathId int) string {
	pathRows, err := db.Query("SELECT polyline FROM paths WHERE id=?", pathId)
	if err != nil {
		panic(err.Error())
	}

	var result string
	for pathRows.Next() {
		//構造体Path型の変数pathを定義
		err := pathRows.Scan(
			&result,
		)

		if err != nil {
			panic(err.Error())
		}
	}

	return result
}


