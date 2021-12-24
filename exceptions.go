package main

import (
	"github.com/pterm/pterm"
)

func startException() {
	pterm.Error.Println("Не найдена точка входа ASL")
}

func emptyStates() {
	pterm.Error.Println("Не найдена стурктура конечных автоматов")
}

func emptyProtoBody() {
	pterm.Error.Println("Отсутсвует структура proto3 файла")
}
func notFindServiceBody() {
	pterm.Error.Println("Не найдены указатели на grpc route")

}

func rpcMethodDoesNotExist() {
	pterm.Error.Println("Не найден задекларированый тип задачи в proto3 файле")
}

func rpcRequestMessageTypeDoesNotExist(rpcName string) {
	pterm.Error.Println("Не найден тип структуры сообщения для запроса", rpcName)
}

func rpcResponseMessageTypeDoesNotExist(rpcName string) {
	pterm.Error.Println("Не найден тип структуры сообщения для ответа", rpcName)
}

func messageBodyDoesNotExist(messageStructName string) {
	pterm.Error.Println("Не найдено тип структуры сообщения", messageStructName, "в proto3 файле")
}

func fieldDoesNotExist(field string, messageType string) {
	pterm.Error.Println("В переданной структуре нет необходимого поля:", field+".", "MessageType:", messageType)
}

func invalidFieldType(field string, actual string, expected string, messageType string) {
	pterm.Error.Println(
		"Не соотвествие типов для поля:", field+". Передан тип:",
		actual+". Ожидается:", expected+".", "MessageType:", messageType)
}
