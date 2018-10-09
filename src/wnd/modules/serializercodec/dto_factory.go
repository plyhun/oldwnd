package serializercodec

import (
	"wnd/api"
	"wnd/api/events"
	"wnd/utils/log"

	"reflect"
	//"fmt"
)

func toDto(object interface{}) (string, interface{}) {
	log.Tracef("%#v", reflect.TypeOf(object).String())

	return reflect.TypeOf(object).String(), object
}

func fromName(name string) interface{} {
	log.Tracef("%v", name)

	switch name {
	case "*events.Chunk", "events.Chunk":
		return new(events.Chunk)
	case "*events.Chunks", "events.Chunks":
		return new(events.Chunks)
	case "*events.Blocks", "events.Blocks":
		return new(events.Blocks)
	case "*api.BlockDefinition", "api.BlockDefinition":
		return new(api.BlockDefinition)
	case "*api.BlockInChunk", "api.BlockInChunk":
		return new(api.BlockInChunk)
	case "*api.Block", "api.Block":
		return new(api.Block)
	case "*api.Chunk":
		return new(api.Chunk)
	case "*api.BiomeData":
		return new(api.BiomeData)
	case "*api.Universe":
		return new(api.Universe)
	case "*api.World":
		return new(api.World)
	case "*events.World", "events.World":
		return new(events.World)
	case "*events.EntityMove", "events.EntityMove":
		return new(events.EntityMove)
	case "*events.EntityPosition", "events.EntityPosition":
		return new(events.EntityPosition)
	default:
		log.Errorf("unsupported: %v", name)
	}

	return nil
}

func fromNameWithCount(name string, count int) interface{} {
	log.Tracef("%v / %d", name, count)

	return nil
}
