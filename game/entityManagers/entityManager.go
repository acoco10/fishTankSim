package entityManagers

var EntityId int

var EntityCatalogue map[int]any

func RegisterEntity(ent any) {
	EntityCatalogue[EntityId] = ent
	EntityId++
}

func GetEntity(id int) any {
	return EntityCatalogue[id]
}
