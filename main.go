package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{Body: string(getRatesPerUserPerProject()), StatusCode: 200}, nil
}

// UserProjectRate is a struct that contains the user id, the project id, the email, the project code, the project name and the rate
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

// getRatesPerUserPerProject returns the rates per user per project in json format and in a slice of UserProjectRate
func getRatesPerUserPerProject() []byte {
	var rates []UserProjectRate
	users := getHarvestActiveUsers()
	for _, user := range users {
		idUser := int(user.(map[string]interface{})["id"].(float64))
		email := user.(map[string]interface{})["email"].(string)
		ratesPerUserPerProject := getHarvestRatesPerUserPerProject(idUser)
		for _, projectRate := range ratesPerUserPerProject {
			idRate := int(projectRate.(map[string]interface{})["id"].(float64))
			project := projectRate.(map[string]interface{})["project"]
			projectCode := project.(map[string]interface{})["code"].(string)
			projectName := project.(map[string]interface{})["name"].(string)
			projectIsBillable := project.(map[string]interface{})["is_billable"].(bool)
			if !projectIsBillable {
				continue // if the project is non billable we ignore it and we dont add it to the results
			}
			amount := -1.0
			if projectRate.(map[string]interface{})["hourly_rate"] != nil {
				amount = projectRate.(map[string]interface{})["hourly_rate"].(float64)
			} else {
				continue // if rate value is nill we ignore it and we dont add it to the results
			}
			userProjectRate := UserProjectRate{idUser, idRate, email, projectCode, projectName, amount}
			rates = append(rates, userProjectRate)

		}
	}
	jsonRates, _ := json.Marshal(rates)
	return jsonRates
}
