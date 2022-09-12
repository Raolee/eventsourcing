package command

import "eventsourcing/item"

func RequestBindToCharacter(*item.Event) *State   {}
func FailedBindToCharacter(*item.Event) *State    {}
func SuccessBindToCharacter(*item.Event) *State   {}
func RequestUnbindToCharacter(*item.Event) *State {}
func FailedUnbindToCharacter(*item.Event) *State  {}
func SuccessUnbindToCharacter(*item.Event) *State {}
