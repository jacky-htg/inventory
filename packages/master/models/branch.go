package models

//Branch : struct of Branch
type Branch struct {
	ID      uint32
	Code    string
	Name    string
	Address string
	Type    string
	Company Company
}
