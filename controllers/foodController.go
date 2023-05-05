package controllers

import (
	"Restaurant_Management_Backend/database"
	"Restaurant_Management_Backend/models"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/bluesuncorp/validator.v10"
	"gopkg.in/mgo.v2/bson"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		foodId := c.Param("food_id") //"food_id"에 해당하는 값을 반환
		var food models.Food         //model.Food struct 타입의 food 변수 생성

		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food) // food 변수에 decode해서 담는다.
		defer cancel()
		if err != nil { //error 있으면
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching the food item"})
		}
		c.JSON(http.StatusOK, food) //error가 없으면 StatusOK, food 데이터 전달
	}
}

//context 작업 지시시 작업 가능 시간, 작업 취소등의 조건을 지시, 새로운 고루틴 작업 시작시 일정 시간 동작을 지시하거나 외부에서 작업을 취소할때 사용
//context.Background() or context.TODO() 빈 context 생성, TODO()는 어떤 context를 사용할지 불명확할때 사용
//var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)를 통해 타임아웃이 설정된 context 생성
//err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)에 타임 아웃 설정된 ctx context를 사용한다.
//즉 타임 아웃 설정된 시간까지 모든 작업이 완료되고 defer cancel()을 통해 cacel() 함수가 호출되지 않으면 ctx context를 사용한 곳은 전부 같이 자동으로 cancel 된다.
//mongoDB는 이런 식으로 query를 할때 context를 전달하여 지연이나 기타 이유로 DB와의 connection이 끊어 지지 않고 유지되어 부하를 발생시키는 부분에 사용하여 자동으로 cancel되게 한다.
//bson(binary Json):JSON 문서를 바이너리로 인코딩한 포맷, 주로 JSON 형태로 데이터를 저장하거나 네트워크 전송하는 용도
//bson.M:순서가 없는 map 형태, 순서를 유지하지 않음 bson.D:하나의 BSON Document 순서가 중요한 경우 bson.A:하나의 arrary 형태, bson.E:D타입 내부에 사용하는 하나의 element

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var menu models.Menu
		var food models.Food

		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		validationErr := validate.Struct(food)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		err := menuCollection.Findone(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
		defer cancel()
		if err != nil {
			msg := fmt.Sprintf("menu was not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		food.Created_at, _ = time.Parse(time.RFC3339, time.Now()).Format(time.RFC3339)
		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now()).Format(time.RFC3339)
		food.ID = primitive.NewObjectID()
		food.Food_id = food.ID.Hex()
		var num = toFixed(*food.Price, 2)
		food.Price = &num

		result, insertErr := foodCollection.InsertOne(ctx, food)
		if insertErr != nil {
			msg := fmt.Sprintf("Food item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, result)
	}
}

func round(num float64) int {

}

func toFixed(num float64, precision int) float64 {

}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

//context.Background()로 생성한 context를 context.WithTimeout() 에 넣고
// duration을 설정해 ctx를 생성하고 duration을 넘어가면 ctx를 cancel
//FindOne() : FindOne executes a find command and returns a SingleResult for one document in the collection. ctx에서 bson.m{} 값을 찾아 반환
//bson.M{}은 json처럼 key-value 형식의 데이터 값을 넣어줄 수 있도록 하는 역할
//Decode(&food): Decode will unmarshal the current event document into val
//c.JSON(code int, obj any): JSON serializes the given struct as JSON into the response body. It also sets the Content-Type as "application/json".
//gin.H is defined as type H map[string]interface{}. You can index it just like a map.
