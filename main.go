package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{Body: string(getRatesPerUserPerProject()), StatusCode: 200}, nil
	//return events.APIGatewayProxyResponse{Body: "hellos", StatusCode: 200}, nil
}

type UserProjectRate struct {
	IdUser      int
	IdRate      int
	Email       string
	ProjectCode string
	ProjectName string
	Rate        float64
}

func main() {
	lambda.Start(HandleRequest)
	//fmt.Print(string(getRatesPerUserPerProject()))

}

func getRatesPerUserPerProject() []byte {
	var rates []UserProjectRate
	users := getHarvestActiveUsers()
	for _, user := range users {
		idUser := int(user.(map[string]interface{})["id"].(float64))
		email := user.(map[string]interface{})["email"].(string)
		/**
		ratesPerUser := getHarvestRatesPerUser(idUser)
		for _, rate := range ratesPerUser {
			idRate := int(rate.(map[string]interface{})["id"].(float64))
			amount := rate.(map[string]interface{})["amount"].(float64)
			if rate.(map[string]interface{})["end_date"] != nil {
				continue
			}
			userProjectRate := UserProjectRate{idUser, idRate, email, "", "", amount}
			rates = append(rates, userProjectRate)
		}
		**/

		ratesPerUserPerProject := getHarvestRatesPerUserPerProject(idUser)
		for _, projectRate := range ratesPerUserPerProject {
			idRate := int(projectRate.(map[string]interface{})["id"].(float64))
			project := projectRate.(map[string]interface{})["project"]
			projectCode := project.(map[string]interface{})["code"].(string)
			projectName := project.(map[string]interface{})["name"].(string)
			projectIsBillable := project.(map[string]interface{})["is_billable"].(bool)
			if !projectIsBillable {
				continue
			}
			amount := -1.0
			if projectRate.(map[string]interface{})["hourly_rate"] != nil {
				amount = projectRate.(map[string]interface{})["hourly_rate"].(float64)
			}
			userProjectRate := UserProjectRate{idUser, idRate, email, projectCode, projectName, amount}
			rates = append(rates, userProjectRate)

		}
	}

	jsonRates, _ := json.Marshal(rates)
	return jsonRates

}

//aws lambda create-function --function-name harversRatesPerUserPerProject --runtime go1.x --zip-file fileb://function.zip --handler my-lambda-binary --role arn:aws:iam::374208052150:role/service-role/lambdaTestRole
/**
func getRatesFromHarvest() []byte {
	var rates []UserRate
	users := getHarvestActiveUsers()
	for _, user := range users {

		idUser := int(user.(map[string]interface{})["id"].(float64))
		email := user.(map[string]interface{})["email"].(string)

		ratesPerUser := getHarvestURatesPerUser(idUser)
		if len(ratesPerUser) == 0 {
			continue
		}

		for _, rate := range ratesPerUser {
			idRate := int(rate.(map[string]interface{})["id"].(float64))
			amount := rate.(map[string]interface{})["amount"].(float64)
			if rate.(map[string]interface{})["end_date"] != nil {
				continue
			}
			userRate := UserRate{idUser, idRate, email, amount}
			rates = append(rates, userRate)
		}
	}
	jsonRates, _ := json.Marshal(rates)
	fmt.Println(jsonRates)
	return jsonRates
}**/
