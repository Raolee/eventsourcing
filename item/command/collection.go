package command

import "eventsourcing/item"

func RequestCollect(*item.Event) *State {}
func FailedCollect(*item.Event) *State  {}
func SuccessCollect(*item.Event) *State {}
