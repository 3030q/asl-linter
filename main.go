package main

import (
	"encoding/json"
	"fmt"
	"github.com/pterm/pterm"
	"io/ioutil"
	"os"
	"time"
)

func parseAsl(fileName string) map[string]interface{} {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		pterm.Error.Println(err)
	} else {
		pterm.Success.Println("Successfully Opened", fileName)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	defer jsonFile.Close()

	return result
}

func main() {
	startConsoleLine()
	var amazonJsonFileName string
	var protoFileName string
	fmt.Print("PlayBook Name: ")
	fmt.Scanln(&amazonJsonFileName)
	fmt.Print("Proto3 File Name: ")
	fmt.Scanln(&protoFileName)
	fmt.Println()

	introSpinner, _ := pterm.DefaultSpinner.WithRemoveWhenDone(true).Start("Doing something important...")
	time.Sleep(time.Second)
	for i := 3; i > 0; i-- {
		if i > 1 {
			introSpinner.UpdateText("Doing something important... ")
		} else {
			introSpinner.UpdateText("Doing something important... ")
		}
		time.Sleep(time.Second)
	}
	introSpinner.Stop()

	amazonJsonFile := parseAsl(amazonJsonFileName + ".json")
	startAtInterface := amazonJsonFile["StartAt"]
	if startAtInterface == nil {
		startException()
		return
	}
	startAtString := startAtInterface.(string)

	states := amazonJsonFile["States"].(map[string]interface{})
	if states == nil {
		emptyStates()
		return
	}

	//По сути response -> request структура
	var steps []string
	flag := true
	var nextStep string
	for {
		if flag {
			steps = append(steps, startAtString)
			actualStep := states[startAtString].(map[string]interface{})
			if actualStep["End"] != nil {
				break
			}
			steps = append(steps, actualStep["Next"].(string))
			nextStep = actualStep["Next"].(string)
			flag = false
			continue
		}
		actualStep := states[nextStep].(map[string]interface{})
		if actualStep["End"] != nil {
			break
		}
		steps = append(steps, actualStep["Next"].(string))
		nextStep = actualStep["Next"].(string)
	}

	//Парсинг прото файла
	protoStructure := protoParse(protoFileName)
	if protoStructure == nil {
		return
	}
	protoBodyStructure := protoStructure["ProtoBody"].([]interface{})
	if protoBodyStructure == nil {
		emptyProtoBody()
		return
	}
	var serviceStructure map[string]interface{}
	for _, value := range protoBodyStructure {
		if value.(map[string]interface{})["ServiceName"] != nil {
			serviceStructure = value.(map[string]interface{})
		}
	}
	serviceBody := serviceStructure["ServiceBody"].([]interface{})
	if serviceBody == nil {
		notFindServiceBody()
		return
	}

	actualStateVariable := map[string]string{"start": "true"}
	for _, step := range steps {
		var rpcStepStructure map[string]interface{}
		for _, value := range serviceBody {
			if value.(map[string]interface{})["RPCName"] == step {
				rpcStepStructure = value.(map[string]interface{})
			}
		}
		if rpcStepStructure == nil {
			rpcMethodDoesNotExist()
			return
		}
		//Обработка запроса
		rpcRequest := rpcStepStructure["RPCRequest"].(map[string]interface{})
		if rpcRequest == nil {
			rpcRequestMessageTypeDoesNotExist(step)
		}
		rpcRequestMessageType := rpcRequest["MessageType"].(string)

		var messageRequestStructure map[string]interface{}
		for _, value := range protoBodyStructure {
			if value.(map[string]interface{})["MessageName"] == rpcRequestMessageType {
				messageRequestStructure = value.(map[string]interface{})
			}
		}
		if messageRequestStructure == nil {
			messageBodyDoesNotExist(rpcRequestMessageType)
			return
		}
		messageRequestBody := messageRequestStructure["MessageBody"].([]interface{})

		//Дикт где ключ-название переменной, значение - тип переменной
		requestVariables := map[string]string{}
		for _, value := range messageRequestBody {
			requestVariables[value.(map[string]interface{})["FieldName"].(string)] = value.(map[string]interface{})["Type"].(string)
		}
		//Проверка на соотвестивие запрашиваемых данных и данных из актуального стейта

		for key, value := range requestVariables {
			//Нулевое вхождение
			if actualStateVariable["start"] == "true" {
				actualStateVariable["start"] = "false"
				break
			}
			_, ok := actualStateVariable[key]
			if !ok {
				fieldDoesNotExist(key, rpcRequestMessageType)
				return
			}
			if actualStateVariable[key] != value {
				invalidFieldType(key, actualStateVariable[key], value, rpcRequestMessageType)
				return
			}
		}
		//Обработка ответа -> актуализирование переменных для следующего вхождения
		rpcResponse := rpcStepStructure["RPCResponse"].(map[string]interface{})
		if rpcResponse == nil {
			rpcResponseMessageTypeDoesNotExist(step)
			return
		}
		rpcResponseMessageType := rpcResponse["MessageType"].(string)
		var messageResponseStructure map[string]interface{}
		for _, value := range protoBodyStructure {
			if value.(map[string]interface{})["MessageName"] == rpcResponseMessageType {
				messageResponseStructure = value.(map[string]interface{})
			}
		}
		if messageResponseStructure == nil {
			messageBodyDoesNotExist(rpcResponseMessageType)
			return
		}
		messageResponseBody := messageResponseStructure["MessageBody"].([]interface{})
		responseVariables := map[string]string{}
		for _, value := range messageResponseBody {
			responseVariables[value.(map[string]interface{})["FieldName"].(string)] = value.(map[string]interface{})["Type"].(string)
		}
		//Обогощаем стейт выходными данными, если ключ уже существует перезаписываем новыми данными
		for key, value := range responseVariables {
			actualStateVariable[key] = value
		}
	}
	pterm.Println()
	pterm.Success.Println("Все готово! Все в порядке! You are awesome :)")
}

func startConsoleLine() {
	introScreen()
}
