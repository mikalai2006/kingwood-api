package model

type StatNode struct {
	GroupCountry []GroupNodeCountry `json:"groupCountry" bson:"groupCountry"`
	GroupType    []GroupNodeType    `json:"groupNodeType" bson:"groupNodeType"`
}
type GroupNodeCountry struct {
	ID    string `json:"id" bson:"_id"`
	Count int    `json:"count" bson:"count"`
}
type GroupNodeType struct {
	ID    string `json:"id" bson:"_id"`
	Count int    `json:"count" bson:"count"`
}
