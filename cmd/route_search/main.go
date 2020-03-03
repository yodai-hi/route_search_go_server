package main

import (
    . "../../internal/apps/route_search/db_helper"
    . "../../internal/apps/route_search/entities"
    "fmt"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    _ "github.com/go-sql-driver/mysql"
    . "github.com/twpayne/go-polyline"
    "log"
    "math"
    "net/http"
    "os"
    "strconv"
    "strings"
    "time"
)

var(
    baseUrl = "~assets"
)


//メインスレッド
func main() {
    InitDBConnection()
    router := gin.Default()
    router.LoadHTMLGlob("templates/test.html")

    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "OPTIONS"},
        AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        AllowOriginFunc: func(origin string) bool {
            return true
        },
        MaxAge: 15 * time.Second,
    }))

    router.GET("/", func(context *gin.Context) {
        file, _ := os.Open("./test/logo.png")
        context.HTML(http.StatusOK, "test.html", gin.H{"video_url": file})
    })

    router.GET("/all_spots/", apiAllSpot)

    router.GET("/video_locations/", apiVideoLocations)

    _ = router.Run(":9000")
}

func apiAllSpot(context *gin.Context) {
    var points AllSpot
    points.Spots  = FetchAllPoi()
    for poiNum := range points.Spots{
        if points.Spots[poiNum].ImageUrl1 != ""{
            points.Spots[poiNum].ImageUrl1 = baseUrl+points.Spots[poiNum].ImageUrl1
        }
        if points.Spots[poiNum].ImageUrl2 != ""{
            points.Spots[poiNum].ImageUrl2 = baseUrl+points.Spots[poiNum].ImageUrl2

        }
        if points.Spots[poiNum].ImageUrl3 != ""{
            points.Spots[poiNum].ImageUrl3 = baseUrl+points.Spots[poiNum].ImageUrl3

        }
        if points.Spots[poiNum].ImageUrl4 != ""{
            points.Spots[poiNum].ImageUrl4 = baseUrl+points.Spots[poiNum].ImageUrl4

        }
        if points.Spots[poiNum].ImageUrl5 != ""{
            points.Spots[poiNum].ImageUrl5 = baseUrl+points.Spots[poiNum].ImageUrl5

        }
        if points.Spots[poiNum].SpotVideo != ""{
            points.Spots[poiNum].SpotVideo = baseUrl+points.Spots[poiNum].SpotVideo
        }

    }
    context.JSON(http.StatusOK, points)
}

func apiVideoLocations(context *gin.Context) {
    //入力は poi をKEYにしたquery文字列
    query, _ := context.GetQuery("spots")
    inputs :=  strings.Split(query, ",")

    graph := GenerateGraph()
    var output OutputRoutes
    output.Cost = 0
    var outputRoute []RouteData

    for inputNum := range inputs {
       if inputNum > len(inputs)-2 {
           break
       }

       var startPoI int
       var destPoI  int

       startPoI, _ = strconv.Atoi(inputs[inputNum])
       destPoI, _ = strconv.Atoi(inputs[inputNum+1])

       best, err := graph.Shortest(startPoI, destPoI)
       if err != nil {
           log.Fatal(err)
       }

       output.Cost += int(best.Distance)
       var idList []int
       for pathNum := range best.Path {
           if pathNum > len(best.Path)-2 {
               break
           }
           id := FetchPathId(best.Path[pathNum], best.Path[pathNum+1])
           idList = append(idList, id)
       }

        fmt.Print(idList)

       var videos []Video
       for routeNum := range idList {
           id := idList[routeNum]
           path := FetchPathInfo(id)
           if path.Transport=="bus"{
               busVideos := FetchBusVideoData(id)
               videos = append(videos, busVideos...)
           }else{
               var beforeId int
               var afterId int
               fmt.Print("\n")
               poly :=  FetchPolyline(id)
               if poly == ""{
                   walkVideo := Video{
                       Id:-1, PathId:id, VideoType:"", VideoName:"", VideoUrl:""}
                   videos = append(videos, walkVideo)
                   continue
               }

               if routeNum>0{
                   beforeId = idList[routeNum-1]
               }else{
                   beforeId = -1
               }

               if routeNum+1<=len(idList)-1{
                   afterId = idList[routeNum+1]
               }else{
                   afterId = -1
               }

               walkVideo := FetchWalkVideoData(beforeId, id, afterId)
               videos = append(videos, walkVideo)
           }
       }

        var route RouteData
        for videoNum := range videos {
            video := videos[videoNum]
            route.Id  = video.PathId
            routeStatus := FetchPointId(video.PathId)
            route.StartPoI = routeStatus[0]
            route.DestinationPoI = routeStatus[1]
            if route.VideoURL!=""{
                route.VideoURL = baseUrl+video.VideoUrl
            }
            polyline := FetchPolyline(route.Id)
            if polyline != "" {
                var locations []PolyLocations
                buf := []byte(polyline)
                coords, _, _ := DecodeCoords(buf)
                for locNum := range coords{
                    var location PolyLocations
                    location.DurationSeconds = 0
                    location.Latitude = math.Round(coords[locNum][0] * 100000)/100000
                    location.Longitude = math.Round(coords[locNum][1] * 100000)/100000
                    locations = append(locations, location)
                }
                route.Locations = locations
            }else{
                route.Locations = []PolyLocations{}
            }
            outputRoute = append(outputRoute, route)
        }

    }

    context.JSON(http.StatusOK, outputRoute)
}
