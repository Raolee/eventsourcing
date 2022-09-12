package command

import "eventsourcing/item"

func RequestBurn(*item.Event) *State {}
func FailedBurn(*item.Event) *State  {}
func SuccessBurn(*item.Event) *State {}
